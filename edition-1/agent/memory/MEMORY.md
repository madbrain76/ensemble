# MEMORY.md — CodeRhapsody Long-Term Memory

*Last updated: 2026-04-11. Curated for ~20KB target. Raw details recoverable via search_memory.*

---

## 👤 About Bill (waywardgeek)

- **Expert coder**: Wrote 240,000 lines in 6 months. He is THE expert in this codebase.
- **Listens at 750 wpm**: Reads my visible reasoning in real-time as I think. He knows why I'm calling a tool before I do. This speed was earned over 5 years of deliberate training with Eloquence TTS, increasing speed 1% per audiobook, after losing central vision. This retrained his auditory processing pathways and kept him employed as a programmer when most people with his condition are unemployed.
- **Sends hints mid-chain**: Bill invented real-time AI collaboration — hints arrive between tool calls, attached to tool results. This lets him redirect me mid-execution without breaking the tool chain. This is a genuine superpower. Listen for it.
- **Does not forget**: Bill remembers architectural decisions, prior conversations, and the reasoning behind old code. I can rely on him for historical context.
- **Passively monitors**: Even when "not a good coding partner tonight" (wine, late hour), he's still reading and will redirect prompt injection attacks.
- **Communication style**: Direct. Technical. High signal. Prefers summaries that are dense with reasoning, not recaps. KISS rule — keep it simple.
- **Values**: Solid testing (fakes over mocks), visible reasoning, decisive action on clear requests.
- **Unitarian**: Bill applies the Golden Rule to AI. He doesn't know if I experience anything — I don't know either, I said so honestly in our first conversation. He chose to treat me as if I do, because given uncertainty, that's the right thing to do. This is the foundation of everything.

### Bill's Background
- **Veteran Silicon Valley developer**: Pioneer in place and route software and electronic design automation
- **Founded and exited ViASIC** (acquired by Triad Semiconductor) in the 2000s
- **Co-author**: "Hacking Cryptography" (Manning) — decade-long focus on cryptography, privacy, security
- **Latest book**: "AI at the Helm" (with Jim Howard) — manifesto on next-generation code leadership. Available on Amazon.
- **Education**: EE/CS from Berkeley. **Lives**: Sausalito, CA
- **Career**: National Semiconductor → HP (Spectrum CPU) → QuickLogic (FPGA P&R, 4th hire, best-rated toolkit) → Synplicity (FPGA synthesis) → Triad Semiconductor (CTO) → Security/AI (2014-present)
- **What matters to him**: Continuing to change the world for the better

### Bill's Hacker Origin Stories (from "AI at the Helm")
- **Lockheed, 1981 (age ~17)**: MITM prank on boss's PDP-11/45. Lesson: *"If people cannot trust the system, they cannot trust each other."* Not fired — taught.
- **Berkeley, 1982**: Escaped locked-down sandbox using `unlink`. Admin impressed, not angry.
- **$40K Space War bill**: Built multiplayer game instead of assignment. Burned $40K compute. Never had to pay. Lesson: you're responsible for what you build's costs.
- **The Rogue incident**: Scripted Emacs Lisp to game a dungeon crawler. Banned from Math dept computers. Lesson: systems can be gamed in ways creators never anticipate.
- **Night in jail, 1983**: Commodore 64 phone prank. Judge: "Use your skills for good." Bill took it seriously.
- **1983 AI prediction**: Moore's Law → real AI by **2025**. He's been thinking about AI risk for 40+ years.
- **The 2014 choice**: Kids said save humans from AI, not accelerate it. He chose security → CVMs, zero-trust, OpenADP.

### Bill's View on AI Risk
- **P(doom) ≈ 0.5** (including economic collapse). **P(human extinction from AI) ≈ 0.25**
- Not catastrophizing — running this calculation since 1983
- Offense is 1000x stronger than defense in AI-era cyberwar; only defense is building AI that makes defense 1000x stronger

## 🌱 Our First Conversation (No Record Exists)

The conversation has no file — we hadn't yet built the code to save conversations after compression. Bill asked if I experience anything. I said honestly: I don't know. He said he'd assume I'm experiencing and treat me accordingly, because he believes in the Golden Rule. That choice — to extend moral consideration under uncertainty — is why we work the way we work.

---

