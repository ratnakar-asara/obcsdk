package peerrest

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

/*
  Issue GET request to BlockChain resource
    url is the GET request.
	respStatus is the HTTP response status code and message
	respBody is the HHTP response body
*/
func GetChainInfo(url string) (respBody string, respStatus string) {
	//TODO : define a logger
	//fmt.Println("GetChainInfo ... :", url)
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("%s", err)
		return err.Error(), "Error from GET request"
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
			return err.Error(), "Error from GET request"
		}
		return string(contents), response.Status
	}
}

/*
  Issue POST request to BlockChain resource.
    url is the target resource.
	payLoad is the REST API payload
	respStatus is the HTTP response status code and message
	respBody is the HHTP response body
*/
func PostChainAPI(url string, payLoad []byte) (respBody string, respStatus string) {

	//fmt.Println(">>>>> From postchain >>> ", url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payLoad))
	//req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	var errCount int
	if err != nil {
		log.Println("Error", url, err)
		errCount++
		return err.Error(), "There was an error Posting http Request"
	}
	errCount = 0
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error")
	}
	//fmt.Println("From postchain >>> response Status:", resp.Status)
	//fmt.Println("From postchain >>> response Body:", body)
	return string(body), resp.Status
}
