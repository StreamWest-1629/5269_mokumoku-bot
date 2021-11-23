package cheerleading

import (
	"math/rand"
)

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
