package main

import (
	"fmt"
	"obcsdk/chaincode"
	"obcsdk/peernetwork"
	"strconv"
	"time"
)

var peerNetworkSetup peernetwork.PeerNetwork
var AVal, BVal, curAVal, curBVal, invokeValue int64

var argA = []string{"a"}
var argB = []string{"b"}

const (
	INVOKE_COUNT = 20
	TOTAL_NODES  = 4
)

var url string

func setupNetwork() {
	fmt.Println("Creating a local docker network")
	//peernetwork.SetupLocalNetwork(TOTAL_NODES, true)
	//peernetwork.PrintNetworkDetails()
	peernetwork.GetNC_Local()
	peerNetworkSetup = chaincode.InitNetwork()
	chaincode.InitChainCodes()
	chaincode.RegisterUsers()

	//get a URL details to get info n chainstats/transactions/blocks etc.
	aPeer, _ := peernetwork.APeer(chaincode.ThisNetwork)
	//fmt.Println(aPeer)
	url = "http://" + aPeer.PeerDetails["ip"] + ":" + aPeer.PeerDetails["port"]
	fmt.Println(url)
	//chaincode.NetworkPeers(url)
	//chaincode.Chain_Stats(url)
}

func sleep(secs int64) {
	time.Sleep(time.Second * time.Duration(secs))
}

func deployChaincode() {
	example := "example02"
	var funcArgs = []string{example, "init"}
	var args = []string{argA[0], strconv.FormatInt(AVal, 10), argB[0], strconv.FormatInt(BVal, 10)}

	fmt.Println("\n######## Deploying chaincode ")
	chaincode.Deploy(funcArgs, args)
	sleep(60)
}

func invokeChaincode() (res1, res2 int64) {
	fmt.Println("\n######## Invoke on chaincode ")
	arg1Construct := []string{"example02", "invoke"}
	arg2Construct := []string{"a", "b", strconv.FormatInt(invokeValue, 10)}

	invRes, _ := chaincode.Invoke(arg1Construct, arg2Construct)
	//fmt.Println("\n Invoke response: ", invRes)

	//TODO : Can we avoid this to make them more generic?
	curAVal = curAVal - int64(invokeValue)
	curBVal = curBVal + int64(invokeValue)
	fmt.Println("\n Invoke Transaction ID: ", invRes)
	//fmt.Println(fmt.Sprintf("\n  Values after Invoke A = %d B= %d", curAVal,curBVal))

	return curAVal, curBVal
}

func invokeChaincodeOnHost() (res1, res2 int64) {
	fmt.Println("\n######## Invoke on chaincode ")
	arg1Construct := []string{"example02", "invoke", "PEER1"}
	arg2Construct := []string{"a", "b", strconv.FormatInt(invokeValue, 10)}

	invRes, _ := chaincode.InvokeOnPeer(arg1Construct, arg2Construct)
	//fmt.Println("\n Invoke response: ", invRes)

	//TODO : Can we avoid this to make them more generic?
	curAVal = curAVal - int64(invokeValue)
	curBVal = curBVal + int64(invokeValue)
	fmt.Println("\n Invoke Transaction ID: ", invRes)
	//fmt.Println(fmt.Sprintf("\n  Values after Invoke A = %d B= %d", curAVal,curBVal))

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
	//fmt.Println(fmt.Sprintf("\n  Values after Query A = %d B= %d", val1,val2))
	return val1, val2
}

func Testcase1() {

	// STEP1 - Deploy chaaincode and several Invokes
	invokeValue = 1
	AVal = 100
	BVal = 200
	curAVal = AVal
	curBVal = BVal

	//Deploy the chaincode
	deployChaincode()

	//Multiple Invokes
	for i := 1; i <= 12; i++ {
		fmt.Println("############## Invoke Iteration:", i)
		_, _ = invokeChaincode()
	}
	sleep(30)
	//_, _ = queryChaincode();
	getBlocksHeight()
	_ = QueryOnHostTest()
	peernetwork.StopPeerLocal(peerNetworkSetup, "PEER3")
	sleep(15)
	for i := 1; i <= 20; i++ {
		fmt.Println("############## PEER3 is Stopped ,  Invoke Iteration:", i)
		//_, _ = invokeChaincode();
		_, _ = invokeChaincodeOnHost()
	}
	sleep(15)
	chaincode.RegisterUsers2()
	QueryOnHostTest2()
	getBlocksHeight2()
	peernetwork.StartPeerLocal(peerNetworkSetup, "PEER3")
	sleep(15)
	for i := 1; i <= 20; i++ {
		fmt.Println("##############  PEER3 Started , Invoke Iteration:", i)
		_, _ = invokeChaincodeOnHost()
	}
	sleep(30)
	chaincode.RegisterUsers2()
	_ = QueryOnHostTest()
	getBlocksHeight()
	fmt.Println("######## Testcase execution DONE")
}

