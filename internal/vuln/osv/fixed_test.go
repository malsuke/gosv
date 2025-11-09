package vuln

import (
	"testing"

	"github.com/malsuke/govs/internal/api/osv"
	"github.com/malsuke/govs/internal/misc/ptr"
)

func TestExtractFixedCommit(t *testing.T) {
	v := &osv.OsvVulnerability{
		Affected: &[]osv.OsvAffected{
			{
				Ranges: &[]osv.OsvRange{
					{
						Events: &[]osv.OsvEvent{
							{Fixed: ptr.String("fix-001")},
							{Fixed: ptr.String("fix-002")},
						},
					},
				},
			},
		},
	}

	if got, want := ExtractFixedCommit(v), "fix-001"; got != want {
		t.Fatalf("ExtractFixedCommit() = %s, want %s", got, want)
	}
}

func TestExtractFixedCommitEmpty(t *testing.T) {
	if got := ExtractFixedCommit(&osv.OsvVulnerability{}); got != "" {
		t.Fatalf("ExtractFixedCommit() = %s, want empty", got)
	}
}
