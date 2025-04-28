package main

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	Username  string `json:"username"`
	Role      string `json:"role"`
	IPAddress string `json:"ip"`
	Timezone  string `json:"timezone"`
	LoginTime string `json:"login_time"`
	jwt.RegisteredClaims
}

func main() {
	privKeyBytes, err := os.ReadFile("fake.key")
	if err != nil {
		panic("cannot read fake.key: " + err.Error())
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(privKeyBytes)
	if err != nil {
		panic("cannot parse fake.key: " + err.Error())
	}

	now := time.Now().UTC()
	claims := Claims{
		Username:  "admin",
		Role:      "admin",
		IPAddress: "127.0.0.1",
		Timezone:  "UTC",
		LoginTime: now.Format(time.RFC3339),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(1 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(privKey)
	if err != nil {
		panic("signing error: " + err.Error())
	}

	fmt.Println(tokenString)
}
