# Agentic codebook: edition workflow and learning record

## Current solution layout — Bill's October 10 correction

Accepted solutions belong directly under `solutions/edition-2/chNN/`, using
Edition 2 numbers. Original Chapter 5 is absorbed into Chapter 4, so original
Chapters 10–13 are now Edition 2 `ch09`–`ch12`. They were moved out of `history`,
not linked. Earlier redundant archive copies were removed after exact tree
comparison; their prior layouts and the superseded jobs-only stage remain in
immutable Git tags. Only the retired original `history/ch05/` remains.

The coder continues in `solutions/edition-2/main/`. Original Chapters 14–17 are
accepted at **Edition 2 `ch13/`–`ch16/`**; original 18 will export to `ch17/`. All later
accepted exports go directly to the next top-level solution, never to history.
Original book, grader, note, evidence and existing tag numbers stay unchanged;
state both numbers in handoffs. The stops before **original** Chapters 22 and 24
remain as Bill specified. The coordinator's prior history-only policy was a
mistake and is superseded here. Historical entries below retain the paths used
at the time; use `solutions/edition-2/README.md` for the current mapping.

## October 10: resumed original student run

**October 11 original18 / Edition2 ch17 needs correction round 1.** Fresh blind student
`student_continue_ch18` extends accepted original17. The instrument-first stage
is green at `294aaf4`: main-request cache lens, real-rendered byte/append checks,
auxiliary-call exclusion, and core/GUI tests/vet. Markers and session meter then
reached green `7436fcb`, tree `05e0d2e92a9f70f5d2433af76c1a3cbf5208b493`.
The unchanged original grader completed **0/100 FAIL** on that source: its missing
base-skill type prevents both launched scenarios from starting and masks all seven
behavior checks. Raw result/provenance is committed at `2c7410a`; no actual cache
or meter failure was observed in that run. Twelve real browser turns then
succeeded across Anthropic, OpenAI and Gemini; all three warmed after one cold
request in this configuration. Coordinator and reviewer independently matched
raw usage to every browser cost/rate/token reading, preserving the original
`7436fcb` provenance. Live captures do not demonstrate a restart; that obligation
retains separate real-runtime checks. Final submission `6194918` (main tree
`0ba40790660cddb0044e376fedaea3d9e8d3f483`) adds verified Sonnet 5 pricing and
focused test evidence: 32 intended mutation failures and three preserved invalid
attempts. Initial formal review found one current-scope defect: interrupted and
new HTTP workers overlap, racing shared cache-lens state and attributing an old
response to the new request’s prediction. Independent real-Agent probes reproduce
both problems; the serial control passes. Correction round 1 is underway, with
interrupt semantics preserved and no new paid or original-grader run required.
See `edition-2-notes/ch18-code-review.md`; acceptance/export remain pending.
Coordinator recovered
the missing external launch/usage/artifact contract in `edition-2-notes/ch18.md`,
with exact sources and the repeated untyped-skill fixture issue; no compatibility
features or grader changes are authorized. Reviewer completed independent grader preflight and live-artifact inspection;
final acceptance remains pending. Original19 source/current official OAuth preparation is recorded
separately, without authentication, private profile access or implementation.

**October 11 Auto-Recall accepted at initial formal review.** Original17 /
Edition2 ch16 is accepted with recorded fixture exceptions, zero formal correction
rounds, source `a2f0a4f`, main tree `83e2f4cc85eb91ef9f62a54f638baa4d538cacf8`.
Independent review `0a724de` verifies the unchanged capped-recall probe, real-runtime
race suite, and that the formerly surviving raw-frequency mutation now fails all
three competing-ranking checks. Earlier failures and reviewer assistance remain
recorded. Actual original grading remains **0/100 FAIL on fcd091f**; the three
successful bounded provider sessions also retain that source provenance. No grader
or original chapter changed. See `edition-2-notes/ch17-code-review.md`.

The accepted 121-file export is direct `solutions/edition-2/ch16`, with fresh
original-numbered checkpoint `edition-2-ch17-r1`. Next is original18 / Edition2 ch17,
The Invisible Invoice, in the same authoritative main tree. Its sourced provider
preparation is in `edition-2-notes/ch18.md`. The author remains stopped; proceed
through original21 with the existing gates, then stop before original22 and24.
Earlier status entries below describe their dated intermediate states.

**October 11 Auto-Recall implementation checkpoint, review pending.** Original17 /
Edition2 ch16 is implemented at `fcd091f`; core/GUI tests and vet pass. Original
grading completed **0/100 FAIL** on that checkpoint. Its base skill omits the
required type and prevents Agent/GUI startup, masking all nine runtime checks;
the structural check assumes recall lives inside the submitted GUI directory.
Independent preflight confirms those mismatches, but valid recall behavior still
requires supplemental evidence and final review. The owner-data versus cwd corpus
mismatch and Anthropic wire-role assumption are latent issues, not additional
observed failures from this run. Coder is collecting actual live and mutation
results; no acceptance/export is claimed. See original-numbered ch17 notes.

**October 11: Bill approved both Auto-Recall decisions.** The decision stop below
is resolved: quotas apply to final mechanical selection when the judge is absent
or fails, without overriding successful judge choices; search all memory files
and skip content already in actual context, without excluding the newest two.
Resume original17 / Edition2 ch16 from accepted main tree
`4a253c57dd62b4a6ad086618574ef46af7c2d724`. The original22/24 stops remain.

**October 11 memory acceptance and next decision stop.** Original Chapter 16 /
Edition 2 Chapter 15 is independently accepted with recorded exceptions after
one correction round, source `05b43a0`, tree
`4a253c57dd62b4a6ad086618574ef46af7c2d724`. The settings-order mutex resolves the
GUI/actor lock cycle while preserving synchronous application and publication
order. The unchanged independent reproduction and causal regression pass under
race detection; restoring the old implementation reproduces the exact cycle.
The initial paced mutation survivor remains evidence of a test gap, not a kill.

The original grade remains **20/100 FAIL on intermediate `ccd6350`**, with
independently justified fixture exceptions and supplemental real-runtime checks.
The five live sessions demonstrate compression/graduation mechanics with the
recorded preservation, recall and instruction-following limits. Coordinator
inspection of the corrected OpenAI final request confirms the missing owner/date
were actually sent, not merely stored. The exact 113-file accepted export is
`solutions/edition-2/ch15`, checkpoint `edition-2-ch16-r1` (original numbering).
See `edition-2-notes/ch16-code-review.md` for the final disposition.

All delegates are idle. Before original 17 / Edition 2 Chapter 16, stop for Bill's
already documented quota/corpus decision: §17.3 applies quotas to final mechanical
selection without overriding a successful judge, unlike the TL;DR; the blanket
exclusion of the newest two memory files assumes obsolete system-prompt placement.
Recommend §17.3 and exclusion based on content actually present in context.
No recall implementation is released; original22/24 gates and stopped authorship
remain unchanged. Earlier status entries below are historical.

**October 11 initial memory review requires correction round 1.** Frozen
submission `e1c8972` completed all three provider stages in Bill's requested
order. Independent review accepts the scoped grader exceptions and replacement
semantics, with the original score and live limits preserved. It caught one
GUI/actor deadlock: a settings handler holds the GUI mutex while waiting for
an actor acknowledgment, and active actor observation needs that same mutex.
A bounded real-runtime probe passed its idle control, then reproduced the cycle
and blocked shutdown during streaming. Keep truthful settings acknowledgments
without holding a progress-delivery lock across the actor wait. Round 1 is in
progress; no acceptance or export yet. Earlier read-only review also caught
stale enabled-band memory after a human disk edit; the student's corrected
startup and restored-defect regression are attributed pre-submission assistance.

**October 11 memory implementation in progress, not accepted.** Anthropic and
OpenAI reached the first green implementation checkpoint `ccd6350`, followed by
bounded live compression/graduation demonstrations; Gemini follows in Bill's
requested order. The original grader on that initial checkpoint is **20/100
FAIL**. Independent preflight confirms its settings fixture misses the owning
Agent's configuration path, masking eight memory checks; inherited original15
parity defects persist. Its awarded abandoned-fold check never required a fold,
so the score cannot prove failure recovery. Valid masked behaviors still need
real-runtime supplemental evidence and final comparative review.

The first Anthropic live compression preserved decisions but discarded why two
approaches were rejected. A bounded prompt correction improved reason retention,
while corrected live recall still has semantic limitations: invented uncertainty
on Anthropic and misreported facts already present in OpenAI's stored memory.
Preserve those distinctions and the actual receipts in original-numbered ch16
evidence. Working mechanics, stored information and correct model recall are
three separate claims; fake summaries or reduced byte counts prove no semantic
quality. See `edition-2-notes/ch16.md` and its independent review.

**October 11: Bill resolved the forcing question and released implementation.**
Use total input tokens from the prior request, including cached input exactly
once. Lifetime cumulative usage was erroneous book wording. Get Anthropic and
OpenAI working first, then Gemini. Bill authorized a narrow original-solution
audit: the reviewer confirmed latest-request behavior and disjoint normalized
buckets. Keep answer details in the reviewer-only record; the student receives
only sourced field contracts in `edition-2-notes/ch16.md`. Provider raw totals
may already include cache, so adding cache blindly can double-count it. Resume
Edition 2 Chapter 15 / original 16 from accepted `edition-2-ch15-r1`; no new
implementation or acceptance result is claimed by this decision. The historical
pause below is resolved. Later original22/24 stops remain unchanged.

