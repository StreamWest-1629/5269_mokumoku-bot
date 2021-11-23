package main

import (
	"app/apis"
	_ "app/apis/discord"

	// _ "app/apis/twitter"
	"log"
	"os"

	"app/apis/heroku"
)

func init() {
	if _, exist := os.LookupEnv("DEBUGMODE"); exist {
		log.Println("[DEBUG MODE!]")
	} else {
		log.Println("[RELEASE MODE!]")
	}
}

func main() {
	defer apis.Finalize()
	heroku.HerokuRouter()
	// for {
	// 	time.Sleep(time.Hour)
	// }
}
