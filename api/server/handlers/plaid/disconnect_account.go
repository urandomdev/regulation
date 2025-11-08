package plaid

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"regulation/internal/ent/account"
	"regulation/internal/protocol"
	"regulation/server/services/request_context"
)

// DisconnectAccount soft-deletes an account (used during onboarding to remove incorrectly linked accounts)
// @Route DELETE /plaid/accounts/:id
func (h *Handler) DisconnectAccount(ctx fiber.Ctx) (*DisconnectAccountResponse, error) {
	session := request_context.Session(ctx)

	// Parse account ID from URL
	accountIDStr := ctx.Params("id")
	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		return nil, protocol.ErrorResponse{
			Code:    protocol.InvalidParametersError,
			Message: "invalid account id",
		}
	}

	// Verify account belongs to user and soft delete
	result, err := h.db.Account.
		Update().
		Where(account.ID(accountID)).
		Where(account.UserID(session.UserID)).
		SetIsActive(false).
		Save(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to disconnect account: %w", err)
	}

	if result == 0 {
		return nil, protocol.ErrorResponse{
			Code:    protocol.NotFoundError,
			Message: "account not found",
		}
	}

	return &DisconnectAccountResponse{
		Success: true,
	}, nil
}

type DisconnectAccountResponse struct {
	Success bool `cbor:"success" json:"success"`
}
