package chaincode

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"obcsdk/peernetwork"
)

var ThisNetwork peernetwork.PeerNetwork
var Peers = ThisNetwork.Peers
var ChainCodeDetails, Versions map[string]string
var LibCC peernetwork.LibChainCodes

const invokeOnPeerUsage = ("iAPIArgs0 := []string{\"example02\", \"invoke\", \"<PEER_IP_ADDRESS>\" + \"(optional)<tagName>\"}" +
	"invArgs0 := []string{\"a\", \"b\", \"500\"} " +

	"chaincode.Invoke(iAPIArgs0, invArgs0)}")
const invokeAsUserUsage = ("\niAPIArgs0 := []string{\"example02\", \"invoke\", \"<Registered_USER_NAME>\" + \"(optional)<tagName>\"}" +
	"invArgs0 := []string{\"a\", \"b\", \"500\"} " +

	"chaincode.Invoke(iAPIArgs0, invArgs0)}")

/**
  initializes users on network using data supplied in NetworkCredentials.json file
*/
func InitNetwork() peernetwork.PeerNetwork {

	ThisNetwork = peernetwork.LoadNetwork()
	return ThisNetwork
}

/**
   initializes chaincodes on network using information supplied in CC_Collections.json file
*/
func InitChainCodes() {
	LibCC = peernetwork.InitializeChainCodes()
}

/*
  initializes network based on files in directory utils
*/
func Init() {
	InitNetwork()
	InitChainCodes()
}

/*
   Registers each user on the network based on the content of ThisNetwork.Peers.
*/
func RegisterUsers() {
	fmt.Println("\nCalling Register ")

	//testuser := peernetwork.AUser(ThisNetwork)
	Peers := ThisNetwork.Peers
	i := 0
	for i < len(Peers) {

		userList := ThisNetwork.Peers[i].UserData
		for user, secret := range userList {
			url := "https://" + Peers[i].PeerDetails["ip"] + ":" + Peers[i].PeerDetails["port"]
			msgStr := fmt.Sprintf("\nRegistering %s with password %s on %s using %s", user, secret, Peers[i].PeerDetails["name"], url)
			fmt.Println(msgStr)
			register(url, user, secret)
		}
		fmt.Println("Done Registering ", len(userList), "users on ", Peers[i].PeerDetails["name"])
		i++
	}
}
func RegisterCustomUsers() {
	fmt.Println("\nCalling RegisterCustomUsers ")

	Peers := ThisNetwork.Peers

	for i := 0; i < len(Peers) ; i++ {

		userList := ThisNetwork.Peers[i].UserData
		for user, secret := range userList {
			url := "https://" + Peers[i].PeerDetails["ip"] + ":" + Peers[i].PeerDetails["port"]
			msgStr := fmt.Sprintf("\nRegistering %s with password %s on %s using %s", user, secret, Peers[i].PeerDetails["name"], url)
			fmt.Println(msgStr)
			register(url, user, secret)
			if i == len(Peers)-1 {
				//if
				user = "dashboarduser_type0_efeeb83216"
				secret = "12211933b3"
				//TODO: Remove the hardcoding
				msgStr = fmt.Sprintf("\nRegistering %s with password %s on %s using %s", user, secret, Peers[i].PeerDetails["name"], url)
				fmt.Println(msgStr)
				register(url, user, secret)
				user = "dashboarduser_type0_fa08214e3b"
				secret = "460c3190dc"
				//TODO: Remove the hardcoding
				msgStr = fmt.Sprintf("\nRegistering %s with password %s on %s using %s", user, secret, Peers[i].PeerDetails["name"], url)
				fmt.Println(msgStr)
				register(url, user, secret)
				user = "dashboarduser_type0_e00e125cf9"
				secret = "fe1a324f86"
				//TODO: Remove the hardcoding
				msgStr = fmt.Sprintf("\nRegistering %s with password %s on %s using %s", user, secret, Peers[i].PeerDetails["name"], url)
				fmt.Println(msgStr)
				register(url, user, secret)
				user = "dashboarduser_type0_e0ee60d5af"
				secret = "bc6911cfd0"
			}
		}
		fmt.Println("Done Registering ", len(userList), "users on ", Peers[i].PeerDetails["name"])
	}
}

