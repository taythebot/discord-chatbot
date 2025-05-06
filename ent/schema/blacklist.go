package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Blacklist holds the schema definition for the Blacklist entity.
type Blacklist struct {
	ent.Schema
}

// Fields of the Blacklist.
func (Blacklist) Fields() []ent.Field {
	return []ent.Field{
		field.String("user_id").
			NotEmpty().
			Unique().
			Immutable(),
		field.String("reason").
			NotEmpty(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Blacklist.
func (Blacklist) Edges() []ent.Edge {
	return nil
}
