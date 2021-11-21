package bot

type (
	EventArgs struct {
		MokuMoku           VoiceConn
		Random             TextConn
		ToDo               TextConn
		BranchIgnore       map[string]interface{}
		MuteIgnore         map[string]interface{}
		MinLaunchMembers   int
		MinContinueMembers int
	}

	MsgArgs struct {
		Title, Description, Footer string
	}
)