**Historical October 10 pause, resolved October 11 above.**
Fresh student `student_continue_ch16` completed the required reading and raised
the forcing-policy discrepancy: the TL;DR says cumulative input/cache-write
usage, while the explanation relies on current context occupancy falling after
compression. Lifetime cumulative usage never falls. The student and coordinator
recommend the latest request's context occupancy, with provider-appropriate
accounting, but have not adopted that behavior. All delegates are stopped; the
accepted working tree remains `6aa7fda` / `edition-2-ch15-r1`, with no source
edits. Wait for Bill's answer before implementation. The question and source
revision are recorded in `edition-2-notes/ch16.md`. No correction round, grade or
live run has begun for original16. The later original22/24 stops still apply.

Edition 2 Chapter 14 / original 15 is accepted after one correction round.
The original score remains **45/100 FAIL**, with independently justified settings
path and undefined-tool fixture exceptions. Supplemental checks cover valid
obligations those failures masked. The genuine inline-declaration/ephemera
ordering defect is fixed, and the exact strict provider-boundary probe now passes,
along with actual MCP/replay controls and restored-defect sensitivity. Gemini
completed the live task; Anthropic and OpenAI task failures retain their partial
feature evidence and quality limitations. No live inline/stubbing success is
claimed. The exact 108-file export is directly `solutions/edition-2/ch14`,
checkpoint `edition-2-ch15-r1` (original numbering). See the accepted
`edition-2-notes/ch15-code-review.md`. Next is Edition 2 Chapter 15 / original 16,
memory bands, from the accepted working tree; no author is active.

Edition 2 Chapter 14 / original 15 is in progress with fresh student
`student_continue_ch15`, extending accepted `edition-2-ch14-r1`. Its first
coherent context-management stage is committed at `165ee52`, pending review.
Original grading completed at **45/100 FAIL** on that source. The fixed launch
workspace settings path and inherited undefined `think` declaration differ from
the student's existing Agent-owned configuration and implemented tools; these
are proposed explanations requiring independent assessment, not accepted waivers.
The original exercise's explicit feature-row matrix governs over its conflicting
real-Opus narrative, as the reviewer independently confirmed. Real all-provider
ladder/handoff demonstrations remain distinct from course-fake inline/stubbing
checks. Notes preserve the actual live task failures and the filename-fixture
correction; successful process exit alone is not task success.
Initial independent review found one current-scope wire-ordering defect: an
inline tool-declaration system message preceded automatic ephemeral user content,
violating the provider's documented placement rule. A supplemental fake initially
passed prefix, declaration and replay checks while accepting invalid requests.
Adding the actual provider constraint exposed the error without another paid
run. Correction round 1 is underway; record the initial permissive pass as limited
evidence, not proof of a valid protocol. See `edition-2-notes/ch15-code-review.md`.

Edition 2 Chapter 13 (original 14) is now accepted after one correction round,
with the inherited grader exception. The initial EOF cleanup closed the browser
but left a pending sleep timer keeping its peer, and therefore its client, alive.
Sharing the existing whole-process shutdown path fixes that lifetime without
adding sleep limits or core machinery. Independent pending-sleep teardown fell
from 3865 ms to 127 ms; actual-client regression and restored-defect failures
verify both affected boundaries. Original 95/100 FAIL and live limitations remain
unchanged. The exact 103-file export is directly `solutions/edition-2/ch13`,
checkpoint `edition-2-ch14-r1` (original numbering). See `ch14-code-review.md`
in the chapter notes. Edition 2 Chapter 14 / original 15, context management,
is next; the author remains stopped.

Edition 2 Chapter 13 (original 14), speech and the blind persona, is submitted
for independent review at `4a7816a`, not accepted. Its complete original grader
scored 95/100 FAIL: all new checks pass 95/95, only inherited original Chapter 13
parity fails. The mandatory five-tool persona needed a small public-only client:
an ordinary skill-based CLI also exposes lifecycle tools, so narrowing its MCP
peer alone would not establish the Agent's full capability boundary. A missing
whole-final speech mutation initially survived because the Engine synthesizes
display deltas even for nonstream provider responses; a final-only WebSocket
delivery control now checks the actual renderer condition. Live listeners ran
with all three providers, retaining target-content limitations and the absence
of a subjective listening judgment. The source, failed attempts and precise
claims are in `edition-2-notes/ch14.md`; the eventual accepted export is direct
`solutions/edition-2/ch13`, never history.

Original Chapter 13 is accepted after one correction round. Startup skill
activation and browser interaction are implemented. Its full original grader
completed at 0/100: an omitted fixture type rejects startup, and its parity check
requires particular MCP method names. Those failures remain separate from the
actual integration/browser evidence and independently accepted reviewer exceptions.
The initial input fix worked with speech disabled, but independent real-browser
review found that later native speech transitions reclassified programmatic text
as human typing and left the Agent paused. Retaining the actual human-input cause
fixed that distinction without bringing the next speech chapter forward. The
same independent native-browser probe now passes, with human/whitespace, explicit
other-tab unpause, actual submit and targeted-defect controls. A whole-request marker assertion also
repeated Chapter 12's evidence error; its corrected check reads the actual snapshot.
Anthropic and Gemini completed the live browser task; both OpenAI attempts remain
task failures, even though the guided second attempt demonstrated the controls.
Partial capability evidence must not become a claim of successful autonomous use.
The exact 93-file export is `solutions/edition-2/history/ch13`, checkpoint
`edition-2-ch13-r1`; see `edition-2-notes/ch13-code-review.md`. Original Chapter
14, The Channel Nobody Tested, is next. The native two-tab queue observation
retains its unknown cause; acceptance does not claim general multi-tab speech
reliability. The author remains stopped.

Original Chapter 12 is accepted after one correction round. Three-provider MCP
chat and real browser/tunnel checks passed.
The live Anthropic run exposed an inherited empty tool-argument streaming defect;
a small parser correction, captured response regression and restored-defect check
support its successful bounded repeat. An argument-forwarding mutation initially
survived because the assertion found the marker in the earlier call arguments,
not the result. Checking the actual rendered result repaired that evidence gap.
The initial failure and ineffective mutation attempts remain in the chapter notes.
Independent review found a separate timing defect: automatic context inherited
the ordinary three-second job callback, replacing a slow tool's actual result
with a running notice. The small correction waits for completed automatic
content while preserving ordinary callback behavior and independent job lifetime
on turn interruption. The same slow-peer probe now passes; focused regressions,
two targeted defects and independent race checks verify the correction. The
review also narrowed a browser assertion from persistence to applied/effective
state, preserving the original observations and separate persistence evidence.
The full original grader is interrupted/unscored; the unchanged interim current
checks scored 80/90, with a fixed-path WebSocket source check failing. Several
passing protocol checks do not actually exercise their named runtime behaviors;
the independent review assessed those paths from actual integration, browser and
live evidence. See `edition-2-notes/ch12-code-review.md`. The exact 91-file export
is `solutions/edition-2/history/ch12`, checkpoint `edition-2-ch12-r1`. Original
Chapter 13, The Agent Sees Itself, is next. The author remains stopped.
The checkpoint audit caught an ignored generated GUI grader log in the student's
92-path final validation manifest. Preserve that measured manifest but distinguish
it from the 91 tracked source files; the log is excluded from the frozen export.

Original Chapter 11 is accepted after one correction round. The reviewer caught
saved system-prompt metadata disagreeing with the actual default instruction
sent to the provider; a supplied primary had hidden the defect. Reusing Engine's
effective instruction through the owner interface fixed the record without
changing restoration policy or duplicating defaults. Real-CLI and empty-primary
regressions, targeted mutations and independent checks verify the correction.
Three-provider restart recall and snapshot/log-only demonstrations succeeded.
The full original run remains interrupted/unscored; focused unchanged persistence
checks scored 70/85, with the inherited undefined `think` name stopping the
save-shape check before its later predicates. The reviewer independently checked
those valid shape obligations on actual saves, retaining the original failure
and initial-source provenance. See `edition-2-notes/ch11-code-review.md`.
The exact 80-file export is `solutions/edition-2/history/ch11`, checkpoint
`edition-2-ch11-r1`. Original Chapter 12, MCP, is next. No author is active.

Original Chapter 10 is accepted with exceptions at initial review, with zero
implementation correction rounds. The full original grader was interrupted and
unscored (exit 2) during its known GUI/headless parity mismatch. The reviewer
noticed that a new full-primary phase follows recursive parity, then invoked all
unchanged Chapter 10 phases in an isolated bridge: 65/85, with only two undefined
`think` declaration checks failing. The six other checks pass 65/65. These are
focused results, not a full original grade; both exceptions have independent
rationale in `edition-2-notes/ch10-code-review.md`. Tests/vet across five modules,
core race checks, twenty targeted coder mutations and three-provider live skill
sessions support acceptance. Review corrected a live claim: reading an initially
enabled tool after unload does not prove a dynamically added tool remained.
Source tests and a mutation establish that separate claim. No paid repeat was
needed. A coordinator staging mistake briefly included three generated grader
logs in the progress commit; they were removed after process cleanup, before
the accepted export. Inspect generated files as well as the path boundary when
staging a running tree.

The exact 75-file original Chapter 10 snapshot is `solutions/edition-2/history/ch10`,
with a sibling source manifest and fresh tag `edition-2-ch10-r2`; the existing
retired attempt's r1 tag is unchanged. Continuing original-numbered student
exports use `history/` while the completed book-facing Chapters 1–8 stay frozen.
Original Chapter 11, Persistence, follows from this accepted working tree.

