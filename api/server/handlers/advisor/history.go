package advisor

import (
	"context"
	"fmt"
	"math"
	"strings"

	advisorsvc "regulation/internal/advisor"
	"regulation/internal/ent"
	"regulation/internal/ent/account"
	"regulation/internal/ent/user"
	"regulation/internal/protocol"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

const (
	defaultHistoryCurrency        = "KRW"
	maxHistoryTransactionsPerAcct = 25
	maxHistoryExpenses            = 50
)

// BudgetPlanFromHistoryRequest asks the server to build a budget plan using
// persisted user history.
type BudgetPlanFromHistoryRequest struct {
	Username                string  `cbor:"username"`
	Currency                string  `cbor:"currency"`
	TargetGoal              string  `cbor:"target_goal"`
	DesiredReductionPercent float64 `cbor:"desired_reduction_percent"`
}

// Validate implements the ro.Validatable contract.
func (r *BudgetPlanFromHistoryRequest) Validate() error {
	if strings.TrimSpace(r.Username) == "" {
		return protocol.ErrorResponse{
			Code:    protocol.InvalidRequest,
			Message: "username is required",
		}
	}

	if r.DesiredReductionPercent < 0 {
		return protocol.ErrorResponse{
			Code:    protocol.InvalidRequest,
			Message: "desired_reduction_percent cannot be negative",
		}
	}

	return nil
}

// BudgetPlanFromHistoryResponse contains the derived plan plus the expenses
// that were sent to the advisor model.
type BudgetPlanFromHistoryResponse struct {
	Username string                         `cbor:"username"`
	Expenses []advisorsvc.BudgetExpense     `cbor:"expenses"`
	Plan     *advisorsvc.BudgetPlanResponse `cbor:"plan"`
}

// BudgetPlanFromHistory looks up a user by username, builds a spending
// snapshot from their account history, and feeds it into the advisor service.
func (h *Handler) BudgetPlanFromHistory(ctx fiber.Ctx, req *BudgetPlanFromHistoryRequest) (*BudgetPlanFromHistoryResponse, error) {
	if h.db == nil {
		return nil, protocol.ErrorResponse{
			Code:    protocol.InternalError,
			Message: "database client not configured",
		}
	}

	if h.service == nil {
		return nil, protocol.ErrorResponse{
			Code:    protocol.InternalError,
			Message: "budget advisor is currently unavailable",
		}
	}

	username := strings.TrimSpace(req.Username)
	userEntity, err := h.db.User.Query().
		Where(user.NicknameEQ(username)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, protocol.ErrorResponse{
				Code:    protocol.InvalidRequest,
				Message: fmt.Sprintf("user %q not found", username),
			}
		}
		return nil, fmt.Errorf("failed to load user %q: %w", username, err)
	}

	expenses, err := h.expensesFromHistory(ctx, userEntity.ID)
	if err != nil {
		return nil, err
	}

	currency := strings.TrimSpace(req.Currency)
	if currency == "" {
		currency = defaultHistoryCurrency
	}

	metadata := map[string]string{
		"username": username,
		"user_id":  userEntity.ID.String(),
	}

	planReq := &advisorsvc.BudgetPlanRequest{
		Currency:                currency,
		Expenses:                expenses,
		TargetGoal:              req.TargetGoal,
		DesiredReductionPercent: req.DesiredReductionPercent,
		Metadata:                metadata,
	}

	plan, err := h.service.GeneratePlan(ctx, planReq)
	if err != nil {
		return nil, err
	}

	return &BudgetPlanFromHistoryResponse{
		Username: username,
		Expenses: expenses,
		Plan:     plan,
	}, nil
}

func (h *Handler) expensesFromHistory(ctx context.Context, userID uuid.UUID) ([]advisorsvc.BudgetExpense, error) {
	accounts, err := h.db.Account.Query().
		Where(account.UserIDEQ(userID)).
		WithTransactions(func(q *ent.TransactionQuery) {
			q.Limit(maxHistoryTransactionsPerAcct)
		}).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load accounts: %w", err)
	}

	var expenses []advisorsvc.BudgetExpense
	for _, acct := range accounts {
		for _, tx := range acct.Edges.Transactions {
			amount := centsToFloat(tx.Amount)
			if amount <= 0 {
				continue
			}

			expenses = append(expenses, advisorsvc.BudgetExpense{
				Category: acct.Name,
				Amount:   amount,
				Notes:    fmt.Sprintf("account type: %s", acct.Type),
			})

			if len(expenses) >= maxHistoryExpenses {
				return expenses, nil
			}
		}
	}

	if len(expenses) == 0 {
		return nil, protocol.ErrorResponse{
			Code:    protocol.InvalidRequest,
			Message: "no transaction history found for user",
		}
	}

	return expenses, nil
}

func centsToFloat(cents int64) float64 {
	return math.Abs(float64(cents)) / 100
}
