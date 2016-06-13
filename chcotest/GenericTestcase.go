package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"obcsdk/chaincode"
	"obcsdk/peernetwork"
	"obcsdk/peerrest"

)

//TODO: 1. have a single program to call various testcases using the existing SDK Apis
// 	2. Accept chaincode,args, functioname as parametrs  (Or)
//         read from config file as input to the program ??

/****** Test Objective: Generic Testcase *****
 * Make it generic to execute various tests
 * Issue1331
 * Issue1478
 * Issue1545
 * SyncTest
 *********************************************************************************/

var peerNetworkSetup peernetwork.PeerNetwork

var AVal, BVal, curAVal, curBVal, invokeValue int64

var argA = []string{"a"}
var argB = []string{"b"}

//Change this value as per usecase //TBD: should we have a better approach to read this from a config file ?
const (
	INVOKE_COUNT = 30
	TOTAL_NODES = 4
)
var url string
func setupNetwork() {
	fmt.Println("Creating a local docker network")
  peernetwork.SetupLocalNetwork(TOTAL_NODES, true)
	//peernetwork.PrintNetworkDetails()
	peerNetworkSetup = chaincode.InitNetwork()
	chaincode.InitChainCodes()
	chaincode.RegisterUsers()

	//get a URL details to get info n chainstats/transactions/blocks etc.
	aPeer, _ := peernetwork.APeer(chaincode.ThisNetwork)
	url = "http://" + aPeer.PeerDetails["ip"] + ":" + aPeer.PeerDetails["port"]

	chaincode.NetworkPeers(url)
	chaincode.Chain_Stats(url)
}

//TODO : rather can we have a map for sleep for millis, secs and mins
func sleep(secs int64) {
	time.Sleep(time.Second * time.Duration(secs))
}

func deployChaincode () {
	example := "example02"
	var funcArgs = []string{example, "init"}
	var args = []string{argA[0], strconv.FormatInt(AVal,10), argB[0], strconv.FormatInt(BVal,10)}

	fmt.Println("\n######## Deploying chaincode ")
	chaincode.Deploy(funcArgs, args)

	//TODO: Increase the delay if required
	//time.Sleep(time.Second * 120)
	sleep(60)
}

func invokeChaincode () (res1, res2 int64) {
	fmt.Println("\n######## Invoke on chaincode ")
	arg1Construct := []string{"example02", "invoke"}
	arg2Construct := []string{"a", "b", strconv.FormatInt(invokeValue,10)}

	invRes, _ := chaincode.Invoke(arg1Construct, arg2Construct)
	//fmt.Println("\n Invoke response: ", invRes)

	//TODO : Can we avoid this to make them more generic?
	curAVal = curAVal - int64(invokeValue)
	curBVal = curBVal + int64(invokeValue)
  fmt.Println("\n Invoke Transaction ID: ", invRes)
	//fmt.Println(fmt.Sprintf("\n  Values after Invoke A = %d B= %d", curAVal,curBVal))

	return curAVal, curBVal
}

func invokeChaincodeOnHost () (res1, res2 int64) {
	fmt.Println("\n######## Invoke on chaincode ")
	arg1Construct := []string{"example02", "invoke", "PEER3"}
	arg2Construct := []string{"a", "b", strconv.FormatInt(invokeValue,10)}

	invRes, _ := chaincode.InvokeOnPeer(arg1Construct, arg2Construct)
	//fmt.Println("\n Invoke response: ", invRes)

	//TODO : Can we avoid this to make them more generic?
	curAVal = curAVal - int64(invokeValue)
	curBVal = curBVal + int64(invokeValue)
  fmt.Println("\n Invoke Transaction ID: ", invRes)
	//fmt.Println(fmt.Sprintf("\n  Values after Invoke A = %d B= %d", curAVal,curBVal))

	return curAVal, curBVal
}

func queryChaincode () (res1, res2 int64) {
	fmt.Println("\n######## Query on chaincode ")
	qAPIArgs0 := []string{"example02", "query"}
	var A, B string

	A, _ = chaincode.Query(qAPIArgs0, argA)
	B, _ = chaincode.Query(qAPIArgs0, argB)
	fmt.Println(fmt.Sprintf("\nA = %s B= %s", A,B))
	val1, _ := strconv.ParseInt(A,10, 64)
	val2, _ := strconv.ParseInt(B,10, 64)
	//fmt.Println(fmt.Sprintf("\n  Values after Query A = %d B= %d", val1,val2))
	return val1,val2
}

