package discord

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type (
	Chat         discordgo.Channel
	TextChannel  Chat
	VoiceChannel Chat
)

func (ch *Chat) MakePrivate() error {

	if everyone, err := ch.__FindEveryone(); err != nil {
		return err
	} else if err := session.ChannelPermissionSet(
		ch.ID,
		ownUserId,
		discordgo.PermissionOverwriteTypeMember,
		discordgo.PermissionViewChannel,
		0,
	); err != nil {

	} else if err := session.ChannelPermissionSet(
		ch.ID,
		everyone.ID,
		discordgo.PermissionOverwriteTypeRole,
		0,
		discordgo.PermissionViewChannel,
	); err != nil {
		return errors.New("cannot make guild channel private: " + err.Error())
	}

	return nil
}

func (ch *Chat) Delete() {
	if _, err := session.ChannelDelete(ch.ID); err != nil {
		fmt.Println(err.Error())
	}
}

func (tc *TextChannel) GetID() string      { return tc.ID }
func (tc *TextChannel) MakePrivate() error { return (*Chat)(tc).MakePrivate() }
func (tc *TextChannel) Delete()            { (*Chat)(tc).Delete() }
func (tc *TextChannel) AllowAccess(memberId string) error {
	return session.ChannelPermissionSet(
		tc.ID,
		memberId,
		discordgo.PermissionOverwriteTypeMember,
		discordgo.PermissionViewChannel|discordgo.PermissionSendMessages,
		0,
	)
}

func (tc *TextChannel) Println(msg string) {
	session.ChannelMessageSend(tc.ID, msg)
}

func (vc *VoiceChannel) GetID() string      { return vc.ID }
func (vc *VoiceChannel) MakePrivate() error { return (*Chat)(vc).MakePrivate() }
func (vc *VoiceChannel) Delete()            { (*Chat)(vc).Delete() }
func (vc *VoiceChannel) AllowAccess(memberId string) error {
	return session.ChannelPermissionSet(
		vc.ID,
		memberId,
		discordgo.PermissionOverwriteTypeMember,
		discordgo.PermissionViewChannel|discordgo.PermissionVoiceConnect|discordgo.PermissionVoiceSpeak,
		0,
	)
}

func (vc *VoiceChannel) MakeEveryoneMute(mute bool) error {
	if mute {
		if everyone, err := (*Chat)(vc).__FindEveryone(); err != nil {
			return err
		} else if err := session.ChannelPermissionSet(
			vc.ID,
			everyone.ID,
			discordgo.PermissionOverwriteTypeRole,
			0,
			discordgo.PermissionVoiceSpeak,
		); err != nil {
			return errors.New("cannot make guild channel private: " + err.Error())
		} else {
			return nil
		}
	} else {
		if everyone, err := (*Chat)(vc).__FindEveryone(); err != nil {
			return err
		} else if err := session.ChannelPermissionSet(
			vc.ID,
			everyone.ID,
			discordgo.PermissionOverwriteTypeRole,
			discordgo.PermissionVoiceSpeak,
			0,
		); err != nil {
			return errors.New("cannot make guild channel private: " + err.Error())
		} else {
			return nil
		}
	}
}

func (vc *VoiceChannel) MoveToHere(memberId string) error {
	return session.GuildMemberMove(vc.GuildID, memberId, &vc.ID)
}

func (vc *VoiceChannel) JoinMemberIds() []string {
	members := []string{}

	state.RLock()
	defer state.RUnlock()

	guild, err := state.Guild(vc.GuildID)
	if err != nil {
		return members
	}

	for i := range guild.VoiceStates {
		if guild.VoiceStates[i].ChannelID == vc.ID {
			members = append(members, guild.VoiceStates[i].UserID)
		}
	}

	return members
}

func (ch *Chat) __FindEveryone() (g *discordgo.Role, err error) {
	if g, exist := SearchGuild(ch.GuildID); !exist {
		return nil, errors.New("cannot found member")
	} else if everyone, err := g.__findEveryone(); err != nil {
		return nil, errors.New("cannot found role:" + err.Error())
	} else {
		return everyone, nil
	}
}
