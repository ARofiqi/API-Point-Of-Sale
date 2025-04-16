package dto

import (
	"aro-shop/models"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ProductResponse struct {
	ID        uuid.UUID       `json:"id"`
	Name      string          `json:"name"`
	Price     float64         `json:"price"`
	Stock     int             `json:"stock"`
	Category  models.Category `json:"category"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type ProductRequest struct {
	Name       string    `json:"name" validate:"required"`
	Price      float64   `json:"price" validate:"required,gt=0"`
	CategoryID uuid.UUID `json:"category_id" validate:"required"`
	Stock      int       `json:"stock" validate:"required,gte=0"`
}

func ConvertToProductResponse(product models.Product) ProductResponse {
	return ProductResponse{
		ID:        product.ID,
		Name:      product.Name,
		Price:     product.Price,
		Stock:     product.Stock,
		Category:  product.Category,
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}
}

var Validate = validator.New()
