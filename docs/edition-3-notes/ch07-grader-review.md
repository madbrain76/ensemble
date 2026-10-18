# Edition 3 Chapter 7 — independent thin-adapter review

## Verdict: PASS

Accept this exact frozen **CLI adapter**, with no required source/test changes. It resolves the entry audit’s executable/core/GUI separation gap by calling the existing Chapter 7 API explicitly. This is not a Chapter 7 student grade, browser acceptance, or permission to launch a student. Root still owns receipt inspection, guide updates, and any later separately assigned implementation/review work.

No student Chapter 7 implementation exists. No source/manuscript/guidance/grader/reference/publication edits were made by this reviewer. Writes were restricted to this report, `ch07-grader-review-evidence/`, and owned `/tmp/cadence-ch07-adapter-review-independent` scratch. Prior entry evidence is unchanged.

## Frozen identity

Main HEAD: `23c75ede7944d9a912ae74d3f683354d163f0c86`.

| Input | SHA-256 |
|---|---|
| `cmd/grade-edition3/main.go` (5,341 bytes) | `1041fe571947f1a17d7e04160a363869c874e9264363c799e0b7d42afa9f06fe` |
| `cmd/grade-edition3/main_test.go` (10,957 bytes) | `55ef77d6191213a9cc11c30dacbe27e94d201c47e5909dbce1c73ac4266525da` |
| `/tmp/cadence-ch07-grader` (10,796,258 bytes) | `b72d12ccf29c9d56eeea0109710bb8930923085f9fd68cd80dcf211bf3b10826` |
| unchanged `internal/gradeedition2/ch07.go` | `2dc39413f88b1fd28059c84b52e9b48661be6438696f921fabe2ddc08d2ddb2d` |
| unchanged public Chapter 7 (692 lines) | `68d02d8e1170c2912676b8f162847f8fac25c87216514ed06adbd5f0ea011b47` |

Accepted student remains clean at `f595c787675e2a6343053ec42ea6c7e1e02a930e`, tree `6fb36168020a22b8993412812987338aa4b6e010`; `edition-3-ch06-r1` remains its accepted tag. See `input-pins.json` and `preservation.json` for full hashes, tag object, source pins, and contemporaneous preservation checks.

## Read and source review

Fully reloaded current AGENTS, coding/reviewer skills, Edition 3 execution/readiness, historical workflow/carryover, voice/writing procedure, and current Chapter 7, with bounded truncation recovery. Read the complete prior entry verdict/index and complete adapter source/tests/diff. `read-ledger.json` records actual read ranges, corrected endpoints, hashes, and line counts. Current Edition 3 student-first procedure overrides historical author-first instructions.

The production diff adds GUI flags, Chapter 7 root/model requirements and pre-build validation, and the exact API call:

```go
gradeedition2.Chapter7(bin, *core, *gui, *guiModel, *consumer, *workflow,
    *streamExercise, *streamModel, example, models)
```

The positional executable target is not substituted for the core root. GUI and inherited roots/models and parsed consumer/media maps are forwarded in the actual API order. A valid library root with a nested `cmd/ch01` executable target need not be flattened. Fresh positive measurement used a separately compiled executable and explicit source/module roots; it was not a new nested-layout build demonstration or a predecessor Chapter 7 grade.

The old `grade.Build`, chapter dispatches, report serialization/text output, and exit-status suffix are byte-identical after removing the one new `else if *chapter == 7` branch (`source-review.json`). Older-chapter GUI flag rejection handles explicitly supplied empty values as well as nonempty values. Existing consumer/media JSON readers retain their validation behavior. No grader algorithm, score threshold, hidden API, or student behavior requirement was added.

## Fresh independent command and guard measurements

Commands run sequentially; each real run has exact argv, cwd, duration, exit, environment controls, stdout/stderr hashes, score and `passed` in `*.run.json`, with matching safe stdout/stderr. No timeout wrapper canceled tests. Build/runtime environment used local toolchain, `GOPROXY=off`, owned TMPDIR, and removed ambient API-key/token variable names without reading or recording their values. Tests used existing fake local transports only.

Fresh checks:

- `go test -count=1 -race ./cmd/grade-edition3`: exit 0, 1.893 s.
- `go vet ./cmd/grade-edition3`: exit 0, 0.123 s.
- `git diff --check -- cmd/grade-edition3`: exit 0, 0.019 s.
- Independent Go sentinel compiled successfully. Its positive run actually wrote `EXECUTED` and exited 99, proving a launch would be observable.
- **64 argument cases × two targets = 128 exact frozen-binary invocations**. One target was the compiled sentinel; the other a module target with an observed `go` build command. All expected exits, **zero executable launches and zero builds**.
- Three additional positional-target guards (missing, empty, multiple targets) each exited 2 before build/launch. Total independent guarded invocations: **131**.

