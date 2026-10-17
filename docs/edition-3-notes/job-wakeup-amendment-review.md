# Independent text review: uniform jobs and foreground wakeup amendment

**Verdict: PASS. No blocking findings; no author revision requested.**

This is the root-owned, single formal TEXT review of the expressly authorized Chapter 4/5 amendment, including the user's queued-tail clarification supplied before the verdict. It is not implementation review, student acceptance, or the ordinary post-student author pass. No second reviewer or child was created. The timing exception applies only to this requested amendment.

## Exact verdict target

Root confirmed these replacement frozen pins in handoff `main-17`; independently checked against current bytes and all **28** author-inventory entries:

| Artifact | SHA-256 |
| --- | --- |
| `book/edition-3/chapter-04.md` | `030c9a7ceb6c18f56596e840bb9df404fa418ea5f22c77bc7712b7ec1e6d8c55` |
| `book/edition-3/chapter-05.md` | `4f4c2595ef65a08511d5cc68553f853feaf854e4cb0fe7d4d9094d4ae7f5a344` |
| `docs/edition-3-notes/job-wakeup-amendment-author.md` | `98c91c1b8a7546e077a359bd384e790a3e7062b8cdc28074191b1c2946c3ce6a` |
| Author `authored-files.sha256` | `6072683d89a958edd43dafe62c5a2c73fecdc0f400bd0833a20ffec55ed87b7c` |
| Author `queued-tail-clarification.diff` | `6ee49e8fe3e664fd17bf210ebe2afd7d7861b027483336a4e87e2a5ffd815240` |

The latter two files are under `docs/edition-3-notes/job-wakeup-amendment-evidence/`. The original Chapter 5 pin `0aa8a206…`, original author pin `7f5bfc71…`, and original 14-entry inventory `60953e0b…` were verified earlier but superseded before a formal verdict. Their retained history is not this verdict's target. The clarification was a direct scope instruction, not an author revision prompted by a completed review.

## Findings and disposition

### 1. Uniform jobs and existing limits: satisfactory

Chapter 4 lines 104–111 explicitly put quick file reads, commands, and application-registered tools through the same ordinary-job path, with later messaging/join/MCP adapters inheriting that path rather than requiring their frameworks now. Lines 264–271 distinguish a control's target from a newly allocated job; lines 295–310 distinguish shared observation from stdin capability and require immediate retained terminal reads without rerunning handlers.

Lines 383–437 retain the existing defaults and `tool_limits` name, maximum-silence—not job-timeout—meaning, output bounds, literal-next-call consumption, control/error consumption, same-response ordering, call-specific precedence, and second-setter reset. The additions do not introduce an unnecessary renamed setter or a scheduler framework. Retained-output and cursor teaching remains in context. These are useful clarifications, not a replacement Chapter 4 contract.

### 2. Prompt wakeup and real question delivery: satisfactory

Chapter 5 lines 505–511 preserve the actual addressed question and sender/recipient/request-turn metadata, separate from lossy observation and the queued blocking `Ask`/`Submit` contract. Lines 605–646 explicitly cover initial observations and later `wait_for_job`, honest original-call-id results, a stable existing job handle, caps and cursor retention, immediate available/terminal reads, one return per observation, preservation of both message and completion, and retiring timers to prevent ghost returns.

The actor receiving a hint is not presented as sufficient: lines 535–540 require waking the foreground observation worker too. Lines 586–592 preserve the already-sent HTTP boundary. Lines 724–728 and 813–820 keep wakeup distinct from canceling work, Interrupt, shutdown, and killing a job. The next model request must contain the actual message and correlated result; a mere acknowledgment is not enough.

### 3. Same-batch delay and queued tail: satisfactory after the authorized clarification

The initially read draft left room for the retained all-results/serial-batch rule to delay a hint behind a later long observation. No formal verdict was issued on it. Final lines 624–637 close that gap explicitly: preserve original assistant calls/signatures and completed results; report truthful observed state for dispatched work; pair every undispatched call with its exact-id interrupted/not-executed error; invoke no tail handler, effect, or new job allocation; never automatically replay the tail; and never hold the question behind another queued long wait.

