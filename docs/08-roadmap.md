# 08 Roadmap

Timelines are estimates for a small team and will change.

## Phase 0: Foundations (weeks 0-4)

- [ ] Finalize name, license, and repo structure
- [ ] Evaluate and choose scan engines (Nmap, Greenbone, Nuclei, ZAP, Trivy)
- [ ] Define finding schema and API contract
- [ ] Decide backend language and deployment baseline
- [ ] UI prototype of dashboard and finding page

## Phase 1: MVP (months 1-3)

Goal: a usable on-demand scanner with AI explanations.

- [ ] Target management with authorization (DNS TXT / file / attestation)
- [ ] Network scan profile (Nmap + Nuclei, Greenbone integration)
- [ ] Web scan profile (ZAP + Nuclei)
- [ ] Normalizer and unified finding model
- [ ] CVE enrichment: NVD, CISA KEV, EPSS
- [ ] Basic risk scoring
- [ ] Live scan view and findings list
- [ ] AI explanations and fix text (local or cloud model)
- [ ] MCP server v0: `list_targets`, `start_scan`, `get_scan_status`, `get_findings`, `explain_finding`, `retest`
- [ ] Docker Compose install and first-run wizard
- [ ] Audit log and basic RBAC (admin / analyst / viewer)

**Exit criteria:** a new user installs and gets a useful, explained scan result in under 10 minutes.

## Phase 2: v1 (months 4-6)

Goal: trustworthy and workflow-ready.

- [ ] Verification engine (Confirmed / Likely / Unverified)
- [ ] Authenticated scans (SSH, WinRM, SNMP) and credential vault
- [ ] Scan diffing and targeted re-tests
- [ ] Fix-centric grouping ("action plan" view)
- [ ] Reports: executive, technical (PDF, DOCX, CSV, SARIF)
- [ ] Scheduled scans
- [ ] Ticketing: Jira and ServiceNow
- [ ] SSO (OIDC/SAML) and full RBAC
- [ ] MCP: `get_remediation`, `compare_scans`, `export_report`, scoped tokens, approval gates

**Exit criteria:** false-positive rate on Confirmed findings under 5% on a benchmark lab; teams can run a full find-fix-retest cycle.

## Phase 3: v2 (months 7-10)

Goal: breadth and intelligence.

- [ ] Attack-path graph and chokepoint analysis
- [ ] Cloud scanning (AWS, Azure, GCP)
- [ ] Container, IaC, and code scanning; SBOM tracking
- [ ] Compliance packs (CIS, PCI-DSS, NIST 800-53, ISO 27001)
- [ ] Generated remediation scripts: Ansible, PowerShell, Terraform
- [ ] Distributed engines with mTLS pairing; Helm chart
- [ ] Bastion as MCP client (CMDB, Slack, GitHub)
- [ ] External attack surface discovery

## Phase 4: v3 (months 11+)

Goal: autonomy with guardrails and ecosystem.

- [ ] AI-planned scans with human approval
- [ ] PR-based remediation for code and IaC
- [ ] Air-gapped appliance with local LLM bundle
- [ ] Optional lightweight agent for roaming devices
- [ ] Plugin SDK for community checks
- [ ] Multi-tenant mode for MSSPs
- [ ] Marketplace for check packs and compliance content

## Milestone summary

| Milestone | Target | Deliverable |
|---|---|---|
| M0 | Week 4 | Specs, schema, UI prototype |
| M1 | Month 3 | MVP (scan, explain, MCP v0) |
| M2 | Month 6 | v1 (verified findings, auth scans, reports, tickets) |
| M3 | Month 10 | v2 (attack paths, cloud, compliance) |
| M4 | Month 12+ | v3 (agentic, appliance, plugins) |

## Risks and mitigations

| Risk | Mitigation |
|---|---|
| Low check coverage vs incumbents | Build on Greenbone/Nuclei/Trivy; accept community content; focus on accuracy and fix quality |
| False positives undermine trust | Verification engine; conservative labeling; benchmark lab |
| Scanner misuse | Authorization proofs, scope-locking, audit logs, no destructive checks |
| AI hallucination | Strict grounding; AI cannot create findings; evidence citations |
| Licensing conflicts of wrapped engines | License review in Phase 0; isolate GPL components as separate processes |
| Scope creep | Phase gates with explicit exit criteria |
