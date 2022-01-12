package mokumoku

import (
	"app/bot/cheerleading"
	"fmt"
	"time"
)

func (e *Event) MokuMoku() (isStopped bool) {

	footer := time.Now().Add(JST + MokuMokuMinute).Format("休憩時間は15:04頃からです！")
	e.Talk(cheerleading.MokuMokuBegin, MsgBeginMokuMoku.Description, footer, true)

	numMember := e.EventArgs.MokuMoku.GetNumJoining()
	e.GroupConn.SetStateMessage(fmt.Sprint(numMember+1, "人が作業中！"))

	halfCall := time.NewTimer(MokuMokuMinute / 2)
	endCall := time.NewTimer(MokuMokuMinute - 30*time.Second)
	timer := time.NewTimer(MokuMokuMinute)

	prev := time.Now()

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

						moveToHere := event.ToChatId == e.EventArgs.MokuMoku.GetID()
						moveFromHere := event.FromChatId == e.EventArgs.MokuMoku.GetID()
						numMember := e.EventArgs.MokuMoku.GetNumJoining()

						// join into voice channel
						if cur := time.Now(); moveToHere && !moveFromHere &&
							cur.Sub(prev) > 20*time.Second {

							prev = cur
							e.Talk(
								cheerleading.JoiningDuringMokuMoku,
								"今は作業時間です。張り切って作業をしましょう。",
								footer, false)
						}

						e.GroupConn.SetStateMessage(fmt.Sprint(numMember+1, "人が作業中！"))

						// check continue event
						return numMember < e.EventArgs.MinContinueMembers
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
