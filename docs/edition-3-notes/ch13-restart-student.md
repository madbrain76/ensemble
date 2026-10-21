# Chapter 13 restart — student

## Current status — released R4 validation passed; HOLD for remaining exercises

**The released changed-input gates passed; source and frozen CLI remain on HOLD.** Fresh GUI race/vet exited 0. Chapter 13 exited **0**, parsed **`passed:true`**, **100/100**; the first Chapter 1 run on the same verified CLI also exited **0**, parsed **`passed:true`**, **100/100**, including separate genuine 62-second header/body waits and complete usage. Unchanged root/driver/application/JS gates were carried, not replayed. Earlier failed receipts remain failed and untouched. Native audibility and real-model blind-listener exercises are still unperformed. An owned controller is now prepared, pinned and syntax/pure-guard checked, but both operational phases remain entirely unexecuted on HOLD for root review/correction and a separate execution release.

## Round 1 held record

**Production correction Round 1/10 complete; source on HOLD for root + same independent reviewer.** The initial compile-only candidate was reviewed and this round addressed the six public correction items without changing the Agent/History/job architecture. No private reviewer/grader material, original solutions, earlier snapshots, children, commits, installs, live-server restart, or manuscript/guidance edits were used.

Assistance: root clarified that the current Chapter 13 §13.9 contract takes precedence over its mistaken initial wrapper-variable brief. Root also authorized retaining the single `pause_while_speaking:true` preference in the older controlled-JS fixture after independently checking the actual defaults and unchanged cancellation assertion. Both are within Round 1, not additional correction rounds.

## Ownership and corrections

- Added executable root `tts-harness.sh`. Normal invocation is `./tts-harness.sh "PROMPT" [--type-during-turn "TEXT"]`, using `LLM_BASE_URL`, `WORKSPACE`, `TTS_LOG` (a **file**), and optional `STREAMING=off`. `--describe` documents the required fields and separate optional attach/JSONL/listener entry. Source/module/assets paths are resolved independently of caller cwd. The launcher uses the real GUI, an installed private Chromium, and an asynchronous synthesizer fake; it does not create a target Engine or fake narrator. Runtime process/profile/log files are owned separately from the live GUI. **Normal browser launch has not been executed.**
- Hub admission now includes failed-tool error cards while successful tool-result bodies remain automatically silent. The existing provider-error event path remains. Thought/error phrases get audible prefixes at the shared page channel. Durable Speak, full readable tool content, and opaque exclusion remain in their existing owners.
- Unchanged retained partial tails flush at actual subsequent part/source boundaries. Raw-offset speech cursors remain responsible for preventing later-final duplication. Provider/tool error boundaries use the same ordering.
- Cancelled tasks discard pending segments for the cancelled task's turns rather than immediately restarting them. Native callbacks are fenced by the active utterance as well as task identity; stale callbacks cannot advance the next chunk.
- Page transcript recording now happens at per-segment admission, not when all future task segments are enqueued. Quiet/audio-suppressed paths still admit hearable segments there. Cancellation/completion records remain visible without disclosing pending text. JSONL includes page sequence, source, turn identity, timing, utterance, and speech-channel audio marker; it is not native/hearing proof.
- Driver submission/wait tracks selected Agent, observed prompt, turn ID, and public terminal outcome rather than a global ended counter. Error, interruption, timeout, success, and loss of ownership are distinct. EOF and listener completion settle owned work or report failure/incompleteness. Draft cleanup checks input revision/text and cannot clear a subsequent human edit. Final capture/write errors propagate. Listener tools remain exactly `hear`, `type`, `submit`, `wait`, `sleep`; no DOM/eval/target-Engine capability is exposed to its model.

## Actual touched paths

Round 1, under `solutions/edition-3/main`:

- `gui/hub.go`, `gui/projection.go`, `gui/speech.go`, `gui/web/gui/app.js`
- `gui/client_round1_test.cjs`, `gui/ch08_round1_test.cjs`
- New `gui/ch13_correction_test.go`, `gui/ch13_correction_test.cjs`
- `examples/ch13-speech/main.go`; new `driver_test.go`, `launch.go`, `launch_test.go`
- New `tts-harness.sh`

Initial held work retained: `gui/server.go`, `gui/settings.go`, `gui/web/gui/index.html`, `gui/speech_text.go`, `gui/speech_text_test.go`, and the speech consumer's `go.mod`, `go.sum`, `main_test.go`. Also this authorized student note. No root adapter or other role's notes were edited.

## Before/after and setup evidence

Pre-fix focused checks demonstrated: failed-tool error speech skipped; unchanged partial tail overtaken; pending speech promoted after cancellation; future transcript text disclosed before admission; absent audible thought/error labels; duplicate same-task callbacks advancing chunks; and driver error/interruption reported as successful completion. Those were behavioral failures, not compile success claims.

Setup failures were kept distinct: early test attempts used incorrect existing API/field names; the plain local-provider fixture initially omitted required usage data; a VM fixture lacked `URLSearchParams`/timer setup; and the process test initially mistook signal termination for failure to reap. These were corrected in test setup. The actual-History test also exposed **two distinct existing public provider-error events**, not double-processing of one event; the test now checks admission against those actual events rather than imposing deduplication.

The older JS suite initially exited **1**. Two assertions encoded superseded behavior (Read Aloud clearing HUMAN typing, and disabled audio skipping channel admission); those assertions were updated narrowly to Chapter 13's contract. A subsequent captured run still had **`js_existing_exit=1`** because its cancellation test omitted the new pause-while-speaking preference. The enclosing shell's zero was **not** an aggregate pass. After the separately authorized fixture preference was supplied, the unchanged stale-callback/unpause assertion and the suite passed (`js_existing_exit=0`). The earlier failure is not erased by that rerun.

## Executed checks

All checks used the same isolated `env -i` setup: owned HOME/TMPDIR, installed Go/Node plus `/Users/bill/go/bin`, `GOMAXPROCS=2`, `GOFLAGS='-mod=readonly -p=1'`, `GOTOOLCHAIN=local`, `GOPROXY=off`, `GOSUMDB=off`, `GOMODCACHE=/Users/bill/go/pkg/mod`, `GOCACHE=/Users/bill/Library/Caches/go-build`. No real keys or real `LLM_BASE_URL` were inherited. Only released loopback test servers/transports were used. Disk remained about 22 GiB free. The single owned directory is `/tmp/ensemble-ch13-student.Vb0IEn` (about 160 MiB including initial and Round 1 compile outputs); no source/cache copies or shared-cache cleanup.

Final focused commands and exits (sequential; no full test campaign):

```
(cd gui && go test -vet=off -count=1 -run '^Test(Ch13.*|Speech.*|Ch7R1WebSocketControls|Ch7R1SpeechRetainsAllAcceptedPhrases|Ch7R1SpeechCancelClaimAndOtherTabRelease)$' .) # 0
(cd gui && node ch13_correction_test.cjs) # 0, ten named controlled-production-JS checks
(cd gui && node ch08_round1_test.cjs) # final 0, after the failures described above
(cd examples/ch13-speech && go test -vet=off -count=1 -run '^Test(TargetFailureIsNotCompletion|DriverSettlement|JSONLEOFAndFinalCapture|TranscriptWriteErrorAndUnsubmittedDraft|NormalHarnessWiring|OwnedProcessStops|ReadinessReportsOwnedExit|NormalGuardBeforeLaunch)$' .) # 0
```

These exercised actual Go owners, a local plain-provider test, the existing GUI WebSocket controls, a local CDP-transport fixture, controlled production JS, and owned shell-process cleanup. They did not exercise an actual browser or speech device. No zero-selection result is offered as evidence.

From the unrelated owned `/tmp` cwd, the absolute `tts-harness.sh --describe` exited **0**, and JSON validation of the five §13.9 description fields exited **0**. Normal-entry missing-input guard exited **1 as expected** (`WORKSPACE is required`) before GUI/browser launch.

Final explicit-module compilation, with `$out=/tmp/ensemble-ch13-student.Vb0IEn/round1-bin`, all exits **0**, aggregate **0**:

```
go build -o "$out/ch01" ./cmd/ch01
(cd gui && go build -o "$out/gui" ./cmd)
(cd ch06 && go build -o "$out/" ./...)
(cd ch07 && go build -o "$out/" ./...)
(cd examples/ch04-consumer && go build -o "$out/ch04-consumer" .)
(cd examples/ch13-speech && go build -o "$out/ch13-speech" .)
(cd examples/ch13-speech && go test -vet=off -c -o "$out/ch13-speech.test" .)
(cd gui && go test -vet=off -c -o "$out/gui.test" .)
go test -vet=off -c -o "$out/main.test" .
node --check gui/web/gui/app.js
sh -n tts-harness.sh
```

## Remaining checks / limits

No private browser/Chrome launch, actual normal-entry provider exercise, native speech/sound, blind listener, real model/key/discovery, grading, race/vet, npm installation, microphone, OS focus/volume, or live-8088 restart was performed. The launcher is compiled and its pure wiring/guards/process helper checked, but actual browser startup, trusted input, synthesizer injection, complete runtime teardown, and external-provider end-to-end behavior remain unverified. Native/hearing and listener comprehension still require separately released exercises. Source is held; no further production changes or tests are underway.

## Round 2/10 correction record

Scope and ownership:

1. Removed the launcher's unsupported `tts_delay_ms` key; it now uses existing GUI settings/defaults. No GUI settings machinery was added. `TestCh13LauncherSettingsLoad` extracts the actual launcher payload and passes its file through `NewServerWithSettings`, including the strict settings-load owner.
2. Split target observation from interactive success. Normal one-shot mode accepts an observed/captured target `error` as successful harness observation, records the actual target outcome, and still fails on interruption, timeout, transport, capture, or cleanup failure. Interactive `wait` and listener settlement still require target success. Speech text remains exclusively from page replies; no narrator or fabricated utterance was added.
3. One stdin peer owner now starts before interactive runtime startup. Its independent reader cancels the shared context on EOF while JSONL sleep/wait, listener Ask, or a page request is active. Normal one-shot mode deliberately has no stdin peer. The owner uses a private nonblocking stdin duplicate so a pending read can be closed even when the original inherited descriptor remains open. EOF does not clear human input authority; existing draft/turn ownership checks remain. Blocked page cancellation closes only that owned connection; resulting capture/transport loss remains a failure, not successful observation.
4. `TestPublicContractJSON` now checks actual `--describe` capabilities and the exact five listener tool names, not an absent invented Version field. It was explicitly selected in pre-fix and corrected runs. The old immediate-EOF test was changed to the clarified interactive lifetime contract; it is not evidence that the old settlement-on-EOF behavior was valid.

