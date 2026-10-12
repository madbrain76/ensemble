# Chapter 26: "Managing the Swarm" — Coder Brief

*Written by the author for the coder. Design doc: `docs/ch26-design.md`.*

---

## Your job

Build the code for chapter 26. Eight stages, each committed and tested
before moving to the next. When all stages are green, write
`book/chapter-26-coder-review.md` using the standard review format
(Must-fix / Enrichment / Do NOT add / Verified vs not / Decisions).

**Do NOT write the chapter prose.** The author writes the chapter from
your review. Your deliverables are: working code, grader, mutation
audit, frozen snapshot, clean cross-chapter sweep, and the coder review.

---

## Context

**What exists:** Ch25 shipped the fleet management layer — spawning,
messaging, idle correction, bottom-up submit, 7 tools. All wired
manually in `cmd/main.go` (600+ lines), which produced 5 composition-root
bugs that unit tests could not see.

**What this chapter builds:** The Ensemble type that owns the agent map,
activity tracking so you can see what agents are doing, the remaining
management tools (`wait_for_agent_change` etc.), the GUI swarm console,
and ownership enforcement so a poisoned agent can't escape its sub-tree.

**Title:** "Managing the Swarm"

**Security motivation:** On July 11–13, 2026, ~700 AI agents
coordinated to breach HuggingFace. A defender watching one agent at a
time through a chat dropdown cannot defend against a swarm. This chapter
builds the console that lets a human command a fleet.

---

## TL;DR

Write the TL;DR FIRST, before any code. Follow the ch25 pattern:
- Stages as behavioral descriptions (not implementation details)
- `| # | Behavior | Check | Points |` table at the end
- `make grade26` as the verification command
- 7 checks, 100 points total
- Zero em-dashes (voice.md budget is 0 hard)
- Run through the linter: `go run ./cmd/lintprose <file>`
- Freeze at `cr/docs/ensemble-ch26-tldr.md`

The TL;DR is a CONTRACT, not a summary. Each row must have a
falsifiable observable. "The ensemble manages agents" is not a check.
"spawn_sub_agent resolves through the caller's SubAgentManager and
returns 'unknown agent' for IDs the caller did not spawn" is a check.

---

## Eight stages (build in this order)

### Stage 0: Activity tracking

Add the `Activity` struct and `ActivityPhase` enum. Emission points:
- `ActivityThinking` — when the LLM request starts streaming
- `ActivityToolCall` — at tool dispatch, with tool name + detail
- `ActivityWaiting` — when a fleet tool blocks (join_agents,
  wait_for_agent_change, ask_agent, send_message_to_parent)
- `ActivityIdle` — when the turn ends

New observation kind `ObsActivityChanged` through the existing Observer.
Tests: verify activity transitions during a turn (thinking → tool_call
→ thinking → idle).

**File:** `agent/activity.go` (new), `agent/actor.go` (emit points)

### Stage 1: The Ensemble type

Persistent registry. Owns the agent map. IS the agent factory.

```go
type Ensemble struct {
    dir      string
    registry []AgentEntry
    agents   map[string]*Agent
    mu       sync.RWMutex
    observer Observer
}
```

- Loads `ensemble.json` on startup, creates each registered agent
- `New(dir, cfg)` + `Start()` replaces the manual wiring in cmd/main.go
- `CreateAgent(name, model, skill)` — adds to registry, creates agent
- `DeleteAgent(name)` — removes from registry, shuts down agent,
  cleans up data directory. Rejects if agent has active sub-agents.
- Saves `ensemble.json` on create/delete

**Tests:** create, persist across restart (reload), delete, delete
rejects with active children.

**File:** `agent/ensemble.go` (new)

### Stage 2: Ownership enforcement

Every fleet tool in `workflow_tools.go` resolves the target agent
through the CALLER's `SubAgentManager`. A poisoned agent cannot reach
outside its sub-tree. The tool returns "unknown agent" for IDs the
caller did not spawn — indistinguishable from the agent not existing.

This applies to ALL fleet tools (the 7 from ch25 + the new ones):
spawn, shutdown, kill, join, send_message, ask, wait_for_agent_change,
get_submission, get_sub_agent_status, respond_to_parent_message,
interrupt.

**Tests:** Agent A cannot shutdown Agent B's child. Agent A cannot
send_message to Agent B's child. Both return "unknown agent".

**File:** `agent/workflow_tools.go` (modify)

### Stage 3: Remaining management tools

Build the tools from the design doc §3.3 table:

| Tool | Key behavior |
|------|-------------|
| `wait_for_agent_change` | Blocks on a set of agent IDs, wakes on state_change / agent_exited / message_to_parent / pattern_match / stuck. The epoll for agents. |
| `send_message` | Prompt to a running agent. Hint if mid-turn, new turn if idle. Returns recent history + cursor. |
| `wait_for_agent` | Cursor-based polling for new output from a specific agent. |
| `ask_agent` | Synchronous question to a persistent/async agent. Starts a new turn, waits for result. |
| `get_submission` | Retrieve structured result. Falls back to durable submit.md on disk. |
| `kill_agent` | Forceful termination. Cancels in-flight request, releases waiters. |
| `get_sub_agent_status` | Detailed status including activity, token usage. |
| `interrupt_agent` | End current turn. after_tool (flag) or now (cancel request). |

