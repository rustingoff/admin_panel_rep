package model

import (
	_ "gopkg.in/go-playground/validator.v9"
	"gorm.io/gorm"
)

type Client struct {
	gorm.Model

	NrContract string `json:"nr_contract" validate:"required"`
	FirstName  string `json:"first_name" validate:"required"`
	LastName   string `json:"last_name" validate:"required"`

	IDNP uint64 `json:"idnp" gorm:"Column:idnp" validate:"required,min=13"`

	Phone string `json:"phone" validate:"required,min=9,max=13"`

	Sum  uint64 `json:"sum" validate:"required,min=0"`
	Time uint64 `json:"time" validate:"required,min=0"`

	SignDate string `json:"sign_date" validate:"required"`

	MonthlyRate string `json:"monthly_rate" validate:"required,min=0"`
}

type ClientUpdate struct {
	NrContract string `json:"nr_contract"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`

	IDNP uint64 `json:"idnp" gorm:"Column:idnp"`

	Phone string `json:"phone"`

	Sum  uint64 `json:"sum"`
	Time uint64 `json:"time"`

	SignDate string `json:"sign_date"`
}
