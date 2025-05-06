package command

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"

	"github.com/taythebot/discord_chatbot"
	"github.com/taythebot/discord_chatbot/ent"
)

// ChatModel changes the model for a chat
var ChatModel = Command{
	Name:        "chat-model",
	Description: "Changes the model for a chat",
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
			var model, prompt string
			for _, option := range options {
				switch option.Name {
				case "model":
					model = option.StringValue()
				case "prompt":
					prompt = option.StringValue()
				}
			}

			// Check if current channel is active
			channelItem, err := r.DB.Channel.Get(ctx, i.ChannelID)
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

				return fmt.Errorf("failed to get total channel count: %w", err)
			}

			if prompt == "" {
				prompt = channelItem.Prompt
			}

			// Update channel in database
			err = r.DB.Channel.
				UpdateOneID(i.ChannelID).
				SetModel(model).
				SetPrompt(prompt).
				Exec(ctx)
			if err != nil {
				return fmt.Errorf("failed to update channel in database: %w", err)
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{&discordgo.MessageEmbed{
						Type:  discordgo.EmbedTypeRich,
						Title: "Model",
						Color: discord_chatbot.GREEN,
						Fields: []*discordgo.MessageEmbedField{
							{Name: "Model", Value: model, Inline: true},
							{Name: "Channel", Value: fmt.Sprintf("<#%s>", i.ChannelID), Inline: true},
							{Name: "System Prompt", Value: prompt},
						},
					}},
				},
			})
		}

		return nil
	},
}
