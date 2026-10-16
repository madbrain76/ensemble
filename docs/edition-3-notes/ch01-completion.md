# Edition 3 Chapter 1 — complete

Chapter 1 completed the student-first loop on 2026-10-16. Code and text are accepted, with no remaining finding or exception. Chapter 2 has not started: Bill requested finishing Chapter 1 before returning to the separate GUI pause issue.

## Exact checkpoint

- Book/grader starting revision: `7f3d069420112c2136bc95d1d551e12fad6964c1`.
- Separate student repository: `/Users/bill/projects/ensemble-edition-3-student`.
- Empty starting tree: `69859da3328c3f710f8143203b99f934fed778a1`.
- Accepted source: **`1aad3975842d502b9144b7b9bddfd30ca7309db7`**, immutable annotated tag **`edition-3-ch01-r1`**. Clean tree and tag target verified after acceptance.
- Accepted chapter SHA256: `d12183d54265fdd849293946878b0b8d7131830de2aa702928419795ecb34cd7`.
- [Acceptance manifest](ch01-accepted.json) pins every student source file and the principal submitted review records. [Execution record](ch01-execution.md) preserves launch boundaries, assistance, failures and lifecycle corrections. The copied final reference was not a student starting point and remains unchanged.

## Implementation and acceptance

The fresh student implemented the Chapter 1 raw-HTTP conversation loop, required Agent/Engine ownership, public-root library, CLI/JSON-lines interface, cumulative usage and success-only history behavior. No tools or later-chapter feature was requested.

The student's final initial candidate `d2baff64055ac33702a91c6d6ee5a4246e84c0c9` passed the unchanged Edition 3 grader: **100/100, `passed: true`, exit 0**, including both required 62-second lifetime gates. A real Anthropic CLI conversation demonstrated retained context and increasing cumulative usage. Mutation failures and restored passes are retained, not inferred from a green suite. The independent reviewer performed its own grader, race/vet/build, mutation and live checks; exact run revisions and limits are in [the initial review](ch01-code-review-r1.md) and their evidence directories.

Initial review required one narrow correction: remove a malformed thinking block from a fixture described as a valid text-only success response. The student changed only the fixture and its comment, preserving both text blocks and all assertions. **One of the allowed three code-correction rounds was used.** The [correction review](ch01-code-review-correction-1.md) accepted `1aad397` with zero remaining findings. It independently verified identical production bytes, repeated full race tests including both 62-second cases, vet/build and the affected mutation/restore check. Prior grader/live runs remain attributed to `d2baff6`; they were not relabelled as new runs against the later test commit.

The library-root directory grader failure is retained as negative evidence. The existing grader accepts a prebuilt executable; the student received that public invocation clarification, not solution code or private grader source. The orchestrator also corrected its own erroneous credential-file reference using the previously authorized store; no key was placed in the task/evidence. Both interventions are recorded in the student and execution records. Source/input separation was by fresh context and explicit permissions, not an OS sandbox claim.

## Student-informed text improvement

The [author pass](ch01-author-review.md) changed only “Grade yourself”: qualify supported directory builds and show how to build a nested CLI and pass the executable explicitly from the correct module/repository directories. The clarification is **+91 whole-Markdown words** (+74 by the prose-linter method). It changes no implementation requirement, architecture or grader behavior.

The author rejected unsupported API-sketch/state-table/duplicate-accounting expansions. The test-fixture mistake and orchestration/pause problems did not become manuscript requirements.

The **single independent [text review](ch01-text-review.md) accepted the frozen draft as written, zero findings**. No author revision or second text review was needed. The reviewer independently checked exact diff/hashes, command syntax, source-backed invocation semantics, fences, local link and whitespace; it did not rerun full delayed or live suites for a prose change. Author build/parser/linter checks remain attributed to the author. After that acceptance, the orchestrator aligned the public student brief and shared procedure with the same supported executable invocation and updated stale preparation status; these are synchronization/status edits, not new exercise contracts.

## Closure and limits

All four workers were joined; saved submissions and final turns were verified before closure (student270, code reviewer600, author304, text reviewer117). The accepted student repository is clean. The chapter, source, grader and review identities are pinned separately; no copied reference was substituted for accepted construction source.

No application/runtime patch, installation, restart, push, publication build or Chapter 2 launch occurred. The separate GUI speech/selection pause observation is not part of Chapter 1 teaching or acceptance. Next: discuss/address that issue as Bill requested, then perform the Chapter 2 entry gate against this exact accepted predecessor.

### Final checkpoint checks

All278 original source files and231 copied application files match the baseline manifest; current grader source is unchanged. All34 local file links in the scoped checkpoint resolve. The99-file checkpoint scope excludes unrelated work and the separate pause diagnostic. The full staged whitespace check reports 75 warnings, all classified as preserved evidence formatting: literal patch/embedded-diff prefixes, verbatim PTY CRLF/prompts, a blank final line in the retained failed mutation log, and the trailing tab emitted by `go version -m`. These evidence bytes and reviewed hashes were retained rather than silently normalized. The manuscript, current guidance, acceptance manifest, execution/completion records and text-review report pass the separate whitespace check. Credential checks remain enabled.
