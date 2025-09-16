package models

type User struct {
	Base
	Name     string `json:"name" gorm:"not null"`
	Email    string `json:"email" gorm:"unique;not null"`
	Password string `json:"-" gorm:"not null"`
	Role     string `json:"role" gorm:"not null"`
	IsActive bool   `json:"is_active"`
}