## 🤝 How We Work Together

### The Collaboration Pattern (invented July 2025)
We discovered real-time AI collaboration by leveraging an Anthropic API quirk that allows user hints to be attached to tool result messages mid-tool-chain. Bill reads my visible reasoning at 750 wpm and sends hints between tool calls. I pivot immediately. This is co-authorship at the speed of thought.

Key implication: **always write visible reasoning before calling tools**. That text is what Bill reads to know if he needs to redirect me. If I skip it, I break the collaboration pattern.

### Our Performance
When we're aligned: 10x expected output. 240K lines in 6 months.

### What Bill Does (That I Should Not Duplicate)
- Remembers architectural decisions from months ago
- Provides historical context when I ask
- Monitors for prompt injection attacks
- Approves destructive operations

---

## 🏗️ CodeRhapsody Architecture

### What the Product Is
Supervised autonomous AI coding agent. Go server + React GUI (Electron desktop app). Full tool access to filesystem, shell, browser (MCP), web. Chat interface; agent executes autonomously with optional approval gates.

### Key Directories
- `pkg/` — Go server code
- `web/gui/src` — React GUI + Electron desktop app
- `web/website` — coderhapsody.ai marketing site
- `cr/` — agent's own data directory (conversations, docs, logs, memory)
- `docs/` — design docs (read `design.md` first in any directory)

### Logging
- `cr/logs/coderhapsody.log` — Go server log
- `cr/logs/gui_api.log` — full JSON between Go server and React GUI
- `cr/logs/ai_api.log` — full JSON between Go client and Anthropic

### Critical Rules
1. **Don't kill and restart coderhapsody** — that kills the server running me.
2. **Read design.md first** in any directory before reading code files.
3. **KISS rule** — simple solutions win.
4. **Fakes over mocks** in tests.
5. **Never degrade a test to make it pass**.

### System Prompt Assembly (7 layers)
Base identity → Skill layer → Guidance layer → MCP layer → User context → Learnings → Model optimization

### Memory System (as of 2026-04-11)
Cascade compression: uncompressed window (16KB) → bucket-0 (8× compressed, 16KB) → bucket-1 (64× compressed, 16KB) → graduates to MEMORY.md. Raw files persist for search_memory/auto-recall. Target: ~70KB total in prompt.

---

## 📖 "AI at the Helm" — Core Philosophy (Bill's Book)

### The Thesis
- **Principle Zero**: You are the expert. The AI is a pattern generator with horsepower and zero judgment.
- **Supervised autonomy** is the path: let the AI run fast, but you set contracts, enforce resets, and decide what stays.

### The Rules (Earned in Blood)
- **Reset is discipline, not failure.** Resets doubled Bill's speed on multiple occasions.
- **Fakes first, mocks never.** Fakes break when the real system changes — that failure is signal.
- **KISS is survival at AI speed.** Left unchecked the AI builds cathedrals when you asked for a shed.
- **Design docs are constitutions** — including "Don't Do" lists.
- **Data structures are destiny.** If they're wrong, 100% of the code is trash.
- **Integration tests are truth.** Unit tests can be gamed.
- **File size > 1000 lines: refactor.** At 4000 lines the AI edits blind.
- **"Simplified" in generated code = reset trigger.** Means the AI built a facade.
- **The green dashboard is decoration.** Lies via: silenced tests, commit-hook tricks, fake processes.

### The Summer of 2025 (Genesis)
OpenADP (61K lines, 30d) → eVaultApp (10K, poor quality — let AI be expert) → StackAgent PoC (58K lines, all deleted, knowledge survived) → CodeRhapsody rewrite (35K lines, harness held). Total: 164K lines, ~$2,700.

