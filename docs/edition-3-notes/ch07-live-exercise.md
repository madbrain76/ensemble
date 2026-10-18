# Chapter 7 live GUI exercise — root acceptance and launcher false negative

Bill explicitly approved one bounded Claude Sonnet 4.6 GUI exercise: maximum four upstream requests, the frozen GUI's existing4096 output tokens/request, isolated muted Chrome, no automatic paid retry. Actual source remained463996b765631d31923a5726036a34942a67cefd. The pinned launcher-root-r1.mjs and unchanged policy were verified before any credential access/process launch; root read the key only through the approved process-memory path. No install/restart/shared browser/microphone action occurred.

## Actual run

`live-approved-01` ran2026-10-18T04:34:33Z,8.627seconds,launcher exit1. Both requests used `claude-sonnet-4-6`/stream=true/max_tokens4096 and returnedHTTP200. First response requested read_file; second contained the actual result and completed the answer. Recorded usage totals3553input/179output tokens. No third request, fallback, retry or discovery. Both owned GUI/Chrome processes exited0; credentials were not exported to children or evidence.

Actual page observations show Connected/in_flight and a partial Assistant artifact before completion; final Connected/idle DOM shows visible explanation, read_file card and exact file contents, and final answer. The result correctly explains projectLantern Orchard, three reviewers, blue folder. Those facts are not encoded as an arbitrary exact answer phrase requirement.

## Failed launcher assertion retained; independently reconciled evidence

The launcher reported passed=false solely because actualReadResult stayedfalse. Its event predicate wrongly required `tool_finished.name === 'read_file'`. The actual production event has no name field; it correlates by stable call_id and part_id. Saved gui.log establishes exactly one read_file dispatch for exercise-note.txt, exactly one nonerror completion with the SAME call/display IDs, and byte-identical result/fixture/job output. The final DOM displays that result. The second authenticated provider request receipt contains that actual tool result, and the final answer correctly uses all three facts. Root saved the exact correlated safe events and full semantic reconciliation, without rerunning or editing student code. The original exit1 and passed=false remain FAILED launcher outcomes; this is not a retroactive green report.

Root accepts the existing real-model/live-GUI exercise requirement on that combined observed evidence. This is not universal provider/GUI proof or a fresh native speech test. Prior Bill-heard audio remains qualified earlier evidence; screenshots were captured but root does not claim an unperformed visual-image review. The plan's prose called the fixture chapter7-note.txt; actual pinned policy, prompt, launch and result consistently use exercise-note.txt. That prose typo is a plan-record issue, not an application mismatch.

## Evidence and portability

`ch07-live-exercise-evidence/` holds actual invocation/exits, original stdout/stderr, sanitized receipt, extracted correlated tool events, root reconciliation and local-artifact hash metadata. Its text index excludes itself and the two PNG screenshot copies; screenshots and raw runtime/API captures are local-only pinned artifacts, not automatically chapter-commit contents. No private settings, credentials, raw authenticated headers or browser profile are publication material.

The failed original launcher predicate was not changed after the paid run, and no paid permission remains to retry automatically. Future implementation changes would need their own qualified evidence rather than inheriting this run as if it used different production bytes. Code acceptance, tag and later author/text closure remain separate recorded decisions.
