# Original25 / Edition2ch22 completion: reviewer preparation

Reviewer `reviewer_astra`, GPT-6 Astra/high, October 11, 2026.
**Preparation only; implementation and formal review pending; zero rounds.**
Mandatory `.agents/skills/ensemble-coding/SKILL.md` applies. Current full-source
reading is recorded in [ch22 review](ch22-code-review.md) and [ch25 notes](ch25.md).
No original answer, provider or grader execution was used for this preflight.

## Missing public-input facts

Source `internal/grade/ch25_harness.go`, revision
`009d74c70b108681b2107b4c3a08685508567cee`, SHA-256
`1896697aef507f31c7955c359594501dd644f70c16fcd388f9f955a6e4bb9bc6`.
The shared builder selects `./cmd/` when present. Sessions launch without
arguments in fresh working directories, with `LLM_BASE_URL` and dummy
`OPENAI_API_KEY`/`ANTHROPIC_API_KEY`; no model or skill override is supplied.
Input is JSONL `{"kind":"prompt","text":"..."}`. A JSON `assistant` string
ends the observed turn, including an empty string. The harness waits90 seconds,
closes stdin and allows20 seconds for shutdown/save before killing the process.
Legacy terminal-envelope differences must remain distinct from fleet behavior.

The seven tools receive these public argument shapes:

| Tool | Arguments |
| --- | --- |
| `spawn_sub_agent` | `agent_id`, `task` strings |
| `join_agents` | `agent_ids` array of strings |
| `agents_status` | empty object |
| `shutdown_agent` | `agent_id` string |
| `send_message_to_parent` | `content` string |
| `respond_to_parent_message` | `agent_id`, `response` strings; `message_id` omitted for the single pending message |
| `submit` | `result` string |

These facts were sent to the coordinator for a narrow student-safe erratum,
not as answer code or permission to expose grader internals to the student.

## The two recorded wording conflicts

Original4 explicitly specifies one-shot next-call `tool_limits`, with no target
argument. Source `book/chapter-04.md`, revision
`0253737c06197d6f9b35e442916b7475c1a59a99`, SHA-256
`44e03445c850c939e331cb25c174e98df692063f824f9cedb748efa3ea1fee84`.
Original25 calls it a per-tool-name setting and says that a missing timer prevents
message wakeups, despite depicting independent timer/message select cases.

The actual25 harness calls `tool_limits` with only `ai_callback_delay:60`, then
`run_command`; no per-name target or persistent configuration is provided. It
delivers another ordinary prompt during that command, before the timer expires,
and expects a provider request with the message while the command remains
unfinished. Later actual tool output must still arrive. This supports the
existing next-call interface and an independent message wake trigger.

Recommendation to the coordinator: preserve accepted next-call semantics and
record a narrow prose correction rather than adding a per-tool policy. Timer
expiry and incoming messages are separate causes of ending a wait; neither
means killing the job. This is sourced contract preparation, not an implemented
result or advance acceptance. Current original22 review takes priority, and
original26 remains conditional at the25 review boundary.

## First-green comparative preflight, October 11

Frozen source `7175d33ebfd14c4d2807580a17320786a131a507`, main tree
`9528cdd0dbfda11a7c3b3b77bf33223a4dade259`. This is pre-submission feedback,
not a formal disposition; zero correction rounds. Read the student's stage
review, current implementation, and matching original fleet, messaging, idle,
spawn, workflow-tool code and tests only after this implementation existed.
The original grader actually scored **0/100 FAIL** in0.610s: all seven checks
stopped before inference with `provider model is not set`. No fleet behavior
ran; latent terminal-envelope differences are not observed failures. Final
exception assessment awaits the complete student review and runtime evidence.

Two independent real-runtime probes found implementation defects and were sent
through the coordinator without old-answer code:

- A child with its parent's accepted ID (`main`) generates the same first
  question ID as the root. Its question overwrites the root's pending human
  question; after answering the child, answering the root fails with `no matching
  pending question`. The same test with distinct child/root names passes.
  [Probe](../../solutions/edition-2/evidence/ch25/reviewer-question-collision_test.go.txt),
  [failure and control](../../solutions/edition-2/evidence/ch25/reviewer-question-collision.txt).
- A worker that exhausts its original request plus three wakes escalates and
  stops, but still appears as `running` in the parent's fleet status. The probe
  waits until the child's escalation also causes the root's bounded four
  requests and human escalation: eight actual HTTP requests, then the stale
  child status. This is current lifecycle truthfulness (§25.2/§25.4), not a
  request for Chapter26 activity states.
  [Probe](../../solutions/edition-2/evidence/ch25/reviewer-idle-status_test.go.txt),
  [failure](../../solutions/edition-2/evidence/ch25/reviewer-idle-status.txt).

Both probes use real Ensemble, Agents, managers, actors and persistence; only
external HTTP inference is faked. The coder owns corrections and affected
regressions. No production file was edited by this reviewer. Final source,
mutation sensitivity, live evidence, full counts and disposition remain pending.
