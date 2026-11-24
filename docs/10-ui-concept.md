# 10 UI Concept

Goal: clean, fast, and understandable by someone who has never used a scanner. Nexpose-level familiarity with a modern look.

## Global

- Left navigation: Home, Targets, Scans, Findings, Attack paths, Reports, Integrations, Settings
- Global search (keyboard shortcut) across assets, findings, CVEs
- Dark and light themes, accessible color contrast
- Persistent **AI panel** on the right, with context of the current page
- Every list exportable; every action available via API

## Screens

### Home
- Overall risk score and trend
- **Top 5 actions** (the fixes that remove the most risk)
- Recent scans and their status
- Exposure map (internet-facing vs internal)
- Findings by confidence and age

### Targets
- List with authorization status (pending, verified, attested)
- Add-target wizard with verification instructions

### New scan
- Choose target and profile (Quick, Full, Web, Compliance, Safe)
- Optional: credentials, schedule, exclusions
- Estimated duration and intrusiveness shown before start

### Live scan view
- Real-time graph of hosts and services
- Findings stream with confidence labels
- Pause / resume / cancel

### Results: action plan
- Fix cards: "Upgrade OpenSSH on 6 hosts: removes 11 findings, breaks 5 attack paths"
- Toggle between By fix, By asset, By CVE
- Filters: severity, confidence, exposure, owner, tag

### Finding detail
- Summary and plain-language explanation
- Risk breakdown (CVSS, EPSS, KEV, exposure, asset role)
- Evidence viewer (request/response, banner, config)
- Fix tabs: Linux, Windows, Cloud, Code, plus Copy and Download
- **Re-test** button with live result
- History: first seen, last seen, status changes, comments

### Attack-path view
- Interactive graph of assets and chained findings
- Highlighted chokepoint fix
- Click a node to see the finding

### Compare scans
- New, fixed, regressed, unchanged with trend chart

### Reports
- Executive, technical, compliance; preview and download
- Scheduled report delivery

### Settings
- Users and roles, SSO, API tokens (including MCP tokens)
- AI mode (local, cloud, off) and redaction policy
- Engines (health, pairing, versions)
- Integrations (ticketing, chat, CMDB)
- Audit log viewer

## Interaction principles

1. Show the answer first (what to fix), details second.
2. Never show a number without a way to see where it came from.
3. Progressive disclosure: simple by default, power on demand.
4. Every finding is one click from evidence and one click from a fix.

## Prototype plan

A clickable prototype of Home, Results, and Finding detail is planned for Phase 0 (see [roadmap](08-roadmap.md)).

A first-pass UI for these three screens is in [ui](../ui/README.md): plain HTML and CSS, open `ui/index.html` in a browser. It uses sample data, and only the navigation between the three screens works so far.
