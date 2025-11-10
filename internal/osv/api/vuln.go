package osvapi

import (
	"context"
	"fmt"

	"github.com/malsuke/govs/internal/common/cve"
	"github.com/malsuke/govs/internal/common/ptr"
)

var OsvAPIBaseURL = "https://api.osv.dev"

/**
 * GitHubリポジトリのURLからCVE番号のリストを取得する
 */
func GetCveIDListFromGithubURL(repoUrl string) ([]string, error) {
	vulns, err := fetchAffectedVulnerabilities(repoUrl)
	if err != nil {
		return nil, err
	}

	cveList := make([]string, 0, len(vulns))
	for _, vuln := range vulns {
		if cve := extractCVEFromVulnerability(&vuln); cve != "" {
			cveList = append(cveList, cve)
		}
	}

	return cveList, nil
}

/**
 * GitHubリポジトリのURLからCVE番号形式のOSV脆弱性情報を取得する
 */
func GetCveVulnerabilityListFromGithubURL(repoUrl string) ([]OsvVulnerability, error) {
	vulns, err := fetchAffectedVulnerabilities(repoUrl)
	if err != nil {
		return nil, err
	}

	result := make([]OsvVulnerability, 0, len(vulns))
	for _, vuln := range vulns {
		if extractCVEFromVulnerability(&vuln) != "" {
			result = append(result, vuln)
		}
	}

	return result, nil
}

func GetVulnerabilityByCVE(cveID string) (*OsvVulnerability, error) {
	if !cve.IsValidCVEFormat(cveID) {
		return nil, fmt.Errorf("invalid CVE format: %s", cveID)
	}

	client, err := NewClientWithResponses(OsvAPIBaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	resp, err := client.OSVGetVulnByIdWithResponse(context.Background(), cveID)
	if err != nil {
		return nil, fmt.Errorf("failed to call API: %w", err)
	}

	if resp.StatusCode() != 200 || resp.JSON200 == nil {
		if resp.JSONDefault != nil && resp.JSONDefault.Message != nil {
			return nil, fmt.Errorf("unexpected API response: %s", *resp.JSONDefault.Message)
		}
		return nil, fmt.Errorf("unexpected API response: %v", resp.StatusCode())
	}

	return resp.JSON200, nil
}

func extractCVEFromAliases(aliases *[]string) string {
	if aliases == nil {
		return ""
	}
	for _, alias := range *aliases {
		if cve.IsValidCVEFormat(alias) {
			return alias
		}
	}
	return ""
}

func fetchAffectedVulnerabilities(repoUrl string) ([]OsvVulnerability, error) {
	client, err := NewClientWithResponses(OsvAPIBaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	resp, err := client.OSVQueryAffectedWithResponse(context.Background(),
		OSVQueryAffectedJSONRequestBody{
			Package: &OsvPackage{
				Name:      ptr.Ptr(repoUrl),
				Ecosystem: ptr.Ptr("GIT"),
			},
		})
	if err != nil {
		return nil, fmt.Errorf("failed to call API: %w", err)
	}

	if resp.StatusCode() != 200 || resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected API response: %v", resp.StatusCode())
	}

	vulns := resp.JSON200.Vulns
	if vulns == nil || len(*vulns) == 0 {
		return []OsvVulnerability{}, nil
	}

	copied := make([]OsvVulnerability, len(*vulns))
	copy(copied, *vulns)
	return copied, nil
}

func extractCVEFromVulnerability(v *OsvVulnerability) string {
	if v == nil {
		return ""
	}
	if v.Id != nil && cve.IsValidCVEFormat(*v.Id) {
		return *v.Id
	}
	return extractCVEFromAliases(v.Aliases)
}
