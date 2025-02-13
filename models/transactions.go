package models

import (
	"time"
)

type Transaction struct {
	ID    uint              `gorm:"primaryKey" json:"id"`
	Date  time.Time         `json:"date" gorm:"autoCreateTime"`
	Total float64           `json:"total"`
	Items []TransactionItem `json:"items" gorm:"foreignKey:TransactionID"`
}

type TransactionItem struct {
	ID            uint    `gorm:"primaryKey" json:"id"`
	TransactionID uint    `json:"transaction_id" gorm:"index"`
	ProductID     uint    `json:"product_id" validate:"required"`
	Quantity      int     `json:"quantity" validate:"required,min=1"`
	SubTotal      float64 `json:"sub_total"`
}
