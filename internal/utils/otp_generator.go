package utils

import (
	"fmt"
	"math/rand/v2"
)

func GenerateOTP() []string {
	otp := make([]string, 0, 6)

	for _, value := range rand.Perm(6) {
		otp = append(otp, fmt.Sprintf("%v", value))
	}

	return otp
}
