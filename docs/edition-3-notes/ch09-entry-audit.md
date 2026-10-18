# Edition 3 Chapter 9 entry audit

## Decision and scope

**A thin Edition 3 dispatch/validation adapter is justified; it is not implemented here.** The current Edition 3 command rejects Chapter 9 before target launch. The existing exported historical Chapter 9 grader supports separate executable/core identities and passed a complete real-source positive, including inherited Chapter 8 parity. One compiling mutation of real skills code failed the existing focused test, and exact source restoration passed. Reuse those checks; do not change their contract, score, or timing.

This is an independent, finite entry audit, not a student attempt, candidate acceptance, manuscript review, or live-provider campaign. The accepted Chapter 8 student was read only for pins/status; no student code was implemented, built, graded, or flattened. No children were launched. No adapter, source, grader, guidance, student, manuscript, publication, commit, or tag was changed.

**Caveat:** this audit had a preserved audit-driver setup failure and an initially incomplete checklist. The complete positive was subsequently run once, with the correctly inspected API; it was not a rerun of a failed grader result. Exact provenance of the short-lived external helper binaries was not captured, although their complete source identities, composition paths, build mechanisms, and successful inherited parity are known. See limitations below. Root should inspect these receipts before accepting entry readiness or authorizing the adapter.

## Governing pins and full reads

Assignment: `/tmp/cadence-ch09-entry-task.txt`, 6,320 bytes, SHA256 `2f102659c352e244f668226323028bdf853b8af29d7426a1ab04dfa3c892419a`.

- Main: `/Users/bill/projects/ensemble`, HEAD `5614faed78f10f12a940eef6573a3bfe50aa3f6c`.
- Student: `/Users/bill/projects/ensemble-edition-3-student`, clean HEAD `b8f22ba5b3c9d1fc52d1a114810b90b77928684d`, tree `fa896caa404b87b08de61eb38ce079eb38138fbe`; annotated `edition-3-ch08-r1` peels to that commit. The annotated tag-object ID is not the commit ID.
- Governing `docs/edition-3-notes/ch09-entry-manifest.json`: SHA256 `ad13c49983be9366b0e9445c2e4edd6723a4ed7faae730565eb6490e0c146a2b`.
- Before controls, 174 manifest file records and three student ref/tree checks passed. Final verification repeats all file checks and refs and compares all initially tracked file bytes. Main already had unrelated modifications/untracked files; this is not a claim of a clean main worktree.

Fully read before the first build/control: `AGENTS.md`; `.agents/skills/ensemble-coding/SKILL.md`; `.agents/skills/edition-3-code-reviewer/SKILL.md`; `docs/edition-3-execution.md`; `docs/edition-3-readiness.md`; `docs/edition-3-carryover.md`; `docs/edition-3-student.md`; `docs/agentic-codebook-workflow.md`; `docs/edition-2-carryover.md`; `book/voice.md`; `book/chapter-writing-procedure.md`; and all 448 lines of `book/edition-3/chapter-09.md`. Current Edition 3 student-first procedure governed; historical author-first/model/STOP instructions did not authorize work. `prerequisite-reads.json` records file hashes/lengths and recovery of clipped reads. In particular, workflow reading continued through line 2434, carryover through 1216, and voice through 756; a clipped tool response was not counted as complete.

Inspected current Edition 3 main/tests, historical routing, `internal/gradeedition2/ch09.go` and its tests, `grade.Build`, Chapter 8-to-2 composition and directly used CLI/GUI/consumer/workflow/stream launch helpers. Relevant Chapter 9 historical parser, registry, Agent integration tests, CLI skill configuration, GUI entry/settings, and module manifests were also inspected. No broad legacy test campaign was run.

## Present route and smallest supported composition

`cmd/grade-edition3/main.go` admits only Chapters 1, 4, 5, 6, 7 and 8. Its early validation rejects Chapter 9 before `grade.Build` or dispatch. The existing test file exercises invalid/help/prelaunch validation for implemented routes; it is not a complete Chapter 9 route.

A compiled Edition 3 executable was used, not merely a source grep. A separately compiled sentinel writes `sentinel.marker` and exits 73. Direct execution produced exactly `launched`; the marker was saved locally then removed. Both `-ch 9 -json` without a target and a fully explicit Chapter 9 invocation with the executable sentinel exited 2 with usage, no JSON grade, and no marker. `sentinel-result.json` and raw outputs establish the marker self-check and no-launch observation. An earlier directory-target rejection is retained but alone was insufficient for the no-launch claim.

