package heroku

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func HerokuRouter() {

	fmt.Println("heroku initializing...")

	port, _ := os.LookupEnv("PORT")
	url, _ := os.LookupEnv("HEROKU_URL")

	http.HandleFunc("/heroku/heartbeat", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("heart beat listened")
		w.WriteHeader(200)
	})

	fmt.Print("begin server listening...")
	go func() {
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			fmt.Println(err.Error())
		}
	}()
	fmt.Println("ended!")
	client := http.DefaultClient
	client.Timeout = 10 * time.Second

	url += "heroku/heartbeat"
	fmt.Println("heroku initializing ended!")

	time.Sleep(30 * time.Second)

	for {
		fmt.Print("heroku heart beat listening(" + url + ")...")
		if _, err := client.Get(url); err != nil {
			fmt.Println("failed to send heart beat: " + err.Error())
		}
		time.Sleep(20 * time.Minute)
	}
}
