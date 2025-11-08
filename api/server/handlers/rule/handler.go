package rule

import (
	"regulation/internal/ent"
	"regulation/internal/session"
)

// Handler handles rule-related HTTP requests
type Handler struct {
	db             *ent.Client
	sessionManager *session.Manager
}

// New creates a new rule handler
func New(db *ent.Client, sessionManager *session.Manager) *Handler {
	return &Handler{
		db:             db,
		sessionManager: sessionManager,
	}
}
