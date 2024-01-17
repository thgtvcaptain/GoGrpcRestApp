package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
	"os"
	"strings"
)

func checkPassword(hashedPassword string, normalPassword string) int {
	decrypted := decryptMessage(hashedPassword)
	n := strings.Compare(decrypted, normalPassword)

	return n
}

func encryptMessage(plainText string) string {

	byteText := []byte(plainText)
	cipherText := make([]byte, len(plainText))

	// must be 16 char or 24 char or 32 char
	secretKey := []byte(os.Getenv("SECRET_KEY"))

	block, _ := aes.NewCipher(secretKey)

	nonce := make([]byte, block.BlockSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	nonceString := hex.EncodeToString(nonce)

	stream := cipher.NewCTR(block, nonce)

	stream.XORKeyStream(cipherText, byteText)

	return nonceString + hex.EncodeToString(cipherText)
}

func decryptMessage(cipherText string) string {

	secretKey := []byte(os.Getenv("SECRET_KEY"))

	block, _ := aes.NewCipher(secretKey)

	nonceSize := block.BlockSize()
	nonce, extracted := cipherText[:nonceSize*2], cipherText[nonceSize*2:]

	byteNonce, _ := hex.DecodeString(nonce)

	stream := cipher.NewCTR(block, byteNonce)

	extractedDecoded, _ := hex.DecodeString(extracted)
	plainText := make([]byte, len(extractedDecoded))

	stream.XORKeyStream(plainText, extractedDecoded)

	plainTextString := string(plainText)

	return plainTextString
}
