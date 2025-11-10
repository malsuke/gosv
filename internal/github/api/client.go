package ghapi

import (
	"net/http"

	"github.com/google/go-github/v77/github"
)

type Client struct {
	github *github.Client
}

func NewClient(token string, httpClient *http.Client) *Client {
	ghClient := github.NewClient(httpClient)
	if token != "" {
		ghClient = ghClient.WithAuthToken(token)
	}

	return &Client{github: ghClient}
}

func NewClientFromGitHubClient(client *github.Client) *Client {
	if client == nil {
		client = github.NewClient(nil)
	}
	return &Client{github: client}
}

func (c *Client) GetGithubClient() *github.Client {
	return c.github
}
