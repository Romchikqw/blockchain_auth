#!/bin/bash

ORG=$1
CHANNEL_BLOCK=./channel-artifacts/mychannel.block


PEER_CONTAINER=peer0.${ORG}.example.com
MSPID="$(tr '[:lower:]' '[:upper:]' <<< ${ORG:0:1})${ORG:1}MSP"
PEER_PORT="7051"
[[ "$ORG" == "org2" ]] && PEER_PORT="9051"

echo "✅ Копирую admin MSP для $ORG..."
docker cp organizations/peerOrganizations/${ORG}.example.com/users/Admin@${ORG}.example.com/msp ${PEER_CONTAINER}:/admin-msp

echo "✅ Копирую блок канала..."
docker cp "$CHANNEL_BLOCK" ${PEER_CONTAINER}:/mychannel.block

echo "✅ Подключаюсь к контейнеру $PEER_CONTAINER..."

docker exec -e CORE_PEER_LOCALMSPID=$MSPID \
            -e CORE_PEER_MSPCONFIGPATH=/admin-msp \
            -e CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/tls/ca.crt \
            -e CORE_PEER_TLS_ENABLED=true \
            -e CORE_PEER_ADDRESS=${PEER_CONTAINER}:${PEER_PORT} \
            -it $PEER_CONTAINER \
            peer channel join -b /mychannel.block
