# Chapter 10 code-review missing-controls completion

## Verdict: qualified / partial — NOT ACCEPTED

This finite follow-up does **not** complete the whole-candidate independent code review. The first authorized public-control command compiled and ran, but exited 1 because three private-fixture inputs were unsuitable. It stopped the sequence. **No candidate defect was established.** No candidate change, model-catalog expansion, relaxed validator, or student correction brief is justified by these failures.

Some previously missing assertions now have valid positive evidence, including real persisted empty/whitespace primary instructions and preservation of an existing regular file after a real Save failure. The rich three-renderer/control group is not green; the independent GUI collision control and compile-valid mutation/restoration sequence did not run. They cannot be credited from source reading, prior grades, the original failed probes, or the unchanged-copy check.

Root owns inspection, join, whole-code acceptance and any subsequent formal round. The student remained frozen and idle. There is no automatic continuation.

## Scope, pins and carried work

- Task: `/tmp/cadence-ch10-code-review-completion-task.txt`, 6,971 bytes, SHA256 `fb8cfe55bd957f304d998179d35ed38645eb901dbb6fd3ccc434e825f4661e0e`.
- Input: `docs/edition-3-notes/ch10-code-review-completion-input.json`.
- Frozen candidate: unchanged `729bc919e734405307b4ac0060e65c40e33652e5`; all 137 tracked-file bytes, tree identity, frozen binaries and source archives verified before and after.
- The 35 input pins and 33 protected tracked-file pins were verified. Historical/governing full reads were carried only where unchanged; the new parent verification, original report/failures, correct public chapter, actual API declarations and relevant saved probe sources were read for this task.
- Original code-review report and inventory remain byte-for-byte unchanged. The original inventory is 116 members / 1,138,356 payload bytes, SHA256 `d781ef13cd6cdfafe8dcc32dcc50ffde64dda611d0faa7472bda0a53acc6e932`. Its report and every member were rehashed. All seven older carried inventories were also verified through the existing read-only verifier.
- The corrected same-CLI Chapter 10 100/100, single fresh Chapter 1 100/100 with both actual 62-second delays, and ten root/GUI/consumer module test/vet results are carried by their exact command/output and source pins. **None was rerun.** Original failures and qualifications remain preserved.

See `carried-controls-and-failures.json`, `verification-before.json`, `verification-after.json`, and the original-seal verification receipts in the new evidence directory. Carrying green gates is not whole-code acceptance.

## Correct public contract and corrections to the original review

The governing public source is **`book/edition-3/chapter-10.md`**. This corrects the original report's inaccurate `books/edition-3/10-persistence.md` citation without rewriting that sealed report.

The actual chapter says:

> If `context` is non-null, install it as the checkpoint and apply only events with `seq > as_of`. ... If `context` is `null`, rebuild the projection from the full `log` ...

It does **not** say that a zero frontier discards a non-null cached context. The earlier unexecuted private expectation to that effect was unsupported and was retired before this task's command. The new zero-frontier test was intended to retain a non-null base plus a newer hint, but its own capitalized state made that fixture invalid; that intended assertion remains unreached.

The same chapter requires staging and validation before replacing live history, malformed/truncated saves to fail rather than become a new conversation, atomic replacement rather than in-place truncation, current active configuration rather than historical metadata overwrite, and frozen primary preservation without falling back merely because a selected instruction body is empty.

The actual public API supplies `History.ReadSave`, JSONL `Dump`/`Load`, `Context`, `Append(Event) error`, `Engine.Render(*Context)`, `Agent.Save`/`Load`, and projection-only `Redact` events. This review adds no invented `Replay` or compaction API requirement.

## Owned extraction and fixed control plan

A fresh owned directory was created and populated **from the pinned source TAR**, not from the live checkout. All 164 TAR members were checked to be regular files or directories, with no absolute/traversing paths or links. The 137 regular-file names and bytes matched the manifest. Archive file modes were preserved; final verification also confirms Git executable-class agreement. No archive was exported into this evidence package.

This resolves the previous checkout-copy method deviation for the new work only; it does not erase that deviation from the original attempt.

The single versioned correction made these private-harness changes before execution:

1. Normal Agents receive nonblank valid bootstrap instructions, `MaxTokens=128`, explicit dummy credentials, loopback-only HTTP guards, and isolated Workspace/DataDir/SavePath.
2. Empty/whitespace instructions use the existing primary-skill selection path after valid bootstrap configuration, then real save/cold startup/render operations. They do not request an invalid ordinary `NewAgent(System="")` constructor.
3. The zero-frontier expectation follows the actual non-null checkpoint rule and avoids imposing nil-versus-empty collection representation.
4. Media references point to an owned actual PNG and synthetic PDF. No real provider acceptance is claimed.
5. The unrun GUI harness uses the actual `gui` package, four-argument `Run` returning `int`, verified flags and WebSocket fields, and explicit loopback provider arguments rather than generic `LLM_BASE_URL`.

