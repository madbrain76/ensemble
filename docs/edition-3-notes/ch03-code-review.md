# Chapter 3 initial code-acceptance review

## Verdict: CHANGES REQUIRED — one blocking implementation finding

This is the **initial** independent code review, not Chapter 3 acceptance and not a correction-round review. Candidate: `0eb6fdffe62e662b0ab96982b01287afba19ff42`; accepted parent: `89cc874b3c32e65ee43f44a72ad9ca48c9dd3323`. The root retains authorization, correction-round accounting, student launch, and author timing. No student contact or coaching occurred during this review.

**R1 — blocking:** Gemini ordinary visible text carrying a `thoughtSignature`, without `thought: true`, is classified as opaque. Consequently a real nonblank final answer can become empty CLI success, and a valid signed explanation can be rejected before tool dispatch. This also removes ordinary text from cross-model conversation projection. Independently reproduced using synthetic, nonsecret fixtures through both the public API and the actual pinned CLI. The real Gemini evidence corroborates the terminal-answer manifestation.

The existing grader composition passed independently at 100/100, but those scores do not establish R1 correctness. No second blocking implementation defect was established in the broader review. A fallback-ID **test coverage** gap was found and closed in reviewer scratch; the candidate itself passes the added check.

## Frozen inputs and method

- Reference HEAD remains `f7530a6643ca00bd2f59a8da60021c567cd7d174`; no tracked reference edits. Prior full governing reads remain valid: every source/governing hash retained in the entry-audit manifest was rechecked and matched. The current reviewer skill and pertinent actual Chapter 2/3 contract passages were reread.
- Student repository was clean at the pinned candidate. Verified all **20 source hashes** and **146 evidence hashes** in `ch03-candidate.json`, the chapter hash, pinned executable hash, and source-archive hash. The archive's 20 file contents match the tracked candidate files. Read the full 268-line student report, SHA-256 `e2501bb153ddb005b2afb361634a7d68868c86a48b2027b492f9a9e627f19487`.
- Inspected the public library/owner chain, Agent loop, transport/parser/renderers, History, CLI and tool implementations, Chapter 3 tests, inherited tests and relevant graders. Ran independent baseline tests, both graders, local actual-CLI fixtures, public-API fixtures and scratch sensitivity controls. No reference/student source changes, installs, commits, runtime restarts, subagents or paid provider requests.
- Both graders used the **same manifest-pinned explicit executable**, whose full hash/path is in this review's evidence manifest. They did not use directory discovery or a newly invented Chapter 3 Edition 3 flag. Directory build fallback is `./cmd/` when that package is `main`, otherwise root `.`; no named-subcommand discovery is implied.
- The supplied executable was built from final source before the commit and carries predecessor-plus-modified VCS metadata. It is not described here as a byte-identical rebuild at the final commit. Go tests compiled archive copies; that is source testing, not executable-byte identity proof.
- Durable receipts are under `ch03-code-review-evidence/`. Test source and local-server driver are saved as nonexecutable `.txt` evidence with synthetic values only. Source mutations were confined to temporary copies.

## Public basis and R1 causality

### Public obligation, not a new feature

Chapter 3's ordinary synchronous loop must return the assistant's final answer after tool results. Its selected-model explanation rule, especially `chapter-03.md` around lines 211–228, requires **ordinary visible nonwhitespace text before the first tool call in that same response**. Hidden thinking or a readable reasoning summary cannot substitute; later text and earlier-conversation prose cannot substitute. Whole-batch refusal must occur before any handler runs, preserve the assistant response, provide a matching refusal result for every call, and permit an explained retry within the existing finite budget.

Chapter 3's actual-human exercise requires a conversation, work with tools, and a follow-up through the real CLI for each supported vendor. Disk/test success alone does not establish that the conversation returned an answer. Inherited Chapter 2 obligations distinguish ordinary portable text from opaque exact-provenance material: vendor, API surface and exact source model constrain signature/reasoning replay; private material must not leak into another model or a human-readable answer.

