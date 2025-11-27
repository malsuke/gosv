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

// ListAllReleases fetches all releases across all pages.
func (c *Client) ListAllReleases(ctx context.Context, opts ReleaseListOptions) ([]*github.RepositoryRelease, error) {
	if err := c.ensureRepositoryContext(); err != nil {
		return nil, err
	}

	allReleases := make([]*github.RepositoryRelease, 0)
	listOpts := opts.ListOptions
	if listOpts.PerPage == 0 {
		listOpts.PerPage = 100
	}

	for {
		releases, resp, err := c.github.Repositories.ListReleases(ctx, c.Owner, c.Name, &listOpts)
		if err != nil {
			return nil, fmt.Errorf("failed to list releases: %w", err)
		}

		if !opts.ExcludePreRelease {
			allReleases = append(allReleases, releases...)
		} else {
			for _, release := range releases {
				if release != nil && !release.GetPrerelease() {
					allReleases = append(allReleases, release)
				}
			}
		}

		if resp.NextPage == 0 {
			break
		}
		listOpts.Page = resp.NextPage
	}

	return allReleases, nil
}

// ListAllStableReleases fetches all stable releases across all pages.
func (c *Client) ListAllStableReleases(ctx context.Context) ([]*github.RepositoryRelease, error) {
	return c.ListAllReleases(ctx, ReleaseListOptions{
		ExcludePreRelease: true,
		ListOptions:       github.ListOptions{PerPage: 100},
	})
}
