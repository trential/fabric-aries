# /bin/sh

export PATH=${PWD}/bin:$PATH

export CORE_PEER_TLS_ENABLED=true
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
export ORG1_CA=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export ORG2_CA=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
export ORDERER_ENDPOINT=localhost:7050
export CHANNEL_NAME=mychannel

export CC_NAME=ssi
export CC_PATH=../chaincode
export CC_LANG=golang


case $1 in
    "1")
        export FABRIC_CFG_PATH=./config
        export CORE_PEER_LOCALMSPID=Org1MSP
        export CORE_PEER_TLS_ROOTCERT_FILE=$ORG1_CA
        export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
        export CORE_PEER_ADDRESS=localhost:7051
    ;;
    "2")
        export FABRIC_CFG_PATH=./config
        export CORE_PEER_LOCALMSPID=Org2MSP
        export CORE_PEER_TLS_ROOTCERT_FILE=$ORG2_CA
        export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
        export CORE_PEER_ADDRESS=localhost:8051
    ;;
esac

