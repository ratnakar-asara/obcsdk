package main

import (
	"fmt"
	"github.com/hyperledger/fabric/obc-test/obcsdk/chaincode"
	"github.com/hyperledger/fabric/obc-test/obcsdk/peernetwork"
	"time"
)

func main() {

	peernetwork.PrintNetworkDetails()
	chaincode.Init()
	chaincode.RegisterUsers()

	//get a URL details to get info n chainstats/transactions/blocks etc.
	aPeer := peernetwork.APeer(chaincode.ThisNetwork)
	url := "http://" + aPeer.PeerDetails["ip"] + ":" + aPeer.PeerDetails["port"]

	//does not work when launching a peer locally
	//fmt.Println("Peers on network ")
	chaincode.Network_Peers(url)

	fmt.Println("Blockchain: GET Chain  ....")
	chaincode.Chain_Stats(url)

	chaincode.User_Registration_Status(url, "lukas")

	chaincode.User_Registration_Status(url, "nishi")

	chaincode.User_Registration_ecertDetail(url, "lukas")

	fmt.Println("\nPOST/Chaincode: Deploying chaincode at the beginning ....")
	dAPIArgs0 := []string{"example02", "init"}
	depArgs0 := []string{"a", "20000", "b", "9000"}
	chaincode.Deploy(dAPIArgs0, depArgs0)
	//fmt.Println("From Deploy error ", err)

	time.Sleep(20000 * time.Millisecond)

	fmt.Println("\nPOST/Chaincode : Querying a and b after a deploy  ")
	qAPIArgs0 := []string{"example02", "query"}
	qArgsa := []string{"a"}
	_, _ = chaincode.Query(qAPIArgs0, qArgsa)
	qArgsb := []string{"b"}
	_, _ = chaincode.Query(qAPIArgs0, qArgsb)

	fmt.Println("\nPOST/Chaincode : Invoke on a and b after a deploy >>>>>>>>>>> ")
	iAPIArgs0 := []string{"example02", "invoke"}
	invArgs0 := []string{"a", "b", "500"}
	invRes, _ := chaincode.Invoke(iAPIArgs0, invArgs0)
	fmt.Println("\nFrom Invoke invRes ", invRes)

	fmt.Println("Sleeping 5secs for invoke to complete on ledger")

	time.Sleep(5000 * time.Millisecond)

	fmt.Println("\nBlockchain: Get Chain  ....")
	chaincode.Chain_Stats(url)

	fmt.Println("\nPOST/Chaincode: Querying a and b after invoke >>>>>>>>>>> ")
	_, _ = chaincode.Query(qAPIArgs0, qArgsa)
	_, _ = chaincode.Query(qAPIArgs0, qArgsb)

	fmt.Println("\nBlockchain: GET Chain  ....")
	response2 := chaincode.Monitor_ChainHeight(url)

	fmt.Println("\nChain Height", chaincode.Monitor_ChainHeight(url))

	fmt.Println("\nBlock: GET/Chain/Blocks/")
	chaincode.Block_Stats(url, response2)

	fmt.Println("\nBlockchain: Getting Transaction detail for   ....", invRes)

	time.Sleep(50000 * time.Millisecond)

	//fmt.Println("\nTransactions: GET/transactions/" + invRes)
	chaincode.Transaction_Detail(url, invRes)

	fmt.Println("\nBlockchain: GET Chain .... ")
	time.Sleep(10000 * time.Millisecond)
	chaincode.Chain_Stats(url)

	/*** let's call deploy with tagName */
	deployUsingTagName()

}
