package plaid

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"regulation/internal/categorizer"
	"regulation/internal/config"
	"regulation/internal/ent"
	entaccount "regulation/internal/ent/account"
	entitem "regulation/internal/ent/item"
	entsynccursor "regulation/internal/ent/synccursor"
	enttransaction "regulation/internal/ent/transaction"
	"regulation/server/services/rule"
)

// SyncService handles transaction synchronization from Plaid
type SyncService struct {
	plaidClient        Client
	entClient          *ent.Client
	ruleEngine         *rule.Engine
	categorizerService *categorizer.Service
}

// NewSyncService creates a new transaction sync service
func NewSyncService(plaidClient Client, entClient *ent.Client, cfg *config.Config, categorizerSvc *categorizer.Service) *SyncService {
	return &SyncService{
		plaidClient:        plaidClient,
		entClient:          entClient,
		ruleEngine:         rule.NewEngine(entClient, cfg),
		categorizerService: categorizerSvc,
	}
}

// SyncItemTransactions synchronizes transactions for a single Item
func (s *SyncService) SyncItemTransactions(ctx context.Context, itemID uuid.UUID) (int, error) {
	// Load Item with access token
	item, err := s.entClient.Item.
		Query().
		Where(entitem.ID(itemID)).
		Where(entitem.IsActive(true)).
		Only(ctx)
	if err != nil {
		// Record failure in cursor
		_ = s.recordSyncError(ctx, itemID, fmt.Errorf("failed to load item: %w", err))
		return 0, fmt.Errorf("failed to load item: %w", err)
	}

	// Load or create sync cursor
	cursor, err := s.getOrCreateCursor(ctx, itemID)
	if err != nil {
		_ = s.recordSyncError(ctx, itemID, fmt.Errorf("failed to get sync cursor: %w", err))
		return 0, fmt.Errorf("failed to get sync cursor: %w", err)
	}

	totalSynced := 0
	currentCursor := cursor.Cursor
	hasMore := true

	// Keep syncing while there are more transactions
	for hasMore {
		// Call Plaid API to sync transactions
		result, err := s.plaidClient.SyncTransactions(ctx, item.AccessToken, currentCursor)
		if err != nil {
			_ = s.recordSyncError(ctx, itemID, fmt.Errorf("failed to sync transactions: %w", err))
			return totalSynced, fmt.Errorf("failed to sync transactions: %w", err)
		}

		// Process added transactions
		if err := s.processAddedTransactions(ctx, result.Added); err != nil {
			_ = s.recordSyncError(ctx, itemID, fmt.Errorf("failed to process added transactions: %w", err))
			return totalSynced, fmt.Errorf("failed to process added transactions: %w", err)
		}
		totalSynced += len(result.Added)

		// Process modified transactions
		if err := s.processModifiedTransactions(ctx, result.Modified); err != nil {
			_ = s.recordSyncError(ctx, itemID, fmt.Errorf("failed to process modified transactions: %w", err))
			return totalSynced, fmt.Errorf("failed to process modified transactions: %w", err)
		}

		// Process removed transactions
		if err := s.processRemovedTransactions(ctx, result.Removed); err != nil {
			_ = s.recordSyncError(ctx, itemID, fmt.Errorf("failed to process removed transactions: %w", err))
			return totalSynced, fmt.Errorf("failed to process removed transactions: %w", err)
		}

		// Update cursor
		currentCursor = result.Cursor
		hasMore = result.HasMore

		log.Info().
			Str("item_id", itemID.String()).
			Int("added", len(result.Added)).
			Int("modified", len(result.Modified)).
			Int("removed", len(result.Removed)).
			Bool("has_more", hasMore).
			Msg("Transaction sync batch complete")
	}

	// Sync account balances after syncing all transactions
	if err := s.syncAccountBalances(ctx, item.AccessToken); err != nil {
		// Log but don't fail the whole sync if balance update fails
		log.Error().
			Err(err).
			Str("item_id", itemID.String()).
			Msg("Failed to sync account balances")
	}

	// Save the final cursor and mark success
	if err := s.recordSyncSuccess(ctx, itemID, currentCursor); err != nil {
		return totalSynced, fmt.Errorf("failed to update cursor: %w", err)
	}

	return totalSynced, nil
}

