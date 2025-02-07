package models

type Transaction struct {
	ID        int     `json:"id"`
	ProductID int     `json:"product_id" validate:"required,gt=0"`
	Quantity  int     `json:"quantity" validate:"required,gt=0"`
	Total     float64 `json:"total"`
}
