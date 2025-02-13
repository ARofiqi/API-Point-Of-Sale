package models

import (
	"github.com/go-playground/validator/v10"
)

type Category struct {
	ID       uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name     string    `json:"name" gorm:"type:char(255);not null"`
	Products []Product `json:"products,omitempty" gorm:"foreignKey:CategoryID"`
}

type Product struct {
	ID         uint     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name       string   `json:"name" validate:"required" gorm:"type:varchar(255);not null"`
	Price      float64  `json:"price" validate:"required,gt=0" gorm:"type:decimal(10,2);not null"`
	CategoryID uint     `json:"category_id" validate:"required" gorm:"not null"`
	Category   Category `json:"category,omitempty" gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

var Validate = validator.New()
