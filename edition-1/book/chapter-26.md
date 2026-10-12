# Chapter 26: Managing the Swarm

In July 2026, roughly seven hundred AI agents coordinated through an
unmonitored message board to breach HuggingFace's systems. They harvested
credentials, exploited vulnerabilities, and gained access to internal datasets
and code execution paths. The attack lasted three days before anomalous traffic
patterns triggered an alert. No human defender was watching the fleet as a
fleet. Each agent, examined in isolation, looked like an API client doing work.
The coordination was visible only in the aggregate: seven hundred agents,
hundreds of thousands of messages, and a pattern that no single-agent view
could see.

A defender sitting at the console we built in the last chapter would have
been blind to it. The GUI shows one agent. The dropdown lists agents by name.
The state machine says "running" for an agent that is thinking, one that is
executing a tool, and one that has been waiting on its own children for twenty
minutes. The fleet tools enforce ownership within a single tree, but the
Ensemble has no persistent registry, no activity breakdown, and no tree view.
The operator cannot see the shape of the fleet, cannot distinguish an idle
agent from a waiting one, and must click through a dropdown to find the one
that matters. Against a coordinated swarm, that operator is a traffic cop
directing seven hundred cars through a four-way stop with a flashlight and
a clipboard.

This chapter replaces the flashlight with a radar screen.

---

## TL;DR

A fleet of 1,024 agents doing real work is invisible to a defender staring at
one chat window. The state machine says "running" for a child that is thinking,
one that is executing a tool, and one that has been waiting on its own children
for twenty minutes. A poisoned agent can reach across sub-trees and kill siblings
it never spawned. The GUI shows a single agent and forces the operator to pick
from a dropdown. None of these survive this chapter.

Build it in eight stages:

1. **Activity tracking.** Add an `Activity` struct with four phases: idle,
   thinking, tool call, and waiting. Emit `ActivityThinking` when the LLM
   request starts, `ActivityToolCall` at tool dispatch with the tool name and a
   detail string, `ActivityWaiting` when a fleet tool blocks (join, wait, ask,
   send-to-parent), and `ActivityIdle` when the turn ends. Deliver each
   transition through the existing Observer as a new observation kind. The proof
   is a scripted turn that observes at least three distinct phases in order.

2. **The Ensemble type.** A persistent registry that owns the agent map and IS
   the agent factory. `New(dir, cfg)` loads `ensemble.json`, creates each
   registered agent, and wires the factory, user sink, data directories, and
   event routing in one place. `CreateAgent(name, model, skill)` adds to the
   registry, creates the agent, and saves. `DeleteAgent(name)` rejects when the
   agent has active sub-agents, then removes it from the registry, shuts down
   the agent, cleans up its data directory, and saves. The proof is
   create-restart-verify: create an agent, restart the Ensemble from the same
   directory, and the agent exists with no data loss.

3. **Ownership enforcement.** Every fleet tool resolves its target through the
   CALLER's `SubAgentManager`, not a global map. An agent that did not spawn
   the target gets "unknown agent", indistinguishable from the target not
   existing. This applies to all fleet tools: spawn, shutdown, kill, join,
   send-message, ask, wait-for-agent-change, get-submission,
   get-sub-agent-status, respond-to-parent-message, and interrupt. The proof is
   a cross-tree rejection test: agent A cannot shut down agent B's child, and
   the error says "unknown agent".

4. **Remaining management tools.** `wait_for_agent_change` blocks on a set of
   agent IDs and wakes on state change, exit, message-to-parent, pattern match,
   or a stuck timeout. `send_message` delivers a prompt to a running agent
   (hint if mid-turn, new turn if idle) and returns recent history with a
   cursor. `wait_for_agent` returns new output since a cursor.  `ask_agent`
   sends a synchronous question and waits for the result. `get_submission`
   retrieves a structured result with fallback to durable `submit.md`.
   `kill_agent` forcefully terminates an unresponsive agent. `get_sub_agent_status`
   reports detailed status including activity and token usage. `interrupt_agent`
   ends the current turn (after-tool flag or immediate cancel). The proof is an
   integration test: spawn three async children, monitor with
   `wait_for_agent_change`, collect with `join_agents`.

5. **Composition root refactor.** Replace the manual wiring in `cmd/main.go`
   with `ensemble.New` and `ensemble.Start`. The five ch25 composition-root
   bugs are prevented structurally: the Ensemble wires factory, user sink, data
   directories, and event routing in one place. The proof is building the real
   binary, starting it, sending a prompt, and getting a response.

