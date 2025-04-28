package abac

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func (s *SmartContract) VerifyAccessToken(ctx contractapi.TransactionContextInterface, token string) (bool, error) {
	data, err := ctx.GetStub().GetState("access_token_" + token)
	if err != nil || data == nil {
		return false, nil
	}

	var at AccessToken
	if err := json.Unmarshal(data, &at); err != nil {
		return false, err
	}

	exp, err := time.Parse(time.RFC3339, at.ExpiresAt)
	if err != nil {
		return false, fmt.Errorf("invalid expiration format")
	}

	if time.Now().After(exp) {
		at.Valid = false
		updated, _ := json.Marshal(at)
		_ = ctx.GetStub().PutState("access_token_"+token, updated)
		return false, nil
	}

	if !at.Valid {
		return false, nil
	}

	return true, nil
}

func (s *SmartContract) GetAccessToken(ctx contractapi.TransactionContextInterface, token string) (*AccessToken, error) {
	data, err := ctx.GetStub().GetState("access_token_" + token)
	if err != nil || data == nil {
		return nil, fmt.Errorf("token not found")
	}
	var at AccessToken
	err = json.Unmarshal(data, &at)
	if err != nil {
		return nil, err
	}
	return &at, nil
}
