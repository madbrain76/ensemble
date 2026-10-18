# Edition 3 Chapter 7 — student implementation report

## Disposition: frozen candidate; NOT accepted / Chapter 7 not complete

The final clean student candidate is **r4**, commit
`48e1d12d9c396b5988c0ce6eccdef67e581872ff`, tree
`0f5f7fbe4273e91e1be3aa13d0715a99f671eaf9`, in
`/Users/bill/projects/ensemble-edition-3-student`.

**Public Chapter 7: 80/100, `passed:false`.** Its `pause-holds-tools` check reports
`pause test error: pause allowed 3 tool starts before other-tab unpause`.
The corrected cumulative Chapter 6 parity check passes. The same frozen CLI's
Chapter 1 gate passes 100/100, including both zero-point 62-second delayed-header
and delayed-body checks. Final core/GUI race tests, vet, inherited module
build/test/vet, and controlled production-JavaScript queue tests pass.

This is a failed candidate submission for independent review, not an acceptance
request disguised by a partial score. Native speech completion/listening quality
and successful recognition are also **not established**. Parent owns independent
review, any public-grader reconciliation, acceptance, and tags. No acceptance tag
was made.

## Exact identities and reproducibility

- Accepted predecessor (unchanged): `f595c787675e2a6343053ec42ea6c7e1e02a930e`.
- Immutable `edition-3-ch06-r1` tag object: `84d0bce5f2525b406affa64fb1a8eb70711b1e06`;
  its peeled commit is the predecessor above, not the tag-object SHA.
- Frozen CLI: `.ch07-local/frozen-r4/ensemble-cli`, SHA256
  `4beda0aa805195229c5987a025a2b794b4fe3e77a6526ff194927282fddd85cb`.
- Raw-probe GUI: `.ch07-local/ensemble-gui-frozen-r4`, SHA256
  `39ebcd3e42a57d8972b6278fb8cc66021e38afd41afa2cd230b93499e6b13a81`.
- Source archive: `.ch07-local/frozen-r4/candidate-r4-source.tar.gz`, SHA256
  `7beef80a0af5282c793954d738f1df9da51c510051f68bc2ab5454278f304624`.
- Full per-file hashes, build artifacts, predecessor and clean status are in
  `ch07-source-manifest.json` and `ch07-student-evidence/candidate-r4-manifest.json`.
- `ch07-evidence-manifest.json` indexes portable evidence and report hashes;
  local binary/archive/screenshot paths and hashes are separately indexed.

Public commands, working directories, complete deliberately minimal child
environments, timestamps, and exit statuses are recorded verbatim in
`chapter7-r4-invocation.json` and `chapter1-r4-invocation.json` under the evidence
directory. These use the supported public `go run ./cmd/grade-edition3 -ch 7`
adapter with `-core`, inherited `-consumer`/`-consumer-case`, `-workflow`,
`-stream-exercise`, `-stream-model claude-sonnet-4-6`, `-gui`,
`-gui-model claude-sonnet-4-6`, explicit declared media rows, and the exact frozen
CLI. The Chapter 1 invocation uses that same CLI. No unsupported per-chapter
modes were invented. Public grader HEAD is
`23c75ede7944d9a912ae74d3f683354d163f0c86`; recorded worktree entrypoint/test
hashes match the final hash-only check in `public-grader-final-fingerprint.json`.
No grader source/case bodies or private audit material were read.

Important artifact distinction: the public grader builds a GUI from the frozen
module source path; the student did not retain the hash of its internal temporary
GUI binary. The GUI hash above belongs to the independently built `-trimpath`
artifact used by the student's real protocol/effect probes. Do not conflate them.

## Prerequisites, scope, and ownership