The 64-case matrix includes every required Chapter 7 flag missing/empty; all five roots absent, not a module, and containing nonregular `go.mod`; whitespace-only model identifiers; nonexistent, malformed, wrong-schema and directory-valued consumer/media files; Chapters 1/4/5/6 with either GUI flag explicitly empty/nonempty; unsupported chapters; unknown flag; and short/long/Chapter-7 help. Help exits 0 without launch, invalid inputs exit 2. Full stderr and argv for each case are in `validation-matrix.json` and `positional-guards.json`; no score or `passed` is fabricated for parser/build/setup cases.

## Fresh real Chapter 7 routing measurements

All use the exact frozen binary. The historical positive is the entry auditor’s privately owned matching `solutions/ch07` copy, not the frozen final Edition 3 solution and never student input. Its 66 source files were verified against the entry manifest. Core executable SHA-256 is `57e5518de31243502da58617a6eb3c1651b80b1db92d06aed3ad7a6a449fdeee`.

| Run | Exit | Score | passed | Seconds | Meaning |
|---|---:|---:|---|---:|---|
| `positive` | 0 | 100/100 | true | 13.794 | Real explicit CLI route, independent executable/core/GUI/inherited arguments |
| `wrong-gui-root` | 1 | 75/100 | false | 13.430 | Supplied alternate GUI binds to another core; actual dependency diagnostic |
| `nondefault-gui-model` | 0 | 100/100 | true | 13.408 | Supplied unlisted GUI model reaches plain fallback, observable diagnostic change |
| `wrong-stream-model` | 1 | 80/100 | false | 13.417 | Supplied inherited stream model reaches delivery checks and fails them |
| `custom-case-media` | 1 | 80/100 | false | 13.433 | Supplied consumer result and all three media model identifiers appear in actual inherited diagnostics |

Positive command (cwd `/Users/bill/projects/ensemble`):

```sh
/tmp/cadence-ch07-grader -json -ch 7 \
  -core /tmp/cadence-ch07-entry-auditor/positive \
  -consumer /tmp/cadence-ch07-entry-auditor/positive/ch05 \
  -workflow /tmp/cadence-ch07-entry-auditor/positive/ch06 \
  -stream-exercise /tmp/cadence-ch07-entry-auditor/positive/ch07 \
  -stream-model claude-sonnet-4-6 \
  -gui /tmp/cadence-ch07-entry-auditor/positive/gui \
  -gui-model claude-sonnet-4-6 \
  -consumer-case /tmp/cadence-ch07-entry-auditor/consumer-case.json \
  -media-models /tmp/cadence-ch07-entry-auditor/media-models.json \
  /tmp/cadence-ch07-entry-auditor/bin/positive-core
```

Forwarding distinctions:

1. **GUI root:** a private copy of the GUI had only its `go.mod` replacement redirected to a second private identical core. It compiled successfully first (exit 0, 0.780 s). Wire assertions passed, then `websocket-streams` explicitly reported that the GUI resolved `ensemble` to the alternate core rather than the supplied submission. This is a real module-binding distinction, not a failed server setup or timeout.
2. **GUI model:** an unlisted identifier produced 20 messages/32 log lines versus the positive’s 69/81; both passed. This is routing sensitivity to the existing plain-fallback behavior, **not** a behavioral rejection and not evidence that unknown models should be prohibited. No model declarations or application behavior were mutated.
3. **Stream model:** only the stream identifier changed. Existing inherited checks reported two text deltas instead of at least eight, one tool-call delta instead of at least three, and streaming/plain delivery equality. GUI checks still passed. The adapter did not merely parse and discard `-stream-model`.
4. **Consumer/media maps:** only parsed map/case inputs changed. Diagnostics named the intentionally different consumer result and each supplied unlisted Anthropic/OpenAI/Gemini media identifier. Both independent meanings reached the existing inherited graders in one finite run.

These are **adapter forwarding controls**, not a new application mutation campaign or replacement for the earlier replay/pause sensitivity audit. The original fixture was never edited in this review. The alternate GUI binding remains isolated in new scratch; nothing was applied back to student/reference/main source.

## Older-route preservation and carried evidence

| Fresh corrected route | Exit | Score | passed | Seconds |
|---|---:|---:|---|---:|
| Chapter 4 | 0 | 160/160 | true | 8.712 |
| Chapter 5 | 0 | 100/100 | true | 11.853 |
| Chapter 6 | 0 | 100/100 | true | 12.171 |

