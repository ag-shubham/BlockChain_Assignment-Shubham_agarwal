# Hyperledger Fabric Supply Chain Network

## Overview
This project sets up a Hyperledger Fabric network for a supply chain system, featuring two organizations: Org4 and Org5. The network facilitates secure communication and data sharing between these organizations.

## Features
- **Two Organizations**: Org4 and Org5 with separate peers.
- **Shared Orderer Service**: Ensures consistent transaction ordering.
- **Communication Channel**: Named `mychannel` for inter-organizational data sharing.

## Architecture
The network architecture consists of:
- **Organizations**: Org4 and Org5
- **Peers**: Peer0 for each organization
- **Orderer**: A shared service to manage transactions
- **Channel**: `mychannel` connecting both organizations

## Getting Started

### Prerequisites
- Docker and Docker Compose installed
- Go (Golang) installed for chaincode development
- Hyperledger Fabric binaries available

### Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/ag-shubham/BlockChain_Assignment-Shubham_agarwal.git

   cd BlockChain_Assignment-Shubham_agarwal
   ```

2. Go to the `test-network-assignment` directory:
    ```bash
    cd test-network-assignment
    ```

3. Certificate Generation Using CryptoGen Tool:
    
    1. Set the Config Path
    ```bash
    export PATH=${PWD}/../bin:${PWD}:$PATH
    export DOCKER_SOCK=/var/run/docker.sock
    IMAGE_TAG=latest docker-compose -f compose/compose-test-net.yaml -f compose/docker/docker-compose-test-net.yaml up
    ```
    2. Generate Certificates for Org4 and Org5
    ```bash
    cryptogen generate --config=./organizations/cryptogen/crypto-config-org4.yaml --output="organizations"
    cryptogen generate --config=./organizations/cryptogen/crypto-config-org5.yaml --output="organizations"
    ```
    3. Generate Certificates for Orderer
    ```bash
    cryptogen generate --config=./organizations/cryptogen/crypto-config-orderer.yaml --output="organizations"
    ```

4. Start docker containers:
    ```bash
    export DOCKER_SOCK=/var/run/docker.sock
    IMAGE_TAG=latest docker-compose -f compose/compose-test-net.yaml -f compose/docker/docker-compose-test-net.yaml up
    ```

5. Creating Channel

    1. Set the Config Path
    ```bash
    export PATH=${PWD}/../bin:${PWD}:$PATH
    export FABRIC_CFG_PATH=${PWD}/configtx
    export CHANNEL_NAME=mychannel
    ```
    2. Create the System Genesis  Block and Channel Genesis block
    ```bash
    configtxgen -profile TwoOrgsApplicationGenesis -outputBlock ./channel-artifacts/${CHANNEL_NAME}.block -channelID $CHANNEL_NAME
    ```
    3. Convert Block to JSON format to understand the data inside it
    ```bash
    configtxgen -inspectBlock ./channel-artifacts/mychannel.block > dump.json
    ```
    4. Copy some prerequisites
    ```bash
    cp ../config/core.yaml ./configtx/.
    ```
    5. Create the Channel
    ```bash
    export ORDERER_CA=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
    export ORDERER_ADMIN_TLS_SIGN_CERT=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/server.crt
    export ORDERER_ADMIN_TLS_PRIVATE_KEY=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/server.key

    osnadmin channel join --channel-id $CHANNEL_NAME --config-block ./channel-artifacts/${CHANNEL_NAME}.block -o localhost:7053 --ca-file "$ORDERER_CA" --client-cert "$ORDERER_ADMIN_TLS_SIGN_CERT" --client-key "$ORDERER_ADMIN_TLS_PRIVATE_KEY"
    ```
    6. Joining Peers to the Channel
    ```bash
    source ./scripts/setOrgPeerContext.sh 1
    peer channel join -b ./channel-artifacts/mychannel.block
    source ./scripts/setOrgPeerContext.sh 2
    peer channel join -b ./channel-artifacts/mychannel.block
    ```
    7. Update Anchor Peers for Org4 and Org5
    ```bash
    source ./scripts/setOrgPeerContext.sh 1
    docker exec cli ./scripts/setAnchorPeer.sh 1 $CHANNEL_NAME

    source ./scripts/setOrgPeerContext.sh 2
    docker exec cli ./scripts/setAnchorPeer.sh 2 $CHANNEL_NAME
    ```

6. Install Chaincode

    1. Update the environment variable to configure the use of GoLang Chaincode
    ```bash
    source ./scripts/setFabCarGolangContext.sh
    export FABRIC_CFG_PATH=$PWD/../config/
    export FABRIC_CFG_PATH=${PWD}/configtx
    export CHANNEL_NAME=mychannel
    export PATH=${PWD}/../bin:${PWD}:$PATH
    ```
    2. Package the Chaincode
    ```bash
    source ./scripts/setOrgPeerContext.sh 1
    peer lifecycle chaincode package fabcar.tar.gz --path ${CC_SRC_PATH} --lang ${CC_RUNTIME_LANGUAGE} --label fabcar_${VERSION}
    ```
    3. Install the Chaincode on peer0.org4
    ```bash
    peer lifecycle chaincode install fabcar.tar.gz
    ```
    4. Install the Chaincode on peer0.org5
    ```bash
    source ./scripts/setOrgPeerContext.sh 2
    peer lifecycle chaincode install fabcar.tar.gz  
    ```
    5. Query Installed Chaincode
    ```bash
    peer lifecycle chaincode queryinstalled 2>&1 | tee outfile
    ```
    6. Set the Chaincode Package ID
    ```bash
    source ./scripts/setPackageID.sh outfile
    ```
    7. Approve Chaincode Definition for Org4
    ```bash
    source ./scripts/setOrgPeerContext.sh 1
    peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA --channelID $CHANNEL_NAME --name fabcar --version ${VERSION} --init-required --package-id ${PACKAGE_ID} --sequence ${VERSION}
    ```
    8. Approve Chaincode Definition for Org5
    ```bash
    source ./scripts/setOrgPeerContext.sh 2
    peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA --channelID $CHANNEL_NAME --name fabcar --version ${VERSION} --init-required --package-id ${PACKAGE_ID} --sequence ${VERSION}
    ```
    9. Check Commit Readiness
    ```bash
    source ./scripts/setOrgPeerContext.sh 1
    peer lifecycle chaincode checkcommitreadiness --channelID $CHANNEL_NAME --name fabcar --version ${VERSION} --sequence ${VERSION} --output json --init-required
    ```
    10. Set the peer address for identifying the endorsing peers
    ```bash
    source ./scripts/setPeerConnectionParam.sh 1 2
    ```
    11. Commit Chaincode Definition
    ```bash
    peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA --channelID $CHANNEL_NAME --name fabcar $PEER_CONN_PARAMS --version ${VERSION} --sequence ${VERSION} --init-required
    ```
    12. Query Committed Chaincode
    ```bash
    source ./scripts/setOrgPeerContext.sh 1
    peer lifecycle chaincode querycommitted --channelID $CHANNEL_NAME --name fabcar
    ```

7. Invoke Chaincode

    1. The chaincode fabcar, has to be initialized before executing other transactions
    ```bash
    source ./scripts/setPeerConnectionParam.sh 1 2

    source ./scripts/setOrgPeerContext.sh 1

    peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C $CHANNEL_NAME -n fabcar $PEER_CONN_PARAMS --isInit -c '{"function":"initLedger","Args":[]}'
    ```
    2. Create a new car
    ```bash
    source ./scripts/setOrgPeerContext.sh 1

    peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n fabcar $PEER_CONN_PARAMS -c '{"function":"CreateCar","Args":["CAR11","Tata","Safari", "Red", "Ron"]}'
    ```
    3. Query the car
    ```bash
    peer chaincode query -C $CHANNEL_NAME -n fabcar -c '{"Args":["queryAllCars"]}'
    ```



    

