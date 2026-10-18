# Prepared native exercise — NOT authorized, NOT executed, NOT execution-ready

## Status and exact identities

This is a separate replacement **proposal**, not a replacement of the original evidence. The older launcher and all failed/incomplete native observations remain byte-for-byte intact. Only `node --check` was run on the new launcher; no imports, preflight, GUI, Chrome, CDP, utterance, sound, or microphone operation was executed.

Frozen source `48e1d12d9c396b5988c0ce6eccdef67e581872ff`, tree `0f5f7fbe4273e91e1be3aa13d0715a99f671eaf9`, all **117** source hashes; GUI `39ebcd3e42a57d8972b6278fb8cc66021e38afd41afa2cd230b93499e6b13a81`. `native-policy.json` pins all 117 sources, the GUI binary, and all seven expected embedded asset files. A future run verifies source cleanliness/hashes and binary/source asset identities before any child process. Actual HTTP-served bytes cannot be read before a fresh GUI exists: the proposed ordering verifies its actual PID/listener first, then all served asset bytes, **before launching Chrome or loading a page**. This causality constraint is explicit; source bytes are not falsely called an observed served response.

Launcher SHA256: `2d582d6fb4caa5704e07a44a6d98c2799cd8e7b2bdb9f03e8880f84decfddcd6`. Policy SHA256: `2bc352bd307417d7392125b1b9c5f993ac81443ba2e705bae0f8e445981da9dc`.

## Known voice blocker — fail closed before launch

The exact production `speech.js` constructs `SpeechSynthesisUtterance(text)` and never sets `utterance.voice`. `window.speechEvents` records queue/native callbacks but not the resolved voice/service. An advertised default with `localService:true` alone cannot prove the voice actually selected for the production utterance. The launcher therefore has `verifyActualNativeVoiceBinding = null` and unconditionally stops before fixture/GUI/Chrome launch. There is **no CLI or environment override**. Root must resolve this evidentiary/policy blocker independently before any differently pinned launcher could execute. This preparation does not propose querying private OS settings, changing voice configuration, injecting a voice, monkey-patching native APIs, replacing callbacks, or relaxing the ban on remote service. No such mechanism is claimed to exist in the frozen app.

The page also has an existing visibility-change handler that may submit an empty muted native utterance. This is not removed/instrumented; no activation/focus/visibility manipulation is planned. This caveat remains for later root assessment and reinforces the pre-launch stop.

## Exact proposed command (do not execute)

```sh
/usr/bin/env -i PATH=/usr/bin:/bin:/usr/sbin:/sbin HOME=/var/empty TMPDIR=/tmp /opt/homebrew/bin/node /Users/bill/projects/ensemble/docs/edition-3-notes/ch07-student-harness-reconciliation-evidence/native-owned-probe-proposed.mjs --explicit-later-root-authorization --expected-script-sha256=2d582d6fb4caa5704e07a44a6d98c2799cd8e7b2bdb9f03e8880f84decfddcd6 --expected-policy-sha256=2bc352bd307417d7392125b1b9c5f993ac81443ba2e705bae0f8e445981da9dc
```

It is currently a blocked-preflight command even if the flag is supplied. Flags are assertions for a future reviewed launch, not permission. Any later approved launch must retain its console stdout/stderr in a newly created owned log directory; early script/policy pin refusals occur before child processes. The launcher itself retains a partial receipt when its known voice blocker is reached. A headed window would additionally require explicit later permission and `--headed-separately-authorized`. Audible output would additionally require separate human consent and `--audible-separately-authorized`; that flag also requires headed approval. No mode is authorized now. Human listening is neither implied by a callback nor asserted by script.

## Proposed ownership proof, before attachment

1. Fresh `/tmp/cadence-ch07-owned-native-*` scratch directory; fresh `mkdtemp` Chrome profile. No reuse/search of existing profiles/pages/CDP services; never port 8088.
2. Fixture binds a new loopback ephemeral listener in the exact launcher PID. `/bin/ps -p PID -o pid,ppid,command` and `/usr/sbin/lsof -nP -iTCP:PORT -sTCP:LISTEN -Fpcn` must identify that PID and only loopback listeners.
3. GUI gets a newly reserved/released port. Verify the exact spawned child PID, expected executable/argv and exclusive loopback listener **before any HTTP request**. A stolen port or daemonized/different PID fails; no substitution.
4. Chrome's endpoint file is only a location hint. Before CDP attachment, verify the exact spawned PID/argv includes the new profile and debugging flags, and `lsof` proves that same PID exclusively owns the actual CDP port on loopback. Recheck PID/command after allocation. Stop if unavailable/ambiguous, including a different child owning the listener.
5. Attach only to that verified browser, create a new isolated context and background target, and scope page/CDP/network events to its new session. No `Page.bringToFront`, target activation, focus, shared CDP or existing target enumeration. Creating a headed Chrome process could itself have OS windowing effects; no claim of headed non-disruption has been established.

