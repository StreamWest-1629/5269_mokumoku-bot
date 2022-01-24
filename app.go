package main

import (
	"app/apis"
	_ "app/apis/discord"
	_ "app/apis/twitter"
	_ "app/toys/random2char"
	"fmt"
	"os"

	"app/apis/heroku"
)

func init() {
	if _, exist := os.LookupEnv("DEBUGMODE"); exist {
		fmt.Println("[DEBUG MODE!]")
	}
}

func main() {
	defer apis.Finalize()
	heroku.HerokuRouter()
}
