package random2char

import (
	"log"
	"os"
	"time"
)

func init() {

	bot := TweetBot{
		ConsumerKey:    os.Getenv("RANDOM2CHARLOOP_APIKEY"),
		ConsumerSecret: os.Getenv("RANDOM2CHARLOOP_APIKEY_SECRET"),
		AccessToken:    os.Getenv("RANDOM2CHARLOOP_ACCESSTOKEN"),
		AccessSecret:   os.Getenv("RANDOM2CHARLOOP_ACCESSTOKEN_SECRET"),
	}
	bot.Init()
	ticker := time.Tick(time.Minute)

	go func() {
		for range ticker {
			tweets := MakeText()
			log.Println("tweeting in @2char_looping: " + tweets)
			if err := bot.Tweet(tweets); err != nil {
				log.Println("failed to tweeting: " + err.Error())
			}
		}
	}()
}
