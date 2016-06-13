package main

import (
	"fmt"
	"obcsdk/chaincode"
	"obcsdk/peernetwork"
	"strconv"
	"sync"
	"time"
)

var peerNetworkSetup peernetwork.PeerNetwork
var AVal, BVal, curAVal, curBVal, invokeValue int64
var argA = []string{"a"}
var argB = []string{"b"}

var data string
var counter int64

var url string
const(
	THREAD_COUNT = 4
	TOTAL_NODES = 4
	TRX_COUNT = 20000
)
func setupNetwork() {
	fmt.Println("Creating a local docker network")
	peernetwork.SetupLocalNetwork(TOTAL_NODES, true)
	//peernetwork.PrintNetworkDetails()
	peerNetworkSetup = chaincode.InitNetwork()
	chaincode.InitChainCodes()
	chaincode.RegisterUsers()
}

//TODO : rather can we have a map for sleep for millis, secs and mins
func sleep(secs int64) {
	time.Sleep(time.Second * time.Duration(secs))
}

func deployChaincode(done chan bool) {
	example := "mycc"
	var funcArgs = []string{example, "init"}
	var args = []string{argA[0], data, argB[0], "0"}

	chaincode.Deploy(funcArgs, args)

	sleep(40)
	done <- true
}

func invokeChaincodeOnPeer(peerName string) {
	counter++
	fmt.Println("Iteration# ["+ strconv.FormatInt(counter, 10)+"] On "+peerName)

	arg1Construct := []string{"mycc", "invoke", peerName}
	arg2Construct := []string{"a" + strconv.FormatInt(counter, 10), data, "b"}

	_,_ = chaincode.InvokeOnPeer(arg1Construct, arg2Construct) //invRes
}

func queryChaincode() (res1, res2 string) { //int64) {
	var qargA = []string{"a" + strconv.FormatInt(counter, 10)}
	qAPIArgs0 := []string{"mycc", "query"}
	A, _ := chaincode.Query(qAPIArgs0, qargA)
	B, _ := chaincode.Query(qAPIArgs0, []string{"b"})
	return A, B
}

//TODO: Be cautious, race conditions might occur
var wg sync.WaitGroup

func main() {
	//done chan int
	done := make(chan bool, 1)
	wg.Add(THREAD_COUNT)
	counter = 0
	data = "!function(a){\"function\"==typeof define&&define.amd?define([\"jquery\"],a):\"object\"==typeof exports?module.exports=a:a(jQuery)}(function(a){function b(b){var g=b||window.event,h=i.call(arguments,1),j=0,l=0,m=0,n=0,o=0,p=0;if(b=a.event.fix(g),b.type=\"mousewheel\",\"detail\"in g&&(m=-1*g.detail),\"wheelDelta\"in g&&(m=g.wheelDelta),\"wheelDeltaY\"in g&&(m=g.wheelDeltaY),\"wheelDeltaX\"in g&&(l=-1*g.wheelDeltaX),\"axis\"in g&&g.axis===g.HORIZONTAL_AXIS&&(l=-1*m,m=0),j=0===m?l:m,\"deltaY\"in g&&(m=-1*g.deltaY,j=m),\"deltaX\"in g&&(l=g.deltaX,0===m&&(j=-1*l)),0!==m||0!==l){if(1===g.deltaMode){var q=a.data(this,\"mousewheel-line-height\");j*=q,m*=q,l*=q}else if(2===g.deltaMode){var r=a.data(this,\"mousewheel-page-height\");j*=r,m*=r,l*=r}if(n=Math.max(Math.abs(m),Math.abs(l)),(!f||f>n)&&(f=n,d(g,n)&&(f/=40)),d(g,n)&&(j/=40,l/=40,m/=40),j=Math[j>=1?\"floor\":\"ceil\"](j/f),l=Math[l>=1?\"floor\":\"ceil\"](l/f),m=Math[m>=1?\"floor\":\"ceil\"](m/f),k.settings.normalizeOffset&&this.getBoundingClientRect){var s=this.getBoundingClientRect();o=b.123456789"
	// Setup the network based on the NetworkCredentials.json provided
	setupNetwork()
	//Deploy chaincode
	deployChaincode(done)

	// time to messure overall execution of the testcase
	defer timeTrack(time.Now(), "LedgerStressFourClientTwoPeer execution")
	//TODO; Can we directly call goroutines here ???
	InvokeThreads()
	wg.Wait()
}

func InvokeMultiThreads() {
	for i := 1; i <= THREAD_COUNT; i++ {
		go func(val int) {
			for j := 1; j <= TRX_COUNT/THREAD_COUNT; j++ {
				fmt.Printf("\n============== CLIENT%d ==============\n", (THREAD_COUNT % val))
				invokeChaincodeOnPeer("PEER"+strconv((THREAD_COUNT % val)+1))
			}
			wg.Done()
		}(i)
	}
}
//Invokes loop
func InvokeThreads() {

	go func() {
		for i := 1; i <= 5000; i++ {
			fmt.Println("============== CLIENT1 ==============")
			invokeChaincodeOnPeer("PEER0")
		}
		wg.Done()
	}()

	go func() {
		for i := 1; i <= 5000; i++ {
			fmt.Println("============== CLIENT2 ==============")
			invokeChaincodeOnPeer("PEER1")
		}
		wg.Done()
	}()
	go func() {
		for i := 1; i <= 5000; i++ {
			fmt.Println("============== CLIENT3 ==============")
			invokeChaincodeOnPeer("PEER3")
		}
		wg.Done()
	}()

	go func() {
		for i := 1; i <= 5000; i++ {
			fmt.Println("============== CLIENT4 ==============")
			invokeChaincodeOnPeer("PEER4")
		}
		wg.Done()
	}()
}
func displayChainHeight(){
	startValue := 3
	height := 0
	var urlStr string
	for i:=0;i<TOTAL_NODES;i++ {
		urlStr = "http://172.17.0."+strconv.Itoa(startValue+i)+":5000"
		height = chaincode.Monitor_ChainHeight(urlStr)
		fmt.Println("################ Chaincode Height on "+urlStr+" is : ", height)
	}
}
func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	// Should we mask this delay ?
	//sleep(10)

	val1, val2 := queryChaincode()
	var exitCounter = 0;
	for val2 != strconv.FormatInt(counter, 10) && exitCounter < 3 {
		fmt.Printf("\n########### Peers are not in sync ? Check again after 5 sec")
		sleep(10)
		_, val2 = queryChaincode()
		exitCounter ++;
	}
  displayChainHeight()

	fmt.Printf("\n########### After Query Vals\n A = %s \nCounter = %s", val1, val2)
	fmt.Printf("\n\n################# %s took %s \n\n", name, elapsed)
}
