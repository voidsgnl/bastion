# 0002 Job queue

Status: Open

## Context

The orchestrator distributes scan jobs to engines. See [07 Architecture](../07-architecture.md) and open question 7.

## Options

- NATS
- Redis
- RabbitMQ

## Decision

Not decided.

## Consequences

To be filled in. The choice affects the one-command install, since the queue ships inside the Docker Compose stack.