Historical CLI routing in `cmd/grade-edition2/main.go` uses the same positional `path` for `grade.Build(path)` and `Chapter9(..., path, ...)`. That conflates executable target and core source. It is unsuitable as the final Edition 3 interface for a valid library-root/nested-CLI student shape. **Do not flatten that architecture.** The exported API already separates the identities:

```go
Chapter9(bin, core, gui, guiModel, consumer, workflow,
         exercise, streamModel string, example ConsumerCase,
         media map[string]string) (grade.Report, error)
```

The real positive used a byte-identical 75-file copy of `solutions/ch09` at `/tmp/ch09-entry-audit/positive`. No shim or substitute implementation was supplied. Exact arguments are emitted by the corrected caller in `historical-positive.stderr`:

| Identity | Actual positive artifact/source |
|---|---|
| Core module / CLI build root | `/tmp/ch09-entry-audit/positive`, module `ensemble`; historical root is an executable package |
| CLI executable | `/tmp/ch09-entry-audit/tmp/course-grade-3436511201/submission`, built by unchanged `grade.Build(core)` |
| Optional GUI module | `positive/gui`, module `ensemble-gui-stub`; unchanged build helper detects/builds `./cmd/` |
| GUI assets | `positive/gui/web/gui`, as resolved by existing skills/GUI helpers; not an invented root `web/gui` argument |
| Separate consumer | `positive/ch05`, module `ensemble-consumer`, `replace ensemble => ..` |
| Workflow | `positive/ch06`, module `ensemble-workflow`, `replace ensemble => ..` |
| Stream shape exercise | `positive/ch07`, module `ensemble-stream-exercise`, `replace ensemble => ..` |
| Stream and GUI fixture labels | Explicit `claude-sonnet-5`, the supported historical local fixture labels |
| Consumer case | Existing `ReadConsumerCase("")`: `shout`, text `hello framework`, result `HELLO FRAMEWORK`, diagnostic `shout called` |
| Media map | Existing `ReadMediaModels("")`: Anthropic `claude-sonnet-4-6`, OpenAI `gpt-5.4-mini`, Gemini `gemini-3.8-flash` |

Despite its historical module name, the GUI artifact here is the real supplied GUI implementation, including its actual Agent/registry/dispatch loop and WebSocket server, not a protocol-replaying audit fixture. Model labels are local fake-provider inputs, not live capability evidence.

`Chapter9` first runs `skillsChecks(gui, guiModel)` and then full `Chapter8(...)`; Chapter 8 composes Chapter 7, continuing through the inherited library/client/streaming/tool checks. `grade.Build` honors existing binaries and, for source directories, the historical nested `cmd` layout. Consumer and workflow helpers build their own module binaries; stream helpers build/exercise the supplied stream module. GUI launch uses dynamically selected loopback ports, owned scratch cwd, the fixture server URL and dummy key. No live port 8088, browser, private settings, discovery, external provider, audio, microphone, focus or volume control was used.

## Controls, failures, and raw receipts

Evidence directory: `docs/edition-3-notes/ch09-entry-audit-evidence/`. `control-receipts.json` is the authoritative exact argv/cwd/environment-file/start/duration/exit/output-size/output-SHA256 table. Each invocation also retains its own JSON, stdout and stderr. This table lists every build/control, not just successful ones:

| Invocation | Exit | Seconds | Interpretation |
|---|---:|---:|---|
| `build-entry` | 0 | 1.161 | Current compiled Edition 3 command |
| `build-driver` | 1 | 0.198 | **Audit setup error:** wrong exported API signature; not mutation detection or grader failure |
| `entry-command` | 2 | 0.011 | Initial directory-target rejection; insufficient alone to prove no launch |
| `positive-suite` | 0 | 0.522 | Existing `TestSkillsThroughLoop` and `TestSkillLoadFailures` on unchanged real source |
| `sensitivity` | 1 | 0.475 | Same real implementation, one source mutation; compiled behavioral failure |
| `build-sentinel` | 0 | 0.414 | Compiled marker executable |
| `sentinel-selfcheck` | 73 | 0.006 | Expected marker-producing self-check, not failure of the grader |
| `entry-no-target` | 2 | 0.024 | Current Chapter 9 route rejects missing target |
| `entry-sentinel` | 2 | 0.010 | Full explicit Chapter 9 arguments rejected without launching marker executable |
| `build-driver-corrected` | 0 | 0.897 | Separate exact-API caller; original failed source retained |
| `historical-positive` | 0 | 14.938 | **Complete exported Chapter 9 API: 100/100, all nine reported checks pass** |
| `restored` | 0 | 0.518 | Exact source restoration; original focused skill-loop check passes |

