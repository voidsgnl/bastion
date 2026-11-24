# 11 Configuration and Deployment

## Install goals

- One command to a working console
- Wizard for the first scan
- No manual database or queue setup

## Install options

| Option | Command / artifact | Use |
|---|---|---|
| Docker Compose | `docker compose up -d` | Labs, small teams, single node |
| Helm chart | `helm install bastion ...` | Kubernetes |
| Appliance image (OVA / cloud image) | Import and boot | Nexpose-like experience |
| Air-gapped bundle | Offline archive with feeds and local model | Isolated networks |

## First-run wizard

1. Create admin account
2. Choose AI mode: local model, cloud model, or off
3. Add first target and verify authorization
4. Pick a profile and launch the first scan

## Scan engines

- Console ships with a built-in local engine.
- Additional engines: run one container per network segment, pair with a token:

```
docker run -d --name bastion-engine \
  -e BASTION_CONSOLE=https://console.example.org \
  -e BASTION_PAIRING_TOKEN=<token> \
  bastion/engine
```

- Pairing uses short-lived tokens and establishes mTLS.

## Credentials

- Enter once in the encrypted vault; assign to targets or profiles
- Types: SSH key/password, WinRM, SNMPv3, cloud role, API tokens, web login recipes
- Test-credential button before scanning

## Config-as-code (YAML)

```yaml
version: 1
targets:
  - name: corp-dmz
    type: cidr
    value: 10.0.4.0/24
    owner: infra-team
    criticality: high
profiles:
  nightly-safe:
    base: safe
    exclusions:
      - 10.0.4.10
schedules:
  - target: corp-dmz
    profile: nightly-safe
    cron: "0 2 * * *"
integrations:
  jira:
    project: SEC
    create_for: [critical, high]
ai:
  mode: local
  model: llama3.1
  redaction: strict
```

## CLI (planned)

```
bastion target add 10.0.4.0/24 --verify attest
bastion scan start corp-dmz --profile safe --wait
bastion findings list --confidence confirmed --sort risk
bastion fix get f_2210 --format ansible
bastion retest f_2210
bastion report export s_881 --type executive --format pdf
```

## API

- REST (OpenAPI spec published) and webhooks (`scan.completed`, `finding.new`, `finding.fixed`)
- Same permissions model as the UI and MCP

## Updates

- Signed application updates and signed feed bundles
- Online auto-update or offline import
- Zero-downtime upgrades for distributed engines (engines can lag one version)

## Backup and retention

- Single command backup of database, evidence, and vault (encrypted)
- Configurable retention for scans and evidence
