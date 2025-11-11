package domain

import (
	"fmt"

	gh "github.com/malsuke/govs/internal/github/domain"
	osvapi "github.com/malsuke/govs/internal/osv/api"
)

// ExtractRepository pulls the first GitHub repository found in the vulnerability metadata.
func ExtractRepository(v *osvapi.OsvVulnerability) (gh.Repository, error) {
	if v == nil {
		return gh.Repository{}, fmt.Errorf("vulnerability is nil")
	}

	if v.Affected != nil {
		for _, affected := range *v.Affected {
			if repo, ok := repositoryFromAffected(&affected); ok {
				return repo, nil
			}
		}
	}

	return gh.Repository{}, fmt.Errorf("repository information not found in vulnerability %q", safeVulnerabilityID(v))
}

func repositoryFromAffected(affected *osvapi.OsvAffected) (gh.Repository, bool) {
	if affected == nil {
		return gh.Repository{}, false
	}

	if repo, ok := repositoryFromPackage(affected.Package); ok {
		return repo, true
	}

	if affected.Ranges != nil {
		for _, r := range *affected.Ranges {
			if repo, ok := repositoryFromRange(&r); ok {
				return repo, true
			}
		}
	}

	return gh.Repository{}, false
}

func repositoryFromPackage(pkg *osvapi.OsvPackage) (gh.Repository, bool) {
	if pkg == nil || pkg.Name == nil {
		return gh.Repository{}, false
	}

	if pkg.Ecosystem != nil && *pkg.Ecosystem != "GIT" {
		return gh.Repository{}, false
	}

	repo, err := gh.ParseRepository(*pkg.Name)
	if err != nil {
		return gh.Repository{}, false
	}

	return repo, true
}

func repositoryFromRange(r *osvapi.OsvRange) (gh.Repository, bool) {
	if r == nil || r.Repo == nil {
		return gh.Repository{}, false
	}

	repo, err := gh.ParseRepository(*r.Repo)
	if err != nil {
		return gh.Repository{}, false
	}

	return repo, true
}

func safeVulnerabilityID(v *osvapi.OsvVulnerability) string {
	if v == nil || v.Id == nil {
		return "<unknown>"
	}
	return *v.Id
}
