package osv

import (
	"context"
	"fmt"

	"github.com/malsuke/govs/internal/misc/cve"
)

func strptr(s string) *string {
	return &s
}

const OSV_API_BASEURL = "https://api.osv.dev"

/**
 * GitHubリポジトリのURLからCVE番号のリストを取得する
 */
func GetCveIDListFromGithubURL(repoUrl string) ([]string, error) {
	client, err := NewClientWithResponses(OSV_API_BASEURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	resp, err := client.OSVQueryAffectedWithResponse(context.Background(),
		OSVQueryAffectedJSONRequestBody{
			Package: &OsvPackage{
				Name:      &repoUrl,
				Ecosystem: strptr("GIT"),
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
		return []string{}, nil
	}

	cveList := make([]string, 0, len(*vulns))
	for _, vuln := range *vulns {
		if vuln.Id != nil && cve.IsValidCVEFormat(*vuln.Id) {
			cveList = append(cveList, *vuln.Id)
			continue
		}

		cve := extractCVEFromAliases(vuln.Aliases)
		if cve != "" {
			cveList = append(cveList, cve)
		}
	}

	return cveList, nil
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
