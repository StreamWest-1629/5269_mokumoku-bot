package cheerleading

type (
	Profile struct {
		ConnectionIdPairs map[string]string `json:"connections"`
		IconFrom          string            `json:"icon-from"`
		NameFrom          string            `json:"name-from"`
		URLs              []string          `json:"urls"`
	}

	Talkset [9][]Talk

	Talk struct {
		FileName string `json:"filename"`
		Text     string `json:"text"`
	}

	TalkCommand int

	VoiceBank struct {
		Profile Profile `json:"profile"`
		Talkset Talkset `json:"voices"`
	}
)

var (
	Voicebanks []VoiceBank
)

const (
	MokuMokuLaunch TalkCommand = iota
	MokuMokuBegin
	MokuMokuHalfDone
	MokuMokuMostlyEnded
	BreakingBegin
	JoiningDuringMokuMoku
	JoiningDuringBreaking
	Cheerleading
	Advertise
)
