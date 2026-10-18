# Chapter 7 — sole authorized author revision

## Status and authority

This is the one authorized author revision after the completed, joined formal text review returned **REQUIRED_REVISION**. It addresses R1–R4 only. Root directly verifies this revision; there is no second formal review and this report does not declare final text acceptance.

Authority: `/tmp/cadence-ch07-author-r1-task.txt`, 4,117 bytes, SHA-256 `f634c957913949a718f32070de2f1d667f9f4a025f952065562d8a8f895814c9`. The entire [formal review](ch07-text-review.md) was read; its SHA-256 is `f93c83469b7fbeae0c1f578b33f82127c998caf0536a40da0789a94f102568f9`, and its evidence index is `b9d083dc47870f8d75bb0550afe304fe2ac2dda15fc1019faf09c018eae21511`.

Writes are confined to the chapter, this new report, and `ch07-author-r1-evidence/`. The original [author report](ch07-author-review.md) and entire `ch07-author-evidence/` remain byte-for-byte unchanged, including their false descriptions. This report supersedes those descriptions explicitly below; it does not silently repair the old seal. The old index's external chapter entry now identifies the historical reviewed draft, while its internal `chapter-after.md` remains that exact draft.

## R1–R3: finite chapter edits

**R1 — admission ownership.** Changed only “a batch the engine had refused” to “a batch the core had refused.” The requested-versus-admitted lesson, live/replay scope, visible refusal results, and actual-admission meaning of `ToolDispatched` remain. No ownership/API redesign or private replay implementation is prescribed.

**R2 — readable versus opaque speech.** Removed “Opaque media still needs only its stable spoken description.” The full readable tool-result/error lesson remains, as do the chapter's existing signed/redacted/empty-content distinctions. The removed sentence was not a supported universal exception: the original author's claim that it merely retained an existing opaque-media exception was wrong. This revision does not require new descriptions or speech for previously silent opaque controls.

**R3 — paired gates.** The existing Chapter 7 command now writes JSON and captures its exit status. Immediately beside it, the recipe runs Chapter 1 with the identical `/tmp/ch7-cli` and captures that exit separately. The text requires the same frozen binary, no rebuild or GUI substitution between gates, exit zero for **each**, and `passed: true` with 100/100 in **each** report. Required zero-point failures cannot be overridden by a score. The two genuine delayed-response phases remain genuine; the reader is told to allow about two minutes for those waits. The manual GUI `./cmd` recipe, module roots and working directories are unchanged.

The exact [revision patch](ch07-author-r1-evidence/chapter-r1.patch) and same-basename [baseline](ch07-author-r1-evidence/baseline/chapter-07.md) / [final](ch07-author-r1-evidence/final/chapter-07.md) copies show the complete chapter change. No other polish, cleanup, chapter audit narrative, C1 policy change or implementation change was made.

## R4: corrections to the original author's evidence claims

### Fixture — original report line 28

The original report's “protocol-completeness and step-boundary mismatches” and “intended completed-call boundary” account is false. The [fixture acceptance](ch07-pause-fixture-accepted.json) records exactly one legacy argument-name repair in `driveGUIPause`: **`delay` → `ai_callback_delay`**. Assertions, timing and scores were unchanged.

The [harness reconciliation](ch07-student-harness-reconciliation.md) records an authorized repaired Chapter 7 run on unchanged frozen r4 source/CLI: 100/100, `passed:true`, exit 0. The original r4 Chapter 7 report stays 80/100 and `passed:false`; earlier failed attempts remain failed. Fixture-only acceptance, repaired grading and later student-code acceptance are separate decisions. The same-byte r4 Chapter 1 result was carried evidence, not a new Chapter 1 run in that reconciliation. This author ran neither gate.

### Paid/live — original report line 32

The original account of a different-ID call, held completed calls, paused cancellation and resumption does **not** describe this evidence and is withdrawn.

The [live exercise](ch07-live-exercise.md) used corrected source `463996b765631d31923a5726036a34942a67cefd` for an ordinary streamed `read_file` interaction. Exactly two Claude Sonnet 4.6 requests returned HTTP 200. The launcher remains **exit 1, `passed:false`**. Its predicate required `tool_finished.name === 'read_file'`, but the actual completion event has no `name` field.

