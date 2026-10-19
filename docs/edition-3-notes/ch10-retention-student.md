# Chapter 10 retention amendment — student submission

**UNACCEPTED candidate; independent root delta review pending.** Zero of three formal
correction rounds used. This is an incremental amendment of accepted Chapter 10,
not a new acceptance or Chapter 11 integration.

## Exact position and inputs

- Worktree: `/Users/bill/projects/ensemble-edition-3-ch10-retention` only.
- Accepted base: `729bc919e734405307b4ac0060e65c40e33652e5`, tree
  `025ef3744b42c8f9daec34e76661edc40583ca6a` (137 members).
- Source-only detached candidate: `b03c28fdb1586f16dfb53849f8c763576ebb39ae`, tree
  `b2f86ba9fd48ff9a923a8a924a9612549aa762c0` (142 members), clean source.
- Task: 8,426 bytes, SHA256
  `088816696501bbaf931c482f33883ea011901017f50a42e59469c49c49e8254f`.
- Public Chapter 10 SHA256:
  `1c2ac0b080e05b53e34152f5669dfc1aff5157c0e06f76d8b4d48fe8079d9992`.
- Read all seven public prerequisites/current chapter; clipped combined reads were
  recovered with separate/bounded reads. `public-read-provenance.json` records exact
  paths, sizes, line counts and hashes. No earlier chapters were additionally needed.
- No paused Chapter 11 sibling, future chapter, historical solution, private grader,
  root report, runtime agent tree or old local evidence was inspected. No children,
  web search, GUI skill, external advice, parent questions or second planning gate.
  Brief file-ownership/test note preceded direct implementation.

## Actual ownership and behavior

`internal/common/retention.go` defines validated default policy (200, full).
Agent owns a detached runtime policy under its own mutex; `Config()` also detaches
its policy pointer. History continues to own canonical facts, snapshot/reducer state,
and archive validation/encoding. `internal/agent/save.go` and the normal-exit actor
branch order cleanup after successful recovery/first save or normal-exit save.
Application-owned `Config.TrimDiagnostics` callbacks name only actual open log writers.
Explicit Save is not an extra diagnostic cleanup boundary.

`diagnostics.go` frames every Write as one JSON record, including multiline payloads.
Retention scans offsets rather than retaining the whole file in memory, selects the
newest complete records in order, syncs and atomically renames a sibling replacement,
and switches the live descriptor under the same mutex as append. Known earlier log
prefixes are recognized; unknown legacy bytes remain one indivisible record. Failed
load/save never invokes cleanup. Controlled replacement failure preserves old bytes
and the usable writer; cleanup errors are surfaced without invalidating recovery.
There is no continuous disk cap and no sweep of unrelated files.

CLI owns api.log; GUI owns debug.log, api.log and gui.log. Both read the retention
keys from the existing cr/settings.json path before Agent startup. GUI partial
settings updates validate and persist before publishing runtime state. Missing keys
use defaults; explicit zero and dialogue_only survive update/reconnect/restart.
No provider credential or transport override is added to this settings representation.

History's full mode preserves all available original Events. Dialogue-only saves
fold the available canonical tail into an authoritative full Context, then write an
ordered `{source_seq, role, text}` projection and no covered original Events.
Human turns become user records; ordinary signed text remains visible; tool results,
tool calls, opaque/private replay parts and trace metadata never enter the projection.
Carried dialogue records survive later saves and mode switches. This does not mutate
live History, fabricate Events or erase private state legitimately kept in Context.
Post-anchor canonical recovery tails are validated intact and applied once; future
sequence allocation respects the maximum frontier. Archive loading never replaces
runtime Config or policy. Invalid representations fail before installation.

Snapshotless rebuild refuses dialogue mode and incomplete canonical streams. The
implementation requires a canonical stream starting at sequence 1 with no gaps (or
an empty initial stream); it does not infer missing events merely from a higher as_of.
Snapshot-backed trimmed archives remain valid. Missing originals cannot be recreated
by switching back to full. `CHAPTER10.md` documents these boundaries and caveats.

## Fresh validation on the frozen candidate

Commands ran sequentially in an env-i shell with owned HOME/TMPDIR, installed Go
1.25.4 darwin/arm64, GOMAXPROCS=2, GOFLAGS='-mod=readonly -p=1', GOTOOLCHAIN=local,
GOPROXY=off, GOSUMDB=off, GOWORK=off and GOENV=off. No inherited provider credentials
or generic LLM_BASE_URL; existing tests use their local fixtures/dummy settings.
Free disk was checked before batches and at natural boundaries; it remained above
14 GiB, far above the 1 GiB stop floor. Shared caches were not removed.

| Module | Exact Go command | Result |
|---|---|---|
| core | `go test -race ./... -skip '^TestRequestLifetime$' -count=1 -timeout=3m` | PASS, exit 0 |
| core | `go vet ./...` | PASS, exit 0 |
| gui | `go test -race ./... -count=1 -timeout=3m` | PASS, exit 0 |
| gui | `go vet ./...` | PASS, exit 0 |

