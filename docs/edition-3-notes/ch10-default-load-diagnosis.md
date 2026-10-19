# Chapter 10 default-load diagnosis — frozen candidate 729bc919

## Bounded finding and decision boundary

**The measured mismatch is one terminal LF in the active primary instruction, not restoration of the saved model/vendor/prompt/tools.** Both frozen clients send the requested current model, the current primary body including its terminal LF, and the current three-tool capability. Their rewritten save metadata matches that actual request. The historical `default-load` helper compares the effective system text against a literal without the LF and emits a combined model/prompt error.

One causal CLI contrast removed only that final LF from the owned primary fixture. The exact system predicate changed from false to true; current model, vendor, tools, restored user history, and metadata agreement remained correct. This diagnoses the representation-sensitive predicate exercised here. It is not a full grade replay, an acceptance decision, or proof of all Chapter 10 behavior.

**Classification:** a fixture/oracle expectation mismatch for this current-configuration assertion; the evidence does not support treating this predicate failure as historical-config restoration by the student. There is a root-owned compatibility/policy question about the expected primary-body representation. The public text does not explicitly require trimming a terminal LF, and this review does not invent that requirement—or a universal prohibition on startup normalization.

The sealed candidate remains blocked: the required Chapter 10 result is still **80/100, `passed: false`, exit 1**. **Chapter 1 remains NOT RUN.** No student source, grader, reference, public guide, historical score, or previous seal was changed. Root alone decides any correction, feedback, authorized subsequent gate, or acceptance.

## Authority, provenance, and prerequisites

Task: `/tmp/cadence-ch10-default-load-diagnosis-task.txt`, 6917 bytes, SHA-256 `363b8f481121333a3402eddf4dc603c84f5f782c6b40964d89ff38d04cb77ce6`. Governing input: `docs/edition-3-notes/ch10-default-load-diagnosis-input.json`.

Before the valid probes, the reviewer verified 24 pinned inputs, 33 protected tracked files, all 137 frozen candidate files, both frozen binaries, both locations of the candidate archive, the blocked evidence seal, and the four earlier review/entry inventories. Unchanged prior complete governing/craft/historical reads were carried by exact hash comparison; changed execution/student guides, current Chapter 10, the new student reports/source manifest/seal/parent verification, and the relevant exact implementation and grader code were read. The source manifest equals the diagnosis input's candidate object. Full before/after pin receipts are retained.

The frozen executables were used directly, not rebuilt:

- CLI: SHA-256 `a18f824da32a225bd6960ebf5d6128420da2a065d640ea38269a31fa2a13cc33`.
- GUI: SHA-256 `fb506081f98c09d2e6da628edbb6cb922f067f9db24983ae856d2d1222dcbf53`.
- Required checker identity is separately reverified in the final receipt; the checker was **not executed** for this diagnosis.

Actual full commit/tree/parent identities, archive hashes, source hashes, and executable paths are in `verification-before.json` and `corrected-v2/verification-before.json`. The current manifest/Git identities are authoritative; the earlier student report's mistranscribed predecessor hash is not adopted. The blocked seal is the actual `seal.json`, 35042 bytes, SHA-256 `37eecf6fa8424a471ba80c98af75478f6353d7b9b1b4fc3e1fa3d52af19b95b7`; its 146 members verify. Its pre-existing TAR remains untouched and local-only. None of those packaging/report qualifications is treated as a product defect.

## What the historical assertion actually expects

`internal/gradeedition2/ch10.go`, `checkPersistenceCurrentConfig` (lines 473 onward), writes a current primary skill declaring `read_file` and ending with:

```text
CURRENT_RUNTIME_PROMPT\n
```

Here `\n` denotes one actual LF byte in the fixture, not two literal characters. The helper selects OpenAI, model `gpt-5-course`, and user input `CURRENT_WORD`. It checks that prior user words occur once, then tests:

