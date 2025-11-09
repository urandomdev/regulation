package categorizer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/shared"
)

const (
	modelName          = shared.ChatModelGPT5Mini
	defaultWorkerCount = 5                      // Default number of concurrent workers
	maxWorkerCount     = 20                     // Maximum concurrent workers
	defaultMaxRetries  = 3                      // Default number of retry attempts
	retryDelay         = 500 * time.Millisecond // Delay between retries
)

// Service handles transaction categorization using GPT-5-mini
type Service struct {
	client      *openai.Client
	configured  bool
	workerCount int
	maxRetries  int
}

type serviceOptions struct {
	requestOptions []option.RequestOption
	workerCount    int
	maxRetries     int
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

// WithWorkerCount sets the number of concurrent workers for batch categorization
func WithWorkerCount(count int) Option {
	return func(opts *serviceOptions) {
		if count > 0 && count <= maxWorkerCount {
			opts.workerCount = count
		}
	}
}

// WithMaxRetries sets the maximum number of retry attempts for failed categorizations
func WithMaxRetries(retries int) Option {
	return func(opts *serviceOptions) {
		if retries >= 0 {
			opts.maxRetries = retries
		}
	}
}

// NewService constructs a new categorization Service
func NewService(apiKey string, opts ...Option) *Service {
	cleanKey := strings.TrimSpace(apiKey)

	var cfg serviceOptions
	cfg.workerCount = defaultWorkerCount // Set default
	cfg.maxRetries = defaultMaxRetries   // Set default
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
		client:      &client,
		configured:  cleanKey != "",
		workerCount: cfg.workerCount,
		maxRetries:  cfg.maxRetries,
	}
}

// CategorizationRequest contains transaction data for categorization
type CategorizationRequest struct {
	MerchantName    string   `json:"merchant_name"`
	TransactionName string   `json:"transaction_name"`
	Amount          float64  `json:"amount"`
	PlaidCategories []string `json:"plaid_categories,omitempty"`
}

// CategorizationResponse contains the categorization result
type CategorizationResponse struct {
	Category   CategoryType `json:"category"`
	Confidence string       `json:"confidence"`
	Reasoning  string       `json:"reasoning"`
}

// Categorize uses GPT-5-mini to categorize a transaction with retry logic
func (s *Service) Categorize(ctx context.Context, req *CategorizationRequest) (*CategorizationResponse, error) {
	if !s.configured {
		return nil, errors.New("categorizer service is not configured")
	}

	if req == nil {
		return nil, errors.New("request cannot be nil")
	}

	var lastErr error
	for attempt := 0; attempt <= s.maxRetries; attempt++ {
		// Add delay for retries (but not on first attempt)
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(retryDelay * time.Duration(attempt)):
				// Exponential backoff
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
						Name:   "transaction_category",
						Strict: openai.Bool(true),
						Schema: categorizationSchema(),
					},
				},
			},
		})
		if err != nil {
			lastErr = fmt.Errorf("openai chat completion failed: %w", err)
			continue // Retry
		}

		if len(completion.Choices) == 0 {
			lastErr = errors.New("openai completion returned no choices")
			continue // Retry
		}

		content := strings.TrimSpace(completion.Choices[0].Message.Content)
		result, err := parseCategorizationResponse(content)
		if err != nil {
			lastErr = err
			continue // Retry
		}

		return result, nil
	}

	return nil, fmt.Errorf("categorization failed after %d attempts: %w", s.maxRetries+1, lastErr)
}

const systemPrompt = `You are a financial transaction categorization assistant.
Your task is to accurately categorize transactions into one of the following categories:
- Dining: Restaurants, cafes, food delivery
- Groceries: Supermarkets, grocery stores, food stores
- Transport: Gas stations, parking, public transit, ride shares, vehicle expenses
- Shopping: Retail stores, clothing, electronics, general merchandise
- Subscriptions: Recurring services like streaming, music, software, memberships
- Entertainment: Movies, concerts, sporting events, recreation
- Bills: Utilities, internet, phone, insurance, rent, bank fees
- Transfer: Internal transfers between accounts (should be excluded from budgeting)
- Misc: Anything that doesn't fit the above categories

Analyze the merchant name, transaction description, amount, and any provided category hints to determine the most appropriate category.
Provide your confidence level (high, medium, low) and brief reasoning for your categorization.`