The [correlated events](ch07-live-exercise-evidence/live-approved-01-correlated-tool-events.json) show one dispatch and nonerror completion with the same call ID, `toolu_013yPgvEruEBrGiQTHcM5GT4`, and the same `part_id`, 3. The [root reconciliation](ch07-live-exercise-evidence/live-approved-01-root-reconciliation.json) establishes byte-identical fixture/job/result content, visible final-DOM result, and that the second authenticated request contained the actual tool result. The final answer used the Lantern Orchard / three reviewers / blue folder facts. Root accepted the combined semantic evidence without a rerun or production change. That does not turn the original failed launcher into a pass and does not prove a paused-call cancellation/resumption scenario. It is not native listening evidence.

### Native and corrected-source controlled checks — original report line 33

The original “provider/native Electron” and “controlled Electron” descriptions are false. The [native report](ch07-native-exercise.md) records disposable **headed Chrome**, earlier source `48e1d12d9c396b5988c0ce6eccdef67e581872ff`, a fixed local provider and the unchanged production speech queue. There was no paid request, credential read or microphone operation. Default voice-service locality remained unknown and explicitly consented to; this is not OS-network-isolation proof.

That report ends with successful native callback observations in attempt 3 and human listening still pending. The separate [attempt-4 summary](ch07-native-exercise-evidence/attempt-4-summary.json) records the repeat's native sequence, cleanup and still-pending hearing observation. The later [human-listening receipt](ch07-native-exercise-evidence/human-listening.json) records Bill's exact quote: `I heard "Quiet Gerdens Rest"`. The intended recovery phrase was “Quiet gardens rest.” This is a qualified hearing observation, not a transcription-accuracy result or a claim that he heard the first canceled phrase, “Harbor lights shine softly.” The earlier failed attempts remain intact.

The corrected-source R2 checks in the [correction code review](ch07-code-review-r1.md) exercised the controlled production-page/DOM/utterance-input boundary. They are not fresh native audio, microphone or paid-provider tests. Earlier human-heard evidence carries only the unchanged speech mechanisms; corrected page-handler controls have their own controlled evidence. The later paid read-file exercise is a third, distinct scope. None was rerun by this author.

## Actual checks and preservation

The [read ledger](ch07-author-r1-evidence/read-ledger.json) pins full current guidance, chapter, review, public brief and directly used accepted reports/primary receipts. The four longer craft prerequisites were reloaded in full bounded ranges rather than relying on compressed earlier context.

Both same-basename prose lint commands exited 0. Baseline lint counted 4,604 prose words; final counted 4,657. Both retain 20 negation forms (warning) and a long person-gap warning (2,343 baseline / 2,334 final). Short closers remain 35; the baseline warning becomes within the final paragraph-count threshold. These are recorded advisories, not a warning-free style claim or authorization to expand this revision. See [lint runs](ch07-author-r1-evidence/lint-runs.json) and their full outputs.

The [mechanical check receipt](ch07-author-r1-evidence/checks.json) records exact diff scope, unchanged headings/links/manual GUI passage, balanced fences, shell parsing without execution, and successful `git --no-pager diff --check`. No application tests, grading commands, provider calls, browser/native exercises, production/student/grader/guide edits, publication rebuilds, commits, tags or children were performed. Running the explicitly requested prose linter is not an application grading run.

All 12 publication/build paths were excluded using **only names** from root's saved observation metadata, before preservation assertions. Their current bytes were not read, hashed, chased or restored. The [preservation receipt](ch07-author-r1-evidence/preservation-after.json) compares 38 non-publication inputs against the revision's own before snapshot, including all original author and formal-review files. The historical author's “30 protected paths unchanged” claim is time-qualified to its earlier pass, not reasserted against moving publication work.

The original author's setup/inspection failures remain in its immutable evidence. This revision's separate [process note](ch07-author-r1-evidence/revision-process.md) records its own mistaken prerequisite path, copied hash assertion, metadata-key assumption, and retracted provisional commentary. Those setup failures stopped before chapter writes. All tool runs completed synchronously; no child or background job remains from this revision.

The [new exact index](ch07-author-r1-evidence/index.json) inventories the revision report, revised chapter and every new evidence artifact except itself. It contains safe text, hashes, static-check receipts and chapter copies—not raw provider payloads, credentials, images or audio. Root owns direct verification and final acceptance.
