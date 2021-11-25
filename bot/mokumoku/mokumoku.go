package mokumoku

import (
	"app/bot/cheerleading"
	"time"
)

func (e *Event) MokuMoku() bool {

	footer := time.Now().Add(JST + MokuMokuMinute).Format("休憩時間は15:04頃からです！")
	e.Talk(cheerleading.MokuMokuBegin, MsgBeginMokuMoku.Description, footer, true)

	halfCall := time.NewTimer(MokuMokuMinute / 2)
	endCall := time.NewTimer(MokuMokuMinute - 30*time.Second)
	timer := time.NewTimer(MokuMokuMinute)

	for isContinue := true; isContinue; {
		select {
		case <-timer.C:
			isContinue = false
		case <-halfCall.C:
			e.Talk(cheerleading.MokuMokuHalfDone, "", "", false)
		case <-endCall.C:
			e.Talk(cheerleading.MokuMokuMostlyEnded, "", "", false)
		case event := <-e.eventListener:
			if func() bool {
				defer event.Release()

				switch event := event.(type) {
				case onClose:
					return true
				case *onCheckMute:
					if _, exist := e.EventArgs.MuteIgnore[event.MemberId]; !exist {
						event.result <- event.ToChatId == e.EventArgs.MokuMoku.GetID()

						// check continue event
						return e.EventArgs.MokuMoku.GetNumJoining() < e.EventArgs.MinContinueMembers
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

	return false

}
