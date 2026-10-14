# Background compaction: chapter dependency audit

October 14, 2026. Independent dependency reviewer, GPT-6 Astra/high.
Read-only audit of baseline `56ec71cb` in `/Users/bill/projects/ensemble-compaction`.
This note is a recommendation for Bill's placement decision, not acceptance of
new implementation, manuscript, graders or snapshots. No graders or live/runtime
operations were run by this reviewer; no source files were edited.

## Recommendation and decision

Adopt the proposed move of current Chapters 15 and 16 together after current
Chapter 23, provided Bill approves two consequences: stage accounting's
memory/recall-specific obligations at their new introduction points, and produce
new reviewed cumulative solution exports rather than rename the old ones.
The move puts the sandbox, child lifecycle and fleet inspection before the
new background compressor that uses them. Auto-Recall still follows memory.
There is no unavoidable earlier exercise dependency that requires memory or
recall before the fleet chapters.

| Current chapter | Proposed chapter | Material |
| --- | --- | --- |
| 0–14 | unchanged | Through recorded context management |
| 17 | 15 | The Invisible Invoice |
| 18 | 16 | Leaving Anthropic for OpenAI |
| 19 | 17 | The Crossover; still ungraded |
| 20 | 18 | Giving the Agent the Web, and a Status Line |
| 21 | 19 | The Sandbox |
| 22 | 20 | The Ensemble; public reuse and child lifecycle |
| 23 | 21 | Managing the Swarm |
| 15 | 22 | Forgetting on Purpose; revised background implementation |
| 16 | 23 | Auto-Recall |
| 24–26 | unchanged | Optional route, reserved crossover, Gateway |

The narrow approval question is whether Bill wants this mapping with staged
accounting and newly validated cumulative exports replacing the current
book-facing `solutions/ch15`–`ch23` references, while old source/manifests/results
remain reachable through immutable checkpoints. Carry the revised cumulative
implementation into Chapters 24 and 26 and validate their affected behavior too.
This is reviewed maintenance/reconstruction, not a fresh blind student experiment.
No new blind build is implied by the current maintenance assignment.

## Material dependencies

1. **Accounting has real forward obligations that must move, not disappear.**
   Current `book/chapter-17.md:49` includes compressor/judge usage in Engine's
   accumulator, with separate latest-main counts; `:70` excludes those requests
   from main-request cache comparisons. `:88` requires the surviving-memory
   fallback after a handoff is compressed. `:425`, `:433` and `:447` explain recall
   clearing, folded boundaries and expected prefix rewrites. Preserve main-request
   accounting, stable-prefix inspection and Chapter 14 handoffs in new15. Complete
   the memory-boundary and compressor attribution contracts in new22; complete
   judge attribution and recall integration in new23. This changes their teaching
   stage, not their eventual requirement. Current `book/chapter-20.md:116` and
   `:935` similarly need staged wording in new18: session versus latest-main usage
   works before auxiliary requests exist and must remain correct after they do.

   The numerical invoice grader does not currently prove all those stronger
   claims: `internal/gradeedition2/ch17.go:257` explicitly states that its marker
   fixture has no ephemeral tail, handoff or memory. Separate implementation
   tests cover them, including `agent/ensemble/cache_test.go:193` (surviving
   memory), `:275` (auxiliary model identity), and `:342` (unpriced auxiliary
   model). Preserve or adapt these later-stage obligations; do not claim an
   unchanged invoice score proves the new child accounting path.

2. **Authentication has bounded explanatory references, not a required memory
   API.** `book/chapter-18.md:258` describes the cuts/memory/recall of current
   14–16 as existing context-mutation controls. In new16 retain the already
   implemented Chapter 14 controls and point forward to new22/23 for memory and
   recall. The authentication, billing-route and Responses exercises do not
   need those features to precede them. Preserve dated historical claims and
   demonstrations as history, not fresh runs of the reordered solution.