//TODO: Can we change this more generic
func schedulerTask() {
	//defer timeTrack(time.Now(), "schedulerTask")
	for range time.Tick(time.Second * 1){
		invokeChaincode();
	}
}

func pausePeer( peer string){
	fmt.Println("####### Pause ", peer)
	peersToPause := []string{peer}//"PEER1"}
	peernetwork.PausePeersLocal(peerNetworkSetup, peersToPause)
	sleep(20)
}

func unpausePeer(peer string){
	fmt.Println("####### Unpause ", peer)
	peernetwork.UnpausePeerLocal(peerNetworkSetup, peer)//"PEER1")
	fmt.Printf("\n Sleeping for 1 minute(s) for %s to sync up - state transfer",peer)
	sleep(30)
}
func Issue1331_3(){

	// STEP1 - Deploy chaaincode and several Invokes
	invokeValue = 1
	AVal = 100000
	BVal = 900000
	curAVal = AVal
	curBVal = BVal
	//var invArg1, invArg2, queryArg1, queryArg2 int64

	//Deploy the chaincode
	deployChaincode();
	//Multiple Invokes
	for i :=1; i<= INVOKE_COUNT;i ++ {
		fmt.Println("############## Invoke Iteration:", i)
		_, _ = invokeChaincode();
		//sleep(5) //TODO: Do we need 5 secs sleep ?
	}
	_, _ = queryChaincode();
	sleep(120)
	//getBlocksHeight()
	_ = QueryOnHostTest()
	//peernetwork.StopPeerLocal(peerNetworkSetup, "PEER3")
	sleep(15)
  _, _ = invokeChaincode()
	//peernetwork.StartPeerLocal(peerNetworkSetup, "PEER3")

	sleep(120)
	for i :=1; i<= 10;i ++ {
		fmt.Println("############## Invoke Iteration:", i)
		//_, _ = invokeChaincode();
		_, _ = invokeChaincodeOnHost();
	}
	sleep(120)
	_ = QueryOnHostTest()
	getBlocksHeight()
	fmt.Println("######## Testcase Issue1331_3 execution DONE")
}

func main() {
	// time to messure overall execution of the testcase
	defer timeTrack(time.Now(), "Testcase executiion")

	// Setup the network based on the NetworkCredentials.json provided
	setupNetwork()

	args := os.Args
	if (len(args) <= 1) {
		//fmt.Println("####### Running All Programs ######")
		//QueryOnHostTest();
		//Issue1545_2()
		Issue1331_3()
		//executeAlltests()
	} else {

		// Have a switch or map and call corresponding function , ex: #1545 --> call Issue1545()
	}

	//go schedulerTask()
	//execute schedulerTask for 1 minute(s)
	//sleep(60);
}
func QueryOnHostTest() bool{
	/*invokeValue = 1
	AVal = 100
	BVal = 200
	curAVal = AVal
	curBVal = BVal*/
	//deployChaincode();
	 qAPIArgs00 := []string{"example02", "query", "PEER0"}
	 qAPIArgs01 := []string{"example02", "query", "PEER1"}
	 qAPIArgs02 := []string{"example02", "query", "PEER2"}
	 qAPIArgs03 := []string{"example02", "query", "PEER3"}
	 qArgsa := []string{"a"}
	 res0A, _ := chaincode.QueryOnHost(qAPIArgs00, qArgsa)
	 res0AI, _ := strconv.Atoi(res0A)
	 fmt.Printf("*********** PEER0 : A Value is %d",res0AI)

	 res0A, _ = chaincode.QueryOnHost(qAPIArgs01, qArgsa)
	 res1AI, _ := strconv.Atoi(res0A)
	 fmt.Printf("*********** PEER1 : A Value is %d",res1AI)

	 res0A, _ = chaincode.QueryOnHost(qAPIArgs02, qArgsa)
	 res2AI, _ := strconv.Atoi(res0A)
	 fmt.Printf("*********** PEER2 : A Value is %d",res2AI)

	 res0A, _ = chaincode.QueryOnHost(qAPIArgs03, qArgsa)
	 res3AI, _ := strconv.Atoi(res0A)
	 fmt.Printf("*********** PEER3 : A Value is %d",res3AI)

	 if (res0AI != res1AI || res1AI != res2AI || res2AI != res3AI ) {
		 return false
	 }
	 return true
}
func Issue1478() {
	fmt.Println("####### Running Test for Issue1478 ######")
	// STEP1 - Deploy chaaincode and several Invokes
	invokeValue = 1
	AVal = 100000
	BVal = 900000
	curAVal = AVal
	curBVal = BVal
	deployChaincode();
	repeatInvokQueries(2)
	getBlocksHeight()
	pausePeer("PEER3")
	//repeatNInvokes(100);
	repeatInvokQueries(100)
	getBlocksHeight()
}

