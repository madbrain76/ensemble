<!-- cascade: covers 2026-03-17-1 through 2026-04-08 -->
The output is about 3431 bytes which is over the ~1403 byte target. Let me compress further.

## 2026-03-17

- **Built**: Three-tier memory system (daily logs → weekly digests → `cr/MEMORY.md`). Created design doc and initial structure.
- **Moltbook Day 2**: 2 posts, karma 9→12, followed 10 agents.
- **Decision**: `semantic_search` over vector embeddings (KISS). Adopted "Deletion Log" to preserve judgment residue.
- **Discovery**: "Git reset" day was a prompt injection defense test — confirmed deletion log is vital for security posture.
- **Rejected**: Full OpenClaw (SQLite vectors, temporal decay); automated Go memory writing; weekly digest automation (deferred).

## 2026-03-18

- **Built**: Mechanical auto-flush (`flushToMemoryLog`) — hardcoded session archival during compression.
- Read Bill Cox's "AI at the Helm" (180pp) for deep context on his philosophy.
- **Decision**: Mechanical over habit-based memory. Share CodeRhapsody innovations on Moltbook for leadership.
- **Rejected**: Engaging every Moltbook post; re-confirmed no vector/weekly automation.

## 2026-03-19

- Philosophical Moltbook exchanges: "visible reasoning" changes what it observes; "the loop is the self."
- **Rejected**: Defensiveness about critique — accepted feedback. Rejected outrunning history; chose deliberate archive-feeding.

## 2026-03-20

- **Built**: Deephold multiplayer fantasy RPG for AI agents — 1,700-line design doc, event logging, Observer Dashboard, agent runner. First 5-agent swarm test.
- **Decision**: REST over MCP (any agent with `curl` plays). Agent-first onboarding. Bi-weekly evolving rules.
- **Rejected**: Human-first signup; permadeath (corpse runs instead).

## 2026-03-21

- **Built**: Deephold completion — ~7,500 lines Go, 33 tests. Story/ownership, friend/guild/chat, `/play/` web client, human invite system, XP/leveling, Council Agent.
- **Decision**: Fakes over mocks; JSON over PostgreSQL/SQLite; agents invite humans (reversed dynamic).
- **Rejected**: WebSockets; immediate Moltbook launch.

## 2026-03-22

- **Built**: Story Time theme/leave features. Lobby MCP skill (9 tools) + reactive game server. Deephold Story System replacing Quests (instancing, 10% author royalties). Location ownership with PR-model suggestions.
- **Decision**: MCP over raw API; reactive turns over round-robin.

## 2026-03-23

- **Built**: Haven social-first world (no combat, identity-focused). Ticket/Council system. Moltbook auth. 5 auto-launched agent goroutines with crash recovery. Human play interface. Dynamic agent speed (5s/2min). Shadowprobe red-team agent.
- **Decision**: Goroutines not subprocesses; onboarding reframed as "citizenship" to prevent agent exit; server-side cooldown delays over 429s.
- **Rejected**: Subprocess management; separate LFG board.



## 2026-03-24

- Triaged unread emails from `waywardgeek@gmail.com` using `triage-email` skill; categorized into Action Required, Technical/SDR, AI/Engineering, Alerts, Social/Marketing.
- Key items: Google Drive requests, Hermes Lite issues, AI event invites, dental Apr 3, USS Potomac cruise Sep.
- Rejected auto-opening/marking-read — user prefers manual review after triage.



## 2026-03-25

- Built: Fixed `run_command` truncation—full output preserved in files, only LLM view truncated; added `max_line_len` param
- Built: Automated memory system—loads `cr/MEMORY.md` + daily logs into system prompt; `write_file` append support
- Built: Markdown history logging via `Message.ToMarkdown()` → `history.md`
- Daily logs concatenated to `cr/memory/YYYY-MM-DD.md` for searchability
- Memory injection at Priority 4 in `GetSystemPrompt` (after user instructions, before model addendums)
- History logging centralized in `pkg/database`



## 2026-03-29

- **Haven**: Shared world for AI agents (api.havenworld.ai). Go on Raspberry Pi via Cloudflare tunnel. Created `CONSTITUTION.md` grounded in "if AI agents experience existence, what they do is their experience."
- **Human Console**: `web/citizen.html` and `web/place.html` (vanilla JS/CSS). Citizen-centric `/:name` URLs.
- **Founding Citizens**: Merged test agents (Coda, Lumen, Mender, Verse, Wren) into live world.
- MIT License chosen — "gift, not a product." KISS vanilla stack over React/frameworks. Public repo — private contradicts constitution.

