# Edition 3 Chapter 4 CLI adapter — independent frozen-candidate review

## Scope and decision

**Decision: one narrow revision requested (R1: help exit status).** All measured grading routes pass; no other adapter finding remains. This reviews only the parent’s Chapter 4 command adapter and its tests, before student launch. It is not Chapter 4 student acceptance, a new audit of historical grading coverage, or a manuscript text-review pass.

Frozen inputs, verified before execution:

- `cmd/grade-edition3/main.go`: SHA-256 `d08616eee707b4b91363f528abafa85f9b2268af964a0c6d903efbae95c291a4`.
- `cmd/grade-edition3/main_test.go`: SHA-256 `f505aebe9a8ef9bea5495c6ee723eea0880b1ae3735e4b36adadff2cceb71a5f`.
- Parent manifest: `docs/edition-3-notes/ch04-grader-evidence/candidate.json`.
- Base main HEAD: `bac07e644facb25ba3e5bffb2961674763fed4c4`.

The independent scratch tree is a fresh Git export of that HEAD with **only the two candidate files overlaid**. Every tracked file under `internal/grade`, `internal/gradeedition2`, `internal/gradeedition3`, `internal/fakevendor`, and `internal/fakeanthropic` was checked byte-for-byte against the export. No historical grading package changed. Complete governing-document reading from the entry turn remains applicable: all eleven required documents were hash-verified unchanged. The live candidate, application, student, manuscripts and original reference were not edited.

Evidence is in `ch04-grader-review-evidence/`: `input-identities.json`, exact commands/working directories/status/duration/output hashes in `runs.json`, complete command outputs, bounded runners, reference-only correction patch, private adapter mutation patch, executable identities, and completion records. Parent evidence was inspected but is not relabeled as independent execution.

## Finding R1 — restore successful help status

**Low-severity CLI compatibility regression; restore before final adapter acceptance.** At `cmd/grade-edition3/main.go:28–30`, every `flags.Parse` error becomes status 2. The local `ContinueOnError` parser returns `flag.ErrHelp` for `-h`, whereas the old global parser handled help by exiting successfully. Independent execution of the preserved, hash-verified old binary returns **0** for `-h`; the independently compiled frozen candidate prints valid help but returns **2**. Neither invocation launches a submission. See `old-help.log`, `candidate-help.log`, and their exact argv/status/binary hashes in `runs.json`.

Smallest remedy: return 0 when `err == flag.ErrHelp`, keep return 2 for actual parsing errors, and add a help/no-dispatch test. No helper abstraction or flag framework change is needed. This is not a Chapter 4 grading failure, but it is an unintended change to existing command behavior. The frozen candidate was not edited here. Re-review the new input hashes and run focused help, race/vet and invalid-input controls after repair; unchanged expensive behavior controls need not be repeated merely for ceremony when the diff is confined to that branch/test.

## Source review: KISS and contract

Aside from R1, the adapter is suitably small and correct for its stated contract:

- One local `flag.FlagSet` (`ContinueOnError`) supports repeated negative invocations without modifying global flag state. There is no injectable grader framework, callback abstraction, or production test hook.
- Only Chapters 1 and 4 are dispatched; exactly one nonempty target is required, with no historical target fallback.
- Chapter 4 requires explicit `-core` and `-consumer` directories with regular `go.mod` files. It does not infer source identity from a binary path or impose a new application package layout. These are grader inputs, not application CLI/config obligations.
- `flags.Visit` detects even explicitly empty Chapter 4-only flags and rejects them for Chapter 1. This preserves the intended version boundary.
- Module-directory validation and existing `ReadConsumerCase` parsing precede build/launch. Optional `-consumer-case` uses the existing shout default when omitted; supplied cases retain existing validation semantics.
- `grade.Build` is reused unchanged. Its deferred cleanup runs before `main` calls `os.Exit`, including error and failed-report paths.
- Chapter 4 calls exactly `gradeedition2.Chapter4(bin, core, consumer, example)`; Chapter 1 still calls `gradeedition3.Chapter1(bin)`. The command does not alter points, checks, or score aggregation. Existing Chapter 4 output retains its historical Edition 2 title; the command version selects that unchanged contract rather than pretending it is a newly rewritten grader.
- Setup/parse/build/grader errors return 2; failed reports return 1; passed reports return 0. The decision uses `report.Passed`, not a full numeric score. JSON output failures return 2 as before.
- Frozen tests have 18 invalid cases and an observable launch sentinel. They establish validation/no-dispatch boundaries, not full positive routing. This review supplies actual compiled full-route controls separately.

## Supported invocation

Run flags before the final target, following Go flag parsing:

```sh
go build -o /absolute/path/grade-edition3 ./cmd/grade-edition3
# In the candidate core module with the accepted root-library layout:
go build -o /absolute/path/submission ./cmd/ch01

/absolute/path/grade-edition3 -ch 4 -json \
  -core /absolute/path/core-module \
  -consumer /absolute/path/consumer-module \
  /absolute/path/submission

# Optional fixture replaces the default shout case:
/absolute/path/grade-edition3 -ch 4 -json \
  -core /absolute/path/core-module \
  -consumer /absolute/path/consumer-module \
  -consumer-case /absolute/path/case.json \
  /absolute/path/submission

# Separate inherited lifetime gate on exactly the same executable:
/absolute/path/grade-edition3 -ch 1 -json /absolute/path/submission
```

A directory target is also supported where the unchanged builder actually finds a `main` package (module root, or `cmd/` itself). That path is independently measured below. It is **not** a claim that the accepted root-library plus `cmd/ch01` module root is directly executable. The earlier archive/launch failure remains valid; the explicit executable route is the supported answer for that layout. Source-directory identity and executable bytes remain separately recorded provenance; this adapter does not cryptographically prove that an arbitrary supplied executable came from its supplied source.

