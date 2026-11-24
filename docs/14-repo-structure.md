# 14 Repository Structure

Planned layout for the Bastion monorepo. It follows the components in [07 Architecture](07-architecture.md) and the phases in [08 Roadmap](08-roadmap.md). Folders are added when their phase starts, not all at once.

```
bastion/
├── README.md                      (exists)
├── CONTRIBUTING.md                (exists)
├── LICENSE                        (blocked on open question #4)
├── SECURITY.md                    how to report vulnerabilities in Bastion itself
├── CODE_OF_CONDUCT.md
├── CHANGELOG.md
├── Makefile                       build, test, lint, run-dev
├── .gitignore  .editorconfig  .env.example
│
├── .github/
│   ├── workflows/                 ci.yml, release.yml, codeql.yml, docs-lint.yml
│   ├── ISSUE_TEMPLATE/            bug_report.yml, feature_request.yml, config.yml
│   ├── PULL_REQUEST_TEMPLATE.md
│   ├── CODEOWNERS
│   └── dependabot.yml
│
├── docs/                          (exists: 01-vision ... 13-open-questions)
│   ├── schema/finding.schema.json (exists; single source of truth)
│   ├── 14-repo-structure.md       this file
│   └── adr/                       one short file per decision
│
├── api/                           contracts shared by every component
│   ├── openapi.yaml               REST API
│   ├── proto/engine/v1/engine.proto   console <-> engine (gRPC, mTLS)
│   └── mcp/tools.json             MCP tool definitions from docs/06
│
├── server/                        console backend (one deployable)
│   ├── cmd/bastion-server/
│   ├── internal/
│   │   ├── gateway/  auth/  rbac/  mcp/
│   │   ├── targets/  scans/  orchestrator/
│   │   ├── normalizer/  verifier/  risk/  attackpath/
│   │   ├── findings/  remediation/
│   │   ├── ai/  vault/  reports/  audit/
│   │   ├── integrations/          jira/, servicenow/, slack/, github/
│   │   ├── feeds/                 NVD, KEV, EPSS loaders
│   │   └── store/                 postgres/, opensearch/, objects/
│   ├── migrations/
│   └── tests/
│
├── engine/                        scan worker, deployed per network segment
│   ├── cmd/bastion-engine/
│   ├── internal/                  pairing/, runner/, parsers/
│   ├── adapters/                  nmap/, nuclei/, zap/, openvas/, trivy/, prowler/
│   └── Dockerfile
│
├── cli/cmd/bastion/               target, scan, findings, fix, retest, report
│
├── web/                           React + TypeScript app (Phase 1)
│   └── src/                       app/, pages/, components/, graph/, api/, styles/
│
├── ui/                            (exists) static prototype: index.html, results.html, finding.html
│
├── ai/                            prompts/, guardrails/, evals/
├── feeds/                         offline feed bundle builder
├── compliance/                    cis/, pci-dss/, nist-800-53/, iso-27001/  (Phase 3)
├── config/                        examples/bastion.example.yaml, profiles/*.yaml
├── deploy/                        compose/, helm/bastion/, appliance/, airgap/
├── testing/                       lab/, fixtures/, e2e/
└── scripts/                       dev-setup.sh, gen-api.sh, lint.sh, release.sh
```

## When each folder is added

| Phase | Add |
|---|---|
| 0 (now) | `LICENSE`, `.github/`, `docs/adr/`, `api/`, `ui/` |
| 1 (MVP) | `server/` (targets, scans, normalizer, risk, ai, mcp, audit, auth), `engine/` (nmap, nuclei, zap), `web/`, `deploy/compose/`, `config/`, `cli/`, `feeds/` |
| 2 (v1) | `verifier/`, `vault/`, `reports/`, `integrations/`, `testing/lab/`, SSO |
| 3 (v2) | `attackpath/`, `compliance/`, `deploy/helm/`, more adapters (trivy, prowler) |
| 4 (v3) | `deploy/appliance/`, `deploy/airgap/`, plugin SDK |

## Decisions this layout leaves open

- **Backend language** ([ADR 0001](adr/0001-backend-language.md)). The insides of `server/` and `engine/` are written Go-style. The top-level folders are the same for Python.
- **`ui/` vs `web/`.** Keep `ui/` as the prototype and build the real app in `web/`. Delete `ui/` once `web/` matches it.
- **Modular monolith first.** Start the services as modules under `server/internal/` in one binary, to keep install at one command. Only `engine/` is a separate deployable from day one.
- **Monorepo.** One repo, because the API contract, UI, engine and docs change together.
