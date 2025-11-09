package vuln

import (
	"testing"

	"github.com/malsuke/govs/internal/api/osv"
	"github.com/malsuke/govs/internal/misc/ptr"
)

func TestExtractIntroducedCommit(t *testing.T) {
	v := &osv.OsvVulnerability{
		Affected: &[]osv.OsvAffected{
			{
				Ranges: &[]osv.OsvRange{
					{
						Events: &[]osv.OsvEvent{
							{Introduced: ptr.String("commit-001")},
							{Introduced: ptr.String("commit-002")},
						},
					},
				},
			},
		},
	}

	if got, want := ExtractIntroducedCommit(v), "commit-001"; got != want {
		t.Fatalf("ExtractIntroducedCommit() = %s, want %s", got, want)
	}
}

func TestExtractIntroducedCommitEmpty(t *testing.T) {
	if got := ExtractIntroducedCommit(&osv.OsvVulnerability{}); got != "" {
		t.Fatalf("ExtractIntroducedCommit() = %s, want empty", got)
	}
}
