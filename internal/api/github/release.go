package gh

import (
	"context"
	"fmt"

	"github.com/google/go-github/v77/github"
)

type ReleaseListOptions struct {
	ExcludePreRelease bool
	ListOptions       *github.ListOptions
}

/**
 * PreReleaseを除外してStableなReleaseのリストを取得する
 */
func (s *RepositoryState) GetStableReleaseList() error {
	opts := &ReleaseListOptions{
		ExcludePreRelease: true,
		ListOptions:       &github.ListOptions{},
	}
	return s.GetReleaseList(opts)
}

/**
 * オプションを指定してReleaseのリストを取得する
 */
func (s *RepositoryState) GetReleaseList(opts *ReleaseListOptions) error {
	if opts == nil {
		opts = &ReleaseListOptions{
			ExcludePreRelease: false,
			ListOptions:       &github.ListOptions{},
		}
	}

	if opts.ListOptions == nil {
		opts.ListOptions = &github.ListOptions{}
	}

	var err error
	var allReleases []*github.RepositoryRelease

	allReleases, _, err = s.Client.Repositories.ListReleases(context.Background(), s.Owner, s.Repo, opts.ListOptions)
	if err != nil {
		return fmt.Errorf("failed to list releases: %w", err)
	}

	if opts.ExcludePreRelease {
		stableReleases := make([]*github.RepositoryRelease, 0, len(allReleases))
		for _, release := range allReleases {
			if release.Prerelease == nil || !*release.Prerelease {
				stableReleases = append(stableReleases, release)
			}
		}
		s.Releases = stableReleases
	} else {
		s.Releases = allReleases
	}

	return nil
}
