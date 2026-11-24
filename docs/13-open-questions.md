# 13 Open Questions

Decisions to make in Phase 0.

## Product

1. **Name:** keep "Bastion" (check trademark conflicts) or rename?
2. **Audience priority:** SMB/solo admins, enterprise blue teams, or MSSPs first?
3. **Hosted offering:** self-hosted only, or also a managed console?

## Licensing and business

4. **License:** Apache-2.0 / MIT core, AGPL core, or source-available? Which features are paid (compliance packs, multi-site engines, SSO, support)?
5. **Engine license compatibility:** which wrapped tools (for example GPL components) must be isolated as separate processes, and what does that mean for distribution?

## Technical

6. **Backend language:** Go (static binaries, concurrency) or Python (speed, scanning-ecosystem familiarity)?
7. **Queue:** NATS vs Redis vs RabbitMQ.
8. **Verification engine design:** how to build safe proof checks per vulnerability class without becoming an exploit framework.
9. **Feed strategy:** which sources can be mirrored and redistributed legally; update cadence.
10. **Local model baseline:** which open models are good enough for explanation and fix generation on commodity hardware?
11. **Import formats:** priority order for ingesting external scanner output (Nessus, Nmap XML, Trivy, SARIF).

## Security and compliance

12. **Authorization proof for internal ranges:** is attestation sufficient, or require network-level proof (for example agent callback)?
13. **Data residency** for evidence and AI prompts in cloud mode.
14. **Certifications** to target (SOC 2, ISO 27001) if offering a hosted service.

## Process

15. **Benchmark lab:** assemble a vulnerable test environment (for example Metasploitable, DVWA, intentionally vulnerable cloud labs) to measure false-positive and detection rates.
16. **Contribution model:** open governance and community check contributions.
