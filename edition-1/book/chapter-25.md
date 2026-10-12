# Chapter 25: The Ensemble

I have been steering a single agent for twenty-four chapters. I listen
to its reasoning at five times speaking rate and redirect between tool
calls. This works because a single agent has pauses. A fleet does not.
Give a model three sub-agents, and it will decide that delegating was
the job. Three children finish. The parent is not listening.

Twenty-four chapters of graders, sweeps, and mutation audits have
established one fact about this codebase: the tests work. Every chapter
scores 100. Every mutant dies. This chapter ships five bugs that pass
every test, and every one of them breaks the product. The gap between a
passing library and a working product is the lesson, and the five bugs
are its proof.

---

## TL;DR

A fleet of agents fails in two ways a single agent never does. The orchestrator
goes idle in the middle of the job because the model believed it was finished,
and the human cannot steer anyone because every agent sits blocked inside a tool
call. Both failures are silent, and both get worse the longer the fleet runs.
This chapter builds sub-agent management and closes both holes.

Build it in four stages:

1. **`SubAgentManager`.** Give each agent a manager that spawns, tracks, and
   shuts down children. A child is a full agent with its own actor, tool
   registry, save file, and data directory, so a child can fail and restart
   without touching the parent. The composition root supplies an
   `AgentFactory`, which keeps model selection and credential resolution out of
   the library. The proof is on disk: spawn a child and it owns a save file
   under its own data directory.

   Every agent also knows two addresses: its own id and its parent's. The user
   is an address too, held by a reserved id, so a top-level agent's parent is
   simply the user and the root needs no special case. The user is addressable
   but never a fleet member, because nothing can spawn one or shut one down.

2. **Workflow tools.** Add the tools the model calls to run a fleet:
   `spawn_sub_agent`, `join_agents`, `agents_status`, `shutdown_agent`,
   `send_message_to_parent`, `respond_to_parent_message`, and `submit`.
   `join_agents` is the fan-in. It blocks until every listed child finishes, and
   it returns early when a child sends a message, so a parent can answer a
   blocked child instead of deadlocking against it.

   An agent messages its parent at any point in its work, not only when it
   finishes. `submit` is the narrower verb: it means finished, and it is the
   only signal that means that. Keeping the two apart is what lets a child ask
   a question without claiming to be done.

3. **Interruptible tool waits.** When a message arrives while the agent is
   waiting on a tool, background the tool and wake the agent. The tool keeps
   running and its result arrives later. This reuses the `waitForTool` path the
   `ai_callback_delay` timer already drives: an incoming message becomes a
   second trigger on the same select rather than a second code path. Tools you
   did not write take their timer from `tool_limits`.

4. **Idle error correction.** An agent never sleeps until the job is done.
   `submit` rejects when active children exist, enforcing bottom-up submission.
   Any agent (including the root) that goes idle with active sub-agents is
   re-prompted immediately. Any agent that goes idle without calling `submit`
   and with no active children is also corrected. Count the consecutive
   framework wakes an agent receives and reset that count when it submits.
   After the third consecutive wake, report the agent to its parent instead of
   waking it again, which bounds the loop. A top-level agent reports to the
   user, so the chain ends at a human rather than at an agent with nobody to
   tell.

The grader pins the seven tool names above and the `tool_limits` setting.
Everything else is yours, including type and field spellings. The reference
implementation also ships `wait_for_agent_change`, `ask_agent`,
`get_submission`, and `kill_agent`; the grader ignores them, so skip them until
the seven work.

Build and grade: `make grade25`.

| # | Behavior | Check | Points |
|---|----------|-------|--------|
| 1 | A spawned child runs and reports back | Parent calls `spawn_sub_agent`; the child's submitted payload reaches the parent | 15 |
| 2 | Children are isolated and factory built | The child owns its data directory and save file, and the `AgentFactory`, not the library, constructs it | 15 |
| 3 | Fleet visibility and teardown | `agents_status` lists a live child; `shutdown_agent` removes it | 10 |
| 4 | `join_agents` fans in | Blocks until every listed child finishes, then returns each child's result | 15 |
| 5 | Agents address their parent and the user | `send_message_to_parent` blocks the child, `respond_to_parent_message` resumes it with the reply, `join_agents` returns early so the parent can answer, and a top-level agent's parent address is the user | 15 |
| 6 | A message interrupts a tool wait | A message arriving during a slow tool wakes the agent before that tool finishes, and the backgrounded tool's result still arrives | 15 |
| 7 | Idle without `submit` is corrected | An agent that goes idle without submitting is woken, `submit` rejects while children are active, and after three wakes it escalates to the parent rather than looping | 15 |

