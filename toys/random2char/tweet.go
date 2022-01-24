package random2char

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/dghubble/oauth1"
)

type TweetBot struct {
	*http.Client
	ConsumerKey    string
	ConsumerSecret string
	AccessToken    string
	AccessSecret   string
}

func (keys *TweetBot) Init() {
	config := oauth1.NewConfig(keys.ConsumerKey, keys.ConsumerSecret)
	token := oauth1.NewToken(keys.AccessToken, keys.AccessSecret)
	keys.Client = config.Client(oauth1.NoContext, token)
}

func (bot *TweetBot) Tweet(msg string) error {
	b, _ := json.Marshal(map[string]string{"text": msg})
	_, err := bot.Post("https://api.twitter.com/2/tweets", "application/json", bytes.NewBuffer(b))
	return err
}
