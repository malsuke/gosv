package models

import "github.com/google/go-github/v77/github"

/**
 * 脆弱性の原因となったコミットとPRを保持する構造体
 */
type Suspected struct {
	Commit      *[]github.Commit
	PullRequest *[]github.PullRequest
}
