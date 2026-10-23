# Edition 3 student inputs

Read this with the shared [execution procedure](edition-3-execution.md), the complete mandatory `.agents/skills/ensemble-coding/SKILL.md`, and the assigned chapter before starting and after compaction. These current paths/phases replace historical edition instructions. Architecture, KISS, testing, documentation and review obligations remain.

## Current boundary

Chapters 1–19 are accepted. Chapter 19 main and its 255-file snapshot are held after code, local/native validation, the bounded required CLI behavior and minimal editorial acceptance. All Chapter 19 workers are closed. Its existing workspace-confinement architecture remains; it is not a chroot redesign. Root fixture repairs did not add student production rounds. Earlier failed results remain preserved and qualified. Chapter 20 is next, but its contract, prerequisites and grading route have not yet been assessed and no Chapter 20 student is launched. Bill authorizes live tests for the remaining book without routine permission requests; root still owns bounded operational plans, and account consent, installation/restart, devices and publication remain separate boundaries. Do not replay completed Chapter 19 work or treat these permissions as permission for a student to operate outside its assigned scope.
Earlier chapter invocation sections below remain historical/public guidance,
not authority to replay completed work or reuse live-operation permissions.

## Permitted inputs

- The assigned `book/edition-3/chapter-NN.md`, Chapter 00 for orientation, and earlier chapters as needed.
- Your own working tree in `/Users/bill/projects/ensemble/solutions/edition-3/main` and the specific accepted predecessor snapshot supplied by the orchestrator under `solutions/edition-3/chNN`. Chapter 1 starts empty; other snapshots are not automatically permitted inputs.
- The current role, shared procedure, this public brief, and narrow public-contract clarifications explicitly supplied by the orchestrator. Mandatory governing instructions may contain incidental historical path/revision references; reading those instructions is allowed, but following those pointers into old solutions or repair evidence is not.
- Public dependency/provider documentation needed for the assigned exercise.

Apart from your assigned workspace and permitted predecessor, do not read `agent/`, other `solutions/` trees, `edition-1/`, future chapters, grader implementation, historical answer/repair notes, or their Git blobs to solve the exercise. Do not read the coordinator's readiness/carryover/evidence records unless a specific student-safe excerpt is supplied. The historical skill's direction to read Edition 2 carryover is replaced for the fresh Edition 3 student by this brief; authors/reviewers still read their complete governing documents. Do not follow answer-revealing links from permitted documents. These are input-separation rules, not a claim of OS sandboxing.

## Chapter boundary

Chapter 1 starts from an empty source tree. Later chapters extend your exact accepted predecessor; do not copy the full reference. The orchestrator supplies the assigned chapter revision, accepted predecessor and grader revision. The orchestrator owns snapshots and repository commits; do not create a nested repository or edit other main-repository files. Submit one concise author-facing chapter review, not an evidence bundle. No later capabilities early. Implement required behavior with the simplest adequate design, with useful explanations of non-obvious choices.

Run the assigned chapter grader with an explicit positional target. Chapter 1 accepts a supported solution directory or a built CLI executable. If your public root is a library with a nested CLI, follow the chapter's “Grade yourself” build instructions rather than grading the library directory. For the executable, from `/Users/bill/projects/ensemble` (the repository-root Go module, not `agent/`):

```sh
go run ./cmd/grade-edition3 -ch 1 -json /absolute/path/to/your/student-cli
```

Substitute the executable you actually built; its package layout is your implementation choice.

For Chapter 2, follow its revised exercise and run both existing graders against
that same built executable from the book repository root:

```sh
go run ./cmd/grade-edition2 -ch 2 -json /absolute/path/to/your/student-cli
go run ./cmd/grade-edition3 -ch 1 -json /absolute/path/to/your/student-cli
```

Both must pass and exit successfully. The second command preserves inherited
Chapter 1 obligations; it does not replace the Chapter 2 pass.

For Chapter 3, build the actual CLI package and use the same executable for:

```sh
go run ./cmd/grade-edition2 -ch 3 -json /absolute/path/to/your/student-cli
go run ./cmd/grade-edition3 -ch 1 -json /absolute/path/to/your/student-cli
```

Both must pass. The Chapter 3 grader includes Chapter 2 parity; the second
command retains the inherited lifetime checks. Follow the chapter’s public-API
tests and source review for the selected-model explanation policy: grader
scores alone do not establish it. The chapter leaves the policy-selection
interface to the implementation; no hidden activation flag is required.