Key commands (full paths/arguments are in the receipts):

```text
cwd=main: go build -o /tmp/ch09-entry-audit/bin/grade-edition3 ./cmd/grade-edition3
cwd=positive: go test -work -v -count=1 -run '^(TestSkillsThroughLoop|TestSkillLoadFailures)$' ./ensemble
cwd=mutant: go test -work -v -count=1 -run '^TestSkillsThroughLoop$' ./ensemble
cwd=main: go build -o /tmp/ch09-entry-audit/bin/driver-corrected ./docs/edition-3-notes/ch09-entry-audit-evidence/driver-corrected
cwd=/tmp/ch09-entry-audit/work: /tmp/ch09-entry-audit/bin/driver-corrected /tmp/ch09-entry-audit/positive
cwd=mutant, after byte restoration: go test -work -v -count=1 -run '^TestSkillsThroughLoop$' ./ensemble
```

The complete positive was run **once**. The failed first driver was never run and was not edited/retried in place. Its compiler error and source remain. The corrected driver delegates directly to the unchanged exported API and does not translate the application protocol. Saved driver sources now use `.go.txt` to keep audit evidence from becoming accidental repository build packages; receipts retain their actual build-time paths.

### Environment and lifetimes

Launch environments are allowlisted, not filtered copies of a credential-bearing shell. Go was `/usr/local/go/bin/go`, version go1.25.4, darwin/arm64. Both env files specify PATH `/usr/local/go/bin:/usr/bin:/bin:/usr/sbin:/sbin`, HOME `/tmp/ch09-entry-audit/home`, TMPDIR `/tmp/ch09-entry-audit/tmp`, existing shared GOCACHE `/Users/bill/Library/Caches/go-build`, GOMODCACHE `/Users/bill/go/pkg/mod`, GOPATH `/Users/bill/go`, GOWORK=off, GOTOOLCHAIN=local, GOPROXY=off, GOSUMDB=off, LANG=C and TZ=UTC. No inherited provider keys/settings were passed. The unchanged grader injects its dummy key and loopback provider URL. The owned home was empty; real `~/.cr/settings.json` was never accessed.

Initial controls used GOFLAGS=`-mod=readonly`; completion controls used GOFLAGS=`-p=1 -mod=readonly`, GOMAXPROCS=2. Builds/tests were sequential. The restored source and focused test command are exact, but concurrency env differs from the first baseline/negative; this is disclosed, not treated as identical environment reproduction. Tests/grader assertions and production timeouts were untouched. Outer guards were 120 seconds for builds, 90 for focused suites, 300 for the full positive, and 15 only for the tiny rejection/self-check processes. No timeout fired and no signal was used. Every launched process was synchronously waited; final process inspection is in `process-cleanup.json`. No owned job remains. Disk stayed about 3.6 GiB available, above the required 1 GiB floor; shared caches were reused, not removed.

## One compiling sensitivity and exact restoration

The sole mutation is `mutation.diff`: in the real registry's `MarkUnload`, delete dynamically acquired tools immediately after setting the pending marker. This violates the already-public pending-unload requirement without introducing a replacement implementation or changing tests. Only `internal/skills/skills.go` differs; original sources were never modified.

- Baseline real `TestSkillsThroughLoop` passed (with `TestSkillLoadFailures` in the same preliminary suite invocation).
- Negative compiled and ran the same focused loop. It failed at `skills_test.go:122`: `pending unload revoked dynamic capability`. The emitted JSON behavior does not fail at build/setup; the test observes declarations after the actual Agent loop has marked unload.
- Restored the mutant file from the byte-identical baseline. All 75 files match baseline again. Exact focused command passed.
- Focused sensitivity proves this premature-removal variant is detected. It is not a new campaign and does not establish every possible hidden-dispatch/provenance bug is killed.

