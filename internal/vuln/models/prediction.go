package vuln

import "github.com/google/go-github/v77/github"

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
