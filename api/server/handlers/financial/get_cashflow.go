package financial

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"

	"regulation/internal/ent/account"
	"regulation/internal/ent/transaction"
	"regulation/server/services/request_context"
)

// GetCashflow retrieves income, spend, and NET for a date range
// @Route POST /financial/cashflow
func (h *Handler) GetCashflow(ctx fiber.Ctx, req *GetCashflowRequest) (*GetCashflowResponse, error) {
	session := request_context.Session(ctx)

	// Set default date range to current month if not provided
	now := time.Now()
	start := req.Start
	end := req.End

	if start == nil {
		monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		start = &monthStart
	}

	if end == nil {
		monthEnd := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location()).Add(-time.Second)
		end = &monthEnd
	}

	// Get all transactions for user's accounts within date range
	transactions, err := h.db.Transaction.
		Query().
		Where(transaction.HasAccountWith(
			account.UserID(session.UserID),
		)).
		Where(transaction.DateGTE(*start)).
		Where(transaction.DateLTE(*end)).
		Where(transaction.Pending(false)). // Only count non-pending transactions
		All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to query transactions: %w", err)
	}

	// Calculate income and spend
	var totalIncome int64
	var totalSpend int64

	for _, txn := range transactions {
		if txn.Amount < 0 {
			// Negative amount = credit (income)
			totalIncome += -txn.Amount // Convert to positive
		} else {
			// Positive amount = debit (spend)
			totalSpend += txn.Amount
		}
	}

	// Calculate NET
	net := totalIncome - totalSpend

	return &GetCashflowResponse{
		TotalIncome: totalIncome,
		TotalSpend:  totalSpend,
		Net:         net,
		Start:       *start,
		End:         *end,
	}, nil
}

type GetCashflowRequest struct {
	Start *time.Time `cbor:"start,omitempty" json:"start,omitempty"`
	End   *time.Time `cbor:"end,omitempty" json:"end,omitempty"`
}

func (r *GetCashflowRequest) Validate() error {
	return nil // No validation needed for optional time fields
}

type GetCashflowResponse struct {
	TotalIncome int64     `cbor:"total_income" json:"total_income"`
	TotalSpend  int64     `cbor:"total_spend" json:"total_spend"`
	Net         int64     `cbor:"net" json:"net"`
	Start       time.Time `cbor:"start" json:"start"`
	End         time.Time `cbor:"end" json:"end"`
}
