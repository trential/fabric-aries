# /bin/bash

source ./scripts/env.sh 1

# make create application tx
peer channel create -f artifacts/channel.tx -c $CHANNEL_NAME -o $ORDERER_ENDPOINT --tls --cafile $ORDERER_CA

# org1 join channel
peer channel join -b mychannel.block

# update anchor peer
peer channel update -f artifacts/org1-anchor.tx -c $CHANNEL_NAME -o $ORDERER_ENDPOINT --tls --cafile $ORDERER_CA

source ./scripts/env.sh 2

# org2 join channel
peer channel join -b mychannel.block

# update anchor peer
peer channel update -f artifacts/org2-anchor.tx -c $CHANNEL_NAME -o $ORDERER_ENDPOINT --tls --cafile $ORDERER_CA