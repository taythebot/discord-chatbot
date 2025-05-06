package discord_chatbot

import (
	"github.com/bwmarrin/discordgo"
)

var InvalidChannelEmbed = &discordgo.MessageEmbed{
	Type:        discordgo.EmbedTypeRich,
	Title:       "Invalid Channel",
	Description: "This channel has not been enabled for bot usage.\n\nUse `/chat-new` with the `channel` option to enable this channel.",
	Color:       RED,
}

var BlacklistedUserEmbed = &discordgo.MessageEmbed{
	Type:        discordgo.EmbedTypeRich,
	Title:       "User Blacklisted",
	Description: "You have been blacklisted from this bot. Contact the admin for help.",
	Color:       RED,
}

var DisabledGuildEmbed = &discordgo.MessageEmbed{
	Type:        discordgo.EmbedTypeRich,
	Title:       "Guild Disabled",
	Description: "Bot usage in this guild has been disabled. Contact the admin for help.",
	Color:       RED,
}
