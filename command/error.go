package command

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	json "github.com/goccy/go-json"
	"github.com/rs/zerolog"

	"github.com/taythebot/discord_chatbot"
)

// errorHandler for failed interactions
func errorHandler(err error, s *discordgo.Session, i *discordgo.InteractionCreate, logger zerolog.Logger) {
	logger.Error().Err(err).Msg("Interaction failed")

	fields := []*discordgo.MessageEmbedField{
		{Name: "Error", Value: fmt.Sprintf("```json\n%s\n```", err)},
		{Name: "Interaction ID", Value: i.ID},
	}

	options := i.ApplicationCommandData().Options
	if len(options) > 0 {
		b, err := json.Marshal(options)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to marshal options data")
		} else {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:  "Options",
				Value: fmt.Sprintf("```json\n%s\n```", b),
			})
		}
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{{
				Type:   discordgo.EmbedTypeRich,
				Title:  "Interaction Failed",
				Color:  discord_chatbot.RED,
				Fields: fields,
			}},
		},
	})
}
