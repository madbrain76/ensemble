# Chapter 27 Brief — Gateway (for Author)

## One-sentence thesis

The gateway is core infrastructure: a standalone edge service at the trust
boundary that routes messages between agents over Noise-KK encrypted
WebSocket connections, using group-based access control where capabilities
flow from membership and parents can only narrow what they grant children.

## Why this chapter exists

Everything Bill wants to build — Discord bot, home automation, scheduler,
self-improvement loop — depends on the gateway. It is the foundation that
turns a local agent into a networked one. Without it, agents can only talk
to whoever is at the keyboard.

## What the student builds

1. **Noise-KK mutual authentication** — end-to-end encryption beneath
   Cloudflare's TLS termination. Both sides prove identity via ed25519 keys.
   No passwords, no API keys, no certificates.

2. **Gateway server** — accepts connections, enforces group-based routing,
   queues messages for offline agents, resumes delivery on reconnect.

3. **Public client library** — the `agent/gateway/` package with a `Handler`
   interface anyone can implement. This is what makes the gateway
   infrastructure rather than just a program.

4. **Interviewer agent** — a standalone binary that connects to Ensemble
   through the gateway, interviews it about friction points, and writes
   a findings doc. Proves the whole stack works: two separate processes
   talking over Noise-KK through a gateway.

## Key design decisions to explain

- **Groups as ACLs, not static routes.** Agents declare contactability by
  group. The gateway checks membership on delivery. This is dynamic — agents
  control their own inbound at runtime.

- **Narrow-only inheritance.** Parents can grant children a subset of their
  own groups. Same pattern as sandbox (ch23). The system bootstraps from
  config and is self-protecting from there.

- **Gateway admin is MCP tools.** No separate admin CLI. Agents manage the
  gateway through tool calls tunneled over the same Noise-KK WebSocket.
  Group administration requires the `group-admin` capability — the system
  protects itself.

- **Handler interface, not a framework.** The client library defines a
  callback interface. Ensemble is just one implementation. A home automation
  agent, a Discord bot, an interviewer — they all implement the same
  interface. This is what makes it a library.

- **Ensemble is a client, not the gateway.** The gateway owns no agents.
  It is a router at the trust boundary. Ensemble connects to it the same
  way any other agent does.

## Connections to prior chapters

- **Ch23 (Sandbox)**: narrow-only capability inheritance is the same pattern
- **Ch25 (Sub-agents)**: parent-child relationship, but now across process
  boundaries via the gateway
- **Ch26 (Swarm management)**: the gateway extends fleet management across
  the network
- **Ch12 (MCP)**: MCP-over-WebSocket tunneling, now used for gateway admin

## What comes next

- **Connectors chapter**: Discord, webhooks, home automation — they connect
  to the gateway and produce prompts
- **Scheduler chapter**: cron-like job runner that fires prompts into the
  gateway's router
- **Self-improvement**: the interviewer exercise is the seed of this
