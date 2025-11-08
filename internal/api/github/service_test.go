package gh

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/google/go-github/v77/github"
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
				if r.Method != "GET" {
					t.Errorf("Expected GET request, got %s", r.Method)
				}
				expectedPath := "/repos/owner/repo/commits/abcdef123456/pulls"
				if r.URL.Path != expectedPath {
					t.Errorf("Request path = %q, want %q", r.URL.Path, expectedPath)
				}
				fmt.Fprint(w, `[{"number": 123}]`)
			},
			wantPRID: 123,
			wantErr:  false,
		},
		{
			name:       "no pull request found",
			repoURL:    mustParseURL("https://github.com/owner/repo"),
			commitHash: "fedcba654321",
			handler: func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `[]`)
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
			wantErrString: "invalid repo URL path: /owner",
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var client *github.Client

			if tt.handler != nil {
				server := httptest.NewServer(tt.handler)
				defer server.Close()

				client = github.NewClient(server.Client())
				url, err := url.Parse(server.URL + "/")
				if err != nil {
					t.Fatalf("Failed to parse server URL: %v", err)
				}
				client.BaseURL = url
			} else {
				client = github.NewClient(nil)
			}

			prID, err := GetPullRequestIDFromCommitHash(client, tt.repoURL, tt.commitHash)

			if (err != nil) != tt.wantErr {
				t.Fatalf("GetPullRequestIDFromCommitHash() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if !strings.Contains(err.Error(), tt.wantErrString) {
					t.Errorf("GetPullRequestIDFromCommitHash() error = %q, want to contain %q", err.Error(), tt.wantErrString)
				}
			}
			if prID != tt.wantPRID {
				t.Errorf("GetPullRequestIDFromCommitHash() = %v, want %v", prID, tt.wantPRID)
			}
		})
	}
}
