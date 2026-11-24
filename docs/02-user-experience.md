# 02 User Experience

## The core loop

```
Add target -> Pick profile -> Scan -> Action plan -> Fix -> Re-test -> Report
```

## 1. Add target

Paste an IP range, CIDR, domain, URL, cloud account, or repository.

Before any scan, Bastion establishes authorization:

- **External domains:** DNS TXT record or a file uploaded to a well-known path
- **Internal ranges:** signed-in attestation ("I am authorized to test this scope") recorded in the audit log
- **Cloud accounts:** read-only credentials or role the user provisioned (itself proof of access)
- **Repos:** OAuth or deploy-key access

## 2. Pick a profile (no configuration required)

| Profile | Use |
|---|---|
| **Quick** | Fast discovery, top ports, known-exploited vulnerabilities |
| **Full** | Deep network and service scan, all supported checks |
| **Web app** | Crawl and test a site or API (OWASP Top 10, auth, misconfig) |
| **Compliance** | CIS, PCI-DSS, NIST benchmark checks |
| **Safe mode** | Strictly non-intrusive checks for production systems |
| **Custom** | Compose your own (power users) |

## 3. Scan

A live view shows hosts discovered, services identified, and findings appearing in real time, with progress and an ETA. Scans can be paused, resumed, or cancelled.

## 4. Results as an action plan

Instead of a flat table of thousands of rows:

- Headline: "Fix these 5 things first and you remove about 80% of your real risk."
- Findings grouped by **the fix**, not the CVE (one patch can close 40 findings).
- Each finding shows:
  - What it is, in plain language
  - Why it matters on this specific asset
  - Evidence (request/response, banner, config line)
  - Confidence: Confirmed, Likely, or Unverified
  - Exact fix with tabs for Linux, Windows, Cloud, and Code
  - **Re-test** button

## 5. Ask the AI

A side panel is present on every screen with the page context preloaded. Example prompts:

- "Which of these can an attacker chain together?"
- "Write the Ansible playbook for these 12 hosts."
- "Explain this finding to my manager."
- "What should I scan next?"

## 6. Fix and re-test

After remediation, Re-test runs only the checks for the selected findings (seconds, not hours). Results update the finding to Fixed or Still Vulnerable, with a diff since the last scan.

## 7. Report

One click for:

- **Executive summary** (risk trend, top actions, MTTR)
- **Technical report** (full evidence and remediation)
- **Compliance pack** (control-by-control, auditor-ready)

Formats: PDF, DOCX, CSV, JSON, SARIF.

## Scan diffing

"Since last scan: 14 fixed, 3 new, 2 regressed." Every scan is comparable to any previous scan of the same target.

## First-run experience

1. Install with one command
2. Wizard: create admin account, choose AI mode (local, cloud, or off), add first target
3. First scan launches, with results in about 5 minutes for a quick profile
