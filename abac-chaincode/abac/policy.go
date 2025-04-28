package abac

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func (s *SmartContract) AddPolicy(ctx contractapi.TransactionContextInterface, policyID string, policyJSON string) error {
	var pol Policy
	if err := json.Unmarshal([]byte(policyJSON), &pol); err != nil {
		return fmt.Errorf("invalid policy JSON: %v", err)
	}

	pol.ID = policyID

	polBytes, err := json.Marshal(pol)
	if err != nil {
		return fmt.Errorf("failed to marshal policy: %v", err)
	}

	return ctx.GetStub().PutState("policy:"+policyID, polBytes)
}

func (s *SmartContract) GetPolicy(ctx contractapi.TransactionContextInterface, policyID string) (*Policy, error) {
	polBytes, err := ctx.GetStub().GetState("policy:" + policyID)
	if err != nil || polBytes == nil {
		return nil, fmt.Errorf("policy not found: %s", policyID)
	}

	var pol Policy
	if err := json.Unmarshal(polBytes, &pol); err != nil {
		return nil, fmt.Errorf("failed to parse policy: %v", err)
	}

	return &pol, nil
}
