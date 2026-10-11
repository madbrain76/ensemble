# Original23 / Edition2ch21: independent review

Reviewer `reviewer_astra`, GPT-6 Astra/high, October11,2026.
**Preparation only. Zero formal correction rounds; no acceptance claimed.**
Mandatory `.agents/skills/ensemble-coding/SKILL.md` read in full; governing
carryover, workflow, AGENTS and approved chapter map apply.

Reread all808 lines of original `book/chapter-23.md`, including examples,
revision `e509c2ddc237f799844a6efe6ebdb744289ce295`, SHA-256
`3b19960c74cf4d01e225b9c60bc175b199cc95ec01cad754b97d9a3dcf4bd061`.
Read [student-safe notes](ch23.md). Accepted predecessor is original21+22,
main tree `a801628b1dd631e5ea20ad03ea7d255584728d31`, exported at `27acbed`
and immutable `edition-2-ch22-r1`. Review is for sandbox capabilities only;
original24 reuse and original25 fleet implementation remain later assignments.
No original23 answer was inspected before student implementation.

The previously supplied launch/public-input erratum remains byte-current:
`internal/grade/ch23_harness.go` SHA-256
`24b7e13a90868cdf7b1f2fd1a3f4bf710f9b5ebffc84594bcfd5bb2c9d3cf80c`.
The original grader plants a library probe in its supplied tree, so coordinator
will grade an isolated exact copy, preserving authoritative source. Its library
package-name and terminal-envelope assumptions require actual-result assessment;
no exception is preaccepted. Missing interface facts were supplied narrowly,
without grader internals or answer code to the student.

Earlier availability preflight established executable `/usr/bin/sandbox-exec`
on Darwin and exit0 from an allow-default `/usr/bin/true` probe. This proves
availability only, not confinement. No provider or original-grader call is made
by this reviewer. Credential tests should use fake markers; existing profile
secrets remain outside reviewer access.

Final review will assess both in-process file operations and real spawned
commands/descendants, with legitimate reads/writes/toolchain work alongside
escape/network denials. It will trace credential environment/log boundaries,
actual registry removal and web-skill withholding, unconditional child narrowing
separately from its diagnostic, host-path-safe failures, actual ownership and
unsupported-platform refusal. Those are the original chapter's claims, not an
expanded security certification or a demand for future spawning features.
Source, own-predecessor/reference comparisons, student explanations, unchanged
grading, live use and targeted mutation evidence remain pending frozen handoff.


## First-green source preflight: actual cache-boundary defect

Read frozen `fcb98a9459b328314f0a7ff06efbb5e7c85fe07e`, main tree
`f2f8ff00229a4e198988b8756c8702b42b34ab1c`, the student's implementation
checkpoint, and matching first-edition sandbox source after implementation was
available. Final review/live/mutations remain pending. Zero formal rounds used.

The original grader's actual result is **0/100 FAIL**,0.623s, without timeout.
Six behavioral sessions stop at `provider model is not set`: its launch supplies
no model. The child-spec probe imports the module root, which is a command,
instead of the public `ensemble/` package. These are independently confirmed
fixture incompatibilities, not proof that sandbox behavior passed. No exception
is decided before the complete student review and supplemental evidence.

The student `TestSandbox*` suite passes independently with `-race` from an exact
frozen archive (exit0,1.252s), including real OS commands and host-network control.
An additional focused probe nevertheless found a real current-scope defect:
Seatbelt's xcrun cache rule combines the Darwin temp subpath with a basename-only
regex. It admits matching names at arbitrary nested depth, not only the actual
cache location. A fake-marker file named `xcrun_db-private-marker` in an unrelated
sibling `t.TempDir()` was readable through the real Agent, job and `run_command`;
`cat` returned `REVIEWER_OUTSIDE_CACHE_MARKER` and exit0. No private file, profile,
credential or provider was used. The chapter explicitly warns against shared
temporary-directory access; this is a boundary defect rather than a new policy.

Reproduction and output are retained as reviewer-owned
[probe](../../solutions/edition-2/evidence/ch23/reviewer-cache-boundary-test.go.txt)
and [failure](../../solutions/edition-2/evidence/ch23/reviewer-cache-boundary.txt).
Coordinator received the finding and relayed the required behavior: narrow the
allowance to actual toolchain cache paths, reject unrelated nested matching
files, and retain Python/Git allowed-work controls. No original answer code or
implementation recipe was supplied. Student main was not edited. This is
attributed pre-submission assistance, not a completed formal correction round.

A second real-runtime probe found a diagnostic leak at the same frozen source.
`run_command` resolves a relative `cwd` to an absolute host path before the
sandbox sees the request. With caller `cwd: "escape"` and a local symlink of that
name pointing outside, refusal correctly prevents execution but exposes the
canonical sandbox root in both error and tool result. The chapter requires the
caller's own spelling. Coordinator received this narrow finding plus the need
for a legitimate-cwd positive control. Preserved
[probe](../../solutions/edition-2/evidence/ch23/reviewer-cwd-refusal-test.go.txt)
and [failure](../../solutions/edition-2/evidence/ch23/reviewer-cwd-refusal.txt)
use fake paths only. The coordinator reports that the student independently
identified the same cwd leak at this point; retain both discoveries rather than
crediting only reviewer assistance. Still pre-submission, zero formal rounds.
