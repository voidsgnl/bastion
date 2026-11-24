# 01 Vision

## Summary

Bastion is an on-demand vulnerability and exposure scanner for blue teams. Point it at a network, web service, cloud account, or code repository; it finds problems, proves they are real, ranks them by actual risk, and explains exactly how to fix them.

## Goals

1. **Scan-to-fix in minutes.** From target entry to a prioritized, actionable plan with minimal configuration.
2. **High-confidence findings.** Verification before reporting, so teams stop chasing false positives.
3. **Context-aware prioritization.** Rank by real-world exploitability and asset context, not raw CVSS.
4. **Built-in AI analyst.** Explain, plan, and generate fixes, always grounded in scanner evidence.
5. **AI-native integration.** Expose everything over MCP so assistants such as Claude can run and interpret scans.
6. **Nexpose-level ease of use.** Simple install, sane defaults, clean UI, config-as-code for power users.
7. **One platform, many surfaces.** Network, web/API, config, cloud, containers, code, external exposure.

## Non-goals (initially)

- Always-on, continuous monitoring or SIEM/EDR replacement. Bastion is used on demand and on schedule, not as a live sensor.
- Offensive tooling: no exploitation frameworks, no destructive payloads by default.
- Beating incumbents on raw plugin count in v1.
- Building new scan engines from scratch where proven open-source ones exist.

## Personas

| Persona | Needs |
|---|---|
| **Blue team analyst** | Accurate findings, evidence, fast re-tests |
| **Sysadmin / DevOps engineer** | Copy-paste fixes (Ansible, Terraform, shell), minimal noise |
| **Security manager / CISO** | Risk trend, MTTR, executive summaries |
| **Auditor / compliance lead** | CIS, PCI-DSS, NIST, ISO mappings and exportable evidence |
| **Small business / solo admin** | Zero-config, plain-language guidance |
| **AI agent (via MCP)** | Structured, safe, scoped tool access |

## Design principles

1. **Facts from scanners, reasoning from AI.** The AI never invents findings. Every claim links to evidence.
2. **Safe by default.** Non-destructive checks; intrusive checks need explicit approval.
3. **Authorized scope only.** Targets must be proven or attested before scanning.
4. **Explain everything.** Every score and prioritization decision is inspectable.
5. **Wrap, do not reinvent.** Orchestrate best-in-class open-source engines and normalize their output.
6. **Works offline.** Air-gapped deployments with local models and offline feeds.
7. **Automation parity.** Everything in the UI is available through API, CLI, and MCP.

## Success metrics

- Time from install to first useful scan: under 10 minutes
- False-positive rate on verified (Confirmed) findings: under 5%
- Share of findings with a generated, tested fix: over 80%
- Median time from finding to verified re-test: under 1 day for teams using the workflow
