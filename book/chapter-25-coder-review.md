# Chapter 25 coder review: "The Ensemble"

Written by the coder for the author. The code is finished, graded and
frozen. This document says what exists, what it cost, and which parts
are worth narrating.

`make grade25` scores 100/100. `make grade25-audit` kills 7 of 7.
The full cross-chapter sweep is clean: every chapter 100/100, chapter 5
at 120/120, `breached=0`.

---

## 1. What shipped, in commit order

Sixteen commits. The four stages the contract names, then the bugs that
building them exposed.

| Commit | What |
|---|---|
| `1dcabdd` | Stage 1: `SubAgentManager`, `AgentFactory` |
| `79ef5c3` | Addressing: every agent has an id and a parent id, and the user is addressable |
| `29fef70` | Messaging: agent to parent, agent to user, and the reply path |
| `b03d501` | Stage 2: the seven workflow tools, async turns, the join fan-in |
| `a678ea3` | Stage 3: a message ends a tool wait, and the tool still delivers |
| `9f8c48c` | Stage 4: wake an agent that stopped without submitting, then escalate |
| `c78d217` | Bug: no agent factory in the composition root |
| `dd49868` | Bug: no user sink, factory never propagated, `message_id` required |
| `4c8c42f` | Grader, five of seven checks |
| `130fe98` | Grader, all seven |
| `6516ed4` | Bug: input arriving during a turn was not read until the turn ended |
| `889e5e1` | Mutation audit, `make grade25`, `make grade25-audit` |
| `fa10657` | Frozen `solutions/ch25` |
| `f8dcc83` | Regression: a package counter cost chapter 5 ten points |
| `c19f82d` | Resync the snapshot after that fix |
| `70e9e3f` | Chapter 19 grader: the event log path was relative |

---

## 2. The chapter's best material: five bugs with one shape

This is the spine I would build the chapter on, because the five
defects found while writing it are not five unrelated mistakes. They are
one mistake wearing five hats.

**Every one of them was invisible to a passing unit test, because each
test wired by hand the exact thing the product forgot to wire.**

1. **No agent factory in the composition root.** Spawning returned
   `no agent factory configured` in the shipped binary. Every fleet test
   passed, because every fleet test sets its own factory.
2. **No user sink in the composition root.** The human was addressable in
   principle and unreachable in practice. Escalation from a top-level
   agent had nowhere to arrive, which is precisely where a human most
   needs to hear.
3. **The factory was never propagated to children.** A child's manager
   held nil, so **the fleet was permanently one level deep**. Found only
   because the messaging probe is three levels deep. A two-level test
   would have passed forever.
4. **`respond_to_parent_message` demanded a `message_id`** copied out of
   prose the model had to read. Now optional: an omitted id answers the
   only pending message.
5. **The command loop read stdin on the goroutine that waited for the
   turn.** Stage 3 built a wake mechanism the terminal could not reach.

Number 5 has the best line in it. The existing comment claimed slash
commands could steer "a turn that is already running". They could not.
The structure forbade what the comment promised, and the comment had
been sitting there being wrong in a file everyone had read.

The honest summary: **a library test proves the library. Only the
binary proves the product.** The grader drives the real binary over its
real stdin, which is why it found all five.

---

## 3. The design decisions worth narrating

**The user gets an address, not an entry in the agent map.** Addressing
is uniform (a top-level agent's parent simply *is* the user, so the root
needs no special case) while lifecycle is not (nothing spawns or shuts
down a human). Putting the user in the children map would force an
"unless it is the user" guard into every loop that walks the fleet.

**Messaging is the general verb; submit is the narrow one.** The first
draft scoped the parent link to submit time, which quietly made *asking
a question* require claiming to be done. Bill caught it.

**A reply is a message, not a string.** An answer with no sender attached
cannot be followed up. The human's address is configurable, so it can be
a name rather than the role.

**Interrupting is not killing.** A message ends the *wait*, not the work.
The tool keeps running and still delivers its result later. `Interrupt`
refuses when nothing is waiting, which keeps the flag off backgrounded
jobs where a later waiter would read it as a message that never came.

**Idle correction is reactive, with no timer and no polling.** Going idle
is already an observed event. A *failed* turn is deliberately not
corrected: that is a real error, not a model lapse, and waking it just
repeats the error.

---

## 4. Two failures worth a paragraph each

**A package counter cost chapter 5 ten points.** Stage 3 added
`var messageSeq atomic.Uint64`. The chapter 5 grader caught it, 120 down
to 110, and **no chapter 25 test could have**: every one of them runs a
single fleet, and a package counter is only wrong when two fleets share
a process. The cross-chapter sweep is the only thing in the build that
can see that class of defect. This is the second time this book has paid
for this exact regression.

**A bisect indicted innocent code.** Chapter 19 scored 85 against the
live tree and 100 against four consecutive worktrees, which framed a
one-line change as the culprit. The diff could not possibly touch event
logs, and that contradiction was the clue. **Every worktree path is
absolute**; the sweep calls graders with relative paths; chapter 19's
harness set `cmd.Dir` to the student directory and then passed a
relative log path. Same tree, absolute 100, relative 85. How a test is
invoked is part of the test.

---

## 5. Must-fix for the author

Nothing. The code is green and frozen. These are accuracy notes:

- The chapter should not claim the sandbox, skills or MCP chapters
  needed changes. They did not.
- If the chapter says the fleet is "arbitrarily deep", that is now true
  but was not before `dd49868`. The depth limit was a bug, not a design.
- The TL;DR is frozen at `cr/docs/ensemble-ch25-tldr.md`: seven checks,
  100 points, verified against the real linter. Per Bill's ruling the
  grader engineer owns it; it should be dropped into the chapter
  unchanged.

## 6. Enrichment, if there is room

- The `ch14-work-*` processes the sweep warns about are leaks from
  chapter 14, not chapter 25 children. Worth a footnote only.
- The mutation audit has an honest oddity: mutant 2 kills six checks, not
  one. The breadth **is** the finding, because a child handed its
  parent's data directory takes zero turns. The one scenario that spawns
  no child survives, which is how we know the grader did not simply fall
  over. A sharper mutant was tried and killed the same six.

## 7. Do NOT add

- **Do not describe a general message router.** There are exactly two
  addresses, parent and user. Routing between arbitrary peers belongs to
  the chapter that gives the Ensemble the agent map.
- **Do not claim the framework knows whether a sub-agent succeeded.**
  It observes `submitted`, which is a claim to check, not a fact.
- **Do not present the five bugs as sloppiness.** They are the natural
  result of testing a library instead of a product, which is the point.
- **Do not say the escalation threshold is tunable.** It is pinned at
  three and tested by its own assertion, because a test written in terms
  of the constant moves with it and cannot detect a change to it.

## 8. Facts verified, and facts not

**Verified by running it:**
- `make grade25` 100/100; `make grade25-audit` 7/7, 0 problems.
- Full sweep clean, `breached=0`, chapter 5 at 120/120.
- `solutions/ch25` byte-identical to `agent/`, graded independently.
- The whole `agent` module passes under `-race`, 13 packages.
- Chapter 19: same tree, absolute 100, relative 85, before the fix.

**Not verified:**
- Behaviour on Linux or Windows. Everything here was measured on
  macOS arm64.
- Any claim about how many agents the design scales to. Nothing above
  four concurrent children was ever run.

## 9. Status

Code complete, graded, frozen, swept. Ready for prose.