A signature attached to a public text part is not, by itself, a `thought: true` flag. The preserved live metadata makes that distinction concrete; no speculative provider capability or model-name assumption is needed.

### Source path

1. `internal/engine/parse.go:213–244` reads each Gemini part, including `text`, `thought` and `thoughtSignature`. Tool calls take a separate branch which preserves the call signature.
2. At **lines 240–243**, `block.Thought || len(block.Signature) > 0 || block.Text == nil` classifies any signed non-tool text as `Type: "opaque"`. An ordinary signed text part therefore never becomes neutral visible text.
3. `internal/agent/agent.go`, `Ask`, collects and recognizes only `Type == "text"` as a visible answer/explanation. When there are no calls, it returns the accumulated string with nil error, even though this wrongly classified response contributed an empty string.
4. `cmd/ch01/main.go` treats nil error as successful protocol/chat completion. The protocol writes an empty assistant string; chat prints no actual answer and exits successfully. When signed ordinary prose precedes a call, the same misclassification instead falsely activates the explanation refusal path.
5. The renderer's exact-provenance filtering correctly prevents opaque signed material from crossing models. Because the ordinary text was mislabeled opaque, that protective filter also removes the public text. Weakening that filter would be the wrong correction.

This classification branch is inherited from Chapter 2. Acceptance is cumulative: inheritance explains the origin but does not excuse the demonstrated Chapter 3/live-conversation failure.

### Independent reproduction

`public-api-fixtures.go.txt` and `gemini-regression.txt` record public-API tests with `gemini-3.8-flash`, invented harmless text/signature strings, measured usage, and normal constructor requirements:

- Both an **absent** thought flag and an explicitly **false** flag fail: the visible text is empty; ordinary text disappears when rendering for another Gemini model, Anthropic Messages, or OpenAI Chat Completions.
- Same-model dump/load/render preserves the synthetic signed material. Genuine `thought: true` material stays hidden and does not leak across targets. Those assertions pass; the fix must preserve them.
- A signed ordinary explanation before `write_file` is wrongly refused, the requested marker is absent, and the signed terminal answer is empty. This is the same classifier defect, not a second policy-selection defect.
- Measured usage remains intact in this fixture.

`local-cli-results.json` records the **actual manifest executable** against a local controlled Gemini provider. Both protocol and `chat`, each with two user turns, receive nonblank signed ordinary terminal text and exit **0 without displaying it**. No provider access or private payload is involved. This is independent of the student's tests and of the old grader's coverage.

### Bounded required correction

Correct the ordinary-text/hidden-thought classification and the necessary neutral representation/rendering path so that:

1. Signed ordinary text without `thought: true` remains ordinary visible text, reaches the final API/CLI answer, and can qualify as a same-response explanation before the first call.
2. Its signature remains opaque and bound to its original exact source vendor/surface/model and original part association. Same-model replay preserves it without duplicate visible text or malformed provider parts.
3. Cross-model/surface rendering retains the ordinary public text but omits the signature and genuine hidden reasoning. Do not solve this by dumping opaque JSON/reasoning into answers, stripping every signature, broadening model equality, or adding a model-name exception.
4. Genuine hidden thought remains unable to satisfy the explanation rule. Existing refused-batch no-effect behavior, per-call pairing, shared ten-round budget, normal nonselected behavior, usage and History replay remain intact.
5. Add regression coverage for these signed-ordinary-text cases, final answers through the real CLI, and exact-source replay/isolation. The current signed **tool-call** test is useful but does not cover signed **ordinary text**.

This review supplies a failing behavioral boundary, not a production patch or dictated internal representation. No new public flag/API/configuration framework, empty-answer policy for every conceivable provider response, universal explanation rule, higher token cap, provider switch, stream-order implementation or grader rewrite is requested.

