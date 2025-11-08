package advisor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	openai "github.com/openai/openai-go/v3" // imported as openai
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/shared"

	"regulation/internal/protocol"
)

const (
	modelName = shared.ChatModelGPT5Mini
)

// Service wraps the logic to talk to the GPT-5-mini model.
type Service struct {
	client     *openai.Client
	configured bool
}

type serviceOptions struct {
	requestOptions []option.RequestOption
}

// Option configures the Service.
type Option func(*serviceOptions)

// WithHTTPClient injects a custom *http.Client (handy for tests).
func WithHTTPClient(client *http.Client) Option {
	return func(opts *serviceOptions) {
		if client != nil {
			opts.requestOptions = append(opts.requestOptions, option.WithHTTPClient(client))
		}
	}
}

// WithBaseURL overrides the OpenAI base URL.
func WithBaseURL(baseURL string) Option {
	return func(opts *serviceOptions) {
		if baseURL != "" {
			opts.requestOptions = append(opts.requestOptions, option.WithBaseURL(baseURL))
		}
	}
}

// NewService constructs a new Service.
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

// BudgetPlanRequest is passed from the handler into the service.
type BudgetPlanRequest struct {
	Currency                string            `json:"currency" cbor:"currency"`
	Expenses                []BudgetExpense   `json:"expenses" cbor:"expenses"`
	TargetGoal              string            `json:"target_goal,omitempty" cbor:"target_goal,omitempty"`
	DesiredReductionPercent float64           `json:"desired_reduction_percent,omitempty" cbor:"desired_reduction_percent,omitempty"`
	Metadata                map[string]string `json:"metadata,omitempty" cbor:"metadata,omitempty"`
}

// BudgetExpense represents a single expense line item.
type BudgetExpense struct {
	Category string  `json:"category" cbor:"category"`
	Amount   float64 `json:"amount" cbor:"amount"`
	Notes    string  `json:"notes,omitempty" cbor:"notes,omitempty"`
}