// SyncAllUserItems synchronizes all active items for a user
func (s *SyncService) SyncAllUserItems(ctx context.Context, userID uuid.UUID) (int, error) {
	// Load all active Items for user
	items, err := s.entClient.Item.
		Query().
		Where(entitem.UserID(userID)).
		Where(entitem.IsActive(true)).
		All(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to load items: %w", err)
	}

	totalSynced := 0
	for _, item := range items {
		synced, err := s.SyncItemTransactions(ctx, item.ID)
		if err != nil {
			log.Error().
				Err(err).
				Str("item_id", item.ID.String()).
				Msg("Failed to sync item transactions")
			continue // Continue with other items
		}
		totalSynced += synced
	}

	return totalSynced, nil
}

// processAddedTransactions processes newly added transactions
func (s *SyncService) processAddedTransactions(ctx context.Context, transactions []Transaction) error {
	if len(transactions) == 0 {
		return nil
	}

	log.Info().
		Int("transaction_count", len(transactions)).
		Msg("[SYNC] Processing added transactions")

	// Step 1: Build categorization requests for batch processing
	requests := make([]*categorizer.CategorizationRequest, len(transactions))
	for i, tx := range transactions {
		requests[i] = &categorizer.CategorizationRequest{
			MerchantName:    tx.MerchantName,
			TransactionName: tx.Name,
			Amount:          float64(tx.Amount),
			PlaidCategories: tx.Categories,
		}
	}

	// Step 2: Batch categorize all transactions concurrently
	responses, errors := s.categorizerService.CategorizeBatch(ctx, requests)

	// Step 3: Process each transaction with its categorization result
	for i, tx := range transactions {
		// Get Account
		account, err := s.entClient.Account.
			Query().
			Where(entaccount.PlaidID(tx.AccountID)).
			Only(ctx)
		if err != nil {
			log.Error().
				Err(err).
				Str("account_id", tx.AccountID).
				Msg("Failed to find account for transaction")
			continue
		}

		// Get categorization result or use fallback
		var categoryResp *categorizer.CategorizationResponse
		if errors[i] != nil {
			log.Error().
				Err(errors[i]).
				Str("transaction_id", tx.TransactionID).
				Msg("Failed to categorize transaction with GPT, using fallback")
			// Fallback to rule-based categorization
			categoryResp = &categorizer.CategorizationResponse{
				Category:   fallbackCategorize(tx.Categories),
				Confidence: "low",
				Reasoning:  "Fallback categorization due to GPT error",
			}
		} else {
			categoryResp = responses[i]
		}

		category := categoryResp.Category

		log.Info().
			Str("transaction_id", tx.TransactionID).
			Str("transaction_name", tx.Name).
			Str("merchant_name", tx.MerchantName).
			Str("assigned_category", string(category)).
			Str("confidence", categoryResp.Confidence).
			Strs("plaid_categories", tx.Categories).
			Int64("amount_cents", tx.Amount).
			Bool("pending", tx.Pending).
			Msg("[SYNC] Categorized transaction")

		// Upsert transaction
		err = s.entClient.Transaction.
			Create().
			SetID(uuid.New()).
			SetAccountID(account.ID).
			SetPlaidID(tx.TransactionID).
			SetAmount(tx.Amount).
			SetDate(tx.Date).
			SetName(tx.Name).
			SetNillableMerchantName(&tx.MerchantName).
			SetCategory(string(category)).
			SetPlaidCategories(tx.Categories).
			SetPending(tx.Pending).
			SetNillablePaymentChannel(&tx.PaymentChannel).
			OnConflictColumns(enttransaction.FieldPlaidID).
			UpdateNewValues().
			Exec(ctx)

		if err != nil {
			log.Error().
				Err(err).
				Str("transaction_id", tx.TransactionID).
				Msg("Failed to save transaction")
			continue
		}

		// NEW: Process rules for non-pending transactions
		if !tx.Pending {
			log.Info().
				Str("transaction_id", tx.TransactionID).
				Str("transaction_name", tx.Name).
				Int64("amount", tx.Amount).
				Str("category", string(category)).
				Msg("[SYNC] Transaction is non-pending, triggering rule processing")

			// Fetch the saved transaction entity
			entTx, err := s.entClient.Transaction.
				Query().
				Where(enttransaction.PlaidID(tx.TransactionID)).
				Only(ctx)

			if err != nil {
				log.Error().
					Err(err).
					Str("transaction_id", tx.TransactionID).
					Msg("[SYNC] Failed to fetch transaction for rule processing")
				continue
			}

			log.Debug().
				Str("transaction_id", tx.TransactionID).
				Str("ent_transaction_id", entTx.ID.String()).
				Msg("[SYNC] Fetched transaction entity, calling rule engine")

			// Evaluate and execute rules
			if err := s.ruleEngine.ProcessTransaction(ctx, entTx); err != nil {
				log.Error().
					Err(err).
					Str("transaction_id", tx.TransactionID).
					Msg("[SYNC] Failed to process rules for transaction")
				// Don't fail the sync, just log and continue
			} else {
				log.Info().
					Str("transaction_id", tx.TransactionID).
					Msg("[SYNC] Successfully processed rules for transaction")
			}
		} else {
			log.Debug().
				Str("transaction_id", tx.TransactionID).
				Msg("[SYNC] Transaction is pending, skipping rule processing")
		}
	}

	return nil
}