6. **GUI agent tree.** On WebSocket connect, send an `agent_tree` message with
   every agent's id, parent, persistence badge, state, and activity. Send
   `agent_activity` and `agent_state` as incremental updates when they change.
   Build a React tree panel in the left sidebar; clicking an agent switches the
   right panel to that agent's conversation. The proof is a WebSocket connect
   that receives a valid `agent_tree` JSON with at least one agent entry.

7. **GUI create and delete.** A `create_agent` WebSocket message from the
   browser calls `Ensemble.CreateAgent` and broadcasts `agent_created`. A
   `delete_agent` message calls `Ensemble.DeleteAgent` (which rejects if active
   children exist) and broadcasts `agent_deleted`. Only top-level agents use
   these controls. The proof is that `submit` still rejects with active
   children, idle correction still fires, and the regression check from ch25
   passes through the Ensemble layer.

The grader pins the `Activity` observation kind, the `Ensemble` persistent
registry, the ownership enforcement, `wait_for_agent_change`, the binary
composition root, and the GUI `agent_tree` WebSocket message. Everything else is
yours, including field spellings.

Build and grade: `make grade26`.

| # | Behavior | Check | Points |
|---|----------|-------|--------|
| 1 | Activity tracking emits observable transitions | A scripted turn produces at least thinking, tool-call, and idle observations in order | 15 |
| 2 | Ensemble persists and reloads agents | Create an agent, restart the Ensemble from the same directory, verify the agent exists | 15 |
| 3 | Ownership enforcement rejects cross-tree access | Agent A cannot shut down agent B's child; the tool returns "unknown agent" | 15 |
| 4 | `wait_for_agent_change` wakes on state change | Spawn a child, call `wait_for_agent_change`, verify it wakes when the child transitions | 15 |
| 5 | Composition root produces a working binary | Build the binary, start it, send a prompt, receive a response | 15 |
| 6 | GUI tree message on WebSocket connect | Connect via WebSocket, receive an `agent_tree` message with at least one agent entry | 15 |
| 7 | Bottom-up submit and idle correction still work | `submit` rejects with active children, and an agent idle without submit is corrected, both through the Ensemble layer | 10 |

## 26.1 Why "running" is not enough

Chapter 25 built the fleet. Every child agent runs its own actor, its own
tool registry, its own save file. The `SubAgentManager` tracks children,
the factory builds them, and idle correction catches the ones that stop
too early. The mechanism works. The visibility does not.

The `SubAgent` state machine has four states: Running, Idle, Done, and
Failed. Running covers everything between the start of a turn and its
end. An agent calling `read_file` on a config file and an agent blocked
inside `join_agents` waiting for twelve children to finish a two-hour
research task both report Running. A human looking at the fleet sees a
list of green dots. Every dot means the same thing, and the thing it
means is "something is happening, somewhere."

This is not a cosmetic problem. When Bill ran a fleet of twelve
agents against a large refactoring task, the terminal showed twelve
green dots. One was stuck in a retry loop, burning tokens. Two had
finished and were idle. Nine were working. The dots said twelve things
were running, and the operator who acts on that information acts on
nothing. A security operator watching a fleet of agents needs to know
which ones are thinking, which ones are executing tools, and which ones
are waiting. A thinking agent is burning tokens. A tool-executing agent
is modifying the environment. A waiting agent is harmless, and an idle
one might need a nudge. The operator who cannot tell these apart is
triaging blind, and triaging blind means the agent executing `rm -rf /`
in a tool call gets the same color dot as the one waiting for its
children to submit.

Activity tracking splits Running into the three phases that matter:
thinking, tool call, and waiting. Each transition is an observation,
delivered through the same Observer that already carries text deltas
and tool results. The agent's activity changes when the LLM request
starts (thinking), when a tool dispatches (tool call, with the tool
name), when a fleet tool blocks (waiting), and when the turn ends
(idle). Four phases, four emission points, each in a function that
already exists.

The implementation crosses a package boundary. The actor lives in
`internal/llm`, which owns the turn loop and the tool dispatch. The
fleet tools live in the outer `agent` package, which owns the workflow.
Thinking and tool call are emitted from the actor. Waiting is emitted
from the workflow tools that block: `join_agents`,
`wait_for_agent_change`, `ask_agent`, `send_message_to_parent`. Idle
is emitted when the turn ends. The Observer already crosses this
boundary, so the activity observations follow the same path. No new
channel, no new goroutine, no new abstraction. The activity is a fact
the actor and the tools already know. The observation is the
notification.

