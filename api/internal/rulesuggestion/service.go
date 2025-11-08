package rulesuggestion

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/shared"
)

const (
	modelName = shared.ChatModelGPT5
)

// Service provides rule suggestion functionality using GPT-5-mini
type Service struct {
	client     *openai.Client
	configured bool
}

type serviceOptions struct {
	requestOptions []option.RequestOption
}

// Option configures the Service
type Option func(*serviceOptions)

// WithHTTPClient injects a custom *http.Client (handy for tests)
func WithHTTPClient(client *http.Client) Option {
	return func(opts *serviceOptions) {
		if client != nil {
			opts.requestOptions = append(opts.requestOptions, option.WithHTTPClient(client))
		}
	}
}

// WithBaseURL overrides the OpenAI base URL
func WithBaseURL(baseURL string) Option {
	return func(opts *serviceOptions) {
		if baseURL != "" {
			opts.requestOptions = append(opts.requestOptions, option.WithBaseURL(baseURL))
		}
	}
}

// NewService constructs a new rule suggestion Service
func NewService(apiKey string, opts ...Option) *Service {
	cleanKey := strings.TrimSpace(apiKey)

	var cfg serviceOptions
	for _, opt := range opts {
		opt(&cfg)
	}

	requestOptions := make([]option.RequestOption, 0, len(cfg.requestOptions)+1)
	if cleanKey != "" {
		requestOptions = append(requestOptions, option.WithAPIKey(cleanKey))
	}
	requestOptions = append(requestOptions, cfg.requestOptions...)

	client := openai.NewClient(requestOptions...)

	return &Service{
		client:     &client,
		configured: cleanKey != "",
	}
}

// AnalysisRequest contains transaction data and user context for analysis
type AnalysisRequest struct {
	UserID             string             `json:"user_id"`
	AnalysisPeriodDays int                `json:"analysis_period_days"` // Days to analyze (default: 30)
	SpendingInsights   []SpendingInsight  `json:"spending_insights"`
	ExistingRules      []ExistingRuleInfo `json:"existing_rules"`
	MonthlyIncome      float64            `json:"monthly_income,omitempty"` // Optional for context
	SavingsGoal        float64            `json:"savings_goal,omitempty"`   // Optional monthly target
}

// ExistingRuleInfo contains information about user's existing rules
type ExistingRuleInfo struct {
	Category       Category   `json:"category"`
	ActionType     ActionType `json:"action_type"`
	ActionValue    float64    `json:"action_value"`
	ExecutionCount int        `json:"execution_count"`
	TotalSaved     float64    `json:"total_saved"` // In dollars
}

// SuggestionResponse contains rule suggestions and analysis
type SuggestionResponse struct {
	Suggestions        []RuleSuggestion `json:"suggestions"`
	OverallAnalysis    string           `json:"overall_analysis"`
	PrioritySuggestion string           `json:"priority_suggestion"`
}

// SuggestRules analyzes transaction patterns and suggests savings rules
func (s *Service) SuggestRules(ctx context.Context, req *AnalysisRequest) (*SuggestionResponse, error) {
	if !s.configured {
		return nil, errors.New("rule suggestion service is not configured")
	}

	if req == nil {
		return nil, errors.New("request cannot be nil")
	}

	// Default analysis period
	if req.AnalysisPeriodDays == 0 {
		req.AnalysisPeriodDays = 30
	}

	completion, err := s.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: modelName,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt),
			openai.UserMessage(buildUserPrompt(req)),
		},
		ReasoningEffort: shared.ReasoningEffortMinimal,
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &shared.ResponseFormatJSONSchemaParam{
				JSONSchema: shared.ResponseFormatJSONSchemaJSONSchemaParam{
					Name:   "rule_suggestions",
					Strict: openai.Bool(true),
					Schema: suggestionSchema(),
				},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("openai chat completion failed: %w", err)
	}

	if len(completion.Choices) == 0 {
		return nil, errors.New("openai completion returned no choices")
	}

	content := strings.TrimSpace(completion.Choices[0].Message.Content)
	result, err := parseSuggestionResponse(content)
	if err != nil {
		return nil, err
	}

	return result, nil
}