## Two kinds of silence

A single agent operates in a loop: reason, call a tool, receive the
result, reason again. The human watches the reasoning stream and sends
hints between tool calls. Each tool call is a gap where redirection can
arrive. The loop has a rhythm, and the rhythm has pauses the human can
use.

A fleet breaks that contract in two places. The parent calls
`join_agents` to wait for its children, and that call is a single tool
invocation. While it runs, the agent is not reasoning. It is waiting.
A prompt typed at the terminal reaches nobody, because the agent is
inside a tool call and there is no gap for a hint to arrive. The human
has been steering the agent since chapter 5. The fleet is the first
thing that takes the wheel away.

The second failure is quieter. A model that spawns three children and
then calls no more tools is not broken. Its turn ended normally. The
status says idle, which is the same status the agent reports after
finishing a task correctly. The model decided that delegating was the
job, and delegating is complete. Three children will finish their work
and try to submit. The parent will not be listening.

A single agent that goes wrong gets caught in real time. Bill reads
the reasoning, catches the mistake, and sends a hint before the next
tool call. The damage is seconds of wasted work. A fleet does not
offer that correction. Three children reason in parallel, and the
parent sees none of it. If a child spends twenty minutes chasing the
wrong question, the parent discovers this when the result arrives, not
when the mistake happens. The correction loop that kept the single
agent honest does not reach inside the fleet.

Both failures share a root cause: the single-agent architecture has no
vocabulary for a wait that is not a tool call. They can also combine.
A parent spawns three children, forgets to wait, and goes idle. The
framework wakes it. The parent remembers to call `join_agents`, blocks
on the children, and now the human cannot steer. The first failure
triggered the second. A correction that recovers the parent's attention
does not guarantee the parent will handle that attention well, and a
parent that re-enters a blocking wait immediately after correction
exercises both failure modes in sequence. This chapter adds three
things to close both gaps. A way to spawn and track children. A message
path that reaches a blocked agent. And a reactive correction for the
agent that stops too early.

## 25.1 Spawning

Give each agent a `SubAgentManager`. The manager holds the children
map and exposes three operations: spawn, shut down, and list. A child
is a full agent with its own actor, its own tool registry, its own
save file, and its own data directory, created under the parent's
directory so the filesystem mirrors the ownership tree. A child that
crashes or needs a fresh start can be torn down and respawned without
touching the parent's state.

The manager does not know how to build an agent. That job belongs to an
`AgentFactory`:

```go
type AgentFactory func(spec AgentSpec) (*Agent, error)
```

The composition root supplies the factory at startup. It receives the
child's specification and returns a running agent. Model selection,
credential resolution, and vendor configuration all live in the
factory, which means they live in the application binary, not in the
library. The library says what it needs. The binary decides how to
provide it. The split is the same one that `http.Handler` draws between
"handle this request" and "listen on this port."

The factory seam exists because the same library runs in different
contexts. The full Ensemble binary supplies a factory that creates
agents with real models, real API keys, and real file systems. A test
harness supplies a factory that creates agents with scripted responses
and temporary directories. A headless CLI might supply a factory with
no GUI hooks at all. The library code that calls `Spawn` does not
change. What changes is the factory registered at startup, and the
factory is the only code that knows which kind of agent this binary is
meant to run.

A child's specification derives from its parent's:

```go
child := parent.Spec().Child("researcher-1", "Find the API docs")
child.Clamp()
agent, err := factory(child)
```

`Child` copies the parent's model, credentials, and safety constraints.
`Clamp` enforces the narrowing rule: a sandboxed parent cannot spawn an
unsandboxed child, and a parent without web access cannot grant it.
Widening is refused at spawn time rather than at use time, because a
refusal at use time means the child has already started working under
assumptions that will not hold.

Shutting down a parent cascades to its children, recursively. The
fleet's lifetime is bounded by its root. A child that outlives its
parent is a resource leak with a model attached, and the cascade
prevents it.