func (r *BudgetPlanRequest) Validate() error {
	if r == nil {
		return protocol.ErrorResponse{
			Code:    protocol.InvalidRequest,
			Message: "request body is required",
		}
	}

	if r.Currency == "" {
		return protocol.ErrorResponse{
			Code:    protocol.InvalidRequest,
			Message: "currency is required",
		}
	}

	if len(r.Expenses) == 0 {
		return protocol.ErrorResponse{
			Code:    protocol.InvalidRequest,
			Message: "at least one expense entry is required",
		}
	}

	for i, expense := range r.Expenses {
		if expense.Category == "" {
			return protocol.ErrorResponse{
				Code:    protocol.InvalidRequest,
				Message: fmt.Sprintf("expense[%d] category is required", i),
			}
		}

		if expense.Amount <= 0 {
			return protocol.ErrorResponse{
				Code:    protocol.InvalidRequest,
				Message: fmt.Sprintf("expense[%d] amount must be greater than zero", i),
			}
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

// BudgetPlanResponse mirrors the JSON schema enforced via GPT.
type BudgetPlanResponse struct {
	RecommendedReductionAmount float64                  `json:"recommended_reduction_amount" cbor:"recommended_reduction_amount"`
	GoalStatement              string                   `json:"goal_statement" cbor:"goal_statement"`
	MonthlyTarget              float64                  `json:"monthly_target" cbor:"monthly_target"`
	ConfidenceNote             string                   `json:"confidence_note,omitempty" cbor:"confidence_note,omitempty"`
	Categories                 []CategoryRecommendation `json:"categories" cbor:"categories"`
}

// CategoryRecommendation breaks the reduction down per category.
type CategoryRecommendation struct {
	Category string  `json:"category" cbor:"category"`
	ReduceBy float64 `json:"reduce_by" cbor:"reduce_by"`
}

// GeneratePlan orchestrates the LLM call.
func (s *Service) GeneratePlan(ctx context.Context, req *BudgetPlanRequest) (*BudgetPlanResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	if !s.configured {
		return nil, protocol.ErrorResponse{
			Code:    protocol.InternalError,
			Message: "budget advisor is not configured",
		}
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
					Name:   "budget_plan",
					Strict: openai.Bool(true),
					Schema: planSchema(),
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
	plan, err := parsePlan(content)
	if err != nil {
		return nil, err
	}

	return plan, nil
}

const systemPrompt = `You are a financial coach helping Korean customers adjust next month's spending using last month's ledger.
Always respond with actionable, empathetic explanations and realistic reduction targets.
Respect the currency provided by the user and stay within achievable percentages (10-40% reduction unless explicitly told otherwise).
Only include fields defined in the JSON schema. Do not add additional fields like actions, tracking_hint, or primary_goal.`

func buildUserPrompt(req *BudgetPlanRequest) string {
	var sb strings.Builder
	total := 0.0

	sb.WriteString("Analyze the following expense list from last month.\n")
	sb.WriteString("Provide concrete reduction amount and a motivating goal for next month.\n")
	sb.WriteString("Expenses:\n")

	for _, expense := range req.Expenses {
		total += expense.Amount
		sb.WriteString(fmt.Sprintf("- %s: %.2f %s", expense.Category, expense.Amount, req.Currency))
		if expense.Notes != "" {
			sb.WriteString(fmt.Sprintf(" (%s)", expense.Notes))
		}
		sb.WriteString("\n")
	}

	sb.WriteString(fmt.Sprintf("Total spent: %.2f %s\n", total, req.Currency))

	if req.TargetGoal != "" {
		sb.WriteString("User goal: " + req.TargetGoal + "\n")
	}

	if req.DesiredReductionPercent > 0 {
		sb.WriteString(fmt.Sprintf("Preferred reduction percent: %.1f%%\n", req.DesiredReductionPercent))
	}

	if len(req.Metadata) > 0 {
		sb.WriteString("Metadata:\n")
		for key, value := range req.Metadata {
			sb.WriteString(fmt.Sprintf("- %s: %s\n", key, value))
		}
	}

	sb.WriteString("Return the JSON schema that was provided exactly.")

	return sb.String()
}

func planSchema() map[string]any {
	return map[string]any{
		"type":                 "object",
		"additionalProperties": false,
		"required": []string{
			"recommended_reduction_amount",
			"goal_statement",
			"monthly_target",
			"confidence_note",
			"categories",
		},
		"properties": map[string]any{
			"recommended_reduction_amount": map[string]any{
				"type":        "number",
				"description": "How much to reduce next month's total spend (numeric value, same currency as input).",
				"minimum":     0,
			},
			"goal_statement": map[string]any{
				"type":        "string",
				"description": "Short motivational statement describing the savings goal.",
			},
			"monthly_target": map[string]any{
				"type":        "number",
				"description": "Suggested target total spend for next month.",
				"minimum":     0,
			},
			"confidence_note": map[string]any{
				"type":        "string",
				"description": "Optional context about how realistic the plan is.",
			},
			"categories": map[string]any{
				"type": "array",
				"items": map[string]any{
					"type":                 "object",
					"required":             []string{"category", "reduce_by"},
					"additionalProperties": false,
					"properties": map[string]any{
						"category": map[string]any{
							"type": "string",
						},
						"reduce_by": map[string]any{
							"type":    "number",
							"minimum": 0,
						},
					},
				},
				"minItems": 1,
			},
		},
	}
}

func parsePlan(raw string) (*BudgetPlanResponse, error) {
	clean := strings.TrimSpace(raw)
	if clean == "" {
		return nil, errors.New("openai completion response empty")
	}

	plan := &BudgetPlanResponse{}
	if err := json.Unmarshal([]byte(clean), plan); err == nil {
		return plan, nil
	}

	start := strings.Index(clean, "{")
	end := strings.LastIndex(clean, "}")
	if start == -1 || end == -1 || end <= start {
		return nil, fmt.Errorf("failed to parse plan: %s", raw)
	}

	if err := json.Unmarshal([]byte(clean[start:end+1]), plan); err != nil {
		return nil, fmt.Errorf("failed to decode plan JSON: %w", err)
	}

	return plan, nil
}
