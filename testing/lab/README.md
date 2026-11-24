# Local scanner lab

This lab exposes a plain Nginx page on `127.0.0.1:18080` for the opt-in Nmap
runner. The port is bound to loopback and is not published on external interfaces.

Start and stop the target:

```bash
docker compose -f testing/lab/docker-compose.lab.yml up -d
docker compose -f testing/lab/docker-compose.lab.yml down
```

In another terminal, with Nmap installed, run:

```bash
make run-nmap-lab
```

Open <http://localhost:8080> and use **New scan**. The UI records an explicit
attestation for `127.0.0.1/32`; the Nmap runner independently rejects any other
target or profile. It uses a fixed TCP connect probe of port `18080`. It does not
run version detection, NSE scripts, vulnerability checks, or caller-provided
commands. The resulting open-port observation is informational, not a vulnerability
finding.

Stop the server with Ctrl+C and stop the target with the Compose command above.
