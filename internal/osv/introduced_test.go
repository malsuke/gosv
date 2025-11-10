package osv

import (
	"testing"

	"github.com/malsuke/govs/internal/common/ptr"
	osvapi "github.com/malsuke/govs/internal/osv/api"
)

func TestExtractIntroducedCommit(t *testing.T) {
	v := &osvapi.OsvVulnerability{
		Affected: &[]osvapi.OsvAffected{
			{
				Ranges: &[]osvapi.OsvRange{
					{
						Events: &[]osvapi.OsvEvent{
							{Introduced: ptr.Ptr("commit-001")},
							{Introduced: ptr.Ptr("commit-002")},
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
	if got := ExtractIntroducedCommit(&osvapi.OsvVulnerability{}); got != "" {
		t.Fatalf("ExtractIntroducedCommit() = %s, want empty", got)
	}
}
