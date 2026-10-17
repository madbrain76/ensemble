# Preserved failures and evidence-credit rules

All observed test/build/grader commands have complete adjacent metadata/stdout/stderr receipts. A failed setup/read/guard is not a behavioral negative, and an expected behavioral negative is not silently converted to a pass.

## Student history (retained, not reviewer work)

The complete R1 student report/ledger preserves: the nonexistent `Context.Turns` compile error; omitted required positive MaxTokens; and a refused wrong Gemini edit anchor. The shell after that refusal still tested old implementation bytes, so the misleading `corrected-boundary.*` files are old-code failures. The student separately identifies the actual applied correction and its passing gates. The reviewer verified the new/old manifests and independently measured the original/corrected/restored production behavior. No student code was changed here.

All 13 actual live exchanges carried by the R1 manifest belong to the original candidate. The original Anthropic final-answer compliance failure remains failed; later successful original-candidate observations do not turn it into a pass or supply new R1 live proof. No new paid calls or key/settings reads occurred.

## Reviewer setup/bookkeeping failures — uncredited

- Informational inventory/report filename guesses were wrong and produced read errors. Exact task/receipt paths were used for verification and tests.
- Copying the previous recorder preserved its read-only mode. The environment edit was refused; only the new owned copy's permissions were changed. No gate ran under that failed setup.
- `input-verification` stopped on incorrectly transcribed report/manifest hash literals. The actual task, receipt and files agreed. The retained initial verifier shows the mistake; the qualified verifier derives the values from the independently hash-verified candidate receipt.
- `input-verification-r2` stopped on the historical manifest's size-key difference. The qualified verifier supports the actual recorded `bytes`/`size_bytes`/`size` schemas and still requires an integer exact size plus SHA256. `input-verification-r3` succeeded before testing/extraction credit.
- An informational progress display used `exit_code`/`elapsed_seconds`; the recorder uses `exit`/`duration_seconds`. This did not affect commands or their saved metadata.
- The first `actual-cli-gated-and-stricter-assertion` stopped before target launch on another incorrect copied core hash. `cli-pin-recovery.json` proves the authoritative source-manifest pin and actual frozen binary agreed. The failed script/diff are retained. The qualified rerun changes only the exact core pin and a fresh owned output path; no public or stricter assertion was weakened.
- Therefore the initial `integration-sequence` remains exit 1 and lacks a completed aggregate success. Its already-completed shape result is valid. The later qualified CLI receipt is separate. `integration-results.json` was reconstructed after recovery from those completed receipts and explicitly says no command was rerun.
- The first `final-integrity` stopped on an assumed internal receiver/type name in a reviewer comparison. The finalizer also used the wrong receipt-path schema. The retained initial finalizer is preserved; `final-integrity-r2` completed all checks before the GUI interruption.
- The report generator initially assumed a nonexistent `times_ms` field and stopped before writing the report. `write-report-initial.py.txt` preserves that mistake. The completed report uses the actual `times` keys from `gated-cli.json` and does not claim an unmeasured CLI persistence assertion.

## Intentional behavioral negatives — credited only for their stated purpose

- `independent-original-boundaries`: exit 1 after successful compilation of the pinned original `stream.go` in the scratch control tree. The independent test source was unchanged. This is the causal negative for the actual Gemini defect, not a setup failure or a new student submission.
- Corrected and byte-restored controls both exit 0; their compiled test executables are byte-identical. All 99 corrected source-file hashes match after restoration.
- `actual-cli-gated-and-stricter-assertion-r2`: candidate exit 0 and every public CLI check true, but `hint_in_answer=false` and checker exit 1. This intentionally scripted non-echoing answer does not satisfy the stronger launcher assertion. It is not credited as final-answer/model compliance. It still demonstrates actual early flush, hint admission, next-request-only delivery, local tool effect and trace sanitization.

## GUI restart recovery

The saved interruption occurred after all behavioral commands and successful `final-integrity-r2`, before the R1 report/checksum index existed. The same reviewer/task/context was resumed. All 28 pre-interruption receipt-output hashes were verified and no relevant process was running. No behavioral gate or student work was replayed. Only metadata reconciliation, instruction refresh, hash-only `recovery-final-integrity`, and report/index creation followed. The parent's private recovery backup was not read or imported. The disk-space increase was of unknown cause; the reviewer performed no cleanup.
