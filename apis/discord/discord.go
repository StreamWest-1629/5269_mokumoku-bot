package discord

import (
	"app/apis"
	"fmt"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const envKey = "DISCORD_TOKEN"

var (
	session *discordgo.Session
)

func init() {

	// initialize session
	if token, exist := os.LookupEnv(envKey); !exist {
		panic("cannot found token in environment values")
	} else if s, err := discordgo.New("Bot " + token); err != nil || s == nil {
		panic("cannot catch session: " + err.Error())
	} else {
		session = s
	}

	// register event listener
	session.AddHandler(onVoiceStateUpdate)
	session.AddHandler(onMessageCreate)

	// open connection
	if err := session.Open(); err != nil {
		panic("cannot open discord bot connection: " + err.Error())
	} else {
		fmt.Println("discord initialize successed")

		// set finalizer
		apis.Finalizer = append(apis.Finalizer, func() {
			session.Close()
		})
	}
}

func onVoiceStateUpdate(_ *discordgo.Session, updated *discordgo.VoiceStateUpdate) {
	fmt.Println("voice state changed")
	checkGuildRegistered(updated.GuildID)
}

func onMessageCreate(_ *discordgo.Session, updated *discordgo.MessageCreate) {
	fmt.Println("message arrived")
	guild := checkGuildRegistered(updated.GuildID)

	if !updated.Author.Bot && !updated.Author.System {
		if strings.Index(updated.Content, "/mokumoku update") != -1 {
			guild.InitializeChannels()
		}
	}
}

func checkGuildRegistered(guildId string) *Guild {
	guildsLock.Lock()
	defer guildsLock.Unlock()

	if guild, exist := guilds[guildId]; exist {
		return guild
	} else if guild, err := RegisterGuild(guildId); err != nil {
		fmt.Println("cannot make guild registered: " + err.Error())
		return nil
	} else {
		return guild
	}
}
