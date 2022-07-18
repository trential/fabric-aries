# /bin/sh
. ./scripts/env.sh

cryptogen generate --output organizations --config crypto-config.yaml

# Generate system channel
configtxgen -outputBlock artifacts/genesis.block -channelID syschannel -profile TwoOrgsOrdererGenesis

# Genreate application channel
./bin/configtxgen -outputCreateChannelTx artifacts/channel.tx -profile TwoOrgsChannel -channelID mychannel

# generate anchor update tx
./bin/configtxgen -outputAnchorPeersUpdate artifacts/org1-anchor.tx -profile TwoOrgsChannel -channelID mychannel -asOrg Org1MSP
./bin/configtxgen -outputAnchorPeersUpdate artifacts/org2-anchor.tx -profile TwoOrgsChannel -channelID mychannel -asOrg Org2MSP