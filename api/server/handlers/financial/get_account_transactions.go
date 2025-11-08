package financial

import (
	"errors"
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"regulation/internal/ent"
	"regulation/internal/ent/account"
	"regulation/internal/ent/transaction"
	"regulation/server/handlers/models"
	"regulation/server/services/request_context"
)

// GetAccountTransactions retrieves transactions for a specific account
// @Route POST /financial/accounts/:id/transactions
func (h *Handler) GetAccountTransactions(ctx fiber.Ctx, req *GetAccountTransactionsRequest) (*GetAccountTransactionsResponse, error) {
	session := request_context.Session(ctx)

	// Get account ID from URL params
	accountIDStr := ctx.Params("id")
	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid account ID: %w", err)
	}

	// Verify the account belongs to the user
	acc, err := h.db.Account.
		Query().
		Where(account.ID(accountID)).
		Where(account.UserID(session.UserID)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("account not found")
		}
		return nil, fmt.Errorf("failed to query account: %w", err)
	}

	// Query transactions for this account
	query := h.db.Transaction.
		Query().
		Where(transaction.AccountID(acc.ID))

	// Apply optional date range filters
	if req.Start != nil {
		query = query.Where(transaction.DateGTE(*req.Start))
	}

	if req.End != nil {
		query = query.Where(transaction.DateLTE(*req.End))
	}

	// Get total count
	total, err := query.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count transactions: %w", err)
	}

	// Apply pagination and ordering
	transactions, err := query.
		Order(ent.Desc(transaction.FieldDate)).
		Limit(req.Limit).
		Offset(req.Offset).
		All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to query transactions: %w", err)
	}

	// Convert to response models
	responses := make([]models.Transaction, 0, len(transactions))
	for _, txn := range transactions {
		merchantPtr := &txn.MerchantName
		if txn.MerchantName == "" {
			merchantPtr = nil
		}

		channelPtr := &txn.PaymentChannel
		if txn.PaymentChannel == "" {
			channelPtr = nil
		}

		responses = append(responses, models.Transaction{
			ID:             txn.ID,
			AccountID:      txn.AccountID,
			Amount:         txn.Amount,
			Date:           txn.Date,
			Name:           txn.Name,
			MerchantName:   merchantPtr,
			Category:       txn.Category,
			Pending:        txn.Pending,
			PaymentChannel: channelPtr,
		})
	}

	return &GetAccountTransactionsResponse{
		Transactions: responses,
		Total:        total,
	}, nil
}

type GetAccountTransactionsRequest struct {
	Start  *time.Time `cbor:"start,omitempty" json:"start,omitempty"`
	End    *time.Time `cbor:"end,omitempty" json:"end,omitempty"`
	Limit  int        `cbor:"limit" json:"limit"`
	Offset int        `cbor:"offset" json:"offset"`
}

func (r *GetAccountTransactionsRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Limit, validation.Min(1), validation.Max(500)),
		validation.Field(&r.Offset, validation.Min(0)),
	)
}

type GetAccountTransactionsResponse struct {
	Transactions []models.Transaction `cbor:"transactions" json:"transactions"`
	Total        int                  `cbor:"total" json:"total"`
}
