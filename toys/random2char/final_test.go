package random2char

import (
	"log"
	"os"
	"testing"
	"time"
)

func TestFinal(t *testing.T) {

	t.Setenv("RANDOM2CHARLOOP_APIKEY", "b5MMmmZ1M3sySBVzUbHI0pIxP")
	t.Setenv("RANDOM2CHARLOOP_APIKEY_SECRET", "Tkh3WIWgul5NWVHeRyPbwGSOmCSECHzDoI4uXt6sScs75TWaAY")
	t.Setenv("RANDOM2CHARLOOP_ACCESSTOKEN", "1485574317471842309-DH52Ujg2OdFjAml2lQTgkj1UU7maQT")
	t.Setenv("RANDOM2CHARLOOP_ACCESSTOKEN_SECRET", "6WYnJDHMK8oLcJzLo02of6iUTXj43Vn3NmK4UNf0rJ8JQ")

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

	for {
		time.Sleep(10 * time.Minute)
	}
}
