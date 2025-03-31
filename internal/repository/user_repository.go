package repository

import (
	"errors"
	"fmt"
	"project/internal/model"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)


type UserRepoI interface {
	CreateUser(string, string, string) error
}

type UserReposiotry struct {
	db *gorm.DB
}


func NewUserRepository(db *gorm.DB) UserRepoI{
	return &UserReposiotry{
		db: db,
	}
}


func (u *UserReposiotry) CreateUser(username string, email string, password string) error{
	user := &model.User{
        Username: username,
        Email:    email,
        Password: password,
    }

	validate := validator.New()
	err := validate.Struct(user)
	if err != nil{
		return fmt.Errorf("validation failed %s", err)
	}

	var existingUser model.User
	err = u.db.Where("email = ?", email).First(&existingUser).Error
	if err == nil {
		return fmt.Errorf("user already exists %s", email)
	}else if !errors.Is(err, gorm.ErrRecordNotFound){
		return fmt.Errorf("failed to check existing user: %w", err)
	}

	err = u.db.Create(user).Error
	if err != nil {
		return fmt.Errorf("cannot create a user %s", err)
	}

	fmt.Println("User was created successful")
	return nil
}