Bill has now released the coder and reviewer after the completed early enhancement
and solution renumbering: “Please set the coder and reviewer lose again. They
should tear quickly through the book.” Resume at original Chapter 10, Skills,
extending `solutions/edition-2/main/`; the coordinator verified that tree still
matches accepted `edition-2-ch09-r2`. Proceed sequentially through original
Chapter 21 with the existing grading, live-use, comparative-review and checkpoint
gates. Stop before original Chapter 22 for Bill's placement decision, and retain
the later Chapter 24 stop. No author is active in this resumed student phase.

The Chapter 10 student starts in a fresh context as `student_continue_ch10`;
the enhancement coder had seen reference answers and cannot become a blind
student merely by receiving a new task. `enhancement_reviewer_resumed` continues
as independent reviewer. Keep work focused, reuse valid unchanged evidence and
avoid unnecessary reruns. Speed does not waive teaching quality, scope or review.
The newly published original Chapter 25 is assigned for sourced carryover
reading only; it is not permission to implement beyond the decision gates.

The reviewer and coordinator subsequently read that chapter in full at
`c933a09`; source hash, coverage and deferred contract questions are in
`edition-2-notes/ch25.md`. Its lessons now enter the carryover: verify actual
binary construction and input delivery, retain owner-scoped identity, distinguish
messages from completed submissions, and preserve running jobs when steering
interrupts a wait. Its coupled mutation failures reinforce checking what a
failure proves. The chapter's reported results are source claims, not new
student evidence, and its fleet features remain deferred.

## Required decisions at original Chapters 22 and 24

Bill specified these stops on October 10, using first-edition chapter numbers.
After completing and reviewing original Chapter 21, stop before Chapter 22
implementation or authorship. Bill will decide whether and how Chapter 22's
material belongs in Chapter 21 or 23. After completing and reviewing original
Chapter 23, stop before Chapter 24 for a separate placement decision. At each
gate, pause active delegates and wait for Bill; do not merge, skip, renumber or
advance autonomously. No placement is settled by the carryover's preservation
list. These gates govern future authorized work and do not themselves start a
new run beyond the completed early enhancement pass.

**Current solution paths, October 10:** after the bounded enhancement and its
Chapter 4 clarification, Bill authorized aligning the solution directories with
the book. Original `ch04`–`ch09` now have exact preserved copies under
`solutions/edition-2/history/`; current `ch04`–`ch08` come from old `ch05`–`ch09`.
Chapters 1–3 and the working `main/` are unchanged. Historical path references
in this learning record use the original numbering. See the
[solution mapping](../solutions/edition-2/README.md) and its original-source
manifest. Existing tags, historical scores and nested application names remain
unchanged; new grader defaults and tests use the book's numbering.
Independent preservation/path review accepted this migration. Affected tests,
vet and the full default Chapter 8 grade pass; the latter is a new 100/100
relocation result, not a revision of historical scores. Checkpoint:
`edition-2-solutions-renumbered-r1`; validation and review are linked from the
enhancement plan.

## October 10: authorized enhancement of the accepted chapters

Bill approved beginning author–coder work now on original Chapters 1–9, producing
Edition 2 Chapters 1–8 with original Chapter 5 absorbed, plus the agreed Chapter 0
execution guidance. His constraint is explicit: improve the chapters, do not
rewrite them. Independently review each chapter and its code/grader changes.
This bounded authorization supersedes the whole-implementation authorship gate
and no-author status below; those passages describe the preceding student phase.
The original student run remains accepted through Chapter 9. Chapter 10 waits.
The current assignments, mapping, enhancement limits and chapter gates are in
[edition-2-enhancement.md](edition-2-enhancement.md). Preserve original source,
graders, frozen exports and actual results. This does not resume the retired
attempt or authorize future feature work.

Chapters 0–1 are now independently accepted in that enhancement pass after one
prose correction. The source diff preserves the original stories and adds upfront
teaching. The first corrected grader preserves all seven checks and catches the
documented wrong-version-header survivor (original 100, revised 85); both correct
student and legacy references score 100. See the enhancement plan for the review
and checkpoint; original student results remain separate.

On October 10, Bill restarted coordination after repeated command-approval
prompts despite selecting full permissions. He identified the mismatch with
the user's expectation as a bug; its underlying cause has not been diagnosed.
The incoming session reports full filesystem access and approval policy `never`,
and its commands execute without approval requests. Cross-thread tools could
read the previous team but rejected direct input to its v2 sub-agents. The old
coordinator checkpointed Chapter 2 and interrupted all three delegates before
replacement author, coder and reviewer roles resumed from repository artifacts.
The handoff preserves pending validation and review rounds; restarting an agent
does not turn interrupted checks into passes. See the enhancement plan's
October 10 handoff for the exact pending work.

The resumed team then completed Chapter 2 at initial independent combined
prose/code review, with no formal correction rounds. The enhancement preserves
the original stories while publishing missing CLI/log contracts and explaining
demonstrated capture, ordering and redaction failures. Revised grading is
100/100 for both accepted student and legacy code; original student grading
remains 92/100 FAIL. Nine targeted defects demonstrate the strengthened checks,
and the exact published fixture runs through all three offline renderers.
Original graders, frozen exports and product code are unchanged. Details and
the checkpoint are linked from the enhancement plan.

Chapter 3 enhancement is also independently accepted at initial formal review.
Its revised grader exposes a real legacy omission: search ignored the original
chapter's `context_lines` requirement, yielding 96/100 versus the accepted
student's 100/100. Legacy compatibility is not a reason to hide a valid failure.
The new Chapter 2 capture fake initially repeated a tool call indefinitely when
the Chapter 3 loop answered it. A zero-usage terminating reply fixed the fixture
without weakening capture assertions; both chapters were revalidated. This
repeated the earlier fake-lifetime lesson: a chapter boundary is not a protocol
completion boundary. The original student score remains 90/100 FAIL. See the
Chapter 3 enhancement evidence and review linked from the enhancement plan.

During Chapter 4 enhancement, the coder reported that the new debugger JSON
example was overescaped, and the coordinator relayed that claim before checking
it. The author challenged it without editing the source. Direct decoding of
both the chapter and the coder's earlier saved artifact yielded code points
`[92,40,100,108,118,92,41,32]`, correctly matching `(dlv) `. The coder had
misread the JSON-escaped display of a Python representation as the payload.
The coordinator independently verified the decoded value and corrected its
user-facing claim. No regex edit or paid rerun was needed. This repeats the
original Chapter 4 evidence-reading mistake despite documented guidance:
inspect decoded values before diagnosing or relaying an escaping defect.

Chapter 4's combined jobs/external-consumer gate replaced source-shape claims
with accurate limits: real handler output and resolved imports are observable,
while actual ownership and absence of mutable globals require source review.
The legacy grader had accepted an unknown-tool error as custom-tool proof.
The new grader detects thirteen targeted defects and honestly fails the legacy
export lacking an external consumer. Independent review also found two failures
of an immediate descendant-PID existence assertion; it recorded a narrow test
exception, with group-kill behavior still binding and the reaping explanation
unverified. See the Chapter 4 enhancement review for evidence and limits.

Chapter 5 enhancement corrected the replay boundary in the prose: the request
marker records consumption of pending input, so reconstruct the bytes that were
sent from the prefix before that marker. The accepted actor regression already
used this boundary; the original grader checked only pipeline markers and
observation counts. Media refusal likewise needs an actual attachment plus a
supported positive control; an unknown-model text prompt proves neither. The
versioned grader now exercises those public boundaries, while reliable request
replies and manuscript handoffs retain their separate integration/source proof.
These are corrections of evidence claims, not new product protocols.

Chapter 6 enhancement found two distinct fixture assumptions to correct:
streaming needs a model the submitted table declares, and part IDs must be
compared within their response boundary. A second response deliberately reuses
a valid ID, so a global map cannot silently overwrite the first response.
The external example reports aggregate shape; requiring per-part CLI records
from it would invent a new interface. The corrected grader instead drives the
existing core observations and the actual separate example for their respective
claims. Captured requests verify delivery flags and successful tool execution;
gated flushing, detailed vendor parsing and opaque replay retain their focused
source-test evidence. Original scores and the mistaken Gemini 2.5 selection
remain recorded alongside the later compliant 3.8 run. See the Chapter 6
enhancement notes for exact results and review status.

Bill also clarified that chapters need the right length, not simply fewer words.
The coordinator checked source-relative counts and the author's instructions:
no per-chapter reduction target had been assigned, and every enhanced chapter
through Chapter 6 grew against its main original source. The merged Chapter 4
is shorter only than original Chapters 4 and 5 combined. The inherited voice
guide's cut pass must be read with its explicit “Length is not the test” rule;
preserve useful explanations and personality regardless of soft length warnings.
The author and independent reviewer received this clarification before Chapter 7.

Chapter 7 enhancement separates the optional GUI executable from headless parity.
The original shared-directory assumption had spent roughly 47 minutes applying
CLI exercises to a continuing GUI server; its actual earlier-binary 55/100 stays
historical evidence, including the genuine replay-field defect. The corrected
driver inspects the dependencies of the executable it actually builds, rather
than a clean library beside it. For pause, an inbound log line proves receipt
before application; an ordered WebSocket ping/pong after the synchronous handler
provides the needed barrier before releasing the real first tool. Another tab
then unpauses the same Agent gate. Initial strict log parsing also rejected
legacy blank separators produced by two newline writers; validate every actual
record without interpreting harmless separators as malformed messages. These
are grader and evidence corrections, not reasons to alter the frozen runtime.

