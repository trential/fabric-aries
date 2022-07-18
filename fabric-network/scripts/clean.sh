# /bin/bash

docker rm -f \
    orderer.example.com \
    ca_org1 \
    peer0.org1.example.com \
    couchdb1 \
    ca_org2 \
    peer0.org2.example.com \
    couchdb2

docker rm -f $(docker ps -a -f 'name=dev-peer0*' -q)
docker image rm $(docker image ls -f "reference=dev-peer0.org*" -q)

docker volume prune
docker network prune

rm -rf organizations artifacts