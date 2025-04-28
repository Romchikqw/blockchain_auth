package abac

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"strings"
)

type Claims struct {
	Username  string `json:"username"`
	Role      string `json:"role"`
	IPAddress string `json:"ip"`
	Timezone  string `json:"timezone"`
	LoginTime string `json:"login_time"`
	jwt.RegisteredClaims
}

const PublicKeyPEM = `
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAuPucejcLqWAd3hDH0DQx
CTF4IaCIy7ghauUpops2DsHiq13yrUvR7+ZYHSEVYKhgKB210sqS5kzFI+jGpORN
E39VOsVBItEubexCgXB7IjaghjWdxpJyzMufPlkmETQN82fmEdvQWrPvmT6j449V
Q4UzY2zhCakcQGQKmip7m7K/IRPa5yuUFRk176IQMaQAvZJbFWft1kqYGo/TLEil
hkKK+8t7g97/8OaD6g5wBYF8wfFdPmik7juDxnZWJPVeNtshJ5SBbtqRhg15FPcF
1WEpNdNNqP3+gjV1akT/EapaHJQKo4VuKukY672IhB4yvHAFmZBT6nu8g2XeabuJ
YQIDAQAB
-----END PUBLIC KEY-----
`

func parsePublicKey() (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(PublicKeyPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pub.(*rsa.PublicKey), nil
}

func VerifyJWT(tokenString string) (*Claims, error) {
	tokenString = strings.TrimSpace(tokenString)
	pubKey, err := parsePublicKey()
	if err != nil {
		return nil, err
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return pubKey, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}
	return claims, nil
}