## Whole-chapter and inherited review

### Selected-model activation and refusal protocol

The accepted public seam is present: `Config.RequireToolExplanation` and documented CLI `LLM_REQUIRE_TOOL_EXPLANATION=true|false|unset`. CLI invalid values fail before sending. The Agent uses the configured boolean rather than vendor/model-name heuristics; an ordinary arbitrary model string works in local controls. `Config.EnableTools` leaves ordinary library Agents empty; the CLI enables the six built-ins. That opt-in is separate from the explanation policy.

The Agent scans the complete synchronous response before dispatch. It records whether qualifying text existed before the first call, assigns missing call IDs, preserves the assistant response, and then either executes calls or records matched refusal results without dispatch. The refusal explains that nothing ran and requests explanation before retry. Tool failures remain results, and the next provider request receives the committed facts. All sends, including policy correction, consume the same `maxRounds = 10` loop; there is no extra policy-retry budget.

Independent actual-CLI controls (not just source inference):

| Case | Observed result |
|---|---|
| Policy true; missing, whitespace-only, opaque-only, late, or prior-turn prose | Entire six-call batch refused; six paired error results; at the next request the seed file is unchanged and neither write nor command marker exists; explained retry subsequently executes and completes |
| Proper same-response explanation before calls | Six calls execute, no refusals, final answer returned |
| Policy false or unset, same arbitrary model and missing explanation | Six calls execute; no universal-rule regression |
| Invalid policy value | Nonzero CLI exit, zero provider requests |
| Repeated invalid six-call batches | Exactly ten requests, sixty paired refusals, nonzero exit, no command/write effect |

Missing/blank/opaque/late cases each made three requests; the prior-turn case made four including its previous text-only turn. Successful retries really changed disk; they were not simulated by a no-op refusal-only agent. File effects are directly observed. The source's single gate precedes every handler, covering read/search/list dispatch as well; those read-only handlers have no independent filesystem marker in this CLI fixture. Chapter 6 real streaming arrival order was deliberately not required or tested.

### Tools, declarations and command contract

Source and tests support six registered tools: `list_directory`, `read_file`, `search_files`, `write_file`, `edit_file`, `run_command`. The registry is per-Agent; tools have real parent Agents. Tool definitions and provider wires are checked by the candidate tests, inherited grader and retained real-vendor evidence. Gemini uses a typed compatible Schema projection rather than copying unsupported JSON Schema keywords wholesale.

Reviewed/tested behavior includes deterministic directory listing, bounded/ranged reading, recursive regex search with context, guarded create/overwrite/append, exact unique-anchor edit refusal, unknown/malformed argument errors, strict object decoding, and recovery after error results. The documented unique-anchor edit policy is an allowed public choice, not a new reviewer threshold. The final inherited grader separately exercises search context and append.

`run_command` preserves stdout, stderr and shell exit status separately in its JSON result. An ordinary nonzero shell status is feedback, not itself a transport/tool failure. Missing command, cancellation and launch/completion failures remain errors. It preserves ordinary environment entries while filtering the four supported provider-key names, drains bounded retained output, uses a process group for cancellation, and does not impose an arbitrary 30-second command deadline. The post-cancellation wait/drain bound is not a generation/command lifetime timeout. Source/test checks and the independent credential-filter mutation support these statements; this is not a sandbox/security claim beyond the specified behavior.

### Ownership, History, provider and lifetime regressions

The public root constructs Agents, Engine and History; children retain parent interfaces, and the Agent alone orchestrates response/call/result facts. Engine parses/renders/sends and accounts usage but does not own tools or mutate History. History remains the append/dump/load/project owner; CLI is an adapter rather than the tool-loop owner. Public API tests and the preserved empty-registry behavior cover unsupported unexpected calls without leaving them unanswered.