3. **Sandbox and fleet must precede the new compressor.** Current21 supplies
   child confinement; current22 supplies actual child construction, submission,
   bottom-up completion, messaging and bounded idle correction; current23
   supplies visible inspection, activity, result retrieval and lifetime controls.
   Their order is retained as new19–21. They do not require memory first.
   `book/chapter-22.md:84` and `:93` make factory-enabled submission behavior
   explicit while ordinary embedders retain ordinary final replies. The new
   default internal factory must not accidentally turn every existing embedder
   into a submit-required fleet. Teach the distinction between construction
   availability and public delegation/submission enablement in new20.

4. **The optional submission schema belongs with fleet creation, before memory
   consumes it.** Current `book/chapter-22.md:487` specifies `submit(result)`;
   `book/chapter-23.md:297` specifies result retrieval and durable `submit.md`.
   The current API is explicitly string-shaped:
   `agent/internal/common/fleet.go:38`, `:54`, `:60`, and `:79` expose status
   result, retrieval, spawn and submission. `agent/internal/fleet/management.go:66`
   also decodes a string fallback. Teach and verify optional launch schemas,
   validated completion and structured result preservation in new20/21; new22
   then uses the memory schema. Default text callers and historical saved text
   must retain their documented compatibility. A string convenience method must
   not bypass an explicitly custom schema. This is the authorized handoff's
   scope, not a proposal to redesign all messaging or build a generic workflow.

5. **Recall is a genuine successor of memory.** Current `book/chapter-16.md:5`,
   `:90`, `:100` and `:106` require the memory bands, archive corpus and exclusion
   of content already actually present in context. Moving the pair preserves
   that prerequisite. The judge is expressly one stateless call (`:119`, `:393`),
   so a compressor implemented as a child does not authorize converting the
   recall judge into another Agent. Review its actual interaction with late
   memory commits, context deduplication and turn cancellation when new23 is
   validated. Preserve Bill's quota/all-memory decisions.

6. **Historical discussions are not forward implementation requirements.**
   Current20's constructor table mentions `NewJudge` and `recall.New`
   (`book/chapter-20.md:730`), but `:737` explicitly labels the numbering as
   original. Retain that documented diagnosis; do not mechanically replace
   original chapter numbers or delete the story because its repaired features
   now arrive later. Current14's forward sentence at `book/chapter-14.md:864`
   should point to new22. The earlier historical crossover remains a historical
   demonstration, with its own source identity, when it becomes new17.

## Grader chains, routing and fixtures

There is no recursive memory→recall→invoice parity chain. Current memory calls
`Chapter14` directly at `internal/gradeedition2/ch15.go:20`; retain that context
parity when memory becomes new22. Replacing it with new21 merely because that
chapter precedes it would silently change the check's meaning. Current16,17,18,
20,21 and22 use feature rubrics rather than recursively running the preceding
chapter. Current23 reuses `fleetMessaging` and `fleetIdle` from22 at
`internal/gradeedition2/ch23.go:27`, then checks its own GUI/management behavior.
Current22 composes original24 and25 weights and reuses the actual model-switch
session from20 (`internal/gradeedition2/ch22.go:11`, `:30`); this is preserved
as new20 depending on new18 support, not on memory.

Helpers such as `recallWrite`, `recallReadEvents` and `memoryConversation` are
shared grader utilities, not student prerequisites. Their use in invoice,
authentication, sandbox or fleet drivers does not imply those features must
exist in the submission. The fleet peer explicitly disables recall with
`EN_RECALL_DISABLED=1` (`internal/gradeedition2/ch22_fleet.go:49`). Preserve
external-provider fakes and the real earlier implementation beneath them.

The memory driver assumes serial completion: `internal/gradeedition2/ch15.go:86`
constructs fixed request expectations, and its checks inspect completed summary
results. Adapt these to observable child/candidate completion barriers without
discarding the eleven existing memory weights. Add honest coverage for the
authorized background/schema behavior. A total-request ordinal is unsafe once
parent and compressor requests overlap. Separate the three identities:
ordinary parent requests, submit-only memory workers, and tool-free judges.
Do not identify every tool-free final/terminal request as recall.