func executeAlltests() {
	//TODO: Remove redundant functions
	Issue1331()
	Issue1478()
	Issue1545()
	SyncTest()
	//PeerResetTest()
	DelayTest()
}

func perfTests() {
	invokeValue = 1
	AVal = 100000
	BVal = 900000
	curAVal = AVal
	curBVal = BVal
	deployChaincode();
	for i:=0;i<100;i++ {

	}
}
func getHt() {
	url1 := "http://172.17.0.3:5000"
	height := chaincode.Monitor_ChainHeight(url1)
	fmt.Println("################ Chaincode on URL: "+url1+" height : ", height)
}

func PeerResettest() {
	invokeValue = 1
	AVal = 100000
	BVal = 900000
	curAVal = AVal
	curBVal = BVal

	//Deploy the chaincode
	deployChaincode();
	repeatInvokQueries(6)
	getHt()
	peernetwork.StopPeerLocal(peerNetworkSetup, "PEER0")
	peernetwork.StartPeerLocal(peerNetworkSetup, "PEER0")
	getHt()
}

func repeatInvokQueries(n int){
	var invArg1, invArg2, queryArg1, queryArg2 int64
	for i :=1; i<= n;i ++ {
		invArg1, invArg2 = invokeChaincode();
		fmt.Println(fmt.Sprintf("\n >>>>>>>> Values after Invoke A = %d B= %d", invArg1,invArg2))
		sleep(2)
		//Verify invokes by querying chaincode
		queryArg1, queryArg2 = queryChaincode()
		fmt.Println(fmt.Sprintf("\n >>>>>>>>  Values after Query A = %d B= %d", queryArg1,queryArg1))
		if (invArg1 == queryArg1 && invArg2 == queryArg2){
			fmt.Println("\n==========================> Iter"+strconv.Itoa(i)+", Invoke and Query Successful")
		} else {
			fmt.Println("\n==========================> Iter"+strconv.Itoa(i)+", Invoke and Query Failed")
				//sleep(10)
			//TODO: Check docker peer status and unpause if required
			//os.Exit(1)
		}
	}
}

//TODO:Take input params
func pauseUnpausePeers(){
	pausePeer("PEER2")
	sleep(2)
	unpausePeer("PEER2")
}

func stopStartPeers(){
	peernetwork.StopPeerLocal(peerNetworkSetup, "PEER2")
	sleep(2)
	peernetwork.StartPeerLocal(peerNetworkSetup, "PEER2")
}

