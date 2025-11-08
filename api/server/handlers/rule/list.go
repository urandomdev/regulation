package rule

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"regulation/internal/ent"
	entrule "regulation/internal/ent/rule"
	"regulation/internal/protocol"
	"regulation/server/services/request_context"
)

// ListRules retrieves all rules for the current user
// @Route GET /rules
func (h *Handler) ListRules(ctx fiber.Ctx) (*ListRulesResponse, error) {
	session := request_context.Session(ctx)

	// Query all rules for user, ordered by priority
	rules, err := h.db.Rule.
		Query().
		Where(entrule.UserID(session.UserID)).
		Order(ent.Asc(entrule.FieldPriority)).
		All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to list rules: %w", err)
	}

	// Convert to response format
	response := make([]RuleResponse, len(rules))
	for i, rule := range rules {
		response[i] = RuleResponse{
			ID:              rule.ID,
			Name:            rule.Name,
			Category:        string(rule.Category),
			MinAmountCents:  rule.MinAmountCents,
			MaxAmountCents:  rule.MaxAmountCents,
			ActionType:      string(rule.ActionType),
			ActionValue:     rule.ActionValue,
			TargetAccountID: rule.TargetAccountID,
			Priority:        rule.Priority,
			IsActive:        rule.IsActive,
			ExecutionCount:  rule.ExecutionCount,
			TotalSavedCents: rule.TotalSavedCents,
			CreatedAt:       rule.CreatedAt,
			UpdatedAt:       rule.UpdatedAt,
		}
	}

	return &ListRulesResponse{
		Rules: response,
		Total: len(response),
	}, nil
}

// GetRule retrieves a specific rule by ID
// @Route GET /rules/:id
func (h *Handler) GetRule(ctx fiber.Ctx) (*RuleResponse, error) {
	session := request_context.Session(ctx)

	// Parse rule ID from path
	ruleID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return nil, protocol.ErrorResponse{
			Code:    protocol.InvalidParametersError,
			Message: "invalid rule ID format",
		}
	}

	// Query rule, ensuring it belongs to the user
	rule, err := h.db.Rule.
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
		return nil, fmt.Errorf("failed to get rule: %w", err)
	}

	return &RuleResponse{
		ID:              rule.ID,
		Name:            rule.Name,
		Category:        string(rule.Category),
		MinAmountCents:  rule.MinAmountCents,
		MaxAmountCents:  rule.MaxAmountCents,
		ActionType:      string(rule.ActionType),
		ActionValue:     rule.ActionValue,
		TargetAccountID: rule.TargetAccountID,
		Priority:        rule.Priority,
		IsActive:        rule.IsActive,
		ExecutionCount:  rule.ExecutionCount,
		TotalSavedCents: rule.TotalSavedCents,
		CreatedAt:       rule.CreatedAt,
		UpdatedAt:       rule.UpdatedAt,
	}, nil
}

type ListRulesResponse struct {
	Rules []RuleResponse `cbor:"rules" json:"rules"`
	Total int            `cbor:"total" json:"total"`
}

type RuleResponse struct {
	ID              uuid.UUID `cbor:"id" json:"id"`
	Name            string    `cbor:"name" json:"name"`
	Category        string    `cbor:"category" json:"category"`
	MinAmountCents  *int64    `cbor:"min_amount_cents,omitempty" json:"min_amount_cents,omitempty"`
	MaxAmountCents  *int64    `cbor:"max_amount_cents,omitempty" json:"max_amount_cents,omitempty"`
	ActionType      string    `cbor:"action_type" json:"action_type"`
	ActionValue     float64   `cbor:"action_value" json:"action_value"`
	TargetAccountID uuid.UUID `cbor:"target_account_id" json:"target_account_id"`
	Priority        int       `cbor:"priority" json:"priority"`
	IsActive        bool      `cbor:"is_active" json:"is_active"`
	ExecutionCount  int       `cbor:"execution_count" json:"execution_count"`
	TotalSavedCents int64     `cbor:"total_saved_cents" json:"total_saved_cents"`
	CreatedAt       any       `cbor:"created_at" json:"created_at"`
	UpdatedAt       any       `cbor:"updated_at" json:"updated_at"`
}
