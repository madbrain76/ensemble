# Edition 3 Chapter 9 — frozen local candidate

## Identity and scope

- Predecessor: `edition-3-ch08-r1`, `b8f22ba5b3c9d1fc52d1a114810b90b77928684d`.
- Final candidate: `251ce88fb1677897dbcc6b33661096c5cb39cd8e`; tree `a2add41c306d567c05d98ae760911bc13c18b75a`. Worktree clean.
- Exact tracked source, archives, executable SHA-256s, construction targets and public input pins: `ch09-student-source-manifest.json`.
- Portable receipts: `ch09-student-evidence-manifest.json`. `.ch09-local` retains local executables and both source archives; binaries/archives are deliberately not copied into portable UTF-8 evidence.
- Local work only. No browser, paid provider, credentials, children, branch/tag promotion, push, or acceptance claim. No private grader/audit/history/reference inputs were read. The public Chapter 9 document and its current Edition 3 recipe were the contract.

## Ownership and behavior

The root constructs a per-Agent registry, installs lifecycle declarations, renders the primary, then publishes the Agent. The Agent sees an injected lifecycle interface, not the registry implementation. The registry has immutable discovered definitions and separately locked loaded/enabled/offered/provenance/pending/variable state; it has a parent Agent. No mutable application globals were added. CLI and GUI translate the three chapter environment variables into explicit `SkillSpec` input. No-spec creation preserves inherited registration.

Frontmatter names, not directory names, identify skills. Parsing supports the specified scalar and indented dash-list fields, detects malformed input with paths, and leaves later Markdown rules in the body. Activation validates the complete dependency closure first; dependencies precede dependents. Declaration filtering and pre-control dispatch gating share enabled state. Hidden calls do not allocate jobs or execute effects. The lifecycle uses ordinary paired tool history, with actual `tool_called` sequence numbers and first-enabling sources. Repeated loads preserve those sources. Unload is pending only; dynamic writes remain executable. The primary text stays fixed. Rendering snapshots callbacks/state under lock and invokes callbacks after unlock, preserving unknown variables. Existing engine request rendering and GUI history/wire paths are reused.

## Final-source validation

All commands, exact cwd/environment, start time, duration, exit, stdout/stderr and disk headroom are in individual JSON receipts. Receipt-wrapper processes receive PATH/HOME/TMPDIR plus offline Go settings only, not provider environment. No suite was interrupted or shortened. Guards required at least 1 GiB free; receipts show actual headroom.

- Root `go test -count=1 ./...`: pass, 129.18s. Root race suite: pass, 151.83s. Root vet: pass.
- GUI, consumer, workflow, and stream modules: tests/race/vet all exit zero. Consumer and stream report no test files; their build and adapter checks are the applicable coverage, not invented unit-test counts.
- Five final executable builds passed, with explicit module/target identity and unchanged pre/post-adapter hashes.
- Chapter 9 adapter: 100/100, `passed: true`, exit zero. Chapter 1 adapter: 100/100, `passed: true`, exit zero, same final CLI.
- Focused tests cover parsing, duplicate identities, initial/dynamic dependency behavior, missing/cyclic/unknown/reachable-vs-hidden loads and atomicity, real write/process effects, hidden controls, repeat/provenance, pending dynamic write, Agent isolation, reentrant callbacks, concurrent snapshots, legacy tools, and fixed system/dynamic declarations over actual local HTTP requests.
- The GUI test connects a real local WebSocket through the production hub, validates launch-variable translation, hidden write refusal, dynamic rendered instructions in the normal wire transcript, pending unload, and a real post-unload write. This is local wire integration, **not** browser evidence.
- Two final-source disposable controls are compile-valid behavior negatives: removing dispatch admission/lookup guards makes the hidden writer actually run; overwriting first-source metadata fails the provenance assertion. Each exact-restored copy passes, and every tracked file of both restored copies matches the final source manifest. See `sensitivity-summary-final2.json` and scripts for complete mutations.

## Preserved failures and repairs

