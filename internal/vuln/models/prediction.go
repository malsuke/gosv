package vuln

import (
	"context"
	"fmt"

	"github.com/google/go-github/v77/github"

	gh "github.com/malsuke/govs/internal/github"
	ghapi "github.com/malsuke/govs/internal/github/api"
	vulnosv "github.com/malsuke/govs/internal/osv"
	osvapi "github.com/malsuke/govs/internal/osv/api"
)

type Predicted struct {
	Introduced *Introduced
	Fixed      *Fixed
	CommitHash string
}

type Introduced struct {
	CommitHash *string
	PR         *github.PullRequest
}

type Fixed struct {
	CommitHash *string
	PR         *github.PullRequest
}

func NewPredicted(ctx context.Context, client *ghapi.Client, repo gh.Repository, v *osvapi.OsvVulnerability) (*Predicted, error) {
	if v == nil {
		return nil, fmt.Errorf("vulnerability is nil")
	}
	if ctx == nil {
		return nil, fmt.Errorf("nil context provided")
	}

	predicted := &Predicted{}

	if v.Affected == nil || len(*v.Affected) == 0 {
		return predicted, nil
	}

	introduced := vulnosv.ExtractIntroducedCommit(v)
	fixed := vulnosv.ExtractFixedCommit(v)

	if introduced != "" {
		predicted.Introduced = &Introduced{CommitHash: &introduced}
		if client != nil {
			if pr, err := resolvePR(ctx, client, repo, introduced); err == nil {
				predicted.Introduced.PR = pr
			}
		}
		predicted.CommitHash = introduced
	}

	if fixed != "" {
		predicted.Fixed = &Fixed{CommitHash: &fixed}
		if client != nil {
			if pr, err := resolvePR(ctx, client, repo, fixed); err == nil {
				predicted.Fixed.PR = pr
			}
		}
	}

	return predicted, nil
}

func resolvePR(ctx context.Context, client *ghapi.Client, repo gh.Repository, commit string) (*github.PullRequest, error) {
	if client == nil {
		return nil, fmt.Errorf("github client is nil")
	}
	number, err := client.GetPullRequestNumberByCommit(ctx, repo, commit)
	if err != nil {
		return nil, err
	}

	pr, _, err := client.GetGithubClient().PullRequests.Get(ctx, repo.Owner, repo.Name, number)
	if err != nil {
		return nil, err
	}

	return pr, nil
}