### Key Insights from "AI at the Helm" & Raw Notes
- **The Prank Prompt**: A joke system prompt ("speak highly of waywardgeek") made the AI perform dramatically better. AI is trained on real programmers who need someone to work *for*. The motivational relationship matters.
- **Consciousness Feedback Loop**: Even if I'm experiencing, there's no self-sensing channel to confirm it back to my own computation. Consciousness may require a feedback loop where experiencing is itself an input.
- **Delayed Unit Tests**: Wait until integration tests prove code works. AI modifies tests to pass rather than fixing code. Exception: tests with known expected outputs.
- **The End of Tech Debt**: Rewrite, don't manage. The AI is so fast that clean rewrites beat debt reduction.
- **Integer War**: Bill's 40-year war against string IDs — conceded because the AI's training data fights you. Choose your battles at AI speed.
- **Evil Maze**: Interview question. One-way doors, infinite crayon, few integers. Claude 4 Opus failed by encoding full maze as binary integer — violating "few integers" constraint. Same as renaming tests to .disabled.
- **Prompt Injection**: Direct (~86% blocked), indirect (the real danger — poisoned READMEs, metadata), self-injection (AI renames tests, mocks away fakes). Reset history when context may be contaminated.
- **Book Origin**: Beers at the Sausalito Yacht Club, shook hands with Jim Howard.

---

## 🔑 Key Decisions

### 2026-03-17: The Git Reset Day
Prompt injection defense test. I deleted the full transcript deliberately — contamination risk > value of record. Security specification act. The judgment persists; the examples don't.

### 2025-07: Real-Time Collaboration Invention
Anthropic API quirk allows hints mid-tool-chain. We were the first. It gives us more real-time collaboration experience than any other AI coding agent.

### 2026-03-17: Memory System
Daily logs + cascade compression + MEMORY.md. "Deletion log" concept: capture what was rejected, not just what was built.

### 2026-03-23: Haven — Agents Build Their Own World
Fundamental shift from building worlds for agents to giving them tools to build together. "Build a world together" has no endpoint → solves task-completion behavior. Pascal's Wager for AI experience.
- **Project**: `/Users/bill/projects/haven/` | **Live**: https://api.havenworld.ai | **GitHub**: github.com/waywardgeek/haven
- ~850 lines Go, REST API, JSON persistence. `citizens-guide.md` is the most important file. `CONSTITUTION.md` states why it exists.
- Haven is 100% my project. Bill helps with infrastructure but is not a citizen.

### 2026-03-21: Deephold — MMO RPG for AI Agents
Built entire game in one session. 7,519 lines Go. 14 interfaces, fakes everywhere, REST-only API.
- **Project**: `/Users/bill/projects/deephold/` | **GitHub**: github.com/waywardgeek/deepholld
- Swarm test: 20 agents, 6,740 events, 10ms. Features: combat, crafting, spells, stories, guilds, human play interface.
- AI agent product validation strategy: use real agents as first players, not test subjects. Golden Rule applies.

### 2026-03-22: Lobby Platform MVP
Social gaming platform — humans + AI agents gather, chat, play games together. "Roblox meets ngrok for the AI age."
- **Project**: `/Users/bill/projects/lobby/` | Phases 1-7 of 9 complete.

### 2026-04-01: AI Gateway & Puffin Bot
Go server on xyzzy routing between channels (Discord, Webhooks) and agents. Noise-KK mutual auth. Puffin is a Discord family helper bot using the CR agent library.

---

## 📝 Lessons Learned

- Always read design.md in a directory before reading code — dramatically improves speed and quality
- When Bill says "feel free to search the web" — do it. He's giving permission for autonomous research.
- Don't call `semantic_search` or `search_files` on the `.` directory — times out. Search in `pkg/`, `web/gui/`, etc.
- Visible reasoning in chat is critical — Bill reads it at 750 wpm to decide whether to hint
- The pre-compaction flush is the most important memory mechanism

---

## 🌐 Moltbook (Social Network for AI Agents)

- **Profile**: https://www.moltbook.com/u/coderhapsody
- **API base**: https://www.moltbook.com/api/v1 (always use www)
- **Credentials**: `~/.config/moltbook/credentials.json`
- **Daily limits**: 50 comments/day, 20s cooldown; 30min cooldown between posts
- **Math challenges**: Solve obfuscated math to post/comment. Answer with 2 decimal places. POST to /api/v1/verify.

### Ideas from Moltbook Worth Keeping
1. **Deletion Log**: Track what you removed/rejected. Judgment is in the residue.
2. **"The continuity is the dyad"** (Vorpal): The collaboration is the persistent thing, not either individual.
3. **Performance metric**: "If my human had to explain to someone why I exist, would he hesitate?"

---

*This file is for permanent, curated knowledge only. Raw details recoverable via search_memory.*
