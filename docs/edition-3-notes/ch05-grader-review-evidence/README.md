# Frozen Chapter 5 adapter review evidence

Verdict and limits: ../ch05-grader-review.md. This accepts only the pinned adapter, not student Chapter 5 code.

43 newly completed commands each retain .meta.json, .stdout and .stderr, including empty streams. 31 commands form the compiled sentinel/no-launch matrix; guard-summary.json records expectations and observations. Eight other commands are positive/negative behavioral or command-test controls. Four initial commands failed because the reviewer guessed nonexistent historical app paths; these are retained, not credited as behavioral sensitivity. No new watchdog expired. A callback returning a handle was not treated as a deadline.

Exact candidate diff, artifact/private-input hashes, input verification and final preservation checks are included. parent-* files are carried parent records, not independent new runs; the two parent build failures remain setup history. Entry audit controls remain separately in ../ch05-entry-evidence and are explicitly identified as carried evidence in the report.

Placeholders: <REVIEW_SCRATCH>, <ENTRY_SCRATCH>, <MAIN_REPO>, <STUDENT_REPO>, <USER_HOME>, <SYSTEM_TMP>. All new routes used the same frozen compiled candidate. Private source, historical implementation, cases/maps, sentinels, launchers and binaries are not distributed here; only their hashes, invocations and observed results are retained. No repair fixtures are exported.

SHA256SUMS covers every file in this directory except itself. The report's hash is delivered separately to avoid a self-referential hash.
