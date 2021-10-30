package mokumoku

import (
	"app/bot"
	"fmt"
	"os"
	"time"
)

type (
	Event struct {
		bot.GroupConn
		*bot.EventArgs
		eventListener chan __Event
		OnClose       func()
	}

	__Event interface {
		Release()
	}

	onClose     struct{}
	onCheckMute struct {
		MemberId, FromChatId, ToChatId string
		result                         chan bool
	}
)

func (onClose) Release()     {}
func (onCheckMute) Release() {}

var (
	MokuMokuMinute = 52 * time.Minute
	BreakingMinute = 17 * time.Minute
	JST            = 9 * time.Hour
)

const (
	MokuMokuExplain = "```\n" +
		"このもくもく会はボットが管理します。\n" +
		"はじめの52分がもくもく時間でミュートとなり、つぎの17分が休憩時間で3-5人のルームに振り分けられる、というルーチンで進められます。\n" +
		"ボイスチャットの移動を含め、すべてボットが管理するので安心して作業（52分）と休憩（17分）のルーチンをお楽しみください！！\n" +
		"```"
	MokuMokuBegining = "もくもく会さぎょう部はじめます！！予定時間は15時04分頃までです！！頑張ってください！！"
	BreakingBegining = "もくもく会やすみ時間はじまります！！予定時間は15時04分頃までです！！しっかりやすんで次のもくもくに備えましょう！！"
	MokuMokuEnded    = "さぎょうお疲れさまでした！！また是非いらしてください！！"
)

func init() {
	if _, exist := os.LookupEnv("DEBUG"); exist {
		MokuMokuMinute /= 60
		BreakingMinute /= 60
	}
}

func LaunchEvent(conn bot.GroupConn, whole *bot.EventArgs) *Event {

	if len(whole.MokuMoku.JoinMemberIds()) > 0 {

		event := (&Event{
			GroupConn:     conn,
			eventListener: make(chan __Event),
			EventArgs:     whole,
		})
		go event.routine()

		return event
	} else {
		return nil
	}

}

func (e *Event) Close() {
	e.eventListener <- onClose{}
}

func (e *Event) CheckMute(memberId, fromChatId, toChatId string) bool {

	voice := &onCheckMute{
		MemberId:   memberId,
		FromChatId: fromChatId,
		ToChatId:   toChatId,
		result:     make(chan bool, 1),
	}

	defer close(voice.result)

	e.eventListener <- voice
	return <-voice.result

}

func (e *Event) onClose() {
	e.OnClose()
}

func (e *Event) routine() {
	e.EventArgs.Random.Println(MokuMokuExplain)
	for !e.routineOnce() {
	}
	e.onClose()
	e.EventArgs.Random.Println(MokuMokuEnded)
	fmt.Println("mokumoku event closed")
}

func (e *Event) routineOnce() (isClosed bool) {

	whole := e.EventArgs

	// mokumoku
	fmt.Println("Begin mokumoku time")
	whole.Random.Println(time.Now().Add(JST + MokuMokuMinute).Format(MokuMokuBegining))
	timer := time.NewTimer(MokuMokuMinute)

	for i, members := 0, whole.MokuMoku.JoinMemberIds(); i < len(members); i++ {
		if _, exist := whole.MuteIgnore[members[i]]; !exist {
			e.MemberMute(members[i], true)
		}
	}

	for isContinue := true; isContinue; {
		select {
		case <-timer.C:
			isContinue = false
		case event := <-e.eventListener:
			if func() bool {
				defer event.Release()

				switch event := event.(type) {
				case onClose:
					return true
				case *onCheckMute:
					if _, exist := whole.MuteIgnore[event.MemberId]; !exist {
						event.result <- event.ToChatId == whole.MokuMoku.GetID()

						// check continue event
						return whole.MokuMoku.GetNumJoining() < whole.MinContinueMembers
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

	if len(whole.MokuMoku.JoinMemberIds()) < whole.MinContinueMembers {
		return true
	}

	// breaking
	fmt.Println("Begin breaking time")
	whole.Random.Println(time.Now().Add(JST + BreakingMinute).Format(BreakingBegining))
	timer = time.NewTimer(BreakingMinute)

	branches, err := bot.SpreadBranches(e.GroupConn, whole)
	if err != nil {
		fmt.Println(err.Error())
		return true
	}

	for isContinue := true; isContinue; {
		select {
		case <-timer.C:
			isContinue = false
		case event := <-e.eventListener:
			if func() bool {
				defer event.Release()

				switch event := event.(type) {
				case onClose:
					return true
				case *onCheckMute:
					if _, exist := whole.MuteIgnore[event.MemberId]; !exist {
						event.result <- event.ToChatId == whole.MokuMoku.GetID()
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

	if err := branches.ClearBranches(e.GroupConn, whole); err != nil {
		fmt.Println(err.Error())
		return true
	}

	// check member
	return len(whole.MokuMoku.JoinMemberIds()) < whole.MinContinueMembers
}
