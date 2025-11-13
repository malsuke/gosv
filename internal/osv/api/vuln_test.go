package osvapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/malsuke/govs/internal/common/ptr"
)

func withTempOSVServer(t *testing.T, handler http.HandlerFunc) func() {
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	orig := OsvAPIBaseURL
	OsvAPIBaseURL = server.URL

	return func() {
		OsvAPIBaseURL = orig
	}
}

func TestGetVulnerabilityByCVE_Success(t *testing.T) {
	cleanup := withTempOSVServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/vulns/CVE-2023-0001":
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(OsvVulnerability{Id: ptr.Ptr("CVE-2023-0001")}); err != nil {
				t.Fatalf("failed to encode response: %v", err)
			}
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})
	defer cleanup()

	ctx := context.Background()

	vuln, err := GetVulnerabilityByCVE(ctx, "CVE-2023-0001")
	if err != nil {
		t.Fatalf("GetVulnerabilityByCVE() error = %v", err)
	}
	if vuln == nil || vuln.Id == nil || *vuln.Id != "CVE-2023-0001" {
		t.Fatalf("GetVulnerabilityByCVE() = %v, want CVE-2023-0001", vuln)
	}
}

func TestGetVulnerabilityByCVE_InvalidCVE(t *testing.T) {
	ctx := context.Background()
	if _, err := GetVulnerabilityByCVE(ctx, "invalid-cve"); err == nil {
		t.Fatalf("GetVulnerabilityByCVE() error = nil, want non-nil")
	}
}

func TestGetVulnerabilityByCVE_APIError(t *testing.T) {
	cleanup := withTempOSVServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(RpcStatus{Message: ptr.Ptr("not found")}); err != nil {
			t.Fatalf("failed to encode response: %v", err)
		}
	})
	defer cleanup()

	ctx := context.Background()

	if _, err := GetVulnerabilityByCVE(ctx, "CVE-2023-0001"); err == nil {
		t.Fatalf("GetVulnerabilityByCVE() error = nil, want non-nil")
	}
}

func TestListCVEIDsByRepository_Success(t *testing.T) {
	cleanup := withTempOSVServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/query":
			w.Header().Set("Content-Type", "application/json")
			response := V1VulnerabilityList{
				Vulns: &[]OsvVulnerability{
					{Id: ptr.Ptr("CVE-2023-0001")},
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

	ctx := context.Background()

	got, err := ListCVEIDsByRepository(ctx, "example", "repo")
	if err != nil {
		t.Fatalf("ListCVEIDsByRepository() error = %v", err)
	}

	want := []string{"CVE-2023-0001", "CVE-2023-0002"}
	if len(got) != len(want) {
		t.Fatalf("ListCVEIDsByRepository() = %v, want %v", got, want)
	}
	for i, id := range want {
		if got[i] != id {
			t.Fatalf("ListCVEIDsByRepository()[%d] = %v, want %v", i, got[i], id)
		}
	}
}

func TestListCVEIDsByRepository_Empty(t *testing.T) {
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

	ctx := context.Background()

	got, err := ListCVEIDsByRepository(ctx, "example", "repo")
	if err != nil {
		t.Fatalf("ListCVEIDsByRepository() error = %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("ListCVEIDsByRepository() = %v, want empty", got)
	}
}

func TestListCVEIDsByRepository_APIError(t *testing.T) {
	cleanup := withTempOSVServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	defer cleanup()

	ctx := context.Background()

	if _, err := ListCVEIDsByRepository(ctx, "example", "repo"); err == nil {
		t.Fatalf("ListCVEIDsByRepository() error = nil, want non-nil")
	}
}
