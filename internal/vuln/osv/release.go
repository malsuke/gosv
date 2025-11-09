package vuln

import (
	"sort"

	"github.com/malsuke/govs/internal/api/osv"
)

/**
 * OSVの脆弱性情報から、影響を受けるバージョン情報を取得する
 */
func CollectReleaseVersions(v *osv.OsvVulnerability) []string {
	if v == nil {
		return nil
	}

	seen := make(map[string]struct{})
	add := func(version string) {
		if version == "" {
			return
		}
		if _, ok := seen[version]; ok {
			return
		}
		seen[version] = struct{}{}
	}

	if v.Affected != nil {
		for _, affected := range *v.Affected {
			if affected.Versions != nil {
				for _, version := range *affected.Versions {
					add(version)
				}
			}
		}
	}

	result := make([]string, 0, len(seen))
	for version := range seen {
		result = append(result, version)
	}
	sort.Strings(result)
	return result
}

// EarliestReleaseVersion returns the lexicographically smallest version collected from the vulnerability.
// It falls back to empty string if no versions are available.
func EarliestReleaseVersion(v *osv.OsvVulnerability) string {
	releases := CollectReleaseVersions(v)
	if len(releases) == 0 {
		return ""
	}
	return releases[0]
}
