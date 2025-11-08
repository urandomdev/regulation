package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.Nil).Default(MustUUIDv7),

		field.String("email").Unique(),
		field.Bytes("password").Sensitive(),

		field.String("nickname"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
