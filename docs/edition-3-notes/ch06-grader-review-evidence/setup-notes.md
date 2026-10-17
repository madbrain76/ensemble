# Setup attempts (uncredited)

1. Initial read command invoked `git diff --stat` with an interactive pager, job 1452. It launched `/usr/bin/less` (PID 16984, parent git 16983, owned shell 16981). Only those owned setup processes were terminated with `kill -TERM 16984 16983 16981`; no test/grader was running. All later Git inspection uses `--no-pager` or piped subprocess output.
2. Initial Python `tempfile.mkdtemp` selected macOS TMPDIR, outside authorized `/tmp`. That owned directory was empty and removed by `Path.rmdir()`. It was replaced using explicit `dir='/tmp'` before any builds or tests.

Neither setup attempt is credited as runtime evidence. Root's earlier guard failure and test-only refinement remain separately preserved; no grading occurred under those failed setup attempts.

3. `rg` was unavailable during an informational source search (`/bin/sh: rg: command not found`); the piped command returned 0 because `head` was last. Search was recovered with `grep`. No test/grade was involved.
4. First independent runtime-gate harness stopped on its own overly specific diagnostic assertion. The candidate returned exit 2, `grader: unexpected end of JSON input`, and did not launch the marker. Initial script/receipts/61 records including sentinel selfcheck are retained as `initial-gates.py`, `initial-runtime-gates.json`, and `gates.*`; not credited as a completed matrix. Corrected assertion matches the actual existing CLI error contract. No production fix occurred.
5. Before any grading, source inspection caught that `measure.py` required JSON but omitted `-json`. The unrun draft is preserved as `initial-measure.py`; explicit `-json` was added to all measured routes. No grading occurred under this draft.

6. A later read-only `git show HEAD:cmd/grade-edition3/main.go` also opened the configured pager (job 1482). Owned pager/git/shell PIDs 17678/17675/17669 were terminated; the independently running sequential grader job was not touched. Base source/test/model reads were recovered with `git --no-pager show` and `cat`.

7. An observation call used an incorrect job handle (`wait_for_job(1479)`) and returned `unknown job 1479`. This did not cancel or rerun any process. Completion is recovered from the recorder's saved terminal receipts and owned-process inspection.
