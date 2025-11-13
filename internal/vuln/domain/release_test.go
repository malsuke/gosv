package domain

import (
	"testing"

	"github.com/google/go-github/v77/github"
	osvapi "github.com/malsuke/govs/internal/osv/api"
)

func TestFindPreviousRelease(t *testing.T) {
	versions := []string{"1.1.0"}
	v := &osvapi.OsvVulnerability{
		Affected: &[]osvapi.OsvAffected{
			{
				Versions: &versions,
			},
		},
	}

	releases := []*github.RepositoryRelease{
		{TagName: github.Ptr("v1.1.0")},
		{TagName: github.Ptr("1.0.0")},
	}

	prev, err := FindPreviousRelease(v, releases)
	if err != nil {
		t.Fatalf("FindPreviousRelease() returned unexpected error: %v", err)
	}
	if prev == nil {
		t.Fatalf("FindPreviousRelease() returned nil release")
	}
	if got, want := prev.GetTagName(), "1.0.0"; got != want {
		t.Fatalf("FindPreviousRelease() = %s, want %s", got, want)
	}
}

func TestFindPreviousRelease_NoMatchingRelease(t *testing.T) {
	versions := []string{"1.2.0"}
	v := &osvapi.OsvVulnerability{
		Affected: &[]osvapi.OsvAffected{
			{
				Versions: &versions,
			},
		},
	}

	releases := []*github.RepositoryRelease{
		{TagName: github.Ptr("v1.1.0")},
	}

	if prev, err := FindPreviousRelease(v, releases); err == nil || prev != nil {
		t.Fatalf("FindPreviousRelease() = (%v, %v), want (nil, error)", prev, err)
	}
}
