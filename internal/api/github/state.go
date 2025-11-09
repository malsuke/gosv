package gh

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/google/go-github/v77/github"
)

type RepositoryState struct {
	RepoUrl  *url.URL
	Owner    string
	Repo     string
	Releases []*github.RepositoryRelease
}

func NewRepositoryState(token string, repoUrl *url.URL) (*RepositoryState, error) {
	owner, repo, err := parseGitHubURL(repoUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse GitHub repository URL: %w", err)
	}

	return &RepositoryState{
		RepoUrl: repoUrl,
		Owner:   owner,
		Repo:    repo,
	}, nil
}

func parseGitHubURL(u *url.URL) (owner, repo string, err error) {
	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid GitHub repository URL: %s", u.String())
	}
	owner = parts[0]
	repo = parts[1]
	repo = strings.TrimSuffix(repo, ".git")
	return owner, repo, nil
}
