package store

import (
	"errors"
	"sync"
	"time"
)

var ErrNotFound = errors.New("not found")

type Target struct {
	ID                       string    `json:"id"`
	Name                     string    `json:"name"`
	Type                     string    `json:"type"`
	Value                    string    `json:"value"`
	Authorized               bool      `json:"authorized"`
	AuthorizationAttestation string    `json:"authorization_attestation,omitempty"`
	AuthorizedAt             time.Time `json:"authorized_at,omitempty"`
}

type Scan struct {
	ID        string    `json:"id"`
	TargetID  string    `json:"target_id"`
	Profile   string    `json:"profile"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type Finding struct {
	ID           string     `json:"id"`
	ScanID       string     `json:"scan_id"`
	AssetID      string     `json:"asset_id"`
	SourceEngine string     `json:"source_engine,omitempty"`
	Title        string     `json:"title"`
	Category     string     `json:"category"`
	Severity     string     `json:"severity"`
	Confidence   string     `json:"confidence"`
	Status       string     `json:"status"`
	Evidence     []Evidence `json:"evidence"`
	Risk         Risk       `json:"risk"`
	FirstSeen    time.Time  `json:"first_seen"`
	LastSeen     time.Time  `json:"last_seen"`
}

type Evidence struct {
	Type    string `json:"type"`
	Summary string `json:"summary"`
}

type Risk struct {
	Score       float64 `json:"score"`
	Explanation string  `json:"explanation"`
}

type Memory struct {
	mu       sync.RWMutex
	targets  map[string]Target
	scans    map[string]Scan
	findings map[string][]Finding
	seq      int
}

func NewMemory() *Memory {
	return &Memory{
		targets:  make(map[string]Target),
		scans:    make(map[string]Scan),
		findings: make(map[string][]Finding),
	}
}

func (m *Memory) next(prefix string) string {
	m.seq++
	return prefix + "_" + time.Now().UTC().Format("150405") + string(rune('a'+m.seq%26))
}

func (m *Memory) CreateTarget(target Target) Target {
	m.mu.Lock()
	defer m.mu.Unlock()
	target.ID = m.next("t")
	m.targets[target.ID] = target
	return target
}

func (m *Memory) ListTargets() []Target {
	m.mu.RLock()
	defer m.mu.RUnlock()
	targets := make([]Target, 0, len(m.targets))
	for _, target := range m.targets {
		targets = append(targets, target)
	}
	return targets
}

func (m *Memory) GetTarget(id string) (Target, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	target, ok := m.targets[id]
	if !ok {
		return Target{}, ErrNotFound
	}
	return target, nil
}

func (m *Memory) AttestTarget(id, statement string) (Target, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	target, ok := m.targets[id]
	if !ok {
		return Target{}, ErrNotFound
	}
	target.Authorized = true
	target.AuthorizationAttestation = statement
	target.AuthorizedAt = time.Now().UTC()
	m.targets[id] = target
	return target, nil
}

func (m *Memory) CreateScan(target Target, profile string) Scan {
	m.mu.Lock()
	defer m.mu.Unlock()
	scan := Scan{
		ID:        m.next("s"),
		TargetID:  target.ID,
		Profile:   profile,
		Status:    "running",
		CreatedAt: time.Now().UTC(),
	}
	m.scans[scan.ID] = scan
	m.findings[scan.ID] = []Finding{}
	return scan
}

func (m *Memory) CompleteScan(id string, findings []Finding) (Scan, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	scan, ok := m.scans[id]
	if !ok {
		return Scan{}, ErrNotFound
	}
	scan.Status = "complete"
	m.scans[id] = scan
	m.findings[id] = append([]Finding(nil), findings...)
	return scan, nil
}

func (m *Memory) FailScan(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	scan, ok := m.scans[id]
	if !ok {
		return ErrNotFound
	}
	scan.Status = "failed"
	m.scans[id] = scan
	return nil
}

func (m *Memory) GetScan(id string) (Scan, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	scan, ok := m.scans[id]
	if !ok {
		return Scan{}, ErrNotFound
	}
	return scan, nil
}

func (m *Memory) Findings(scanID string) []Finding {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return append([]Finding(nil), m.findings[scanID]...)
}
