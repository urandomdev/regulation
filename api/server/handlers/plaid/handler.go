package plaid

import (
	"regulation/internal/ent"
	"regulation/internal/session"
	"regulation/server/services/plaid"
)

// Handler handles Plaid-related requests
type Handler struct {
	db             *ent.Client
	sessionManager *session.Manager
	plaidClient    plaid.PlaidClient
	syncService    *plaid.SyncService
}

// New creates a new Plaid handler
func New(
	db *ent.Client,
	sessionManager *session.Manager,
	plaidClient plaid.PlaidClient,
	syncService *plaid.SyncService,
) *Handler {
	return &Handler{
		db:             db,
		sessionManager: sessionManager,
		plaidClient:    plaidClient,
		syncService:    syncService,
	}
}
