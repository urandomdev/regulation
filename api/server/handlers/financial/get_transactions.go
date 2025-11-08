package financial

import (
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

// GetTransactions retrieves transactions for the authenticated user with filtering and pagination
// @Route POST /financial/transactions
func (h *Handler) GetTransactions(ctx fiber.Ctx, req *GetTransactionsRequest) (*GetTransactionsResponse, error) {
	session := request_context.Session(ctx)

	// Build query for user's transactions via their accounts
	query := h.db.Transaction.
		Query().
		Where(transaction.HasAccountWith(
			account.UserID(session.UserID),
		))

	// Apply optional date range filters
	if req.Start != nil {
		query = query.Where(transaction.DateGTE(*req.Start))
	}

	if req.End != nil {
		query = query.Where(transaction.DateLTE(*req.End))
	}

	// Apply account filter if provided
	if req.AccountID != nil {
		query = query.Where(transaction.AccountID(*req.AccountID))
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

	return &GetTransactionsResponse{
		Transactions: responses,
		Total:        total,
	}, nil
}

type GetTransactionsRequest struct {
	Start     *time.Time `cbor:"start,omitempty" json:"start,omitempty"`
	End       *time.Time `cbor:"end,omitempty" json:"end,omitempty"`
	AccountID *uuid.UUID `cbor:"account_id,omitempty" json:"account_id,omitempty"`
	Limit     int        `cbor:"limit" json:"limit"`
	Offset    int        `cbor:"offset" json:"offset"`
}

func (r *GetTransactionsRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Limit, validation.Min(1), validation.Max(500)),
		validation.Field(&r.Offset, validation.Min(0)),
	)
}

type GetTransactionsResponse struct {
	Transactions []models.Transaction `cbor:"transactions" json:"transactions"`
	Total        int                  `cbor:"total" json:"total"`
}
