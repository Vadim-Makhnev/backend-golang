package service

import (
	"errors"
	"project/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	service repository.UserRepoI
}

func NewAuthService(service repository.UserRepoI) *AuthService{
	return &AuthService{
		service: service,
	}
}


func (a *AuthService) CreateUser(username string, email string, password string) error{

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return errors.New("cannot hash password")
	}
	err = a.service.CreateUser(username, email , string(hashedPassword))
	if err != nil{
		return errors.New("cannot create a user")
	}
	return nil
}