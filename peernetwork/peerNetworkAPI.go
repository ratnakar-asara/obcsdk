package peernetwork

import (
	"fmt"
	"errors"
	"log"
	"strings"
	"time"
	"os/exec"
)

/*
  prints the content of the network: peers, users, and chaincodes.
*/
func PrintNetworkDetails() {

	ThisNetwork := LoadNetwork()
	Peers := ThisNetwork.Peers
	i := 0
	for i < len(Peers) {

		msgStr := fmt.Sprintf("ip: %s port: %s name %s ", Peers[i].PeerDetails["ip"], Peers[i].PeerDetails["port"], Peers[i].PeerDetails["name"])
		fmt.Println(msgStr)
		userList := ThisNetwork.Peers[i].UserData
		fmt.Println("Users:")
		for user, secret := range userList {

			fmt.Println(user, secret)
		}
		i++
	}
	fmt.Println("Available Chaincodes :")
	libChainCodes := InitializeChainCodes()
	for k, v := range libChainCodes.ChainCodes {
		fmt.Println("\nChaincode :", k)
		fmt.Println("\nDetail :\n")
		for i, j := range v.Detail {
			msgStr := fmt.Sprintf("user: %s secret: %s", i, j)
			fmt.Println(msgStr)
		}
		fmt.Println("\n")
	}
}

/*
 Gets ChainCode detail for a given chaincode name
  Takes two arguments
	1. name (string)			- name of the chaincode as specified in CC_Collections.json file
	2. lcc (LibChainCodes)		- LibChainCodes struct having current collection of all chaincodes loaded in the network.
  Returns:
 	1. ccDetail map[string]string  	- chaincode details of the chaincode requested as a map of key/value pairs.
	2. Versions map[string]string   - versioning or tagging details on the chaincode requested as a map of key/value pairs
*/
func GetCCDetailByName(name string, lcc LibChainCodes) (ccDetail map[string]string, versions map[string]string, err error) {
	var errStr string
	var err1 error
	for k, v := range lcc.ChainCodes {
		if strings.Contains(k, name) {
			return v.Detail, v.Versions, err1
		}
	}
	//no more chaincodes construct error string and empty maps
	errStr = fmt.Sprintf("chaincode %s does not exist on the network", name)
	//need to check for this
	j := make(map[string]string)
	return j, j, errors.New(errStr)
}

/** utility functions to aid users in getting to a valid URL on network
 ** to post chaincode rest API
 **/

/*
  gets any one running peer from 'thisNetwork' as set by chaincode.Init()
*/
func APeer(thisNetwork PeerNetwork) (thisPeer *Peer, err error) {
	//thisNetwork := LoadNetwork()
	Peers := thisNetwork.Peers
	var aPeer *Peer
	var errStr string
	//get any peer that has at a minimum one userData and one peerDetails
	for peerIter := range Peers {
		if (len(Peers[peerIter].UserData) > 0) && (len(Peers[peerIter].PeerDetails) > 0) {
			if Peers[peerIter].State == 0 ||  Peers[peerIter].State == 2 || Peers[peerIter].State == 4 {
				aPeer = &Peers[peerIter]
			}
		}
	}
	if aPeer != nil {
		return (aPeer), errors.New("")
	} else {
		errStr = fmt.Sprintf("Not found valid running peer on network")
		return aPeer, errors.New(errStr)
	}

}

/*
  gets IP address of a Peer given it's name on the entire network.
*/
func IPPeer(thisNetwork PeerNetwork, peername string) (IP string, err error) {

	//fmt.Println("Values inside AUserFromNetwork ", ip, port, user)
	Peers := thisNetwork.Peers
	var aPeer *Peer
	var errStr string
	peerFullName, _ := GetFullPeerName(thisNetwork, peername)
	//get any peer that has at a minimum one userData and one peerDetails
	for peerIter := range Peers {
		if (len(Peers[peerIter].UserData) > 0) && (len(Peers[peerIter].PeerDetails) > 0) {
			if Peers[peerIter].PeerDetails["name"] == peerFullName {
				aPeer = &Peers[peerIter]
			}
		}
	}
	if aPeer != nil {
		return (aPeer.PeerDetails["IP"]), errors.New("")
	} else {
		errStr = fmt.Sprintf("Not found %s peer on network", peername)
		return aPeer.PeerDetails["IP"], errors.New(errStr)
	}
}

