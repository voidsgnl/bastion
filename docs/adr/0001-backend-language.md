# 0001 Backend language

Status: Accepted

## Context

The console backend and scan workers need a language. See [07 Architecture](../07-architecture.md) and open question 6.

## Options

| Option | For | Against |
|---|---|---|
| Go | Concurrency, static binaries, simple container images | Slower to prototype |
| Python (FastAPI) | Faster prototyping, familiar to the scanning ecosystem | Packaging and concurrency need more care |

## Decision

Go is accepted for the first implementation of both the console backend and scan
worker. It provides static binaries, straightforward concurrency, and small
container images. FastAPI remains available for isolated integrations if a later
prototype proves Go inappropriate.

## Consequences

When decided, update the `server/` and `engine/` layout in [14 Repository Structure](../14-repo-structure.md).