func RegisterUsers2() {
	fmt.Println("\nCalling Register ")

	//testuser := peernetwork.AUser(ThisNetwork)
	Peers := ThisNetwork.Peers
	for i:= 0;i < len(Peers)-2;i++ {

		userList := ThisNetwork.Peers[i].UserData
		for user, secret := range userList {
			url := "https://" + Peers[i].PeerDetails["ip"] + ":" + Peers[i].PeerDetails["port"]
			msgStr := fmt.Sprintf("\nRegistering %s with password %s on %s using %s", user, secret, Peers[i].PeerDetails["name"], url)
			fmt.Println(msgStr)
			register(url, user, secret)
		}
		fmt.Println("Done Registering ", len(userList), "users on ", Peers[i].PeerDetails["name"])
	}
}
/*
   deploys a chaincode in the fabric to later execute functions on this deployed chaincode
   Takes two arguments
 	 A. args []string
	   	1.ccName (string)			- name of the chaincode as specified in CC_Collections.json file
		2.funcName (string)			- name of the function to call from chaincode specification
									"init" for chaincodeexample02
		3.tagName(string)(optional)		- tag a deployment to support something like versioning

 	B. depargs []string				- actual arguments passed to initialize chaincode inside the fabric.

		Sample Code:
		dAPIArgs0 := []string{"example02", "init"}
		depArgs0 := []string{"a", "20000", "b", "9000"}

		var depRes string
		var err error
		depRes, err := chaincode.Deploy(dAPIArgs0, depArgs0)
*/
func Deploy(args []string, depargs []string) error {

	if (len(args) < 2) || (len(args) > 3) {
		return errors.New("Deploy : Incorrect number of arguments. Expecting 2 or 3")
	}
	ccName := args[0]
	funcName := args[1]
	var tagName string
	if len(args) == 2 {
		tagName = ""
	} else if len(args) == 3 {
		tagName = args[2]
	}
	dargs := depargs
	var err error
	ChainCodeDetails, Versions, err = peernetwork.GetCCDetailByName(ccName, LibCC)
	if err != nil {
		fmt.Println("Inside deploy: ", err)
		//log.Fatal("No Chain Code Details, we cannot proceed")
		return errors.New("No Chain Code Details we cannot proceed")
	}
	if strings.Contains(ChainCodeDetails["deployed"], "true") {
		fmt.Println("\n\n ** Already deployed ..")
		fmt.Println(" skipping deploy...")
	} else {
		msgStr := fmt.Sprintf("\n** Initializing and deploying chaincode %s on network with args %s\n", ChainCodeDetails["path"], dargs)
		fmt.Println(msgStr)
		restCallName := "deploy"
		peer, auser := peernetwork.AUserFromNetwork(ThisNetwork)
		//fmt.Println("Value in State : ", peer.State)
		//fmt.Println("Value in State : ", peer.PeerDetails["state"])
		//aPeer.PeerDetails["ip"], aPeer.PeerDetails["port"]
		url := "https://" + peer.PeerDetails["ip"] + ":" + peer.PeerDetails["port"]
		txId := changeState(url, ChainCodeDetails["path"], restCallName, dargs, auser, funcName)
		//storing the value of most recently deployed chaincode inside chaincode details if no tagname or versioning
		ChainCodeDetails["dep_txid"] = txId
		if len(tagName) != 0 {
			Versions[tagName] = txId
		}
	}
	return err
}

