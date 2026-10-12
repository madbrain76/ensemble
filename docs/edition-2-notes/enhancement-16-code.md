# Edition 2 Chapter 16 grader: coherent stage, review pending

Coder: GPT-6 Astra/high. Mandatory `.agents/skills/ensemble-coding/SKILL.md`,
full carryover/workflow and current enhancement/map reread after resumption.
This stage concerns original Chapter 17, Auto-Recall. No frozen runtime,
original grader, active application or historical model identity was changed.

The separate command retains all ten original weights (100 total, no parity).
It launches the actual optional GUI, observes completed turns and successful
EOF, uses a typed primary with real tools, and reads recorded recall from the
existing four-field save Log. `-recall-data` selects only the grader corpus root:
`cr` for the accepted application, `.` for the legacy application's cwd layout.
It is not a new application flag. The external Anthropic HTTP peer keeps actual
Agent, indexing, judging, rendering and persistence machinery in the path.
Current fixtures select the already supported Sonnet-5.

## Measured stage

Evidence: `solutions/edition-2/evidence/ch17/enhancement16/`.

- `accepted-first-green.json` and its `-run.json`: actual full 100/100,
  exit 0, 16.563 seconds; command `go run ./cmd/grade-edition2 -ch 16 -json`.
- `legacy-first-full.json` and its `-run.json`: actual full 60/100,
  exit 1, 16.327 seconds; command adds `-gui solutions/ch17 -recall-data .
  solutions/ch17`. This is not a composed score. The quota scenario encountered
  an extra memory-graduation request; this masks its intended recall assertion
  and is not yet classified as a recall defect. Refresh and retained/capped
  attachment failures are separately reported in the raw checks.
- Initial `accepted-initial.json` is actual 10/100: the driver assumed
  `CH02_LOG`, while GUI launch sets its own event-log filename. Nine runtime
  checks stopped at the missing observation file. Reading the existing public
  save Log corrected this fixture; no runtime defect was established.
- `accepted-save-log.json` is actual 95/100. The supposed stop-word-only trap
  used heading “Archived note”, while the query legitimately included “archived”.
  Accepted indexing includes heading text. The neutral “Record” heading removes
  that accidental competing hit; the initial result remains preserved.
- `accepted-focused.txt` passes the fast actual-runtime suite after that fix.

The original shared fake holds its capture mutex while calling Route. A held
judge therefore also blocked main requests and capture inspection. The narrow
external peer captures under a short lock, then waits outside it. Its full
15-second timeout scenario observes main completion while the judge is still
held; release occurs on cleanup. This is a fixture correction, not a runtime
protocol or shared-fake change.

## Remaining work and limits

This is a first-green checkpoint, not chapter acceptance. Exact published
examples, targeted sensitivity, final source binding and independent review
remain pending. Check quota fixture compatibility before interpreting its legacy
failure. Ranking compares competing hits without demanding exact scores; runtime
parent relationships, all-provider transport, interruption and semantic relevance
remain separately reviewed source/live claims. Historical original 0/100 and
three successful live sessions retain their original provenance, including the
historical Anthropic model identity. No fresh paid result is claimed.