func getBlocksHeight2() {
	startValue := 3
	height := 0
	var urlStr string
	for i := 0; i < 3; i++ {
		urlStr = "http://172.17.0." + strconv.Itoa(startValue+i) + ":5000"
		height = chaincode.Monitor_ChainHeight(urlStr)
		fmt.Println("################ Chaincode Height on "+urlStr+" is : ", height)
	}
}
func QueryOnHostTest2() {
	qAPIArgs00 := []string{"example02", "query", "PEER0"}
	qAPIArgs01 := []string{"example02", "query", "PEER1"}
	qAPIArgs02 := []string{"example02", "query", "PEER2"}
	qArgsa := []string{"a"}
	res0A, _ := chaincode.QueryOnHost(qAPIArgs00, qArgsa)
	res0AI, _ := strconv.Atoi(res0A)
	fmt.Printf("\n*********** PEER0 : A Value is %d, While PEER3 is down", res0AI)

	res0A, _ = chaincode.QueryOnHost(qAPIArgs01, qArgsa)
	res1AI, _ := strconv.Atoi(res0A)
	fmt.Printf("\n*********** PEER1 : A Value is %d, While PEER3 is down", res1AI)

	res0A, _ = chaincode.QueryOnHost(qAPIArgs02, qArgsa)
	res2AI, _ := strconv.Atoi(res0A)
	fmt.Printf("\n*********** PEER2 : A Value is %d, While PEER3 is down", res2AI)
}

func main() {
	// time to messure overall execution of the testcase
	defer timeTrack(time.Now(), "Testcase executiion")

	// Setup the network based on the NetworkCredentials.json provided
	setupNetwork()
	Testcase1()
	//Testcase2()
}

func Testcase2() {
	// STEP1 - Deploy chaaincode and several Invokes
	invokeValue = 1
	AVal = 100
	BVal = 200
	curAVal = AVal
	curBVal = BVal
	peernetwork.StopPeerLocal(peerNetworkSetup, "PEER2")
	peernetwork.StopPeerLocal(peerNetworkSetup, "PEER3")
	//Deploy the chaincode
	deployChaincode()
	peernetwork.StartPeerLocal(peerNetworkSetup, "PEER2")
	sleep(15)
	for i := 1; i <= 20; i++ {
		fmt.Println("############## Invoke Iteration:", i)
		_, _ = invokeChaincode()
	}
	QueryOnHostTest2()
}

func QueryOnHostTest() bool {
	qAPIArgs00 := []string{"example02", "query", "PEER0"}
	qAPIArgs01 := []string{"example02", "query", "PEER1"}
	qAPIArgs02 := []string{"example02", "query", "PEER2"}
	qAPIArgs03 := []string{"example02", "query", "PEER3"}
	qArgsa := []string{"a"}
	res0A, _ := chaincode.QueryOnHost(qAPIArgs00, qArgsa)
	res0AI, _ := strconv.Atoi(res0A)
	fmt.Printf("\n\n*********** PEER0 : A Value is %d", res0AI)

	res0A, _ = chaincode.QueryOnHost(qAPIArgs01, qArgsa)
	res1AI, _ := strconv.Atoi(res0A)
	fmt.Printf("\n\n*********** PEER1 : A Value is %d", res1AI)

	res0A, _ = chaincode.QueryOnHost(qAPIArgs02, qArgsa)
	res2AI, _ := strconv.Atoi(res0A)
	fmt.Printf("\n\n*********** PEER2 : A Value is %d", res2AI)

	res0A, _ = chaincode.QueryOnHost(qAPIArgs03, qArgsa)
	res3AI, _ := strconv.Atoi(res0A)
	fmt.Printf("\n\n*********** PEER3 : A Value is %d", res3AI)

	if res0AI != res1AI || res1AI != res2AI || res2AI != res3AI {
		return false
	}
	return true
}

func getBlocksHeight() {
	startValue := 3
	height := 0
	var urlStr string
	for i := 0; i < TOTAL_NODES; i++ {
		urlStr = "http://172.17.0." + strconv.Itoa(startValue+i) + ":5000"
		height = chaincode.Monitor_ChainHeight(urlStr)
		fmt.Println("################ Chaincode Height on "+urlStr+" is : ", height)
	}
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)

	fmt.Printf("\n################# %s took %s \n", name, elapsed)
	fmt.Println("################# Execution Completed #################")
}
