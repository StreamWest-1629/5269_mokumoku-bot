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

	fmt.Println("ended!")
	go http.ListenAndServe(":"+port, nil)
	client := http.DefaultClient

	time.Sleep(30 * time.Second)

	for {
		fmt.Print("heroku heart beat listening...")
		if _, err := client.Get(url + "/heroku/heartbeat"); err != nil {
			fmt.Println("failed to send heart beat: " + err.Error())
		}
		time.Sleep(20 * time.Minute)
	}
}
