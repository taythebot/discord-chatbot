package handler

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

// Ready handler on Discord connect
func (r *Registry) Ready(s *discordgo.Session, _ *discordgo.Ready) {
	log.Info().Str("handler", "ready").Msgf("Connected to Discord as %v#%v", s.State.User.Username, s.State.User.Discriminator)
}
