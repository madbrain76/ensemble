# Chapter 6 restart author note

## Bounded manuscript change

Read current AGENTS, coding/Edition 3 author skills, shared lean execution
procedure, voice and chapter-writing guidance, the full current chapter, and
the restart student/code-review notes. Historical workflow archives were waived.
Preserved the opening, voice, reasons, scope and historical evidence. Most
corrected requirements were already explicit; no broad rewrite was warranted.

Four local clarifications add about 210 words:

- Request table: native token-limit parameter names differ. A rejected request
  is distinct from a stream-parser failure.
- Section 6.4: the protocol CLI forwards response-boundary observations even
  when it supplies explicit IDs; tool display chunks carry arriving name
  fragments as well as arguments, without repeating accumulated names.
- Section 6.8: OpenAI's complete tool-result group moves only in its disposable
  request projection; chronological History retains the busy hint's position.

Checked these against accepted `solutions/edition-3/ch06/internal/engine/`
`render.go`, `sse.go`, `engine.go`, `internal/common/actor.go`, and
`internal/cli/{cli,chat}.go`. Added no API, acceptance criterion, process-history
feature, or new live-result claim. Accepted grades/provider evidence is carried,
not rerun. In particular, the capped Gemini answer remains incomplete despite
all 137 native characters being delivered/saved; neither full instruction
compliance nor a typing benchmark follows. Independent live generations agreeing
does not prove response-local shape integrity.

## Prose checks and scope

Ran `go run ./cmd/lintprose` on baseline/final copies both named `chapter-06.md`
under `/tmp/e3-ch06-author/{baseline,final}/`; outputs are
`/tmp/e3-ch06-author/{baseline,final}-lint.txt`. Both exit 0, **with warnings**:
negation forms stay 19 (limit 11); person gap grows from 5,109 to 5,319 words
(limit 1,200). Short closers stay 39, within the respective 42/43 thresholds.
All other checks remain passing. No style overhaul was made to silence inherited
warnings. `git diff --check` passes for the manuscript and this note.

Only this note and `book/edition-3/chapter-06.md` were written. No source,
grader, other chapter, guidance, tests, provider calls, builds, children, or
commit; the existing prose linter was the only Go execution. No formal text
review commissioned: the root retains that single review.

## Sole response to the formal text review

Read `ch06-restart-text-review.md` and addressed its one attribution/representation
condition. Relative to the first author version, the response changes 44 lines
in and 32 out (about 90 net words), localized to the TL;DR and §§6.3–6.6.
Retained the reference sketches, labeled their API/callback/taxonomy/capability
claims consistently, and stated the restart's actual `Render(Context)` and
`Parse(*http.Response, ParseCallbacks) ([]Event, error)` signatures. Replaced
the boolean-forces-loss/fabrication and bitmask-minimum claim with the accepted
SSE/plain capability plus native-parser alternative. Checked declarations in
`internal/common/{common,actor}.go`, capability configuration in `config.go`,
and complete Gemini tool-object emission in `internal/engine/sse.go`.

The four earlier additions remain intact. Parent ownership, one parsing entry
point, shared final construction, incremental notifications, response-local
identity, authoritative finals, saving, CLI boundaries/tool chunks, and
plain/unknown fallback requirements remain intact. No source or other chapter
changed; no new API, operational claim, test campaign or second review.

Because manuscript input changed, reran the existing prose linter on
`/tmp/e3-ch06-author/revised-final/chapter-06.md`; output is
`/tmp/e3-ch06-author/revised-final-lint.txt`. Preserved baseline and earlier-final
copies/logs. Exit 0 **with warnings**: negations 17 (baseline/earlier final 19;
limit 11), person gap 5,409 (baseline 5,109; earlier final 5,319; limit 1,200).
Short closers 41, within 43; all other checks pass. Manuscript
`git diff --check` passes. Root directly verifies this response; no second
formal review was commissioned.
