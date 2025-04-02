package repository

import (
	"errors"
	"fmt"
	"project/internal/model"
	"time"

	"gorm.io/gorm"
)


type UserRepoI interface {
	CreateUser(username string,email string, password string) error
	Login(email string, password string ) (*model.User, error)
	AddToken(token string, UserId uint) error
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

	var existingUser model.User
	err := u.db.Where("email = ?", email).First(&existingUser).Error
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

func (u *UserReposiotry) Login(email string, password string) (*model.User ,error){
	var existingUser model.User
	err := u.db.Where("email = ?", email).First(&existingUser).Error
	if err != nil{
		return nil, fmt.Errorf("user is not found")
	}
		return &existingUser, nil
}

func (u *UserReposiotry) AddToken(token string, UserId uint) error{
	var existingToken model.Token
	err := u.db.Where("user_id = ?", UserId).First(&existingToken).Error
	if err == nil {
		if err := u.db.Delete(&existingToken).Error; err != nil{
			return fmt.Errorf("failed to delete token")
		}
	}

	newToken := model.Token{
        Token:     token,
        ExpiresIn: time.Now().Add(7 * 24 * time.Hour), // Устанавливаем срок действия токена
        UserID:    UserId,
    }


	if err := u.db.Create(&newToken).Error; err != nil {
        return fmt.Errorf("failed add token to the databse: %v", err)
        
    }
	return nil
}