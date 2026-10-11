# Original24 / Edition2ch22 opening stage: reviewer preparation

Reviewer `reviewer_astra`, GPT-6 Astra/high, October 11, 2026.
**Preparation only; implementation and formal review pending; zero rounds.**
Mandatory `.agents/skills/ensemble-coding/SKILL.md` applies. Full chapter/read
coverage and independent merge agreement are recorded in [ch22 review](ch22-code-review.md).
No original answer, provider or grader execution was used for this preflight.

Inspected launch/public-consumer contracts in `internal/grade/ch24_harness.go`,
revision `6896df7faca8b049f80293c19a944436e1656007`, SHA-256
`5631a794d043a79ac43b22178666fcd5d9e6b797ddb195afe4741c3adee10de5`.
Wire compatibility reuses original22's GUI/stdin/model-switch harness. Embedded
deployment launches `--port PORT` from a fresh unrelated cwd, with the same
LLM/skill environment inputs as22 and no disk asset override. The core-only
entry and legacy terminal-envelope limitations therefore remain relevant.

External probes assume one module with `<base>/gui`, a `gui.AgentHooks`
callback bundle and `gui.New`; the temporary module supplies only one root-module
replace. The probes spell `Assets.ReadFile("web/artifact-scroll.js")` and
`web/renderers.js`, `ServeWS`, `ServeMCP(net.Listener)`, and a specific root
event/log vocabulary. The supplied outside consumer mounts only the component.
Its WebSocket sends `type:subscribe`, then `type:prompt` with text, and observes
`type:message` with text. Its bare TCP MCP relay sends a JSONL `tools/list`
request and expects a JSON-RPC result/error with the same id before any browser
connects. These are observations of the harness, not required implementation
names or authority to introduce callback bags.

Both external programs build together: an incompatible GUI probe can prevent
the independent headless dependency checks from running. The checks also assume
`./cmd/virtual-user`, `<base>/mcpws`, and literal root `func NewAgent(` as a
positive control. Fixed package/module/constructor assumptions cannot override
Bill's separate optional GUI module and actual-parent architecture. Review any
actual failed result independently rather than pre-awarding exceptions here.

Preserve all valid obligations: outside consumers construct public inputs and
exercise replay/prompt/relay paths; components serve without the application
shell; embedded deployment and explicit disk override work; core/MCP transitive
dependencies exclude GUI/WebSocket implementation with positive controls; the
actual model switch survives unchanged. No additional feature is authorized by
the grader's reference-shaped fixture. Original24 remains separately reviewed
before25, with only the accepted combined source exported as Edition2ch22.