func buildUserPrompt(req *CategorizationRequest) string {
	var sb strings.Builder

	sb.WriteString("Categorize the following transaction:\n\n")

	if req.MerchantName != "" {
		sb.WriteString(fmt.Sprintf("Merchant: %s\n", req.MerchantName))
	}

	if req.TransactionName != "" {
		sb.WriteString(fmt.Sprintf("Description: %s\n", req.TransactionName))
	}

	sb.WriteString(fmt.Sprintf("Amount: $%.2f\n", req.Amount))

	if len(req.PlaidCategories) > 0 {
		sb.WriteString(fmt.Sprintf("Plaid Categories (hints): %s\n", strings.Join(req.PlaidCategories, " > ")))
	}

	sb.WriteString("\nReturn the category, confidence level, and your reasoning.")

	return sb.String()
}

func categorizationSchema() map[string]any {
	return map[string]any{
		"type":                 "object",
		"additionalProperties": false,
		"required": []string{
			"category",
			"confidence",
			"reasoning",
		},
		"properties": map[string]any{
			"category": map[string]any{
				"type": "string",
				"enum": []string{
					"Dining",
					"Groceries",
					"Transport",
					"Shopping",
					"Subscriptions",
					"Entertainment",
					"Bills",
					"Transfer",
					"Misc",
				},
				"description": "The transaction category",
			},
			"confidence": map[string]any{
				"type": "string",
				"enum": []string{
					"high",
					"medium",
					"low",
				},
				"description": "Confidence level in the categorization",
			},
			"reasoning": map[string]any{
				"type":        "string",
				"description": "Brief explanation for the categorization decision",
			},
		},
	}
}

func parseCategorizationResponse(raw string) (*CategorizationResponse, error) {
	clean := strings.TrimSpace(raw)
	if clean == "" {
		return nil, errors.New("openai completion response empty")
	}

	result := &CategorizationResponse{}
	if err := json.Unmarshal([]byte(clean), result); err == nil {
		return result, nil
	}

	// Try to extract JSON from markdown code blocks or other wrapping
	start := strings.Index(clean, "{")
	end := strings.LastIndex(clean, "}")
	if start == -1 || end == -1 || end <= start {
		return nil, fmt.Errorf("failed to parse categorization response: %s", raw)
	}

	if err := json.Unmarshal([]byte(clean[start:end+1]), result); err != nil {
		return nil, fmt.Errorf("failed to decode categorization JSON: %w", err)
	}

	return result, nil
}

// BatchRequest contains a categorization request with an index for ordering
type BatchRequest struct {
	Index   int
	Request *CategorizationRequest
}

// BatchResult contains a categorization result with the original index
type BatchResult struct {
	Index    int
	Response *CategorizationResponse
	Error    error
}

// CategorizeBatch categorizes multiple transactions concurrently using a worker pool
// Returns results in the same order as the input requests
func (s *Service) CategorizeBatch(ctx context.Context, requests []*CategorizationRequest) ([]*CategorizationResponse, []error) {
	if !s.configured {
		err := errors.New("categorizer service is not configured")
		errs := make([]error, len(requests))
		for i := range errs {
			errs[i] = err
		}
		return nil, errs
	}

	if len(requests) == 0 {
		return []*CategorizationResponse{}, []error{}
	}

	// Create channels for work distribution and result collection
	jobs := make(chan BatchRequest, len(requests))
	results := make(chan BatchResult, len(requests))

	// Start worker pool
	var wg sync.WaitGroup
	workerCount := s.workerCount
	if workerCount > len(requests) {
		workerCount = len(requests) // Don't create more workers than requests
	}

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go s.categorizationWorker(ctx, jobs, results, &wg)
	}

	// Send all jobs to the channel
	for i, req := range requests {
		jobs <- BatchRequest{
			Index:   i,
			Request: req,
		}
	}
	close(jobs)

	// Wait for all workers to complete in a separate goroutine
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	responses := make([]*CategorizationResponse, len(requests))
	errs := make([]error, len(requests))

	for result := range results {
		responses[result.Index] = result.Response
		errs[result.Index] = result.Error
	}

	return responses, errs
}

// categorizationWorker processes categorization requests from the jobs channel
func (s *Service) categorizationWorker(ctx context.Context, jobs <-chan BatchRequest, results chan<- BatchResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		// Check context cancellation
		select {
		case <-ctx.Done():
			results <- BatchResult{
				Index:    job.Index,
				Response: nil,
				Error:    ctx.Err(),
			}
			continue
		default:
		}

		// Categorize the transaction
		response, err := s.Categorize(ctx, job.Request)

		results <- BatchResult{
			Index:    job.Index,
			Response: response,
			Error:    err,
		}
	}
}