/*
  gets any one user from any Peer on the entire network.
*/
func AUserFromNetwork(thisNetwork PeerNetwork) (thisPeer *Peer, user string) {

	//fmt.Println("Values inside AUserFromNetwork ", ip, port, user)
	var u string
	aPeer, _ := APeer(thisNetwork)
	users := aPeer.UserData

	//fmt.Println(" ip ", aPeer.PeerDetails["ip"])
	for u, _ = range users {
		break
	}
	//fmt.Println(" ip ", aPeer.UserData["ip"])
	//fmt.Println(" ip ", user)
	//return aPeer.PeerDetails["ip"], aPeer.PeerDetails["port"], u
	return aPeer, u
}

/*
  finds any one user associated with the given peer
*/
func AUserFromAPeer(thisPeer Peer) (ip string, port string, user string) {

	//var aPeer *Peer
	aPeer := thisPeer
	var curUser string
	userList := aPeer.UserData
	for curUser, _ = range userList {
		break
	}
	//fmt.Println(" ip ", aPeer.UserData["ip"])
	//fmt.Println(" ip ", user)
	return aPeer.PeerDetails["ip"], aPeer.PeerDetails["port"], curUser
}

/*
 gets a user from a Peer with the given IP on the PeerNetwork
*/
func AUserFromThisPeer(thisNetwork PeerNetwork, host string) (ip string, port string, user string, err error) {

	//var aPeer *Peer
	Peers := thisNetwork.Peers
	var aPeer *Peer
	var u string
	var errStr string
	var err1 error

	//get a random peer that has at a minimum one userData and one peerDetails
	for peerIter := range Peers {
		if Peers[peerIter].State == 0 {
				if (strings.Contains(host, ":")) {
					if strings.Contains(Peers[peerIter].PeerDetails["ip"], host) {
						aPeer = &Peers[peerIter]
					}
				}else { //host: "vp1"
					if strings.Contains(Peers[peerIter].PeerDetails["name"], host) {
						//fmt.Println("Inside name IP resolution")
						aPeer = &Peers[peerIter]
					}
				}
		}
	}

	//fmt.Println(" * aPeer ", *aPeer)
	if aPeer != nil {
		users := aPeer.UserData
		for u, _ = range users {
			break
		}
		return aPeer.PeerDetails["ip"], aPeer.PeerDetails["port"], u, err1
	}else {
			errStr = fmt.Sprintf("%s, Not found on network", host)
			return "", "", "", errors.New(errStr)
	}
}


/*
  finds the peer address corresponding to a given user
    thisNetwork as set by chaincode.init
	ip, port are the address of the peer.
	user is the user details: name, credential.
	err	is an error message, or nil if no error occurred.
*/
func PeerOfThisUser(thisNetwork PeerNetwork, username string) (ip string, port string, user string, err error) {

	//var aPeer *Peer
	Peers := thisNetwork.Peers
	var aPeer *Peer
	var errStr string
	var err1 error
	//fmt.Println("Inside function")
	//get a random peer that has at a minimum one userData and one peerDetails
	for peerIter := range Peers {
		if (len(Peers[peerIter].UserData) > 0) && (len(Peers[peerIter].PeerDetails) > 0) {
			if Peers[peerIter].State == 0 {
				if _, ok := Peers[peerIter].UserData[username]; ok {
					fmt.Println("Found %s in network", username)
					aPeer = &Peers[peerIter]
				}
			}
		}
	}
	if aPeer == nil {
		errStr = fmt.Sprintf("%s, Not found on network", username)
		return "", "", "", errors.New(errStr)
	} else {
		return aPeer.PeerDetails["ip"], aPeer.PeerDetails["port"], username, err1
	}
}

