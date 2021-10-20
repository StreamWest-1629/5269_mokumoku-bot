package main

import (
	"app/apis"
	_ "app/apis/discord"
	"fmt"
	"os"

	"app/apis/heroku"
)

func init() {
	if _, exist := os.LookupEnv("DEBUG"); exist {
		fmt.Println("[DEBUG MODE!]")
	}
}

func main() {
	defer apis.Finalize()
	heroku.HerokuRouter()
	// for {
	// 	time.Sleep(time.Hour)
	// }
}
