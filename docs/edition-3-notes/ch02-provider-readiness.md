# Chapter 2 — finite direct-provider readiness receipt

## Scope and verdict

Original pass observed **2026-10-16, 20:19–20:22 UTC** (host clock); the explicitly authorized Chat Completions follow-up completed at **20:29:48 UTC**. Repository HEAD recorded by the original read and reverified during the provenance audit: `5351583208988f8370937cd0494727645efa8b46`. See the correction record below; the initial receipt transcribed this incorrectly.

**Setup result: all three chapter interfaces (Anthropic Messages, OpenAI Chat Completions, Gemini generateContent) returned `OK`; the extra parent-requested OpenAI Responses interface also returned `OK`. Four generation attempts total across both turns.** This establishes credential/model/interface usability at the observation time only. It is **not** Chapter 2 student CLI acceptance, implementation acceptance, a grader result, or reviewer-independent source-code acceptance. No student/app source was inspected to establish correctness, and no student program was run.

Guidance consulted: root `AGENTS.md`, the complete `.agents/skills/ensemble-coding/SKILL.md`, `docs/edition-3-execution.md`, the public Chapter 2 provider/live-use contract in `book/edition-3/chapter-02.md` (including §§2.2, 2.6, 2.9), and the public student/readiness documents. This receipt tests **native OpenAI Responses** as expressly requested. The previous characterization of the chapter's wording as historical was an unsupported interpretation and is withdrawn. The parent subsequently confirmed that the chapter teaches Chat Completions and acknowledged that the parent's original native-Responses contract assertion/instruction was mistaken. The original Responses check is extra parent-requested setup, **not a chapter requirement**. The separately authorized Chat Completions follow-up checks the actual chapter interface; neither check establishes student implementation acceptance. Chapter 2 still requires the student's own typed CLI chat, real provider reply, durable events, and clean-process replay demonstration for its supported providers.

## Authorization and handling

Bill authorized use of exactly these fields in `~/.cr/settings.json`:

| Vendor | Credential field reference (never the value) | Authentication used |
|---|---|---|
| Anthropic | `directClaudeAPIKey` | `x-api-key`; `anthropic-version: 2023-06-01` |
| OpenAI | `directOpenAIAPIKey` | `Authorization: Bearer …` |
| Gemini | `directGeminiAPIKey` | `x-goog-api-key` |

The JSON store was parsed locally; only those three fields were selected for requests. Credential values remained in process memory and authenticated request headers: no values were printed, passed in command arguments, copied into artifacts, or sent to documentation tools. No credential-bearing request/header logs were enabled. No broad secret search, GUI/live session access, OAuth/Codex fallback, or proxy was used. Direct Python HTTPS used normal TLS validation and explicitly disabled proxy handling; generation redirects were also disabled. The public documentation fetches were unauthenticated.

Only this receipt and sanitized JSON under `docs/edition-3-notes/ch02-entry-evidence/` were written by this task. No application/student/manuscript/grader changes, deployment, restart, commit, or push were performed.

## Current authenticated model discovery

Each listing was fetched once directly, returned **HTTP 200**, and reported no next page (`has_more: false`). A successful listing is **not inference proof**.

| Vendor | Discovery URL | Selected ID present in listing | Selection rationale |
|---|---|---|---|
| Anthropic | <https://api.anthropic.com/v1/models?limit=1000> | `claude-haiku-4-5-20251001` | Economical Haiku text model; official pricing lists Haiku 4.5 at $1/input MTok and $5/output MTok. Avoids the repo-prohibited 5.5 family. |
| OpenAI | <https://api.openai.com/v1/models> | `gpt-4.1-nano` | Low-cost, non-reasoning text model; model page lists $0.10/input MTok and $0.40/output MTok and Responses support. **Page marks it deprecated**; it was nevertheless discoverable and callable at this observation. |
| Gemini | <https://generativelanguage.googleapis.com/v1beta/models?pageSize=1000> | `models/gemini-2.5-flash-lite` | Listing declares `generateContent`; official model page describes a cost-efficient model for lightweight tasks. The page also recommends newer models for new applications. |

