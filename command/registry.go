package command

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/taythebot/discord_chatbot"
	"github.com/taythebot/discord_chatbot/ent"
	"github.com/taythebot/discord_chatbot/ent/blacklist"
	"github.com/taythebot/discord_chatbot/ent/guild"
	"github.com/taythebot/discord_chatbot/provider"
)

// Registry holding all commands
type Registry struct {
	DB       *ent.Client
	Provider *provider.Client
	AdminIDs []string
	Models   []discord_chatbot.Model
	Commands []Command
}

// Command global struct for making commands
type Command struct {
	Name        string
	Description string
	AdminOnly   bool
	GuildOnly   bool
	Options     []*discordgo.ApplicationCommandOption
	Handler     func(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, r *Registry, logger zerolog.Logger) error
}

// NewRegistry for commands
func NewRegistry(db *ent.Client, provider *provider.Client, adminIDs []string, models []discord_chatbot.Model) *Registry {
	return &Registry{
		DB:       db,
		Provider: provider,
		AdminIDs: adminIDs,
		Models:   models,
		Commands: []Command{
			GuildEnable,
			GuildDisable,
			GuildList,
			BlacklistAdd,
			BlacklistRemove,
			BlacklistList,
			Balance,
			ChatNew,
			ChatInfo,
			ChatModel,
			ChatPurge,
		},
	}
}

// Register all slash commands and handlers from registry with Discord
func (r *Registry) RegisterAll(s *discordgo.Session, appID string) error {
	// Open websocket connection
	log.Debug().Msg("Opening websocket connection to Discord")

	if err := s.Open(); err != nil {
		return fmt.Errorf("failed to open websocket connection: %w", err)
	}

	// Register commands
	commands := make([]*discordgo.ApplicationCommand, 0, len(r.Commands))
	for _, command := range r.Commands {
		cmd := &discordgo.ApplicationCommand{
			Name:        command.Name,
			Description: command.Description,
			Options:     command.Options,
		}

		commands = append(commands, cmd)
	}
	if _, err := s.ApplicationCommandBulkOverwrite(appID, "", commands); err != nil {
		return fmt.Errorf("failed to bulk register commands: %w", err)
	}

	log.Info().Msgf("Registered %d commands", len(commands))

	// Handler command interactions
	s.AddHandler(r.HandleInteraction)

	return nil
}

// HandleInteraction handles command interactions
func (r *Registry) HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var (
		ctx         = context.TODO()
		commandName = i.ApplicationCommandData().Name
	)
	defer ctx.Done()

	// Get user ID from interaction
	var userID string
	if i.Member != nil {
		userID = i.Member.User.ID
	} else {
		userID = i.User.ID
	}

	logger := log.
		With().
		Str("interaction", i.ID).
		Str("guild", i.GuildID).
		Str("channel", i.ChannelID).
		Str("user", userID).
		Logger()

	logger.Info().Msgf("Received interaction '%s' for command '%s'", i.ID, commandName)

	// Find command in registry
	var command *Command
	for _, cmd := range r.Commands {
		if cmd.Name == commandName {
			command = &cmd
			break
		}
	}
	if command == nil {
		errorHandler(fmt.Errorf("failed to find command handler for '%s'", i.ApplicationCommandData().Name), s, i, logger)
		return
	}

	// Check guildOnly
	if command.GuildOnly {
		if i.GuildID == "" {
			errorHandler(errors.New("command is only available in guilds"), s, i, logger)
			return
		}

		// Check if guild is active
		result, err := r.DB.Guild.
			Query().
			Select(guild.FieldActive).
			Where(guild.ID(i.GuildID)).
			First(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{{
							Type:        discordgo.EmbedTypeRich,
							Title:       "Guild Not Enabled",
							Description: "Bot usage has not been enabled yet in this guild. Contact the admin for help.",
							Color:       discord_chatbot.RED,
						}},
					},
				})
			} else {
				errorHandler(fmt.Errorf("failed to check guild: %w", err), s, i, logger)
			}

			return
		}

		if !result.Active {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{discord_chatbot.DisabledGuildEmbed},
				},
			})
			return
		}
	}

	// Check permissions
	if command.AdminOnly && !slices.Contains(r.AdminIDs, userID) {
		errorHandler(errors.New("permission denied: user is not admin"), s, i, logger)
		return
	} else {
		// Check blacklist
		exists, err := r.DB.Blacklist.
			Query().
			Where(blacklist.UserID(userID)).
			Exist(ctx)
		if err != nil {
			errorHandler(fmt.Errorf("failed to check blacklist: %w", err), s, i, logger)
			return
		}

		if exists {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{discord_chatbot.BlacklistedUserEmbed},
				},
			})
			return
		}
	}

	// Execute command handler
	if err := command.Handler(ctx, s, i, r, logger); err != nil {
		errorHandler(err, s, i, logger)
	}
}
