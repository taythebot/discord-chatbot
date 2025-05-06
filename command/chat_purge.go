package command

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"

	"github.com/taythebot/discord_chatbot"
	"github.com/taythebot/discord_chatbot/ent/channel"
	"github.com/taythebot/discord_chatbot/ent/message"
)

// ChatPurge removes all previous chat history from database
var ChatPurge = Command{
	Name:        "chat-purge",
	Description: "Removes all previous channel messages from database",
	GuildOnly:   true,
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:         "channel",
			Description:  "Use another channel",
			Type:         discordgo.ApplicationCommandOptionChannel,
			ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildText},
		},
	},
	Handler: func(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, r *Registry, logger zerolog.Logger) error {
		options := i.ApplicationCommandData().Options

		channelID := i.ChannelID
		for _, option := range options {
			if option.Name == "channel" {
				channelID = option.ChannelValue(nil).ID
			}
		}

		// Check if channel exists
		exists, err := r.DB.Channel.
			Query().
			Where(channel.ID(channelID)).
			Exist(ctx)
		if err != nil {
			return fmt.Errorf("failed to get channel from database: %w", err)
		} else if !exists {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{discord_chatbot.InvalidChannelEmbed},
				},
			})

			return nil
		}

		// Delete messages
		affected, err := r.DB.Message.
			Delete().
			Where(message.ChannelID(channelID)).
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed to delete messages from database: %w", err)
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{{
					Type:        discordgo.EmbedTypeRich,
					Title:       "Chat Purged",
					Description: fmt.Sprintf("Deleted %d messages from channel <#%s> in the database", affected, channelID),
					Color:       discord_chatbot.GREEN,
				}},
			},
		})

		return nil
	},
}
