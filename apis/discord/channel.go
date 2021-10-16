package discord

import "github.com/bwmarrin/discordgo"

type (
	TextChat  discordgo.Channel
	VoiceChat discordgo.Channel
	TVChat    struct {
		textChat  TextChat
		voiceChat VoiceChat
	}
)