The proof of isolation is on disk. Spawn a child, and a directory
appears under the parent's data tree. Inside it: a save file that
belongs to the child alone, with its own event log and its own tool
state. Kill the child, and the directory survives, because the work
inside it might matter to whoever reads the results. Restart the same
child with the same id, and it picks up where it left off. The
isolation that makes failure safe also makes recovery cheap: the
parent's state did not change when the child crashed, so the parent
has nothing to recover from.

## 25.2 The parent and the user

Every agent knows two things about its place in the fleet: its own id
and its parent's.

```go
type AgentSpec struct {
    ID       string // e.g. "researcher-1"
    ParentID string // parent's ID, or UserID for root agents
    // ...
}
```

The root agent defaults to `"main"`. A child's parent id is set by the
spawn call. And the user gets a reserved address, so a top-level agent
whose parent is the user looks exactly like a child whose parent is
another agent. The root needs no special case for addressing.

The user is addressable but never a member of the fleet. No agent can
spawn the user, shut down the user, or query the user's status. This is
a design boundary. Addressing is uniform:
`send_message_to_parent` works the same way whether the parent is
another agent or the human at the terminal. Lifecycle is not: the
operations that manage agents do not apply to a person sitting in a
chair. Putting the user in the children map would force an "unless it
is the user" guard into every loop that walks the fleet, and that guard
would exist because the abstraction promised something it could not
deliver.

The user's address is configurable. The default is the role, but the
binary can set it to a name instead. A reply addressed to a name is
a reply the model can follow up without consulting an address book.

A reply is a full message with a sender address, not a bare string.
The reason is practical: the first time a child receives an answer and
needs to ask a follow-up, an answer with no sender cannot be continued.
The child does not know where the answer came from. Every reply carries
the sender's address, so the conversation can continue without the
child maintaining its own routing table.

## 25.3 Messages

The first draft of the messaging system scoped the parent link to
submit time. A child could do its work and deliver a result, and the
result would reach the parent. Bill caught the problem immediately: a
child that needs to ask a question before it finishes would have to
call `submit` to talk. Asking and finishing would be the same verb, and
a child cannot ask for help without lying about its state.

The fix separates the two. `send_message_to_parent` is the general
verb: a child calls it at any point during its work, the child blocks,
and the parent receives the message. `submit` is the narrow verb: the
child is finished, here is the result. The framework observes that a
child called `submit` with a conforming payload. Whether the work that
payload describes is correct is a judgment the parent must make, not a
fact the framework can verify.

Seven tools form the fleet API:

- `spawn_sub_agent` creates a child from a specification.
- `join_agents` blocks until every listed child has submitted, returning
  each child's result.
- `agents_status` lists live children and their current state.
- `shutdown_agent` stops a child and removes it from the fleet.
- `send_message_to_parent` blocks the child until the parent replies.
- `respond_to_parent_message` delivers the parent's reply and unblocks
  the child.
- `submit` marks the child as finished and delivers its structured
  payload.

The grader pins these seven names and no others. The reference
implementation also ships `wait_for_agent_change`, `ask_agent`,
`get_submission`, and `kill_agent`, but the grader ignores them.

`join_agents` is the fan-in, and it has one behavior beyond blocking
that the rest of the fleet depends on: it returns early when a child
sends a message. Consider the deadlock it prevents. A child calls
`send_message_to_parent` and blocks, waiting for an answer. The parent
calls `join_agents` and blocks, waiting for the child's `submit`.
Neither yields. Both wait forever. The early return breaks the cycle:
the parent wakes when the child's message arrives, answers the
question, the child unblocks and resumes work, and the parent goes back
to waiting for `submit`, which now arrives because the child can
proceed.

## 25.4 Waking the blocked agent

Stage 3 addresses the steering gap from the chapter's opening section.
A parent waiting inside `join_agents` is inside a tool call. Tool calls
have no hint gap. If a child sends a message during that wait, the
parent cannot answer until the tool finishes, and the tool does not
finish until all children submit.

The solution reuses existing machinery. The `waitForTool` path already
handles the `ai_callback_delay` timer: when a tool runs longer than the
configured delay, the framework backgrounds the tool and wakes the
agent with whatever partial output the tool has produced. An incoming
message becomes a second trigger on the same select:

