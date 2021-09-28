package model

import (
	_ "gopkg.in/go-playground/validator.v9"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	FirstName string `json:"first_name" validate:"required,min=2,max=32"`
	LastName  string `json:"last_name" validate:"required,min=2,max=32"`

	Username string `json:"username" validate:"required,min=4,max=20"`
	Password string `json:"password" validate:"required,min=8,max=64"`

	Email string `json:"email" validate:"min=0,max=32,email"`
}
