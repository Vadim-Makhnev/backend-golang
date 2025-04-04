package middleware

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Проверяем, что JWT_SECRET установлен
        secretKey := os.Getenv("JWT_SECRET")
        if secretKey == "" {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error": "JWT_SECRET is not set",
            })
        }

        // Извлекаем access_token из cookie
        accessToken := c.Cookies("access_token")
        if accessToken == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Missing access token in cookie",
            })
        }

        // Парсим токен
        token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
            // Проверяем метод подписи
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            return []byte(secretKey), nil
        })

        // Обрабатываем ошибки парсинга токена
        if err != nil {
            if errors.Is(err, jwt.ErrTokenMalformed) {
                return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                    "error": "Token is malformed",
                })
            } else if errors.Is(err, jwt.ErrTokenExpired) {
                return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                    "error": "Token has expired",
                })
            } else if errors.Is(err, jwt.ErrTokenNotValidYet) {
                return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                    "error": "Token is not valid yet",
                })
            } else {
                return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                    "error": "Invalid token",
                })
            }
        }

        // Проверяем валидность токена
        if !token.Valid {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Invalid token",
            })
        }

        // Извлекаем claims
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Invalid token claims",
            })
        }

        // Проверяем, истёк ли токен
        exp, expExists := claims["exp"].(float64)
        if !expExists {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Expiration time not found in token",
            })
        }

        if time.Now().Unix() > int64(exp) {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Access token has expired",
            })
        }

        // Извлекаем user ID из токена
        userID, ok := claims["user_id"].(string)
        if !ok {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "User ID not found or invalid in token",
            })
        }

        c.Locals("user_id", userID) // Сохраняем user ID в контексте

        return c.Next()
    }
}