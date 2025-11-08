package advisor

import (
	advisorsvc "regulation/internal/advisor"
	"regulation/internal/ent"
)

type Handler struct {
	service *advisorsvc.Service
	db      *ent.Client
}

func New(db *ent.Client, service *advisorsvc.Service) *Handler {
	return &Handler{
		service: service,
		db:      db,
	}
}