For Chapter 4, the Edition 3 command exposes the chapter's existing grading
checks with separate executable, core-source and external-consumer identities.
Build your actual CLI; do not turn a public library root into `main` to suit a
builder. From the book repository root run:

```sh
go run ./cmd/grade-edition3 -ch 4 -json -core /absolute/path/to/your/core-module -consumer /absolute/path/to/your/consumer-module /absolute/path/to/your/student-cli
go run ./cmd/grade-edition3 -ch 1 -json /absolute/path/to/your/student-cli
```

For a consumer other than the default `shout` example, add
`-consumer-case /absolute/path/to/case.json` before the executable, using the
chapter's public case schema. The first report retains the chapter's 160-point
jobs/inherited/consumer checks; the second retains the lifetime obligation on
that same executable. Both must pass. The core path must be the actual module
resolved by the consumer, not a similarly named copy. Required focused tests,
source ownership review, job lifecycle/cursor/limit checks and inherited
selected-model policy evidence remain necessary beyond those reports.

For Chapter 5, build the actual CLI and keep the core, external consumer and
three-Agent workflow as explicit module paths. Run from the book repository root:

```sh
go run ./cmd/grade-edition3 -ch 5 -json -core /absolute/path/to/your/core-module -consumer /absolute/path/to/your/consumer-module -workflow /absolute/path/to/your/workflow-module /absolute/path/to/your/student-cli
go run ./cmd/grade-edition3 -ch 1 -json /absolute/path/to/your/student-cli
```

Put optional `-consumer-case /absolute/case.json` and
`-media-models /absolute/media-map.json` before the executable. The media mapping
is a JSON object with exactly the nonempty `anthropic`, `openai`, and `gemini`
model IDs selecting your explicit capability rows; it selects grader inputs,
not live-provider capabilities. The case uses the chapter's public schema.
Both reports must pass on that same executable. The first retains Chapter 5's
seven checks and 100-point total, including Chapter 4 parity; the second
retains inherited lifetime checks. Public actor/reply/observer contracts,
interruption and job behavior, round-budget semantics, media loading, actual
workflow tool effects and human-interface exercises still require their own
focused tests and independent acceptance beyond those scores.

For Chapter 6, build the actual CLI and identify the core, inherited consumer,
workflow, and separate public-API streaming exercise modules explicitly. From
the book repository root run:

```sh
go run ./cmd/grade-edition3 -ch 6 -json -core /absolute/path/to/your/core-module -consumer /absolute/path/to/your/consumer-module -workflow /absolute/path/to/your/workflow-module -stream-exercise /absolute/path/to/your/stream-exercise-module -stream-model YOUR_DECLARED_STREAMING_MODEL /absolute/path/to/your/student-cli
go run ./cmd/grade-edition3 -ch 1 -json /absolute/path/to/your/student-cli
```

The streaming model is required and must select an explicit capability row in
your core; the adapter does not invent one. Put optional `-consumer-case` and
`-media-models` JSON paths before the executable as above. The first report
retains the existing Chapter 6 checks and Chapter 5 parity; the second retains
the lifetime obligation on the same binary. Both must pass and exit zero.
Actual arrival order, early human CLI delivery, SSE framing, source-bound replay,
and inherited job wake-up behavior still require the chapter's public-interface
tests, source review and independent acceptance beyond those reports.

For Chapter 7, keep the accepted core CLI and inherited modules separate from the new GUI module. Build the actual CLI target, then run from the book repository root:

```sh
go run ./cmd/grade-edition3 -ch 7 -json \
  -core /absolute/path/to/core-module \
  -consumer /absolute/path/to/consumer-module \
  -workflow /absolute/path/to/workflow-module \
  -stream-exercise /absolute/path/to/stream-exercise-module \
  -stream-model YOUR_DECLARED_STREAMING_MODEL \
  -gui /absolute/path/to/gui-module \
  -gui-model YOUR_DECLARED_GUI_MODEL \
  /absolute/path/to/student-cli
go run ./cmd/grade-edition3 -ch 1 -json /absolute/path/to/SAME/student-cli
```

All five module roots and both model identifiers are explicit; flags precede the executable. Optional `-consumer-case` and `-media-models` paths are forwarded unchanged to the inherited checks. Use models declared by your implementation, not historical defaults or invented feature rows. A passing WebSocket grade does not substitute for the chapter's actual browser, native-speech, replay, and pause exercises. Run isolated servers and browser contexts without disrupting the orchestrator's live GUI. The same frozen core CLI must pass Chapter 7 and Chapter 1, including the real delayed-response phases.

