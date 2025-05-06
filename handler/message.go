package handler

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/taythebot/discord_chatbot"
	"github.com/taythebot/discord_chatbot/ent"
	"github.com/taythebot/discord_chatbot/ent/blacklist"
	"github.com/taythebot/discord_chatbot/ent/message"
	"github.com/taythebot/discord_chatbot/provider"
)

// messageErrorHandler for failed message processing
func messageErrorHandler(err error, s *discordgo.Session, m *discordgo.MessageCreate, logger zerolog.Logger) {
	logger.Error().Err(err).Msg("Failed to process message")

	s.ChannelMessageSendEmbedReply(m.ChannelID, &discordgo.MessageEmbed{
		Type:  discordgo.EmbedTypeRich,
		Title: "MessageCreate Error",
		Color: discord_chatbot.RED,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Error", Value: fmt.Sprintf("```json\n%s\n```", err)},
			{Name: "Message ID", Value: m.ID},
		},
	}, m.Reference())
}

// Function to split a message into chunks of a specified size, considering new lines
func splitMessage(message string, maxLength int) []string {
	var chunks []string
	for len(message) > maxLength {
		// Find the last newline or space within the maxLength limit
		lastNewline := strings.LastIndex(message[:maxLength], "\n")
		lastSpace := strings.LastIndex(message[:maxLength], " ")

		// Choose the best split point: prefer newline over space
		splitPoint := lastNewline
		if lastNewline == -1 || (lastSpace > lastNewline) {
			splitPoint = lastSpace
		}

		if splitPoint == -1 {
			splitPoint = maxLength // No newline or space found, split at maxLength
		}

		chunks = append(chunks, message[:splitPoint])
		message = message[splitPoint:]
	}
	chunks = append(chunks, message) // Add the remaining part
	return chunks
}

// MessageCreate handler on new messages
func (r *Registry) MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	ctx := context.TODO()
	defer ctx.Done()

	// Ignore self, bots, and DMs
	if m.Author.ID == s.State.User.ID || m.Author.Bot || m.GuildID == "" {
		return
	}

	// Check if replied to
	shouldProcess := m.ReferencedMessage != nil && m.ReferencedMessage.Author.ID == s.State.User.ID

	// Check if pinged
	if !shouldProcess {
		for _, user := range m.Mentions {
			if user.ID == s.State.User.ID {
				shouldProcess = true
				break
			}
		}
	}

	if !shouldProcess {
		return
	}

	logger := log.With().
		Str("handler", "messageCreate").
		Str("msg", m.ID).
		Str("channel", m.ChannelID).
		Str("guild", m.GuildID).
		Str("user", m.Author.ID).
		Logger()

	logger.Debug().Msgf("Processing message '%s'", m.ID)

	// Check if from blacklisted user
	blacklistExists, err := r.DB.Blacklist.
		Query().
		Where(blacklist.UserID(m.Author.ID)).
		Exist(ctx)
	if err != nil {
		messageErrorHandler(fmt.Errorf("failed to get user from blacklist: %w", err), s, m, logger)
		return
	} else if blacklistExists {
		s.ChannelMessageSendEmbedReply(m.ChannelID, discord_chatbot.BlacklistedUserEmbed, m.Reference())
		return
	}

	// Check if valid channelItem
	channelItem, err := r.DB.Channel.Get(ctx, m.ChannelID)
	if err != nil {
		if ent.IsNotFound(err) {
			s.ChannelMessageSendEmbedReply(m.ChannelID, discord_chatbot.InvalidChannelEmbed, m.Reference())
			return
		}

		messageErrorHandler(fmt.Errorf("failed to get channel from database: %w", err), s, m, logger)
		return
	}

	if channelItem == nil {
		logger.Debug().Msg("Ignoring message from invalid channel")
		return
	}

	// Get all chats
	messages, err := r.DB.Message.
		Query().
		Select(message.FieldUserName, message.FieldContent).
		Where(message.ChannelID(m.ChannelID)).
		Order(ent.Asc(message.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		messageErrorHandler(fmt.Errorf("failed to get previous messages from database: %w", err), s, m, logger)
		return
	}

	// Request chat completion
	request := provider.PostChatCompletionsRequest{
		Model:            channelItem.Model,
		Stream:           false,
		Messages:         make([]provider.PostChatCompletionsRequestMessage, 0, len(messages)),
		Temperature:      0.7,
		MaxTokens:        2000,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
	}

	// Add system prompt
	request.Messages = append(request.Messages, provider.PostChatCompletionsRequestMessage{
		Role:    provider.RoleSystem,
		Content: channelItem.Prompt,
	})

	// Add previous messages
	for _, message := range messages {
		role := provider.RoleAssistant
		if message.MessageID == "" {
			role = provider.RoleUser
		}

		request.Messages = append(request.Messages, provider.PostChatCompletionsRequestMessage{
			Role:    role,
			Content: fmt.Sprintf("%s: %s", message.UserName, message.Content),
		})
	}

	// Get new reply
	content := m.ContentWithMentionsReplaced()
	request.Messages = append(request.Messages, provider.PostChatCompletionsRequestMessage{
		Role:    provider.RoleUser,
		Content: fmt.Sprintf("%s: %s", m.Author.GlobalName, content),
	})

	response, err := r.Provider.PostChatCompletions(ctx, logger, request)
	if err != nil {
		messageErrorHandler(fmt.Errorf("failed to get chat completion: %w", err), s, m, logger)
		return
	}

	if len(response.Choices) < 1 {
		messageErrorHandler(errors.New("failed to get chat completion: no choices provided"), s, m, logger)
		return
	}

	// Reply to user
	reply := response.Choices[0].Message.Content

	if len(reply) > 2000 {
		// Split the reply into chunks
		chunks := splitMessage(reply, 2000)
		for _, chunk := range chunks {
			if _, err := s.ChannelMessageSendReply(m.ChannelID, chunk, m.Reference()); err != nil {
				messageErrorHandler(fmt.Errorf("failed to send message: %w", err), s, m, logger)
				return
			}
		}
	} else {
		if _, err := s.ChannelMessageSendReply(m.ChannelID, reply, m.Reference()); err != nil {
			messageErrorHandler(fmt.Errorf("failed to send message: %w", err), s, m, logger)
			return
		}
	}

	// Save messages
	err = r.DB.Message.CreateBulk(
		r.DB.Message.
			Create().
			SetMessageID(m.ID).
			SetUserID(m.Author.ID).
			SetUserName(m.Author.GlobalName).
			SetContent(content).
			SetChannelID(m.ChannelID),
		r.DB.Message.
			Create().
			SetContent(reply).
			SetChannelID(m.ChannelID),
	).Exec(ctx)
	if err != nil {
		messageErrorHandler(fmt.Errorf("failed to save messages in database: %w", err), s, m, logger)
		return
	}
}
