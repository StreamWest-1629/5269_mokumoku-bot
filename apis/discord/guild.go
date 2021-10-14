package discord

import (
	"errors"
	"fmt"
	"sync"

	"github.com/bwmarrin/discordgo"
)

type (
	Guild struct {
		guild      *discordgo.Guild
		rooms      []MokuMoku
		categoryID string
		randomChat *TextChat
		todoChat   *TextChat
		wholeVoice *VoiceChat
	}

	MokuMoku struct {
		TextChat
		VoiceChat
	}
)

const (
	MokuMokuCategory   = "もくもくエリア"
	RandomChatName     = "おしゃべり"
	RandomChatTopic    = "おしゃべりをするところです。作業しながら思ったことや感じたことなどメモしておくのによいです。"
	ToDoChatName       = "すること"
	ToDoChatTopic      = "今日することをメモするところです。資格試験勉強から朝食を食べるまで、今日することをメモするのに使ってください。書くかどうかは任せます。"
	WholeVoiceChatName = "もくもく"
)

var (
	guilds     = map[string]*Guild{}
	guildsLock = sync.Mutex{}
)

func RegisterGuild(guildId string) (guild *Guild, err error) {

	// check error and make instance
	if _, exist := guilds[guildId]; exist {
		return nil, errors.New("guild instance has already registered")
	} else if g, err := session.Guild(guildId); err != nil {
		return nil, errors.New("cannot make discord guild instance: " + err.Error())
	} else {

		fmt.Println("make guild registed successfully")

		guild = &Guild{
			guild: g,
			rooms: []MokuMoku{},
		}

		// register
		guilds[guildId] = guild
	}

	if err := guild.InitializeChannels(); err != nil {
		fmt.Println(err.Error())
	}

	return guild, nil
}

func (g *Guild) InitializeChannels() (err error) {

	// update guild
	if g.guild, err = session.Guild(g.guild.ID); err != nil {
		return err
	} else if chs, err := session.GuildChannels(g.guild.ID); err != nil {
		return err
	} else {
		g.guild.Channels = append(g.guild.Channels, chs...)
	}

	// find category channel
	if ch, err := g.__findChannel(MokuMokuCategory, discordgo.ChannelTypeGuildCategory, ""); err != nil {
		return errors.New("cannot make category: " + err.Error())
	} else {
		g.categoryID = ch.ID
	}

	if ch, err := g.__findChannel(RandomChatName, discordgo.ChannelTypeGuildText, RandomChatTopic); err != nil {
		return errors.New("cannot make channel: " + err.Error())
	} else {
		g.randomChat = (*TextChat)(ch)
	}

	if ch, err := g.__findChannel(ToDoChatName, discordgo.ChannelTypeGuildText, ToDoChatTopic); err != nil {
		return errors.New("cannot make channel: " + err.Error())
	} else {
		g.todoChat = (*TextChat)(ch)
	}

	if ch, err := g.__findChannel(WholeVoiceChatName, discordgo.ChannelTypeGuildVoice, ""); err != nil {
		return errors.New("cannot make channel: " + err.Error())
	} else {
		g.wholeVoice = (*VoiceChat)(ch)
	}

	// todo: delete later
	session.ChannelMessageSend(g.randomChat.ID, "OK, Channels are found!")
	return nil
}

func (g *Guild) __findChannel(findName string, findType discordgo.ChannelType, topic string) (*discordgo.Channel, error) {
	switch findType {
	case discordgo.ChannelTypeGuildCategory:

		// find category
		for i := range g.guild.Channels {
			if ch := g.guild.Channels[i]; ch != nil {
				if ch.Type == findType {
					if ch.Name == findName {
						return ch, nil
					}
				}
			}
		}

		// make category
		if ch, err := session.GuildChannelCreateComplex(g.guild.ID, discordgo.GuildChannelCreateData{
			Name: findName,
			Type: findType,
		}); err != nil {
			return nil, errors.New("cannot make new category: " + err.Error())
		} else {
			return ch, nil
		}

	case discordgo.ChannelTypeGuildText, discordgo.ChannelTypeGuildVoice:

		// find channel
		for i := range g.guild.Channels {
			if ch := g.guild.Channels[i]; ch != nil {
				if ch.ParentID == g.categoryID && ch.Type == findType {
					if ch.Name == findName {
						return ch, nil
					}
				}
			}
		}

		// make channel
		if ch, err := session.GuildChannelCreateComplex(g.guild.ID, discordgo.GuildChannelCreateData{
			Name:     findName,
			Topic:    topic,
			Type:     findType,
			ParentID: g.categoryID,
		}); err != nil {
			return nil, errors.New("cannot make new channel: " + err.Error())
		} else {
			return ch, nil
		}
	}

	return nil, errors.New("unknown type")
}
