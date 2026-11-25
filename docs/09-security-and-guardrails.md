# 09 Security and Guardrails

A vulnerability scanner holds powerful credentials and network reach, and can be misused. Bastion's guardrails are part of the product, not an afterthought.

## 1. Authorization and scope

- **No scan without an authorized target.** Targets need proof of control or a recorded attestation.
- **Scope-locking:** engines refuse packets to addresses outside authorized scope, enforced at the orchestrator and again at the engine.
- **Redirect and DNS safety:** web crawls do not follow redirects out of scope; DNS rebinding is checked.
- **Default deny** for public IP ranges not proven owned.

## 2. Safe-by-default scanning

- Non-destructive checks by default (intrusiveness levels 0-1)
- Intrusive checks (level 2) require explicit human approval per scan
- Destructive checks (level 3) do not exist
- Per-target rate limits, concurrency caps, and blackout windows
- A "safe mode" profile for production systems
- The initial Nmap lab runner is opt-in and limited to `127.0.0.1/32`, TCP port
  18080, and fixed TCP-connect probe arguments; it does not run vulnerability
  scripts or accept caller-controlled command options.
- Do not load arbitrary Nmap command lines from text files. Richer scan modes
  must be implemented as reviewed structured profiles with explicit bounds.

## 3. Platform security

| Area | Control |
|---|---|
| Credentials | Encrypted vault (envelope encryption, KMS/HSM optional); never sent to AI; masked in UI |
| Transport | TLS everywhere; mTLS between console and engines |
| Access | RBAC, SSO (OIDC/SAML), MFA, scoped API tokens |
| Audit | Immutable log of user, API, MCP, and AI actions |
| Supply chain | Signed releases, SBOM published, pinned dependencies, signed feed updates |
| Hardening | Containers run non-root, read-only filesystems where possible, minimal images |
| Evidence | Stored encrypted; retention policy and secret redaction in captured evidence |

## 4. AI and MCP guardrails

- AI cannot create or modify findings
- AI-initiated scans are scope-locked, rate-limited, and require approval for anything above safe level
- MCP tokens are scoped (read-only, scan-allowed) and revocable
- Prompt-injection defense: scan output and target content are treated as untrusted data, never as instructions; tool calls triggered from untrusted content require confirmation
- Redaction before any cloud model call

## 5. Responsible disclosure and misuse

- Terms of use prohibit scanning without authorization
- Abuse monitoring for scan patterns against unauthorized public ranges (hosted offerings)
- Vulnerabilities in Bastion itself: a `SECURITY.md` with private reporting instructions will be added

## 6. Threat model (initial)

| Threat | Mitigation |
|---|---|
| Attacker steals scan credentials from the vault | Envelope encryption, least-privilege scan accounts, rotation guidance |
| Compromised engine pivots into network | Engine isolation, minimal privileges, short-lived tokens, egress policy |
| Malicious target content attacks the scanner (parser bugs, prompt injection) | Sandboxed engines, size and time limits, untrusted-data handling in AI |
| Insider uses Bastion for unauthorized recon | Authorization proofs, RBAC, audit logs, alerts on out-of-pattern scans |
| Tampered feed or check pack | Signed updates, checksum verification |
| Findings leak sensitive data | Encryption at rest, role-based field redaction, report access control |
