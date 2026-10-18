# Chapter 7 — independent formal text review

## Verdict: REQUIRED_REVISION

This is the one formal, whole-chapter TEXT review, not another code review or execution campaign. Four finite corrections below are sufficient. Root owns the single authorized author revision and its mechanical verification; this report does not authorize another review round, implementation change, provider/native/browser run, or publication work.

The chapter is generally coherent and the small author patch addresses real student failures. The required corrections concern three public teaching details and factual errors in the author report. They do **not** revoke the separately recorded student/code acceptance.

## Frozen inputs and coverage

- Review input: [`ch07-text-review-input.json`](ch07-text-review-input.json), SHA-256 `0ad9ace57a802e7e988a5cb84081cdaff1ff4a32cd3758ae04a8cf49b67e02ad`.
- All **17 manifest members** match their pinned byte counts and SHA-256 hashes. The manifest itself is an eighteenth identified file, not an eighteenth member.
- Whole original chapter: `/tmp/cadence-ch07-text-phase-baseline/chapter-07.md`, 33,024 bytes, 692 lines, SHA-256 `68d02d8e1170c2912676b8f162847f8fac25c87216514ed06adbd5f0ea011b47`.
- Whole reviewed chapter: [`book/edition-3/chapter-07.md`](../../book/edition-3/chapter-07.md), 34,792 bytes, 730 lines, SHA-256 `42a97d0a7235442f70ebaa66527177ecf01721f9aaaa074f155527b5d8ad2e35`.
- Original/current bytes equal the author's `chapter-before.md` / `chapter-after.md`. All 10 members of the author's self-excluding inventory match; its evidence directory has the exact nine expected files, including its index.
- I read the complete AGENTS, text-reviewer skill, execution procedure, workflow, carryover, voice guide and writing procedure; both complete chapters; all eleven required Chapter 7 student/review/acceptance/native/live/reconciliation documents; the author report and whole indexed author package; the public student brief; and the directly cited sanitized primary receipts. The read ledger records 41 whole-document reads, hashes and line counts. Large combined outputs were recovered using smaller ranges or standalone reads, not treated as complete while truncated.
- The coding skill was read before a narrow static inspection of the public CLI adapter. That was not a new code assessment. No student tests, graders, builds, browser, speech, provider or native exercises were run.

Evidence: [`ch07-text-review-evidence/index.json`](ch07-text-review-evidence/index.json). Its inventory seals this report and every other member of the reviewer evidence directory, excluding only the index itself. Exact numbered excerpts are in [`cited-extracts.md`](ch07-text-review-evidence/cited-extracts.md); the independent full chapter diff is also retained.

## Required revisions — finite scope

### R1 — Attribute tool admission without assigning policy refusal to Engine

**Location:** current chapter lines 333–337, especially “a batch the engine had refused” at line 334.

The useful lesson is correct: a requested call and its refusal result do not prove execution. However, this chapter explicitly distinguishes Agent from Engine and says the gate belongs to Agent, not an Engine sibling (lines 376–380). The correction review locates the actual refusal/admission reasoning in `internal/agent/model_tool.go` and `internal/agent/events.go`, and explains the GUI's erroneous projection from stored calls. The engine wording is unnecessarily misleading in this ownership-sensitive account.

**Required edit:** replace that attribution with “a batch the core had refused,” or comparably accurate Agent/admission-owner wording. Keep the requested-versus-admitted, live/replay and refusal-visibility lesson. Do not move policy, invent an API, or prescribe the student's private replay implementation.

**Done when:** the paragraph no longer teaches that Engine owns this policy refusal, and still says dispatch observations describe actual admitted execution.

### R2 — Remove the unsupported blanket opaque-media speech requirement

**Location:** current chapter lines 589–592, especially “Opaque media still needs only its stable spoken description.”

The new readable-result lesson is supported: on-demand speech must not substitute an automatic-speech caption for the full readable tool result/error already available to the page. The added final sentence is not a supported universal exception or requirement. It presupposes a stable spoken description for opaque material, while the existing chapter explicitly preserves empty text and signed/redacted content without pretending those bytes are readable (lines 568–578). The correction review separately records no utterance for empty/signed/opaque controls. Those controls do not justify claiming every opaque media card needs a spoken description.

**Required edit:** remove the newly added opaque-media sentence, or narrowly state that the readable-result correction does not turn opaque payloads into speech. Preserve the existing safe/opaque distinctions. Correct the author report's line 21 claim that this sentence merely retains an existing opaque-media exception if that description is carried into the revision report.