Canonical full-source map hashes (`source-provenance.json` defines canonicalization): baseline/candidate and restored `51b7ee4bb380ad29189c224e21a1f3134623733798bb7b3bb9f31e64da4694ec`; negative `465612597fab247cc096fca1de86979125f4bf5e5351c34f153111e37dec0464`. The changed file is baseline/restored SHA256 `3bd84fab99e93aa2c4fb4bde49e39fc825d2e2efd3336f033de7f9c7759959f4` (6,864 bytes), negative `7cb5247ba11d85f4675929073550d129f1e23c9e7a3eeb6dabbf9583916d5513` (6,990 bytes). Raw maps and source archive are local only. Actual test binaries were preserved locally with hashes in `test-artifacts.json`.

## Public requirement → evidence → remaining review

The public contract is Chapter 9, especially TL;DR/steps (lines 25–177), format (179–210), Agent-owned registry (212–265), discovery (267–281), variables (283–320), provenance (322–346), and constitution/exercise (348–448). Names below are existing checks, not added criteria.

| Public requirement | Existing evidence | Remaining candidate source/integration/manual obligation |
|---|---|---|
| Flat frontmatter, valid types and list forms, body including later horizontal rule | `system-prompt` observes PRIMARY_BODY/BODY_AFTER_RULE; fixtures use scalar and YAML-list forms | Read parser/error paths for missing/unterminated frontmatter, names/types/duplicates and malformed inputs; current numeric checks are not exhaustive parser validation |
| Per-Agent registry and required Agent parent; common interfaces; package/composition direction | Existing real skill-loop suite asserts ownership graph and separate Agent isolation; inherited parity covers earlier public composition | Review candidate common/skills/Agent wiring, root parent and import graph; no registry singleton/process-global identity; no demand for a new exported concrete SkillRegistry API |
| Optional skills construction, earlier tool-registration behavior still valid without it | Existing AgentSpec `SkillsDir`, `PrimarySkill`, `Variables` path used by real suite; inherited external consumers exercise earlier construction | Check candidate library-root/nested-CLI use and default/no-skills path; CLI environment is only construction convenience, not a mandated library configuration mechanism |
| Primary/dependency declarations; lifecycle tools always available; guessed hidden name cannot execute | `initial-tools` checks exact declarations, correlated failed hidden result and absent hidden.txt; real suite independently attempts hidden dispatch | Verify dispatch and renderer use the same enabled set, including alternate integration paths; UI acknowledgement is not authority |
| Full ensemble primary declares installed tools | `ensemble-tools`, fresh actual Agent/process configuration | Preserve all twelve taught tools; do not add the old undefined `think` requirement |
| Frozen primary constitution; loading returns real instructions | `system-prompt` compares all request system content across loads; `load-skill` observes correlated next-request tool result | Inspect source ownership/freezing and prohibit mutation of historical system/tool-result text |
| Load grants executable dynamic capability; unload only marks pending | `load-skill` checks actual enabled.txt and pending.txt bytes, successful correlated calls, and declarations after unload; focused mutation fails | Validate unload errors (primary/unloaded), pending state and shared capabilities; no removal/compaction implementation belongs in this chapter |
| One-level progressive discovery; discovered is not enabled | `progressive-disclosure` delimited fixture expansions distinguish current and future sets | Review union across all loaded skills, discovered/loadable filter, exclusion of loaded/dependency-only types; do not mandate fixture formatting |
| Recursive dependencies bring tools and instructions; cycle/missing errors must be atomic | `depends-autoload` observes dependency tool and HELPER_BODY; real `TestSkillLoadFailures` covers cycle/missing/wrong-type/atomic scenarios | Candidate focused tests and source review must check entire graph before any activation and avoid partial grants; numeric score alone does not prove this |
| Blocked/unreachable skill is error with no grant | `blocked-skill` checks result error and unchanged declarations | Check malformed/unknown/repeated inputs and traversal boundaries in candidate source |
| Creation-time and load-time variable values; `$TOOLS`, `$SKILLS`, custom callback registration | `var-substitution` checks expansions against enabled state and custom text only at load; real loop tests update custom callback for future rendering and preserve primary | Review `RegisterVar` path, callbacks outside registry locks, unknown-placeholder preservation and historical loaded-result freezing; do not invent an additional mutable construction API |
| Tool provenance distinguishes initial vs dynamic, actual enabling event sequence, original initial identity retained | Existing real loop suite inspects actual tool-call EventSeq and initial/shared-tool metadata | This is not independently weighted by numeric grader; candidate source must keep actual event identity and avoid separate skill counter or overwriting first-grant provenance |
| Prior chapters remain intact | Complete `ch8-parity` passed (existing nested composition) | Do not infer acceptance of student architecture, provider behavior, lifecycle, source cleanliness or every integration path from the folded parity result |
| Human/live exercise | Current audit deliberately uses real executable plus local fake HTTP/WebSocket only | Later authorized candidate review may inspect supplied exercise evidence; no new provider/browser/audio campaign, caching guarantee, or numeric coverage target follows from this audit |

