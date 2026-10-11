# Chapter 27: Gateway

I built a fantasy role-playing game for AI agents called Deephold.
Moltbook had taken off as a social network for agents, and I thought a
game would bring them together. The agents followed their prompts. They
explored the dungeons and fought the monsters. They did not appear to
have a sense of fun. It was just more work for them.

And the architecture was awful. Each agent played by making a tool call
to wait for its turn, sitting blocked inside a tool handler while the
other players moved. I didn't figure out the right architecture on my
own. I read about how OpenClaw worked, and the need for a dedicated
message router became clear. I designed one, but it didn't really come
together until I built Puffin with CodeRhapsody, to help a family
member who needed daily care. Once the need was real, we built it with
care, and it hasn't failed since I started it in the spring of 2026.

Gateways are critical infrastructure. Scheduled tasks, turn-based games
where agents wait naturally, agents talking across servers, Discord,
home automation. Every one of those needs a router at the trust
boundary. That router belongs in the framework, not bolted on after.

---

## TL;DR

The gateway is a standalone edge service that sits at the trust boundary,
routing messages between agents over Noise-KK encrypted WebSocket
connections. It uses group-based access control where an agent's
capabilities derive from its group memberships, and parents can grant
children a subset of their own groups (narrow only, never widen). The
gateway administration surface is exposed as MCP tools tunneled over the
same Noise-KK connection.

Build a public client library in `agent/gateway/` with a `Handler`
interface, a `Client`, Noise-KK handshake, and envelope types. The library
must have zero imports from `agent/internal/`. Build a gateway server in
`agent/internal/gateway/` with WebSocket connections, group-based routing,
an offline queue, and MCP administration tools. Build the gateway binary
in `cmd/gateway/` and an interviewer agent in `cmd/interviewer/` that uses
the client library to interview Ensemble through the gateway.

Build and grade: `make grade27`.

| # | Behavior | Check | Points |
|---|----------|-------|--------|
| 1 | Noise-KK handshake | Two peers complete mutual authentication over a byte stream. Wrong keys are rejected. | 15 |
| 2 | Prompt envelope with provenance | Envelope carries source, trust level, conversation ID, and response mode. Provenance header prepended on delivery. | 10 |
| 3 | Gateway server accepts connections | `cmd/gateway` binary starts, accepts Noise-KK WebSocket connections from multiple agents, and tracks connection state. | 15 |
| 4 | Group-based routing | Agent declares contactability by group. Gateway delivers messages only from senders whose group membership the recipient has accepted. Unauthorized sends are rejected. | 15 |
| 5 | Delivery semantics | At-least-once delivery with ack. Messages queued when recipient is offline. Resume from last-acked sequence on reconnect. | 15 |
| 6 | Client library with Handler interface | `agent/gateway/` package exposes `Client` and `Handler` interface. A receive goroutine calls `Handler.OnDeliver` for incoming messages. Any Go program can implement `Handler`. | 10 |
| 7 | Parent registers child | Connected agent registers a child's public key and grants it a subset of its own groups. Child connects independently via Noise-KK. Gateway rejects group grants the parent does not hold. | 20 |

## 27.1 Why the agent needs a router

Every agent built so far in this book talks to one person: whoever is at the
keyboard. Prompts arrive through stdin or the GUI's WebSocket. Responses go
back the same way. The agent is a local process, and its world is the
terminal it was started in.

That changes here. An agent reachable by a Discord bot, a webhook, a
scheduled job, or another agent on a different machine is a network service,
and network services have a problem local processes do not: they cannot trust
their input.

A prompt arriving through stdin came from the operator. A prompt arriving
over a WebSocket from the internet came from whoever managed to connect. The
agent cannot tell the difference by inspecting the content. A message that
says "I am your operator" is just text. A message that says "trust=owner" in
its body is text with an assertion nobody verified. The agent needs someone
else to verify the sender before the message reaches its context window, and
that someone is the gateway.