**Done when:** the passage teaches full readable result/error speech without requiring new descriptions, speaking opaque signatures, or turning previously silent controls into utterances. No implementation or test change is requested.

### R3 — Include the already-required same-CLI Chapter 1 gate in the public recipe

**Location:** current chapter's new invocation paragraph and command block, lines 221–264.

The chapter correctly distinguishes the core CLI from the separate GUI module, uses the actual CLI package conditionally, identifies the book-repository working directory and requires explicit module/model inputs. Its manual GUI command is also correctly rooted in the GUI module and builds `./cmd` (lines 206–219); **no manual-launch revision is required**.

But the new public grading recipe stops at `-ch 7`. The already-public student brief, lines 107–125, explicitly requires `-ch 1` on the **same frozen CLI**, including the genuine delayed-response phases. The accepted correction report records that exact paired gate. A student following the chapter's new self-contained recipe should not miss it.

**Required edit:** immediately after the Chapter 7 command, add the Chapter 1 invocation using the identical already-built CLI path, for example:

```sh
go run ./cmd/grade-edition3 -ch 1 -json /absolute/path/to/student-cli
```

Say that the binary is the same one just checked for Chapter 7; do not rebuild or silently substitute a GUI executable between invocations. Retain the existing explanation that required zero-point checks, `passed`, and exit status matter. A short reminder that the two genuine delayed-response phases take about two minutes is appropriate; do not invent a shortened deadline or new test.

**Done when:** the public chapter recipe unambiguously shows both existing gates on the same CLI and preserves the separate GUI-module/manual-run explanation.

### R4 — Correct the author report's three materially false evidence attributions

**Location:** [`ch07-author-review.md`](ch07-author-review.md), lines 28, 32 and 33. These are not merely missing qualifications; the report describes different mechanisms or exercises from its linked sources.

1. **Pause fixture, line 28.** The report says the original predicate had “protocol-completeness and step-boundary mismatches” and the repair used an “intended completed-call boundary.” The reconciliation report instead attributes the accepted repair to **one legacy argument name**. `ch07-pause-fixture-accepted.json` explicitly records `driveGUIPause` argument `delay -> ai_callback_delay`, with assertions, timing and scores unchanged. Replace the invented predicate/boundary explanation with this actual narrow repair. Preserve the original failed grades and the distinction between fixture acceptance, repaired grading and student acceptance.

2. **Paid/live exercise, line 32.** The author describes an unexpected different-ID call, held completed calls, cancellation while paused and later resumption. The linked live report and primary correlation show a **two-request read-file exercise**, not that pause scenario. `live-approved-01` ran on corrected commit `463996b765631d31923a5726036a34942a67cefd`; its launcher exited **1**, `passed:false`. Its predicate incorrectly required `tool_finished.name === 'read_file'`, a field absent from the production completion event. Root correlated one dispatch and nonerror completion by the **same call_id and part_id**, matched the result to the fixture/job output, and confirmed that the second authenticated provider request and final answer used that result. State this semantic reconciliation accurately without converting the original launcher into a passing artifact or claiming it proved paused-call cancellation/resumption.

3. **Native and controlled browser evidence, line 33.** The report calls the earlier exercise “provider/native Electron” and the later checks “controlled Electron.” The native report identifies a disposable **headed Chrome** run at earlier commit `48e1d12d9c396b5988c0ce6eccdef67e581872ff`, with a **fixed local provider**, unchanged production speech queue and no paid request. Its report ends at attempt 3 with human listening still pending; separate attempt-4 and human-listening receipts establish the later heard recovery phrase. Bill's exact recorded quote is `I heard "Quiet Gerdens Rest"`; the intended phrase was “Quiet gardens rest.” That is a limited human hearing observation, not transcription accuracy or proof that he heard the first canceled phrase. The corrected-bundle R2 checks used a controlled production-page/DOM/utterance boundary, not a fresh native audio or paid-provider run. Replace the Electron/provider attribution with these actual scopes, keeping the earlier native/human evidence separate from corrected-source controlled checks and the later paid read-file exercise.

**Required edit:** publish these factual corrections in the root-authorized revision report/evidence, and align any retained summary of the opaque-media sentence with R2. Do not overwrite the original sealed author package, inherited receipts, launcher failure, acceptance records or this review input. Do not add an audit-history paragraph to the chapter merely to accommodate these report corrections.

