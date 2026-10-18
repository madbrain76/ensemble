# Chapter 7 — bounded real-model GUI exercise plan

## State and ownership

**Prepared and locally validated; real-provider/browser/key execution NOT authorized or performed.** Chapter 7 corrected source `463996b765631d31923a5726036a34942a67cefd` has passed independent code re-review. This plan closes the remaining existing real-interface + real-model exercise requirement, not another implementation round. No student source, manuscript, grader, model capability, or application setting changed during this preparation.

The same student prepared the policy, launcher and real-executable local controls. Its model transport failed after completed control attempt 03 (`stream ID 583; INTERNAL_ERROR`), then failed again during packaging recovery (`stream ID 13; INTERNAL_ERROR`). The second interrupted turn ended in saved error event 1600, not a final submission. Root retained completed receipts, did not rerun the student because of those transport failures, and took ownership of final harness inspection/packaging. Root then measured and repaired one independent harness-only socket-timeout defect, described below. Do not describe the interrupted plan turn as submitted or the paid exercise as complete.

## Proposed single exercise

Use the unchanged frozen GUI and its actual embedded assets in a disposable, explicitly muted, headless Chrome profile. Model: **Anthropic `claude-sonnet-4-6`**, endpoint `https://api.anthropic.com/v1/messages`. One conversation asks it to read a harmless local `chapter7-note.txt` and summarize its three facts: project Lantern Orchard; three reviewers needed; blue folder. No browser UI campaign, hints/reconnect/native speech repetition, other providers, discovery, or microphone exercise.

Observe real streaming content in the page before final completion, a real `read_file` job/result, return of the actual file contents in a subsequent model request, and a meaningful final answer. Capture partial/final screenshots and safe receipts. Root must assess the answer's meaning; no arbitrary exact response phrase is required. The DOM result is an actual browser observation, not a wire-event-only substitute.

**Spending/observation bounds:** maximum four upstream requests total; the application's existing **4096 output-token maximum per request** remains unchanged (maximum 16384 allowed output tokens across four requests, not expected usage). Root withdrew its initially proposed 1024 cap rather than rewrite requests or add an application feature. Ninety-second observation/forwarding window, 120-second enclosing watchdog plus bounded owned-process cleanup. No automatic paid retry, redirects, provider/model substitution, fallback, or discovery. Stop on the first invalid request, refusal/error, unavailable model, scope breach, or failed observation, preserving partial evidence. Earlier chapters' paid calls are not evidence for this candidate.

## Exact proposed command (DO NOT RUN until explicit paid/browser/key approval)

From `/Users/bill/projects/ensemble`:

```sh
node docs/edition-3-notes/ch07-live-exercise-plan-evidence/launcher-root-r1.mjs live \
  --label=live-approved-01 \
  --launcher-sha256=4fc67b564f0af1610397b08cd565a31fa40968c4abb90e06a95c44f6ea1e4f43 \
  --policy-sha256=633498cb220fd1f88e157334f09dd870ac9389b02937670fe418ad4574dcadc3 \
  --explicit-later-root-paid-browser-authorization
```

The flag is a fail-closed mechanism, not authorization in itself. `policy.json` pins all 120 source files, frozen GUI, eight embedded assets, provider route/model/tool declaration, fixture and spending limits. Startup verifies source/GUI before credentials or processes, served assets before browser interaction, and distinct GUI/CDP listeners and spawned PID/command/profile before attaching. Root inspected the complete final launcher and policy; they remain unexecuted in live mode. Future mutation requires a new hash and review, not silently bypassing the pin.

## Credential and isolation boundary

Only after approval, load `~/.cr/settings.json` field `directClaudeAPIKey` inside the Node proxy process. Do not print settings or key, place it in arguments or temporary files, or pass it to children. The GUI receives a dummy key and only the loopback proxy endpoint; child environments are explicitly constructed. Key stays in process memory for the authenticated upstream request. Native audio stays muted and read-aloud unchecked; no microphone, focus manipulation, shared/default Chrome profile, personal browser, live port 8088, live GUI navigation/restart, or installation.

Accepted request bodies and response frames remain byte-identical; tool blocks are held until their harmless file-read scope is verified, then released unchanged. An out-of-scope tool stops the disposable exercise before actionable delivery, not a product restriction or fabricated successful model result. Atomic pre-forwarding reservations enforce at most four requests, including concurrent fourth/fifth arrivals. Non-200/redirect responses fail instead of being followed. Request/response hashes and usage are recorded; arbitrary upstream error bodies and key headers are not exported. The browser uses owned loopback request interception and DNS/background-network suppression; these flags are not an OS sandbox or a claim of complete network privacy.

## Completed local proof and retained failures

Student control attempt 01 failed an assertion because later stream-end bookkeeping replaced the first `TOOL_SCOPE` stop reason; no forbidden effect occurred. Attempt 02 failed because one SSE output path bypassed the newly added output-hash accumulator; this was broken harness instrumentation, not proof of actual body rewriting. Both receipts and launcher/policy versions are preserved. Attempt 03 passed six controls on the actual frozen GUI: streamed fake-provider forwarding, actual file read and next-request content; unchanged accepted request/response bytes; synchronized fourth/fifth request boundary; wrong model and token cap rejected before upstream; wrong GUI pin rejected before launch; and out-of-scope tool stopped with zero actual tool effects. Stats: zero credential reads/external attempts/browser launches, two owned GUI exits 0.

Root's full launcher read found `req.setTimeout(3000)` persisted on the socket after request-body completion. A local Node reproduction received body at 6ms, reset the client at 3009ms and attempted a valid response at 3507ms. This undermined the intended 90-second provider observation. Root preserved `launcher.mjs`/attempt03 unchanged, created `launcher-root-r1.mjs`, disarmed only the completed-body timeout, and delayed the existing positive fake-provider first headers by 3500ms. All six actual-executable controls passed (7.453s, exit 0), with zero keys/external/browser attempts and two GUI exits 0. A separate copy restoring only the old timeout failed the same delayed positive behavior (5.832s, exit 1, `POSITIVE_TWO_REQUESTS`), with owned GUI exit 0. Wrapper exit 0 merely records that expected negative failure; the negative itself remains failed. This is a measured harness repair, not a student correction round or new HTTP timeout in the product.

Root preservation receipt verifies all 120 source files, six artifacts and archive, 138 old R1 portable evidence entries, 31 frozen re-review inputs and twelve publication pins; student remains clean and main index empty. No paid/browser/native claim follows from local controls. Final `index.json` lists portable text evidence and this external report, excludes itself to avoid recursive hashing, and distinguishes local-only workspace files. Old evidence remains untouched.

## After approval and execution

Run once under the exact approved bounds. Preserve all actual outcomes and partial/failure receipts; verify cleanup, actual source/artifact identity, provider status/count/usage, live DOM/tool/result semantics, and final answer. A failed run is not silently retried or converted to a pass. Root accepts the chapter only after this remaining gate is genuinely resolved, then closes code workers and proceeds to the bounded author/text phase. No chapter tag or author edits are justified by this plan alone.
