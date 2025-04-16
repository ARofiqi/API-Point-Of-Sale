package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Product struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name       string    `json:"name" validate:"required" gorm:"type:varchar(255);not null"`
	Price      float64   `json:"price" validate:"required,gt=0" gorm:"type:numeric(10,2);not null"`
	Stock      int       `json:"stock" validate:"required,gte=0" gorm:"not null"`
	CategoryID uuid.UUID `json:"category_id" gorm:"type:uuid;not null;index"`
	Category   Category  `json:"category" gorm:"foreignKey:CategoryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

var Validate = validator.New()
