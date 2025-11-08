# Rule Suggestion Module

This module provides AI-powered savings rule suggestions by analyzing transaction patterns and user spending behavior using GPT-5-mini.

## Overview

The rule suggestion service analyzes a user's transaction history, identifies spending patterns, and suggests personalized savings rules to help prevent overspending and build healthy savings habits.

## Purpose

**Problem:** Users may not know which savings rules to create or how aggressive their savings should be.

**Solution:** AI analyzes spending patterns and suggests 3-5 actionable rules based on:
- Transaction frequency and amounts by category
- Existing rules to avoid duplication
- User's financial context (income, savings goals)
- Behavioral psychology of habit formation
- Realistic savings targets without financial stress

## Architecture

### Files

- **types.go** - Data structures for suggestions, insights, and rule definitions
- **service.go** - GPT-5-mini powered suggestion service

### Rule Types

1. **Multiply Action** - Save a percentage of each transaction
   - Example: `1.5x` means save 50% extra on each transaction
   - Good for: Variable spending categories (dining, shopping)

2. **Fixed Action** - Save a fixed amount per transaction
   - Example: `$5` per transaction
   - Good for: Frequent small purchases (coffee, snacks)

### Impact Levels

- **Gentle** (5-15%): Easy to maintain, builds habits
- **Moderate** (15-30%): Balanced savings approach
- **Aggressive** (30%+): Higher savings for motivated users

## Data Structures

### RuleSuggestion

```go
type RuleSuggestion struct {
    Name             string     // "Dining 2x Save"
    Category         Category   // "Dining", "Shopping", etc.
    ActionType       ActionType // "multiply" or "fixed"
    ActionValue      float64    // Multiplier or dollar amount
    MinAmountCents   *int64     // Optional minimum filter
    MaxAmountCents   *int64     // Optional maximum filter
    EstimatedSavings float64    // Projected monthly savings
    Confidence       string     // "high", "medium", "low"
    Reasoning        string     // Why this rule is suggested
    ImpactLevel      string     // "gentle", "moderate", "aggressive"
}
```

### SpendingInsight

```go
type SpendingInsight struct {
    Category           Category
    TotalSpent         float64
    TransactionCount   int
    AverageTransaction float64
    MonthlyTrend       string  // "increasing", "stable", "decreasing"
    HighestTransaction float64
    LowestTransaction  float64
}
```

## Usage

### Basic Usage

```go
import "regulation/internal/rulesuggestion"

// Initialize service
service := rulesuggestion.NewService(openaiAPIKey)

// Prepare transaction data
transactions := []rulesuggestion.TransactionSummary{
    {Category: "Dining", Amount: 1500, Date: time.Now()},      // $15.00
    {Category: "Dining", Amount: 2500, Date: time.Now()},      // $25.00
    {Category: "Shopping", Amount: 7500, Date: time.Now()},    // $75.00
    // ... more transactions
}

// Existing rules (to avoid duplication)
existingRules := []rulesuggestion.ExistingRuleInfo{
    {
        Category:       rulesuggestion.CategoryGroceries,
        ActionType:     rulesuggestion.ActionMultiply,
        ActionValue:    1.5,
        ExecutionCount: 20,
        TotalSaved:     150.00,
    },
}

// Get suggestions
ctx := context.Background()
suggestions, err := service.AnalyzeTransactionsAndSuggestRules(
    ctx,
    userID,
    transactions,
    existingRules,
    5000.00, // monthly income
    500.00,  // savings goal
)
```

### Advanced Usage

```go
// Manual analysis request with pre-calculated insights
insights := []rulesuggestion.SpendingInsight{
    {
        Category:           rulesuggestion.CategoryDining,
        TotalSpent:         450.00,
        TransactionCount:   30,
        AverageTransaction: 15.00,
        MonthlyTrend:       "increasing",
        HighestTransaction: 75.00,
        LowestTransaction:  5.00,
    },
    {
        Category:           rulesuggestion.CategoryTransport,
        TotalSpent:         320.00,
        TransactionCount:   20,
        AverageTransaction: 16.00,
        MonthlyTrend:       "stable",
        HighestTransaction: 45.00,
        LowestTransaction:  10.00,
    },
}

req := &rulesuggestion.AnalysisRequest{
    UserID:            userID,
    AnalysisPeriodDays: 30,
    SpendingInsights:  insights,
    ExistingRules:     existingRules,
    MonthlyIncome:     5000.00,
    SavingsGoal:       500.00,
}

response, err := service.SuggestRules(ctx, req)
```

### Response Structure

```go
type SuggestionResponse struct {
    Suggestions        []RuleSuggestion
    OverallAnalysis    string  // General spending analysis
    PrioritySuggestion string  // Which rule to implement first
}
```

## Example Output

