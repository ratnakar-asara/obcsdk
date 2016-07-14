# obcsdk
A test framework for testing Blockchain with GoSDK (written in Go Language)

**_PAGE CONSTRUCTION IS UNDER PROGRESS_**

**How to test the programs:**
* Clone to the src folder in where go is installed
```
  $ cd $GOROOT/src
  $ git clone https://github.com/ratnakar-asara/obcsdk.git
```
* Follow instructions mentioned [here](https://github.com/hyperledger/fabric/blob/master/docs/dev-setup/devnet-setup.md) to setup peer network
* Change the credentials in [NetworkCredentials.json](https://github.com/ratnakar-asara/obcsdk/blob/master/util/NetworkCredentials.json) accordingly
* Go to the chcotest folder and execute the tests
```
  $ cd obcsdk/chcotest
  $ NETWORK="LOCAL" go run LedgerStressOneCliOnePeer.go Utils.go
```



