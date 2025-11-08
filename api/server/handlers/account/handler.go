package account

import (
	"regulation/internal/ent"
	"regulation/internal/session"
)

type Handler struct {
	db             *ent.Client
	sessionManager *session.Manager
}

func New(db *ent.Client, sessionManager *session.Manager) *Handler {
	return &Handler{
		db:             db,
		sessionManager: sessionManager,
	}
}
