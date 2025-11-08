package notification

import (
	"regulation/internal/config"
	"regulation/internal/ent"
)

type Handler struct {
	config *config.Config
	db     *ent.Client
}

func New(config *config.Config, db *ent.Client) *Handler {
	return &Handler{
		config: config,
		db:     db,
	}
}
