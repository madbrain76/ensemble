# Chapter 5 restart — implementation notes

The Agent now owns one mailbox-driven conversation and all runtime History writes. Engine remains a one-request HTTP child. Ask/Submit retain separate, reliable private replies; live admission is classified by the actor, not by an adapter’s possibly stale view of its buffered reply. Observer delivery is bounded and lossy, so it is not a completion protocol.

A hint releases an observation wait, not the job or HTTP request. Dispatched work keeps its real handle, cursor, output, and outcome. Undispatched batch tails get exact-ID interruption results without effects or replay. Frozen HTTP requests stay frozen; admitted hints survive text-only responses and are consumed by a later request.

Two lifecycle distinctions matter. Explicit Go-only kill contains output and reports that Go cannot stop the handler; it does not wait for that handler. Shutdown does wait for real work. Shutdown also joins an admitted killed-fact acknowledgment independently of process/handler settlement, so final save cannot omit an external Kill’s terminal fact.

The four supervision tools never own a `tool.job`. Their call arguments and truthful reports identify the target; ordinary-job metadata and top-level terminal facts remain intact.

The separate historical `ch06` module uses three ordinary Agents and application-owned synchronized manuscript state. Author, editor, and reviewer must perform their actual role-tool effects. Reliable private completion advances the pipeline; text-only replies or observation labels cannot manufacture a completed stage.

Media support is explicit per model/type. Unknown models remain usable for text, but have no assumed media capabilities. Supported local media supplies actual bytes; unsupported media fails before HTTP and names the model/type.

## Local validation

Both public grades passed with process exit 0 and parsed `passed:true`, using the same frozen CLI: Chapter 5 (all seven reported checks) and Chapter 1 (including real 62-second header and body delays).

The full root race/vet and module checks passed before the final one-field supervision-schema correction. After that correction, focused race tests for serialized supervision, explicit Go kill, and final-save terminal delivery passed. GUI and consumer have no test files: their test results are compile-only. The workflow’s real role-effect/handoff integration test passed.

Local success is not a claim of final acceptance. The subsequent, separately authorized live execution is recorded below.

## Single authorized live execution

The actual CLI chat path passed with four Gemini POSTs. A real shell job printed its start marker and kept running while a hint was admitted, `/interrupt` ended the turn, and a new turn correctly answered that the remembered marker was **COMET**. That answer preceded the old job’s terminal report by about 11 seconds. The same job retained its late output, finished with exit code 0, and was reaped; EOF then closed the CLI cleanly. These were root-supplied inputs through the real chat adapter. The inspected runner deliberately held delivery of the genuine hint-bearing HTTP response until interrupt acceptance: this demonstrates a controlled schedule, not natural provider latency.

The run stopped at the workflow author’s first response. Gemini returned HTTP 200 with `MAX_TOKENS`, no content parts or tool calls, and usage containing prompt/total/thought counts but no `candidatesTokenCount`. The adapter reported `response missing usage`, and the workflow exited 1. No author tool effect or editor/reviewer handoff occurred. That was not a successful live workflow. These facts do not establish why the content was empty, and omitted candidate usage is unknown, not zero. No production repair was justified by this response; nothing was patched or retried.

Five provider POSTs were made in total, all to the approved Gemini model with native output limit 1024: four CLI and one workflow. The image request and local unsupported-video probe were not reached. The runner exited 1 and cleaned up its owned work; partial data remains local. Production, the frozen binaries, and the inspected runner were unchanged. At that stop, live workflow and media validation remained incomplete.

## Authorized remaining exercises

A separately approved B→C-only run exited 0; A was not rerun. The unchanged workflow used `claude-sonnet-5-5` for nine POSTs. It stored a 100-word ocean-at-dawn draft, measured that text through `word_count`, read and edited the stored draft, and obtained a reasoned acceptance of the actual edited text. Reliable completion advanced author, editor, then reviewer. The first Claude response had `max_tokens` and empty `write_draft` arguments; the tool truthfully returned `text is required`. The model subsequently issued a distinct successful call in the same conversation. This ordinary tool-error recovery was retained, not concealed or implemented as a runner/HTTP retry.

One `gemini-3.8-flash` POST then received the actual local PNG bytes and correctly described a blue square on the left and a red circle on the right. The local `gpt-5.4` video probe named the unsupported model/type and made zero HTTP calls. Workflow, image driver, and video probe exited cleanly, and the runner closed its owned work.

The two authorized executions used 15 provider POSTs in total (five earlier plus ten remaining), within the combined ceiling of 18; every request retained native output limit 1024. Only the temporary remaining runner's provider forwarding/environment was adapted under explicit authorization. Production and frozen binaries were unchanged. Earlier failures and private raw outputs remain local; no automatic retry or further campaign occurred. Final acceptance remains root-owned.