/*
 changes state of a chaincode by passing arguments to BlockChain REST API invoke.
 Takes two arguments
 	 A. args []string
	    1.ccName (string)			- name of the chaincode as specified in CC_Collections.json file
		2.funcName (string)		- name of the function to call from chaincode specification
								"invoke" for chaincodeexample02
		3.tagName(string)(optional)	- tag a deployment to support something like versioning

	B. invargs []string			- actual arguments passed to change the state of chaincode inside the fabric.

		Sample Code:
		iAPIArgs0 := []string{"example02", "invoke"}
		invArgs0 := []string{"a", "b", "500"}

		var invRes string
		var err error
		invRes,err := chaincode.Invoke(iAPIArgs0, invArgs0)}
*/
func Invoke(args []string, invokeargs []string) (id string, err error) {

	if (len(args) < 2) || (len(args) > 3) {
		fmt.Println("Invoke : Incorrect number of arguments. Expecting 2")
		return "", errors.New("Invoke : Incorrect number of arguments. Expecting 2")
	}
	ccName := args[0]
	funcName := args[1]
	var tagName string
	if len(args) == 2 {
		tagName = ""
	} else if len(args) == 3 {
		tagName = args[2]
	}
	invargs := invokeargs
	//fmt.Println("Inside invoke .....")
	var err1 error
	ChainCodeDetails, Versions, err1 = peernetwork.GetCCDetailByName(ccName, LibCC)
	if err1 != nil {
		fmt.Println("Inside invoke: ", err1)
		log.Fatal("No Chain Code Details we cannot proceed")
		return "", errors.New("No Chain Code Details we cannot proceed")
	}
	restCallName := "invoke"
	aPeer, _ := peernetwork.APeer(ThisNetwork)
	//fmt.Println(aPeer.PeerDetails["ip"], aPeer.PeerDetails["port"])
	ip, port, auser := peernetwork.AUserFromAPeer(*aPeer)
	url := "https://" + ip + ":" + port
	//msgStr0 := fmt.Sprintf("\n** Calling %s on chaincode %s with args %s on  %s as %s\n", funcName, ccName, invargs, url, auser)
	//fmt.Println(msgStr0)
	var txId string
	if len(tagName) != 0 {
		txId = changeState(url, Versions[tagName], restCallName, invargs, auser, funcName)
	} else {
		txId = changeState(url, (ChainCodeDetails["dep_txid"]), restCallName, invargs, auser, funcName)
	}
	//fmt.Println("\n\n\n*** END Invoking as  ***\n\n", auser, " on a single peer\n\n")
	return txId, errors.New("")
}

/*
 changes state of a chaincode on a specific peer by passing arguments to REST API call
 Takes two arguments
	A. args []string
	   	1. ccName (string)				- name of the chaincode as specified in CC_Collections.json file
		2. funcName(string)				- name of the function to call from chaincode specification
										"invoke" for chaincodeexample02
		3. host (string)				- hostname or ipaddress to call invoke from
		4. tagName(string)(optional)			- tag the invocation to support something like versioning

	B. invargs []string					- actual arguments passed to change the state of chaincode inside the fabric.

		Sample Code:
		iAPIArgs0 := []string{"example02", "invoke", "127.0.0.1"}
		invArgs0 := []string{"a", "b", "500"}

		var invRes string
		var err error
		invRes,err := chaincode.Invoke(iAPIArgs0, invArgs0)}
*/
func InvokeOnPeer(args []string, invokeargs []string) (id string, err error) {

	//fmt.Println("Inside InvokeOnPeer .....")
	if (len(args) < 3) || (len(args) > 4) {
		fmt.Println("InvokeOnPeer : Incorrect number of arguments. Expecting 3 or 4 in invokeAPI arguments")
		fmt.Println(invokeOnPeerUsage)
		return "", errors.New("InvokeOPeer : Incorrect number of arguments. Expecting 3 or 4 in function arguments")
	}
	ccName := args[0]
	funcName := args[1]
	host := args[2]
	var tagName string
	if len(args) == 3 {
		tagName = ""
	} else if len(args) == 4 {
		tagName = args[3]
	}
	invargs := invokeargs
        restCallName := "invoke"
	var err1 error
	var txId string
	ChainCodeDetails, Versions, err1 = peernetwork.GetCCDetailByName(ccName, LibCC)
	if err1 != nil {
		fmt.Println("Inside InvokeOnPeer: ", err1)
		log.Fatal("No Chain Code Details we cannot proceed")
		return "", errors.New("No Chain Code Details we cannot proceed")
	}

	ip, port, auser, err2 := peernetwork.AUserFromThisPeer(ThisNetwork, host)
	if err2 != nil {
		fmt.Println("Inside invoke3: ", err2)
		return "", err2
	} else {
		url := "https://" + ip + ":" + port
		//msgStr0 := fmt.Sprintf("\n** Calling %s on chaincode %s with args %s on  %s as %s on %s\n", funcName, ccName, invargs, url, auser, host)
		//fmt.Println(msgStr0)
		if (len(tagName) > 0) {
			txId = changeState(url, Versions[tagName], restCallName, invargs, auser, funcName)
		}else {
		        txId = changeState(url, (ChainCodeDetails["dep_txid"]), restCallName, invargs, auser, funcName)
		}
		return txId, errors.New("")
	}
}

