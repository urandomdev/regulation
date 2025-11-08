package recommendation

import (
	"regulation/internal/ent"
	"regulation/internal/rulesuggestion"
	"regulation/internal/session"
)

// Handler handles rule recommendation requests
type Handler struct {
	db                *ent.Client
	sessionManager    *session.Manager
	suggestionService *rulesuggestion.Service
}

// New creates a new recommendation handler
func New(db *ent.Client, sessionManager *session.Manager, suggestionService *rulesuggestion.Service) *Handler {
	return &Handler{
		db:                db,
		sessionManager:    sessionManager,
		suggestionService: suggestionService,
	}
}
