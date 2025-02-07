package models

type Transaction struct {
	ID    int               `json:"id"`
	Date  string            `json:"date"`
	Total float64           `json:"total"`
	Items []TransactionItem `json:"items"`
}

type TransactionItem struct {
	ID            int     `json:"id"`
	TransactionID int     `json:"transaction_id"`
	ProductID     int     `json:"product_id" validate:"required"`
	Quantity      int     `json:"quantity" validate:"required,min=1"`
	SubTotal      float64 `json:"sub_total"`
}
