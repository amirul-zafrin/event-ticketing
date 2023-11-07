package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
	Admin     bool   `json:"admin"`
}
