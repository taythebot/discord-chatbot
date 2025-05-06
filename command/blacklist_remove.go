package command

import (
	"context"
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"

	"github.com/taythebot/discord_chatbot"
	"github.com/taythebot/discord_chatbot/ent/blacklist"
)

// BlacklistRemove removes a user from the blacklist
var BlacklistRemove = Command{
	Name:        "blacklist-remove",
	Description: "Remove a user from the blacklist",
	AdminOnly:   true,
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "user-id",
			Description: "ID for user",
			Required:    true,
		},
	},
	Handler: func(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, r *Registry, logger zerolog.Logger) error {
		options := i.ApplicationCommandData().Options

		// Get user-id option
		var userID string
		for _, option := range options {
			if option.Name == "user-id" {
				userID = option.StringValue()
			}
		}
		if userID == "" {
			return errors.New("failed to get option 'user-id'")
		}

		// Create blacklist in database
		_, err := r.DB.Blacklist.
			Delete().
			Where(blacklist.UserIDEQ(userID)).
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed to create blacklist in database: %w", err)
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{{
					Type:        discordgo.EmbedTypeRich,
					Title:       "User Whitelisted",
					Description: fmt.Sprintf("User %s removed", userID),
					Color:       discord_chatbot.GREEN,
				}},
			},
		})

		return nil
	},
}
