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

var (
	MsgBeginEvent = bot.MsgArgs{
		Title: "🔥 もくもく会をはじめます！",
		Description: "作業と休憩のルーチンを機械的に管理することでより効率的な作業の支援を行います。\n" +
			"作業中はボットが全員のマイクミュートを行うので静かに作業し、休憩中は3-5人のルームに振り分けて和気あいあいとお話しいただければと思います。",
	}
	MsgBeginMokuMoku = bot.MsgArgs{
		Title: "🚀 作業時間が始まります！",
		Description: "作業は52分間です。ボットが参加者のマイクミュートをします。\n" +
			"途中参加もできるのでぜひ来てください。",
	}
	MsgBeginBreaking = bot.MsgArgs{
		Title: "❤️ 休憩時間が始まります！",
		Description: "休憩は17分間です。ボットがメンバーの振り分けを行います。\n" +
			"休憩はプライベートチャットで行うので途中参加の方は次の作業時間が始まるまでそのままお待ちください。",
	}
	MsgEndEvent = bot.MsgArgs{
		Title: "😴 作業お疲れさまでした！",
		Description: "次回もぜひ好きな時間にボイスチャットに入ってもくもく会を始めてください。\n" +
			"ボットが24時間監視しているので、もくもく会は24時間行うことができます。\n" +
			"作業お疲れさまでした！",
	}
)

func init() {
	if _, exist := os.LookupEnv("DEBUGMODE"); exist {
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
	e.EventArgs.Random.Println(&MsgBeginEvent)
	for !e.routineOnce() {
	}
	e.onClose()
	e.EventArgs.Random.Println(&MsgEndEvent)
	fmt.Println("mokumoku event closed")
}

func (e *Event) routineOnce() (isClosed bool) {

	whole := e.EventArgs

	// mokumoku
	fmt.Println("Begin mokumoku time")
	msg := MsgBeginMokuMoku
	msg.Footer = time.Now().Add(JST + MokuMokuMinute).Format("休憩時間は15:04頃からです！")
	whole.Random.Println(&msg)
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
	msg = MsgBeginBreaking
	msg.Footer = time.Now().Add(JST + BreakingMinute).Format("作業時間は15:04頃からです！")
	whole.Random.Println(&msg)
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
