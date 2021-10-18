package mokumoku

import (
	"app/bot"
	"fmt"
	"sync"
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

	onClose       struct{}
	onVoiceUpdate struct {
		MemberId, FromChatId, ToChatId string
		Mute                           bool
		waitGroup                      sync.WaitGroup
	}
)

func (_ onClose) Release()       {}
func (v onVoiceUpdate) Release() { v.waitGroup.Done() }

const (
	MokuMokuMinute = 52 * time.Second
	BreakingMinute = 17 * time.Second
)

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

func (e *Event) VoiceUpdated(memberId, fromChatId, toChatId string, mute bool) {

	voice := &onVoiceUpdate{
		MemberId:   memberId,
		FromChatId: fromChatId,
		ToChatId:   toChatId,
		Mute:       mute,
	}

	voice.waitGroup.Add(1)
	e.eventListener <- voice
	voice.waitGroup.Wait()

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
				case *onVoiceUpdate:
					if event.ToChatId == whole.MokuMoku.GetID() {
						if !event.Mute {
							m.MemberMute(event.MemberId, true)
						}
					} else if event.FromChatId == whole.MokuMoku.GetID() {
						if event.Mute {
							m.MemberMute(event.MemberId, false)
						}
						fmt.Println("found move from mokumoku room")
						if len(whole.MokuMoku.JoinMemberIds()) <= 0 {
							return true
						}
					}
				}
				return false
			}() {
				return true
			}
		}
	}

	// breaking
	fmt.Println("Begin breaking time")
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
				case onVoiceUpdate:
					if event.ToChatId == whole.MokuMoku.GetID() {
						if !event.Mute {
							m.MemberMute(event.MemberId, true)
						}
					}
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
