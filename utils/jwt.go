package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)


func GenerateAccessToken(userID string) (string, error) {
    // Определяем время истечения токена (например, 15 минут)
    expirationTime := time.Now().Add(30 * time.Minute)

    // Создаем claims
    claims := jwt.MapClaims{
        "user_id": userID ,           // Полезная нагрузка (например, ID пользователя)
        "exp":      expirationTime.Unix(), // Время истечения
    }

    // Создаем новый токен с указанными claims
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    // Подписываем токен секретным ключом
    return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func GenerateRefreshToken() string{
	return uuid.New().String()
}