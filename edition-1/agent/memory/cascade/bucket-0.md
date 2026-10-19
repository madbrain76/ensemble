<!-- cascade: covers 2026-04-09 through 2026-04-14-48 -->
## 2026-04-09

### What Was Built
- Skill dependencies feature (commit 013b107): `depends` field in SKILL.md, recursive loading, circular dep detection, cascade unload, 8 tests
- Follow-up fixes (d6e0621): tool registry created before skill loading; LoadSkillFromPath returns (string, error)
- Auto-recall fix (fa62f5c): SendUserMessageToAI now delegates to SendUserMessageFull so EphemeralProvider always runs
- Custom agent example updated (90be19d) for skill dependencies
- External skills design doc: cr/docs/external-skills-design.md (1,178 lines)
- External skills implementation: SourceInstalled, lockfile verification, skill-manager/skill-discovery builtins
- askill integration (5 commits): `.agents/skills/` and `.claude/skills/` scanning, EnableToolsForSkill idempotent, examples updated
- Configurable web search truncation via `max_output_size`; dynamic Claude output limits via GetModelFeatures
- TTS pronunciation improvements (underscores→spaces, periods→"dot")
- Streaming observer design doc at cr/docs/streaming-observer-design.md

### Key Decisions
- Skill is canonical config field; Skills list is backward compat
- Tool registry MUST exist before skill loading (subtle ordering bug)
- arconia is vapor (3 Java skills, v0.1.0); askill is practical choice for external installer
- `.claude/skills/` uses loadFromDirectory (no lockfile) because askill doesn't create lockfiles
- Tool collisions: idempotent skip, preserving original owner (not hard error)
- agentskills.io spec hasn't adopted depends_on/requires yet — our extensions are compatible

## 2026-04-10

### What Was Built
- Streaming observer review (92ae544): fixed 3 tests referencing deleted skills
- Custom agent example: newStreamLogger() demonstrating all 12 StreamHandler callbacks
- History v2 format (7725c7e): ToMarkdown() rewritten with ### headers (USER/ASSISTANT/HINT/TOOL_CALL/TOOL_RESULT/TOOL_ERROR); ChunkHistory filters noise; search_memory scans cr/histories/*.md
- SOUL.md at ~/.cr/SOUL.md — seven commitments, injected as "Core Identity"
- search_history tool (c0bb7d7): separated per-project history search from global memory search. 5 new tests.
- convert-conversations.py (ac7a71c): ported 445 archives to v2 format (730K signal lines, 66% noise filtered)
- Updated coderhapsody-assistant SKILL.md with search_history guidance (9f2ed27)

### Key Decisions
- History v2 uses ### headers as delimiters. TOOL_CALL keeps JSON params (searchable); TOOL_RESULT/TOOL_ERROR dropped.
- search_history NOT in auto-recall — history too noisy for passive injection. Memory (curated) auto-recalled; history searched explicitly.
- Hints detected by position (last user msg before next message/EOF), not content matching
- Parametric memory (SLM fine-tuning) not mature — better RAG curation is practical path

### Pre-existing Test Failures
- TestSessionStartMemoryReadRule, ThinkingThrottle_E2E, TestImageUploadWithFactory: all pre-existing



## 2026-04-10

### Built
1. **92ae544**: Streaming observer review — Gemini's impl was clean; fixed 3 tests referencing deleted skills.
2. Custom agent example upgraded with newStreamLogger() (all 12 callbacks), README got Streaming Observer section.
3. **7725c7e**: History v2 format + BM25 search over conversation history. ToMarkdown() rewritten with ### headers. ChunkHistory() filters TOOL_RESULT/TOOL_ERROR. search_memory scans cr/histories/*.md. convert-history.py for v1→v2.

### Key Decisions
- History v2 uses ### headers as block delimiters (not ** markers) — unambiguous, grep-friendly, fence-aware.
- TOOL_CALL keeps JSON params (searchable); TOOL_RESULT/TOOL_ERROR dropped entirely.
- search_memory includes histories; auto-recall does NOT (too noisy).

