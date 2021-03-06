package mokumoku

import (
	"app/bot"
	"app/bot/cheerleading"
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
		Title: "đĨ ããããäŧãã¯ãããžãīŧ",
		Description: "äŊæĨ­ã¨äŧæŠãŽãĢãŧããŗãæŠæĸ°įãĢįŽĄįãããã¨ã§ããåšįįãĒäŊæĨ­ãŽæ¯æ´ãčĄããžãã\n" +
			"äŊæĨ­ä¸­ã¯ããããå¨åĄãŽãã¤ã¯ããĨãŧããčĄããŽã§éããĢäŊæĨ­ããäŧæŠä¸­ã¯3-5äēēãŽãĢãŧã ãĢæ¯ãåããĻåæ°ããããã¨ãčŠąãããã ããã°ã¨æããžãã",
	}
	MsgBeginMokuMoku = bot.MsgArgs{
		Title: "đ äŊæĨ­æéãå§ãžããžãīŧ",
		Description: "äŊæĨ­ã¯52åéã§ããããããåå čãŽãã¤ã¯ããĨãŧããããžãã\n" +
			"éä¸­åå ãã§ãããŽã§ãã˛æĨãĻãã ããã",
	}
	MsgBeginBreaking = bot.MsgArgs{
		Title: "â¤ī¸ äŧæŠæéãå§ãžããžãīŧ",
		Description: "äŧæŠã¯17åéã§ãããããããĄãŗããŧãŽæ¯ãåããčĄããžãã\n" +
			"äŧæŠã¯ããŠã¤ããŧãããŖããã§čĄããŽã§éä¸­åå ãŽæšã¯æŦĄãŽäŊæĨ­æéãå§ãžããžã§ããŽãžãžãåžãĄãã ããã",
	}
	MsgEndEvent = bot.MsgArgs{
		Title: "đ´ äŊæĨ­ãį˛ãããžã§ããīŧ",
		Description: "æŦĄåããã˛åĨŊããĒæéãĢãã¤ãšããŖãããĢåĨãŖãĻããããäŧãå§ããĻãã ããã\n" +
			"ãããã24æéįŖčĻããĻãããŽã§ãããããäŧã¯24æéčĄããã¨ãã§ããžãã\n" +
			"äŊæĨ­ãį˛ãããžã§ããīŧ",
	}
	MsgMostlyEndedBreaking = bot.MsgArgs{
		Title: "äŧæŠæéįĩäē30į§åã§ã",
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
	e.EventArgs.Random.Println(&MsgEndEvent)
	e.GroupConn.SetStateMessage("ããããã¸ãŠãã")
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

	defer e.onClose()

	for !e.routineOnce() {
		e.cheerleader = cheerleading.RandomCheerleader()
	}
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
