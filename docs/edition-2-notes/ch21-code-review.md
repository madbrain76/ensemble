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

## First frozen stage: independent probes, before final submission

Inspected source `5a00b98f489bc06b533ac612145e1102e9199a60`, main tree
`02e9bb14883f4dddc0e76c844656369bd91f6973`. No formal disposition yet; final
student review, live results and targeted mutation work remain pending.

The focused real-runtime HTTP tests and shipped-skill check pass independently
under race detection,1.143s, together with a small external-peer primary probe.
That probe confirms a departure from Rule5: constructing a primary with an HTTP
MCP declaration performs initialize/initialized/discovery and exposes its selected
tool. Existing primary loading already starts every newly enabled skill after
the root attaches MCP. This coherent inherited behavior may warrant an explicit
exception to the chapter's historical silent-no-connect trap rather than a
deliberate regression. The coordinator asked the student to explain the departure
in its review; no waiver or implementation change is imposed before that response.

The independent **no-receive-broadcast mutant survives** both new HTTP tests,
exit0 in0.025s. Removing only `t.ready.Broadcast()` preserves their apparent
unload success because `client.close` closes its Done channel before transport
shutdown. That channel is not evidence that the already-parked HTTP `Recv` exits.
The production wakeup is present; this is a sensitivity gap, not an observed
production failure. Return the finding for a focused receive-lifetime check and
restored-defect proof before final submission. No additional runtime capability
or paid run is requested. This is attributed pre-submission review assistance,
zero formal correction rounds. See [probe receipt](../../solutions/edition-2/evidence/ch21/reviewer-initial-probes.json)
and [primary probe](../../solutions/edition-2/evidence/ch21/reviewer-primary-probe.go.txt).
Both isolated scratch checks were removed afterward.

The original grader's completion wait also disagrees with the existing GUI
terminal: the harness requires JSON `assistant`/`error` on stdout, while this
command prints raw deltas/newlines and errors on stderr. Visible remote calls
can complete while the harness waits90seconds. The eventual score and exact
failure remain for the final evidence review; this source observation does not
turn visible tool output into proof of every masked rubric assertion.

Initial [comparative counts](../../solutions/edition-2/evidence/ch21/reviewer-counts.json)
are tied to this frozen stage. Original21 has no original20 snapshot predecessor,
so original19→21 includes the crossover's unrelated model, cache, actor and GUI
repairs. Do not attribute all that growth to the HTTP feature. Relevant original
HTTP transport/tests and integration were read separately. Final counts and
comparative judgment will follow the submitted source, not the moving worktree.
Mandatory handoff: `.agents/skills/ensemble-coding/SKILL.md`.
