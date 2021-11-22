package twitter

import (
	"log"
	"testing"
)

func TestTwitter(t *testing.T) {

	if img, err := Connection.GetIconUrl("streamwest1629"); err != nil {
		log.Println(err.Error())
	} else {
		log.Println(img)
	}
}
