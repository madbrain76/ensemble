# Chapter 10 — remaining independent controls

## Verdict: PASS recommendation for the first whole-candidate code review

The five authorized remaining controls completed with expected exits **0, 0, 0, 1, 0**. Combined with the unchanged, verified prior source review, required same-CLI grades, ten module gates and valid partial controls, the previously identified review gaps are now closed. No supported candidate must-fix was found. **Root owns final inspection, join, code acceptance and any formal-round decision. Code acceptance does not establish operational exercise completion.**

Candidate remains frozen at `729bc919e734405307b4ac0060e65c40e33652e5`. No student, grader, catalog, public-guidance or manuscript edit occurred. No grade, lifetime, standard suite, historical compatibility or adapter control was rerun. No automatic further work is scheduled.

## Scope and provenance

Task: `/tmp/cadence-ch10-code-review-remaining-controls-task.txt`, 6,471 bytes, SHA256 `87b4595d2d5a7d9793d05c296375879751875e030a3fd9e7dac123cd202d025a`. Coordinator input: `docs/edition-3-notes/ch10-code-review-remaining-input.json`.

The task and latest partial report were read fully. Governing/source reads carried only against unchanged hashes. Direct checks verified all **39 inputs, 137 candidate files, 33 protected files, two frozen binaries**, the source TARs, repository identities and nine older inventories. The current contract remains `book/edition-3/chapter-10.md`: a non-null checkpoint is installed and only strictly newer events applied; null context requests full replay. No stronger zero-frontier policy, new Replay/compaction API, catalog fallback or relaxed admission rule was invented.

The retained archive-extracted copy was checked against the pinned TAR and Git executable classes before use. Its exact 137 original files/modes were verified again before/after the mutation, restored positive and cleanup. Only new private probe files and the explicitly authorized owned-copy mutant were written. A new runtime HOME/TMPDIR isolated this run from earlier local artifacts.

### Exact private changes

The public probe differs from the immutable partial probe only in these authorized values:

- `"Idle"` → `"idle"` at the two matching fixture/assertion sites;
- `gpt-5-course` → existing image-capability row `gpt-4.1-2025-04-14`;
- `gemini-2.5-pro` → existing capability row `gemini-3.8-flash`.

The emitted unredacted image is supported by these existing rows. The PDF tool result is projection-redacted before rendering. Rich data and assertions were retained. The GUI probe and atomic sentinel function remain byte-for-byte unchanged. `fixture-pins.json` and `public-values.diff.txt` record the exact proof.

The GUI inspection confirmed the four-argument `Run` returning `int`, `-port 0`, `-save`, `-no-stream`, explicit loopback base URL, and real wire fields. Although durable `turn_ended` capture only clears hub state, the live `ensemble.TurnEnded` observation explicitly emits the expected WebSocket message. No GUI adjustment was needed.

## Carried evidence — not rerun or rewritten

The corrected same-CLI Chapter 10 100/100, single fresh Chapter 1 100/100 with both actual 62-second quiet delays, and ten root/GUI/consumer test/vet gates remain carried through unchanged source, binary, command/output and inventory pins. The partial valid Anthropic replay modes, malformed-load preservation and valid persisted primary lengths 0/3/18 also remain carried.

The atomic baseline was specifically pinned to the partial receipt containing an actual permission-denied temporary-file creation error and `sentinel unchanged`. Its exact test function and `REVIEW_ORIGINAL_REGULAR_FILE\n` sentinel did not change. The partial command's overall exit 1 is preserved; only its actual atomic PASS is carried as the baseline.

The original incomplete report/inventory and partial report/inventory remain immutable. Their prior failures, incorrect assumptions and qualifications are not erased. The parent-verifier timestamp refresh/recovery record was read and hashed; no parent verifier writer was executed. See `carry.json` and the before/after direct verification receipts.

## Executed controls

Every command used `/opt/homebrew/bin/go test -count=1 -run SELECTOR -v -timeout=... .` on the verified owned copy. Exact cwd, complete allowlist, PID/process group, timing, deadline, exit and stdout/stderr identities are in `command-receipts.json`.

| # | Exact selector | Exit | Elapsed seconds | Result |
|---|---|---:|---:|---|
| 1 | `^TestReviewReplayModesAndRichRender$/anthropic$/zero-cached-plus-tail$` | 0 | 1.886254 | Actual selected RUN/PASS; valid non-null zero-frontier checkpoint retained and newer tail applied |
| 2 | `^TestReviewReplayModesAndRichRender$/(openai\|gemini)$` | 0 | 0.357480 | Both branches and all twelve intended replay subtests passed |
| 3 | `^TestReviewLegacyCollisionOverWebSocket$` | 0 | 0.543797 | Real GUI entry point, WebSocket replay/reconnect/restart and one fake-provider continuation passed |
| 4 | `^TestReviewAtomicSavePreservesFile$` | **1 expected** | 0.818705 | Compiled mutant reached `ATOMIC_SENTINEL_CORRUPTED`, not a compile/panic/setup failure |
| 5 | `^TestReviewAtomicSavePreservesFile$` | 0 | 0.258187 | All original bytes/modes restored; same focused sentinel positive actually ran and passed |