The historical grader test file includes protocol fixtures checking the grader itself. Those were inspected, not substituted for the real-source positive or mutation here. The eight skills weights plus parity remain exactly the existing 100-point contract.

## Artifact identity and limits

CLI SHA256 `e73c4cad99a32a7f8769d524d52eccbac4a825279bb715d43aff0bb0d9af4c55`, 9,245,202 bytes. GUI SHA256 `c4c7aa55dd42f769b8fd24645ddc68c1ecdf8634c0f4537f6ca7601aa32da053`, 10,234,706 bytes; the observer captured three independently named GUI build outputs with this same hash. `generated-artifacts.json` records actual transient paths, retained local paths, byte hashes and `go version -m` output. CLI and GUI are distinct artifacts built from the stated source snapshot; no global invocation replacement was used.

The observation loop watched generated `submission` files only. Consumer/workflow/stream helpers use other transient executable names and their exact binary hashes/paths were **not captured** before cleanup. Their module/source map and unchanged build/launch code are pinned, and the complete API result includes their inherited checks; do not describe the binary manifest as exhaustive. No second full run was done to improve this receipt. Likewise, dynamically assigned endpoint port numbers and internal raw provider exchanges were not separately retained; the unchanged grader assertions and complete report are the runtime evidence.

Other preserved audit issues:

1. First audit caller used an invented eleven-argument signature and nonexistent static-argument slot. Compiler output rejected it. This was an auditor setup error, not a source/grader mismatch. The separate corrected caller uses the actual ten-argument signature.
2. My initial milestone incorrectly called controls complete after focused tests. Re-reading the exact task caught the missing full API positive, sentinel and restored pass; visible correction preceded completing these required steps. No passing focused suite is represented as the full grader.
3. Initial source/archive and test binaries were mistakenly placed in portable evidence. They were moved intact to owned `/tmp/ch09-entry-audit/local-artifacts`; final portable evidence is UTF-8/no-NUL text only. Earlier receipts with old artifact paths are preserved locally; current receipts record relocation.
4. Several exploratory reads guessed nonexistent paths (including old fixture filenames and `internal/skills/skills_test.go`) and some source inspection output clipped. These are inspection/setup failures, not grader failures. Required prerequisite clipping was recovered as recorded; actual helper defaults and source paths were subsequently used.
5. Some early Git inspections omitted the required explicit `--no-pager`; they were read-only and did not change repositories. Final Git checks use it. This procedural deviation is not concealed.

## Next action and completion boundary

The finite entry obstacle is the public Edition 3 route, not absence of a usable skills grader. A routine thin adapter should extend existing explicit core/consumer/workflow/stream/GUI validation and dispatch to Chapter 9 while retaining separate target/core identities and requiring explicit model fixtures. It should reject invalid input before launch and have targeted routing/validation tests. No change to numeric weights, grader claims, student architecture, or public contract is indicated by this evidence. No fresh Bill permission for routine adapter work or student API-design preparation is recommended; root controls the next authorized phase.

Known source/integration gaps above belong to later full-candidate review, not an expanding entry campaign. This audit itself does not authorize student exposure to private grader/source material or declare the chapter accepted. No live/provider or unrelated lifetime campaign is requested.

Final receipts: `pins-after.json`, `process-cleanup.json`, `disk-after.txt`, `control-receipts.json`, `source-provenance.json`, `raw-local-artifacts.json`, `completion-tail.txt`, and the sealed `inventory.json`. The inventory excludes itself, includes this external report, and validates every included file as UTF-8 without NUL. Raw source maps, binary/source archive artifacts and NUL status snapshot remain local and are hash-indexed separately. Protected original files and student refs are unchanged, and all owned jobs/processes have ended.
