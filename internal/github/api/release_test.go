package ghapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-github/v77/github"
	gh "github.com/malsuke/govs/internal/github/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetReleaseList(t *testing.T) {
	tests := []struct {
		name             string
		opts             *ReleaseListOptions
		handler          http.HandlerFunc
		wantReleaseCount int
		wantPreRelease   bool
		wantErr          bool
		wantErrString    string
	}{
		{
			name: "ExcludePreRelease true - stable releases only",
			opts: &ReleaseListOptions{
				ExcludePreRelease: true,
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "GET", r.Method)
				assert.Equal(t, "/repos/owner/repo/releases", r.URL.Path)

				tag1 := "v1.0.0"
				tag2 := "v1.1.0-beta"
				tag3 := "v1.2.0"
				prerelease1 := false
				prerelease2 := true
				prerelease3 := false

				releases := []*github.RepositoryRelease{
					{
						TagName:    &tag1,
						Prerelease: &prerelease1,
					},
					{
						TagName:    &tag2,
						Prerelease: &prerelease2,
					},
					{
						TagName:    &tag3,
						Prerelease: &prerelease3,
					},
				}
				w.Header().Set("Content-Type", "application/json")
				require.NoError(t, json.NewEncoder(w).Encode(releases))
			},
			wantReleaseCount: 2,
			wantPreRelease:   false,
			wantErr:          false,
		},
		{
			name: "ExcludePreRelease false - all releases",
			opts: &ReleaseListOptions{
				ExcludePreRelease: false,
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				tag1 := "v1.0.0"
				tag2 := "v1.1.0-beta"
				tag3 := "v1.2.0"
				prerelease1 := false
				prerelease2 := true
				prerelease3 := false

				releases := []*github.RepositoryRelease{
					{
						TagName:    &tag1,
						Prerelease: &prerelease1,
					},
					{
						TagName:    &tag2,
						Prerelease: &prerelease2,
					},
					{
						TagName:    &tag3,
						Prerelease: &prerelease3,
					},
				}
				w.Header().Set("Content-Type", "application/json")
				require.NoError(t, json.NewEncoder(w).Encode(releases))
			},
			wantReleaseCount: 3,
			wantPreRelease:   true,
			wantErr:          false,
		},
		{
			name: "opts is nil - default behavior",
			opts: nil,
			handler: func(w http.ResponseWriter, r *http.Request) {
				tag1 := "v1.0.0"
				tag2 := "v1.1.0-beta"
				prerelease1 := false
				prerelease2 := true

				releases := []*github.RepositoryRelease{
					{
						TagName:    &tag1,
						Prerelease: &prerelease1,
					},
					{
						TagName:    &tag2,
						Prerelease: &prerelease2,
					},
				}
				w.Header().Set("Content-Type", "application/json")
				require.NoError(t, json.NewEncoder(w).Encode(releases))
			},
			wantReleaseCount: 2,
			wantPreRelease:   true,
			wantErr:          false,
		},
		{
			name: "ListOptions is nil - should use default",
			opts: &ReleaseListOptions{
				ExcludePreRelease: true,
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				tag := "v1.0.0"
				prerelease := false

				releases := []*github.RepositoryRelease{
					{
						TagName:    &tag,
						Prerelease: &prerelease,
					},
				}
				w.Header().Set("Content-Type", "application/json")
				require.NoError(t, json.NewEncoder(w).Encode(releases))
			},
			wantReleaseCount: 1,
			wantPreRelease:   false,
			wantErr:          false,
		},
		{
			name: "Prerelease is nil - should be treated as stable",
			opts: &ReleaseListOptions{
				ExcludePreRelease: true,
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				tag1 := "v1.0.0"
				tag2 := "v1.1.0"
				prerelease1 := false

				releases := []*github.RepositoryRelease{
					{
						TagName:    &tag1,
						Prerelease: &prerelease1,
					},
					{
						TagName:    &tag2,
						Prerelease: nil,
					},
				}
				w.Header().Set("Content-Type", "application/json")
				require.NoError(t, json.NewEncoder(w).Encode(releases))
			},
			wantReleaseCount: 2,
			wantPreRelease:   false,
			wantErr:          false,
		},
		{
			name: "github api error",
			opts: &ReleaseListOptions{
				ExcludePreRelease: false,
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			wantReleaseCount: 0,
			wantErr:          true,
			wantErrString:    "failed to list releases",
		},
		{
			name: "empty releases list",
			opts: &ReleaseListOptions{
				ExcludePreRelease: true,
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				releases := []*github.RepositoryRelease{}
				w.Header().Set("Content-Type", "application/json")
				require.NoError(t, json.NewEncoder(w).Encode(releases))
			},
			wantReleaseCount: 0,
			wantPreRelease:   false,
			wantErr:          false,
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

			repo := &gh.Repository{Owner: "owner", Name: "repo"}
			opts := ReleaseListOptions{}
			if tt.opts != nil {
				opts = *tt.opts
			}

			releases, err := NewClientFromGitHubClient(client).
				ListReleases(context.Background(), repo, opts)

			if tt.wantErr {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.wantErrString)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, repo.Releases)
			assert.Len(t, releases, tt.wantReleaseCount)
			assert.Len(t, repo.Releases, len(releases))

			if tt.wantReleaseCount > 0 {
				hasPreRelease := false
				for _, release := range releases {
					if release.Prerelease != nil && *release.Prerelease {
						hasPreRelease = true
						break
					}
				}
				assert.Equal(t, tt.wantPreRelease, hasPreRelease)
			}
		})
	}

	t.Run("repository is nil", func(t *testing.T) {
		client := NewClient("", nil)
		_, err := client.ListReleases(context.Background(), nil, ReleaseListOptions{})
		require.Error(t, err)
		assert.ErrorContains(t, err, "repository is nil")
	})
}

