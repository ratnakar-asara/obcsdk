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

/********** Test Objective : Ledger Stress with 1 Peer and 1 Client ************
*
*   Setup: 4 node peer network with security enabled
*   1. Deploy chaincode https://goo.gl/TysS79
*   2. Invoke 1 Million trxns
*      After each 10K transactions, sleep for 1 min, StateTransfer to take place
*      All transactions takes place on single peer with single client
*   3. Check the chain height and no of transactions successful and Pass tests
*			 If results matches else Fail the test
*
* USAGE: NETWORK="LOCAL" go run LedgerStressOneCliOnePeer.go Utils.go
*  This NETWORK env value could be LOCAL or Z
*
*********************************************************************/
var peerNetworkSetup peernetwork.PeerNetwork
var AVal, BVal, curAVal, curBVal, invokeValue int64
var argA = []string{"a"}
var argB = []string{"counter"}
var wg sync.WaitGroup
var counter int64

const (
	//TODO: change value to 30000000, for inserting 30 million records
	TRX_COUNT = 1000000 // 1 Million
)

func initNetwork() {
	logger("========= Init Network =========")
	peernetwork.GetNC_Local()
	peerNetworkSetup = chaincode.InitNetwork()
	chaincode.InitChainCodes()
	logger("========= Register Users =========")
	chaincode.RegisterUsers()
}

// Utility function to invoke on chaincode available @ http://urlmin.com/4r76d
func invokeChaincode() {
	counter++
	arg1 := []string{CHAINCODE_NAME, INVOKE}
	arg2 := []string{"a" + strconv.FormatInt(counter, 10), DATA, "counter"}
	_, _ = chaincode.Invoke(arg1, arg2)
}

//Repeated Invokes based on the thread and Transaction count configuration
func invokeLoop() {

	go func() {
		curTime := time.Now()
		for i := 1; i <= TRX_COUNT; i++ {
			if counter > 1 && counter%1000 == 0 {
				elapsed := time.Since(curTime)
				logger(fmt.Sprintf("=========>>>>>> Iteration# %d Time: %s", counter, elapsed))
				sleep(30) //TODO: should we remove this delay ?
				curTime = time.Now()
			}
			invokeChaincode()
		}
		wg.Done()
	}()
}

//Execution starts from here ...
func main() {
	initLogger("LedgerStressOneCliOnePeer1Million")
	//TODO:Add support similar to GNU getopts, http://goo.gl/Cp6cIg
	if len(os.Args) < 1 {
		logger("Usage: go run LedgerStressOneCliOnePeer1Million.go Utils.go")
		return
	}

	// time to messure overall execution of the testcase
	defer TimeTracker(time.Now(), "Total execution time for LedgerStressOneCliOnePeer1Million.go ")

	//maintained counter variable to compare with final query value
	counter = 0

	//done chan int
	done := make(chan bool, 1)
	wg.Add(1)
	// Setup the network based on the NetworkCredentials.json provided
	initNetwork()

	//Deploy chaincode
	deployChaincode(done)
	logger("========= Transacations execution stated  =========")
	invokeLoop()
	wg.Wait()
	logger("========= Transacations execution ended  =========")
	tearDown(counter)
}
