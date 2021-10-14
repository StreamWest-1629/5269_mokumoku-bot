package main

import (
	"time"

	"app/apis"
	_ "app/apis/discord"
)

func main() {
	defer apis.Finalize()
	for {
		time.Sleep(time.Minute)
	}
}