/*
 changes state of a chaincode using a specific user credential
  Takes two arguments
 	A. args []string
	   	1. ccName (string)				- name of the chaincode as specified in CC_Collections.json file
		2. funcName(string)				- name of the function to call from chaincode specification
										"invoke" for chaincodeexample02
		3. user (string)				- login name of a registered user
		4. tagName(string)(optional)			- tag the invocation to support something like versioning

	B. invargs []string					- actual arguments passed to change the state of chaincode inside the fabric.

		Sample Code:
		iAPIArgs0 := []string{"example02", "invoke", "jim"}
		invArgs0 := []string{"a", "b", "500"}

		var invRes string
		var err error
		invRes,err := chaincode.Invoke(iAPIArgs0, invArgs0)}
*/
func InvokeAsUser(args []string, invokeargs []string) (id string, err error) {
	if (len(args) < 3) || (len(args) > 4) {
		fmt.Println("InvokeAsUser : Incorrect number of arguments. Expecting 3 or 4 in invokeAPI arguments")
		fmt.Println(invokeAsUserUsage)
		return "", errors.New("InvokeAsUser : Incorrect number of arguments. Expecting 3 or 4 number in invokeAPI arguments")
	}
	ccName := args[0]
	funcName := args[1]
	userName := args[2]
	var tagName string
	if len(args) == 3 {
		tagName = ""
	} else if len(args) == 4 {
		tagName = args[3]
	}
	invargs := invokeargs
	var err1 error
	ChainCodeDetails, Versions, err1 = peernetwork.GetCCDetailByName(ccName, LibCC)
	if err1 != nil {
		fmt.Println("Inside InvokeAsUser: ", err1)
		log.Fatal("No Chain Code Details we cannot proceed")
		return "", errors.New("No Chain Code Details we cannot proceed")
	}
	restCallName := "invoke"
	ip, port, auser, err2 := peernetwork.PeerOfThisUser(ThisNetwork, userName)
	if err2 != nil {
		fmt.Println("Inside InvokeAsUser: ", err2)
		return "", err2
	} else {
		url := "https://" + ip + ":" + port
		//msgStr0 := fmt.Sprintf("\n** Calling %s on chaincode %s with args %s on  %s as %s\n", funcName, ccName, invargs, url, auser)
		//fmt.Println(msgStr0)
		var txId string
		if len(tagName) > 0 {
			txId = changeState(url, Versions[tagName], restCallName, invargs, auser, funcName)
		}else {
			txId = changeState(url, ChainCodeDetails["dep_txid"], restCallName, invargs, auser, funcName)
		}
		return txId, errors.New("")
	}
}

/*
  Query fetches the value of the arguments supplied to query function from the fabric.
  Takes two arguments
 	A. args []string
	   	1. ccName (string)				- name of the chaincode as specified in CC_Collections.json file
		2. funcName(string)				- name of the function to call from chaincode specification
										"query" for chaincodeexample02
		3. tagName(string)(optional)	- tag the invocation to support something like versioning

	B. qargs []string					- actual arguments passed to get the values as stored inside fabric.

		Sample Code:
		qAPIArgs0 := []string{"example02", "query"}
		qArgsa := []string{"a"}

		var queryRes string
		var err error
		queryRes,err := chaincode.Query(qAPIArgs0, qArgsa)
*/
func Query(args []string, queryArgs []string) (id string, err error) {

	if (len(args) < 2) || (len(args) > 3) {
		return "", errors.New("Incorrect number of arguments. Expecting 2")
	}
	ccName := args[0]
	funcName := args[1]
	var tagName string
	if len(args) == 2 {
		tagName = ""
	} else if len(args) == 3 {
		tagName = args[2]
	}
	qargs := queryArgs
	var err1 error

	ChainCodeDetails, Versions, err1 = peernetwork.GetCCDetailByName(ccName, LibCC)
	if err1 != nil {
		fmt.Println("Inside Query: ", err1)
		fmt.Println("No Chain Code Details we cannot proceed")
		return "", errors.New("No Chain Code Details we cannot proceed")
	}
	restCallName := "query"
	//ip, port, auser := peernetwork.AUserFromNetwork(ThisNetwork)
	//url := "https://" + ip + ":" + port
	peer, auser := peernetwork.AUserFromNetwork(ThisNetwork)
	//fmt.Println("Value in State : ", peer.State)
	//aPeer.PeerDetails["ip"], aPeer.PeerDetails["port"]
	url := "https://" + peer.PeerDetails["ip"] + ":" + peer.PeerDetails["port"]

	var txId string
	msgStr0 := fmt.Sprintf("\n** Calling %s on chaincode %s with args %s on  %s as %s\n", funcName, ccName, queryArgs, url, auser)
	fmt.Println(msgStr0)

	if len(tagName) != 0 {
		txId = readState(url, Versions[tagName], restCallName, qargs, auser, funcName)
	} else {
		txId = readState(url, (ChainCodeDetails["dep_txid"]), restCallName, qargs, auser, funcName)
	}

	return txId, errors.New("")
}


