package command

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"

	"github.com/taythebot/discord_chatbot"
)

// GuildEnable enables the bot in a guild
var GuildEnable = Command{
	Name:        "guild-enable",
	Description: "Enables the bot in a guild",
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

		// Find guild
		guild, err := s.Guild(guildID)
		if err != nil {
			return fmt.Errorf("failed to get guild: %w", err)
		}

		// Create guild in database
		err = r.DB.Guild.
			Create().
			SetID(guild.ID).
			SetName(guild.Name).
			SetActive(true).
			OnConflict().
			UpdateActive().
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed to create guild in database: %w", err)
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{{
					Type:        discordgo.EmbedTypeRich,
					Title:       "Guild Enabled",
					Description: fmt.Sprintf("Guild '%s' (%s) enabled", guild.Name, guild.ID),
					Color:       discord_chatbot.GREEN,
				}},
			},
		})

		return nil
	},
}
