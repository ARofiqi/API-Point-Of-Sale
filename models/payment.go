package models

import (
	"time"

	"github.com/google/uuid"
)

type PaymentMethod struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string    `json:"name" validate:"required" gorm:"unique;type:varchar(50);not null;unique"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Payment struct {
	ID              uuid.UUID     `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TransactionID   uuid.UUID     `json:"transaction_id" gorm:"type:uuid;not null;unique"`
	Transaction     Transaction   `json:"transaction" gorm:"foreignKey:TransactionID"`
	PaymentMethodID uuid.UUID     `json:"payment_method_id" gorm:"not null"`
	AmountPaid      float64       `json:"amount_paid" gorm:"type:numeric(10,2);not null"`
	PaidAt          time.Time     `json:"paid_at" gorm:"autoCreateTime"`
	PaymentStatus   string        `json:"payment_status" gorm:"type:varchar(20);default:'pending';not null"`
	PaymentMethod   PaymentMethod `json:"payment_method" gorm:"foreignKey:PaymentMethodID"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
}
