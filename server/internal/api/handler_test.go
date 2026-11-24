package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/voidsgnl/bastion/server/internal/store"
)

func post(t *testing.T, handler http.Handler, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	raw, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(raw))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	return response
}

func TestAuthorizedTargetCompletesFixtureScan(t *testing.T) {
	handler := NewHandler(store.NewMemory())

	targetResponse := post(t, handler, "/api/v1/targets", map[string]any{
		"name": "lab", "type": "cidr", "value": "10.0.0.0/24",
	})
	if targetResponse.Code != http.StatusCreated {
		t.Fatalf("create target status = %d", targetResponse.Code)
	}
	var target store.Target
	if err := json.Unmarshal(targetResponse.Body.Bytes(), &target); err != nil {
		t.Fatal(err)
	}
	if target.Authorized {
		t.Fatal("new target must not be authorized")
	}
	attestResponse := post(t, handler, "/api/v1/targets/"+target.ID+"/attestation", map[string]string{
		"statement": "I am authorized to assess this lab network.",
	})
	if attestResponse.Code != http.StatusOK {
		t.Fatalf("attest target status = %d", attestResponse.Code)
	}

	scanResponse := post(t, handler, "/api/v1/scans", map[string]any{
		"target_id": target.ID, "profile": "safe",
	})
	if scanResponse.Code != http.StatusCreated {
		t.Fatalf("create scan status = %d", scanResponse.Code)
	}

	var scan store.Scan
	if err := json.Unmarshal(scanResponse.Body.Bytes(), &scan); err != nil {
		t.Fatal(err)
	}
	if scan.Status != "complete" {
		t.Fatalf("scan status = %q", scan.Status)
	}
}

func TestUnauthorizedTargetCannotBeScanned(t *testing.T) {
	handler := NewHandler(store.NewMemory())
	targetResponse := post(t, handler, "/api/v1/targets", map[string]any{
		"name": "unknown", "type": "domain", "value": "example.com",
	})
	var target store.Target
	if err := json.Unmarshal(targetResponse.Body.Bytes(), &target); err != nil {
		t.Fatal(err)
	}

	response := post(t, handler, "/api/v1/scans", map[string]any{
		"target_id": target.ID, "profile": "safe",
	})
	if response.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusForbidden)
	}
}

func TestClientCannotMarkTargetAuthorizedOnCreation(t *testing.T) {
	handler := NewHandler(store.NewMemory())
	response := post(t, handler, "/api/v1/targets", map[string]any{
		"name": "unknown", "type": "domain", "value": "example.com", "authorized": true,
	})
	var target store.Target
	if err := json.Unmarshal(response.Body.Bytes(), &target); err != nil {
		t.Fatal(err)
	}
	if target.Authorized {
		t.Fatal("client-supplied authorized flag must be ignored")
	}
}

func TestCIDRTargetMustBeCanonicalNetwork(t *testing.T) {
	handler := NewHandler(store.NewMemory())
	response := post(t, handler, "/api/v1/targets", map[string]string{
		"name": "lab", "type": "cidr", "value": "10.0.0.7/24",
	})
	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusBadRequest)
	}
}

func TestBlankAttestationIsRejected(t *testing.T) {
	handler := NewHandler(store.NewMemory())
	targetResponse := post(t, handler, "/api/v1/targets", map[string]string{
		"name": "lab", "type": "domain", "value": "lab.example",
	})
	var target store.Target
	if err := json.Unmarshal(targetResponse.Body.Bytes(), &target); err != nil {
		t.Fatal(err)
	}

	response := post(t, handler, "/api/v1/targets/"+target.ID+"/attestation", map[string]string{
		"statement": "   ",
	})
	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusBadRequest)
	}
}

func TestScanFindingsRoute(t *testing.T) {
	handler := NewHandler(store.NewMemory())
	targetResponse := post(t, handler, "/api/v1/targets", map[string]any{
		"name": "lab", "type": "domain", "value": "lab.example",
	})
	var target store.Target
	if err := json.Unmarshal(targetResponse.Body.Bytes(), &target); err != nil {
		t.Fatal(err)
	}
	attestResponse := post(t, handler, "/api/v1/targets/"+target.ID+"/attestation", map[string]string{
		"statement": "Authorized test target.",
	})
	if attestResponse.Code != http.StatusOK {
		t.Fatalf("attest target status = %d", attestResponse.Code)
	}
	scanResponse := post(t, handler, "/api/v1/scans", map[string]any{
		"target_id": target.ID, "profile": "safe",
	})
	var scan store.Scan
	if err := json.Unmarshal(scanResponse.Body.Bytes(), &scan); err != nil {
		t.Fatal(err)
	}

	request := httptest.NewRequest(http.MethodGet, "/api/v1/scans/"+scan.ID+"/findings", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	var findings []store.Finding
	if err := json.Unmarshal(response.Body.Bytes(), &findings); err != nil {
		t.Fatal(err)
	}
	if len(findings) != 1 {
		t.Fatalf("findings = %d, want 1", len(findings))
	}
}
