package crypto

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCryptography(t *testing.T) {
	var metrics []byte
	key := "test"

	hash := GetHash(metrics, key)

	h := hmac.New(sha256.New, []byte(key))

	require.Equal(t, hash, base64.StdEncoding.EncodeToString(h.Sum(nil)))
}

func Example() {
	var metrics []byte
	key := "test"

	fmt.Println(GetHash(metrics, key))
	// Output:
	// rXEUjHnyGrnuxR6lx90rZoeS98DTU0rmayL3HGFSP7M=
}

func TestCipherToPemString(t *testing.T) {
	cipher := []byte("test message")
	expectedPem := "-----BEGIN MESSAGE-----\ndGVzdCBtZXNzYWdl\n-----END MESSAGE-----\n"

	result := cipherToPemString(cipher)

	if !bytes.Equal([]byte(result), []byte(expectedPem)) {
		t.Errorf("Expected:\n%s\nGot:\n%s", expectedPem, result)
	}
}

func TestPemStringToCipher(t *testing.T) {
	// Создаем тестовую строку PEM
	testCipher := []byte("some encrypted data")
	pemBlock := &pem.Block{
		Type:  "MESSAGE",
		Bytes: testCipher,
	}
	pemEncoded := string(pem.EncodeToMemory(pemBlock))

	// Декодируем с помощью функции
	result := pemStringToCipher(pemEncoded)

	// Проверяем, что полученный результат соответствует исходным данным
	if string(result) != string(testCipher) {
		t.Errorf("Expected %s, but got %s", testCipher, result)
	}
}

// Функция для создания тестового RSA-ключа
func generateTestRSAKey() ([]byte, error) {
	// Генерация RSA ключа
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	// Печать ключа в формате PEM
	privKeyBytes := x509.MarshalPKCS1PrivateKey(privKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privKeyBytes,
	}
	return pem.EncodeToMemory(block), nil
}

func TestConvertBytesToPublicKey(t *testing.T) {
	// Example PEM-encoded RSA public key
	pemKey := `-----BEGIN CERTIFICATE-----
MIIFMDCCAxigAwIBAgICBnowDQYJKoZIhvcNAQELBQAwKDELMAkGA1UEBhMCUlUx
GTAXBgNVBAoTEFlhbmRleC5QcmFrdGlrdW0wHhcNMjMxMDE2MTY0NDI1WhcNMzMx
MDE2MTY0NDI1WjAoMQswCQYDVQQGEwJSVTEZMBcGA1UEChMQWWFuZGV4LlByYWt0
aWt1bTCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoCggIBAMbo1k9s3EWQ8hz8
T5RYHIjL7HW+dee32LeowGddPPAfC2oocA8bbOmdTPhVKypsdpNoILJFupT/Co5L
FKDR0hu+vqk+zXfamgxe4bESGBH0PMjVPY9T3I24pBV60QKsKQ8/h4zIAyrCmzlA
TPag8uvhxYtMlte763GcHfqBFTZsP1gy3W+DQ0Hsn6q8c8rDI2i8JMuYj8VlPDlC
k64Iso3ssKmKq1vUDKf1yVPdXZlAkmdt2NuHu6QpHrKbHnTkGp+gjZyHUKwdhGzp
iGC0uxPVpCwo11C9zjtSRn2o84i7z7xr+mu7+t48+nHI6uPFk98y9my4hViXKsK5
D99aTLAi9Uil9u3qPSRcTn8/lZHklewYUGlCszEBAdbR8TaRfybgXSouTxjiwpHV
GME/wlPgGU6uzjdrMGeFZEJR0AWPqPV3wPpOexUd7mYccpW5tGN0seqqbZttjkA0
zEIE520mIfiTMpRx4V9gxlAKqQXag6SdgvJk8OM3UHqyCo3zuG/xMASrhTVvhOEv
8rOOgjOuQLU9sMqXFX05ejzWq7S0kEc595fyyy98vxxPxG268U7rHREFXpmXJ04Z
m8pWv+1h6kniR6J0fKiuI8QoRag+RalgGf5KkFvtabecf4kZOpqzi1d7nvfkL+Dy
v04S89ozP0f7L38wbEaqmHv8w3AFAgMBAAGjZDBiMA4GA1UdDwEB/wQEAwIHgDAd
BgNVHSUEFjAUBggrBgEFBQcDAgYIKwYBBQUHAwEwDgYDVR0OBAcEBQECAwQGMCEG
A1UdEQQaMBiHBH8AAAGHEAAAAAAAAAAAAAAAAAAAAAEwDQYJKoZIhvcNAQELBQAD
ggIBAJY1NrLgWrYAJij5Z/wpvXfiMb9zyRatggkMom7bPUoghLfmCfF+cQ0EqLKY
7bbjEsB7YRWslapKlqf+DlaSQxIvNjiuY2wpXOKoxrXEDfcZqAMOh6NAbUhfZ+Rp
9MTCvEkxokMMQy5C6RvFRwahAXGickw8SXVDJpb1lgkvbcp+jURbz63NOcMOV+Wj
k6VIqllv/S5LmaEYCFtnsWy4GsM3zeQqT7xsKEMU1IcYrIjbZ/ydY7kGNGz3WiMf
3OWsFEBuFKHcy+jDp9NHQ3EKmkGLQAxUgk0Q2TyuNNumBGgqKyGr784q3cMwNxxJ
hSpbuGf2VbV+cQnYS6+f3GoUWruFyeQZ47QuUaoQ7Z2F5KLownF1gD8i05c8UQH/
t5hU15Y+QqtWe+6MSbz7J/FrM/7sO0lqDQvtE85a565/QZvMZyIWxo1HcRfc5kgT
4G5xOt7jLo8nBPQJq/g1Vg8XbQM7IPunuWU4Mldt3CIGGugIJrNN3YTF8X6ydYBi
ttqYBSSordiEaPwKwX7J2E0zwr48Kgn52lZ+yDZV1Gjpxx8shvbdLkpKSIx91PKO
94rOyT2nJObbvFtDo/J4dP/Pz4hnNtjs+ByDsinB0mygwVD6TdEtzj33nzEBSv5T
fSVgV1aWNMXG3XNnXHYIKiE5iL/hCnJ/tkcM5meXXNZAFzvm
-----END CERTIFICATE-----`
	keyBytes := []byte(pemKey)

	publicKey, err := convertBytesToPublicKey(keyBytes)

	assert.NoError(t, err)
	assert.NotNil(t, publicKey)
	assert.IsType(t, &rsa.PublicKey{}, publicKey)
}

func TestConvertBytesToPrivateKey(t *testing.T) {
	// Example PEM-encoded RSA private key
	pemKey := `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`
	keyBytes := []byte(pemKey)

	privateKey, err := convertBytesToPrivateKey(keyBytes)

	assert.NoError(t, err)
	assert.NotNil(t, privateKey)
	assert.IsType(t, &rsa.PrivateKey{}, privateKey)
}
