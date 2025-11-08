package plaid

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"regulation/server/services/request_context"
)

// SyncTransactions manually triggers transaction synchronization
// @Route POST /plaid/sync-transactions
func (h *Handler) SyncTransactions(ctx fiber.Ctx, req *SyncTransactionsRequest) error {
	session := request_context.Session(ctx)

	var syncedCount int
	var err error

	// Sync specific item or all items
	if req.ItemID != nil {
		// Sync specific item
		syncedCount, err = h.syncService.SyncItemTransactions(ctx, *req.ItemID)
	} else {
		// Sync all user items
		syncedCount, err = h.syncService.SyncAllUserItems(ctx, session.UserID)
	}

	if err != nil {
		return fmt.Errorf("failed to sync transactions: %w", err)
	}

	return ctx.JSON(SyncTransactionsResponse{
		SyncedCount: syncedCount,
	})
}

type SyncTransactionsRequest struct {
	ItemID *uuid.UUID `cbor:"item_id,omitempty" json:"item_id,omitempty"`
}

type SyncTransactionsResponse struct {
	SyncedCount int `cbor:"synced_count" json:"synced_count"`
}