Round 2 touched paths, relative to `solutions/edition-3/main`:

- `examples/ch13-speech/launch.go`, `main.go`, `main_test.go`, `driver_test.go`
- New `examples/ch13-speech/peer.go`, `round2_test.go`
- New `gui/ch13_harness_settings_test.go`
- This existing authorized student note

No other production source, GUI speech/typing logic, core ownership, modules/dependencies, wrapper, guidance, manuscript, or reviewer material changed this round. Assistance: only the public Round 2 correction brief; no new scope question or delegated worker.

### Actual failure/pass evidence

Evidence directory remains `/tmp/ensemble-ch13-student.Vb0IEn`; these are small execution outputs, not a source/cache copy or audit bundle.

- `r2-before-settings.log`, exit **1**: actual settings load rejected `tts_delay_ms`.
- `r2-before-harness.log`, exit **1**: explicitly selected `TestPublicContractJSON` failed on absent Version; normal target-error observation failed; EOF did not stop active sleep, wait, or listener Ask promptly. Normal closed-stdin observation and fatal capture failure already passed. Those original failures remain failures.
- `r2-after-settings.log`, exit **0**: the same actual settings-load boundary passed after removing the key. Later peer-only changes did not change that input/result.
- `r2-after-harness.log`, exit **0**: the original correction selections passed after the bounded implementation.
- `r2-lifetime-regressions.log`, exit **1**: additional directly relevant lifecycle checks found inherited process stdin could keep its reader blocked during shutdown. The bounded test killed and reaped only its owned helper; this was a genuine lifecycle failure, not a passing test or setup failure.
- `r2-final-focused.log`, exit **0**: the same real-process reader check passed after the private pollable-duplicate fix, along with strict interactive outcomes, immediate EOF, finite sleep while input stayed open, EOF during active sleep/wait/listener Ask/page I/O, normal closed stdin, and fatal transcript-write/capture errors.
- `r2-listener-qualified.log`, exit **0**: reran only the changed listener test qualification. The blocked loopback request was required to be a chat-completion POST, not discovery/configuration. It used the actual listener Agent Ask and observed request cancellation. The test supplied no model response and made no external provider/model call.

One scripted source edit stopped on a mismatched anchor before writing the main-session change. It was corrected before behavioral reruns; it was not counted as a test pass. There were no dependency/setup-download blockers or zero-selection passes.

Commands (under the same clean offline environment recorded above):

```
(cd gui && go test -vet=off -count=1 -v -run '^TestCh13LauncherSettingsLoad$' .)
# before 1; corrected 0

(cd examples/ch13-speech && go test -vet=off -count=1 -v -run '^Test(PublicContractJSON|NormalOneShotErrorObservation|JSONLPeerEOFDuringOperation|ListenerPeerEOFDuringAsk|NormalClosedStdinStillObserves|NormalCaptureFailureStillFails)$' .)
# before 1; corrected 0

(cd examples/ch13-speech && go test -vet=off -count=1 -v -run '^Test(TargetFailureIsNotCompletion|DriverSettlement|JSONLImmediateEOFStopsSession|TranscriptWriteErrorAndUnsubmittedDraft|PeerOperationsWhileOpen|PeerEOFDuringPageCall|TargetOutcomeWriteFailure|PeerStopsReaderWithOpenProcessStdin)$' .)
# reader shutdown failure 1, preserved in r2-lifetime-regressions.log

(cd examples/ch13-speech && go test -vet=off -count=1 -v -run '^Test(PublicContractJSON|NormalOneShotErrorObservation|JSONLPeerEOFDuringOperation|ListenerPeerEOFDuringAsk|NormalClosedStdinStillObserves|NormalCaptureFailureStillFails|TargetFailureIsNotCompletion|DriverSettlement|JSONLImmediateEOFStopsSession|TranscriptWriteErrorAndUnsubmittedDraft|PeerOperationsWhileOpen|PeerEOFDuringPageCall|TargetOutcomeWriteFailure|PeerStopsReaderWithOpenProcessStdin)$' .)
# final 0; named selections and subcases are in r2-final-focused.log

(cd examples/ch13-speech && go test -vet=off -count=1 -v -run '^TestListenerPeerEOFDuringAsk$' .)
# qualified request-path rerun 0
```

Final compilation reused `$out=/tmp/ensemble-ch13-student.Vb0IEn/round1-bin` rather than creating another binary set. `r2-compile.log` records every command exit **0**, aggregate **0**:

```
go build -o "$out/ch01" ./cmd/ch01
(cd gui && go build -o "$out/gui" ./cmd)
(cd ch06 && go build -o "$out/" ./...)
(cd ch07 && go build -o "$out/" ./...)
(cd examples/ch04-consumer && go build -o "$out/ch04-consumer" .)
(cd examples/ch13-speech && go build -o "$out/ch13-speech" .)
(cd examples/ch13-speech && go test -vet=off -c -o "$out/ch13-speech.test" .)
(cd gui && go test -vet=off -c -o "$out/gui.test" .)
go test -vet=off -c -o "$out/main.test" .
node --check gui/web/gui/app.js
sh -n tts-harness.sh
```

Disk was still approximately 22 GiB free. Same owned HOME/TMP, offline module settings/shared caches, and no inherited real keys or `LLM_BASE_URL`. No browser/Chrome/native/audio, external provider/model/listener exercise, grading, frozen candidate CLI, race/vet campaign, install, discovery, child agent, live restart, commit/snapshot, or shared cleanup occurred. The helper process and blocked loopback transport were only the narrowly released lifetime tests.

Remaining uncertainty: normal launcher/browser startup and the actual public provider-error exercise still require separately released execution. Transport-fixture capture and settings-owner tests do not prove browser speech, and blocked listener Ask cancellation does not prove hearing or comprehension. Source is now held; no further production changes or checks are underway.

## First full local validation — stopped at Chapter 13

R2 source/focused evidence was independently accepted **for validation**, not chapter acceptance. No source or fixture edits were made during this release. New owned directory: `/tmp/ensemble-ch13-validation.At5CrJ`. Same clean offline Go/shared-cache environment, fresh owned HOME/TMP, no real keys or inherited `LLM_BASE_URL`. Disk floor was checked before each executed gate; available space stayed approximately 22 GiB. No dependency installation or manual artifact/shared-cache cleanup.

Local gates, sequential, all exit **0** (`local-exits.log` and the correspondingly named `.log` files):

- Root, GUI, ch13-speech: each ran `go test -race -vet=off -count=1 ./...`, then `go vet ./...`.
- ch06, ch07, ch04 consumer: each ran `go test -vet=off -count=1 ./...`, then `go vet ./...`. ch07 and ch04 reported **no test files**: compilation evidence only, not behavioral passes.
- `node gui/ch13_correction_test.cjs` and `node gui/ch08_round1_test.cjs` (executed from `gui` with local filenames), once each, exit **0**. No new JS campaign.
- Local batch shell exit **0**. Individual logs: `root-race`, `root-vet`, `gui-race`, `gui-vet`, `speech-race`, `speech-vet`, `ch06-tests`, `ch06-vet`, `ch07-tests`, `ch07-vet`, `ch04-consumer-tests`, `ch04-consumer-vet`, `js-ch13`, `js-existing`.

One frozen CLI was then built with `go build -o /tmp/ensemble-ch13-validation.At5CrJ/ch01-frozen ./cmd/ch01`, exit **0**, and made non-writable. SHA-256:

`1bb274f4bc9037b128c2886ee73732b0bb09c3667f067ddef9fbd83da9ebe34b`

`frozen.sha256` and `frozen-build.log` preserve that receipt. Public grader `-h` exited **0**; only command-line help was read, not grader source. Owned `media-models.json` used exactly root's map: anthropic=`claude-sonnet-4-6`, openai=`gpt-5.4`, gemini=`gemini-flash-3.8`.

From book root `/Users/bill/projects/ensemble`, the **single** authorized full grade was:

```
go run ./cmd/grade-edition3 -ch 13 -json \
  -core /Users/bill/projects/ensemble/solutions/edition-3/main \
  -consumer /Users/bill/projects/ensemble/solutions/edition-3/main/examples/ch04-consumer \
  -workflow /Users/bill/projects/ensemble/solutions/edition-3/main/ch06 \
  -stream-exercise /Users/bill/projects/ensemble/solutions/edition-3/main/ch07 \
  -gui /Users/bill/projects/ensemble/solutions/edition-3/main/gui \
  -stream-model claude-sonnet-4-6 -gui-model claude-sonnet-4-6 \
  -media-models /tmp/ensemble-ch13-validation.At5CrJ/media-models.json \
  /tmp/ensemble-ch13-validation.At5CrJ/ch01-frozen
```

Result: process/shell exit **1**, emitted JSON **`passed:false`**, score **0/100**. Raw evidence: `ch13.stdout.json`, `ch13.stderr` (`exit status 1`), `ch13.exit`, `ch13.started`, `ch13.finished`, all in the owned validation directory.

The public report's speech cases all report harness **`page command read failed`**, with provider-request checks showing **0**. These are failed startup/runtime observations; empty transcripts do not establish filter, ordering, or native-speech behavior. Separate parity detail: `gemini supported image ... model "gemini-flash-3.8" does not support image media`. Root supplied that exact mapping and directed that it be treated as an **unresolved setup mismatch**, not automatically a production defect. No cause beyond the public diagnostics was established.

