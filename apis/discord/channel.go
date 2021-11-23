package discord

import (
	"app/bot"
	"encoding/binary"
	"errors"
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

type (
	Chat         discordgo.Channel
	TextChannel  Chat
	VoiceChannel struct {
		*Chat
		connection *discordgo.VoiceConnection
	}
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

func (tc *TextChannel) Println(msgArgs *bot.MsgArgs) {
	session.ChannelMessageSendEmbed(tc.ID, &discordgo.MessageEmbed{
		Type:        discordgo.EmbedTypeArticle,
		Title:       msgArgs.Title,
		Description: msgArgs.Description,
		Footer: &discordgo.MessageEmbedFooter{
			Text: msgArgs.Footer,
		},
	})
}

func (vc *VoiceChannel) GetID() string      { return vc.ID }
func (vc *VoiceChannel) MakePrivate() error { return vc.Chat.MakePrivate() }
func (vc *VoiceChannel) Delete()            { vc.Chat.Delete() }
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
		if everyone, err := vc.Chat.__FindEveryone(); err != nil {
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
		if everyone, err := vc.Chat.__FindEveryone(); err != nil {
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

func (vc *VoiceChannel) GetNumJoining() int {

	numMember := 0

	state.RLock()
	defer state.RUnlock()

	guild, err := state.Guild(vc.GuildID)
	if err != nil {
		return numMember
	}

	for i := range guild.VoiceStates {
		if guild.VoiceStates[i].ChannelID == vc.ID {
			numMember++
		}
	}

	return numMember
}

func (vc *VoiceChannel) PlaySound(path string) error {

	var (
		buffer []byte
	)

	if err := func() error {
		opuslen := int16(0)
		file, err := os.Open(path + ".dca")

		if err != nil {
			return errors.New("cannot open dca file: " + err.Error())
		}

		defer file.Close()

		if err := binary.Read(file, binary.LittleEndian, &opuslen); err != nil {
			return errors.New("cannot read dca file size: " + err.Error())
		}

		buffer = make([]byte, opuslen)
		if err = binary.Read(file, binary.LittleEndian, &buffer); err != nil {
			return errors.New("cannot read dca file: " + err.Error())
		}

		return nil

	}(); err != nil {
		return err
	}

	vc.connection.Speaking(true)
	vc.connection.OpusSend <- buffer
	vc.connection.Speaking(false)

	return nil
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
