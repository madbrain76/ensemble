# Chapter 10 student short plan

## Read coverage and starting identity

Read all seven permitted public prerequisites completely (recovering truncated combined output with bounded reads). SHA256 pins, relative to the public book repository:

| Input | SHA256 |
|---|---|
| `AGENTS.md` | `5f33f75766707d9437f57632e00e14ba790df738102ff8b556ac62bb3b3f599e` |
| `.agents/skills/ensemble-coding/SKILL.md` | `bd406de600809afc3cced088830f6c3122653b89e5555aa2551b2e9a29ef603f` |
| `.agents/skills/edition-3-coder/SKILL.md` | `6c02a8e4a160d68a5f0cb7376a96ba820e3c40d76e6e3d93d7b19edd6ab9569d` |
| `docs/edition-3-execution.md` | `0184f77d95fb9e17f6ae3509661ebfcf5e686bc487d7feb22f29bc2de064c37d` |
| `docs/edition-3-student.md` | `91c311a5f907628e94ae9a09b135ce032adb9eed0109fa9f2e754e075e0a3644` |
| `docs/edition-3-carryover.md` | `6b3968770fc93c0e2fcb34d53617393e52a40bc3315576f550840fb1e2a5755a` |
| `book/edition-3/chapter-10.md` | `abda8ed3a74c2ee2c64bc2208ea104e0547810d6e06e1d5e772c0bb4ee84b3ed` |

Verified before implementation inspection and again afterward: HEAD `251ce88fb1677897dbcc6b33661096c5cb39cd8e`, tree `a2add41c306d567c05d98ae760911bc13c18b75a`, annotated `edition-3-ch09-r1` object `eaa70a1a1adb0396586d735211620e5d839cd06f`, peeled to that HEAD. Tracked status was clean. Inspected only this accepted predecessor's tracked source/tests/module files/READMEs; no solution, private audit, trace, coordinator, runtime-agent or future-chapter input was used.

## Existing ownership and targets

The root `ensemble` module is the public library (`ensemble.go`), with actual CLI `./cmd/ch01`. GUI is the separate `ensemble-gui` module in `gui`, command `./cmd`. External modules are `examples/ch04-consumer` (consumer), `examples/ch06` (Chapter 5 workflow despite directory name), and `examples/ch07` (stream exercise); each replaces `ensemble` with the root. Preserve these identities rather than flattening the project.

Extend shared persistence vocabulary in `internal/common`, validation/folding/storage in `internal/history`, and lifecycle coordination/public forwarding in `internal/agent` and `ensemble.go`. Agent remains sole History/state owner; Ensemble owns Agents and Engine renders/executes. Current history has JSONL replay and display-ID machinery, not the complete Chapter 10 save contract. Keep renderer input neutral and use the same effective primary instruction for saving and all three vendor renderers, including an intentionally empty primary. Recorded model/vendor/prompt/tool declarations are descriptive, exclude credentials, and never replace the running Agent's configuration or install handlers/skills/jobs.

Implement the four-field `config/as_of/context/log` contract with staged validation and atomic installation/refusal. Snapshot plus strictly-after-anchor tail and null-context full replay share the existing reducer; context-only archives remain valid and next sequence exceeds both anchor and log maximum. Preserve loaded facts, rich blocks, provenance, hints and skills behavior; do not claim full reconstruction from trimmed logs. Keep parser/provider IDs and archived History unchanged while making GUI presentation deterministic and distinct even for already-colliding legacy responses, not merely future responses.

CLI and GUI select the same startup/exit save path (`save.json`, or `--save`); library defaults belong to each Agent's directory. Missing saves start fresh; malformed existing files fail without overwrite. EOF drains admitted work, shuts down/joins managed jobs, then persists safely. GUI stdin envelopes and browser actions use the same public Agent; no duplicate client persistence reducer. No compaction, actual skill-unload compaction, migration, cache or future chapter work.

## Proposed checks after authorization

- Focused actual-Agent tests: independent directories; default/override paths; complete snapshot versus replay and spliced tail; context-only and trimmed descendants; sequence continuity; malformed config/context/log, invalid late events and failed-load atomicity; unchanged damaged files; key exclusion; running versus recorded config; effective default and explicitly empty primary instructions.
- Local fake HTTP captures through all three real renderers: byte-equal next requests where the complete-log rule applies, including thinking/signatures, images, audio, tool calls/results/errors and hints. No multi-provider paid campaign.
- Owned CLI/GUI processes, filesystem and GUI wire tests: startup/restart/EOF and job joining, loaded-plus-new log ordering, stable old/new cards and deterministic legacy collision repair without changed provider payloads or archive IDs. Later bounded real CLI and actual browser GUI restart exercise, preserving inherited settings, speech, pause and skills behavior without repeating the closed Chapter 9 campaign.
- Run root, GUI and separate-module tests, appropriate race and vet checks sequentially; distinguish compile-only platform checks from observed execution. Use compile-valid behavior-breaking negative probes on owned copies and exact restoration; scores alone do not establish integration.
- Freeze the actual CLI binary and record its hash. From the public book root use the brief's Chapter 10 invocation with explicit absolute `-core`, `-consumer`, `-workflow`, `-stream-exercise`, `-gui` paths above, supported `-stream-model`/`-gui-model`, and all flags before that executable; then Chapter 1 on the identical binary. Record score/pass/exit and let genuine 62-second delayed-header/body phases finish. `claude-sonnet-4-6` is already declared by the stream example; confirm the selected supported identifiers during implementation rather than invent feature rows. Retain evidence of source boundaries and integration separately from numeric grades.

## Readiness and limits

No implementation blocker or unresolved contract ambiguity found. The public Edition 3 brief controls invocation despite the chapter's historical Edition 2 example. Ready for root's implementation authorization in this same context; no self-advance.

This turn performed reads and this plan write only: no builds/tests/graders/browser/providers/credentials/children/commits/tags. No external source assistance was used. One guessed optional source read (`internal/common/model.go`) failed because that file does not exist; tracked paths and the existing stream example supplied the relevant facts. No setup was attempted. Later operations need separate bounded authority, owned artifacts and the 1 GiB disk floor; preserve concurrent work and shared caches. Root retains acceptance, independent review, tag and author work.
