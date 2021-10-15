package bot

type (
	MokuMokuRoomId int
	MsgFlag        int
	BotState       int

	GroupRepository interface {
		InitializeChannels() error
		GetMokuMoku() (room VoiceChat, err error)
		MakeBranch(name string) (room Branch, err error)
		ClearBranch()
		Println(flag MsgFlag, msg string)
	}

	TextChat interface {
		SendMessage(msg string)
	}

	VoiceChat interface {
		MoveToHere(member Member) error
		Joining() []Member
	}

	Member interface {
		GetUsername() string
		UserMute(mute bool)
	}

	Branch interface {
		GetMokuMokuId() MokuMokuRoomId
		GetName() string
		TextChat
		VoiceChat
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