Stopped immediately after the failed grade; only its public JSON was read, then root was asked for direction. Root confirmed stop/hold. **No Chapter 1 run** (its prerequisite did not pass), no grade retry, separate browser rehearsal/diagnosis, source/fixture fix, or native/listener proposal. The authorized grader-triggered isolated browser/fake-provider sessions are the only browser sessions attempted in this release; no native voice/audio, microphone, external real provider/key/model, live 8088/shared browser, or discovery exercise.

Cleanup qualification: the grader invocation returned, and the public report records harness subprocess exits. No student manual cleanup or artifact deletion was performed. Existing owned artifacts were left untouched. Browser-child teardown/grader-managed cleanup was **not independently verified** after the failure; no such success is claimed. Source remains frozen on HOLD for root and the same reviewer's inspection of the existing evidence. Earlier R1/R2 failure receipts and qualifications remain unchanged.

## Round 3/10 — narrow launch/diagnostic correction, compiling HOLD

Public findings supplied by root + retained reviewer: held launch reached `Page.enable`/script injection but timed out in `Page.navigate` before a provider request (16.62 s run). Their one-variable `--use-mock-keychain` comparison, with the same other flags/GUI/controlled synth, completed navigation, trusted submission, one local dummy POST, error-channel capture and exit-0 cleanup in 1.94 s. These are **supplied observations**, not student reruns. They support compatibility on this host, not a proven OS-keychain cause or universal portability. No grader defect is asserted. Root corrected its model-map spelling to `gemini-3.8-flash`; the old failed map/grade receipt remains untouched and no model table changed.

Actual edits only:
- `examples/ch13-speech/launch.go`: add `--use-mock-keychain` to the existing owned headless/profile launch; all other flags, timeouts, cleanup and controlled synth retained.
- `examples/ch13-speech/main.go`: annotate command errors with the CDP method; retain underlying read/write causes with `%w`. No parameter or protocol error-body dump. Existing cancellation, target, transcript and draft owners unchanged.
- New `examples/ch13-speech/page_test.go`: focused local-WebSocket checks for read/write causes, method/stage attribution, parameter/protocol-body exclusion, and cancellation cause. This is the explicit diagnostic improvement requested by root, not a newly invented prior blocker.
- This student note. No other source/dependency/framework changes.

Evidence under existing owned `/tmp/ensemble-ch13-validation.At5CrJ/r3`, with fresh owned HOME/TMP, same clean offline flags/shared caches, no real keys/inherited `LLM_BASE_URL`; disk checks stayed about 22 GiB free:

```
# From examples/ch13-speech; new diagnostics against held generic errors:
go test -vet=off -count=1 -v -run '^TestPageCommand(FailureMetadata|CancelledCause)$' .
# exit 1, expected missing method/cause failures; diagnostics-before.log

# After the bounded correction:
go test -vet=off -count=1 -v -run '^Test(PageCommandFailureMetadata|PageCommandCancelledCause|NormalCaptureFailureStillFails|PeerEOFDuringPageCall)$' .
# exit 0; focused.log

go build -o /tmp/ensemble-ch13-validation.At5CrJ/r3/ch13-speech .
go test -vet=off -c -o /tmp/ensemble-ch13-validation.At5CrJ/r3/ch13-speech.test .
# both exit 0; build.log, test-compile.log; aggregate 0 in exits.log
```

No setup failures. Reused the supplied real launch comparison and earlier unchanged-module validation; no new Chrome/native/fake-provider campaign, grade, full race/vet, application-module replay, frozen-CLI rebuild, key/discovery activity, child agent, install, cleanup, commit or snapshot. No private diagnostic source/fixtures were inspected. Source held after compiling; changed-input validation requires root + the same reviewer's release. All earlier failed receipts remain failed.

## R3 changed-input validation — stopped, no production edits

Fresh output subdirectory: `/tmp/ensemble-ch13-validation.At5CrJ/r3-validation`, with owned HOME/TMP, same clean offline settings/shared caches, no real keys/inherited `LLM_BASE_URL`, and a floor check before each gate (about 22 GiB free). Only the changed module was replayed:

```
(cd examples/ch13-speech && go test -race ./...)  # exit 0, 2.621 s
(cd examples/ch13-speech && go vet ./...)       # exit 0
shasum -a 256 /tmp/ensemble-ch13-validation.At5CrJ/ch01-frozen  # exit 0, exact match
```

Frozen CLI SHA-256 remains `1bb274f4bc9037b128c2886ee73732b0bb09c3667f067ddef9fbd83da9ebe34b`; no rebuild. Local batch shell exit **0**. Logs: `speech-race.log`, `speech-vet.log`, `frozen-hash.log`, `exits.log`. Root/GUI/inherited-module and JS results were carried forward, not rerun.

New `media-models-v2.json` contains anthropic=`claude-sonnet-4-6`, openai=`gpt-5.4`, gemini=`gemini-3.8-flash`. Original wrong map and failed grade remain untouched. Exactly one new grade, from book root `/Users/bill/projects/ensemble`, with `$root=/Users/bill/projects/ensemble/solutions/edition-3/main` and `$work=/tmp/ensemble-ch13-validation.At5CrJ/r3-validation`:

```
go run ./cmd/grade-edition3 -ch 13 -json \
  -core "$root" -consumer "$root/examples/ch04-consumer" \
  -workflow "$root/ch06" -stream-exercise "$root/ch07" -gui "$root/gui" \
  -stream-model claude-sonnet-4-6 -gui-model claude-sonnet-4-6 \
  -media-models "$work/media-models-v2.json" \
  /tmp/ensemble-ch13-validation.At5CrJ/ch01-frozen
```

Process/shell exit **1**; public JSON **`passed:false`**, **95/100**. Evidence: `ch13.stdout.json`, `ch13.stderr` (`exit status 1`), `ch13.exit`, `ch13.started`, `ch13.finished`. The sole failed check is `tts-boundaries`: **`want one complete utterance "Second paragraph", got0`**. Harness launch, discrimination, markup filtering, fragment buffering, provider-nonstream delivery, typing gate, error speech, identifiers, and Chapter 12 parity all reported passed. This resolves the launch/map blockers for these grader cases, not overall chapter acceptance or native/hearing proof.

The process failure stopped the sequence before the success-check parser/Chapter 1 branch. Only the emitted public JSON was read afterward. **No Chapter 1 execution**, retry, further diagnosis, source/fixture edit, or native/listener proposal. Root notified of the blocked gate and asked to direct any later correction after retained-reviewer inspection. No grader source accessed.

Cleanup qualification: grader returned and reported successful harness-run checks; browser-child teardown was not independently inspected after the failure and is not claimed verified. No student manual cleanup/artifact deletion. Only authorized grader-owned local browser/controlled-synth/fake-HTTP sessions ran; no separate rehearsal, native/audio/MIC, real provider/key/model, discovery exercise, live/shared browser, children, install, commit or snapshot. Source and existing artifacts remain on HOLD.

## Round 4/10 — line-local paragraph preservation, compiling HOLD

Actual edits: `gui/speech_text.go` restricts only `speechMarker` and `speechRule` whitespace to `[ \t]`, with a short explanation that `\s` can consume paragraph newlines. No identifier/link/fence/cue logic, raw cursor handling, hub/speaker/browser ownership, R1 cancellation/admission, R2 EOF, or R3 launch/diagnostics changed. New `gui/speech_text_paragraph_test.go` uses fresh prose and the real `cleanSpeech` / `speechText.update` helpers; no Hub scaffolding, private fixture, or grader source. This note is the only other touched path.

Valid before evidence: tests compiled and ran, exit **1**, demonstrating merged paragraphs before headings, quote/bullet and numbered markers, and under coalesced updates. Rule, soft-wrap, and paragraph-at-a-time controls passed before; the latter demonstrates why delivery grouping could mask the defect. The supplied reviewer's non-compiling probe is **not** credited as negative evidence. After the two-pattern fix, all new cases and existing selected formatter tests passed, including exact expected phrase order/content and empty repeated final flushes. No setup failures.

Evidence: `/tmp/ensemble-ch13-validation.At5CrJ/r4/{formatter-before.log,formatter-after.log,build.log,test-compile.log,exits.log}`. Commands from the GUI module:

```
go test -vet=off -count=1 -v -run '^TestSpeech(MarkupParagraphBoundaries|ParagraphGrouping)$' .
# before exit 1; actual test-body failures

go test -vet=off -count=1 -v -run '^TestSpeech(TextFiltering|StreamMatchesFinalAndFlushesOnce|StableSentenceAndIncompleteMarkup|SemanticSources|MarkupParagraphBoundaries|ParagraphGrouping)$' .
# after exit 0

go build -o /tmp/ensemble-ch13-validation.At5CrJ/r4/gui ./cmd
go test -vet=off -c -o /tmp/ensemble-ch13-validation.At5CrJ/r4/gui.test .
# both exit 0; final focused/compile aggregate 0
```

Same clean offline environment/shared caches and owned HOME/TMP; floor checked before commands, about 22 GiB free. Assistance was the public R4 correction brief only. No full race/vet, grade, Chrome/native/provider/key activity, other-module replay, frozen-CLI rebuild, child agent, install, cleanup, commit/snapshot, or guidance/manuscript/grader edit. All prior receipts remain untouched and retain their original outcomes. Source held after compilation for root + the same reviewer; no further execution or correction is underway.

## R4 changed-input validation — released gates passed, source HOLD

Fresh artifacts: `/tmp/ensemble-ch13-validation.At5CrJ/r4-validation`. No source/fixture edits during validation. Same clean offline/shared-cache setup with owned HOME/TMP, no real keys/inherited `LLM_BASE_URL`, and >1 GiB floor checked before each gate (about 22 GiB available). **Fresh** GUI `go test -race ./...` and `go vet ./...` exited **0**; GUI package results were 5.958 s / cmd 8.788 s. **Carried, not replayed:** unchanged root and ch13-speech race/vet, ch06 checks, ch07/ch04 compile-only test commands and vet, and the two controlled-JS checks. Earlier no-test qualifications remain.

