package auth

import (
	"math/rand"
)

func GenerateSecretKey() string {
	const charset = "aAbBcCdDeEfFgGhHiIjJkKlLmMnNoOpPqQrRsStTuUvVwWxXyYzZ0123456789"
	result := make([]byte, 20)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}