// processModifiedTransactions processes modified transactions
func (s *SyncService) processModifiedTransactions(ctx context.Context, transactions []Transaction) error {
	if len(transactions) == 0 {
		return nil
	}

	log.Info().
		Int("transaction_count", len(transactions)).
		Msg("[SYNC] Processing modified transactions")

	// Step 1: Build categorization requests for batch processing
	requests := make([]*categorizer.CategorizationRequest, len(transactions))
	for i, tx := range transactions {
		requests[i] = &categorizer.CategorizationRequest{
			MerchantName:    tx.MerchantName,
			TransactionName: tx.Name,
			Amount:          float64(tx.Amount),
			PlaidCategories: tx.Categories,
		}
	}

	// Step 2: Batch categorize all transactions concurrently
	responses, errors := s.categorizerService.CategorizeBatch(ctx, requests)

	// Step 3: Process each transaction with its categorization result
	for i, tx := range transactions {
		// Get Account
		account, err := s.entClient.Account.
			Query().
			Where(entaccount.PlaidID(tx.AccountID)).
			Only(ctx)
		if err != nil {
			log.Error().
				Err(err).
				Str("account_id", tx.AccountID).
				Msg("Failed to find account for transaction")
			continue
		}

		// Get categorization result or use fallback
		var categoryResp *categorizer.CategorizationResponse
		if errors[i] != nil {
			log.Error().
				Err(errors[i]).
				Str("transaction_id", tx.TransactionID).
				Msg("Failed to categorize transaction with GPT, using fallback")
			// Fallback to rule-based categorization
			categoryResp = &categorizer.CategorizationResponse{
				Category:   fallbackCategorize(tx.Categories),
				Confidence: "low",
				Reasoning:  "Fallback categorization due to GPT error",
			}
		} else {
			categoryResp = responses[i]
		}

		category := categoryResp.Category

		log.Info().
			Str("transaction_id", tx.TransactionID).
			Str("transaction_name", tx.Name).
			Str("merchant_name", tx.MerchantName).
			Str("assigned_category", string(category)).
			Str("confidence", categoryResp.Confidence).
			Strs("plaid_categories", tx.Categories).
			Int64("amount_cents", tx.Amount).
			Bool("pending", tx.Pending).
			Msg("[SYNC] Categorized modified transaction")

		// Update existing transaction
		err = s.entClient.Transaction.
			Update().
			Where(enttransaction.PlaidID(tx.TransactionID)).
			SetAccountID(account.ID).
			SetAmount(tx.Amount).
			SetDate(tx.Date).
			SetName(tx.Name).
			SetNillableMerchantName(&tx.MerchantName).
			SetCategory(string(category)).
			SetPlaidCategories(tx.Categories).
			SetPending(tx.Pending).
			SetNillablePaymentChannel(&tx.PaymentChannel).
			Exec(ctx)

		if err != nil {
			log.Error().
				Err(err).
				Str("transaction_id", tx.TransactionID).
				Msg("Failed to update transaction")
			continue
		}

		// Process rules for non-pending transactions
		// This is critical for when transactions change from pending to non-pending
		if !tx.Pending {
			log.Info().
				Str("transaction_id", tx.TransactionID).
				Str("transaction_name", tx.Name).
				Int64("amount", tx.Amount).
				Str("category", string(category)).
				Msg("[SYNC] Modified transaction is non-pending, triggering rule processing")

			// Fetch the updated transaction entity
			entTx, err := s.entClient.Transaction.
				Query().
				Where(enttransaction.PlaidID(tx.TransactionID)).
				Only(ctx)

			if err != nil {
				log.Error().
					Err(err).
					Str("transaction_id", tx.TransactionID).
					Msg("[SYNC] Failed to fetch transaction for rule processing")
				continue
			}

			log.Debug().
				Str("transaction_id", tx.TransactionID).
				Str("ent_transaction_id", entTx.ID.String()).
				Msg("[SYNC] Fetched modified transaction entity, calling rule engine")

			// Evaluate and execute rules
			if err := s.ruleEngine.ProcessTransaction(ctx, entTx); err != nil {
				log.Error().
					Err(err).
					Str("transaction_id", tx.TransactionID).
					Msg("[SYNC] Failed to process rules for modified transaction")
				// Don't fail the sync, just log and continue
			} else {
				log.Info().
					Str("transaction_id", tx.TransactionID).
					Msg("[SYNC] Successfully processed rules for modified transaction")
			}
		} else {
			log.Debug().
				Str("transaction_id", tx.TransactionID).
				Msg("[SYNC] Modified transaction is pending, skipping rule processing")
		}
	}

	return nil
}

