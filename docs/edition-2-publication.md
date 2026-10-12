# Edition 2 publication and layout migration

Bill authorized this sequence on October 12, 2026: finish the Gateway chapter,
archive Edition 1, promote Edition 2, and prepare its first publication. The
crossover has begun; Chapter 25 deliberately remains a placeholder. Producing
Edition 3 with Ensemble remains future work.

## Pending migration

No files have moved yet. Gateway implementation and manuscript/grader enhancement
are accepted, the latter at independent review `63709086`. The move is released.

| Current path | Destination |
| --- | --- |
| `agent/` | `edition-1/agent/` |
| First-edition content in `book/` | `edition-1/book/` |
| First-edition `solutions/chNN/` | `edition-1/solutions/chNN/` |
| `solutions/edition-2/main/` | `agent/` |
| `book/edition-2/` | `book/` |
| Other contents of `solutions/edition-2/` | `solutions/` |
| Original `cmd/grade/` | `edition-1/cmd/grade/` |

Preserve the original grading module and its support under `edition-1/`, with
the original module path, so it can run from that directory. The current revised
graders still use `internal/grade` and the fake providers. Retain that shared
support at the top level as well; moving the old command does not justify a
rewrite of the grading library. Keep `cmd/grade-edition2` as the explicit current
command and make the top-level Makefile select it.

Preserve original build scripts, Makefile and README in the archive. Shared craft
guides, build assets and book tooling also remain available for the current book.
Archive old chapter-specific maintenance scripts; distinguish retired Edition 2
attempt scripts from the current exporter rather than silently making obsolete
checks look current.

The move preserves tracked blob identities for the original edition and frozen
solutions. Bill's dirty application files and local untracked material move with
their directories, without becoming accepted source or being staged. Historical
manifests, scores and receipts retain their source paths and revisions. This map
explains their new locations. There are no symlinks or new chapter numbers.

## Assignments and final checks

- Coordinator: physical moves, preservation check, current README and agent
  instructions, chapter paths, external launcher, final epilogue and assembly.
- Author: preface, front matter, author's note and narrow Chapter 0 corrections;
  preserve the documented voice and distinguish completed work from crossover.
- Coder: current grading defaults, affected test paths and publication tooling.
- Reviewer: independent Gateway acceptance, then relocated entry points,
  preservation and publication prose/artifact checks.

Verify relocated default grading and original archive grading, affected module
checks, runnable example paths and local links. Do not repeat paid demonstrations
for a path-only move. Build the assembled Markdown, EPUB, print interior and
cover-prefixed PDF. Preserve the established EPUB navigation behavior. Inspect
chapter coverage, edition metadata, page references, code layout and images.

The new epilogue is the coordinator's own first-person account from Codex, written
after the manuscript is complete. The archived epilogue's narrator running inside
Ensemble is a different contribution and stays intact. Local publication files
and commits are preparation; no external upload or push has occurred.
