# Edition 3 Chapter 1 — initial independent code review

## Disposition

**Needs correction — one narrow test-fixture finding (R1).** No production-code defect was found in the reviewed Chapter 1 paths. This is the **initial** formal code review, not a correction/re-review round; zero of the three permitted correction rounds have been used. The earlier prerequisite-readiness turn was not a review. No author/text review or acceptance is implied.

The candidate passes the independently rerun vet/race suite, unchanged executable-target grader, and a separate reviewer-observed live chat. Those passes do not erase the faithful-boundary-fake requirement. Correct the successful-response fixture without weakening the chapter's assertions or adding later-chapter behavior.

## Exact inputs and provenance

- Student: **`d2baff64055ac33702a91c6d6ee5a4246e84c0c9`**, clean `/Users/bill/projects/ensemble-edition-3-student`.
- Empty predecessor: `69859da3328c3f710f8143203b99f934fed778a1`.
- Chapter/grader: **`7f3d069420112c2136bc95d1d551e12fad6964c1`**, `/Users/bill/projects/ensemble`.
- Chapter SHA256: `ab016baaff3f6a9d83ecffcbfc3772cd692c9d33c06d9fa1c52ad9a8708d25df`.
- Student review SHA256: `910cfc73d500da90b16c3049bc91c3cb9c486a166950eaf6cc74ef8908fd882e`.
- Read all required current prerequisites in the preparation turn; then read the complete candidate source/tests, README/module, student review/evidence, assistance record, Chapter 1 grader path and relevant historical implementations. Current chapter controls; no private-readiness coverage claim was substituted for source review.
- Built/tested a `git archive` of the exact candidate. After all experiments, every archived tracked file was compared byte-for-byte against that commit. Student HEAD and clean status were reverified. Full source/grader and binary hashes: [provenance.json](ch01-code-review-evidence/provenance.json).

## R1 — make the success fixture faithful to this chapter

**Required test correction; student `d2baff64055ac33702a91c6d6ee5a4246e84c0c9`, `ensemble_test.go:54–55`.** `TestConversation` returns a supposedly valid non-text response block:

```json
{"type":"thinking","thinking":"private","text":"not an answer"}
```

The adjacent comment calls this a “valid non-text block.” It is instead an invented successful Chapter 1 provider response: the chapter says replies are text here, with other content deferred to Chapter 3 (`book/edition-3/chapter-01.md:233–249` at the pinned book revision), and explicitly says to leave thinking out (`:590–600`). The fixture also puts an answer-looking `text` member on the alleged thinking block, rather than grounding that asserted success shape in the chapter or live evidence. The recorded live sessions do not establish this shape.

The useful HTTP-boundary fake should still exercise the real public API, multiple **text** blocks, exact history, measured usage and request headers. Make its normal successful response a faithful Chapter 1 response and correct the misleading comment. If an unexpected/malformed-block robustness case is retained, distinguish it explicitly from valid provider output and from required Chapter 1 behavior; do not build thinking support or add a new acceptance gate. The simplest in-scope correction is to keep the existing two text blocks as the success fixture. Do not remove their exact concatenation assertion or weaken the first-block mutation check.

This is not a request to change Engine's sensible `block.Type == "text"` filter (`internal/engine/engine.go:102–106`). It is the existing coding instruction's requirement that a fake and implementation not certify an invented protocol together (`.agents/skills/ensemble-coding/SKILL.md:196–217`). Deliberately malformed responses in `TestFailuresDoNotCommit` are correctly labeled failure cases and are not this finding.

## Production, architecture and scope assessment

All student citations in this section refer to the exact candidate above.

