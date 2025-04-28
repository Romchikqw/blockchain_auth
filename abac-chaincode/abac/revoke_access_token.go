package abac

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func (s *SmartContract) RevokeAccessToken(ctx contractapi.TransactionContextInterface, token string) error {
	data, err := ctx.GetStub().GetState("access_token_" + token)
	if err != nil || data == nil {
		return fmt.Errorf("token not found")
	}

	var at AccessToken
	if err := json.Unmarshal(data, &at); err != nil {
		return err
	}

	at.Valid = false
	updated, _ := json.Marshal(at)
	return ctx.GetStub().PutState("access_token_"+token, updated)
}