For Chapter 8, use the same separate executable and five module identities. After building the actual CLI target, run from the book repository root:

```sh
go run ./cmd/grade-edition3 -ch 8 -json \
  -core /absolute/path/to/core-module \
  -consumer /absolute/path/to/consumer-module \
  -workflow /absolute/path/to/workflow-module \
  -stream-exercise /absolute/path/to/stream-exercise-module \
  -stream-model YOUR_DECLARED_STREAMING_MODEL \
  -gui /absolute/path/to/gui-module \
  -gui-model YOUR_DECLARED_GUI_MODEL \
  /absolute/path/to/student-cli
go run ./cmd/grade-edition3 -ch 1 -json /absolute/path/to/SAME/student-cli
```

Put optional `-consumer-case` and `-media-models` paths before the executable as above. Both reports must pass and exit zero on the same frozen CLI, including the real delayed-header/body phases. Settings snapshots, restart persistence, and Chapter 7 parity do not replace Chapter 8's public source/integration obligations: effective next-turn budget use, validation/save ordering, and actual browser controls and speech-rate propagation still need evidence. No particular budget-setting API is prescribed.

For Chapter 9, use the accepted Chapter 8 predecessor and the same separate executable and module identities. After building your actual CLI target, run from the book repository root:

```sh
go run ./cmd/grade-edition3 -ch 9 -json \
  -core /absolute/path/to/core-module \
  -consumer /absolute/path/to/consumer-module \
  -workflow /absolute/path/to/workflow-module \
  -stream-exercise /absolute/path/to/stream-exercise-module \
  -stream-model YOUR_DECLARED_STREAMING_MODEL \
  -gui /absolute/path/to/gui-module \
  -gui-model YOUR_DECLARED_GUI_MODEL \
  /absolute/path/to/student-cli
go run ./cmd/grade-edition3 -ch 1 -json /absolute/path/to/SAME/student-cli
```

Optional `-consumer-case` and `-media-models` paths precede the executable. Declare models supported by your implementation rather than adding feature rows to satisfy historical defaults. Both reports must pass and exit zero on the same frozen CLI, including the real delayed-header/body phases. The skills and inherited settings checks do not establish every public obligation: full-graph validation before activation, failed-load atomicity, Agent isolation, variable expansion outside the Agent mutex, and durable enabling-event provenance also require source and integration evidence. `unload_skill` records pending removal only in this chapter; it does not implement future compaction. Preserve the chapter's real GUI/manual exercise without expanding it into a provider campaign.

For Chapter 10, extend the accepted Chapter 9 predecessor. Build your actual CLI target, then use the same explicit artifact identities from the book repository root:

```sh
go run ./cmd/grade-edition3 -ch 10 -json \
  -core /absolute/path/to/core-module \
  -consumer /absolute/path/to/consumer-module \
  -workflow /absolute/path/to/workflow-module \
  -stream-exercise /absolute/path/to/stream-exercise-module \
  -stream-model YOUR_DECLARED_STREAMING_MODEL \
  -gui /absolute/path/to/gui-module \
  -gui-model YOUR_DECLARED_GUI_MODEL \
  /absolute/path/to/student-cli
go run ./cmd/grade-edition3 -ch 1 -json /absolute/path/to/SAME/student-cli
```

Optional `-consumer-case` and `-media-models` paths go before the executable. Both reports must pass and exit zero on the same frozen CLI, allowing the actual delayed-header and delayed-body lifetime phases to finish. The persistence checks and inherited Chapter 9 parity do not establish every public obligation: rich replay through all three vendor renderers, atomic load/refusal, effective-running configuration boundaries, and stable GUI addresses across restart—including already-colliding legacy archives—also require source and integration evidence. Preserve provider payloads and parser-local IDs; no new public identity API or future compaction implementation is prescribed. Use owned save paths for probes and the chapter's real CLI/GUI restart exercise, without importing private histories or expanding it into a multi-provider paid campaign.

For Chapter 11, extend accepted Chapter 10 and use the same command and separate artifact/model flags above with `-ch 11`, followed by `-ch 1` on the **same frozen CLI**. Optional case/media flags still precede the executable. Require both parsed pass and exit zero, preserving the genuine delayed-header/body waits. The MCP checks and inherited persistence parity are finite: review every taught transport, request correlation, declared-skill startup and tool filtering, lifecycle/cancellation, and fresh generic automatic context (including valid empty refresh). `gui_snapshot` and `tts_queue` remain on-demand tools. The browser endpoint must be enabled by a declared skill, not an empty-skill fallback. Complete the chapter’s bounded real-browser exercise in an isolated owned session, without a provider campaign or live-GUI disruption.

