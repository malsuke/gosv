package osv

import (
	"testing"

	"github.com/malsuke/govs/internal/common/ptr"
	osvapi "github.com/malsuke/govs/internal/osv/api"
)

func TestExtractFixedCommit(t *testing.T) {
	v := &osvapi.OsvVulnerability{
		Affected: &[]osvapi.OsvAffected{
			{
				Ranges: &[]osvapi.OsvRange{
					{
						Events: &[]osvapi.OsvEvent{
							{Fixed: ptr.Ptr("fix-001")},
							{Fixed: ptr.Ptr("fix-002")},
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
	if got := ExtractFixedCommit(&osvapi.OsvVulnerability{}); got != "" {
		t.Fatalf("ExtractFixedCommit() = %s, want empty", got)
	}
}
