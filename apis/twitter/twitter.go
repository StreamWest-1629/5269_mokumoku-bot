package twitter

import (
	"app/bot/cheerleading"
	"fmt"
	"log"
	"os"

	"github.com/dghubble/go-twitter/twitter"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

var (
	Connection Connect
)

const (
	keyKey    = "TWITTER_KEY"
	secretKey = "TWITTER_SECRET"
	authURL   = "https://api.twitter.com/oauth2/token"
)

func init() {

	// launch twitter client
	if key, exist := os.LookupEnv(keyKey); !exist {
		fmt.Println("cannot found key in environment values")
	} else if secret, exist := os.LookupEnv(secretKey); !exist {
		fmt.Println("cannot found secret in environment values")
	} else {
		config := &clientcredentials.Config{
			ClientID:     key,
			ClientSecret: secret,
			TokenURL:     authURL,
		}

		if httpClient := config.Client(oauth2.NoContext); httpClient == nil {
			log.Fatalln("cannot login twitter API")
		} else if client := twitter.NewClient(httpClient); client == nil {
			log.Fatalln("cannot make twitter API client instance")
		} else if client.Users == nil {
			log.Fatalln("cannot make twitter User API client instance")
		} else {
			Connection = Connect{
				client: client,
			}

			cheerleading.Connections["twitter"] = &Connection

		}
	}
}
