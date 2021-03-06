package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"obcsdk/chaincode"
)

// A Utility program, contains several utility methods that can be used across
// test programs
const (
	CHAINCODE_NAME = "mycc"
	INIT           = "init"
	INVOKE         = "invoke"
	QUERY          = "query"
	DATA           = "Yh1WWZlw1gGd2qyMNaHqBCt4zuBrnT4cvZ5iMXRRM3YBMXLZmmvyVr0ybWfiX4N3UMliEVA0d1dfTxvKs0EnHAKQe4zcoGVLzMHd8jPQlR5ww3wHeSUGOutios16lxfuQTdnsFcxhXLiGwp83ahyBomdmJ3igAYTyYw2bwXqhBeL9fa6CTK43M2QjgFhQtlcpsh7XMcUWnjJhvMHAyH67Z8Ugke6U8GQMO5aF1Oph0B2HlIQUaHMq2i6wKN8ZXyx7CCPr7lKnIVWk4zn0MLZ16LstNErrmsGeo188Rdx5Yyw04TE2OSPSsaQSDO6KrDlHYnT2DahsrY3rt3WLfBZBrUGhr9orpigPxhKq1zzXdhwKEzZ0mi6tdPqSzMKna7O9STstf2aFdrnsoovOm8SwDoOiyqfT5fc0ifVZSytVNeKE1C1eHn8FztytU2itAl1yDYSfTZQv42tnVgDjWcLe2JR1FpfexVlcB8RUhSiyoThSIFHDBZg8xyULPmp4e6acOfKfW2BXh1IDtGR87nBWqmytTOZrPoXRPq2QXiUjZS2HflHJzB0giDbWEeoZoMeF11364Xzmo0iWsBw0TQ2cHapS4cR49IoEDWkC6AJgRaNb79s6vythxX9CqfMKxIpqYAbm3UAZRS7QU7MiZu2qG3xBIEegpTrkVNneprtlgh3uTSVZ2n2JTWgexMcpPsk0ILh10157SooK2P8F5RcOVrjfFoTGF3QJTC2jhuobG3PIXs5yBHdELe5yXSEUqUm2ioOGznORmVBkkaY4lP025SG1GNPnydEV9GdnMCPbrgg91UebkiZsBMM21TZFbUqP70FDAzMWZKHDkDKCPoO7b8EPXrz3qkyaIWBymSlLt6FNPcT3NkkTfg7wl4DZYDvXA2EYu0riJvaWon12KWt9aOoXig7Jh4wiaE1BgB3j5gsqKmUZTuU9op5IXSk92EIqB2zSM9XRp9W2I0yLX1KWGVkkv2OIsdTlDKIWQS9q1W8OFKuFKxbAEaQwhc7Q5Mm"
)

var logEnabled bool
var logFile *os.File

// Called in teardown methods to messure and display over all execution time
func TimeTracker(start time.Time, info string) {
	elapsed := time.Since(start)
	logger(fmt.Sprintf("========= %s is %s", info, elapsed))
	closeLogger()
}

func getChainHeight(url string) int {
	height := chaincode.Monitor_ChainHeight(url)
	logger(fmt.Sprintf("=========  Chaincode Height on "+url+" is : %d", height))
	return height
}

// This is a helper function to generate a random string of the requested length
// This is to make each Deploy transaction unique
func RandomString(strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

// Utility function to deploy chaincode available @ http://urlmin.com/4r76d
func deployChaincode(done chan bool) {
	var funcArgs = []string{CHAINCODE_NAME, INIT}
	var args = []string{argA[0], RandomString(1024), argB[0], "0"}
	//call chaincode deploy function to do actual deployment
	chaincode.Deploy(funcArgs, args)
	logger("<<<<<< Deploy needs time, Let's sleep for 60 secs >>>>>>")
	sleep(60)
	done <- true
}

// Utility function to invoke on chaincode available @ http://urlmin.com/4r76d
/*func invokeChaincode(counter int64) {
	arg1 := []string{CHAINCODE_NAME, INVOKE}
	arg2 := []string{"a" + strconv.FormatInt(counter, 10), data, "counter"}
	_, _ = chaincode.Invoke(arg1, arg2)
}*/

// Utility function to query on chaincode available @ http://urlmin.com/4r76d
func queryChaincode(counter int64) (res1, res2 string) {
	var arg1 = []string{CHAINCODE_NAME, QUERY}
	var arg2 = []string{"a" + strconv.FormatInt(counter, 10)}

	val, _ := chaincode.Query(arg1, arg2)
	counterArg, _ := chaincode.Query(arg1, []string{"counter"})
	return val, counterArg
}

//TODO : These values should be configurable for different environments
var LocalUsers = []string{"test_user3", "test_user4", "test_user5", "test_user6", "test_user7"}
var ZUsers = []string{"dashboarduser_type0_efeeb83216", "dashboarduser_type0_fa08214e3b", "dashboarduser_type0_e00e125cf9", "dashboarduser_type0_e0ee60d5af"}

var LocalPeers = []string{"PEER0", "PEER1", "PEER2", "PEER3"}
var ZPeers = []string{"vp0", "vp1", "vp2", "vp3"}

//Get the user names based on network environment Z/Local
func getUser(userNumber int) string {
	if os.Getenv("NETWORK") == "Z" {
		return ZUsers[userNumber]
	} else {
		return LocalUsers[userNumber]
	}
}

//Get the peer name based on network environment Z/Local
func getPeer(peerNumber int) string {
	if os.Getenv("NETWORK") == "Z" {
		return ZPeers[peerNumber]
	} else {
		return LocalPeers[peerNumber]
	}
}

func sleep(secs int64) {
	time.Sleep(time.Second * time.Duration(secs))
}

func initLogger(fileName string) {
	layout := "Jan__2_2006"
	// Format Now with the layout const.
	t := time.Now()
	res := t.Format(layout)
	var err error
	logFile, err = os.OpenFile(res+"-"+fileName+".txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("error opening file: %v", err))
	}

	logEnabled = true
	log.SetOutput(logFile)
	//log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetFlags(log.LstdFlags)
}

func logger(printStmt string) {
	fmt.Println(printStmt)
	if !logEnabled {
		return
	}
	//TODO: Should we disable logging ?
	log.Println(printStmt)
}

func closeLogger() {
	if logEnabled && logFile != nil {
		logFile.Close()
	}
}

//Cleanup methods to display useful information
func tearDown(counter int64) {
	logger("....... State transfer is happening, Lets take a nap for 2 mins ......")
	// TODO: Change this value when invokes are in millions ?
	sleep(120)
	val1, val2 := queryChaincode(counter)
	logger(fmt.Sprintf("========= After Query values a%d = %s,  counter = %s\n", counter, val1, val2))

	newVal, err := strconv.ParseInt(val2, 10, 64)

	if err != nil {
		logger(fmt.Sprintf("Failed to convert %d to int64\n Error: %s\n", val2, err))
	}

	//TODO: Block size again depends on the Block configuration in pbft config file
	//Test passes when 2 * block height match with total transactions, else fails
	if newVal == counter {
		logger(fmt.Sprintf("######### Inserted %d records #########\n", counter))
		logger("######### TEST PASSED #########")
	} else {
		logger("######### TEST FAILED #########")
	}

}
