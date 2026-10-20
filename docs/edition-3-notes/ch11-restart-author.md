# Chapter 11 restart author

Status: **saved FINAL / HOLD for root's one formal text reviewer**. No formal
review performed here; at most one authorized response remains.

## Bounded changes

Only `book/edition-3/chapter-11.md` and this new note were authored. Required
current guidance, full chapter and both restart notes were read; truncated
output was recovered. Public invocation reference: `docs/edition-3-student.md`.

Exact diff: `/tmp/ch11-restart-author-lint/chapter.diff`, against
`/tmp/ch11-author-before.mfE3eY/chapter-11.md`. Five hunks:

- TL;DR: point to the current Edition 3 invocation in the exercise.
- Ephemeral Tools: teach once that every provider disclosure surface excludes
  automatic tools, including Chapter 9's deferred catalogue and skill-loading
  `tool_reference` blocks. Round 2 exposed the full-capability projection bypass;
  dispatch refusal alone was insufficient. Retain complete capability facts for
  automatic collection, ordinary discovery and canonical history.
- Exercise build: use explicit `./cmd/ch01` instead of the library root and show
  the separate GUI module build at `./cmd`; state the working directory.
- Exercise grade: use the public `grade-edition3` command with explicit module
  and model flags before the required executable; Chapter 11 then Chapter 1 only
  after success on the same unchanged CLI; genuine 62-second header/body waits.
  Remove obsolete omitted-target advice; retain optional case/media flags.
- Taking It for a Spin: correct the same executable build target.

NO CHANGE to already-clear public API, wire/ID/envelope, grants, same-page
reload, snapshot bounds, automatic lifetime, cancellation and grader-observable
requirements. No implementation prescriptions added to justify grader repairs.
Preface and historical receipt tail were checked byte-identical. No new live,
browser, native-audio or universal-provider claims were added.

## Checks

Root explicitly authorized the cached offline lint helper despite the broader
no-build boundary. Sequential commands (same `chapter-11.md` basename):

```sh
env GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off GOFLAGS='-mod=readonly -p=1' GOMAXPROCS=2 go run ./cmd/lintprose /tmp/ch11-author-before.mfE3eY/chapter-11.md
env GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off GOFLAGS='-mod=readonly -p=1' GOMAXPROCS=2 go run ./cmd/lintprose book/edition-3/chapter-11.md
git --no-pager diff --check -- book/edition-3/chapter-11.md docs/edition-3-notes/ch11-restart-author.md
```

Baseline and final lint: exit 0, no hard findings. Outputs:
`/tmp/ch11-restart-author-lint/baseline.txt` and
`/tmp/ch11-restart-author-lint/final.txt`. Retained soft person-gap finding after
L66: 2246 → 2363 words; no cosmetic rewrite. Whitespace check passed. Disk remained
about 1.4 GiB free, above the 1 GiB floor. No application builds, source/test
campaigns, providers, browser/native operations, children, publication or commits.
