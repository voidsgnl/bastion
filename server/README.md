# Bastion server

The Phase 1 console backend is a Go implementation that currently provides an in-memory
vertical slice:

- health and target APIs
- explicit recorded authorization attestation before fixture scans
- canonical CIDR target validation and IP-in-scope checks
- an injectable scan-runner interface with a clearly labeled synthetic fixture runner
- scan status and normalized finding APIs
- static UI serving

The backend now requires a bootstrapped operator login for API access. This is a
single-operator prototype session, not multi-user identity or RBAC. Target
attestation remains a separate scope control: login identifies the operator but
does not prove ownership of every target. Attestation is still a recorded
statement, not independent ownership verification. The Nmap runner independently
restricts itself to the local lab target.

## Run

```bash
cd server
export BASTION_ADMIN_PASSWORD='use-a-unique-password-of-at-least-12-characters'
go run ./cmd/bastion-server -ui-dir ../ui
```

The default listener is `127.0.0.1:8080`. Set `BASTION_ADMIN_USER` to change the
default username (`admin`). The server refuses to start without the password.
For non-loopback listening, provide `-tls-cert` and `-tls-key`; plaintext
non-loopback listeners are rejected.

Log in to the API to receive an HttpOnly, SameSite=Strict session cookie:

```bash
curl -i -c bastion-cookies.txt -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"YOUR_PASSWORD"}' \
  http://127.0.0.1:8080/api/v1/auth/login
```

Use that cookie for API requests. `POST /api/v1/auth/logout` revokes the session;
`GET /api/v1/auth/me` returns the current operator. Sessions expire after 12
hours and are cleared when the server restarts. `/health` remains public.

The default `fixture` engine is synthetic and performs no network activity. To use
the restricted Nmap lab runner, follow [the local lab instructions](../testing/lab/README.md)
and start the server with:

```bash
go run ./cmd/bastion-server -ui-dir ../ui -engine nmap-lab
```

The Nmap runner only accepts `127.0.0.1/32` with the `safe` profile and probes
TCP port 18080 using fixed arguments. Nmap must be installed locally. It does not
perform vulnerability checks.

## Persistent scan history

Copy `.env.example` to `.env`, replace the local development password, and export
the matching `BASTION_DATABASE_URL` in the shell where the server will run. Start
PostgreSQL with `make db-up`, then run:

```bash
make run-postgres
```

The server applies embedded schema migrations at startup. The database port is
published only on loopback. `make db-down` stops PostgreSQL without deleting its
volume. The default `memory` store remains useful for tests and disposable demos;
it does not survive a restart.

## Test

```bash
cd server
go test ./...
```
