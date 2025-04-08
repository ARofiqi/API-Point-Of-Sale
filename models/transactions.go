package models

import (
	"time"
)

type PaymentMethod struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"type:varchar(50);not null;unique"`
}

type Payment struct {
	ID              uint          `json:"id" gorm:"primaryKey"`
	TransactionID   uint          `json:"transaction_id" gorm:"unique;not null"`
	Transaction     *Transaction  `json:"transaction" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	PaymentMethodID uint          `json:"payment_method_id" gorm:"not null"`
	PaymentMethod   PaymentMethod `json:"payment_method" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	AmountPaid      float64       `json:"amount_paid" gorm:"type:numeric(10,2);not null"`
	PaymentStatus   string        `json:"payment_status" gorm:"type:varchar(20);not null;default:'pending'"`
	PaidAt          *time.Time    `json:"paid_at"`
}

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