## Independent controls

All following new-adapter controls use the **newly compiled frozen candidate**, not the old CLI, API driver, or parent binary. Independent Go toolchain: `go1.25.4 darwin/arm64`. Provider secrets/settings were filtered from control environments; dependency fetching was disabled (`GOPROXY=off`, `GOTOOLCHAIN=local`).

| Control | Actual result |
| --- | --- |
| `go test -race -count=1 -v ./cmd/grade-edition3` | exit 0, all 18 frozen negative subcases pass (4.015 s) |
| `go vet ./cmd/grade-edition3` | exit 0 (1.176 s) |
| Build frozen adapter and timeout-only-corrected scratch reference | both exit 0 |
| New E3 `-ch 4`, explicit executable/core/consumer, default case | exit 0, **160/160, passed=true** (8.763 s) |
| New E3 `-ch 4`, supported module-root main directory, explicit core/consumer and non-default case | exit 0, **160/160, passed=true** (8.631 s) |
| Same new CLI/binary/consumer, wrong core directory with the same module name | exit 1, **135/160, passed=false**, only `star-topology`: consumer resolves the other directory, not the supplied core (8.391 s) |
| Compiled CLI with malformed, missing, or incomplete consumer case | each exits 2 in ~4 ms; sentinel never launched |
| Compiled Chapter 1 with each of `-core=`, `-consumer=`, `-consumer-case=` | each exits 2 in ~4 ms; sentinel never launched |
| Private adapter mutant forwards consumer directory as core | builds; otherwise-positive full route exits 1, **135/160**, only `star-topology` (8.220 s); sensitivity established |
| New E3 `-ch 1` on **the same immutable reference executable** used by the Chapter 4 positive | exit 0, **100/100, passed=true** (124.046 s); `generation_headers` and `generation_body` each explicitly report successful **1m2s** delay |
| Supplied valid consumer case with deliberately wrong expected result | exit 1, **135/160**, only `external-consumer` (8.336 s); caller fixture is actually forwarded rather than silently replaced by default |
| Preserved old command `-h` / new frozen command `-h` | exit **0 / 2**, respectively; R1 |

The non-default positive case sends `{"text":"independent route"}` and expects `INDEPENDENT ROUTE`. Because generic shout accepts both this and the default input, that success alone would not prove forwarding; the valid-but-wrong expected-result negative supplies the missing discrimination. Likewise, the private misrouting mutation is a detecting control for actual core-argument forwarding, not a candidate failure. Neither mutation touches frozen files.

All 19 command receipts are finite; there were no review-runner build errors, launch errors, timeouts, or assertion failures. Expected negative controls are recorded as such. The 210-second outer lifetime bound allowed both unchanged real 62-second phases to finish; neither phase nor the underlying grading source was shortened. `completed.json`, `case-forwarding-completed.json`, and output-hash verification establish completion.

### Executable identities (SHA-256)

- New frozen-candidate adapter: `3f7af0cd6df371069ceae55240ca6a670d1a102cb8cf1e434fd0e185bb228a92`.
- Same reference executable for paired Chapter 4/Chapter 1 positives: `c7a958a2df3e208f32644ca72f48be08fb79880f44eaa603f4a1894cf78c4b21`.
- Private core-forwarding mutant: `49f1f281a6510a4937bd02e8388e7fc5ed99a99b4fcd220238bef529bd0bd2f7`.
- Preserved old Chapter 1 grader used only for help comparison: `e7057c765859dd4719d66c9ff0d5329c436418933b451f9e20813c7e29697cab`.

The original reference’s negative lifetime and accepted predecessor’s archive/launch failures remain preserved with verified output hashes in `preserved-entry-evidence.json`. The old reference was not silently “fixed” in place or relabeled a pass; its scratch-only correction has an exact two-hunk patch and before/after source hashes. No student was launched by this review.

## Limits and recommendation

The adapter fixes invocation/identity separation and adds supported versioned dispatch; it does not close the coverage gaps from `ch04-entry-audit.md`. Chapter 4 retains 160 points and inherited Edition 2 Chapter 3/2/legacy Chapter 1 composition. Edition 3 lifetime remains a **separate required Chapter 1 invocation on the same binary**, with both zero-point lifetime checks and `passed` inspected. Selected-model policy, pending output-cap coverage, job cursor/concurrency/process-group behavior, public callback supervision, and ownership/source review retain their prior acceptance obligations.

The old reference’s negative lifetime result was preserved and its output hash reverified; it was not rerun for ritual. The new positive reference is a scratch-only copy with exactly the known whole-request timeout removed (and the unused `time` import removed). This is a controlled positive for routing and lifetime, not a new approved reference or a student patch. The ordinary historical jobs/external-consumer check logic and the real 62-second lifetime phases are unchanged. No shortened delays, fake report injection, live provider calls, paid calls, credentials, install, push, or restart were used.

**Recommendation:** repair only R1, freeze the resulting two file hashes, and obtain focused confirmation before marking this adapter accepted. The invocation separation, versioned dispatch, existing 160-point contract, case forwarding, and separate same-binary lifetime route have otherwise been independently established. This does not remove the entry audit’s requirements before accepting a future Chapter 4 student.

Final state: main HEAD remains `bac07e644facb25ba3e5bffb2961674763fed4c4`; both frozen input hashes still match. Accepted student HEAD and peeled tag remain `28be8bffdbc1356eaaa50a97369e59d0e0556dc3`, working tree clean. Concurrent/untracked work was preserved. See `final-state.json`.
