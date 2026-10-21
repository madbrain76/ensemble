# Chapter 12 restart — formal post-student text review

## Verdict: REQUIRED_REVISION

One material correction is required in `book/edition-3/chapter-12.md`, under **TL;DR → Exercise**, in the command beginning `go run ./cmd/grade-edition3 -ch 12 -json`. The reviewed actual-path grading example reverses the workflow and streaming modules (`-workflow "$work/ch07"` and `-stream-exercise "$work/ch06"`):

```sh
  -workflow "$work/ch06" \
  -stream-exercise "$work/ch07" \
```

These are the required replacements, not the current text. Independently checked accepted source: `solutions/edition-3/ch12/ch06/main.go:196–211,274–278` implements the author/editor/reviewer manuscript workflow; `solutions/edition-3/ch12/ch07/main.go:38–42,81–89` implements the streaming exercise and its `stream_stats`, delta and first-delta observations. The current command sends each grader flag the other's artifact. `docs/edition-3-student.md` supplies generic module paths, not evidence for the reversed concrete mapping. No guessed `commands.log` is evidence for this finding.

**Minimal acceptance condition:** swap only those two module-path suffixes. The root should verify the same author's one response against these source locations. No second formal review, source/grader change, rerun of accepted grades, or additional teaching requirement is requested.

## Review scope and findings

Read the current text-reviewer skill and its full governing prerequisites: `AGENTS.md`, the author skill, `docs/edition-3-execution.md` (including the current lean procedure), workflow, carryover, `book/voice.md`, and `book/chapter-writing-procedure.md`. Read the complete current chapter and its actual diff from `/tmp/ch12-author-before.82Fja7/chapter-12.md`; consulted `ch12-restart-author.md`, the final student/code-review qualifications, the public student brief, and accepted `solutions/edition-3/ch12` source. This is the single formal post-student text review, not a reopening of code or grader acceptance.

No other required text revision emerged:

- **Public startup/API and recipes:** the exported aliases, `ConfigFromEnv`, `SkillSpec.Directory`, `StartupSkills`, and `App.NewAgent` match `ensemble.go`, `config.go`, and the common configuration/skill types. The GUI imports the public core. The core `./cmd/ch01` and separate GUI `./cmd` build paths exist; `--gui-debug`, `--skills-dir`, `--save`, and flags-before-`chat` agree with the actual client parsers. The positional grading target, separate artifact/model flags, same-built-CLI Chapter 1 follow-up, and genuine delayed waits are correctly described apart from the two paths above. These are read-only source checks, not claims of newly successful builds or commands.
- **Cold metadata and connected discovery (Skill-based MCP lifecycle, “The bundled browser attaches after Agent construction.”):** the narrow explanation agrees with `internal/cli/skills.go`, `gui/cmd/skills.go`, `internal/mcp/{browser,manager}.go`, and `internal/skills/skills.go`: known disconnected browser metadata permits admission, connected discovery supplies the actual catalog, overrides survive configuration, attachment validates schema/automatic identity, and explicit skill names control grants. It does not license invented tool capabilities.
- **Unload/disclosure (Skill-based MCP lifecycle, “Withdrawn external declarations must disappear even from Anthropic's deferred catalog”):** the distinction is correct. `internal/engine/engine.go` excludes disabled external tools from provider projection; ordinary dormant compiled tools retain their disclosure behavior. Skill instructions can remain pending until compaction while external grants disappear and an ungranted peer closes. The revised unload response in `internal/tool/skills.go` supports this wording.
- **R3 observation (Taking it for a spin, “The accepted Edition 3 run separately changed font size from 16 to 20 using selectors observed in snapshots.”):** the final student R3 record and code-review acceptance support correlated font input/readback, unload, withdrawn catalog declarations in the following request, final response and EOF. That live scenario did not attempt a subsequent refused tool call; refusal is supported separately by grade/regression evidence, not this R3 live observation. The paragraph explicitly avoids a fresh native walkthrough, exact speech-hearing claim, or pixel certification. It does not turn earlier native evidence into R3 evidence.
- **Prose and preservation:** the additions clarify demonstrated interfaces and failures rather than inventing features. Bill's opening, existing SYSTEM/on-demand alignment, historical debugging stories, qualified earlier evidence and historical counts remain intact. No whole-chapter style campaign or speculative cuts are warranted.

## Existing lint and limits

Read the author's saved lint output at `/Users/bill/projects/cadence-memories/children/e3-restart-ch12-author/io/568`: exit 0, all hard ceilings met. The person-gap warnings are soft: baseline 3408 after line 165, final 3480 after line 204. They do not require invented first-person material or unrelated restructuring. No lint was rerun.

Only this review note was written. No manuscript, source, guide, grader, publication or dedication work was changed; no children, builds, tests, grades, browser/network/key/native operations, commits, push or cleanup were performed.

## Note-only factual erratum

After submission, corrected chapter locations to section titles and exact current-text anchors, and removed the overstatement that the R3 live scenario demonstrated a later tool refusal. The formal verdict and sole module-path correction are unchanged. This is a correction to this note, not a second text review or verification of the author response; no new checks were run.
