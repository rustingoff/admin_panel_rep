package model

import (
	"gorm.io/gorm"
	"time"
)

type Client struct {
	gorm.Model

	NrContract string `json:"nr_contract" validate:"required"`
	FirstName  string `json:"first_name" validate:"required"`
	LastName   string `json:"last_name" validate:"required"`

	IDNP uint64 `json:"idnp" gorm:"Column:idnp" validate:"required,min=13,max=13"`

	Phone string `json:"phone" validate:"required,min=10,max=13"`

	Sum  uint64 `json:"sum" validate:"required,min=0"`
	Time uint64 `json:"time" validate:"required,min=0"`

	SignDate time.Time `json:"sign_date" validate:"required"`

	MonthlyRate float64 `json:"monthly_rate" validate:"required,min=0"`
}

type ClientUpdate struct {
	NrContract string `json:"nr_contract"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`

	IDNP uint64 `json:"idnp" gorm:"Column:idnp" validate:"min=13,max=13"`

	Phone string `json:"phone" validate:"min=10,max=13"`

	Sum  uint64 `json:"sum" validate:"required,min=0"`
	Time uint64 `json:"time" validate:"required,min=0"`

	SignDate    time.Time `json:"sign_date" validate:"required"`
	MonthlyRate float64   `json:"monthly_rate" validate:"required,min=0"`
}
