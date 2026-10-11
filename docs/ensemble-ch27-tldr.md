# Chapter 27 — Gateway

## TL;DR

The gateway is core infrastructure: a standalone edge service that sits at
the trust boundary, routing messages between agents over Noise-KK encrypted
WebSocket connections. It uses group-based access control where an agent's
capabilities derive from its group memberships, and parents can grant children
a subset of their own groups (narrow only, never widen). The gateway
administration surface is exposed as MCP tools tunneled over the same
Noise-KK connection.

### Deliverables

| # | Deliverable | Observable | Points |
|---|------------|------------|--------|
| 1 | **Noise-KK handshake** | Two peers complete mutual authentication over a byte stream. Wrong keys are rejected. | 15 |
| 2 | **Prompt envelope with provenance** | Envelope carries source, trust level (derived from sender's group), conversation ID, and response mode. Provenance header is prepended on delivery. | 10 |
| 3 | **Gateway server accepts connections** | `cmd/gateway` binary starts, accepts Noise-KK WebSocket connections from multiple agents, and tracks connection state. | 15 |
| 4 | **Group-based routing** | Agent declares contactability by group. Gateway delivers messages only from senders whose group membership the recipient has accepted. Unauthorized sends are rejected. | 15 |
| 5 | **Delivery semantics** | At-least-once delivery with ack. Messages queued when recipient is offline. Resume from last-acked sequence on reconnect. | 15 |
| 6 | **Client library with Handler interface** | `agent/gateway/` package exposes `Client` and `Handler` interface. A receive goroutine calls `Handler.OnDeliver` for incoming messages. Any Go program can implement `Handler`. | 10 |
| 7 | **Parent registers child** | Connected agent registers a child's public key and grants it a subset of its own groups. Child connects independently via Noise-KK. Gateway rejects group grants the parent does not hold. | 20 |

**Total: 100 points.** `make grade27`

### What to build

**Public client library** (`agent/gateway/`):
- `Handler` interface with `OnDeliver`, `OnError`, `OnConnected`, `OnDisconnected`
- `Client` with `Connect`, `Send`, `Respond`, `RegisterChild`, `DeclareContactable`, `Close`
- Noise-KK handshake adapted from OpenADP (`noise_nk.go`)
- Envelope type with trust levels derived from group membership
- Zero imports from `agent/internal/`

**Gateway server** (`agent/internal/gateway/`):
- WebSocket server accepting Noise-KK connections
- Config: agents (public keys, base group memberships), groups (capabilities, trust level)
- Group-based routing: deliver only if sender is in a group the recipient accepts
- Connection manager with ack tracking, offline queue, resume
- MCP tool surface for runtime administration (declare contactable, register child, grant membership)
- JSON persistence for dynamic state (runtime memberships, contactability)

**Gateway binary** (`agent/cmd/gateway/main.go`):
- Wires server + config + persistence
- Graceful shutdown

**Interviewer agent** (`agent/cmd/interviewer/main.go`):
- Standalone binary using the client library
- Connects to gateway, sends interview questions to Ensemble
- Writes findings to a markdown doc

### What NOT to build

- Channels/connectors (Discord, webhooks, email) — next chapter
- Scheduler — next chapter
- Ensemble actor-loop integration — separate concern, deferred
- Web dashboard — CLI + logs
- TLS certificate management — Cloudflare handles it