```json
{
  "suggestions": [
    {
      "name": "Dining Gentle Save",
      "category": "Dining",
      "action_type": "multiply",
      "action_value": 1.2,
      "estimated_savings": 90.00,
      "confidence": "high",
      "reasoning": "You spend $450/month on dining with frequent transactions. A 20% savings rule would save ~$90/month without major lifestyle changes.",
      "impact_level": "gentle"
    },
    {
      "name": "Coffee Fixed Save",
      "category": "Dining",
      "action_type": "fixed",
      "action_value": 2.00,
      "max_amount_cents": 1000,
      "estimated_savings": 40.00,
      "confidence": "medium",
      "reasoning": "Save $2 on small purchases under $10. Based on 20 coffee/snack transactions per month.",
      "impact_level": "gentle"
    },
    {
      "name": "Transport Moderate Save",
      "category": "Transport",
      "action_type": "multiply",
      "action_value": 1.5,
      "estimated_savings": 160.00,
      "confidence": "high",
      "reasoning": "Transportation spending is stable at $320/month. A 50% savings rule could save $160/month.",
      "impact_level": "moderate"
    }
  ],
  "overall_analysis": "Your spending is concentrated in Dining ($450) and Transport ($320). You have good savings potential in both categories without major lifestyle changes.",
  "priority_suggestion": "Start with the 'Dining Gentle Save' rule. It's easy to maintain and will save ~$90/month, getting you 18% toward your $500 monthly goal."
}
```

## Integration with Existing Rule System

The suggestions map directly to the rule schema:

```go
// Convert suggestion to rule creation
rule := &ent.RuleCreate{
    Name:            suggestion.Name,
    Category:        string(suggestion.Category),
    ActionType:      string(suggestion.ActionType),
    ActionValue:     suggestion.ActionValue,
    MinAmountCents:  suggestion.MinAmountCents,
    MaxAmountCents:  suggestion.MaxAmountCents,
    TargetAccountID: userSavingsAccountID, // User's savings account
    IsActive:        true,
    Priority:        0,
}
```

## GPT Prompt Design

### System Prompt

The system prompt instructs GPT-5-mini to:
- Act as a financial advisor specializing in savings rules
- Consider behavioral psychology of habit formation
- Provide realistic suggestions without causing financial stress
- Mix gentle, moderate, and aggressive options
- Focus on high-impact categories

### User Prompt

Includes:
- Spending breakdown by category with trends
- Existing rules to avoid duplication
- User's income and savings goals
- Clear request for 3-5 actionable suggestions

### Structured Output

Uses strict JSON schema to ensure:
- Valid categories and action types
- Positive values for amounts and savings
- Required fields (name, reasoning, confidence)
- 3-5 suggestions per response

## Use Cases

### 1. New User Onboarding
```go
// User just linked their first account
suggestions, err := service.AnalyzeTransactionsAndSuggestRules(
    ctx, userID, recentTransactions, []ExistingRuleInfo{}, 0, 0,
)
// Show suggestions in onboarding flow
```

### 2. Periodic Optimization
```go
// Monthly or quarterly review
// Analyze updated spending patterns
// Suggest new rules or adjustments to existing ones
```

### 3. Savings Goal Assistant
```go
// User sets savings goal but doesn't know how to reach it
// Service suggests specific rules to achieve the goal
```

### 4. Overspending Alert
```go
// Detect increased spending in a category
// Suggest protective rules to curb spending
```

## Performance Considerations

- **Latency**: ~1-2 seconds per suggestion request
- **Cost**: GPT-5-mini is cost-effective for complex analysis
- **Caching**: Consider caching suggestions for 24-48 hours
- **Batch Analysis**: Process multiple users in background jobs

## Configuration

Required in `config/config.json`:
```json
{
  "openai": {
    "api_key": "sk-..."
  }
}
```

## Best Practices

### 1. Data Quality
- Use at least 30 days of transaction data
- Filter out pending transactions
- Exclude internal transfers

### 2. Context Matters
- Always pass existing rules to avoid duplication
- Include income/goals when available
- Provide trend information for better suggestions

### 3. User Experience
- Present suggestions with clear reasoning
- Allow users to adjust parameters before creating
- Show estimated savings prominently
- Start with gentle rules for beginners

### 4. Testing
- Test with various spending patterns
- Verify suggestions don't exceed user's capacity
- Check edge cases (very low/high spending)

## Future Enhancements

1. **Historical Trend Analysis** - Track spending changes over time
2. **Rule Performance Tracking** - Learn from successful vs unsuccessful rules
3. **Personalization** - Adapt suggestions based on user behavior
4. **Goal-Based Suggestions** - Reverse engineer rules from savings goals
5. **Category Prioritization** - AI ranks categories by savings potential
6. **A/B Testing** - Test different suggestion strategies
7. **Notification Integration** - Alert users when new suggestions available
8. **Multi-Account Support** - Suggest rules across multiple accounts

## Example Integration Flow

```go
// 1. User requests suggestions
func (h *Handler) GetRuleSuggestions(ctx context.Context, userID uuid.UUID) (*SuggestionResponse, error) {
    // 2. Load user's transactions (last 30 days)
    transactions := h.loadUserTransactions(ctx, userID, 30)

    // 3. Load existing rules
    existingRules := h.loadExistingRules(ctx, userID)

    // 4. Get user context (optional)
    income, goal := h.getUserFinancialContext(ctx, userID)

    // 5. Generate suggestions
    suggestions, err := h.suggestionService.AnalyzeTransactionsAndSuggestRules(
        ctx, userID.String(), transactions, existingRules, income, goal,
    )

    // 6. Return to user
    return suggestions, err
}

// 7. User selects a suggestion and creates rule
func (h *Handler) CreateRuleFromSuggestion(ctx context.Context, userID uuid.UUID, suggestion RuleSuggestion) error {
    // Convert suggestion to rule
    // Validate user has savings account
    // Create rule in database
    // Return success
}
```

This module provides intelligent, personalized savings rule suggestions to help users build better financial habits!
