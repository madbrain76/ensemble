# Chapter 1 — raw HTTP conversation

This module starts from the current Edition 3 Chapter 1, not a previous solution.
It uses only the Go standard library. No tools, retries, streaming, persistence,
provider framework, or background conversation machinery are introduced.

## Ownership rationale

The public `ensemble` package is a small composition root. It owns a logger and
its Agents, constructing both directions of the Agent/Engine relationship before
returning an Agent. `internal/common` declares shared values and interfaces only;
`internal/agent` and `internal/engine` never import each other.

Agent holds configuration and completed history. Engine holds the HTTP client
and measured usage, looking up configuration and diagnostics through its actual
Agent parent (and that parent's Ensemble). A failed request never commits an
unanswered user turn. Missing usage is not zero; measured usage from an empty
answer is still counted. History and usage accessors return copies, not mutable
internal state. Calls on a single Agent and creation within an Ensemble are
synchronous; separate Agents have independent conversation and accounting.

The HTTP client has no total response or header deadline. Standard transport
connection/TLS protections remain, while the caller's context governs model
work, including reading the body. Redirects are rejected to avoid forwarding
credentials. Errors omit remote bodies, URLs, and raw transport diagnostics.

## Use

Set `ANTHROPIC_API_KEY` and `ANTHROPIC_MODEL`; discover the model with the
provider's `GET /v1/models` endpoint rather than guessing. `ANTHROPIC_BASE_URL`
defaults to `https://api.anthropic.com`, and can select a local HTTP fake.

```sh
go build -o /tmp/ch01 ./cmd/ch01
printf '%s\n' '{"user":"Hello"}' | /tmp/ch01
/tmp/ch01 chat
```

No arguments selects JSON-lines input/output: one `assistant` object per input,
then cumulative `usage` at EOF. Errors go to stderr and exit nonzero without
retrying. `chat` prints human prompts and per-turn usage; Ctrl-D ends it.
The root accepts an application logger (`ensemble.New(logger)`), then
`app.NewAgent(ensemble.Config{...})` creates a conversation for synchronous
`Ask(ctx, text)` calls. A nil logger discards diagnostics; no global logger is used.

## Checks

```sh
go test -race ./...
go vet ./...
```

Tests use real local HTTP servers with fake credentials. They cover wire shape,
parent/configuration/logger isolation, replay and failed-turn retention, usage,
sanitized errors, caller cancellation, redirects, and the CLI's line/EOF rules.
JSON input has no arbitrary physical-line size limit; a final unterminated line
is processed at EOF. The default JSON mode does not print chat prompts.