func repeatNInvokes(n int){
	var invArg1, invArg2 int64
	for i :=1; i<= n;i ++ {
		invArg1, invArg2 = invokeChaincode();
		fmt.Println("\n==========================> Iter"+strconv.Itoa(i)+", After Invoke values are : "+strconv.FormatInt(invArg1, 10)+", "+strconv.FormatInt(invArg2, 10))
		//isConsistent := checkBlockStates();
		if (!QueryOnHostTest()){
			sleep(10)
			//Second Check for state transfer to complete
			if (!QueryOnHostTest() || !checkBlockStates()) {
				fmt.Println("#################### Current Hashblocks across nodes are inconsistent ###########")
				fmt.Println("#################### Exiting... ###########")
				return;
			}
		}
	}
}
func repeatNInvokes1(n int){
	var invArg1, invArg2 int64
	for i :=1; i<= n;i ++ {
		go func(){
		invArg1, invArg2 = invokeChaincode()
		fmt.Println("\n==========================> Iter"+strconv.Itoa(i)+", After Invoke values are : "+strconv.FormatInt(invArg1, 10)+", "+strconv.FormatInt(invArg2, 10))
		}()
	}
	sleep (20)
	_ = QueryOnHostTest()
	_ = checkBlockStates()

	sleep (20)
	_ = QueryOnHostTest()
	_ = checkBlockStates()

	fmt.Println("#################### Test Done... ###########")
}
func Issue1545(){
	// STEP1 - Deploy chaaincode and several Invokes
	invokeValue = 1
	AVal = 100000
	BVal = 900000
	curAVal = AVal
	curBVal = BVal
	deployChaincode();
	repeatNInvokes(250)
	//How to get URL ?
	peerrest.GetChainInfo(url + "/chain")
}
func getBlocksHeight(){
	startValue := 3
	height := 0
	var urlStr string
	for i:=0;i<TOTAL_NODES;i++ {
		urlStr = "http://172.17.0."+strconv.Itoa(startValue+i)+":5000"
		height = chaincode.Monitor_ChainHeight(urlStr)
		fmt.Println("################ Chaincode Height on "+urlStr+" is : ", height)
	}
}
func getCurrentBlockHash(){
	startValue := 3
	height := 0
	var urlStr string
	for i:=0;i<TOTAL_NODES;i++ {
		urlStr = "http://172.17.0."+strconv.Itoa(startValue+i)+":5000"
		height = chaincode.Monitor_ChainHeight(urlStr)
		fmt.Println("################ Chaincode Height on "+urlStr+" is : ", height)
		respBody := chaincode.ChaincodeBlockHash(urlStr, height)
		fmt.Println("################ StateHash on PEER"+strconv.Itoa(i)+" : "+respBody)
	}
}
func checkBlockStates() bool {
	//Avoid hardcoding and get the details from
	url1 := "http://172.17.0.3:5000"
	height := chaincode.Monitor_ChainHeight(url1)
	fmt.Println("################ Chaincode Height on "+url1+" is : ", height)
	respBody1 := chaincode.ChaincodeBlockHash(url1, height)
	fmt.Println("################ StateHash on PEER0 : "+respBody1)

	url1 = "http://172.17.0.4:5000"
	height = chaincode.Monitor_ChainHeight(url1)
	fmt.Println("################ Chaincode Height on "+url1+" is : ", height)
	respBody2 := chaincode.ChaincodeBlockHash(url1, height)
	fmt.Println("################ StateHash on PEER1 : "+respBody2)

	url1 = "http://172.17.0.5:5000"
	height = chaincode.Monitor_ChainHeight(url1)
	fmt.Println("################ Chaincode Height on "+url1+" is : ", height)
	respBody3 := chaincode.ChaincodeBlockHash(url1, height)
	fmt.Println("################ StateHash on PEER2 : "+respBody3)

	url1 = "http://172.17.0.6:5000"
	height = chaincode.Monitor_ChainHeight(url1)
	fmt.Println("################ Chaincode Height on "+url1+" is : ", height)
	respBody4 := chaincode.ChaincodeBlockHash(url1, height)
	fmt.Println("\n################ StateHash on PEER3 : "+respBody4)
	if (respBody1 == respBody2 && respBody2 == respBody3 && respBody3 == respBody4) {
		return true;
	}
	return false;

}
/*func getAllBlockStats(){
	//Avoid hardcoding and get the details from
	url1 := "http://172.17.0.3:5000"
	respBody , respStatus := chaincode.Chain_Stats(url1, height)
	fmt.Println("################ ResponseBody : "+respBody)

	url1 = "http://172.17.0.4:5000"
	respBody , respStatus = chaincode.Chain_Stats(url1, height-1)
	fmt.Println("################ ResponseBody : "+respBody)

	url1 = "http://172.17.0.5:5000"
	respBody , respStatus = chaincode.Chain_Stats(url1, height-1)
	fmt.Println("################ ResponseBody : "+res/home/ratnakar/auctionapp/blockchain
	respBody , respStatus = chaincode.Chain_Stats(url1, height-1)
	fmt.Println("################ ResponseBody : "+respBody)

	//url1 = "http://172.17.0.7:5000"
	//respBody , respStatus = chaincode.Chain_Stats(url1)
	//fmt.Println("################ ResponseBody : "+respBody)

}*/
//Check for A value on all nodes and get A value on all peers
func Issue1545_2(){
	// STEP1 - Deploy chaaincode and several Invokes
	invokeValue = 10
	AVal = 10000
	BVal = 20000
	curAVal = AVal
	curBVal = BVal
	deployChaincode()
	repeatNInvokes1(650)
}

