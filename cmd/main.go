package main

import (
	"context"
	"fmt"
	"log"
	"time"

	ghapi "github.com/malsuke/govs/internal/github/api"
	"github.com/malsuke/govs/internal/github/domain"
	"github.com/malsuke/govs/pkg/vuln"
)

func main() {
	ctx := context.Background()
	repoURL := "https://github.com/kubernetes/kubernetes"
	cveIDs, err := vuln.ListCVEIDsByGitHubURL(ctx, repoURL)
	if err != nil {
		log.Fatalf("failed to list CVEs: %v", err)
	}
	fmt.Println(cveIDs)

	layout := "Jan 2, 2006"
	start := "Dec 10, 2016"
	end := "Dec 13, 2016"

	startDate, err := time.Parse(layout, start)
	if err != nil {
		panic(err)
	}

	endDate, err := time.Parse(layout, end)
	if err != nil {
		panic(err)
	}

	repo := domain.Repository{
		Owner: "kubernetes",
		Name:  "kubernetes",
	}

	items, err := ghapi.NewClient("", nil).SearchMergedPullRequests(ctx, repo, startDate, endDate)
	if err != nil {
		panic(err)
	}

	fmt.Println(items)

	for _, item := range items {
		fmt.Println(item.GetNumber())
	}
}

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"net/url"
// 	"os"

// 	"github.com/google/go-github/v77/github"
// 	ghapi "github.com/malsuke/govs/internal/github/api"
// 	gh "github.com/malsuke/govs/internal/github/domain"
// 	osvapi "github.com/malsuke/govs/internal/osv/api"
// 	osvdomain "github.com/malsuke/govs/internal/osv/domain"
// )

// func main() {
// 	token := os.Getenv("GITHUB_TOKEN")
// 	ctx := context.Background()

// 	repoURL, err := url.Parse("https://github.com/kubernetes/kubernetes")
// 	if err != nil {
// 		log.Fatalf("failed to parse repository URL: %v", err)
// 	}

// 	client := ghapi.NewClient(token, nil)
// 	repo, err := gh.ParseRepositoryURL(repoURL)
// 	if err != nil {
// 		log.Fatalf("failed to parse repository: %v", err)
// 	}

// 	items, err := client.SearchMergedPullRequests(ctx, repo, start, end)
// 	if err != nil {
// 		log.Fatalf("failed to search merged pull requests: %v", err)
// 	}

// 	if len(items) == 0 {
// 		fmt.Println("no merged pull requests found in the specified range")
// 		return
// 	}

// 	for _, issue := range items {
// 		printPullRequestNumber(ctx, client, repo, issue)
// 	}

// 	cveID := "CVE-2023-2727"
// 	vuln, err := osvapi.GetVulnerabilityByCVE(cveID)
// 	if err != nil {
// 		log.Fatalf("failed to fetch vulnerability: %v", err)
// 	}

// 	versions := osvdomain.CollectReleaseVersions(vuln)
// 	fmt.Printf("CVE %s affected versions: %v\n", cveID, versions)
// 	fmt.Printf("earliest affected version: %s\n", osvdomain.EarliestReleaseVersion(vuln))
// }

// func printPullRequestNumber(ctx context.Context, client *ghapi.Client, repo gh.Repository, issue *github.Issue) {
// 	if issue == nil {
// 		return
// 	}

// 	prNumber, err := resolvePullRequestNumber(ctx, client, repo, issue)
// 	if err != nil {
// 		log.Printf("failed to resolve pull request number: %v", err)
// 		return
// 	}

// 	fmt.Println(prNumber)
// }

// func resolvePullRequestNumber(ctx context.Context, client *ghapi.Client, repo gh.Repository, issue *github.Issue) (int, error) {
// 	if issue.PullRequestLinks == nil {
// 		return 0, fmt.Errorf("issue %d is not a pull request", issue.GetNumber())
// 	}
// 	if issue.Number != nil {
// 		return issue.GetNumber(), nil
// 	}

// 	pr, _, err := client.GetGithubClient().PullRequests.Get(ctx, repo.Owner, repo.Name, issue.GetNumber())
// 	if err != nil {
// 		return 0, fmt.Errorf("failed to fetch pull request details: %w", err)
// 	}
// 	return pr.GetNumber(), nil
// }
