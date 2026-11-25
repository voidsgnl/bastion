package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/voidsgnl/bastion/server/internal/engine"
	"github.com/voidsgnl/bastion/server/internal/scope"
	"github.com/voidsgnl/bastion/server/internal/store"
)

type Handler struct {
	store  store.Repository
	runner engine.Runner
	mux    *http.ServeMux
	static http.Handler
}

func NewHandler(memory *store.Memory) *Handler {
	return NewHandlerWithRunner(memory, engine.FixtureRunner{})
}

func NewHandlerWithRunner(memory *store.Memory, runner engine.Runner) *Handler {
	return NewHandlerWithRepository(memory, runner)
}

func NewHandlerWithRepository(repository store.Repository, runner engine.Runner) *Handler {
	h := &Handler{store: repository, runner: runner, mux: http.NewServeMux()}
	h.mux.HandleFunc("GET /health", h.health)
	h.mux.HandleFunc("GET /api/v1/targets", h.listTargets)
	h.mux.HandleFunc("POST /api/v1/targets", h.createTarget)
	h.mux.HandleFunc("POST /api/v1/targets/{target_id}/attestation", h.attestTarget)
	h.mux.HandleFunc("POST /api/v1/scans", h.createScan)
	h.mux.HandleFunc("POST /api/v1/scans/{scan_id}/retest", h.retestScan)
	h.mux.HandleFunc("GET /api/v1/scans/{scan_id}", h.getScan)
	h.mux.HandleFunc("GET /api/v1/scans/{scan_id}/findings", h.getScanFindings)
	h.mux.HandleFunc("GET /api/v1/scans/{scan_id}/compare/{other_scan_id}", h.compareScans)
	h.mux.HandleFunc("GET /api/v1/findings", h.getFindings)
	return h
}

func (h *Handler) SetStatic(handler http.Handler) {
	h.static = handler
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.static != nil && (r.URL.Path == "/" || r.URL.Path == "/index.html" || r.URL.Path == "/styles.css" || r.URL.Path == "/app.js" || r.URL.Path == "/results.js" || strings.HasSuffix(r.URL.Path, ".html")) {
		h.static.ServeHTTP(w, r)
		return
	}
	h.mux.ServeHTTP(w, r)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) listTargets(w http.ResponseWriter, _ *http.Request) {
	targets, err := h.store.ListTargets()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, targets)
}

func (h *Handler) createTarget(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name  string `json:"name"`
		Type  string `json:"type"`
		Value string `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	if input.Name == "" || input.Type == "" || input.Value == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "name, type, and value are required"})
		return
	}
	switch input.Type {
	case "cidr":
		prefix, err := scope.ParseCIDR(input.Value)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		input.Value = prefix.String()
	case "domain", "url", "repository", "cloud_account":
	default:
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "unsupported target type"})
		return
	}
	target, err := h.store.CreateTarget(store.Target{
		Name: input.Name, Type: input.Type, Value: input.Value,
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, target)
}

func (h *Handler) attestTarget(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Statement string `json:"statement"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	if strings.TrimSpace(input.Statement) == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "attestation statement is required"})
		return
	}
	target, err := h.store.AttestTarget(r.PathValue("target_id"), strings.TrimSpace(input.Statement))
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "target not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, target)
}

func (h *Handler) createScan(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TargetID string `json:"target_id"`
		Profile  string `json:"profile"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	target, err := h.store.GetTarget(input.TargetID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "target not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if !target.Authorized {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "target is not authorized"})
		return
	}
	if input.Profile == "" {
		input.Profile = "safe"
	}
	scan, err := h.runScan(r.Context(), target, input.Profile)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error(), "scan_id": scan.ID})
		return
	}
	writeJSON(w, http.StatusCreated, scan)
}

func (h *Handler) runScan(ctx context.Context, target store.Target, profile string) (store.Scan, error) {
	scan, err := h.store.CreateScan(target, profile)
	if err != nil {
		return store.Scan{}, fmt.Errorf("create scan: %w", err)
	}
	findings, err := h.runner.Run(ctx, scan, target)
	if err != nil {
		if failErr := h.store.FailScan(scan.ID); failErr != nil {
			return scan, fmt.Errorf("scan runner failed (%v) and status update failed: %w", err, failErr)
		}
		return scan, fmt.Errorf("scan runner failed: %w", err)
	}
	completed, err := h.store.CompleteScan(scan.ID, findings)
	if err != nil {
		if failErr := h.store.FailScan(scan.ID); failErr != nil {
			return scan, fmt.Errorf("persist scan findings failed (%v) and status update failed: %w", err, failErr)
		}
		return scan, fmt.Errorf("persist scan findings: %w", err)
	}
	return completed, nil
}

func (h *Handler) retestScan(w http.ResponseWriter, r *http.Request) {
	previous, err := h.store.GetScan(r.PathValue("scan_id"))
	if err != nil {
		h.writeRepositoryError(w, "scan not found", err)
		return
	}
	if previous.Status != "complete" {
		writeJSON(w, http.StatusConflict, map[string]string{"error": "only completed scans can be retested"})
		return
	}
	if previous.Profile != "safe" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "only safe-profile scans can be retested"})
		return
	}
	target, err := h.store.GetTarget(previous.TargetID)
	if err != nil {
		h.writeRepositoryError(w, "target not found", err)
		return
	}
	if !target.Authorized {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "target is not authorized"})
		return
	}
	scan, err := h.runScan(r.Context(), target, previous.Profile)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error(), "scan_id": scan.ID})
		return
	}
	writeJSON(w, http.StatusCreated, scan)
}

func (h *Handler) getScan(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("scan_id")
	scan, err := h.store.GetScan(id)
	if err != nil {
		h.writeRepositoryError(w, "scan not found", err)
		return
	}
	writeJSON(w, http.StatusOK, scan)
}

func (h *Handler) getScanFindings(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("scan_id")
	findings, err := h.store.Findings(id)
	if err != nil {
		h.writeRepositoryError(w, "scan not found", err)
		return
	}
	writeJSON(w, http.StatusOK, findings)
}

func (h *Handler) getFindings(w http.ResponseWriter, r *http.Request) {
	scanID := r.URL.Query().Get("scan_id")
	if scanID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "scan_id is required"})
		return
	}
	findings, err := h.store.Findings(scanID)
	if err != nil {
		h.writeRepositoryError(w, "scan not found", err)
		return
	}
	writeJSON(w, http.StatusOK, findings)
}

func (h *Handler) compareScans(w http.ResponseWriter, r *http.Request) {
	baseline, err := h.store.GetScan(r.PathValue("scan_id"))
	if err != nil {
		h.writeRepositoryError(w, "baseline scan not found", err)
		return
	}
	current, err := h.store.GetScan(r.PathValue("other_scan_id"))
	if err != nil {
		h.writeRepositoryError(w, "current scan not found", err)
		return
	}
	if baseline.Status != "complete" || current.Status != "complete" {
		writeJSON(w, http.StatusConflict, map[string]string{"error": "both scans must be complete"})
		return
	}
	if baseline.TargetID != current.TargetID {
		writeJSON(w, http.StatusConflict, map[string]string{"error": store.ErrDifferentScope.Error()})
		return
	}
	baselineFindings, err := h.store.Findings(baseline.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	currentFindings, err := h.store.Findings(current.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, store.CompareFindings(baseline.ID, current.ID, baselineFindings, currentFindings))
}

func (h *Handler) writeRepositoryError(w http.ResponseWriter, notFoundMessage string, err error) {
	if errors.Is(err, store.ErrNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": notFoundMessage})
		return
	}
	writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
}