// processRemovedTransactions processes removed transactions
func (s *SyncService) processRemovedTransactions(ctx context.Context, transactionIDs []string) error {
	for _, txID := range transactionIDs {
		// Delete transaction
		_, err := s.entClient.Transaction.
			Delete().
			Where(enttransaction.PlaidID(txID)).
			Exec(ctx)

		if err != nil {
			log.Error().
				Err(err).
				Str("transaction_id", txID).
				Msg("Failed to delete transaction")
			continue
		}
	}

	return nil
}

// getOrCreateCursor gets or creates a sync cursor for an item
func (s *SyncService) getOrCreateCursor(ctx context.Context, itemID uuid.UUID) (*ent.SyncCursor, error) {
	cursor, err := s.entClient.SyncCursor.
		Query().
		Where(entsynccursor.ItemID(itemID)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			// Create new cursor
			return s.entClient.SyncCursor.
				Create().
				SetItemID(itemID).
				SetCursor("").
				Save(ctx)
		}
		return nil, err
	}

	return cursor, nil
}

// updateCursor updates the sync cursor for an item
func (s *SyncService) updateCursor(ctx context.Context, itemID uuid.UUID, newCursor string) error {
	return s.entClient.SyncCursor.
		Update().
		Where(entsynccursor.ItemID(itemID)).
		SetCursor(newCursor).
		Exec(ctx)
}

// recordSyncSuccess records a successful sync and resets error counters
func (s *SyncService) recordSyncSuccess(ctx context.Context, itemID uuid.UUID, newCursor string) error {
	return s.entClient.SyncCursor.
		Update().
		Where(entsynccursor.ItemID(itemID)).
		SetCursor(newCursor).
		SetLastError("").
		SetConsecutiveFailures(0).
		Exec(ctx)
}

// recordSyncError records a sync error and increments failure counter
func (s *SyncService) recordSyncError(ctx context.Context, itemID uuid.UUID, err error) error {
	// Get current cursor to increment consecutive failures
	cursor, queryErr := s.entClient.SyncCursor.
		Query().
		Where(entsynccursor.ItemID(itemID)).
		Only(ctx)

	if queryErr != nil {
		// If we can't query the cursor, just log and return
		log.Error().
			Err(queryErr).
			Str("item_id", itemID.String()).
			Msg("Failed to query cursor for error recording")
		return queryErr
	}

	// Increment consecutive failures
	newFailures := cursor.ConsecutiveFailures + 1

	return s.entClient.SyncCursor.
		Update().
		Where(entsynccursor.ItemID(itemID)).
		SetLastError(err.Error()).
		SetConsecutiveFailures(newFailures).
		Exec(ctx)
}