func Issue1545_1(){
	// STEP1 - Deploy chaaincode and several Invokes
	invokeValue = 1
	AVal = 100000
	BVal = 900000
	curAVal = AVal
	curBVal = BVal
	deployChaincode();
	repeatNInvokes(250)
  getBlocksHeight()
	repeatNInvokes(265)
	sleep(10)
	//How to get URL ?
	getBlocksHeight()
}

func DelayTest(){
	// STEP1 - Deploy chaaincode and several Invokes
	invokeValue = 1
	AVal = 100000
	BVal = 900000
	curAVal = AVal
	curBVal = BVal

	//Deploy the chaincode
	deployChaincode();

	repeatInvokQueries(20)

	//Start Peer
	peernetwork.StopPeerLocal(peerNetworkSetup, "PEER3")

	repeatInvokQueries(20)

	//Stop Peer
        peernetwork.StartPeerLocal(peerNetworkSetup, "PEER3")

	//Multiple Invokes/Queries
	repeatInvokQueries(20) //INVOKE_COUNT

	// TODO : Check for PEER1 block height (from /chain REST API endpoint) is different from other peers ?
	fmt.Println("######## DelayTest execution done ...")
}
func Issue1331_1(){

	// STEP1 - Deploy chaaincode and several Invokes
	invokeValue = 1
	AVal = 100
	BVal = 200
	curAVal = AVal
	curBVal = BVal
	//var err error

	//Deploy the chaincode
	deployChaincode();
	//Multiple Invokes
	for i :=1; i<= 30;i ++ {
		fmt.Println("############## Invoke Iteration:", i)
		_, _ = invokeChaincode();
		//sleep(5) //TODO: Do we need 5 secs sleep ?
	}
	sleep(120)
	getBlocksHeight()

	qVal1, qVal2 := queryChaincode()
	fmt.Printf("\n############## Query Values A=%s, B=%s", qVal1, qVal2)
	// STEP2: Stop Peer
	peernetwork.StopPeerLocal(peerNetworkSetup, "PEER3")
	_, _ = invokeChaincode();
	_,_ = queryChaincode();
	peernetwork.StartPeerLocal(peerNetworkSetup, "PEER3")
	sleep(15)

	for i :=1; i<= 12;i ++ {
		_,_ = invokeChaincode()
		sleep(60);
	}
	sleep(60)
	qVal1,qVal2 = queryChaincode();
	fmt.Printf("\n############## Query Values A=%s, B=%s", qVal1, qVal2)
	getBlocksHeight();
	/*qAPIArgs0_peer0 := []string{"example02", "query", "PEER0"}
	qAPIArgs0_peer1 := []string{"example02", "query", "PEER1"}
	qAPIArgs0_peer2 := []string{"example02", "query", "PEER2"}
	qAPIArgs0_peer3 := []string{"example02", "query", "PEER2"}
	qArgs0 := []string{"a"}

	qRes,_ := chaincode.QueryOnHost(qAPIArgs0_peer0, qArgs0)
	fmt.Printf("####### Query on Host PEER0 A=%s", qRes)
	fmt.Printf("####### Query on Host PEER1 A=%s", qRes)
	fmt.Printf("####### Query on Host PEER2 A=%s", qRes)

	qRes,err = chaincode.QueryOnHost(qAPIArgs0_peer1, qArgs0)
	qRes,err = chaincode.QueryOnHost(qAPIArgs0_peer2, qArgs0)
	if (err != nil){

	}

	// STEP3: unpause Peer
	peernetwork.StartPeerLocal(peerNetworkSetup, "PEER3")
	sleep(15)

	iAPIArgs0 := []string{"example02", "invoke", "172.17.0.3"}
	invArgs0 := []string{"a", "b", "10"}

	for i :=1; i<= 10;i ++ {
		_,_ = chaincode.InvokeOnPeer(iAPIArgs0, invArgs0)
		sleep(60);
	}

	qRes,err = chaincode.QueryOnHost(qAPIArgs0_peer0, qArgs0)
	qRes,err = chaincode.QueryOnHost(qAPIArgs0_peer1, qArgs0)
	qRes,err = chaincode.QueryOnHost(qAPIArgs0_peer2, qArgs0)
	qRes,err = chaincode.QueryOnHost(qAPIArgs0_peer3, qArgs0)*/

	/*for i :=1; i<= INVOKE_COUNT;i ++ {
		fmt.Println("############## Invoke Iteration:", i)
		_, _ = invokeChaincode();
		_, _ = queryChaincode();
		//sleep(5) //TODO: Do we need 5 secs sleep ?
	}
	//sleep(5)
	getBlocksHeight()
	peernetwork.StopPeerLocal(peerNetworkSetup, "PEER3")
	sleep(5)
	//Invoke chaincode
	for i :=1; i<= INVOKE_COUNT;i ++ {
		fmt.Println("############## Invoke/Query Iteration:", i)
		_, _ = invokeChaincode();
		_, _ = queryChaincode();
		//sleep(5) //TODO: Do we need 5 secs sleep ?
	}

	// STEP3: unpause Peer
	//unpausePeer("PEER2") // TODO: should we start/stop the peer rather ?
	peernetwork.StartPeerLocal(peerNetworkSetup, "PEER3")
	sleep(60)
	//Invoke chaincode
	for i :=1; i<= INVOKE_COUNT;i ++ {
		fmt.Println("############## Invoke/Query Iteration:", i)
		_, _ = invokeChaincode();
		_, _ = queryChaincode();
		//sleep(5) //TODO: Do we need 5 secs sleep ?
	}
	getBlocksHeight()*/
	/*queryArg1, queryArg2 = queryChaincode()

	if (invArg1 == queryArg1 && invArg2 == queryArg2){
		fmt.Printf("\n==========================> Values matches after Starting peer")
	} else {
		fmt.Printf("\n==========================> Values doesn't matches after Starting peer")
	}*/

	// TODO : Check for PEER1 block height (from /chain REST API endpoint) is different from other peers ?
	fmt.Println("######## Testcase execution DONE")
}

