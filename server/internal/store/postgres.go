package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"sort"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/voidsgnl/bastion/server/migrations"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres(ctx context.Context, databaseURL string) (*Postgres, error) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("open PostgreSQL connection: %w", err)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("connect to PostgreSQL: %w", err)
	}
	repository := &Postgres{db: db}
	if err := repository.migrate(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("migrate PostgreSQL store: %w", err)
	}
	return repository, nil
}

func (p *Postgres) Close() error {
	return p.db.Close()
}

func (p *Postgres) migrate(ctx context.Context) error {
	if _, err := p.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS bastion_schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
		)`); err != nil {
		return err
	}
	files, err := fs.ReadDir(migrations.Files, ".")
	if err != nil {
		return err
	}
	sort.Slice(files, func(i, j int) bool { return files[i].Name() < files[j].Name() })
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}
		var applied bool
		if err := p.db.QueryRowContext(ctx,
			"SELECT EXISTS (SELECT 1 FROM bastion_schema_migrations WHERE version = $1)",
			file.Name(),
		).Scan(&applied); err != nil {
			return err
		}
		if applied {
			continue
		}
		content, err := fs.ReadFile(migrations.Files, file.Name())
		if err != nil {
			return err
		}
		tx, err := p.db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		for _, statement := range strings.Split(string(content), ";") {
			statement = strings.TrimSpace(statement)
			if statement == "" {
				continue
			}
			if _, err := tx.ExecContext(ctx, statement); err != nil {
				_ = tx.Rollback()
				return fmt.Errorf("%s: %w", file.Name(), err)
			}
		}
		if _, err := tx.ExecContext(ctx,
			"INSERT INTO bastion_schema_migrations (version) VALUES ($1)",
			file.Name(),
		); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("record migration %s: %w", file.Name(), err)
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit migration %s: %w", file.Name(), err)
		}
	}
	return nil
}

func (p *Postgres) CreateTarget(target Target) (Target, error) {
	id, err := newID("t")
	if err != nil {
		return Target{}, err
	}
	target.ID = id
	_, err = p.db.Exec(`
		INSERT INTO targets (id, name, type, value, authorized)
		VALUES ($1, $2, $3, $4, FALSE)`,
		target.ID, target.Name, target.Type, target.Value,
	)
	if err != nil {
		return Target{}, fmt.Errorf("create target: %w", err)
	}
	return target, nil
}

func (p *Postgres) ListTargets() ([]Target, error) {
	rows, err := p.db.Query(`
		SELECT id, name, type, value, authorized, authorization_attestation, authorized_at
		FROM targets ORDER BY name, id`)
	if err != nil {
		return nil, fmt.Errorf("list targets: %w", err)
	}
	defer rows.Close()
	targets := make([]Target, 0)
	for rows.Next() {
		target, err := scanTarget(rows)
		if err != nil {
			return nil, fmt.Errorf("read target: %w", err)
		}
		targets = append(targets, target)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate targets: %w", err)
	}
	return targets, nil
}

func (p *Postgres) GetTarget(id string) (Target, error) {
	target, err := scanTarget(p.db.QueryRow(`
		SELECT id, name, type, value, authorized, authorization_attestation, authorized_at
		FROM targets WHERE id = $1`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return Target{}, ErrNotFound
	}
	if err != nil {
		return Target{}, fmt.Errorf("get target: %w", err)
	}
	return target, nil
}

type rowScanner interface {
	Scan(...any) error
}

func scanTarget(row rowScanner) (Target, error) {
	var target Target
	var attestation sql.NullString
	var authorizedAt sql.NullTime
	err := row.Scan(
		&target.ID, &target.Name, &target.Type, &target.Value, &target.Authorized,
		&attestation, &authorizedAt,
	)
	if err != nil {
		return Target{}, err
	}
	if attestation.Valid {
		target.AuthorizationAttestation = attestation.String
	}
	if authorizedAt.Valid {
		target.AuthorizedAt = authorizedAt.Time
	}
	return target, nil
}

func (p *Postgres) AttestTarget(id, statement string) (Target, error) {
	result, err := p.db.Exec(`
		UPDATE targets
		SET authorized = TRUE, authorization_attestation = $2, authorized_at = now()
		WHERE id = $1`, id, statement)
	if err != nil {
		return Target{}, fmt.Errorf("record target attestation: %w", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return Target{}, fmt.Errorf("check target attestation: %w", err)
	}
	if count == 0 {
		return Target{}, ErrNotFound
	}
	return p.GetTarget(id)
}

func (p *Postgres) CreateScan(target Target, profile string) (Scan, error) {
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
	_, err = p.db.Exec(`
		INSERT INTO scans (id, target_id, profile, status, created_at)
		VALUES ($1, $2, $3, $4, $5)`,
		scan.ID, scan.TargetID, scan.Profile, scan.Status, scan.CreatedAt)
	if err != nil {
		return Scan{}, fmt.Errorf("create scan: %w", err)
	}
	return scan, nil
}

func (p *Postgres) CompleteScan(id string, findings []Finding) (Scan, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return Scan{}, fmt.Errorf("begin scan completion: %w", err)
	}
	var targetID string
	if err := tx.QueryRow("SELECT target_id FROM scans WHERE id = $1 FOR UPDATE", id).Scan(&targetID); err != nil {
		_ = tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			return Scan{}, ErrNotFound
		}
		return Scan{}, fmt.Errorf("load scan for completion: %w", err)
	}

	now := time.Now().UTC()
	for i := range findings {
		finding := &findings[i]
		if finding.ID == "" {
			_ = tx.Rollback()
			return Scan{}, fmt.Errorf("finding ID is required")
		}
		finding.ScanID = id
		if finding.FirstSeen.IsZero() {
			finding.FirstSeen = now
		}
		var previousFirstSeen sql.NullTime
		err := tx.QueryRow(`
			SELECT MIN(first_seen) FROM findings
			WHERE target_id = $1 AND finding_id = $2`,
			targetID, finding.ID).Scan(&previousFirstSeen)
		if err != nil {
			_ = tx.Rollback()
			return Scan{}, fmt.Errorf("load finding history: %w", err)
		}
		if previousFirstSeen.Valid && previousFirstSeen.Time.Before(finding.FirstSeen) {
			finding.FirstSeen = previousFirstSeen.Time
		}
		finding.LastSeen = now
		evidence, err := json.Marshal(finding.Evidence)
		if err != nil {
			_ = tx.Rollback()
			return Scan{}, fmt.Errorf("encode finding evidence: %w", err)
		}
		risk, err := json.Marshal(finding.Risk)
		if err != nil {
			_ = tx.Rollback()
			return Scan{}, fmt.Errorf("encode finding risk: %w", err)
		}
		_, err = tx.Exec(`
			INSERT INTO findings (
				scan_id, finding_id, target_id, asset_id, source_engine, title, category,
				severity, confidence, status, evidence, risk, first_seen, last_seen
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
			ON CONFLICT (scan_id, finding_id) DO UPDATE SET
				asset_id = EXCLUDED.asset_id,
				source_engine = EXCLUDED.source_engine,
				title = EXCLUDED.title,
				category = EXCLUDED.category,
				severity = EXCLUDED.severity,
				confidence = EXCLUDED.confidence,
				status = EXCLUDED.status,
				evidence = EXCLUDED.evidence,
				risk = EXCLUDED.risk,
				first_seen = EXCLUDED.first_seen,
				last_seen = EXCLUDED.last_seen`,
			id, finding.ID, targetID, finding.AssetID, finding.SourceEngine, finding.Title,
			finding.Category, finding.Severity, finding.Confidence, finding.Status,
			string(evidence), string(risk), finding.FirstSeen, finding.LastSeen)
		if err != nil {
			_ = tx.Rollback()
			return Scan{}, fmt.Errorf("persist finding %s: %w", finding.ID, err)
		}
	}
	if _, err := tx.Exec(
		"UPDATE scans SET status = 'complete', completed_at = $2 WHERE id = $1",
		id, now,
	); err != nil {
		_ = tx.Rollback()
		return Scan{}, fmt.Errorf("mark scan complete: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return Scan{}, fmt.Errorf("commit scan completion: %w", err)
	}
	return p.GetScan(id)
}

func (p *Postgres) FailScan(id string) error {
	result, err := p.db.Exec(
		"UPDATE scans SET status = 'failed', completed_at = now() WHERE id = $1 AND status = 'running'",
		id,
	)
	if err != nil {
		return fmt.Errorf("mark scan failed: %w", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check failed scan update: %w", err)
	}
	if count == 0 {
		var exists bool
		if err := p.db.QueryRow("SELECT EXISTS (SELECT 1 FROM scans WHERE id = $1)", id).Scan(&exists); err != nil {
			return fmt.Errorf("check scan before marking failed: %w", err)
		}
		if !exists {
			return ErrNotFound
		}
	}
	return nil
}

func (p *Postgres) GetScan(id string) (Scan, error) {
	var scan Scan
	var completedAt sql.NullTime
	err := p.db.QueryRow(`
		SELECT id, target_id, profile, status, created_at, completed_at
		FROM scans WHERE id = $1`, id,
	).Scan(&scan.ID, &scan.TargetID, &scan.Profile, &scan.Status, &scan.CreatedAt, &completedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return Scan{}, ErrNotFound
	}
	if err != nil {
		return Scan{}, fmt.Errorf("get scan: %w", err)
	}
	if completedAt.Valid {
		scan.CompletedAt = completedAt.Time
	}
	return scan, nil
}

func (p *Postgres) Findings(scanID string) ([]Finding, error) {
	var exists bool
	if err := p.db.QueryRow("SELECT EXISTS (SELECT 1 FROM scans WHERE id = $1)", scanID).Scan(&exists); err != nil {
		return nil, fmt.Errorf("check scan before listing findings: %w", err)
	}
	if !exists {
		return nil, ErrNotFound
	}
	rows, err := p.db.Query(`
		SELECT finding_id, scan_id, asset_id, source_engine, title, category, severity,
			confidence, status, evidence, risk, first_seen, last_seen
		FROM findings WHERE scan_id = $1 ORDER BY finding_id`, scanID)
	if err != nil {
		return nil, fmt.Errorf("list findings: %w", err)
	}
	defer rows.Close()

	findings := make([]Finding, 0)
	for rows.Next() {
		var finding Finding
		var evidence, risk []byte
		if err := rows.Scan(
			&finding.ID, &finding.ScanID, &finding.AssetID, &finding.SourceEngine,
			&finding.Title, &finding.Category, &finding.Severity, &finding.Confidence,
			&finding.Status, &evidence, &risk, &finding.FirstSeen, &finding.LastSeen,
		); err != nil {
			return nil, fmt.Errorf("read finding: %w", err)
		}
		if err := json.Unmarshal(evidence, &finding.Evidence); err != nil {
			return nil, fmt.Errorf("decode finding evidence: %w", err)
		}
		if err := json.Unmarshal(risk, &finding.Risk); err != nil {
			return nil, fmt.Errorf("decode finding risk: %w", err)
		}
		findings = append(findings, finding)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate findings: %w", err)
	}
	return findings, nil
}

var _ Repository = (*Postgres)(nil)