// fallbackCategorize provides rule-based categorization as a fallback when GPT fails
func fallbackCategorize(plaidCategories []string) categorizer.CategoryType {
	if len(plaidCategories) == 0 {
		return categorizer.CategoryMisc
	}

	// Get the primary category (first element)
	primary := strings.ToLower(plaidCategories[0])

	// Check subcategories for more specific matching
	var secondary string
	if len(plaidCategories) > 1 {
		secondary = strings.ToLower(plaidCategories[1])
	}

	// Map Plaid categories to app categories
	switch {
	// Transfer - ignore for rules (internal transfers)
	case strings.Contains(primary, "transfer"):
		return categorizer.CategoryTransfer
	case strings.Contains(primary, "payment"):
		return categorizer.CategoryTransfer

	// Dining
	case strings.Contains(primary, "food and drink"):
		// Check if it's groceries
		if strings.Contains(secondary, "groceries") || strings.Contains(secondary, "supermarket") {
			return categorizer.CategoryGroceries
		}
		return categorizer.CategoryDining
	case strings.Contains(primary, "restaurants"):
		return categorizer.CategoryDining

	// Groceries
	case strings.Contains(primary, "groceries"):
		return categorizer.CategoryGroceries
	case strings.Contains(secondary, "supermarket"):
		return categorizer.CategoryGroceries

	// Transport
	case strings.Contains(primary, "transportation"):
		return categorizer.CategoryTransport
	case strings.Contains(primary, "travel"):
		return categorizer.CategoryTransport
	case strings.Contains(secondary, "gas"):
		return categorizer.CategoryTransport
	case strings.Contains(secondary, "parking"):
		return categorizer.CategoryTransport
	case strings.Contains(secondary, "public transit"):
		return categorizer.CategoryTransport
	case strings.Contains(secondary, "ride share"):
		return categorizer.CategoryTransport

	// Shopping
	case strings.Contains(primary, "shops"):
		return categorizer.CategoryShopping
	case strings.Contains(primary, "retail"):
		return categorizer.CategoryShopping
	case strings.Contains(secondary, "clothing"):
		return categorizer.CategoryShopping
	case strings.Contains(secondary, "electronics"):
		return categorizer.CategoryShopping

	// Subscriptions
	case strings.Contains(secondary, "subscription"):
		return categorizer.CategorySubscriptions
	case strings.Contains(primary, "service"):
		// Check if it's a recurring service
		if strings.Contains(secondary, "streaming") ||
			strings.Contains(secondary, "music") ||
			strings.Contains(secondary, "software") {
			return categorizer.CategorySubscriptions
		}
		return categorizer.CategoryMisc

	// Entertainment
	case strings.Contains(primary, "recreation"):
		return categorizer.CategoryEntertainment
	case strings.Contains(primary, "entertainment"):
		return categorizer.CategoryEntertainment
	case strings.Contains(secondary, "movie"):
		return categorizer.CategoryEntertainment
	case strings.Contains(secondary, "concert"):
		return categorizer.CategoryEntertainment
	case strings.Contains(secondary, "sporting"):
		return categorizer.CategoryEntertainment

	// Bills
	case strings.Contains(primary, "bank fees"):
		return categorizer.CategoryBills
	case strings.Contains(primary, "interest"):
		return categorizer.CategoryBills
	case strings.Contains(secondary, "utilities"):
		return categorizer.CategoryBills
	case strings.Contains(secondary, "internet"):
		return categorizer.CategoryBills
	case strings.Contains(secondary, "phone"):
		return categorizer.CategoryBills
	case strings.Contains(secondary, "insurance"):
		return categorizer.CategoryBills
	case strings.Contains(secondary, "rent"):
		return categorizer.CategoryBills

	default:
		return categorizer.CategoryMisc
	}
}

// syncAccountBalances fetches and updates account balances from Plaid
func (s *SyncService) syncAccountBalances(ctx context.Context, accessToken string) error {
	// Get accounts from Plaid with current balances
	accounts, err := s.plaidClient.GetAccounts(ctx, accessToken)
	if err != nil {
		return fmt.Errorf("failed to get accounts from Plaid: %w", err)
	}

	// Update each account's balance in the database
	for _, acc := range accounts {
		err := s.entClient.Account.
			Update().
			Where(entaccount.PlaidID(acc.AccountID)).
			SetCurrentBalance(acc.BalanceCurrent).
			SetNillableAvailableBalance(acc.BalanceAvailable).
			Exec(ctx)

		if err != nil {
			log.Error().
				Err(err).
				Str("account_id", acc.AccountID).
				Msg("Failed to update account balance")
			// Continue with other accounts even if one fails
			continue
		}

		log.Debug().
			Str("account_id", acc.AccountID).
			Int64("current_balance", acc.BalanceCurrent).
			Msg("Account balance updated")
	}

	return nil
}