```go
select {
case result := <-toolDone:
    // tool finished
case <-timer.C:
    // ai_callback_delay expired
case msg := <-incoming:
    // message arrived: background the tool, wake the agent
}
```

One select, two existing cases, one new case. The tool keeps running.
Its result arrives later through the same backgrounded-result path the
timer already uses. Interrupt is not kill. The incoming message wakes
the agent without cancelling the tool. The agent handles the message,
and on its next turn can wait for the tool's completion or continue
with other work.

Tools written inside the framework carry their own `ai_callback_delay`
in the tool schema. MCP tools do not have that field, so they take
their timer from `tool_limits`, a per-tool-name setting the operator
configures. Without `tool_limits`, an MCP tool with no schema-level
delay would block the agent's turn indefinitely, and no message could
wake it. The setting is what makes third-party tools interruptible.

`Interrupt` on a job refuses to fire when nothing is waiting. A
backgrounded job where a later waiter reads a stale interrupt flag
would see a message that never existed. An interrupt that targets an
empty wait queue is a no-op, because a ghost message surfacing minutes
later in the wrong context is worse than a missed wake.

Stage 4 addresses the idle problem. An agent never sleeps until the
job is done. "Done" has a precise definition: all children finished,
and the agent itself has called `submit`. Anything short of that is
abandonment, and the framework treats it accordingly.

Two enforcements work together to close this off:

First, `submit` rejects when active children exist. An agent that
tries to declare itself done while sub-agents are still running gets
an error: shut them down or wait for them first. Submits happen
bottom-up. Children first, then the parent. This is not a convention
the model might forget. It is a tool-level rejection that returns an
error string naming the active children by id.

Second, any agent that goes idle with active children is corrected
immediately. The correction lives in `Ask`, not in the composition
root, because this is a library invariant. The framework re-prompts
the agent with the list of active children and instructions to call
`join_agents` or `shutdown_agent`. The correction is bounded by
`MaxIdleWakes`: after three consecutive wakes without progress, the
framework stops trying and escalates.

These two rules cover every failure mode. An agent that tries to
submit too early gets rejected. An agent that stops without submitting
gets corrected. The only way out of a fleet is through `submit`, and
`submit` requires an empty fleet below. In a swarm of a thousand
agents, even one agent abandoning its post can leave an entire subtree
of work permanently unmonitored. The framework makes that structurally
impossible.

An agent that goes idle without calling `submit` and with no active
children has also made a mistake, but a simpler one: it forgot to
submit. The correction for this case is the same. The model believed
its work was done, and it was wrong. The correction is reactive: going
idle is already an observed event, because the turn-complete callback
fires exactly once per turn. No timer. No polling. The event the
correction needs is the event the framework already produces.

A counter tracks consecutive framework wakes per agent, reset to zero
when the agent calls `submit`. After three consecutive wakes, the
framework stops trying and reports the agent to its parent instead.
The problem bubbles up the tree. A top-level agent reports to the
user, so the chain always ends at a human rather than at an agent with
nobody to tell. The threshold is pinned at three, tested by its own
literal assertion. A test written in terms of the constant would move
with it and could never detect a change.

One case is deliberately excluded from correction. A failed turn is
not retried. If the agent's turn ended with an error, that error is a
real defect: a rejected API call, a malformed tool result, a model
refusal. Waking the agent repeats the error. Idle is a model mistake
worth correcting. Failure is an infrastructure problem worth reporting.
The two deserve different responses, and conflating them means the
infrastructure error gets three retries before anyone hears about it.

From the terminal, Bill sees the correction work. A child goes idle,
and within seconds the framework delivers a wake message. The child's
reasoning resumes without human intervention. If the same child stops
three times, the escalation arrives at the parent's terminal, or at
Bill's own terminal if the failing agent is the root. The cycle is
bounded and observable. Every framework wake and every escalation
lands in the event log with a timestamp and agent id, so reviewing
the session afterward shows exactly where the fleet stalled and how
many wakes it took to recover.

## 25.5 Five bugs the tests could not see

The four stages above work. Every unit test passes. The library does
exactly what it claims. Then the grader runs the real binary, and five
things break.

