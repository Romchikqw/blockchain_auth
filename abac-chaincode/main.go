package main

import (
	"chaincode/abac"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"log"
)

func main() {
	contract := new(abac.SmartContract)

	chaincode, err := contractapi.NewChaincode(contract)
	if err != nil {
		log.Fatalf("Error creating chaincode: %v", err)
	}
	if err := chaincode.Start(); err != nil {
		log.Fatalf("Error starting chaincode: %v", err)
	}
}
