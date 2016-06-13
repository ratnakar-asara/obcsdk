rm -rf /var/hyperledger/*
cp -r /opt/gopath/src/github.com/hyperledger/fabric/obcsdk /opt/gopath/src/
docker rm -f $(docker ps -aq)
/opt/gopath/src/github.com/hyperledger/fabric/obcsdk/automation/local_fabric.sh -n 1 -s
