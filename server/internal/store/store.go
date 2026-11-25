package store

import (
	"sync"
	"time"
)

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
	ID          string    `json:"id"`
	TargetID    string    `json:"target_id"`
	Profile     string    `json:"profile"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
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
}

func NewMemory() *Memory {
	return &Memory{
		targets:  make(map[string]Target),
		scans:    make(map[string]Scan),
		findings: make(map[string][]Finding),
	}
}

func (m *Memory) CreateTarget(target Target) (Target, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	id, err := newID("t")
	if err != nil {
		return Target{}, err
	}
	target.ID = id
	m.targets[target.ID] = target
	return target, nil
}

func (m *Memory) ListTargets() ([]Target, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	targets := make([]Target, 0, len(m.targets))
	for _, target := range m.targets {
		targets = append(targets, target)
	}
	return targets, nil
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

func (m *Memory) CreateScan(target Target, profile string) (Scan, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	id, err := newID("s")
	if err != nil {
		return Scan{}, err
	}
	scan := Scan{
		ID:        id,
		TargetID:  target.ID,
		Profile:   profile,
		Status:    "running",
		CreatedAt: time.Now().UTC(),
	}
	m.scans[scan.ID] = scan
	m.findings[scan.ID] = []Finding{}
	return scan, nil
}

func (m *Memory) CompleteScan(id string, findings []Finding) (Scan, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	scan, ok := m.scans[id]
	if !ok {
		return Scan{}, ErrNotFound
	}
	now := time.Now().UTC()
	for i := range findings {
		if findings[i].FirstSeen.IsZero() {
			findings[i].FirstSeen = now
		}
		for oldScanID, oldFindings := range m.findings {
			oldScan := m.scans[oldScanID]
			if oldScan.TargetID != scan.TargetID {
				continue
			}
			for _, oldFinding := range oldFindings {
				if oldFinding.ID == findings[i].ID && oldFinding.FirstSeen.Before(findings[i].FirstSeen) {
					findings[i].FirstSeen = oldFinding.FirstSeen
				}
			}
		}
		findings[i].ScanID = id
		findings[i].LastSeen = now
	}
	scan.Status = "complete"
	scan.CompletedAt = now
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
	if scan.Status != "running" {
		return nil
	}
	scan.Status = "failed"
	scan.CompletedAt = time.Now().UTC()
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

func (m *Memory) Findings(scanID string) ([]Finding, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	findings, ok := m.findings[scanID]
	if !ok {
		return nil, ErrNotFound
	}
	return append([]Finding(nil), findings...), nil
}
