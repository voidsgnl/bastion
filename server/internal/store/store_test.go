package store

import (
	"errors"
	"testing"
	"time"
)

func TestCompleteScanPreservesFindingIdentityAndFirstSeen(t *testing.T) {
	memory := NewMemory()
	target, err := memory.CreateTarget(Target{Name: "lab", Type: "domain", Value: "lab.example"})
	if err != nil {
		t.Fatal(err)
	}
	firstSeen := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)

	firstScan, err := memory.CreateScan(target, "safe")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := memory.CompleteScan(firstScan.ID, []Finding{{
		ID: "stable-finding", FirstSeen: firstSeen,
	}}); err != nil {
		t.Fatal(err)
	}

	secondScan, err := memory.CreateScan(target, "safe")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := memory.CompleteScan(secondScan.ID, []Finding{{
		ID: "stable-finding", FirstSeen: time.Now().UTC(),
	}}); err != nil {
		t.Fatal(err)
	}
	secondFindings, err := memory.Findings(secondScan.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(secondFindings) != 1 {
		t.Fatalf("findings = %d, want 1", len(secondFindings))
	}
	if secondFindings[0].ID != "stable-finding" {
		t.Fatalf("finding ID = %q", secondFindings[0].ID)
	}
	if !secondFindings[0].FirstSeen.Equal(firstSeen) {
		t.Fatalf("first_seen = %s, want %s", secondFindings[0].FirstSeen, firstSeen)
	}
	if secondFindings[0].LastSeen.Before(firstSeen) {
		t.Fatalf("last_seen %s precedes first_seen %s", secondFindings[0].LastSeen, firstSeen)
	}
}

func TestFindingsRequiresExistingScan(t *testing.T) {
	if _, err := NewMemory().Findings("missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("error = %v, want ErrNotFound", err)
	}
}

func TestFailScanDoesNotOverwriteCompletedScan(t *testing.T) {
	memory := NewMemory()
	target, err := memory.CreateTarget(Target{Name: "lab", Type: "domain", Value: "lab.example"})
	if err != nil {
		t.Fatal(err)
	}
	scan, err := memory.CreateScan(target, "safe")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := memory.CompleteScan(scan.ID, nil); err != nil {
		t.Fatal(err)
	}
	if err := memory.FailScan(scan.ID); err != nil {
		t.Fatal(err)
	}
	stored, err := memory.GetScan(scan.ID)
	if err != nil {
		t.Fatal(err)
	}
	if stored.Status != "complete" {
		t.Fatalf("scan status = %q, want complete", stored.Status)
	}
}

func TestCompareFindingsClassifiesStableIDs(t *testing.T) {
	baseline := []Finding{
		{ID: "still-open"},
		{ID: "now-fixed"},
	}
	current := []Finding{
		{ID: "still-open"},
		{ID: "new-finding"},
	}

	comparison := CompareFindings("scan-old", "scan-new", baseline, current)
	if comparison.BaselineScanID != "scan-old" || comparison.CurrentScanID != "scan-new" {
		t.Fatalf("comparison IDs = %s/%s", comparison.BaselineScanID, comparison.CurrentScanID)
	}
	if len(comparison.New) != 1 || comparison.New[0].ID != "new-finding" {
		t.Fatalf("new findings = %#v", comparison.New)
	}
	if len(comparison.Fixed) != 1 || comparison.Fixed[0].ID != "now-fixed" {
		t.Fatalf("fixed findings = %#v", comparison.Fixed)
	}
	if len(comparison.Unchanged) != 1 || comparison.Unchanged[0].ID != "still-open" {
		t.Fatalf("unchanged findings = %#v", comparison.Unchanged)
	}
}
