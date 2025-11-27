package domain

import (
	"sort"
	"strings"
	"time"

	"github.com/google/go-github/v77/github"
)

// NormalizeVersion removes common prefixes from a version string for comparison.
func NormalizeVersion(version string) string {
	return strings.TrimPrefix(strings.TrimSpace(version), "v")
}

// FindReleaseByVersion searches for a release whose tag or name matches the provided version.
func FindReleaseByVersion(releases []*github.RepositoryRelease, version string) *github.RepositoryRelease {
	if len(releases) == 0 || version == "" {
		return nil
	}

	target := NormalizeVersion(version)
	for _, release := range releases {
		if release == nil {
			continue
		}
		if NormalizeVersion(release.GetTagName()) == target || NormalizeVersion(release.GetName()) == target {
			return release
		}
	}
	return nil
}

// SortReleasesByTime returns a copy of releases sorted by their publish/creation time in ascending order.
func SortReleasesByTime(releases []*github.RepositoryRelease) []*github.RepositoryRelease {
	sorted := make([]*github.RepositoryRelease, 0, len(releases))
	sorted = append(sorted, releases...)

	sort.SliceStable(sorted, func(i, j int) bool {
		ti, _ := ReleaseTime(sorted[i])
		tj, _ := ReleaseTime(sorted[j])
		return ti.Before(tj)
	})

	return sorted
}

// ReleaseTime extracts the best effort timestamp of a release.
func ReleaseTime(release *github.RepositoryRelease) (time.Time, bool) {
	if release == nil {
		return time.Time{}, false
	}

	if release.PublishedAt != nil && !release.PublishedAt.Time.IsZero() {
		return release.PublishedAt.Time, true
	}

	if release.CreatedAt != nil && !release.CreatedAt.Time.IsZero() {
		return release.CreatedAt.Time, true
	}

	return time.Time{}, false
}
