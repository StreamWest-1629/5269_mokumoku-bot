package discord

import (
	"app/bot"
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type (
	Guild struct {
		guild      *discordgo.Guild
		categoryID string
		wholeCache *bot.WholeChats
	}
)

const (
	CategoryName = "もくもくエリア"
	MokuMokuName = "もくもく"
	RandomName   = "おしゃべり"
	RandomTopic  = "おしゃべりをするところです。作業しながら思ったことや感じたことなどメモしておくのによいです。"
	ToDoName     = "すること"
	ToDoTopic    = "今日することをメモするところです。資格試験勉強から朝食を食べるまで、今日することをメモするのに使ってください。書くかどうかは任せます。"
)

func SearchGuild(guildId string) (guild *Guild, exist bool) {

	// search guild
	g, err := state.Guild(guildId)
	var c *discordgo.Channel

	if err != nil {
		return nil, false
	} else if c, err = __findCategory(guildId, CategoryName); err != nil {
		if c, err = session.GuildChannelCreate(guildId, CategoryName, discordgo.ChannelTypeGuildCategory); err != nil {
			return nil, false
		} else {
			fmt.Println("discord guild's category added")
		}
	}

	return &Guild{
		guild:      g,
		categoryID: c.ID,
	}, true
}

func (g *Guild) ID() string {
	return g.guild.ID
}

func (g *Guild) Name() string {
	return g.guild.Name
}

func (g *Guild) MakeTextChat(name, topic string) (vc bot.TextConn, err error) {

	// make text chat
	if ch, err := session.GuildChannelCreateComplex(
		g.guild.ID,
		discordgo.GuildChannelCreateData{
			Name:     name,
			Type:     discordgo.ChannelTypeGuildText,
			Topic:    topic,
			ParentID: g.categoryID,
		}); err != nil {
		return nil, errors.New("failed to make a discord text chat: " + err.Error())
	} else {
		return (*TextChannel)(ch), nil
	}
}

func (g *Guild) MakeVoiceChat(name string) (vc bot.VoiceConn, err error) {
	// make voice chat
	if ch, err := session.GuildChannelCreateComplex(
		g.guild.ID,
		discordgo.GuildChannelCreateData{
			Name:     name,
			Type:     discordgo.ChannelTypeGuildVoice,
			ParentID: g.categoryID,
		}); err != nil {
		return nil, errors.New("failed to make a discord voice chat: " + err.Error())
	} else {
		return (*VoiceChannel)(ch), nil
	}
}

func (g *Guild) GetWholeChats() (whole *bot.WholeChats) {

	makeChan := func() (whole *bot.WholeChats) {
		if whole, err := g.__makeChannels(); err != nil {
			fmt.Println("cannot make whole chat instance: " + err.Error())
			return nil
		} else {
			return whole
		}
	}

	if g.wholeCache != nil {

		if _, err := state.Channel(g.wholeCache.MokuMoku.(*VoiceChannel).ID); err != nil {
			return makeChan()
		} else if _, err := state.Channel(g.wholeCache.Random.(*TextChannel).ID); err != nil {
			return makeChan()
		} else if _, err := state.Channel(g.wholeCache.ToDo.(*TextChannel).ID); err != nil {
			return makeChan()
		} else {
			return g.wholeCache
		}
	}

	if chats, err := g.__makeChannels(); err != nil {

	} else {
		return chats
	}

	return
}

func (g *Guild) MemberMute(memberId string, mute bool) {
	if err := session.GuildMemberMute(g.ID(), memberId, mute); err != nil {
		fmt.Println("cannot changed mute: " + err.Error())
	}
}

func (g *Guild) __makeChannels() (*bot.WholeChats, error) {

	var (
		MokuMoku *VoiceChannel = nil
		Random   *TextChannel  = nil
		ToDo     *TextChannel  = nil
	)

	state.RLock()
	defer state.RUnlock()

	// search channels
	for _, ch := range g.guild.Channels {
		switch ch.Type {
		case discordgo.ChannelTypeGuildVoice:
			switch ch.Name {
			case MokuMokuName:
				MokuMoku = (*VoiceChannel)(ch)
			}
		case discordgo.ChannelTypeGuildText:
			switch ch.Name {
			case RandomName:
				Random = (*TextChannel)(ch)
			case ToDoName:
				ToDo = (*TextChannel)(ch)
			}
		}
	}

	everyone, _ := g.__findEveryone()
	me := g.__findMe()

	// make channels
	if MokuMoku == nil {
		if ch, err := session.GuildChannelCreateComplex(g.ID(), discordgo.GuildChannelCreateData{
			Name:     MokuMokuName,
			Type:     discordgo.ChannelTypeGuildVoice,
			ParentID: g.categoryID,
			PermissionOverwrites: []*discordgo.PermissionOverwrite{
				{
					ID:   everyone.ID,
					Type: discordgo.PermissionOverwriteTypeRole,
					Allow: discordgo.PermissionViewChannel |
						discordgo.PermissionVoiceConnect,
					Deny: discordgo.PermissionVoiceSpeak,
				},
				{
					ID:    me.User.ID,
					Type:  discordgo.PermissionOverwriteTypeMember,
					Allow: discordgo.PermissionVoiceSpeak,
				},
			},
		}); err != nil {
			return nil, err
		} else {
			MokuMoku = (*VoiceChannel)(ch)
		}
	}

	if ToDo == nil {
		if ch, err := session.GuildChannelCreateComplex(g.ID(), discordgo.GuildChannelCreateData{
			Name:     ToDoName,
			Type:     discordgo.ChannelTypeGuildText,
			ParentID: g.categoryID,
		}); err != nil {
			return nil, err
		} else {
			ToDo = (*TextChannel)(ch)
		}
	}

	if Random == nil {
		if ch, err := session.GuildChannelCreateComplex(g.ID(), discordgo.GuildChannelCreateData{
			Name:     RandomName,
			Type:     discordgo.ChannelTypeGuildText,
			ParentID: g.categoryID,
		}); err != nil {
			return nil, err
		} else {
			Random = (*TextChannel)(ch)
		}
	}

	return &bot.WholeChats{
		MokuMoku: MokuMoku,
		Random:   Random,
		ToDo:     ToDo,
	}, nil
}

func (g *Guild) __findMe() *discordgo.Member {
	member, _ := session.GuildMember(g.ID(), ownUserId)
	return member
}

func (g *Guild) __findEveryone() (role *discordgo.Role, err error) {
	if roles, err := session.GuildRoles(g.ID()); err != nil {
		return nil, err
	} else {
		for i := range roles {
			if roles[i].Name == "@everyone" {
				return roles[i], nil
			}
		}
		return nil, errors.New("cannot found everyone role")
	}
}

func __findCategory(guildId, findName string) (*discordgo.Channel, error) {

	const Category = discordgo.ChannelTypeGuildCategory

	// get channels
	chs, err := session.GuildChannels(guildId)
	if err != nil {
		return nil, err
	}

	// search category
	for i := range chs {
		if chs[i] != nil {
			if ch := *chs[i]; ch.Type == Category && ch.Name == findName {
				return chs[i], nil
			}
		}
	}

	// make category
	if ch, err := session.GuildChannelCreateComplex(
		guildId,
		discordgo.GuildChannelCreateData{
			Name: findName,
			Type: Category,
		},
	); err != nil {
		return nil, errors.New("cannot make new category: " + err.Error())
	} else {
		return ch, nil
	}
}
