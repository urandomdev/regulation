package rulesuggestion

// Category represents transaction categories
type Category string

const (
	CategoryDining        Category = "Dining"
	CategoryGroceries     Category = "Groceries"
	CategoryTransport     Category = "Transport"
	CategoryShopping      Category = "Shopping"
	CategorySubscriptions Category = "Subscriptions"
	CategoryEntertainment Category = "Entertainment"
	CategoryBills         Category = "Bills"
	CategoryMisc          Category = "Misc"
)

// ActionType represents the type of savings action
type ActionType string

const (
	ActionMultiply ActionType = "multiply" // Multiply transaction amount
	ActionFixed    ActionType = "fixed"    // Fixed dollar amount
)

// RuleSuggestion represents a suggested savings rule
type RuleSuggestion struct {
	Name             string     `json:"name"`
	Category         Category   `json:"category"`
	ActionType       ActionType `json:"action_type"`
	ActionValue      float64    `json:"action_value"`
	MinAmountCents   *int64     `json:"min_amount_cents,omitempty"`
	MaxAmountCents   *int64     `json:"max_amount_cents,omitempty"`
	EstimatedSavings float64    `json:"estimated_savings"` // Projected monthly savings in dollars
	Confidence       string     `json:"confidence"`        // "high", "medium", "low"
	Reasoning        string     `json:"reasoning"`         // Why this rule is suggested
	ImpactLevel      string     `json:"impact_level"`      // "aggressive", "moderate", "gentle"
}

// SpendingInsight provides analysis of user spending patterns
type SpendingInsight struct {
	Category           Category `json:"category"`
	TotalSpent         float64  `json:"total_spent"` // Total in dollars
	TransactionCount   int      `json:"transaction_count"`
	AverageTransaction float64  `json:"average_transaction"`
	MonthlyTrend       string   `json:"monthly_trend"` // "increasing", "stable", "decreasing"
	HighestTransaction float64  `json:"highest_transaction"`
	LowestTransaction  float64  `json:"lowest_transaction"`
}

// AllCategories returns all available categories
func AllCategories() []Category {
	return []Category{
		CategoryDining,
		CategoryGroceries,
		CategoryTransport,
		CategoryShopping,
		CategorySubscriptions,
		CategoryEntertainment,
		CategoryBills,
		CategoryMisc,
	}
}
