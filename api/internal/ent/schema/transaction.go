package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// Transaction holds the schema definition for the Transaction entity.
// Represents a financial transaction imported from Plaid.
type Transaction struct {
	ent.Schema
}

// Fields of the Transaction.
func (Transaction) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.UUID("account_id", uuid.UUID{}).
			Comment("FK to Account"),
		field.String("plaid_id").
			NotEmpty().
			Unique().
			Comment("Plaid's unique identifier for this transaction"),
		field.Int64("amount").
			Comment("Transaction amount in cents (positive for debit, negative for credit)"),
		field.Time("date").
			Comment("Transaction date"),
		field.String("name").
			NotEmpty().
			Comment("Transaction description/name"),
		field.String("merchant_name").
			Optional().
			Comment("Merchant name if available"),
		field.String("category").
			Default("Misc").
			Comment("App category: Dining, Shopping, Transport, etc."),
		field.JSON("plaid_categories", []string{}).
			Optional().
			Comment("Original Plaid category array"),
		field.Bool("pending").
			Default(false).
			Comment("Whether transaction is still pending"),
		field.String("payment_channel").
			Optional().
			Comment("Payment channel: online, in store, etc."),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Transaction.
func (Transaction) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("account", Account.Type).
			Ref("transactions").
			Field("account_id").
			Required().
			Unique(),
	}
}

// Indexes of the Transaction.
func (Transaction) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("account_id", "date"),
		index.Fields("plaid_id").Unique(),
		index.Fields("category"),
		index.Fields("pending"),
	}
}