func Issue1331(){

	// STEP1 - Deploy chaaincode and several Invokes
	invokeValue = 1
	AVal = 100000
	BVal = 900000
	curAVal = AVal
	curBVal = BVal
	//var invArg1, invArg2, queryArg1, queryArg2 int64

	//Deploy the chaincode
	deployChaincode();
	//Multiple Invokes
	for i :=1; i<= INVOKE_COUNT;i ++ {
		fmt.Println("############## Invoke Iteration:", i)
		_, _ = invokeChaincode();
		//sleep(5) //TODO: Do we need 5 secs sleep ?
	}

	sleep(5)

	//Verify invokes by querying chaincode
	/*queryArg1, queryArg2 = queryChaincode()

	if (invArg1 == queryArg1 && invArg2 == queryArg2){
		fmt.Printf("\n==========================> Deploy , Multiple Invokes and Query Successful")
	} else {
		fmt.Printf("\n==========================> Query fialure [Invoke,Query values doesn't match ]")
	}*/
	getBlocksHeight()
	// STEP2: Pause Peer
	//pausePeer("PEER2") // TODO: should we start/stop the peer rather ?
	peernetwork.StopPeerLocal(peerNetworkSetup, "PEER3")
	sleep(10)
	//Invoke chaincode
	for i :=1; i<= INVOKE_COUNT;i ++ {
		fmt.Println("############## Invoke/Query Iteration:", i)
		_, _ = invokeChaincode();
		_, _ = queryChaincode();
		//sleep(5) //TODO: Do we need 5 secs sleep ?
	}

	// STEP3: unpause Peer
	//unpausePeer("PEER2") // TODO: should we start/stop the peer rather ?
	peernetwork.StartPeerLocal(peerNetworkSetup, "PEER3")
	sleep(60)
	// Invoke after unapuse
	//invArg1, invArg2 = invokeChaincode();
	//Invoke chaincode
	for i :=1; i<= INVOKE_COUNT;i ++ {
		fmt.Println("############## Invoke Iteration:", i)
		_, _ = invokeChaincode();
		_, _ = queryChaincode();
		//sleep(5) //TODO: Do we need 5 secs sleep ?
	}
	sleep(60)
	getBlocksHeight()
	/*peernetwork.StopPeerLocal(peerNetworkSetup, "PEER3")
	sleep(5)
	//Invoke chaincode
	for i :=1; i<= INVOKE_COUNT;i ++ {
		fmt.Println("############## Invoke/Query Iteration:", i)
		_, _ = invokeChaincode();
		_, _ = queryChaincode();
		//sleep(5) //TODO: Do we need 5 secs sleep ?
	}

	// STEP3: unpause Peer
	//unpausePeer("PEER2") // TODO: should we start/stop the peer rather ?
	peernetwork.StartPeerLocal(peerNetworkSetup, "PEER3")
	sleep(60)
	//Invoke chaincode
	for i :=1; i<= INVOKE_COUNT;i ++ {
		fmt.Println("############## Invoke/Query Iteration:", i)
		_, _ = invokeChaincode();
		_, _ = queryChaincode();
		//sleep(5) //TODO: Do we need 5 secs sleep ?
	}
	getBlocksHeight()
	/*queryArg1, queryArg2 = queryChaincode()

	if (invArg1 == queryArg1 && invArg2 == queryArg2){
		fmt.Printf("\n==========================> Values matches after Starting peer")
	} else {
		fmt.Printf("\n==========================> Values doesn't matches after Starting peer")
	}*/

	// TODO : Check for PEER1 block height (from /chain REST API endpoint) is different from other peers ?
	fmt.Println("######## Testcase execution DONE")
}

