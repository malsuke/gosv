package osvapi

import (
	"context"
	"fmt"

	gh "github.com/malsuke/govs/internal/github/domain"

	"github.com/malsuke/govs/internal/common/cve"
	"github.com/malsuke/govs/internal/common/ptr"
)

var OsvAPIBaseURL = "https://api.osv.dev"

// ListCVEIDsByRepository は GitHub リポジトリに紐づく CVE ID の一覧を取得する。
func ListCVEIDsByRepository(ctx context.Context, repo gh.Repository) ([]string, error) {
	vulns, err := listVulnerabilitiesWithCVE(ctx, repo)
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(vulns))
	for _, vuln := range vulns {
		if cveID := ExtractCVEFromVulnerability(&vuln); cveID != "" {
			ids = append(ids, cveID)
		}
	}

	return ids, nil
}

// ListCVEVulnerabilitiesByRepository は GitHub リポジトリに紐づく CVE 形式の OSV 脆弱性一覧を取得する。
func ListCVEVulnerabilitiesByRepository(ctx context.Context, repo gh.Repository) ([]OsvVulnerability, error) {
	return listVulnerabilitiesWithCVE(ctx, repo)
}

// GetVulnerabilityByCVE は CVE ID から脆弱性を取得する。
func GetVulnerabilityByCVE(ctx context.Context, cveID string) (*OsvVulnerability, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context must not be nil")
	}
	if !cve.IsValidCVEFormat(cveID) {
		return nil, fmt.Errorf("invalid CVE format: %s", cveID)
	}

	client, err := NewClientWithResponses(OsvAPIBaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	resp, err := client.OSVGetVulnByIdWithResponse(ctx, cveID)
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

// ExtractCVEFromVulnerability は脆弱性情報から CVE ID を抽出する。
func ExtractCVEFromVulnerability(v *OsvVulnerability) string {
	if v == nil {
		return ""
	}
	if v.Id != nil && cve.IsValidCVEFormat(*v.Id) {
		return *v.Id
	}
	return extractCVEFromAliases(v.Aliases)
}

func listVulnerabilitiesWithCVE(ctx context.Context, repo gh.Repository) ([]OsvVulnerability, error) {
	vulns, err := fetchAffectedVulnerabilities(ctx, repo)
	if err != nil {
		return nil, err
	}

	filtered := make([]OsvVulnerability, 0, len(vulns))
	for _, vuln := range vulns {
		if ExtractCVEFromVulnerability(&vuln) != "" {
			filtered = append(filtered, vuln)
		}
	}

	return filtered, nil
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

func fetchAffectedVulnerabilities(ctx context.Context, repo gh.Repository) ([]OsvVulnerability, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context must not be nil")
	}

	canonicalURL, err := repo.CanonicalGitURL()
	if err != nil {
		return nil, err
	}

	client, err := NewClientWithResponses(OsvAPIBaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	resp, err := client.OSVQueryAffectedWithResponse(ctx,
		OSVQueryAffectedJSONRequestBody{
			Package: &OsvPackage{
				Name:      ptr.Ptr(canonicalURL),
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
