package peernetwork

import (
	"encoding/json"
	"fmt"
	"io"
)

type userData struct {
	USER   string `json:"username"`
	SECRET string `json:"secret"`
}

type peerHTTP struct {
	NAME string `json:"name"`
	IP   string `json:"api-host"`
	PORT string `json:"api-port"`
}

type peerGRPC struct {
	IP   string `json:"api-host"`
	PORT string `json:"api-port"`
}

/*************
type networkCredentials struct {
	PEERLIST []peerHTTP `json:"PeerData"`
	USERLIST []userData `json:"UserData"`
	PEERGRPC []peerGRPC `json:"PeerGrpc"`
	NAME     string     `json:"Name"`
}
*********************/
type networkCredentials struct {
	PEERHTTP []peerHTTP `json:"PeerData"`
	USERDATA []userData `json:"UserData"`
	PEERGRPC []peerGRPC `json:"PeerGrpc"`
	NAME     string     `json:"Name"`
}

type chainCodeData struct {
	NAME string `json:"name"`
	TYPE string `json:"type"`
	PATH string `json:"path"`
	//DEP_TXID       string `json:"dep_txid"`
	//DEPLOYED       string `json:"deployed"`
}

//func main() {

//file, err := os.OpenFile("C:/Go/src/obcpeer-test/util/Userdata.txt", os.O_RDONLY)
// For read access.
//file, err := os.Open("C:/Go/src/obcpeer-test/util/PNStruct.json")
//this will not get a reader
//file, err := ioutil.ReadFile("C:/Go/src/obcpeer-test/util/Userdata.txt")
//if err != nil {
//	log.Fatal(err)
//}
//fmt.Println("Unmarshalling user data for ")
//NC, err := UnmarshalNCData(file)
//if err != nil {
//	fmt.Println("Error in unmarshalling")
//}
//	i := 0

//pl := NC.PEERLIST
//	 for  i < len (NC.PEERLIST){
//		 fmt.Println(NC.PEERLIST[i])
//		 i++
//	 }
//	 i = 0
//	 //ul := NC.USERLIST
//	 for  i < len (NC.USERLIST){
//		 fmt.Println(NC.USERLIST[i])
//		 i++
//	 }
//	 fmt.Println(NC.NAME)
//	}

/*
  converts input stream to NetworkCredentials
	reader is an open file
*/
func unmarshalNetworkCredentials(reader io.Reader) (networkCredentials, error) {

	decoder := json.NewDecoder(reader)
	//fmt.Println("Inside Unmarshal Network Credentials JSONREADER")
	var NC networkCredentials
	err := decoder.Decode(&NC)
	if err != nil {
		fmt.Println("Decoding error \t Error:", err)
	}
	return NC, err
}

/*
  converts input stream to chaincodes.
*/
func unmarshalChainCodes(reader io.Reader) ([]*chainCodeData, error) {

	decoder := json.NewDecoder(reader)
	//fmt.Println("Inside Unmarshal ChainCode JSONREADER")
	var ChainCodeCollection []*chainCodeData
	err := decoder.Decode(&ChainCodeCollection)
	if err != nil {
		fmt.Println("Error in decoding ")
	}
	return ChainCodeCollection, err
}