These are economical available selections, not a claim that they are the globally cheapest models or durable future defaults. No alternative models were probed after selection. Full sanitized ID listings and timestamps: [provider-discovery.json](ch02-entry-evidence/provider-discovery.json).

## Primary documentation verification

All direct origin fetches below returned **HTTP 200**. Redirect destinations, timestamps, and bounded relevant excerpts are preserved in [provider-docs.json](ch02-entry-evidence/provider-docs.json) and [provider-doc-details.json](ch02-entry-evidence/provider-doc-details.json). Excerpts are deliberately bounded, not complete page archives; a null excerpt means that exact text search did not match, not that a feature is absent.

| Official source URL | Verified content |
|---|---|
| <https://platform.claude.com/docs/en/api/messages/create> | `POST /v1/messages`, structured messages, `max_tokens`, API-key/version headers. |
| <https://developers.openai.com/api/reference/resources/responses/methods/create> | Native Responses POST, string `input`, minimum `max_output_tokens` of 16, Bearer authentication, `store: false` option. |
| <https://ai.google.dev/api/generate-content> | `models.generateContent`, v1beta model-path endpoint, generation configuration/output-token controls. |
| <https://ai.google.dev/gemini-api/docs/api-key> | API-key header authentication via `x-goog-api-key`. |
| <https://platform.claude.com/docs/en/about-claude/models/overview> → <https://platform.claude.com/docs/en/models/overview> | Current Claude overview; mentions Haiku 4.5 capabilities. Newer 5.5 examples were not used. |
| <https://platform.claude.com/docs/en/about-claude/pricing> | Haiku 4.5 economical pricing. |
| <https://developers.openai.com/api/docs/models/gpt-4.1-nano> | Model economics, lack of a reasoning step, Responses endpoint, deprecation notice. |
| <https://ai.google.dev/gemini-api/docs/models/gemini-2.5-flash-lite> | Model identity and lightweight/cost-efficient use; guidance toward newer models. |

The initial Firecrawl attempts for the three API references failed at the **tool transport/access layer**: anonymous/keyless access unavailable. They yielded **no origin HTTP status**, so they are not represented as origin HTTP failures or successful documentation reads. The same pages were then successfully fetched directly, without credentials. The parent was notified of the tool blocker and acknowledged the safe direct-doc continuation and selected model IDs.

## Actual minimal inference — original three-call pass

The original pass made **three generation requests, one per vendor**, with no retry or spend loop and no alternative endpoint/model probing. The subsequent explicit parent authorization added exactly one OpenAI Chat Completions request, bringing the final total to **four** (Anthropic 1, OpenAI 2 across distinct interfaces, Gemini 1). Each user prompt was exactly `Reply OK`. Each request capped output at 16 tokens, was non-streaming, used no tools or conversation history, and left optional thinking/reasoning at the selected model's default. OpenAI additionally used `store: false`; Gemini requested one candidate. No assertion of zero provider retention is implied.

| Vendor / interface | Selected model → returned model | Actual HTTP / completion | Exact bounded output | Reported tokens | Elapsed |
|---|---|---|---|---|---|
| Anthropic Messages — `POST https://api.anthropic.com/v1/messages` | `claude-haiku-4-5-20251001` → same | **200**, `end_turn` | `OK` (not truncated) | 9 input / 4 output | 0.841 s |
| OpenAI native Responses — `POST https://api.openai.com/v1/responses` | `gpt-4.1-nano` → `gpt-4.1-nano-2025-04-14` | **200**, `completed` | `OK` (not truncated) | 9 input / 2 output; 0 reasoning | 4.067 s |
| Gemini generateContent — `POST https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash-lite:generateContent` | `gemini-2.5-flash-lite` → same | **200**, `STOP` | `OK` (not truncated) | 3 prompt / 1 candidate | 0.440 s |

