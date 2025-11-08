package rule

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"regulation/internal/ent"
	entrule "regulation/internal/ent/rule"
	"regulation/internal/protocol"
	"regulation/server/services/request_context"
)

// UpdateRule updates an existing rule
// @Route PATCH /rules/:id
func (h *Handler) UpdateRule(ctx fiber.Ctx, req *UpdateRuleRequest) (*RuleResponse, error) {
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

	// Build update query
	update := h.db.Rule.UpdateOneID(ruleID)

	if req.Name != nil {
		update.SetName(*req.Name)
	}
	if req.Category != nil {
		update.SetCategory(*req.Category)
	}
	if req.MinAmountCents != nil {
		update.SetNillableMinAmountCents(req.MinAmountCents)
	}
	if req.MaxAmountCents != nil {
		update.SetNillableMaxAmountCents(req.MaxAmountCents)
	}
	if req.ActionType != nil {
		update.SetActionType(*req.ActionType)
	}
	if req.ActionValue != nil {
		update.SetActionValue(*req.ActionValue)
	}
	if req.Priority != nil {
		update.SetPriority(*req.Priority)
	}

	// Execute update
	rule, err := update.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to update rule: %w", err)
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

// ToggleRule toggles a rule's active status
// @Route PATCH /rules/:id/toggle
func (h *Handler) ToggleRule(ctx fiber.Ctx) (*RuleResponse, error) {
	session := request_context.Session(ctx)

	// Parse rule ID from path
	ruleID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return nil, protocol.ErrorResponse{
			Code:    protocol.InvalidParametersError,
			Message: "invalid rule ID format",
		}
	}

	// Get current rule
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

	// Toggle active status
	rule, err = h.db.Rule.
		UpdateOneID(ruleID).
		SetIsActive(!rule.IsActive).
		Save(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to toggle rule: %w", err)
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

type UpdateRuleRequest struct {
	Name           *string             `cbor:"name,omitempty" json:"name,omitempty"`
	Category       *entrule.Category   `cbor:"category,omitempty" json:"category,omitempty"`
	MinAmountCents *int64              `cbor:"min_amount_cents,omitempty" json:"min_amount_cents,omitempty"`
	MaxAmountCents *int64              `cbor:"max_amount_cents,omitempty" json:"max_amount_cents,omitempty"`
	ActionType     *entrule.ActionType `cbor:"action_type,omitempty" json:"action_type,omitempty"`
	ActionValue    *float64            `cbor:"action_value,omitempty" json:"action_value,omitempty"`
	Priority       *int                `cbor:"priority,omitempty" json:"priority,omitempty"`
}

func (r *UpdateRuleRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Name, validation.RuneLength(1, 100)),
		validation.Field(&r.Category, validation.In(
			entrule.CategoryDining,
			entrule.CategoryGroceries,
			entrule.CategoryTransport,
			entrule.CategoryShopping,
			entrule.CategorySubscriptions,
			entrule.CategoryEntertainment,
			entrule.CategoryBills,
			entrule.CategoryMisc,
		)),
		validation.Field(&r.MinAmountCents, validation.Min(0)),
		validation.Field(&r.MaxAmountCents, validation.Min(0)),
		validation.Field(&r.ActionType, validation.In(
			entrule.ActionTypeMultiply,
			entrule.ActionTypeFixed,
		)),
		validation.Field(&r.ActionValue, validation.Min(0.0)),
		validation.Field(&r.Priority, validation.Min(0)),
	)
}
