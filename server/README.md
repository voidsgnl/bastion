# Bastion server

The Phase 1 console backend is a Go implementation that currently provides an in-memory
vertical slice:

- health and target APIs
- explicit recorded authorization attestation before fixture scans
- canonical CIDR target validation and IP-in-scope checks
- an injectable scan-runner interface with a clearly labeled synthetic fixture runner
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

The default `fixture` engine is synthetic and performs no network activity. To use
the restricted Nmap lab runner, follow [the local lab instructions](../testing/lab/README.md)
and start the server with:

```bash
go run ./cmd/bastion-server -ui-dir ../ui -engine nmap-lab
```

The Nmap runner only accepts `127.0.0.1/32` with the `safe` profile and probes
TCP port 18080 using fixed arguments. Nmap must be installed locally. It does not
perform vulnerability checks.

## Test

```bash
cd server
go test ./...
```