```go
current["model"] != "gpt-5-course" ||
    persistenceSystem(current) != "CURRENT_RUNTIME_PROMPT"
```

The expected system literal has **22 UTF-8 bytes**, with no final LF. `persistenceSystem` extracts/joins the system-message text; it does not trim it. Failure of either side produces “restored historical model/vendor/system prompt instead of current runtime config.” That message does not independently establish a vendor mismatch, a saved-config overwrite, or which half of the conjunction failed.

The helper subsequently requires exactly `load_skill`, `read_file`, and `unload_skill`, and verifies the rewritten save against the request. Those requirements remain intact in this diagnosis; no assertion or score was edited.

## Relevant contract and source trace

### Public obligations

- `book/edition-3/chapter-10.md:23–42`: `config` records request settings; `system_prompt` is the effective renderer instruction; vendor uses its existing disk spelling; tool declarations are data, not handlers.
- Chapter 10, hard rule 7 (`87–89`): **“Config is a record, not a restore.”** The running Agent's own model, vendor, prompt, and tools win.
- Chapter 10 `177–195`: configured versus effective instructions must not make metadata lie; saving and renderers must share the effective decision; intentionally empty primary instructions must remain empty; a new Agent supplies current primary/tools rather than restoring capability from declarations.
- `book/edition-3/chapter-09.md:77–81`: the primary body and initial dependencies become the system prompt rendered at creation; the two skill lifecycle tools are implicit.
- Chapter 9 also requires the initial primary instructions to remain fixed across later skill loads. That is not cited as a separate, explicit rule about trimming terminal LF during initial parsing.

Neither cited chapter specifies the helper's no-terminal-LF literal as the mandatory normalization policy. Preserving the selected primary's body is consistent with the stated effective-instruction path. The causal probe demonstrates actual behavior; it does not establish that every implementation must preserve or trim every kind of whitespace.

### Frozen implementation

- `internal/skills/parse.go`, `parse`: normalizes CRLF to LF, separates frontmatter/body at the delimiter, and assigns `d.body = body` without trimming that body. Header scalar/list parsing is separate from body handling.
- `internal/skills/registry.go:119–132`: rendering joins bodies and performs one-pass variable substitution; it does not trim the result. This fixture has one body and no substitution.
- `ensemble.go:151–157`: `LoadInitial` supplies the current system text, `SetInitialSystem` installs it, and only then `LoadStartup` runs.
- `internal/agent/save.go:34–56`: startup reads the selected archive into history and restores conversation state/display high-water handling. It does not assign archive config into the active Agent config.
- `internal/agent/save.go:75–81`: saving records the active model/vendor/system and current declarations.
- `internal/engine`'s OpenAI rendering uses the active configuration for the system message and model. The request captures below confirm that path at the actual executable boundary.
- `cmd/ch01/main.go` and `gui/main.go` construct the current config/primary from launch inputs, use default `save.json`, and share the public Agent startup path. GUI's independent preferences path is not merged into the conversation save.

This is a narrow trace for the reported mismatch, not a new whole-code architecture or lifecycle review.

## Fixture provenance and finite controls

### The save is deliberately minimal and synthetic

The owned 1284-byte seed is a handcrafted, valid public four-field snapshot-only archive, **not the original grader's captured save** and not a replay of its multi-start seed campaign. Its SHA-256 is `6e8a04338d56109dcc947f46cf356ee6cafe0c39763ae0764446b8cb8e75601a`.

It has `as_of: 6`, six neutral conversation turns, `log: []`, and deliberately opposing recorded config: vendor `anthropic`, model `fake-model`, prompt `RECORDED_PROMPT_MUST_NOT_RESTORE`, and a recorded-only `write_file` declaration. The prior user texts are `FIRST_WORD`, `SECOND_WORD`, and `THIRD_WORD`. This makes accidental metadata restoration observable without performing extra historical sessions. It proves the tested snapshot-default-load path, not every path in the original grade.

