package gh

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseRepositoryURL(t *testing.T) {
	t.Run("success without git suffix", func(t *testing.T) {
		parsed, err := url.Parse("https://github.com/owner/repo")
		require.NoError(t, err)

		repo, err := ParseRepositoryURL(parsed)
		require.NoError(t, err)

		assert.Equal(t, "owner", repo.Owner)
		assert.Equal(t, "repo", repo.Name)
		assert.Nil(t, repo.Releases)
	})

	t.Run("success with git suffix", func(t *testing.T) {
		parsed, err := url.Parse("https://github.com/owner/repo.git")
		require.NoError(t, err)

		repo, err := ParseRepositoryURL(parsed)
		require.NoError(t, err)

		assert.Equal(t, "repo", repo.Name)
	})

	t.Run("invalid url path", func(t *testing.T) {
		parsed, err := url.Parse("https://github.com/onlyowner")
		require.NoError(t, err)

		_, err = ParseRepositoryURL(parsed)
		require.Error(t, err)
		assert.ErrorContains(t, err, "invalid GitHub repository URL")
	})

	t.Run("nil url", func(t *testing.T) {
		_, err := ParseRepositoryURL(nil)
		require.Error(t, err)
	})
}

func TestParseRepository(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo, err := ParseRepository("https://github.com/owner/repo.git")
		require.NoError(t, err)
		assert.Equal(t, "owner", repo.Owner)
		assert.Equal(t, "repo", repo.Name)
	})

	t.Run("invalid url string", func(t *testing.T) {
		_, err := ParseRepository("://invalid-url")
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to parse GitHub repository URL")
	})
}