The same non-rebuilt `/tmp/ensemble-ch13-validation.At5CrJ/ch01-frozen` was verified before grading: SHA-256 `1bb274f4bc9037b128c2886ee73732b0bb09c3667f067ddef9fbd83da9ebe34b`. Verification exit **0**. Reused the inspected, unchanged `r3-validation/media-models-v2.json`: anthropic `claude-sonnet-4-6`, openai `gpt-5.4`, gemini `gemini-3.8-flash`.

From book root `/Users/bill/projects/ensemble`, with `$root=/Users/bill/projects/ensemble/solutions/edition-3/main`, exactly one Chapter 13 run:

```
go run ./cmd/grade-edition3 -ch 13 -json \
  -core "$root" -consumer "$root/examples/ch04-consumer" \
  -workflow "$root/ch06" -stream-exercise "$root/ch07" -gui "$root/gui" \
  -stream-model claude-sonnet-4-6 -gui-model claude-sonnet-4-6 \
  -media-models /tmp/ensemble-ch13-validation.At5CrJ/r3-validation/media-models-v2.json \
  /tmp/ensemble-ch13-validation.At5CrJ/ch01-frozen
```

Chapter 13 process/shell **0**, parsed JSON **`passed:true`**, **100/100**; JSON pass-check exit **0**. All public checks, including boundaries and Chapter 12 parity, passed. Only then ran the **first Chapter 1 for this candidate**, once:

```
go run ./cmd/grade-edition3 -ch 1 -json /tmp/ensemble-ch13-validation.At5CrJ/ch01-frozen
```

Chapter 1 process/shell **0**, parsed JSON **`passed:true`**, **100/100**; pass-check exit **0**. Report confirms separate **1m2s header** and **1m2s body** delays, complete answers and usage, five rounds/five requests, preserved history, and correct cumulative usage **511 input / 65 output**. Started `2026-10-21T05:22:34Z`, finished `05:24:39Z` (125 s). No outer deadline or observation cancellation; no shortened wait/usage settings. These checks prove the reported waits, not arbitrary longer deadlines.

Receipts: `gui-race.log`, `gui-vet.log`, `frozen-hash.log`, `exits.log`; for each of `ch13` and `ch1`: `.stdout.json`, `.stderr`, `.exit`, `.started`, `.finished`, `.report-check.log`. Both stderr files were empty. Earlier failed 0/100 and 95/100 grades, wrong/corrected maps, and correction evidence remain untouched and retain their original outcomes.

Qualifications: Chapter 13 uses grader-owned isolated Chrome, controlled synth and local fake HTTP. This is speech-channel evidence, **not native hearing or real-model listener proof**. Its provider-nonstream check does not by itself establish separate final-only renderer coverage. Both grader processes returned successfully; no independent post-run enumeration verified browser-child teardown, and no manual cleanup was performed. No separate browser rehearsal, native/audio/MIC, real provider/model/key/discovery exercise, child agent, install, commit/snapshot, or live/shared browser restart. Source/frozen CLI stay held.

### Remaining native + blind-listener work — proposal only

1. Obtain a separate bounded release naming target and distinct listener models/vendors/base URLs, approved credential injection, and time/token/spend limits. Use an owned private GUI/browser, not live 8088 or a shared browser. Reuse existing attach mode; native run must have **no injected synth fake**. No microphone or OS focus/volume manipulation is proposed.
2. Have Bill confirm actual native audibility during short real-model GUI interactions. Exercise streaming and independently observed final-only delivery, semantic cues/paragraphs, Interrupt/Speak and full readable tools, and independent speech/typing holds through existing controls. Record native callbacks and Bill's observation separately from the page transcript; stop on an operational failure rather than silently substituting quiet mode.
3. Run the separate blind listener through exactly **hear/type/submit/wait/sleep**, receiving only admitted page speech/status and operating actual GUI input. No target Engine, provider-response shortcut, DOM/debug capability, or visual answer feed. Preserve listener responses and observed target outcomes. Transcript-based `hear` is speech-channel comprehension evidence, not acoustic recognition.
4. Stop at the approved cap or failure; no autonomous retry/correction campaign. Preserve concise results and verify shutdown of only the owned peer/browser/GUI processes under that release. Review evidence before any further work.

No scripts were created and none of these proposed native/listener actions were executed or self-authorized.

## Operational controller preparation — UNEXECUTED HOLD for root correction/review

Owned directory: `/private/tmp/ensemble-ch13-ops.Gk2JOo` (`/tmp/ensemble-ch13-ops.Gk2JOo` is its alias). Complete executable controller: `operate.mjs`; one existing attach-driver path, no alternative client framework or source copy. Root has requested HOLD for a consolidated controller-only correction; **this candidate is not execution-approved**.

Final SHA-256 pins:
- controller `operate.mjs`: `c795f546316a730aa1179c0f38ec99d0c0d1826f3c5b471b8b332e513d463dd8`
- `pins.json`: `1f378ccab16bc46776f0d7b879f350850282be8ee6fef54d7676fa08a68b61b9`
- `bin/gui`: `fa4fee74c3474129eac5ea64a9a80b6aa65156fbf9390385934dfefd5b9960e6`
- `bin/ch13-speech`: `d1f5d4df82542711cc9853fe51021b95b6796f222e581c84aa6c79e184a22c2f`

Fresh preparation only: built held GUI `go build -o "$D/bin/gui" ./cmd` from `main/gui` and driver `go build -o "$D/bin/ch13-speech" .` from its module, **once each, both exit 0**. Native source cap/route inspected: Anthropic `/v1/messages`, `max_tokens:1024`, version `2023-06-01`. Existing supported GUI flags/settings and native Chrome flags were read, not guessed. Source/assets/runtime inputs were hashed before builds, checked unchanged afterward; final pins include both binaries and controller (72 files, no copies/caches duplicated). Final `node --check` and controller `check` both exit **0**. Logs: `gui-build.log`, `ch13-speech-build.log`, `syntax.log`, `pure-check.log`, `prepare-exits.log`. Clean offline environment/shared caches, owned HOME/TMP, floor >1 GiB (about 20 GiB free). Preparation-only fixes aligned the hear envelope with the actual driver, canonicalized `/tmp`, and retained fixed Go limits in child environments; no production/test changes or gate reruns.

Prepared scope, **not observed execution**: A uses the actual native page without synth injection, one planned fake POST and three short utterances, both hold orders/Speak/Escape observation, upper bounds four utterances/90 seconds speech/300 seconds session. B is separately gated on a reviewed same-pin local receipt: quiet actual GUI plus distinct listener Agent/context using only the five tools, one model task, atomic 4 target/8 listener/12 total pre-upstream POST bounds, unchanged 1024 native caps, one discovery GET/key-read attempt, 240 active/300 absolute seconds, no retries/redirects/body rewriting or model substitution. Real credentials are not supplied to children. The controller records actual exits and partial failures; cleanup claims are limited to tracked groups/listeners. Runtime correctness of these paths remains unvalidated.

Coverage qualification: complete Cards passed directly to `Hub.narrate` establish carried **owner-level full-part admission**. The cited old browser late-delta block has preceding deltas and is **not** direct-final browser proof. Live frame ordering is only to be observed/reported, with no interception or manufactured native/final-only result. No hearing, voice-locality or acoustic-comprehension success is claimed. Listener content omissions remain distinct from operational completion; no exact fruit/answer assertion.

Exact held commands (canonical `$D=/private/tmp/ensemble-ch13-ops.Gk2JOo`, `$P=1f378ccab16bc46776f0d7b879f350850282be8ee6fef54d7676fa08a68b61b9`):

```
# Preparation check was executed; pure/local only:
env -i HOME="$D/home" TMPDIR="$D/tmp" PATH=/usr/bin:/bin:/usr/sbin "$D/operate.mjs" check "$P"
# NOT EXECUTED; separate root release and visible native notice required:
env -i HOME="$D/home" TMPDIR="$D/tmp" PATH=/usr/bin:/bin:/usr/sbin "$D/operate.mjs" native "$P"
# NOT EXECUTED; only after reviewed local success and separate paid release:
env -i HOME="$D/home" TMPDIR="$D/tmp" PATH=/usr/bin:/bin:/usr/sbin "$D/operate.mjs" listener "$P" NATIVE_RECEIPT_SHA256
```

**Zero operational execution:** no GUI/browser/native/fake-provider/listener session, external request, discovery, credential read, MIC, focus/volume operation, live/shared browser, cleanup, source/test/fixture/grader/manuscript/guide edit, child agent, snapshot or commit. Prior gates and failures are carried with their original qualifications, not repeated or overwritten. Source and the original frozen CLI remain held and unchanged. Only this note and the owned preparation directory changed. Await consolidated root correction/review; no more checks or revisions underway.

## Operational controller v2 — bounded correction, UNEXECUTED HOLD

New files in the same owned directory `/private/tmp/ensemble-ch13-ops.Gk2JOo`:
- `operate-v2.mjs` SHA-256 `03fffc8d78ec253a8b9fa355886bb08656a9c9d40df1618fd02f24e19ae17f6b`
- `pins-v2.json` SHA-256 `2ff080757e607c320d03d24401e3962203df38c86d9d0af3d4a72683128e57dc`
- `operate-v2.diff`: exact unified diff from preserved `operate.mjs`.

Only the requested controller delta: provider guard/transport/discovery/listener-validation/non-200 failures invoke the existing global abort owner, preserve the first failure’s status/count/error metadata, deny subsequent reservations/upstream starts, and abort outstanding requests. Non-200 responses are no longer forwarded as if the attempt could continue. Removed HTTPS `timeout:20000` and its timeout handler; the existing 240s paid timer is armed before discovery, with the 300s absolute cap retained. Added `--mute-audio` only to quiet-listener Chrome. Added a 60s setup timer, cleared after owned GUI/CDP/page readiness and before native/discovery work. Busy predicate, native callbacks, no-synth-replacement policy, and precise final-only/hearing/locality limitations are unchanged. No new framework, operational mode or review campaign.

