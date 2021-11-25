package discord

import (
	"app/apis"
	"app/bot/mokumoku"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
)

const envKey = "DISCORD_TOKEN"
const idKey = "DISCORD_CLIENT"

var (
	session         *discordgo.Session
	state           *discordgo.State
	ownUserId       string
	mokumokuRunning = map[string]*mokumoku.Event{}
	lock            = sync.Mutex{}
)

func init() {

	fmt.Print("discord initializing...")

	if _, exist := os.LookupEnv("DEBUGMODE"); exist {
		CategoryName += "-DEV"
	}
	// initialize session
	if token, exist := os.LookupEnv(envKey); !exist {
		fmt.Println("cannot found token in environment values")
	} else if own, exist := os.LookupEnv(idKey); !exist {
		fmt.Println("cannot found crient id in environment values")
	} else if s, err := discordgo.New("Bot " + token); err != nil || s == nil {
		fmt.Println("cannot catch session: " + err.Error())
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
		fmt.Println("cannot open discord bot connection: " + err.Error())
	}

	// set finalizer
	apis.Finalizer = append(apis.Finalizer, func() {
		session.Close()
	})

	fmt.Println("ended!")
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
	lock.Lock()
	if event, running := mokumokuRunning[updated.GuildID]; running {
		lock.Unlock()
		before := ""
		if updated.BeforeUpdate != nil {
			before = updated.BeforeUpdate.ChannelID
		}

		// check mute
		if updated.ChannelID != "" {
			if mute := event.CheckMute(updated.UserID, before, updated.ChannelID); updated.Mute != mute {
				session.GuildMemberMute(updated.GuildID, updated.UserID, mute)
			}
		}

	} else if guild, exist := SearchGuild(updated.GuildID); exist {
		defer lock.Unlock()
		args := guild.GetWholeChats()

		if mokumoku.CheckLaunchEvent(args) {

			// voice join
			vc, _ := session.ChannelVoiceJoin(
				args.MokuMoku.(*VoiceChannel).GuildID,
				args.MokuMoku.(*VoiceChannel).ID, false, false)
			args.MokuMoku.(*VoiceChannel).conn = vc

			if ev := mokumoku.LaunchEvent(guild, args); ev != nil {
				// event begin to run
				ev.OnClose = func() {
					session.ChannelVoiceJoin(guild.guild.ID, "", false, false)
					delete(mokumokuRunning, guild.ID())
				}
				mokumokuRunning[guild.ID()] = ev

			} else {
				session.ChannelVoiceJoin(guild.guild.ID, "", false, false)
			}

		} else if updated.Mute {
			defer lock.Unlock()
			// bitween test and debug environment
			if ch, err := session.Channel(updated.ChannelID); err == nil && ch.Name != MokuMokuName {
				// release mute
				session.GuildMemberMute(updated.GuildID, updated.UserID, false)
			}
		}
	}
}
