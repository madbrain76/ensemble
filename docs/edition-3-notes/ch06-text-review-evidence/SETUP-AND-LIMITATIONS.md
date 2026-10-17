# Setup and limitations

The initial inventory command unsuccessfully tried to cat six guessed author metadata filenames: INDEX.md, commands.json, initial-state.json, authoring-checks.json, final-validation.json, SETUP-AND-FAILURES.md. Each returned “No such file or directory”. All were read-only probes; the enclosing shell returned 0 because its final worktree-list command succeeded. These are not author failures or chapter-check failures. The actual inventory was then used.

Earlier large read_file displays of AGENTS.md and ensemble-coding/SKILL.md were capped; both were recovered in contiguous explicit line ranges. Historical workflow was read through line 2434; ranges beyond EOF returned empty output, not missing content.

No code tests, provider calls, key reads, publication builds, live restart, staging, commits, cleanup, or changes outside the allowed reviewer report/evidence were performed. Old provider proofs remain old-commit evidence.

Some long explanatory grep displays were capped after the reports themselves had been fully read. The exact grep output was saved as report-qualification-excerpts.txt (116 lines) and recovered in complete 1–75 and 76–116 ranges. Chapter excerpts and successful validator output were also recovered. This is output-size recovery, not a behavioral rerun.

The independent text validator completed on its first attempt, exit 0. The six guessed-filename cat probes were the only unsuccessful reviewer setup probes; there were no failed lint runs or failed chapter validation assertions. The author’s retained failed setup/lint-path probe, two refused prewrite anchors, and off-by-one diff validator are separate historical author receipts; no author formal revision had been used when this review began.
