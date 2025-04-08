package models

import (
	"github.com/go-playground/validator/v10"
)

type Category struct {
	ID   uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name string `json:"name" gorm:"type:varchar(255);not null"`
}

type Product struct {
	ID         uint     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name       string   `json:"name" validate:"required" gorm:"type:varchar(255);not null"`
	Price      float64  `json:"price" validate:"required,gt=0" gorm:"type:numeric(10,2);not null"`
	CategoryID uint     `gorm:"not null;index"`
	Category   Category `json:"category" gorm:"foreignKey:CategoryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type ProductResponse struct {
	ID       uint     `json:"id"`
	Name     string   `json:"name"`
	Price    float64  `json:"price"`
	Category Category `json:"category"`
}

type ProductRequest struct {
	Name       string  `json:"name" validate:"required"`
	Price      float64 `json:"price" validate:"required,gt=0"`
	CategoryID uint    `json:"category_id" validate:"required"`
}

func ConvertToProductResponse(product Product) ProductResponse {
	return ProductResponse{
		ID:       product.ID,
		Name:     product.Name,
		Price:    product.Price,
		Category: product.Category,
	}
}

var Validate = validator.New()
