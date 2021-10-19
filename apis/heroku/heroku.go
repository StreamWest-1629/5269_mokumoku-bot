package heroku

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func init() {

	fmt.Print("heroku initializing...")

	port, _ := os.LookupEnv("PORT")
	url, _ := os.LookupEnv("HEROKU_URL")

	http.HandleFunc("/heroku/heartbeat", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("heart beat listened")
		w.WriteHeader(200)
	})

	go fmt.Println("server listen ended: " + http.ListenAndServe(":"+port, nil).Error())
	client := http.DefaultClient

	time.Sleep(30 * time.Second)
	url += "heroku/heartbeat"

	fmt.Println("ended!")

	for {
		fmt.Print("heroku heart beat listening(" + url + ")...")
		if _, err := client.Get(url); err != nil {
			fmt.Println("failed to send heart beat: " + err.Error())
		}
		time.Sleep(20 * time.Minute)
	}
}
