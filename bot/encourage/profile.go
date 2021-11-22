package encourage

type (
	CheerProfile struct {
		ConnectionIdPairs map[string]string `json:"connections"`
		IconFrom          string            `json:"icon-from"`
		NameFrom          string            `json:"name-from"`
		Link              string            `json:"link"`
	}

	CheerTalkset [9][]CheerTalk

	CheerTalk struct {
		FileName string `json:"filename"`
		Text     string `json:"text"`
	}
)
