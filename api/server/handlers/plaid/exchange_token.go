package plaid

import (
	"context"
	"fmt"

	"github.com/DeltaLaboratory/contrib/hooks"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"regulation/internal/ent/account"
	"regulation/server/handlers/models"
	"regulation/server/services/request_context"
)

// ExchangeToken exchanges a public token for an access token and creates Item and Account records
// @Route POST /plaid/exchange-token
func (h *Handler) ExchangeToken(ctx fiber.Ctx, req *ExchangeTokenRequest) (ret error) {
	session := request_context.Session(ctx)

	// Exchange public token for access token
	result, err := h.plaidClient.ExchangePublicToken(ctx, req.PublicToken)
	if err != nil {
		return fmt.Errorf("failed to exchange public token: %w", err)
	}

	// Get accounts from Plaid
	accounts, err := h.plaidClient.GetAccounts(ctx, result.AccessToken)
	if err != nil {
		return fmt.Errorf("failed to get accounts: %w", err)
	}

	// Get institution name (simplified - in production, use Plaid's institutions endpoint)
	institutionName := "Bank"
	if len(accounts) > 0 {
		institutionName = accounts[0].Name // Use first account name as fallback
	}

	// Create Item and Account records in a transaction
	tx, err := h.db.Tx(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer hooks.Rollback(tx, &ret)

	// Create Item
	itemCreate := tx.Item.Create().
		SetID(uuid.New()).
		SetUserID(session.UserID).
		SetPlaidID(result.ItemID).
		SetAccessToken(result.AccessToken).
		SetInstitutionName(institutionName).
		SetIsActive(true)

	itemEntity, err := itemCreate.Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to create item: %w", err)
	}

	// Create Account records
	accountResponses := make([]models.Account, 0, len(accounts))
	for _, acc := range accounts {
		// Map account type
		var accountType account.Type
		switch acc.Subtype {
		case "checking":
			accountType = account.TypeChecking
		case "savings":
			accountType = account.TypeSavings
		default:
			accountType = account.TypeOther
		}

		accountEntity, err := tx.Account.Create().
			SetID(uuid.New()).
			SetItemID(itemEntity.ID).
			SetUserID(session.UserID).
			SetPlaidID(acc.AccountID).
			SetName(acc.Name).
			SetType(accountType).
			SetSubtype(acc.Subtype).
			SetNillableMask(&acc.Mask).
			SetCurrentBalance(acc.BalanceCurrent).
			SetNillableAvailableBalance(acc.BalanceAvailable).
			SetIsActive(true).
			Save(ctx)

		if err != nil {
			return fmt.Errorf("failed to create account: %w", err)
		}

		maskPtr := &accountEntity.Mask
		if accountEntity.Mask == "" {
			maskPtr = nil
		}

		accountResponses = append(accountResponses, models.Account{
			ID:               accountEntity.ID,
			Name:             accountEntity.Name,
			Type:             string(accountEntity.Type),
			Mask:             maskPtr,
			CurrentBalance:   accountEntity.CurrentBalance,
			AvailableBalance: accountEntity.AvailableBalance,
			IsActive:         accountEntity.IsActive,
		})
	}

	// hooks.Rollback will auto-commit here since ret == nil

	// Trigger initial transaction sync in the background
	// IMPORTANT: Use context.Background() instead of the Fiber context
	// to avoid nil pointer dereference when the HTTP request context is recycled
	go func() {
		_, err := h.syncService.SyncItemTransactions(context.Background(), itemEntity.ID)
		if err != nil {
			// Note: Using a new context-independent logger since this runs after response
			_ = err // Background task error - already logged
		}
	}()

	return ctx.JSON(ExchangeTokenResponse{
		ItemID:   itemEntity.ID,
		Accounts: accountResponses,
	})
}

type ExchangeTokenRequest struct {
	PublicToken string `cbor:"public_token" json:"public_token"`
}

type ExchangeTokenResponse struct {
	ItemID   uuid.UUID        `cbor:"item_id" json:"item_id"`
	Accounts []models.Account `cbor:"accounts" json:"accounts"`
}
