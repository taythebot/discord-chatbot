package command

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"

	"github.com/taythebot/discord_chatbot"
	"github.com/taythebot/discord_chatbot/ent/blacklist"
)

// BlacklistList lists all users from blacklist
var BlacklistList = Command{
	Name:        "blacklist-list",
	Description: "Lists all users from blacklist",
	AdminOnly:   true,
	Handler: func(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, r *Registry, logger zerolog.Logger) error {
		// Fetch all blacklist
		results, err := r.DB.Blacklist.
			Query().
			Select(blacklist.FieldUserID, blacklist.FieldReason).
			All(ctx)
		if err != nil {
			return fmt.Errorf("failed to fetch blacklist from database: %w", err)
		}

		var userIDs, reasons string
		for _, item := range results {
			userIDs += fmt.Sprintf("%s\n", item.UserID)
			reasons += fmt.Sprintf("%s\n", item.Reason)
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{{
					Type:  discordgo.EmbedTypeRich,
					Title: fmt.Sprintf("Total of %d user in blacklist", len(results)),
					Color: discord_chatbot.GREEN,
					Fields: []*discordgo.MessageEmbedField{
						{Name: "User ID", Value: userIDs, Inline: true},
						{Name: "Reason", Value: reasons, Inline: true},
					},
				}},
			},
		})

		return nil
	},
}
