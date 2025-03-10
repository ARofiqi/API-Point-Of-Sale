package models

import (
	"time"
)

type Transaction struct {
	ID    uint              `json:"id" gorm:"primaryKey"`
	Date  time.Time         `json:"date" gorm:"autoCreateTime"`
	Total float64           `json:"total" gorm:"type:numeric(10,2);not null"`
	Items []TransactionItem `json:"items" gorm:"foreignKey:TransactionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type TransactionItem struct {
	ID            uint    `json:"id" gorm:"primaryKey"`
	TransactionID uint    `json:"transaction_id" gorm:"index;not null"`
	ProductID     uint    `json:"product_id" validate:"required" gorm:"not null"`
	Quantity      int     `json:"quantity" validate:"required,min=1" gorm:"not null"`
	SubTotal      float64 `json:"sub_total" gorm:"type:numeric(10,2);not null"`
}
