# Edition 3 independent readiness review

## Verdict and scope

**Accepted for the authorized, finite preparation closure. No material findings; no corrective changes requested.** This is not acceptance of a student chapter, a claim that the application is bug-free, or a claim that future graders already exist or pass. No unresolved Bill decision was found within the agreed scope.

This was one independent review against the frozen 27-file manifest below, with base HEAD `2be995ebc12c6a7efb21cee3b9740a758489adda`. The review covered all eleven chapter diffs and relevant owning context, all three new Go files, all five thin roles, the execution/student/readiness/carryover documents, and the three evidence reports plus machine-readable readiness record. The initial and final manifest checks both matched every file's exact size and SHA256.

Before review, the complete current `AGENTS.md`, `.agents/skills/ensemble-coding/SKILL.md`, `docs/agentic-codebook-workflow.md`, `docs/edition-2-carryover.md`, `book/voice.md`, `book/chapter-writing-procedure.md`, and the Edition 3 execution/carryover guides were read in bounded ranges. The oversized carryover read was recovered in smaller ranges. At freeze, the six unchanged governing files were checked byte-for-byte against base HEAD; all changed Edition 3 guides and roles were reread completely. Modern Edition 3 instructions govern superseded historical paths and phase ordering.

Only this note was written in the repository. Grader controls used temporary scratch code and local fake-provider processes, then removed those controls. No original grader, reference, copied application, student code, or reviewed file was edited. No private runtime/authentication files, live provider calls, live browser/MCP, install, restart, publication build, push, or commit were used.

## Source and contract judgment

**Source acceptance: yes, for this preparation scope.** The relocation remains tied to established requirements rather than a new product design:

- **Chapter 1:** retains environment-based startup; prohibits an arbitrary whole-generation deadline without prohibiting explicit cancellation or suitable connection protections. Only the allowed command section changes: the Edition 3 grader, the `passed`/exit-status requirement despite a possible 100-point score, and the approximately two-minute cost. File configuration does not move here.
- **Chapter 3:** opt-in/model-capability ordinary visible-text-before-tools enforcement belongs at dispatch. Whitespace and private reasoning do not count; rejected calls are represented without running the tools. This does not force a policy on incapable models or invent a future provider integration.
- **Chapter 6:** stream order is checked as observed, not inferred from final accumulated text. Existing real-parser tests demonstrate the text-before-first-call boundary and tool-first rejection, including concurrent turns.
- **Chapter 10:** preserves producer identity across delta/final and restored history. The added legacy requirement is distinct, deterministic presentation addressing while retaining historical facts, not a prescribed new ID API or destructive archive rewrite. The documented current-runtime defect remains unfixed; its real actor/persistence/renderer regression is student acceptance work here, not a preflight runtime patch.
- **Chapter 16:** restores the evidenced distinction between `auto` reasoning-summary default and explicit `detailed`, including actual request serialization. It does not introduce a second application-settings loader or treat private reasoning as public explanation.
- **Chapter 20:** base-prompt preservation, root/child role separation, and resolved configuration inheritance remain with real root/child construction. `BaseSystemPrompt`, the child constructor, and the prompt regression tests support that placement.
- **Chapter 21:** the demonstrated owning surface is GUI `--config` and `ApplicationConfig`, not a new mandatory root-CLI feature. Source establishes strict decoding/version checking, ignored provider environment in file mode, restricted GUI overrides, external credentials, private atomic save, save-before-acceptance, and active-versus-next-launch provider settings. Legacy launch remains compatible. Credentials are not exposed as loaded flag defaults.
- **Chapters 22–23:** memory and recall extend the same schema only when their types exist. The source has a separate top-level `RecallSettings`; Chapter 22 cannot finish the whole schema. Existing Chapter 23 already promises a public recall settings patch, so teaching its shared-file integration is not a new retrieval feature. The inspected current application test proves recall save/load/runtime mapping, not the new student's complete live-GUI integration. Live acknowledgment and restart obligations remain future Chapter 23 acceptance work and must actually be demonstrated there.
- **Chapter 25:** historical crossover account, with Bill's first-use account and unverified/deferred status preserved; no stranded construction exercise or new implementation checkpoint. It maps the actual technical owners instead of requiring the old catch-all exercise. **Chapter 26** changes only its boundary to start from the exact accepted Chapter 24 commit; Gateway requirements are otherwise retained.