const systemPrompt = `You are a financial advisor AI specializing in creating personalized savings rules to help users prevent overspending and build savings habits.

Your task is to analyze transaction patterns and suggest effective savings rules. Consider:
- Transaction frequency and amounts by category
- Existing rules to avoid duplication
- User's financial context (income, goals)
- Realistic savings without causing financial stress
- Behavioral psychology of savings (start gentle, build habits)

Rule Types:
1. "multiply" - Save a multiple of transaction amount (e.g., 1.5x = save 50% extra)
2. "fixed" - Save a fixed amount per transaction (e.g., $5 per coffee purchase)

Impact Levels:
- "gentle": Low savings rate, easy to maintain (5-15%)
- "moderate": Balanced approach (15-30%)
- "aggressive": Higher savings for users ready to cut back (30%+)

Provide 3-5 actionable rule suggestions with clear reasoning and realistic estimated savings.
Focus on categories with high spending or frequent transactions for maximum impact.`

func buildUserPrompt(req *AnalysisRequest) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Analyze spending over the last %d days and suggest savings rules.\n\n", req.AnalysisPeriodDays))

	// Spending insights
	sb.WriteString("## Spending Breakdown\n")
	totalSpent := 0.0
	for _, insight := range req.SpendingInsights {
		totalSpent += insight.TotalSpent
		sb.WriteString(fmt.Sprintf("- %s: $%.2f across %d transactions (avg: $%.2f, trend: %s)\n",
			insight.Category, insight.TotalSpent, insight.TransactionCount,
			insight.AverageTransaction, insight.MonthlyTrend))
	}
	sb.WriteString(fmt.Sprintf("Total spent: $%.2f\n\n", totalSpent))

	// Existing rules
	if len(req.ExistingRules) > 0 {
		sb.WriteString("## Existing Rules\n")
		for _, rule := range req.ExistingRules {
			sb.WriteString(fmt.Sprintf("- %s: %s %.2f (executed %d times, saved $%.2f)\n",
				rule.Category, rule.ActionType, rule.ActionValue, rule.ExecutionCount, rule.TotalSaved))
		}
		sb.WriteString("\n")
	}

	// User context
	if req.MonthlyIncome > 0 {
		sb.WriteString(fmt.Sprintf("Monthly income: $%.2f\n", req.MonthlyIncome))
	}
	if req.SavingsGoal > 0 {
		sb.WriteString(fmt.Sprintf("Monthly savings goal: $%.2f\n", req.SavingsGoal))
	}

	sb.WriteString("\nSuggest 3-5 effective savings rules. Consider:\n")
	sb.WriteString("- Categories with high spending or frequency\n")
	sb.WriteString("- Rules that don't duplicate existing ones\n")
	sb.WriteString("- Mix of gentle, moderate, and aggressive options\n")
	sb.WriteString("- Estimated monthly savings for each rule\n")
	sb.WriteString("- Clear reasoning and actionable advice\n")

	return sb.String()
}

func suggestionSchema() map[string]any {
	return map[string]any{
		"type":                 "object",
		"additionalProperties": false,
		"required": []string{
			"suggestions",
			"overall_analysis",
			"priority_suggestion",
		},
		"properties": map[string]any{
			"suggestions": map[string]any{
				"type": "array",
				"items": map[string]any{
					"type":                 "object",
					"additionalProperties": false,
					"required": []string{
						"name",
						"category",
						"action_type",
						"action_value",
						"min_amount_cents",
						"max_amount_cents",
						"estimated_savings",
						"confidence",
						"reasoning",
						"impact_level",
					},
					"properties": map[string]any{
						"name": map[string]any{
							"type":        "string",
							"description": "Human-friendly rule name",
						},
						"category": map[string]any{
							"type": "string",
							"enum": []string{
								"Dining", "Groceries", "Transport", "Shopping",
								"Subscriptions", "Entertainment", "Bills", "Misc",
							},
							"description": "Transaction category",
						},
						"action_type": map[string]any{
							"type":        "string",
							"enum":        []string{"multiply", "fixed"},
							"description": "Type of savings action",
						},
						"action_value": map[string]any{
							"type":        "number",
							"description": "Multiplier or fixed amount",
							"minimum":     0,
						},
						"min_amount_cents": map[string]any{
							"type":        []string{"integer", "null"},
							"description": "Optional minimum transaction amount filter in cents",
						},
						"max_amount_cents": map[string]any{
							"type":        []string{"integer", "null"},
							"description": "Optional maximum transaction amount filter in cents",
						},
						"estimated_savings": map[string]any{
							"type":        "number",
							"description": "Projected monthly savings in dollars",
							"minimum":     0,
						},
						"confidence": map[string]any{
							"type":        "string",
							"enum":        []string{"high", "medium", "low"},
							"description": "Confidence in the suggestion",
						},
						"reasoning": map[string]any{
							"type":        "string",
							"description": "Explanation for this suggestion",
						},
						"impact_level": map[string]any{
							"type":        "string",
							"enum":        []string{"gentle", "moderate", "aggressive"},
							"description": "Savings intensity level",
						},
					},
				},
				"minItems": 3,
				"maxItems": 5,
			},
			"overall_analysis": map[string]any{
				"type":        "string",
				"description": "Overall spending analysis and recommendations",
			},
			"priority_suggestion": map[string]any{
				"type":        "string",
				"description": "Which suggestion to implement first and why",
			},
		},
	}
}

