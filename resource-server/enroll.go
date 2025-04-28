package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

func main() {
	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		log.Fatalf("❌ Failed to create wallet: %v", err)
	}

	// Удалим, если уже существует (по желанию можно закомментировать)
	if wallet.Exists("appUser") {
		fmt.Println("⚠️  appUser already enrolled. Removing and re-enrolling...")
		_ = wallet.Remove("appUser")
	}

	certPath := filepath.Join(
		"../fabric-samples-main/test-network/organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/signcerts",
		"cert.pem",
	)
	keyDir := filepath.Join(
		"../fabric-samples-main/test-network/organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/keystore",
	)

	files, err := filepath.Glob(filepath.Join(keyDir, "*_sk"))
	if err != nil || len(files) == 0 {
		log.Fatalf("❌ Private key not found: %v", err)
	}
	keyPath := files[0]

	certBytes, err := os.ReadFile(certPath)
	if err != nil {
		log.Fatalf("❌ Failed to read cert: %v", err)
	}

	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		log.Fatalf("❌ Failed to read key: %v", err)
	}

	identity := gateway.NewX509Identity("Org1MSP", string(certBytes), string(keyBytes))

	err = wallet.Put("appUser", identity)
	if err != nil {
		log.Fatalf("❌ Failed to put identity: %v", err)
	}

	fmt.Println("✅ Successfully enrolled appUser")
}