These are proposed code controls, not ownership facts observed in this follow-up. `ps`/`lsof` paths/output parsing and Chrome flags remain runtime-unvalidated.

## Exact fixed local provider/network policy

One request only: `POST /v1/messages`, no query, at most 65,536 bytes, at most 2 seconds to receive. Exact recursive JSON equality to `native-expected-request.json`: model `claude-sonnet-4-6`, fixed system string, 4096 tokens, exact ten tool declarations, no `stream` key, one user text **“Show the two fixed local speech samples.”** The policy is derived from the owned historical fixture request, with only that user text and non-streaming representation changed; it has not been validated by a new live call. Any extra/mismatching/malformed/oversized/timed-out request immediately fails the whole proposed exercise. There is no provider retry, model/endpoint change, broad permissive peer, or fallback response.

Response: exactly two text blocks, with these fixed harmless samples:

- **Harbor lights shine softly.**
- **Quiet gardens rest.**

No tools, opaque data, remote references or additional model turns. GUI explicitly uses the owned loopback origin, fixed model/vendor, `--no-stream` and fresh workspace. Child environment is an explicit closed map, with only a literal `LOCAL_FIXTURE_NOT_A_CREDENTIAL` marker and fresh HOME/TMPDIR; no inherited secrets or real credentials. Inspection helpers disable global/system Git configuration. No API keys, browser settings stores, shared caches or credential files are read by the launcher.

Before Chrome birth, compare every served asset byte/hash with the frozen source pin. Then allow only one GET per pinned page-asset path plus optional favicon, and one exact owned `/ws` connection. Reject unexpected/duplicate traffic. Chrome background/sync/component/extension/DNS restrictions are explicit. Fetch interception and WebSocket observation are target-scoped, **not an OS firewall** and not proof against all browser/OS service network activity; they do not establish speech locality. The unresolved native service gate therefore blocks all launch now.

## Finite actual-page scenario proposed after all blockers/permissions

No production code, callbacks, queue, `speechSynthesis.speak/cancel`, voice fields or visibility behavior is replaced. Only existing page controls and the read-only `speechEvents` diagnostics are used. CDP `userGesture:true` is a synthetic gesture designation, not a human gesture/hearing claim.

1. Disable auto-read using the real existing checkbox before the single fixed prompt. Require connected/idle, exactly the fixed final assistant artifacts, and exactly one accepted provider request. Missing or merged artifacts cause failure, not adaptive rewriting.
2. Click both fixed artifact speaker buttons through the real page. Record both `queued` events. Require the first phrase's genuine current `onstart`, native `speaking:true`, and acknowledged global pause. The second phrase must not have started.
3. Immediately use the real read-aloud-off control, invoking production cancellation. Require recorded cancel, native `speaking:false` and `pending:false`, and pause release. Canceled utterance callbacks are deliberately generation-guarded by production; **do not invent a native `oncancel` or demand a fictitious canceled error callback**. Observe a 300 ms quiet interval with no leaked second start.
4. Queue the second phrase anew through its existing speaker control. Require its new queue admission, genuine `onstart`, and pause. Require its genuine **`onend`**, not merely start/error, followed by native idle and pause release. Stop immediately on a current native `error`, page exception, timeout, unexpected network, state mismatch, or unavailable condition. Do not repeat or change voices/mode/endpoints.
5. Native callback/state observations do not establish subjective intelligibility/listening quality. A separately consenting human would need to record their own hearing; no script assertion substitutes. No microphone or recognition activity is included.

Bounds: startup/listener/asset checks each finite under a **55-second overall exercise deadline**; first start 3 s; cancellation release 2 s; quiet interval 300 ms; recovery start 3 s; subsequent end/release 6 s. One attempt only, no fallback. The first native start must still be active when observed; a too-fast/missed first start is failure rather than claimed cancellation proof.

## Cleanup and receipts

On success, first failure, or deadline: preserve partial request/command/identity/native/state/network receipts. Use only the actual page cancel control, new owned context/browser, tracked spawned PIDs, and local fixture. Request shutdown, wait for actual exit, escalate only those still-live spawned children from TERM to KILL. Cleanup budget **10 seconds**, maximum planned total **65 seconds**. If an OS process still will not exit, retain its exact PID and an explicit cleanup-incomplete failure; do not claim successful cleanup or kill unrelated processes. No shared cache cleanup, old artifact deletion or profile reuse. The fresh scratch/profile remains retained for root-controlled later disposition. No cleanup behavior was executed here.

## Remaining limitations

Actual voice/service proof is absent and blocks launch. Proposed runtime controls have only syntax validation, not application/browser/native execution. No final-r4 native start/cancel/recovery/end, permission, audible quality or recognition success was observed. Old incomplete speech evidence and the earlier unmuted-probe authorization gap remain unchanged. A successful repaired public adapter does not satisfy these native/human gates or whole-code acceptance.
