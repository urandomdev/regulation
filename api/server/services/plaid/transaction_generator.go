package plaid

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// TransactionGenerator triggers Plaid's /transactions/refresh every 15 seconds
// This works with Plaid's test accounts (like user_transactions_dynamic) which
// automatically generate new transactions when /transactions/refresh is called.
type TransactionGenerator struct {
	plaidClient Client
	accessToken string
	accountID   string

	running  bool
	stopChan chan struct{}
	mu       sync.Mutex
	interval time.Duration
}

// NewTransactionGenerator creates a new transaction generator
// Note: Use with test accounts like user_transactions_dynamic for automatic transaction generation
func NewTransactionGenerator(plaidClient Client, accessToken, accountID string) *TransactionGenerator {
	return &TransactionGenerator{
		plaidClient: plaidClient,
		accessToken: accessToken,
		accountID:   accountID,
		interval:    15 * time.Second,
	}
}

// Start begins triggering transaction refreshes every 15 seconds
// Note: This will only work with special Plaid test accounts like user_transactions_dynamic
func (g *TransactionGenerator) Start(ctx context.Context) error {
	g.mu.Lock()
	if g.running {
		g.mu.Unlock()
		return fmt.Errorf("transaction generator already running")
	}

	g.running = true
	g.stopChan = make(chan struct{})
	g.mu.Unlock()

	log.Info().
		Str("account_id", g.accountID).
		Dur("interval", g.interval).
		Msg("Transaction generator started - refreshing transactions every 15 seconds")

	// Start the refresh loop in a goroutine
	go g.refreshLoop(ctx)

	return nil
}

// Stop stops the transaction generator
func (g *TransactionGenerator) Stop() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if !g.running {
		return fmt.Errorf("transaction generator not running")
	}

	if g.stopChan != nil {
		close(g.stopChan)
	}
	g.running = false

	log.Info().Msg("Transaction generator stopped")

	return nil
}

// IsRunning returns whether the generator is currently running
func (g *TransactionGenerator) IsRunning() bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.running
}

// refreshLoop continuously calls /transactions/refresh every 15 seconds
func (g *TransactionGenerator) refreshLoop(ctx context.Context) {
	ticker := time.NewTicker(g.interval)
	defer ticker.Stop()

	// Do an initial refresh immediately
	if err := g.refresh(ctx); err != nil {
		log.Error().
			Err(err).
			Msg("Initial transaction refresh failed")
	}

	for {
		select {
		case <-g.stopChan:
			log.Info().Msg("Transaction generator refresh loop stopped")
			return
		case <-ctx.Done():
			log.Info().Msg("Transaction generator context cancelled")
			return
		case <-ticker.C:
			if err := g.refresh(ctx); err != nil {
				log.Error().
					Err(err).
					Msg("Failed to refresh transactions")
			}
		}
	}
}

// refresh calls Plaid's /transactions/refresh endpoint
func (g *TransactionGenerator) refresh(ctx context.Context) error {
	log.Info().
		Str("account_id", g.accountID).
		Msg("Refreshing transactions...")

	if err := g.plaidClient.RefreshTransactions(ctx, g.accessToken); err != nil {
		return fmt.Errorf("failed to refresh transactions: %w", err)
	}

	log.Info().
		Str("account_id", g.accountID).
		Msg("Transactions refreshed successfully - new transactions should be available via sync")

	return nil
}
