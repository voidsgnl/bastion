package migrations

import "embed"

// Files contains the SQL migrations applied by the PostgreSQL store at startup.
//
//go:embed *.sql
var Files embed.FS
