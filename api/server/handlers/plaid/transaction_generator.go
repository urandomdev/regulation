package plaid

import (
	"fmt"
	"sync"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	plaidservice "regulation/server/services/plaid"
	"regulation/server/services/request_context"
)

// TransactionGeneratorHandler handles API endpoints for the transaction generator
type TransactionGeneratorHandler struct {
	generators map[uuid.UUID]*plaidservice.TransactionGenerator // map of item_id -> generator
	mu         sync.RWMutex
	plaidSvc   *Handler
}

// NewTransactionGeneratorHandler creates a new handler for transaction generation
func NewTransactionGeneratorHandler(plaidHandler *Handler) *TransactionGeneratorHandler {
	return &TransactionGeneratorHandler{
		generators: make(map[uuid.UUID]*plaidservice.TransactionGenerator),
		plaidSvc:   plaidHandler,
	}
}

// StartGenerator starts the transaction generator for a specific item
// @Route POST /plaid/generator/start
func (h *TransactionGeneratorHandler) StartGenerator(ctx fiber.Ctx, req *StartGeneratorRequest) error {
	session := request_context.Session(ctx)

	// Verify the item belongs to the user
	item, err := h.plaidSvc.db.Item.Get(ctx, req.ItemID)
	if err != nil {
		return fmt.Errorf("item not found: %w", err)
	}

	if item.UserID != session.UserID {
		return fiber.NewError(fiber.StatusForbidden, "item does not belong to user")
	}

	// Get the account
	account, err := h.plaidSvc.db.Account.Get(ctx, req.AccountID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}

	if account.ItemID != req.ItemID {
		return fiber.NewError(fiber.StatusBadRequest, "account does not belong to item")
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	// Check if generator already exists and is running
	if gen, exists := h.generators[req.ItemID]; exists && gen.IsRunning() {
		return fiber.NewError(fiber.StatusConflict, "generator already running for this item")
	}

	// Create new generator
	generator := plaidservice.NewTransactionGenerator(
		h.plaidSvc.plaidClient,
		item.AccessToken,
		account.PlaidID,
	)

	if err := generator.Start(ctx); err != nil {
		return fmt.Errorf("failed to start generator: %w", err)
	}

	h.generators[req.ItemID] = generator

	return ctx.JSON(StartGeneratorResponse{
		Message: "Transaction generator started",
		ItemID:  req.ItemID,
		Running: true,
	})
}

// StopGenerator stops the transaction generator for a specific item
// @Route POST /plaid/generator/stop
func (h *TransactionGeneratorHandler) StopGenerator(ctx fiber.Ctx, req *StopGeneratorRequest) error {
	session := request_context.Session(ctx)

	// Verify the item belongs to the user
	item, err := h.plaidSvc.db.Item.Get(ctx, req.ItemID)
	if err != nil {
		return fmt.Errorf("item not found: %w", err)
	}

	if item.UserID != session.UserID {
		return fiber.NewError(fiber.StatusForbidden, "item does not belong to user")
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	generator, exists := h.generators[req.ItemID]
	if !exists || !generator.IsRunning() {
		return fiber.NewError(fiber.StatusNotFound, "no running generator for this item")
	}

	if err := generator.Stop(); err != nil {
		return fmt.Errorf("failed to stop generator: %w", err)
	}

	delete(h.generators, req.ItemID)

	return ctx.JSON(StopGeneratorResponse{
		Message: "Transaction generator stopped",
		ItemID:  req.ItemID,
		Running: false,
	})
}

// GetGeneratorStatus gets the status of the transaction generator for an item
// @Route GET /plaid/generator/status/:item_id
func (h *TransactionGeneratorHandler) GetGeneratorStatus(ctx fiber.Ctx) error {
	session := request_context.Session(ctx)

	itemIDStr := ctx.Params("item_id")
	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid item_id")
	}

	// Verify the item belongs to the user
	item, err := h.plaidSvc.db.Item.Get(ctx, itemID)
	if err != nil {
		return fmt.Errorf("item not found: %w", err)
	}

	if item.UserID != session.UserID {
		return fiber.NewError(fiber.StatusForbidden, "item does not belong to user")
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	generator, exists := h.generators[itemID]
	running := exists && generator.IsRunning()

	return ctx.JSON(GeneratorStatusResponse{
		ItemID:  itemID,
		Running: running,
	})
}

type StartGeneratorRequest struct {
	ItemID    uuid.UUID `cbor:"item_id" json:"item_id"`
	AccountID uuid.UUID `cbor:"account_id" json:"account_id"`
}

type StartGeneratorResponse struct {
	Message string    `cbor:"message" json:"message"`
	ItemID  uuid.UUID `cbor:"item_id" json:"item_id"`
	Running bool      `cbor:"running" json:"running"`
}

type StopGeneratorRequest struct {
	ItemID uuid.UUID `cbor:"item_id" json:"item_id"`
}

type StopGeneratorResponse struct {
	Message string    `cbor:"message" json:"message"`
	ItemID  uuid.UUID `cbor:"item_id" json:"item_id"`
	Running bool      `cbor:"running" json:"running"`
}

type GeneratorStatusResponse struct {
	ItemID  uuid.UUID `cbor:"item_id" json:"item_id"`
	Running bool      `cbor:"running" json:"running"`
}
