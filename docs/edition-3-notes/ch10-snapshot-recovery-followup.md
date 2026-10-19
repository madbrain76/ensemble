# Chapter 10 boundary review — snapshot-based recovery retention

Bill requested this review at the next natural boundary between chapters; do not interrupt current Chapter 2 work. Continue autonomously afterward.

Proposed policy: preserve the current valid snapshot as the prior/backup snapshot, write the new snapshot, and only after success trim the recovery event log to retain all events needed since the backup snapshot. Load the newest valid snapshot first; on failure, try the backup and replay its retained suffix. Retention follows the recovery boundary rather than an arbitrary 200-event count.

At review, distinguish the existing `diagnostic_max_records` default (diagnostic records) from canonical recovery events and the dialogue-only archive policy. Check atomic replacement, snapshot validation, save/trim failures, interruption between steps, sequence/as_of replay exactly once, and preservation of sufficient events for backup recovery. Do not assume copying corrupt current bytes over a valid backup is safe. Determine the smallest coherent teaching change and its chapter ownership; no new implementation or acceptance is claimed here.

## Boundary review after restarted Chapter 2

Read current Chapter 10 in full. Its canonical log lives inside `SaveFile`; an independent continuously appended journal is introduced in Chapter 14. The 200-record setting currently controls diagnostics, not canonical recovery events.

Adopt backup-anchored recovery retention when aligning the public contracts before Chapter 10/14 execution:

- Chapter 10: keep the previous **validated** save under a backup suffix; fully encode the new save before atomic replacement. Try the primary first, then validate and load the backup on primary failure, with a visible recovery diagnostic. Never rotate a corrupt primary over a usable backup. If neither is valid, refuse without overwriting either. A missing primary with a usable backup is recovery, not a fresh start.
- Preserve coherent snapshot/as_of pairs and per-Agent paths. Failed encoding, backup creation or replacement must leave a valid recovery option; distinguish process interruption from guaranteed power-loss durability. Check first save, successive saves, corrupt newest, corrupt/missing backup and interrupted rotation. No new persistence coordinator or database is needed.
- Chapter 10 alone can fall back to the previous saved state. It cannot recover newer events stored only inside a corrupt primary. State that limit rather than claim an absent independent journal.
- Chapter 14: retain the complete independent canonical journal **after the backup snapshot's anchor**, not merely after the newest anchor. Publish the new snapshot and validated backup before trimming. Load newest-or-backup, then replay the retained tail above the chosen anchor exactly once. Do not trim required events after a failed save, rotation or recovery; serialized writer appends must remain ordered. Replace arbitrary canonical recovery-event counts with this correctness boundary.
- Full archival history remains an explicit separate choice; dialogue-only remains a derived archive, never replay input. Neither restores discarded raw history. Diagnostic traces have no replay anchor, so the configurable diagnostic-record limit is a separate policy—not a competing canonical-event retention rule.

This is the requested review and entry direction, not a new implementation or tested recovery claim. Align Chapter 10's current fatal-primary rule, exercise fixtures and prose before its new student starts; align Chapter 14's journal truncation and existing `log_retention` then. The manuscript remains unchanged at this boundary; continue Chapter 3 now.
