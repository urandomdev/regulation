package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// SyncCursor holds the schema definition for the SyncCursor entity.
// Tracks the sync state for each Item to enable incremental transaction fetching.
type SyncCursor struct {
	ent.Schema
}

// Fields of the SyncCursor.
func (SyncCursor) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.UUID("item_id", uuid.UUID{}).
			Unique().
			Comment("FK to Item"),
		field.String("cursor").
			Default("").
			Comment("Plaid sync cursor for incremental updates"),
		field.Time("last_sync_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("Last sync attempt timestamp"),
		field.String("last_error").
			Default("").
			Comment("Most recent error message, empty if successful"),
		field.Int("consecutive_failures").
			Default(0).
			Comment("Number of consecutive failed sync attempts"),
	}
}

// Edges of the SyncCursor.
func (SyncCursor) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("item", Item.Type).
			Ref("sync_cursor").
			Field("item_id").
			Required().
			Unique(),
	}
}

// Indexes of the SyncCursor.
func (SyncCursor) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("item_id").Unique(),
	}
}
