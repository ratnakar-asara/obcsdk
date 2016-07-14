/**** TEST OBJECTIVE: Consensus: Threshold not met *************************************
*   Setup: 5 node local docker peer network with security
*   0. Deploy chaincodeexample02 with 100, 200 as initial args
*   1. PAUSE PEER1, PEER2 PEER3   FOR CONSENSUS NOT TO HAPPEN on a 5 peer network
*   2. Send ONE INVOKE REQUESTS
*   3. Query for A and B on unpaused peer(PEER0) to get initial values for A and B as 100, 200 respectivley
*   4. Unpause PEER1
*   5. Assuming Consensus happens
*   8. Do A Query ON PEER1
*   9. Get updated results after three invokes on PEER1 in step8
****************************************/

package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/obcsdk/chaincode"
	"github.com/hyperledger/fabric/obcsdk/peernetwork"
)

func main() {

	var MyNetwork peernetwork.PeerNetwork

	fmt.Println("Using an existing local docker network")
	peernetwork.PrintNetworkDetails()
	MyNetwork = chaincode.InitNetwork()
	chaincode.InitChainCodes()
	//chaincode.Init()
	chaincode.RegisterUsers()

	//get a URL details to get info n chainstats/transactions/blocks etc.
	aPeer, _ := peernetwork.APeer(chaincode.ThisNetwork)
	url := "http://" + aPeer.PeerDetails["ip"] + ":" + aPeer.PeerDetails["port"]

	//does not on localhosts
	//fmt.Println("Peers on network ")
	chaincode.NetworkPeers(url)

	fmt.Println("\nPOST/Chaincode: Deploying chaincode at the beginning ....")
	dAPIArgs0 := []string{"example02", "init"}
	depArgs0 := []string{"a", "1000", "b", "2000"}
	chaincode.Deploy(dAPIArgs0, depArgs0)
	//fmt.Println("From Deploy error ", err)

	var resa, resb string
	var inita, initb, curra, currb, j, resaI, resbI int
	inita = 1000
	initb = 2000
	curra = inita
	currb = initb

	time.Sleep(60000 * time.Millisecond)
	fmt.Println("\nPOST/Chaincode: Querying a and b after invoke >>>>>>>>>>> ")
	qAPIArgs0 := []string{"example02", "query"}
	qArgsa := []string{"a"}
	qArgsb := []string{"b"}
	chaincode.Query(qAPIArgs0, qArgsa)
	chaincode.Query(qAPIArgs0, qArgsb)

	fmt.Println("******************************")
	fmt.Println("PAUSING PEER1 and PEER2 PEER3 .. To Test Consensus")
	fmt.Println("******************************")

	peersToStartStop := []string{"PEER1"}
	peernetwork.PausePeersLocal(MyNetwork, peersToStartStop)

	j = 0
	for j < 1 {
		iAPIArgs0 := []string{"example02", "invoke"}
		invArgs0 := []string{"a", "b", "1"}
		invRes, _ := chaincode.Invoke(iAPIArgs0, invArgs0)
		fmt.Println("\nFrom Invoke invRes ", invRes)
		curra = curra - 1
		currb = currb + 1
	}

	//fmt.Println("\nPOST/Chaincode: Querying a and b after invoke >>>>>>>>>>> ")
	time.Sleep(60000 * time.Millisecond)
	qAPIArgs00 := []string{"example02", "query", "PEER0"}
	resa, _ = chaincode.QueryOnHost(qAPIArgs00, qArgsa)
	resb, _ = chaincode.QueryOnHost(qAPIArgs00, qArgsb)

	resaI, _ = strconv.Atoi(resa)
	resbI, _ = strconv.Atoi(resb)

	if (curra == resaI) && (currb == resbI) {
		fmt.Println("Results in a and b after bringing up peers match on PEER0: TEST FAIL")
		valueStr := fmt.Sprintf(" curra : %d, currb : %d, resa : %d , resb : %d", curra, currb, resaI, resbI)
		fmt.Println(valueStr)
	} else {
		fmt.Println("******************************")
		fmt.Println("RESULTS DO NOT MATCH AS EXPECTED ON PEER0, since three PEERS ARE DOWN")
		fmt.Println("******************************")
	}

	fmt.Println("UNPAUSING PEER1, ... To Test Consensus STATE TRANSFER")
	peernetwork.UnpausePeerLocal(MyNetwork, "PEER1")
	fmt.Println("Sleeping for 2 minutes for PEER1 to sync up - state transfer")
	time.Sleep(120000 * time.Millisecond)

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
		fmt.Println("RESULTS DO ... MATCH AS EXPECTED ON PEER0, since ONLY TWO PEERS ARE DOWN, CONSENSUS CAN TAKE PLACE")
	}

}