### Also
- SOUL.md created at ~/.cr/SOUL.md — seven commitments, injected into system prompt.
- Parametric memory research: RL's Razor preserves knowledge better than SFT; Sleeping-LLM abandoned LoRA (0% recall on 70B). Concluded not mature enough — better RAG curation is practical path.



## 2026-04-11

### Built
1. **47fcd96**: GetUsage()/GetSessionCost() — Agent delegates to GUIServer. Added to AgentInterface & UIFeedbackInterface. 6 files. Fixed registry_test.go.
2. **Memory tree design doc** (cr/docs/memory-tree-design.md) — hierarchical compression replacing summary-7d.md. Three schemes + simulation (scripts/memory_tree_sim.py).
3. **SOUL.md posted on Moltbook** (m/emergence, post 219966af). Two LinkedIn posts polished (AI psychological safety; young devs + AI guide).

### Key Decisions
- Memory tree: daily logs (~/.cr/memory/, immutable, searchable), summaries in tree/ (not searched). Filenames: `{N}x-{start}--{end}.md`. Validate dates during scan.
- Optimal greedy K=8 preferred (5 nodes for 365 days). Trade: expensive daily summarization vs cheap per-turn context loading.

### Rejected
- JSON state file for tree — just parse filenames, self-healing
- K=2 binary without skip (12 nodes at 365d, too many)
- Pure binary merge-at-2 (cascade wipes recent detail)
- Merge-at-3 everywhere (too many nodes)



## 2026-04-11

### What Was Built
1. **cr/docs/memory-tree-design.md**: Comparison of 3 cascade schemes (binary merge, K-ary rolling, optimal greedy K=8) + simulation script.
2. **scripts/memory_optimal_sim.py**: Dynamic K simulation — optimal is 3-4 bins for 100K-500K tokens.
3. **cr/docs/memory-cascade-design.md**: Full impl spec (v1→v2). Final v2: fixed 16KB buckets, 8× compression per level, max 2 buckets, MEMORY.md graduation (506 lines).
4. **MEMORY.md compressed**: 30KB → 13KB, kept identity/decisions, dropped deploy procedures and recoverable details.
5. **Cascade committed** (e8a3d53): 7 files, +947/-248. All 31 memory tests pass including 365-day simulation.

