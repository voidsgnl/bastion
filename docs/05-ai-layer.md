# 05 AI Layer

## Design rule

> Scanners find facts. The AI explains, prioritizes, and plans. It never invents findings, and every claim links to evidence.

## Capabilities

| Capability | Description |
|---|---|
| **Triage** | De-duplicate and group findings by root cause and fix |
| **Explain** | Plain-language explanation tailored to engineer, manager, or auditor |
| **Context reasoning** | "Why this matters on this asset" using role, exposure, and neighbors |
| **Next-step planning** | Suggest follow-up scans ("port 8443 looks like an admin panel, run the web profile?") |
| **Fix generation** | Shell, Ansible, PowerShell, Terraform, config diffs, WAF rules, PR drafts |
| **Fix review** | Assess risk of a proposed change and suggest rollback |
| **Attack-path narration** | Explain chained risks in prose |
| **Reporting** | Draft executive and technical summaries |
| **Q&A** | Chat panel with page-aware context |

## Grounding and anti-hallucination

1. The AI receives structured findings (JSON) and retrieved knowledge-base passages, not free-form guesses.
2. Every statement about a vulnerability must cite a finding ID, evidence item, or knowledge-base entry.
3. The AI cannot create or modify findings. Only scan engines and verification checks can.
4. Generated fixes are labeled *AI-generated* and shown with a diff and risk note before use.
5. Confidence labels (Confirmed / Likely / Unverified) are set by the verification engine, not by the model.

## Retrieval layer

Retrieval over:

- Vulnerability knowledge base (NVD, KEV, EPSS, vendor advisories)
- Remediation playbooks and vendor hardening guides
- The organization's own past findings, fixes, and exceptions

## Model options

| Mode | Description |
|---|---|
| **Local** | Ollama or similar; nothing leaves the network (air-gapped friendly) |
| **Cloud** | Hosted models via API (Claude or others), with redaction |
| **Off** | Full scanner functionality with no AI features |

The gateway is model-agnostic and configurable per deployment.

## Data protection

- Redact IPs, hostnames, credentials, and secrets before any cloud call (configurable policy, reversible mapping kept locally)
- Never send credential vault contents to a model
- Log every prompt and response metadata to the audit trail

## Agentic behavior (later phases)

An AI agent may propose and, with approval, execute scan plans:

- Propose scan -> human approves -> engine runs -> AI summarizes
- Intrusive actions always require human approval
- All AI-initiated actions are scope-locked and audited
