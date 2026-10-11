# Independent review: Edition 2 solution numbering

Reviewer: `enhancement_reviewer_resumed`, October 10, 2026. Mandatory skill:
`.agents/skills/ensemble-coding/SKILL.md`; full skill and carryover read. Bill's
explicit preservation and renumbering authorization governs this change.
Pre-move baseline: `c2df11b5f86c2abbcc14a827e59d263f766ee8e6`.
Final decision: accepted at initial independent review; no correction round.

## Preservation and mapping

Independently compared each destination's complete file set, Git blob identity
and executable mode with the pre-move commit, rather than relying on the copy
operation or its manifest alone. All comparisons pass:

| Original snapshot | Preserved archive | Reader-facing snapshot | Files |
| --- | --- | --- | ---: |
| `ch04` | `history/ch04` | none: jobs-only historical stage | 25 |
| `ch05` | `history/ch05` | `ch04` | 32 |
| `ch06` | `history/ch06` | `ch05` | 45 |
| `ch07` | `history/ch07` | `ch06` | 59 |
| `ch08` | `history/ch08` | `ch07` | 66 |
| `ch09` | `history/ch09` | `ch08` | 69 |

Paths are under `solutions/edition-2/`. No missing, additional or changed file
was found within any snapshot. Chapters 1–3 and `main` also match their original
file sets, bytes and modes. The old reader-facing `ch09` directory is absent.
Independently recomputed all six sorted relative-path/SHA-256 manifest digests
and checked the original Git tree IDs in `history/mapping.json`; all match.
Nested application names and their `replace ensemble => ..` directives retain
their exact original bytes.

## Executable paths and evidence

Read the entire executable diff: remove the command's five special-case/offset
lines and change one fixture path in each of the Chapter 4–8 integration tests.
The simple default is now `chNN` for Chapter N. No behavioral predicate, runtime
implementation, API, comment contract or test coverage changes. No framework or
new test infrastructure is needed for this path substitution.

Read the coder's raw command statuses, test output, full grade JSON and module
resolution records. All six final grader-file SHA-256 values match the current
files. Affected tests pass (grader package 9.451 seconds), vet is clean, and
`go run ./cmd/grade-edition2 -ch 8 -json`, with no path overrides, scores actual
**100/100 PASS**. All 17 nested modules resolve to their own relocated core.
Independently ran `go list -m -json ensemble` in new `ch04/ch05` and `ch08/gui`;
both resolve correctly. This verifies the smallest and final relocated public
boundaries without duplicating the full behavioral suite.

These are relocation results; historical grades and exceptions remain unchanged.
No paid call, mock, runtime mutation or repeated historical grader was used.

## Documentation and final decision

Read the current solution README, historical archive README, mapping manifest,
workflow/plan notices and all three changed manuscript passages. They distinguish
book numbering from original student numbering and retain the nested consumer,
workflow and streaming application names. The current default commands agree
with the simplified grader. The archive warns that the copied prior README's
relative paths and status describe its original location and date.

Verified the prior README is byte-identical to its pre-move version. Each of
the 23 changed historical notes is also byte-identical after removing its one
new path-era notice. Their original commands, hashes, grades, exceptions and
claims are not rewritten as later execution. Historical run evidence is
unchanged. All 60 local file links in the nine manuscripts and two new READMEs
resolve; scoped whitespace checks pass.

No preservation, mapping, executable-path or historical-reference ambiguity
blocks acceptance. Runtime content and existing behavior are unchanged; the
path changes make the reader's chapter selection direct. The coordinator owns
the commit and any fresh tag. This review authorizes no later chapter work.
