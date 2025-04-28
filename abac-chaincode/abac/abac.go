package abac

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-contract-api-go/metadata"
	"strings"
	"time"
)

type SmartContract struct {
	contractapi.Contract
}

type Attribute struct {
	Role       string `json:"role"`
	Department string `json:"department"`
}

type Policy struct {
	ID           string   `json:"id"`
	Resources    []string `json:"resources"`
	AllowedRoles []string `json:"allowed_roles"`
	AllowedIP    string   `json:"ip"`
	AllowedHours string   `json:"allowed_hours"`
}
type AccessToken struct {
	Token      string `json:"token"`
	Username   string `json:"username"`
	ResourceID string `json:"resource_id"`
	IssuedAt   string `json:"issued_at"`
	ExpiresAt  string `json:"expires_at"`
	Valid      bool   `json:"valid"`
}

func generateUUIDv4() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	// set version (4) and variant bits
	b[6] = (b[6] & 0x0f) | 0x40 // Version 4
	b[8] = (b[8] & 0x3f) | 0x80 // Variant

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	return nil
}

func (s *SmartContract) CheckAccess(ctx contractapi.TransactionContextInterface, jwtToken string, resourceID string) (string, error) {
	claims, err := VerifyJWT(jwtToken)
	if err != nil {
		return "8", fmt.Errorf("JWT verification failed: %v", err)
	}

	username := claims.Username
	ip := claims.IPAddress
	loginTime := claims.LoginTime

	attrBytes, err := ctx.GetStub().GetState("attributes:" + username)
	if err != nil || attrBytes == nil {
		return "7", fmt.Errorf("user attributes not found")
	}
	var attrs Attribute
	if err := json.Unmarshal(attrBytes, &attrs); err != nil {
		return "6", fmt.Errorf("failed to parse attributes: %v", err)
	}

	iter, err := ctx.GetStub().GetStateByRange("policy:", "policy;")
	if err != nil {
		return "5", fmt.Errorf("failed to get policies: %v", err)
	}
	defer iter.Close()

	var matchedPolicy *Policy

	for iter.HasNext() {
		entry, err := iter.Next()
		if err != nil {
			continue
		}

		var pol Policy
		if err := json.Unmarshal(entry.Value, &pol); err != nil {
			continue
		}

		for _, res := range pol.Resources {
			if res == resourceID {
				matchedPolicy = &pol
				break
			}
		}
		if matchedPolicy != nil {
			break
		}
	}

	if matchedPolicy == nil {
		return "1", fmt.Errorf("no matching policy for resource: %s", resourceID)
	}

	roleAllowed := false
	for _, r := range matchedPolicy.AllowedRoles {
		if r == attrs.Role {
			roleAllowed = true
			break
		}
	}

	if !roleAllowed {
		return "2", nil
	}
	fmt.Printf("Comparing IPs: tokenIP=%q, policyIP=%q\n", ip, matchedPolicy.AllowedIP)

	if matchedPolicy.AllowedIP != "" && matchedPolicy.AllowedIP != ip {
		return "3", nil
	}

	if matchedPolicy.AllowedHours != "" {
		var from, to string
		normalized := strings.ReplaceAll(matchedPolicy.AllowedHours, "–", "-") // normalize long dash
		fmt.Sscanf(normalized, "%5s-%5s", &from, &to)
		t, _ := time.Parse(time.RFC3339, loginTime)
		hhmm := t.Format("15:04")

		if hhmm < from || hhmm > to {
			return "4", nil
		}

	}

	accessToken := username + "-" + loginTime
	at := AccessToken{
		Token:      accessToken,
		Username:   username,
		ResourceID: resourceID,
		IssuedAt:   time.Now().UTC().Format(time.RFC3339),
		ExpiresAt:  time.Now().UTC().Add(15 * time.Minute).Format(time.RFC3339),
		Valid:      true,
	}

	data, err := json.Marshal(at)
	if err != nil {
		return "9", fmt.Errorf("failed to marshal access token")
	}

	err = ctx.GetStub().PutState("access_token_"+accessToken, data)
	if err != nil {
		return "10", fmt.Errorf("failed to store access token")
	}

	return accessToken, nil
}

func (s *SmartContract) GetInfo() metadata.InfoMetadata {
	return metadata.InfoMetadata{
		Title:       "ABACChaincode",
		Description: "Chaincode for Attribute-Based Access Control",
		Version:     "0.0.1",
	}
}
