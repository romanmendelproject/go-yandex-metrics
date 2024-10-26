// Модуль для получения хеша
package crypto

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"os"
)

// GetHash получение GetHash по ключу
func GetHash(metrics []byte, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(metrics)

	sum := h.Sum(nil)
	return base64.StdEncoding.EncodeToString(sum)
}

func Encrypt(publicKeyPath, plainText string) (string, error) {
	bytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return "", err
	}

	publicKey, err := convertBytesToPublicKey(bytes)
	if err != nil {
		return "", err
	}

	cipher, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(plainText))
	if err != nil {
		return "", err
	}

	return cipherToPemString(cipher), nil
}

func convertBytesToPublicKey(keyBytes []byte) (*rsa.PublicKey, error) {
	var err error

	block, _ := pem.Decode(keyBytes)
	blockBytes := block.Bytes

	cert, err := x509.ParseCertificate(blockBytes)
	if err != nil {
		return nil, err
	}

	return cert.PublicKey.(*rsa.PublicKey), nil
}

func cipherToPemString(cipher []byte) string {
	return string(
		pem.EncodeToMemory(
			&pem.Block{
				Type:  "MESSAGE",
				Bytes: cipher,
			},
		),
	)
}

func Decrypt(privateKeyPath, encryptedMessage string) (string, error) {
	bytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return "", err
	}

	privateKey, err := convertBytesToPrivateKey(bytes)
	if err != nil {
		return "", err
	}

	plainMessage, err := rsa.DecryptPKCS1v15(
		rand.Reader,
		privateKey,
		pemStringToCipher(encryptedMessage),
	)

	return string(plainMessage), err
}

func convertBytesToPrivateKey(keyBytes []byte) (*rsa.PrivateKey, error) {
	var err error

	block, _ := pem.Decode(keyBytes)
	blockBytes := block.Bytes

	privateKey, err := x509.ParsePKCS1PrivateKey(blockBytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func pemStringToCipher(encryptedMessage string) []byte {
	b, _ := pem.Decode([]byte(encryptedMessage))

	return b.Bytes
}
