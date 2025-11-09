package vuln

import (
	"testing"

	"github.com/malsuke/govs/internal/api/osv"
)

func TestCollectReleaseVersions(t *testing.T) {
	vuln := &osv.OsvVulnerability{
		Affected: &[]osv.OsvAffected{
			{
				Versions: &[]string{"1.0.0", "1.0.1"},
			},
			{
				Versions: &[]string{"2.0.0"},
			},
		},
	}

	got := CollectReleaseVersions(vuln)
	want := []string{"1.0.0", "1.0.1", "2.0.0"}

	if len(got) != len(want) {
		t.Fatalf("CollectReleaseVersions() = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("CollectReleaseVersions()[%d] = %s, want %s", i, got[i], want[i])
		}
	}
}

func TestCollectReleaseVersions_NilVulnerability(t *testing.T) {
	if versions := CollectReleaseVersions(nil); versions != nil {
		t.Fatalf("CollectReleaseVersions(nil) = %v, want nil", versions)
	}
}

func TestEarliestReleaseVersion(t *testing.T) {
	vuln := &osv.OsvVulnerability{
		Affected: &[]osv.OsvAffected{
			{
				Versions: &[]string{"1.2.0", "1.1.0"},
			},
		},
	}

	got := EarliestReleaseVersion(vuln)
	want := "1.1.0"

	if got != want {
		t.Fatalf("EarliestReleaseVersion() = %s, want %s", got, want)
	}
}

func TestEarliestReleaseVersion_NoVersions(t *testing.T) {
	vuln := &osv.OsvVulnerability{}

	if got := EarliestReleaseVersion(vuln); got != "" {
		t.Fatalf("EarliestReleaseVersion() = %s, want empty string", got)
	}
}