The command's numeric dispatch cannot be changed alone:

- `cmd/grade-edition2/main.go:31`–`:38`: help text and chapter-specific options.
- `:41`–`:47`: valid/ungraded chapters and implicit `solutions/chNN` selection.
  New17 is ungraded; new19 is graded. Optional24 and reserved25 remain ungraded.
- `:76`–`:84`: executable routing. Current14,15,16,17,18,20 use the GUI entry;
  after mapping that is14,22,23,15,16,18. New19,20,21 use core CLI and separately
  built GUI probes where needed. Gateway26 retains its own command.
- `:93`–`:155`: feature dispatch and report identities. Keep feature identity
  separate from predecessor arithmetic. Shared `internal/grade` names such as
  `Ch16Checks`, `Ch17Checks`, `Ch24Checks` and `Ch25Checks` identify original
  rubrics and are not wholesale rename targets.

Root reports baseline current14/15/16 grading on explicit `./agent` at100/100.
Root also reports that old22/23 CLI grading stops before checks because
`grade.Build` chooses the present `agent/cmd/` directory without a Go program;
a temporary runner using a CLI built from `agent` obtained200/200 and100/100
from the unchanged feature graders. Those are coordinator-supplied observations,
not this reviewer's executed checks. Preserve the launch failure separately and
repair routing before claiming the final advertised commands work.

## Cumulative exports and historical access

Old17 already contains memory and recall, for example
`solutions/ch17/internal/common/types.go:247`, `:282`, `:559`. Old23 constructs
both services (`solutions/ch23/ensemble/ensemble.go:253`, `:305`). Renaming old17
to new15 would place future features in the reference and would not test that a
reader can follow the new sequence. Its old grade cannot become new-order
acceptance. The same problem carries into old24/26 even though their numbers do
not move.

Recommended concrete approach, conditional on Bill's decision:

1. Keep all existing immutable tags, original-numbered notes, raw evidence,
   manifests and first-edition sources unchanged. Record a dated old↔new table
   and identify which checkpoint supplies each historical reference.
2. Reconstruct cumulative current15–21 in isolated source stages from accepted14,
   bringing in only each stage's capabilities. Then add revised memory as22 and
   recall as23. These may reuse reviewed source; explicitly label the work
   maintenance reconstruction with prior-answer exposure. Validate each stage's
   public exercise and affected inherited behavior, plus independent scope,
   architecture and comparative review. A final full-feature `agent` pass cannot
   establish those intermediate boundaries.
3. Produce exact exports/manifests from committed, accepted staged source and new
   immutable revision tags. Replace book-facing directories only after approval
   and acceptance. Extend/revalidate the changed cumulative implementation for
   optional24 and Gateway26; retain25 as a placeholder without an accepted export.
   No symlinks, fabricated student runs or new history directory are required.
4. Make historical grading explicit through the preserved checkpoint checkout
   and its command, with an explicit submission path. Current grading selects
   the newly accepted mapping. Until replacement exports exist, require an
   explicit submission for affected current chapters or fail clearly rather than
   silently grade the old default path. A separate historical selector in the
   current binary is optional; immutable checkout commands already preserve the
   exact old grader and fixture identity without maintaining duplicate engines.

The top-level preservation rule forbids editing frozen snapshots during normal
development (`.agents/skills/ensemble-coding/SKILL.md:10` and
`book/chapter-writing-procedure.md:38`). The workflow permits reviewed
implementation corrections but rejects silent frozen edits
(`docs/agentic-codebook-workflow.md:2103`). Its explicit earlier replacement
precedent is the combined ch20 export (`:849`), replaced only after independent
acceptance while its old tag stayed immutable. The procedure's historical
revision-export mechanism (`book/chapter-writing-procedure.md:568`) likewise
does not turn renamed source into validation. Bill's approval should explicitly
cover the proposed current reference replacement and maintenance review status.

## Change map for the authorized later stages

- Teaching: `book/chapter-14.md` forward pointer; current
  `book/chapter-{15,16,17,18,19,20,21,22,23}.md` filenames, headings, section
  references, examples, commands and staged contracts; inspect24/26 references.
  Preserve Bill's `book/preface.md` verbatim.
