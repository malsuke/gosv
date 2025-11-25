package ghapi

import (
	"context"
	"fmt"

	"github.com/google/go-github/v77/github"
)

func (c *Client) GetCommit(ctx context.Context, commitHash string) (*github.Commit, error) {
	if err := c.ensureRepositoryContext(); err != nil {
		return nil, err
	}
	if commitHash == "" {
		return nil, fmt.Errorf("commit hash must be provided")
	}

	commit, _, err := c.github.Git.GetCommit(ctx, c.Owner, c.Name, commitHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get commit: %w", err)
	}

	return commit, nil
}

