package model

import (
	"time"

	"gorm.io/gorm"
)


type Token struct {
	gorm.Model
	Token string `gorm:"required,not null"`
	ExpiresIn time.Time `gorm:"required, not null"`

	UserID uint `gorm:"unique;not null"`
}