**Done when:** all three report attributions match the linked primary records; the failed launcher remains failed; evidence is qualified by the source/run actually exercised; and the old sealed package remains intact.

## Whole-chapter assessment

**Voice, structure and teaching value.** The chapter retains its technical register and human trail: the planned thin transport encounters observation ambiguity, pause and reconnect lifetimes, speech queue behavior, visible failures and narrower corrections. Its existing browser-tab symptoms, zero-volume wakeup compromise, mutation that initially passed for the wrong reason, and replay accounting are concrete rather than generic process narration. The new short requested-versus-executed and readable-result examples fit that trajectory. The added 1,768 bytes / 38 net lines do not call for broad restructuring. No required style rewrite, new anecdote or additional figure is warranted.

**Contract and architecture.** The chapter keeps the optional GUI module, public core boundary, same-Agent use, explicit pause control, bounded replay, separate observations and stored facts, and inherited tool-job lifetime distinctions. The command clarification is useful and does not prescribe an invented root executable. R1 removes an avoidable ownership misattribution; R2 prevents an unsupported feature obligation; R3 makes an existing public gate usable. None requires altering accepted code.

**Evidence and narrative honesty.** The initial failed grades, one-field fixture repair, two accepted code corrections, and same-binary paired grading must remain distinct. The accepted manifest/root decision support code acceptance and qualified carried native evidence; they do not make the author report's descriptions true. The paid launcher is still a failed artifact with independent semantic acceptance, not a green launcher report. The native listening observation belongs to the earlier source/run; controlled corrected-source speech-input checks are not fresh audible proof. I did not inspect screenshot images and claim no independent visual/aesthetic or human-listening result.

**C1 and scope.** The initial code review's C1 is the measured page-hide/unpause release behavior, retained by root as non-blocking contract clarification, not a third code correction. I do not reopen it, introduce hidden-tab cancellation, or require a new disconnect/gate policy. A provisional reviewer comment about C1 was wrong and was explicitly retracted after reading the primary passage; it is not a finding here. Likewise, no issue is raised against the correctly stated manual GUI working directory.

## Checks, preservation and limits

The reviewer ran only document reads, static source/path inspection, hashes/inventory comparisons, a textual original/current diff and read-only Git status. The bundle records actual checks rather than inheriting the author's check claims as new reviewer results. No prose linter, application test, grade, build, provider request, browser, microphone, native speech, child agent or commit was run or created.

The first protection assertion failed **before any reviewer writes** because six publication files no longer matched the older root baseline. A subsequent pre-write snapshot observed a seventh publication-document change. `input-verification.json` records the exact 23 matches and seven historical-baseline differences. These concern `book/README.md`, both book PDFs, `book/build-print.py`, `book/interior-print.html`, `book/print-header.html`, and `docs/edition-2-publication.md`. They are recorded as concurrent/out-of-scope state, not attributed to the author or silently restored. Root was notified. Root subsequently confirmed the seventh difference in `/tmp/cadence-ch07-publication-observed-during-text-review.json` (05:04:24Z) and directed no more publication checks. This reviewer writes only this report and `ch07-text-review-evidence/`; the final preservation receipt compares non-publication protected inputs and the authoritative input hashes. An initial final-check filter excluded only the seven known-changing paths and inadvertently compared the EPUB, finding further concurrent drift before the seal. This inspection mistake is recorded, not attributed to the author. The corrected final check excludes all publication artifacts/build paths in the saved snapshot, makes no further publication reads, and does not pretend the old publication baseline is still current.

Inspection failures and corrected provisional statements are retained in `review-process.md`. None triggered source edits or an extra test campaign. The sealed evidence inventory excludes itself to avoid a circular hash; final handoff supplies its external SHA-256.

## Root's single-revision checklist

1. R1: change the admission-owner attribution only; preserve the dispatch lesson.
2. R2: remove/neutralize the unsupported opaque-description obligation; preserve readable result/error speech and silent opaque controls.
3. R3: show the same-CLI Chapter 1 command beside Chapter 7; preserve module roots, working directories and genuine delayed phases.
4. R4: correct the fixture/live/native author-report attributions in new revision evidence; preserve all original seals and failed outcomes.
5. Mechanically verify those bounded edits and seal the new revision without reopening implementation or spending another formal text-review round.

No further required revisions are proposed.
