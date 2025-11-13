package domain

import (
	"fmt"

	"github.com/malsuke/govs/internal/github/domain"
	osvapi "github.com/malsuke/govs/internal/osv/api"
)

// ExtractRepository pulls the first GitHub repository found in the vulnerability metadata.
func ExtractRepository(v *osvapi.OsvVulnerability) (string, string, error) {
	if v == nil {
		return "", "", fmt.Errorf("vulnerability is nil")
	}

	if v.Affected != nil {
		for _, affected := range *v.Affected {
			if owner, name, ok := repositoryFromAffected(&affected); ok {
				return owner, name, nil
			}
		}
	}

	return "", "", fmt.Errorf("repository information not found in vulnerability %q", safeVulnerabilityID(v))
}

func repositoryFromAffected(affected *osvapi.OsvAffected) (string, string, bool) {
	if affected == nil {
		return "", "", false
	}

	if owner, name, ok := repositoryFromPackage(affected.Package); ok {
		return owner, name, true
	}

	if affected.Ranges != nil {
		for _, r := range *affected.Ranges {
			if owner, name, ok := repositoryFromRange(&r); ok {
				return owner, name, true
			}
		}
	}

	return "", "", false
}

func repositoryFromPackage(pkg *osvapi.OsvPackage) (string, string, bool) {
	if pkg == nil || pkg.Name == nil {
		return "", "", false
	}

	if pkg.Ecosystem != nil && *pkg.Ecosystem != "GIT" {
		return "", "", false
	}

	owner, name, err := domain.ParseRepository(*pkg.Name)
	if err != nil {
		return "", "", false
	}

	return owner, name, true
}

func repositoryFromRange(r *osvapi.OsvRange) (string, string, bool) {
	if r == nil || r.Repo == nil {
		return "", "", false
	}

	owner, name, err := domain.ParseRepository(*r.Repo)
	if err != nil {
		return "", "", false
	}

	return owner, name, true
}

func safeVulnerabilityID(v *osvapi.OsvVulnerability) string {
	if v == nil || v.Id == nil {
		return "<unknown>"
	}
	return *v.Id
}
