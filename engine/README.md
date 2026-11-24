# engine

Scan worker. Runs one container per network segment and pairs with the console using a short-lived token over mTLS (docs/11). Each folder under `adapters/` wraps one scan tool and outputs the common finding schema. Phase 1: nmap, nuclei, zap. Later: openvas, trivy, prowler.