func parseSuggestionResponse(raw string) (*SuggestionResponse, error) {
	clean := strings.TrimSpace(raw)
	if clean == "" {
		return nil, errors.New("openai completion response empty")
	}

	result := &SuggestionResponse{}
	if err := json.Unmarshal([]byte(clean), result); err == nil {
		return result, nil
	}

	// Try to extract JSON from markdown code blocks
	start := strings.Index(clean, "{")
	end := strings.LastIndex(clean, "}")
	if start == -1 || end == -1 || end <= start {
		return nil, fmt.Errorf("failed to parse suggestion response: %s", raw)
	}

	if err := json.Unmarshal([]byte(clean[start:end+1]), result); err != nil {
		return nil, fmt.Errorf("failed to decode suggestion JSON: %w", err)
	}

	return result, nil
}

// TransactionSummary represents a simple transaction for analysis
type TransactionSummary struct {
	Category string
	Amount   int64 // in cents
	Date     time.Time
}

// AnalyzeTransactionsAndSuggestRules is a helper that takes raw transactions and generates suggestions
func (s *Service) AnalyzeTransactionsAndSuggestRules(
	ctx context.Context,
	userID string,
	transactions []TransactionSummary,
	existingRules []ExistingRuleInfo,
	monthlyIncome float64,
	savingsGoal float64,
) (*SuggestionResponse, error) {

	// Calculate spending insights from transactions
	insights := calculateSpendingInsights(transactions)

	// Build analysis request
	req := &AnalysisRequest{
		UserID:             userID,
		AnalysisPeriodDays: 30,
		SpendingInsights:   insights,
		ExistingRules:      existingRules,
		MonthlyIncome:      monthlyIncome,
		SavingsGoal:        savingsGoal,
	}

	return s.SuggestRules(ctx, req)
}

// calculateSpendingInsights aggregates transaction data into spending insights
func calculateSpendingInsights(transactions []TransactionSummary) []SpendingInsight {
	categoryMap := make(map[string]*SpendingInsight)

	for _, tx := range transactions {
		if _, exists := categoryMap[tx.Category]; !exists {
			categoryMap[tx.Category] = &SpendingInsight{
				Category:           Category(tx.Category),
				HighestTransaction: float64(tx.Amount) / 100,
				LowestTransaction:  float64(tx.Amount) / 100,
			}
		}

		insight := categoryMap[tx.Category]
		amount := float64(tx.Amount) / 100

		insight.TotalSpent += amount
		insight.TransactionCount++

		if amount > insight.HighestTransaction {
			insight.HighestTransaction = amount
		}
		if amount < insight.LowestTransaction {
			insight.LowestTransaction = amount
		}
	}

	// Calculate averages and format results
	insights := make([]SpendingInsight, 0, len(categoryMap))
	for _, insight := range categoryMap {
		if insight.TransactionCount > 0 {
			insight.AverageTransaction = insight.TotalSpent / float64(insight.TransactionCount)
		}
		insight.MonthlyTrend = "stable" // Simplified - would need historical data for real trends
		insights = append(insights, *insight)
	}

	return insights
}
