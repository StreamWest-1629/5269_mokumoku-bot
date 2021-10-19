package main

import (
	"app/apis"
	_ "app/apis/discord"

	"app/apis/heroku"
)

func main() {
	defer apis.Finalize()
	heroku.HerokuRouter()
	// for {
	// 	time.Sleep(time.Hour)
	// }
}
