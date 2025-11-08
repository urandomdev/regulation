package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// SavingsTransfer holds the schema definition for the SavingsTransfer entity.
// Represents a suggested or executed savings transfer (virtual ledger).
type SavingsTransfer struct {
	ent.Schema
}

// Fields of the SavingsTransfer.
func (SavingsTransfer) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.UUID("rule_execution_id", uuid.UUID{}).
			Comment("FK to RuleExecution that created this transfer"),
		field.UUID("user_id", uuid.UUID{}).
			Comment("FK to User"),

		// Transfer details
		field.UUID("source_account_id", uuid.UUID{}).
			Comment("FK to Account (checking)"),
		field.UUID("target_account_id", uuid.UUID{}).
			Comment("FK to Account (savings)"),
		field.Int64("amount_cents").
			Comment("Amount to transfer in cents"),

		// Status
		field.Enum("status").
			Values("suggested", "approved", "executed", "cancelled", "failed").
			Default("suggested").
			Comment("Transfer status"),
		field.String("plaid_transfer_id").
			Optional().
			Comment("Plaid transfer ID if real transfer was executed"),
		field.String("error_message").
			Optional().
			Comment("Error details if failed"),

		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("executed_at").
			Optional().
			Nillable().
			Comment("When transfer was executed"),
	}
}

// Edges of the SavingsTransfer.
func (SavingsTransfer) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("rule_execution", RuleExecution.Type).
			Ref("transfer").
			Field("rule_execution_id").
			Required().
			Unique(),
		edge.From("user", User.Type).
			Ref("savings_transfers").
			Field("user_id").
			Required().
			Unique(),
		edge.From("source_account", Account.Type).
			Ref("outgoing_transfers").
			Field("source_account_id").
			Required().
			Unique(),
		edge.From("target_account", Account.Type).
			Ref("incoming_transfers").
			Field("target_account_id").
			Required().
			Unique(),
	}
}

// Indexes of the SavingsTransfer.
func (SavingsTransfer) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "status"),
		index.Fields("created_at"),
		index.Fields("rule_execution_id").Unique(),
	}
}
