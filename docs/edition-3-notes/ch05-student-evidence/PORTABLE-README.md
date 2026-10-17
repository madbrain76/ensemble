# Chapter 5 finite-candidate evidence

This directory contains exact sanitized UTF-8 text copies of receipts, tested capture scripts, command/exit logs, failure records, source/binary identities and reports. It does not contain executables, source archives, PNGs, interpreter caches, credentials, private raw histories or provider opaque payloads. See the adjacent `ch05-student-evidence-manifest.json` for hashes of every copied file, exclusions and retained local artifact paths.

Run `python3 verify-portable.py` from any location to verify this directory and the adjacent report. It performs no network call, candidate execution or credential lookup. Do not rerun historical live commands: their one-shot authorizations have been consumed.

## Decisive receipts

- `final-candidate-state.json`, `frozen/identity.json`, build-information text: clean commit/tree, exact source archive and CLI/workflow hashes.
- `resume-r2/gate-exits.txt`, root race/coverage/module logs: final production-byte gates. Source did not change after these tests; frozen binaries were rebuilt only for clean VCS stamps before grading/live work.
- `frozen/grade-ch05-r2.json`, `frozen/grade-ch01-r2.json`, exit/command receipts: both100/100, exit0, same frozen CLI.
- `resume-r2/mutations.json` and individual logs: six valid detections; earlier disk-blocked records are not detections.
- `live-r1/`: original authorized run remains failed overall. Human executable exited0; four Anthropic requests; actual guarded process and retained output survived. `live-r1-human-observed-facts.json` and `guard-route-r4/results.json` reconcile the faulty prose-sentinel assertion without repeating inference. Never label the original launcher PASS.
- `live-remaining-r1/`: separately authorized workflow/image/local-unsupported run, exit0, no human; seven OpenAI and one Gemini requests, all200. Actual named role-tool effects and image orientation checked, not just done markers/files.
- `local-remaining-r4/`: actual frozen-executable local rehearsal. It is the no-human route evidence. The `unittest.mock` scenario-replacement route check in `guard-route-controls-r4.py` is retained but **uncredited**. Other direct guard, saved-hash and actual HTTP cap controls retain their own value.
- Exact credential selector fields are `directClaudeAPIKey`, `directOpenAIAPIKey`, `directGeminiAPIKey` in approved `~/.cr/settings.json`. Earlier wrong-field scripts are preserved as errors, never authorized live; no actual key value is present.

The final tested launcher `launch.py` has SHA25668066167187e3618beba4d9c8a49b59ee5f2ba105286dc3ee23235960c4cbfad. Its `png()` function regenerates the image fixture; actual PNG remains outside Git, hash057d5c091556b8846631c6c1e00c276ee1033a6d5b2c93cc078a38429e62605e. Older scripts/checklists are historical receipts, not instructions to execute or claim current verification.

The report is a chronological record with a final-status header. Original absolute paths remain in commands for historical accuracy; this portable index uses relative evidence paths. Source review/acceptance/tag belongs to the parent; this is a candidate handoff, not acceptance.
