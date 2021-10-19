package main

import (
	"net"
	"os"
	"time"

	"app/apis"
	_ "app/apis/discord"
)

func main() {
	defer apis.Finalize()
	port, _ := os.LookupEnv("PORT")
	net.Listen("tcp", ":"+port)

	for {
		time.Sleep(time.Minute)
	}
}