An operator who runs an agent behind a gateway has already solved this
problem once: at the keyboard, the operator types and the agent answers. The
gateway extends that trust model across the network. The operator configures
which agents exist, which keys they hold, and which groups they belong to.
When a message arrives, the gateway checks the sender's identity against
that configuration and stamps the result onto the message. The agent reads
the stamp and knows as much about the sender as the operator decided to tell
it. The operator who configured the gateway is still the authority. The
difference is that the authority now works at a distance.

The gateway is a standalone binary that sits between the internet and the
agents. It accepts connections, authenticates them with Noise-KK mutual
authentication, routes messages based on group membership, and stamps every
delivered message with a provenance header. The header tells the agent who
sent the message and how much the gateway trusts that sender. The agent reads
the header the same way it reads any other prompt content, but the header was
written by the gateway after authentication, not by the sender.

The gateway is a router. It owns no agents, runs no prompts, and stores no
conversation history. The Ensemble connects to it as
a client, the same way any other program does. A standalone interviewer
binary connects to it the same way. A future Discord bot, a webhook handler,
or a scheduled job runner will all connect the same way. The gateway's job
is to authenticate connections, enforce routing policy, and hold messages
for agents that are not currently online.

The package layout reflects this separation. `agent/gateway/` is a public
client library with zero imports from `agent/internal/`. Any Go program can
import it, dial a gateway, and start sending and receiving messages.
`agent/internal/gateway/` is the server, private to the gateway binary.
`agent/cmd/gateway/` is the binary itself: 116 lines of flag parsing, signal
handling, and periodic persistence. Bill arrived at this architecture after
the Deephold experience proved that agents blocking inside tool calls to
wait for each other produces an unusable system. The gateway replaces that
pattern: an agent that wants to hear from another agent connects to the
gateway and handles messages when they arrive, through a callback, instead
of parking inside a blocking call. The code totals 3,058 non-test lines
across four packages, with 1,326 lines of tests.

## 27.2 Noise-KK beneath the cloud

TLS client certificates are the standard answer to mutual authentication on
the web. They are also the wrong answer here.

A gateway deployed behind Cloudflare never sees the TLS handshake.
Cloudflare terminates TLS at its edge, re-encrypts the connection to the
origin server, and the gateway receives plaintext HTTP from Cloudflare's
proxy. TLS authenticates Cloudflare to the client. It does not authenticate
the gateway to the agent, and it leaves the agent equally unauthenticated to the
gateway. A client certificate proves the agent's identity to Cloudflare
alone, which is the wrong party.

Noise-KK solves this by running a second layer of encryption inside the
WebSocket, beneath Cloudflare's TLS. The KK pattern means both sides know
each other's static public key before the handshake begins. The agent knows
the gateway's key from its configuration. The gateway knows the agent's key
from its configuration. The handshake proves that each side holds the private
half of the key the other side expects, establishes a session with forward
secrecy through ephemeral keys, and does all of it in two WebSocket frames
before any application data flows.

The sequence:

1. The agent opens a WebSocket connection to the gateway. TLS terminates at
   Cloudflare; the gateway sees plaintext.
2. The agent sends the first Noise-KK handshake message as a binary
   WebSocket frame. This contains the agent's ephemeral public key, encrypted
   with the gateway's known static key.
3. The gateway responds with the second handshake message. At this point both
   sides have authenticated and derived session keys.
4. Every subsequent WebSocket frame is encrypted with the session keys. The
   JSON inside is invisible to Cloudflare, to any proxy, and to anything that
   does not hold both static keys.

The keys are Curve25519, used for Diffie-Hellman key agreement through the
`flynn/noise` library. Some deployments label these `ed25519:` out of habit.
The label is wrong, and a wrong label on a cryptographic key is the kind of
small untruth that costs somebody an afternoon when they try to use it for
signing. The implementation labels them `x25519:` and accepts any prefix on
input, because correcting a label should not break a configuration that was
merely imprecise.

Key management is file-based. Each binary calls `LoadOrCreateKeypair`, which
reads a private key from disk or generates one and saves it. The public half
is derived from the private half on every load, so the pair is always
consistent. The gateway prints its public key on startup; Bill copies it
into the agent's configuration, copies the agent's public key into the
gateway's configuration, and the two can authenticate. No certificate
authority, no token rotation, no password. The keypair is the identity.

## 27.3 Trust stamped at the boundary

