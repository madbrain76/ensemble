# Chapter 8 entry reconciliation

## Corrected conclusion and scope

**A valid complete existing Chapter 8 API route is demonstrated: 100/100, `Passed=true`, exit 0.** The missing Edition 3 publication dispatch is a small adapter omission, not evidence that Chapter 8 needs a new grading framework or a prescribed student API. This is entry/adapter evidence, not student Chapter 8 acceptance; no student was launched or changed.

The original audit's scope/authorization claims were wrong. Its carryover appendix and extra named deliverables exceeded the original limited-output assignment. Root preserved the exact appended guide/suffix in `ch08-entry-scope-correction-evidence/` and restored only that owned appendix to the pinned guide bytes. This follow-up did **not** reappend or edit any guide or original audit file. The original inventory's carryover member is now a historical snapshot; all other inventoried members remain unchanged. This reconciliation supersedes the original fresh-Bill-approval and imposed-public-API-preparation recommendations. Bill already authorized autonomous continuation: no new routine permission gate is warranted. Next-turn budget behavior is the student's existing implementation obligation; architecture and API spelling remain the student's choice, not an entry prerequisite requiring a mutable public setter.

Only this report and `ch08-entry-reconciliation-evidence/` were added to the repository in this follow-up. All other work was owned local scratch. No repository source, grader, manuscript, guide, student, settings, or key edits; no children, browser/native GUI, paid service, live-8088 control, shared cleanup, or commit. The restored governing prerequisites were carried forward only after all 19 pinned input hashes matched.

## Complete historical route

Source: `internal/gradeedition2/ch08.go:18-30` exports

```go
Chapter8(bin, core, gui, guiModel, consumer, workflow, exercise, streamModel,
         example, media)
```

The scratch caller invokes that function unchanged. Its full mode is **not** a settings-only replacement. It executes the existing settings checks plus inherited Chapter 7 parity through the supplied identities. No standalone Chapter 1 lifetime rerun or entry-wide legacy matrix was launched.

Explicit identities (repository root `R=/Users/bill/projects/ensemble`; owned scratch `S=/var/folders/p_/11vk6y114xq7k9g5r_pc3yv80000gn/T/ch08-reconciliation.zwpmj1an`):

| Role | Actual identity/options |
|---|---|
| Executable | Built `go build -o "$S/bin/ch08-cli" .` with cwd `$R/solutions/ch08`; finite root executable, not GUI server |
| Core | `$R/solutions/ch08` |
| GUI | `$R/solutions/ch08/gui`; existing grader builds `./cmd/` |
| External consumer | `$R/solutions/ch08/ch05` |
| Workflow | `$R/solutions/ch08/ch06` |
| Stream exercise | `$R/solutions/ch08/ch07` |
| GUI and stream model | Explicit `claude-sonnet-4-6` for each |
| Consumer fixture | `shout`, arguments `{"text":"hello framework"}`, result `HELLO FRAMEWORK`, diagnostic `shout called` |
| Media fixtures | anthropic `claude-sonnet-4-6`; openai `gpt-5.4-mini`; gemini `gemini-3.8-flash` |

Caller build target: `go build -o "$S/bin/driver" .`, cwd `$S/driver` (owned module). Exact full invocation, cwd `$S/work`:

```sh
"$S/bin/driver" -mode full -bin "$S/bin/ch08-cli" \
  -core "$R/solutions/ch08" -gui "$R/solutions/ch08/gui" \
  -consumer "$R/solutions/ch08/ch05" -workflow "$R/solutions/ch08/ch06" \
  -stream "$R/solutions/ch08/ch07" \
  -stream-model claude-sonnet-4-6 -gui-model claude-sonnet-4-6 \
  -consumer-case "$S/consumer-case.json" -media-models "$S/media-models.json"
```

Exact absolute argv/cwd, exit, duration and process-group state: `historical-full-corrected.receipt.json`. Report: `historical-full-corrected.stdout`. **100/100; all five checks pass; exit 0; 15.63 seconds; no timeout.** Both explicit builds succeeded. CLI SHA-256: `1b9e8483141be6c1f26ad4cdfce72fb575c764811bd67cc737ab62756953189a`. Corrected caller SHA-256: `208f20dc53261f4c9f249c6177b5a6144169cf2f8b964338d72f03dbceda4e3c`. `builds.jsonl` records the actual build argv, cwd, compile exit, size and SHA-256 of internally built GUI/consumer/workflow/stream/helper artifacts as well; protected source identities are recorded in `source-baseline.json`.

**Invalid first attempt retained, not counted as a valid negative:** the audit caller initially put `guiModel` last instead of fourth and omitted explicit `GOPATH` from its scrubbed environment. That attempt returned 80/100, `Passed=false`, exit 1, with misrouted consumer/workflow/stream identities and an unavailable debugger presence check. `historical-full-first-invalid.*`, `main-first-invalid.go.txt`, and `clean-env-first-invalid.json` preserve it. Correcting only the owned caller argument order and exposing existing `/Users/bill/go` as `GOPATH` produced the complete positive above. No historical implementation/grader repair occurred, and no debugger was installed. The initial failure was invocation/environment error, not a demonstrated implementation contract mismatch.

