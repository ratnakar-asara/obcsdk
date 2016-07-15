package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"obcsdk/chaincode"
	"obcsdk/peernetwork"
)

/********** Test Objective : Ledger Stress with 2 Clients and 1 Peer ************
*
*   Setup: 4 node peer network with security enabled
*   1. Deploy chaincode https://goo.gl/TysS79
*   2. Invoke 10K transactions from each client, simultaneously on a single peer
*   3. Check if the counter value(20000) matches with query on "counter"
*
* USAGE: NETWORK="LOCAL" go run LedgerStressOneCliOnePeer.go Utils.go
*  This NETWORK env value could be LOCAL or Z
*********************************************************************/
var peerNetworkSetup peernetwork.PeerNetwork
var AVal, BVal, curAVal, curBVal, invokeValue int64
var argA = []string{"a"}
var argB = []string{"counter"}
var counter int64
var wg sync.WaitGroup

const (
	TRX_COUNT = 20000
	CLIENTS   = 2
)

func initNetwork() {
	logger("========= Init Network =========")
	//peernetwork.GetNC_Local()
	peerNetworkSetup = chaincode.InitNetwork()
	chaincode.InitChainCodes()
	logger("========= Register Users =========")
	chaincode.RegisterCustomUsers()
}

func invokeChaincode(user string) {
	counter++
	arg1Construct := []string{CHAINCODE_NAME, "invoke", user}
	arg2Construct := []string{"a" + strconv.FormatInt(counter, 10), DATA, "counter"}

	_, _ = chaincode.InvokeAsUser(arg1Construct, arg2Construct)
}

func Init() {
	//initialize
	done := make(chan bool, 1)
	counter = 0
	wg.Add(CLIENTS)
	// Setup the network based on the NetworkCredentials.json provided
	initNetwork()

	//Deploy chaincode
	deployChaincode(done)
}

func InvokeMultiThreads() {
	curTime := time.Now()
	go func() {
		for i := 1; i <= TRX_COUNT/CLIENTS; i++ {
			if counter%1000 == 0 {
				elapsed := time.Since(curTime)
				logger(fmt.Sprintf("=========>>>>>> Iteration# %d Time: %s CLIENT-1", counter, elapsed))
				sleep(60)
				curTime = time.Now()
			}
			invokeChaincode(getUser(0))
			//invokeChaincode("user_type1_1c1b79bb45") //For Z Machine
			//invokeChaincode("test_user3")//For local Machine
		}
		wg.Done()
	}()
	go func() {
		for i := 1; i <= TRX_COUNT/CLIENTS; i++ {
			if counter%1000 == 0 {
				elapsed := time.Since(curTime)
				logger(fmt.Sprintf("=========>>>>>> Iteration# %d Time: %s CLIENT-2", counter, elapsed))
				sleep(60)
				curTime = time.Now()
			}
			invokeChaincode(getUser(1))
			//invokeChaincode("dashboarduser_type0_efeeb83216") //For Z -Machine
			//invokeChaincode("test_user4")
		}
		wg.Done()
	}()
}

//Cleanup methods to display useful information
func tearDown() {
	logger("....... State transfer is happening, Lets take a nap for 2 mins ......")
	sleep(120)
	val1, val2 := queryChaincode(counter)
	logger(fmt.Sprintf("========= After Query values a%d = %s,  counter = %s\n", counter, val1, val2))

	newVal, err := strconv.ParseInt(val2, 10, 64)
	if err != nil {
		logger(fmt.Sprintf("Failed to convert %s to int64\n Error: %s", val2, err))
	}

	//TODO: Block size again depends on the Block configuration in pbft config file
	//Test passes when 2 * block height match with total transactions, else fails
	if newVal == counter {
		logger(fmt.Sprintf("######### Inserted %d records #########\n", TRX_COUNT))
		logger("######### TEST PASSED #########")
	} else {
		logger("######### TEST FAILED #########")
	}
}

//Execution starts here ...
func main() {
	initLogger("LedgerStressTwoCliOnePeer")
	//TODO:Add support similar to GNU getopts, http://goo.gl/Cp6cIg
	if len(os.Args) < 1 {
		logger("Usage: go run LedgerStressTwoCliOnePeer.go Utils.go")
		return
	}
	//TODO: Have a regular expression to check if the give argument is correct format
	/*if !strings.Contains(os.Args[1], "http://") {
		logger("Error: Argument submitted is not right format ex: http://127.0.0.1:5000 ")
		return;
	}*/
	//Get the URL
	//url := os.Args[1]
	// time to messure overall execution of the testcase
	defer TimeTracker(time.Now(), "Total execution time for LedgerStressOneCliOnePeer.go ")

	Init()
	logger("========= Transacations execution stated  =========")
	InvokeMultiThreads()
	wg.Wait()
	logger("========= Transacations execution ended  =========")
	tearDown() //url
}