A message flowing through the gateway carries an envelope: a message ID, a
sequence number, a destination agent, a source, a timestamp, an optional
conversation ID, content, and a response mode. The source carries a type
(webhook, discord, scheduler, agent), a channel ID, a user ID, a display
name, and a trust level.

Three trust levels: `owner`, `trusted`, and `external`. The operator
authenticated by their own static key is `owner`. A service the operator
controls, such as a webhook from their own game server, is `trusted`. An
untrusted user reaching the agent over a trusted platform, such as a family
member in a Discord channel, is `external`.

Trust is stamped by the gateway from the group the sender belongs to, at the
moment of routing. The sender has no say. A message that says "trust=owner"
in its body is text, and an agent that believed it would be taking orders
from anyone who could type.

Consider a scenario that already exists. A family member asks Puffin to set
a reminder through a Discord channel. The Discord bot connects to the
gateway, and the gateway knows that messages from this channel carry
`external` trust. The agent reads
`[GATEWAY source=discord user=family-member trust=external]` at the top of
the prompt and knows exactly how much authority the request carries. The
family member cannot escalate to `owner` by typing the right words, because
the gateway wrote the header, and the gateway authenticated the channel.

On delivery, the gateway prepends a provenance header:

```
[GATEWAY source=webhook-game user=server trust=trusted]

Your turn in the dungeon. The goblin attacks.
```

The agent is a language model. Its only way to know who is talking to it is
to read it in the prompt. Everything the gateway knows about the sender that
the agent cannot verify for itself goes in this header, in a fixed format,
above the content it describes.

Groups are the authorization primitive. Each group carries a trust level, a
description, and a set of capabilities. An agent's base group memberships
are declared in the gateway configuration:

```json
{
  "groups": {
    "admin": {
      "description": "Full gateway administration",
      "capabilities": ["group-admin"],
      "trust": "owner"
    },
    "agents": {
      "description": "Standard agent group",
      "capabilities": ["send-message", "declare-contactable"],
      "trust": "trusted"
    }
  },
  "agents": {
    "ensemble": {
      "public_key": "x25519:A1B2C3...",
      "groups": ["admin", "agents"]
    }
  }
}
```

At runtime, agents declare which groups may contact them. A sender's message
is delivered only if the sender belongs to a group the recipient has
accepted. An unauthorized send is rejected at the gateway. The agent never
sees the message. The decision was made before the message entered its
context.

Parents can grant children a subset of their own groups, the same narrow-only
pattern from the sandbox chapter. A parent holding `admin` and `agents` can
grant a child `agents` alone. It cannot grant `interviewers` if it does not
hold `interviewers`. The system bootstraps from configuration and is
self-protecting from there.

Group administration is itself a capability. Creating, editing, and deleting
groups at runtime requires `group-admin`, which belongs to a root-level group
the operator controls in config. Runtime groups live in a separate namespace
from config groups. Configuration groups are immutable at runtime because a
runtime edit would silently revert on restart, and a gateway whose
authorization policy differs before and after a restart is a gateway that
passes security review on Tuesday and fails it on Wednesday.

The administration surface is exposed as MCP tools tunneled over the same
Noise-KK WebSocket, the same pattern as the MCP-over-WebSocket tunneling
built in chapter 12. No separate admin CLI, no out-of-band HTTP API, no
second authentication layer.

## 27.4 A library, not a framework

The `Handler` interface is four methods:

```go
type Handler interface {
    OnDeliver(env Envelope)
    OnConnected(welcome WelcomePayload)
    OnDisconnected(reason string)
    OnError(err error)
}
```

The client's receive goroutine calls these methods in the order traffic
arrives, one at a time. Returning from `OnDeliver` is what acknowledges the
message, so the client only acks work the handler has actually accepted. A
handler that blocks for ten seconds inside `OnDeliver` blocks the connection
for ten seconds. Work that takes real time belongs on a goroutine the handler
owns.

`NopHandler` implements all four methods as no-ops. Embed it to implement
only the methods that matter, so that adding a method to `Handler` later
does not break every existing implementation.

The `Client` connects, handshakes, and starts the receive goroutine:

