package bot

type (
	MemberConn interface {
		MakeMute(mute bool) error
		GetID() string
	}

	ChatConn interface {
		GetID() string
		MakePrivate() error
		MakeMemberAllow(memberId string) error
	}

	TextConn interface {
		ChatConn
		Println(msgs ...interface{})
	}

	VoiceConn interface {
		ChatConn
		MoveToHere(memberId string) error
	}

	VoiceChat struct {
		VoiceConn
		parent *BranchGroup
	}
)
