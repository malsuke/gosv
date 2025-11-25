package ghapi

import (
	"context"
	"fmt"

	"github.com/google/go-github/v77/github"
	"github.com/malsuke/govs/pkg/vuln/models"
)

func (c *Client) GetRepository(ctx context.Context) (*github.Repository, error) {
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
	if err := c.ensureRepositoryContext(); err != nil {
		return nil, err
	}

	repo, err := c.GetRepository(ctx)
	if err != nil {
		return nil, err
	}

	releases, err := c.ListAllStableReleases(ctx)
	if err != nil {
		return nil, err
	}

	return models.NewRepository(repo, releases), nil
}