The corrected probes match the exact relevant historical inputs: primary body `CURRENT_RUNTIME_PROMPT` plus LF, explicit `tools: read_file`, current model `gpt-5-course`, and user `CURRENT_WORD`. Non-semantic header description/order and the synthetic prior dialogue/reply differ from the historical seed. Those differences are not concealed as a byte-for-byte replay.

### Launch and transport discipline

The only real-executable launches were one invalid pre-request attempt and the three root-authorized corrected launches below. No build, full grade, Chapter 1, package/race/lifetime test, legacy reference run, browser, WebSocket campaign, or paid request ran. The GUI used `-port 0`; its observed listener was `127.0.0.1:59603`, not 8088. One owned ephemeral loopback HTTP server served the three valid requests using deterministic dummy replies. It checked a fixed dummy Authorization value without exporting headers.

Each launch had an explicit environment, owned HOME/TMPDIR/save directory, a 60-second deadline, separate stdout/stderr, and a recorded PID/process group. `LLM_BASE_URL` and `OPENAI_BASE_URL` intentionally selected that exact owned loopback server. Generic provider overrides were not inherited. Exact argv, cwd, allowlists, stdin, durations, exits, and output hashes are in the per-launch receipts.

**Environment deviation:** the runner set `GOMAXPROCS=2`, `GOFLAGS='-mod=readonly -p=1'`, and the required Go/Delve/system PATH, but omitted `GOTOOLCHAIN=local`, `GOPROXY=off`, and `GOSUMDB=off` from its client allowlist. No Go invocation/build or tool call occurred; the omission is retained as a procedural qualification, not repaired by cosmetic reruns or hidden in reconstructed environment receipts.

### Valid observations

| Control | Exit / duration | Actual model | Actual system | Historical model predicate | Historical exact-system predicate |
|---|---|---|---|---|---|
| Corrected CLI default-load | 0 / 0.032673 s | JSON string `gpt-5-course` | JSON string `CURRENT_RUNTIME_PROMPT\n`, 23 UTF-8 bytes | true | **false** |
| Corrected GUI default-load | 0 / 0.027124 s | JSON string `gpt-5-course` | JSON string `CURRENT_RUNTIME_PROMPT\n`, 23 UTF-8 bytes | true | **false** |
| One CLI contrast, terminal LF removed | 0 / 0.022912 s | JSON string `gpt-5-course` | JSON string `CURRENT_RUNTIME_PROMPT`, 22 UTF-8 bytes | true | **true** |

The 23-byte system hex is:

```text
43555252454e545f52554e54494d455f50524f4d50540a
```

The expected and contrasting 22-byte hex is the same sequence without final `0a`.

For every valid request:

- Transport was `POST /v1/chat/completions`, with OpenAI message/tool schema and `stream: false`. The request has no asserted `vendor` JSON field; provider identification is the actual endpoint/schema plus current configuration.
- Tool names were exactly `load_skill`, `read_file`, `unload_skill`; recorded `write_file` was not installed.
- The restored user texts `FIRST_WORD`, `SECOND_WORD`, and `THIRD_WORD`, plus new `CURRENT_WORD`, each occurred exactly once as user text.
- Saved `config.model` was the current model string; saved `config.vendor` was the JSON string `openai`, not a numeric enum or the saved seed's `anthropic`; saved `config.system_prompt` exactly matched the wire string including the LF when present.
- Saved declaration names matched the current wire tools. All saves retained the four expected top-level fields, reached `as_of: 10`, contained four new log entries and eight context turns.

No observed active field reverted to the seed's opposing config. Thus the broad failure message is not an accurate causal description of these default-load observations.

### Causal contrast and its limits

The corrected CLI and contrast use the same frozen binary, cwd, environment, user input, seed-save bytes, server endpoint, and deterministic reply. The sole changed input is removal of the last byte `0a` from the owned primary skill file. The file shrinks from 113 to 112 bytes:

