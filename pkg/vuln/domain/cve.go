package domain

import (
	"context"

	"github.com/google/go-github/v77/github"
	ghapi "github.com/malsuke/govs/internal/github/api"
	"github.com/malsuke/govs/internal/github/domain"
	vulndomain "github.com/malsuke/govs/internal/vuln/domain"
	vulnservice "github.com/malsuke/govs/internal/vuln/service"
	"github.com/malsuke/govs/pkg/vuln/models"
)

func ListVulnerabilitiesByGitHubURL(ctx context.Context, repoURL string) ([]vulndomain.Vulnerability, error) {
	owner, name, err := domain.ParseRepository(repoURL)
	if err != nil {
		return nil, err
	}
	return vulnservice.FetchVulnerabilitiesByRepository(ctx, owner, name)
}

func ListCVEIDsByGitHubURL(ctx context.Context, repoURL string) ([]string, error) {
	owner, name, err := domain.ParseRepository(repoURL)
	if err != nil {
		return nil, err
	}
	return vulnservice.ListCVEIDsByRepository(ctx, owner, name)
}

func FindSuspectedPullRequestsByCVE(ctx context.Context, client *ghapi.Client, cveID string) (*vulnservice.SuspectedPullRequests, error) {
	return vulnservice.FindSuspectedPullRequestsByCVE(ctx, client, cveID)
}

func ConvertPredictedDomainToPredicted(ctx context.Context, client *ghapi.Client, predictedDomain *vulndomain.Predicted) (*models.Predicted, error) {
	if predictedDomain == nil {
		return nil, nil
	}

	predicted := &models.Predicted{}

	if predictedDomain.Introduced != nil {
		var commit *github.Commit
		if predictedDomain.Introduced.CommitHash != nil && *predictedDomain.Introduced.CommitHash != "" {
			var err error
			commit, err = client.GetCommit(ctx, *predictedDomain.Introduced.CommitHash)
			if err != nil {
				commit = nil
			}
		}
		predicted.Introduced = models.NewCommitMatch(
			predictedDomain.Introduced.CommitHash,
			commit,
			predictedDomain.Introduced.PR,
		)
	}

	if predictedDomain.Fixed != nil {
		var commit *github.Commit
		if predictedDomain.Fixed.CommitHash != nil && *predictedDomain.Fixed.CommitHash != "" {
			var err error
			commit, err = client.GetCommit(ctx, *predictedDomain.Fixed.CommitHash)
			if err != nil {
				commit = nil
			}
		}
		predicted.Fixed = models.NewCommitMatch(
			predictedDomain.Fixed.CommitHash,
			commit,
			predictedDomain.Fixed.PR,
		)
	}

	return predicted, nil
}

func ConvertSuspectedPullRequestsToSuspected(ctx context.Context, client *ghapi.Client, suspectedPRs *vulnservice.SuspectedPullRequests) (*models.Suspected, error) {
	if suspectedPRs == nil || len(suspectedPRs.MergedPullRequests) == 0 {
		return nil, nil
	}

	prs, err := client.GetPullRequestsFromIssues(ctx, suspectedPRs.MergedPullRequests)
	if err != nil {
		return nil, err
	}

	if len(prs) == 0 {
		return nil, nil
	}

	commits, err := client.GetCommitsFromPullRequests(ctx, prs)
	if err != nil {
		return nil, err
	}

	var prPtr *[]github.PullRequest
	var commitPtr *[]github.Commit
	if len(prs) > 0 {
		prSlice := make([]github.PullRequest, len(prs))
		for i, pr := range prs {
			prSlice[i] = *pr
		}
		prPtr = &prSlice
	}
	if len(commits) > 0 {
		commitSlice := make([]github.Commit, len(commits))
		for i, commit := range commits {
			commitSlice[i] = *commit
		}
		commitPtr = &commitSlice
	}

	return models.NewSuspected(commitPtr, prPtr), nil
}
