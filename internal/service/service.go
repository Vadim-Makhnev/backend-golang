package service

import (
	"fmt"
	"project/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo repository.UserRepoI
}

func NewAuthService(service repository.UserRepoI) *AuthService{
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

func (a *AuthService) LoginUser(email string, password string) error{
	user, err := a.repo.Login(email, password)
	if err != nil {
		return fmt.Errorf("failed to login user: %w", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return fmt.Errorf("invalid credentials: %w", err)
	}
	return nil
}