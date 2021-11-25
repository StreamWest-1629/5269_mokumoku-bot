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

	i := len(e.EventArgs.MokuMoku.JoinMemberIds())
	prev := time.Now().Add(-20 * time.Second)

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
						if event.ToChatId == e.EventArgs.MokuMoku.GetID() {
							l := len(e.EventArgs.MokuMoku.JoinMemberIds())
							cur := time.Now()

							if l > i && cur.Sub(prev) > 20*time.Second {
								prev = cur
								e.Talk(
									cheerleading.JoiningDuringMokuMoku,
									"作業していた方は今休憩しています。しばらくお待ちください。",
									footer, false)
							}

							i = l
						}

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
