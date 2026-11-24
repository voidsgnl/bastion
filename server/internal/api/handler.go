package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/voidsgnl/bastion/server/internal/scope"
	"github.com/voidsgnl/bastion/server/internal/store"
)

type Handler struct {
	store  *store.Memory
	mux    *http.ServeMux
	static http.Handler
}

func NewHandler(memory *store.Memory) *Handler {
	h := &Handler{store: memory, mux: http.NewServeMux()}
	h.mux.HandleFunc("GET /health", h.health)
	h.mux.HandleFunc("GET /api/v1/targets", h.listTargets)
	h.mux.HandleFunc("POST /api/v1/targets", h.createTarget)
	h.mux.HandleFunc("POST /api/v1/targets/{target_id}/attestation", h.attestTarget)
	h.mux.HandleFunc("POST /api/v1/scans", h.createScan)
	h.mux.HandleFunc("GET /api/v1/scans/{scan_id}", h.getScan)
	h.mux.HandleFunc("GET /api/v1/scans/{scan_id}/findings", h.getScanFindings)
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
	writeJSON(w, http.StatusOK, h.store.ListTargets())
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
	target := h.store.CreateTarget(store.Target{
		Name: input.Name, Type: input.Type, Value: input.Value,
	})
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
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "target not found"})
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
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "target not found"})
		return
	}
	if !target.Authorized {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "target is not authorized"})
		return
	}
	if input.Profile == "" {
		input.Profile = "safe"
	}
	writeJSON(w, http.StatusCreated, h.store.CreateScan(target, input.Profile))
}

func (h *Handler) getScan(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("scan_id")
	scan, err := h.store.GetScan(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "scan not found"})
		return
	}
	writeJSON(w, http.StatusOK, scan)
}

func (h *Handler) getScanFindings(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("scan_id")
	if _, err := h.store.GetScan(id); err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "scan not found"})
		return
	}
	writeJSON(w, http.StatusOK, h.store.Findings(id))
}

func (h *Handler) getFindings(w http.ResponseWriter, r *http.Request) {
	scanID := r.URL.Query().Get("scan_id")
	if scanID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "scan_id is required"})
		return
	}
	writeJSON(w, http.StatusOK, h.store.Findings(scanID))
}
