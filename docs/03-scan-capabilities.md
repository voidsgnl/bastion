# 03 Scan Capabilities

Bastion orchestrates proven open-source engines and normalizes their output into one finding schema ([schema](schema/finding.schema.json)). Engine choices below are initial candidates and may change after evaluation.

## Coverage matrix

| Area | Checks | Candidate engines |
|---|---|---|
| **Network** | Host discovery, port scan, service/version detection, CVE matching, default credentials, weak TLS/SSH ciphers, exposed management interfaces | Nmap, Greenbone/OpenVAS, Nuclei |
| **Web apps and APIs** | OWASP Top 10, auth/session flaws, security headers, misconfigurations, OpenAPI/GraphQL import, authenticated crawling | OWASP ZAP, Nuclei, custom crawler |
| **Configuration and hardening** | CIS benchmarks, SSH/WinRM authenticated checks, patch levels | OpenSCAP, custom authenticated checks |
| **Cloud** | Misconfigurations, public storage, over-permissive IAM, logging gaps | Prowler, ScoutSuite-style checks |
| **Containers and code** | Image CVEs, dependency CVEs, secrets, IaC misconfig, SBOM generation | Trivy, Grype, Semgrep, Gitleaks, Syft |
| **External exposure** | Subdomain enumeration, forgotten services, expiring certificates | Subfinder, httpx, certificate checks |

## Scan modes

- **Unauthenticated**: what an outsider sees.
- **Authenticated**: SSH, WinRM, SNMP, cloud role. Gives accurate patch and config state with fewer false positives.
- **Agentless first.** An optional lightweight agent for roaming or remote devices is a later phase.

## Verification levels

| Label | Meaning |
|---|---|
| **Confirmed** | A safe proof check demonstrated the issue (version read directly, config observed, behavior reproduced) |
| **Likely** | Strong indicators (banner/version match) but no direct proof |
| **Unverified** | Signature match only; flagged for manual review |

## Intrusiveness levels

| Level | Behavior | Default |
|---|---|---|
| 0 Passive | Reads only | allowed |
| 1 Safe active | Standard probes, no state change | allowed |
| 2 Intrusive | Brute force, fuzzing, possible service impact | requires explicit approval |
| 3 Destructive | Not supported | never |

## Knowledge base

Updated on a schedule or on demand from:

- NVD / CVE
- CISA KEV (known exploited vulnerabilities)
- EPSS (exploit prediction scores)
- Vendor advisories (Microsoft, Red Hat, Ubuntu, Debian, Cisco, and others)
- OVAL and CIS content
- Nuclei community templates

Feeds can be mirrored for offline use.
