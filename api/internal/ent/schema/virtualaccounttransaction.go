package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// VirtualAccountTransaction holds the schema definition for the VirtualAccountTransaction entity.
type VirtualAccountTransaction struct {
	ent.Schema
}

// Fields of the VirtualAccountTransaction.
func (VirtualAccountTransaction) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.Nil).Default(MustUUIDv7),

		field.Int("adjusted_dollars"),
		field.Int("adjusted_cents"),

		field.UUID("virtual_account_id", uuid.Nil),
	}
}

// Edges of the VirtualAccountTransaction.
func (VirtualAccountTransaction) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("virtual_account", VirtualAccount.Type).
			Field("virtual_account_id").
			Ref("transactions").
			Unique().Required(),
	}
}
