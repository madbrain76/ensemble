# Follow-up persistence lesson requested by Bill

After Chapter 10 closure, Bill requested that this operational lesson be included in the persistence chapter. This note preserves the request for a bounded author amendment; it is not a manuscript revision or an implemented-feature claim.

Observed in the running v2 environment: diagnostic `api.log` and `gui.log` occupied about 3.47 GiB and 1.66 GiB respectively. The Go build cache separately occupied about 9 GiB. Bill authorized clearing only the build cache with `go clean -cache`; root did so, restoring roughly 10 GiB free space without removing modules, source, evidence, logs, memories, saves, or journals.

Teaching points for the persistence chapter, kept brief and separate from new student requirements:

- Durable recovery state and disposable diagnostic traces serve different purposes. Do not use saved-history retention or compaction controls to solve diagnostic-log disk growth.
- Diagnostic retention should be a separate configurable policy, counting complete records rather than arbitrary lines/fragments. Bill proposed the most recent 200 records; suggested semantics are default 200, positive values retain that many, zero retains all, negatives invalid. This feature is not yet implemented.
- Prefer retention at startup and after a successful durable checkpoint, never before the save is secured. A save failure must not trigger destructive diagnostic cleanup.
- Coordinate trimming with the writer. Live files opened without append mode can resume at old offsets after external truncation, leaving holes or corrupt boundaries. Replacing an open file's pathname may also leave the running process holding the old inode, delaying reclaimed space.
- Preserve `save.json`, recovery journals, memories, and worker state. Diagnostic trimming is not conversation compaction, and neither establishes semantic fidelity by itself.
- Caches are regenerable but build headroom and cold-cache cost still matter. Do not delete shared caches without authorization.

Scope of the future author change: one short practical passage in Chapter 10, based on these observed operational facts; no invented claim that the accepted student implemented this policy, no retroactive exercise/grader requirement, no long harness diary, and no reopening unrelated chapter polish. Root remains orchestrator; use a fresh bounded author and the appropriate review before committing manuscript changes.

The separately authorized current-source GUI rebuild installed at `/Users/bill/ensemble-edition2-codex/bin/ensemble-gui` is recorded in `/tmp/cadence-current-gui-rebuild-receipt.json`. It does NOT implement the proposed retention setting, does not change live settings, and did not restart the running process. Build/help checks passed, not a full runtime regression suite. The first root rebuild overlapped an already-started student control due to a preflight/message acknowledgment race; the student preserved the valid first mutation/restoration result and a wrapper-signature failure, with the second control unrun. Both build jobs are finished and the student is held for Bill's restart.
