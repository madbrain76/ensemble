# Chapter 2 restart — corrected early code review

## Verdict and correction

**HOLD. Recommend a fresh public-only attempt from accepted Chapter 1 (`c5c7be07`), subject to the parent’s decision.** The confirmed failures span the shared event model, reducer, accounting, replay configuration, and CLI—not one localized bug. They concern explicit public requirements, not additional scope.

**My earlier source descriptions and several citations were wrong. Source drift is not established.** This note replaces those claims. Current code has **six** event types, block-level redaction (`Drop`, `Stub`, `Mask`, `Truncate`), and `Materialize` attaches pending blocks to the last user message or appends a message—not the system prompt. `Engine.Send` returns `(ParseResult, error)` and `Vendor.Render` returns `(Request, error)`; there is no `ValidateResponse` function. Response-supplied model IDs are preserved. Original common files are `types.go` and `enums.go`, not `event.go`.

## Confirmed blockers

Citations below are relative to the rejected candidate’s `main` directory. The necessary files were reread directly by absolute path at `solutions/edition-3/main` before root restored that location to accepted Chapter 1. Root reports the exact rejected candidate preserved unchanged at `/var/folders/p_/11vk6y114xq7k9g5r_pc3yv80000gn/T/ensemble-ch02-rejected-attempt1-4ifxgihv/main`; that is now the source path for this review. The restoration/preservation was root action, not a student hold violation. The restored Chapter 1 tree was not rereviewed as the rejected candidate.

1. **Alternate event/JSONL contract.** `internal/common/common.go:111–118` declares `MessageAdded`, `EphemeralAdded`, `EphemeralConsumed`, `Redaction`, `UsageObserved`, and `Compaction`. Its `Event` declaration (`:145–152`) has `ID`, `Type`, `Content`, `Redaction`, and `Usage`, rather than the public sequence/time and causal payload envelopes. `internal/history/history.go:27–58` accepts this alternate vocabulary and otherwise returns `"unknown event type"`. Public facts such as `message_received`, `request_sent`, and `response_ended` therefore cannot be loaded as specified.

2. **Missing turn-state and historical-usage projection.** `common.go:100–106` declares Context with only `System`, `Messages`, `Tools`, and `PendingEphemeral`. `history.go:103–173` folds messages, ephemera, and block edits but not usage or the required turn state. `internal/agent/agent.go:96–98` separately delegates `Usage()` to Engine and `Load()` to History. The live/historical separation is appropriate, but the historical usage projection and turn-state reducer are absent.

3. **Missing four-bucket accounting.** The actual declaration is `type Usage struct { Input int64; Output int64 }` (`common.go:27–31`); `UsageObservation` also has only input/output measurements plus provenance (`:138–143`). Engine accumulation is exclusively `e.usage.Input += ...` and `e.usage.Output += ...` (`internal/engine/engine.go:97–101`). There is no representation or accumulation of the required disjoint cache-write/cache-read categories.

4. **CLI does not implement the public file-based protocol.** `cmd/ch01/main.go:45–50` recognizes only one-argument `chat` or `render`, otherwise returning `"usage: ch01 [chat|render]"`; thus `render LOG` and standalone `dump` are rejected. The full CLI has no `CH02_LOG` persistence path. Instead, `stateCommand` loads an event array from stdin and dumps an `{"events":...}` wrapper (`:123–150`). Its ephemeral branch decodes `ensemble.Content` and returns `a.Apply(...)` without an acknowledgement (`:131–136`), rather than accepting the public string directive and emitting the required acknowledgement.

5. **Resolved replay target is not separate from routing.** `Config` has only one model field (`common.go:10–19`), and CLI configuration reads the provider routing model but not `LLM_RESOLVED_MODEL` (`cmd/ch01/main.go:27–44`). Each vendor’s `provenance()` uses `Model: v.parent.Parent().Config().Model`: Anthropic `internal/anthropic/anthropic.go:17–19`, OpenAI `internal/openai/openai.go:16–18`, Google `internal/google/google.go:18–20`. Renderers compare replay material against that tuple. Parsers **do** preserve response-supplied model IDs (Anthropic `:133–135`, OpenAI `:153–155`, Google `:134–136`); the blocker is the absent independent target identity, not an Engine overwrite.

## Architecture and comparative judgment

The mandatory star is a positive: root composition wires Agent, History, Engine, and the three vendors (`ensemble.go:48–66`); Engine owns HTTP, and shared interfaces expose the actual ownership links (`common.go:167–207`). Comments explain these boundaries. `Materialize` and `Replayable` are concrete cross-renderer reuse (`common.go:233–255`), not evidence by themselves of unjustified machinery.

Direct comparison with original `solutions/ch02` supports the contract findings: `internal/common/types.go:23–28,133–161` represents four usage categories, causal event envelopes, and a Context with turn and historical usage. That reference is comparison material, not replacement code to give the student. Current richer provenance/wire requirements justify additional implementation, and removing the obsolete blanket timeout is appropriate (`internal/engine/engine.go:22–28`). **No independent material-bloat finding or line-quota objection is established.** The fresh-attempt recommendation rests on foundational public-contract mismatch.

**Re-execution TLDR:** execute the documented causal-log/context/accounting/CLI contract instead of substituting a smaller protocol; retain the accepted ownership architecture. No reference-derived patch is supplied.

## Pending validation

No builds, tests, grades, providers, or browser checks were run for this review or correction. Compilation, Chapter 2 tests, grades, and live/GUI validation remain pending. Only this note was revised. Retain the same reviewer context for any later changed-area/final-results check; no second review campaign is proposed.

## Root decision

The root directly read the candidate and original Chapter 2 implementation and confirms the five corrected contract blockers above. The first attempt is rejected before costly validation. This is correction round 1 of at most 3, addressed by a fresh public-only execution from accepted Chapter 1 after the small TL;DR clarification. It is not a claim that the chapter body was ambiguous or that smaller code is necessarily better. The earlier conflicting root source descriptions were not sufficient evidence of a student hold violation. The rejected source is retained locally at the path above; the evolving workspace has been restored to the nine accepted Chapter 1 files. The original Chapter 1 snapshot is unchanged. No Chapter 2 tests, live exercises, or acceptance are claimed for the rejected attempt.
