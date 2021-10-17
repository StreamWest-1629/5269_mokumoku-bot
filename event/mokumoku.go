package event

import (
	"app/bot"
	"time"
)

type (
	MokuMoku struct {
		*bot.WholeChats
		VoiceStateChange chan struct {
			MemberId, ChatId string
		}
		voiceState map[string]int
		close      chan error
	}
)

const (
	MokuMokuMinute = 52 * time.Minute
	BreakingMinute = 17 * time.Minute
)

func LaunchMokuMoku(conn bot.GroupConn) (mokumoku *MokuMoku, err error) {

	// initialize
	mokumoku.WholeChats, err = conn.Initialize()
	if err != nil {
		return nil, err
	}
	mokumoku.VoiceStateChange = make(chan struct {
		MemberId string
		ChatId   string
	})
	members, err := conn.GetMemberIds()

	// set voice state

	// routine
	go func() {
		random := mokumoku.Random
		for {

			////////////////////////
			// breaking time
			////////////////////////
			branches, err := bot.SpreadBranch(conn)
			if err != nil {
				random.Println("failed to make branch rooms: " + err.Error())
			}
			time.Sleep(BreakingMinute)
			branches.ClearBranch()
		}
	}()

}
