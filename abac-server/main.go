package main

import (
	"crypto/rsa"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/joho/godotenv"
)

var publicKey *rsa.PublicKey

func main() {
	_ = godotenv.Load()

	keyData, err := ioutil.ReadFile("public.key")
	if err != nil {
		log.Fatalf("Error reading public key: %v", err)
	}
	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(keyData)
	if err != nil {
		log.Fatalf("Error parsing public key: %v", err)
	}

	r := gin.Default()

	r.GET("/check-access", checkAccessHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}
	fmt.Printf("ABAC server running on port %s\n", port)
	r.Run(":" + port)
}

func checkAccessHandler(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	resourceID := c.Query("resource")
	if resourceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing resource ID"})
		return
	}

	accessToken, err := queryBlockchain(tokenString, resourceID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"expires_in":   "15m",
	})
}

func queryBlockchain(token string, resourceID string) (string, error) {
	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		return "", err
	}

	ccpPath := "connection-org1.yaml"

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, "appUser"),
	)
	if err != nil {
		return "", err
	}
	defer gw.Close()

	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		return "", err
	}

	contract := network.GetContract("abac")

	result, err := contract.SubmitTransaction("CheckAccess", token, resourceID)
	if err != nil {
		return "", err
	}

	return string(result), nil
}
