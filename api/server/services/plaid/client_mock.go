package plaid

import (
	"context"
)

// MockPlaidClient is a mock implementation of PlaidClient for testing
type MockPlaidClient struct {
	CreateLinkTokenFn     func(ctx context.Context, userID string) (string, error)
	ExchangePublicTokenFn func(ctx context.Context, publicToken string) (*TokenExchangeResult, error)
	GetAccountsFn         func(ctx context.Context, accessToken string) ([]Account, error)
	SyncTransactionsFn    func(ctx context.Context, accessToken, cursor string) (*TransactionSyncResult, error)
}

// NewMockPlaidClient creates a new mock Plaid client with default behaviors
func NewMockPlaidClient() *MockPlaidClient {
	return &MockPlaidClient{
		CreateLinkTokenFn: func(_ context.Context, userID string) (string, error) {
			return "link-sandbox-test-token-" + userID, nil
		},
		ExchangePublicTokenFn: func(_ context.Context, _ string) (*TokenExchangeResult, error) {
			return &TokenExchangeResult{
				AccessToken: "access-sandbox-test-token",
				ItemID:      "item-sandbox-test-id",
			}, nil
		},
		GetAccountsFn: func(_ context.Context, _ string) ([]Account, error) {
			available := int64(150000) // $1500.00
			return []Account{
				{
					AccountID:        "account-checking-123",
					Name:             "Plaid Checking",
					Type:             "depository",
					Subtype:          "checking",
					Mask:             "0000",
					BalanceCurrent:   200000, // $2000.00 in cents
					BalanceAvailable: &available,
				},
				{
					AccountID:        "account-savings-456",
					Name:             "Plaid Savings",
					Type:             "depository",
					Subtype:          "savings",
					Mask:             "1111",
					BalanceCurrent:   500000, // $5000.00 in cents
					BalanceAvailable: nil,
				},
			}, nil
		},
		SyncTransactionsFn: func(_ context.Context, _ string, cursor string) (*TransactionSyncResult, error) {
			// Return empty result if cursor is provided (simulating no new transactions)
			if cursor != "" {
				return &TransactionSyncResult{
					Added:    []Transaction{},
					Modified: []Transaction{},
					Removed:  []string{},
					Cursor:   cursor,
					HasMore:  false,
				}, nil
			}

			// Return sample transactions on first sync
			return &TransactionSyncResult{
				Added: []Transaction{
					{
						TransactionID:  "tx-1",
						AccountID:      "account-checking-123",
						Amount:         2500, // $25.00
						Name:           "Starbucks",
						MerchantName:   "Starbucks",
						Categories:     []string{"Food and Drink", "Restaurants", "Coffee Shop"},
						Pending:        false,
						PaymentChannel: "in store",
					},
					{
						TransactionID:  "tx-2",
						AccountID:      "account-checking-123",
						Amount:         5000, // $50.00
						Name:           "Amazon",
						MerchantName:   "Amazon",
						Categories:     []string{"Shops", "Digital Purchase"},
						Pending:        false,
						PaymentChannel: "online",
					},
				},
				Modified: []Transaction{},
				Removed:  []string{},
				Cursor:   "mock-cursor-1",
				HasMore:  false,
			}, nil
		},
	}
}

// CreateLinkToken creates a mock link token
func (m *MockPlaidClient) CreateLinkToken(ctx context.Context, userID string) (string, error) {
	if m.CreateLinkTokenFn != nil {
		return m.CreateLinkTokenFn(ctx, userID)
	}
	return "link-sandbox-test-token", nil
}

// ExchangePublicToken exchanges a mock public token
func (m *MockPlaidClient) ExchangePublicToken(ctx context.Context, publicToken string) (*TokenExchangeResult, error) {
	if m.ExchangePublicTokenFn != nil {
		return m.ExchangePublicTokenFn(ctx, publicToken)
	}
	return &TokenExchangeResult{
		AccessToken: "access-sandbox-test-token",
		ItemID:      "item-sandbox-test-id",
	}, nil
}

// GetAccounts retrieves mock accounts
func (m *MockPlaidClient) GetAccounts(ctx context.Context, accessToken string) ([]Account, error) {
	if m.GetAccountsFn != nil {
		return m.GetAccountsFn(ctx, accessToken)
	}
	return []Account{}, nil
}

// SyncTransactions retrieves mock transactions
func (m *MockPlaidClient) SyncTransactions(ctx context.Context, accessToken, cursor string) (*TransactionSyncResult, error) {
	if m.SyncTransactionsFn != nil {
		return m.SyncTransactionsFn(ctx, accessToken, cursor)
	}
	return &TransactionSyncResult{
		Added:    []Transaction{},
		Modified: []Transaction{},
		Removed:  []string{},
		Cursor:   "mock-cursor",
		HasMore:  false,
	}, nil
}
