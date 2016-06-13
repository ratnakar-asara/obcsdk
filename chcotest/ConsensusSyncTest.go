package main
/******************** Testing Objective: Consensu Sync Test ********
*   Setup: 5 node local docker peer network with security
*   0. Deploy chaincode example02 with 100, 200 as initial args
*   1. PAUSE ONE PEER1 
*   2. Send ONE INVOKE REQUEST
*   3. Unpause Paused PEER1
*   4. Do A Query ON same PEER0 and PEER1 as in step3
*   5. Verify query results match on PEER0 and PEER1 after invoke 
*********************************************************************/


import (
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/obcsdk/chaincode"
	"github.com/hyperledger/fabric/obcsdk/peernetwork"
)

var peerNetworkSetup peernetwork.PeerNetwork
var dAPIArgs0 = []string{"example02", "init"}

var curAVal, curBVal, invokeValue int
var argA = []string{"a"}
var argB = []string{"b"}

func setupNetwork() {
	
	fmt.Println("Creating a local docker network")

	peernetwork.PrintNetworkDetails()
	peerNetworkSetup = chaincode.InitNetwork()
	chaincode.InitChainCodes()
	chaincode.RegisterUsers()

	//get a URL details to get info n chainstats/transactions/blocks etc.
	aPeer, _ := peernetwork.APeer(chaincode.ThisNetwork)
	url := "http://" + aPeer.PeerDetails["ip"] + ":" + aPeer.PeerDetails["port"]

	chaincode.NetworkPeers(url)
	chaincode.Chain_Stats(url)
}

func sleep(secs int64) {
	time.Sleep(time.Second * time.Duration(secs))
}

func deployChaincode () {
	var AVal, BVal int64
	AVal = 100
	BVal = 900
	example := "example02"
	var funcArgs = []string{example, "init"}
	var args = []string{argA[0], strconv.FormatInt(AVal,10), argB[0], strconv.FormatInt(BVal,10)}
	
	fmt.Println("\n######## Deploying chaincode ")
	chaincode.Deploy(funcArgs, args)

	//TODO: Increase the delay if required
	//time.Sleep(time.Second * 120)
	sleep(120)
}

func invokeChaincode () (res1, res2 int) {
	fmt.Println("\n######## Invoke on chaincode ")
	arg1Construct := []string{"example02", "invoke"}
	arg2Construct := []string{"a", "b", strconv.Itoa(invokeValue)}

	invRes, _ := chaincode.Invoke(arg1Construct, arg2Construct)
	fmt.Println("\n Invoke response: ", invRes)

	//TODO : Can we avoid this to make them more generic?
	curAVal = curAVal - invokeValue
	curBVal = curBVal + invokeValue

	return curAVal, curBVal
}

func queryChaincode () (res1, res2 int) {
	fmt.Println("\n######## Query on chaincode ")
	qAPIArgs0 := []string{"example02", "query"}
	var A, B string

	A, _ = chaincode.Query(qAPIArgs0, argA)
	B, _ = chaincode.Query(qAPIArgs0, argB)
	fmt.Println(fmt.Sprintf("\nA = %s B= %s", A,B))
	val1, _ := strconv.Atoi(A)
	val2, _ := strconv.Atoi(B)
	return val1,val2
}


func pausePeer(){
	fmt.Println("####### Pause PEER1")
	peersToPause := []string{"PEER1"}
	peernetwork.PausePeersLocal(peerNetworkSetup, peersToPause)
	sleep(60)
}

func unpausePeer(){
	fmt.Println("####### Unpause PEER1")
	peernetwork.UnpausePeerLocal(peerNetworkSetup, "PEER1")
	fmt.Println("Sleeping for 2 minutes for PEER1 to sync up - state transfer")
	sleep(120)
}

//Change this functionality
func chaincodeQueryOnHost() {
	fmt.Println("####### Querying a and b on PEER0 and PEER1 ")
	qAPIArgs00 := []string{"example02", "query", "PEER0"}
	qAPIArgs01 := []string{"example02", "query", "PEER1"}

	res0A, _ := chaincode.QueryOnHost(qAPIArgs00, argA)
	res0B, _ := chaincode.QueryOnHost(qAPIArgs00, argB)

	res0AI, _ := strconv.Atoi(res0A)
	res0BI, _ := strconv.Atoi(res0B)

	if (curAVal == res0AI) && (curBVal == res0BI) {
		fmt.Println("Results in a and b match with Invoke values on PEER0:")
		valueStr := fmt.Sprintf(" AVal : %d, BVal : %d, resa : %d , resb : %d", curAVal, curBVal, res0AI, res0BI)
		fmt.Println(valueStr)
	} else {
		fmt.Println("******************************")
		fmt.Println("RESULTS DO NOT MATCH AS EXPECTED on PEER0")

		fmt.Println("******************************")
	}
	res1A, _ := chaincode.QueryOnHost(qAPIArgs01, argA)
	res1B, _ := chaincode.QueryOnHost(qAPIArgs01, argB)

	res1AI, _ := strconv.Atoi(res1A)
	res1BI, _ := strconv.Atoi(res1B)
	if (curAVal == res1AI) && (curBVal == res1BI) {
		fmt.Println("Results in a and b match with Invoke values on PEER1:")
		valueStr := fmt.Sprintf(" AVal : %d, BVal : %d, resa : %d , resb : %d", curAVal, curBVal, res1AI, res1BI)
		fmt.Println(valueStr)
	} else {
		fmt.Println("******************************")
		fmt.Println("RESULTS DO NOT MATCH AS EXPECTED on PEER1")

		fmt.Println("******************************")
	}
}

func main() {
	// Setup the network based on the Network credentials
	setupNetwork();

	deployChaincode();
	invokeChaincode();
	queryChaincode();
	
	//TODO: Pass which PEER to be paused
	pausePeer()
	
	//TODO: Add loop if required
	invokeChaincode();
	
	unpausePeer()
	//TODO : Change this as required 
	chaincodeQueryOnHost()

	
}

