package twitter

import (
	"errors"

	"github.com/dghubble/go-twitter/twitter"
)

type Connect struct {
	client *twitter.Client
}

func GetIconUrl(id string) (string, error) {
	return Connection.GetIconUrl(id)
}

func (c *Connect) GetIconUrl(id string) (string, error) {
	if c.client == nil {
		return "", errors.New("???")
	}

	if user, _, err := c.client.Users.Show(&twitter.UserShowParams{
		ScreenName: id,
	}); err != nil {
		return "", errors.New("cannot get user infomation from twitter: " + err.Error())
	} else {
		return user.ProfileImageURL, nil
	}
}

func (c *Connect) GetName(id string) (string, error) {
	if user, _, err := c.client.Users.Show(&twitter.UserShowParams{
		ScreenName: id,
	}); err != nil {
		return "", errors.New("cannot get user infomation from twitter: " + err.Error())
	} else {
		return user.Name, nil
	}
}
