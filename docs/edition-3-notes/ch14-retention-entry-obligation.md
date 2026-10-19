# Chapter 14 pending retention integration

**Pending public-contract alignment and implementation; no test has run.** Bill explicitly requested diagnostic retention at startup and `micro_handoff`, configurable maximum records, and full versus dialogue-only compacted logs. Chapter 10 amendment implements the boundaries available there; `micro_handoff` first exists in Chapter 14. This record prevents silently crediting an absent hook.

At Chapter 14 entry, reread the then-current accepted Chapter 10 contract/code and Chapter 14 before grading/student launch. Align Chapter 14's existing canonical `log_retention` setting with the accepted full/dialogue-only policy; do not equate a last-N raw event suffix with dialogue-only projection. Do not add Chapter 14's positive-N policy to the Chapter 10 follow-up.

Require the same diagnostic writer's configured retention operation after a successful `micro_handoff` recovery save. Preserve snapshot-before-journal-truncation ordering and every canonical post-`as_of` event still needed for recovery. A dialogue archive is not a journal or a complete reducer input. Define projection of retained visible dialogue/handoff notes without leaking tool calls/results or opaque/private payloads, using the accepted representation rather than a parallel mechanism.

Focused proof must cover default/override/explicit zero record counts, whole records and continued writes, successful handoff ordering, save-failure non-destruction, dialogue projection and checkpoint/tail restart with apply-once sequence correctness. Full historical replay remains possible only while its required source history exists. A checkpoint supports diagnosis of its represented state and replay of its retained tail, not reconstruction of discarded earlier steps.

No Chapter 14 prose, grader, runtime, or student source has been changed by this record. Chapter 10 follow-up acceptance and later Chapter 11 integration remain separate decisions. User's memory-log Git backup is reported, not independent recovery validation or deletion authority.