- Navigation and build declarations: `book/chapters.md`, `book/README.md`,
  root `README.md`, `Makefile`, and current mapping statements in
  `docs/edition-2-chapter-map.md`, `docs/edition-2-enhancement.md`,
  `docs/edition-2-publication.md`, workflow/carryover/current instructions.
  Add current clarification above dated history; do not globally rewrite it.
- Grading: `cmd/grade-edition2/main.go`, affected feature families under
  `internal/gradeedition2/ch15*`, `ch16*`, `ch17*`, `ch18*`, `ch20*`, `ch21*`,
  `ch22*`, `ch23*`, plus affected26 and build-selection support/tests.
  File renaming is a maintenance choice; dispatch and published identities must
  agree. Preserve original `edition-1` graders and shared original rubric meaning.
- Public/runtime contracts: `agent/internal/common/{types,fleet}.go`,
  `agent/internal/fleet/`, `agent/internal/memory/`, Agent/Engine request and
  completion paths, `agent/ensemble/`, public GUI result/inspection paths and
  source tests affected by `Engine.Compress` or string submissions. This is a
  dependency list, not authority for speculative APIs.
- References: `solutions/README.md`, approved replacement `solutions/chNN/`
  exports and sibling manifests; preserve old tags and historical evidence URLs.
  `scripts/edition2/export_snapshot.py` exports committed `agent` by default and
  accepts historical prefixes; it creates neither acceptance nor replacement
  authorization (`:5`, `:73`).
- Publication inputs: inspect `scripts/build-ebook.sh`, `book/build-print.py`
  and the generated assembly's chapter order when renumbering is accepted.
  Existing EPUB/PDF/assembled Markdown remain their earlier publication artifact;
  no publication build or upload occurred in this audit.

Reading coverage: full handoff; AGENTS; mandatory coding skill; full workflow,
carryover, voice and chapter procedure; full current Chapters14–23; chapter map
and solution README; command dispatch; relevant feature grader, driver, fixture,
public-interface and frozen-source sections cited above. After one compaction,
the full mandatory skill and all four governing/craft documents were reloaded
before completing this review. This is not a claim to have reread every runtime
file or every grader line, nor a formal implementation acceptance.

## Follow-up: approved mapping and build-selection preparation

October 14, resumed at coordinator baseline `7a362ec2`. The coordinator reports
Bill approved the proposed mapping, staged accounting and validated cumulative
export replacement. The decision paragraphs above remain the original audit;
they no longer represent an outstanding permission gate. New exports still need
their implementation, validation and independent review before acceptance.

The narrow build repair is justified by the existing advertised selection
contract, without changing any chapter obligation. `internal/grade/harness.go:86`
says to find a main package first at `cmd/`, then root; `:104` actually selects
any directory named `cmd`. Independent metadata inspection with
`go list -e -json ./cmd/` from `agent/` returns no package name and the explicit
error `no Go files in .../agent/cmd`. This is a command-container directory,
not the CLI package.

Recommended implementation: ask Go for the candidate's package metadata, select
`./cmd/` when it identifies a main package, otherwise build the root. Keep the
existing single selected build; a selected main package's compilation error
must fail instead of triggering a build of a healthy root. Go's error-tolerant
metadata mode can still identify `main` when dependencies are invalid. Avoid
recursive command discovery or choosing among arbitrary child executables;
Gateway already passes its explicit command directory.

Focused real-module tests are warranted because this selector serves earlier
graders as well as current22/23. Check root-only, a `cmd` container with a nested
command, preferred `cmd` main beside a healthy root, a non-main `cmd` package,
and a selected `cmd` main that fails compilation despite a healthy root. Execute
successful binaries and compare distinct output so a successful build of the
wrong target cannot pass. Preserve the direct-file input behavior. These tests
need neither mocks nor model inference. No existing `internal/grade/*_test.go`
files were found during preparation. This is a strategy review, not approval of
an implementation that has not yet been submitted.

