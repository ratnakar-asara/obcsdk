package peernetwork

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
    "os/exec"
	"strconv"

	//"strings"
	//"github.com/pkg/sftp"
	//"golang.org/x/crypto/ssh"
)

const (
	RUNNING      = 0
	STOPPED      = 1
	STARTED      = 2
	PAUSED       = 3
	UNPAUSED     = 4
	NOTRESPONDIN = 5
)

type Peer struct {
	PeerDetails map[string]string
	UserData    map[string]string
	State       int
}

type PeerNetwork struct {
	Peers []Peer
	Name  string
}

type LibChainCodes struct {
	ChainCodes map[string]ChainCode
}

type ChainCode struct {
	Detail   map[string]string
	Versions map[string]string
}

var peerNetwork PeerNetwork

const USER = "ibmadmin"
const PASSWORD = "m0115jan"

//HOST = "urania"
const IP = "9.37.136.147"

//NEW_IP = "9.42.91.158"

func SetupLocalNetwork(numPeers int, sec bool){

    var cmd *exec.Cmd
    goroot := os.Getenv("GOROOT")
    pwd, _ := os.Getwd()
    fmt.Println("Initially ", pwd)
    os.Chdir(pwd + "/../automation/")
    pwd, _ = os.Getwd()
    fmt.Println("After change dir ", pwd)
    script := pwd + "/local_fabric.sh"
    arg0 :=  "-n"
    arg1 := strconv.Itoa(numPeers)
    arg2 := "-s"

    //cmdStr := script + arg0 + arg1 + arg2
    cmdStr := script + arg0 + arg1
    fmt.Println("cmd ", cmdStr)
    //cmd := exec.Command("/bin/bash", cmdStr)
    //cmd := exec.Command("sudo", script, arg0, arg1, arg2)
    cmd = exec.Command("sudo", script, arg0, arg1, arg2, "")

    var stdoutBuf bytes.Buffer
    cmd.Stdout = &stdoutBuf
     err := cmd.Run()
     if err != nil {
        log.Fatal(err)
     }
    fmt.Printf("in all caps: \n", stdoutBuf.String())

    GetNC_Local()
    os.Chdir(goroot + "/src/obcsdk/chcotest")
    pwd, _ = os.Getwd()
    fmt.Println("After change back dir ", pwd)
}

func GetNC_Local() {
        //goroot := os.Getenv("GOROOT")
        //pwd, _ := os.Getwd()
        //fileName := pwd + "/networkcredentials"
        fileName := "../automation/networkcredentials"

        inputfile, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	outfile, err := os.Create("../util/NetworkCredentials.json")
	if err != nil {
		fmt.Println("Error in creating NetworkCredentials file ", err)
	}

	_, err = io.Copy(outfile, inputfile)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Copied contents")
	//log.Println(inputfile)

	outfile.Close()

}

/*
  creates network as defined in NetworkCredentials.json, distributing users evenly among the peers of the network.
*/
func LoadNetwork() PeerNetwork {

	//p, n := initializePeers()

	p, n := initializePeers()

	peerNetwork := PeerNetwork{Peers: p, Name: n}
	return peerNetwork
}

/*
  reads CC_Collection.json and returns a library of chain codes.
*/
func InitializeChainCodes() LibChainCodes {
	pwd, _ := os.Getwd()
	file, err := os.Open(pwd+"/../util/CC_Collection.json")
	if err != nil {
		log.Fatal("Error in opening CC_Collection.json file ")
	}

	poolChainCode, err := unmarshalChainCodes(file)
	if err != nil {
		log.Fatal("Error in unmarshalling")
	}

	//make a map to hold each chaincode detail
	ChCos := make(map[string]ChainCode)
	for i := 0; i < len(poolChainCode); i++ {
		//construct a map for each chaincode detail
		detail := make(map[string]string)
		detail["type"] = poolChainCode[i].TYPE
		detail["path"] = poolChainCode[i].PATH
		//detail["dep_txid"] = poolChainCode[i].DEP_TXID
		//detail["deployed"] = poolChainCode[i].DEPLOYED

		versions := make(map[string]string)
		CC := ChainCode{Detail: detail, Versions: versions}
		//add the structure to map of chaincodes
		ChCos[poolChainCode[i].NAME] = CC
	}
	//finally add this map - collection of chaincode detail to library of chaincodes
	libChainCodes := LibChainCodes{ChainCodes: ChCos}
	return libChainCodes
}

func initializePeers() (peers []Peer, name string) {

	fmt.Println("Getting and Initializing Peer details from network")
	peerDetails, userDetails, Name := initNetworkCredentials()
	//userDetails := initializeUsers()
	numOfPeersOnNetwork := len(peerDetails)
	numOfUsersOnNetwork := len(userDetails)
	fmt.Println("Num of Peers", numOfPeersOnNetwork)
	fmt.Println("Num of Users", numOfUsersOnNetwork)
	fmt.Println("Name of network", Name)

	allPeers := make([]Peer, numOfPeersOnNetwork)

	//factor := numOfUsersOnNetwork / numOfPeersOnNetwork
	remainder := numOfUsersOnNetwork % numOfPeersOnNetwork
	k := 0
	//for each peerDetail we construct a new peer evenly distributing the list of users
	for i :=0; i < numOfPeersOnNetwork; i++{

		aPeerDetail := make(map[string]string)
		aPeerDetail["ip"] = peerDetails[i].IP
		aPeerDetail["port"] = peerDetails[i].PORT
		aPeerDetail["name"] = peerDetails[i].NAME

		userInfo := make(map[string]string)
		userInfo[userDetails[i].USER] = userDetails[i].SECRET

		aPeer := new(Peer)
		aPeer.PeerDetails = aPeerDetail
		aPeer.UserData = userInfo
		aPeer.State = RUNNING
		allPeers[i] = *aPeer
	}
	//do we have any left over users details
	if remainder > 0 {
		for m := 0; m < remainder; m++ {
			allPeers[m].UserData[userDetails[k].USER] = userDetails[k].SECRET
			k++;
		}
	}

	return allPeers, Name
}

func initNetworkCredentials() ([]peerHTTP, []userData, string) {

	pwd, _ := os.Getwd()
        fmt.Println("CWD Inside initNetworkCredentials:", pwd)
	file, err := os.Open(pwd+"/../util/NetworkCredentials.json")

	if err != nil {
		fmt.Println("Error in opening NetworkCredentials file ", err)
		log.Fatal("Error in opening Network Credential json File")

	}
	networkCredentials, err := unmarshalNetworkCredentials(file)
	if err != nil {
		log.Fatal("Error in unmarshalling")
	}
	//peerdata := make(map[string]string)
	peerData := networkCredentials.PEERHTTP
	userData := networkCredentials.USERDATA
	name := networkCredentials.NAME
	return peerData, userData, name
}