The 8,395-byte task matched SHA256
`8d34f80dfce2dfc506a2dfd1b3f3a18129279227846a738914fd85619d882d2f`.
The public Chapter 7 matched SHA256
`68d02d8e1170c2912676b8f162847f8fac25c87216514ed06adbd5f0ea011b47`.
I read the mandatory public book instructions, book skill, Edition 3 execution
procedure, student brief, carryover, and full Chapter 7; abbreviated tool reads
were recovered in bounded ranges. I verified the clean predecessor/tag, sent an
ownership/test plan, and waited for parent approval before edits/browser work.
I acted as student implementer, not author/grader, and spawned no children.

Core ownership remains in the existing Root/Agent/jobs/tools layers:

- Actor-owned mutex-protected pause gate with cancellable pending model/tool
  boundaries; synchronous public `Pause`, `Unpause`, `Paused` control. The actor
  continues processing mail, hints, interruption, completion, and shutdown.
- Agent-assigned persistent display identities, separate from the unchanged
  provider-parser local `PartID` grammar. Existing Chapter 5/6 observations remain
  compatible. Immutable finalized history is committed before final observations.
- Identified synchronous root observers receive defensive event copies. Added
  tool-dispatch/tool-finish/turn-ended facts remain presentation-neutral.
- Synchronized immutable history suffix snapshots; no transport network ownership
  or GUI tool/job execution moved into core.
- Retained processes preserve real lifetime/outcomes across paused wait and
  interrupt. A new root prompt is not replaced by browser-side execution.

The optional `gui/` module imports public core APIs and embeds `gui/web/gui/`.
It provides `gui/cmd`, the documented launch flags/environment precedence,
terminal prompt/busy-hint forwarding with EOF detach, and local debug/event/API/
WebSocket logs. Its bounded, single-writer WebSocket peers disconnect on overflow
rather than blocking actor callbacks. Ping/control and fragmented-message support
were added after independently reproducing actual r1 transport failures.

The hub captures durable suffixes, projects finalized artifacts, retains active
prefixes, cleans up finals while tools can still be pending, and makes replay/live
admission atomic. The window is last hour or last 100 renderable records, whichever
is larger. Display IDs persist across repeated parser/call IDs. Exact-once tool
admission bookkeeping covers live-before-commit, late capture, and replay. Final
subscription envelopes are cloned and explicitly `replay:true`; genuine live
progress and stored projection objects stay unmarked. Snapshot `from` and `to`
use one sampled state. History types/content/order/identities are not suppressed
or changed to hide tool-start counts.

The browser has safe DOM/Markdown/code rendering, structured tool views, a shared
scroll helper, per-artifact phrase queues, real streaming append, final-drain and
reconnect-generation handling, a real speech queue and per-artifact speakers,
opaque/readable thinking separation, draft/voice pause causes, explicit global
release, and terminal-state-aware prompt/hint routing. Browser interruption waits
for turn acknowledgement before releasing its local pause cause.

## Final test results and attribution

| Evidence | Actual outcome / scope |
|---|---|
| `chapter7-r4-report.json` | **FAIL 80/100**: all-observation wire, event replay, GUI log and corrected Chapter 6 parity pass; pause gate fails with three starts. |
| `chapter1-r4-report.json` | PASS 100/100 on exact final CLI; both 62-second zero-point lifetime checks pass. No shorter artificial deadline stopped them. |
| `frozen-r4-matrix.log` | PASS: core `go test -race ./...` (~140s), core vet; GUI race/vet; inherited consumer/workflow/streaming module race/test/vet/build; production JS queue checks; clean diff/status before and after. |
| `protocol-frozen-r4.json` | Actual frozen GUI: ping/pong and fragmented subscribe pass. Real loopback process/socket, not mocked callbacks. |
| `pause-boundary-r4-fresh.json` | Fresh Agent: 0 replayed starts; before other-tab unpause 1 live admission, first file effect present, second absent; afterward 2 live admissions and second file effect present. |
| `pause-boundary-r4-warm.json` | Warm history: 2 historical dispatches plus 1 live admission; second real process effect remains absent until explicit other-tab unpause, then appears. Raw replay markers, full wire rows and owned process receipts retained. |
| `mutation-controls-confirmed.json` | Three compile-valid intended failures and exact restoration/pass: bypass pause gate; truncate larger replay window to 100; omit partial cleanup. Earlier compile/setup/anchor failures are separate, not sensitivity proof. |
| `dispatch-exactly-once-red.log`, `dispatch-r3-final-green.log` | Old behavior compiled and failed join-gap, duplicate-live and late-duplicate assertions; corrected behavior passes. |
| `replay-marker-valid-red.log`, `replay-marker-snapshot-green.log` | Old behavior compiled and failed explicit final/partial/state replay markers; clone-only marker/snapshot change passes full GUI race/vet. |
| `browser-frozen-r1-results.json` and safe `runtime/` logs | Real isolated Chrome 155/CDP browser observations on **r1 GUI bytes**, not r4 runtime proof; see below. |

