package service

import (
	"context"
	"fmt"

	"github.com/google/go-github/v77/github"
	ghapi "github.com/malsuke/govs/internal/github/api"
	osvapi "github.com/malsuke/govs/internal/osv/api"
	osvdomain "github.com/malsuke/govs/internal/osv/domain"
	vulndomain "github.com/malsuke/govs/internal/vuln/domain"
)

func NewPredicted(ctx context.Context, client *ghapi.Client, v *osvapi.OsvVulnerability) (*vulndomain.Predicted, error) {
	if v == nil {
		return nil, fmt.Errorf("vulnerability is nil")
	}
	if ctx == nil {
		return nil, fmt.Errorf("nil context provided")
	}

	predicted := &vulndomain.Predicted{}

	if v.Affected == nil || len(*v.Affected) == 0 {
		return predicted, nil
	}

	introduced := osvdomain.ExtractIntroducedCommit(v)
	fixed := osvdomain.ExtractFixedCommit(v)

	if introduced != "" {
		predicted.Introduced = &vulndomain.Introduced{CommitHash: &introduced}
		if client != nil {
			if pr, err := resolvePR(ctx, client, introduced); err == nil {
				predicted.Introduced.PR = pr
			}
		}
		predicted.CommitHash = introduced
	}

	if fixed != "" {
		predicted.Fixed = &vulndomain.Fixed{CommitHash: &fixed}
		if client != nil {
			if pr, err := resolvePR(ctx, client, fixed); err == nil {
				predicted.Fixed.PR = pr
			}
		}
	}

	return predicted, nil
}

func resolvePR(ctx context.Context, client *ghapi.Client, commit string) (*github.PullRequest, error) {
	if client == nil {
		return nil, fmt.Errorf("github client is nil")
	}
	number, err := client.GetPullRequestNumberByCommit(ctx, commit)
	if err != nil {
		return nil, err
	}

	pr, _, err := client.GetGithubClient().PullRequests.Get(ctx, client.Owner, client.Name, number)
	if err != nil {
		return nil, err
	}

	return pr, nil
}