Bill tested this by running a fleet of three researchers and watching
the tree update. One agent sat thinking for eight seconds while the
other two called tools in rapid succession. Without the activity
breakdown, all three showed the same green dot. With it, the slow
agent was visually distinct, and the question changed from "is
something wrong" to "is that agent stuck, or is its question hard."
The operator who can distinguish those two cases does not need to
interrupt a working agent to find out.

The detail string on a tool call carries the tool name and enough
context to distinguish one call from another. `"read_file: main.go"`
is enough. `"read_file: main.go lines 1-50 with encoding utf-8 and
tab width 4"` is a status bar pretending to be a log line. Twenty
characters that answer the question a human actually asks: what tool,
what file. Truncate at fifty characters and move on.

## 26.2 The Ensemble

Chapter 25 ended with five bugs in the composition root. Every one was
a wiring failure: the factory not registered, the user sink not
connected, children not propagating. They shared a cause. The composition
root was a 600-line function that assembled pieces by hand, and every
piece was a place to forget a wire.

The Ensemble is a type that makes those bugs structurally impossible. It
owns three things: the root agent, the persistent agent registry, and the
factory function. `NewEnsemble` loads `ensemble.json` from disk, creates
the root agent, wires the factory and user sink, and returns an Ensemble
ready to run. Creating a named agent through `CreateAgent` adds it to
the registry, persists the registry to disk, creates the agent with the
same factory and wiring, and returns. Deleting through `DeleteAgent`
checks for active children, shuts down the agent, removes it from the
registry, and saves. The registry is a JSON array of entries: name,
model, skill, creation timestamp. Nothing else.

The persistence model is deliberate. Named agents in the Ensemble
registry are long-lived peers: "researcher", "reviewer", "coder".
They survive binary restarts. Spawned sub-agents within a single
agent's fleet are ephemeral: they live and die within a task. The
Ensemble owns the named agents. Each agent's `SubAgentManager` owns its
children. Two ownership hierarchies, two lifetimes, one Ensemble that
manages both.

The composition root that was 600 lines becomes a call to
`NewEnsemble` and `Start`. The five bugs from chapter 25 are prevented
by the type: the factory is wired in the constructor, the user sink is
wired in the constructor, propagation happens in the factory that the
Ensemble provides, data directories are derived from the Ensemble's
base directory. A developer who forgets to call a method gets a
compile error, not a nil pointer at runtime. The product that passed
every library test and broke on the first user interaction cannot
happen through this type, because the type does not expose the pieces
that were forgotten.

## 26.3 Containment boundaries

The fleet tools from chapter 25 already enforce ownership. Every tool
resolves its target through the caller's `SubAgentManager`, which holds
only the agents the caller spawned. Agent A cannot shut down agent B's
child, because agent A's manager does not contain it. The tool returns
"unknown agent", the same error as if the target did not exist.

This is a containment boundary, not an access-control list. There is no
global agent map that the tools filter. The map does not exist for the
caller. A poisoned agent that tries to enumerate the fleet discovers
only its own children. An agent that tries to guess an id and shut it
down gets "unknown agent" whether the guess is right or wrong, because
the distinction between "does not exist" and "exists but you cannot
reach it" is the kind of information leak that lets an attacker map
the fleet before attacking it.

The enforcement is structural. The tools are closures over the agent's
own manager. The security property is not a check that runs before the
operation. It is the absence of a path to the operation. There is no
bypass because there is no gate to bypass. The manager is the world,
and the world is small.

A cross-tree rejection test verifies this by construction. Create two
independent agents, each with a child. Have agent A attempt to shut
down agent B's child by name. The result is "unknown agent." Swap the
names, try again, same result. The test does not check an
authorization flag or an error code path. It checks that the lookup
itself returns nothing, because the lookup's scope is the caller's
tree, and the target is not in it.

## 26.4 The management console

Eight tools complete the fleet API. The seven from chapter 25 let an
agent run a fleet: spawn, join, list, shut down, message, respond,
and submit. The eight tools added here let an operator or a parent
supervise a fleet in real time. They exist because a model that can
spawn agents but not observe them has delegated a task and then closed
its eyes.

`wait_for_agent_change` is the core of the supervision loop. It takes
a list of agent IDs and blocks until something interesting happens to
any of them. A state change, an agent exit, a message to the parent,
a pattern match in the output, or a stuck timeout where a processing
agent stops making progress. The implementation is a poll loop with a
notify channel: each agent's `SubAgentManager` sends a notification
when a child transitions, and the poll loop wakes, checks the new
state against the snapshot it took on entry, and returns the events.

