package abac

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func (s *SmartContract) AddAttribute(ctx contractapi.TransactionContextInterface, username string, attrJSON string) error {
	var attr Attribute
	if err := json.Unmarshal([]byte(attrJSON), &attr); err != nil {
		return fmt.Errorf("invalid attribute JSON: %v", err)
	}

	attrBytes, err := json.Marshal(attr)
	if err != nil {
		return fmt.Errorf("failed to marshal attributes: %v", err)
	}

	return ctx.GetStub().PutState("attributes:"+username, attrBytes)
}

func (s *SmartContract) GetAttribute(ctx contractapi.TransactionContextInterface, username string) (*Attribute, error) {
	attrBytes, err := ctx.GetStub().GetState("attributes:" + username)
	if err != nil || attrBytes == nil {
		return nil, fmt.Errorf("attribute not found for user: %s", username)
	}

	var attr Attribute
	if err := json.Unmarshal(attrBytes, &attr); err != nil {
		return nil, fmt.Errorf("failed to parse attribute: %v", err)
	}

	return &attr, nil
}
