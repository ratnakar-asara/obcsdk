package peerrest

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"crypto/tls"
)

/*
  Issue GET request to BlockChain resource
    url is the GET request.
	respStatus is the HTTP response status code and message
	respBody is the HHTP response body
*/
func GetChainInfo(url string) (respBody string, respStatus string) {
	//TODO : define a logger

        tr := &http.Transport{
	         TLSClientConfig:    &tls.Config{RootCAs: nil},
	         DisableCompression: true,
        }
        client := &http.Client{Transport: tr}
        response, err := client.Get(url)
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

	//fmt.Println(">>>>> From secure postchain >>> ", url)
        tr := &http.Transport{
	         TLSClientConfig:    &tls.Config{RootCAs: nil},
	         DisableCompression: true,
        }
        client := &http.Client{Transport: tr}
	response, err := client.Post(url, "json", bytes.NewBuffer(payLoad))

	var errCount int
	if err != nil {
		log.Println("Error", url, err)
		errCount++
		return err.Error(), "There was an error Posting http Request"
	}
	errCount = 0
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error")
	}
	//fmt.Println("From secure postchain >>> response Status:", response.Status)
	//fmt.Println("From secure postchain >>> response Body:", body)
	return string(body), response.Status
}
