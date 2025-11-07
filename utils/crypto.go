package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"math/big"
)

var (
	EMAIL_SENDING_KEY = "EmailSendingMediorbit"
)

func Sha256(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}

// MakeEncryptKey a 32-byte key using SHA-256
func MakeEncryptKey(key string) []byte {
	// Use SHA-256 to derive a 32-byte key from the secret key string
	hash := sha256.New()
	hash.Write([]byte(key))
	return hash.Sum(nil)
}

// Encrypt encrypts the given data using AES-GCM.
func Encrypt(data string, key string) (string, error) {
	keyBytes := MakeEncryptKey(key)
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	dataBytes := []byte(data)
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	ciphertext := gcm.Seal(nonce, nonce, dataBytes, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts the given encrypted data using AES-GCM.
func Decrypt(data string, key string) (string, error) {
	keyBytes := MakeEncryptKey(key)
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	ciphertext, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := ciphertext[:gcm.NonceSize()]
	plaintext, err := gcm.Open(nil, nonce, ciphertext[gcm.NonceSize():], nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

func EncryptEmailPassword(password string) (string, error) {
	encryptedPassword, err := Encrypt(password, EMAIL_SENDING_KEY)
	if err != nil {
		return "", err
	}
	return encryptedPassword, nil
}

func DecryptEmailPassword(encryptedPassword string) (string, error) {
	password, err := Decrypt(encryptedPassword, EMAIL_SENDING_KEY)
	if err != nil {
		return "", err
	}
	return password, nil
}

func EncryptId(id int64) string {
	zipNum := ZipNumber(id)
	newId := big.NewInt(id)
	newId.Mul(newId, big.NewInt(10))
	newId.Add(newId, big.NewInt(int64(zipNum)))
	return Base10ToBase58(newId)
}

func DecryptId(encryptedId string) int64 {
	decodeId := Base58ToBase10(encryptedId)
	decodeId.Div(decodeId, big.NewInt(10))
	return decodeId.Int64()
}

func MD5(data string) string {
	hash := md5.Sum([]byte(data)) // Calculate the MD5 hash of the string
	return hex.EncodeToString(hash[:])
}
