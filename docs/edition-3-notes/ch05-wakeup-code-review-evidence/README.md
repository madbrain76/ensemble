# Independent whole-candidate amendment review evidence

Portable UTF-8 records for `ch05-wakeup-code-review.md`. `command-index.json`
classifies all 48 child command results: passing regression gates, valid behavioral
controls, source/stamp provenance, and retained uncredited setup failures. Each
command retains exact logical argv/cwd, actual exit, elapsed time, watchdog state,
and complete stdout/stderr. No watchdog expired.

Path substitutions: `<MAIN>`, `<STUDENT>`, `<REVIEW_TMP>`,
`<STUDENT_EVIDENCE_TMP>`, `<ACCEPTED_R1_TMP>`, `<AGENT_RECORDS>`, `<USER>` and
`<SYSTEM_TMP_PATH>` remove host-specific roots. They are not download links.
Scratch source/binaries/private probe code remain outside Git; their hashes and
candidate identities are recorded. No repair code is supplied.

`identity.json` records verified inputs before review; `final-verification.json`
checks preservation afterward. `prior-preservation.json` keeps accepted prior
work separate. `retained-local-facts.json` inspects already-recorded local
executable proof, not a new provider or executable campaign. Setup correction
receipts explain the uncredited fixture failures. Source resolution and
successful no-test compilation alone are not behavioral sensitivity.

`MANIFEST.json` hashes every other evidence file except itself and SHA256SUMS;
SHA256SUMS additionally includes MANIFEST.json. No binaries or opaque live
provider history are exported. Existing student/reviewer reports are unchanged.