Historical CLI topology is source-qualified rather than redundantly rerun: `cmd/grade-edition2/main.go` uses its historical path/build selection and forwards `path` as core to `Chapter8` at line 134. The explicit API caller separates finite executable identity from that source identity and optional GUI. Edition 3 already separates positional target from `-core`, but currently dispatches only through Chapter 7 (`cmd/grade-edition3/main.go:105-120`). Its `-ch 8` rejection/no-launch was measured in the original audit; no duplicate sentinel run was needed here.

## One real historical GUI persistence mutation

A bounded copy of `solutions/ch08/gui` was placed at `$S/real-gui`; only its scratch `go.mod` replacement was made absolute to the unchanged historical core. Old, negative and restored all used the existing `settingsChecks` via a scratch-only linkname bridge, each building the actual copied GUI `./cmd/`. This bridge is not a new grader or public API requirement. No wire-only server substitutes for the real GUI in this sequence.

Single mutation at `solutions/ch08/gui/settings.go:91` in the **scratch copy only**:

```diff
- err = os.Rename(path+".tmp", path)
+ err = os.Remove(path+".tmp") // OWNED AUDIT MUTANT: discard durable update
```

This discards the temporary settings file instead of publishing it, leaving the existing update/merge/broadcast path otherwise intact. One preparation anchor assertion failed before any source modification, compile or behavioral run; `mutation-preparation-failure.json` records that tooling error and the preserved first script. It is not a compile-negative or behavioral result.

| Real GUI stage | Compile | Settings score / Passed | Driver exit | Behavioral outcome |
|---|---|---|---|---|
| Old | 0 | 80/80 / true | 0 | All four settings checks pass |
| Persistence-disabled | 0 | 65/80 / false | 1 | Only `settings-persist` fails: restarted theme expected `system`, got `dark` |
| Byte-restored | 0 | 80/80 / true | 0 | All four settings checks pass again |

Before/restored source SHA-256: `6e2e90c4d61677a1dc7f667bb945335bd179e145043d2050eb663756d55364de`; negative: `f96dcb01c0bdbd3607b81998091d810cb78a75f38fc877bfa693522e21f069c6`. Old/restored compiled GUI SHA-256: `a36bcdbce768ffedbfe7f289f41260ab93264834fa8b2d04edc2d7d57e01d4f1`; negative: `aab4b21ad0acb6b25bed724fa00a24df75408f192454fc76ef1beb651a58c1ce`. Thus the negative is compiled behavioral sensitivity, not a build failure. Full source snapshots, compile receipts, actual reports and restoration hashes are retained in the evidence directory. The original two wire-only mutant receipts remain limited non-production proxies; this real historical GUI sequence supplies the persistence sensitivity they did not establish.

## Coverage and minimal next action

The existing complete grader covers settings roundtrip, subscription, broadcast, sparse retained values including explicit false/zero, actual process-restart persistence, and inherited Chapter 7 parity. The real mutation establishes sensitivity to loss of durable settings in production-path historical GUI code. A score alone does not establish the student's effective per-turn tool budget, next-turn semantics, ownership/implementation quality, or future browser control behavior. Those remain ordinary code/browser acceptance obligations of the existing chapter contract—not grounds for a new framework, a malformed-input/write-failure campaign, or an imposed preparation API.

**Recommendation: only the minimal thin Edition 3 Chapter 8 adapter**, preserving the existing explicit executable/core/GUI/consumer/workflow/stream/model/fixture identities and forwarding them in the verified order to the existing exported route. Extend the current publication chapter dispatch/validation/help consistently; do not replace inherited grading or rearchitect the student's public core. Existing autonomous authorization covers routine continuation; this follow-up itself makes no adapter edit and launches no student.

One narrow objectively measured completeness gap remains: the original wire-only mutant that omitted initial `max_tool_rounds` still passed because the initial snapshot check does not require that key. An existing-driver key-presence check is a proportionate possible correction; the result does not justify wholesale tests/frameworks or treating every unmeasured requirement as an entry blocker.

## Final preservation and inventory

`verification.json`: all 19 governing inputs, 120 student files, 25 protected tracked files and 234 follow-up source-baseline files unchanged; main/student commits pinned; student clean; no unexpected original-bundle differences. Only root's intentionally restored guide differs from the original sealed historical inventory. All eight subprocess receipts are terminal, no timeout, and no owned process group remains. Controls were serial, credential-scrubbed, local-fixture-only; module downloads disabled; Go builds single-worker. `GOPATH` exposes the existing debugger presence fixture, not credentials.

Final resource/process/status receipt and complete evidence hash inventory accompany this report. Owned local sources/binaries are deliberately retained for root verification, not deleted: `local-artifacts.json` inventories 81 files / 51,305,915 bytes under `$S` (including five local-only binaries). Internally temporary grader directories were removed by the existing grader itself. No shared or unrelated cleanup occurred. `inventory.json` hashes this report and every evidence file except itself; `tail.txt` is the final completion marker. All owned jobs have ended before submission.
