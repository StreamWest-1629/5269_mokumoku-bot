package discord

import (
	"app/apis"
	"app/bot/mokumoku"
	"fmt"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const envKey = "DISCORD_TOKEN"
const idKey = "DISCORD_CLIENT"

var (
	session         *discordgo.Session
	state           *discordgo.State
	ownUserId       string
	mokumokuRunning = map[string]*mokumoku.Event{}
)

func init() {

	// initialize session
	if token, exist := os.LookupEnv(envKey); !exist {
		panic("cannot found token in environment values")
	} else if own, exist := os.LookupEnv(idKey); !exist {
		panic("cannot found crient id in environment values")
	} else if s, err := discordgo.New("Bot " + token); err != nil || s == nil {
		panic("cannot catch session: " + err.Error())
	} else {
		session, ownUserId = s, own
	}

	// set intents
	session.Identify.Intents = discordgo.IntentsAll

	// register event listener
	session.AddHandler(onMessageCreate)
	session.AddHandler(onVoiceStateUpdate)

	// make state availabled
	state = discordgo.NewState()
	state.TrackEmojis = false
	session.StateEnabled, session.State = true, state

	// open connection
	if err := session.Open(); err != nil {
		panic("cannot open discord bot connection: " + err.Error())
	}

	// set finalizer
	apis.Finalizer = append(apis.Finalizer, func() {
		session.Close()
	})

	fmt.Println("discord initialize successed")
}

func onMessageCreate(_ *discordgo.Session, created *discordgo.MessageCreate) {
	if strings.HasPrefix(created.Content, "/mokumoku update") {
		if _, running := mokumokuRunning[created.GuildID]; running {
			guild, _ := SearchGuild(created.GuildID)
			guild.GetWholeChats()
		}
	}
}

func onVoiceStateUpdate(_ *discordgo.Session, updated *discordgo.VoiceStateUpdate) {
	if event, running := mokumokuRunning[updated.GuildID]; running {
		before := ""
		fmt.Println("found running discord' voice state updates")
		fmt.Println("launch event: ", *event)
		if updated.BeforeUpdate != nil {
			before = updated.BeforeUpdate.ChannelID
		}
		event.VoiceUpdated(updated.UserID, before, updated.ChannelID, updated.Mute)
	} else if guild, exist := SearchGuild(updated.GuildID); exist {
		if ev := mokumoku.LaunchEvent(guild); ev != nil {
			mokumokuRunning[guild.ID()] = ev
			ev.OnClose = func() {
				delete(mokumokuRunning, guild.ID())
			}
		}
	}
}
