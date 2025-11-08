package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// Account holds the schema definition for the Account entity.
// Represents an individual bank account linked through Plaid.
type Account struct {
	ent.Schema
}

// Fields of the Account.
func (Account) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.UUID("item_id", uuid.UUID{}).
			Comment("FK to Item"),
		field.UUID("user_id", uuid.UUID{}).
			Comment("FK to User for quick queries"),
		field.String("plaid_id").
			NotEmpty().
			Unique().
			Comment("Plaid's unique identifier for this account"),
		field.String("name").
			NotEmpty().
			Comment("User-friendly account name from bank"),
		field.Enum("type").
			Values("checking", "savings", "credit", "other").
			Default("checking").
			Comment("Type of account"),
		field.String("subtype").
			Optional().
			Comment("Detailed account subtype from Plaid"),
		field.String("mask").
			Optional().
			Comment("Last 4 digits of account number"),
		field.Int64("current_balance").
			Default(0).
			Comment("Current balance in cents"),
		field.Int64("available_balance").
			Optional().
			Nillable().
			Comment("Available balance in cents (may be nil)"),
		field.Bool("is_active").
			Default(true).
			Comment("Whether this account is still active"),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Account.
func (Account) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("item", Item.Type).
			Ref("accounts").
			Field("item_id").
			Required().
			Unique(),
		edge.From("user", User.Type).
			Ref("accounts").
			Field("user_id").
			Required().
			Unique(),
		edge.To("transactions", Transaction.Type),
		edge.To("target_rules", Rule.Type),
		edge.To("outgoing_transfers", SavingsTransfer.Type),
		edge.To("incoming_transfers", SavingsTransfer.Type),
	}
}

// Indexes of the Account.
func (Account) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "is_active"),
		index.Fields("plaid_id").Unique(),
		index.Fields("item_id"),
	}
}
