# Edition 3 Chapter 2 — entry audit

## Disposition and boundary

**Entry alignment needed; no Chapter 2 implementation accepted or launched.** Audited book/grader HEAD `5351583208988f8370937cd0494727645efa8b46` against the unchanged predecessor `/Users/bill/projects/ensemble-edition-3-student`, HEAD/tag `1aad3975842d502b9144b7b9bddfd30ca7309db7` / `edition-3-ch01-r1`. This is independent, private author/reviewer evidence, not student guidance. Do not expose reference sources or mutation details to a future student.

Completed the required full policy/skill/workflow/carryover/voice/procedure/current-chapter reads, recovering omitted output. Current Edition 3 procedure supersedes historical phase/path instructions. Only this audit note was written in the repository; tests, binaries, relocated reference copies and mutations stayed in `/tmp/e3-ch02-entry-audit.f52GpR`. Tracked chapter, grader, reference and student files were not changed. No credentials, live-provider traffic, GUI/live-session inspection, install, restart, push or publication build.

## 1. Exact supported targets and inherited gate

`internal/grade/harness.go:Build` accepts an executable, or builds a directory. It selects `./cmd` only when **that directory itself** has a top-level `package main`; otherwise it builds `.`. It does not discover `cmd/ch01` or `cmd/ch02`. A library-root build can produce a non-executable archive without a useful build failure. The accepted predecessor already documented this in `ch01-code-review-r1.md` (R2).

The public-library + nested-CLI architecture is valid and must not be flattened to accommodate this helper. From the grader repository root, the directory-target form is:

```sh
go run ./cmd/grade-edition2 -ch 2 -json /absolute/student/actual-command-directory
```

The existing predecessor's command directory is `/Users/bill/projects/ensemble-edition-3-student/cmd/ch01`. `cmd/ch02` is also supported if that is the eventual actual command package; the grader does not require either name. Passing the student's module root is not the supported recipe for this layout.

Prefer one explicitly built artifact for both gates. This is a **future gate recipe**, not a claim that the unlaunched predecessor implements Chapter 2:

```sh
student=/Users/bill/projects/ensemble-edition-3-student
cli_pkg=./cmd/ch01  # current actual package; change only if the implementation relocates it
out=$(mktemp -d)
(cd "$student" && go build -o "$out/ch02" "$cli_pkg")
cd /Users/bill/projects/ensemble
go run ./cmd/grade-edition2 -ch 2 -json "$out/ch02"
go run ./cmd/grade-edition3 -ch 1 -json "$out/ch02"
```

Both commands must exit zero and report **`passed:true`**, not merely 100/100. `internal/grade/report.go:NewTitledReport` includes zero-point failures in `Passed`. `cmd/grade-edition3/main.go` implements **only Chapter 1**; `-ch 2` is usage error, not a Chapter 2 gate.

Why both: `internal/gradeedition2/ch02.go:Chapter2` replaces inherited parity with `grade.Run` + **Edition 2** `Chapter1`. It preserves conversation/usage/exact Anthropic-version checks, but never calls Edition 3 `Chapter1`. The latter adds 62-second delayed **non-streaming** headers/body checks (`internal/gradeedition3/ch01.go`). These are not SSE tests. A pre-Edition-3 reference with a 60-second client timeout still passes Chapter 2 today (evidence below). Any arbitrary longer generation/body deadline remains source-review territory even after the two probes pass.

**Align now:** Chapter 2's exercise still gives the ambiguous directory command and §2.9 says `go build -o ch02 .` (`book/edition-3/chapter-02.md:1257,1294`). Publish the explicit command-package/executable recipe and same-artifact Edition 3 Chapter 1 gate in the entry instructions; do not require a root `main`, a new module split, or a new Chapter 2 grader. Gate the Chapter 2 executable, not a stale separate Chapter 1 binary.

## 2. What the current Chapter 2 grader proves

