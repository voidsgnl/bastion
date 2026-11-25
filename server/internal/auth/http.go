package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

const sessionCookie = "bastion_session"

type Manager struct {
	username     string
	passwordHash [32]byte
	ttl          time.Duration
	mu           sync.Mutex
	sessions     map[string]time.Time
}

func NewManager(username, password string, ttl time.Duration) (*Manager, error) {
	if strings.TrimSpace(username) == "" || len(password) < 12 {
		return nil, errors.New("auth requires a username and a password of at least 12 characters")
	}
	if ttl <= 0 {
		ttl = 12 * time.Hour
	}
	return &Manager{
		username: username, passwordHash: sha256.Sum256([]byte(password)), ttl: ttl,
		sessions: make(map[string]time.Time),
	}, nil
}

// Middleware provides single-operator login and protects all versioned API routes.
// Session state is intentionally process-local for this first deployment slice.
func (m *Manager) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/auth/login" && r.Method == http.MethodPost {
			m.login(w, r)
			return
		}
		if r.URL.Path == "/api/v1/auth/logout" && r.Method == http.MethodPost {
			m.logout(w, r)
			return
		}
		if r.URL.Path == "/api/v1/auth/me" && r.Method == http.MethodGet {
			if !m.authenticated(r) {
				writeError(w, http.StatusUnauthorized, "login required")
				return
			}
			writeJSON(w, http.StatusOK, map[string]string{"username": m.username})
			return
		}
		if strings.HasPrefix(r.URL.Path, "/api/v1/") && !m.authenticated(r) {
			writeError(w, http.StatusUnauthorized, "login required")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (m *Manager) login(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&credentials); err != nil {
		writeError(w, http.StatusBadRequest, "invalid login request")
		return
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		writeError(w, http.StatusBadRequest, "invalid login request")
		return
	}
	userHash := sha256.Sum256([]byte(credentials.Username))
	wantUserHash := sha256.Sum256([]byte(m.username))
	passHash := sha256.Sum256([]byte(credentials.Password))
	if subtle.ConstantTimeCompare(userHash[:], wantUserHash[:]) != 1 ||
		subtle.ConstantTimeCompare(passHash[:], m.passwordHash[:]) != 1 {
		writeError(w, http.StatusUnauthorized, "invalid username or password")
		return
	}

	var tokenBytes [32]byte
	if _, err := rand.Read(tokenBytes[:]); err != nil {
		writeError(w, http.StatusInternalServerError, "could not create session")
		return
	}
	token := hex.EncodeToString(tokenBytes[:])
	expires := time.Now().Add(m.ttl)
	m.mu.Lock()
	for existing, expiry := range m.sessions {
		if time.Now().After(expiry) {
			delete(m.sessions, existing)
		}
	}
	m.sessions[token] = expires
	m.mu.Unlock()
	http.SetCookie(w, &http.Cookie{
		Name: sessionCookie, Value: token, Path: "/", Expires: expires,
		MaxAge: int(m.ttl.Seconds()), HttpOnly: true, Secure: r.TLS != nil,
		SameSite: http.SameSiteStrictMode,
	})
	writeJSON(w, http.StatusOK, map[string]any{"username": m.username, "expires_at": expires.UTC()})
}

func (m *Manager) logout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(sessionCookie); err == nil {
		m.mu.Lock()
		delete(m.sessions, cookie.Value)
		m.mu.Unlock()
	}
	http.SetCookie(w, &http.Cookie{
		Name: sessionCookie, Value: "", Path: "/", MaxAge: -1,
		HttpOnly: true, Secure: r.TLS != nil, SameSite: http.SameSiteStrictMode,
	})
	w.WriteHeader(http.StatusNoContent)
}

func (m *Manager) authenticated(r *http.Request) bool {
	cookie, err := r.Cookie(sessionCookie)
	if err != nil || len(cookie.Value) != 64 {
		return false
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	expires, ok := m.sessions[cookie.Value]
	if !ok {
		return false
	}
	if time.Now().After(expires) {
		delete(m.sessions, cookie.Value)
		return false
	}
	return true
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
