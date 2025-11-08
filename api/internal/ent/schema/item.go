package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// Item holds the schema definition for the Item entity.
// Represents a connection to a financial institution via Plaid.
type Item struct {
	ent.Schema
}

// Fields of the Item.
func (Item) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.UUID("user_id", uuid.UUID{}).
			Comment("FK to User who owns this connection"),
		field.String("plaid_id").
			NotEmpty().
			Unique().
			Comment("Plaid's unique identifier for this item"),
		field.String("access_token").
			NotEmpty().
			Sensitive().
			Comment("Encrypted Plaid access token for API calls"),
		field.String("institution_name").
			NotEmpty().
			Comment("Name of the financial institution"),
		field.Bool("is_active").
			Default(true).
			Comment("Whether this connection is still active"),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Item.
func (Item) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("items").
			Field("user_id").
			Required().
			Unique(),
		edge.To("accounts", Account.Type),
		edge.To("sync_cursor", SyncCursor.Type).
			Unique(),
	}
}

// Indexes of the Item.
func (Item) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "is_active"),
		index.Fields("plaid_id").Unique(),
	}
}
