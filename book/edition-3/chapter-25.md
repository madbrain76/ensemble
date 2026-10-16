# Chapter 25: Crossing Over to Ensemble

The second-edition agent became the place where Bill and Cadence worked on the
book. Daily use exposed problems that the construction tests had missed. The
[crossover record](../../docs/edition-2-notes/ch25-crossover.md) preserves the
changes and their evidence; they are not fresh Edition 3 results.

A speech-speed slider stopped at 2×, and the native voice rate was nonlinear.
The correction changed both the control and saved-setting validation, with an
approximate mapping above 2×. Building the fix was not installing it: the
launcher still used the old embedded-assets executable. Later, a server restart
left an open browser running old JavaScript. Binary installation, server
restart and browser reload are distinct steps, each requiring observation.

A 60-second HTTP client deadline also killed legitimate model responses. The
fix removed that total deadline, not explicit cancellation or connection
protections. A short job callback is a progress interval, not permission to
terminate the work. Chapter 1 now makes that distinction before jobs exist.

Automatic browser snapshots and speech-queue observations generated noisy
per-round activity. Making those two tools on-demand solved the actual problem;
removing general ephemeral context would have discarded a useful mechanism.
Chapters 11–12 now teach explicit observation and confirmation.

Bill also needed visible explanations before actions. Private or encrypted
reasoning was not a substitute. The implemented model capability requires
ordinary non-whitespace assistant text before the first call in the same batch;
rejected batches produce corrective results without dispatching their tools.
This checks presence and order, not whether an explanation is good. Honest
rationale remains the agent's responsibility.

Application settings moved into one explicit, versioned configuration file,
with credentials referenced separately. File mode isolates provider settings
from inherited environment; legacy startup remains available for tests.
Provider changes require restart, while ordinary GUI and memory settings apply
immediately. Persist a valid update before accepting it in memory, and keep
secrets out of flag defaults and help output. Chapter 21 now owns this
application integration and its GUI startup path. Chapters 22 and 23 extend
the same strict schema and immediate-application contract when memory and
recall exist. Chapter 16 owns provider capability and reasoning choices;
Chapter 20 owns inherited running configuration and role-specific instructions.
The explanation-before-tools checks belong in Chapters 3 and 6, where dispatch
and actual stream order first become observable. None is a new exercise here.

Other repairs are already in the preceding chapters: root conversation does
not require worker submission; completed handoffs retire old opaque replay
from active context; compaction runs as a restricted child; and cleanup waits
for the child's final save. Preserving archive facts and reducing the working
context are different operations.

The artifact-identity investigation found a separate restart failure: a new
response reused an old display ID and entered its card. The messages remained
in the archive. Chapters 7 and 10 clarify the existing identity obligation;
that clarification is not a claim that the copied baseline or legacy saved
collisions have been repaired.

This chapter remains a historical crossover account, with no independent
construction exercise, grader receipt or code checkpoint. Its demonstrated
repairs belong in the earlier construction chapters. Chapter 26 continues
from the exact accepted Chapter 24 code, not a hidden Chapter 25 repair pass.
Edition 3 must earn its own execution results.
