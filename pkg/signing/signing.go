package signing

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

var signingKey []byte = []byte("password")

func ValidateToken(tokenFromAgent string) (int, error) {
	tkn, err := jwt.Parse(tokenFromAgent, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})
	if err != nil {
		return -1, err
	}

	if claims, ok := tkn.Claims.(jwt.MapClaims); ok && tkn.Valid {
		return int(claims["runnerId"].(float64)), nil
	}

	return -1, fmt.Errorf("invalid signature")
}

func GenerateToken(runnerId int) (string, error) {
	// Generate new  JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"runnerId": runnerId,
	})

	// Update struct
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}

	return tokenString, err
}
