package discord

import (
	"app/bot"
	"errors"
	"sync"

	"github.com/bwmarrin/discordgo"
)

type (
	Guild struct {
		*discordgo.Guild
		*bot.BranchGroup
		categoryID string
		everyone   *discordgo.Role
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
		guild = &Guild{
			Guild: g,
		}
	}

	// make instance usecases
	if guild.everyone, err = guild.__findEveryone(); err != nil {
		return nil, errors.New("cannot find everyone role: " + err.Error())
	}

	// register
	guilds[guildId] = guild

	return guild, nil
}

func (g *Guild) Update() error {
	if chs, err := session.GuildChannels(g.ID); err != nil {
		return err
	} else {
		g.Channels = chs
		return nil
	}
}

func (g *Guild) Initialize() (whole *bot.WholeChats, err error) {

	// update guild
	g.Update()

	whole = &bot.WholeChats{}

	// find category channel
	if ch, err := g.__findChannel(MokuMokuCategory, discordgo.ChannelTypeGuildCategory, ""); err != nil {
		return nil, errors.New("cannot make category: " + err.Error())
	} else {
		g.categoryID = ch.ID
	}

	if ch, err := g.__findChannel(RandomChatName, discordgo.ChannelTypeGuildText, RandomChatTopic); err != nil {
		return nil, errors.New("cannot make channel: " + err.Error())
	} else {
		whole.Random = &TextChat{
			Channel: ch,
			parent:  g,
		}
	}

	if ch, err := g.__findChannel(ToDoChatName, discordgo.ChannelTypeGuildText, ToDoChatTopic); err != nil {
		return nil, errors.New("cannot make channel: " + err.Error())
	} else {
		whole.Todo = &TextChat{
			Channel: ch,
			parent:  g,
		}
	}

	if ch, err := g.__findChannel(WholeVoiceChatName, discordgo.ChannelTypeGuildVoice, ""); err != nil {
		return nil, errors.New("cannot make channel: " + err.Error())
	} else {
		whole.MokuMoku = &VoiceChat{
			Channel: ch,
			parent:  g,
		}
	}

	return whole, nil
}

func (g *Guild) VoiceState(memberId string) (chatId *string) {
	if state, err := state.VoiceState(g.ID, memberId); err != nil {
		return nil
	} else {
		return &state.ChannelID
	}
}

func (g *Guild) GetMember(userId string) (bot.MemberConn, error) {
	member, err := session.GuildMember(g.ID, userId)
	return (*Member)(member), err
}

func (g *Guild) GetMemberIds() ([]string, error) {

	members := []string{}

	for begins := ""; true; {
		if mem, err := session.GuildMembers(g.ID, begins, 1000); err != nil {
			return nil, errors.New("cannot get members: " + err.Error())
		} else {

			for i := range mem {
				members = append(members, mem[i].User.ID)
			}
			if len(mem) < 1000 {
				return members, nil
			} else {
				begins = mem[len(mem)-1].User.ID
			}
		}
	}
	return nil, errors.New("unexcepted function end")
}

func (g *Guild) MakeTextChat(name, description string) (bot.TextConn, error) {
	if ch, err := session.GuildChannelCreateComplex(g.ID, discordgo.GuildChannelCreateData{
		Name:     name,
		Topic:    description,
		Type:     discordgo.ChannelTypeGuildText,
		ParentID: g.categoryID,
	}); err != nil {
		return nil, errors.New("cannot make text chat: " + err.Error())
	} else {
		return &TextChat{
			Channel: ch,
			parent:  g,
		}, nil
	}
}

func (g *Guild) MakeVoiceChat(name string) (bot.VoiceConn, error) {
	if ch, err := session.GuildChannelCreateComplex(g.ID, discordgo.GuildChannelCreateData{
		Name:     name,
		Type:     discordgo.ChannelTypeGuildVoice,
		ParentID: g.categoryID,
	}); err != nil {
		return nil, errors.New("cannot make voice chat: " + err.Error())
	} else {
		return &VoiceChat{
			Channel: ch,
			parent:  g,
		}, nil
	}
}

func (g *Guild) RemoveChat(chatID string) {
	// find category
	for i := range g.Channels {
		if ch := g.Channels[i]; ch != nil {
			if ch.ID == chatID && ch.ParentID == g.categoryID {
				session.ChannelDelete(ch.ID)
			}
		}
	}
}

func (g *Guild) __findEveryone() (role *discordgo.Role, err error) {
	if roles, err := session.GuildRoles(g.ID); err != nil {
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

func (g *Guild) __findChannel(findName string, findType discordgo.ChannelType, topic string) (*discordgo.Channel, error) {
	switch findType {
	case discordgo.ChannelTypeGuildCategory:

		// find category
		for i := range g.Channels {
			if ch := g.Channels[i]; ch != nil {
				if ch.Type == findType {
					if ch.Name == findName {
						return ch, nil
					}
				}
			}
		}

		// make category
		if ch, err := session.GuildChannelCreateComplex(g.ID, discordgo.GuildChannelCreateData{
			Name: findName,
			Type: findType,
		}); err != nil {
			return nil, errors.New("cannot make new category: " + err.Error())
		} else {
			return ch, nil
		}

	case discordgo.ChannelTypeGuildText, discordgo.ChannelTypeGuildVoice:

		// find channel
		for i := range g.Channels {
			if ch := g.Channels[i]; ch != nil {
				if ch.ParentID == g.categoryID && ch.Type == findType {
					if ch.Name == findName {
						return ch, nil
					}
				}
			}
		}

		// make channel
		if ch, err := session.GuildChannelCreateComplex(g.ID, discordgo.GuildChannelCreateData{
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
