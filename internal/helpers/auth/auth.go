package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	config "github.com/ochiengotieno304/oneotp/internal/configs"
)

func GenerateToken(accountID string) (string, error) {
	configs, err := config.LoadConfig()
	if err != nil {
		return "", err
	}

	mySigningKey := []byte(configs.JWTSecretKey)

	type MyCustomClaims struct {
		AccountID string `json:"account_id"`
		jwt.RegisteredClaims
	}

	// Create claims with multiple fields populated
	claims := MyCustomClaims{
		accountID,
		jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "oneotp",
			ID:        "1",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}

	return ss, nil
}
