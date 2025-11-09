package gh

import (
	"context"
	"fmt"

	"github.com/google/go-github/v77/github"
)

type ReleaseListOptions struct {
	ExcludePreRelease bool
	ListOptions       github.ListOptions
}

func (c *Client) ListReleases(ctx context.Context, repo *Repository, opts ReleaseListOptions) ([]*github.RepositoryRelease, error) {
	if ctx == nil {
		return nil, fmt.Errorf("nil context provided")
	}
	if c == nil || c.github == nil {
		return nil, fmt.Errorf("github client is not configured")
	}
	if repo == nil {
		return nil, fmt.Errorf("repository is nil")
	}

	listOpts := opts.ListOptions
	releases, _, err := c.github.Repositories.ListReleases(ctx, repo.Owner, repo.Name, &listOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to list releases: %w", err)
	}

	if !opts.ExcludePreRelease {
		repo.Releases = releases
		return releases, nil
	}

	filtered := make([]*github.RepositoryRelease, 0, len(releases))
	for _, release := range releases {
		if release == nil {
			continue
		}
		if release.GetPrerelease() {
			continue
		}
		filtered = append(filtered, release)
	}

	repo.Releases = filtered
	return filtered, nil
}

func (c *Client) ListStableReleases(ctx context.Context, repo *Repository, listOpts github.ListOptions) ([]*github.RepositoryRelease, error) {
	return c.ListReleases(ctx, repo, ReleaseListOptions{
		ExcludePreRelease: true,
		ListOptions:       listOpts,
	})
}