/*Gets the peer details corresponding to a given peer-name
state if running/stopped/suspended:0/1/2
err	is an error message, or nil if no error occurred.
*/
func GetPeerState(thisNetwork PeerNetwork, peername string) (currPeer *Peer, err error) {

	//var aPeer *Peer
	Peers := thisNetwork.Peers
	var aPeer *Peer
	var errStr string
	fullName, _ := GetFullPeerName(thisNetwork, peername)
	for peerIter := range Peers {
		if (len(Peers[peerIter].UserData) > 0) && (len(Peers[peerIter].PeerDetails) > 0) {
			if strings.Contains(Peers[peerIter].PeerDetails["name"], fullName) {
				aPeer = &Peers[peerIter]
			}
		}
	}

	if aPeer == nil {
		errStr = fmt.Sprintf("%s, Not found on network", peername)
		emptyPD := new(Peer)
		return emptyPD, errors.New(errStr)
	} else {
		return aPeer, errors.New("")
	}
}

/*
  sets the peer details corresponding to a given peer-name
  state if running/stopped/suspended:0/1/2
	err	is an error message, or nil if no error occurred.
*/
func SetPeerState(thisNetwork PeerNetwork, peername string, curstate int) (peerDetails map[string]string, err error) {

	//var aPeer *Peer
	Peers := thisNetwork.Peers
	var aPeer *Peer
	var errStr string
	//fmt.Println("Inside function")
	//get a peerDetails from peername
	fullName, _ := GetFullPeerName(thisNetwork, peername)
	for peerIter := range Peers {
		if (len(Peers[peerIter].UserData) > 0) && (len(Peers[peerIter].PeerDetails) > 0) {
			if strings.Contains(Peers[peerIter].PeerDetails["name"], fullName) {
				aPeer = &Peers[peerIter]
			}
		}
	}

	if aPeer == nil {
		errStr = fmt.Sprintf("%s, Not found on network", peername)
		emptyPD := make(map[string]string)
		return emptyPD, errors.New(errStr)
	} else {
		aPeer.State = curstate
		fmt.Println("curstate", curstate)
		return aPeer.PeerDetails, errors.New("")
	}
}


func PausePeersLocal(thisNetwork PeerNetwork, peers []string) {

	i := 0
	for i < len(peers) {
		cmd := "docker pause " + peers[i]
		out, err := exec.Command("/bin/sh", "-c", cmd).Output()
                if (err != nil) {
                     fmt.Println("Could not Pause %s", peers[i])
                     log.Fatal(err)
                }
                fmt.Printf("peer %s", out)
		time.Sleep(5000 * time.Millisecond)
		SetPeerState(thisNetwork, peers[i], 3)
		i++
	}
}

func PausePeerLocal(thisNetwork PeerNetwork, peer string) {

	cmd := "docker pause " + peer
        out, err := exec.Command("/bin/sh", "-c", cmd).Output()
        if (err != nil) {
           fmt.Println("Could not Pause %s", peer)
           log.Fatal(err)
        }
        fmt.Printf("Paused peer %s", out)
	time.Sleep(5000 * time.Millisecond)
	SetPeerState(thisNetwork, peer, 3)

}

func UnpausePeersLocal(thisNetwork PeerNetwork, peers []string) {

	i := 0
	for i < len(peers) {
		cmd := "docker unpause " + peers[i]
                out, err := exec.Command("/bin/sh", "-c", cmd).Output()
                if (err != nil) {
                   fmt.Println("Could not Unpause %s ", peers[i])
                   log.Fatal(err)
                }
                fmt.Printf("Unpaused peer %s", out)
		exec.Command(cmd)
		time.Sleep(5000 * time.Millisecond)
		SetPeerState(thisNetwork, peers[i], 4)
		i++
	}
}


