package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/voidsgnl/bastion/server/internal/engine"
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

type failingRunner struct{}

func (failingRunner) Run(context.Context, store.Scan, store.Target) ([]store.Finding, error) {
	return nil, context.DeadlineExceeded
}

type sequenceRunner struct {
	results [][]store.Finding
	next    int
}

type failingCompleteRepository struct {
	store.Repository
}

func (failingCompleteRepository) CompleteScan(string, []store.Finding) (store.Scan, error) {
	return store.Scan{}, errors.New("simulated persistence failure")
}

func (r *sequenceRunner) Run(_ context.Context, _ store.Scan, _ store.Target) ([]store.Finding, error) {
	if r.next >= len(r.results) {
		return nil, context.DeadlineExceeded
	}
	result := append([]store.Finding(nil), r.results[r.next]...)
	r.next++
	return result, nil
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
	findingsResponse := httptest.NewRecorder()
	handler.ServeHTTP(findingsResponse, httptest.NewRequest(http.MethodGet,
		"/api/v1/scans/"+scan.ID+"/findings", nil))
	var findings []store.Finding
	if err := json.Unmarshal(findingsResponse.Body.Bytes(), &findings); err != nil {
		t.Fatal(err)
	}
	if len(findings) != 1 || findings[0].SourceEngine != "fixture" {
		t.Fatalf("fixture findings not returned: %#v", findings)
	}
	if findings[0].Confidence != "unverified" {
		t.Fatalf("fixture confidence = %q, want unverified", findings[0].Confidence)
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

func TestRunnerFailureIsReportedAndRecorded(t *testing.T) {
	handler := NewHandlerWithRunner(store.NewMemory(), failingRunner{})
	targetResponse := post(t, handler, "/api/v1/targets", map[string]string{
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

	response := post(t, handler, "/api/v1/scans", map[string]string{
		"target_id": target.ID, "profile": "safe",
	})
	if response.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusInternalServerError)
	}
	var body struct {
		ScanID string `json:"scan_id"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	scanResponse := httptest.NewRecorder()
	handler.ServeHTTP(scanResponse, httptest.NewRequest(http.MethodGet, "/api/v1/scans/"+body.ScanID, nil))
	var scan store.Scan
	if err := json.Unmarshal(scanResponse.Body.Bytes(), &scan); err != nil {
		t.Fatal(err)
	}
	if scan.Status != "failed" {
		t.Fatalf("scan status = %q, want failed", scan.Status)
	}
}

func TestFindingPersistenceFailureMarksScanFailed(t *testing.T) {
	memory := store.NewMemory()
	handler := NewHandlerWithRepository(
		failingCompleteRepository{Repository: memory},
		engine.FixtureRunner{},
	)
	targetResponse := post(t, handler, "/api/v1/targets", map[string]string{
		"name": "lab", "type": "domain", "value": "lab.example",
	})
	var target store.Target
	if err := json.Unmarshal(targetResponse.Body.Bytes(), &target); err != nil {
		t.Fatal(err)
	}
	if response := post(t, handler, "/api/v1/targets/"+target.ID+"/attestation",
		map[string]string{"statement": "Authorized test target."}); response.Code != http.StatusOK {
		t.Fatalf("attestation status = %d", response.Code)
	}

	response := post(t, handler, "/api/v1/scans", map[string]string{
		"target_id": target.ID, "profile": "safe",
	})
	if response.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusInternalServerError)
	}
	var body struct {
		ScanID string `json:"scan_id"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	scan, err := memory.GetScan(body.ScanID)
	if err != nil {
		t.Fatal(err)
	}
	if scan.Status != "failed" {
		t.Fatalf("scan status = %q, want failed", scan.Status)
	}
}

type fakeNmapExecutor struct {
	output []byte
	called bool
}

func (e *fakeNmapExecutor) Run(_ context.Context, _ string, _ ...string) ([]byte, error) {
	e.called = true
	return e.output, nil
}

func TestNmapLabRunnerCompletesAPIFlowWithoutNetworkProbe(t *testing.T) {
	executor := &fakeNmapExecutor{output: []byte(`<nmaprun><host><ports><port protocol="tcp" portid="18080"><state state="open"/></port></ports></host></nmaprun>`)}
	runner := engine.NmapRunner{Executor: executor}
	handler := NewHandlerWithRunner(store.NewMemory(), runner)

	targetResponse := post(t, handler, "/api/v1/targets", map[string]string{
		"name": "local lab", "type": "cidr", "value": "127.0.0.1/32",
	})
	var target store.Target
	if err := json.Unmarshal(targetResponse.Body.Bytes(), &target); err != nil {
		t.Fatal(err)
	}
	attestationResponse := post(t, handler, "/api/v1/targets/"+target.ID+"/attestation", map[string]string{
		"statement": "This is my loopback lab.",
	})
	if attestationResponse.Code != http.StatusOK {
		t.Fatalf("attestation status = %d", attestationResponse.Code)
	}

	scanResponse := post(t, handler, "/api/v1/scans", map[string]string{
		"target_id": target.ID, "profile": "safe",
	})
	if scanResponse.Code != http.StatusCreated {
		t.Fatalf("scan status = %d: %s", scanResponse.Code, scanResponse.Body.String())
	}
	if !executor.called {
		t.Fatal("expected Nmap executor to be called")
	}
	var scan store.Scan
	if err := json.Unmarshal(scanResponse.Body.Bytes(), &scan); err != nil {
		t.Fatal(err)
	}

	findingsResponse := httptest.NewRecorder()
	handler.ServeHTTP(findingsResponse, httptest.NewRequest(http.MethodGet,
		"/api/v1/scans/"+scan.ID+"/findings", nil))
	var findings []store.Finding
	if err := json.Unmarshal(findingsResponse.Body.Bytes(), &findings); err != nil {
		t.Fatal(err)
	}
	if len(findings) != 1 || findings[0].SourceEngine != "nmap" {
		t.Fatalf("normalized findings = %#v", findings)
	}
	if findings[0].Severity != "info" {
		t.Fatalf("severity = %q, want info", findings[0].Severity)
	}
}

func TestSafeRetestAndComparisonLifecycle(t *testing.T) {
	firstSeen := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	runner := &sequenceRunner{results: [][]store.Finding{
		{
			{ID: "stable", Title: "still open", FirstSeen: firstSeen},
			{ID: "resolved", Title: "now fixed", FirstSeen: firstSeen},
		},
		{
			{ID: "stable", Title: "still open"},
			{ID: "new", Title: "new observation"},
		},
	}}
	handler := NewHandlerWithRunner(store.NewMemory(), runner)
	targetResponse := post(t, handler, "/api/v1/targets", map[string]string{
		"name": "lab", "type": "domain", "value": "lab.example",
	})
	var target store.Target
	if err := json.Unmarshal(targetResponse.Body.Bytes(), &target); err != nil {
		t.Fatal(err)
	}
	attestationResponse := post(t, handler, "/api/v1/targets/"+target.ID+"/attestation", map[string]string{
		"statement": "Authorized test target.",
	})
	if attestationResponse.Code != http.StatusOK {
		t.Fatalf("attestation status = %d", attestationResponse.Code)
	}
	firstResponse := post(t, handler, "/api/v1/scans", map[string]string{
		"target_id": target.ID, "profile": "safe",
	})
	if firstResponse.Code != http.StatusCreated {
		t.Fatalf("first scan status = %d", firstResponse.Code)
	}
	var first store.Scan
	if err := json.Unmarshal(firstResponse.Body.Bytes(), &first); err != nil {
		t.Fatal(err)
	}

	retestResponse := post(t, handler, "/api/v1/scans/"+first.ID+"/retest", map[string]string{})
	if retestResponse.Code != http.StatusCreated {
		t.Fatalf("retest status = %d: %s", retestResponse.Code, retestResponse.Body.String())
	}
	var second store.Scan
	if err := json.Unmarshal(retestResponse.Body.Bytes(), &second); err != nil {
		t.Fatal(err)
	}

	compareRequest := httptest.NewRequest(http.MethodGet,
		"/api/v1/scans/"+first.ID+"/compare/"+second.ID, nil)
	compareResponse := httptest.NewRecorder()
	handler.ServeHTTP(compareResponse, compareRequest)
	if compareResponse.Code != http.StatusOK {
		t.Fatalf("compare status = %d: %s", compareResponse.Code, compareResponse.Body.String())
	}
	var comparison store.Comparison
	if err := json.Unmarshal(compareResponse.Body.Bytes(), &comparison); err != nil {
		t.Fatal(err)
	}
	if len(comparison.New) != 1 || comparison.New[0].ID != "new" {
		t.Fatalf("new findings = %#v", comparison.New)
	}
	if len(comparison.Fixed) != 1 || comparison.Fixed[0].ID != "resolved" {
		t.Fatalf("fixed findings = %#v", comparison.Fixed)
	}
	if len(comparison.Unchanged) != 1 || comparison.Unchanged[0].ID != "stable" {
		t.Fatalf("unchanged findings = %#v", comparison.Unchanged)
	}
	secondFindings, err := handler.store.Findings(second.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !secondFindings[0].FirstSeen.Equal(firstSeen) {
		t.Fatalf("retest first_seen = %s, want %s", secondFindings[0].FirstSeen, firstSeen)
	}
}

func TestRetestRejectsNonSafeProfile(t *testing.T) {
	handler := NewHandler(store.NewMemory())
	targetResponse := post(t, handler, "/api/v1/targets", map[string]string{
		"name": "lab", "type": "domain", "value": "lab.example",
	})
	var target store.Target
	if err := json.Unmarshal(targetResponse.Body.Bytes(), &target); err != nil {
		t.Fatal(err)
	}
	post(t, handler, "/api/v1/targets/"+target.ID+"/attestation", map[string]string{
		"statement": "Authorized test target.",
	})
	scanResponse := post(t, handler, "/api/v1/scans", map[string]string{
		"target_id": target.ID, "profile": "full",
	})
	var scan store.Scan
	if err := json.Unmarshal(scanResponse.Body.Bytes(), &scan); err != nil {
		t.Fatal(err)
	}
	response := post(t, handler, "/api/v1/scans/"+scan.ID+"/retest", map[string]string{})
	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusBadRequest)
	}
}

var _ engine.Runner = failingRunner{}
