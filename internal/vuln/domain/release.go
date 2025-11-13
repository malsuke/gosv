package domain

import (
	"fmt"
	"strings"

	"github.com/google/go-github/v77/github"
	osvapi "github.com/malsuke/govs/internal/osv/api"
	osvdomain "github.com/malsuke/govs/internal/osv/domain"
)

func FindPreviousRelease(vulnerability *osvapi.OsvVulnerability, releases []*github.RepositoryRelease) (*github.RepositoryRelease, error) {
	if vulnerability == nil {
		return nil, fmt.Errorf("vulnerability is nil")
	}
	if len(releases) == 0 {
		return nil, fmt.Errorf("releases are empty")
	}

	affectedVersions := osvdomain.CollectReleaseVersions(vulnerability)
	if len(affectedVersions) == 0 {
		return nil, fmt.Errorf("vulnerability has no affected versions")
	}

	for _, version := range affectedVersions {
		if prev := previousRelease(releases, version); prev != nil {
			return prev, nil
		}
	}

	return nil, fmt.Errorf("no matching release found for affected versions: %v", affectedVersions)
}

func previousRelease(releases []*github.RepositoryRelease, targetVersion string) *github.RepositoryRelease {
	if len(releases) == 0 {
		return nil
	}

	for idx, release := range releases {
		if release == nil {
			continue
		}
		if matchVersion(targetVersion, release.GetTagName()) {
			if idx+1 < len(releases) {
				return releases[idx+1]
			}
			return nil
		}
	}

	return nil
}

func matchVersion(version, tag string) bool {
	if version == "" || tag == "" {
		return false
	}
	return normalizeVersion(version) == normalizeVersion(tag)
}

func normalizeVersion(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "refs/tags/")
	if len(s) > 0 && (s[0] == 'v' || s[0] == 'V') {
		s = s[1:]
	}
	return s
}
