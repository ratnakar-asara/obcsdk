package main

import (
	"fmt"
	"strconv"
	"time"

	"obcsdk/chaincode"
	"obcsdk/peernetwork"
)

func main() {

	var MyNetwork peernetwork.PeerNetwork

	fmt.Println("Creating a local docker network")
	peernetwork.SetupLocalNetwork(5, true)

	time.Sleep(10000 * time.Millisecond)
	peernetwork.PrintNetworkDetails()
	MyNetwork = chaincode.InitNetwork()
	chaincode.InitChainCodes()
	//chaincode.Init()
	chaincode.RegisterUsers()

	//get a URL details to get info n chainstats/transactions/blocks etc.
	aPeer, _ := peernetwork.APeer(chaincode.ThisNetwork)
	url := "http://" + aPeer.PeerDetails["ip"] + ":" + aPeer.PeerDetails["port"]

        //barry


	//does not on localhosts
	//fmt.Println("Peers on network ")
	chaincode.NetworkPeers(url)

	//chaincode.User_Registration_Status(url, "jim")

	//fmt.Println("Blockchain: GET Chain  ....")
	chaincode.Chain_Stats(url)

	//chaincode.User_Registration_Status(url, "lukas")

	//chaincode.User_Registration_Status(url, "nishi")

	//chaincode.User_Registration_ecertDetail(url, "lukas")

	fmt.Println("\nPOST/Chaincode: Deploying chaincode at the beginning ....")
	dAPIArgs0 := []string{"example02", "init"}
	depArgs0 := []string{"a", "100000", "b", "90000"}
	chaincode.Deploy(dAPIArgs0, depArgs0)
	//fmt.Println("From Deploy error ", err)

	var resa, resb string
	var inita, initb, curra, currb, j, resaI, resbI int
	inita = 100000
	initb = 90000
	curra = inita
	currb = initb

	time.Sleep(60000 * time.Millisecond);
	fmt.Println("\nPOST/Chaincode: Querying a and b after deploy >>>>>>>>>>> ")
	qAPIArgs0 := []string{"example02", "query"}
	qArgsa := []string{"a"}
	qArgsb := []string{"b"}
	A, _ := chaincode.Query(qAPIArgs0, qArgsa)
	B, _ := chaincode.Query(qAPIArgs0, qArgsb)
	myStr := fmt.Sprintf("\nA = %s B= %s", A,B)
	fmt.Println(myStr)


	fmt.Println("******************************")
	fmt.Println("PAUSING PEER1 and PEER2 .. To Test Consensus")
	fmt.Println("******************************")

	peersToStartStop := []string{"PEER1", "PEER3", "PEER4"}
	peernetwork.PausePeersLocal(MyNetwork, peersToStartStop)

	j = 0
	for j < 1500 {
		iAPIArgs0 := []string{"example02", "invoke"}
		invArgs0 := []string{"a", "b", "1"}
		invRes, _ := chaincode.Invoke(iAPIArgs0, invArgs0)
		fmt.Println("\nFrom Invoke invRes ", invRes)
		curra = curra - 1
		currb = currb + 1

/*******************************************
		fmt.Println("******************************")
		fmt.Println("\nPOST/Chaincode : UNPAUSE ONE PEER >>>>>>>>>>> Transactions are queued ")
		fmt.Println("******************************")

		fmt.Println("UNPAUSING PEER1 ONLY.. To Test Consensus")
		peernetwork.StartPeerLocal(MyNetwork, "PEER1")
		time.Sleep(30000 * time.Millisecond)

		resa, _ = chaincode.Query(qAPIArgs0, qArgsa)
		resb, _ = chaincode.Query(qAPIArgs0, qArgsb)

		resaI, _ = strconv.Atoi(resa)
		resbI, _ = strconv.Atoi(resb)
		if (curra == resaI) && (currb == resbI) {
			fmt.Println("Results in a and b after bringing up peers match :")
			valueStr := fmt.Sprintf(" curra : %d, currb : %d, resa : %d , resb : %d", curra, currb, resaI, resbI)
			fmt.Println(valueStr)
		} else {
			fmt.Println("******************************")
			fmt.Println("RESULTS DO NOT MATCH ")
			fmt.Println("******************************")
		}

		fmt.Println("PAUSING PEER1, ONLY.. To Test Consensus")
		peernetwork.StopPeerLocal(MyNetwork, "PEER1")
		//time.Sleep(30000 * time.Millisecond)
************************************************/
		j++

	}
	fmt.Println("UNPAUSING PEER1, ... To Test Consensus STATE TRANSFER")
	peernetwork.UnpausePeerLocal(MyNetwork, "PEER1")
	fmt.Println("Sleeping for 2 minutes for PEER1 to sync up - state transfer")
	//time.Sleep(240000 * time.Millisecond)

	//fmt.Println("\nPOST/Chaincode: Querying a and b after invoke >>>>>>>>>>> ")
	//qAPIArgs01 := []string{"example02", "query", "PEER1"}
	qAPIArgs00 := []string{"example02", "query", "PEER0"}
	//qArgsa = []string{"a"}
	//qArgsb = []string{"b"}

	resa, _ = chaincode.QueryOnHost(qAPIArgs00, qArgsa)
	resb, _ = chaincode.QueryOnHost(qAPIArgs00, qArgsb)

	resaI, _ = strconv.Atoi(resa)
	resbI, _ = strconv.Atoi(resb)

	if (curra == resaI) && (currb == resbI) {
		fmt.Println("Results in a and b after bringing up peers match :")
		valueStr := fmt.Sprintf(" curra : %d, currb : %d, resa : %d , resb : %d", curra, currb, resaI, resbI)
		fmt.Println(valueStr)
	} else {
		fmt.Println("******************************")
		fmt.Println("RESULTS DO NOT MATCH AS EXPECTED ON VP2")
		fmt.Println("******************************")
	}

}

/***********************
  fmt.Println("\nBlockchain: Get Chain  ....")
  chaincode.Chain_Stats(url)


	fmt.Println("\nBlockchain: GET Chain  ....")
	response2 := chaincode.Monitor_ChainHeight(url)

	fmt.Println("\nChain Height", chaincode.Monitor_ChainHeight(url))

	fmt.Println("\nBlock: GET/Chain/Blocks/")
	chaincode.Block_Stats(url, response2)

	//time.Sleep(80000 * time.Millisecond);

	fmt.Println("\nTransactions: GET/transactions/" + invRes)
	chaincode.Transaction_Detail(url, invRes)
*********************/

//CCEx02, _, err := peernetwork.GetCCDetailByName("example02", chaincode.LibCC)
//if err != nil {
//	fmt.Println(err)
//}
/******************************
	qAPIArgs1 := []string{"example05", "query"}
	qurArgs1 := []string {CCEx02["path"], "b"}
	chaincode.Query(qAPIArgs1, qurArgs1)
        ************************************/

//deployUsingTagName()

//peernetwork.TearDownRemoteNetwork()
