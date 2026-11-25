package store

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
)

var (
	ErrNotFound       = errors.New("not found")
	ErrDifferentScope = errors.New("scans belong to different targets")
)

type Repository interface {
	CreateTarget(Target) (Target, error)
	ListTargets() ([]Target, error)
	GetTarget(string) (Target, error)
	AttestTarget(string, string) (Target, error)
	CreateScan(Target, string) (Scan, error)
	CompleteScan(string, []Finding) (Scan, error)
	FailScan(string) error
	GetScan(string) (Scan, error)
	Findings(string) ([]Finding, error)
}

type Comparison struct {
	BaselineScanID string    `json:"baseline_scan_id"`
	CurrentScanID  string    `json:"current_scan_id"`
	New            []Finding `json:"new"`
	Fixed          []Finding `json:"fixed"`
	Unchanged      []Finding `json:"unchanged"`
}

func CompareFindings(baselineScanID, currentScanID string, baseline, current []Finding) Comparison {
	baselineByID := make(map[string]Finding, len(baseline))
	currentByID := make(map[string]Finding, len(current))
	for _, finding := range baseline {
		baselineByID[finding.ID] = finding
	}
	for _, finding := range current {
		currentByID[finding.ID] = finding
	}

	result := Comparison{
		BaselineScanID: baselineScanID,
		CurrentScanID:  currentScanID,
		New:            make([]Finding, 0),
		Fixed:          make([]Finding, 0),
		Unchanged:      make([]Finding, 0),
	}
	for id, finding := range currentByID {
		if _, exists := baselineByID[id]; exists {
			result.Unchanged = append(result.Unchanged, finding)
		} else {
			result.New = append(result.New, finding)
		}
	}
	for id, finding := range baselineByID {
		if _, exists := currentByID[id]; !exists {
			result.Fixed = append(result.Fixed, finding)
		}
	}
	sortFindings(result.New)
	sortFindings(result.Fixed)
	sortFindings(result.Unchanged)
	return result
}

func newID(prefix string) (string, error) {
	var random [16]byte
	if _, err := rand.Read(random[:]); err != nil {
		return "", fmt.Errorf("generate %s id: %w", prefix, err)
	}
	return prefix + "_" + hex.EncodeToString(random[:]), nil
}

func sortFindings(findings []Finding) {
	sort.Slice(findings, func(i, j int) bool {
		return findings[i].ID < findings[j].ID
	})
}
