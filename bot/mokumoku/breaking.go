package mokumoku

import (
	"app/bot"
	"app/bot/cheerleading"
	"fmt"
	"time"
)

func (e *Event) Breaking() bool {

	footer := time.Now().Add(JST + BreakingMinute).Format("作業時間は15:04頃からです！")
	e.Talk(cheerleading.BreakingBegin, MsgBeginBreaking.Description, footer, false)

	branches, err := bot.SpreadBranches(e.GroupConn, e.EventArgs)
	if err != nil {
		fmt.Println(err.Error())
		return true
	}

	timer := time.NewTimer(BreakingMinute)
	endCall := time.NewTimer(BreakingMinute - (30 * time.Second))

	prev := time.Now().Add(-20 * time.Second)

	for isContinue := true; isContinue; {
		select {
		case <-timer.C:
			isContinue = false
		case <-endCall.C:
			for i := range branches {
				branches[i].TextConn.Println(&MsgMostlyEndedBreaking)
			}

		case event := <-e.eventListener:
			if func() bool {
				defer event.Release()

				switch event := event.(type) {
				case onClose:
					return true
				case *onCheckMute:
					if _, exist := e.EventArgs.MuteIgnore[event.MemberId]; !exist {

						event.result <- event.ToChatId == e.EventArgs.MokuMoku.GetID()
						if cur := time.Now(); event.ToChatId == e.EventArgs.MokuMoku.GetID() && cur.Sub(prev) > 20*time.Second {
							prev = cur
							e.Talk(
								cheerleading.JoiningDuringBreaking,
								"作業していた方は今休憩しています。しばらくお待ちください。",
								footer, false)
						}

					} else {
						event.result <- false
					}
				}
				return false
			}() {
				return true
			}
		}
	}

	if err := branches.ClearBranches(e.GroupConn, e.EventArgs); err != nil {
		fmt.Println(err.Error())
		return true
	}

	return false
}
