package rule

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"regulation/internal/ent"
	entaccount "regulation/internal/ent/account"
	entrule "regulation/internal/ent/rule"
	"regulation/internal/protocol"
	"regulation/server/services/request_context"
)

// CreateRule creates a new savings rule
// @Route POST /rules
func (h *Handler) CreateRule(ctx fiber.Ctx, req *CreateRuleRequest) (*CreateRuleResponse, error) {
	session := request_context.Session(ctx)

	// Validate that target account exists and belongs to user
	targetAccount, err := h.db.Account.
		Query().
		Where(
			entaccount.ID(req.TargetAccountID),
			entaccount.UserID(session.UserID),
			entaccount.IsActive(true),
		).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, protocol.ErrorResponse{
				Code:    protocol.NotFoundError,
				Message: "target account not found or not owned by user",
			}
		}
		return nil, fmt.Errorf("failed to validate target account: %w", err)
	}

	// Validate that target account is a savings account
	if targetAccount.Type != entaccount.TypeSavings {
		return nil, protocol.ErrorResponse{
			Code:    protocol.InvalidParametersError,
			Message: "target account must be a savings account",
		}
	}

	// Create the rule
	rule, err := h.db.Rule.
		Create().
		SetUserID(session.UserID).
		SetName(req.Name).
		SetCategory(req.Category).
		SetNillableMinAmountCents(req.MinAmountCents).
		SetNillableMaxAmountCents(req.MaxAmountCents).
		SetActionType(req.ActionType).
		SetActionValue(req.ActionValue).
		SetTargetAccountID(req.TargetAccountID).
		SetPriority(req.Priority).
		Save(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to create rule: %w", err)
	}

	return &CreateRuleResponse{
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

type CreateRuleRequest struct {
	Name            string             `cbor:"name" json:"name"`
	Category        entrule.Category   `cbor:"category" json:"category"`
	MinAmountCents  *int64             `cbor:"min_amount_cents,omitempty" json:"min_amount_cents,omitempty"`
	MaxAmountCents  *int64             `cbor:"max_amount_cents,omitempty" json:"max_amount_cents,omitempty"`
	ActionType      entrule.ActionType `cbor:"action_type" json:"action_type"`
	ActionValue     float64            `cbor:"action_value" json:"action_value"`
	TargetAccountID uuid.UUID          `cbor:"target_account_id" json:"target_account_id"`
	Priority        int                `cbor:"priority" json:"priority"`
}

func (r *CreateRuleRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Name, validation.Required, validation.RuneLength(1, 100)),
		validation.Field(&r.Category, validation.Required, validation.In(
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
		validation.Field(&r.ActionType, validation.Required, validation.In(
			entrule.ActionTypeMultiply,
			entrule.ActionTypeFixed,
		)),
		validation.Field(&r.ActionValue, validation.Required, validation.Min(0.0)),
		validation.Field(&r.TargetAccountID, validation.Required),
		validation.Field(&r.Priority, validation.Min(0)),
	)
}

type CreateRuleResponse struct {
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
