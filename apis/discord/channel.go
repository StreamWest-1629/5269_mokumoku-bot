package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type (
	Member discordgo.Member
	__Chat struct {
		*discordgo.Channel
		parent *Guild
	}
	TextChat  __Chat
	VoiceChat __Chat
)

const (
	RoleType = discordgo.PermissionOverwriteTypeRole
)

func (m *Member) MakeMute(mute bool) error {
	return session.GuildMemberMute(m.GuildID, m.User.ID, mute)
}

func (m *Member) GetID() string {
	return m.User.ID
}

func (ch *TextChat) Println(msgs ...interface{}) {
	session.ChannelMessageSend(ch.ID, fmt.Sprint(msgs...))

}

func (ch *__Chat) GetID() string {
	return ch.ID
}

func (ch *__Chat) MakePrivate() error {
	return session.ChannelPermissionSet(
		ch.ID,
		ch.parent.everyone.ID,
		discordgo.PermissionOverwriteTypeRole,
		0,
		discordgo.PermissionViewChannel,
	)
}

func (ch *TextChat) GetID() string      { return (*__Chat)(ch).GetID() }
func (ch *TextChat) MakePrivate() error { return (*__Chat)(ch).MakePrivate() }
func (ch *TextChat) MakeMemberAllow(memberId string) error {
	return session.ChannelPermissionSet(
		ch.ID, memberId,
		discordgo.PermissionOverwriteTypeMember,
		discordgo.PermissionViewChannel|
			discordgo.PermissionSendMessages,
		0,
	)
}

func (ch *VoiceChat) GetID() string      { return (*__Chat)(ch).GetID() }
func (ch *VoiceChat) MakePrivate() error { return (*__Chat)(ch).MakePrivate() }
func (ch *VoiceChat) MakeMemberAllow(memberId string) error {
	return session.ChannelPermissionSet(
		ch.ID, memberId,
		discordgo.PermissionOverwriteTypeMember,
		discordgo.PermissionViewChannel|
			discordgo.PermissionVoiceConnect|
			discordgo.PermissionVoiceSpeak,
		0,
	)
}

func (ch *VoiceChat) MoveToHere(memberId string) error {
	return session.ChannelVoiceJoinManual(ch.GuildID, memberId, false, false)
}
