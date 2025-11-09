package gh

import (
	"github.com/google/go-github/v77/github"
)

type Client struct {
	Client *github.Client
}

func NewClientWrapper(token string) *Client {
	if token == "" {
		return &Client{
			Client: github.NewClient(nil),
		}
	}

	return &Client{
		Client: github.NewClient(nil).WithAuthToken(token),
	}
}
