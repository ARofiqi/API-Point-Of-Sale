package models

import "github.com/go-playground/validator/v10"

type Product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name" validate:"required"`
	Price    float64 `json:"price" validate:"required,gt=0"`
	Category string  `json:"category" validate:"required"`
}

var Validate = validator.New()
