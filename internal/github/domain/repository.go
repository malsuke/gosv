package domain

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/google/go-github/v77/github"
)

type Repository struct {
	Owner    string
	Name     string
	Releases []*github.RepositoryRelease
}

func ParseRepositoryURL(u *url.URL) (Repository, error) {
	if u == nil {
		return Repository{}, fmt.Errorf("invalid GitHub repository URL: <nil>")
	}

	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) < 2 {
		return Repository{}, fmt.Errorf("invalid GitHub repository URL: %s", u.String())
	}

	repo := strings.TrimSuffix(parts[1], ".git")

	return Repository{
		Owner: parts[0],
		Name:  repo,
	}, nil
}

func ParseRepository(rawURL string) (Repository, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return Repository{}, fmt.Errorf("failed to parse GitHub repository URL: %w", err)
	}
	return ParseRepositoryURL(parsed)
}