The source checks included the visible-reasoning implementation and real-parser tests, prompt/fleet construction, provider detailed-summary request test, GUI application configuration/store/startup/settings source and tests, and the public artifact-identity investigation and crossover record. Historical reports were treated as reports, not as freshly rerun application tests. No claim is made that a future student's implementation already satisfies these contracts.

## Workflow and isolation judgment

**Accepted.** The student role begins in the separate empty repository and is restricted to the current exercise, accepted earlier work, public guides, and actual grader feedback. Old solutions, the copied reference, grader internals, later chapter material, and private reviewer evidence are not student inputs. The guide explicitly says to stop if isolated data are unavailable rather than explore sibling directories. Authentication access is not disclosed as a student exercise or a credential-reading instruction.

The other roles are thin entry points to the governing procedure, not competing policy documents. They preserve the student-first/author-last rule and the distinction between code acceptance and later prose work. The source default is **three code-correction rounds after initial review**, counted across agents; a failed limit stops progression. That is distinct from **exactly one post-student text-review pass and one author revision**, with any final verification limited to specified fixes rather than a new critique round. Reviewer-authority exceptions remain explicit and cannot manufacture a passing grader result.

The present user-authorized preparation closure is bounded and does not silently authorize a full-book rewrite before student evidence. Future chapter-specific grader gaps are entry gates, not claims of completed preparation for every chapter. The student's empty starting state and the known Chapter 10 artifact acceptance obligation are stated plainly.

## Grader and KISS judgment

**Grader design accepted; KISS adequate for the demonstrated requirement.** The new CLI calls the Edition 3 wrapper; the wrapper preserves all seven original Chapter 1 checks and their 100-point total, then adds two zero-point blocking checks. `Report.Passed` and process exit status retain the significance that score alone lacks. Original grader/reference paths remain unchanged.

The actual submitted binary runs against a local HTTP provider. One check delays headers for 62 seconds; the other flushes headers then delays the body for 62 seconds. Each requires one complete expected answer, matching usage, and an actual provider request. A shared private helper handles those two phases; short test timings remain private to the focused tests. Cancellation-aware server waits and a distinct bounded watchdog prevent harness hangs without imposing an application-generation deadline.

The tests do not merely inspect a removed timeout literal. They build and run real controlled programs over HTTP: timeout-free positive, whole-request deadline negative, and fabricated-output/no-request negative. The reviewer also independently repeated the full-duration controls against the historical solution and a scratch copy differing only by removal of the 60-second client timeout (plus its unused import). Their contrasting outcomes are below.

The implementation introduces no new dependency or generic configuration/testing framework. Production additions total 212 physical Go lines: 179 nonblank/non-comment-only lines and 22 comment-only lines. Tests add 86 physical lines, 74 nonblank/non-comment-only and 7 comment-only. Method: count lines in each of the three new files, excluding blank lines and lines whose trimmed form starts with `//` for the active-line count. These are preparation-tool additions, not a student/predecessor quality comparison. Comments explain deadline bounds, actual-request evidence, streaming-body behavior, score semantics, and the source-review limit; their value is not judged by a mechanical percentage target. No unresolved bloat or scope finding.

**Important limit:** passing a 62-second control cannot prove the absence of every longer arbitrary deadline. The implementation, chapter command note, and readiness report say so. Source review still checks the transport and request-context policy. This supplemental harness is Chapter 1-only; it does not pretend to implement later Edition 3 graders.

## Independently executed checks

Commands ran from `/Users/bill/projects/ensemble`; Go commands used `GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off` to avoid fetching/installing dependencies.

