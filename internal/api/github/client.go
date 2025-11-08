package gh

import (
	"github.com/google/go-github/v77/github"
)

func NewClientWrapper(token string) *github.Client {
	if token == "" {
		return github.NewClient(nil)
	}
	return github.NewClient(nil).WithAuthToken(token)
}
