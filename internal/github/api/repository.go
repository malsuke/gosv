package ghapi

import (
	"context"
	"fmt"

	"github.com/google/go-github/v77/github"
	"github.com/malsuke/govs/pkg/vuln/models"
)

func (c *Client) GetRepository(ctx context.Context) (*github.Repository, error) {
	if ctx == nil {
		return nil, fmt.Errorf("nil context provided")
	}
	if c == nil || c.github == nil {
		return nil, fmt.Errorf("github client is not configured")
	}
	if err := c.ensureRepositoryContext(); err != nil {
		return nil, err
	}

	repo, _, err := c.github.Repositories.Get(ctx, c.Owner, c.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository: %w", err)
	}

	return repo, nil
}

func (c *Client) GetRepositoryWithReleases(ctx context.Context, opts *github.ListOptions) (*models.Repository, error) {
	if ctx == nil {
		return nil, fmt.Errorf("nil context provided")
	}
	if c == nil || c.github == nil {
		return nil, fmt.Errorf("github client is not configured")
	}
	if err := c.ensureRepositoryContext(); err != nil {
		return nil, err
	}

	repo, err := c.GetRepository(ctx)
	if err != nil {
		return nil, err
	}

	listOpts := github.ListOptions{}
	if opts != nil {
		listOpts = *opts
	}

	releases, err := c.ListStableReleases(ctx, listOpts)
	if err != nil {
		return nil, err
	}

	return models.NewRepository(repo, releases), nil
}