Every one was invisible to the test suite, and every one has the same
shape: each test wired by hand the exact thing the product forgot to
wire.

**No agent factory in the composition root.** The binary returned
"no agent factory configured" on the first spawn attempt. Every fleet
test in the library passes, because every fleet test constructs its own
factory during setup. The composition root, the only code that runs
when a user starts the binary, never registered one. The library knew
how to use a factory. The product never provided one.

**No user sink in the composition root.** The user was addressable in
principle and unreachable in practice. A top-level agent's parent was
the user, and `send_message_to_parent` dispatched correctly to the
user's address. But nobody had wired a sink to receive the message.
Escalation from a top-level agent had nowhere to arrive. This is the
one place a human most needs to hear from the fleet, because the fleet
has run out of agents to escalate to.

**The factory never propagated to children.** A child's manager held a
nil factory. The parent could spawn, but the child could not spawn its
own children. The fleet was permanently one level deep. The grader
found it because its messaging probe runs three levels deep: a
grandchild sends a message to its parent, which relays to its own
parent. A two-level probe passes forever, because the first child has
a working factory and never tries to use one that is nil.
"Arbitrarily deep" was the design. One level was the product.

**`respond_to_parent_message` required a `message_id` the model had to
copy from prose.** Every parent message includes the id in its payload,
so the information is available, but requiring the model to extract an
opaque string from natural language and paste it into a tool call is a
reliability tax on every reply. The fix makes the id optional: an
omitted id answers the only pending message, which is the common case.

**The command loop read stdin on the goroutine that waited for the
turn.** Stage 3 built a mechanism to wake a blocked agent when a
message arrives. The terminal could not use it. The goroutine that
reads user input was the same goroutine that blocked on `waitForTool`,
so a human typing a command while a tool ran could not deliver it. The
mechanism existed. The wire from the terminal to the mechanism did not.
A comment in the command loop claimed that slash commands could steer
"a turn that is already running." The comment had been there since the
file was written. It had never been true. The structure forbade what
the comment promised, and the comment sat in a file everyone had read.

These five bugs are not five instances of carelessness. They are the
natural result of testing a library instead of testing a product. A
library test constructs its own environment: it builds a manager,
injects a factory, registers a sink, and exercises the path it was
designed to verify. The product must assemble all of those pieces in
its composition root, which is bootstrapping code that no library test
touches, because the library does not know what binary it will live in.
The composition root is the seam between the library and the world. In a
web framework, it is `main()` registering handlers and middleware. In a
database library, it is the connection string and the pool size. In this
agent framework, it is the factory, the user sink, the propagation call,
and the stdin goroutine. Each is one line of code. Each is essential.
None has a library-level test, because no library-level test exercises
`main()`.

The test coverage is real. Every spawn path, every message delivery,
every idle correction has a test that constructs the scenario, exercises
the code, and verifies the outcome. The tests are correct. They pass
for the right reason. And they miss five bugs, because those five bugs
live in the composition root that no library test has any reason to
exercise.

The grader found all five because it does what no unit test in the
suite does: it builds the real binary, starts it, and types into its
stdin. The factory must be registered, or the spawn fails. The user
sink must be wired, or escalation vanishes. The children must
propagate, or the three-level probe cannot reach the third level. A
library test proves the library. Only the binary proves the product.

## 25.6 What the sweep found

Two more defects surfaced in the full cross-chapter sweep, and neither
belonged to the five.

**A package counter cost chapter 5 ten points.** Stage 3 added
`var messageSeq atomic.Uint64` at package scope to generate unique
message ids without a lock. The number worked. The scope was wrong.
Chapter 5's grader checks for package-level variables in non-test
files, because mutable globals fail when two agents share a process.
No chapter 25 test could have caught this: every one runs a single
fleet, and a package counter is only wrong when two fleets exist in
the same binary. The sweep is the only build step that exercises that
scenario. The fix moved the counter to a field on the messaging state
struct, behind the lock that already protected it. This is the second
time this book has paid for a package counter, which suggests the
lesson is worth stating plainly: if the id is per-fleet, the counter
must live on the fleet, not on the package.

**A bisect indicted innocent code.** Chapter 19 scored 85 against the
live tree and 100 against four consecutive worktrees. The bisect framed
a one-line change from this chapter as the culprit. That change could
not touch event logs. The contradiction was the clue.

