// Модуль для получения хеша
package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

// GetHash получение GetHash по ключу
func GetHash(metrics []byte, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(metrics)

	sum := h.Sum(nil)
	return base64.StdEncoding.EncodeToString(sum)
}
