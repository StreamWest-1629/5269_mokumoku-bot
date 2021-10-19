package main

import (
	"app/apis"
	_ "app/apis/discord"
	// _ "app/apis/heroku"
)

func main() {
	defer apis.Finalize()
}
