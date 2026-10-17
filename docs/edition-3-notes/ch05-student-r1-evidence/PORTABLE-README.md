# Chapter 5 correction round 1 portable evidence

This is new evidence for clean commit `6c4561cbd6cbae080d952da34a5fcd588a76920f`, not a replacement for original candidate `00489b05` artifacts. No paid calls occurred in this round. See adjacent `ch05-student-review-r1.md` for findings, fixes, commands, limitations, and all preserved mistakes.

Verify after relocation by keeping this directory, `ch05-student-r1-evidence-manifest.json`, `ch05-student-r1-source-manifest.json`, and the report as siblings, then run:

```
python3 ch05-student-r1-evidence/verify-portable.py
```

The verifier hashes text and report/source-index files only; it never executes the candidate, loads credentials, invokes a grader, or makes network calls.

Decisive receipts:

- `frozen/identity.json`, `frozen/cli-build-info.txt`, `frozen/workflow-build-info.txt`: exact clean source/archive/binary identity. The two build-info files are also inside `frozen/`.
- `correction.patch`: exactly the two production files and public regression test added in this round.
- `sensitivity/source-proof.json`, `sensitivity/commands.json`, `sensitivity/results.json`: original compile success; behavioral R1/R2/R3 detections with original production restored only in scratch. All eight R3 layouts, including the favorable original ordering pass, remain in `sensitivity/original-R3.txt`.
- `sensitivity/corrected-focused.txt`: corrected focused race pass.
- `gates/commands.json`, `gates/coverage.out`, `gates/coverage-summary.txt`: complete frozen inherited/race/vet/module commands and exits; 86.2% statement coverage. Coverage files are within `gates/`.
- `local-interface-r2/summary.json`: actual corrected CLI/workflow held-HTTP hint checks with real workflow tools/files and one dispatched batch per role. These are loopback fake-provider receipts, not live inference.
- `frozen/grade-commands-and-exits.json`, `frozen/grade-ch05.json`, `frozen/grade-ch01.json`: both exact-binary public routes passed 100/100, exit0. Both grade JSON files are within `frozen/`. Chapter1 reports actual 62s header and 62s body delays.
- `coordinator-exits.json`: top-level exit ledger, including intermediate fixture/capture failures. Those failures are not sensitivity successes.
- `focused-first.txt`, `local-interface/`, `local-interface-attempt1.py`: preserved report-format and acknowledgement-whitespace mistakes.
- `original-preservation.json`: all 341 original portable texts, 341 original local texts, 14 excluded local artifacts and 46 original source snapshot files still match their original hashes.

Binaries and source tar archives are excluded from the Git text directory. Their exact hashes, sizes, and outside-Git local paths are in the index. The source archive contains all 47 tracked source files; the source manifest also provides per-file hashes. Original paid evidence remains attributed only to the original candidate and was not rerun or relabeled.

No private grader/audit/probe source or historical solution was used. The parent's public findings and approved bounded plan guided these independently written public-boundary regressions. Whole-candidate re-review and acceptance/tagging remain with the parent.
