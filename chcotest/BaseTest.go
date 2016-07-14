package main

import (
	"fmt"
	"time"

	"obcsdk/chaincode"
	"obcsdk/peernetwork"
)

var peerNetworkSetup peernetwork.PeerNetwork

func setupNetwork() {
	fmt.Println("Creating a local docker network")
	peernetwork.SetupLocalNetwork(4, true)
	peernetwork.PrintNetworkDetails()
	peerNetworkSetup = chaincode.InitNetwork()
	chaincode.InitChainCodes()
	chaincode.RegisterUsers()
}
func main() {
	setupNetwork()

	fmt.Println("\nPOST/Chaincode: Deploying chaincode at the beginning ....")
	dAPIArgs0 := []string{"mycc", "init"}
	depArgs0 := []string{"a", "ihglkjfdjkghfdlkjhgk", "b", "90000"}
	chaincode.Deploy(dAPIArgs0, depArgs0)

	time.Sleep(time.Second * 60)
	fmt.Println("\nPOST/Chaincode: Querying a and b after deploy >>>>>>>>>>> ")
	qAPIArgs0 := []string{"mycc", "query"}
	qArgsa := []string{"a"}
	qArgsb := []string{"b"}
	A, _ := chaincode.Query(qAPIArgs0, qArgsa)
	B, _ := chaincode.Query(qAPIArgs0, qArgsb)
	myStr := fmt.Sprintf("\nA = %s B= %s", A, B)
	fmt.Println(myStr)

}
