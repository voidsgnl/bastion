# 06 MCP Server

Bastion exposes its capabilities through the Model Context Protocol so Claude and other MCP-capable assistants can operate it. Bastion can also act as an MCP **client**.

## Bastion as MCP server

### Tools

| Tool                                        | Purpose                                                       | Notes                                                             |
| ------------------------------------------- | ------------------------------------------------------------- | ----------------------------------------------------------------- |
| `list_targets`                            | List authorized targets                                       | read-only                                                         |
| `add_target(target, type)`                | Add a target and start authorization                          | returns verification instructions                                 |
| `start_scan(target_id, profile, options)` | Launch a scan                                                 | requires authorized target; intrusive options need human approval |
| `get_scan_status(scan_id)`                | Progress, hosts found, ETA                                    | read-only                                                         |
| `get_findings(scan_id, filters)`          | Findings by severity, confidence, asset, tag                  | read-only, paginated                                              |
| `explain_finding(finding_id, audience)`   | Evidence, impact, context                                     | audience: engineer / manager / auditor                            |
| `get_attack_paths(scan_id)`               | Chained risks and best chokepoint fix                         | read-only                                                         |
| `get_remediation(finding_id, format)`     | Fix as shell / ansible / powershell / terraform / ticket / pr | format-specific                                                   |
| `retest(finding_id)`                      | Re-run verification for one finding                           | safe-level checks only                                            |
| `compare_scans(scan_a, scan_b)`           | New / fixed / regressed                                       | read-only                                                         |
| `export_report(scan_id, type, format)`    | executive / technical / compliance; pdf / docx / csv / sarif  | returns file reference                                            |
| `suggest_next_scan(scan_id)`              | AI-proposed follow-up scans                                   | proposal only                                                     |

### Resources

- `bastion://targets` - authorized scope
- `bastion://scans/{id}` - scan summary
- `bastion://findings/{id}` - finding with evidence
- `bastion://kb/cve/{id}` - enriched CVE entry (CVSS, EPSS, KEV)

### Prompts

- `triage_scan` - walk through a scan and produce a prioritized plan
- `write_exec_summary` - management summary for a scan
- `plan_remediation_sprint` - group fixes into a time-boxed work plan

## Safety model for MCP access

- **Scope-locked:** tools only accept targets already authorized.
- **Scoped tokens:** each MCP client gets a token with specific permissions (for example read-only, or scan-allowed).
- **Approval gates:** `start_scan` with intrusiveness level 2 returns `approval_required`; a human approves in the UI.
- **Rate limits** per token and per target.
- **Audit log:** every call recorded with client identity, arguments, and result summary.
- **No destructive tools exist.** There is no tool for exploitation or destructive actions.

## Transport and auth

- stdio for local use and streamable HTTP for remote
- OAuth 2.1 or API tokens; mTLS optional for on-prem

## Bastion as MCP client

Planned integrations:

| Integration            | Use                                          |
| ---------------------- | -------------------------------------------- |
| CMDB / asset inventory | Enrich asset owner, criticality, environment |
| Jira / ServiceNow      | Create and sync remediation tickets          |
| Slack / Teams          | Notifications and approvals                  |
| GitHub / GitLab        | Open fix pull requests                       |
| Cloud providers        | Read-only inventory                          |

## Example session

```
User (in Claude): Scan 10.0.4.0/24 with the safe profile and tell me what to fix first.

Claude -> start_scan(target_id="t_104", profile="safe")
Claude -> get_scan_status("s_881")   ... polls until complete
Claude -> get_findings("s_881", {confidence: "confirmed", sort: "risk"})
Claude -> get_attack_paths("s_881")
Claude -> get_remediation("f_2210", "ansible")

Claude: Fix these 4 items first. Upgrading OpenSSH on 6 hosts and removing the exposed
Jenkins admin panel breaks 7 of 9 attack paths. Here is the Ansible playbook...
```