Fresh checks only, both exit **0**: canonical Node `--check operate-v2.mjs`; `operate-v2.mjs check` with the new pins hash. The pure checks include refusal to advance a POST budget after abort. Logs: `syntax-v2.log`, `pure-check-v2.log`, `prepare-v2-exits.log`. Clean environment; 1 GiB floor checked (~20 GiB free). These are syntax/pure checks, **not live failure-path, timeout, browser, native or model observations**.

Exact held command pattern (`D=/private/tmp/ensemble-ch13-ops.Gk2JOo`, `P=2ff080757e607c320d03d24401e3962203df38c86d9d0af3d4a72683128e57dc`):
```
# Executed: pure/local preparation only.
env -i HOME="$D/home" TMPDIR="$D/tmp" PATH=/usr/bin:/bin:/usr/sbin "$D/operate-v2.mjs" check "$P"
# NOT EXECUTED; each operational phase still needs its separate root release.
env -i HOME="$D/home" TMPDIR="$D/tmp" PATH=/usr/bin:/bin:/usr/sbin "$D/operate-v2.mjs" native "$P"
env -i HOME="$D/home" TMPDIR="$D/tmp" PATH=/usr/bin:/bin:/usr/sbin "$D/operate-v2.mjs" listener "$P" NATIVE_RECEIPT_SHA256
```

Old script/pins, prior note content, production source, binaries, frozen CLI and all prior receipts preserved. No rebuild, gate replay, browser/native/provider/listener session, network/discovery/key access, cleanup or delegated child. Await root inspection; overall acceptance remains pending.

## One released native-local attempt — FAILED, paid authority CLOSED

Executed **once** with the reviewed v2 controller/pins, live PTY stdin and clean environment:
```
D=/private/tmp/ensemble-ch13-ops.Gk2JOo
P=2ff080757e607c320d03d24401e3962203df38c86d9d0af3d4a72683128e57dc
env -i HOME="$D/home" TMPDIR="$D/tmp" PATH=/usr/bin:/bin:/usr/sbin "$D/operate-v2.mjs" native "$P"
```
Visible native-audio notice was given; floor check passed (21,576,776 KiB available). **Process exit 1; receipt status `failed`.** Started `2026-10-21T06:35:46.423Z`, setup ready `06:35:47.473Z`, finished `06:36:07.543Z` (about 1.05s setup / 21.12s total). Primary recorded error: `bounded observation expired`. A subsequent provider-failure entry says `operational failure`; this is not established as an independent provider defect. No diagnosis campaign undertaken.

Actual receipt counts: **one local fake target POST**, model `claude-sonnet-4-6`, `max_tokens:1024`, no tools; **zero listener POSTs, zero credential reads and zero discovery GETs**. No paid/listener phase was attempted. Configured limits remained 4 utterances / 90s speech / 300s session / 60s setup. The failed receipt does not contain a successful native-observation block or measured native utterance count/duration; do not replace that absence with a claim of zero speech, successful callbacks, human hearing, acoustic comprehension or voice locality.

Owned cleanup as recorded: driver PID 97611 exit **1**, Chrome PID 97583 exit **0**, GUI PID 97579 exit **0**; all three tracked groups report `group_gone:true`; proxy port 51178 closed; `cleanup_complete:true` within the controller’s stated tracked-group/listener scope. No broader OS cleanup or untracked-process claim. Profiles/raw logs retained; no manual/shared cleanup performed.

Artifacts preserved under `$D`: `native-attempt.stdout`, `native-attempt.stderr` (empty), `native-attempt.exit` (`1`), and `native/receipt.json` plus the existing native run directory. Receipt SHA-256: `2f1d060178e7fd884708db21b3c5f20961d1ca20e08656ae0400535a8579a8b0`.

Stopped at failure. No retry/fallback, extra browser attempt, source/controller edit, rebuild/gate, MIC/OS focus/volume operation, key/discovery/paid/listener run, delegated child or commit. Prior source/local-grade results and failures remain unchanged. **Return HOLD; all paid authority remains CLOSED; overall acceptance remains pending.**

## Controller v3 first-native-wait diagnostic — preparation ONLY, UNEXECUTED HOLD

Public root/reviewer classification: the first wait failed; saved frames establish prompt admission/claim/done, but absent page snapshots leave the cause unknown. **No product finding and no production Round 5.** Old native failure remains failed and untouched.

New versioned artifacts in `/private/tmp/ensemble-ch13-ops.Gk2JOo`:
- `operate-v3.mjs` SHA-256 `0f7a7315394403349f4e6b27b178bd997ac8cf93541624119bafbe72df63c24f`
- `pins-v3.json` SHA-256 `e614af7507d5a056f66fe7091537edf4bab5b138dd3eb9c4faa20ceac057ac6e`
- `operate-v3.diff`: exact narrow diff from preserved v2.

Only first-wait diagnostics and version paths changed. Actual `snapshot(page)` results already sampled by that wait are saved in full (`tts`, `channel`, `status`) to `native-v3/native-first-wait.jsonl` on value change. On expiry/error, the latest actual sample is synchronously persisted before the outer abort/cleanup, with its observation timestamp, save timestamp and error. It is explicitly `timeout-last-sample`, **not a fabricated fresh observation**; no extra CDP read or added wait. If no sample was obtained, null remains explicit. Original predicate/20s wait, flags, speech text/rate, synth behavior, caps/deadlines, transport and cleanup are unchanged.

Future native and listener run paths are explicitly `native-v3` and `listener-v3`. Listener prerequisite reads only `native-v3/receipt.json` and still requires successful status, matching v3 pins, receipt hash and cleanup. The preserved `native/receipt.json` failure cannot satisfy it.

Fresh syntax and pure pin/guard checks only: both **exit 0**. Logs `syntax-v3.log`, `pure-check-v3.log`, `prepare-v3-exits.log`; clean environment and >1GiB floor (~20GiB free). Original v2 and old failed receipt hashes were checked unchanged. No runtime snapshot capture is claimed by these checks.

Exact held commands (`D=/private/tmp/ensemble-ch13-ops.Gk2JOo`, `P=e614af7507d5a056f66fe7091537edf4bab5b138dd3eb9c4faa20ceac057ac6e`):
```
# Executed: pure/local only.
env -i HOME="$D/home" TMPDIR="$D/tmp" PATH=/usr/bin:/bin:/usr/sbin "$D/operate-v3.mjs" check "$P"
# NOT EXECUTED; awaits a separate root diagnostic-run release.
env -i HOME="$D/home" TMPDIR="$D/tmp" PATH=/usr/bin:/bin:/usr/sbin "$D/operate-v3.mjs" native "$P"
```

No new binaries, source edits, gates, browser/model/native run, voice discovery, GET/key/paid access, delegated child or cleanup. All old scripts/pins/results/note content preserved. **HOLD for root inspection; paid/listener authority remains closed.**

## One released native-v3 diagnostic invocation — FAILED, observations recovered

Ran the exact reviewed v3 native command once, with live PTY stdin, clean environment and visible native-audio notice. Initial free space **21,040,552 KiB** (>1GiB). **Process exit 1 / receipt `failed`**, primary error `bounded observation expired`. Started `2026-10-21T07:04:30.965Z`, setup ready `07:04:32.202Z`, finished `07:04:52.302Z` (~1.237s setup / 21.337s total). No retry, variant or further live inspection.

Actual saved first-wait observations: `native-v3/native-first-wait.jsonl` contains 705 sampled rows and one `timeout-last-sample` record (697,296 bytes). All sampled client IDs were `2`; current speech was always null, pending/delayed/requests/recent arrays always empty, and native `speaking/pending/paused` all false. Final actual sample at `07:04:52.228Z`, saved at `07:04:52.229Z` before abort: agent-1 ready, turn1 `in_flight`, outcome `pending`, busy/paused/human_typing false, draft empty; native channel cursor2 contains only typing pause and clear/resume. No successful native callback observation is established for this sampled page/client. This is **not** a claim of system-wide zero speech, human hearing, voice locality or an underlying cause/product defect. The later planned native pause/Speak/Escape sequence was not reached.

Diagnostic limitation: `tts.observed_at` changes every poll, so whole-snapshot comparison retained otherwise identical rows rather than meaningfully compacting them. Full data and the final sample were preserved; no controller fix or re-execution was made. Inspection was limited to saved receipt/process output and a read-only summary of this saved JSONL after the inline read clipped. No diagnosis campaign.

Actual counters: **one local fake target POST**, `claude-sonnet-4-6`, native `max_tokens:1024`, tools[]; listener POSTs0, proxy credential reads0, model-discovery GETs0. No paid/listener phase. The receipt’s secondary `operational failure` provider entry accompanies abort; no independent provider-failure classification is made. Original caps unchanged: <=4 short utterances/90s speech/300s session/60s setup/fake cap4; this failed observation does not demonstrate the complete native exercise.

Recorded owned cleanup: driver PID97927 exit1, Chrome PID97908 exit0, GUI PID97904 exit0; all tracked groups gone, proxy51251 closed, `cleanup_complete:true` within the stated tracked-group/listener scope only. Profiles and traces retained; no manual/shared or broader OS cleanup claim.

Artifacts under `/private/tmp/ensemble-ch13-ops.Gk2JOo`: `native-v3-attempt.stdout`, empty `native-v3-attempt.stderr`, `native-v3-attempt.exit`, and `native-v3/receipt.json` plus the full trace/run directory. Receipt SHA-256 `4e226a248d69d73f6ab23e6bc6a6719b5a254ed20807c2f823cf5cc23df7c024`. Old native failure, controllers/pins, source/binaries, gates and prior receipts unchanged.

**HOLD.** No source/controller/binary changes, extra attempts, gates, MIC/OS focus/volume operation, real provider/key/discovery/paid/listener work, delegated child or commit. Paid/listener authority remains closed; overall acceptance remains pending.

## Production correction Round 5/10 — submit/release ordering, compiling HOLD

Public classification applies to native-v3: its own automatic typing release invalidated newly requested prompt speech. **Do not retroactively assign that cause to native-v2; its done-trace cause remains unknown.** No native rerun was performed for this correction.