func UnpausePeerLocal(thisNetwork PeerNetwork, peer string) {

	cmd := "docker unpause " + peer
        out, err := exec.Command("/bin/sh", "-c", cmd).Output()
        if (err != nil) {
           fmt.Println("Could not Unpause %s ", peer)
           log.Fatal(err)
        }
        fmt.Printf("Paused peer %s", out)
	time.Sleep(5000 * time.Millisecond)
	SetPeerState(thisNetwork, peer, 4)

}


func StopPeersLocal(thisNetwork PeerNetwork, peers []string) {

	i := 0
	for i < len(peers) {
		cmd := "docker stop " + peers[i]
                out, err := exec.Command("/bin/sh", "-c", cmd).Output()
                if (err != nil) {
                   fmt.Println("Could not Stop %s successfully", peers[i])
                   log.Fatal(err)
                }
                fmt.Printf("Stopped peer %s successfully", out)
		exec.Command(cmd)
		time.Sleep(5000 * time.Millisecond)
		SetPeerState(thisNetwork, peers[i], 1)
		i++
	}
}

func StartPeersLocal(thisNetwork PeerNetwork, peers []string) {

	i := 0
	for i < len(peers) {
		cmd := "docker start " + peers[i]
                out, err := exec.Command("/bin/sh", "-c", cmd).Output()
                if (err != nil) {
                   fmt.Println("Could not Start %s successfully", peers[i])
                   log.Fatal(err)
                }
                fmt.Printf("Started peer %s successfully", out)
		exec.Command(cmd)
		time.Sleep(5000 * time.Millisecond)
		SetPeerState(thisNetwork, peers[i], 2)
		i++
	}
}
func StartPeerLocal(thisNetwork PeerNetwork, peer string) {

	cmd := "docker start " + peer
  out, err := exec.Command("/bin/sh", "-c", cmd).Output()
  if (err != nil) {
           fmt.Println("Could not Start %s successfully", peer)
           log.Fatal(err)
  }
  fmt.Printf("Started peer %s successfully", out)
	time.Sleep(5000 * time.Millisecond)
	SetPeerState(thisNetwork, peer, 2)
}

func StopPeerLocal(thisNetwork PeerNetwork, peer string) {

	cmd := "docker stop " + peer
        out, err := exec.Command("/bin/sh", "-c", cmd).Output()
        if (err != nil) {
           fmt.Println("Could not Stop %s ", peer)
           log.Fatal(err)
        }
        fmt.Printf("Stopped peer %s successfully", out)
	time.Sleep(5000 * time.Millisecond)
	SetPeerState(thisNetwork, peer, 1)
}

func GetFullPeerName(thisNetwork PeerNetwork, shortname string) (name string, err error) {
	Peers := thisNetwork.Peers
	var aPeer *Peer
	var errStr string
	//fmt.Println("Inside function")
	//get a peerDetails from peername
	for peerIter := range Peers {
		if (len(Peers[peerIter].UserData) > 0) && (len(Peers[peerIter].PeerDetails) > 0) {
			if strings.Contains(Peers[peerIter].PeerDetails["name"], shortname) {
				aPeer = &Peers[peerIter]
			}
		}
	}

	if aPeer == nil {
		errStr = fmt.Sprintf("%s, Not found on network", shortname)
		return "", errors.New(errStr)
	} else {
		return aPeer.PeerDetails["name"], errors.New("")

	}
}

func AddAPeerNetwork() {

}

/********************
type PeerNetworks struct {
	PNetworks      []PeerNetwork
}


func AddAPeerNetwork() {

}

func DeleteAPeerNetwork() {

}

func AddUserOnAPeer(){

}

func RemoveUserOnAPeer(){

}


func LoadNetworkByName(name string) PeerNetwork {

  networks := LoadPeerNetworks()
	pnetworks := networks.PNetworks
	for peerIter := range pnetworks {
		//fmt.Println(pnetworks[peerIter].Name)
		if strings.Contains(pnetworks[peerIter].Name, name) {
			return pnetworks[peerIter]
		}
	}
	//return *new(PeerNetwork)
}
*********************/
