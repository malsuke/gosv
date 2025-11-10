package osv

import osvapi "github.com/malsuke/govs/internal/osv/api"

// ExtractIntroducedCommit returns the first non-empty introduced commit hash from the vulnerability.
func ExtractIntroducedCommit(v *osvapi.OsvVulnerability) string {
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
				if event.Introduced != nil && *event.Introduced != "" {
					return *event.Introduced
				}
			}
		}
	}

	return ""
}
