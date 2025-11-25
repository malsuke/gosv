package ghapi

import (
	"net/http"

	"github.com/google/go-github/v77/github"
	"github.com/malsuke/govs/internal/github/domain"
)

type Client struct {
	Owner  string
	Name   string
	github *github.Client
}

func NewClient(token string, repo string, httpClient *http.Client) (*Client, error) {
	owner, name, err := domain.ParseRepository(repo)
	if err != nil {
		return nil, err
	}

	var ghClient *github.Client
	if httpClient != nil {
		ghClient = github.NewClient(httpClient)
	} else {
		ghClient = github.NewClient(nil)
	}

	if token != "" {
		ghClient = ghClient.WithAuthToken(token)
	}

	return &Client{Owner: owner, Name: name, github: ghClient}, nil
}

func NewClientFromGitHubClient(owner, name string, client *github.Client) (*Client, error) {
	if err := domain.ValidateRepository(owner, name); err != nil {
		return nil, err
	}

	if client == nil {
		client = github.NewClient(nil)
	}
	return &Client{Owner: owner, Name: name, github: client}, nil
}

func (c *Client) GetGithubClient() *github.Client {
	return c.github
}

func (c *Client) ensureRepositoryContext() error {
	return domain.ValidateRepository(c.Owner, c.Name)
}