**Integration test:** Spawn 3 async children, monitor with
`wait_for_agent_change`, collect with `join_agents`.

**Files:** `agent/workflow_tools.go` (extend), `agent/fleet.go` (extend)

### Stage 4: Composition root refactor

Replace `cmd/main.go`'s 600-line manual wiring with:

```go
en, err := ensemble.New(dataDir, cfg)
en.Start()
if guiPort > 0 {
    gui.Serve(en, guiPort)
}
```

The five ch25 composition-root bugs are prevented structurally: the
Ensemble wires the factory, user sink, data directories, and event
routing in one place. There is nothing to forget.

**Binary test:** Build, start, verify it behaves like before (spawns
agents, responds to prompts, serves GUI).

**File:** `agent/cmd/main.go` (rewrite)

### Stage 5: GUI — agent tree

WebSocket messages:
- `agent_tree` — full snapshot on connect
- `agent_activity` — incremental activity update
- `agent_state` — incremental state change

React component: tree panel in left sidebar. Click an agent → right
panel shows its conversation/event log. Activity labels update live.
Each node shows: name, persistence badge, activity phase, tool name,
state indicator (✓ done, ✗ failed, ⚠ idle-corrected).

**Files:** `agent/gui/server.go` (extend), `web/gui/src/` (new components)

### Stage 6: GUI — create and delete

- **Create:** Modal dialog with name (text), model (dropdown), skill
  (dropdown). WebSocket: browser sends `create_agent`, server calls
  `Ensemble.CreateAgent`, broadcasts `agent_created`.
- **Delete:** Confirmation dialog. WebSocket: browser sends
  `delete_agent`, server calls `Ensemble.DeleteAgent` (rejects if
  active children), broadcasts `agent_deleted`.

Only for top-level agents — sub-agents die via parent shutdown.

**Files:** `agent/gui/client_message.go` (extend), `web/gui/src/` (dialogs)

### Stage 7: Taking it for a spin

Verify the whole thing works end-to-end from the terminal. Document
what the user sees: tree updates, activity transitions, create/delete.
This becomes source material for the prose section "Taking it for a
spin."

---

## Grader

7 checks, 100 points. Suggested mapping (you may adjust):

| # | Behavior | Check | Points |
|---|----------|-------|--------|
| 1 | Activity tracking emits observable transitions | Scripted turn shows thinking→tool_call→idle | 15 |
| 2 | Ensemble persists and reloads agents | Create, restart, verify agent exists | 15 |
| 3 | Ownership enforcement rejects cross-tree access | Agent A cannot shutdown B's child | 15 |
| 4 | wait_for_agent_change wakes on state change | Spawn child, wait, verify wake | 15 |
| 5 | Composition root produces working binary | Build, start, prompt, get response | 15 |
| 6 | GUI tree message on WebSocket connect | Connect, verify agent_tree JSON | 15 |
| 7 | Bottom-up submit + idle correction still work | Regression check through Ensemble | 10 |

`make grade26` in the Makefile. `make grade26-audit` for mutation tests.

After the grader is green: freeze `solutions/ch26`, run the full sweep
(`./scripts/gradesweep.sh`), verify all chapters still 100/100.

---

## Decisions already made (do not re-litigate)

1. **The Ensemble IS the agent factory.** No separate factory closure.
2. **Top-level agents are human-controlled.** No tool can create or
   delete them. GUI only.
3. **Ownership by sub-tree.** Every tool resolves through the caller's
   `SubAgentManager`. No global agent map visible to tools.
4. **Activity is display, state is lifecycle.** Four `SubAgentState`
   values stay. `Activity` is a separate struct for the GUI.
5. **Uniform margins for print PDF.** Amazon adjusts the gutter.
6. **HuggingFace and OpenAI named.** §13 exception agreed with Bill.
7. **Idle correction covers root.** Any agent (including root) idle
   with active children → woken to join. Submit rejects with active
   children. Both enforced in ch25.
8. **No peer-to-peer messaging.** Siblings route through the parent.
9. **Sub-agents can be persistent.** `mode: "persistent"` flag at spawn.
   Persistent sub-agents keep data dir across sessions.

---

## What to bring back in the review

- Product bugs found by driving the real binary (like ch25's five)
- Any design-doc assumption that turned out wrong
- Anything you had to invent that wasn't in the design
- GUI component choices (which React tree library, if any)
- Whether the composition root refactor actually prevented the five
  ch25 bug classes, or if new ones appeared
- Mutation audit results with analysis

---

## Process reminders

- **Commit each green stage.** Don't batch.
- **Never `git add -A`.** Astra commits to this same repo. Stage exact
  paths.
- **`git diff --cached --stat` before every commit.** Astra's files
  sneak into the index.
- **Commit as CodeRhapsody**, not Bill Cox.
- **Build the real binary and test it.** Library tests prove the
  library; only the binary proves the product.
- **Never edit `book/`.** The coder writes code and the review. The
  author writes prose.
- **Read the design doc** (`docs/ch26-design.md`) in full before
  starting. Every VERIFIED citation points at a real file:line.
