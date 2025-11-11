package models

import (
	"github.com/google/go-github/v77/github"
)

type Predicted struct {
	Introduced *CommitMatch
	Fixed      *CommitMatch
}

type CommitMatch struct {
	CommitHash  *string
	Commit      *github.Commit
	PullRequest *github.PullRequest
}

func NewCommitMatch(commitHash *string, commit *github.Commit, pullRequest *github.PullRequest) *CommitMatch {
	return &CommitMatch{
		CommitHash:  commitHash,
		Commit:      commit,
		PullRequest: pullRequest,
	}
}
