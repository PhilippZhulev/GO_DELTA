package helpers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"

	"github.com/PhilippZhulev/delta/internal/app/confiiguration"
)

// Hesh ...
// Протокол хеша
type Hesh struct{}

var config = confiiguration.NewConfig()

// HashPassword ...
// захешировать пароль
func (h *Hesh) HashPassword(password string) string {

	hesh := hmac.New(sha256.New, []byte(config.Salt))
	hesh.Write([]byte(password))

	return base64.StdEncoding.EncodeToString(hesh.Sum(nil))
}

// CheckPasswordHash ...
// Сравнить пароль
// Сравнивает хеш пароля в базе с хешем
// Поступившего пароля
func (h *Hesh) CheckPasswordHash(password, hash string) bool {
	expected := h.HashPassword(password)
	return hmac.Equal([]byte(hash), []byte(expected))
}