For later reconstruction, use the existing export tool and short source-stage
commits instead of a new generator. Start from14, then validate new15 accounting,
16 authentication,17 historical crossover's bounded integration,18 web/status,
19 sandbox,20 fleet/schema,21 management,22 background memory, and23 recall.
At each stage verify that the new capability is present and later memory/recall
services have not accidentally arrived in an earlier cumulative tree. Carry
the resulting23 into optional24 and Gateway26, preserving their existing route
and Gateway exercises and adding only affected integration checks. Do not make
unrelated live repetitions or an imagined chapter25 acceptance part of this
work. Detailed source selection waits until the behavioral stages and prose
contracts are accepted; the sequence itself is now approved.

## Independent build-selection review

October 14, reviewer disposition: **accepted for the narrow build repair**, with
zero correction rounds and no blocking findings. This does not accept later
submission-schema, compaction or chapter work. Reviewed against `7a362ec2`:

- `internal/grade/harness.go`, SHA-256
  `43f93da9b67504fc9cb69a88572fdc3ebd382463a2d5e31c6703fd422c01fb7c`.
- `internal/grade/harness_test.go`, SHA-256
  `05877258c04632ae5ba4b14eecb1f3d1307b1f37b27c091b85d20851ddcb9335`.

The implementation queries Go package metadata only when `cmd` is a directory,
selects it only for `main`, and keeps exactly one subsequent build. Metadata
failure returns an error and cleans the temporary directory; selected command
compilation failure is never replaced with a successful root build. Direct-file
input and root-only behavior are unchanged. Go, rather than a hand-written source
scanner, determines the package under the effective build constraints. No
recursive command search, new dependencies, chapter assertions or runtime changes
were introduced.

For the `Build` function, counted from its declaration to the following function
comment, physical lines grow31→43 and active lines27→36; active means nonblank,
non-`//` lines. The nine active lines implement the necessary metadata query,
error cleanup and selection. The explanatory comment states why directory
existence is insufficient and why a broken main must remain selected. The five
real-module scenarios execute successful binaries and distinguish their output;
the broken-main case requires the actual undefined-symbol compiler diagnostic.
This is small, behavior-based coverage of the reported failure and retained
selection semantics, without mocks or a test framework.

Independent verification copied the exact submitted `Build` function and
unchanged test file into an isolated temporary Go module, adding only its
standard-library imports. It did not modify shared source or invoke a grader:

| Isolated source | Result | Intended observation |
| --- | --- | --- |
| Submitted implementation | PASS,0.319s | All five scenarios |
| Always select root | FAIL,0.328s | Preferred command executes root; broken preferred command wrongly succeeds |
| Original directory-existence predicate | FAIL,0.186s | Command container has no Go files; library archive cannot execute |

The failures occur at the intended runtime/assertion boundaries, not because the
mutated harness cannot compile. Temporary probe modules were removed. These
mutations establish the selector's tested distinctions, not exhaustive coverage
of every possible Go toolchain failure.

The reviewer also tried a temporary main importing an unavailable module, with
network lookup disabled. Both `go list -f '{{.Name}}' ./cmd/` and its `-e` form
returned `main` and exit0 in this environment. Thus a missing-import test would
broaden compilation-error examples but would not establish sensitivity to `-e`
here. The existing undefined-symbol failure already protects against hiding a
selected command's compile failure; no additional blocker is raised.

Reviewed the coordinator's complete raw reports and invocation metadata at
`/tmp/ensemble-compaction-baseline/grader-22-buildfix{,-run}.json` and
`grader-23-buildfix{,-run}.json`: advertised explicit-`./agent` commands exit0,
Chapter22 scores200/200 in2.489s and Chapter23 scores100/100 in1.125s, all rows
passed. Those are coordinator-executed checks, independently inspected here.
The preceding command-selection failures remain historical baseline results.
Broader package regression checks were still running when this disposition was
written; their eventual outcome belongs in the coordinator's progress record.
No model inference or live application operation was performed for this review.
