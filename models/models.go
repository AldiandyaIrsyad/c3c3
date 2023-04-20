package models

import "github.com/jinzhu/gorm"

type Role string

const (
	Admin   Role = "admin"
	Creator Role = "creator"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Role     Role   `gorm:"not null"`
}

type Product struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string
	Price       float64 `gorm:"not null"`
	UserID      uint    `gorm:"not null"`
}
