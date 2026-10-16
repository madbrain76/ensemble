# Edition 3 crossover closure audit

Read-only audit against repository HEAD `2be995ebc12c6a7efb21cee3b9740a758489adda`. (The orchestrator corrected the original revision transcription against `git rev-parse HEAD` after submission; the student baseline and preface hashes below also match fresh reads.) The working tree already contained Edition 3 preparation/execution changes. This report is the only file written by this audit; no manuscript, grader, runtime, student worktree, or historical evidence was changed. No tests, live-provider probes, or GUI reproductions were run here. Test names below identify inspected evidence, not fresh passes.

## Decision

Distribute the construction obligations, and keep Chapter 25 historical with **no independent exercise**. This is a placement decision already authorized by the request, not a new product decision for Bill.

The smallest dependency-correct ownership map is:

- **Chapter 3:** the visible-explanation policy at the tool dispatch boundary.
- **Chapter 6:** streaming-order enforcement, distinct visible/summary/opaque channels, and part assembly.
- **Chapter 7:** the browser's part-identity consumption; do not reteach the parser there.
- **Chapter 10:** replay/restart identity, including the recorded legacy duplicate-ID regression.
- **Chapter 16:** provider-route/model policy and configurable reasoning summaries.
- **Chapter 20:** policy-bearing root/child construction and base-prompt preservation.
- **Chapter 21:** application file mode, shared configured GUI construction, and Application Settings, using only types available at this point.
- **Chapter 22:** add the existing memory configuration fields once their types/behavior exist.
- **Chapter 23:** add the existing recall fields once the recall contract exists.
- **Chapter 25:** the historical crossover, its evidence, and a short map to those owners.

The small Chapter 23 extension matters: the demonstrated `ApplicationConfig` contains `ensemble.RecallSettings`, while Chapter 23 introduces the public recall patch/defaults. Calling Chapter 22 the final dependency of the *whole* schema would still leave a forward dependency. Chapter 22 is the memory dependency, not the recall dependency.

This map does **not** require a new generic core configuration API, a root-CLI file-mode feature, or moving file configuration into Chapter 1. The demonstrated file-mode interface is `agent/gui/cmd/main.go`; preserve that scope unless a separate product expansion is requested.

## What was read

The required governing documents were read in bounded ranges, recovering truncated ranges: `AGENTS.md`, `.agents/skills/ensemble-coding/SKILL.md`, `docs/agentic-codebook-workflow.md`, `docs/edition-2-carryover.md`, `book/voice.md`, `book/chapter-writing-procedure.md`, `docs/edition-3-execution.md`, and `docs/edition-3-carryover.md`.

Crossover evidence: the full `docs/edition-2-notes/ch25-crossover.md` and `docs/edition-2-notes/artifact-identity-investigation.md`. Relevant full Edition 3 chapter reads: 01, 03, 06, 07, 08, 10, 16, 20, 21, 22, 23, and 25. The new Edition 3 student brief and role skills were also inspected. Current implementation and tests were read only to establish existing contracts and evidence limits.

Neither preface was edited. At the end of inspection both `book/preface.md` and `book/edition-3/preface.md` had SHA-1 `f80a65996d4e710835e502e7518b2e4f14d59dd5`.

## 1. Already integrated teaching versus actual gaps

| Boundary | Already present in Edition 3 | Remaining crossover work |
| --- | --- | --- |
| Ch3, tools | Strict tool arguments; unknown/bad calls become model-visible results; usable conversation after a recoverable error. | The opt-in visible-explanation gate, all-or-nothing pre-dispatch rejection, preserving the rejected assistant turn, and bounded corrective resend. |
| Ch6, streaming | Provider-local versus application part identity; summary/opaque separation; summary-only delivery flags; one parser contract. | Preserve observed arrival order for the explanation gate; align fixtures with nonempty artifact assembly and empty-final suppression. Do not add a second summary-delivery essay. |
| Ch7, GUI | Browser treats `part_id` as the address of a card; replay/live reconnect protocol and sequence watermark. | Test the existing identity consumption against the restart/legacy scenario, coordinated with Ch10. |
| Ch10, persistence | Explicit promise that old saved IDs and new live parts do not collide; save/load and restart identity are already exercise obligations. | A test and acceptance statement for **pre-existing** duplicate IDs, not merely a higher next counter. Current implementation is not demonstrated to meet the promise. |
| Ch16, provider route | Distinct plan/API credentials, no fallback, request-shape restrictions, provider-neutral core, credential rereads. | Explicit named-model capability activation for the gate; allow configured detailed summary mode rather than teaching only a literal `auto`; preserve actual route policy and native response shape. |
| Ch20, fleet | Shared factory, sink plumbing, prompt ownership, interruptible input, child lifecycle, scoped IDs. | Carry the effective model policy through child creation without repeatedly mutating the base system prompt. |
| Ch21, GUI/fleet | Public constructor, shared Ensemble bootstrap, engine-free GUI adapters, settings policy boundary. | The demonstrated versioned application file, explicit file/legacy entry selection, environment isolation, limited launch overrides, persistent Application Settings, and usable UI. |
| Ch22, memory | Memory policy sources, runtime updates, lifetime and persistence ownership. | Extend the Ch21 application configuration with memory fields; round-trip/apply them through the existing memory policy path. |
| Ch23, recall | Structured recall patch/defaults and distinct immediate `recall` operation. | Extend the same application file with its already-demonstrated recall settings, now without future types. |
| Ch25, crossover | Honest historical account and evidence links. | Remove the assertion that whole-application configuration must remain here; replace it with the resolved earlier ownership map. |

