package main

import (
	"app/apis"
	_ "app/apis/discord"

	// _ "app/apis/heroku"
	"time"
)

func main() {
	defer apis.Finalize()
	// sc := make(chan os.Signal, 1)
	// // signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	// <-sc
	for {
		time.Sleep(time.Hour)
	}
}
