package plaid

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/plaid/plaid-go/v35/plaid"

	"regulation/internal/config"
)

// plaidClientImpl is the real implementation using Plaid SDK
type plaidClientImpl struct {
	client *plaid.APIClient
	config *config.PlaidConfig
}

// NewPlaidClient creates a new Plaid client using the official SDK
func NewPlaidClient(cfg *config.PlaidConfig) (PlaidClient, error) {
	if cfg == nil {
		return nil, errors.New("plaid config is nil")
	}

	configuration := plaid.NewConfiguration()
	configuration.AddDefaultHeader("PLAID-CLIENT-ID", cfg.GetClientID())
	configuration.AddDefaultHeader("PLAID-SECRET", cfg.GetSecret())

	// Set environment
	env := cfg.GetEnvironment()
	switch env {
	case "sandbox":
		configuration.UseEnvironment(plaid.Sandbox)
	case "development":
		configuration.UseEnvironment(plaid.Sandbox) // Use sandbox for development
	case "production":
		configuration.UseEnvironment(plaid.Production)
	default:
		return nil, fmt.Errorf("invalid plaid environment: %s", env)
	}

	return &plaidClientImpl{
		client: plaid.NewAPIClient(configuration),
		config: cfg,
	}, nil
}

// CreateLinkToken creates a link token for initializing Plaid Link
func (p *plaidClientImpl) CreateLinkToken(ctx context.Context, userID string) (string, error) {
	user := plaid.LinkTokenCreateRequestUser{
		ClientUserId: userID,
	}

	request := plaid.NewLinkTokenCreateRequest(
		"Who Are You To Spend",
		"en",
		[]plaid.CountryCode{plaid.COUNTRYCODE_US},
		user,
	)

	// Request transactions product for transaction data
	request.SetProducts([]plaid.Products{plaid.PRODUCTS_TRANSACTIONS})

	// Filter to only checking and savings accounts
	request.SetAccountFilters(plaid.LinkTokenAccountFilters{
		Depository: &plaid.DepositoryFilter{
			AccountSubtypes: []plaid.DepositoryAccountSubtype{
				plaid.DEPOSITORYACCOUNTSUBTYPE_CHECKING,
				plaid.DEPOSITORYACCOUNTSUBTYPE_SAVINGS,
			},
		},
	})

	resp, _, err := p.client.PlaidApi.LinkTokenCreate(ctx).
		LinkTokenCreateRequest(*request).
		Execute()
	if err != nil {
		return "", fmt.Errorf("failed to create link token: %w", err)
	}

	return resp.GetLinkToken(), nil
}

// ExchangePublicToken exchanges a public token for an access token
func (p *plaidClientImpl) ExchangePublicToken(ctx context.Context, publicToken string) (*TokenExchangeResult, error) {
	request := plaid.NewItemPublicTokenExchangeRequest(publicToken)

	resp, _, err := p.client.PlaidApi.ItemPublicTokenExchange(ctx).
		ItemPublicTokenExchangeRequest(*request).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to exchange public token: %w", err)
	}

	return &TokenExchangeResult{
		AccessToken: resp.GetAccessToken(),
		ItemID:      resp.GetItemId(),
	}, nil
}

// GetAccounts retrieves all accounts for an access token
func (p *plaidClientImpl) GetAccounts(ctx context.Context, accessToken string) ([]Account, error) {
	request := plaid.NewAccountsGetRequest(accessToken)

	resp, _, err := p.client.PlaidApi.AccountsGet(ctx).
		AccountsGetRequest(*request).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get accounts: %w", err)
	}

	accounts := make([]Account, 0, len(resp.GetAccounts()))
	for _, acc := range resp.GetAccounts() {
		balances := acc.GetBalances()

		var balanceAvailable *int64
		if available, ok := balances.GetAvailableOk(); ok && available != nil {
			cents := int64(*available * 100)
			balanceAvailable = &cents
		}

		accounts = append(accounts, Account{
			AccountID:        acc.GetAccountId(),
			Name:             acc.GetName(),
			Type:             string(acc.GetType()),
			Subtype:          string(acc.GetSubtype()),
			Mask:             acc.GetMask(),
			BalanceCurrent:   int64(balances.GetCurrent() * 100), // convert to cents
			BalanceAvailable: balanceAvailable,
		})
	}

	return accounts, nil
}

// SyncTransactions retrieves transactions using the sync endpoint
func (p *plaidClientImpl) SyncTransactions(ctx context.Context, accessToken, cursor string) (*TransactionSyncResult, error) {
	request := plaid.NewTransactionsSyncRequest(accessToken)
	if cursor != "" {
		request.SetCursor(cursor)
	}

	resp, _, err := p.client.PlaidApi.TransactionsSync(ctx).
		TransactionsSyncRequest(*request).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to sync transactions: %w", err)
	}

	result := &TransactionSyncResult{
		Added:    make([]Transaction, 0, len(resp.GetAdded())),
		Modified: make([]Transaction, 0, len(resp.GetModified())),
		Removed:  make([]string, 0, len(resp.GetRemoved())),
		Cursor:   resp.GetNextCursor(),
		HasMore:  resp.GetHasMore(),
	}

	// Process added transactions
	for _, tx := range resp.GetAdded() {
		result.Added = append(result.Added, p.convertTransaction(tx))
	}

	// Process modified transactions
	for _, tx := range resp.GetModified() {
		result.Modified = append(result.Modified, p.convertTransaction(tx))
	}

	// Process removed transactions
	for _, removed := range resp.GetRemoved() {
		result.Removed = append(result.Removed, removed.GetTransactionId())
	}

	return result, nil
}

// convertTransaction converts a Plaid transaction to our internal Transaction type
func (p *plaidClientImpl) convertTransaction(tx plaid.Transaction) Transaction {
	// Parse date
	date, _ := time.Parse("2006-01-02", tx.GetDate())

	// Get merchant name if available
	merchantName := ""
	if merchant, ok := tx.GetMerchantNameOk(); ok && merchant != nil {
		merchantName = *merchant
	}

	// Get payment channel if available
	paymentChannel := ""
	if channel, ok := tx.GetPaymentChannelOk(); ok && channel != nil {
		paymentChannel = *channel
	}

	return Transaction{
		TransactionID:  tx.GetTransactionId(),
		AccountID:      tx.GetAccountId(),
		Amount:         int64(tx.GetAmount() * 100), // convert to cents
		Date:           date,
		Name:           tx.GetName(),
		MerchantName:   merchantName,
		Categories:     tx.GetCategory(),
		Pending:        tx.GetPending(),
		PaymentChannel: paymentChannel,
	}
}
