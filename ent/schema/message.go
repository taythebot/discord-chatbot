package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Message holds the schema definition for the Message entity.
type Message struct {
	ent.Schema
}

// Fields of the Message.
func (Message) Fields() []ent.Field {
	return []ent.Field{
		field.String("message_id").
			Unique().
			Immutable().
			Optional(),
		field.String("user_id").
			Immutable().
			Optional(),
		field.String("user_name").
			Immutable().
			Optional(),
		field.String("content").
			NotEmpty().
			Immutable(),
		field.String("channel_id").
			NotEmpty().
			Immutable(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
	}
}

// Edges of the Message.
func (Message) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", Channel.Type).
			Ref("messages").
			Field("channel_id").
			Unique().
			Required().
			Immutable(),
	}
}
