# API contracts

Contracts shared by the server, engine, CLI, web app and MCP clients. These are stubs until Phase 1.

| File | Contract |
|---|---|
| `openapi.yaml` | REST API |
| `proto/engine/v1/engine.proto` | Console to engine (gRPC over mTLS) |
| `mcp/tools.json` | MCP tool list. Full definitions are in [docs/06-mcp-server.md](../docs/06-mcp-server.md) |

The finding format lives in [docs/schema/finding.schema.json](../docs/schema/finding.schema.json).
