package gh

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

func TestGetPullRequestIDFromCommitHash(t *testing.T) {
	mustParseURL := func(rawURL string) url.URL {
		u, err := url.Parse(rawURL)
		if err != nil {
			t.Fatalf("test setup: failed to parse URL: %v", err)
		}
		return *u
	}

	tests := []struct {
		name          string
		repoURL       url.URL
		commitHash    string
		handler       http.HandlerFunc
		wantPRID      int
		wantErr       bool
		wantErrString string
	}{
		{
			name:       "success",
			repoURL:    mustParseURL("https://github.com/owner/repo"),
			commitHash: "abcdef123456",
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "GET", r.Method)
				assert.Equal(t, "/repos/owner/repo/commits/abcdef123456/pulls", r.URL.Path)

				prNumber := 123
				repoID := int64(12345)

				prs := []*github.PullRequest{
					{
						Number: &prNumber,
						Base: &github.PullRequestBranch{
							Repo: &github.Repository{
								ID: &repoID,
							},
						},
					},
				}
				w.Header().Set("Content-Type", "application/json")
				require.NoError(t, json.NewEncoder(w).Encode(prs))
			},
			wantPRID: 123,
			wantErr:  false,
		},
		{
			name:       "no pull request found",
			repoURL:    mustParseURL("https://github.com/owner/repo"),
			commitHash: "fedcba654321",
			handler: func(w http.ResponseWriter, r *http.Request) {
				prs := []*github.PullRequest{}
				w.Header().Set("Content-Type", "application/json")
				require.NoError(t, json.NewEncoder(w).Encode(prs))
			},
			wantPRID:      0,
			wantErr:       true,
			wantErrString: "no pull request found for commit hash fedcba654321",
		},
		{
			name:          "invalid repo path",
			repoURL:       mustParseURL("https://github.com/owner"),
			commitHash:    "abcdef123456",
			handler:       nil,
			wantPRID:      0,
			wantErr:       true,
			wantErrString: "invalid GitHub repository URL: https://github.com/owner",
		},
		{
			name:       "github api error",
			repoURL:    mustParseURL("https://github.com/owner/repo"),
			commitHash: "abcdef123456",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			wantPRID:      0,
			wantErr:       true,
			wantErrString: "failed to list pull requests with commit",
		},
		{
			name:       "real world case vuejs/core",
			repoURL:    mustParseURL("https://github.com/vuejs/core"),
			commitHash: "079010a38cfff4c49e0a13e54ebff0c189a4d5dc",
			handler: func(w http.ResponseWriter, r *http.Request) {
				expectedPath := "/repos/vuejs/core/commits/079010a38cfff4c49e0a13e54ebff0c189a4d5dc/pulls"
				if r.URL.Path != expectedPath {
					t.Errorf("Request path = %q, want %q", r.URL.Path, expectedPath)
				}

				prNumber := 13974
				repoID := int64(137078487)

				prs := []*github.PullRequest{
					{
						Number: &prNumber,
						Base: &github.PullRequestBranch{
							Repo: &github.Repository{
								ID: &repoID,
							},
						},
					},
				}
				jsonBytes, err := json.Marshal(prs)
				if err != nil {
					t.Fatalf("failed to marshal test data: %v", err)
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write(jsonBytes)
			},
			wantPRID: 13974,
			wantErr:  false,
		},
		{
			name:       "multiple pull requests - returns first one",
			repoURL:    mustParseURL("https://github.com/owner/repo"),
			commitHash: "multiple123456",
			handler: func(w http.ResponseWriter, r *http.Request) {
				prNumber1 := 100
				prNumber2 := 200
				repoID := int64(12345)

				prs := []*github.PullRequest{
					{
						Number: &prNumber1,
						Base: &github.PullRequestBranch{
							Repo: &github.Repository{
								ID: &repoID,
							},
						},
					},
					{
						Number: &prNumber2,
						Base: &github.PullRequestBranch{
							Repo: &github.Repository{
								ID: &repoID,
							},
						},
					},
				}
				w.Header().Set("Content-Type", "application/json")
				require.NoError(t, json.NewEncoder(w).Encode(prs))
			},
			wantPRID: 100,
			wantErr:  false,
		},
		{
			name:       "repository with .git suffix",
			repoURL:    mustParseURL("https://github.com/owner/repo.git"),
			commitHash: "commit123456",
			handler: func(w http.ResponseWriter, r *http.Request) {
				expectedPath := "/repos/owner/repo/commits/commit123456/pulls"
				assert.Equal(t, expectedPath, r.URL.Path)

				prNumber := 999
				repoID := int64(12345)

				prs := []*github.PullRequest{
					{
						Number: &prNumber,
						Base: &github.PullRequestBranch{
							Repo: &github.Repository{
								ID: &repoID,
							},
						},
					},
				}
				w.Header().Set("Content-Type", "application/json")
				require.NoError(t, json.NewEncoder(w).Encode(prs))
			},
			wantPRID: 999,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := github.NewClient(nil)

			if tt.handler != nil {
				server := httptest.NewServer(tt.handler)
				defer server.Close()

				client = github.NewClient(server.Client())
				baseURL, err := url.Parse(server.URL + "/")
				require.NoError(t, err)
				client.BaseURL = baseURL
			}

			prID, err := GetPullRequestIDFromCommitHash(client, tt.repoURL, tt.commitHash)

			if tt.wantErr {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.wantErrString)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantPRID, prID)
		})
	}
}

func TestClientGetPullRequestNumberByCommit(t *testing.T) {
	ctx := context.Background()
	repo := Repository{Owner: "owner", Name: "repo"}

	t.Run("success", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			prNumber := 42
			resp := []*github.PullRequest{
				{Number: &prNumber},
			}
			require.NoError(t, json.NewEncoder(w).Encode(resp))
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		defer server.Close()

		client := github.NewClient(server.Client())
		baseURL, err := url.Parse(server.URL + "/")
		require.NoError(t, err)
		client.BaseURL = baseURL

		number, err := NewClientFromGitHubClient(client).GetPullRequestNumberByCommit(ctx, repo, "hash")
		require.NoError(t, err)
		assert.Equal(t, 42, number)
	})

	t.Run("no pull request found", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			resp := []*github.PullRequest{}
			require.NoError(t, json.NewEncoder(w).Encode(resp))
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		defer server.Close()

		client := github.NewClient(server.Client())
		baseURL, err := url.Parse(server.URL + "/")
		require.NoError(t, err)
		client.BaseURL = baseURL

		_, err = NewClientFromGitHubClient(client).GetPullRequestNumberByCommit(ctx, repo, "hash")
		require.Error(t, err)
		assert.ErrorContains(t, err, "no pull request found for commit hash")
	})

	t.Run("nil context", func(t *testing.T) {
		client := NewClient("", nil)
		var nilCtx context.Context
		_, err := client.GetPullRequestNumberByCommit(nilCtx, repo, "hash")
		require.Error(t, err)
		assert.ErrorContains(t, err, "nil context provided")
	})

	t.Run("nil client", func(t *testing.T) {
		var client *Client
		_, err := client.GetPullRequestNumberByCommit(ctx, repo, "hash")
		require.Error(t, err)
		assert.ErrorContains(t, err, "github client is not configured")
	})
}