The full before/after probe diffs and exact final probe bytes are retained. No probe was edited after the command returned. The originally planned four commands were: public controls; GUI legacy/reconnect/restart/new-output control; one premature-truncation mutant against the regular-file sentinel; exact restoration and that same focused positive. Commands 2–4 were conditional on the prior positives and were not reached.

## Actual executed command and environment

Exactly one control command executed:

```text
/opt/homebrew/bin/go test -count=1 -run ^TestReview -v -timeout=90s .
```

It ran in the fresh archive-extracted root module. Exact absolute cwd, PID/process group, complete environment, timestamps, duration, deadline, stdout/stderr byte counts and SHA256 values are in `command-receipts.json`.

- Exit: **1**.
- Duration: **2.583182 seconds**; test output reports 0.203 seconds of test execution.
- Administrative deadline: 150 seconds; Go test deadline: 90 seconds.
- No timeout, signal, cancellation or retry.
- The wrapper caught the failure and returned 0 administratively; that is **not** a control PASS.

The allowlist includes `GOMAXPROCS=2`, `GOFLAGS='-mod=readonly -p=1'`, `GOTOOLCHAIN=local`, `GOPROXY=off`, `GOSUMDB=off`, installed Go/Delve paths and owned HOME/TMPDIR. No credential/provider environment or generic `LLM_BASE_URL` was inherited. Test configuration uses explicitly fake keys and local guards only. No unexpected provider request was reported. The captured payload is the real native Anthropic renderer output from synthetic context, not a remote-provider response or paid-provider acceptance claim.

## Failure-by-failure classification

### 1. Zero-frontier/non-null Anthropic checkpoint — invalid state literal

**Saved observation:** `invalid context turn`, at the private `reviewRead` call for `anthropic/zero-cached-plus-tail`.

**Private input:** the draft assigned `expected.Turn = "Idle"` and serialized that capitalized value into the cache. The actual canonical output generated by the candidate says `"turn":"idle"`.

**Source evidence:** `internal/common/common.go:486–490` declares `type TurnState string` and `Idle TurnState = "idle"`. `internal/history/save.go` validates the cached turn against the declared state domain before installation.