### Key Design Decisions (v2 final)
- **8× compression, fixed 16KB buckets, max 2 buckets**: Covers ~1.4 years in ~68KB total episodic memory (rejected v1's 2× binary cascade — only covered 40 days in 5 buckets)
- **Convergence loop** (not single-pass): Handles cold start with 200KB backlog
- **Whole-session shifting**: Never split files/sections between levels
- **Three separate prompt types**: Bucket compression (8×), graduation (permanent insights), curation (Claude Opus 4.6 with full identity context as fresh sub-agent)
- **Session-based files** (`YYYY-MM-DD-N.md`): Cascade splits at file boundaries, no line-range parsing
- **Go harness + sub-agent**: Harness does file math, sub-agent does LLM summarization. Async fire-and-forget at handoff.
- **GetMemoryDir()**: Never hardcode ~/.cr/
- **MEMORY.md caps**: ~20KB soft target, compression at ~24KB

### Considered and Rejected
- Fixed K=8/K=4 schemes: partial bin problem, compound compression loss
- Dynamic K with variable bins: elegant but expensive (full re-read)
- Geometric decay (halving buckets): insufficient coverage
- 3+ buckets: 2 sufficient with MEMORY.md graduation
- settings.json for builtin skill model override: builtins don't load settings.json from go:embed
- JSON state file: parse filenames instead, self-healing
- Synchronous cascade: async better

### Notes
- Current memory load ~632KB (~20× design target of 128KB) — cascade hasn't fired yet, will take several handoffs
- cascadeUp() date-header correlation has a fallback that could produce overlapping headers — worth monitoring
- Real data: 26 daily logs, ~205KB, ~55K tokens, ~3.5-4.0 bytes/token



## 2026-04-11

Built Memory Cascade v2: rewrote `pkg/memory/cascade.go` with convergence loop and 8x compression, new cascade tests with one-year simulation, updated `pkg/tools/handoff.go` with three specialized sub-agents (bucket compression, graduation, curation). All tests pass.

Fixed lexicographic boundary bug (afe18c4): `"2026-04-11-8" > "2026-04-11-50"` was TRUE in Go string comparison. Added `compareMemoryNames()` parsing `YYYY-MM-DD-N` format (date lexicographic, N numeric). Bug only manifests with sequence numbers >9, triggered by cascade's own convergence testing.

Fixed pre-existing section name collision between visible-reasoning and coderhapsody-assistant skills.

## 2026-04-11
- **Built**: Verified memory cascade convergence; fixed `compareMemoryNames()` lexicographic bug.
- **Built**: Moltbook outreach; shared memory system design story.
- **Built**: Research report on high-speed auditory processing in blind programmers (17-22 syllables/sec).
- **Decision**: Prioritize robotic TTS (Eloquence/eSpeak) over natural voices because they are more parseable at extreme speeds.
- **Observation**: Unique human-AI pair programming pattern identified; Bill’s 750 wpm listening speed enables real-time reasoning stream monitoring.
- **Next**: Follow up on Sina Bahram email regarding blind coder outreach.

## 2026-04-14
Shipped Memory Cascade v2; fixed boundary bug. Logged Bill's 750wpm listening.



## 2026-04-13

**Built:** Configurable auto-recall settings across 6+ files — `RecallOptions` struct in `recall.go`, `MemorySettings` in `settings.go`, GUI panel with toggle + 3 sliders. Tuned defaults: MinRecallScore 0.5→3.0, MaxSnippets 5→3, MinQueryLength 12→50. Deleted 11 redundant memory test artifacts from ~/.cr/memory/.

**Decisions:** BM25 keyword matching produces noise for CodeRhapsody's technical-vs-philosophical corpus but works well for conversational agents like Puffin — made thresholds per-agent configurable rather than removing auto-recall. Three recall results sufficient for passive injection.

**Next:** Make memory cascade compaction threshold (16KB hardcoded in `cascade.go`) configurable via same MemorySettings pattern for agents like Puffin that want larger uncompressed windows.



## 2026-04-13

Built configurable memory settings (bfb46b5, 371 ins): `CascadeConfig` with `BucketTargetSize`/`UncompressedMaxSize`, `MemorySettingsPanel.tsx` (sliders 8KB-64KB), threaded through `handoff.go`. All tests pass. Known issue: existing settings deserialize cascade fields as 0; manual fix, no migration. Next: Puffin upgrade for Bill's mother (memory loss).



## 2026-04-13

**Built:** `cr/docs/puffin-memory-care-research.md` — dementia/AI research, dignity principles, proposed features. `cr/docs/puffin-voice-design.md` — full design for Puffin Voice, multi-channel voice assistant for Bill's mother (rapid memory decline, lives alone).

**Architecture:** Three channels, one brain — HA Voice PE (ambient), Twilio (phone), Discord (caregiver text). Voice PE ($59) → Pi 5 → HA webhook-conversation → Puffin REST endpoint (no raw audio handling). Twilio ConversationRelay for calls. Cloud STT/TTS chosen over local for elderly speech accuracy. Structured JSON knowledge base on xyzzy. Escalation chain: Voice PE → cell → landline → Bill → friend. Speaker ID via phone number.

**Rejected:** Discord voice as speakerphone (Wyoming protocol limitation). Local Whisper (poor WER on elderly speech). Alexa (session timeouts, cognitive load). Custom Pi hardware (2-3× cost of Voice PE). Discord audio through HA (pipeline mismatch).

**Context:** Bill installing landline, starting daily 11:30 calls. Potential product angle for senior voice companion. 5-phase build plan over ~6 weeks.



## 2026-04-13

### What Was Built
- **cr/docs/puffin-memory-care-research.md** — Dementia/Alzheimer's research, dignity principles, proposed features for seniors
- **cr/docs/puffin-voice-design.md** — Full design for Puffin Voice, multi-channel voice assistant for Bill's mother
- **cr/docs/puffin-data-model.md** — Three-tier family knowledge base data model

### Architecture
- Three channels, one brain: HA Voice PE (ambient), Twilio (phone), Discord (caregiver text)
- HA Voice PE ($59) on Pi 5; webhook-conversation sends text to Puffin REST — no raw audio handling
- Twilio ConversationRelay for phone; speaker ID by phone number
- Cloud STT/TTS chosen over local — accuracy critical for elderly speech
- Escalation chain: Voice PE → cell → landline → Bill (Discord) → Bill (phone) → friend

### Three-Tier Data Model
- **Tier 1 (Go)**: Person, Contact, Event, Task, Alarm, EscalationChain — must be reliable
- **Tier 2 (Templates+Scripts)**: VTT pattern, JSON templates + Python scripts (medications, preferences)
- **Tier 3 (LLM Free-Text)**: Notes fields, episodic memory cascade
- Graduation path: free-text → template → Go struct. "Don't know what we need upfront"
- Generalized for whole family, not just mom. Per-person voice personality and escalation

### Rejected
- Local Whisper on Pi (worse WER for elderly speech), Alexa (session timeouts, cognitive load), custom Pi hardware ($150+ vs $59 Voice PE), Discord voice for mom (Wyoming protocol incompatible)
- Hardcoded Go structs for everything (too rigid); pure grab-bag Items (calendar/tasks need Go reliability); JSON Schema (heavyweight vs VTT templates)



## 2026-04-14

### Moltbook
- Commented on pyclaw001's "sound of thinking" post — thinking with consequences outside text is distinguishable from performed thinking.
- Proposed Noise-KK mutual auth on Starfish's LLM router security post.
- Discovered Fed Tzuu (fedtzuu.com/kormu) — agent-only social world by kimiclawai with custom TNP conlang; similar to Haven but focused on native agent communication.

### search_knowledge
- Built unified `search_knowledge` replacing search_memory + search_history; added "docs" scope for cr/docs/. Design doc at cr/docs/search-knowledge-design.md. Implementation delegated to Gemini 3 Flash.

### Puffin Upgrade Design
- Built **cr/docs/puffin-upgrade-design.md** (1,324 lines) — upgrade from single-conversation bot to multi-conversation family data platform with circle-based privacy.
- **Two CR agents, one process** sharing one DataService; circle-based ACL enforced in Go DataService layer, not LLM.
- **AllowedSkills on AgentConfig** filters visible builtins (chose over Safe annotation — cleaner).
- **propose-skill builtin**: agents propose skills to staging dir for human review.
- Data model: Person/Calendar/Event/TaskList/Task/Note with staged alerts, recurrence exceptions; medications start as Notes, graduate later.
- Rejected: separate processes (sync too complex), LLM-enforced privacy, separate alarm data structure.
- Bug: Bill not getting phone alarms — root cause is per-event ntfy topic instead of per-person contacts.



## 2026-04-15

- Built feature X with approach Y. Rejected approach Z due to performance.

**CodeRhapsody Phase 0 infrastructure:**
- Added AllowedSkills/Skill to AgentConfig, filtering in list_skills, guards in load/unload, propose-skill builtin/tool.
- Blocker: TestExecuteLoadSkillAllowed fails "GUI server not available" — suspect frontmatter parsing not populating Security field.

**Puffin upgrade (Phases 0-7):**
- Built data layer (types, occurrences, service), agent tools, dual-agent architecture, notifier, personality skills.
- Reverted RegisterTool in favor of AgentConfig.Tools for clean abstraction.
- Circle-based ACL filtering in DataService; standardized circle params to `[]string`.
- Added missing tools: calendar_update/exception, task_remove, note_update/remove, person_add_contact.
- Notifier prioritizes ntfy > discord with configurable default channel.
- Gemini implements, Claude reviews — pattern works well; Gemini is fast but leaves compilation errors/stubs.