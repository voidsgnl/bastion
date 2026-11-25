package store

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"
)

func TestPostgresPersistsScansAcrossReopen(t *testing.T) {
	databaseURL := os.Getenv("BASTION_TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("set BASTION_TEST_DATABASE_URL to run PostgreSQL integration test")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	cleanupDB, err := sql.Open("pgx", databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanupDB.Close()

	repository, err := NewPostgres(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	target, err := repository.CreateTarget(Target{
		Name: "integration-test", Type: "domain", Value: "integration.invalid",
	})
	if err != nil {
		_ = repository.Close()
		t.Fatal(err)
	}
	defer func() {
		_, _ = cleanupDB.ExecContext(context.Background(), "DELETE FROM targets WHERE id = $1", target.ID)
		_ = repository.Close()
	}()
	if _, err := repository.AttestTarget(target.ID, "Local PostgreSQL integration test."); err != nil {
		t.Fatal(err)
	}

	first, err := repository.CreateScan(target, "safe")
	if err != nil {
		t.Fatal(err)
	}
	initialTime := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	if _, err := repository.CompleteScan(first.ID, []Finding{{
		ID: "stable-finding", FirstSeen: initialTime, Evidence: []Evidence{}, Risk: Risk{},
	}}); err != nil {
		t.Fatal(err)
	}
	if err := repository.Close(); err != nil {
		t.Fatal(err)
	}

	reopened, err := NewPostgres(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	repository = reopened
	second, err := repository.CreateScan(target, "safe")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := repository.CompleteScan(second.ID, []Finding{{
		ID: "stable-finding", Evidence: []Evidence{}, Risk: Risk{},
	}}); err != nil {
		t.Fatal(err)
	}
	findings, err := repository.Findings(second.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 1 || !findings[0].FirstSeen.Equal(initialTime) {
		t.Fatalf("persisted first_seen not preserved: %#v", findings)
	}
}
