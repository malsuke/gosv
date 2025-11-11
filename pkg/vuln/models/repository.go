package models

import (
	"github.com/google/go-github/v77/github"
)

type Repository struct {
	Repository                *github.Repository
	ReleasesWithoutPreRelease []*github.RepositoryRelease
}

func NewRepository(
	repository *github.Repository,
	releasesWithoutPreRelease []*github.RepositoryRelease,
) *Repository {
	return &Repository{
		Repository:                repository,
		ReleasesWithoutPreRelease: releasesWithoutPreRelease,
	}
}
