package domain

import (
	"testing"
	"time"

	"github.com/google/go-github/v77/github"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeVersion(t *testing.T) {
	assert.Equal(t, "1.0.0", NormalizeVersion("v1.0.0"))
	assert.Equal(t, "1.0.0", NormalizeVersion(" 1.0.0 "))
}

func TestFindReleaseByVersion(t *testing.T) {
	releases := []*github.RepositoryRelease{
		{TagName: github.String("v1.2.3")},
		{Name: github.String("1.2.4")},
	}

	assert.NotNil(t, FindReleaseByVersion(releases, "1.2.3"))
	assert.NotNil(t, FindReleaseByVersion(releases, "v1.2.4"))
	assert.Nil(t, FindReleaseByVersion(releases, "2.0.0"))
}

func TestSortReleasesByTime(t *testing.T) {
	t1 := github.Timestamp{Time: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)}
	t2 := github.Timestamp{Time: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC)}
	t3 := github.Timestamp{Time: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)}

	releases := []*github.RepositoryRelease{
		{PublishedAt: &t1},
		{PublishedAt: &t2},
		{PublishedAt: &t3},
	}

	sorted := SortReleasesByTime(releases)

	assert.Equal(t, t3.Time, sorted[0].PublishedAt.Time)
	assert.Equal(t, t1.Time, sorted[1].PublishedAt.Time)
	assert.Equal(t, t2.Time, sorted[2].PublishedAt.Time)
}

func TestReleaseTime(t *testing.T) {
	t1 := github.Timestamp{Time: time.Now()}
	release := &github.RepositoryRelease{PublishedAt: &t1}

	got, ok := ReleaseTime(release)
	assert.True(t, ok)
	assert.Equal(t, t1.Time, got)

	got, ok = ReleaseTime(&github.RepositoryRelease{})
	assert.False(t, ok)
	assert.True(t, got.IsZero())
}