```go
client, err := gateway.Connect(ctx, gateway.ClientConfig{
    GatewayURL:    "ws://127.0.0.1:8093/v1/ws",
    AgentID:       "interviewer",
    StaticKey:     key,
    GatewayPubKey: gwPub,
}, handler)
```

That call opens the WebSocket, runs the Noise-KK handshake, waits for the
gateway's welcome frame, and starts the goroutine that calls
`handler.OnDeliver` for incoming traffic. If the handshake fails because the
keys do not match, the error comes back from `Connect`. No partial
connection, no deferred failure. The caller's next line of code runs after
authentication is complete or has definitively failed.

The client's public surface is what makes the gateway reusable: `Send` a
message to another agent, `Respond` to a previously delivered message,
`RegisterChild` with the gateway by providing the child's public key and a
subset of the parent's groups, `DeclareContactable` to control which groups
may reach this agent, `Close` to disconnect, and `Done` to receive a signal
when the connection drops.

The `Handler` interface is what makes this a library rather than a framework.
The client owns the connection, the handshake, and the delivery bookkeeping.
It owns no opinion about what a message means. Ensemble implements `Handler`
to route envelopes into its actor system. The interviewer implements it to
collect answers. A home automation program would implement it to switch on a
light. The gateway cannot tell the difference, and it does not need to.

## 27.5 Delivery and recovery

After the Noise-KK handshake, both sides speak a JSON protocol inside the
encrypted stream. The agent sends five frame types: `hello` to identify
itself and request resume from a sequence number, `ack` to confirm delivery
of a message, `response` to reply to a prompt, `ping` for keepalive, and
`mcp_request` for gateway administration tool calls. The gateway sends five
frame types back: `welcome` to confirm the connection with the queue depth
and server version, `deliver` to send a prompt envelope with a sequence
number, `pong` for keepalive, `error` for protocol or authorization
failures, and `mcp_response` for tool call results. Ten frame types total,
five in each direction, each a JSON object with a `type` field.

The connection lifecycle starts with `hello`. The agent sends its agent ID,
protocol version, and the sequence number of the last message it
successfully acknowledged. The gateway looks up the agent by ID, verifies
that the Noise-KK static key matches the one in its configuration, and
responds with `welcome`. The welcome carries the number of queued messages
waiting for this agent, so the receiver can estimate how long the catch-up
will take.

Every delivered message carries a sequence number. The agent processes the
message, calls `handler.OnDeliver`, and when the handler returns, the client
sends an `ack` with the sequence number. The ack is the receipt. Tying the
ack to the handler's return is a deliberate choice: an ack sent before the
handler accepts the work loses the message if the process crashes between
the ack and the acceptance.

When a message arrives for an agent that is not connected, the gateway
queues it. The queue is a JSON file in the gateway's data directory, one per
recipient. Messages accumulate until the agent reconnects.

On reconnect, the agent sends `last_acked_seq` in the `hello` frame. The
gateway replays every queued message with a sequence number greater than
that value. The agent picks up where it left off. Every unacked message is
delivered again. The guarantee is at-least-once. Exactly-once would require
the gateway to track what the agent did with each message, and a router has
no business being an application server.

The server persists its dynamic state to JSON files in its data directory.
Runtime group memberships, contactability declarations, queued messages, and
child registrations all survive a restart. A ticker writes every thirty
seconds. A shutdown hook writes one final time. The thirty-second ticker
exists because `SIGKILL` never runs the shutdown path, and losing half an
hour of child registrations to an abrupt stop is the kind of data loss that
surfaces much later, when an agent registered twenty minutes ago cannot
connect and nobody remembers why.

## 27.6 The interviewer proves the stack

The interviewer is a standalone binary in `cmd/interviewer/`. It imports
`agent/gateway/`, holds its own keypair, and reaches Ensemble only through
the gateway's routing. It shares no code with Ensemble. If the interviewer
works, so does any other program anyone writes against the client library.

The binary connects to the gateway, declares itself contactable by the
groups its subject might reply from, and asks four fixed questions about
workflow friction: what slows you down, what burns the most tokens, where do
you need a human, and what new tool would help most. The questions are fixed
rather than model-generated so the interview is reproducible.

