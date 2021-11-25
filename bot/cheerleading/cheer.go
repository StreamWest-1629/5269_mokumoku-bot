package cheerleading

import (
	"app/bot"
	"log"
	"math/rand"
)

var Connections = map[string]bot.ConnectConn{}

func RandomCheerleader() *Cheerleader {
	return &Cheerleaders[rand.Intn(len(Cheerleaders))]
}

func (cheerleader *Cheerleader) RandomTalkPath(cmd TalkCommand) (cheerleading *Cheerleader, talk Talk) {
	if l := len(cheerleader.Talkset[cmd]); l > 0 {
		return cheerleader, cheerleader.Talkset[cmd][rand.Intn(l)]
	} else {
		return RandomCheerleader().RandomTalkPath(cmd)
	}
}

func (cheerleader *Profile) IconURL(forceUpdate bool) string {
	if forceUpdate || cheerleader.iconURL == nil {
		if url, err := Connections[cheerleader.IconFrom].GetIconUrl(cheerleader.ConnectionIdPairs[cheerleader.IconFrom]); err != nil {
			log.Println("cannot get icon url: " + err.Error())
			return ""
		} else {
			return url
		}
	} else {
		return *cheerleader.iconURL
	}
}

func (cheerleader *Profile) Name(forceUpdate bool) string {
	if forceUpdate || cheerleader.name == nil {
		if name, err := Connections[cheerleader.NameFrom].GetName(cheerleader.ConnectionIdPairs[cheerleader.NameFrom]); err != nil {
			log.Println("cannot get name: " + err.Error())
			return ""
		} else {
			return name
		}
	} else {
		return *cheerleader.name
	}
}
