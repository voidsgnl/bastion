# 07 Architecture

## High-level components

```
                 +-----------------------------+
                 |   Web UI (React/TypeScript) |
                 +--------------+--------------+
                                |
  CLI / API clients  ---------->|<---------- MCP clients (Claude etc.)
                                v
                 +-----------------------------+
                 |   API gateway + MCP server  |
                 |   (authn/z, RBAC, audit)    |
                 +---+---------+----------+----+
                     |         |          |
          +----------+   +-----+-----+  +-+-----------+
          | Orchestr.|   | AI service|  | Report svc  |
          | + queue  |   | (gateway, |  | (PDF/DOCX/  |
          +----+-----+   | retrieval)|  |  CSV/SARIF) |
               |         +-----+-----+  +-------------+
               v               |
   +-----------+-------+       v
   | Scan engines      |   Knowledge base
   | (per segment,     |   (NVD, KEV, EPSS,
   |  containerized)   |    advisories, OVAL)
   +-----------+-------+
               v
   Normalizer -> Verification engine -> Risk scorer -> Attack-path graph
               |
               v
   PostgreSQL (assets, findings)  |  OpenSearch (search, history)  |  Object store (evidence)
```

## Components

| Component | Responsibility |
|---|---|
| **Web UI** | Dashboard, scan view, finding pages, attack-path graph, AI panel |
| **API gateway** | REST/gRPC API, auth (SSO, tokens), RBAC, rate limiting |
| **MCP server** | Exposes tools/resources/prompts; shares auth and audit with API |
| **Orchestrator** | Turns a scan request into jobs, schedules on engines, tracks state |
| **Scan engines** | Containerized workers wrapping Nmap, OpenVAS, Nuclei, ZAP, Trivy, Prowler, etc. Deployed per network segment, paired to console with a token over mTLS |
| **Normalizer** | Converts engine output to the common [finding schema](schema/finding.schema.json) |
| **Verification engine** | Safe proof checks that upgrade confidence (Unverified -> Likely -> Confirmed) |
| **Risk scorer** | Combines CVSS, EPSS, KEV, exposure, asset role, controls |
| **Attack-path graph** | Builds graph of assets, services, credentials, and findings; computes chokepoints |
| **AI service** | Model gateway, redaction, retrieval, prompt templates, guardrails |
| **Credential vault** | Encrypted storage for scan credentials; never exposed to AI |
| **Report service** | Renders executive/technical/compliance reports |
| **Audit service** | Immutable log of user, API, MCP, and AI actions |

## Technology stack (proposed)

| Layer | Choice | Reason |
|---|---|---|
| Frontend | React, TypeScript, Tailwind, Cytoscape/D3 | Fast iteration, graph visualization |
| Backend API | Go (or Python/FastAPI) | Concurrency and static binaries (Go) vs. faster prototyping (Python) |
| Scan workers | Go or Python in containers | Easy wrapping of CLI engines |
| Primary DB | PostgreSQL | Relational assets and findings |
| Search/history | OpenSearch | Finding search, diff, trend queries |
| Queue | NATS or Redis/RabbitMQ | Job distribution to engines |
| Evidence store | S3-compatible (MinIO) | Raw responses, screenshots, reports |
| AI | Model-agnostic gateway; Ollama for local | Air-gap support |
| Packaging | Docker Compose, Helm chart, appliance image | Nexpose-like easy install |

A final backend language decision is listed in [open questions](13-open-questions.md).

## Data model (core entities)

- **Organization / User / Role**
- **Target** (scope, authorization status and proof)
- **Asset** (host, service, web app, cloud resource, repo; owner, criticality, tags)
- **Scan** (profile, engine runs, status, timestamps)
- **Finding** (see schema; links to asset, scan, evidence, remediation)
- **Evidence** (request/response, banner, config snippet, screenshot)
- **Remediation** (steps, generated scripts, affected findings)
- **AttackPath** (ordered findings, chokepoint)
- **Ticket** (external reference)
- **AuditEvent**

## Scan data flow

1. User or MCP client requests a scan of an authorized target.
2. Orchestrator validates scope, creates jobs, and sends them to engines.
3. Engines run tools and emit raw results.
4. Normalizer maps results to findings and de-duplicates.
5. Verification engine runs safe proof checks and sets confidence.
6. Risk scorer enriches (KEV, EPSS) and computes risk.
7. Attack-path graph is updated.
8. AI service generates explanations and remediation on demand (or in batch).
9. UI streams updates over WebSocket/SSE.

## Deployment topologies

| Topology | Description |
|---|---|
| **Single node** | Console + engine in one Docker Compose stack (small teams, labs) |
| **Distributed** | One console, multiple engines per segment/site/cloud |
| **Air-gapped** | Offline feed bundles, local model, no outbound traffic |
| **Hosted (later)** | Managed console with customer-side engines |

## Non-functional requirements

- Scan 10k hosts in a distributed setup with horizontal engine scaling
- Idempotent, resumable scan jobs
- Findings searchable and filterable in under 1 second at 1M findings
- Backward-compatible schema migrations
- Full API/MCP parity with UI
