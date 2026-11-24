# Bastion

**An on-demand vulnerability scanner that tells you what to do next.**

Enter a target (IP range, domain, URL, cloud account, or repo), press **Scan**, and get a verified list of problems. Each one comes with a plain-language explanation, an exact fix, and a re-test button. There is no always-on monitoring: it is used like Nexpose or Nessus, with an AI analyst built into every step and an MCP server so any MCP-capable assistant can drive it.

> Status: **planning / documentation phase.** No code yet. See the [roadmap](docs/08-roadmap.md).

## Why Bastion

| Problem with existing scanners | Bastion's answer |
|---|---|
| Thousands of unranked findings | An action plan: "fix these 5 things first" |
| False positives | Every finding is verified and labeled Confirmed / Likely / Unverified |
| CVSS ignores context | Risk = CVSS + EPSS + CISA KEV + exposure + asset role + compensating controls |
| "Here is the CVE" with no fix | Exact remediation: commands, config diffs, Ansible/Terraform/PowerShell, PRs |
| Hard to set up | One-command install, wizard, first scan in about 5 minutes |
| Siloed from AI tooling | Built-in AI analyst plus an MCP server |

**Honest caveat:** Bastion will not beat Nessus or Qualys on raw check count on day one. It aims to win on accuracy, fix guidance, and time from scan to remediation, built on proven open-source engines.

## Documentation

| Doc | Contents |
|---|---|
| [01 Vision](docs/01-vision.md) | Goals, non-goals, personas, principles |
| [02 User experience](docs/02-user-experience.md) | The scan-to-fix workflow |
| [03 Scan capabilities](docs/03-scan-capabilities.md) | What is scanned and with which engines |
| [04 Differentiators](docs/04-differentiators.md) | What makes it better than existing tools |
| [05 AI layer](docs/05-ai-layer.md) | AI analyst design and rules |
| [06 MCP server](docs/06-mcp-server.md) | Tool definitions for MCP clients |
| [07 Architecture](docs/07-architecture.md) | Components, data flow, tech stack |
| [08 Roadmap](docs/08-roadmap.md) | Phases and milestones |
| [09 Security and guardrails](docs/09-security-and-guardrails.md) | Authorization, scope-locking, safe use |
| [10 UI concept](docs/10-ui-concept.md) | Screens and interaction design |
| [11 Configuration and deployment](docs/11-configuration-and-deployment.md) | Install, engines, config-as-code |
| [12 Competitive landscape](docs/12-competitive-landscape.md) | Existing tools and positioning |
| [13 Open questions](docs/13-open-questions.md) | Decisions still to make |
| [Finding schema](docs/schema/finding.schema.json) | Normalized finding format |

## Responsible use

Bastion must only be used against systems you own or are explicitly authorized to test. The product enforces authorization and scope-locking by design. See [docs/09-security-and-guardrails.md](docs/09-security-and-guardrails.md).

## License

To be decided (see [open questions](docs/13-open-questions.md)).
