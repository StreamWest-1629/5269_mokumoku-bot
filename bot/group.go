package bot

// Domain services about group (ex: discord) are defined here.

type (
	GroupConn interface {
		MakePrivateTextChat(name, topic string, allowMemberIds []string) (TextConn, error)
		MakePrivateVoiceChat(name string, allowMemberIds []string) (VoiceConn, error)
		MemberMute(memberId string, mute bool)
	}

	ChatConn interface {
		MakePrivate() error
		AllowAccess(memberId string) error
		Delete()
		GetID() string
	}

	TextConn interface {
		ChatConn
		Println(msgArgs *MsgArgs)
	}

	VoiceConn interface {
		ChatConn
		MakeEveryoneMute(mute bool) error
		MoveToHere(memberId string) error
		JoinMemberIds() []string
		GetNumJoining() int
		PlaySound(pathWithoutExt string) error
	}
)
