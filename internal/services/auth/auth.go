package auth

import (
	"fmt"
	"github.com/go-chi/jwtauth/v5"
	"time"
)

var JWTAuth *jwtauth.JWTAuth

func InitJWTAuth(secretKey []byte) (*jwtauth.JWTAuth, error) {
	JWTAuth = jwtauth.New("HS256", secretKey, nil)
	return JWTAuth, nil
}

func GenerateToken(login string) (string, error) {
	const op = "service.GenerateToken"

	// Создание токена
	_, tokenString, err := JWTAuth.Encode(map[string]interface{}{
		"login": login,
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // Время истечения токена 24 часа
	})
	if err != nil {
		return "", fmt.Errorf("%w: %s", err, op)
	}

	return tokenString, nil
}
