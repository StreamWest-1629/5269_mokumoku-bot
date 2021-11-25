package mokumoku

import (
	"app/bot"
	"app/bot/cheerleading"
	"log"
	"math/rand"
)

func (e *Event) Talk(cmd cheerleading.TalkCommand, description, footer string, forceUpdate bool) {

	cheerleader, talk := e.cheerleader.RandomTalkPath(cmd)
	if cheerleader != nil {
		e.cheerleader = cheerleader

		e.EventArgs.Random.Println(&bot.MsgArgs{
			Title:       talk.Text,
			Description: description,
			Footer:      footer,
			URL:         e.cheerleader.Profile.URLs[rand.Intn(len(e.cheerleader.Profile.URLs))],
			Authorname:  cheerleader.Profile.Name(forceUpdate),
			IconURL:     cheerleader.Profile.IconURL(forceUpdate),
		})
		if err := e.EventArgs.MokuMoku.Playsound(talk.FileName); err != nil {
			log.Println("cannot play sound: " + err.Error())
		}
	}
}