For Chapter 12, extend accepted Chapter 11 and use the same explicit artifact/model flags above with `-ch 12`, followed by `-ch 1` on the **same frozen CLI**. The Chapter 12 route checks startup skill selection, explicitly requested snapshots and speech queues, real click/input effects, unload refusal/peer teardown, and Chapter 11 parity. Browser observations must not be injected automatically; legitimate tool-result history may remain after unload. The fake peer does not replace the bounded real-browser exercise. Preserve Chapter 9's dynamic SYSTEM instructions and provider-specific disclosure rules. Follow the phase authorization: a compiling hold and source review precede full gates and operational work.

The Edition 3 command requires an explicit target and currently supports Chapters 1, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16 and 18. Its Chapter 1 mode preserves the inherited conversation checks and adds delayed-header/body lifetime gates. Check the report's `passed` field and exit status, not score alone: required zero-point checks can fail a 100/100 score. Those real delayed probes take about two minutes, so a 90-second suite timeout is invalid. Omitting the path from the older Edition 2 command silently targets an old solution and is not a valid student result. A grader result does not replace required live exercises, test-sensitivity evidence or independent acceptance. Credentials remain in the authorized provider loader/test process, never prompts, command arguments, evidence or commits.

Submit one concise chapter review identifying the chapter/code/grader revisions, what you built, actual tests and live observations, failures, ambiguity, assistance and useful teaching clarifications. You may explain a suspected grader defect; you may not silently redefine success or edit the grader. The orchestrator supplies any missing public contract without old answer code. Follow the shared bounded review loop.


Chapter 13 uses the same explicit core/consumer/workflow/stream/GUI identities and flags-first ordering as Chapter 12, with `-ch 13`. For the accepted nested layout, build the CLI from `./cmd/ch01`; use `-workflow /absolute/student/ch06` and `-stream-exercise /absolute/student/ch07`, not the reverse. Supply `-core /absolute/student`, `-consumer /absolute/student/examples/ch04-consumer`, `-gui /absolute/student/gui`, and explicit `-stream-model` and `-gui-model` before the executable. The speech harness is discovered under the core directory (prefer `scripts/tts-harness.sh` or `tts-harness.sh`) and must honor the chapter's environment and JSONL contract while driving the real GUI. Run from the book repository with `go run ./cmd/grade-edition3 -ch 13 -json ... /absolute/student-cli`, then Chapter 1 on the SAME frozen CLI after Chapter 13 passes. Both require `passed: true` and process exit zero; allow the actual 62-second header and body checks to finish. The grader's browser/synth boundary is controlled evidence, not native hearing or the blind listener exercise. Operational execution requires the root's bounded release.


Chapter 14 directly grades the **GUI executable or GUI module directory**, not the CLI or core library. Build and freeze the GUI and CLI from the same held source, then run from the book repository:

```sh
go run ./cmd/grade-edition3 -ch 14 -json /absolute/owned-output/ch14-gui
go run ./cmd/grade-edition3 -ch 1 -json /absolute/owned-output/ch14-cli
```

For the accepted layout, build `ch14-cli` from `main` with `go build -o /absolute/owned-output/ch14-cli ./cmd/ch01`, and `ch14-gui` from `main/gui` with `go build -o /absolute/owned-output/ch14-gui ./cmd`. Chapter 14 rejects unrelated `-core`, consumer, workflow, stream, media and GUI/model flags; its fixture owns its models and assets. Require parsed pass and process exit zero for both reports, allowing the genuine Chapter 1 62-second header and body phases. Follow phase authorization: a compiling HOLD and architectural review come before these gates. The finite grade does not replace writer-ordering, failed-save/diagnostic preservation, backup-tail overlap, total-reducer, provider-disclosure and scoped operational evidence required by the public chapter. No numeric canonical-event retention cap replaces the validated-BACKUP boundary.

Chapter 15 also grades the **GUI executable or GUI module** directly. Its cache
artifact paths default to `cr/cache-current.txt` and `cr/cache-prior.txt` in the
owned launch directory; override with `-cache-current`/`-cache-prior` only to
match the implementation's actual path. Supply four independently checked
per-million prices with `-cache-prices INPUT,WRITE,READ,OUTPUT`; the adapter
requires them rather than deriving expected costs from student code. Other
chapters reject these flags. For the dated local fixture tariff, from the book
root after building held GUI/CLI artifacts in their own modules:

