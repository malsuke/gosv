package domain

import (
	"testing"

	"github.com/malsuke/govs/internal/common/ptr"
	gh "github.com/malsuke/govs/internal/github/domain"
	osvapi "github.com/malsuke/govs/internal/osv/api"
	"github.com/stretchr/testify/assert"
)

func TestExtractRepository_FromPackageName(t *testing.T) {
	v := &osvapi.OsvVulnerability{
		Id: ptr.Ptr("CVE-TEST-1"),
		Affected: &[]osvapi.OsvAffected{
			{
				Package: &osvapi.OsvPackage{
					Ecosystem: ptr.Ptr("GIT"),
					Name:      ptr.Ptr("https://github.com/owner/repo"),
				},
			},
		},
	}

	got, err := ExtractRepository(v)
	assert.NoError(t, err)
	assert.Equal(t, gh.Repository{Owner: "owner", Name: "repo"}, got)
}

func TestExtractRepository_FromRangeRepo(t *testing.T) {
	v := &osvapi.OsvVulnerability{
		Id: ptr.Ptr("CVE-TEST-2"),
		Affected: &[]osvapi.OsvAffected{
			{
				Ranges: &[]osvapi.OsvRange{
					{Repo: ptr.Ptr("https://github.com/owner/repo.git")},
				},
			},
		},
	}

	got, err := ExtractRepository(v)
	assert.NoError(t, err)
	assert.Equal(t, gh.Repository{Owner: "owner", Name: "repo"}, got)
}

func TestExtractRepository_NotFound(t *testing.T) {
	v := &osvapi.OsvVulnerability{}

	_, err := ExtractRepository(v)
	assert.Error(t, err)
}
