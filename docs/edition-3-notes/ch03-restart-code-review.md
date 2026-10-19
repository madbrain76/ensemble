# Chapter 3 restart — independent code review

**Final recommendation: accept Chapter 3, subject to root’s acceptance decision. No open source-review blocker remains.** All six early findings and three formal correction rounds are closed; final frozen-candidate gates and the resolved-identity Gemini live exercise have been checked. Earlier sections below preserve checkpoint findings and qualifications; the closing section gives the final evidence scope.

## Scope and evidence

The initial review inspected the held `solutions/edition-3/main` candidate over accepted Chapter 2 `a1b08e12` / `solutions/edition-3/ch02`, current Chapter 3, and permitted original Chapter 3 comparisons. It recommended six bounded corrections before expensive validation. The student’s initial root/GUI build success was attributed, not independently established.

Round-one follow-up inspected only `internal/tool/tool.go`, the related README changes, new `chapter03_tools_test.go`, the inherited `chapter02_round2_test.go` collector adjustment, and these actual recorded outputs under `/tmp/e3-ch03-early.jR0AOc/round1/`:

- `focused-initial.log`: all five new round-one tool tests passed; the inherited ID test failed in its two OpenAI cases because its whole-request walk counted the six new `type:function` declarations as calls with nil IDs.
- `focused-corrected.log`: the five new tool tests and all six vendor/missing-ID cases of `TestRound2CallIDsAreStableRenderOnly` passed (`ok ensemble 0.019s`).
- `build.log`: records core, CLI and GUI build commands and the successful round-one checkpoint.

The reviewer read these outputs and the assertions, but ran **no builds, tests, graders, network or live exercise**. Focused success is not evidence that the full suite or remaining gates passed. Source locations below are candidate-relative unless prefixed otherwise.

## Later public requirement — edit diagnostics

**Approved; no open finding.** This is Bill’s new public requirement, not correction round two. The added paragraph at `book/edition-3/chapter-03.md:426–432` requires zero-count/reread guidance, truthful overlapping-match counts, bounded one-based match-start locations with truncation disclosure, and unchanged bytes on refusal without fuzzy selection.

The diagnostic delta at `internal/tool/tool.go:447–479` counts all exact occurrences by advancing one byte after each match, retains at most eight start-line numbers, and counts newlines strictly before each match start. It therefore handles same-line, overlapping and multiline starts without conflating locations or truncating the total count. Both refusal branches precede the write; the unique-edit/deletion path remains exact. The model-facing declaration (`:553`) and README (`:104–113`) describe the same policy. This small local change introduces no new ownership mechanism or substantial needless complexity.

Inspected `chapter03_tools_test.go:274–339`: eight refusal cases cover zero, same-line repetition, overlap, separate lines, overlapping multiline anchors, newline anchors and both ordinary/overlapping location-list truncation. Assertions check actual counts/locations, truncation only when appropriate, reread guidance, failure status and unchanged file bytes. Three success cases preserve unique multiline replacement, deletion and whole-file deletion.

Read `/tmp/e3-ch03-early.jR0AOc/edit-diagnostics/focused.log`: both new test groups and the existing registry/write-handler focused checks pass (`ok ensemble 0.012s`). Its `build.log` records successful core, CLI and GUI builds. No reviewer rerun or full-suite inference was made. No regression reopening an original finding was found in this limited delta; those six remain closed. No extra campaign is requested.

## Correction round two — refusal identity and local gates

**Approved; ready for bounded live validation.** `internal/tool/tool.go:472,479` now includes the actual supplied path, quoted, in both zero-match and ambiguous-match refusals. `chapter03_tools_test.go:302–304` asserts that identity. Counts, bounded start-line locations, omission disclosure, unchanged-byte checks, exact replacement and no-anchor-excerpt assertions remain intact. No new complexity or regression was found in this tiny delta.

The retained `/tmp/e3-ch03-early.jR0AOc/final/grade-edition2-ch3-cache-corrected.log` records the prior 95/100, `passed:false`, exit 1: the refusal named neither file nor anchor. Round two supplies file identity without echoing the anchor. The separate newly requested diagnostic feature above is **not** a correction round; the correction count at this checkpoint was **two**.

Read actual outputs under `/tmp/e3-ch03-early.jR0AOc/round2/`:

- `focused.log` preserves the intermediate compile failure (`a.str` undefined); current source uses the already validated `path`, not that nonexistent helper.
- `focused-corrected.log` passes the registered write-handler, all eight diagnostic and all three unique-edit/deletion cases (`ok ensemble 0.011s`).
- `root-race.log` passes root packages with tests; `root-vet.log` and `frozen-build.log` contain no diagnostics. Successful command exits for the latter two are reported by the parent, not independently encoded in those empty files.
- `grade-edition2-ch3.log`: **100/100, `passed:true`**, including recovered edit-contract credit and Chapter 2 parity.
- `grade-edition3-ch1.log`: **100/100, `passed:true`**, including actual **1m2s header delay and 1m2s body delay**, complete answers and usage. Parent reports both grader commands exited 0; their JSON outputs substantiate the scores/check outcomes.

No commands were rerun by this reviewer. GUI gates are carried forward unchanged as compile/build-only evidence, not an interactive GUI exercise. Original findings remain closed; live evidence and final root acceptance are still pending.

## Correction round three — Gemini declaration projection

**Approved for full gates; no open finding in this delta.** The actual `/tmp/e3-ch03-early.jR0AOc/gemini-diagnostic/run.log` preserves HTTP 400 / `INVALID_ARGUMENT`: `additionalProperties` is unsupported at each of the six function declarations’ `parameters`. The diagnostic used one reserved POST and ended with CLI exit 1. Prior aggregate passing checks did not establish actual Gemini schema compatibility.

`internal/engine/tools.go:41–61` now unmarshals each current parameter object into a fresh map, removes only its top-level `additionalProperties`, and serializes the projection. All six current schemas have exactly this object-level keyword and primitive property schemas. Common declarations remain intact; Anthropic/OpenAI functions are unchanged. `internal/engine/render.go:410–413` propagates projection errors instead of silently omitting declarations. This is a small Engine-owned dialect adjustment, not a general recursive schema framework. Handler-side unknown-field/type validation remains strict (`internal/tool/tool.go:115–147` and existing handlers).

Inspected `chapter03_schema_test.go`: it verifies six Gemini declarations, unchanged names/descriptions and every other schema field, retained `additionalProperties:false` in common definitions, no common mutation, byte-identical Anthropic/OpenAI schemas and repeat output, and rejection of malformed local arguments without file/shell effects. Read `/tmp/e3-ch03-early.jR0AOc/round3/focused.log`: projection and registry tests pass (`ok ensemble 0.010s`). Core/CLI/GUI build logs contain no diagnostics; success is reported by the submitting run, not inferred solely from empty files. No reviewer test/build rerun occurred.

The specific observed schema finding is closed at source/focused-review level. Full gates and the next bounded Gemini run are still pending; **no claim of live Gemini compatibility** is made. Formal correction count is now **three**; the separate user-added diagnostic requirement remains outside that count. Original six closures are unchanged.

## Six findings — round-one disposition

The locations in this section refer to the round-one review checkpoint, before the later diagnostic insertion shifted subsequent tool-file lines.

1. **CLOSED — provider credentials inherited by shell.** `internal/tool/tool.go:482–496` now removes exactly `LLM_API_KEY`, `ANTHROPIC_API_KEY`, `OPENAI_API_KEY` and `GEMINI_API_KEY` and preserves other environment entries. It explicitly avoids nil-environment inheritance even when filtering leaves no entries. This meets the public requirement at `book/edition-3/chapter-03.md:333–338`. The actual-child-shell regression (`chapter03_tools_test.go:172–207`) checks all four variables are unset using synthetic values, while an ordinary build setting and `GOFLAGS` remain available. No real-key exposure was alleged or tested.

2. **CLOSED — missing read cap/truncation notice.** `internal/tool/tool.go:195–257,523` accepts `max_bytes`, defaults omitted/zero to 32,768, selects the inclusive line range before capping returned bytes, and adds an explicit marker only when actual selected data exceeds the cap. Bounded fragments avoid loading an entire long line merely to truncate it. This addresses the built-in table (`chapter-03.md:318`). Focused assertions (`chapter03_tools_test.go:122–170`) cover range-before-cap after a long first line, exact-size/no-marker results, default/explicit caps, EOF, empty files and malformed options. README wording matches the implementation, including the marker outside the byte cap.

