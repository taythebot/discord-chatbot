package command

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"

	"github.com/taythebot/discord_chatbot"
	"github.com/taythebot/discord_chatbot/ent"
	"github.com/taythebot/discord_chatbot/ent/message"
)

// ChatInfo returns the current chat info for a channel
var ChatInfo = Command{
	Name:        "chat-info",
	Description: "Returns the current chat info for a channel",
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

		// Get channel
		result, err := r.DB.Channel.Get(ctx, channelID)
		if err != nil {
			if ent.IsNotFound(err) {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{discord_chatbot.InvalidChannelEmbed},
					},
				})

				return nil
			}

			return fmt.Errorf("failed to get channel from database: %w", err)
		}

		// Get total messages
		total, err := r.DB.Message.
			Query().
			Where(message.ChannelID(channelID)).
			Count(ctx)
		if err != nil {
			return fmt.Errorf("failed to get total messages count from database: %w", err)
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{{
					Type:  discordgo.EmbedTypeRich,
					Title: "Chat Information",
					Color: discord_chatbot.GREEN,
					Fields: []*discordgo.MessageEmbedField{
						{Name: "Channel", Value: fmt.Sprintf("<#%s>", result.ID), Inline: true},
						{Name: "Model", Value: result.Model, Inline: true},
						{Name: "Total Messages", Value: fmt.Sprintf("%d", total), Inline: true},
						{Name: "System Prompt", Value: result.Prompt},
					},
				}},
			},
		})

		return nil
	},
}
