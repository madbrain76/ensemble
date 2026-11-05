# Fix: unreadable persisted state says which file it rejected

Branch `fix/state-failure-diagnostics`, based on `fix/linux-path-case`. Coordinator, 2026-11-05.

A launch that could not read its own saved state died as `GUI failed:
AggregateException. Check private diagnostics and workspace/data permissions; no
private error text is printed here.` Following that instruction led to
`ensemble.log`, which said only `persistence: invalid save candidate; trying
backup` — no path, no reason — and to empty `debug.log`/`gui.log`, because the
rejected settings document never logged at all. The privacy policy that hides
exception text is right; the gap is that the sanctioned diagnostics channel did
not record what an operator could act on.

Both sites now log the file and the cause. A rejected save candidate records its
path with the inner exception type and message; exhausting both candidates records
the path it gave up on; and a rejected settings document records its path through
a small `Reject` helper so every rejection in that reader reports itself the same
way. Strictness is unchanged: a malformed settings document still refuses to mean
defaults, and save files are still left untouched.

Only the diagnostics channel gained paths, which is where private paths already
belong, so the launcher continues to print exception text by type alone.

Verified by launching against a copy of Bill's stale on-disk state, which is how
the real culprit became visible at all:

```
persistence: <data dir>/agents/ensemble debugger/save.json rejected
    (FormatException: Malformed save document); trying backup
persistence: no valid save candidate under <path>/save.json; files left untouched
```

The settings rejection lands in the owning Agent's own `debug.log`, which is where
`Memory` already logs, rather than in the application `ensemble.log` that the save
path uses: `Agent settings rejected at <data dir>/agents/root/settings.json:
Invalid Agent settings document`. Verified with a deliberately malformed document.
Before this, that same failure produced no log line anywhere.

Each of those lines turns an hour of guessing into one command.