Sources: `internal/grade/ch02_harness.go`, `ch02_checks.go`, `ch02_log.go`, `ch02_ref_checks.go`; `internal/gradeedition2/ch02.go`, `ch02_ref.go`.

The public CLI is `session`, `dump`, `render`; local fake servers exercise all three wire protocols. The overlay retains log round-trip/fresh-process rendering, deterministic replay, redaction, one-shot ephemera, normalized usage, request-shape and response-normalization checks. Edition 2 adds resolved-model provenance (including a response resolving a previously unresolved model), captured response/call model identities and Gemini call-bound signatures, Gemini `functionResponse.name`, Anthropic interleaved result-block ordering, and redaction pressure across all three Ref kinds. Opaque replay tests include exact Gemini call-bound signature preservation for the matching fixture identity; they do not exhaust mismatched provenance/refusal cases. These are useful targeted checks, not exhaustive state-machine or API acceptance.

Scope already correctly narrowed: deferred URI-media rendering is removed, with its four points reassigned to observable preservation work. Three Ref kinds must remain representable; this chapter need not fetch media, implement uploads, replay URI media, run tools, stream responses, add retries, build later GUI behavior, or add production model-resolution machinery. Do not reintroduce the removed `ref-render` requirement.

`book/edition-3/chapter-02.md:1245–1251` correctly limits the Ref CLI claim: dump exposes retained original Ref/text/redaction facts; black-box refusal/redaction checks do not expose exact reduced Ref carryforward. `ch2RefRoundTripCheck` checks three ordered MIME/kind/locator tuples, not the reducer projection. `refRefusal` requires empty stdout and a stderr diagnostic but accepts exit zero when stderr is present, despite its diagnostic describing a nonzero exit. Do not describe this as exhaustive strict-loader or exit-code coverage.

## 3. Source-test/review duties, not a grader rebuild

The reviewer must inspect the actual predecessor-to-Chapter-2 delta and implementation, not substitute the reference's tests for student evidence:

- Preserve a public `Ensemble` construction/import boundary, concrete-parent reachability, library-owned engine/history, per-instance history and configuration, tiny CLI, and separate-module GUI only if present. Keep provider/history/renderer logic out of CLI stdio handling; no application-owned global history filename, manufactured instance identifiers, mock interfaces or speculative framework. A root-module `go test ./...` does not include a separate GUI module.
- Inspect all production clients/transports and cancellation paths for inherited deadline limits; retain useful nonzero error paths, scanner/EOF behavior, root-path OpenAI compatibility and truthful usage. Two latency probes cannot establish absence of every deadline or error-path regression. Prior accepted evidence is not a new Chapter 2 streaming implementation mandate.
- Verify append-only event history versus reduced projection; timestamps/ordering/state transitions; append/load/dump/replay equivalence; custom system configuration; pending ephemera consumption; in-memory session isolation; and fail-closed unknown interpreted version/event/part/actor/vendor/location data. Do not impose historical sequence-number ideas already rejected by carryover.
- Verify **exact Ref kind and locator** through load, reducer, both redaction levels, re-reduction and dump without resolving/fetching it; stub survival; all content-part redaction; and removal of superseded content from the reduced/rendered view without losing the deliberately retained Ref. No need to expose a new debugging API merely for the grader.
- Verify raw opaque capture and exact vendor/model/surface matching, no silent discard or textification, multiple-call/result pairing and ordering, unmatched-result and required-signature refusals. The graded examples are not exhaustive combinations.
- Verify absent-ID synthesis is deterministic, nonempty and distinct for multiple calls, and preserves issued IDs (public §2.5 requires issued IDs preserved and deterministic synthesis from `Seq`; distinct missing IDs are exercised by the existing source regression, not a new ID-format rule). **New control: a missing-ID defect survives the public grader but fails an existing source test.** This duty must be explicit in the review/test checklist; a small missing-ID black-box fixture could be added later by deliberate scope choice, but is not grounds to replace all graders or invent a new interface.

