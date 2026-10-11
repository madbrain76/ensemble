# Original21 / Edition2 Chapter20 — independent review

## October 11 preparation and original-grader preflight

Reviewer `reviewer_astra` read all692 lines of original21, the student-safe
preparation, and reread the complete `.agents/skills/ensemble-coding/SKILL.md`.
The governing carryover/workflow were read in full for the preceding review.
Original20 is accepted; the current coder extends its direct Edition2 `ch19`
predecessor. This is **preparation only**, zero formal correction rounds. No
unfinished source, original answer, credential profile or provider was inspected;
no grader or service request was run.

The original exercise adds HTTP MCP to the existing transport and loadable-skill
seams. Its required behavior includes JSON/SSE POST replies,202 notifications,
session headers, environment authentication with a missing-key diagnostic,
meaningful transport/tool errors, and skill-authorized dispatch/unload. Preserve
the real prior runtime and cancellation/lifetime behavior. The deployed primary
must allow loading the supplied web-search skill; a fixture's primary proves
nothing about that shipped configuration. Fetched content remains tool-result
evidence, not a guarantee of model resistance to prompt injection. No Chapter22
status/accounting features or unrelated HTTP/MCP framework are authorized.

Read all four original21 grader files. The runner builds the supplied directory,
launches `--port`/`--gui-dir`, supplies skill/provider environment, and drives
stdin prompt envelopes while awaiting assistant/error completion. Thus the
optional GUI executable is its intended entry, even though the actual live
exercise may use the CLI. Unlike prior fixtures, its base explicitly declares
`type: primary`; do not reuse the missing-type exception automatically. The
base also names `think`, but these checks neither require that declaration nor
call it, so its name alone is not an observed fixture failure or waiver.

Coverage limits, to assess after the frozen implementation and actual result:

- Before/after declaration checks and real URL calls establish skill loading.
  They do not test guessed unauthorized dispatch, unload, primary no-connect,
  or the deployed primary. The fake advertises exactly the two allowed tools;
  it never presents an extra tool to filter.
- The fake returns single-frame SSE and202 notifications. It does not issue
  session IDs, demand authentication, vary JSON response framing, or exercise
  concurrent close/blocked receive/cancellation. These public obligations need
  their corresponding implementation and focused supplemental evidence.
- Search and nonce checks inspect real calls and actual tool-result bytes.
  The nonce protects against a canned fetch. They do not establish save/replay
  portability or unrestricted live-provider reliability.
- Injected text must appear only inside a tool result, which is a useful
  structural check. Its scripted model cannot demonstrate prompt-injection
  resistance or broader authorization enforcement.
- Both MCP `isError` and HTTP500 diagnostics must be actual error results.
  Matching text alone does not earn those checks. The fake's404 intentionally
  uses `isError:true`, whereas the historical real service could return a404
  page as ordinary content. Do not invent a metadata reinterpretation feature
  solely to hide that distinction.

These are test boundaries, not preemptive findings or a request to expand scope.
No score or reviewer exception is claimed. The source's assertion that skill
gating cannot be mutation-tested does not waive the mandatory targeted-defect
procedure. Meaningful existing inherited evidence can be reused; no whole-book
audit or unnecessary paid repeat follows from these observations.

Source chapter revision is `7159f0ba8708c7582f4c7d64cf09a204d644dd70`, hash
`e1d44feb80a51676e628ed4e5682484325099120d1f5a29aa5c8be060cc359c1`.
The current grader source provenance is recorded below before any student run.

| File | Last file revision | SHA-256 |
| --- | --- | --- |
| `ch21_checks.go` | `a87e6d406f6acd6e24a0885ae48e22a8ed6e51d0` | `e58283417a8e141f5409c29b6917530e574f2dfd86d990ed7bb96206d4c4a171` |
| `ch21_run.go` | `ad3a6dc7ed20d96a785e86d1880c9caac36f7162` | `63979a1d0cea261a95d72c089f814249b63d9fdbd7950f76904e9fb239706f61` |
| `ch21_harness.go` | `a87e6d406f6acd6e24a0885ae48e22a8ed6e51d0` | `f948d70153334739356cbb74ea97514406472d2f5581bca6385e1b88413eb76a` |
| `ch21_fakemcp.go` | `107c59d45ee2977d82910f1c1e734eabb73d538a` | `1fcd849b77794f85a80d23ede777571cd716a8d0bf31c6c3c0ef54b5b9c8293e` |

Wait for frozen source and the complete student review before comparative
assessment. Original22 is the next binding stop after accepted21; original24
and stopped authorship remain. Preserve CodeRhapsody's concurrent work.
Mandatory handoff: `.agents/skills/ensemble-coding/SKILL.md`.
