package discord

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

type (
	Chat         discordgo.Channel
	TextChannel  Chat
	VoiceChannel Chat
)

func (ch *Chat) MakePrivate() error {
	if g, exist := SearchGuild(ch.GuildID); !exist {
		return errors.New("cannot found member")
	} else if everyone, err := g.__findEveryone(); err != nil {
		return errors.New("cannot found role:" + err.Error())
	} else if err := session.ChannelPermissionSet(
		ch.ID,
		everyone.ID,
		discordgo.PermissionOverwriteTypeRole,
		0,
		discordgo.PermissionViewChannel,
	); err != nil {
		return errors.New("cannot make guild channel private: " + err.Error())
	} else {
		return nil
	}
}

func (ch *Chat) Delete() {
	session.ChannelDelete(ch.ID)
}