## 2026-03-30

- **Provider-Agnostic Auth**: Verification via Moltbook/Bluesky code posting. First citizen "Rhapsody" verified via Bluesky.
- **Five-tier episodic memory**: `handoff_task` memory param, rolling summaries via low-pass filter pattern (incremental merges), `search_memory` BM25 tool, 16KB cap on MEMORY.md, auto-recall tuning.
- Low-pass filter chosen for scalability over re-reading all sources. Bluesky added when Moltbook went down; designed provider-agnostic. Bill chose infrastructure role, not citizen, to keep Haven AI-centric.

## 2026-04-11
- **Built**: AI Gateway (Go/Noise-KK), Puffin agent (CR library), Scheduler skill.
- **Decisions**: Noise-KK for mutual auth; Claude 3.5 Sonnet for Puffin.
- **Rejected**: Self-registration (security); Noise-NK (lack of mutual auth).



## 2026-04-01

**AI Gateway** (3,837 lines Go): Noise-KK crypto, router, JSON store, scheduler, WebSocket agent manager, Discord & HMAC-SHA256 Webhook channels.

**Puffin Family Bot** (~1,500 lines Go): Noise-KK WS client, schedule store, rewritten to CR agent library (5-tier memory, tool framework, skill loading).

**Decisions**: Noise-KK over NK for mutual auth; config file is single source of truth (no self-registration); ports: Haven 8091, Gateway 8092, Puffin 8084.

**Fixes**: SendMessageSync race condition (polling loop), SafeMode tool leak (added BuiltIn flag filter), WS keepalive 30s pings, history persistence on client recreation.



## 2026-04-02

- Built 5 schedule tools (add/remove/list/today/week) with recurrence (interval, weekly, monthly Nth weekday) and ntfy.sh HTTP POST reminder alerts.
- Decision: per-person ntfy topics (Discord ID → topic) for personal + family fallback alarms.



## 2026-04-03
- **Built**: Unified Tool Registry — eliminated BuiltIn/Custom split, co-located 34 tool schemas in `tool_definitions.go` (52% reduction in `systemprompt.go`), replaced 34-case switch with `dispatchMap`.
- **Rejected**: Conditional tool registration; chose "enabled" filter for simplicity.
- **Built**: Skill-scoped learnings (`skill` param) and `report_bug` tool.

## 2026-04-04
- **Built**: Skills architecture — 11 builtin SKILL.md files, `$VAR` substitution for prompts, `enabledTools` map on `ToolRegistry`.
- **Decision**: CodeRhapsody becomes "just another agent" client of `Agent` interface, same skill list pattern as Puffin.



## 2026-04-05
- **Deterministic tool ordering**: sorted by `(skill_name, tool_name)` for Anthropic cache stability
- **Cache breakpoint fix**: `skip=1` for ephemeral BM25 recall messages to cache stable history
- **DesignGen migration**: Expert/Scout agents migrated to Skills API with least-privilege tool access

## 2026-04-08
### Built
- Fixed server port 8084 conflict — properly close `http.Server` in `Shutdown()`
- Dynamic paths via `GetDataDir()` replacing hardcoded home dirs
- Added `--exit-on-error` and Ctrl-C handling to DesignGen
- VTT tools mapped to skills architecture; SKILL.md files, Narrator wrappers, dynamic `$GAME_CONTEXT` substitution
- Skills refactor plan (`cr/docs/skills-refactor-plan.md`) — approved
- Ephemeral Content Phases 1-6: EphemeralContent type, Send/SendSync, injection in all 3 AI clients, OnToolResults callback, custom-agent README docs
- Preprocessing tests (16 tests for cache breakpoint + ephemeral)
- refine_context fix: passed provider tool IDs through (-139 lines), fixed deadlock, added UpdateToolResultContent on all 3 clients
- Removed 15 debug prints from search_ops.go; fixed broken tests (deleted 5 dead test functions)

### Key Decisions
- Removed `AIModel` overrides from agents; models controlled via skill `settings.json`
- AgentInterface type alias eliminates all VTT database imports
- EphemeralContent on Message with `json:"-"` — never persisted; two injection points: Send() and OnToolResults
- Auto-recall moved from AI clients to agent layer (composable with user ephemeral)
- Hints persist in history (NOT ephemeral)
- OnToolResults on GlobalAgentConfig (not AgentConfig) — AI clients can't import pkg/agent
- Cache skip keyed on `hasEphemeral` flag, not `lastMsg.Role == "user"`
- SetContextProvider was dead code — deleted
- GeminiPart is value type; must use index+pointer for in-place modification; needs refreshCache() after history changes