3. **CLOSED — merged stdout/stderr.** `internal/tool/tool.go:497–517,528` captures the two streams separately and returns JSON text with `stdout`, `stderr` and `exit_code`, as required at `chapter-03.md:323`. The focused real-shell test verifies distinct outputs plus exit 7, ordinary exit 127, and startup failure. Execution remains synchronous, with caller cancellation and no arbitrary deadline.

4. **CLOSED — independent declaration/handler catalogs.** `internal/tool/tool.go:22–34,36–78,87–110,520–533` now has per-owner entries joining each declaration to its executable closure. Both declaration exposure and dispatch derive from these entries; the parallel handler-name switch/enabled map is gone. This addresses §3.1 (`chapter-03.md:107–143`) without a new registration framework. Owner-local private state, shared common interfaces and root wiring are preserved. Focused checks (`chapter03_tools_test.go:50–120`) exercise all six declarations, explicit-empty behavior on all providers, disabled dispatch/no write effect, defensive copies, Agent isolation, and invalid selections. The other focused tests exercise the registered handlers.

5. **CLOSED — missing context-group separator across files.** `internal/tool/tool.go:351–354` bases separation on prior output rather than a per-file span index, satisfying the built-in table (`chapter-03.md:320`). Exact-output assertions (`chapter03_tools_test.go:209–237`) cover overlapping and separated windows, a cross-file separator, zero-context output and adjacent-window merging.

6. **CLOSED — overwrite refusal omitted size.** `internal/tool/tool.go:390–411` retains exclusive default creation, reports existing-file size and opt-in flags, and handles a failed diagnostic stat without claiming a size. This addresses the described refusal at `chapter-03.md:288–297`, especially `:291–292`. Focused assertions (`chapter03_tools_test.go:239–271`) check the diagnostic and unchanged bytes, append, overwrite, invalid options, and the registered list/edit paths.

## Inherited fixture alignment

The only inherited-test change narrows its recursive collector to actual `messages` (Anthropic/OpenAI) or `contents` (Gemini), excluding the separate declaration catalog (`chapter02_round2_test.go:281–288`). This is legitimate Chapter 3 fixture alignment, not a lowered ID bar. The test still checks repeat-render stability and unchanged History (`:239–242`), exactly two nonempty unique IDs (`:289–290`), retention of supplied IDs, and same-name result pairing by ID (`:292–300`), across all three providers and both supplied/missing-ID cases. The initial failure and corrected result are both retained in the reviewed logs.

## Architectural and reference judgment

The initial whole-source review found appropriate reuse of the accepted provider transport, canonical parts, History reducer, renderers and usage path. Root wiring and owner-local `privateState` remained intact. The loop appended provider events before ordered calls/results; failed later calls did not erase prior results. Explanation gating was an explicit opt-in, whole-batch refusal produced paired failed results, and the sixteenth permitted batch’s actual results received a final consuming request before further dispatch was refused. These initial observations were not revalidated by full execution in this follow-up. Environment policy remains default-off unless `LLM_REQUIRE_TOOL_EXPLANATION` is explicitly true (`config.go:46–52`); that is an allowed opt-in, not a defect.

Direct original comparison remains instructive: `solutions/ch03/internal/tools/tools.go` already joined registry entries (`:20–73`), bounded reads (`:93–137`), supplied useful overwrite refusals (`:188–194`), separated shell streams/filtered credentials (`:251–289`), and separated search groups across files (`:292–362`). The initial omissions were not justified migration costs. Round one restores those public behaviors while retaining useful cancellation, exclusive creation, strict option checks and overlap-aware editing. Its bounded-fragment reader has a concrete long-line/bounded-memory purpose. No substantial needless complexity requiring another redesign was found in these corrections; no LOC target applies.

**Next:** proceed with full required gates for round three, then the authorized bounded Gemini run and remaining chapter review; return for the final results check in this same context. No additional speculative campaign is requested. Final acceptance remains with root.

## Post-round-three live qualification — launcher identity

The final-CLI Gemini attempt is **not a live pass**: `/tmp/e3-ch03-early.jR0AOc/live-gemini-final/run.log` records three POSTs, HTTP **200, 200, 400**, and CLI exit 1. Inspection supports root’s narrow launcher-identity diagnosis, not a fourth production correction. `run_once.py:231` sets `LLM_MODEL` but omits `LLM_RESOLVED_MODEL`; its CLI is the frozen `round3-final/ch03-final`. `config.go:35` reads resolved identity separately. Engine correctly requires exact nonempty target/producer identity before replaying extras (`internal/engine/render.go:12–15,350–375`), rather than trusting a routing alias.

