package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	config "github.com/ochiengotieno304/oneotp/internal/configs"
	"github.com/ochiengotieno304/oneotp/internal/utils"
	"github.com/ochiengotieno304/oneotp/internal/utils/errors"
	"github.com/ochiengotieno304/oneotp/pkg/db/stores"
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

func ValidateRequest(clientID, secret string) error {
	accountStore := stores.NewAccountStore()

	account, err := accountStore.FindAccountByID(clientID)
	if err != nil {
		return err
	}

	decreptedKey, err := utils.Decrypt(account.Credentials.SecretKey)
	if err != nil {
		return err
	}

	verifySecret := secret == decreptedKey
	unverified := account.Status != 200 // status 100 for unverified accounts, 200 for verified, 300 for revoked

	if !verifySecret {
		return errors.ErrSecretVerification
	}

	if unverified {
		return errors.ErrUnverifiedAccount
	}
	return nil
}
