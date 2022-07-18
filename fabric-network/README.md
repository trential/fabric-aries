# Fabric Network

Hyperledger fabric network for aries-fabric integration demo.

- Prerequisite : 
  - Download fabric binaries : `curl -sSL https://bit.ly/2ysbOFE | bash -s -- 2.2.7 1.5.3 -d -s`
- Generate artifacts : `bash scripts/artifacts.sh`
- Create docker network : `docker network create aries_fabric`
- Start all the docker containers : `docker-compose up -d`
- Run channel script : `bash ./scripts/channel.sh`
- Installing chaincode :
  - Insert genesis trustee details in `genesis_trustees.csv`
  - chaincode name, path, language can be changed in `scripts/env.sh`
  - `bash ./scripts/chaincode.sh [version]`
  - example : `bash ./scripts/chaincode.sh 1`
- Generate connection profile for both organizations: `bash scripts/ccp-generate.sh`
  - ccp at : `organizations/peerOrganizations/org[1,2].example.com/connection-org[1,2].json`

## Identity

- Enroll a identity : `bash scripts/identity.sh [1|2] [identity_name] [identity_secret]`
  - Example : `bash scripts/identity.sh 1 trustee_1 pw`

## Clean up

```
$ bash scripts/clean.sh
```