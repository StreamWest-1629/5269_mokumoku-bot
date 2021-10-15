package bot

import (
	"log"
	"sync"
	"time"
)

type (
	Group struct {
		GroupRepository
		lock    sync.Mutex
		state   BotState
		changed chan VoiceStateUpdate
	}
)

const (
	MokuMokuMinutes = 52 * time.Minute
	BreakingMinutes = 17 * time.Minute
)

func (g *Group) LaunchMokuMoku() {

	// initialize channel
	if err := g.InitializeChannels(); err != nil {
		log.Println("cannot initialize channel: " + err.Error())
	}

	var (
		mokumokuVC *VoiceChat
	)

	// set global values
	if vc, err := g.GetMokuMoku(); err != nil {
		log.Println("cannot get mokumoku channel infomation: " + err.Error())
	} else {
		mokumokuVC = vc
	}

	// define some functions
	launchMokuMoku := func() error { // during mokumoku time
		if err := func() error {
			g.lock.Lock()
			defer g.lock.Unlock()

			// change mokumoku mode
			g.state = BotStateMokuMoku
			return nil

		}(); err != nil {
			return err
		}

		// mokumoku time
		time.Sleep(MokuMokuMinutes)
		return nil
	}

	launchBreaking := func() error { // during breaking time
		
		checkContinue := func() bool {
			return len(mokumokuVC.Joining())
		}

		// check continue task and initialize
		if ctn, err := func() (ctn bool, err error) {

			g.lock.Lock()
			defer g.lock.Unlock()

			// check continue
			// if ctn
		}
	}

}
