package handler

import (
	"github.com/bwmarrin/discordgo"

	"github.com/taythebot/discord_chatbot/ent"
	"github.com/taythebot/discord_chatbot/provider"
)

// Registry holding all handlers
type Registry struct {
	DB       *ent.Client
	Provider *provider.Client
}

// NewRegistry creates a new registry instance
func NewRegistry(db *ent.Client, provider *provider.Client) *Registry {
	return &Registry{
		DB:       db,
		Provider: provider,
	}
}

// RegisterAll registers all handlers
func (r *Registry) RegisterAll(s *discordgo.Session) {
	s.AddHandler(r.Ready)
	s.AddHandler(r.MessageCreate)
}
