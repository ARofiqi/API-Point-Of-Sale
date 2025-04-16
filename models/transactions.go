package models

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID         uuid.UUID         `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID     uuid.UUID         `json:"user_id" gorm:"type:uuid;not null"`
	User       User              `json:"user" gorm:"foreignKey:UserID;references:ID"`
	Date       time.Time         `json:"date" gorm:"autoCreateTime"`
	// AmountPaid float64           `json:"amount_paid" gorm:"type:numeric(10,2);not null"`
	Items      []TransactionItem `json:"items" gorm:"foreignKey:TransactionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Payment    *Payment          `json:"payment,omitempty" gorm:"foreignKey:TransactionID"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
}

type TransactionItem struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TransactionID uuid.UUID `json:"transaction_id" gorm:"type:uuid;not null"`
	ProductID     uuid.UUID `json:"product_id" gorm:"type:uuid;not null"`
	Product       Product   `json:"product" gorm:"foreignKey:ProductID"`
	Quantity      int       `json:"quantity" gorm:"not null"`
	SubTotal      float64   `json:"subtotal" gorm:"type:numeric(10,2);not null"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