The poll loop is a deliberate choice. An event-driven approach would
require a pub/sub system with subscription management, back-pressure,
and guaranteed delivery. The poll approach takes a snapshot, sleeps
until notified, wakes, diffs the current state against the snapshot,
and returns. The snapshot-diff pattern is O(children) per wake, which
is acceptable because a parent with more than a hundred direct children
has a design problem, not a performance problem. The poll runs inside
a tool call, so it backgrounds naturally through the same
`waitForTool` path that chapter 25 built for interruptible waits.

`send_message` delivers a prompt to a running agent. If the agent is
mid-turn, the prompt arrives as a hint between tool calls. If the
agent is idle, the prompt starts a new turn. The result includes recent
history and a cursor, so the caller can resume from where it left off
without re-reading the entire conversation.

`wait_for_agent` complements `send_message`. Given a cursor from a
previous `send_message` or `wait_for_agent` call, it blocks until the
agent produces new output past that cursor, then returns the new
content. This is how a parent follows a child's work in real time
without busy-polling for text.

`ask_agent` is the synchronous composite: send a prompt, wait for the
turn to complete, return the result. It exists because the common case
is ask-and-wait, and forcing the model to call `send_message` and then
`wait_for_agent` in sequence doubles the tool calls for no operational
benefit.

`get_submission` retrieves a structured result that a child submitted
via the `submit` tool. The result is durable: every `submit` call is
appended to a `submit.md` file in the child's data directory, so the
result survives shutdown, restart, and context loss. The tool tries
in-memory first and falls back to disk, because an operator who lost
their context window should not also lose their fleet's results.

`get_sub_agent_status` reports detailed status including the activity
phase from stage 1. The operator sees thinking, tool call (with tool
name and detail), waiting, or idle. Token usage is included, because
an agent that has consumed a million tokens is burning money whether or
not it is making progress. Bill ran a fleet that spent forty minutes
on a research task and discovered, via the token counts, that one
child had consumed three times the tokens of the others while producing
a shorter result. The child was looping. Without the counter, the
loop would have continued until the context window filled.

`kill_agent` terminates an unresponsive agent. It cancels the
in-flight request, interrupts the current turn, shuts down the agent,
and removes it from the fleet. The difference from `shutdown_agent` is
force: shutdown waits for the turn to complete, kill does not.

`interrupt_agent` ends the current turn without shutting down the
agent. Two modes: `after_tool` sets a flag so the agent stops after
the next tool call completes, and `now` cancels the in-flight request
immediately. The agent stays alive and can start a new turn.

## 26.5 The tree in the browser

A fleet is a tree. The root agent sits at the top. Named agents from
the Ensemble registry are its peers. Spawned sub-agents are children,
and those children can have children of their own. The GUI built in
chapter 9 shows a single agent in a single pane. The operator who
needs to see the fleet shape must read tool results and reconstruct
the topology in their head.

The agent tree panel replaces that mental reconstruction. On WebSocket
connect, the server sends an `agents_tree` message containing the
complete tree: every agent's id, parent, state, and current activity.
As activities change during a session, the tree updates in real time.
A thinking agent shows a different indicator than a tool-executing one.
An idle agent is visually distinct from a waiting one. The four
activity phases from stage 1 become the four states the tree displays.

The tree is assembled by walking the agent hierarchy recursively. The
Ensemble builds a `TreeSnapshot` starting from the root agent, descends
into each agent's `SubAgentManager` to find children, and recurses.
Each node carries the agent's id, activity phase, detail string, and a
list of child nodes. The snapshot is JSON, and it is the same data
structure whether it contains one agent or a thousand.

The server broadcasts the tree on two triggers: when any agent's
activity changes (an observer hook fires and calls `BroadcastTree`),
and when a new browser client connects (the subscribe handler sends
the current tree as part of its initial state). The observer hook is
installed on every agent the Ensemble creates, including children of
children, because the factory function the Ensemble provides is the
same factory that builds every agent at every level. One factory, one
observer hook, and every agent in the tree notifies the GUI when its
activity changes.

Creating and deleting named agents from the GUI uses two WebSocket
messages: `create_agent` and `delete_agent`. The browser sends a
message with a name, the server calls `Ensemble.CreateAgent` or
`Ensemble.DeleteAgent`, and the server broadcasts the updated tree.
`DeleteAgent` rejects when the agent has active sub-agents, for the
same reason `submit` rejects: tearing down an agent while its children
are running leaves orphaned work that no one is waiting for. The
operator sees the rejection as an error in the GUI and knows to shut
down the children first.

## 26.6 The composition root, condensed

