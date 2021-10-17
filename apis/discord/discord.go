package discord

import (
	"app/apis"
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

const envKey = "DISCORD_TOKEN"

var (
	session *discordgo.Session
	state   *discordgo.State
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
	session.AddHandler(onMessageCreate)

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
	// TODO: CHECK COMMAND
}
