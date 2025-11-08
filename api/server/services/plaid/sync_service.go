package plaid

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"regulation/internal/ent"
	entaccount "regulation/internal/ent/account"
	entitem "regulation/internal/ent/item"
	entsynccursor "regulation/internal/ent/synccursor"
	enttransaction "regulation/internal/ent/transaction"
	"regulation/server/services"
)

// SyncService handles transaction synchronization from Plaid
type SyncService struct {
	plaidClient PlaidClient
	entClient   *ent.Client
}

// NewSyncService creates a new transaction sync service
func NewSyncService(plaidClient PlaidClient, entClient *ent.Client) *SyncService {
	return &SyncService{
		plaidClient: plaidClient,
		entClient:   entClient,
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
		return 0, fmt.Errorf("failed to load item: %w", err)
	}

	// Load or create sync cursor
	cursor, err := s.getOrCreateCursor(ctx, itemID)
	if err != nil {
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
			return totalSynced, fmt.Errorf("failed to sync transactions: %w", err)
		}

		// Process added transactions
		if err := s.processAddedTransactions(ctx, result.Added); err != nil {
			return totalSynced, fmt.Errorf("failed to process added transactions: %w", err)
		}
		totalSynced += len(result.Added)

		// Process modified transactions
		if err := s.processModifiedTransactions(ctx, result.Modified); err != nil {
			return totalSynced, fmt.Errorf("failed to process modified transactions: %w", err)
		}

		// Process removed transactions
		if err := s.processRemovedTransactions(ctx, result.Removed); err != nil {
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

	// Save the final cursor
	if err := s.updateCursor(ctx, itemID, currentCursor); err != nil {
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
	for _, tx := range transactions {
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

		// Categorize transaction
		category := services.CategorizeTransaction(tx.Categories)

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
			OnConflict().
			UpdateNewValues().
			Exec(ctx)

		if err != nil {
			log.Error().
				Err(err).
				Str("transaction_id", tx.TransactionID).
				Msg("Failed to save transaction")
			continue
		}
	}

	return nil
}

// processModifiedTransactions processes modified transactions
func (s *SyncService) processModifiedTransactions(ctx context.Context, transactions []Transaction) error {
	for _, tx := range transactions {
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

		// Categorize transaction
		category := services.CategorizeTransaction(tx.Categories)

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
