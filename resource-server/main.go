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

type Claims struct {
	Username  string `json:"username"`
	IPAddress string `json:"ip"`
	Timezone  string `json:"timezone"`
	LoginTime string `json:"login_time"`
	jwt.RegisteredClaims
}

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
	r.GET("/resource", resourceHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}
	fmt.Printf("Resource server running on port %s\n", port)
	r.Run(":" + port)
}

func resourceHandler(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
		return
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
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
	accessToken := c.Query("access_token")
	if resourceID == "" || accessToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing resource ID or access token"})
		return
	}

	if !strings.HasPrefix(accessToken, claims.Username) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access token does not match user"})
		return
	}

	valid, err := verifyAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !valid {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Access granted to resource: %s", resourceID)})
}

func verifyAccessToken(accessToken string) (bool, error) {
	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		return false, err
	}

	ccpPath := "connection-org1.yaml"
	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, "appUser"),
	)
	if err != nil {
		return false, err
	}
	defer gw.Close()

	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		return false, err
	}

	contract := network.GetContract("abac")

	result, err := contract.EvaluateTransaction("VerifyAccessToken", accessToken)
	if err != nil {
		return false, err
	}

	return string(result) == "true", nil
}
