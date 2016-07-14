package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/obcsdk/chaincode"
	"github.com/hyperledger/fabric/obcsdk/peernetwork"
)

//TODO: 1. have a single program to call various testcases using the existing SDK Apis
// 	2. Accept chaincode,args, functioname as parametrs  (Or)
//         read from config file as input to the program ??

/****** Test Objective: Continous Invokes on chaincode crashes [Issue #1483] *****
 *   0. Setup: N node local docker peer network with security
 *   1. Deploy chaincode example02 with Initial args
 *   2. Invoke and Query repeatedly for N times
 *   3. A scheduler task runs about N minutes to Invoke chaincode per each second
 *********************************************************************************/

var peerNetworkSetup peernetwork.PeerNetwork

var AVal, BVal, curAVal, curBVal, invokeValue int64

var argA = []string{"a"}
var argB = []string{"b"}

func setupNetwork() {
	fmt.Println("Creating a local docker network")

	//peernetwork.PrintNetworkDetails()
	peerNetworkSetup = chaincode.InitNetwork()
	chaincode.InitChainCodes()
	chaincode.RegisterUsers()

	//get a URL details to get info n chainstats/transactions/blocks etc.
	aPeer, _ := peernetwork.APeer(chaincode.ThisNetwork)
	url := "http://" + aPeer.PeerDetails["ip"] + ":" + aPeer.PeerDetails["port"]

	chaincode.NetworkPeers(url)
	chaincode.Chain_Stats(url)
}

//TODO : rather can we have a map for sleep for millis, secs and mins
func sleep(secs int64) {
	time.Sleep(time.Second * time.Duration(secs))
}

func deployChaincode() {
	example := "example02"
	var funcArgs = []string{example, "init"}
	var args = []string{argA[0], strconv.FormatInt(AVal, 10), argB[0], strconv.FormatInt(BVal, 10)}

	fmt.Println("\n######## Deploying chaincode ")
	chaincode.Deploy(funcArgs, args)

	//TODO: Increase the delay if required
	//time.Sleep(time.Second * 120)
	sleep(120)
}

func invokeChaincode() (res1, res2 int64) {
	fmt.Println("\n######## Invoke on chaincode ")
	arg1Construct := []string{"example02", "invoke"}
	arg2Construct := []string{"a", "b", strconv.FormatInt(invokeValue, 10)}

	invRes, _ := chaincode.Invoke(arg1Construct, arg2Construct)
	fmt.Println("\n Invoke response: ", invRes)

	//TODO : Can we avoid this to make them more generic?
	curAVal = curAVal - int64(invokeValue)
	curBVal = curBVal + int64(invokeValue)

	return curAVal, curBVal
}

func queryChaincode() (res1, res2 int64) {
	fmt.Println("\n######## Query on chaincode ")
	qAPIArgs0 := []string{"example02", "query"}
	var A, B string

	A, _ = chaincode.Query(qAPIArgs0, argA)
	B, _ = chaincode.Query(qAPIArgs0, argB)
	fmt.Println(fmt.Sprintf("\nA = %s B= %s", A, B))
	val1, _ := strconv.ParseInt(A, 10, 64)
	val2, _ := strconv.ParseInt(B, 10, 64)
	return val1, val2
}

//TODO: Can we change do this in more generic way
func schedulerTask() {
	//defer timeTrack(time.Now(), "schedulerTask")
	for range time.Tick(time.Second * 1) {
		invokeChaincode()
	}
}

func main() {
	// time to messure overall execution of the testcase
	defer timeTrack(time.Now(), "Testcase executiion")

	// Change values accordingly
	invokeValue = 1
	AVal = 100000
	BVal = 900000
	curAVal = AVal
	curBVal = BVal

	// Setup the network based on the NetworkCredentials.json provided
	setupNetwork()

	//Deploy the chaincode
	deployChaincode()

	var invArg1, invArg2, queryArg1, queryArg2 int64
	for i := 1; i <= 10; i++ {
		invArg1, invArg2 = invokeChaincode()
		sleep(5) // TODO : Do we need 5 secs sleep ?
		queryArg1, queryArg2 = queryChaincode()
		if invArg1 == queryArg1 && invArg2 == queryArg2 {
			fmt.Printf("\n==========================> Iteration %d is Successful", i)
		} else {
			fmt.Printf("\n==========================> Iteration %d is Failed", i)
		}
	}

	fmt.Println("######## repeate Invokes on chaincode for 2 mins")
	go schedulerTask()
	//execute schedulerTask for 1 minute(s)
	sleep(60)
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)

	fmt.Printf("\n################# %s took %s \n", name, elapsed)
	fmt.Println("################# Execution Completed #################")
}
