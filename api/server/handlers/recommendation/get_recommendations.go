package recommendation

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"

	"regulation/internal/ent"
	"regulation/internal/ent/account"
	entrule "regulation/internal/ent/rule"
	"regulation/internal/ent/transaction"
	"regulation/internal/protocol"
	"regulation/internal/rulesuggestion"
	"regulation/server/services/request_context"
)

// RuleSuggestionResponse represents a suggested rule
type RuleSuggestionResponse struct {
	Name             string  `json:"name" cbor:"name"`
	Category         string  `json:"category" cbor:"category"`
	ActionType       string  `json:"action_type" cbor:"action_type"`
	ActionValue      float64 `json:"action_value" cbor:"action_value"`
	MinAmountCents   *int64  `json:"min_amount_cents,omitempty" cbor:"min_amount_cents,omitempty"`
	MaxAmountCents   *int64  `json:"max_amount_cents,omitempty" cbor:"max_amount_cents,omitempty"`
	EstimatedSavings float64 `json:"estimated_savings" cbor:"estimated_savings"`
	Confidence       string  `json:"confidence" cbor:"confidence"`
	Reasoning        string  `json:"reasoning" cbor:"reasoning"`
	ImpactLevel      string  `json:"impact_level" cbor:"impact_level"`
}

// GetRecommendationsResponse contains rule suggestions and analysis
type GetRecommendationsResponse struct {
	Suggestions        []RuleSuggestionResponse `json:"suggestions" cbor:"suggestions"`
	OverallAnalysis    string                   `json:"overall_analysis" cbor:"overall_analysis"`
	PrioritySuggestion string                   `json:"priority_suggestion" cbor:"priority_suggestion"`
	AnalysisPeriodDays int                      `json:"analysis_period_days" cbor:"analysis_period_days"`
	TransactionCount   int                      `json:"transaction_count" cbor:"transaction_count"`
	TotalSpent         float64                  `json:"total_spent" cbor:"total_spent"`
}

// GetRecommendations analyzes user's recent transactions and suggests savings rules
// @Route GET /recommendations
func (h *Handler) GetRecommendations(ctx fiber.Ctx) (*GetRecommendationsResponse, error) {
	session := request_context.Session(ctx)

	// 1. Define analysis period (last 30 days)
	analysisPeriodDays := 30
	startDate := time.Now().AddDate(0, 0, -analysisPeriodDays)

	// 2. Fetch user's transactions from the analysis period
	transactions, err := h.db.Transaction.
		Query().
		Where(
			transaction.HasAccountWith(
				account.UserID(session.UserID),
			),
			transaction.DateGTE(startDate),
			transaction.PendingEQ(false), // Only settled transactions
		).
		Order(ent.Desc(transaction.FieldDate)).
		All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch transactions: %w", err)
	}

	// Check if user has enough transaction data
	if len(transactions) < 5 {
		return nil, protocol.ErrorResponse{
			Code:    protocol.InvalidRequest,
			Message: fmt.Sprintf("insufficient transaction history. need at least 5 transactions, found %d", len(transactions)),
		}
	}

	// 3. Convert transactions to summary format
	txSummaries := make([]rulesuggestion.TransactionSummary, len(transactions))
	totalSpent := int64(0)
	for i, tx := range transactions {
		txSummaries[i] = rulesuggestion.TransactionSummary{
			Category: tx.Category,
			Amount:   tx.Amount,
			Date:     tx.Date,
		}
		totalSpent += tx.Amount
	}

	// 4. Fetch user's existing rules to avoid duplication
	existingRules, err := h.db.Rule.
		Query().
		Where(entrule.UserID(session.UserID)).
		All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch existing rules: %w", err)
	}

	// Convert existing rules to info format
	existingRuleInfos := make([]rulesuggestion.ExistingRuleInfo, len(existingRules))
	for i, rule := range existingRules {
		existingRuleInfos[i] = rulesuggestion.ExistingRuleInfo{
			Category:       rulesuggestion.Category(rule.Category),
			ActionType:     rulesuggestion.ActionType(rule.ActionType),
			ActionValue:    rule.ActionValue,
			ExecutionCount: rule.ExecutionCount,
			TotalSaved:     float64(rule.TotalSavedCents) / 100,
		}
	}

	// 5. Optionally get user's financial context
	// For now, we'll pass 0 for income and savings goal
	// These could be fetched from a user profile table if available
	monthlyIncome := 0.0
	savingsGoal := 0.0

	// 6. Call the suggestion service
	suggestions, err := h.suggestionService.AnalyzeTransactionsAndSuggestRules(
		ctx,
		session.UserID.String(),
		txSummaries,
		existingRuleInfos,
		monthlyIncome,
		savingsGoal,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to generate suggestions: %w", err)
	}

	// 7. Convert service response to handler response
	suggestionResponses := make([]RuleSuggestionResponse, len(suggestions.Suggestions))
	for i, suggestion := range suggestions.Suggestions {
		suggestionResponses[i] = RuleSuggestionResponse{
			Name:             suggestion.Name,
			Category:         string(suggestion.Category),
			ActionType:       string(suggestion.ActionType),
			ActionValue:      suggestion.ActionValue,
			MinAmountCents:   suggestion.MinAmountCents,
			MaxAmountCents:   suggestion.MaxAmountCents,
			EstimatedSavings: suggestion.EstimatedSavings,
			Confidence:       suggestion.Confidence,
			Reasoning:        suggestion.Reasoning,
			ImpactLevel:      suggestion.ImpactLevel,
		}
	}

	return &GetRecommendationsResponse{
		Suggestions:        suggestionResponses,
		OverallAnalysis:    suggestions.OverallAnalysis,
		PrioritySuggestion: suggestions.PrioritySuggestion,
		AnalysisPeriodDays: analysisPeriodDays,
		TransactionCount:   len(transactions),
		TotalSpent:         float64(totalSpent) / 100,
	}, nil
}
