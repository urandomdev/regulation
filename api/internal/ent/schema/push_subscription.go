package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// PushSubscription holds the schema definition for web push subscriptions.
type PushSubscription struct {
	ent.Schema
}

// Fields of the PushSubscription.
func (PushSubscription) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.Nil).
			Default(MustUUIDv7).
			Immutable().
			Comment("Unique identifier for the push subscription"),

		field.UUID("user_id", uuid.Nil).
			Immutable().
			Comment("User who owns this push subscription"),

		field.String("endpoint").
			NotEmpty().
			Comment("Push service endpoint URL"),

		field.String("p256dh").
			NotEmpty().
			Sensitive().
			Comment("P256DH public key for encryption"),

		field.String("auth").
			NotEmpty().
			Sensitive().
			Comment("Authentication secret for encryption"),

		field.Bool("active").
			Default(true).
			Comment("Whether this subscription is active"),

		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("Timestamp when subscription was created"),

		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("Timestamp when subscription was last updated"),
	}
}

// Edges of the PushSubscription.
func (PushSubscription) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("push_subscriptions").
			Field("user_id").
			Required().
			Unique().
			Immutable().
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

// Indexes of the PushSubscription.
func (PushSubscription) Indexes() []ent.Index {
	return []ent.Index{
		// Index for querying active subscriptions by user
		index.Fields("user_id", "active"),

		// Unique constraint on endpoint to prevent duplicate subscriptions
		index.Fields("endpoint").
			Unique(),
	}
}
