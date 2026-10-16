# Chapter 1 — student execution review

## Candidate and independent input boundary

This is my student attempt, **not an acceptance or reviewer verdict**. Independent code review comes next; I did not start Chapter 2.

- Assigned chapter: `book/edition-3/chapter-01.md`, SHA256 `ab016baaff3f6a9d83ecffcbfc3772cd692c9d33c06d9fa1c52ad9a8708d25df`.
- Book and unchanged grader revision: `7f3d069420112c2136bc95d1d551e12fad6964c1`.
- Empty student baseline: `69859da3328c3f710f8143203b99f934fed778a1`.
- Submitted student commit: **`d2baff64055ac33702a91c6d6ee5a4246e84c0c9`**, in `/Users/bill/projects/ensemble-edition-3-student`.
- Implementation commit: `60bf58a7728a96e27c2e13eea1afbf744889cb6e`; final follow-up adds rejection of nil/zero roots, its test, and comments.
- Go: `go1.25.4 darwin/arm64`. Final executable SHA256: `b322897974592bb714fdd8dffb6a32c4d00382117fc6a14930a3006ae509a558`.

I read the assigned chapter, root instructions and both mandatory skills, execution procedure and public student brief, recovering truncated ranges. I did not delegate, inspect solutions, historical carryovers, private review/readiness records, grader implementation or future chapters, or follow answer-bearing links. The parent supplied only the two launch clarifications below and acknowledged my live observations. No book/grader/application source, settings or runtime was changed. Main-repository writes are limited to this review and its evidence directory.

## What I built and learned

The public `ensemble` package assembles the ownership chain. Ensemble retains Agents and a logger; Agent owns configuration, committed history and its Engine; Engine owns its Agent parent, HTTP client and cumulative usage. `internal/common` contains vocabulary/interfaces; Agent and Engine implementations do not import each other. Engine actually reaches diagnostics through `Parent().Parent().Logger()`.

Clients construct `ensemble.New(logger)`, call `NewAgent(Config)`, then `Ask(ctx, text)`; history/configuration/usage accessors return snapshots or values. Agent stages a user message and commits both messages only after a usable answer. Engine counts a decoded response's measured usage even when the answer is empty. Missing counters are errors, while measured zeros are valid. The CLI only reads configuration/input, calls that API and formats output; default JSON-lines and opt-in `chat` are separate.

This made the chapter's central mechanism concrete: I own the whole conversation, replay it on each request, and pay again for that replay. There is no SDK, retry, streaming, persistence, tool loop or future-provider machinery. Requests use caller cancellation without a model-work timeout.

## Commands and results

From the student repository:

```sh
go vet ./...
go test -race -count=1 -v ./...
go build -o /tmp/e3-ch01.VimxYH/ch01-final ./cmd/ch01
```

Final vet and tests exited 0. The race suite took 125.316 seconds, including real 62-second waits before headers and during the body. External-package tests use the actual public root against local HTTP servers; built-CLI tests execute the binary. Coverage includes exact wire/configuration, owner links/services, full history and snapshots, all text blocks, exact usage, missing versus zero usage, HTTP/JSON/empty-answer failures, cancellation in both phases, Agent isolation, invalid roots, >64 KiB input, final input without newline, protocol purity, chat, EOF and nonzero failure exits.

The initially assigned directory command, run from `/Users/bill/projects/ensemble`, failed before checks:

```sh
go run ./cmd/grade-edition3 -ch 1 -json /Users/bill/projects/ensemble-edition-3-student
```

Observed: execution of the generated `submission` artifact failed with `permission denied`; `go run` exited 1 and reported grader exit status 2. I asked for a public build contract instead of inspecting internals. The parent clarified that the directory adapter does not discover `cmd/ch01`, but the same public CLI accepts a prebuilt executable; root `package main` was not an added requirement. I retained the library layout and ran, from the main repository:

```sh
go run ./cmd/grade-edition3 -ch 1 -json /tmp/e3-ch01.VimxYH/ch01-final
```

**Final result: exit 0, `passed: true`, 100/100, all seven scored checks and both zero-point lifetime gates passed.** The grader observed five requests for five turns and exact totals input=476/output=65. An earlier executable-target run at `60bf58a` also passed and is labeled separately.

## Sensitivity and honest corrections

Nine temporary production mutations at `60bf58a` each produced a specific failing assertion, not a build failure: first-text-block-only, overwriting usage, omitting previous history, omitting root registration, discarding measured usage on an empty answer, treating missing usage as zero, stdout debug pollution, ignoring caller cancellation, and a 60-second HTTP-client timeout. `sensitivity.txt` retains exact replacements, targeted `go test -count=1 -v -run … .` commands, outputs and exits. All mutations were restored exactly before final work. The final follow-up changes only the public invalid-root guard and comments/tests; the full final suite was rerun afterward.

My first draft mistakenly decoded usage into plain integers. Re-reading the explicit missing-versus-zero contract caught that omission before submission; pointer counters and regression cases now distinguish it. I also noticed and fixed the nil/zero-root construction edge case. A Git pager interrupted one inspection command; I closed only that pager and used `--no-pager` afterward. Ordinary baseline tests did not fail; recorded test failures are intentional mutations.

## Separate live exercise

The supplied `anthropic-primary.key` path was absent: initial discovery stopped with `FileNotFoundError`, before any API call. I did not search other credential locations. The parent corrected authorization to the mandatory skill's exact settings file and `directClaudeAPIKey` field. A local process selected that field without printing it, injected it only into the intended CLI environment, and sanitized transcript output. `GET https://api.anthropic.com/v1/models` succeeded and listed `claude-sonnet-5`, which I used. No key was placed in arguments, prompts, evidence or commits; a final in-process comparison found the credential value absent from tracked student files and evidence.

I drove the actual built `chat` interface through a PTY, not a helper API or fake server. On the final commit, I introduced **Juniper Arc**, asked why sailboats heel, then asked the boat's name: the third reply correctly recalled it. A fourth turn asked about persistence. Per-turn input counts were **81, 138, 203, 251**; cumulative input grew **81 → 219 → 422 → 673**, and cumulative output ended at **313**. EOF exited 0. The model described conversation-only memory, consistent with my code, but its self-description is not independent proof of provider retention practices.

For the prompt experiment, an isolated copy of `60bf58a` changed only the chapter's system string to terse, faintly annoyed, and fond of Waywardgeek. Two live replies were one sentence, mentioned Waywardgeek and sounded mildly annoyed. This demonstrates local prompt influence, not improved coding quality or a general causal claim. That variant was not committed. I shared these observations with the parent, who acknowledged them; no outside friend independently observed the session.

## Teaching clarifications I would propose

1. Document both the directory builder's recognized layouts and the prebuilt-executable invocation beside the grading command. A legitimate public library plus `cmd/ch01` otherwise fails as a misleading permission error.
2. Add a tiny state table: successful answer commits history and usage; empty answer preserves measured usage only; missing usage is not measured zero. The prose is explicit, but my initial decoder still missed the distinction.
3. Show one minimal public-root construction/API sketch alongside the ownership diagram, without supplying an implementation. Choosing the public facade and wiring direction took more interpretation than the HTTP exchange.
4. Label live token lines as per-turn versus cumulative and keep the prompt experiment's anecdotal status explicit. Actual counts and a changed style were observable; persistent learning and better code were not.

Evidence is under `docs/edition-3-notes/ch01-student-evidence/`; its README indexes final versus intermediate artifacts. Student source/tests are committed and the repository is clean, with no binary or credential. Remaining disposition: independent review and owner acceptance, not further student implementation.
