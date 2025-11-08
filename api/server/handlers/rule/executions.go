package rule

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"regulation/internal/ent"
	entrule "regulation/internal/ent/rule"
	entruleexecution "regulation/internal/ent/ruleexecution"
	"regulation/internal/protocol"
	"regulation/server/services/request_context"
)

// GetRuleExecutions retrieves execution history for a specific rule
// @Route GET /rules/:id/executions
func (h *Handler) GetRuleExecutions(ctx fiber.Ctx) (*RuleExecutionsResponse, error) {
	session := request_context.Session(ctx)

	// Parse rule ID from path
	ruleID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return nil, protocol.ErrorResponse{
			Code:    protocol.InvalidParametersError,
			Message: "invalid rule ID format",
		}
	}

	// Verify rule exists and belongs to user
	_, err = h.db.Rule.
		Query().
		Where(
			entrule.ID(ruleID),
			entrule.UserID(session.UserID),
		).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, protocol.ErrorResponse{
				Code:    protocol.NotFoundError,
				Message: "rule not found",
			}
		}
		return nil, fmt.Errorf("failed to verify rule: %w", err)
	}

	// Query executions with related transaction data
	executions, err := h.db.RuleExecution.
		Query().
		Where(entruleexecution.RuleID(ruleID)).
		WithTransaction().
		Order(ent.Desc(entruleexecution.FieldCreatedAt)).
		Limit(100). // Limit to most recent 100
		All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get rule executions: %w", err)
	}

	// Convert to response format
	response := make([]RuleExecutionResponse, len(executions))
	for i, exec := range executions {
		resp := RuleExecutionResponse{
			ID:              exec.ID,
			RuleID:          exec.RuleID,
			TransactionID:   exec.TransactionID,
			AmountCents:     exec.AmountCents,
			SourceAccountID: exec.SourceAccountID,
			TargetAccountID: exec.TargetAccountID,
			Status:          string(exec.Status),
			ErrorMessage:    &exec.ErrorMessage,
			CreatedAt:       exec.CreatedAt,
			CompletedAt:     exec.CompletedAt,
		}

		// Add transaction details if available
		if exec.Edges.Transaction != nil {
			tx := exec.Edges.Transaction
			resp.Transaction = &TransactionSummary{
				Name:         tx.Name,
				MerchantName: &tx.MerchantName,
				Amount:       tx.Amount,
				Date:         tx.Date,
				Category:     tx.Category,
			}
		}

		response[i] = resp
	}

	return &RuleExecutionsResponse{
		Executions: response,
		Total:      len(response),
	}, nil
}

type RuleExecutionsResponse struct {
	Executions []RuleExecutionResponse `cbor:"executions" json:"executions"`
	Total      int                     `cbor:"total" json:"total"`
}

type RuleExecutionResponse struct {
	ID              uuid.UUID           `cbor:"id" json:"id"`
	RuleID          uuid.UUID           `cbor:"rule_id" json:"rule_id"`
	TransactionID   uuid.UUID           `cbor:"transaction_id" json:"transaction_id"`
	AmountCents     int64               `cbor:"amount_cents" json:"amount_cents"`
	SourceAccountID uuid.UUID           `cbor:"source_account_id" json:"source_account_id"`
	TargetAccountID uuid.UUID           `cbor:"target_account_id" json:"target_account_id"`
	Status          string              `cbor:"status" json:"status"`
	ErrorMessage    *string             `cbor:"error_message,omitempty" json:"error_message,omitempty"`
	CreatedAt       any                 `cbor:"created_at" json:"created_at"`
	CompletedAt     any                 `cbor:"completed_at,omitempty" json:"completed_at,omitempty"`
	Transaction     *TransactionSummary `cbor:"transaction,omitempty" json:"transaction,omitempty"`
}

type TransactionSummary struct {
	Name         string  `cbor:"name" json:"name"`
	MerchantName *string `cbor:"merchant_name,omitempty" json:"merchant_name,omitempty"`
	Amount       int64   `cbor:"amount" json:"amount"`
	Date         any     `cbor:"date" json:"date"`
	Category     string  `cbor:"category" json:"category"`
}