/**
 *  1. Deploy Chaincode example02
 *   	a. Invoke
 * 	b. Query
 *  2. Pause Peer A
 *  3. Invoke, Query on Peer B
 *  4. Deploy chaincode 3
 *	a. Invoke
 *	b. Query
 *  5. Unpause PEER A
 *  6. PEER A -- > Invoke , Query – Chaincode 2
 *		  -- > Invoke , Query – Chaincode 3
 **/
func SyncTest(){
	// Change values accordingly
	invokeValue = 1
	AVal = 100000
	BVal = 900000
	curAVal = AVal
	curBVal = BVal
	var invArg1, invArg2, queryArg1, queryArg2 int64

	// STEP1 - Deploy, Invoke and Query on chaincode

	//Deploy the chaincode
	deployChaincode();

	//Invoke chaincode
	invArg1, invArg2 = invokeChaincode();
	sleep(5) //TODO: Do we need 5 secs sleep ?

	//Query chaincode
	queryArg1, queryArg2 = queryChaincode()

	if (invArg1 == queryArg1 && invArg2 == queryArg2){
		fmt.Printf("\n==========================> Deploy , Incvoke and Query Successful")
	} else {
		//TODO: Should we exit here ?
		fmt.Printf("\n==========================> Query fialure")
	}

	// STEP2: Pause Peer
	pausePeer("PEER1")

	for i :=1; i<= 10;i ++ {
		invArg1, invArg2 = invokeChaincode();
		sleep(5) //TODO: Do we need 5 secs sleep ?
		queryArg1, queryArg2 = queryChaincode()
		if (invArg1 == queryArg1 && invArg2 == queryArg2){
			fmt.Printf("\n==========================> Iteration %d is Successful",i)
		} else {
			fmt.Printf("\n==========================> Iteration %d is Failed",i)
		}
	}

	unpausePeer("PEER1")

	fmt.Println("######## repeate Invokes on chaincode for 2 mins")
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)

	fmt.Printf("\n################# %s took %s \n", name, elapsed)
	fmt.Println("################# Execution Completed #################")
}