Generation observation interval: **20:22:01.287–20:22:06.639 UTC**. Each request used a generous 600-second socket timeout, not a short blanket generation deadline. Tool polling at three-second intervals was observation only; running jobs were collected to actual completion. No timeout occurred. The output artifact bound was 160 characters; actual outputs were two characters each.

Exact nonsecret request bodies, statuses, usage, completion states, selected/returned IDs, and timestamps: [provider-generation.json](ch02-entry-evidence/provider-generation.json). Usage is provider-reported; no billing invoice or actual charged amount was verified.

## Limits and next owner

- This receipt provides direct-provider setup evidence, **not independent source review**. No code correctness verdict is offered.
- Model listing and documentation retrieval are separate from the successful tiny inference evidence above. A tiny `OK` response does not test meaningful chat history, multi-turn behavior, streaming, tool calls, rate limits, larger budgets, persistence, error recovery, or model-family compatibility.
- These calls bypassed any student adapter/CLI. Student-owned commands and captured output must still establish all public Chapter 2 live-use and replay requirements. No acceptance gate was advanced by this task.
- Availability is account-, region-, quota-, and time-dependent. The selected OpenAI model is documented as deprecated; Gemini's page steers new applications toward newer models. Future users should rediscover supported IDs rather than treating this receipt as a permanent default.
- No grader/private expectations were consulted, no configuration was changed, and no credential values or provider account identifiers were retained in these artifacts.
- Finite task is complete: do not rerun these generation probes as part of this receipt. There are no pending generation jobs or vendor blockers.

## Provenance correction and evidence audit

This section corrects the initial receipt after the parent challenged its repository HEAD and description of the current public contract. **The provenance audit itself made no network requests or additional inference.** After the discrepancy was reconciled, parent messages `main-6` and `main-7` explicitly authorized one additional OpenAI Chat Completions request after primary-doc verification, expanding the original one-per-vendor limit only for that single call. This correction turn changed this Markdown receipt and added two separate Chat Completions evidence files; the four original JSON artifacts remain byte-for-byte unchanged.

### Repository and chapter source

- Working directory explicitly used for the original repository commands and the correction audit: `/Users/bill/projects/ensemble`.
- The original saved command output `/Users/bill/projects/cadence-memories/children/e3-ch02-provider-readiness/io/61`, line 1, records HEAD `5351583208988f8370937cd0494727645efa8b46`. The correction-time `git rev-parse HEAD` returns the same value. The initial receipt's `cd9ab4a9d407541c2f9c51f40f5fc4c0f2c62e09545d` was not supported by that output. **Root cause: an unchecked manual transcription in the receipt, not an observed repository revision change.** No evidence establishes where that erroneous string originated.
- Exact file read: `/Users/bill/projects/ensemble/book/edition-3/chapter-02.md`. Its SHA-256 is `1a67ee9e3b6ae5f0f776a30f58908c6d53b53be2d8136ef44661728a677dc2d2` (a SHA-256 digest, not SHA-1); size 71,337 bytes. The audit verified that the complete current file text is an exact substring of the original full chapter read in saved output `io/61`.
- The parent initially asserted that the current chapter explicitly teaches native Responses in §§2.2/2.6/2.9, but then independently rechecked the absolute file and HEAD blob and acknowledged that assertion was wrong (`main-6`, confirmed in `main-7`). The file returned to this agent by both direct Python/file reads and the file-reading tool, at the hash above, contains `OpenAI, its own message with a **tool** role:` at line 801; the OpenAI parse table uses `choices[0].message.content`, `.tool_calls[]`, and `finish_reason` at line 872; and the surface list contains `chat_completions` at line 1129. A literal search found zero occurrences of `native Responses`. The parent confirmed these observations and explicitly identified Chat Completions as the actual chapter interface. **The discrepancy is resolved; the chapter was not changed.**
- **Root cause of the contract-description error:** the initial receipt labeled observed wording “historical” without source evidence establishing that status or reconciling it with the requested native Responses check. That interpretive claim is withdrawn. The direct Responses request remains verified as extra parent-requested setup, not a chapter requirement. The parent owns the erroneous native-Responses instruction; this agent owns the erroneous HEAD transcription and unsupported historical framing. No chapter or implementation file was modified.

