package plaid

import (
	"context"
	"time"
)

// Client defines the interface for interacting with Plaid API.
// This interface allows for easy mocking in tests.
type Client interface {
	// CreateLinkToken creates a link token for initializing Plaid Link
	CreateLinkToken(ctx context.Context, userID string) (string, error)

	// ExchangePublicToken exchanges a public token for an access token
	ExchangePublicToken(ctx context.Context, publicToken string) (*TokenExchangeResult, error)

	// GetAccounts retrieves all accounts for an access token
	GetAccounts(ctx context.Context, accessToken string) ([]Account, error)

	// SyncTransactions retrieves transactions using the sync endpoint
	SyncTransactions(ctx context.Context, accessToken, cursor string) (*TransactionSyncResult, error)
}

// TokenExchangeResult contains the result of exchanging a public token
type TokenExchangeResult struct {
	AccessToken string
	ItemID      string
}

// Account represents a bank account from Plaid
type Account struct {
	AccountID        string
	Name             string
	Type             string
	Subtype          string
	Mask             string
	BalanceCurrent   int64  // in cents
	BalanceAvailable *int64 // in cents, may be nil
}

// Transaction represents a transaction from Plaid
type Transaction struct {
	TransactionID  string
	AccountID      string
	Amount         int64 // in cents (positive = debit, negative = credit)
	Date           time.Time
	Name           string
	MerchantName   string
	Categories     []string
	Pending        bool
	PaymentChannel string
}

// TransactionSyncResult contains the result of syncing transactions
type TransactionSyncResult struct {
	Added    []Transaction
	Modified []Transaction
	Removed  []string // transaction IDs that were removed
	Cursor   string
	HasMore  bool
}