The first two selectors execute necessary parent setup/postlude, but no unrelated Anthropic child selector, primary/malformed/atomic baseline campaign, grade or standard suite was replayed. Actual RUN/PASS lines were required; no empty selection was credited. All five stderr files are empty. No timeout, cancellation, unexpected execution failure or behavioral retry occurred.

### Replay, typed context and native rendering

The completed OpenAI and Gemini branches cover checkpoint-only, null-context full replay at nonzero/zero frontiers, non-null zero-frontier plus tail, ignored retained prefix and strictly newer tail. They assert repeated context/loading stability, monotonic append frontier, cold native-render equality where applicable, JSONL replacement, detached views and unchanged raw event facts. The selected Anthropic zero-frontier case completes the gap left by the prior malformed state literal.

The synthetic history includes multipart text, owned real PNG reference, opaque/signature/provenance annotations, tool call/result, redaction/summary, ephemera and usage. Native renderer output is retained in stdout. This is real Agent/history/rendering behavior with synthetic data, **not real provider acceptance**; no paid or external provider request was made.

### GUI legacy identity and correlation

The real GUI entry point ran with owned stdin/EOF, loopback HTTP, ephemeral port and actual WebSockets; no browser was opened. Two old tool calls/results and multiple old parts shared raw display ID 1. Observed presentation IDs were:

```text
FIRST  = "archive:2:0"
SECOND = "archive:5:0"
call-A = "archive:2:1"
call-B = "archive:5:1"
new continuation = 3
```

The control verified distinct old text/tool identities, call/result ID association, identical replay identities on reconnect, a new continuation not reusing the old presentation identities, stability after a real saved restart, and byte-structural equality of the archived event prefix. Exactly one loopback fake-provider request was observed. These are protocol/runtime facts, not a browser visual claim or execution of serialized old handlers.

### Compile-valid sensitivity and exact restoration

The only behavioral mutation inserted this statement in owned `internal/agent/save.go`, before safe temporary-file creation:

```go
_ = os.WriteFile(path, nil, 0600) // REVIEW MUTANT: premature in-place truncation
```

The public Save then reached the real non-writable-parent failure path after corrupting the existing writable file. Actual output:

```text
ATOMIC_SENTINEL_CORRUPTED: got="" want="REVIEW_ORIGINAL_REGULAR_FILE\n"
Save error: .../.save-...: permission denied
```

The negative was compiled and failed the intended unchanged assertion. All 137 original bytes and archive modes were restored, not just the changed file. The same focused positive then reported actual Save permission failure with the sentinel unchanged. This closes the required **carried baseline green → compiled behavioral red → exact-restored focused green** chain. `mutation.diff.txt`, `mutation-pins.json`, `restoration.json` and commands 4/5 are the evidence; a clean hash alone was not credited as a restored runtime result.

## Environment, failures and cleanup

The complete allowlist records installed Go/Delve paths, owned HOME/TMPDIR, `GOMAXPROCS=2`, `GOFLAGS='-mod=readonly -p=1'`, `GOTOOLCHAIN=local`, `GOPROXY=off`, and `GOSUMDB=off`. No credential/provider environment or generic `LLM_BASE_URL` was inherited. The GUI's key is an explicit dummy scoped inside its test. Execution remained sequential and above the 1 GiB floor.

One **pre-control reviewer verifier error** is retained: an older inventory was initially addressed as `ch10-student-final-evidence` instead of the sealed receipt's `ch10-student-final-reconciliation-evidence`. Zero commands ran in that attempt. Only that read-only verifier pathname was corrected; probes, selectors and deadlines did not change. The original driver, traceback, zero-command receipts and recovery record are retained. This is neither a student failure nor an executed-control retry.

After all commands, no owned process/group remained. The two newly added reviewer files were hashed **durably before deletion**. The reused source extraction is again exactly the original 137 files and modes; the new TMPDIR is empty. No shared cache, old runtime HOME/TMPDIR or unrelated state was cleaned. The original student, binaries, archives, protected files, old reports and seals remain unchanged. No delegated children, browser, paid provider, real key, source commit or guidance change occurred.

The closed UTF-8/NUL-free inventory includes this report, excludes only itself, and contains no source TAR, binary, settings payload or private raw capture. Final process, preservation, restoration and cleanup facts are in `final-receipt.json`.

## Decision boundary

The previously unmet controls now close without weakening assertions or changing the candidate. My recommendation for this **same first whole-candidate code review is PASS**, based on the combined pinned source review and completed evidence—not merely the grade scores. Root must inspect this package and join the reviewer before deciding acceptance. No operational exercise completion is inferred, and no additional work is scheduled after submission.
