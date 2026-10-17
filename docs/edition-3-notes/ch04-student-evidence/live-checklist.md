# Chapter 4 live exercise checklist (saved before any live request)

## Entry configuration
- Core: actual `cmd/ch01` binary, `chat` subcommand, ordinary stdin prompts (two turns, then EOF); selected-model explanation gate enabled.
- Exact models: Anthropic `claude-sonnet-4-6` / Messages; OpenAI `gpt-4.1-2025-04-14` / Chat Completions; Gemini `gemini-3.8-flash` / generateContent. No substitution on errors.
- Root endpoints: `https://api.anthropic.com`, `https://api.openai.com`, `https://generativelanguage.googleapis.com`. Engine adds existing native surface path.
- Explicit environment: `LLM_VENDOR`, `LLM_BASE_URL`, `LLM_MODEL`, `LLM_API_KEY`, `LLM_REQUIRE_TOOL_EXPLANATION=true`, `ENSEMBLE_WORKSPACE`, `ENSEMBLE_DATA_DIR`, `CH02_LOG`; preserve normal build PATH/HOME/SDK environment. Each exercise gets distinct disposable workspace/data directories.
- Credentials: launcher reads only needed `directClaudeAPIKey`, `directOpenAIAPIKey`, `directGeminiAPIKey` fields from `~/.cr/settings.json` in memory and injects CLI environment. No key copies, prints, prompts, tool arguments, or source env files. Existing command runner strips all supported credential environment keys.
- External consumer separately maps all provider, gate, workspace/data, and log variables in its own public-API main; uses root-owned timestamped stderr logger through actual ToolContext.Agent.Parent().Logger(), and shutdown before Dump.
- Provenance: record candidate Git tree/commit and binary SHA256; request/response target metadata, process exit, visible stdout/stderr, and neutral events with opaque values omitted. Never emit settings, keys, or raw provider wire payloads.

## Feature → bounded exercise
|Feature|Live exercise|Additional deterministic evidence|
|---|---|---|
|Ordinary dispatch, persistence, caps, one-shot consumption|Anthropic writes and reads a large file under tool_limits, follows persisted output; OpenAI/Gemini file calls|Focused file/custom/failure and cursor tests; public bigoutput/toollimits checks|
|Delayed work, pattern wakeup, consumed cursors|Anthropic command emits READY, waits for line input, then reports result; follow via wait/send_input|Real PTY tests for already-consumed pattern, completed wait, zero replay|
|PTY/debugger|Anthropic runs prebuilt tiny Go fixture under `/Users/bill/go/bin/dlv exec`, sends breakpoint/continue/print/quit|Public debugger and real PTY environment/size/newline tests|
|Two-turn continuity and repair|Each accessible core provider repairs `greeting.txt` from `helo` to `hello` in turn 1, verifies that file and remembers session token in turn 2|Provider-boundary tests for replay/explanation/refusal/retry|
|Explicit kill and process groups|OpenAI/Gemini start retained shell job and kill it; inspect killed state|Real descendant-kill and blocked-waiter tests|
|Still-running EOF cleanup|Anthropic final turn starts `sleep 60` with immediate observation, does not wait/kill it, answers and reaches EOF; require `tool_returned` running then `job_killed` shutdown and process absence|Public shutdown and real shutdown test; do not count already-finished jobs|
|External registered tool|Consumer calls `shout`, returns uppercase, retains builtins, emits timestamped diagnostic|Consumer public grader + registration/owner/late-completion tests|
|Lifetime, star ownership, terminal/cursor concurrency|Do not claim live output proves source ownership/races|Full suite, race suite, Chapter1 lifetime gate, source review|

Run only these bounded cases. An access/model error is an observed limitation, not permission to downgrade. Record requested vs returned model attribution. Any missing live behavior is explicitly identified and supported only by the stated deterministic evidence unless a minimal corrective rerun is needed.
