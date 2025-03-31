package model

import "gorm.io/gorm"


type User struct {
    gorm.Model
    Username string `gorm:"unique;not null" validate:"required, min=5, max=20"`
    Email    string `gorm:"unique;not null" validate:"required, email"`
    Password string `gorm:"not null" validate:"required, min=7,max=20"`
}