package bot

import (
	"errors"
	"sync"
	"time"
)

type (
	Event struct {
		GroupConn
		Rooms
		lock    sync.Mutex
		state   BotState
		changed chan VoiceStateUpdate
	}
)

const (
	MokuMokuMinutes = 52 * time.Minute
	BreakingMinutes = 17 * time.Minute
)

func Launch(conn GroupConn) (event *Event, err error) {
	if conn == nil {
		return nil, errors.New("conn argument is nil")
	}

	if mokumoku, err := conn.GetMokuMoku(); err != nil {
		return nil, errors.New("cannot find mokumoku voice chat: " + err.Error())
	} else {

		// initialize member
		return &Event{
			GroupConn: conn,
			state:     BotStateMokuMoku,
			Rooms: Rooms{
				Members:  Members{},
				MokuMoku: mokumoku,
				branches: []Branch{},
			},
		}, nil
	}
}