Production delta is only `gui/web/gui/app.js`’s existing `submit` helper: retain disconnected/empty drafts; clear this draft’s HUMAN cause and publish the effective automatic release **before** publishing the message. An independent speech hold still prevents automatic release. Agent/Hub epoch semantics, explicit Continue, Interrupt, new-turn/stale-claim/callback fences and speech backends are unchanged. Hint shares this helper/message-part projection; no separate state mechanism was added.

Focused actual-owner proof added in `gui/ch13_submit_test.go` with `gui/ch13_submit_bridge.cjs`: production input/form/WS-receive handlers in the existing controlled VM fixture generate commands sent over a real GUI WebSocket to a real Agent/Hub. A local HTTP handler holds model output. Each submit command takes effect before the next; claims wait until submission commands have completed. Only actual server frames are returned to JS—no injected preclaimed speech frame. A real duplicate claim plus a pause acknowledgement checks one prompt admission. This is controlled JS/local-owner evidence, **not Chrome/native/hearing evidence**.

Before proof: the valid executed regression failed with `input -> unpause`, request epoch1, claim epoch2, `speech request expired`, and prompt channel admission count0. After: `unpause -> input`, request/claim epoch2, `speech_claimed`, exactly one prompt admission. Two earlier attempts were **test-adapter setup failures**, not product negative evidence: the bridge mistakenly called the fixture’s convenience replay wrapper with a wire envelope. The bridge now uses production `socket.onmessage` and forwards original WS JSON unchanged. The initial speculation about lost Wire fields was not the cause.

Added narrow cases in `gui/ch13_correction_test.cjs` for disconnected drafts, whitespace-only HUMAN hold, and retaining an active speech cause during submit; existing segment/callback/cancellation/source checks also passed. Actual touched source/test paths: the four files above. No core/server/speaker/controller change or framework.

Commands/results, clean offline env/shared caches, GOMAXPROCS2/readonly-p1/local toolchain, owned HOME/TMP under `/private/tmp/ensemble-ch13-ops.Gk2JOo/r5`, >1GiB floor:
- `go test -vet=off -count=1 -v -run '^TestCh13SubmitPromptSurvivesOwnRelease$' .`: valid **before exit1** (`submit-before-receive.log`); setup attempts remain separately in `submit-before.log` and `submit-before-raw.log`.
- `go test -vet=off -count=1 -v -run '^Test(Ch13SubmitPromptSurvivesOwnRelease|Ch7R1SpeechCancelClaimAndOtherTabRelease)$' .`: **after exit0**, both tests ran (`submit-after.log`).
- `node ch13_correction_test.cjs`: **exit0**, 13 cases (`controlled-js.log`).
- GUI module `go build -o "$work/gui-r5-compile" ./cmd` and `go test -vet=off -c -o "$work/gui-r5.test" .`: **exit0** each.
- `node --check web/gui/app.js`, `ch13_submit_bridge.cjs`, and `ch13_correction_test.cjs`: **exit0** each. Per-command logs and `exits.log` retained in `$work`.

No full race/vet/suite/grade, Chrome/native, real provider/key/discovery, other-module replay, frozen CLI rebuild, delegated child, snapshot/commit or manual/shared cleanup. Old R4 binaries/controllers/pins/native failures were not altered. Existing operational pins still describe the old held inputs and were not repinned for this changed page asset. **HOLD for root and same reviewer before changed-input validation; paid/listener authority remains closed.**

### R5 HOLD confirmation — latest steering

The Round 5 delta and qualified receipts immediately above are the final held record. No further tests, packaging, broader GUI race/vet/grade, controller preparation, native or paid/listener operation was started after the latest HOLD instruction. Source remains held for root and the same reviewer; all broader execution is closed pending renewed explicit release.

## R5 release 1629 — fresh validation PASS; v4 preparation-only HOLD

Executed after root clarified that release1629 superseded the earlier HOLD message. No production/source edits during validation.

**Fresh gates**, sequential clean offline env/shared caches, GOMAXPROCS2/readonly-p1/local toolchain, owned HOME/TMP, >1GiB floor:
- GUI `go test -race ./...` **exit0** (gui and gui/cmd); `go vet ./...` **exit0**.
- Recorded JS suites: `node ch08_round1_test.cjs` **exit0** (9 cases); `node ch13_correction_test.cjs` **exit0** (13 cases). No separate R5 focused-Go replay.
- Frozen CLI hash verified unchanged: `/tmp/ensemble-ch13-validation.At5CrJ/ch01-frozen`, SHA256 `1bb274f4bc9037b128c2886ee73732b0bb09c3667f067ddef9fbd83da9ebe34b`.
- **ONE Chapter13 grade: process0, parsed passed:true,100/100; all reported checks pass.** Book-root command was the same explicit `go run ./cmd/grade-edition3 -ch 13 -json` with absolute core `main`, consumer `main/examples/ch04-consumer`, workflow `main/ch06`, stream-exercise `main/ch07`, gui `main/gui` under `/Users/bill/projects/ensemble/solutions/edition-3/`; stream/gui models `claude-sonnet-4-6`; media map `/tmp/ensemble-ch13-validation.At5CrJ/r3-validation/media-models-v2.json`; all flags before the unchanged frozen CLI. Map: anthropic `claude-sonnet-4-6`, openai `gpt-5.4`, gemini `gemini-3.8-flash`.

Logs: `/private/tmp/ensemble-ch13-ops.Gk2JOo/r5-validation/{local-exits.log,gui-race.log,gui-vet.log,js-ch08.log,js-ch13.log,frozen-hash.log,ch13.stdout.json,ch13.stderr,ch13.exit,ch13.report-check.log,ch13.started,ch13.finished}`. Grader returned normally; no independent post-grade OS orphan scan or manual/shared cleanup. Grader-owned controlled synth/fake HTTP/browser evidence is not native audibility/hearing evidence; provider-nonstream is not standalone final-only renderer proof.

**Carried, not rerun:** unchanged root/driver/application gates (including earlier no-tests=compile-only qualifications), and the FIRST genuine R4 Chapter1 result on this identical CLI: process0/parsed passed:true/100, separate62s header/body waits and full five-round usage. No Chapter1 repeat and no frozen CLI rebuild.

Conditional preparation completed ONLY after all gates passed:
- ONE GUI `go build -o "$work/gui-r5" ./cmd` from GUI module, **exit0**. New path `/private/tmp/ensemble-ch13-ops.Gk2JOo/r5-validation/gui-r5`; SHA256 `fa4fee74c3474129eac5ea64a9a80b6aa65156fbf9390385934dfefd5b9960e6` (byte-identical to retained GUI binary; R5 production change is a served asset). Old GUI and unchanged driver retained; driver SHA256 `d1f5d4df82542711cc9853fe51021b95b6796f222e581c84aa6c79e184a22c2f`.
- `operate-v4.mjs` SHA256 `dfe653dec28151b38dd7d61939d765a907422ad2b98eb714b33a7b10c2617662`.
- `pins-v4.json` SHA256 `172b7cb7178658e32a149357135cdc2e74d83172eb7f71592ea817796213ad90`;76 pinned inputs. Only expected held-input drift was R5 `app.js` (new SHA256 `27edd705210ea38df68053c6ee2e92386655429bf804f5c7b16bb4079350657c`).
- Exact diff: `operate-v4.diff` in the same owned directory. Only version/pin references, new GUI path, fresh `native-v4`/`listener-v4` paths and authorized comparison-only omission of `tts.observed_at` changed. Full original samples and final last-sample record remain preserved; predicate/behavior/timing/flags/prompt/rate/budget/abort/cleanup unchanged. Listener prerequisite points only to a successful matching `native-v4` receipt.
- Syntax/pure checks **exit0**: `syntax-v4.log`, `pure-check-v4.log`, `preparation-exits.log` in `$work`. No operational phase ran.

Exact held native command (NOT EXECUTED / not authorization):
```
D=/private/tmp/ensemble-ch13-ops.Gk2JOo
P=172b7cb7178658e32a149357135cdc2e74d83172eb7f71592ea817796213ad90
env -i HOME="$D/home" TMPDIR="$D/tmp" PATH=/usr/bin:/bin:/usr/sbin "$D/operate-v4.mjs" native "$P"
```

All old maps/scripts/pins/binaries, failed grades and native-v2/v3 receipts preserved. No independent browser/native/real-provider/key/model-discovery/paid/listener activity, delegated child, source archive, guide/manuscript/grader edit, commit or shared cleanup. **HOLD for root inspection of gates and controller diff; native/paid execution remains CLOSED; no overall acceptance claimed.**

## ONE released native-v4 operation — FAILED; callback observations, cleanup qualification; HOLD

Executed exactly once with live PTY stdin and visible audio notice:
`env -i HOME="$D/home" TMPDIR="$D/tmp" PATH=/usr/bin:/bin:/usr/sbin "$D/operate-v4.mjs" native 172b7cb7178658e32a149357135cdc2e74d83172eb7f71592ea817796213ad90`, `D=/private/tmp/ensemble-ch13-ops.Gk2JOo`. Start disk20,570,380KiB free (>1GiB). **Process1 / receipt failed**; no retry, variant or use of unused slots.

Receipt `native-v4/receipt.json`, SHA256 `e4bc5793f2dd55df60d9d2bf6f887acdb093a8c3c594b01fb2e050f0c8e290f9`; outer `native-v4-attempt.stdout/.stderr/.exit` retained (stderr empty). Started08:55:48.470Z; ready08:55:49.627Z (~1.157s); finished08:56:09.717Z (~21.247s), 2026-10-21. Primary error: `bounded observation expired`; provider failure entry `operational failure` recorded during abort.

Actual saved state (`native-v4/native-first-wait.jsonl`,5 changes + final timeout-last-sample): one human-prompt utterance admitted, speech_id1/part_id1; native speaking true while current; actual page `native-start` at1792572949704 and `chunk-ended`/`done` at1792572951959 (~2.255s callback span). Channel records typing pause/clear, speech pause, one utterance, completion done and clear/resume. Current/recent text members are bounded-text objects. Controller nevertheless remained in its first wait and timed out; later planned Speak/overlapping-pause/Escape steps were not reached. No fix or diagnosis campaign performed. These are callback/channel observations—not human hearing, acoustic comprehension, voice locality, or successful overall native operation.

