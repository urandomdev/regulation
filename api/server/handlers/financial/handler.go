package financial

import (
	"regulation/internal/ent"
	"regulation/internal/session"
)

// Handler manages financial data queries for dashboard
type Handler struct {
	db             *ent.Client
	sessionManager *session.Manager
}

// New creates a new financial handler
func New(db *ent.Client, sessionManager *session.Manager) *Handler {
	return &Handler{
		db:             db,
		sessionManager: sessionManager,
	}
}
