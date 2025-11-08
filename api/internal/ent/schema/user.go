package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
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

		field.UUID("custody_account_id", uuid.Nil).Optional().Nillable(),

		field.String("nickname"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("custody_account", User.Type).
			Field("custody_account_id").
			Unique().
			From("user"),
		edge.To("items", Item.Type),
		edge.To("accounts", Account.Type),
		edge.To("push_subscriptions", PushSubscription.Type),
		edge.To("rules", Rule.Type),
		edge.To("rule_executions", RuleExecution.Type),
		edge.To("savings_transfers", SavingsTransfer.Type),
	}
}
