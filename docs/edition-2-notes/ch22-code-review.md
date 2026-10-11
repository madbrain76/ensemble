# Original 22 / Edition 2 ch20: independent review

Reviewer `reviewer_astra` (GPT-6 Astra/high), October 11, 2026.
**Preparation only; implementation review pending. Zero correction rounds used.**
Read the governing documents and mandatory
`.agents/skills/ensemble-coding/SKILL.md`. No first-edition answer code was
inspected for this assignment; no grader or provider call was made.

## Placement and source coverage

Independently agree with the coordinator's mapping, authorized by Bill and
committed in `6ec35e2`: original21+22 becomes Edition2ch20, original23 becomes
ch21, and original24+25 becomes ch22. Status/accounting completes introspection
before the sandbox; public frontend reuse belongs beside application/fleet
composition. Review original24 before implementing25, retaining every capability
and original-numbered result. Fleet GUI selection remains deferred to26; inspect
its completion at the25 review boundary. This is placement, not authorship or
acceptance of any new implementation.

Read every line, including examples, of these current chapters:

| Original | Lines | File revision | SHA-256 |
| --- | ---: | --- | --- |
| 22 | 586 | `7855ab179562b94c353efda07b53302528ca7c92` | `2589a196ad87504ab7193475a3ffba0e5961d57dd5a35ed4cc4e9a2888208118` |
| 23 | 808 | `e509c2ddc237f799844a6efe6ebdb744289ce295` | `3b19960c74cf4d01e225b9c60bc175b199cc95ec01cad754b97d9a3dcf4bd061` |
| 24 | 565 | `2ca63636c8ba592ea0ed190fb0a3c4731e2b0848` | `c6ac58046ed8e5496c44931e20d126d36081325fc06bfc27f4276c73424976b1` |
| 25 | 582 | `009d74c70b108681b2107b4c3a08685508567cee` | `eb09ec25fe8a681479310ffe22f4cf19d0b4eba34bf507dbbfd8181f7302616f` |

Read the existing24/25 notes and subsequently created student-safe22 notes.
The25 notes still correctly identify unresolved `tool_limits` wording and the
timer/message-wake contradiction. Those are later contract questions, not22
requirements. Do not recreate flat roots, callback bags or coupled GUI modules
to imitate the historical repair narratives.

## Original22 grader preflight

Source inspection confirms the narrow launch/report erratum recorded in
[the student-safe notes](ch22.md). The behavioral harness
`internal/grade/ch22_harness.go` is at
`6f6279367b4e75f47e8eb24568c928fa80dcdc5f`, SHA-256
`df7884a797d8056098473b2de137d7ec295dc58fefd6a5d945834269a3f68dd0`.
It launches `--port PORT` in a fresh cwd with the five documented LLM/skill
environment fields; it sends JSONL `kind:prompt`, awaits nonempty JSONL
`assistant`, and switches models through `/ws` subscribe/update_settings.
No GUI asset override is supplied. Its prompt scanner has no local deadline:
the coordinator must bound the unchanged original run externally. The known
human-terminal envelope mismatch therefore remains relevant, without any
automatic waiver of22 behavior or prediction of a score.

The shared builder `internal/grade/harness.go`, revision
`c24b8e84da112ef66f729c0cdc60413ee6fb547d`, SHA-256
`a54ae4976b4a58dd7a2cd0c4d2fe86f8abda9ee990874669759192bd5d144eff`,
prefers `./cmd/` when that directory exists, otherwise the module root; it has
no optional-GUI entry override. Source checks still require a source tree.

`internal/grade/ch22_status.go`, at the same revision as the22 harness, SHA-256
`c3b4ebc4a78ad58cfe75ef350b0586ff7f95f34cc292ad2e3f2bfcc25c4f12af`,
requires the text labels/order published in the erratum, including lowercase
`model:` and `session` for discovery. It observes the actual returned tool text
in a subsequent provider request. Expected cost uses the models observed on
the wire, avoiding an assumed timing boundary for asynchronous model changes.

The run invokes original21 parity before the new status session, so inherited
failures/timeouts must remain distinguishable from22 findings. Valid status,
per-model accounting, ownership, composition and observer-lifetime obligations
remain binding. Formal disposition awaits frozen code, the student review and
actual original/live/mutation evidence; no new exception is accepted here.
