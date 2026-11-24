# 04 Differentiators

What Bastion does that incumbents do poorly or not at all.

## 1. Verification, not just detection

For each finding Bastion runs a safe, non-destructive proof check and labels the result **Confirmed**, **Likely**, or **Unverified**. This targets the number one scanner complaint: false positives.

## 2. Context-aware risk scoring

Inputs:

- CVSS base score
- EPSS probability
- CISA KEV membership
- Internet exposure
- Asset role and criticality (database, domain controller, dev box)
- Compensating controls (WAF, segmentation, EDR) when known

Output: a single risk score with a visible explanation. "Critical here, ignorable there."

## 3. Fix-centric results

Findings are grouped by remediation, not by CVE. "Upgrade OpenSSL on these 23 hosts" is one task, not 40 rows.

## 4. Attack-path chaining

Links findings into realistic paths, for example: exposed admin panel -> default credential -> reachable database. Identifies the single fix that breaks the most paths.

## 5. Generated remediation

Per finding, as applicable:

- Patch command per OS and package manager
- Config diff
- WAF/virtual-patch rule
- Ansible, PowerShell, Terraform, or shell snippet
- Pull request for code and IaC issues
- Rollback note and risk of change

## 6. Scan diffing and fast re-tests

Compare any two scans. Re-test touches only the selected findings, so verification takes seconds.

## 7. AI analyst with evidence grounding

An assistant that explains, plans, and drafts, constrained to cite scan evidence. See [05 AI layer](05-ai-layer.md).

## 8. MCP-native

Any MCP-capable assistant can start scans, read findings, request fixes, and re-test, inside scope and approval guardrails. See [06 MCP server](06-mcp-server.md).

## 9. Easy configuration

One-command install, wizard, sane profiles, and config-as-code. See [11 Configuration](11-configuration-and-deployment.md).

## 10. Plain-language reporting

Executive, technical, and auditor reports generated on demand, tailored to the reader.

## Where Bastion will not (initially) win

- Raw plugin count versus Nessus/Qualys
- Enterprise-scale asset management versus InsightVM/Tenable.io
- Mature compliance content libraries

The plan is to close these gaps over time via open-source feeds and community content.