### JSON provenance verification

The audit parsed each original saved stdout JSON line and compared the complete resulting records structurally to its artifact, including URLs, timestamps, model IDs, statuses, output, usage, and request bodies where present. Every comparison matched exactly. None of these artifacts contains the erroneous repository HEAD or the receipt's historical-contract claim. They establish the original tool-run results, not an independent wire capture or source-code acceptance.

| Unchanged artifact | Original saved stdout under `.../children/e3-ch02-provider-readiness/io/` | SHA-256 |
|---|---|---|
| `ch02-entry-evidence/provider-discovery.json` | `71`: all 3 discovery records match | `b42a6339d5b84371eb538598f8565d8b5b895a2d87843eff7cb0b79e7282b041` |
| `ch02-entry-evidence/provider-generation.json` | `90`: all 3 generation records match | `7d3fdc7d5a7f5ecdc27451caacc148aeaa7e1f7bfa6071ba4a8fb41c4af0e47c` |
| `ch02-entry-evidence/provider-docs.json` | `81`: all 6 direct documentation records match; Firecrawl failures separately recorded in `73`–`75` | `868bf0c0596313b9e3d777ff36ca530a968592e943d15da95015070f4185a8af` |
| `ch02-entry-evidence/provider-doc-details.json` | `86`: all 4 additional documentation records match | `e2e38601ebdbeb391d4ee0b5c6f49d5ad3ad2ff44ce4e8bf4b3567d5cbf793f6` |

The original three successful HTTP 200 / `OK` results are preserved, with their actual selected/returned IDs and original timestamps. No replayed output is represented as a new request. This correction does not advance any student, implementation, or independent-review acceptance gate.

## Explicitly authorized actual-chapter OpenAI follow-up

After the local provenance audit, parent messages `main-6` and `main-7` acknowledged the parent's incorrect native-Responses instruction and authorized **exactly one additional** minimal OpenAI Chat Completions request. This explicit scope expansion superseded the correction turn's no-new-inference restriction only for this call. No additional discovery, other-vendor generation, retry, or model probing occurred.

Primary documentation was fetched directly without credentials at <https://developers.openai.com/api/reference/resources/chat/subresources/completions/methods/create>, **HTTP 200** at 20:29:30 UTC. It confirms `POST /chat/completions`, `messages`, Bearer authentication, and the `max_completion_tokens` output bound. See [provider-openai-chat-docs.json](ch02-entry-evidence/provider-openai-chat-docs.json).

| Item | Actual result |
|---|---|
| Interface | `POST https://api.openai.com/v1/chat/completions` — actual Chapter 2 OpenAI setup interface |
| Model | Already discovered `gpt-4.1-nano`; returned `gpt-4.1-nano-2025-04-14` |
| Credential reference | `directOpenAIAPIKey`, memory/header only; no value retained |
| Request | User prompt exactly `Reply OK`; `max_completion_tokens: 16`; `n: 1`; `stream: false`; `store: false` |
| Attempt and transport | Exactly 1 additional attempt; direct HTTPS, proxy/redirects disabled; 600-second socket timeout; no retry |
| HTTP/completion/output | **200**, `stop`, exactly `OK`, not truncated |
| Usage | 9 prompt / 1 completion token; 10 total; 0 reasoning tokens |
| Timing | 2026-10-16T20:29:46.953561+00:00 through 20:29:48.308579+00:00; 1.355 seconds |

Full sanitized record: [provider-openai-chat-generation.json](ch02-entry-evidence/provider-openai-chat-generation.json). Final generation total across both turns: **4 attempts** — Anthropic Messages 1, OpenAI Responses 1 (extra parent-requested setup), OpenAI Chat Completions 1 (actual chapter interface), Gemini generateContent 1. All four returned HTTP 200 and `OK`. This remains only provider setup evidence, not student CLI or independent implementation acceptance. No further probes are authorized or pending under this finite task.
