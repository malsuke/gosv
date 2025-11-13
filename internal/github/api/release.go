package ghapi

import (
	"context"
	"fmt"

	"github.com/google/go-github/v77/github"
)

type ReleaseListOptions struct {
	ExcludePreRelease bool
	ListOptions       github.ListOptions
}

func (c *Client) ListReleases(ctx context.Context, opts ReleaseListOptions) ([]*github.RepositoryRelease, error) {
	if ctx == nil {
		return nil, fmt.Errorf("nil context provided")
	}
	if c == nil || c.github == nil {
		return nil, fmt.Errorf("github client is not configured")
	}
	if err := c.ensureRepositoryContext(); err != nil {
		return nil, err
	}

	listOpts := opts.ListOptions
	releases, _, err := c.github.Repositories.ListReleases(ctx, c.Owner, c.Name, &listOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to list releases: %w", err)
	}

	if !opts.ExcludePreRelease {
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

	return filtered, nil
}

func (c *Client) ListStableReleases(ctx context.Context, listOpts github.ListOptions) ([]*github.RepositoryRelease, error) {
	return c.ListReleases(ctx, ReleaseListOptions{
		ExcludePreRelease: true,
		ListOptions:       listOpts,
	})
}