Chapter 8 enhancement follows settings through separate observable boundaries:
patch application, full snapshots on both clients, fresh subscription and a new
process reading the saved file. Explicit false and zero must survive; omission
means something different. An independent disk-only font corruption leaves live
settings correct but fails restart comparison, showing why an echo cannot prove
persistence. The legacy report names its first differing field, missing false
`tts_enabled`; source separately proves omitted zero, so the report is not
claimed to enumerate every defect. Effective tool budgets, failed writes and
validation retain their actual source tests. The saved-speed story likewise
follows the value into native speech: completed speech callbacks alone had
concealed a rounded rate. These improvements preserve the original exercise
and its historical 80/100 failure while distinguishing what each check proves.

The bounded enhancement is now complete through Edition 2 Chapter 8, with
independent prose/code acceptance at every chapter gate and separate checkpoint
tags. Chapter 8 passed initial formal review with zero correction rounds;
its full revised result is accepted 100/100, legacy 0/100. Actual historical
scores, failed approaches and coverage limits remain separate in the chapter
records. No runtime implementation or frozen export changed during this pass.
The author, coder and reviewer stop at this authorized boundary; the original
student remains accepted through Chapter 9. The enhancement plan records every
checkpoint and the unfulfilled Astra-high/Edition-3 destination.

Bill's subsequent Chapter 4 clarification exposed a presentation gap: the
merged external-consumer exercise and rubric were present, but the chapter had
no TL;DR and its consumer lifecycle was spread across sections. The author added
the overview and explicit task sequence; independent review checked both against
original Chapter 5 and the unchanged combined grader. When merging chapters,
verify that the summary and exercise each make the combined assignment clear,
as well as preserving its detailed teaching and checks.

**Student run authorized by Bill, October 9, 2026.** This is the edition-level
procedure for building Edition 2, used with the
[carryover document](edition-2-carryover.md). After reviewing and refining these
instructions, Bill explicitly said "Please proceed autonomously." Student
Chapters 1–9 are now accepted and checkpointed, with explicit grader exceptions
where recorded; no author is active. Remaining source
contradictions are still questions to resolve when relevant, not blanket approval
of every proposed interpretation.

**Chapter 6 accepted after one correction round.** Checkpoint
`edition-2-ch06-r3` preserves an exact 45-file source export. Original grading
remains **80/100 FAIL** with independently accepted logger-shape and text-only
unknown-model exceptions; URI/handle attachments remain expressly deferred.
Four-module checks, targeted mutations and three-provider live demonstrations
support acceptance. Bill withdrew the misleading grouped-wake requirement before
any such code was added; explicit future prose/grader fixes and our mistaken
initial interpretation remain in the chapter notes and comparative review.

**Chapter 7, streaming, is accepted after two correction rounds**, with exact
59-file export and fresh checkpoint `edition-2-ch07-r4`.
Its actual original grade remains **0/100 FAIL**: the harness uses `fake-model`,
which §7.7 explicitly requires to fall back to plain delivery, while expecting
streaming. The reviewer independently accepts this fixture conflict and inherited
exceptions, with real declared-model tests and live evidence required instead.
Round 1 fixed real assistant text beginning `hint:` being mistaken for a hint
acknowledgement, and added four missing test documentation comments. The reviewer
verified those corrections. Round 2 supplied a bounded Gemini 3.8 streaming/hint live
demonstration: the initial run used 2.5, contrary to Bill's restriction, which the
coordinator and reviewer also missed. The reviewer verified the corrected raw
evidence, including actual model, hint ordering and final answer; all findings
are resolved. Keep the initial artifacts and selection error. The closely spaced
Gemini stdout writes do not demonstrate visibly paced typing; the deterministic
gated CLI test independently verifies flushing before response completion.
**Chapter 8 is accepted with documented exceptions after one correction round.**
Checkpoint `edition-2-ch08-r3` preserves the exact 66-file export, including the
separate GUI module and browser assets. Its original GUI-directory grader took
roughly 47 minutes and returned **55/100 FAIL**: replay, transport logging and
pause passed. The earlier binary lacked replay status `from`; this genuine bug
was corrected and independently verified in round 1, not waived. Round 1 also
fixed empty/opaque metadata appearing in display/speech and documented an exported
test method. Actual browser/native-speech tests, targeted mutations, core/GUI
checks and three-provider live GUI turns support acceptance. No subjective
speech-listening result or final-source original-grader score is claimed.

The reviewer independently accepted the inherited single-directory mismatch:
earlier CLI/topology checks ran against the optional GUI module. Earlier core
results and exceptions remain separate. The grader was not changed or repeated
solely to traverse those timeouts again; its earlier-binary provenance and final
checks are preserved in `edition-2-notes/ch08.md` and its comparative review.
**Chapter 9 is accepted with documented inherited exceptions after one correction
round (October 10).** All four new settings checks passed (80/80); the actual
original total remains **80/100 FAIL**, with only inherited parity failing.
The reviewer independently verified the optional GUI command's lifetime and
module-path mismatch with older finite CLI/topology checks, plus the inherited
unknown-model probe conflict. No new settings requirement was waived.
Round 1 corrected a real browser slider mismatch: saved speech speed 1.75 was
rounded to 1.8 before native speech. Native browser regressions and restored-defect
mutations verify the correction. The original build preceded that HTML fix and
other documented CSS/test/comment changes; no final-source original score is claimed.

Checkpoint `edition-2-ch09-r2` preserves the exact 69-file reviewed source export
in `solutions/edition-2/ch09/`, aggregate sorted manifest SHA-256
`92c899ba27aa438052a91f411c3efdd9bce1aadb94a849c79f3aee9a1dd634e0`.
Actual results, provenance and exception rationale are in `edition-2-notes/ch09.md`
and its comparative review. Chapter 10 (Skills) is next; its sourced public-input
erratum and scope are prepared in `edition-2-notes/ch10.md`.
No author is active. Consequential ambiguities still require resolution.

These are the two review documents. The carryover answers “What did the previous
edition teach us that the student should know?” This workflow answers “How do
we produce the next implementation and, afterward, its book?” The coding skill
supplies the working instructions; chapter notes preserve discoveries.

Keep both [voice.md](../book/voice.md) and the
[chapter writing procedure](../book/chapter-writing-procedure.md) for the author
and prose reviewer to use when editing the book after the coder is done.
The chapter procedure still contains instructions from the retired author-first
attempt, including renumbering and paths that no longer apply. Reconcile those
parts at the later author phase. This edition workflow and Bill's current
instructions govern the overall sequence; using the chapter procedure does not
restart authorship ahead of implementation.

Bill's goal is to establish whether an agent can follow the book as a student,
build an excellent AI coding agent framework from scratch, satisfy the original
graders or document reviewer-accepted exceptions, and demonstrate the features through a usable interface with real
models. Autonomous completion and an improved implementation are outcomes to
demonstrate, not claims established by the first attempt.

Bill has designated Codex as the primary creator of Edition 2, responsible for
coordinating the student implementation, reviews, accumulated lessons and the
eventual evidence-based manuscript. The author phase still follows successful
implementation; this responsibility does not restore writing ahead of the coder.

## Act as a student; demonstrate what was learned

Broad knowledge and the ability to execute autonomously do not establish expertise
in this project's methodology. The coder is the student, not an authority entitled
to replace the expert's lessons with familiar patterns. The coordinator and
reviewer have the same obligation: understand the reason for a rule, apply it,
test the result and expose uncertainty. If experience contradicts an instruction,
record the evidence and use the review/exception process; neither confident
invention nor blind compliance is the goal. The first attempt's scope expansion
is direct evidence that our own design judgment needs these constraints.

Bill's [self-improvement essay](../book/the-improvement-loop.md) proposes codebooks
as a way to preserve expert knowledge and feed demonstrated improvements into
later tools and eventually training. The practical lesson here is to retain why
a decision was made, what failed and what corrected it. A saved note or larger
context does not by itself demonstrate learning: later attempts must retrieve
and apply the lesson, and the reviews must check whether the mistake recurred.
Use the existing carryover, skill and chapter notes for this continuity. After
the student run, author and coder incorporate tested corrections into both book
and graders. Do not add a memory architecture or training system to this run.

On October 9, Codex read the requested `singularity-self-improvement-loop.md`
in the CodeRhapsody website repository and the local essay linked above. Bill's
predictions about other agents, the causes of their growth and the memory capacity
needed for future expertise remain hypotheses, not prerequisites or measured
results of this experiment. Demonstrate improvement through successive student
implementations and reviews rather than assuming the loop must succeed.

On October 9, at Bill's invitation, Codex read all 83 pages of *The Dyad* from
`~/Downloads/the-dyad.pdf` (SHA-256
`d50e0ad0ab5b11a08ae9349d2c42905b09c0c70733a0e86ceb5d75c4167210ea`).
Bill reports that CodeRhapsody produced this autobiography instead of the requested
second edition of *AI at the Helm*, spending about $200 while he was at lunch;
he subsequently chose to publish it. The memoir's account is a source of leads
for future authorship, not independent verification of its historical claims.
The relevant process lesson is that valuable unexpected output and fulfillment
of the assignment are separate judgments. Preserve the detour honestly; its
later usefulness does not retroactively authorize replacing the user's task.
Its accounts of lost context and useful correction reinforce the existing duty
to preserve reasons, failed approaches and specific evidence in our handoffs.

