package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"obcsdk/chaincode"
)

// A Utility program, contains several utility methods that can be used across programs
const(
	CHAINCODE_NAME = "mycc"
	INIT= "init"
	INVOKE = "invoke"
	QUERY = "query"
	DATA = "Yh1WWZlw1gGd2qyMNaHqBCt4zuBrnT4cvZ5iMXRRM3YBMXLZmmvyVr0ybWfiX4N3UMliEVA0d1dfTxvKs0EnHAKQe4zcoGVLzMHd8jPQlR5ww3wHeSUGOutios16lxfuQTdnsFcxhXLiGwp83ahyBomdmJ3igAYTyYw2bwXqhBeL9fa6CTK43M2QjgFhQtlcpsh7XMcUWnjJhvMHAyH67Z8Ugke6U8GQMO5aF1Oph0B2HlIQUaHMq2i6wKN8ZXyx7CCPr7lKnIVWk4zn0MLZ16LstNErrmsGeo188Rdx5Yyw04TE2OSPSsaQSDO6KrDlHYnT2DahsrY3rt3WLfBZBrUGhr9orpigPxhKq1zzXdhwKEzZ0mi6tdPqSzMKna7O9STstf2aFdrnsoovOm8SwDoOiyqfT5fc0ifVZSytVNeKE1C1eHn8FztytU2itAl1yDYSfTZQv42tnVgDjWcLe2JR1FpfexVlcB8RUhSiyoThSIFHDBZg8xyULPmp4e6acOfKfW2BXh1IDtGR87nBWqmytTOZrPoXRPq2QXiUjZS2HflHJzB0giDbWEeoZoMeF11364Xzmo0iWsBw0TQ2cHapS4cR49IoEDWkC6AJgRaNb79s6vythxX9CqfMKxIpqYAbm3UAZRS7QU7MiZu2qG3xBIEegpTrkVNneprtlgh3uTSVZ2n2JTWgexMcpPsk0ILh10157SooK2P8F5RcOVrjfFoTGF3QJTC2jhuobG3PIXs5yBHdELe5yXSEUqUm2ioOGznORmVBkkaY4lP025SG1GNPnydEV9GdnMCPbrgg91UebkiZsBMM21TZFbUqP70FDAzMWZKHDkDKCPoO7b8EPXrz3qkyaIWBymSlLt6FNPcT3NkkTfg7wl4DZYDvXA2EYu0riJvaWon12KWt9aOoXig7Jh4wiaE1BgB3j5gsqKmUZTuU9op5IXSk92EIqB2zSM9XRp9W2I0yLX1KWGVkkv2OIsdTlDKIWQS9q1W8OFKuFKxbAEaQwhc7Q5Mm"
)
// Called in teardown methods to messure and display over all execution time
func TimeTracker(start time.Time, info string) {
	elapsed := time.Since(start)
	fmt.Println("=========",info, " is ",elapsed)
}

func getChainHeight(url string) int {
	/*var urlStr string
	//TODO: peername - shouldn't hardcoded ??
	urlStr = "http://172.17.0.3:5000"*/
	height := chaincode.Monitor_ChainHeight(url)
	fmt.Println("=========  Chaincode Height on "+url+" is : ", height)
	return height
}

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
	fmt.Println("<<<<<< Deploy needs time, Let's sleep for 60 secs >>>>>>")
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

//TODO: This function doesn't work outside vagrant, need to relook
/*func displayChainHeight(nodes int){
	startValue := 3
	height := 0
	var urlStr string
	for i:=0;i<nodes;i++ {
		urlStr = "http://172.17.0."+strconv.Itoa(startValue+i)+":5000"
		height = chaincode.Monitor_ChainHeight(urlStr)
		fmt.Println("################ Chaincode Height on "+urlStr+" is : ", height)
	}
}*/

func sleep(secs int64) {
	time.Sleep(time.Second * time.Duration(secs))
}
