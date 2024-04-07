package util

import (
	"crypto/sha256"
	"encoding/base64"
)

// Функция для генерации короткого идентификатора на основе URL
func GenerateShortID(url string) string {
	hash := sha256.New()
	hash.Write([]byte(url))
	// Преобразуем хэш в строку и кодируем в base64, чтобы получить короткий идентификатор
	shortID := base64.URLEncoding.EncodeToString(hash.Sum(nil))
	return shortID[:8]
}
