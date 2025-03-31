package service

import (
	"errors"
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
		return errors.New("cannot hash password")
	}
	err = a.repo.CreateUser(username, email , string(hashedPassword))
	if err != nil{
		return errors.New("cannot create a user")
	}
	return nil
}

func (a *AuthService) LoginUser(email string, password string) error{
	user, err := a.repo.Login(email, password)
	if err != nil {
		return fmt.Errorf("cannot login %s", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return fmt.Errorf("cannot to compare password %s", err)
	}
	return nil
}