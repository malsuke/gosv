package ghapi

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/google/go-github/v77/github"
	gh "github.com/malsuke/govs/internal/github/domain"
)

func (c *Client) GetPullRequestNumberByCommit(ctx context.Context, repo gh.Repository, commitHash string) (int, error) {
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

func (c *Client) SearchMergedPullRequests(ctx context.Context, repo gh.Repository, start time.Time, end time.Time) ([]*github.Issue, error) {
	if ctx == nil {
		return nil, fmt.Errorf("nil context provided")
	}
	if c == nil || c.github == nil {
		return nil, fmt.Errorf("github client is not configured")
	}
	if repo.Owner == "" || repo.Name == "" {
		return nil, fmt.Errorf("repository owner and name must be provided")
	}
	if start.IsZero() || end.IsZero() {
		return nil, fmt.Errorf("time range must be provided")
	}
	if end.Before(start) {
		return nil, fmt.Errorf("end time must not be before start time")
	}

	query := fmt.Sprintf(
		"repo:%s/%s is:pr is:merged merged:%s..%s",
		repo.Owner,
		repo.Name,
		start.Format(time.RFC3339),
		end.Format(time.RFC3339),
	)

	result, _, err := c.github.Search.Issues(ctx, query, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to search merged pull requests: %w", err)
	}
	if result == nil {
		return nil, nil
	}
	return result.Issues, nil
}

func GetPullRequestIDFromCommitHash(client *github.Client, repoURL url.URL, commitHash string) (int, error) {
	repo, err := gh.ParseRepositoryURL(&repoURL)
	if err != nil {
		return 0, err
	}

	return NewClientFromGitHubClient(client).
		GetPullRequestNumberByCommit(context.Background(), repo, commitHash)
}