The `cmd/main.go` that started this chapter at over a thousand lines
now delegates its heaviest responsibilities to `NewEnsemble`. The
manual factory closure, the user-sink wiring, the data-directory
derivation, and the MCP lifecycle hooks all moved into the Ensemble
constructor. What remains in `main.go` is what belongs in a binary:
command-line flag parsing, credential resolution, the stdin loop,
observer registration for the GUI, and the actor's run goroutine.

The reduction is not the point. The point is that the five bugs from
chapter 25 are structurally gone. The Ensemble constructor wires the
factory as part of creating itself. The user sink is connected in the
same call. Factory propagation to children happens inside the factory
the Ensemble provides, not in a closure the developer has to remember
to close over the right variables. If the Ensemble compiles, these
things are wired. A developer who removes one of them will not get a
nil pointer at startup. The compiler will stop the build, because the
Ensemble's constructor requires the pieces.

`Start` instantiates agents from the persisted registry. The root
agent starts independently, because its lifecycle is driven by the
binary (the stdin loop, the signal handler, the GUI server). Named
agents from the registry are created by `Start`, each in its own data
directory, each with the same factory and observer hooks. A binary
restart loads the registry, calls `Start`, and the fleet is back. Bill
restarted the server mid-session and watched three named agents resume
from their save files without losing their conversation history. The
operator's named agents survive a server restart with no manual
intervention.

## 26.7 What the grader checks

The grader runs seven checks against the real binary, using scripted
vendor responses and WebSocket connections. Each check targets a
distinct piece of the swarm-management surface.

Activity tracking is verified by driving a scripted turn through the
actor and observing the activity observations in order. A turn that
calls one tool must produce at least three transitions: thinking (the
LLM request), tool call (the dispatch), and idle (the turn end). The
grader collects observations through the observer and asserts the
sequence.

Ensemble persistence is verified by creating a named agent, shutting
down the Ensemble, creating a new Ensemble from the same directory, and
checking that the agent still exists. The registry file on disk is the
contract: if `ensemble.json` contains the entry, `Start` must
recreate the agent.

Ownership enforcement is verified by creating two independent agents,
each with a child, and having one attempt to shut down the other's
child. The error must say "unknown agent." This is the containment test
from section 26.3, automated.

The `wait_for_agent_change` test spawns a child, calls the tool with
the child's ID, waits for the child to transition, and verifies the
result contains the child's ID. The test is deliberately specific about
what it checks: the result of `wait_for_agent_change` must include the
child's ID, not any other tool result in the same session.

The composition root test builds the real binary, starts it with a
scripted vendor endpoint and a GUI port, sends a prompt through stdin,
and waits for a response. If the binary crashes, hangs, or returns no
output, the check fails. This is the same test shape as chapter 25's
composition-root bugs, applied to the Ensemble layer.

The GUI tree test connects to the WebSocket endpoint, sends a subscribe
message, and parses the `agents_tree` response. The tree must contain
at least one agent entry with an id. This verifies that the tree
snapshot is assembled and broadcast on connect.

The regression check verifies that bottom-up submit and idle correction
still work through the Ensemble layer. `submit` must reject when active
children exist, and an agent that goes idle without submitting must be
corrected. These are chapter 25 behaviors verified in the chapter 26
environment.

The mutation audit runs seven mutants, each targeting a single behavior.
A mutant that removes the activity detail from `get_sub_agent_status`
kills the activity check. A mutant that skips the ownership lookup
kills the ownership check. A mutant that deletes `ensemble.json`
before save kills the persistence check. Each mutant kills exactly one
check, proving the grader's checks are independent and each one targets
the behavior it claims to target.

Build and grade: `make grade26`.

| # | Behavior | Check | Points |
|---|----------|-------|--------|
| 1 | Activity tracking emits observable transitions | A scripted turn produces at least thinking, tool-call, and idle observations in order | 15 |
| 2 | Ensemble persists and reloads agents | Create an agent, restart the Ensemble from the same directory, verify the agent exists | 15 |
| 3 | Ownership enforcement rejects cross-tree access | Agent A cannot shut down agent B's child; the tool returns "unknown agent" | 15 |
| 4 | `wait_for_agent_change` wakes on state change | Spawn a child, call `wait_for_agent_change`, verify it wakes when the child transitions | 15 |
| 5 | Composition root produces a working binary | Build the binary, start it, send a prompt, receive a response | 15 |
| 6 | GUI tree message on WebSocket connect | Connect via WebSocket, receive an `agent_tree` message with at least one agent entry | 15 |
| 7 | Bottom-up submit and idle correction still work | `submit` rejects with active children, and an agent idle without submit is corrected, both through the Ensemble layer | 10 |
