package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// Rule holds the schema definition for the Rule entity.
// Represents a user-defined savings rule that triggers on transactions.
type Rule struct {
	ent.Schema
}

// Fields of the Rule.
func (Rule) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.UUID("user_id", uuid.UUID{}).
			Comment("FK to User who owns this rule"),
		field.String("name").
			NotEmpty().
			Comment("User-friendly rule name (e.g., 'Dining 2x Save')"),

		// Condition: Category-based matching
		field.Enum("category").
			Values(
				"Dining",
				"Groceries",
				"Transport",
				"Shopping",
				"Subscriptions",
				"Entertainment",
				"Bills",
				"Misc",
			).
			Comment("Transaction category to match"),

		// Condition: Amount range (optional)
		field.Int64("min_amount_cents").
			Optional().
			Nillable().
			Comment("Minimum transaction amount to trigger (cents)"),
		field.Int64("max_amount_cents").
			Optional().
			Nillable().
			Comment("Maximum transaction amount to trigger (cents)"),

		// Action: What to do when rule matches
		field.Enum("action_type").
			Values("multiply", "fixed").
			Comment("Type of savings action: multiply transaction amount or fixed amount"),
		field.Float("action_value").
			Comment("Multiplier (e.g., 1.5 for 1.5x) or fixed amount in dollars"),

		// Target account for savings
		field.UUID("target_account_id", uuid.UUID{}).
			Comment("FK to Account where savings should be transferred"),

		// Control
		field.Bool("is_active").
			Default(true).
			Comment("Whether this rule is currently enabled"),
		field.Int("priority").
			Default(0).
			Comment("Rule priority (lower number = higher priority)"),

		// Stats
		field.Int("execution_count").
			Default(0).
			Comment("Number of times this rule has been triggered"),
		field.Int64("total_saved_cents").
			Default(0).
			Comment("Total amount saved by this rule (cents)"),

		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Rule.
func (Rule) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("rules").
			Field("user_id").
			Required().
			Unique(),
		edge.From("target_account", Account.Type).
			Ref("target_rules").
			Field("target_account_id").
			Required().
			Unique(),
		edge.To("executions", RuleExecution.Type),
	}
}

// Indexes of the Rule.
func (Rule) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "is_active"),
		index.Fields("category"),
		index.Fields("priority"),
	}
}