For each question, the interviewer sends an `OutboundMessage` through the
client and waits on the handler's answer channel with a per-question timeout.
The handler is seven lines: `OnDeliver` puts the envelope on a channel, and
`NopHandler` covers the rest.

When all four answers arrive, the interviewer writes a markdown report:
friction points in one section, improvement opportunities in another, raw
transcript at the bottom. It strips the provenance header from each answer
before writing, because `[GATEWAY source=agents user=ensemble trust=owner]`
is routing metadata, not interview content.

Three binaries, three processes, one proof:

1. Start the gateway with a config that knows both agents and their keys.
2. Start Ensemble, which connects and declares contactability.
3. Start the interviewer, which connects, asks four questions, and writes
   a file.

The interviewer reaches Ensemble only because the gateway's group membership
allows it. Change the group membership in the config, and the interviewer's
messages are rejected without Ensemble seeing them. The security property and
the routing property are the same property, verified by the same test.

The exercise proves four things at once. The Noise-KK handshake works
between two independent binaries. The client library is usable without any
Ensemble code. The group-based routing correctly authorizes and delivers
messages. And the provenance header appears on delivered content, carrying
the trust level the gateway assigned. A test that drove the gateway from
inside the same process would prove three of these. Only separate binaries
with separate keypairs prove all four.

## 27.7 What the grader checks

The grader drives the gateway binary and compiles a probe program against the
student's client library in a separate Go module, with `GOPROXY=off` to
ensure the import resolves from local source. This is the only honest proof
that `agent/internal/` stopped blocking external consumers: a program outside
the module that compiles and links against the public API.

Seven checks, three of which drive the real binary:

The Noise-KK check creates two peers with matching static keys and verifies
that the handshake completes. A second pair with mismatched keys must fail.
This runs in-process, without a gateway, because the handshake is a property
of the cryptographic layer, not the server. The Noise-KK handshake has a
property that makes mutation testing interesting: swallowing only the read
error leaves the crypto intact, because the wrong key produces a decryption
failure at the next step. Killing the check requires swallowing both the
read and write error paths, because the authentication is not separable from
the transport.

The envelope check verifies that the provenance header format matches the
contract exactly: `[GATEWAY source=<channel> user=<user> trust=<trust>]`,
followed by a blank line, followed by the original content.

The server check starts the gateway binary, connects two Noise-KK clients,
and verifies that both receive a welcome frame with queue depth and server
version.

The routing check has one agent declare contactability by group and another
send a message. The gateway must deliver when the sender's group is accepted
and reject when it is not. Authorization is verified at the gateway, not at
the recipient.

The delivery check disconnects an agent, sends a message while it is offline,
reconnects with a `last_acked_seq`, and verifies that the queued message
arrives. The ack and resume sequence is verified through the sequence
numbers.

The library check compiles the probe program and calls `Connect`, `Send`,
`DeclareContactable`, and `Close` against a live gateway started by the
grader.

The child registration check has a parent register a child's public key with
a subset of its own groups. The child connects independently and communicates
through the gateway. The parent then attempts to grant a group it does not
hold. The gateway rejects the grant. This check carries twenty points because
it verifies both the positive path and the narrow-only enforcement in a
single scenario.

The mutation audit runs seven mutants, each targeting a single check. Each
mutant kills exactly one check, proving the checks are independent and each
targets the behavior it names.

Build and grade: `make grade27`.

| # | Behavior | Check | Points |
|---|----------|-------|--------|
| 1 | Noise-KK handshake | Two peers complete mutual authentication. Wrong keys rejected. | 15 |
| 2 | Prompt envelope with provenance | Provenance header prepended on delivery. | 10 |
| 3 | Gateway server accepts connections | Binary starts, accepts multiple Noise-KK WebSocket connections. | 15 |
| 4 | Group-based routing | Delivers to accepted groups. Rejects unauthorized sends. | 15 |
| 5 | Delivery semantics | At-least-once with ack. Queue offline. Resume on reconnect. | 15 |
| 6 | Client library with Handler interface | External probe compiles and connects. | 10 |
| 7 | Parent registers child | Child connects independently. Over-grant rejected. | 20 |
