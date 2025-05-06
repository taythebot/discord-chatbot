package command

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"

	"github.com/taythebot/discord_chatbot"
)

// Balance checks provider balance
var Balance = Command{
	Name:        "balance",
	Description: "Check provider balance",
	// AdminOnly:   true,
	Handler: func(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, r *Registry, logger zerolog.Logger) error {
		balance, err := r.Provider.PostCheckBalance(ctx, logger)
		if err != nil {
			return fmt.Errorf("failed to check balance: %w", err)
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{{
					Type:        discordgo.EmbedTypeRich,
					Title:       "Current Balance",
					Description: fmt.Sprintf("%s NANO", balance.NanoBalance),
					Color:       discord_chatbot.GREEN,
				}},
			},
		})

		return nil
	},
}