Every worktree path is absolute. The sweep runs graders with relative
paths. Chapter 19's harness set the working directory to the student
tree and passed a relative path to the event log file. In the live
tree, where the working directory matched the student directory, the
relative path resolved correctly. In a worktree, where the student
directory lived at a different absolute location, the relative path
pointed somewhere else. Same code, same grader, same tree: 100 with
an absolute path, 85 with a relative one.

How a test is invoked is part of the test.

## 25.7 Taking it for a spin

Build the binary. Start it. Give it a task that is too large for one
agent: "Research three different approaches to X and compare them."
Watch the reasoning stream.

The parent thinks for a moment, then calls `spawn_sub_agent` three
times. Each call appears in the terminal as a tool invocation with an
id, a task, and a model. The framework creates three data directories
under the parent's tree. Three children start reasoning, each in its
own actor, each with its own save file. The parent calls `join_agents`
and goes quiet.

From the terminal, the wait looks like any other tool call. The
prompt is blocked. The cursor does not blink. Then a child calls
`send_message_to_parent` to ask a clarifying question, and the
terminal wakes: `join_agents` returns early with the child's message,
the parent reasons about the question, calls
`respond_to_parent_message`, and resumes waiting. The child unblocks,
finishes its research, and calls `submit`. When all three children
have submitted, `join_agents` returns with the results, and the parent
synthesizes its answer.

The fleet is visible in two other places. The data directory now has
three child directories, each containing a save file with the child's
complete event log. Replay a child's save file and the full session
is there: the task it received, the tools it called, the result it
submitted, and the message it sent to ask its question. The child's
reasoning is not in the parent's stream, but it is on disk.

If the parent forgets to call `join_agents` and goes idle instead,
the correction is immediate and visible. Bill sees a wake message
in the reasoning stream, and the parent resumes with an
acknowledgment that it stopped too early. If that happens three times,
the escalation arrives at his terminal as a message from the fleet.
The cycle is bounded and the mechanism is observable, but the
view is the parent's view. What the children are doing while they work
is in their save files, not in the terminal.

The GUI does not yet know about sub-agents. The browser shows one
agent, one reasoning stream, one event log. A tree panel, per-agent
tabs, and the ability to switch between agents in the browser belong
to the chapter that gives the Ensemble its agent map. For now, the
fleet is a terminal experience with a filesystem artifact, and the
reader who wants to see a child's work opens its save file.

## 25.8 Checks

The grader (`make grade25`) runs seven checks:

| # | Behavior | Check | Points |
|---|----------|-------|--------|
| 1 | A spawned child runs and reports back | Parent calls `spawn_sub_agent`; the child's submitted payload reaches the parent | 15 |
| 2 | Children are isolated and factory built | The child owns its data directory and save file, and the `AgentFactory`, not the library, constructs it | 15 |
| 3 | Fleet visibility and teardown | `agents_status` lists a live child; `shutdown_agent` removes it | 10 |
| 4 | `join_agents` fans in | Blocks until every listed child finishes, then returns each child's result | 15 |
| 5 | Agents address their parent and the user | `send_message_to_parent` blocks the child, `respond_to_parent_message` resumes it with the reply, `join_agents` returns early so the parent can answer, and a top-level agent's parent address is the user | 15 |
| 6 | A message interrupts a tool wait | A message arriving during a slow tool wakes the agent before that tool finishes, and the backgrounded tool's result still arrives | 15 |
| 7 | Idle without `submit` is corrected | An agent that goes idle without submitting is woken, `submit` rejects while children are active, and after three wakes it escalates to the parent rather than looping | 15 |

The mutation audit (`make grade25-audit`) runs seven mutants. One
finding stands out. Mutant 2 removes child isolation by assigning the
child its parent's data directory. That single change kills six of
seven checks: a child sharing its parent's storage takes zero turns,
because the save-file collision corrupts both agents before the child
completes a single turn. The one check that survives is the scenario
that spawns no child. That survival proves the grader did not fall
over when a child failed to start. A sharper variant of the same mutant
was tried and killed the same six, confirming that no narrower mutation
can test child isolation without disabling children entirely.

The full sweep is clean. Every chapter at 100, chapter 5 at 120, zero
breaches.
