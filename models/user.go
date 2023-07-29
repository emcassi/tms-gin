package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:"not null" json:"name"`
	Email    string `gorm:"not null;unique" json:"email"`
	Password string `gorm:"not null" json:"password"`
	Avatar   string `json:"avatar"`
}
