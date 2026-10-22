# Chapter 17 restart author pass

**HOLD for the single formal text review.** No author response round used;
maximum one bounded response remains. No student task, grader, code snapshot,
or live operation belongs to this editorial pass. Accepted Chapter 16
(`5110e28ae0d8d53d41ee84280dcd87e5a33c9ec6`) remains the code predecessor.

Read the full current AGENTS, edition-3-author skill and its mandatory workflow,
Edition 2 carryover, voice and writing-procedure documents; current Edition 3
execution guide/carryover; and all 745 original chapter lines. Clipped sections
were recovered. Current lean/student-first scope overrides dated campaign rules.
Read the correct historical original20/Edition2 Chapter19 notes, `ch20.md` and
`ch20-code-review.md`, not the old recall chapter's `ch17*` notes.

## Minimal clarification

- Added a short restart boundary: accepted Chapter 16 remains the predecessor;
  this chapter creates no implementation task, grader, or snapshot, and reports
  no fresh live crossover or private-memory migration.
- Labeled the existing “Try the same task” shell walkthrough as a historical
  reproduction using the second-edition `solutions/ch17/` reference, not a blind
  Edition 3 assignment or proof of restarted-main capability. Fresh live access
  requires separate authorization; profile/sign-in references belong to that
  historical implementation. Commands and recorded prompt are unchanged.
- Preserved Bill's opening, quotations, technical narrative, failures and caveats.
  October 11 remains a five-request, 27.46-second CSV/Decimal repair with six
  unchanged tests, ordinary progress, zero reported reasoning/cache counts and
  no summary deltas: not GUI/steering evidence or Edition 3 authorship.

## Checks

Read the complete `cmd/lintprose/main.go` and the Makefile lint target. The saved
baseline matched the original chapter (`cmp -s`, exit 0). Both commands ran from
`/Users/bill/projects/ensemble`, using the same basename `chapter-17.md`:

```sh
GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go run ./cmd/lintprose /private/tmp/ch17-author-baseline/chapter-17.md
GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go run ./cmd/lintprose book/edition-3/chapter-17.md
```

Both exited 0. The existing soft `person gap (after L3)` warning remains:
5,038 baseline / 5,082 final words against 1,200. No new warning category;
all other checks pass. Outputs: `/private/tmp/ch17-author-checks.UZqMa0/`
(`baseline-lint.txt`, `final-lint.txt`). Scoped `git --no-pager diff --check`
exited 0. No unrelated style cleanup or implementation tests were added.

Only this note and `book/edition-3/chapter-17.md` were edited. No code, guides,
graders, snapshots, credentials, raw auth/save/wire evidence, runtime or private
memory were changed or operated. No children, provider/browser/native calls,
installation, commit or push. Bill's live-OAuth deferral remains a deferral,
not fresh evidence or transferable authority. Root review and formal text review
remain pending; maintain steady HOLD.