func TestGetStableReleaseList(t *testing.T) {
	tests := []struct {
		name             string
		handler          http.HandlerFunc
		wantReleaseCount int
		wantErr          bool
		wantErrString    string
	}{
		{
			name: "success - excludes prereleases",
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/repos/owner/repo/releases", r.URL.Path)

				tag1 := "v1.0.0"
				tag2 := "v1.1.0-beta"
				tag3 := "v1.2.0"
				tag4 := "v1.3.0-alpha"
				prerelease1 := false
				prerelease2 := true
				prerelease3 := false
				prerelease4 := true

				releases := []*github.RepositoryRelease{
					{
						TagName:    &tag1,
						Prerelease: &prerelease1,
					},
					{
						TagName:    &tag2,
						Prerelease: &prerelease2,
					},
					{
						TagName:    &tag3,
						Prerelease: &prerelease3,
					},
					{
						TagName:    &tag4,
						Prerelease: &prerelease4,
					},
				}
				w.Header().Set("Content-Type", "application/json")
				require.NoError(t, json.NewEncoder(w).Encode(releases))
			},
			wantReleaseCount: 2,
			wantErr:          false,
		},
		{
			name: "all prereleases - returns empty",
			handler: func(w http.ResponseWriter, r *http.Request) {
				tag1 := "v1.0.0-beta"
				tag2 := "v1.1.0-alpha"
				prerelease1 := true
				prerelease2 := true

				releases := []*github.RepositoryRelease{
					{
						TagName:    &tag1,
						Prerelease: &prerelease1,
					},
					{
						TagName:    &tag2,
						Prerelease: &prerelease2,
					},
				}
				w.Header().Set("Content-Type", "application/json")
				require.NoError(t, json.NewEncoder(w).Encode(releases))
			},
			wantReleaseCount: 0,
			wantErr:          false,
		},
		{
			name: "github api error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			wantReleaseCount: 0,
			wantErr:          true,
			wantErrString:    "failed to list releases",
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

			repo := &gh.Repository{Owner: "owner", Name: "repo"}

			releases, err := NewClientFromGitHubClient(client).
				ListStableReleases(context.Background(), repo, github.ListOptions{})

			if tt.wantErr {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.wantErrString)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, repo.Releases)
			assert.Len(t, releases, tt.wantReleaseCount)
			assert.Len(t, repo.Releases, len(releases))

			for i, release := range releases {
				if release.Prerelease != nil && *release.Prerelease {
					assert.Failf(t, "unexpected prerelease", "release[%d] is prerelease", i)
				}
			}
		})
	}

	t.Run("repository is nil", func(t *testing.T) {
		client := NewClient("", nil)
		_, err := client.ListStableReleases(context.Background(), nil, github.ListOptions{})
		require.Error(t, err)
		assert.ErrorContains(t, err, "repository is nil")
	})
}
