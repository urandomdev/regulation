package plaid

import (
	"context"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"regulation/internal/ent"
	entitem "regulation/internal/ent/item"
)

// SyncWorker handles periodic background syncing of all active Plaid items
type SyncWorker struct {
	syncService  *SyncService
	entClient    *ent.Client
	syncInterval time.Duration
	minBackoff   time.Duration
	maxBackoff   time.Duration
}

// NewSyncWorker creates a new sync worker
func NewSyncWorker(syncService *SyncService, entClient *ent.Client, syncInterval time.Duration) *SyncWorker {
	return &SyncWorker{
		syncService:  syncService,
		entClient:    entClient,
		syncInterval: syncInterval,
		minBackoff:   1 * time.Second,
		maxBackoff:   60 * time.Second,
	}
}

// Start begins the sync worker loop
func (w *SyncWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.syncInterval)
	defer ticker.Stop()

	log.Info().
		Dur("interval", w.syncInterval).
		Msg("Plaid sync worker started")

	// Run initial sync immediately
	w.runSyncCycle(ctx)

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("Plaid sync worker stopping")
			return
		case <-ticker.C:
			w.runSyncCycle(ctx)
		}
	}
}

// runSyncCycle performs one complete sync cycle for all active items
func (w *SyncWorker) runSyncCycle(ctx context.Context) {
	log.Debug().Msg("Starting sync cycle for all active items")

	// Query all active items
	items, err := w.entClient.Item.
		Query().
		Where(entitem.IsActive(true)).
		All(ctx)

	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to query active items for sync")
		return
	}

	if len(items) == 0 {
		log.Debug().Msg("No active items to sync")
		return
	}

	log.Info().
		Int("item_count", len(items)).
		Msg("Syncing active items")

	// Sync each item with retry logic
	successCount := 0
	errorCount := 0

	for _, item := range items {
		if err := w.syncItemWithRetry(ctx, item.ID); err != nil {
			log.Error().
				Err(err).
				Str("item_id", item.ID.String()).
				Str("institution", item.InstitutionName).
				Msg("Failed to sync item after retries")
			errorCount++
		} else {
			log.Debug().
				Str("item_id", item.ID.String()).
				Str("institution", item.InstitutionName).
				Msg("Successfully synced item")
			successCount++
		}
	}

	log.Info().
		Int("success", successCount).
		Int("errors", errorCount).
		Int("total", len(items)).
		Msg("Sync cycle completed")
}

// syncItemWithRetry attempts to sync an item with exponential backoff retry
func (w *SyncWorker) syncItemWithRetry(ctx context.Context, itemID uuid.UUID) error {
	attempt := 0

	for {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		_, err := w.syncService.SyncItemTransactions(ctx, itemID)
		if err == nil {
			// Success!
			if attempt > 0 {
				log.Info().
					Str("item_id", itemID.String()).
					Int("attempts", attempt+1).
					Msg("Item sync succeeded after retries")
			}
			return nil
		}

		// Calculate backoff delay with exponential increase
		backoff := w.calculateBackoff(attempt)

		log.Warn().
			Err(err).
			Str("item_id", itemID.String()).
			Int("attempt", attempt+1).
			Dur("retry_in", backoff).
			Msg("Item sync failed, will retry")

		// Wait before retry (with context cancellation check)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoff):
			// Continue to next retry
		}

		attempt++
	}
}

// calculateBackoff calculates exponential backoff delay
func (w *SyncWorker) calculateBackoff(attempt int) time.Duration {
	// Exponential backoff: min * 2^attempt, capped at max
	backoff := time.Duration(float64(w.minBackoff) * math.Pow(2, float64(attempt)))

	if backoff > w.maxBackoff {
		backoff = w.maxBackoff
	}

	return backoff
}