```sh
go run ./cmd/grade-edition3 -ch 15 -json -cache-prices 3,3.75,0.3,15 /absolute/output/ch15-gui
go run ./cmd/grade-edition3 -ch 1 -json /absolute/output/ch15-cli
```

Preserve all Chapter 14 instruction, recovery, and writer-boundary behavior.
The Chapter 15 grade is finite fake-provider evidence, not a real cache grant,
browser inspection, price verification, or architecture acceptance. Follow the
current phase authorization before running it or paid exercises.

Chapter 16 grades the **GUI executable or GUI module** directly; it takes no
core/parity/model/media/cache flags. From the book repository root, after the
held source is reviewed and the validation phase is released:

```sh
go run ./cmd/grade-edition3 -ch 16 -json /absolute/output/ch16-gui
go run ./cmd/grade-edition3 -ch 1 -json /absolute/output/ch16-cli
```

Build the CLI from accepted `main` with `go build -o /absolute/output/ch16-cli
./cmd/ch01`, and the separate GUI from `main/gui` with `go build -o
/absolute/output/ch16-gui ./cmd`. Both command entrypoints delegate auth commands
to the same owner; the offline GUI-direct exercise invokes `auth login
--print-url` on its supplied executable. The public chapter declares the
synthetic capability row and local override names. Require parsed pass and
exit zero for both reports, with at least 360 seconds allowed for Chapter 1's
real 62-second header and body checks. Keep capability routing conservative;
provider labels alone do not prove Responses, summary or cache support.
The fake OAuth/Responses grade is not personal-account authorization, proof
of live entitlement, native speech, or complete security coverage.


## Chapter 18 restart invocation

Chapter 17 is a field report, so extend accepted Chapter 16 in `main`;
`../ch16` is the permitted predecessor. Build the actual nested CLI and
separate GUI, without changing the public library layout:

```sh
(cd solutions/edition-3/main && go build -o /private/tmp/ensemble-ch18 ./cmd/ch01)
(cd solutions/edition-3/main/gui && go build -o /private/tmp/ensemble-ch18-gui ./cmd)
go run ./cmd/grade-edition3 -ch 18 -json \
  -core "$PWD/solutions/edition-3/main" -gui /private/tmp/ensemble-ch18-gui \
  -status-prices 3,3.75,0.3,15 /private/tmp/ensemble-ch18
go run ./cmd/grade-edition3 -ch 1 -json /private/tmp/ensemble-ch18
```

All flags precede the explicit CLI target. Supply the four independent
input/write/read/output rates per million tokens. The chapter documents
these dated Sonnet 4.6 fixture rates and the optional Opus 4.6 tariff
`-status-second-prices 5,6.25,0.5,25`; omitting the latter selects the
unpriced middle-model control. Both reports must have `passed:true` and
exit zero. Allow at least 360 seconds for Chapter 1. The coordinator
releases these gates only after the coherent compiling HOLD and review.
Source review still follows real owners and every shipped composition
root; numerical grading alone is insufficient. Existing `ToolCall`,
`Parent`, and `Status` seams are valid without aliases for reference names.
Live account consent, provider requests, and hosted web operations require
their own explicit authorization; historical examples do not supply it.


## Restarted Chapter 19 public entry brief

Start from accepted Chapter 18 in `solutions/edition-3/main` after root release.
Read the complete current Chapter 19 and required public prerequisites. Extend
existing Config/common.Config with policy and Child/Clamp; preserve the root
library and nested CLI/GUI owners. Sandbox retains its required Agent parent.
Read actual file, process, job-output, state, media, transport and render owners
rather than add a second security or credential framework around them.

Implement public confinement, safe mode, the web-skill gate, environment and
redaction behavior with positive useful-work cases as well as refusals. Share
CLI/GUI policy composition. Public commands use `./cmd/ch01` and flags before
`chat`. Account consent and host/provider demonstrations need separate release;
historical prose is not new operational evidence.

First stop at an early compiling HOLD, not acceptance. Authorized work upon
root release is implementation, focused pure/local tests and affected-module
compilation, sequential/offline with owned HOME/TMP and existing caches. No
full suites, race/vet, grades, Seatbelt/host scenarios, browser/native/provider
or live OAuth work yet. Save actual changed owners, commands/exits, setup
failures and unfinished work concisely. Root and the retained reviewer inspect
held source before broader gates. Initial implementation is round zero; up to
ten production correction rounds are available. Originals, private graders and
reviews, runtime and root context remain outside student inputs.
