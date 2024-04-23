package jwt

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func BuildToken(data map[string]interface{}) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, mapData(data))

	pkey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(os.Getenv("JWT_PRIVATE_KEY")))
	if err != nil {
		return "", err
	}

	tokenString, err := token.SignedString(pkey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func mapData(data map[string]interface{}) jwt.MapClaims {
	mapClaim := jwt.MapClaims{
		"nbf": time.Now().Unix(),
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(3 * time.Hour).Unix(),
		"iss": "SawitPro-Challenge",
		"aud": "SawitPro-Challenge",
	}

	for key, value := range data {
		mapClaim[key] = value
	}

	return mapClaim
}
