# Ported feature: retry a tool call cut off by the output budget

Branch `port/5-truncated-tool-call`, based on `fix/linux-path-case`. Coordinator, 2026-11-03. Ported from `go-checkpoint` (`3429b42c`, tool-call arm).

A call truncated by the cap is not a call, and its partial arguments must never be dispatched. When the provider reports its length verdict while the reply was reaching for a tool, the attempt is repeated once at a doubled cap and the partial result is discarded; a second truncation is reported instead of retried, bounded by `toolRetries` so a persistently capped reply cannot recurse to the ceiling.

Detection lives in the parser, where the finish reason actually exists, and is recorded on `ParsedResponse.ToolCallCapped`. Only the vendor values that mean *budget exhausted* qualify: `length`/`max_tokens` for OpenAI, `max_tokens` for Anthropic, `MAX_TOKENS` for Gemini. A missing reason, a filter stop, or unreadable JSON is not grounds for a retry.
