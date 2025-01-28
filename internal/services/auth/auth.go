package auth

import (
	"fmt"
	"github.com/go-chi/jwtauth/v5"
	"os"
	"time"
)

var JWTAuth *jwtauth.JWTAuth

func InitJWTAuth(secretKey []byte) (*jwtauth.JWTAuth, error) {
	JWTAuth = jwtauth.New("HS256", secretKey, nil)
	return JWTAuth, nil
}

func GenerateToken(login string) (string, error) {
	const op = "service.GenerateToken"

	expStr := os.Getenv("JWT_EXPIRATION")
	expDur, err := time.ParseDuration(expStr)
	if err != nil {
		fmt.Println("Error parsing duration:", err)
		return "", nil
	}

	// Создание токена
	_, tokenString, err := JWTAuth.Encode(map[string]interface{}{
		"login": login,
		"exp":   time.Now().Add(expDur).Unix(),
	})
	if err != nil {
		return "", fmt.Errorf("%w: %s", err, op)
	}

	return tokenString, nil
}
