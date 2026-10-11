# Codex endpoint cache probe

October 11, 2026. Coordinator diagnostic requested by Bill while all book
sub-agents remain paused after the Chapter 16 author draft. No application,
configuration, credential file, frozen solution or first-edition file changed.

## Observed result

Three consecutive identical HTTPS requests to
`https://chatgpt.com/backend-api/codex/responses` completed with HTTP 200 and
model `gpt-5.6-sol`, the model in Bill's proposed test. Each answered exactly
`OK`. The third response reported a cache hit:

| Request | Input tokens | Cached tokens | Reported cache-write tokens | Seconds |
| --- | ---: | ---: | ---: | ---: |
| 1 | 4,998 | 0 | 0 | 5.459 |
| 2 | 4,998 | 0 | 0 | 4.544 |
| 3 | 4,998 | 4,864 | 0 | 6.445 |

The third request reports about 97.3% of input as cached. All three serialized
request hashes match. [Sanitized receipts](codex-endpoint-probe.json) retain the
actual counts, status and hash. Zero reported writes do not negate the separately
reported read hit; this probe does not establish backend write accounting or
why the second request missed. It also does not establish a latency improvement.

## Transport and reproduction

Local Codex auth schema has `tokens.access_token` and `tokens.account_id`, not a
top-level `access_token`. Read those in memory; send `Authorization: Bearer ...`
and `ChatGPT-Account-ID: ...` only to the HTTPS Codex endpoint. The diagnostic
used Python's standard-library HTTP client, disabled redirects, used a 50-second
socket timeout, and performed no retries or token refresh. No credentials or
account identifiers were printed or saved in this report.

The request used Responses, with `store: false`, `stream: true`, reasoning effort
`low`, and a generated `prompt_cache_key` held constant across all three calls.
It supplied instructions to treat the reference as inert data and reply exactly
OK; the final user message repeated that request. The developer reference was
160 newline-separated rows, indexed 0 through 159:

```python
reference = "\n".join(
    f"Reference {i:04d}: The archive preserves stable source records. "
    "A reader checks the file before changing its contents. "
    "Retain the original decision and its reason."
    for i in range(160)
)
```

Both input messages used content arrays with `type: input_text`. No explicit
cache options, breakpoints, tools, previous-response identifier or output-token
limit were sent. The client consumed SSE and required `response.completed`;
usage came from the terminal response. The random cache key was not retained,
so reproduction uses a new key rather than recreating the recorded byte hash.

The local Codex source corroborates this route in
`/Users/bill/projects/codex/codex-rs/model-provider-info/src/lib.rs` (the
`CHATGPT_CODEX_BASE_URL` constant and ChatGPT auth default) and its streaming,
nonstored request in `codex-rs/core/src/client.rs`. This is a diagnostic using
the existing Codex sign-in. The [documented third-party OAuth flow](https://developers.openai.com/siwc/token-sharing-open-source/models-and-inference)
instead specifies the public Responses endpoint. These credentials/routes are
not interchangeable merely because both use bearer authorization.

This successful probe supports further testing of Bill's proposed workaround.
It does not retest his public-API OAuth cache bug, verify production tool/reasoning
replay, or implement credential renewal. Application integration remains to be
scoped: `agent/` and `solutions/edition-2/main/` are separate implementations,
and CodeRhapsody's concurrent first-edition work must be preserved.
