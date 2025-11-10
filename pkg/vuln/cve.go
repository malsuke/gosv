package vuln

import (
	osvapi "github.com/malsuke/govs/internal/osv/api"
	osvsvc "github.com/malsuke/govs/internal/osv/service"
)

func ListCVEsStringByGitHubURL(repoURL string) ([]string, error) {
	cveIDs, err := osvsvc.FetchCveIDsByGitHubURL(repoURL)
	if err != nil {
		return nil, err
	}
	return cveIDs, nil
}

func ListCVEsDetailByGitHubURL(repoURL string) ([]osvapi.OsvVulnerability, error) {
	vulns, err := osvsvc.FetchCveVulnerabilitiesByGitHubURL(repoURL)
	if err != nil {
		return nil, err
	}
	return vulns, nil
}
