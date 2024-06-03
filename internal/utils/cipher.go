package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"os"
)

var bytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

var secretKey string = os.Getenv("ENCRYPTION_KEY")

func encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}
func decode(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

// Encrypt method is to encrypt or hide any classified text
func Encrypt(text string) (string, error) {
	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", err
	}
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, bytes)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return encode(cipherText), nil
}

// Decrypt method is to extract back the encrypted text
func Decrypt(text string) (string, error) {
	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", err
	}
	cipherText := decode(text)
	cfb := cipher.NewCFBDecrypter(block, bytes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}

// func main() {
// 	StringToEncrypt := "Encrypting this string"
// 	// To encrypt the StringToEncrypt
// 	encText, err := encrypt(StringToEncrypt, secretKey)
// 	if err != nil {
// 		fmt.Println("error encrypting your classified text: ", err)
// 	}
// 	fmt.Println(encText)
// 	// To decrypt the original StringToEncrypt
// 	decText, err := decrypt("Li5E8RFcV/EPZY/neyCXQYjrfa/atA==", secretKey)
// 	if err != nil {
// 		fmt.Println("error decrypting your encrypted text: ", err)
// 	}
// 	fmt.Println(decText)
// }