Checked saved `gemini/control/session.jsonl` metadata without exposing payloads: producer is exactly `gemini` / `gemini-3.8-flash` / `generate_content`. Both the ordinary-text signature and the first parallel `read_file` signature survive in History, consistent with parser `:275–283` and reducer `internal/history/history.go:246–247`.

Compared root’s paired offline outputs `/tmp/ch03-signature-diagnosis-kpjg54x7/{unknown,known}.json`: unknown-target output equals `failed-request.json` **structurally**, not byte-for-byte. Known-target output adds only `thoughtSignature` at `contents[1].parts[0]` and `contents[3].parts[0]`; both values exactly equal those saved in History. There are **zero other structural differences**. This substantiates supplying the observed `LLM_RESOLVED_MODEL=gemini-3.8-flash` in the recovery launcher without weakening production provenance rules or changing source. Formal production correction count remains **three**.

No reviewer replay, provider request, test/grade campaign or production edit was performed. This qualification supports the authorized bounded recovery; its live outcome and the frozen-final grade records still require the final results check. It does not retroactively pass the failed attempt or establish live compatibility.

## Final results check and recommendation

**Review complete: recommend root acceptance; no open source-review finding.** Read the final student note and actual records under `/tmp/e3-ch03-early.jR0AOc/round3-final/`: `root-race.log` passes; `grade-ch3.log` is **100/100, passed:true**, including all 13 Chapter 2 checks; `grade-ch1.log` is **100/100, passed:true**, including separate real 62-second header/body delays. The student reports exit 0 for these commands and root/GUI vet and final build. Empty vet/build logs supply no contrary diagnostics but do not themselves encode exit status. `gui-test.log` explicitly has **no test files**: GUI evidence remains compile/build-only, not interactive GUI behavior. These records concern the same frozen `ch03-final`; no gates were replayed for the later harness-only identity change.

Checked `live-gemini-resolved/run.log` and sanitized saved conversation/tool results: **six POSTs, all HTTP 200, CLI exit 0, 13.68 seconds, three completed human turns, four ordered call/result pairs**. The first reads actually return the `helo` defect and the test. The edit reports one exact replacement; current source returns `hello`; the saved test bytes equal the final test. Actual `run_command` result contains passing `go test ./...` output, empty stderr and exit 0. A distinct third response recalls the codename without tools. The workspace contains only its three exercise files. This establishes success for this bounded Gemini exercise with explicit observed producer identity—not universal provider/schema compatibility.

Read earlier `live/run.log` and `live-v2/run.log`, and the clean Claude/OpenAI final responses. Original Claude exercised all six tools and a real test failure→repair→pass, but hit its ten-POST cap; exposed prompt controls and a false-positive recall summary disqualify its recall claim. The separate clean four-POST Claude session demonstrates actual third-turn recall. OpenAI’s six-POST session records repair, unchanged test, successful process result and actual third-turn recall. **Those Claude/OpenAI runs used round-two bytes and are carried on unchanged provider paths, not fresh final-binary live proof.** Initial Gemini schema rejection and later unknown-identity signature rejection remain preserved failures.

Total paid attempts: **31 POSTs (14 Claude, 6 OpenAI, 11 Gemini)**. Formal production correction count is **three**; the added user diagnostic requirement and final launcher identity correction are not additional rounds. No reviewer tests, builds, grades or paid calls were run.

The original Chapter 3 comparison above still supports the final KISS judgment: joined declarations/handlers and owner-local state reuse the existing star architecture; provider projection stays in Engine, shared types remain interfaces/data, and root wiring remains thin. The bounded diagnostic and Gemini projection are justified local additions, not frameworks. No substantial needless complexity remains flagged. **Final acceptance and chapter closure belong to root.**

## Root decision

Accepted after direct source/reference comparison, final gates, and qualified live-result inspection. Three formal correction rounds were used; the user-requested diagnostic addition and proven launcher identity correction were separate. The 27 source/support files are snapshotted in `solutions/edition-3/ch03`. Text closure follows; no additional code or provider campaign is required.
