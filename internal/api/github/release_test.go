package gh

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/google/go-github/v77/github"
)

func TestGetReleaseList(t *testing.T) {
	mustParseURL := func(rawURL string) *url.URL {
		u, err := url.Parse(rawURL)
		if err != nil {
			t.Fatalf("test setup: failed to parse URL: %v", err)
		}
		return u
	}

	tests := []struct {
		name             string
		repoURL          *url.URL
		opts             *ReleaseListOptions
		handler          http.HandlerFunc
		wantReleaseCount int
		wantPreRelease   bool
		wantErr          bool
		wantErrString    string
	}{
		{
			name:    "ExcludePreRelease true - stable releases only",
			repoURL: mustParseURL("https://github.com/owner/repo"),
			opts: &ReleaseListOptions{
				ExcludePreRelease: true,
				ListOptions:       &github.ListOptions{},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Expected GET request, got %s", r.Method)
				}
				expectedPath := "/repos/owner/repo/releases"
				if r.URL.Path != expectedPath {
					t.Errorf("Request path = %q, want %q", r.URL.Path, expectedPath)
				}

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
				jsonBytes, err := json.Marshal(releases)
				if err != nil {
					t.Fatalf("failed to marshal test data: %v", err)
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write(jsonBytes)
			},
			wantReleaseCount: 2, // v1.0.0 and v1.2.0 only
			wantPreRelease:   false,
			wantErr:          false,
		},
		{
			name:    "ExcludePreRelease false - all releases",
			repoURL: mustParseURL("https://github.com/owner/repo"),
			opts: &ReleaseListOptions{
				ExcludePreRelease: false,
				ListOptions:       &github.ListOptions{},
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
				jsonBytes, err := json.Marshal(releases)
				if err != nil {
					t.Fatalf("failed to marshal test data: %v", err)
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write(jsonBytes)
			},
			wantReleaseCount: 3,    // all releases
			wantPreRelease:   true, // includes prerelease
			wantErr:          false,
		},
		{
			name:    "opts is nil - default behavior",
			repoURL: mustParseURL("https://github.com/owner/repo"),
			opts:    nil,
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
				jsonBytes, err := json.Marshal(releases)
				if err != nil {
					t.Fatalf("failed to marshal test data: %v", err)
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write(jsonBytes)
			},
			wantReleaseCount: 2, // all releases (default is false)
			wantPreRelease:   true,
			wantErr:          false,
		},
		{
			name:    "ListOptions is nil - should use default",
			repoURL: mustParseURL("https://github.com/owner/repo"),
			opts: &ReleaseListOptions{
				ExcludePreRelease: true,
				ListOptions:       nil,
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
				jsonBytes, err := json.Marshal(releases)
				if err != nil {
					t.Fatalf("failed to marshal test data: %v", err)
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write(jsonBytes)
			},
			wantReleaseCount: 1,
			wantPreRelease:   false,
			wantErr:          false,
		},
		{
			name:    "Prerelease is nil - should be treated as stable",
			repoURL: mustParseURL("https://github.com/owner/repo"),
			opts: &ReleaseListOptions{
				ExcludePreRelease: true,
				ListOptions:       &github.ListOptions{},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				tag1 := "v1.0.0"
				tag2 := "v1.1.0"
				prerelease1 := false
				// prerelease2 is nil

				releases := []*github.RepositoryRelease{
					{
						TagName:    &tag1,
						Prerelease: &prerelease1,
					},
					{
						TagName:    &tag2,
						Prerelease: nil, // nil should be treated as false
					},
				}
				jsonBytes, err := json.Marshal(releases)
				if err != nil {
					t.Fatalf("failed to marshal test data: %v", err)
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write(jsonBytes)
			},
			wantReleaseCount: 2, // both should be included
			wantPreRelease:   false,
			wantErr:          false,
		},
		{
			name:    "github api error",
			repoURL: mustParseURL("https://github.com/owner/repo"),
			opts: &ReleaseListOptions{
				ExcludePreRelease: false,
				ListOptions:       &github.ListOptions{},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			wantReleaseCount: 0,
			wantErr:          true,
			wantErrString:    "failed to list releases",
		},
		{
			name:    "empty releases list",
			repoURL: mustParseURL("https://github.com/owner/repo"),
			opts: &ReleaseListOptions{
				ExcludePreRelease: true,
				ListOptions:       &github.ListOptions{},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				releases := []*github.RepositoryRelease{}
				jsonBytes, err := json.Marshal(releases)
				if err != nil {
					t.Fatalf("failed to marshal test data: %v", err)
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write(jsonBytes)
			},
			wantReleaseCount: 0,
			wantPreRelease:   false,
			wantErr:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var client *github.Client

			if tt.handler != nil {
				server := httptest.NewServer(tt.handler)
				defer server.Close()

				client = github.NewClient(server.Client())
				baseURL, err := url.Parse(server.URL + "/")
				if err != nil {
					t.Fatalf("Failed to parse server URL: %v", err)
				}
				client.BaseURL = baseURL
			} else {
				client = github.NewClient(nil)
			}

			state := &RepositoryState{
				Client: client,
				Owner:  "owner",
				Repo:   "repo",
			}

			err := state.GetReleaseList(tt.opts)

			if (err != nil) != tt.wantErr {
				t.Fatalf("GetReleaseList() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if !strings.Contains(err.Error(), tt.wantErrString) {
					t.Errorf("GetReleaseList() error = %q, want to contain %q", err.Error(), tt.wantErrString)
				}
				return
			}

			if len(state.Releases) != tt.wantReleaseCount {
				t.Errorf("GetReleaseList() release count = %v, want %v", len(state.Releases), tt.wantReleaseCount)
			}

			if tt.wantReleaseCount > 0 {
				hasPreRelease := false
				for _, release := range state.Releases {
					if release.Prerelease != nil && *release.Prerelease {
						hasPreRelease = true
						break
					}
				}
				if hasPreRelease != tt.wantPreRelease {
					t.Errorf("GetReleaseList() hasPreRelease = %v, want %v", hasPreRelease, tt.wantPreRelease)
				}
			}
		})
	}
}

func TestGetStableReleaseList(t *testing.T) {
	mustParseURL := func(rawURL string) *url.URL {
		u, err := url.Parse(rawURL)
		if err != nil {
			t.Fatalf("test setup: failed to parse URL: %v", err)
		}
		return u
	}

	tests := []struct {
		name             string
		repoURL          *url.URL
		handler          http.HandlerFunc
		wantReleaseCount int
		wantErr          bool
		wantErrString    string
	}{
		{
			name:    "success - excludes prereleases",
			repoURL: mustParseURL("https://github.com/owner/repo"),
			handler: func(w http.ResponseWriter, r *http.Request) {
				expectedPath := "/repos/owner/repo/releases"
				if r.URL.Path != expectedPath {
					t.Errorf("Request path = %q, want %q", r.URL.Path, expectedPath)
				}

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
				jsonBytes, err := json.Marshal(releases)
				if err != nil {
					t.Fatalf("failed to marshal test data: %v", err)
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write(jsonBytes)
			},
			wantReleaseCount: 2, // v1.0.0 and v1.2.0 only
			wantErr:          false,
		},
		{
			name:    "all prereleases - returns empty",
			repoURL: mustParseURL("https://github.com/owner/repo"),
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
				jsonBytes, err := json.Marshal(releases)
				if err != nil {
					t.Fatalf("failed to marshal test data: %v", err)
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write(jsonBytes)
			},
			wantReleaseCount: 0,
			wantErr:          false,
		},
		{
			name:    "github api error",
			repoURL: mustParseURL("https://github.com/owner/repo"),
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
			var client *github.Client

			if tt.handler != nil {
				server := httptest.NewServer(tt.handler)
				defer server.Close()

				client = github.NewClient(server.Client())
				baseURL, err := url.Parse(server.URL + "/")
				if err != nil {
					t.Fatalf("Failed to parse server URL: %v", err)
				}
				client.BaseURL = baseURL
			} else {
				client = github.NewClient(nil)
			}

			state := &RepositoryState{
				Client: client,
				Owner:  "owner",
				Repo:   "repo",
			}

			err := state.GetStableReleaseList()

			if (err != nil) != tt.wantErr {
				t.Fatalf("GetStableReleaseList() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if !strings.Contains(err.Error(), tt.wantErrString) {
					t.Errorf("GetStableReleaseList() error = %q, want to contain %q", err.Error(), tt.wantErrString)
				}
				return
			}

			if len(state.Releases) != tt.wantReleaseCount {
				t.Errorf("GetStableReleaseList() release count = %v, want %v", len(state.Releases), tt.wantReleaseCount)
			}

			for i, release := range state.Releases {
				if release.Prerelease != nil && *release.Prerelease {
					t.Errorf("GetStableReleaseList() release[%d] is prerelease, but should be stable", i)
				}
			}
		})
	}
}
