package command

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"

	"github.com/taythebot/discord_chatbot"
)

// GuildDisable disables the bot in a guild
var GuildDisable = Command{
	Name:        "guild-disable",
	Description: "Disables the bot in a guild",
	AdminOnly:   true,
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "guild-id",
			Description: "ID for guild",
		},
	},
	Handler: func(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, r *Registry, logger zerolog.Logger) error {
		options := i.ApplicationCommandData().Options

		// Get guild-id option
		guildID := i.GuildID
		for _, option := range options {
			if option.Name == "guild-id" {
				guildID = option.StringValue()
			}
		}

		// Update guild in database
		guild, err := r.DB.Guild.
			UpdateOneID(guildID).
			SetActive(false).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to update guild in database: %w", err)
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{{
					Type:        discordgo.EmbedTypeRich,
					Title:       "Guild Disabled",
					Description: fmt.Sprintf("Guild '%s' (%s) disabled", guild.Name, guild.ID),
					Color:       discord_chatbot.GREEN,
				}},
			},
		})

		return nil
	},
}
