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
