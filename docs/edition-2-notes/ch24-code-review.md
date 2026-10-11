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


## First-green comparative preflight

Inspected frozen student `35df549210130f824fdc2a4eb632ea878b2c529b`, main tree
`5b36f28cbc8bac320febe1a9afa810efd27452e3`, after reading the student's stage
review. Only then read matching first-edition implementation, own-predecessor
changes and relevant tests. Final student review/live/mutation assessment remains
pending; no acceptance or formal correction round is claimed.

Actual original grade is **0/100 FAIL**,20.89s, no external timeout. Four checks
stop while assembling fixed root/same-module GUI probes; two launch the headless
CLI with GUI flags and never get a GUI port; the root scan's literal free-function
NewAgent positive control misses the actual public composition-root method.
These match the earlier sourced preflight. Valid obligations remain required;
no exception is awarded merely from these fixture diagnoses.

Independent exact-archive `TestExternalReuse` passes under race detection,
including the separately compiled outside-module tests: public event construction,
replay and prompt through real Agent/Engine with an external HTTP fake;
embedded/disk asset selection; actual Chrome component rendering with only two
modules served and no shell; same-ID concurrent relay clients with reverse-order
replies alongside the core handshake; shutdown; actual headless/GUI executable
transitive dependency negative and positive controls. See
[output](../../solutions/edition-2/evidence/ch24/reviewer-first-green-external.txt).
The archive only links the already installed node_modules as test dependencies;
its implementation bytes come from the frozen commit, not the moving main tree.

An additional small real-network probe passes under `-race`: after a relay
request reaches the browser, TCP EOF removes the parked request; its late browser
reply does not block the next request; browser EOF returns an error preserving
the second caller's original ID. It uses actual Ensemble/Agent/GUI and external
socket peers, with no model call or credential access. Retained
[probe](../../solutions/edition-2/evidence/ch24/reviewer-relay-disconnects-test.go.txt)
and [output](../../solutions/edition-2/evidence/ch24/reviewer-relay-disconnects.txt)
check the new advertised lifetimes without editing student source.

Provisional own-predecessor counts are in
[receipt](../../solutions/edition-2/evidence/ch24/reviewer-counts.json).
Student net production adds216 physical Go lines (187 active,20 comment-only,9
blank) and8 JS (4 active,3 comments,1 blank); tests add351 Go and38 MJS.
Original net production adds24 Go (13 active,11 comments), with no JS growth;
tests add144 Go. Complete before/after totals and lexical method are recorded.
Moves are counted once: the original moved its existing GUI/WebSocket/assets;
the student extracted existing renderers/component from app.js. The original
already had its TCP relay and standalone component, explaining much of the
student's additional Go growth. No core production change or fleet capability
appears in the student's delta. Final KISS/comments/tests disposition waits for
completed evidence; no blocking source finding has emerged at this checkpoint.

Display transcripts normalize trailing whitespace only.