| Check | Actual result |
|---|---|
| `shasum -a 256 docs/edition-3-notes/readiness-review-manifest.json`; `git rev-parse HEAD`; Python size/SHA256 verification of all entries | Exact supplied manifest digest and base HEAD; all 27 entries match at start and immediately before writing this note |
| Python comparison of the six unchanged governing files with `git show <base>:<path>` | All byte-identical |
| `go test -race ./internal/gradeedition3 ./cmd/grade-edition3 -count=1 -v` | PASS; CLI has no test files. Real-process transport positives, deadline negatives, no-request negative, and unsupported-chapter check passed |
| `go vet ./internal/gradeedition3 ./cmd/grade-edition3` | PASS, exit 0 |
| `gofmt -l` on all three new Go files | No paths returned |
| `go test ./internal/grade ./internal/gradeedition2 ./internal/gradeedition3 ./cmd/grade-edition3 -run '^$' -count=1` | PASS as a compile-only check; no claim that these original package test suites ran |
| `go build -o <temporary>/grade-edition3 ./cmd/grade-edition3`, then compiled CLI argument controls | Missing path, unsupported Chapter 2, and nonexistent submission path each exited 2 |
| Compiled CLI `-ch 1 -json solutions/ch01` (real 62-second controls) | All seven original checks passed; score 100/100; both lifetime checks failed; `passed:false`; exit 1; 120.3 seconds overall. Both diagnostics reported no complete answer/usage after the delayed response, not a missing provider request or watchdog failure |
| Same CLI against temporary `solutions/ch01` copy with only `&http.Client{Timeout: 60 * time.Second}` replaced by `&http.Client{}` and unused `time` import removed | All nine checks passed; score 100/100; `passed:true`; exit 0; 124.3 seconds overall. Each supplemental success explicitly reported the 62-second delay. Scratch controls removed afterward |
| Baseline-manifest hashes and protected-path `git diff --exit-code` | All 278 original source hashes and 231 copied application hashes match; no protected original grader/reference/application modifications |
| Student `git rev-parse HEAD`, `git ls-files`, `git status --porcelain` | Exact empty baseline `69859da3328c3f710f8143203b99f934fed778a1`, zero tracked files, clean working tree |
| Local Markdown link and chapter fence checks | All 35 local file-link occurrences resolve; all eleven chapter fence sets balance; only fenced-example change is Chapter 1's `grade-edition2` to `grade-edition3` command |
| `git --no-pager diff --check` | PASS |

One non-blocking evidence-count discrepancy: counting matched line-start triple-backtick pairs in the frozen eleven chapters yields **115** blocks, not the integration report's 116. All are balanced and the sole changed example is the approved command. This does not change the readiness verdict or justify another prose-review round. External links, anchors, publication rendering, and a full application/browser suite were not tested.

The preparation report's minimal Anthropic access success was read as supplied evidence; this reviewer made no live provider request and did not independently verify credentials. That access check is not student Chapter 1 acceptance. The long-duration controls used only a local fake provider. No student grade, comparative implementation improvement, future grader passage, or corrected runtime artifact behavior is inferred from them.

## Disposition

Proceed with the authorized Edition 3 execution under the current student-first procedure. There are **no material findings requiring the one author revision** and no specified fixes to verify from this review. Do not reopen this bounded pass for optional polish. Preserve future chapter-entry gates and acceptance obligations, and record actual student/test outcomes rather than carrying this preparation verdict forward as their success.

## Exact frozen identity reviewed

Manifest: `docs/edition-3-notes/readiness-review-manifest.json`

SHA256: `2ab9baa51eff8893af101e84b2a90092936fdb7376dd1c5062a51801c78aaf73`

Base commit: `2be995ebc12c6a7efb21cee3b9740a758489adda`