The accepted predecessor review and correction (`docs/edition-3-notes/ch01-code-review-r1.md`, `ch01-code-review-correction-1.md`, `ch01-completion.md`) remain inherited evidence, not proof that a future Chapter 2 delta preserves them. The two-command gate plus explicit source-review duties is a bounded route; grader-only acceptance is not.

## 4. Independently run evidence and coverage limits

All runs used an allowlisted environment (`PATH`, `HOME`, `TMPDIR`, `GOPROXY=off`, `GOSUMDB=off`), not inherited provider credentials. Explicit target paths were supplied. Outputs and nonzero controls remain in the scratch directory above.

| Run/control | Observed result |
| --- | --- |
| Root `go test -race -count=1 -timeout=180s -run '^TestChapter2AcceptedSolution$' -coverpkg=./internal/grade,./internal/gradeedition2,./internal/fakevendor,./internal/fakeanthropic -coverprofile=… ./internal/gradeedition2` | Pass. Current test explicitly targets `../../solutions/ch02`, checks 100/100 and `Passed`; no skipped export. |
| Root `go test -count=1 -timeout=60s -run '^TestBuild$' ./internal/grade` | Pass; its fixtures execute root main, root main with a named-command container, preferred direct `cmd/main.go`, root fallback when `cmd` is a library, and a broken preferred command. It does **not** test library-root archive rejection or directly target a named nested command. |
| Root `go test -race -count=1 -timeout=90s -run '^(TestGenerationLifetimeControls|TestHistoricalTimeoutIsExceeded)$' ./internal/gradeedition3` | Pass: short no-deadline controls pass at headers/body; short deadline controls fail; real probe is asserted longer than historical 60 seconds. These are short mechanism tests, not the full two-minute gate on a new Chapter 2 candidate. |
| Built root `cmd/grade-edition2`, `-ch 1 -json /Users/bill/projects/ensemble-edition-3-student/cmd/ch01` | Exit 0, 100/100, `passed:true`. Confirms explicit nested target on unchanged actual predecessor; does not establish Chapter 2 support. |
| Scratch reference copy, only existing `main.go` relocated to `cmd/ch02`; explicit `-ch 2 -json "$scratch/nested/cmd/ch02"` | Exit 0, 100/100, `passed:true`. Confirms named nested command-directory support, not a forced root layout. |
| Missing-ID single-defect scratch copy; explicit `-ch 2 -json "$scratch/missing-id"` | **Surviving negative control:** exit 0, 100/100, `passed:true`. |
| Same mutant, `go test -count=1 -timeout=60s -run '^TestRefCarryforwardAndMissingCallIDs$' ./ensemble` | **Expected failure retained**, exit 1: `missing call IDs were not distinct`. Baseline source test resides at `solutions/ch02/ensemble/replay_test.go:218–257`. |
| Unmutated relocated scratch reference, `go test -race -count=1 -timeout=90s -coverprofile=… ./ensemble` | Pass, 88.2% **public-package-only** statement coverage; not whole engine/history/CLI coverage. |

Coverage interpretation: the broad grader-process profile reports **5.9%**, diluted by unrelated grader chapters. Its relevant files report `internal/gradeedition2/ch02.go` **79.1% (117/148 statements)** and `ch02_ref.go` **75.4% (43/57)**; `internal/grade/ch02_harness.go` **83.1% (172/207)**; `ch02_checks.go` **68.2% (242/355)**; `ch02_ref_checks.go` **45.9% (62/135)**. These measurements exercise graders/fakes, **not the child executable** built by `grade.Build`. Source coverage of a future student implementation needs that module's tests and a named production-package denominator; CLI subprocess coverage requires deliberate instrumentation, not attribution of the parent's coverage. No student Chapter 2 coverage exists yet, and neither 88.2% facade coverage nor 100/100 grader score establishes the inherited acceptance threshold for the whole implementation.

