# Chapter 19 restart — independent formal text review

## Verdict: REQUIRED_REVISION

One formal review of the held `book/edition-3/chapter-19.md`, after code acceptance. The required change is confined to an existing factual passage in §19.5. The three new paragraphs and historical-scoping phrase are accepted as written; they do not need a controller diary, broader claims, or additional requirements. Root owns the sole bounded author response, verification and closure; no second text review is requested.

## Required finding

**R1 — §19.5 names implementation APIs that the accepted Edition 3 code does not have (lines 455–470).** The paragraph says Chapter 16's `CredentialProvider` supplies the token and `Engine.sendRequest` adds the `Authorization` header; the following error paragraph names Chapter 1's `requestFailure`. These are unqualified descriptions of the current implementation, unlike the explicitly labelled path-check sketch earlier in this chapter. They send a reader to nonexistent symbols and overgeneralize the header name across providers.

Evidence: `solutions/edition-3/main/internal/common/auth.go:11–15` declares `Auth.Resolve`; `internal/auth/auth.go:163` implements it. `internal/engine/engine.go:194–203` owns the exchange in `ExchangeRequest` and resolves OpenAI credentials through the parent auth service. Lines 241–252 set provider-specific authentication headers (`x-api-key`, `Authorization`, or `x-goog-api-key`). Lines 267–273 replace a transport failure with a safe generic error (or the context cancellation error), rather than exposing the credential-bearing transport error; the local method is `failure` at line 48.

**Smallest response:** correct those two short passages to describe the existing auth-owner/Engine transport boundary and credential-safe transport errors. Prefer plain owner descriptions if exact method names add no teaching value. Do not introduce alias APIs, change code, add features, or broaden the exercise. Preserve the rationale that credentials and credential-bearing errors must not reach model-visible context.

**Minor correction in the same response:** §19.4's “section 21.3” at line 400 is a stale cross-reference; the root-canonicalization discussion is in this chapter's §19.3. This alone would not warrant withholding acceptance.

## What passes

- The §19.3 addition separates the illustrative `ValidatePath` sketch from the actual rooted file operations and explains why check-then-open alone cannot establish race-safe confinement. It also preserves the distinction between those in-process operations and OS confinement of arbitrary commands and descendants. The trusted host corpus remains outside the tool root. No chroot or remapped `/` is implied.
- The §19.4 addition makes the metadata-only `/var` distinction concrete and identifies the physical selected developer directory plus developer-link allowance as the observed useful combination. It does not claim that experiment proved each allowance individually necessary. The concrete paths explain the demonstrated alias/selector failure without prescribing a universal workstation path or wider content-read exception.
- The §19.11 addition is a dated, finite API-key Sonnet 4.6 observation, not live OAuth/browser validation or a universal security proof. Actual paired results and artifacts support useful file/Python/Git work alongside outside-root, descendant, filtered-token, network and safe-mode controls. “The completed sessions reached EOF” does not say every attempt or wrapper passed. No boilerplate failure catalogue is needed in the chapter.
- “The original profile experiments” appropriately scopes the older operation table. The existing first-person debugging narrative retains its historical attribution; the new observation does not impersonate a new firsthand author experiment or relabel original results.
- The edit preserves the security requirements, ownership rationale, grounded examples and chapter voice. Its prose is load-bearing rather than quota-driven. The opening and Bill's preface are unchanged. The actual diff is three paragraphs (236 added words) plus the historical-scoping phrase, not a redesign.

## Evidence and limits

Read the reviewer and author role files; complete `AGENTS.md`, workflow, Edition 2 carryover, voice, chapter-writing procedure and current Edition 3 execution guide; the full current Edition 3 chapter; final student, code-review and author notes; saved baseline and held final; and their exact diff. Recovered clipped required reads. Current lean-restart instructions govern, not historical campaigns. A read-only comparison confirmed the working chapter equals `/private/tmp/ch19-author.X0Ldgq/final/chapter-19.md`.

Reused the saved `baseline-lint.log/.exit` and `final-lint.log/.exit` under that author directory: both exit zero, hard checks pass, and soft warnings remain. This verdict is based on text and evidence, not lint acceptance. No lint, tests, builds, provider calls, browser/native operations or code validation were rerun.

For the live wording, read only the small existing receipts and verification summaries under `/private/tmp/ch19-live.qGIHh5/` and `/private/tmp/ch19-live-remaining.MPQUHp/`, not raw response bodies or settings. The first wrapper remains FAILED at its 12-POST cap before safe mode reached upstream; three confined turns reached EOF/exit zero and 15 tool pairs were verified. Plain `cat` did not establish an explicit descendant. The follow-up used five POSTs across two completed sessions and did establish Python → shell → `cat` refusal; `cat` ran and its read was denied. Its wrapper remains FAILED because of the newline assertion, while `verification.txt` establishes exact model-supplied/write/read byte agreement, safe-mode declarations and parent-read refusal. Neither wrapper is relabelled green here. These are bounded saved observations, not a fresh reproduction or exhaustive raw-body audit. The final code-review note and execution guide retain the fuller qualified acceptance.

Only this review note was created. No author contact, child delegation, book/source/snapshot/guide/grader edit, commit or closure action was performed.