Focused Go tests additionally cover 270-record replay-window sensitivity, capture
before callback, active-prefix replay, cleanup while tools are pending, bounded
slow-client isolation, single-write-witness replay-before-live ordering, safe
thinking projection, tool-dispatch handoff/turn cleanup, WebSocket control/
continuation limits, pause before model, hints/interruption/resumption, streaming
completion during pause, retained-process survival, and shutdown while paused.
Names and exact source are in `chapter7_test.go`, `gui/hub_test.go`, and
`gui/socket_test.go`.

### Unresolved public pause discrepancy — not a grader-defect conclusion

`r4-pause-discrepancy.json` packages the exact public failure, source/artifact pins,
command/model/environment provenance, and fresh/warm raw controls. The controls
prove their own real process boundaries, not a separate grader's scenario. They
separate tool-request artifacts, historical observations, and actual newly
admitted execution/effects. They **do not establish why the public gate counted
three**, and do not justify claiming the grader is wrong. On parent instruction,
I stopped all speculative source changes and new execution campaigns after r4.
Independent root reconciliation is required before any acceptance decision.

## Browser/native evidence and limitations

Earlier authorized browser work used only a newly launched headless Chrome process
and empty profile, unique non-8088 ports, and an owned scripted local HTTP peer.
The shared/live GUI/CDP target was not attached, restarted, or reused. Owned
processes were shut down and safe fixture/API/wire/DOM logs retained; profile
contents were not copied into evidence. No paid/model-network calls were made.

The final r1 browser run observed: exact served asset/source matching; typing
pause with read-aloud off; cross-tab explicit release winning an unchanged draft;
real partial DOM during a stream and partial reload/final convergence; safe hostile
Markdown/HTML/link handling; unique repeated-turn IDs; opaque/empty exclusion;
current-tool completion while next dispatch is held; interrupted paused wait with
retained process survival and skipped undispatched tool; next turn reuse; denied
microphone permission releasing pause; terminal prompt/busy hint and EOF detach.
Those assertions are in the retained receipt. Browser assets and core hashes are
unchanged from r1 through r4, but that is source carryover, **not an r4 browser run**.
No browser/speech/mic exercise was repeated after parent stopped that campaign.

Controlled tests exercise the actual production PhraseQueue/SpeechQueue classes
with deterministic schedulers/callbacks, including ordering, identity, end/error,
cancel/reset and stale callback rejection. They are not native speech evidence.
One earlier development run observed native start callbacks. Frozen r1's muted
headless attempt did not observe end/error within ten seconds; explicit cancellation
released pause. Native availability was **not** conclusively ruled out. No final-r4
native completion, listening-quality, recognition-success, or full voice usability
claim is made. Early headless runs did not yet pass `--mute-audio`; no human
listening observation was made, so inaudibility of those early attempts is not
asserted. Subsequent bounded probes explicitly used mute.

`native-probe-proposed.mjs` is a separately permission-gated, syntax-checked but
**unexecuted** proposal pinned to the final GUI hash. It refuses default launch,
uses a new empty owned profile/CDP endpoint, an owned local scripted HTTP peer,
non-8088 ports, no focus/frontmost calls and no mic API; it defaults headless/muted,
requires separate explicit headed/audible approvals, restricts page requests to
local origin, requires an advertised default local-service voice, records native
callbacks/state and cancellation with bounded deadlines, and closes only its own
processes. There is no synthetic callback substitution or audible-quality claim.
It has not been runtime-validated. The parent has received the proposal; **no
approval to execute it is implied or requested as part of this submission**.

