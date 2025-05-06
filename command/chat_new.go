package command

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"

	"github.com/taythebot/discord_chatbot"
	"github.com/taythebot/discord_chatbot/ent/channel"
)

// ChatNew creates a new channel for chat
var ChatNew = Command{
	Name:        "chat-new",
	Description: "Creates a new channel for chat",
	GuildOnly:   true,
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:         "model",
			Description:  "AI model to chat with",
			Type:         discordgo.ApplicationCommandOptionString,
			Required:     true,
			Autocomplete: true,
		},
		{
			Name:        "prompt",
			Description: "System prompt",
			Type:        discordgo.ApplicationCommandOptionString,
		},
		{
			Name:         "channel",
			Description:  "Use existing channel",
			Type:         discordgo.ApplicationCommandOptionChannel,
			ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildText},
		},
	},
	Handler: func(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, r *Registry, logger zerolog.Logger) error {
		options := i.ApplicationCommandData().Options

		switch i.Type {
		case discordgo.InteractionApplicationCommandAutocomplete:
			choices := make([]*discordgo.ApplicationCommandOptionChoice, 0, len(r.Models))
			for _, model := range r.Models {
				choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
					Name:  model.Name,
					Value: model.Value,
				})
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionApplicationCommandAutocompleteResult,
				Data: &discordgo.InteractionResponseData{Choices: choices},
			})
		case discordgo.InteractionApplicationCommand:
			var (
				model          string
				prompt         = discord_chatbot.DefaultModelPrompt
				discordChannel *discordgo.Channel
				err            error
			)
			for _, option := range options {
				switch option.Name {
				case "model":
					model = option.StringValue()
				case "prompt":
					prompt = option.StringValue()
				case "channel":
					discordChannel, err = s.Channel(option.ChannelValue(nil).ID)
					if err != nil {
						return fmt.Errorf("failed to resolve channel: %w", err)
					}
				}
			}

			// Check for existing channel
			if discordChannel != nil {
				exists, err := r.DB.Channel.
					Query().
					Where(channel.ID(discordChannel.ID)).
					Exist(ctx)
				if err != nil {
					return fmt.Errorf("failed to check for existing channel: %w", err)
				}

				if exists {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Embeds: []*discordgo.MessageEmbed{{
								Type:        discordgo.EmbedTypeRich,
								Title:       "Existing Channel",
								Description: fmt.Sprintf("The channel <#%s> is already in use.\n\nUse `/chat-info` to see the current channel info\nUse `/chat-model` to change the AI model\nUse `/chat-purge` to purge the chat history", discordChannel.ID),
								Color:       discord_chatbot.RED,
							}},
						},
					})

					return nil
				}
			}

			// Get number of current channels
			n, err := r.DB.Channel.
				Query().
				Where(channel.GuildID(i.GuildID)).
				Count(ctx)
			if err != nil {
				return fmt.Errorf("failed to get total channel count: %w", err)
			}

			// Create new channel
			if discordChannel == nil {
				data := discordgo.GuildChannelCreateData{
					Name:  fmt.Sprintf("chat-bot-%d", n+1),
					Type:  discordgo.ChannelTypeGuildText,
					Topic: fmt.Sprintf("Chat created using model %s by <@%s>", model, i.Member.User.ID),
				}

				discordChannel, err = s.GuildChannelCreateComplex(i.GuildID, data)
				if err != nil {
					return fmt.Errorf("failed to create channel: %w", err)
				}
			}

			// Create channel in database
			err = r.DB.Channel.
				Create().
				SetID(discordChannel.ID).
				SetName(discordChannel.Name).
				SetModel(model).
				SetPrompt(prompt).
				SetGuildID(i.GuildID).
				Exec(ctx)
			if err != nil {
				return fmt.Errorf("failed to create channel in database: %w", err)
			}

			embed := &discordgo.MessageEmbed{
				Type:  discordgo.EmbedTypeRich,
				Title: "New Chat Created",
				Color: discord_chatbot.GREEN,
				Fields: []*discordgo.MessageEmbedField{
					{Name: "Channel", Value: fmt.Sprintf("<#%s>", discordChannel.ID), Inline: true},
					{Name: "Model", Value: model, Inline: true},
					{Name: "System Prompt", Value: prompt},
				},
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{embed},
				},
			})

			// Don't double send message
			if discordChannel.ID != i.ChannelID {
				if _, err := s.ChannelMessageSendEmbed(discordChannel.ID, embed); err != nil {
					logger.Error().Err(err).Msg("Failed to send channel welcome message")
				}
			}
		}
		return nil
	},
}
