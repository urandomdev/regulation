package models

import (
	"time"

	"github.com/google/uuid"
)

// Transaction represents a financial transaction in API responses
type Transaction struct {
	ID             uuid.UUID `cbor:"id" json:"id"`
	AccountID      uuid.UUID `cbor:"account_id" json:"account_id"`
	Amount         int64     `cbor:"amount" json:"amount"`
	Date           time.Time `cbor:"date" json:"date"`
	Name           string    `cbor:"name" json:"name"`
	MerchantName   *string   `cbor:"merchant_name,omitempty" json:"merchant_name,omitempty"`
	Category       string    `cbor:"category" json:"category"`
	Pending        bool      `cbor:"pending" json:"pending"`
	PaymentChannel *string   `cbor:"payment_channel,omitempty" json:"payment_channel,omitempty"`
}
