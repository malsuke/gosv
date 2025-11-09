package osv

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/malsuke/govs/internal/misc/ptr"
)

func withTempOSVServer(t *testing.T, handler http.HandlerFunc) func() {
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	orig := osvAPIBaseURL
	osvAPIBaseURL = server.URL

	return func() {
		osvAPIBaseURL = orig
	}
}

func TestGetVulnerabilityByCVE_Success(t *testing.T) {
	cleanup := withTempOSVServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/vulns/CVE-2023-0001":
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(OsvVulnerability{Id: ptr.String("CVE-2023-0001")}); err != nil {
				t.Fatalf("failed to encode response: %v", err)
			}
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})
	defer cleanup()

	vuln, err := GetVulnerabilityByCVE("CVE-2023-0001")
	if err != nil {
		t.Fatalf("GetVulnerabilityByCVE() error = %v", err)
	}
	if vuln == nil || vuln.Id == nil || *vuln.Id != "CVE-2023-0001" {
		t.Fatalf("GetVulnerabilityByCVE() = %v, want CVE-2023-0001", vuln)
	}
}

func TestGetVulnerabilityByCVE_InvalidCVE(t *testing.T) {
	if _, err := GetVulnerabilityByCVE("invalid-cve"); err == nil {
		t.Fatalf("GetVulnerabilityByCVE() error = nil, want non-nil")
	}
}

func TestGetVulnerabilityByCVE_APIError(t *testing.T) {
	cleanup := withTempOSVServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(RpcStatus{Message: ptr.String("not found")}); err != nil {
			t.Fatalf("failed to encode response: %v", err)
		}
	})
	defer cleanup()

	if _, err := GetVulnerabilityByCVE("CVE-2023-0001"); err == nil {
		t.Fatalf("GetVulnerabilityByCVE() error = nil, want non-nil")
	}
}

func TestGetCveIDListFromGithubURL_Success(t *testing.T) {
	cleanup := withTempOSVServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/query":
			w.Header().Set("Content-Type", "application/json")
			response := V1VulnerabilityList{
				Vulns: &[]OsvVulnerability{
					{Id: ptr.String("CVE-2023-0001")},
					{Aliases: &[]string{"CVE-2023-0002"}},
				},
			}
			if err := json.NewEncoder(w).Encode(response); err != nil {
				t.Fatalf("failed to encode response: %v", err)
			}
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})
	defer cleanup()

	got, err := GetCveIDListFromGithubURL("https://example.com/repo.git")
	if err != nil {
		t.Fatalf("GetCveIDListFromGithubURL() error = %v", err)
	}

	want := []string{"CVE-2023-0001", "CVE-2023-0002"}
	if len(got) != len(want) {
		t.Fatalf("GetCveIDListFromGithubURL() = %v, want %v", got, want)
	}
	for i, id := range want {
		if got[i] != id {
			t.Fatalf("GetCveIDListFromGithubURL()[%d] = %v, want %v", i, got[i], id)
		}
	}
}

func TestGetCveIDListFromGithubURL_Empty(t *testing.T) {
	cleanup := withTempOSVServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/query":
			w.Header().Set("Content-Type", "application/json")
			response := V1VulnerabilityList{}
			if err := json.NewEncoder(w).Encode(response); err != nil {
				t.Fatalf("failed to encode response: %v", err)
			}
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})
	defer cleanup()

	got, err := GetCveIDListFromGithubURL("https://example.com/repo.git")
	if err != nil {
		t.Fatalf("GetCveIDListFromGithubURL() error = %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("GetCveIDListFromGithubURL() = %v, want empty", got)
	}
}

func TestGetCveIDListFromGithubURL_APIError(t *testing.T) {
	cleanup := withTempOSVServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	defer cleanup()

	if _, err := GetCveIDListFromGithubURL("https://example.com/repo.git"); err == nil {
		t.Fatalf("GetCveIDListFromGithubURL() error = nil, want non-nil")
	}
}
