# /bin/bash
source ./scripts/env.sh 1

CC_VERSION=${1-1}
LOG_PATH=/tmp/"${CC_NAME}_log"

# package chaincode
peer lifecycle chaincode package "${CC_NAME}.tar.gz" -p $CC_PATH --label $CC_NAME --lang $CC_LANG

# install chaincode on peer0.org1.example.com
peer lifecycle chaincode install "${CC_NAME}.tar.gz" >&$LOG_PATH

PACKAGE_ID=`cat ${LOG_PATH} | grep "Chaincode code package identifier:" | awk '{split($0,a,"Chaincode code package identifier: "); print a[2]}'`

# approveformyorg
peer lifecycle chaincode approveformyorg -o $ORDERER_ENDPOINT \
        --tls --cafile $ORDERER_CA \
        -C $CHANNEL_NAME \
        -n $CC_NAME \
        --package-id $PACKAGE_ID \
        --sequence $CC_VERSION \
        -v $CC_VERSION \
        --init-required

source ./scripts/env.sh 2

# install chaincode on peer0.org2.example.com
peer lifecycle chaincode install "${CC_NAME}.tar.gz" >&$LOG_PATH

PACKAGE_ID=`cat ${LOG_PATH} | grep "Chaincode code package identifier:" | awk '{split($0,a,"Chaincode code package identifier: "); print a[2]}'`

# approveformyorg
peer lifecycle chaincode approveformyorg -o $ORDERER_ENDPOINT \
        --tls --cafile $ORDERER_CA \
        -C $CHANNEL_NAME \
        -n $CC_NAME \
        --package-id $PACKAGE_ID \
        --sequence $CC_VERSION \
        -v $CC_VERSION \
        --init-required

# query before commiting chaincode
peer lifecycle chaincode checkcommitreadiness -o $ORDERER_ENDPOINT \
        --tls --cafile $ORDERER_CA \
        -C $CHANNEL_NAME \
        -n $CC_NAME \
        --sequence $CC_VERSION\
       -v $CC_VERSION \
       --init-required

# commit chaincode
peer lifecycle chaincode commit -o $ORDERER_ENDPOINT \
        --tls --cafile $ORDERER_CA \
        -C $CHANNEL_NAME \
        -n $CC_NAME \
        --sequence $CC_VERSION \
        -v $CC_VERSION \
        --waitForEvent \
        --peerAddresses $CORE_PEER_ADDRESS \
        --tlsRootCertFiles $CORE_PEER_TLS_ROOTCERT_FILE \
        --peerAddresses "localhost:7051" \
        --tlsRootCertFiles $ORG1_CA \
        --init-required

source ./scripts/env.sh 2

# genesis trustees tx
nym_template='{"alias" : "%s", "dest" : "%s", "verkey": "%s"}'

TX_FILE="genesis_trustees.csv"
REQUEST="["
while IFS="," read -r val_alias val_did val_verkey
do
   nym=$(printf "$nym_template" "$val_alias" "$val_did" "$val_verkey")
   REQUEST="${REQUEST}${nym}," 
done < <(tail -n +2 $TX_FILE)
REQUEST=${REQUEST::-1}
REQUEST="${REQUEST}]"

REQUEST=$(echo -n $REQUEST | jq -Rsa .)
c=$(printf '{"args" : ["_",%s]}' "$REQUEST") 

peer chaincode invoke -n $CC_NAME -C $CHANNEL_NAME -c "$c" \
                            -o $ORDERER_ENDPOINT \
                            --tls --cafile $ORDERER_CA \
                            --peerAddresses $CORE_PEER_ADDRESS \
                            --tlsRootCertFiles $CORE_PEER_TLS_ROOTCERT_FILE \
                            --peerAddresses "localhost:7051" \
                            --tlsRootCertFiles $ORG1_CA \
                            --waitForEvent \
                            --isInit