package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/plaid/plaid-go/v35/plaid"

	"regulation/internal/config"
)

// loadConfigFromFile loads config from a specific file path
func loadConfigFromFile(path string) (*config.Config, error) {
	configFile, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg config.Config
	if err = json.Unmarshal(configFile, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	if cfg.Plaid == nil {
		return nil, fmt.Errorf("plaid config is missing")
	}

	return &cfg, nil
}

// Creates a Plaid sandbox account simulating heavy overspending with transaction history
func main() {
	// Determine config file path
	configPath := filepath.Join("..", "config", "config.json")
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	// Load config from file
	cfg, err := loadConfigFromFile(configPath)
	if err != nil {
		log.Fatalf("Failed to load config from %s: %v", configPath, err)
	}

	clientID := cfg.Plaid.GetClientID()
	secret := cfg.Plaid.GetSecret()

	if clientID == "" || secret == "" {
		log.Fatal("Plaid client_id and secret are required in config file")
	}

	configuration := plaid.NewConfiguration()
	configuration.AddDefaultHeader("PLAID-CLIENT-ID", clientID)
	configuration.AddDefaultHeader("PLAID-SECRET", secret)
	configuration.UseEnvironment(plaid.Sandbox)

	client := plaid.NewAPIClient(configuration)
	ctx := context.Background()

	institutionID := "ins_109508" // First Platypus Bank (sandbox test institution)
	initialProducts := []plaid.Products{plaid.PRODUCTS_TRANSACTIONS}

	fmt.Println("🏦 Creating Plaid Sandbox Account for Overspending Simulation...")
	fmt.Printf("📋 Using config from: %s\n", configPath)
	fmt.Printf("📋 Using Client ID: %s\n", clientID[:10]+"...")

	// Create a public token using sandbox endpoint
	publicTokenReq := plaid.NewSandboxPublicTokenCreateRequest(
		institutionID,
		initialProducts,
	)

	// Set options for test credentials and transaction history
	options := plaid.NewSandboxPublicTokenCreateRequestOptions()

	// Use transactions dynamic user for realistic data
	username := "user_transactions_dynamic"
	options.SetOverrideUsername(username)
	options.SetOverridePassword("pass_good")

	// Set transaction options for 3 months of historical data
	startDate := time.Now().AddDate(0, -3, 0).Format("2006-01-02")
	endDate := time.Now().Format("2006-01-02")

	transactions := plaid.NewSandboxPublicTokenCreateRequestOptionsTransactions()
	transactions.SetStartDate(startDate)
	transactions.SetEndDate(endDate)
	options.SetTransactions(*transactions)

	publicTokenReq.SetOptions(*options)

	publicTokenResp, _, err := client.PlaidApi.SandboxPublicTokenCreate(ctx).
		SandboxPublicTokenCreateRequest(*publicTokenReq).
		Execute()
	if err != nil {
		log.Fatalf("Failed to create public token: %v", err)
	}

	fmt.Printf("✓ Public token created: %s\n", publicTokenResp.GetPublicToken())

	// Exchange public token for access token
	exchangeReq := plaid.NewItemPublicTokenExchangeRequest(publicTokenResp.GetPublicToken())
	exchangeResp, _, err := client.PlaidApi.ItemPublicTokenExchange(ctx).
		ItemPublicTokenExchangeRequest(*exchangeReq).
		Execute()
	if err != nil {
		log.Fatalf("Failed to exchange public token: %v", err)
	}

	accessToken := exchangeResp.GetAccessToken()
	itemID := exchangeResp.GetItemId()

	fmt.Printf("✓ Access token obtained\n")
	fmt.Printf("✓ Item ID: %s\n", itemID)

	// Get accounts to verify
	accountsReq := plaid.NewAccountsGetRequest(accessToken)
	accountsResp, _, err := client.PlaidApi.AccountsGet(ctx).
		AccountsGetRequest(*accountsReq).
		Execute()
	if err != nil {
		log.Fatalf("Failed to get accounts: %v", err)
	}

	fmt.Printf("\n📊 Accounts Created:\n")
	var checkingAccountID string
	for _, account := range accountsResp.GetAccounts() {
		balances := account.GetBalances()
		fmt.Printf("  - %s (%s): $%.2f\n",
			account.GetName(),
			account.GetSubtype(),
			balances.GetCurrent())

		// Store the first checking account for transaction simulation
		if account.GetSubtype() == "checking" && checkingAccountID == "" {
			checkingAccountID = account.GetAccountId()
		}
	}

	// Get initial transactions to verify
	fmt.Printf("\n💸 Fetching transaction history...\n")

	syncReq := plaid.NewTransactionsSyncRequest(accessToken)
	syncResp, _, err := client.PlaidApi.TransactionsSync(ctx).
		TransactionsSyncRequest(*syncReq).
		Execute()
	if err != nil {
		log.Printf("Warning: Failed to fetch transactions: %v", err)
	} else {
		addedTxs := syncResp.GetAdded()
		fmt.Printf("  ✓ %d transactions available\n", len(addedTxs))

		if len(addedTxs) > 0 {
			fmt.Printf("\n📝 Sample transactions:\n")
			// Show first 5 transactions
			for i, tx := range addedTxs {
				if i >= 5 {
					break
				}
				fmt.Printf("  - %s: $%.2f (%s)\n",
					tx.GetName(),
					tx.GetAmount(),
					tx.GetDate())
			}
		}
	}

	// Save credentials to file
	credentials := map[string]interface{}{
		"access_token":        accessToken,
		"item_id":             itemID,
		"public_token":        publicTokenResp.GetPublicToken(),
		"checking_account_id": checkingAccountID,
		"created_at":          time.Now().Format(time.RFC3339),
		"username":            username,
		"institution":         institutionID,
		"transaction_period":  "3_months",
	}

	file, err := os.Create("overspending_account.json")
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(credentials); err != nil {
		log.Fatalf("Failed to write credentials: %v", err)
	}

	fmt.Printf("\n✅ Setup Complete!\n")
	fmt.Printf("✓ Credentials saved to overspending_account.json\n")
	fmt.Printf("\n🎯 Public token to use in your app:\n")
	fmt.Printf("   %s\n", publicTokenResp.GetPublicToken())
	fmt.Printf("\n💡 Access token for API calls:\n")
	fmt.Printf("   %s\n", accessToken)
	fmt.Printf("\n🔧 Checking Account ID for transaction generator:\n")
	fmt.Printf("   %s\n", checkingAccountID)
	fmt.Printf("\n📖 Next steps:\n")
	fmt.Printf("   1. Exchange the public token in your app via POST /plaid/exchange-token\n")
	fmt.Printf("   2. Start the transaction generator to simulate ongoing overspending\n")
	fmt.Printf("   3. See docs/OVERSPENDING_SIMULATION.md for detailed instructions\n")
}
