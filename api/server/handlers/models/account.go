package models

import "github.com/google/uuid"

// Account represents a bank account in API responses
type Account struct {
	ID               uuid.UUID `cbor:"id" json:"id"`
	Name             string    `cbor:"name" json:"name"`
	Type             string    `cbor:"type" json:"type"`
	Mask             *string   `cbor:"mask,omitempty" json:"mask,omitempty"`
	CurrentBalance   int64     `cbor:"current_balance" json:"current_balance"`
	AvailableBalance *int64    `cbor:"available_balance,omitempty" json:"available_balance,omitempty"`
	IsActive         bool      `cbor:"is_active" json:"is_active"`
}
