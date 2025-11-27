package ghapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-github/v77/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient("", "owner/repo", nil)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.NotNil(t, client.github)
	require.NotNil(t, client.github.BaseURL)
	assert.Equal(t, "owner", client.Owner)
	assert.Equal(t, "repo", client.Name)
}

func TestNewClient_InvalidRepository(t *testing.T) {
	client, err := NewClient("", "invalid", nil)
	require.Error(t, err)
	assert.Nil(t, client)
}

func TestNewClientFromGitHubClient(t *testing.T) {
	t.Run("nil github client creates new instance", func(t *testing.T) {
		client, err := NewClientFromGitHubClient("owner", "repo", nil)
		require.NoError(t, err)
		require.NotNil(t, client)
		require.NotNil(t, client.github)
	})

	t.Run("existing github client reused", func(t *testing.T) {
		existing := github.NewClient(nil)
		client, err := NewClientFromGitHubClient("owner", "repo", existing)
		require.NoError(t, err)
		require.NotNil(t, client)
		assert.Same(t, existing, client.github)
	})

	t.Run("invalid repository", func(t *testing.T) {
		client, err := NewClientFromGitHubClient("", "", nil)
		require.Error(t, err)
		assert.Nil(t, client)
	})
}

func TestClientGetRepository(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/repos/owner/repo", r.URL.Path)

			repoID := int64(123)
			require.NoError(t, json.NewEncoder(w).Encode(github.Repository{ID: &repoID}))
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		defer server.Close()

		client := github.NewClient(server.Client())
		baseURL, err := url.Parse(server.URL + "/")
		require.NoError(t, err)
		client.BaseURL = baseURL

		apiClient, err := NewClientFromGitHubClient("owner", "repo", client)
		require.NoError(t, err)

		repo, err := apiClient.GetRepository(ctx)
		require.NoError(t, err)
		require.NotNil(t, repo)
		assert.Equal(t, int64(123), repo.GetID())
	})

	t.Run("github api error", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		defer server.Close()

		client := github.NewClient(server.Client())
		baseURL, err := url.Parse(server.URL + "/")
		require.NoError(t, err)
		client.BaseURL = baseURL

		apiClient, err := NewClientFromGitHubClient("owner", "repo", client)
		require.NoError(t, err)

		_, err = apiClient.GetRepository(ctx)
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to get repository")
	})

	t.Run("nil context", func(t *testing.T) {
		client, err := NewClient("", "owner/repo", nil)
		require.NoError(t, err)

		var nilCtx context.Context
		_, err = client.GetRepository(nilCtx)
		require.Error(t, err)
		assert.ErrorContains(t, err, "nil context provided")
	})

	t.Run("nil client", func(t *testing.T) {
		var client *Client
		_, err := client.GetRepository(ctx)
		require.Error(t, err)
		assert.ErrorContains(t, err, "github client is not configured")
	})

	t.Run("missing repository context", func(t *testing.T) {
		client := &Client{github: github.NewClient(nil)}
		_, err := client.GetRepository(ctx)
		require.Error(t, err)
		assert.ErrorContains(t, err, "repository owner is empty")
	})
}

func TestClientGetRepositoryWithReleases(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/repos/owner/repo":
				repoID := int64(100)
				require.NoError(t, json.NewEncoder(w).Encode(github.Repository{ID: &repoID}))
			case "/repos/owner/repo/releases":
				tag1 := "v1.0.0"
				tag2 := "v1.1.0-beta"
				prereleaseFalse := false
				prereleaseTrue := true
				require.NoError(t, json.NewEncoder(w).Encode([]*github.RepositoryRelease{
					{TagName: &tag1, Prerelease: &prereleaseFalse},
					{TagName: &tag2, Prerelease: &prereleaseTrue},
				}))
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		defer server.Close()

		client := github.NewClient(server.Client())
		baseURL, err := url.Parse(server.URL + "/")
		require.NoError(t, err)
		client.BaseURL = baseURL

		apiClient, err := NewClientFromGitHubClient("owner", "repo", client)
		require.NoError(t, err)

		summary, err := apiClient.GetRepositoryWithReleases(ctx, nil)
		require.NoError(t, err)
		require.NotNil(t, summary)
		require.NotNil(t, summary.Repository)
		assert.Equal(t, int64(100), summary.Repository.GetID())
		assert.Len(t, summary.ReleasesWithoutPreRelease, 1)
		assert.Equal(t, "v1.0.0", summary.ReleasesWithoutPreRelease[0].GetTagName())
	})

	t.Run("release fetch error", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/repos/owner/repo" {
				require.NoError(t, json.NewEncoder(w).Encode(github.Repository{}))
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		defer server.Close()

		client := github.NewClient(server.Client())
		baseURL, err := url.Parse(server.URL + "/")
		require.NoError(t, err)
		client.BaseURL = baseURL

		apiClient, err := NewClientFromGitHubClient("owner", "repo", client)
		require.NoError(t, err)

		_, err = apiClient.GetRepositoryWithReleases(ctx, nil)
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to list releases")
	})

	t.Run("nil context", func(t *testing.T) {
		client, err := NewClient("", "owner/repo", nil)
		require.NoError(t, err)
		var nilCtx context.Context
		_, err = client.GetRepositoryWithReleases(nilCtx, nil)
		require.Error(t, err)
		assert.ErrorContains(t, err, "nil context provided")
	})

	t.Run("nil client", func(t *testing.T) {
		var client *Client
		_, err := client.GetRepositoryWithReleases(ctx, nil)
		require.Error(t, err)
		assert.ErrorContains(t, err, "github client is not configured")
	})
}