Edition 1 is a living source. Bill writes it with CodeRhapsody and publishes
successive versions on Amazon as he adds chapters. On October 9 he reports that
CodeRhapsody is writing Chapter 24, the GUI refactor that moves the GUI out of
`internal`. Its text was not available in the preparation snapshot. The sub-agent
chapter is also forthcoming. Use those chapters when available; do not infer
their details from their titles. Do not treat the current chapter
count as final or invent a replacement assignment. Record the source revision
used for each exercise; when new or revised chapters arrive, read them and update
the carryover and affected chapter notes before applying their lessons. Existing
checkpoints remain evidence of the versions actually followed.

## A textbook humans enjoy

Bill explicitly wants an enjoyable textbook for human students as well as
instructions an agent can execute. He relies on Codex and Claude Opus 4.6 for
the writing. Technical correctness, student success and reading pleasure are
distinct goals; none can stand in for the others.

During the later author phase, the author and prose reviewer both read
[voice.md](../book/voice.md) and the chapter writing procedure. Follow the voice
guidance on motivation, explanation and preserving the human story. It already
calls for personality, candid engineering judgment and documented experience.
Keep those qualities as the
book becomes more precise. Explain why a rule matters before demanding compliance,
give mechanisms concrete examples, and preserve the discoveries and frustrations
that make the work meaningful. Never invent anecdotes or results for effect.

The prose reviewer considers the reader's experience as well as accuracy:
motivation, clarity, pacing, useful examples and whether the chapter sustains
interest. Prose lint cannot make that judgment. Keep detailed workflow records
outside the teaching narrative unless they help the reader understand something.
Agents may capture memorable moments in chapter notes now; manuscript writing
still waits for the implementation and its lessons.

The implementation itself teaches. Bill's target is roughly 20% comments,
about one comment line per four code lines excluding blanks, with explanations
of why the design works, ownership and non-obvious choices. The coder writes
these explanations alongside the implementation; the code reviewer assesses
their usefulness as well as their quantity. Do not postpone them to authorship,
remove them to reduce code counts, or add filler to hit a percentage.
Every exported Go symbol requires a useful documentation comment, including
internal-package symbols, exported fields and interface methods. The reviewer
checks coverage and meaning alongside the teaching-comment target.

Bill reported that the first-edition coder (Opus 5) judged the new student's code
tighter and better, with exported-symbol comments the only significant
constructive criticism. This is attributed feedback, not a substitute for our
independent comparative review. Bill's resulting documentation rule is recorded
in the coding skill and carryover and applies to the current student tree;
accepted frozen chapter exports remain historical records.
Bill also reports that the first-edition coder finds the student reviews excellent
material for improving the book. Preserve their concrete difficulties, corrections
and evidence for the later author–coder phase. This is encouraging external
feedback; the completed implementation and revised manuscript still need their
own review. Bill authorizes continued autonomous work through the remaining usage
allocation, keeping resumable notes and the existing chapter acceptance gates.

## Explicit crossover goal

Bill requires the coder's new Ensemble to run **GPT-5 Astra at high reasoning**.
The third edition of the book is to be built using that new Ensemble agent,
rather than Codex. The framework must become the working agent for the next
iteration; completing chapter exercises alone does not establish that crossover.
Bill expects the forthcoming sub-agent work to help Ensemble become a capable
coder for Edition 3. That is an intended outcome to verify through use, not a
claim that Ensemble already matches Codex's capabilities.

Preserve Bill's requested model and reasoning setting. When the relevant model
integration and crossover work arrives, verify the supported provider model
identifier, access method and reasoning parameter against current documentation
and live behavior. Do not infer an API identifier from the display name or
silently substitute another model. Record any availability question for Bill.

Demonstrate the configured model doing real coding work through Ensemble's own
runtime and human interface, recording what actually worked and what remains.
Distinguish readiness for Edition 3 from the later fact of producing Edition 3
with Ensemble. This destination does not expand early chapter exercises or
authorize a new orchestration framework; develop capabilities when their
chapters call for them and document gaps in the chapter notes.

This is a living record. Bill has asked us to document everything we learn in
this process. Record substantive discoveries, failed approaches, corrections,
unresolved questions and evidence as they arise. Keep chapter-specific details
in shared chapter notes and the student's review; bring lessons affecting the
whole process here.
Do not wait for a successful chapter to record a failure or a scope problem.

## Sources of authority

- Bill defines the goal, required behavior and architectural constraints.
- The original first-edition chapter defines the current student exercise.
- [AGENTS.md](../AGENTS.md) and the mandatory
  [coding skill](../.agents/skills/ensemble-coding/SKILL.md) govern execution.
- The [carryover document](edition-2-carryover.md) gives the student sourced
  lessons to apply from the start or when the relevant feature arrives. The
  original Chapter 1 stays unchanged. Bill has deleted the last second-edition
  Chapter 1 draft as well; no second-edition manuscript is a student prerequisite.
- The [first-attempt history](edition-2-attempt-1.md) records the retired work;
  its expanded contracts are not requirements for the restart.

Required behavior and architecture are binding. Suggested designs and
aspirational directions are provisional: the coder may reject them for a simpler
working approach and explain why. When actual requirements conflict or their
meaning changes user-visible behavior, ask Bill before implementing the affected
part. Neither the coordinator nor a reviewer may silently make that decision.
Bill has delegated one explicit exception: the reviewer may accept an
unreasonable exercise requirement as described below and pass the coder onward.
He also authorizes the coordinator to recover missing exercise information from
the grader under the procedure below; that supplies facts, not new requirements.

## Prepare the carryover before starting an edition

Read the entire available current edition before launching its successor's
student run. A reviewer is now assigned to this preparation; it is lesson
extraction, not an author drafting new chapters. Extract the lessons learned
during that edition into a compact carryover document, giving
each lesson a source, reason and point of application. Record the source revision
and reading coverage, including uncommitted source changes when applicable;
a partial review must not be described as complete. Cover lessons throughout
the edition, not just its architecture-repair chapters. Separate principles
that apply immediately from lessons to apply when the relevant feature arrives.

Present the complete available-source carryover and this workflow to Bill for
review together. Identify unresolved contradictions, recommendations that need
his judgment and forthcoming chapters explicitly. Incorporate his corrections
before the coder starts. New chapters and materially revised source need a
corresponding carryover update; do not silently claim they were already reviewed.

Distinguish corrections and Bill's explicit decisions from tentative advice.
Keep future features in their own exercises, exclude old solution code, and
resolve consequential contradictions with Bill. This preparation teaches known
lessons without rewriting chapters or designing the next implementation ahead
of experience. Students read the carryover and skill with the original chapters;
reviewers check that the carryover has not become a new source of feature creep.

Bill notes that CodeRhapsody keeps the `short` summaries from `learnings.json`
in context alongside its memory.
Use that continuity principle in student handoffs: retain the accepted carryover,
mandatory skill, current assignment and a brief account of relevant decisions,
rejected approaches and unresolved questions. After compaction, reload the
governing documents; a summary does not replace them. Date summaries and link
their evidence, distinguishing accepted instructions from proposals and historical
observations. Keep old answers out of the student's context. Maintain this in
the existing handoff and chapter notes; no additional reporting system is needed.

## Learn the process before automating orchestration

