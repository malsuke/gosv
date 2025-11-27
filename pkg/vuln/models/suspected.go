package models

import "github.com/google/go-github/v77/github"

/**
 * 脆弱性の原因となった可能性のあるコミットとPRをまとめて保持する構造体
 */
type Suspected struct {
	Commit      *[]github.Commit      `json:"commit"`
	PullRequest *[]github.PullRequest `json:"pull_request"`
}

func NewSuspected(commit *[]github.Commit, pullRequest *[]github.PullRequest) *Suspected {
	return &Suspected{
		Commit:      commit,
		PullRequest: pullRequest,
	}
}
