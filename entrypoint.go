package main

import (
	"fmt"
	"strconv"
	"time"

	"app/apis"
	_ "app/apis/discord"
)

func main() {
	defer apis.Finalize()
	fmt.Println("Hay! George!" + strconv.Itoa(len(apis.Finalizer)))
	for {
		time.Sleep(time.Minute)
	}
}