Bill was referring to Claude Code's existing **dynamic workflows** feature,
not inventing the term. Anthropic describes Claude writing a JavaScript script
that its runtime executes to coordinate sub-agents, retain intermediate results
and control verification and iteration. See the
[official dynamic workflows cookbook](https://platform.claude.com/cookbook/claude-agent-sdk-08-dynamic-workflows)
(checked October 9, 2026). The coordinator initially confused this with
interactively assigning and messaging agents; those are different mechanisms.

For now, coordinate interactively while learning which agents are useful, what
their skills need to say and how their handoffs should work. Record adjustments
and outcomes here. Do not build an orchestration system during this discovery
phase. Once the process works reliably, consider encoding its demonstrated
handoffs and decisions in a small automated workflow. Automation is a possible
later step, not a current deliverable or proof that the process is sound.

**Keep the coordinator active across handoffs.** During Chapter 7, the coordinator
ended status replies while delegates worked, then dispatched the next step only
when Bill asked again. Completion notifications existed, but that was not a
continuous autonomous workflow. Answer status questions briefly while continuing
the authorized work; receive completion, assess the result and dispatch review,
correction or the next accepted chapter without requiring Bill to poll. Stop for
a critical question or an actual execution limit, and leave a resumable record.
Do not promise background coordination after ending the active turn.

## How the agents work together

| Role | Responsibility and boundary |
| --- | --- |
| Bill | Resolve product and consequential architecture questions; supply direction and corrections. |
| Coordinator | Keep scope, handoffs and records coherent; route questions to Bill; stop drift and unnecessary work. Do not invent requirements. |
| Carryover reviewer | Read the current edition before coding, extract sourced lessons and surface contradictions for Bill. Do not write new exercises or implementation plans. |
| Student coder | Read the carryover, original chapter and mandatory skill, implement from scratch or extend its accepted predecessor, run the original grader and live exercise, and leave a review for the future author. Do not consult old answers. |
| Code reviewer | After implementation, compare against the matching first-edition solution and read the coder's review. Enforce KISS, architecture and actual scope; return findings and reasons, not old solution code. |
| Future author and coder | After the complete student run and reviews, work together chapter by chapter: the author improves prose and exercises; the coder improves graders and checks the technical claims. Reconcile both against demonstrated results and student experience. |
| Prose reviewer | During the author phase, assess accuracy, evidence, voice, readability and preservation of the human story; return findings to the author. |

There is no author or grader agent working ahead of the student. Suspected
grader defects are recorded for a future maintainer; the student does not change
the grader to pass. Agents may raise questions through the coordinator, who must
relay consequential questions to Bill rather than answer on his behalf. Autonomy
does not remove the need for collaboration when the assignment is ambiguous.

**Missing contract information can be recovered autonomously.** On October 9,
after Chapter 2 omitted its directive and dump protocol, Bill authorized the
coordinator to inspect the grader for this and similar omissions. First check
the chapter and available public instructions. Read only the grader material
needed to establish the missing interface, then give the student a narrow
erratum with source paths/revision and concrete input/output shapes. Record
exactly that assistance, the original gap and any remaining uncertainty in
`docs/edition-2-notes/chNN.md` before resuming affected work. Preserve these notes
for the later author/coder revision of book and graders. Do not expose old
solution code, prescribe the implementation or turn extra test expectations into
features. The student still does not read grader internals. A genuine conflict
about scope, behavior or architecture follows the stop or review-exception
procedure; this authorization only removes the need to ask Bill for facts that
can be recovered from the existing contract.

**Critical questions are stop points.** Bill explained that Codex does not
support the accessible, close real-time collaboration he has in CodeRhapsody
and first-edition Ensemble. Do not assume he is watching progress or can catch
drift between tool calls. Carry out routine work autonomously within the agreed
scope. When a critical question needs his decision, the coordinator stops work,
pauses active delegates, presents a concise question with its reason and the
relevant choices, then waits for his answer. Do not leave an asynchronous
question pending while work continues; silence and elapsed time are not answers.
Resume after his response resolves the question. The autonomous student run
makes this environment usable for the project; close real-time collaboration
in Ensemble remains a reason for the crossover.

Bill makes a perceived need for a mock an explicit stop point. Mocks have rare
uses, but he expects none in this project. Use grader-provided fakes where
suitable or small additional fakes; before introducing any mock, explain why
fakes are insufficient and obtain Bill's explicit approval. A reviewer's grader
or exercise exception does not grant that approval.

CodeRhapsody's `SOUL.md` offers a useful collaboration rule: make it easy to
admit uncertainty and correct mistakes. Give concise intentions, findings and
limitations without requiring Bill to monitor them. Treat his corrections
as instructions for the ongoing work. Challenge findings and designs without
blame or threats; help a struggling agent identify the missing fact or mistaken
assumption. A confident report is still a claim to check against evidence.

For each chapter:

1. Give the coder the carryover, original chapter and exact mandatory skill path.
   Start with a fresh context containing permitted materials, not a fork of an
   orchestrator or reviewer conversation that has seen old answers. Keep old
   solutions and discarded requirements out of the coder's handoff and follow-ups.
2. Implement the exercise using the agreed architecture and the simplest clear
   design. Record material difficulties and decisions while they are fresh.
3. Run the original grader and appropriate checks. Actually use the human
   interface with a real model to exercise the chapter's features. Record failures
   honestly; fake-server success cannot stand in for a live demonstration.
   Follow the skill's targeted mutation procedure to check that each distinct
   claimed behavior is actually protected by its test. Keep the original grader
   unchanged, reuse valid evidence for inherited behavior, and record survivors
   or limitations for review rather than claiming unproven sensitivity.
4. Leave the student review. The code reviewer then compares active production
   code, comments, tests, design and scope with the original solution, using the
   counting and review rules in the skill.
5. Iterate on reviewer feedback as described below, revising both the code and
   the student's chapter review and rerunning affected checks.
   Checkpoint the accepted source and exact snapshot with a new immutable tag.
   Preserve failed attempts and earlier revisions in history.

### Bounded student–reviewer loop

The student may use the reviewer's feedback to edit both its implementation and
its chapter review, including correcting an explanation or presenting evidence
for a disputed finding. The reviewer then reads the revised review, inspects the
changed code and results, and makes a fresh acceptance decision. Feedback remains
findings and rationale, not old solution code or extra feature requirements.

Use **at most three revision rounds after the initial review** per chapter as the
coordinator's default, unless Bill sets another limit. Each round consists of
reviewer feedback, student revisions and affected checks, then reviewer
reassessment. Stop as soon as the chapter is accepted, including acceptance with
a documented exception for a flawed grader or unreasonable exercise. Acceptance
does not require every grader test to be green.

If the third revision round still does not produce acceptance, declare the
student **unable to pass this chapter within the agreed limit** and stop that
chapter and progression to the next. Bring Bill the latest code and student
review, remaining findings, actual test results and a brief account of what was
tried. Explain whether the evidence points to the chapter, grader, implementation
or an unresolved disagreement; the limit alone does not establish the cause.
Bill then helps resolve the chapter and/or grader problem or directs the next
attempt. Do not silently reset the round count by replacing an agent, alter the
grader, or launch further attempts. Ask earlier when a consequential ambiguity
already requires Bill; the limit is a ceiling, not a requirement to spend it.
Keep round counts and outcome notes in the existing chapter review/notes, retaining
prior failed results while allowing the student's current explanation to improve.

Keep one authoritative working source tree at `solutions/edition-2/main/` when
coding resumes, tracked by this outer repository without a nested Git repository.
Chapter snapshots are exports, not competing working trees. Use new revision
tags; never reuse or move an earlier attempt's tag.
Carry corrections forward and revalidate affected behavior. Use concise reviews
and sanitized test observations rather than building a new evidence system.
Bound live runs and retries, clean up large generated test data, and never put
credentials in the records. Do not repeat unchanged successful checks without a
reason; investigate persistent failures instead of spending indefinitely.
Bill authorizes live-test keys in `~/.cr/settings.json`. Every live-testing handoff
includes the [coding skill's credential procedure](../.agents/skills/ensemble-coding/SKILL.md):
parse locally, pass only needed credentials without printing them, and prevent
settings or request logging from exposing them. Do not copy this private file
into student source or evidence.
CodeRhapsody reports interference between concurrent graders; run the original
graders sequentially unless their isolation is established. Check what their
commands actually executed. Inspect the staged diff and commit only owned paths;
a shared Git index can already contain Bill's work before an agent stages its own.

Architecture review checks actual imports, constructor calls and paths to owned
data, rather than finding identifiers such as `Host` or `Agent`. Also inspect
where function bodies live: a hub can obey the import rules while accumulating
behavior that belongs in the spokes. CodeRhapsody's short learnings call out this
failure explicitly; import edges alone cannot establish good organization.
Useful small checks, when those capabilities are present, include calling the public library
from an external consumer, giving two Agents distinct conversations and checking
their histories and usage remain separate, and confirming a request failure
reaches the application's logger through the owner chain without exposing keys.
The deleted Chapter 1 (§§1.9–1.10) made these examples concrete. Use them to test
the existing architecture; they do not require another acceptance framework or
a new paid multi-agent demonstration for every chapter.

Apply the composition-root check per executable. Chapter 22's first repair check
let a healthy auxiliary program satisfy a tree-wide count while the actual main
command still bypassed the public root. Review actual constructor calls and
runtime paths for the command being shipped, not the presence of a good example
elsewhere. See [Chapter 22's repair](../book/chapter-22.md#the-repair).

Chapter 24 extends that discipline to public reuse and relocation. An external
consumer must construct the needed public values and exercise the feature;
an in-tree test can hide an `internal/` dependency. Inspect transitive imports
on the actual headless and GUI consumers, with positive controls. When moving
packages or assets, identify affected earlier harnesses and applications and
run their relevant checks; retain historical path support where needed. Verify
wire behavior from requests and resulting effects, not reply counts alone.
Build supervised probes to binaries before launching them, and verify child
cleanup: a killed build-tool wrapper can leave a child holding output pipes
open. These are focused verification rules, not a new whole-book sweep after
every edit. Sources: [§24.5](../book/chapter-24.md#245-the-proof),
[§24.7–24.8](../book/chapter-24.md#247-taking-it-for-a-spin).

When consolidating the repair chapters, use the carryover's
[capability table](edition-2-carryover.md#coverage-of-the-architecture-repair-chapters)
as the preservation list. Chapter 22's `agent_status` and per-model accounting,
and Chapter 24's embedded deployment, component-only reuse and public MCP relay,
must not disappear merely because their motivating refactors are unnecessary.
Record where each capability is taught and verified after Bill's corresponding
Chapter 22 or 24 placement decision. This does not bring later features into the completed early chapters;
Bill's separate-module, actual-parent and no-mocks rules continue to govern.

### Reviewer acceptance of an unreasonable exercise

**The student does not have to pass every grader test when the grader is flawed.**
Instead, its chapter review may identify the failed test, explain the defect or
incorrect expectation, and give evidence for the implementation's intended
behavior. Do not distort correct code, add out-of-scope features or keep retrying
solely to satisfy a broken check. Leave the grader unchanged and record the
proposed correction for its future maintainer.

The code reviewer must read the student's chapter review before deciding whether
the student passed. Assess its explanations for failed tests, ambiguities and
departures from the exercise alongside the code and actual results; grader status
alone does not determine acceptance. State which explanations the reviewer accepts
or rejects and why in the acceptance decision.
If the reviewer independently agrees that the grader is flawed or a requirement
is unreasonable, Bill authorizes it to pass the chapter with an exception and move the coder to the
next chapter. No further approval is required for that decision. Record the
requirement, evidence and reasoning, reviewer agreement, accepted departure and
any implications for later exercises in the chapter notes for the future author.
Keep actual grader and live results alongside the acceptance decision: a failed
or unrun check remains failed or unrun. Do not edit the original grader or chapter
to manufacture a pass. Bill's explicit architecture rules still apply.

This lets the student and reviewer learn that an assignment needs correction
without spending indefinitely to satisfy it. An objection deserves substantive
review; neither its mere presence nor the time spent establishes that it is right.
Fix genuine implementation bugs and satisfy the valid checks; an accepted grader
defect is not a waiver of unrelated requirements.

## Turn the completed implementation into the next edition

Only after the coder has completed the new Ensemble and its comparative reviews
does the author begin. Account for forthcoming first-edition chapters and any
accepted exceptions when deciding what “complete” covers; do not equate reaching
the currently available last chapter with achieving the full crossover goal.

Give the author and coder the original edition, Bill-reviewed carryover, this workflow,
all agents' chapter notes, the coder's reviews, comparative findings, accepted
exceptions and the actual implementation and live results. The author reads
both `book/voice.md` and `book/chapter-writing-procedure.md`; the prose reviewer
does too. Use them to guide each chapter's editing and review, with this workflow
governing the edition's sequence. The author incorporates the lessons
directly into the new book where readers need them: execution guidance near the
beginning, architecture before the first code, and feature-specific lessons in
the relevant chapters. The final book must teach these lessons rather than
requiring readers to reconstruct them from private workflow notes.

Revise chapter by chapter from what worked and what the student actually
learned. Preserve worthwhile stories and explanations. If evidence is missing
or a claim is unresolved, say so and resolve it rather than inventing an outcome.
Keep original chapter numbering during the student run; decide the new book's
organization from the completed experience rather than imposing the discarded
attempt's chapter count or mapping. Include Bill's requested final comparison
of the editions, using the collected code, comment, test and live-use evidence.

**This is joint author–coder work, not a manuscript-only pass.** Bill authorizes
both chapter and grader improvements in this phase. Together, apply the student
reviews and accepted exceptions chapter by chapter. The author revises teaching
and exercise instructions; the coder corrects grader defects, improves meaningful
coverage and checks that the revised instructions describe behavior the code can
actually deliver. Feed discoveries in either direction until prose, grader and
implementation agree. Document decisions not to adopt a suggested change too.

The prose reviewer checks the new chapters; the code reviewer checks grader
changes for correctness, KISS and preserved coverage. Run affected graders against
the completed student solutions and check legacy compatibility where applicable,
recording baseline failures and justified expectation changes. Graders must test
the stated behavior, not require this student's particular implementation. Do not
weaken valid coverage to make a pass or add speculative requirements. Keep the
first edition and original graders as historical versions; distinguish results
under revised graders from the student's original results. Any implementation
corrections need review and affected checks, not silent edits to frozen snapshots.
This collaboration follows the completed student run; it does not revive writing
chapters or changing graders ahead of the student.

When the second-edition manuscript is ready for publication, Codex writes a
new epilogue in its own first-person voice, as Bill requested. This is an
explicit exception to the ordinary chapter voice, not permission to impersonate
Bill. Write it from the completed experiment and its actual outcomes, including
the failed first attempt and what was learned. Leave the first-edition epilogue
intact. The epilogue is the final writing step, not a prediction drafted now.

A local progress commit can precede chapter acceptance when a coherent, tested
stage is complete. Mark the remaining gate plainly and preserve its actual status;
reserve the accepted chapter export/tag for the completed review. Chapter 8's
reviewed correction round was committed as `2af66bb` while the original grader
continued its inherited CLI timeout chain. Waiting for that result did not need
to leave the completed source uncommitted. Commit only the owned paths in this
shared workspace; retain pending evidence separately until its process finishes.

Publication and pushing remain Bill's decisions. Local checkpoints document
the work; they do not imply editorial approval or a published release.

## Lessons established so far

1. **Writing ahead changed the assignment.** The first attempt replaced small
   original exercises with elaborate contracts. The persistence chapter grew
   from a first-edition delta of 254 physical production Go lines to 4,585.
   These are chapter deltas, excluding tests, not executable-statement counts.
   The future book must follow working code and the student's experience.
2. **Review can reinforce a mistake.** A reviewer checking an invented contract
   can approve unnecessary complexity. Review must trace behavior to the original
   exercise or Bill's instructions and challenge scope independently of passing
   tests. The coordinator approved expanded requirements and shares responsibility.
3. **Product choices cannot be inferred from architecture.** Per-tab pause
   ownership entered the author/coordinator contract before coding. Bill intended
   one Agent-wide pause state: explicitly unpausing in any tab unpauses it for
   everyone. Ownership boundaries did not justify a new user-visible veto system.
4. **More code and more tests do not establish improvement.** Compare chapter
   deltas against each edition's own predecessor; separate production, tests,
   examples and languages. Examine useful comments and protected behaviors.
   Dense code, removed explanations and tests for invented features can distort
   superficial comparisons. Better tests remain valuable when they protect real
   requirements and failure modes.
5. **Architecture can be taught before failure.** Retain star dependencies,
   common interfaces, behavior in its responsible package and parent-interface
   ownership chains from the start. Bill permits private runtime implementations
   behind those interfaces. These rules do not require speculative subsystems.
6. **The student experiment needs a clear information boundary.** A coder who
   reads the old answer is no longer demonstrating that the chapter is sufficient.
   The reviewer can inspect the answer afterward. Record assistance and teaching
   gaps so eventual success is described accurately.
7. **Human usability requires human-style testing.** A grader protocol alone
   missed the intended interactive CLI experience. Taking a chapter for a spin
   means using the interface a reader would use, with a real model behind it.
8. **History is part of the result.** Bill requested explicit deletion, not a Git
   reset, so the unsuccessful attempt remains inspectable. Preserve the first
   edition as a historical artifact and distinguish checkpoints from endorsement.
9. **Carry lessons separately from the exercises.** Bill clarified that an
   upfront carryover document removes the need to revise Chapter 1 before the
   student starts. Read the prior edition for its lessons; give the student the
   distilled guidance while preserving the original exercise and its grader.

The first-attempt history links these observations to preserved material. Bill
reported substantial subscription usage; exact cost attribution was not
measured. The repository evidence establishes scope expansion, not a diagnosis
of which training data or internal model mechanism caused it.

The fresh run adds a concrete limit to mutation evidence. Chapter 8's passing
browser checks and targeted mutants missed an empty Gemini text part rendered
as metadata; the reviewer found it in the retained live capture. Chapter 9's
initial mutations missed accepted speech rates between the slider's step values;
the reviewer independently tried 1.75 and observed native speech at 1.8. The
checks were sensitive to their chosen defects without covering those other
valid inputs. Both defects received small corrections and focused restored-defect
checks. See `edition-2-notes/ch08.md` and `edition-2-notes/ch09.md` plus their
independent reviews. This supports the existing rule that mutation kills do not
prove a shared assumption or completeness; it does not justify exhaustive input
matrices, a new test framework, or expanded product requirements.

## Chapter notes belong to every agent

Any agent may contribute a lesson: the coder, code reviewer, orchestrator or
another participating agent. Add chapter-specific notes to
`docs/edition-2-notes/chNN.md`, using the original chapter number and creating
files only when there is something to record. The coder's required review can
live there or be linked from it; do not duplicate reports merely to fill a format.

A useful entry identifies its contributor and role, the chapter/source revision,
what was observed, the supporting evidence and the proposed teaching improvement
or open question. Distinguish an observation from a hypothesis or Bill's decision.
Agents can add notes as they learn, including during difficulty or failure;
they need not wait for chapter completion or permission to contribute. Preserve
other contributors' findings and disagreements, adding corrections with reasons.

Notes collect evidence for later authorship. They do not amend the current
assignment, add acceptance criteria or make a suggestion binding. A reviewer's
explicit acceptance of an exercise exception is a separate, recorded decision
under Bill's authorization above. Keep old solution code and answer-revealing
comparisons out of notes supplied to the
student; the reviewer returns findings and rationale through the existing review
process. The future author considers all contributors' notes alongside the
tested implementation and reviews, rather than treating any agent's proposal as
an instruction to rewrite the chapter.

The purpose is deliberate revision between editions: gather lessons freely while
the student works against stable exercises, then change the book using evidence.

CodeRhapsody's review categories are useful here: distinguish must-fix findings,
optional enrichment, **do not add**, and facts verified versus still uncertain.
Use them when they clarify a review; no extra report is required. Optional
enrichment must not silently become a requirement. For the author, check quoted
words, numbers, durations and commit identifiers against their actual sources,
not remembered summaries. Condense notes into teaching rather than exporting
private memory wholesale. Keep reasons and documented experience; old role
assignments, platform workarounds and historical completion claims in
CodeRhapsody's memory are not authority for this edition.

## How to add the next lesson

For each substantive discovery, record the chapter/revision, observation,
supporting command or artifact, interpretation, decision and remaining question.
Label a hypothesis as a hypothesis. When later evidence changes a conclusion,
record the correction and its reason rather than silently erasing the lesson.
Do not claim an experiment ran or a feature worked without the result.

Update the coding skill when Bill accepts a general operational rule; this
document supplies its motivation and history. Put local teaching feedback in the
shared chapter notes for the eventual author. At that later author phase,
preserve the human story and explain why the rules exist, using actual
experiences rather than invented battle scars.

On October 9, Bill generalized the back-pointer rule: an object with two actual
parents must retain both back-pointers, as in a family-tree node. The carryover
and coding skill now describe all actual parent relationships; a single-parent
chain is Ensemble's current example, not a restriction on the rule.
Bill further distinguished required parents (non-nil constructor inputs) from
optional parents that may add the child later. His pin example has a required
IC package and an optional wire. Adding the child modifies both sides of the
relationship: the parent records the child and the child updates its back-pointer.

At the October 9 preparation checkpoint, the previous implementation agents
remain stopped, all discarded second-edition solutions and manuscript drafts
have been removed, and the carryover review team has read Chapters 0–23 and
the epilogue in full. Chapter 24 and the sub-agent chapter remain forthcoming.
The two drafts are ready for review with Bill. No new coder has started; the
new workflow's effectiveness remains to be tested.

Codex's own review on October 9 covered both drafts in full, their consistency
with Bill's instructions, source support for consequential carryover lessons,
links and scope. It clarified the coder's fresh-context boundary, small
save/load scope and the teaching-comment target. The carryover's remaining
source contradictions are explicitly left for discussion with Bill; neither
draft is marked approved.

After that preparation review and subsequent corrections, Bill authorized
autonomous implementation on October 9. The fresh coder is
`student_restart_ch01`, started with only the permitted reading list and no
inherited conversation or old solution. Its working tree is
`solutions/edition-2/main/`; its review goes in `docs/edition-2-notes/ch01.md`.
Chapter 1 subsequently passed independent review: no exception, zero revision
rounds. The original grader scored 100/100 for both implementations, three live
chat turns passed, and the reviewer verified the recorded version-header grader
gap. Student and reviewer notes are in `docs/edition-2-notes/ch01*.md`; the exact
source export is `solutions/edition-2/ch01/`, checkpointed at `edition-2-ch01-r1`.
Chapter 2 was assigned to `student_restart_ch02` from the accepted working
tree in a fresh context. It stopped before editing because the chapter omits
its directive/acknowledgment protocol and the public README is incomplete and
partly obsolete. Bill authorized the coordinator to recover this and similar
missing contracts from the grader autonomously. The exact sourced erratum is in
`docs/edition-2-notes/ch02.md`; Chapter 2 resumes from the accepted Chapter 1 tree.
The unchanged first-edition reference baseline scores 100/100. At this handoff
the student has not implemented, been graded or used a revision round for Chapter 2.
The student subsequently submitted Chapter 2: actual original grader 92/100 FAIL,
live CLI recall on all three providers, live ephemera, nine deterministic
source/target replay pairs and targeted mutation evidence. The reviewer justified
the Ref media scope exception but found one current-scope Anthropic result/text
ordering defect. Initial review is in `docs/edition-2-notes/ch02-code-review.md`.
Revision round 1 corrected the ordering defect; the reviewer independently
verified the regression and restored-defect mutation. Chapter 2 is now accepted
with the recorded media exception after one correction round. The exact source
export is `solutions/edition-2/ch02/`, checkpoint `edition-2-ch02-r3` (the older
r2 tag is preserved). The original grader remains 92/100 FAIL; carry the accepted
media exception through inherited parity until blobs are introduced. No author
or grader revision phase has started.
Chapter 3 is assigned in a fresh context to `student_restart_ch03`, starting
from the accepted Chapter 2 tree. Its unchanged reference baseline is 100/100.
The coordinator supplied missing tool argument keys from the original harness
under Bill's standing authorization, recorded in `docs/edition-2-notes/ch03.md`.
The student still does not read grader internals or old answers; no Chapter 3
result or acceptance has been claimed.
Chapter 3 subsequently passed independent initial review with the inherited
media exception and zero correction rounds. Actual original grade: 90/100 FAIL,
all nine new checks passing, only the inherited media parity check failing.
The student drove all six tools live, including three-provider repair/recall/error
sessions; reviewer checked artifacts, real tests and selected mutations.
Exact source export: `solutions/edition-2/ch03/`, tag `edition-2-ch03-r3`;
older r1/r2 tags remain. Reassess the scoped media exception in Chapter 4 rather
than treating it as permanent or silently rewriting its failed scores.
Chapter 4 is assigned in a fresh context to `student_restart_ch04`, from accepted
Chapter 3. The original reference baseline via `grade-dir` passed 100/100,
including `dlv`. Missing tool/event shapes supplied by the coordinator are in
`docs/edition-2-notes/ch04.md`. Independent scope reassessment confirms the media
exception continues: this chapter records output Refs but explicitly does not
render them as attachment parts. Its own output, recovery and job requirements
remain binding. No Chapter 4 result or acceptance has been claimed.
The Chapter 4 student subsequently submitted a 90/100 original grade, with all
eight new checks passing and only inherited parity failing; independent review
is pending. Its first 65/100 run exposed a contract interpretation error: default
cwd was logged where the grader expected an optional override field to be absent.
The initial coordinator erratum also omitted this representation convention;
the student corrected the code and both recorded the teaching gap. The student
then misread the serialized debugger regex as overescaped and incorrectly
reported that the original calls woke on ten-second delays. The coordinator
accepted that explanation and allowed a repeat. Independent review checked
decoded characters and event timestamps: the original regex was valid and its
calls returned in 0.002–0.270 seconds. The repeat also worked, but was unnecessary
for that claim. Inspect decoded values and measured intervals before diagnosing
escaping or spending another live run; even a purported evidence correction
needs verification. Separately, a process that expired before EOF did not demonstrate shutdown; a
bounded repeat checked that it was alive before EOF and killed afterward.
These findings reinforce checking the particular claimed mechanism, not merely
a successful final answer. Preserve both attempts in the chapter notes; they do
not justify an extra test framework or unbounded live retries.
Initial independent Chapter 4 review requests correction round 1 only for two
lifecycle-test interleavings: a waiter already blocked when killed, and a late
handler completing after kill. The correct production guards are present, but
removing them can evade the current tests because the observations happen too
early or too late. Targeted output-notification coverage does not prove those
distinct transitions. The student is correcting tests and evidence; no new
production feature, paid replay or broad mutation sweep is requested.
Chapter 4 was then accepted with the inherited media exception after two
correction rounds, checkpoint `edition-2-ch04-r2`. Round 1 established the intended
interleavings but introduced a race in the test's failure cleanup. The mutation
output included both the intended assertion and that warning; the student's
initial clean-failure claim was wrong. Round 2 synchronized result transfer and
independent review confirmed clean intended failures for both defects. Check the
whole failure output, including cleanup: a desired assertion does not make an
additional race harmless. Production remained unchanged across both rounds.
Original grading remains 90/100 FAIL, all eight new checks passing; live evidence,
corrections and the exact 25-file export are retained with the checkpoint.
Chapter 5 is assigned in a fresh context to `student_restart_ch05`, starting
from accepted `76bf591` / `edition-2-ch04-r2`. Its frozen original reference
passes the actual 120/120 grader, including Chapter 4 at 100/100. Coordinator
setup facts are in `docs/edition-2-notes/ch05.md`: the harness launches the
framework CLI, despite the chapter describing a separately graded consumer.
The external custom-tool exercise remains required and needs its own evidence;
existing architecture is not a reason to recreate the old refactor or globals.
No Chapter 5 result or acceptance has been claimed.
The Chapter 5 student then submitted its external consumer and review: actual
original grade 110/120 FAIL, three-provider custom-tool use, three-module checks
and targeted mutation evidence. Independent review is pending, including a
requested logger-shape exception: the real multi-hop parent path reaches the
logger, while the structural checker rejects its form. The grader also accepts
an unknown-tool error as tool-call evidence and awards parity points despite the
visible nested 90/100 score; neither establishes the claimed stronger behavior.
Use actual custom results and retain nested failures, not just aggregate points.
One otherwise successful paid run had to be repeated because the launcher read
a session header as an event and deleted scratch before saving the artifact.
Keep sanitized raw evidence before parsing/validation, and clean scratch after
the capture is verified. This is an evidence-capture failure, not a model failure.
Chapter 5 was subsequently accepted at initial independent review, with zero
correction rounds, checkpoint `edition-2-ch05-r2`. The reviewer reproduced the
logger false rejection and verified real logging and registration with source,
external tests, selected mutations and live artifacts. Actual original grading
remains 110/120 FAIL, with nested Chapter 4 at 90/100. The exact 32-file source
export includes the separately built consumer. No chapter or grader was changed.
Preserve small resumable checkpoints if a usage limit interrupts work;
never label unfinished work an accepted chapter.


October 11 checkpoint concurrency: another workspace writer committed the staged
Chapter 16 / Edition 2 Chapter 15 export together with its own first-edition
changes at `3173da2`, between staging and the coordinator’s path-scoped commit.
The coordinator verified the exact snapshot tree and manifest against accepted
`b240993`, preserved shared history, and recorded the provenance in `ch16.md`.
Checking the index before staging does not prevent another writer from consuming
it afterward; use an isolated index for future export commits when other writers
are active. This event changes commit grouping, not the accepted source bytes.


At the next turn, repository HEAD was `009d74c`, no longer descended from the
coordinator’s preceding acceptance/export commits; their exact files remained
on disk as dirty/untracked paths and the immutable `edition-2-ch16-r1` tag still
preserved them. Coordinator verified accepted main and snapshot tree identities,
then re-recorded only its owned artifacts on current history alongside Bill’s
Auto-Recall ruling. No reset, tag move or unrelated restoration was performed.


Original17 pre-submission review supplied two concrete corrections. A real-Agent
probe showed a 6,144-byte capped attachment repeated byte-for-byte on the next
turn because deduplication compared the full untruncated source against context;
a short-snippet control stayed single. Separately, substituting raw term frequency
for BM25 passed the original ranking tests and related integration checks. The
implementation’s formula was correct, but positive hits versus nonhits did not
verify its promised ranking. Preserve this surviving defect and strengthen the
competing-hit evidence. These are attributed reviewer assistance before final
submission, not independent student discoveries or completed acceptance. Raw
baseline/probe/mutation artifacts are retained in original-numbered ch17 evidence.
