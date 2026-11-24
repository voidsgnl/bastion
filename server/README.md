# Bastion server

The Phase 1 console backend is a Go implementation that currently provides an in-memory
vertical slice:

- health and target APIs
- explicit recorded authorization attestation before fixture scans
- canonical CIDR target validation and IP-in-scope checks
- scan status and normalized finding APIs
- static UI serving

The attestation is a prototype record only; the server does not authenticate the
caller or independently verify ownership. Real scanning must not be enabled until
identity, authorization proof, and engine-side scope enforcement are implemented.

## Run

```bash
cd server
go run ./cmd/bastion-server -ui-dir ../ui
```

Open <http://localhost:8080>.

## Test

```bash
cd server
go test ./...
```