**Contract assessment:** a save must be valid before live state is replaced. Rejecting an unsupported state is required validation, not a demonstrated violation of checkpoint replay. The valid zero-frontier/non-null replay expectation was never reached. The separate zero-frontier/**null**-context test did pass.

**Disposition:** reviewer-fixture error, uncredited. No case-normalization, alias acceptance, validator weakening, or candidate fix is requested.

### 2. OpenAI rich renderer — undeclared image-capability alias

**Saved observation:** `model "gpt-5-course" does not support image media "image/png"` while preparing the canonical rich rendering, before cold-render comparisons.

**Source evidence:** `internal/engine/models.go` explicitly separates wire support from model capability and has no capability fallback. Its OpenAI image-capability row is `gpt-4.1-2025-04-14`, not the private draft alias. `internal/engine/render.go:136–137` rejects a media kind absent from the selected model's capability bits.

**Contract assessment:** the persistence chapter says to keep typed native requests and existing multimedia admission behavior. It does not authorize turning an arbitrary text/fake model alias into an image model. This failure does not establish loss of saved rich content or a wrong cold replay.

**Disposition:** reviewer model/media fixture error, uncredited. No model catalog change or admission-policy relaxation is requested.

### 3. Gemini rich renderer — undeclared image-capability alias

**Saved observation:** `model "gemini-2.5-pro" does not support image media "image/png"`, at the same canonical rendering stage.

**Source evidence:** the explicit matching Gemini capability row in `internal/engine/models.go` is `gemini-3.8-flash`. The same `render.go` capability guard applies. The private chosen alias did not carry declared image support.

**Contract assessment/disposition:** the same boundary as OpenAI: expected capability enforcement for an unsuitable fixture, not a proven persistence defect. Cold Gemini comparisons were not reached. No candidate/catalog edit is requested.

The original report's unsupported zero-cache hypothesis and these model/state errors must not be converted into a student must-fix. The task and final root steering required stopping, classification from saved evidence/source only, and sealing; no greener rerun was attempted.

## Valid assertions reached and their limits

### Explicitly green Anthropic replay subtests

Five subtests have actual `PASS` lines:

- checkpoint-only archive with cached context and no log;
- nonzero frontier with null context, rebuilding the full log;
- zero frontier with null context, rebuilding the full log;
- retained prefix ignored beneath the checkpoint frontier;
- strictly newer tail applied once, repeated context access stable, and next append beyond the frontier.

For the applicable modes the probe also compared native cold rendering to the live canonical renderer output. The synthetic context contains multipart text, a real owned image reference, opaque thinking/signature/provenance data, a tool call/result, projection-only result redaction and summary, usage and ephemera. Raw-log-prefix preservation is checked. The Anthropic branch continued after its failed child to the JSONL replacement and detached-view checks without reporting another failure; nevertheless the parent branch and whole three-renderer suite are **not green**.

This is not proof of OpenAI/Gemini rich replay or real provider acceptance.

### Persisted empty/whitespace primary and malformed-load preservation

`TestReviewMalformedLoadAtomicityAndFrozenPrompt` and all three children passed: primary lengths **0**, **3** (`" \t\n"`), and **18** (literal text plus LF).

Each starts with a valid bootstrap configuration and selects the exact body through the existing primary-skill interface. The checks exercise real save, public history loading, malformed/trailing input and nonexistent Agent.Load refusal without prior-state replacement, frozen saved primary bytes, cold startup against the saved archive, native rendering without bootstrap fallback, and a separately valid current runtime configuration retaining its own instruction string rather than being overwritten by recorded metadata.

These are valid observations replacing the old invalid-constructor attempt, not a demand to accept invalid ordinary constructors.

### Real existing-file Save failure

`TestReviewAtomicSavePreservesFile` passed and reached the actual filesystem failure path. An owned regular file contained `REVIEW_ORIGINAL_REGULAR_FILE\n`; its parent directory was made non-writable while the existing file remained writable. Public `Agent.Save` returned an actual `.save-*` temporary-file creation permission error, and the existing regular-file bytes remained unchanged. The parent permissions were restored and ownership joined.

This is a valid positive baseline for the intended sensitivity. **It is not a sensitivity result.** No early-truncation mutant was applied or compiled, and no restored runtime control ran.

## Explicitly unmet obligations

- Valid non-null zero-frontier replay assertion: **unreached**, because the fixture's state was invalid.
- OpenAI and Gemini rich cold-render comparisons: **unreached**, because their model/media fixtures were unsuitable.
- Real independent GUI legacy collisions, call/result pairing, reconnect/restart and new generation: **NOT RUN**.
- Compile-valid premature-truncation mutation, behavior-specific negative and focused restored positive: **NOT RUN**.
- Whole-candidate code acceptance: **NOT ESTABLISHED**.

The original baseline and both original failed probe attempts remain immutable. The valid partial results here do not retroactively repair or complete those attempts.

## Failure, carry and cleanup ledger

The new evidence retains the original report's incorrect citation as a standalone correction, the final probe diffs, private input fixtures, complete command receipt/stdout/stderr, native synthetic context/render output, failure classification, archive extraction pins and provenance, all carried command hashes and old failures, and final preservation receipts.

Pre-execution read-only lookup mistakes are also disclosed in `pre-control-ledger.json`: the first saved GUI filename was wrong, and a subsequent `gui` name filter missed the actual `reviewer_legacy_ws_test.go.txt`. No probe or source modification occurred in those failed lookups.

After the single exit-1 command:

- The Go job and sampler finished; the recorded PID/process group has no remaining owned member.
- No production mutation was applied.
- All 137 extracted original files match archive bytes and modes, with Git executable-class agreement.
- Exactly two reviewer-added test files were recorded by hash **durably before cleanup**, then removed from the owned extraction. Its exact regular-file set is again 137.
- The owned TMPDIR is empty. `cleanup-before.json` and `cleanup-after.json` retain exact identities and scope.
- Source archive, frozen student, binaries, public guidance, grader, protected files, original report/evidence and old seals are unchanged.
- No shared cache/unrelated cleanup, delegated child, browser, paid provider, real key, repeated grade, standard-suite replay or additional campaign occurred.
- Headroom remained above 1 GiB; after cleanup it was about 2.42 GB.

Transient executable/process observations are explicitly non-exhaustive. Local archives/binaries and temporary raw artifacts are not exported. The new exact UTF-8/no-NUL inventory includes this report and excludes only itself. No settings.json payload, binary or TAR is part of the portable package.

## Root decision boundary

Keep acceptance withheld on this evidence. The remaining review obligations are incomplete because of reviewer fixtures and the required stop rule, not because this run proved a candidate defect. Any subsequent completion or formal round needs a new explicit root decision. No source fix, assertion weakening, model substitution, probe repair, or rerun is performed or scheduled here.
