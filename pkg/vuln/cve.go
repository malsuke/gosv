package vuln

import (
	"context"
	"fmt"

	gh "github.com/malsuke/govs/internal/github/domain"
	vulndomain "github.com/malsuke/govs/internal/vuln/domain"
	vulnservice "github.com/malsuke/govs/internal/vuln/service"
)

func ListVulnerabilitiesByGitHubURL(ctx context.Context, repoURL string) ([]vulndomain.Vulnerability, error) {
	repo, err := parseRepository(repoURL)
	if err != nil {
		return nil, err
	}
	return vulnservice.FetchVulnerabilitiesByRepository(ctx, repo)
}

func ListCVEIDsByGitHubURL(ctx context.Context, repoURL string) ([]string, error) {
	repo, err := parseRepository(repoURL)
	if err != nil {
		return nil, err
	}
	return vulnservice.ListCVEIDsByRepository(ctx, repo)
}

func parseRepository(repoURL string) (gh.Repository, error) {
	if repoURL == "" {
		return gh.Repository{}, fmt.Errorf("repository URL must not be empty")
	}
	repo, err := gh.ParseRepository(repoURL)
	if err != nil {
		return gh.Repository{}, fmt.Errorf("failed to parse GitHub repository URL: %w", err)
	}
	return repo, nil
}
