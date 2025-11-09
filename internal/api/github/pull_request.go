package gh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/google/go-github/v77/github"
)

func (c *Client) GetPullRequestNumberByCommit(ctx context.Context, repo Repository, commitHash string) (int, error) {
	if ctx == nil {
		return 0, fmt.Errorf("nil context provided")
	}
	if c == nil || c.github == nil {
		return 0, fmt.Errorf("github client is not configured")
	}

	prs, _, err := c.github.PullRequests.ListPullRequestsWithCommit(ctx, repo.Owner, repo.Name, commitHash, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to list pull requests with commit: %w", err)
	}

	if len(prs) == 0 {
		return 0, fmt.Errorf("no pull request found for commit hash %s", commitHash)
	}

	return prs[0].GetNumber(), nil
}

func GetPullRequestIDFromCommitHash(client *github.Client, repoURL url.URL, commitHash string) (int, error) {
	repo, err := ParseRepositoryURL(&repoURL)
	if err != nil {
		return 0, err
	}

	return NewClientFromGitHubClient(client).
		GetPullRequestNumberByCommit(context.Background(), repo, commitHash)
}