1. Initial pre-freeze root compile failed because a scripted field insertion missed the predecessor's spacing. The insertion and slice-based wiring were repaired; the full inherited suite then passed (129.048s reported by Go). Early merged logs remain; their wrapper start/duration was not measured and is not retroactively invented.
2. A metadata `printf` command failed before hashing; corrected immediately. One guessed public brief path was absent; the prescribed path was recovered from the public assignment. A later metadata check guessed a nonexistent predecessor tag; it was corrected to the assignment's annotated edition-3-ch08-r1 and commit ancestry verified. None changed source.
3. First GUI build used nonexistent `./cmd/ensemble-gui`. That exit-one receipt is retained; corrected `./cmd` and all final builds passed.
4. The first dispatch sensitivity mutation caused a nil-job panic, and its second revision did not compile (`RLock` on a Mutex). Both are explicitly uncredited; both were restored and passed. The corrected third control produced the intended real-write failure and restored pass. Final-source controls were then rerun from the second frozen archive.
5. First frozen adapter attempt scored 85/100, `passed: false`, exit one: all skill checks passed, but inherited star-topology parity correctly rejected Agent's concrete sibling registry import. This was an ownership defect, not an accepted exception. Agent now depends on an injected interface; the root still owns concrete construction. The new commit was rebuilt, both adapters rerun, and all local final gates repeated. The first archive, logs and report are preserved and labelled superseded. Its fail-fast command did not run Chapter 1.

## Assistance, interpretation and limits

Assistance was the public operational assignments, root-approved plan, repository AGENTS/coding-skill guidance, current Edition 3 guides, public chapter, and pinned black-box diagnostics. No child/reviewer or private reference input was used. Implementation iterations above are not formal independent-review rounds; that review remains root-owned. The earned design lesson is the concrete sibling-import defect: injection through an interface is necessary even when the registry is per-Agent. The chapter leaves exact exported API spelling and callback error signatures to implementation; the documented `SkillSpec`/read-only `SkillRegistry` API supplies that contract without claiming a new chapter ambiguity. Frontmatter remains deliberately narrower than YAML. A successful finite sensitivity check is not universal proof, and no reproducible-build claim is made; embedded build metadata and exact current binary hashes are retained.

## Browser/provider continuation — NOT EXECUTED

Root must supply its later finite authorization. Use final CLI/GUI hashes from the manifest; do not rebuild silently. Prepare a fresh scratch skill directory with `base` (primary/read_file/offers code-tools), `code-tools` (loadable/write_file/depends helper), and `helper` (dependency/list_directory), matching the public chapter. Launch the GUI from scratch with `EN_SKILLS_DIR`, `EN_PRIMARY_SKILL=base`, and a distinctive `EN_CUSTOM_VAR`; use the frozen GUI executable and root-approved local fixture/model endpoint.

The local-only fixture should attempt a hidden write, load code-tools, write/read a scratch note, unload, then **write again using the dynamic write_file capability**. Capture clean request bodies showing initial versus expanded declarations and identical system text. The real browser should connect through the production GUI, visibly show the skill instructions and pending result in ordinary tool panels, inspect the real note content, and exercise pause/Stop, a later ask, and shutdown with a bounded delayed fixture where needed. Check responsive layout/console through the approved browser tools. Keep screenshots and sanitized console/network evidence distinct from local Go/WebSocket tests. No dummy browser or provider success is claimed here.

Paid CLI provider work, if authorized, is limited to the public chapter graph and a finite root-selected matrix. Current declared identities used locally are claude-sonnet-4-6, gpt-4.1-2025-04-14 and gemini-3.8-flash; root must confirm account availability and budget. Inspect the real note and post-unload dynamic write, retain credential-free request structure and exact exits, and terminate owned processes. No credentials have been opened or calls made by this student.

## Status

Frozen local implementation/evidence complete. No genuine implementation blocker remains. Manual browser/provider validation and acceptance/promotion are intentionally pending root action; there is no new ownership plan required.
