package mokumoku

import (
	"app/bot"
	"app/bot/cheerleading"
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
		cheerleader   *cheerleading.Cheerleader
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

func CheckLaunchEvent(whole *bot.EventArgs) bool {
	return len(whole.MokuMoku.JoinMemberIds()) > 0
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

	for i, members := 0, e.EventArgs.MokuMoku.JoinMemberIds(); i < len(members); i++ {
		if _, exist := e.EventArgs.MuteIgnore[members[i]]; !exist {
			e.MemberMute(members[i], true)
		}
	}

	e.cheerleader = cheerleading.RandomCheerleader()
	e.Talk(cheerleading.MokuMokuLaunch, MsgBeginEvent.Description, "", true)

	for !e.routineOnce() {
		e.cheerleader = cheerleading.RandomCheerleader()
	}

	e.onClose()
	e.EventArgs.Random.Println(&MsgEndEvent)
	fmt.Println("mokumoku event closed")
}

func (e *Event) routineOnce() (isClosed bool) {

	whole := e.EventArgs

	for i, members := 0, e.EventArgs.MokuMoku.JoinMemberIds(); i < len(members); i++ {
		if _, exist := e.EventArgs.MuteIgnore[members[i]]; !exist {
			e.MemberMute(members[i], true)
		}
	}

	if e.MokuMoku() {
		return true
	}

	if len(whole.MokuMoku.JoinMemberIds()) < whole.MinContinueMembers {
		return true
	}

	// breaking
	if e.Breaking() {
		return true
	}

	// check member
	return len(whole.MokuMoku.JoinMemberIds()) < whole.MinContinueMembers
}
