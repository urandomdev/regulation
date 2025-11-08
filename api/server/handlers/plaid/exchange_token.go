package plaid

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"regulation/internal/ent/account"
	"regulation/server/handlers/models"
	"regulation/server/services/request_context"
)

// ExchangeToken exchanges a public token for an access token and creates Item and Account records
// @Route POST /plaid/exchange-token
func (h *Handler) ExchangeToken(ctx fiber.Ctx, req *ExchangeTokenRequest) error {
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
	defer func() {
		if v := recover(); v != nil {
			_ = tx.Rollback()
			panic(v)
		}
	}()

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
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Error().Err(rollbackErr).Msg("Failed to rollback transaction")
		}
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
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Msg("Failed to rollback transaction")
			}
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

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Trigger initial transaction sync in the background
	go func() {
		_, err := h.syncService.SyncItemTransactions(ctx, itemEntity.ID)
		if err != nil {
			log.Error().
				Err(err).
				Str("item_id", itemEntity.ID.String()).
				Msg("Failed to perform initial transaction sync")
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
