# Rolling 7-Day Summary
*Last updated: 2026-04-11*

## What Was Built
- [04-06]: Unified Tool Registry — dispatch maps replace 34-case switch, systemprompt.go -52%.
- [04-07]: Tool Skills — 11 SKILL.md, enabledTools, $VAR sub, report_bug, cache fixes.
- [04-08]: Ephemeral Content 1-6, refine_context 3-bug fix, dead code purge (-1250), CR port fix.
- [04-09]: Skill dependencies — `depends` field, recursive load, cascade unload, 8 tests.
- [04-09]: External skills — askill integration, .agents/ + .claude/skills/ scanning, lockfile verify.
- [04-09]: Auto-recall fix, deleted unused builtins, TTS pronunciation, streaming observer design.

## Key Decisions
- [04-07]: enabledTools nil=all. Order by (skill,name) for cache stability.
- [04-08]: Memory/learnings skill-controlled. Cache skip `hasEphemeral` not `role=="user"`.
- [04-09]: Skill is canonical config; one primary skill, depends handles the rest.
- [04-09]: Tool registry MUST init before skill loading (ordering bug).
- [04-09]: askill over arconia — arconia vapor; askill practical (73 commits, safety scoring).
- [04-09]: .claude/skills/ no lockfile (askill doesn't create them); tool collisions → idempotent skip.

## Considered and Rejected
- [04-07]: ToolGroups — enabledTools simpler.
- [04-08]: Keeping CRSystemPrompt/multi-edit — skill is source of truth.
- [04-09]: arconia as installer — vapor. Hard error on tool collisions — idempotent skip better.

## Active Context
- VTT skills refactor approved, ready for 3-phase implementation.
- Streaming observer design ready (~200 lines). agentskills.io depends proposals open (#100, #110).
- Test hangs: ThinkingThrottle_E2E(3x), TestImageUploadWithFactory, TestSessionStartMemoryReadRule.
