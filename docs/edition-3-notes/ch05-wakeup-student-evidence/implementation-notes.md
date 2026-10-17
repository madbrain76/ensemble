# Approved public wakeup amendment: implementation/assistance record

This file records this student's implementation, not a parent/reviewer source.
Full current mandatory public prerequisites and both pinned chapters were read
before source edits. The parent approved the bounded observation/batch ownership
plan as main-18. This is new public work, not a second formal correction round.
The closed R1–R3 tests remain unchanged in source and run in focused/full gates.

No new architecture: keep Agent's one actor, existing dispatcher, retained jobs,
WaitOptions and Job.Wait, application-owned role routing, and the existing history
reducer. A per-observation child context with a distinct shared cause and a scalar
observation identity separates waiting from execution. Real incoming UserMessage
metadata and Parts are recorded first. The worker exits through the locked job
report/cursor path, stopping its timer, without canceling the retained job or
active turn. Actor then pairs the undispatched tail with exact-ID not-executed
errors and consumes literal-next limits without creating work or arming a skipped
setter. All original calls and matching signatures stay intact. New model calls
can execute while the old job continues. No new special join/set_limits API.

Public-boundary tests compile on exact original6c: new scalar fields/options are
populated through JSON so old production must fail behaviorally, not at compile.
Real held registered handlers and real PTY commands remain live through wakes;
actual next HTTP, same-ID results, real outputs, cap/cursor/terminal behavior,
side effects, limits, dispatch races, short-timer races and saved histories are
asserted. Observer markers synchronize admission only; they are never the proof
that a question or result reached the model.

The first frozen new candidate ccf82d8 passed focused tests but failed the full
inherited PTY input test. Broad immediate-byte delivery had allowed later arrivals
and send_input echo to bypass callback expectations. This real intermediate
compatibility failure is preserved. A concrete question to the parent cited the
public wait_for_job wording and the inherited send_input behavior. Parent main-19
approved a small observation-data option, enabled at wait_for_job only, to inspect
initially available unseen bytes. Later arrivals still recheck pattern/completion/
message/delay rather than forcing return. Ordinary initial/input/kill callbacks
remain unchanged and all remain hint-wakeable. The parent did not author source.

Final source f60ebf3 is clean. Seven original6c controls compile and fail for the
intended public behavior; restoring the complete corrected production diff with
the identical tests passes. Ten corrected focused race repetitions include R1–R3
and inherited PTY compatibility. Full final-byte gate and same-binary grade
receipts are separate from these controls; the historical grader does not cover
the newly amended wakeup contract.

Only new owned main report/evidence/index paths are written; old reports, indexes,
local receipts, source snapshot, artifacts and live evidence remain unchanged.
Original12 paid calls remain00489-only. No paid call, credential lookup, child
agent, GUI, install, restart, push, acceptance tag, or main commit is involved.
