# Edition 2 Chapter 19 technical validation

Status: historical evidence reconciled; exact author examples and independent
combined review pending. Coder: `edition2_grader_coder`, October 11, 2026.

This is original Chapter 20, **The Crossover**, a prose-only chapter. There is
no original or revised grader and no numeric score. No runtime repair list,
new implementation, command dispatch or frozen export is part of this pass.
The mandatory `.agents/skills/ensemble-coding/SKILL.md`, full carryover and
workflow, current enhancement and map were reloaded after compaction. The full
original chapter and complete student/comparative records were read.

## Reused evidence, checked against the actual artifacts

Original source is `2161d2e9f26b6acecfb1457e7053065063944fb4`, chapter SHA-256
`9a6d41326806d229d1dc7913297d7b5f2ef04ab861e1068f67d393625e9700ac`.
The accepted direct `solutions/edition-2/ch19` tree remains
`bd2998e68f976dbaf437b642b131d350e086242f` (144 files). Its original independent
acceptance is `04e450da2f38f880e9db9f2300e49ad19c831872`, zero correction rounds.
The paid experiment used source `fab031529cda55a76d530478f09de0cd08e7db36`,
tree `59010a4c2b333c2f74777b6d76079e8f2b3af90d`; only README changed afterward.

[The read-only reconciliation](../../solutions/edition-2/evidence/ch20/enhancement19/historical-audit.json)
was produced by `python3 solutions/edition-2/evidence/ch20/enhancement19/audit.py`.
It verifies retained hashes and decodes the actual wire, event log and save:

- Five requests and five completed responses select `gpt-6-astra` at high
  reasoning, detailed summary, streaming and `store:false`.
- Two actual reads are followed by failing `run_command`, exact `edit_file`, and
  successful `run_command`. Correlated results contain three failures and three
  errors before the edit, then six passing tests. Applying the recorded unique
  replacement to the original file produces the retained final file exactly.
- The six-test file is byte-identical before and after, SHA-256
  `d0a5eaab886c93918f853cde4e19819f0ab2fc8b03bf8aae504f5c8152081ea9`.
  The existing independent reviewer reran both versions; this pass reuses those
  results rather than claiming another test or model run.
- The saved log matches the complete recorded event log. The original launcher
  records normal exit 0, no timeout, 27.46 seconds and stdin closed after the
  final answer. It did not wait for optional public summaries.
- The trace contains 161 ordinary text deltas and zero public summary deltas.
  All five reported cache-read/write grants and reasoning-token counts are zero.
  The latter is not a measurement of private reasoning. Plan requests omit the
  previously accepted paired explicit-cache controls.

The historical endpoint derives from explicit configuration and transport source;
the retained trace does not independently capture URL or authorization headers.
The signed-in profile was retained, not newly created or logged out by that run.
One bounded CLI repair proves neither interactive steering, GUI use, arbitrary
project reliability, cache savings nor Edition 3 authorship. Chapter 25's actual
ongoing crossover remains separate.

## Existing implementation and sensitivity

The accepted chapter change moves the already selected Astra row onto the
existing Responses renderer. It adds two rationale comments and six active test
lines, with no new runtime abstraction or configuration path. The source test
renders actual requests and checks both transport and configured high effort.
Its existing passing baseline and two intended failures are retained in
`evidence/ch20/baseline.txt`, `astra-chat-route.txt`,
`drop-reasoning-effort.txt` and `mutations.json`: restoring the Chat route fails
`model transport`; replacing high with low fails `configured model/reasoning lost`.
No fresh mutation run or broader inference is claimed. The independent original
model-boundary race result remains 1.027 seconds on its recorded source.

No legacy `solutions/ch20` snapshot exists, so no invented reference delta or
legacy grade is supplied. The original historical repair names and code excerpts
are narrative evidence, not current public API or new Chapter 19 requirements.

## Exact examples

Pending the author's stable sample freeze. The planned reproduction uses the
actual accepted core command, a fresh scratch copy of the retained task inputs
and a typed primary exposing the existing file/command tools. This is explicitly
a new demonstration setup, different from the historical default-tools launch.
Any new execution will use local external issuer/provider peers and will be
reported separately from the actual paid model's coding result.
