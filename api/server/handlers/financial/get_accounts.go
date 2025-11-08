package financial

import (
	"fmt"

	"github.com/gofiber/fiber/v3"

	"regulation/internal/ent/account"
	"regulation/server/handlers/models"
	"regulation/server/services/request_context"
)

// GetAccounts retrieves all accounts for the authenticated user
// @Route GET /financial/accounts
func (h *Handler) GetAccounts(ctx fiber.Ctx) (*GetAccountsResponse, error) {
	session := request_context.Session(ctx)

	// Query all active accounts for user
	accounts, err := h.db.Account.
		Query().
		Where(account.UserID(session.UserID)).
		Where(account.IsActive(true)).
		All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to query accounts: %w", err)
	}

	// Convert to response
	responses := make([]models.Account, 0, len(accounts))
	for _, acc := range accounts {
		maskPtr := &acc.Mask
		if acc.Mask == "" {
			maskPtr = nil
		}

		responses = append(responses, models.Account{
			ID:               acc.ID,
			Name:             acc.Name,
			Type:             string(acc.Type),
			Mask:             maskPtr,
			CurrentBalance:   acc.CurrentBalance,
			AvailableBalance: acc.AvailableBalance,
			IsActive:         acc.IsActive,
		})
	}

	return &GetAccountsResponse{
		Accounts: responses,
	}, nil
}

type GetAccountsResponse struct {
	Accounts []models.Account `cbor:"accounts" json:"accounts"`
}
