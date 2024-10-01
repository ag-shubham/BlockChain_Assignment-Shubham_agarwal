package main 

import (
	"log"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"products/contracts"
)

func main(){
	assetChaincode, err := contractapi.NewChaincode(&contracts.SmartContract{})
	if err != nil {
		log.Panicf("Error creating asset-transfer-basic chaincode: %v", err)
	}

	if err := assetChaincode.Start(); err != nil {
		log.Panicf("Error starting asset-transfer-basic chaincode: %v", err)
	}
}
