package model

import (
	"gorm.io/gorm"
)


type User struct {
    gorm.Model
    Username string `gorm:"unique;not null"`
    Email    string `gorm:"unique;not null"`
    Password string `gorm:"not null"`
    Role string `gorm:"default:'user'"`
    IsSubscribe bool `gorm:"default:false"`

     Token Token `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}