/*
/*
  Query fetches the value of the arguments supplied to query function from the fabric.
  Takes two arguments
 	A. args []string
	  1. ccName (string)				- name of the chaincode as specified in CC_Collections.json file
		2. funcName(string)				- name of the function to call from chaincode specification
		3. host (string)				- hostname or ipaddress to call query
		4. tagName(string)(optional)	- tag the invocation to support something like versioning

	B. qargs []string					- actual arguments passed to get the values as stored inside fabric.

		Sample Code:
		qAPIArgs0 := []string{"example02", "query", "vp2"}
		qArgsa := []string{"a"}

		var queryRes string
		var err error
		queryRes,err := chaincode.Query(qAPIArgs0, qArgsa)
*/



func QueryOnHost(args []string, queryargs []string) (id string, err error) {
	if (len(args) < 3) || (len(args) > 4) {
		fmt.Println("QueryOnHost : Incorrect number of arguments. Expecting 3 or 4 in invokeAPI arguments")
		fmt.Println(invokeOnPeerUsage)
		return "", errors.New("QueryOnHost : Incorrect number of arguments. Expecting 3 or 4 in function arguments")
	}
	ccName := args[0]
	funcName := args[1]
	host := args[2]
	var tagName string
	if len(args) == 3 {
		tagName = ""
	} else if len(args) == 4 {
		tagName = args[3]
	}
	qryargs := queryargs
	var err1 error
	var txId string
	ChainCodeDetails, Versions, err1 = peernetwork.GetCCDetailByName(ccName, LibCC)
	if err1 != nil {
		fmt.Println("Inside QueryOnHost: ", err1)
		log.Fatal("No Chain Code Details we cannot proceed")
		return "", errors.New("No Chain Code Details we cannot proceed")
	}
	restCallName := "query"
	ip, port, auser, err2 := peernetwork.AUserFromThisPeer(ThisNetwork, host)
	if err2 != nil {
		fmt.Println("Inside Query: ", err2)
		return "", err2
	} else {
		url := "https://" + ip + ":" + port
		//msgStr0 := fmt.Sprintf("\n** Calling %s on chaincode %s with args %s on  %s as %s on %s\n", funcName, ccName, qryargs, url, auser, host)
		//fmt.Println(msgStr0)
		if (len(tagName) > 0) {
			txId = changeState(url, Versions[tagName], restCallName, qryargs, auser, funcName)
		}else {
			txId = changeState(url, (ChainCodeDetails["dep_txid"]), restCallName, qryargs, auser, funcName)
		}
		return txId, errors.New("")
	}

}

func GetChainHeight(host string) (ht int, err error) {

			fmt.Println("Inside GetChainHeight chcoAPI.....")
			ip, port, _, err2 := peernetwork.AUserFromThisPeer(ThisNetwork, host)
			if err2 != nil {
				fmt.Println("Inside GetChainHeight: ", err2)
				return -1, err2
			} else {
				url := "https://" + ip + ":" + port
				ht := Monitor_ChainHeight(url)
				return ht, errors.New("")
			}

}