- With LF: `f5fffa60b1ecd62fef9fe4919cbd52fad6fe073563fe6187bcaaac57ce5a4c1c`.
- Without LF: `044baf2c02ac8c1eec49e29d2e42d126332fbf0242d441c858897edb81227851`.

The same input save is restored before the contrast; normal prior-run output files are not input to the configuration decision. Only the system string and corresponding serialized metadata representation change. The body is restored to its original LF form after the contrast. There is no production mutation or test relaxation.

This demonstrates that terminal-LF representation is sufficient to flip the exact historical predicate on this frozen executable. It does not rerun or change the original grade, establish all history equivalence/ownership/identity guarantees, test intentionally empty primary behavior, or prove the original grader's full captured request byte-for-byte.

## Reviewer failures, authorization, and unrun work

All failures remain in this new package; none is overwritten by the corrected subdirectory.

1. Read-only source lookups initially guessed incorrect filenames, and combined outputs sometimes exceeded output caps. Actual paths and bounded relevant sections were read afterward. No runtime evidence is inferred from those failed reads.
2. The first preservation verifier used obsolete `ch10-entry-audit-followup-evidence`; the actual package is `ch10-entry-followup-evidence`. It exited before workspace/server/client creation. The original driver and path-failure receipt are retained.
3. The first owned primary used comma-separated tool names. The frozen CLI correctly rejected `read_file,`, exit 1, with **zero provider requests**. Its GUI and contrast were NOT RUN. This is reviewer fixture failure, not the product default-load failure.
4. That invalid draft also used `gpt-5-course-fast`/`SWITCH_WORD`. Root messages `main-29` and `main-30` explicitly authorized a versioned recovery and alignment to the exact historical tools/model/user inputs. The original fixture, labels, stdout/stderr, receipt, failure, and cleanup remain unchanged. The full corrected driver diff and hashes are recorded before execution.
5. The corrected version ran exactly the remaining CLI, GUI, and one terminal-LF contrast. Total across both versions: **four client launches, three valid provider requests**, within root's amended cap. There was no further failure or retry.
6. The three missing Go-environment variables are disclosed above; no control was rerun to make receipts cosmetically conform.

No private fixture or answer code was sent to the student, who remains idle/frozen. No implementation correction was attempted. No full-grade result is inferred from these micro-probes.

## Evidence, preservation, and cleanup

Portable evidence is the exact closed set in `ch10-default-load-diagnosis-evidence/inventory.json`, including this external report. Highlights:

- Root-level files preserve the original invalid attempt and pre-control verifier failure.
- `recovery-authorization-and-diff-pins.json` and `recovery-driver.diff.txt` record the authorized changes.
- `corrected-v2/*.receipt.json`, stdout/stderr, `observations.json`, and per-control sanitized observations record actual executions and field types/bytes.
- `snapshot-fixture.txt` and the primary fixture texts are demonstrably synthetic. Raw full HTTP bodies, API/debug logs, full rewritten saves, binaries, and archives remain local/hash-only. No raw authorization header or file named `settings.json` is exported.
- Each version has a local-only artifact identity index limited to its owned scratch; no broad old-worktree or user-memory sweep was performed.
- Final receipts reverify the candidate, frozen checker, archives, all pinned/protected files, blocked seal, and earlier inventories; source Git identities/diffs/status remain unchanged and the student stays clean.

All client jobs have exited and been waited. Both owned HTTP servers were shut down and joined; no owned process remains. Both runtime TMPDIRs are empty. The corrected skill fixture is restored. The 1 GiB floor held at measured boundaries (about 3.1 GB free at corrected completion); no shared cache cleanup occurred. No archive or executable is included in this new portable package.

**Root handoff:** the evidence isolates an exact-string representation mismatch and rules out historical-config restoration in the tested paths. Decide the appropriate contract/fixture compatibility treatment before changing any student instruction or grader. Preserve the original blocked grade until a separately authorized gate supplies a different result.
