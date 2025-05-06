package command

import (
	"context"
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"

	"github.com/taythebot/discord_chatbot"
)

// BlacklistAdd adds a user to the blacklist
var BlacklistAdd = Command{
	Name:        "blacklist-add",
	Description: "Adds a user to the blacklist",
	AdminOnly:   true,
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "user-id",
			Description: "ID for user",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "reason",
			Description: "Reason for blacklist",
			Required:    true,
		},
	},
	Handler: func(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, r *Registry, logger zerolog.Logger) error {
		options := i.ApplicationCommandData().Options

		// Get user-id option
		var (
			userID string
			reason string
		)
		for _, option := range options {
			switch option.Name {
			case "user-id":
				userID = option.StringValue()
			case "reason":
				reason = option.StringValue()
			}
		}
		if userID == "" {
			return errors.New("failed to get option 'user-id'")
		}

		// Create blacklist in database
		err := r.DB.Blacklist.
			Create().
			SetUserID(userID).
			SetReason(reason).
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed to create blacklist in database: %w", err)
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{{
					Type:        discordgo.EmbedTypeRich,
					Title:       "User Blacklisted",
					Description: fmt.Sprintf("User %s blacklisted", userID),
					Color:       discord_chatbot.GREEN,
				}},
			},
		})

		return nil
	},
}