Counts: **1 local fake target POST**, model claude-sonnet-4-6, unchanged max_tokens1024, tools[]; listener0, credential reads0, model-discovery GETs0; no real-provider/paid activity. Observed one native start/end sequence; no further utterance or operation authorized. Original4-utterance/90s speech/60s setup/300s session/fake cap4 bounds unchanged.

**Cleanup is not confirmed complete:** receipt says cleanup_complete:false. Driver99613 exit1/group_gone true; GUI99590 exit0/group_gone true; Chrome99594 parent exit0 but group_gone:false at receipt time. Proxy closed. No claim that a live orphan persists or that broader OS teardown succeeded; no extra manual probe/kill/shared cleanup performed. Root notified via question channel for any narrowly scoped follow-up authority.

All earlier failed native runs remain failed. No source/controller/binary/gate/key/discovery/listener changes or further execution; paid/listener authority remains CLOSED. HOLD for root review/direction.

### Native-v4 cleanup follow-up — root-authorized read-only verification

At2026-10-21T09:10:21Z, scoped `ps -g 99594` and `ps -p 99594` returned no members (exit1, no diagnostics); exact private-profile `pgrep` returned no match (exit1); signal-0 process-group existence check reported **absent**. No termination signals, browser attachment, relaunch, shared cleanup or other operational activity occurred. Read-only wrapper exit0; `/private/tmp/ensemble-ch13-ops.Gk2JOo/native-v4-cleanup-verification.log` retains subcommand exits.

This is a later scoped observation, **not** a rewrite of the original failed receipt/cleanup_complete:false and not a broader OS cleanup claim. Native-v4 remains FAILED; paid/listener authority remains CLOSED. HOLD for root’s controller/text-shape review; no further work undertaken.

## Controller v5 — bounded preparation complete, UNEXECUTED HOLD

Resumed the same unfinished preparation after inference interruption; no v5 file had been written before that interruption. It was not a product/test failure or new production round. Held R5 source, passed gates, binaries, old controllers/pins and failed native receipts remain unchanged.

Minimal v4→v5 delta (`/private/tmp/ensemble-ch13-ops.Gk2JOo/operate-v5.diff`):
- `queueTextIs` requires the actual bounded-text object’s complete exact expected text, `truncated:false`, and matching total_characters. Applied only to queue current/recent comparisons; channel event text stays string-based. Native start/done, cancellation, both independent-pause checks, words/rate/caps/flags/deadlines are unchanged.
- Pure positive checks read the retained actual v4 speaking snapshot/current/native-start record. Same-length wrong text, truncated, inconsistent length and plain-string negatives reject. Snapshot itself is pinned (SHA256 `fb35644c8f11b376a4eddcd5e40bb171079d25b89aa1262f06c43b1d00463928`). No fabricated callback or new native run.
- Cleanup polls only the recorded owned group’s absence within the existing1500ms TERM /1000ms KILL graces, not an already-resolved parent promise. Initial stdin/parent grace and existing signal ownership are retained; no later signal after observed disappearance. No OS cleanup framework or new scope. Runtime cleanup remains untested in this preparation.
- Fresh `native-v5`/`listener-v5`; listener prerequisite only a matching successful/clean native-v5 receipt. All older FAILs remain FAILs, including v4’s original cleanup_complete:false. Later scoped absence observation is not retroactive success.

Files/hashes:
- `operate-v5.mjs`: `011ee7621553342c1818e71ad150454ab30bc99804cd33cc33d93acbd6091b64`.
- `pins-v5.json`: `b533f30185bb84d90b71b554f9aa8e321f3472c71de749d91dbadb90d50c7ddc` (78 inputs).
- Syntax `node --check operate-v5.mjs` and existing controller `check` mode with that pins hash: **exit0** each. `syntax-v5.log`, `pure-check-v5.log`, `preparation-v5-exits.log` retained. Clean environment and >1GiB floor; no application gate/rebuild/review campaign.

Exact future held command (NOT EXECUTED / not authorization):
```
D=/private/tmp/ensemble-ch13-ops.Gk2JOo
P=b533f30185bb84d90b71b554f9aa8e321f3472c71de749d91dbadb90d50c7ddc
env -i HOME="$D/home" TMPDIR="$D/tmp" PATH=/usr/bin:/bin:/usr/sbin "$D/operate-v5.mjs" native "$P"
```

No native/browser/provider/key/model-discovery/paid/listener/cleanup execution, children, source/binary/production edits, commit or shared cleanup. **HOLD for root inspection and separate operational release; paid listener CLOSED.**

## ONE released native-v5 operation — FAILED; partial native observation; HOLD

Executed the exact reviewed v5 command once with live PTY stdin and visible audio notice; start disk21,737,496KiB free (>1GiB). **Process1 / receipt failed**, primary `bounded observation expired`; no retry/variant or unused-slot replay.

`/private/tmp/ensemble-ch13-ops.Gk2JOo/native-v5/receipt.json` SHA256 `f79bbff276bf7fbb39a6022553fad5cc1c184a6e21a06026eacb874367547a35`; outer `native-v5-attempt.stdout/.stderr/.exit` retained (stderr empty). Started2026-10-21T11:57:08.723Z; ready09.907 (~1.184s); finished29.243 (~20.520s total). Provider failure entry `operational failure` recorded during abort.

Actual `native-v5/native-first-wait.jsonl` has four change snapshots and no first-wait error row: one prompt admitted (speech_id1/part_id1), speech pause, then native-start at1792583829989 with current state speaking/native.speaking true at the final sampled first-wait state. The schema-correct first wait advanced. A later bounded observation failed; these saved first-wait samples do not establish the exact later failed stage, subsequent native count/duration, completion/cancellation or full independent-pause exercise. No further diagnosis/probe or fix made. Callback/channel evidence is **not** human hearing, locality or acoustic-comprehension proof.

Counts:1 local fake target POST, claude-sonnet-4-6/max_tokens1024/tools[]; listener0, credential reads0, model-discovery GETs0; no real-provider/paid activity. Original4-short-utterance/90s speech/60s setup/300s session/fake cap4 bounds unchanged; unused capacity is not authority.

Receipt cleanup_complete:true for its stated tracked scope: driver1237 exit1/group_gone true; Chrome1210 exit0/group_gone true; GUI1206 exit0/group_gone true; proxy closed. Profiles/raw traces retained. No broader OS cleanup/hearing claim or manual/shared cleanup.

All old failures remain failed. No source/controller/binary/gate change, new build/test, children, commit, key/discovery/paid/listener execution or further native/browser attempt. **HOLD; root will inspect and decide remaining work.**

### Native-v5 classification addendum — root’s saved-wire/source review

Root reports first wait **succeeded**, and controller ordering therefore shows the speech-first overlap and post-speech HUMAN-hold wait completed before the second/on-demand wait. Actual saved frames show speech id2 reason exactly **`on demand`**, claimed then done; the controller compares nonexistent **`ondemand`**. This is a controller literal error, not a production finding. No fix or new run has been made.

Second-utterance native callback state is not retained in the first-wait trace: claim/done frames do not establish second native callback success, human hearing, locality or acoustic comprehension. Overall native-v5 remains FAILED/process1; first native-start evidence and the narrower completed hold observations above are retained without upgrading the overall result. Receipt’s tracked cleanup remains complete (driver1, GUI/Chrome0, all recorded groups gone, proxy closed), not broader OS cleanup. Await one narrow versioned preparation release; no further analysis, probes or execution. HOLD; paid/listener CLOSED.

## Controller v6 — minimal on-demand preparation, UNEXECUTED HOLD

Exact diff `/private/tmp/ensemble-ch13-ops.Gk2JOo/operate-v6.diff`: shared `ON_DEMAND_REASON='on demand'` replaces only the wrong second-wait literal. New pure assertion uses retained/pinned native-v5 `speech_claimed` id2, accepts the actual reason and rejects `ondemand`. Simple `receipt.stage` assignments immediately precede existing native waits; no new observation loop, predicate weakening, timeout, voice/word/rate/flag/cap/cleanup/secret-rule change. Fresh `native-v6`/`listener-v6`, with listener prerequisite only matching successful/clean native-v6.

Controller SHA256 `29bef76e71fc6f94b0ba095fa8ce2a71c1c9aa3929938e3f6ba8cd26349c49cd`; pins-v6 SHA256 `3ac3d1990b57f23f6b6996f86e4eb8c531607e3c3e468757d2fb6c704114a3ea` (80 inputs). Retained frame input hash `8de55329282e035fd28050bdb6d21fee14af24c580c62d925a846c55e77c41e0`.

Syntax exit0; focused added reason-check block exit0 (shared constant and actual pinned frames checked; other guard/application suites not replayed). `syntax-v6.log`, `pure-reason-v6.log`, `preparation-v6-exits.log` retained. First check wrapper had shell-quoting parse exit2 before any controller/check body; corrected invocation only, files unchanged. Clean environment and >1GiB floor. No browser/native/model/key/discovery/listener/paid/cleanup operation, rebuild, production edit or correction round.

Exact future HELD command, NOT executed/authorized:
```
D=/private/tmp/ensemble-ch13-ops.Gk2JOo
P=3ac3d1990b57f23f6b6996f86e4eb8c531607e3c3e468757d2fb6c704114a3ea
env -i HOME="$D/home" TMPDIR="$D/tmp" PATH=/usr/bin:/bin:/usr/sbin "$D/operate-v6.mjs" native "$P"
```
Old controllers/pins/results and failed native-v5 remain unchanged/failed; all passed production gates carry. No new binaries, source/archive, children, commit or shared cleanup. **HOLD for root inspection; operational release and paid listener CLOSED.**

## ONE released native-v6 operation — observed local PASS / process0; HOLD

