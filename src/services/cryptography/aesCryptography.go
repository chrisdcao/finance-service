package cryptography

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"finance-service/utils"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
	"io"
)

// KeyDerivation derives a key from a password and salt using PBKDF2.
func KeyDerivation(password, salt string) []byte {
	return pbkdf2.Key([]byte(password), []byte(salt), 4096, 32, sha256.New)
}

// AES encrypt the wallet id
func Encrypt(toBeEncrypted string) (string, error) {
	password, err := utils.GetEnvVariableFromKey("ENCRYPTION_PASSWORD")
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	salt, err := utils.GetEnvVariableFromKey("ENCRYPTION_SALT")
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	key := KeyDerivation(password, salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	plainTextBytes := []byte(toBeEncrypted)
	cipherText := make([]byte, aes.BlockSize+len(plainTextBytes))
	iv := cipherText[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainTextBytes)

	return base64.URLEncoding.EncodeToString(cipherText), nil
}

// AES decrypt the wallet id
func Decrypt(toBeDecrypted string) (string, error) {
	password, err := utils.GetEnvVariableFromKey("ENCRYPTION_PASSWORD")
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	salt, err := utils.GetEnvVariableFromKey("ENCRYPTION_SALT")
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	key := KeyDerivation(password, salt)
	cipherTextBytes, err := base64.URLEncoding.DecodeString(toBeDecrypted)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(cipherTextBytes) < aes.BlockSize {
		return "", fmt.Errorf("cipherText too short")
	}

	iv := cipherTextBytes[:aes.BlockSize]
	cipherTextBytes = cipherTextBytes[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherTextBytes, cipherTextBytes)

	return string(cipherTextBytes), nil
}
