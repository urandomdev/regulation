package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// VirtualAccount holds the schema definition for the VirtualAccount entity.
type VirtualAccount struct {
	ent.Schema
}

// Fields of the VirtualAccount.
func (VirtualAccount) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.Nil).Default(MustUUIDv7),
		field.Enum("type").Values("checking", "saving"),

		field.String("name").NotEmpty(),

		field.Int("dollars"),
		field.Int("cents"),

		field.UUID("user_id", uuid.Nil),
	}
}

// Edges of the VirtualAccount.
func (VirtualAccount) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("transactions", VirtualAccountTransaction.Type),

		edge.From("user", User.Type).
			Field("user_id").
			Ref("accounts").
			Required().Unique(),
	}
}
