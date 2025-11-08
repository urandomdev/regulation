package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// RuleExecution holds the schema definition for the RuleExecution entity.
// Tracks each time a rule is triggered by a transaction.
type RuleExecution struct {
	ent.Schema
}

// Fields of the RuleExecution.
func (RuleExecution) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.UUID("rule_id", uuid.UUID{}).
			Comment("FK to Rule that was triggered"),
		field.UUID("transaction_id", uuid.UUID{}).
			Comment("FK to Transaction that triggered the rule"),
		field.UUID("user_id", uuid.UUID{}).
			Comment("FK to User (denormalized for quick queries)"),

		// Transfer details
		field.Int64("amount_cents").
			Comment("Calculated savings amount in cents"),
		field.UUID("source_account_id", uuid.UUID{}).
			Comment("FK to Account where transaction occurred"),
		field.UUID("target_account_id", uuid.UUID{}).
			Comment("FK to Account where savings should go"),

		// Status tracking
		field.Enum("status").
			Values("pending", "completed", "failed").
			Default("pending").
			Comment("Execution status"),
		field.String("error_message").
			Optional().
			Comment("Error details if status is failed"),

		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("completed_at").
			Optional().
			Nillable().
			Comment("When the transfer was completed"),
	}
}

// Edges of the RuleExecution.
func (RuleExecution) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("rule", Rule.Type).
			Ref("executions").
			Field("rule_id").
			Required().
			Unique(),
		edge.From("transaction", Transaction.Type).
			Ref("rule_executions").
			Field("transaction_id").
			Required().
			Unique(),
		edge.From("user", User.Type).
			Ref("rule_executions").
			Field("user_id").
			Required().
			Unique(),
		edge.To("transfer", SavingsTransfer.Type).
			Unique(),
	}
}

// Indexes of the RuleExecution.
func (RuleExecution) Indexes() []ent.Index {
	return []ent.Index{
		// Idempotency: prevent duplicate rule executions for same transaction
		index.Fields("rule_id", "transaction_id").Unique(),
		index.Fields("user_id", "created_at"),
		index.Fields("status"),
	}
}
