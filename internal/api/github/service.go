package gh

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/google/go-github/v77/github"
)

func GetPullRequestIDFromCommitHash(client *github.Client, repoURL url.URL, commitHash string) (int, error) {
	pathParts := strings.Split(strings.Trim(repoURL.Path, "/"), "/")
	if len(pathParts) < 2 {
		return 0, fmt.Errorf("invalid repo URL path: %s", repoURL.Path)
	}

	owner := pathParts[0]
	repo := strings.TrimSuffix(pathParts[1], ".git")

	prs, _, err := client.PullRequests.ListPullRequestsWithCommit(context.Background(), owner, repo, commitHash, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to list pull requests with commit: %w", err)
	}

	if len(prs) == 0 {
		return 0, fmt.Errorf("no pull request found for commit hash %s", commitHash)
	}

	return prs[0].GetNumber(), nil
}
