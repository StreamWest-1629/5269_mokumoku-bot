package bot

type (
	MokuMokuRoomId int
	MsgFlag        int
	BotState       int

	GroupConn interface {
		InitializeChannels() error
		GetMokuMoku() (room VoiceChatConn, err error)
		MakeBranch(name string) (room Branch, err error)
		ClearBranch()
		Println(flag MsgFlag, msg string)
	}

	TextChatConn interface {
		SendMessage(msg string)
	}

	VoiceChatConn interface {
		MoveToHere(member MemberConn) error
	}

	MemberConn interface {
		GetID() string
		GetUsername() string
		UserMute(mute bool)
	}

	VoiceStateUpdate struct {
		MemberConn
		MokuMokuRoomId
	}
)

const (
	RoomIDDisconnect MokuMokuRoomId = -iota
	RoomIDMokuMoku
	RoomIDOther
)

const (
	MsgFlagEveryone MsgFlag = 1 << iota
	MsgFlagHere
	MsgFlagAdmin
	MsgFlagNone    MsgFlag = 0
	MsgFlagSendFor MsgFlag = MsgFlagEveryone | MsgFlagHere | MsgFlagAdmin
)

const (
	BotStateStopped BotState = iota
	BotStateMokuMoku
	BotStateBreaking
)

const (
	BotStateInitializing BotState = -(iota + 1)
)
