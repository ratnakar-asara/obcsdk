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
*   1. Deploy chaincode
*   2. Invoke 20K transactions (TODO: Should make this configurable ?)
*      After each 10K transactions, sleep for 1 min, StateTransfer to take place
*      All transactions takes place on single peer with single client
*   3. Check the chain height and no of transactions successful and Pass tests
*			 If results matches else Fail the test
*
* USAGE: go run LedgerStressOneCliOnePeer.go Utils.go http://172.17.0.3:5000
*
*********************************************************************/
var peerNetworkSetup peernetwork.PeerNetwork
var AVal, BVal, curAVal, curBVal, invokeValue int64
var argA = []string{"a"}
var argB = []string{"counter"}

var counter int64
var wg sync.WaitGroup

const (
	TRX_COUNT = 20000 //3000000 Enable for long runs
)

func initNetwork() {
	fmt.Println("========= Init Network =========")
	peernetwork.GetNC_Local()
	peerNetworkSetup = chaincode.InitNetwork()
	chaincode.InitChainCodes()
	fmt.Println("========= Register Users =========")
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
				fmt.Println("=========>>>>>> Iteration#", counter, " Time: ", elapsed)
				sleep(30)
				curTime = time.Now()
			}
			invokeChaincode()
		}
		wg.Done()
	}()
}

//Cleanup methods to display useful information
func tearDown() {
	fmt.Println("....... State transfer is happening, Lets take a nap for two minute ......")
	sleep(120)
	fmt.Println("========= Counter is", counter)
	val1, val2 := queryChaincode(counter)
	fmt.Printf("\n========= After Query values a%d = %s,  counter = %s\n", counter, val1, val2)

	newVal, err := strconv.ParseInt(val2, 10, 64)

	if err != nil {
		fmt.Println("Failed to convert ", val2, " to int64\n Error: ", err)
	}

	//TODO: Block size again depends on the Block configuration in pbft config file
	//Test passes when 2 * block height match with total transactions, else fails
	if newVal == counter {
		fmt.Println("\n######### Inserted ", counter, " records #########\n")
		fmt.Println("######### TEST PASSED #########")
	} else {
		fmt.Println("######### TEST FAILED #########")
	}
}

//Execution starts from here ...
func main() {
	//TODO:Add support similar to GNU getopts, http://goo.gl/Cp6cIg
	if len(os.Args) < 1 {
		fmt.Println("Usage: go run LedgerStressOneCliOnePeer.go Utils.go")
		return
	}
	//TODO: Have a regular expression to check if the give argument is correct format
	/*if !strings.Contains(os.Args[1], "http://") {
		fmt.Println("Error: Argument submitted is not right format ex: http://127.0.0.1:5000 ")
		return;
	}*/
	//Get the URL
	//url = os.Args[1]

	// time to messure overall execution of the testcase
	defer TimeTracker(time.Now(), "Total execution time for LedgerStressOneCliOnePeer.go ")

	//maintained counter variable to compare with final query value
	counter = 0
	wg.Add(1)
	done := make(chan bool, 1)

	// Setup the network based on the NetworkCredentials.json provided
	initNetwork()

	//Deploy chaincode
	deployChaincode(done)
	fmt.Println("========= Transacations execution stated  =========")
	invokeLoop()
	wg.Wait()
	fmt.Println("========= Transacations execution ended  =========")
	tearDown()
}
