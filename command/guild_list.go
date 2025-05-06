package command

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"

	"github.com/taythebot/discord_chatbot"
	"github.com/taythebot/discord_chatbot/ent"
	"github.com/taythebot/discord_chatbot/ent/guild"
)

// GuildList lists all available guilds from the bot config
var GuildList = Command{
	Name:        "guild-list",
	Description: "Lists all guilds from the bot config",
	AdminOnly:   true,
	Handler: func(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, r *Registry, logger zerolog.Logger) error {
		// Fetch all guilds
		var guilds []ent.Guild
		err := r.DB.Guild.
			Query().
			Select(guild.FieldID, guild.FieldName, guild.FieldActive, guild.FieldUpdatedAt, guild.FieldCreatedAt).
			Scan(ctx, &guilds)
		if err != nil {
			return fmt.Errorf("failed to fetch guilds from database: %w", err)
		}

		var ids, names, actives string
		for _, guild := range guilds {
			ids += fmt.Sprintf("%s\n", guild.ID)
			names += fmt.Sprintf("%s\n", guild.Name)
			actives += fmt.Sprintf("%t\n", guild.Active)
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{{
					Type:  discordgo.EmbedTypeRich,
					Title: fmt.Sprintf("Total of %d guilds in bot", len(guilds)),
					Color: discord_chatbot.GREEN,
					Fields: []*discordgo.MessageEmbedField{
						{Name: "ID", Value: ids, Inline: true},
						{Name: "Name", Value: names, Inline: true},
						{Name: "Active", Value: actives, Inline: true},
					},
				}},
			},
		})

		return nil
	},
}
