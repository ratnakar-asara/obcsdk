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
*   2. Invoke 500K trxns from each client, simultaneously on a single peer
*   3. Check if the counter value(5000000) matches with query value "counter"
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
	TRX_COUNT = 1000000 // 1 Million
	CLIENTS   = 2
)

func initNetwork() {
	logger("========= Init Network =========")
	peernetwork.GetNC_Local()
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
				curTime = time.Now()
			}
			//invokeChaincode("test_user3")
			invokeChaincode(getUser(0))
		}
		wg.Done()
	}()
	go func() {
		for i := 1; i <= TRX_COUNT/CLIENTS; i++ {
			if counter%1000 == 0 {
				elapsed := time.Since(curTime)
				logger(fmt.Sprintf("=========>>>>>> Iteration# %d Time: %s CLIENT-2", counter, elapsed))
				curTime = time.Now()
			}
			//invokeChaincode("test_user4")
			invokeChaincode(getUser(1))
		}
		wg.Done()
	}()
}

//Execution starts here ...
func main() {
	initLogger("LedgerStressTwoCliOnePeer500K")
	//TODO:Add support similar to GNU getopts, http://goo.gl/Cp6cIg
	if len(os.Args) < 1 {
		logger("Usage: go run LedgerStressTwoCliOnePeer500K.go Utils.go")
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
	defer TimeTracker(time.Now(), "Total execution time for LedgerStressTwoCliOnePeer500K.go ")

	Init()
	logger("========= Transacations execution stated  =========")
	InvokeMultiThreads()
	wg.Wait()
	logger("========= Transacations execution ended  =========")
	tearDown(counter)
}