Reviewed provider results retain native call IDs and provider-valid result shapes, and fallback IDs are generated before the response fact is committed. Native and synthesized IDs remain associated with result facts. Candidate tests cover vendor declarations/results and signature binding; the reviewer added a two-ID-less-call Gemini test described below. Usage accumulation, measured zero versus missing usage, failed-attempt accounting, ordered replay/redaction, source provenance and opaque isolation continue under inherited tests and the grader composition. R1 is the explicit exception for ordinary signed Gemini text.

The independent race suite includes the inherited long generation phases and cancellation tests. The separate Edition 3 Chapter 1 grader also completed its **real 62-second delayed-header and delayed-body phases**. Missing credentials/configuration checks, sanitized transport diagnostics and no hidden total-generation deadline remain part of the inherited source/test review; not every failure permutation was independently recreated live. No full-book run was attempted.

## Tests and sensitivity: actual outcomes

| Independent run | Result and limits |
|---|---|
| `go test -race -count=1 ./...` in archive scratch | Exit 0; root package approximately 126 seconds; internal/CLI packages compile. This was the pristine baseline test invocation, not the added intentionally failing regression fixture. |
| `go vet ./...` on a separate pristine archive | Exit 0, no diagnostics. |
| `go run ./cmd/grade-edition2 -ch 3 -json <manifest executable>` | Exit 0, `passed:true`, 100/100. Includes inherited Edition 2 Chapter 2 checks. |
| `go run ./cmd/grade-edition3 -ch 1 -json <same executable>` | Exit 0, `passed:true`, 100/100; both actual 62-second phases completed. |
| Independent public-API signed-Gemini fixture | Builds, then fails behavioral assertions described in R1. |
| Actual-CLI local matrix | Ten policy/round scenarios pass their expectations; two Gemini scenarios independently reproduce empty-success failure. |
| Retained Chapter 2 five-control script, independently rerun in scratch | All five mutation failures and all five restorations pass their intended detection criteria. This is not the Chapter 3 set. |
| Eight independently constructed Chapter 3 mutations | Seven detected immediately by existing assertions; fallback-ID collision initially survives the chosen selected-policy test. Every restoration passes. See below. |
| Added public-API fallback-ID check | Candidate passes; collision mutation fails an explicit distinct-and-paired identity assertion; restored candidate passes, including dump/load/render pairing. |

Independent single-defect controls were: remove policy gate; apply it universally; allow late text; allow whitespace-only text; change ten rounds to eleven; collide fallback IDs; preinstall tools in ordinary library Agents; stop filtering supported credential keys. `sensitivity-index.json` intentionally retains the original survivor rather than rewriting it as a pass. `sensitivity-details.txt` contains exact temporary diffs and assertion logs. `fallback-*.txt` records the subsequent baseline/mutation/restoration sequence. Mutation failures were assertions, not failed compilation. The initial aggregate mutation driver therefore exited nonzero on the survivor; that is a truthful coverage finding, not a baseline candidate failure.

The student's **17** retained Chapter 3 mutation records were separately audited by reading every mutated/restored transcript and index: each reports a nonzero mutation run with a test assertion failure, a zero restoration run, and no compile failure. The exact student index lists: `gemini-schema-extra-property`, `gate-disabled`, `opaque-counts-as-explanation`, `blank-counts-as-explanation`, `late-text-counts-as-explanation`, `dispatch-skipped`, `result-id-corrupted`, `extra-round-allowed`, `registry-declarations-empty`, `gemini-result-id-lost`, `overwrite-guard-removed`, `read-cap-removed`, `search-context-removed`, `credential-filter-removed`, `nonzero-command-treated-as-tool-error`, `cli-policy-always-cleared`, `ambiguous-edit-allowed`. **Report-only transcription correction:** universal-policy and default-registry mutations belong to the reviewer's independent controls; strict argument decoding and Gemini call-signature tests were not among these 17 student mutations. `student-mutation-audit.json` records the assertion names. Their original mutation patch generator is not in the 20-file source archive; this review does **not** claim to have independently replayed those exact 17 patches. The independently specified eight controls and additional fallback check provide separate sensitivity evidence.

