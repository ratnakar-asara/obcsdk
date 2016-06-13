package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/obcsdk/peernetwork"
	"github.com/hyperledger/fabric/obcsdk/chaincode"


)

func main() {

  var MyNetwork peernetwork.PeerNetwork

	fmt.Println("Creating a remote network")
	peernetwork.SetupRemoteNetwork(4)

	time.Sleep(10000 * time.Millisecond);
	peernetwork.PrintNetworkDetails()
  MyNetwork = chaincode.InitNetwork()
  chaincode.InitChainCodes()
	//chaincode.Init()
	chaincode.RegisterUsers()

	//get a URL details to get info n chainstats/transactions/blocks etc.
	aPeer,_ := peernetwork.APeer(chaincode.ThisNetwork)
 	url := "http://" + aPeer.PeerDetails["ip"] + ":" + aPeer.PeerDetails["port"]

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

  var inita, initb, curra, currb, i int
	inita = 100000
	initb = 90000
	curra = inita
	currb = initb

	i = 0
  for ( i < 1000)  {

		//time.Sleep(60000 * time.Millisecond);
    fmt.Println("******************************")
		fmt.Println("LOOOP :", i)
		fmt.Println("******************************")

		fmt.Println("\nBlockchain: Get Chain .......... ....")
		chaincode.Chain_Stats(url)

		fmt.Println("******************************")
	  fmt.Println("PAUSING VP1 and VP2 .. To Test Consensus")
		fmt.Println("******************************")

		peersToStartStop := []string{"vp1", "vp2"}
	  peernetwork.StopPeers(MyNetwork, peersToStartStop)
	  //peerDetails1, _ := peernetwork.GetPeerState(MyNetwork, "vp1")
	  //peerDetails2, _ := peernetwork.GetPeerState(MyNetwork, "vp2")
	  //fmt.Println("curstates VP1 is : ",  peerDetails1.State)
	  //fmt.Println("curstates VP2 is : ",  peerDetails2.State)

		fmt.Println("\nPOST/Chaincode : Querying a and b after a deploy  with two peers... should pass")
		qAPIArgs0 := []string{"example02", "query"}
		qArgsa := []string{"a"}
		resa, _  := chaincode.Query(qAPIArgs0, qArgsa)
		qArgsb := []string{"b"}
		resb, _  := chaincode.Query(qAPIArgs0, qArgsb)

	  curra, _ = strconv.Atoi(resa)
		currb, _ = strconv.Atoi(resb)
		fmt.Println("******************************")
		valueStr := fmt.Sprintf("After deploy values in  curra : %d currb : %d", curra, currb)
		fmt.Println(valueStr)
		fmt.Println("******************************")

		fmt.Println("******************************")
		fmt.Println("\nPOST/Chaincode : Invoke on a and b after STOPPING TWO PEERS >>>>>>>>>>> Transactions are queued ")
		fmt.Println("******************************")

		iAPIArgs0 := []string{"example02", "invoke"}
		invArgs0 := []string{"a", "b", "1"}
		invRes, _ := chaincode.Invoke(iAPIArgs0, invArgs0)
		fmt.Println("\nFrom Invoke invRes ", invRes)
	  curra=curra-1
		currb= currb+1

		fmt.Println("******************************")
	  fmt.Println("\nPOST/Chaincode : UNPAUSE TWO PEERS >>>>>>>>>>> Transactions are queued ")
		fmt.Println("******************************")

		fmt.Println("UNPAUSING VP1, VP2 .. To Test Consensus")
		peernetwork.StartPeers(MyNetwork, peersToStartStop)
		time.Sleep(120000 * time.Millisecond);

		fmt.Println("\nPOST/Chaincode: Querying a and b after invoke >>>>>>>>>>> ")
		qAPIArgs0 = []string{"example02", "query"}
		qArgsa = []string{"a"}
		qArgsb = []string{"b"}

		resa, _  = chaincode.Query(qAPIArgs0, qArgsa)
		resb, _  = chaincode.Query(qAPIArgs0, qArgsb)

    resaI, _ :=  strconv.Atoi(resa)
		resbI, _ :=  strconv.Atoi(resb)
   if ( (curra == resaI) && (currb == resbI) ) {
	 		fmt.Println("Results in a and b after bringing up peers match :")
			valueStr := fmt.Sprintf(" curra : %d, currb : %d, resa : %d , resb : %d", curra, currb, resaI, resbI)
			fmt.Println(valueStr)
		}else {
			fmt.Println("******************************")
			fmt.Println("RESULTS DO NOT MATCH AS EXPECTED")
			fmt.Println("******************************")

		}
	  i++
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
