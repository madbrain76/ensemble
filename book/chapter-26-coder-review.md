# Chapter 26 Coder Review

## Must-fix

None.

## Enrichment

1. **wait_for_agent_change poll loop**: The current implementation uses a
   ticker-based poll with 500ms interval. A channel-based approach that wakes
   on the notify channel directly would reduce latency and CPU. Currently
   functional and correct, but the poll interval means state changes are
   detected up to 500ms late. Low priority because the human latency is much
   larger, but worth revisiting if fleet sizes grow.

2. **Tree snapshot walks children recursively on every activity change**: For
   deeply nested fleets (>5 levels, >100 agents), the full tree rebuild on
   every activity observation could be expensive. An incremental delta approach
   (send only the changed node) would scale better. Current approach is correct
   and simple; flag for when fleet sizes grow.

3. **GUI create_agent takes name only**: The WebSocket `create_agent` message
   currently takes only a name. Model and skill default to empty strings, which
   means the Ensemble creates a named agent with no model specified. A
   production GUI would want model and skill pickers. Adequate for the chapter
   since the create/delete mechanism is the point.

## Do NOT add

1. **Per-agent tabs in the GUI**: The tree shows agents and their activity, but
   clicking an agent does not switch the conversation pane. This is explicitly
   deferred to a later chapter. Do not add tab switching here.

2. **Global agent map for cross-tree operations**: The containment boundary is
   the security property. Adding a global map "for convenience" would destroy
   the ownership enforcement that this chapter's grader verifies.

3. **Activity history or timeline**: The activity tracking emits real-time
   observations. There is no historical activity log (which agent was doing
   what at which time). This is a monitoring feature, not an audit trail.
   Deferred.

## Verified facts

- Activity observation emits 4 phases (thinking, tool_call, waiting, idle):
  **VERIFIED** by test and grader check 1 (100/100).
- Ensemble persists and reloads from ensemble.json: **VERIFIED** by
  integration test and grader check 2.
- Ownership enforcement returns "unknown agent" for cross-tree access:
  **VERIFIED** by test and grader check 3.
- wait_for_agent_change wakes on state change: **VERIFIED** by test and
  grader check 4.
- Composition root builds and runs correctly: **VERIFIED** by grader check 5.
- GUI tree message sent on WebSocket connect: **VERIFIED** by grader check 6.
- Bottom-up submit and idle correction work through Ensemble: **VERIFIED** by
  grader check 7.
- Mutation audit 7/7, all mutants kill exactly their expected checks:
  **VERIFIED**.
- Cross-chapter sweep all chapters 100/100, zero breaches: **VERIFIED**.

## Unverified claims

- The HuggingFace breach narrative in the chapter opener. The design doc
  references "the July 2026 HuggingFace breach involving coordinated AI
  agents." This is a story told to motivate the chapter. The details
  (seven hundred agents, three days, message board coordination) are from
  the design doc, not independently verified against a primary source.
  **Flag for author verification before print.**

## Decisions

- Ensemble type lives in `agent/` package (same package as Agent, not
  `agent/ensemble/`). This avoids exporting internal types and keeps the
  import graph simple.
- Activity observations use the existing Observer pattern (no new channels
  or goroutines). The actor emits thinking/tool_call/idle; workflow tools
  emit waiting.
- The poll-based wait_for_agent_change was chosen over pub/sub for
  simplicity. O(children) per wake is acceptable for reasonable fleet sizes.
- Tree snapshots are full rebuilds, not deltas. Simple and correct for
  current fleet sizes.
- main.go went from 1,126 to 999 lines. The Ensemble absorbed factory
  wiring, user sink, MCP lifecycle hooks, and data directory derivation.