## Failures, corrections, assistance and scope disclosure

All earlier artifacts, grades and safe logs remain; r1/r2/r3 were not overwritten.
Student-only freeze commits are listed in the manifests. These are development
iterations, not accepted tags or completed independent-review rounds.

- Initial compile/API assumptions, empty embedded-assets directory, wrong public
  type fields and GUI fixture provenance/configuration errors were retained.
- Core changes initially broke Chapter 6 local `PartID` expectations and removed a
  Chapter 5 tool-result `PartFinal` boundary. Both were restored compatibly;
  subsequent full core race/parity runs pass.
- Browser runs uncovered unbound native animation-frame calls (illegal invocation),
  signed-thinking readable projection, and terminal `interrupted` incorrectly
  routed as busy. Corrected and retained failed DOM/wire evidence.
- The first formal invocation wrongly set a global fail-closed `LLM_BASE_URL`,
  overriding the harness's vendor fixture URL, and omitted an explicit declared
  media row. Its failed Chapter 1/cumulative results are retained. Corrected
  invocation uses only the recorded minimal environment and inherited declarations.
- r1 corrected Chapter 7 still failed (WebSocket EOF); raw ping/fragmented probes
  independently demonstrated transport defects, fixed in r2. r2 still failed on
  extra start counts; exact-once projection defects were independently red/green
  tested and fixed in r3. r3 and r4 still fail the public count check.
- Parent's **public semantic assistance** clarified that cloned subscription
  envelopes may explicitly mark replay while preserving types/content/order/
  identities and leaving live unmarked. No private grader details were supplied.
  r4 implements that clarification; it does not claim that the marker explains
  or fixes the separate public failure.
- Some added test setups were invalid: wrong helper, missing provider provenance,
  stale/nonunique patch anchors, wrong raw-fixture filename. A newly introduced
  stuck unit test lacking provenance was stopped after verifying its owned PID;
  genuine long-lived tests were never killed. Invalid attempts are not counted
  as behavioral sensitivity. All corrected red/green receipts are named explicitly.
- The r2 long matrix overlapped later hub development: its GUI portion is not
  advertised as immutable-r2 whole-tree proof (`r2-matrix-attribution.txt`). Final
  r4 matrix ran entirely on its clean frozen bytes with no source edits.
- **Scope deviation:** early own preliminary logs were accidentally created under
  `reference/edition-3/ch07-student-evidence` instead of authorized
  `docs/edition-3-notes/ch07-student-evidence`. I promptly disclosed and moved only
  those newly created files, removed the newly created empty mistaken directories,
  and edited no pre-existing reference material. `scope-relocation-receipt.json`
  records reconstructed filenames/current hashes and the lack of contemporaneous
  before-move hashes; this is not concealed as perfect boundary compliance.
- Root approved the initial ownership/test plan and later bounded local corrections,
  supplied only public-scope/native-stop guidance and the replay-marker clarification,
  and directed final packaging after the unresolved r4 failure. No children,
  implementation from canonical snapshots, parent private audit, grader case bodies,
  credential reads, installation, push, publication, main source edits/staging/
  commits/tags, or live GUI disruption were used.

## Completion and next owner

All owned launched jobs/processes are finished. Source remains clean at the exact
r4 identity above, and the Chapter 6 tag is unchanged. Only named student evidence
and report paths in main were written, apart from the disclosed corrected initial
log-path deviation. Binary archives/profiles/screenshots remain in excluded owned
student artifact paths with pins rather than being staged in main.

The parent should first reconcile the preserved public pause failure against the
frozen candidate and raw controls, then independently assess full Chapter 7
requirements including remaining browser/native usability evidence. This report
neither accepts the candidate nor asks that the failing gate be waived.
