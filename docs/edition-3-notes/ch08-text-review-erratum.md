# Chapter 8 text-review erratum: R1 attribution

## Scope and status

Root fully inspected the initial text-review report and caught a factual error in its evidence-attribution bullet. This erratum acknowledges and corrects that error. It is **narrow report reconciliation, not a second formal text review**. The formal text-review budget remains **one**. No author revision, chapter/code change, new review commission, lint/grade run, browser/provider exercise, or full chapter reread was performed for this erratum.

The original [text-review report](ch08-text-review.md), its [inventory](ch08-text-review-evidence/inventory.json), and its saved submission remain unchanged. Their hashes are preserved in the [erratum pins](ch08-text-review-erratum-evidence/pins.json). The original seal establishes artifact identity, not the factual correctness of the disputed sentence. This erratum supersedes that attribution; it does not silently rewrite the original record.

## Original statement — incorrect

`docs/edition-3-notes/ch08-text-review.md`, line 37, states:

> The initial independent review's blocker was one-tab TTS cancellation after another tab changed the global setting. The r1 review closes that specific defect and related coverage; the earlier failure remains part of the record.

The asserted production cancellation defect and its alleged repair were not the independent review's finding or R1 correction. That attribution was my error.

## Corrected facts and causal distinction

The initial independent review's **sole blocking finding was test/integration-maintenance R1**: the inherited `gui/test/correction-r1-page.mjs` fixture could not initialize against the current application. Its unchanged run failed with `TypeError: window.matchMedia is not a function` during `setupWorkbench`, before its behavioral assertions. The initial review expressly found **no production-code blocker** in the exercised settings, budget, speech-rate, or ownership seams.

The initial review also identified missing current workbench/root-style boundaries, a stale Transcript-only artifact locator where artifacts now route to Actions, and an obsolete direct read-aloud `.onchange` hook rather than the current accepted-settings handling. Those additional incompatibilities were source findings; the original failing run reached only the initialization error. The broken regression entrypoint therefore did not exercise its retained assertions, including cancellation-generation behavior. This was a real test-maintenance blocker to whole-candidate closure, **not a demonstrated production cancellation defect or native browser/audio failure**.

The initial reviewer adapted an owned diagnostic copy and obtained a pass against identical production app bytes. That diagnostic did not itself fix the submitted fixture or erase its failure. R1 then corrected only the checked-in fixture. The follow-up review records exactly one changed tracked path, `gui/test/correction-r1-page.mjs`, with all other 125 source members matching the parent; production, Go tests, module files, grader inputs, and APIs were unchanged.

The corrected fixture supplies the current boundaries, locates artifacts across their intended panes, and tests cancellation through the real settings sequence: dispatch `change`, inspect the sparse `update_settings` message, verify rollback/no premature cancellation, then verify cancellation on the accepted snapshot. Its passing standalone, real-wire, workbench, and queue checks support **fixture-compatibility closure while retaining the behavioral assertions**. They do not establish that R1 repaired a production cross-tab TTS defect.

**Replacement attribution:**

> The initial independent review's sole blocker was an inherited production-page integration fixture that failed during initialization and had stale boundary, locator, and cancellation-entrypoint assumptions; no production-code blocker was found. The fixture-only R1 review closes that test-maintenance finding with production code unchanged. The initial fixture failure remains part of the record.

### Exact source citations

- [Initial independent review](ch08-code-review.md), lines **3–5**: blocking test/integration-maintenance R1; no production-code blocker; diagnostic distinguished from a browser defect.
- Initial review, lines **33–51**, particularly **45–49**: initialization failure, additional source-identified incompatibilities, unexecuted assertions, identical-production diagnostic, and minimal fixture correction.
- [R1 independent review](ch08-code-review-r1.md), lines **3–5**: narrow fixture-compatibility finding closed, with no browser/native/live acceptance claim.
- R1 review, line **12**: sole changed fixture and unchanged other 125 source members, production, tests, modules, grader inputs, and APIs.
- R1 review, lines **18–27**, particularly **22–24**: repaired boundaries, actual accepted-settings cancellation sequence, current artifact locations, retained assertions, and passing fixture controls.

Exact excerpts and full-file size/SHA-256 pins are saved in the erratum evidence. These are readings of existing reports, not fresh runtime results or an independent inspection of student implementation.

## Effect on the existing manuscript PASS

**The existing manuscript PASS remains unchanged; this is not a new PASS or a second review.** The correction materially changes the report's explanation of what failed and what R1 closed. That false causal explanation must not be used as evidence of a production defect or production repair.

The narrow reconciliation does not expose an unresolved implementation blocker: the initial review classified R1 as fixture maintenance, and the follow-up explicitly closes that same finding with production unchanged. It therefore does not invalidate the acceptance endpoint on which the original text review relied. Nor does this correction identify a manuscript sentence requiring author revision: the mistaken assertion is in the review report itself. The original text review's grading-recipe assessment is not dependent on the asserted cross-tab production repair.

This conclusion is limited to the effect of the identified factual correction. It does not independently revalidate the rest of the report, expand the earlier evidence claims, or certify flawless reviewer process. No chapter reread, lint, grade, browser/provider run, or additional formal pass was used to reach it. Root will verify this erratum before closure.

## Saved tail and seal

The [saved erratum submission](ch08-text-review-erratum-evidence/final-submission.json) records this bounded disposition and the erratum identity. The [erratum inventory](ch08-text-review-erratum-evidence/inventory.json) seals only the new evidence and external erratum, explicitly excluding itself. Source/report identities are checked again before sealing, and the new seal is validated read-only after its final write. The original report and original inventory remain untouched.
