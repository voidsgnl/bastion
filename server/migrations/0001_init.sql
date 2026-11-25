CREATE TABLE targets (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    value TEXT NOT NULL,
    authorized BOOLEAN NOT NULL DEFAULT FALSE,
    authorization_attestation TEXT,
    authorized_at TIMESTAMPTZ
);

CREATE TABLE scans (
    id TEXT PRIMARY KEY,
    target_id TEXT NOT NULL REFERENCES targets(id) ON DELETE CASCADE,
    profile TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('running', 'complete', 'failed')),
    created_at TIMESTAMPTZ NOT NULL,
    completed_at TIMESTAMPTZ
);

CREATE INDEX scans_target_created_idx ON scans (target_id, created_at DESC);

CREATE TABLE findings (
    scan_id TEXT NOT NULL REFERENCES scans(id) ON DELETE CASCADE,
    finding_id TEXT NOT NULL,
    target_id TEXT NOT NULL REFERENCES targets(id) ON DELETE CASCADE,
    asset_id TEXT NOT NULL,
    source_engine TEXT NOT NULL DEFAULT '',
    title TEXT NOT NULL,
    category TEXT NOT NULL,
    severity TEXT NOT NULL,
    confidence TEXT NOT NULL,
    status TEXT NOT NULL,
    evidence JSONB NOT NULL,
    risk JSONB NOT NULL,
    first_seen TIMESTAMPTZ NOT NULL,
    last_seen TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (scan_id, finding_id)
);

CREATE INDEX findings_identity_history_idx ON findings (target_id, finding_id, first_seen);
