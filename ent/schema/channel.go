package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Channel holds the schema definition for the Channel entity.
type Channel struct {
	ent.Schema
}

// Fields of the Channel.
func (Channel) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			NotEmpty().
			Unique().
			Immutable(),
		field.String("name").
			NotEmpty(),
		field.String("model").
			NotEmpty(),
		field.String("prompt").
			NotEmpty(),
		field.String("guild_id").
			NotEmpty().
			Immutable(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Channel.
func (Channel) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", Guild.Type).
			Ref("channels").
			Field("guild_id").
			Unique().
			Required().
			Immutable(),
		edge.To("messages", Message.Type),
	}
}