The final-core/final-gui receipts contain exact argv, environment, PWD, start time,
Go version, disk measurement and output. Completion times are explicitly
parent-observed, not precise measured test ends: the attempted wrapper-footer
insertion missed its anchor. Supervised jobs 154 and 162 both completed exit 0;
the report does not hide or cosmetically rerun this receipt defect.

New real-production tests cover >200 multiline logical records, default/override/zero,
ordering and live/concurrent appends; legacy records; real startup and exit saves;
load/save/replacement/cleanup failures; policy validation and detached ownership;
full/dialogue-only save/load, mixed signed/opaque/tool parts, prior dialogue across
restarts/mode switches, recovery-tail apply-once and frontier behavior; malformed
atomic refusal; preserved canonical facts, Context, ephemera and Engine.Render output;
real CLI restarts and GUI startup/exit logs; GUI settings updates, reconnect snapshots,
restart, persisted zero and failed-setting-save atomicity. Existing Chapter 10 and
other offline package regression tests also ran in the final suites. The GUI settings
checks use the existing Hub/protocol peer fixtures rather than a browser or network
WebSocket campaign; existing offline socket regressions remain in the GUI suite.

### One meaningful negative, exactly restored

An owned source copy under `.retention-local/negative` was changed only to return
success without doing positive-limit retention. It compiled and the real
`TestRetentionDiagnosticRecordsAndConcurrentAppends` assertion failed with 251
records instead of 200 (exit 1). No setup-error or timeout credit. Restoring the exact
original diagnostics.go made the same focused test pass (exit 0). All 142 source
member hashes match the frozen candidate after restoration. Mutation diff, outputs,
exit files and full restoration receipt are included. No second mutant was used.

## Failures, scope deviation, assistance and limitations

`development-index.md` links the preserved compile/fixture failures and repairs.
These include mistaken local names and API/helper assumptions, incomplete provenance
fixtures, a wrong CLI save-path expectation, and a GUI invocation whose wrapper
silently ran core instead. The invalid GUI receipt is explicitly not GUI proof.
A new CLI test briefly wrote synthetic settings in the source tree; those bytes were
moved into local evidence and the test isolated with t.Chdir. The fixture is preserved
under a descriptive .txt name. Failed rg/path probes are also disclosed.

A development complete-core run inadvertently included TestRequestLifetime, which
waits twice for 62 seconds; it passed in a 132.717s package run. This was outside the
intended finite follow-up scope. It is recorded, not relabelled as a requested campaign
or fresh acceptance, and explicitly excluded from the final root race suite. There
were no paid/provider/browser/native exercises, historical graders, unchanged
standalone consumer/workflow/stream campaigns, live 8088 listener, installs, pushes,
publication, tags or main-repository commits. Old acceptance receipts remain old.

Diagnostic framing does not sanitize provider payloads. Memory for cleanup is bounded
by one physical input line, not an absolute byte cap. Unknown legacy text is retained
as a whole record. Cross-process writers sharing the same path, hostile filesystem
changes, directory-fsync power-loss guarantees and multi-file transactional cleanup
are not promised. Dialogue-only is not private-data erasure and cannot recover missing
canonical facts. The conservative contiguous-stream rebuild rule is explicit above.
No relevant application executable was built/frozen; Go test binaries were transient.
This is a source-only candidate, not a new-binary consumer or GUI acceptance claim.

**Chapter 14 micro_handoff/journal retention: UNIMPLEMENTED, deferred.** No fake future
handoff test or implementation. No Chapter 11 merge/replay or integration was attempted.
Root owns independent delta acceptance/tagging and any later student-owned integration.

## Teaching clarified by actual work

1. A diagnostic limit must count writer records, not physical lines, and swapping the
   pathname without switching the live descriptor is insufficient.
2. Startup ordering requires a valid recovery boundary before destructive cleanup;
   absent recovery state needs a first successful save, not merely ENOENT handling.
3. Dialogue-only is a provenance-bearing projection, not a special Event variant.
   Current Context may legitimately retain opaque replay state even when its audit
   projection excludes it; mode changes cannot reconstruct lost originals.
4. as_of is a replay/allocation frontier, not a count of missing canonical facts.
5. Omission and explicit zero differ. The existing partial settings/update/persist
   path must preserve that distinction through restart, not just API defaulting.

## Artifacts and cleanup

Evidence directory: `ch10-retention-student-evidence/`. `candidate.diff`, the complete
source manifest with Git/filesystem modes and SHA256, source identity, public input
provenance, raw development/final/negative receipts, wrapper text and local-artifact
hash index are included. No archive, binary, private capture or settings.json basename
is copied into main evidence. Owned source snapshots/scratch remain under the worktree's
.retention-local, indexed rather than deleted. No shared cache/evidence/sibling cleanup.
All started test jobs have been joined; no child agents or owned test processes remain.
The final source is clean at the candidate pin. The self-excluded inventory includes
this report and all evidence members, each verified UTF-8/no NUL, size/mode/SHA256.
This report and evidence are sealed once; root must verify final submission/tail.
