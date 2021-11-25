package cheerleading

type (
	Profile struct {
		ConnectionIdPairs map[string]string `json:"connections"`
		IconFrom          string            `json:"icon-from"`
		NameFrom          string            `json:"name-from"`
		URLs              []string          `json:"urls"`
		iconURL, name     *string
	}

	Talkset [9][]Talk

	Talk struct {
		FileName string `json:"file"`
		Text     string `json:"text"`
	}

	TalkCommand int

	Cheerleader struct {
		Profile Profile `json:"profile"`
		Talkset Talkset `json:"voices"`
	}
)

var (
	Cheerleaders []Cheerleader
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
