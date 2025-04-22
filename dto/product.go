package dto

import (
	"aro-shop/models"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ProductResponse struct {
	ID          uuid.UUID       `json:"id"`
	Name        string          `json:"name"`
	Price       float64         `json:"price"`
	Description string          `json:"description"`
	Stock       int             `json:"stock"`
	URLImage    string          `json:"url_image"`
	Category    models.Category `json:"category"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type ProductRequest struct {
	Name        string    `json:"name" validate:"required"`
	Price       float64   `json:"price" validate:"required,gt=0"`
	Description string    `json:"description"`
	Stock       int       `json:"stock" validate:"required,gte=0"`
	URLImage    string    `json:"url_image"`
	CategoryID  uuid.UUID `json:"category_id" validate:"required"`
}

func ConvertToProductResponse(product models.Product) ProductResponse {
	return ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		URLImage:    product.URLImage,
		Category:    product.Category,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}

var Validate = validator.New()
