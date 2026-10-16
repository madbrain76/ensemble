# Edition 3 Chapter 1 — correction-round-1 code re-review

## Disposition: accepted

**Accepted candidate:** `1aad3975842d502b9144b7b9bddfd30ca7309db7` in `/Users/bill/projects/ensemble-edition-3-student`; independently verified clean.

**Initially reviewed candidate:** `d2baff64055ac33702a91c6d6ee5a4246e84c0c9`.

**Unchanged chapter/grader revision:** `7f3d069420112c2136bc95d1d551e12fad6964c1`. Chapter SHA256 remains `ab016baaff3f6a9d83ecffcbfc3772cd692c9d33c06d9fa1c52ad9a8708d25df`.

This is the first student correction/re-review round, limited to initial finding R1 and regressions from its correction. **R1 is closed; no required finding remains. No acceptance exception is needed.** The initial production/architecture/public-API/comment/scope/reference review remains applicable through independently verified source identity. This is code-review acceptance, not owner sign-off, text review, author revision, or authorization to begin another chapter.

Initial report preserved unchanged: [ch01-code-review-r1.md](ch01-code-review-r1.md), SHA256 `0649afca4480df50ba929602a5574122e17c565a8d6294dade3ddf9db6e3d314`. Its historical needs-correction disposition is resolved by this new report, not erased. The earlier correction of the reviewer's transcription was not a student correction round.

## R1 closure — exact change and retained sensitivity

At **`1aad3975842d502b9144b7b9bddfd30ca7309db7:ensemble_test.go:54–55`**, the successful fixture now says:

```go
// Two text blocks must be concatenated exactly, without added separators.
fmt.Fprint(w, `{"content":[{"type":"text","text":"hello "},{"type":"text","text":"世界"}],"usage":{"input_tokens":17,"output_tokens":9}}`)
```

The invented thinking block is removed. This is a faithful minimal successful response for the current text-only Chapter 1 contract; it does not introduce thinking support or a new malformed-input requirement. The comment now describes precisely the behavior being exercised.

I independently compared complete Git trees and every tracked file at both revisions, not only the supplied patch. The path set is unchanged. **Only `ensemble_test.go` differs, and only lines 54–55 differ.** Every other test line, including all assertions, remains byte-identical. Specifically, in the revised source:

- Lines 37–51 retain real method/path/header and request-configuration checks.
- Lines 57–65 retain ownership and configuration/root snapshot checks.
- Lines 68–80 retain the four-turn input sequence, exact `hello 世界` answer, full history replay and committed-history checks.
- Lines 82–99 retain request count, cumulative usage, history snapshot, parent logging and Agent isolation checks.

The two text blocks retain their exact original values, including the trailing space in `hello `; usage remains input 17/output 9. No assertion was weakened to match an implementation.

**Independent sensitivity proof:** the revised `TestConversation` passes. In a separate disposable copy, inserting `break` immediately after `text.WriteString(block.Text)` makes that same test fail at `ensemble_test.go:71` with `all text blocks: "hello " <nil>`, exit **1**. This is the intended behavioral failure, not a compilation/setup failure. Restoring the Engine file byte-for-byte makes the focused test pass again, exit **0**. All temporary-copy tracked files were checked against the exact candidate after restoration. Neither the student repository nor its production archive was mutated.

## New independently run checks

I exported `git archive 1aad3975842d502b9144b7b9bddfd30ca7309db7` into a disposable directory and ran:

| Check | Actual result |
|---|---|
| `go build -o "$SCRATCH/ch01" ./cmd/ch01` | Exit 0 |
| `go vet ./...` | Exit 0 |
| `go test -count=1 -v -run '^TestConversation$' .` | Passing control, exit 0 |
| Same focused test with first-block mutant | Intended assertion failure, exit 1 |
| Same focused test after restoration | Pass, exit 0 |
| `go test -race -count=1 -v ./...` | Pass, exit 0; `ensemble` elapsed **125.239 s** |
| Real `TestRequestLifetime/headers` | Pass, **62.00 s** |
| Real `TestRequestLifetime/body` | Pass, **62.01 s** |
| `gofmt -l` on all archived Go files | Exit 0, no filenames emitted |