- **Public composition and ownership:** `ensemble.go:31–69` owns the logger and Agents, validates the root, constructs Agent then Engine and connects the actual relationship before exposing it. `internal/agent/agent.go:16–43` retains the Ensemble parent/config/history/Engine; Engine stores its actual Agent parent and reaches diagnostics through `Parent().Parent().Logger()` (`internal/engine/engine.go:17–35,112`). Shared vocabulary/interfaces are in `internal/common/common.go`; Agent and Engine implementations do not import one another. Required parent and registration relationships are not replaced by callback bags or duplicated logger/config state.
- **Public usability:** root aliases expose configuration, messages, usage and Agent access without an application importing private packages. Both tests use external `ensemble_test`, and CLI imports only the public package. History/Agent-list snapshots avoid leaking mutable backing slices. The sequential API is documented; no unrequested concurrency framework is present.
- **Turn semantics:** Agent stages the user turn and commits user/assistant only after success (`internal/agent/agent.go:45–60`). Engine does raw JSON/HTTP once, uses the supplied context, handles non-200/malformed/missing accounting as failures, concatenates text blocks, and preserves decoded measured cost on an empty answer (`internal/engine/engine.go:38–113`). Pointer usage fields distinguish absent values from measured zero. No retries, SDK, streaming, tools, persistence or speculative Observer/GUI subsystem was added.
- **Lifetime:** `http.Client{}` has no blanket model-response timeout, and `http.NewRequestWithContext` propagates caller cancellation. Source review found no longer arbitrary generation cutoff hidden beyond the 62-second probes. Setup/protocol harness watchdogs are not student generation deadlines.
- **Frontends:** `cmd/ch01/main.go:15–106` supplies one human chat and JSON-lines adapter over the same Agent. Environment-backed key/model and base-URL override, fixed system prompt, line reading, stderr diagnostics, stdout protocol isolation, usage at EOF, nonzero errors and no retry loop are implemented. Independent live use confirmed the human path, not just piped fake-server output.
- **Comments/KISS:** comments explain ownership, two-phase wiring, snapshots, nil-root guard, request lifetime and usage-before-empty-error ordering. Exported public/common API vocabulary has useful documentation. A small current-owner interface set, one shared config, one transport and two simple frontends are justified. No required simplification beyond R1 was found. Numeric size alone is not a reason to remove explanatory comments or error handling.

## Independent verification

Evidence is indexed in [ch01-code-review-evidence/README.md](ch01-code-review-evidence/README.md). These are reviewer runs unless labeled otherwise.

| Check | Result |
| --- | --- |
| Exact-candidate build (`go build -o <scratch>/ch01 ./cmd/ch01`) | Successful; binary hash recorded |
| `go vet ./...` followed by `go test -race -count=1 -v ./...` | **Exit 0**; package `ensemble` 125.261 s; lifetime subtests 62.00 s headers and 62.01 s body |
| `gofmt -l` on all candidate Go source/tests | No filenames reported |
| Unchanged grader with explicit absolute executable target, from main root | **Exit 0; `passed: true`; 100/100**; all seven scored checks and both zero-point `generation_headers`/`generation_body` gates true |
| Grader focused regression/control tests with race detector | **Exit 0**; lifetime pass/fail controls, historical 60-second relationship and preserved Chapter 1 checks passed |
| Nine reviewer mutations against final-commit disposable copies | All rejected by behavioral assertion/test failures, not compile failures |
| Independent real-provider PTY chat | Fresh model discovery, three successful turns, correct recalled name, EOF exit 0 |

The grader observed five requests for five turns, history 1 → 9 messages, and exact cumulative usage input **476**, output **65**. Both lifetime gates completed with a real 1m2s delay. A score of 100 alone was not used as acceptance.

Mutation checks independently killed first-block-only parsing, overwritten cumulative usage, omitted history, missing root registration, lost measured empty-answer usage, missing counters treated as zero, stdout pollution, ignored cancellation, and a one-minute HTTP timeout. The final timeout mutant failed after 62 seconds. These are nine specific sensitivity probes, not a claim of exhaustive mutation coverage. The initial reviewer mutation runner wrongly excluded `cmd/ch01` while excluding the built binary: six probes ran before its fixture-copy error. Its partial log/exit are preserved; the corrected runner completed all nine. That harness error is not attributed to the student.

The grader test selector `TestChapter1` also matched existing Chapter 10–19 tests. They passed incidentally; they are not used as Chapter 1 coverage. No short tool observation interval or 90-second overall shell deadline cut off either the 124-second suite waits or the grader's real delays.

### Live evidence and student-evidence distinction

Reviewer `live.py` discovered available models via real `GET /v1/models`, selected returned `claude-sonnet-5`, and drove the actual exact-candidate executable through a PTY. It introduced **Cobalt Wren**, received that name correctly on turn 2, and asked a context question on turn 3. Cumulative input was **78 → 219 → 430**, cumulative output **82 → 135 → 187**; Ctrl-D exited 0. This is independent live evidence, not the fake grader. The authorized credential field was loaded in memory, never printed or put in arguments, and only supplied to the intended chat child; output was sanitized before writing.

Separately inspected the student's final four-turn **Juniper Arc** transcript: final input/output **673/313**, EOF 0. The student's system-prompt experiment and nine saved mutations are explicitly from `60bf58a7728a96e27c2e13eea1afbf744889cb6e`. The final delta adds nil/zero-root rejection, its test and comments; it does not secretly carry the prompt experiment. I did not rerun that prompt variant and do not claim its saved transcript as my own observation. Its two terse/annoyed Waywardgeek replies support a bounded prompt-style observation, not generally improved coding or persistent learning. There was no separate friend observer; parent acknowledgement and the stated limitation are recorded. Model claims about persistence are not proof of provider retention practices.

## Grader alignment and demonstrated invocation feedback

