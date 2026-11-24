# 12 Competitive Landscape

| Tool | Strengths | Gaps Bastion targets |
|---|---|---|
| **Rapid7 Nexpose / InsightVM** | Ease of use, good reporting, asset focus | Limited AI, prioritization still heavy, commercial only |
| **Tenable Nessus / Tenable.io** | Huge plugin library, industry standard | Noisy results, fix guidance generic, pricey at scale |
| **Qualys VMDR** | Cloud platform, broad coverage, agents | Complexity, cost, UX |
| **Greenbone / OpenVAS** | Open source, solid network checks | Rough UI, setup friction, little prioritization or fix help |
| **Wazuh** | SIEM/XDR plus vulnerability detection | Detects vulnerable installed packages, not a scanner; always-on model |
| **Nuclei** | Fast, community templates | CLI-first, no workflow, no prioritization |
| **OWASP ZAP / Burp** | Strong web testing | Specialized, manual-heavy, not a platform |
| **Trivy / Grype** | Excellent container and dependency scanning | Narrow scope, no network or web coverage |
| **Prowler / ScoutSuite** | Cloud posture checks | Cloud only |

## Positioning

> The ease of Nexpose, the open engines of the community, and an AI analyst that turns findings into fixes.

## Where Bastion competes

1. Time from scan to verified fix
2. Trustworthy, evidence-backed findings
3. AI and MCP integration as a first-class interface
4. Open-core pricing for small and mid-size teams

## Where Bastion complements

Bastion can ingest results from other scanners (SARIF, Nessus, Nmap XML, Trivy JSON) to prioritize and generate fixes, which lowers the barrier to adoption for teams already invested in other tools.

## Build, assemble, or adopt?

If the need is only to scan and report, adopting **Greenbone/OpenVAS** plus **Nuclei** is cheaper than building. Bastion is justified if the verification, prioritization, fix generation, and AI/MCP workflow are the product.