The existing prose is not evidence that a new Edition 3 grader checks these cases. Conversely, an old implementation test is not an Edition 3 chapter pass.

## 2. Exact minimal manuscript placement

### Chapter 3: dispatch owns the gate

At §3.2, immediately before the tool-dispatch walk, add the small optional policy and its failure behavior. In §3.7/exercise invariants, require deterministic enabled/disabled tests. Introduce only the policy information needed at this chapter; do not import the eventual provider or fleet schema merely to teach one dispatch decision. Chapter 16 later selects the policy for named models.

Established contract:

1. When this capability is enabled and an assistant response requests tools, ordinary nonempty user-visible text must precede its tool calls in that response.
2. Encrypted/opaque reasoning and reasoning-summary text do not satisfy it. It is a presence/order check, not an arbitrary minimum-length or semantic-quality classifier.
3. Reject the **whole** nonconforming tool-call batch before dispatch. No requested tool runs, no tool job is opened, and no success-shaped result is fabricated.
4. Preserve the original assistant turn, including opaque data. Pair every rejected call with a model-visible error and corrective instruction, then let the model resend inside the existing tool-round bound. Do not loop indefinitely or turn a recoverable rejection into a destroyed session.
5. Disabled models retain their previous behavior.

Evidence: `agent/internal/agent/visible_reasoning.go`, `agent/internal/agent/actor.go` around the pre-dispatch check, `agent/ensemble/visible_reasoning_test.go`, `agent/ensemble/visible_reasoning_prompt_test.go`; crossover notes lines 331–369.

### Chapter 6: streaming can disprove an apparently valid assembled answer

At §6.8 (record-keepers) and the exercise fixture list, connect arrival order to Chapter 3's policy. A tool-call delta arriving before visible text cannot be excused because later assembly groups text before calls. Keep both final-response validation and observed stream-order validation.

Add/retain fixtures for separate ordinary text, summary-only streams, opaque reasoning, multiple text/reasoning parts, empty final counterparts, and a tool-first stream followed by late text. The existing §6.6–6.7 explanation of capability flags and summary-only streams is already integrated; only align its implementation/grader contract.

Evidence: `agent/internal/agent/actor.go` records `toolFirstStream` before checking `visibleExplanationBeforeTools`; `agent/internal/engine/responses_stream.go`; artifact investigation section on nonempty parts/empty finals.

### Chapters 7 and 10: consume identity, then prove restart stability

Keep parser/allocator teaching in Chapter 6, GUI card addressing in Chapter 7, and persisted/restarted identity in Chapter 10. Add one shared regression to their acceptance map rather than duplicating implementation advice.

At Chapter 10's save/load/replay contract and exercise checks, state that the known legacy saved history must display distinct historical messages/cards and keep new live output separate across another reload. Old and new artifact associations must remain stable. A high-water counter alone is not acceptance of this legacy case.

Do not prescribe a new allocator, sidecar, compound key, migration format, or archive rewrite. The requirement is behavioral. An implementation proposal that sacrifices historical separation/stability or destructively rewrites archives would require a separate decision; merely choosing an implementation that satisfies the existing guarantee does not.

### Chapter 16: explicit provider/model policy

Extend the existing route/model policy discussion and wire fixtures, rather than adding a late crossover mini-chapter:

- Activate the visible-explanation requirement by explicit model capability, not a name-prefix guess or a blanket rule for every route/model.
- Preserve the native credential-route distinctions and already-taught omission of unsupported request knobs. No fallback to a paid key.
- Make `reasoning.summary` configurable, including the demonstrated `detailed` value; preserve the existing `auto` default. Summary receipt remains distinct from visible narration.
- Keep provider reasoning/summary/capability details out of generic UI dispatch.

