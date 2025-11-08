package account

import "regulation/internal/ent"

type Handler struct {
	db *ent.Client
}

func New(db *ent.Client) *Handler {
	return &Handler{
		db: db,
	}
}
