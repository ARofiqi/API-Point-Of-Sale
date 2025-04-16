package dto

import (
	"time"

	"github.com/google/uuid"
)

type TransactionItemRequest struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required,min=1"`
}

type TransactionRequest struct {
	Items           []TransactionItemRequest `json:"items" validate:"required,dive"`
	PaymentMethodID uuid.UUID                `json:"payment_method_id" validate:"required"`
	AmountPaid      float64                  `json:"amount_paid" validate:"required,gt=0"`
}

type TransactionResponse struct {
	ID         uuid.UUID                 `json:"id"`
	User       SimpleUserResponse        `json:"user"`
	Date       time.Time                 `json:"date"`
	// AmountPaid float64                   `json:"amount_paid"`
	Items      []TransactionItemResponse `json:"items,omitempty"`
	Payment    *PaymentResponse          `json:"payment,omitempty"`
	CreatedAt  time.Time                 `json:"created_at"`
	UpdatedAt  time.Time                 `json:"updated_at"`
}

type SimpleUserResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

type TransactionItemResponse struct {
	ID          uuid.UUID `json:"id"`
	ProductID   uuid.UUID `json:"product_id"`
	ProductName string    `json:"product_name"`
	Quantity    int       `json:"quantity"`
	SubTotal    float64   `json:"subtotal"`
}

type PaymentResponse struct {
	ID            uuid.UUID            `json:"id"`
	Method        string               `json:"method"`
	Status        string               `json:"status"`
	PaidAt        time.Time            `json:"paid_at"`
	AmountPaid    float64              `json:"amount_paid"`
	PaymentMethod *PaymentMethodSimple `json:"payment_method,omitempty"`
	CreatedAt     time.Time            `json:"created_at"`
	UpdatedAt     time.Time            `json:"updated_at"`
}

type PaymentMethodSimple struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