The current chapter's native-body example specifies `summary:"auto"`; the small change is a configured value with its default and a detailed-summary test, not a changed global default.

Evidence: `agent/internal/common/models.go` (`RequireVisibleReasoning` for the demonstrated Astra entry); `agent/internal/engine/media.go` (`Responses`, explicit-cache and summary-only stream flags); `agent/internal/engine/capabilities.go` (`responseSummary`); `agent/internal/engine/responses.go` and `responses_stream.go`; `agent/ensemble/codex_test.go`, especially `TestCodexDetailedSummaryFromEnvironment`. These are not evidence for an additional pricing requirement or a new universal temperature rule.

### Chapter 20: child construction must preserve the policy

Add a short invariant next to the shared factory and skill/child prompt composition, plus one fixture in the fleet exercise. The shared root/child construction path selects effective instructions for the chosen model while retaining an unmodified base system prompt. A child's selected model must not inherit an inappropriate parent-only policy or accumulate duplicate policy paragraphs.

Evidence: `agent/internal/engine/media.go` (`BaseSystemPrompt` and `EffectiveSystemPrompt`), `agent/internal/fleet/child.go`, and `TestVisibleReasoningPolicyFollowsChildModelNotParent` in `agent/ensemble/visible_reasoning_prompt_test.go`.

No new mandatory root-CLI file-loading surface follows from this lesson.

### Chapter 21: staged application configuration at the existing constructor

Extend §21.1's public constructor/configuration boundary, §21.6's composition root, and §21.7/§21.9's exercise checks. This is where the versioned application file and its usable GUI policy surface naturally belong; do not put a complete application schema into the early browser-preference store of Chapter 8.

Use the demonstrated shape, with only already-introduced features: provider/route/model/reasoning configuration, workspace/skills/data/save/log/auth locations, GUI defaults, streaming, web-search/sandbox defaults, and port. Preserve the existing boundary between application-owned settings, lower-case disk/UI fields, and provider wire fields. Omit future memory/recall members until their owning chapters.

Preserve the demonstrated entry contract:

- Explicit `--config path` selects file mode. It does not call environment configuration builders.
- Inherited provider/model/endpoint/reasoning/skill/auth-store environment settings do not override the selected file.
- Competing agent/provider launch flags are rejected in file mode.
- Only the established explicit `--port`, `--gui-dir`, and `--gui-debug` launch overrides are accepted/reported there.
- Without `--config`, the existing environment/flag mode remains available and unchanged.
- Validate the existing schema/version, types, URL/model/route constraints and normalized paths. Do not invent automatic file discovery or a general environment-over-file precedence ladder.
- Save settings persistently and apply them through the public construction/runtime path. Restart-required changes are reported honestly; there is no claim that a live provider/skill was invisibly replaced.
- GUI-enabled acceptance includes a usable editable Application Settings view and a real resulting runtime request, not just HTML or an API-shaped response.

Evidence: crossover notes lines 276–323; `agent/gui/application_config.go` (`LoadApplicationConfig`, `Save`, normalization, `RuntimeConfig`); `agent/gui/configured.go`; `agent/gui/cmd/main.go`; `agent/gui/entry_config_test.go` (`TestFileConfiguredEntryIgnoresEnvironment`); `agent/gui/application_settings_test.go` (`TestApplicationSettingsWireAndDisk`). Retain the existing optional-GUI execution distinction; making this a new mandatory core-only CLI feature would be additional scope.

### Chapters 22 and 23: finish the schema only as dependencies arrive

Chapter 22: add the established memory preference members to the application-file/disk/UI projection and exercise round trip, applying them through the memory settings authority already taught here. Do not create another memory-policy source.

Chapter 23: add the already-demonstrated recall settings member/defaults to that same schema and verify effective configuration. `ApplicationConfig.Recall` is in `agent/gui/application_config.go`; the public recall patch/defaults arrive in Chapter 23's code/protocol contract. This is a small schema extension, not a new recall feature or a reason to move recall into Chapter 22.

### Chapter 25 and governing cross-references

Keep the historical crossover narrative and evidence links. Replace §25.2's “whole-application configuration belongs at the end” rationale with a compact ownership recap. Update the TL;DR/ending if they still imply Chapter 25 is the configuration construction stage. Do not create an executable exercise or fabricate an independently passable Chapter 25 grader.

Update `docs/edition-3-carryover.md` row 25 and execution-gate item 3 to record this resolved distribution. Update the Chapter 25 paragraph in `docs/edition-3-execution.md` from an unresolved mapping gate to the accepted map and remaining acceptance evidence. Preserve the linked historical notes rather than rewriting them as if the lessons had been present earlier.