Reviewed the unchanged `cmd/grade-edition3/main.go`, Edition 3 Chapter 1 overlay/tests, inherited Edition 2 Chapter 1 runner, `grade.Build`/harness/checks/report and `fakeanthropic` source at **`7f3d069420112c2136bc95d1d551e12fad6964c1`**. The fake splits text blocks and checks the narrow chapter wire/history/usage contract; the overlay adds both lifetime gates. The report requires every gate to pass. No newly demonstrated scoring mismatch warrants changing this already-accepted readiness grader.

**Nonblocking demonstrated author/build-documentation feedback:** `internal/grade/harness.go:68–120` explicitly accepts a prebuilt executable, but directory discovery tries only root or direct `./cmd/`. A public root library plus `cmd/ch01` is legitimate chapter architecture; `go build -o submission .` creates a non-executable archive, which subsequently produces `permission denied`. I independently repeated the chapter's exact directory command:

```sh
# cwd /Users/bill/projects/ensemble
go run ./cmd/grade-edition3 -ch 1 -json /Users/bill/projects/ensemble-edition-3-student
```

It failed before checks: shell/go-run **exit 1**, reported grader **exit status 2**. The passing run instead passed the actual compiled `cmd/ch01` executable as the explicit positional target to the same unchanged grader. That is supported input, not waived behavior or an imposed root-main layout. Record the recognized layouts and executable command beside the grading instructions in the later author phase; no manuscript/grader change is made here. The student's README already explains the workaround.

The assistance record also documents the parent's incorrect initial credential path and subsequent authorization correction. Neither clarification supplied implementation answer code. The missing-path attempt and directory failure remain failures in the record, not passing evidence.

## Historical comparison and chapter-local cost

Compared the complete promoted **second-edition Chapter 1** source in `solutions/ch01/` at the pinned main HEAD (last change/promotion `564e316f3f8aeac2ec26fe8edfbd8cdf42aa350b`), and the original **first-edition Chapter 1** source at `7eb3e247dccd553b096a450da124815aa167662e:solutions/ch01/main.go`. The latter was verified byte-identical to current `edition-1/solutions/ch01/main.go`; neither later-chapter nor retired-attempt code was substituted. The prior `docs/edition-2-notes/ch01-code-review.md` confirms both historical Chapter 1 implementations began from empty. The first-edition Chapter 1 path is also absent before its original addition.

Counts classify physical Go lines into active (including declarations/imports/braces), comment-only (leading `//`) and blank. Tests are separate; these source files have no block comments. Chapter 1 additions are therefore the same as totals, with zero Chapter 1 production removed at each edition's empty baseline. Promotion/archive moves are not charged as implementation growth.

| Category | Edition 1 | Edition 2 reference | Edition 3 candidate |
| --- | ---: | ---: | ---: |
| Production active / comment-only / blank | 199 / 48 / 31 | 281 / 59 / 29 | **304 / 74 / 41** |
| Production physical | 278 | 369 | **419** |
| Test active / comment-only / blank | 0 / 0 / 0 | 123 / 6 / 4 | **334 / 15 / 13** |
| Test physical | 0 | 133 | **362** |

Edition 3 adds **23 active production lines (8.2%)** and **50 physical lines** over the relevant Edition 2 reference, and 105 active/141 physical over the original single-file reference. Owner/API constraints explain much of the first-edition difference. Compared with Edition 2, current documentation, root validation, snapshots, request-context lifetime and missing-versus-zero accounting justify the modest net growth; Engine is actually smaller (88 versus 92 active lines). The old reference's `http.Client{Timeout: 60 * time.Second}` is specifically not a rule to preserve in Edition 3. Comment-only/active density is **74/304 = 24.3%**, near the chapter's approximate one-per-four aim without padding. Larger tests separately cover real request lifetime, cancellation, CLI behavior and failure/accounting semantics; they are not counted as production bloat. README (47 lines) and go.mod (3) are documentation/module overhead, not hidden production code. No SDK dependencies, generated sources, retained mutants or test framework were introduced into the candidate.

## Handoff and limits

Address **R1 only** for the initial correction, preserving the existing multi-text/history/usage/wire assertions and their sensitivity. Supply an exact clean revised commit, updated student-review provenance and rerun evidence for review. The passing production paths and grader alignment do not need speculative redesign. The documented directory-builder issue belongs to the later author feedback, not a new student architecture rule.

No student contact/coaching, source edits, manuscript changes, grader/reference changes, staging, commit, push, deployment, restart or GUI tools were used by this reviewer. Only this review and its scoped evidence directory were written in the main repository; mutations lived in disposable scratch copies. Unrelated/Codex files were left untouched. No code acceptance or author-loop authorization is given by this report.
