package bot

type (
	WholeChats struct {
		MokuMoku VoiceConn
		Random   TextConn
		ToDo     TextConn
	}

	ChatConn interface {
		MakePrivate() error
		AllowAccess(memberId string) error
		Delete()
		GetID() string
	}

	TextConn interface {
		ChatConn
		Println(msg string)
	}

	VoiceConn interface {
		ChatConn
		MakeEveryoneMute(mute bool) error
		MoveToHere(memberId string) error
		JoinMemberIds() []string
	}
)
