package plaid

import (
	"fmt"

	"github.com/gofiber/fiber/v3"

	"regulation/server/services/request_context"
)

// CreateLinkToken creates a Plaid Link token for account linking
// @Route POST /plaid/create-link-token
func (h *Handler) CreateLinkToken(ctx fiber.Ctx) (*CreateLinkTokenResponse, error) {
	session := request_context.Session(ctx)

	// Create link token with user ID
	linkToken, err := h.plaidClient.CreateLinkToken(ctx, session.UserID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to create link token: %w", err)
	}

	return &CreateLinkTokenResponse{
		LinkToken: linkToken,
	}, nil
}

type CreateLinkTokenResponse struct {
	LinkToken string `cbor:"link_token" json:"link_token"`
}