## 3. Artifact collision: actual blocker, not forced architecture selection

The full investigation records a concrete regression, not merely a concern: a saved history containing three messages whose text parts all used `part_id:1` rendered as concatenated history in one card. The earlier allocator/assembly work did not solve that replay. See `docs/edition-2-notes/artifact-identity-investigation.md`, especially the two cold restarts and conclusion.

Its separate finding that archive content remained distinct matters: do not describe this as text concatenation already stored on disk, nor blame the loaded skill without evidence. It also records that clearing the browser map or using a larger counter alone does not repair duplicate IDs already in the archive.

Current-source inspection is consistent with an unresolved allocation/restoration problem: `agent/internal/agent/actor.go` has the runtime `partID`, creates a fresh runtime without restoring it, and increments it in `displayID`; `agent/internal/agent/save.go`'s `Load` delegates to history recovery without demonstrated allocation-state restoration. This is source-level evidence, **not** a fresh reproduction of the browser failure on this audit date.

Close the preflight policy gap by recording the existing guarantee and the actual legacy fixture. Leave implementation choice to construction, and require a real passing result before claiming the defect resolved. Acceptance must cover both: (a) collision-free old/new parts for clean saved IDs, and (b) separation/stability when the loaded archive already contains duplicate IDs. The second is the known missing case. This can be made an explicit Chapter 10 acceptance blocker without asking Bill to select an allocator architecture now.

## 4. Grader alignment and evidence status

The new Edition 3 harness must receive the relocated obligations at their chapter boundaries. Do not patch the old `cmd/gradechapterNN` programs or reinterpret old passes as new coverage.

Minimum acceptance additions/checks:

| Owner | Deterministic acceptance case |
| --- | --- |
| Ch3 | Enabled/disabled model policy; tool-first response; opaque/summary-only response; preserved rejected assistant data; every call gets paired error; zero tool side effects; corrective resend; bounded fail-closed exhaustion. |
| Ch6 | Tool delta before late visible text remains rejected; summary-only delivery works; multiple parts remain separate; empty final artifacts do not erase nonempty streamed state. |
| Ch7/10 | Reloaded transcript and new live response remain separate; known legacy duplicate-ID archive yields stable distinct cards through another restart. Add the GUI diagnostic only in a GUI-enabled run. |
| Ch16 | Actual native/API route requests; no credential fallback; configured detailed summary arrives; existing default and unsupported-field omissions preserved; only explicitly enabled models receive gate policy. |
| Ch20 | A child's selected model gets the right effective prompt while the original base prompt is retained and no duplicate policy is accumulated. |
| Ch21 | File/legacy selection; noisy inherited environment ignored in file mode; competing flags rejected; only documented explicit overrides; invalid file values rejected; settings persist/reload/apply; captured provider request uses file selection; usable GUI editing when enabled. |
| Ch22/23 | Memory then recall extensions round-trip and affect the established owner, without premature type dependencies or a second settings authority. |
| Ch25 | Historical/no-exercise recording, not an executable feature-grade result. |

The parent reports that Chapter 1's no-short-total-HTTP-timeout coverage is being added to the new Edition 3 harness; this audit did not reimplement or verify it. Likewise, the corrected repo-root grader cwd and empty student baseline `69859da3328c3f710f8143203b99f934fed778a1` are parent-reported execution setup, not work performed here. The audit is not a student execution.

## 5. Decisions and remaining work

### Resolved within existing authorization

- Chapter ownership above; no late Chapter 25 exercise.
- Stage the existing application schema at its constructor and feature owners rather than add future dependencies or invent a core loader.
- Keep exact file-mode environment isolation and unchanged legacy mode.
- Keep the pre-dispatch gate optional and model-selected; count only ordinary visible text and preserve streaming arrival order.
- Treat the known legacy artifact behavior as required acceptance, not as an invitation to choose architecture now.

### No Bill decision is needed to perform those small changes

A Bill question would become necessary only if proposing a genuine behavior/scope change: dropping the legacy stability guarantee, destructively rewriting stored archives, broadening file mode to a mandatory new core-only CLI surface, changing the default summary mode, removing the existing recall field from the final application schema, or replacing the optional-GUI execution policy. None is necessary for this closure.

### Still not done by this read-only audit

The manuscript/support edits and new-harness tests are proposed, not applied. The legacy artifact regression remains an implementation acceptance blocker until demonstrated fixed. The historical tests and migration logs establish prior evidence but do not prove a fresh Edition 3 student implementation. Do not mark all crossover work complete merely because this ownership map is accepted.
