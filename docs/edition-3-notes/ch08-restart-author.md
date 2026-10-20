# Restarted Chapter 8 — bounded post-student author pass

Changed only four existing passages in `book/edition-3/chapter-08.md` after reading the current guidance, complete chapter, current student/code-review notes and relevant carryover map. The accepted architecture, opening, voice, historical explanations and evidence qualifications remain. No new section, framework, acceptance threshold or operational history was added.

- Made the inherited captured-turn budget explanation explicit about the final response's request-local stop-tools instruction: no durable History/signed-content changes or leakage into later turns. Narrow source check: `solutions/edition-3/ch08/internal/agent/actor.go`, `send`.
- Clarified automatic Read aloud cancellation versus explicit Speak (`on demand`), draft preservation, fresh typing surviving an earlier automatic acknowledgement, and explicit any-tab Continue authority. Narrow source check: accepted `gui/web/gui/app.js`, settings application, state/continued handling and speech classification. These distinctions address the demonstrated producer/consumer and ordering mistakes without describing a new protocol framework.
- Replaced the vague higher-rate mapping with the actual native clamp, 0.5–3, while retaining arbitrary fractional control values and the qualification that native rates are not measured speedups. Kept the historical fractional-rate evidence and its earlier-range limit.
- Added one sentence keeping each pane independently scrollable and its follow control clear of transcript text. No visual-quality or observation claim was added.

No tree expansion was needed: the existing section already explicitly requires a data-driven one-node tree and excludes child-Agent lifecycle work. System-theme instructions remain unchanged. Actual dynamic System change-event delivery remains **UNVERIFIED**: changed emulated query values did not establish delivered events. No workaround or dynamic PASS is taught. The truncated visual reply, source-supported overlay correction and verified geometry/click do not become a claim of direct root viewing or hearing.

## Checks

From `/Users/bill/projects/ensemble`, before and after the edit, on the same basename and actual chapter path:

```sh
GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go run ./cmd/lintprose -v book/edition-3/chapter-08.md
```

Both runs passed the hard ceilings (exit 0). Both retain soft warnings: negation forms 11 against allowance 6; person gap after line 11 grows from 3156 to 3253 words. These warnings are reported, not called clean; no invented first-person experience was added to reduce them. Short closers remain 14 of 69 paragraphs, within the ceiling.

`git --no-pager diff --check -- book/edition-3/chapter-08.md docs/edition-3-notes/ch08-restart-author.md` passed. Reviewed the actual four-hunk chapter diff in context. An initial assumed Python-linter path/glob was absent; the existing command record identified `cmd/lintprose`, which was used for both valid runs.

No application tests/builds, browser sessions, credentials/provider calls, publication, staging, commits, code/reference/snapshot/guidance edits, children or evidence bundles. Root owns the subsequent formal text review; this is the initial author submission, not that review or its one permitted response.

## Root prose closure

Root inspected the four-hunk diff, author result and saved independent text review. The single formal review passed with no author revision; its lint counts and source attributions match the actual changes. Both text workers' finals were saved, joined and closed. Code and operational acceptance remain qualified as recorded in the student/code-review notes; this prose pass adds no new execution claim. Chapter 8 is complete for scoped commit. Chapter 9 remains paused for the requested discussion with Bill.