Executed the exact held v6 native command once, live PTY stdin, visible audio notice. Start disk21,217,616KiB free (>1GiB). Receipt `/private/tmp/ensemble-ch13-ops.Gk2JOo/native-v6/receipt.json` SHA256 `ce2879fa5addd020137e4456b05bb1ea86c6b6808ae53131e2452e3951cadaf8`: **observed-pass**, process0. Outer `native-v6-attempt.stdout/.stderr/.exit` retained (stderr empty), plus `page-final.json`, first-wait trace and raw local records. No retry/variant.

2026-10-21 started12:15:13.495Z, ready14.678 (~1.183s setup), finished22.872 (~9.377s total). Actual3 admitted utterances /3 native-start records: initial prompt, intentional explicit Speak of that prompt, then short reply. Initial prompt native-start14761→chunk-ended16957 (~2.196s within timestamp prefix17925849); on-demand speech native-start16981→cancellation17010 (~29ms, not natural completion); reply native-start17041→chunk-ended21780 (~4.739s). These are browser callback/channel measurements, not physical audio-duration/hearing/locality/comprehension proof.

Controller completed its speech-first hold/post-speech HUMAN-hold checks; typing-first speech-held check after releasing only the owned whitespace draft; Escape cancellation and release; target success/driver completion. Final page idle/success, empty draft, human_typing:false, paused:false, busy:false; no current/pending/delayed/request speech, native speaking/pending/paused:false. Second prompt reading is intentional explicit Speak, not a duplicated automatic prompt.

Counts:1 local fake target POST, claude-sonnet-4-6/max_tokens1024/tools[]; listener0, credential reads0, model-discovery GETs0; no real-provider/paid call. Original limits4 short utterances/90s audio/60s setup/300s session/fakecap4 retained. Unused capacity confers no additional authority.

Tracked cleanup_complete:true: driver1550, Chrome1531, GUI1527 all exit0/group_gone:true; proxy closed; profiles/traces retained. No broader OS cleanup claim or manual/shared cleanup.

Qualification: assistant part2 had prior_delta:true. Owner-level complete-Card coverage carries, but this run does **not** establish direct-final browser/native rendering. Human hearing not observed; locality/acoustic comprehension not established. All older failed attempts remain failed. No source/controller/binary/gate change, build/test, children or commit. **HOLD for root inspection; paid/listener execution remains CLOSED despite local prerequisite success.**

### Native-v6 human testimony and scoped closure

Root verified the machine receipt/final states and closed the scoped native gate. Bill reports hearing the ending **“…past the stones”**. This is limited human testimony, separate from the unchanged machine receipt; it does not establish full-utterance hearing, voice locality or acoustic comprehension. Root separately released exactly one existing listener-v6 operation; no native rerun or blanket paid authority follows.

## ONE released listener-v6 operation — FAILED before model POSTs; HOLD

Executed the exact existing env-i listener invocation once, live PTY, pins `3ac3d1990b57f23f6b6996f86e4eb8c531607e3c3e468757d2fb6c704114a3ea` and native receipt `ce2879fa5addd020137e4456b05bb1ea86c6b6808ae53131e2452e3951cadaf8`. Start disk21,148,708KiB (>1GiB). **Process1 / receipt failed**. No repair, retry, fallback, unused-budget replay or additional model activity.

`/private/tmp/ensemble-ch13-ops.Gk2JOo/listener-v6/receipt.json` SHA256 `5fd21352dd55d5dfa1fe717fe43b1f163e45dcb53f76a9c377f407f28de80845`; outer `listener-v6-attempt.stdout/.stderr/.exit` retained (outer stderr empty). Started2026-10-21T12:20:42.315Z; ready43.539 (~1.224s); finished43.948 (~1.633s total). Receipt error `listener driver nonzero exit`; actual retained `listener-v6/driver.log`: **`attach requires a prompt argument`**.

Actual counts: target POSTs0, listener POSTs0, request_audit[]; **one proxy credential read and one models GET occurred** under the approved release. No model tool sequence/hear provenance/listener report/model-token usage was produced, so the blind-listener/content task was not exercised or passed. Discovery response and local process/profile/save artifacts retained; credential file was not manually inspected or exposed. No secret argv/log/file/child-env handling was added or changed.

Tracked cleanup_complete:true: driver1706 exit1/group_gone:true; Chrome1686 and GUI1682 exit0/group_gone:true; proxy closed. No broader OS cleanup claim or manual/shared cleanup. Original4target+8listener/12POST/1024native cap/60setup/240active/300absolute bounds unchanged; unused budget is not authorization. Muted/transcript-only setup, no new native run/MIC/device action.

Bill’s limited native-v6 hearing testimony remains separate and unchanged; prior failed attempts remain failed. No source/controller/binary/gate changes, build, children, commit, auto-repair or further operation. **HOLD for root inspection and decision; no further paid authority.**

## Controller v7 listener invocation/goal preparation — UNEXECUTED HOLD

Root classified listener-v6 as controller-only invocation failure, not production R6. Exact diff `/private/tmp/ensemble-ch13-ops.Gk2JOo/operate-v7.diff`: driver receives `['--attach','--listener',INITIAL_TARGET_PROMPT]`, exactly one positional argument: **“In one short sentence, explain queueChunk, max_tool_rounds, and HTTPServer using inline code markup.”** Goal first waits/hears the already-submitted turn, expressly forbids repeating that initial submission, then asks exactly one follow-up: **“Choose one fruit to mention and repeat the explanation in plain words.”** Truthful speech reporting/omissions and no DOM/history/target-engine/provider-wire access remain required; no exact-answer compliance assertion added.

Fresh listener-v7 path. Native-v6 remains the accepted prerequisite (same pinned receipt/pins); no native replay. Corrected the old same-pins comparison to the accepted native-v6 pins identity so versioning does not falsely invalidate the accepted run. Native path itself remains native-v6/no-overwrite. Five tools,4target+8listener/12POST,1024 cap,60/240/300s limits, muting/proxy/secrets/cleanup unchanged.

Controller SHA256 `a8a3c7cdec6104e4d2004f42f6a4f6b840308a0042a9385e07eeff52e3c05973`; pins-v7 SHA256 `fd282d3ff3b9281cf9846736a55c39ef9cdf1b7e20079771a81b196437089b3f` (83 inputs). Syntax0; focused added pure checks0 against actual held run()/runSession(): one argument, type→submit before listener, goal wait/hear-first plus one follow-up, native-v6 status/pins/cleanup. `syntax-v7.log`, `pure-contract-v7.log`, `preparation-v7-exits.log`. All pinned inputs verified. No application gates/test campaign, source/binary rebuild, browser/native/provider/key/GET/POST/paid/listener session or cleanup.

Exact future HELD command, NOT executed/authorized:
```
D=/private/tmp/ensemble-ch13-ops.Gk2JOo
P=fd282d3ff3b9281cf9846736a55c39ef9cdf1b7e20079771a81b196437089b3f
N=ce2879fa5addd020137e4456b05bb1ea86c6b6808ae53131e2452e3951cadaf8
env -i HOME="$D/home" TMPDIR="$D/tmp" PATH=/usr/bin:/bin:/usr/sbin "$D/operate-v7.mjs" listener "$P" "$N"
```
Prior bytes/results remain unchanged; listener-v6 remains FAILED with its one credential read/models GET and zero POSTs. Root owns missed invocation assumption. No children/commit/shared cleanup. **HOLD for inspection; no further operational/paid authority.**

## ONE released listener-v7 operation — controller FAILED; actual task/report recovered read-only; HOLD

Exact held v7 listener command executed once, live PTY, native-v6 prerequisite. Start disk20,986,744KiB (>1GiB). **Controller process1 / receipt failed**, error `listener report not preserved`; driver actually exited0. No repair/retry/fallback or additional paid work.

Receipt `/private/tmp/ensemble-ch13-ops.Gk2JOo/listener-v7/receipt.json` SHA256 `fed5a93232be9ebf1d2a11f603e9a5fa430bd6cf61d26dbb5ae5568f2f1e11d8`; outer stdout/stderr/exit retained (driver log/outer stderr empty). Started2026-10-21T12:40:05.050Z, ready06.058 (~1.008s), finished22.306 (~17.256s). Original receipt unchanged.

Actual counts:2 target+7 listener=9 model POSTs, all claude-sonnet-4-6/max_tokens1024; one proxy credential read/models GET. Each listener request exposed exactly hear/type/submit/wait/sleep; target tools[]. Native output caps/deadlines/budgets unchanged; no real key inspected or exposed. Saved-response usage: target181 input/114 output; listener11,114 input/542 output; total11,295/656, cache creation/read0. This is token usage, not a USD estimate.

Actual listener tool sequence: **wait → hear → type → submit → wait → hear → final report**. All six tool results non-error; both waits returned success/completed. Hear cursors6 and12 delivered seq[3,4,5,6] and[9,10,11,12], mode transcript. Read-only comparison of saved tool results to page-final events matched event/text/source/turn/part fields. Proxy's existing complete provenance guard also admitted those listener requests. No DOM/history/raw target-wire access was supplied by the five tools.

Report was present at **`listener-v7/tmp/ch13-speech-4047728693/listener.txt`**, despite controller lookup failure. It quotes both admitted assistant responses and identifies **oranges**. Assistant channel explains all three requested concepts, then repeats using oranges/conveyor-belt metaphor. No requested concept/fruit omission apparent in retained speech/report; no arbitrary exact-answer/definition-compliance assertion. Listener paraphrase uses `maxToolRounds`; filtered speech does not prove reconstruction of original identifier punctuation. Report and all raw records remain local, not moved/rewritten to heal the receipt.

Final page: idle/success, turn2, draft empty, human_typing/paused/busy:false. Four channel utterances (two human prompts/two assistant replies), transcript mode, no native-start records; native speaking/pending/paused:false. This is blind speech-channel evidence, **not hearing/locality/acoustic-comprehension proof**.

Tracked cleanup_complete:true: driver1940, Chrome1920, GUI1916 all exit0/group_gone:true; proxy closed. No broader OS cleanup claim/manual shared cleanup. All older failures—including listener-v6—remain failed. Source/controllers/binaries/gates unchanged; no builds/tests/children/commit, key retry/discovery retry or further operation. **HOLD for root inspection; controller overall remains failed/process1, and unused budget is not additional authority.**
