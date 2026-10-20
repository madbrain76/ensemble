# Restarted Chapter 8 — formal text review

**Verdict: PASS.** No required revisions. This is the single bounded formal review of the initial restart author pass; no author response or further text-review round is needed.

## Scope and findings

Read the current AGENTS, text-reviewer and author roles, Edition 3 execution/carryover guidance, complete current chapter, voice guide and chapter-writing procedure, author note, and full current-run student/code-review notes. Applied the current lean scope rather than reopening historical workflow audits. Reviewed the actual four-hunk manuscript diff (19 insertions, 10 deletions) and surrounding text, with narrow source checks in the accepted `solutions/edition-3/ch08` snapshot.

The four changes teach existing demonstrated requirements, not new acceptance criteria:

- **Captured-turn budget, chapter lines 138–145:** correctly keeps the final-response stop-tools instruction request-local, leaving History/signed content intact and preventing leakage into later turns. Accepted `internal/agent/actor.go`, `startNext` and `send`, captures the limit and adds the instruction to request ephemera. This addresses the current-run request-context correction without introducing a second policy owner.
- **Follow control, lines 244–245:** independent scrolling and clearance from transcript text express the demonstrated overlap correction without prescribing a redesign. Accepted `gui/web/gui/style.css` keeps `#follow` in normal flow; current-run notes separately record rendered geometry and click behavior. The sentence makes no direct visual-observation claim.
- **Speech and typing, lines 488–496:** distinguishes automatic Read aloud cancellation from explicit Speak (`on demand`), preserves the draft, protects fresh typing against stale automatic acknowledgement, and retains explicit any-tab Continue authority. These distinctions match accepted `gui/web/gui/app.js` settings, speech classification, state/continued handlers and typing controls, and the demonstrated ordering corrections.
- **Fractional speed, lines 515–521:** distinguishes exact saved/control values from the native 0.5–3 clamp. Accepted HTML uses `step="any"`; every native utterance construction applies the clamp. The older 0.05-step example remains explicitly historical, and the text retains the limitation that its earlier-range check does not prove the extension. Native rate values are not represented as measured speedups.

The opening, personal voice, architecture, failure explanations and useful reasons remain intact. No additional framework, personal experience, tree lifecycle work or unrelated scope was introduced.

## Evidence limits and lint

Actual dynamic System-theme change-event delivery remains **UNVERIFIED**. The incomplete visual-model reply is not a PASS. The review makes no hearing or direct-root-view claim and does not promote historical chapter receipts into restart results.

Checked the author's saved baseline/final lint receipts at `/Users/bill/projects/cadence-memories/children/e3-restart-ch08-author/io/2314` and `2315`: both exit 0, with hard ceilings met. The note accurately retains soft warnings: negation forms 11 against allowance 6; person gap 3156 → 3253 words; short closers unchanged at 14/69. This is not a warning-free result, and those diagnostics do not justify removing useful explanations or inventing first-person material.

Only this review note was written. No manuscript/source edits, children, tests/builds, browser/provider/key operations, publication, staging or commits were performed. Root retains direct final verification.
