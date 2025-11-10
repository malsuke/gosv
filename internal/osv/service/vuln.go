package service

import (
	osvapi "github.com/malsuke/govs/internal/osv/api"
)

// FetchCveIDsByGitHubURL retrieves a list of CVE identifiers associated with the given repository URL.
func FetchCveIDsByGitHubURL(repoURL string) ([]string, error) {
	return osvapi.GetCveIDListFromGithubURL(repoURL)
}

// FetchCveVulnerabilitiesByGitHubURL retrieves detailed OSV vulnerability entries in CVE format for the repository URL.
func FetchCveVulnerabilitiesByGitHubURL(repoURL string) ([]osvapi.OsvVulnerability, error) {
	return osvapi.GetCveVulnerabilityListFromGithubURL(repoURL)
}
