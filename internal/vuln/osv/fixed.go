package vuln

import "github.com/malsuke/govs/internal/api/osv"

// ExtractFixedCommit returns the first non-empty fixed commit hash from the vulnerability.
func ExtractFixedCommit(v *osv.OsvVulnerability) string {
	if v == nil || v.Affected == nil {
		return ""
	}

	for _, affected := range *v.Affected {
		if affected.Ranges == nil {
			continue
		}
		for _, r := range *affected.Ranges {
			if r.Events == nil {
				continue
			}
			for _, event := range *r.Events {
				if event.Fixed != nil && *event.Fixed != "" {
					return *event.Fixed
				}
			}
		}
	}

	return ""
}
