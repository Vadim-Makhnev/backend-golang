package service

import (
	"fmt"
	"project/internal/repository"
	"project/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceI interface{
	CreateUser(username string, email string, password string) error
	LoginUser(c *fiber.Ctx, email string, password string) error
}

type AuthService struct {
	repo repository.UserRepoI
}

func NewAuthService(service repository.UserRepoI) AuthServiceI{
	return &AuthService{
		repo: service,
	}
}


func (a *AuthService) CreateUser(username string, email string, password string) error{

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	err = a.repo.CreateUser(username, email , string(hashedPassword))
	if err != nil{
		fmt.Printf("Error creating user: %v\n", err)
        return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (a *AuthService) LoginUser(c* fiber.Ctx, email string, password string) error{
	user, err := a.repo.Login(email, password)
	if err != nil {
		return fmt.Errorf("failed to login user: %w", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return fmt.Errorf("invalid credentials: %w", err)
	}
	userID1 := fmt.Sprintf("%d", user.ID)
	sign , err := utils.GenerateAccessToken(userID1)
	if err != nil {
		return fmt.Errorf("failed to generate jwt token")
	}

	setAccessTokenCookie(c , sign)
	utils.GenerateRefreshToken()
	return nil
}

func setAccessTokenCookie(c *fiber.Ctx, token string) {
    cookie := fiber.Cookie{
        Name:     "access_token",         // Имя cookie
        Value:    token,                  // Значение токена
        Expires:  time.Now().Add(30 * time.Minute), // Время жизни cookie
        HTTPOnly: true,                   // Защита от доступа через JavaScript
        Secure:   false,                   // Отправка только по HTTPS
        SameSite: "Strict",               // Защита от CSRF
        Path:     "/",                    // Доступность на всех путях
    }
    c.Cookie(&cookie) // Устанавливаем cookie
}