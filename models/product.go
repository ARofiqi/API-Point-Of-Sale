package models

import "github.com/go-playground/validator/v10"

type Product struct {
	ID         int     `json:"id"`
	Name       string  `json:"name" validate:"required"`
	Price      float64 `json:"price" validate:"required,gt=0"`
	CategoryID int    `json:"category_id" validate:"required"`
	Category   string  `json:"category,omitempty"` // Opsional, hanya untuk output
}

var Validate = validator.New()