The serial dispatch boundary settles races without pretending started work was unexecuted. Pending one-shot limits are spent on the literal next call's error, and an unexecuted setter cannot arm settings. Fresh model calls are expressly permitted while the original job still runs. This is a **model-response batch observation boundary**, not actor retirement or a prohibition on answering until job completion. Lines 1054–1063 require a corresponding mixed-batch regression including side effects, setters, dispatch races, and fresh calls. No new scheduler API is prescribed.

### 4. Bill's incident and evidence boundaries: satisfactory

Lines 648–666 attribute the unanswered student questions, timeouts, stuck `join_agents`, and `message_fleet` attempt to Bill. “Perhaps in desperation” remains a possible motive, not a proved one; the text disclaims a reconstructed wire trace. The causal lesson is concrete: parent waits for completion while students need an answer. Increasing a timeout does not resolve that dependency. The parent can answer and resume the same join id. Neither `message_fleet` nor child spawning becomes a Chapter 5 prerequisite.

Lines 1032–1068 request bounded offline checks through real Agent/job interfaces, a held command and registered join-like handler, with a fake provider boundary. A configured 3,600-second delay is not an hour-long test. Separate short controlled race checks address obsolete timers. The existing `not-deaf` marker is expressly insufficient (lines 248–254). Existing grades and historical live demonstrations are not represented as new-feature proof.

### 5. Scope, voice, and prose load: satisfactory

Independently reproduced diff hunks match the author's exact diff: Chapter 4 **28 additions / 5 deletions**, Chapter 5 **132 additions / 9 deletions**, including the **26 added clarification lines**. The amendments stay in the job contract, actor explanation, incident, and bounded exercise obligations. The frozen baseline/final pairs—not unrelated worktree changes—define this judgment.

The additions carry necessary distinctions: context cap versus work lifetime, mailbox admission versus model delivery, executed versus undispatched calls, prompt input versus Interrupt, and observation versus reliable completion. Their repetition across TL;DR, mechanism, and exercise serves different teaching purposes. The held-test example and dependency story explain why the rules matter; the text does not become a feature inventory or an implementation mandate.

Both Bill prefaces remain byte-equal to HEAD. Chapter 4's opening before §4.1, Chapter 5's opening before TL;DR, Chapter 4's historical §4.9 demonstration, and Chapter 5 from §5.11 onward remain byte-equal to the author baselines. The historical examples and voice are preserved.

Author same-basename lint receipts report exit **0**, **zero HARD** findings both before and after. Chapter 5's soft word warning is disclosed: **6,613 → 7,586 → 7,854** words (baseline, initial amendment, clarification). Chapter 4 is **8,369 → 8,603**. Soft negation/person-gap warnings remain. I do not request arbitrary pruning to reach a metric: the added prose is load-bearing, and no bounded clarity defect remains.

## Evidence and limits

The [evidence index](job-wakeup-amendment-review-evidence/index.md) records full required reads, exact checks, and receipts. Both entire chapters were read; the changed Chapter 5 was reread in full after the replacement freeze. Capped reads were recovered. Original Chapter 4/5 student reports, initial code reviews, and correction reports were read as historical context, not new-feature evidence. Chapter 4's old acceptance remains historical; Chapter 5's corrected `6c4561` candidate and parallel old-contract R1–R3 review are not accepted here and do not prove this amendment.

Only reviewer-owned report/evidence files were written. No manuscript, code, grader, guidance, protected reference, old receipt, or other worker report was edited. No commits, staging, push, installation, service restart, credentials, providers, runtime tests, graders, or child agents were used. Author lint receipts were inspected, not presented as independently rerun lint.

Two local receipt-generation failures are retained and disclosed: an incorrect Chapter 4 heading lookup, then incorrectly transcribed replacement expected hashes. The second failed assertion correctly withheld certification despite a self-consistent inventory. Root confirmed the exact actual pins; only reviewer expectations were corrected. Final independent checks pass. A final readback also corrected transcribed guidance lengths and a baseline word count, and moved the two tool-written review documents from the file-tool workspace to the intended repository; the misplaced copies were removed. These were reviewer evidence-handling errors, not manuscript defects or requests for another author revision.

**Disposition:** PASS for the exact frozen text above. No actionable author changes are required. This verdict establishes a coherent, bounded teaching contract; it makes no claim that production, the student candidate, or the graders implement or validate the new behavior.
