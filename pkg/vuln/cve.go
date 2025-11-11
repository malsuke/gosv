package vuln

import (
	"context"

	ghapi "github.com/malsuke/govs/internal/github/api"
	"github.com/malsuke/govs/internal/github/domain"
	vulndomain "github.com/malsuke/govs/internal/vuln/domain"
	vulnservice "github.com/malsuke/govs/internal/vuln/service"
)

func ListVulnerabilitiesByGitHubURL(ctx context.Context, repoURL string) ([]vulndomain.Vulnerability, error) {
	repo, err := domain.ParseRepository(repoURL)
	if err != nil {
		return nil, err
	}
	return vulnservice.FetchVulnerabilitiesByRepository(ctx, repo)
}

func ListCVEIDsByGitHubURL(ctx context.Context, repoURL string) ([]string, error) {
	repo, err := domain.ParseRepository(repoURL)
	if err != nil {
		return nil, err
	}
	return vulnservice.ListCVEIDsByRepository(ctx, repo)
}

func FindSuspectedPullRequestsByCVE(ctx context.Context, client *ghapi.Client, cveID string) (*vulnservice.SuspectedPullRequests, error) {
	return vulnservice.FindSuspectedPullRequestsByCVE(ctx, client, cveID)
}
