# Setup / failure ledger — no success-only accounting

All actual build/test/grader invocations are preserved in `commands.json` and their full stdout/stderr files. A nonzero command is not silently recast as a pass. No candidate implementation was repaired by the reviewer.

## Required candidate finding

`independent-public-final-r2`: exit 1, four Gemini signed/opaque boundary failures; all seven other top-level tests pass under race instrumentation. This is the qualified reproduction. `independent-public-vet-r2`: exit 0. See report R1 and `public_checks_test.go`.

## Expected deliberate-fault failures

`ignore-disable-behavior` and `wrong-final-id-behavior`: each exit 1 after successful compilation and a passing baseline. Both byte-restored controls compile and pass. `mutation-results.json` has source/binary hashes. These are scratch sensitivity defects, not defects alleged in the candidate.

## Stricter answer-compliance assertion

The original student live failure is retained in the immutable student report/command ledger and safe observations, verified through the 474-entry manifest. No live call was repeated here. `actual-cli-gated-hint-and-stricter-assertion-r2` exits 1 on `hint_in_answer`; all the public-behavior checks and candidate process exit pass. This is an explicitly scripted local answer, not an invented claim about model behavior. The separate `public-shape-reproduction` stream/plain exercise passes but does not replace the live failure.

## Reviewer setup / fixture errors (uncredited)

1. `input-verification`: exit 1 from assuming `candidate.source` instead of the actual receipt schema. Original verifier is `verify-inputs-initial.py`. No extraction or test ran before correction; `input-verification-r2` succeeds.
2. `independent-public-parser-checks`: Go requested dependency metadata updates before compilation. `independent-module-tidy` modifies only the owned review module and succeeds.
3. `independent-public-parser-checks-r2`: preserved real Gemini failures plus invalid reviewer assumptions: tool finals expected arrival order instead of provider-index order; Anthropic ID expected compact 1 instead of legal sparse 17; a recovery fixture lacked required usage. Corrections change only fixture/assertion expectations, never candidate code. Initial source snapshots are retained.
4. `independent-public-parser-checks-r3`: compilation fails for the reviewer's nonexistent `SendPreparedWithCallbacks` method. Corrected to the actual public actor `AskWatching` API. `public_checks_nonexistent_method_test.go.txt` retains that source.
5. `independent-public-parser-checks-r4`: Gemini failures plus missing usage on the parser-reset fixture and missing required System in the unknown-model fixture. `public_checks_missing_fixture_fields_test.go.txt` retains that state. Those setup facts are corrected in r5/final; no behavior is credited from an invalid fixture.
6. `streaming-ch05-wakeup-and-interruption`: initial wrapper incorrectly reports `claude-sonnet-4-6` instead of the original fixture's `test-model`. Correct same-model replay protection rejects the opaque payload in the mixed test. The exact-model negative is retained. The fixture wrapper is corrected to preserve `response["model"]`; no production byte or original wake-up assertion changes. R2 passes all ten selected tests. The reviewer's preliminary JSON-order explanation was withdrawn after inspecting the actual assertion and model path.
7. `actual-cli-gated-hint-and-stricter-assertion`: the first observation checker expects a nonexistent `response.` trace label. The actual logger writes event-kind plus sanitized JSON, with a blank kind for unlabeled SSE frames. Initial script/JSON/trace/stdout/stderr are retained. R2 checks the real message-start frame; the stricter final-answer assertion remains false and preserved.
8. `independent-public-final` and `independent-public-vet`: added batch-interrupt test misreads `Jobs.Get`'s second return as bool instead of error; both fail at compilation. `batch_interrupt_initial_test.go.txt` retained. Correct API use yields final-r2's qualified result and vet-r2 exit 0.
9. `final-integrity-check`: all 98 pinned source hashes pass, but a mistaken no-extra-files assertion fails because root tests create owned local `api.log`, `cmd/ch01/api.log`, and `cr/io/*` byproducts. They are retained, not deleted. `*-runtime-extras.json` records all hashes; the 98 source identities are verified separately. Final-integrity-r2 succeeds, including complete restored scratch trees.

## Informational read errors (no grading credit)

- Initial `wc` used `docs/voice.md` / `docs/chapter-writing-procedure.md`; the correct full `book/` documents were then read.
- Early archive/notes/test/helper filename guesses returned file-not-found; the complete assignment, manifests, and actual inventory supplied the authoritative paths before execution.
- An informational final `ls` included an unmatched `*pin*` glob; the prerequisite hash file and complete manifests exist and were verified. This was not a failed gate.
- Combined prerequisite output omitted a voice-guide interval; it was recovered explicitly. A large runtime-extra metadata listing was bounded inline; the complete four JSON inventories and full owned-scratch manifest retain every record without exposing private payload contents.

## Historical failures remain historical

The assignment/student report preserve the initial student 70-point result, actual CLI buffered-input fix, and authorized Anthropic answer-compliance failure. Prior root adapter setup guard and first test-isolation weakness belong to the unchanged accepted adapter evidence, not fresh failures or passes in this review. No formal student correction round was used here.