The complete race suite also passed CLI protocol/chat/error/empty-input checks, public conversation behavior, failure/history/usage handling, cancellation, config validation, measured zero, and invalid-root tests. No short shell deadline interrupted the lifetime checks. The verification runner completed with exit **0** and performed final exact-archive, clean-tree and artifact-preservation checks.

New reviewer evidence and runnable verification script: [ch01-code-review-correction-1-evidence/README.md](ch01-code-review-correction-1-evidence/README.md). The expected mutant failure remains explicit in its own log/exit file; it is not a failure of the unmodified candidate.

## Student correction evidence independently checked

I read [ch01-student-correction-1.md](ch01-student-correction-1.md) and its exact patch, provenance, focused control/mutation/restoration logs, vet/full race-suite outputs and exits, build metadata, and hygiene record. The saved student full suite passed in **125.327 s** with both actual 62-second waits. That is student evidence, distinct from the independent **125.239 s** run above.

I independently recomputed all **21** original-artifact hashes in the student's preservation record, including `ch01-candidate.json`, and both revision hashes for all **nine** tracked paths in its provenance. They match. The student's build metadata records revised VCS revision `1aad3975842d502b9144b7b9bddfd30ca7309db7` and `vcs.modified=false`. The reported dispatch hold is not a code/test failure and required no runtime repair.

## Production identity and prior evidence attribution

All five production Go files and `go.mod` are byte-identical between the initial and corrected candidates:

- `ensemble.go`
- `internal/common/common.go`
- `internal/agent/agent.go`
- `internal/engine/engine.go`
- `cmd/ch01/main.go`
- `go.mod`

`README.md` and `cli_test.go` are also unchanged. Full before/after SHA256s are recorded in the new `provenance.json`. This is **source-byte identity**, not an assertion that binaries built at different VCS revisions/paths have identical bytes.

The following earlier independent results remain attributed to **`d2baff64055ac33702a91c6d6ee5a4246e84c0c9`**, not relabeled as new correction-candidate runs:

- Unchanged grader against the actual prebuilt `cmd/ch01` executable: **exit 0, `passed: true`, 100/100, all seven scored checks and both zero-point lifetime gates passed**.
- Separate real-provider three-turn PTY exercise with fresh model discovery, Cobalt Wren recall, cumulative input/output 430/187 and EOF exit 0.
- Initial grader controls, nine behaviorally rejected mutations, and production/historical-reference review.

No grader rerun, provider calls or credential access occurred in this round. Their omission is deliberate and justified by the two-line test-only delta and complete production/module identity; it is not an unreported failed run or an acceptance exception.

The earlier directory-target permission failure also remains explicit historical evidence. The grader's existing supported executable-target invocation passed unchanged checks. The narrow build-layout documentation feedback in the initial report still stands for the later author process; it does not require a production-layout change or grader weakening, and no new alignment mismatch was introduced here.

## Preservation and phase boundary

The verifier rechecked **64 protected original/student artifacts**, including the initial review and all its evidence, original student review/evidence and candidate record, and the new student correction record/evidence. All remained unchanged. Pinned chapter/grader hashes and HEAD also remained unchanged. The student working tree remains clean at the accepted revision.

Only this new correction review and `ch01-code-review-correction-1-evidence/` were written in the main repository; mutations/builds ran in disposable scratch directories. No student source, book, grader, historical reference or unrelated files were edited. No student coaching/contact, provider spend, runtime repair, push, deployment, restart, broad staging or GUI use occurred.

**Final result: accepted at `1aad3975842d502b9144b7b9bddfd30ca7309db7`; R1 closed; one of the maximum three correction/re-review rounds used. Author phase has not started.**