These use the same historical positive and appropriate inherited flags. They establish the selected old CLI routes still function; they are not fresh acceptance of the current student and do not repeat the Chapter 6 application/lifetime campaign.

Chapter 1 was **not rerun**. Its new-flag rejection was freshly exercised, its entire old dispatch/build/report/exit suffix is unchanged, and its existing Chapter 6 adapter-review receipt was fully inspected and pinned. That historical receipt ran `/tmp/cadence-ch06-grader -json -ch 1 /tmp/cadence-ch05-wakeup-evidence/candidate-r2/frozen/ensemble`, exit 0, 100/100, `passed=true`, 127.147 s, including separate real 62-second header and body checks. It is a prior accepted-predecessor receipt, **not** a fresh run on current `f595…`. `carried-ch1-receipt-pins.json` pins argv/timing/output and historical source/binary identity receipts; the stdout hash is `f5f3b049294ce0f224cdf86853eee291fdcf937ce67daf30deedc7fa4b73019d`. No genuine long HTTP gate was started and canceled.

## Failures, qualifications, and resource handling

- Initial combined instruction reads were capped; bounded rereads recovered them. Current historical carryover is 1,110 lines, and its tail was read after checking the endpoint.
- A source-inspection `git diff` entered its pager. Only that owned pager was terminated; no test was interrupted. Subsequent diffs disabled paging.
- Two guessed baseline filenames were absent; exact task paths were recovered. Publication comparison used only `/tmp/cadence-ch07-publication-preservation-after-body11.json`.
- A textual preservation helper first assumed a `switch` branch and failed. Corrected `else if` comparison succeeded; there was no source change or legacy regression.
- The first legacy Chapter 4 command incorrectly retained Chapter 5’s `-media-models`; it correctly returned usage/exit 2 and no report. This stopped the first sequence. That failed argv/stdout/stderr/metadata are retained as `legacy-ch4.*`; only the three appropriate legacy commands were then rerun. No prior positive or forwarding test was silently overwritten.
- No compile, watchdog, disk, or setup failure is credited as behavioral sensitivity. No hidden “all models must fail” inference is made from the passing GUI fallback control.
- Disk monitoring found no blocker (about 2.2 GiB available during early inspection; minimum captured run free space 3,468,152,832 bytes; final `df` about 3.2 GiB). Other concurrent disk changes are not attributed to this audit. No shared cache/tmp cleanup occurred.

## Boundaries and student evidence still required

The accepted entry diagnosis remains authoritative. This adapter does not close the grader’s coverage limits or change the public one-Agent Chapter 7 contract. Later student/reviewer proof still needs actor-owned History and stable display identity across responses/replay; atomic replay/active partial/live handoff with fixed last-hour OR100-event selection and nontruncated 270-event replay; actual content and from/to state; Agent-wide pause across future/following dispatch while active tools continue; explicit any-tab unpause without stale/local/reconnect veto; inherited hint OBSERVATION wake, continuation, and truthful paired undispatched handling; safe real DOM/readable-versus-opaque projection; real speech queue onerror/cancellation/late-callback behavior and typing cause; separately identified native-page and listening evidence. Existing historical links are not new student proofs. Tabs/tree/workspace remain Chapter 8; disk restoration identity remains Chapter 10.

No live GUI at 8088, browser navigation/refresh, service restart/deploy, provider request, key access, child agent, publication action, staging/commit/tag/push, or new student launch occurred. Only existing grader-owned local servers were exercised.

## Preservation, evidence, and completion

`preservation.json` confirms all twelve current publication-path hashes against the supplied body-11 baseline; all unchanged entry source pins; exact frozen candidate hashes; clean accepted student/tag/tree; and the unchanged 76-file prior entry evidence set plus its report and checksum manifest. The older ten-file publication baseline was not misclassified as drift. Source/manuscript/guidance/grader inputs outside the frozen adapter remain untouched.

Portable evidence lives in `ch07-grader-review-evidence/`. `INDEX.md` enumerates its exact set. `SHA256SUMS` covers this report and all evidence except its own nonrecursive self-hash, supplied in the final handoff. All portable files are UTF-8 text: safe command reports, harness source, minimal binding/source diffs, and inventories. No binaries, archives, pyc, full historical solution payloads, raw private traffic traces, or keys are included. Private binary/source locations and hashes are explicitly local-only.

All grading/build commands have exited, all test job results were collected, and no children or pending tests remain. **PASS applies only to the frozen thin adapter.**
