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
/********** Test Objective : Ledger Stress with 4 Clients, 4 Peers ************
*
*   Setup: 4 node peer network with security enabled
*   1. Deploy chaincode https://goo.gl/TysS79
*   2. Invoke 5K txns from each client simultaneously on 4 different peers
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
	CLIENTS   = 4
)

func initNetwork() {
	fmt.Println("========= Init Network =========")
	//peernetwork.GetNC_Local()
	peerNetworkSetup = chaincode.InitNetwork()
	chaincode.InitChainCodes()
	fmt.Println("========= Register Users =========")
	chaincode.RegisterCustomUsers()
}

func invokeChaincode(peer string) {
	counter++
	arg1Construct := []string{CHAINCODE_NAME, "invoke", peer}
	arg2Construct := []string{"a" + strconv.FormatInt(counter, 10), DATA, "counter"}

	_, _ = chaincode.InvokeOnPeer(arg1Construct, arg2Construct)
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

func InvokeLoop() {
	curTime := time.Now()
	go func() {
		for i := 1; i <= TRX_COUNT/CLIENTS; i++ {
			if counter > 0 && counter%1000 == 0 {
				elapsed := time.Since(curTime)
				fmt.Println("=========>>>>>> Iteration#", counter, " Time: ", elapsed, "CLIENT-1")
				curTime = time.Now()
			}
			//invokeChaincode("PEER0") //For Local testing
			//invokeChaincode("vp0") //For Z -Testing
			invokeChaincode(getPeer(0))
		}
		wg.Done()
	}()
	go func() {
		for i := 1; i <= TRX_COUNT/CLIENTS; i++ {
			if counter > 0 && counter%1000 == 0 {
				elapsed := time.Since(curTime)
				fmt.Println("=========>>>>>> Iteration#", counter, " Time: ", elapsed, "CLIENT-2")
				curTime = time.Now()
			}
			//invokeChaincode("PEER1") //For Local testing
			//invokeChaincode("vp1") //For Z -Testing
			invokeChaincode(getPeer(1))
		}
		wg.Done()
	}()
	go func() {
		for i := 1; i <= TRX_COUNT/CLIENTS; i++ {
			if counter > 0 && counter%1000 == 0 {
				elapsed := time.Since(curTime)
				fmt.Println("=========>>>>>> Iteration#", counter, " Time: ", elapsed, "CLIENT-3")
				curTime = time.Now()
			}
			//invokeChaincode("PEER2") //For Local testing
			//invokeChaincode("vp2") //For Z -Testing
			invokeChaincode(getPeer(2))
		}
		wg.Done()
	}()
	go func() {
		for i := 1; i <= TRX_COUNT/CLIENTS; i++ {
			if counter > 0 && counter%1000 == 0 {
				elapsed := time.Since(curTime)
				fmt.Println("=========>>>>>> Iteration#", counter, " Time: ", elapsed, "CLIENT-4")
				curTime = time.Now()
			}
			//invokeChaincode("PEER3") //For Local testing
			//invokeChaincode("vp3") //For Z -Testing
			invokeChaincode(getPeer(3))

		}
		wg.Done()
	}()
}

//Cleanup methods to display useful information
func tearDown() {
	fmt.Println("....... State transfer is happening, Lets take a nap for 3 mins ......")
	sleep(60)
	val1, val2 := queryChaincode(counter)
	fmt.Printf("\n========= After Query values a%d = %s,  counter = %s\n", counter, val1, val2)

	newVal, err := strconv.ParseInt(val2, 10, 64)

	if err != nil {
		fmt.Println("Failed to convert ", val2, " to int64\n Error: ", err)
	}

	//TODO: Block size again depends on the Block configuration in pbft config file
	//Test passes when 2 * block height match with total transactions, else fails
	if newVal == counter {
		fmt.Println("\n######### Inserted ", TRX_COUNT, " records #########\n")
		fmt.Println("######### TEST PASSED #########")
	} else {
		fmt.Println("######### TEST FAILED #########")
	}

}

//Execution starts here ...
func main() {
	//TODO:Add support similar to GNU getopts, http://goo.gl/Cp6cIg
	if len(os.Args) < 1 {
		fmt.Println("Usage: go run LedgerStressFourCliFourPeer.go Utils.go")
		return
	}
	//TODO: Have a regular expression to check if the give argument is correct format
	/*if !strings.Contains(os.Args[1], "http://") {
		fmt.Println("Error: Argument submitted is not right format ex: http://127.0.0.1:5000 ")
		return;
	}
	//Get the URL
	url := os.Args[1]*/

	// time to messure overall execution of the testcase
	defer TimeTracker(time.Now(), "Total execution time for LedgerStressFourCliFourPeer.go ")

	Init()
	fmt.Println("========= Transacations execution stated  =========")
	InvokeLoop()
	wg.Wait()
	fmt.Println("========= Transacations execution ended  =========")

	tearDown() //url
}