Concrete inherited false green: the passing Chapter 2 accepted-reference test targets `solutions/ch02`; `solutions/ch02/internal/engine/engine.go:36` still sets a 60-second client timeout. This is source-confirmed absence of the Edition 3 gate, not a fresh 62-second failure measurement. Historical Chapter 2 targeted sensitivity results (nine mutations, including opaque/ordering/Ref/capture cases) are documented separately in `docs/edition-2-notes/enhancement-02-code.md` and `enhancement-02-review.md`; they were **not all rerun** here. The new missing-ID survivor supplements rather than erases those results.

Not run: a full Edition 3 62-second-per-phase gate on a future Chapter 2 implementation (none exists), all-repository/all-chapter tests, a new student public-surface smoke/coverage suite, GUI tests, live-provider tests, all historical mutations, arbitrary longer-deadline search by runtime waits. `rg` was unavailable; standard reads/grep were used. An audit-created `git diff` pager was stopped and diff checks repeated with `--no-pager`; final tracked diffs in grader and student repositories were empty.

## Bounded next step

Before student launch, the entry owner should align the public build/target recipe and explicitly require **Edition 2 Chapter 2 + Edition 3 Chapter 1 on the same Chapter 2 artifact**, with `Passed` checks and the source-review checklist above. Resolve the root-build wording without flattening the library. Then launch only the ordinary Chapter 2 exercise from the accepted predecessor when authorized. At completion, collect implementation-local tests/coverage, run those two gates, inspect the delta and public construction surface, and perform the specific Ref/missing-ID/replay/lifetime source checks. No code acceptance, broad grader rewrite, reference repair or unrequested feature is authorized by this audit.

## Follow-up receipt — required real-duration composition controls

Completed the previously omitted full-duration reference controls. This receipt supersedes the earlier **reference-control** limitation (“not a fresh 62-second failure measurement”); it does not turn reference evidence into student acceptance. Evidence: [`ch02-entry-evidence/grader-controls/`](ch02-entry-evidence/grader-controls/), including sanitized commands, stdout/stderr, parsed JSON reports, statuses/elapsed times, source/binary hashes and integrity checks.

| Explicit built target and grader | Actual elapsed | Exit | Score / Passed | Observed lifetime checks |
| --- | ---: | ---: | --- | --- |
| Unchanged historical `solutions/ch02` executable → Edition 3 Chapter 1 | 120.060 s | 1 | 100/100; **false** | `generation_headers` and `generation_body` both fail: no complete answer or matching usage survives the historical timeout. |
| Scratch Chapter 2 executable, only the 60-second blanket client timeout and unused `time` import removed → Edition 2 Chapter 2 | 0.473 s | 0 | 100/100; **true** | Full current Chapter 2 public checks pass. |
| **Same scratch executable**, no rebuild → Edition 3 Chapter 1 | 124.057 s | 0 | 100/100; **true** | Complete answer and matching usage survive both real 62-second header/body delays. |

All four explicit builds exited zero. No outer watchdog expired; the Edition 3 grader's actual `generationDelay = 62 * time.Second` was unchanged. This was sequential real execution, not the short lifetime tests or a source-only inference. Negative exit 1 and failed checks are preserved as expected sensitivity evidence, not repaired or hidden.

Historical executable SHA-256: `02cf107d3db198f65f8b218b1cbc131604b8710c7c521d0c1ad7af442bbedc4d`. Positive executable SHA-256: `f935d0b798db73e495c49443687e9d4b60f518cb98f62d000f9b60500662dc2a`, verified before/after **both** positive invocations. The source manifests prove the historical source was unchanged and the scratch source differs only at `internal/engine/engine.go` by those two permitted removals. No broader repair, other mutation, coverage expansion, protocol/routing change, tracked source edit or student edit was performed. OpenAI remains the Chapter 2 Chat Completions scope; this receipt introduces no Responses requirement.

**Conclusion:** the proposed same-artifact two-command gate composes successfully on this bounded control and detects the inherited timeout defect that Edition 2 Chapter 2 alone misses. Align the entry recipe/gate, retain the listed source-review duties, and require fresh student-specific evidence after authorized implementation. No Chapter 2 student code has been accepted or launched by this audit.
