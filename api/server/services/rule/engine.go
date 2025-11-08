package rule

import (
	"context"
	"fmt"
	"time"

	"github.com/DeltaLaboratory/contrib/hooks"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"regulation/internal/config"
	"regulation/internal/ent"
	entrule "regulation/internal/ent/rule"
	entruleexecution "regulation/internal/ent/ruleexecution"
	savingstransfer "regulation/internal/ent/savingstransfer"
	"regulation/server/services/notification"
)

// Engine handles rule evaluation and execution
type Engine struct {
	db                  *ent.Client
	notificationService *notification.Service
}

// NewEngine creates a new rule engine
func NewEngine(db *ent.Client, cfg *config.Config) *Engine {
	return &Engine{
		db:                  db,
		notificationService: notification.New(cfg, db),
	}
}

// ProcessTransaction evaluates and executes rules for a new transaction
// This is the main entry point called by the sync service
func (e *Engine) ProcessTransaction(ctx context.Context, transaction *ent.Transaction) error {
	// Skip pending transactions
	if transaction.Pending {
		return nil
	}

	// Skip transfers (they're not spending)
	if transaction.Category == "Transfer" {
		return nil
	}

	// Get the account to find user_id
	account, err := transaction.QueryAccount().Only(ctx)
	if err != nil {
		return fmt.Errorf("failed to get account: %w", err)
	}

	// Load all active rules for this user, ordered by priority
	rules, err := e.db.Rule.
		Query().
		Where(
			entrule.UserID(account.UserID),
			entrule.IsActive(true),
		).
		Order(ent.Asc(entrule.FieldPriority)).
		All(ctx)

	if err != nil {
		return fmt.Errorf("failed to load rules: %w", err)
	}

	// Evaluate each rule
	for _, rule := range rules {
		if e.matchesRule(transaction, rule) {
			if err := e.executeRule(ctx, transaction, rule, account); err != nil {
				log.Error().
					Err(err).
					Str("rule_id", rule.ID.String()).
					Str("transaction_id", transaction.ID.String()).
					Msg("Failed to execute rule")
				// Continue with other rules even if one fails
				continue
			}
		}
	}

	return nil
}

// matchesRule checks if a transaction matches a rule's conditions
func (e *Engine) matchesRule(transaction *ent.Transaction, rule *ent.Rule) bool {
	// Check category match
	if transaction.Category != string(rule.Category) {
		return false
	}

	// Check amount range if specified
	txAmount := absInt64(transaction.Amount) // Use absolute value for debit/credit

	if rule.MinAmountCents != nil && txAmount < *rule.MinAmountCents {
		return false
	}

	if rule.MaxAmountCents != nil && txAmount > *rule.MaxAmountCents {
		return false
	}

	// All conditions passed
	return true
}

// executeRule creates a rule execution and savings transfer
func (e *Engine) executeRule(ctx context.Context, transaction *ent.Transaction, rule *ent.Rule, account *ent.Account) (ret error) {
	// Calculate savings amount
	savingsAmount, err := e.calculateSavingsAmount(transaction, rule)
	if err != nil {
		return fmt.Errorf("failed to calculate savings amount: %w", err)
	}

	// Start a database transaction for atomicity
	tx, err := e.db.Tx(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer hooks.Rollback(tx, &ret)

	// Create rule execution record (rely on unique constraint for idempotency)
	execution, err := tx.RuleExecution.
		Create().
		SetRuleID(rule.ID).
		SetTransactionID(transaction.ID).
		SetUserID(account.UserID).
		SetAmountCents(savingsAmount).
		SetSourceAccountID(account.ID).
		SetTargetAccountID(rule.TargetAccountID).
		SetStatus(entruleexecution.StatusPending).
		Save(ctx)

	if err != nil {
		// Check if it's a duplicate (idempotency)
		if ent.IsConstraintError(err) {
			log.Debug().
				Str("rule_id", rule.ID.String()).
				Str("transaction_id", transaction.ID.String()).
				Msg("Rule already executed for this transaction, skipping")
			return nil // Not an error, just skip
		}
		return fmt.Errorf("failed to create rule execution: %w", err)
	}

	// Create savings transfer (virtual ledger)
	_, err = tx.SavingsTransfer.
		Create().
		SetRuleExecutionID(execution.ID).
		SetUserID(account.UserID).
		SetSourceAccountID(account.ID).
		SetTargetAccountID(rule.TargetAccountID).
		SetAmountCents(savingsAmount).
		SetStatus(savingstransfer.StatusSuggested).
		Save(ctx)

	if err != nil {
		return fmt.Errorf("failed to create savings transfer: %w", err)
	}

	// Update rule statistics
	err = tx.Rule.
		UpdateOneID(rule.ID).
		AddExecutionCount(1).
		AddTotalSavedCents(savingsAmount).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to update rule stats: %w", err)
	}

	// Mark execution as completed
	err = tx.RuleExecution.
		UpdateOneID(execution.ID).
		SetStatus(entruleexecution.StatusCompleted).
		SetCompletedAt(time.Now()).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to mark execution completed: %w", err)
	}

	// hooks.Rollback will auto-commit here since ret == nil

	log.Info().
		Str("rule_name", rule.Name).
		Str("transaction_name", transaction.Name).
		Int64("amount_cents", savingsAmount).
		Msg("Rule executed successfully")

	// Send notification to user (don't fail if notification fails)
	e.sendNotification(ctx, account.UserID, rule.Name, savingsAmount)

	return nil
}

// sendNotification sends a push notification to the user about the rule execution
func (e *Engine) sendNotification(ctx context.Context, userID uuid.UUID, ruleName string, amountCents int64) {
	// Format amount as dollars
	amountDollars := float64(amountCents) / 100.0

	payload := &notification.Payload{
		Title: fmt.Sprintf("$%.2f moved to Savings", amountDollars),
		Body:  fmt.Sprintf("Rule: %s", ruleName),
	}

	successCount, errs := e.notificationService.SendToUser(ctx, userID, payload)

	if len(errs) > 0 {
		log.Warn().
			Str("user_id", userID.String()).
			Int("success_count", successCount).
			Int("error_count", len(errs)).
			Msg("Some notifications failed to send")
	} else if successCount > 0 {
		log.Debug().
			Str("user_id", userID.String()).
			Int("count", successCount).
			Msg("Notifications sent successfully")
	}
}

// calculateSavingsAmount calculates how much to save based on rule action
func (e *Engine) calculateSavingsAmount(transaction *ent.Transaction, rule *ent.Rule) (int64, error) {
	txAmount := absInt64(transaction.Amount)

	switch rule.ActionType {
	case entrule.ActionTypeMultiply:
		// Multiply transaction amount by action_value
		return int64(float64(txAmount) * float64(rule.ActionValue)), nil

	case entrule.ActionTypeFixed:
		// Fixed amount in dollars -> convert to cents
		return int64(rule.ActionValue * 100), nil

	default:
		return 0, fmt.Errorf("unknown action type: %s", rule.ActionType)
	}
}

// absInt64 returns absolute value of int64
func absInt64(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}
