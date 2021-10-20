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
)

const (
	MokuMokuBegining = "もくもく会さぎょう部はじめます！！頑張ってください！！"
	BreakingBegining = "もくもく会やすみ時間はじまります！！しっかりやすんで次のもくもくに備えましょう！！"
)

func init() {
	if _, exist := os.LookupEnv("DEBUG"); exist {
		MokuMokuMinute /= 60
		BreakingMinute /= 60
	}
}

func LaunchEvent(conn bot.GroupConn) *Event {

	if len(conn.GetWholeChats().MokuMoku.JoinMemberIds()) > 0 {

		event := (&Event{
			GroupConn:     conn,
			eventListener: make(chan __Event),
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

func (m *Event) routine() {
	for !m.routineOnce() {
	}
	m.onClose()
	fmt.Println("mokumoku event closed")
}

func (m *Event) routineOnce() (isClosed bool) {

	whole := m.GetWholeChats()

	// mokumoku
	fmt.Println("Begin mokumoku time")
	whole.Random.Println(MokuMokuBegining)
	timer := time.NewTimer(MokuMokuMinute)

	for i, members := 0, whole.MokuMoku.JoinMemberIds(); i < len(members); i++ {
		m.MemberMute(members[i], true)
	}

	for isContinue := true; isContinue; {
		select {
		case <-timer.C:
			isContinue = false
		case event := <-m.eventListener:
			if func() bool {
				defer event.Release()

				switch event := event.(type) {
				case onClose:
					return true
				case *onCheckMute:
					event.result <- event.ToChatId == whole.MokuMoku.GetID()
					return len(whole.MokuMoku.JoinMemberIds()) <= 0
				}
				return false
			}() {
				return true
			}
		}
	}

	if len(whole.MokuMoku.JoinMemberIds()) <= 0 {
		return true
	}

	// breaking
	fmt.Println("Begin breaking time")
	whole.Random.Println(BreakingBegining)
	timer = time.NewTimer(BreakingMinute)

	branches, err := bot.SpreadBranches(m.GroupConn)
	if err != nil {
		fmt.Println(err.Error())
		return true
	}

	for isContinue := true; isContinue; {
		select {
		case <-timer.C:
			isContinue = false
		case event := <-m.eventListener:
			if func() bool {
				defer event.Release()

				switch event := event.(type) {
				case onClose:
					return true
				case onCheckMute:
					event.result <- event.ToChatId == whole.MokuMoku.GetID()
				}
				return false
			}() {
				return true
			}
		}
	}

	if err := branches.ClearBranches(m.GroupConn); err != nil {
		fmt.Println(err.Error())
		return true
	}

	// check member
	return len(whole.MokuMoku.JoinMemberIds()) <= 0
}
