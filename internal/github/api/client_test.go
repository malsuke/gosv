package ghapi

import (
	"testing"

	"github.com/google/go-github/v77/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	client := NewClient("", nil)
	require.NotNil(t, client)
	require.NotNil(t, client.github)
	require.NotNil(t, client.github.BaseURL)
}

func TestNewClientFromGitHubClient(t *testing.T) {
	t.Run("nil github client creates new instance", func(t *testing.T) {
		client := NewClientFromGitHubClient(nil)
		require.NotNil(t, client)
		require.NotNil(t, client.github)
	})

	t.Run("existing github client reused", func(t *testing.T) {
		existing := github.NewClient(nil)
		client := NewClientFromGitHubClient(existing)
		require.NotNil(t, client)
		assert.Same(t, existing, client.github)
	})
}
