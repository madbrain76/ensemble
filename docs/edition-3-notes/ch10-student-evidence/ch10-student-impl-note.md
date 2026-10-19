# Chapter 10 implementation handoff

Source owner remains student branch/worktree based on accepted `251ce88` / `edition-3-ch09-r1`. Exact freeze identity and evidence outcomes are recorded in `ch10-student-evidence/local-report.json` and `ch10-student-local-report.txt` after the final gates.

## Owners and changes

- `internal/history/save.go`, `history.go`: exact four-field archive, strict detached decode/validation before installation; context checkpoint plus tail reduction through the existing reducer; log-only restore; retained-prefix validation; anchor/frontier continuation and snapshot-only resave. Original JSONL inspection APIs remain separate. Public vocabulary is `SaveFile`, `SaveConfig`, `ToolDeclaration`.
- `internal/agent/save.go`, actor/constructor wiring: public serialized `Save`/`Load`, busy-load refusal, startup load before admission, automatic shutdown checkpoint, informational effective runtime config/tool declarations, atomic sibling-temp replacement. Runtime credentials/settings/handlers/skills/jobs/workers are never restored. Restored display frontier seeds live allocation.
- Unspecified library DataDir gets an independent temporary directory; caller supplies a stable explicit DataDir for cross-root restore. CLI/GUI explicitly retain working-directory DataDir (including inherited `cr/io` semantics) and default working-directory `save.json`; `--save` overrides both load/save. Clients honor current `LLM_SYSTEM`; primary skills remain authoritative, including empty instruction.
- CLI preserves live/rich/JSONL compatibility; offline JSONL commands use isolated lifecycle storage. Pending trailing ephemera survive EOF. GUI canonical kind/text envelopes share the browser Agent, strict aliases are accepted, EOF closes browser admission and drains before checkpoint/exit. `Ephemeral` is explicit durable one-request context without changing `Hint` semantics.
- GUI replay projects missing context-prefix turns and gives repeated legacy display IDs stable presentation-only identities. It never edits stored facts, tool call IDs or provider messages.

## Verification scope / limitations

Tests use local scripted HTTP fixtures, not real credentials/providers. Browser native speech and real UI behavior remain **not observed**. Provider/browser gates remain unauthorized and pending. No compatibility claim is inferred from Go tests alone. Jobs/processes/workers/active skills are not restored. Concurrent writers to one archive are unsupported. Atomic replacement protects prior target bytes against serialization/write/rename failures, not every power-loss scenario. Files referenced by rich parts remain externally owned.

The first black-box development grader returned 65/100: runtime-config and inherited job-path failures plus an omitted media map. The corrected clients explicitly retain cwd job storage, honor current runtime system input, and the final invocation supplies supported model selections. All failed attempts are retained, not replaced by summaries of success.