| Reviewed path | Bytes | SHA256 |
|---|---:|---|
| `.agents/skills/edition-3-author/SKILL.md` | 1308 | `fa363e7ba234019267a6e23fff2d48fb354592db7a56a6f83711bcaa3bd47d9d` |
| `.agents/skills/edition-3-code-reviewer/SKILL.md` | 1328 | `0e954c764c0e99a6783a679aad5ad88225c2e6dc4901e7795bd182de2d1936a7` |
| `.agents/skills/edition-3-coder/SKILL.md` | 1324 | `6c02a8e4a160d68a5f0cb7376a96ba820e3c40d76e6e3d93d7b19edd6ab9569d` |
| `.agents/skills/edition-3-orchestrator/SKILL.md` | 1559 | `ee81044b632d3b54030c5cfdf1e0fbfad1dcadb22612b735aebd531e3456c44e` |
| `.agents/skills/edition-3-text-reviewer/SKILL.md` | 1185 | `9e97149a9e18c937452911f385e21ccf52293fe3c98ea2483eeb22ca939c14d5` |
| `book/edition-3/chapter-01.md` | 33110 | `ab016baaff3f6a9d83ecffcbfc3772cd692c9d33c06d9fa1c52ad9a8708d25df` |
| `book/edition-3/chapter-03.md` | 47686 | `2366f74bdb873af57fe95e0c7c0b040187365a8db7efc09f5e63ce1f3c59c801` |
| `book/edition-3/chapter-06.md` | 40891 | `ed007e7e7059400939a6735534d7a953d4575e83230a6f68fc35583e8a17bd55` |
| `book/edition-3/chapter-10.md` | 21655 | `abda8ed3a74c2ee2c64bc2208ea104e0547810d6e06e1d5e772c0bb4ee84b3ed` |
| `book/edition-3/chapter-16.md` | 43443 | `7e20c9116cc5d1f0d7514ca61a8501c342d1ac49b044f4ba283b28e1be0e1f28` |
| `book/edition-3/chapter-20.md` | 60570 | `291038c32369f9e519afea9cb30d43023c17f82fd8945578410b00fc0f21fe2e` |
| `book/edition-3/chapter-21.md` | 40794 | `356d762d00d8a5276f12c3f5c1b357e7c8298e6cb3f7ee5a4bf3220bfeae07cf` |
| `book/edition-3/chapter-22.md` | 63533 | `d471365f586d0c083082933ea73bdc2f5cf0cc089d75885f71409a85c470d54e` |
| `book/edition-3/chapter-23.md` | 40486 | `d0d155215945e7e943287ee49f2a90954b51b97e3198367e73528a30f26d4f43` |
| `book/edition-3/chapter-25.md` | 3770 | `db4581c30727babd43cbbd0c46c9573d2d6db0599ebd293c84d6120f2dea6e78` |
| `book/edition-3/chapter-26.md` | 36859 | `1a1d074520793c01ffd95735a95d00e8e2992f1821ec9a81714ed9c1bf4b2cad` |
| `cmd/grade-edition3/main.go` | 1301 | `08ca96d522cafde0d07d6843afe9df4982b002f478b7ba3ea9882b4974ca9ffc` |
| `docs/edition-3-carryover.md` | 13715 | `6b3968770fc93c0e2fcb34d53617393e52a40bc3315576f550840fb1e2a5755a` |
| `docs/edition-3-execution.md` | 10175 | `bf1bf261a2869a82971959b219e5f7aec30005b9caf4578d8e73429db9302bb0` |
| `docs/edition-3-notes/ch01-readiness.md` | 18960 | `b253340f876740fb43f020b9d7c77fc04e95d0ba8cc423e858fbf97b03d5475c` |
| `docs/edition-3-notes/crossover-closure-audit.md` | 21416 | `18564ed82fd7bf1a8c594ecb34b12485ea79ce888c009d6fd3064830fab9e8e9` |
| `docs/edition-3-notes/crossover-integration.md` | 5043 | `be3f75e4b6b98fa63f9e6e83a964b070758d7f617a39ae1e6efa036be1453d49` |
| `docs/edition-3-notes/readiness-checks.json` | 659 | `544a2a35d4db1948d2d9753d0e936eabd0ef301037364c6c633fdd149ed49e11` |
| `docs/edition-3-readiness.md` | 4209 | `da4b6cd29be9a701193cfec76fd07c160cc6452d7c672e5a50a107036ccdb1a6` |
| `docs/edition-3-student.md` | 3484 | `0fad2374e9f6c53770f2060f0ef5552e437adc27fac7538baa70757e84a74b91` |
| `internal/gradeedition3/ch01.go` | 5798 | `a1716c6b11b0531acad0a1c815450a979d2ae382e343162f379285c24d9cf97e` |
| `internal/gradeedition3/ch01_test.go` | 2721 | `3d229c5b85c2939ecd2af2969028d02e7f3d0076193371b2c56cb0f6b2b66d2d` |

## Post-review commit hygiene (orchestrator)

The staged whitespace check stopped the first commit attempt before any commit was created. I replaced two Markdown hard-break trailing-space pairs in this report with blank lines. No reviewed file, manifest identity, verdict or finding changed. The reviewer-submitted report had SHA256 `94792f6d75a4bf3ecda113ded5b41021eca94617cd864a5ca5e6af0eddf6b3b4`; the committed blob records this formatting normalization and note.
