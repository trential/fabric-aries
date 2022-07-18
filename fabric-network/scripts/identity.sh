# /bin/bash

source ./scripts/env.sh


ORG1_CA_URL=localhost:7054
ORG2_CA_URL=localhost:8054

# 1 for Org1
# 2 for Org2
ORG=${1-1}

ORG_CA_URL=""

case $ORG in
    "1")
        ORG_CA_URL=$ORG1_CA_URL
    ;;
    "2")
        ORG_CA_URL=$ORG2_CA_URL
    ;;
esac

export FABRIC_CA_CLIENT_HOME=organizations/peerOrganizations/org${ORG}.example.com/registrar

if [ ! -d $FABRIC_CA_CLIENT_HOME ];then
    fabric-ca-client enroll -u http://admin:adminpw@${ORG_CA_URL}
fi

# user name of the identity
ID=$2
# secret of the identity
SECRET=$3


# register identity
fabric-ca-client identity add $ID --secret $SECRET --type client