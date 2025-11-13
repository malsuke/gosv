package ghapi

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/google/go-github/v77/github"
	"github.com/malsuke/govs/internal/github/domain"
)

func (c *Client) GetPullRequestNumberByCommit(ctx context.Context, commitHash string) (int, error) {
	if ctx == nil {
		return 0, fmt.Errorf("nil context provided")
	}
	if c == nil || c.github == nil {
		return 0, fmt.Errorf("github client is not configured")
	}
	if err := c.ensureRepositoryContext(); err != nil {
		return 0, err
	}
	if commitHash == "" {
		return 0, fmt.Errorf("commit hash must be provided")
	}

	prs, _, err := c.github.PullRequests.ListPullRequestsWithCommit(ctx, c.Owner, c.Name, commitHash, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to list pull requests with commit: %w", err)
	}

	if len(prs) == 0 {
		return 0, fmt.Errorf("no pull request found for commit hash %s", commitHash)
	}

	return prs[0].GetNumber(), nil
}

func (c *Client) SearchMergedPullRequests(ctx context.Context, start time.Time, end time.Time) ([]*github.Issue, error) {
	if ctx == nil {
		return nil, fmt.Errorf("nil context provided")
	}
	if c == nil || c.github == nil {
		return nil, fmt.Errorf("github client is not configured")
	}
	if err := c.ensureRepositoryContext(); err != nil {
		return nil, err
	}
	if start.IsZero() || end.IsZero() {
		return nil, fmt.Errorf("time range must be provided")
	}
	if end.Before(start) {
		return nil, fmt.Errorf("end time must not be before start time")
	}

	query := fmt.Sprintf(
		"repo:%s/%s is:pr is:merged merged:%s..%s",
		c.Owner,
		c.Name,
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

func (c *Client) GetPullRequest(ctx context.Context, prNumber int) (*github.PullRequest, error) {
	if ctx == nil {
		return nil, fmt.Errorf("nil context provided")
	}
	if c == nil || c.github == nil {
		return nil, fmt.Errorf("github client is not configured")
	}
	if err := c.ensureRepositoryContext(); err != nil {
		return nil, err
	}

	pr, _, err := c.github.PullRequests.Get(ctx, c.Owner, c.Name, prNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get pull request: %w", err)
	}

	return pr, nil
}

func (c *Client) GetPullRequestsFromIssues(ctx context.Context, issues []*github.Issue) ([]*github.PullRequest, error) {
	if ctx == nil {
		return nil, fmt.Errorf("nil context provided")
	}
	if c == nil || c.github == nil {
		return nil, fmt.Errorf("github client is not configured")
	}
	if err := c.ensureRepositoryContext(); err != nil {
		return nil, err
	}

	prs := make([]*github.PullRequest, 0, len(issues))
	for _, issue := range issues {
		if issue == nil || issue.PullRequestLinks == nil {
			continue
		}

		prNumber := issue.GetNumber()
		pr, err := c.GetPullRequest(ctx, prNumber)
		if err != nil {
			continue
		}
		prs = append(prs, pr)
	}

	return prs, nil
}

func (c *Client) GetCommitsFromPullRequests(ctx context.Context, prs []*github.PullRequest) ([]*github.Commit, error) {
	if ctx == nil {
		return nil, fmt.Errorf("nil context provided")
	}
	if c == nil || c.github == nil {
		return nil, fmt.Errorf("github client is not configured")
	}
	if err := c.ensureRepositoryContext(); err != nil {
		return nil, err
	}

	commits := make([]*github.Commit, 0)
	for _, pr := range prs {
		if pr == nil {
			continue
		}

		mergeSHA := pr.GetMergeCommitSHA()
		if mergeSHA == "" {
			continue
		}

		commit, err := c.GetCommit(ctx, mergeSHA)
		if err != nil {
			continue
		}
		commits = append(commits, commit)
	}

	return commits, nil
}

func GetPullRequestIDFromCommitHash(client *github.Client, repoURL url.URL, commitHash string) (int, error) {
	owner, name, err := domain.ParseRepositoryURL(&repoURL)
	if err != nil {
		return 0, err
	}

	apiClient, err := NewClientFromGitHubClient(owner, name, client)
	if err != nil {
		return 0, err
	}

	return apiClient.GetPullRequestNumberByCommit(context.Background(), commitHash)
}
