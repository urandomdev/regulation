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
	log.Info().
		Str("transaction_id", transaction.ID.String()).
		Str("transaction_name", transaction.Name).
		Int64("amount_cents", transaction.Amount).
		Str("category", transaction.Category).
		Bool("pending", transaction.Pending).
		Msg("[RULE ENGINE] Processing transaction for rules")

	// Skip pending transactions
	if transaction.Pending {
		log.Debug().
			Str("transaction_id", transaction.ID.String()).
			Msg("[RULE ENGINE] Skipping pending transaction")
		return nil
	}

	// Skip transfers (they're not spending)
	if transaction.Category == "Transfer" {
		log.Debug().
			Str("transaction_id", transaction.ID.String()).
			Msg("[RULE ENGINE] Skipping transfer transaction")
		return nil
	}

	// Get the account to find user_id
	account, err := transaction.QueryAccount().Only(ctx)
	if err != nil {
		return fmt.Errorf("failed to get account: %w", err)
	}

	log.Debug().
		Str("user_id", account.UserID.String()).
		Str("account_id", account.ID.String()).
		Msg("[RULE ENGINE] Loaded account for transaction")

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

	log.Info().
		Str("user_id", account.UserID.String()).
		Int("rule_count", len(rules)).
		Msg("[RULE ENGINE] Loaded active rules for user")

	if len(rules) == 0 {
		log.Warn().
			Str("user_id", account.UserID.String()).
			Msg("[RULE ENGINE] No active rules found for user")
		return nil
	}

	// Log all rules being evaluated
	for i, rule := range rules {
		log.Debug().
			Int("rule_index", i).
			Str("rule_id", rule.ID.String()).
			Str("rule_name", rule.Name).
			Str("rule_category", string(rule.Category)).
			Int("priority", rule.Priority).
			Msg("[RULE ENGINE] Evaluating rule")
	}

	// Evaluate each rule
	matchedCount := 0
	for _, rule := range rules {
		if e.matchesRule(transaction, rule) {
			matchedCount++
			log.Info().
				Str("rule_id", rule.ID.String()).
				Str("rule_name", rule.Name).
				Str("transaction_id", transaction.ID.String()).
				Msg("[RULE ENGINE] Rule matched! Executing...")

			if err := e.executeRule(ctx, transaction, rule, account); err != nil {
				log.Error().
					Err(err).
					Str("rule_id", rule.ID.String()).
					Str("transaction_id", transaction.ID.String()).
					Msg("[RULE ENGINE] Failed to execute rule")
				// Continue with other rules even if one fails
				continue
			}
		}
	}

	log.Info().
		Str("transaction_id", transaction.ID.String()).
		Int("total_rules", len(rules)).
		Int("matched_rules", matchedCount).
		Msg("[RULE ENGINE] Finished processing transaction")

	return nil
}

// matchesRule checks if a transaction matches a rule's conditions
func (e *Engine) matchesRule(transaction *ent.Transaction, rule *ent.Rule) bool {
	txAmount := absInt64(transaction.Amount) // Use absolute value for debit/credit

	log.Debug().
		Str("rule_id", rule.ID.String()).
		Str("rule_name", rule.Name).
		Str("transaction_id", transaction.ID.String()).
		Str("tx_category", transaction.Category).
		Str("rule_category", string(rule.Category)).
		Int64("tx_amount", txAmount).
		Interface("rule_min_amount", rule.MinAmountCents).
		Interface("rule_max_amount", rule.MaxAmountCents).
		Msg("[RULE ENGINE] Checking if transaction matches rule")

	// Check category match
	if transaction.Category != string(rule.Category) {
		log.Debug().
			Str("rule_id", rule.ID.String()).
			Str("rule_name", rule.Name).
			Str("tx_category", transaction.Category).
			Str("rule_category", string(rule.Category)).
			Msg("[RULE ENGINE] Rule does NOT match - category mismatch")
		return false
	}

	// Check amount range if specified
	if rule.MinAmountCents != nil && txAmount < *rule.MinAmountCents {
		log.Debug().
			Str("rule_id", rule.ID.String()).
			Str("rule_name", rule.Name).
			Int64("tx_amount", txAmount).
			Int64("rule_min", *rule.MinAmountCents).
			Msg("[RULE ENGINE] Rule does NOT match - amount below minimum")
		return false
	}

	if rule.MaxAmountCents != nil && txAmount > *rule.MaxAmountCents {
		log.Debug().
			Str("rule_id", rule.ID.String()).
			Str("rule_name", rule.Name).
			Int64("tx_amount", txAmount).
			Int64("rule_max", *rule.MaxAmountCents).
			Msg("[RULE ENGINE] Rule does NOT match - amount above maximum")
		return false
	}

	// All conditions passed
	log.Info().
		Str("rule_id", rule.ID.String()).
		Str("rule_name", rule.Name).
		Str("transaction_id", transaction.ID.String()).
		Str("category", transaction.Category).
		Int64("amount", txAmount).
		Msg("[RULE ENGINE] Rule MATCHES all conditions!")
	return true
}

// executeRule creates a rule execution and savings transfer
func (e *Engine) executeRule(ctx context.Context, transaction *ent.Transaction, rule *ent.Rule, account *ent.Account) (ret error) {
	log.Info().
		Str("rule_id", rule.ID.String()).
		Str("rule_name", rule.Name).
		Str("transaction_id", transaction.ID.String()).
		Str("action_type", string(rule.ActionType)).
		Float64("action_value", rule.ActionValue).
		Msg("[RULE ENGINE] Starting rule execution")

	// Calculate savings amount
	savingsAmount, err := e.calculateSavingsAmount(transaction, rule)
	if err != nil {
		return fmt.Errorf("failed to calculate savings amount: %w", err)
	}

	log.Info().
		Str("rule_id", rule.ID.String()).
		Int64("savings_amount_cents", savingsAmount).
		Float64("savings_amount_dollars", float64(savingsAmount)/100.0).
		Msg("[RULE ENGINE] Calculated savings amount")

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
				Msg("[RULE ENGINE] Rule already executed for this transaction, skipping (idempotent)")
			return nil // Not an error, just skip
		}
		return fmt.Errorf("failed to create rule execution: %w", err)
	}

	log.Debug().
		Str("execution_id", execution.ID.String()).
		Msg("[RULE ENGINE] Created rule execution record")

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

	log.Debug().
		Str("execution_id", execution.ID.String()).
		Msg("[RULE ENGINE] Created savings transfer")

	// Update rule statistics
	err = tx.Rule.
		UpdateOneID(rule.ID).
		AddExecutionCount(1).
		AddTotalSavedCents(savingsAmount).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to update rule stats: %w", err)
	}

	log.Debug().
		Str("rule_id", rule.ID.String()).
		Msg("[RULE ENGINE] Updated rule statistics")

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
		Msg("[RULE ENGINE] Rule executed successfully - committed to database")

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