Review-command mistakes are not hidden candidate failures: an initial unsupported `-student` grader flag produced usage only and was corrected to the positional executable; early scratch fixture edits used the wrong public parser signature and omitted required configuration; the first concurrent vet attempt saw those incomplete fixture edits. The corrected public fixture builds, and vet was rerun on a wholly separate pristine archive. A script/evidence-path lookup and receipt `kind`/`type` mismatch were likewise review setup errors. None is counted as a candidate or grader result. All owned jobs are finished.

## Required live evidence audit and limits

No fresh paid request was made by this reviewer. The supplied sanitized receipts were hash-verified, read and summarized independently in `live-evidence-audit.json`; this is retrospective evidence audit, not a claim of directly observing those network sessions.

- Final protocol and human-chat receipts use Anthropic `claude-sonnet-4-6`, OpenAI `gpt-4.1-2025-04-14` via Chat Completions, and **Gemini `gemini-3.8-flash`**, not the prior chapter's 2.5 model. The earlier model discrepancy is not reused as authority. The root's accepted Gemini 3.8 access receipt remains minimal access proof, not proof that this candidate handles its replies correctly.
- For each successful repair receipt, all six distinct tool names occur and every call has a corresponding same-ID/name return. Baseline disk test exit is 1; independent post-repair test exit is 0; expected source bytes match; test file remains unchanged. Protocol call/return counts are 11/11 for each vendor; chat counts are 12/12 Anthropic, 10/10 OpenAI, 11/11 Gemini.
- Anthropic/OpenAI protocol answers are recorded nonblank. Human-chat answer evidence is carried in the saved transcript/final-response summaries rather than the protocol `assistant_replies` array (which is empty for chat). Their results support working conversations and follow-up memory. OpenAI additionally preserves a real policy refusal and explained retry; that is stronger than merely setting the policy flag in a passing run.
- Gemini's final protocol receipt has two empty answers. The terminal metadata covering **both protocol and chat** shows four terminal responses with nonblank text and a signature but no thought flag; all four were hidden by the local parser. The completed file repairs and CLI exit 0 do not close this human-interface failure.
- The initial Gemini schema-400 attempt and failed repair are retained, not replaced by later success. Its typed-schema correction and later real tool work show that schema acceptance was repaired; that does not fix terminal text classification. Earlier grader/diagnostic failures likewise remain in the student evidence rather than being retrospectively described as passing.
- Human live tests preceded a **Go-comment-only** final delta. The supplied final-source comparison is a reconstruction, not a contemporaneous pre-live source archive. The small comment patch and its attribution were inspected; this is a stated lineage limit, not exact-at-live byte attestation. Final protocol evidence and final binary/source pins have their separately recorded attribution.
- Local fakes, mutation controls and the public-API reproduction are not relabeled as live-provider proof. No token-cap inflation, credential-dependent fallback or extra provider calls were used to diagnose R1.

## Correction gate and handoff

Return R1 for a bounded correction. After a revised frozen candidate, rerun the signed ordinary-text regression (including replay/isolation), the relevant Chapter 2/3 baseline tests, the selected-policy no-effect/retry controls, and the approved two-grader composition on one explicit candidate CLI. Preserve existing successful tool behavior and evidence rather than rewriting the harness or moving to a different model. The candidate currently passes the independent fallback-ID check; retaining equivalent coverage is useful protection, not an additional implementation defect.

The failed Gemini human-answer evidence cannot be represented as successful live closure. Root should authorize the smallest needed successful Gemini CLI revalidation after correction, with clear source/binary lineage, rather than blind repeated calls or wholesale three-vendor/whole-book reruns. This review requests no provider call now and does not spend a correction round itself. No Chapter 3 code acceptance or author timing is granted here.
