# Chapter 8: operational proposal and bounded local preparation

**PLAN + LOCAL HARNESS PREPARATION ONLY.** The frozen R1 is independently code-review accepted per root, not overall chapter acceptance. No browser/CDP, audible/native speech, microphone, provider/key/discovery, rebuild, grade or long-test exercise occurred here. This packet is a request for later separate permission, not that permission.

## Frozen input and prepared commands

Candidate `b8f22ba5b3c9d1fc52d1a114810b90b77928684d`, tree `fa896caa404b87b08de61eb38ce079eb38138fbe`, remains clean and unchanged. Preserved GUI: `.ch08-local/bin/gui-ch08`, SHA-256 `6b513f3b10c5906bf66c6acf89e645903c95b0a397b711d93ae067b5b35ec0ae` (initial-candidate build stamp, not rebuilt). All source and existing artifact/package hashes are checked in the carry receipt.

The two small Python files reuse the exact checked-in `gui/test/fixture.py` Peer and the actual GUI binary. They do not recreate the model loop or browser. Final pins:

- `launcher.py`: `1765ead95d429a7f2b7e3e5c26a342aa571ae24fcf8ebd9dc53e3d51ca6aaa9f`
- `rehearse.py`: `2440eb69c89c72ced13e1821ff82b2809f7a74eb13c2331dfb56924102d17b42`
- `policy.json`: `eef19bf051df623c71d6227f1e375b48098177009e1d801f50f713e3533e6b27`

The later commands, **not executed now**, are:

```sh
cd /Users/bill/projects/ensemble-edition-3-student
# Local visual/native serving, ONLY after the corresponding human approval:
env -i PATH="$PATH" HOME=/Users/bill TMPDIR=/tmp LANG=C PYTHONDONTWRITEBYTECODE=1 \
  python3 -B .ch08-local/operational-prep/launcher.py serve-local \
  --policy /Users/bill/projects/ensemble-edition-3-student/.ch08-local/operational-prep/policy.json \
  --policy-sha256 eef19bf051df623c71d6227f1e375b48098177009e1d801f50f713e3533e6b27 \
  --approval /tmp/ch08-operational-local-approval.json \
  --run-dir /Users/bill/projects/ensemble-edition-3-student/.ch08-local/operational-runs/local-01

# One muted paid turn, ONLY after separate paid/key approval:
env -i PATH="$PATH" HOME=/Users/bill TMPDIR=/tmp LANG=C PYTHONDONTWRITEBYTECODE=1 \
  python3 -B .ch08-local/operational-prep/launcher.py serve-paid \
  --policy /Users/bill/projects/ensemble-edition-3-student/.ch08-local/operational-prep/policy.json \
  --policy-sha256 eef19bf051df623c71d6227f1e375b48098177009e1d801f50f713e3533e6b27 \
  --approval /tmp/ch08-operational-paid-approval.json \
  --run-dir /Users/bill/projects/ensemble-edition-3-student/.ch08-local/operational-runs/paid-01
```

The launcher refuses existing run directories. Only root may provide the approval after real human consent: `approved: true`, exact `policy_sha256` and `launcher_sha256`, a future timezone-aware `expires_utc`, a `consent_reference`, and appropriate `scopes` (`local-visual`, separately `native-audible`, or `paid-muted`). No approval file exists as part of this packet. A JSON flag is a technical interlock, not a substitute for consent or a same-user security boundary.

Serving is approval-gated. The harness has **no browser/CDP/native launcher at all**; unknown launcher flags are rejected. Later browser attachment is a separately permissioned root action, after loading the public GUI-debug skill. Root must first record its newly owned browser executable/version/argv/PID/profile and listener, then attach only to that process's loopback debugging endpoint. Use the run's new `browser-profile`, never a personal/shared profile. No live8088, focus automation, OS volume, preferred voice, voice enumeration or MIC permission.

## Minimal later exercise (stop at first genuine failure)

### A. Actual local GUI — muted

One loopback local session, at most **15 minutes including the native phase**, two pages, one owned GUI restart, and at most four local upstream requests (expected three). Launch the owned visible browser with `--mute-audio` and the reserved disposable profile. Harness-created `cr/settings.json` starts Read aloud off. This phase makes no real-provider requests.

1. Check the actual three panes, independently scrollable content, meaningful single root-Agent tree state, Agents/Settings keyboard tabs, and pointer plus keyboard divider controls. Inspect normal/narrow layout and light/dark/system token contrast; use browser media emulation for system preference, not OS changes. No invented fleet or schema-attack campaign.
2. Send `tools slow` once through the real composer. The inherited local peer streams and runs its two harmless file tools in the owned workspace. Observe thinking/text/tool artifacts, actions-pane placement and result readability. During that one turn exercise a typed-draft pause, then clear it or send one short hint; do not create an extra turn. Check independent scroll positions and visible busy/idle feedback.
3. Use actual settings controls and accepted snapshots across the two pages. Save a non-default theme/font, explicit TTS `false`, arbitrary speed `1.2345`, and budget `0` through the existing sparse settings protocol (budget has no visible control). Check false/zero survive snapshots. Restart **only the owned GUI once** by `kill -USR1 <recorded controller PID>`, then reconnect/reload. Verify snapshot, label, slider, theme/font and persistence, without expecting a new conversation-persistence contract. No automatic speech on replay/reconnect.
4. If needed, submit one local `Hello` turn after restart to create the short fixed text card **“First phrase. Final phrase.”** for the native phase. It uses one local request. Close the owned muted browser process before any audible phase; retain its exit and evidence. Do not switch an unrelated foreground application.

### B. Native speech — separately audible, no provider and no MIC

Only if the human explicitly agrees to audible output: reopen the owned profile visibly without `--mute-audio`, at most **120 seconds and four non-empty utterances**, no extra model prompts. Keep Read aloud off except the short cancellation check. Preserve the default voice and system volume; do not call `getVoices` or change OS/audio settings. Existing visibility wake utterances are empty and volume zero; note them separately rather than count them as intelligibility evidence.

Click the existing text card's Speak button once at each accepted saved speed **1.75**, **1.2345**, and **5**. Before each click capture the accepted snapshot, label and actual range value. Read the production `window.speechEvents`: it records requested setting, the actual native utterance `rate` getter, and native start/end/error events. Expected API rates are 1.75, 1.2345 and **3** (the source maps above 2 monotonically into 2–3). Use tolerance `0.00001`, not exact float equality. Record before asserting; a missing callback or wrong rate remains evidence and stops the exercise. This does not establish measured 5× wall-clock speech speed.

Fourth and final short utterance, at 1.75: enable Read aloud while idle, click Speak, wait for its native start, enter a draft, then disable Read aloud via the real control and wait for accepted `settings_changed`. Require cancellation and no speaking, while the draft keeps the local pause cause. Cancel then clears the draft and releases that remaining cause. Do not claim the Cancel button itself cancels speech; production separates those actions. If the phrase ends before the overlap can be observed, record inconclusive timing and stop rather than adding utterances. Close the owned browser by the window deadline, including on partial failure.

Native events/rates demonstrate API use, **not listening**. Separately record whether an actually present listener heard intelligible speech and a plausible relative-rate change; otherwise mark the audible/listening obligation unmet. No fabricated listening judgment.

### C. One useful paid model turn — separately approved, muted

Fresh owned directory/profile, Read aloud off and browser muted. One `anthropic/claude-sonnet-4-6` turn on the existing Anthropic route, no discovery matrix. The workspace already contains `exercise.txt` with `Chapter 8 boundary check.` Proposed prompt: **“Read exercise.txt using read_file once; quote its one line in one short sentence. Do not run commands or write files.”** Save budget 1 for this fresh instance. This is ordinary Agent guidance and an owned workspace, **not an OS security sandbox**.

At most **four upstream generation requests**, native **4096 max_tokens per request**, a **90-second observation window from first reservation**, zero automatic retries after any failed case. Setup has a separate 60-second allowance; if not ready to prompt, stop. The controller has a 150-second absolute paid lifetime and a 90-second active-observation guard. No fallback vendor, discovery, alternate model, added schema attack or tuned token maximum. The ceiling is expressed in requests/native output limits, not an invented USD estimate.

Only after this exact paid/key permission, the isolated proxy may read only `directClaudeAPIKey` from `~/.cr/settings.json`, lazily after pin/body validation and slot reservation. No current key read occurred. The real key stays in proxy memory and its required HTTPS header to `api.anthropic.com`; never GUI/browser env, argv, URL, logs, screenshots or portable evidence. The GUI receives only an explicit dummy key and proxy loopback base. No unrelated settings access; no claim of physical RAM zeroization. Stop at first transport/provider/contract failure, preserve partial receipts, do not silently retry.

## Forwarding and cleanup controls

Before owned GUI start and before each forwarding request: require exact candidate/tree, clean checkout, all 126 source pins, GUI binary, both script pins and caller-supplied policy hash. Verify all nine actually served assets against source hashes before browser readiness. The proxy validates only POST `/v1/messages`, exact model, native integer `max_tokens: 4096`, and streaming; no body rewriting. An irrevocable slot is reserved under a lock **before** key/upstream work. Request bodies are byte-preserved. The paid origin is fixed HTTPS `api.anthropic.com`; redirects/non-2xx terminate, never follow/retry. GUI inbound auth is discarded. The paid read timeout uses the remaining 90-second budget, not a short provider-work deadline. Controller expiry shuts active sockets before waiting for handlers.

GUI argv/cwd/PID, allowlisted environment names, dummy-key status, profile and proxy listener are recorded before readiness/attachment; served asset pins and GUI listener follow validation. Each process exit, partial result and forwarding counter is retained. SIGTERM stops only the owned controller; it closes sockets and its own GUI, waits, and uses a kill fallback only on that owned PID. No global process kill or cache cleanup. Browser cleanup remains root's responsibility because this launcher never creates one.

## Preparation actually executed

Full 6,418-byte task read/hash verified; retained public prerequisite reads qualified by unchanged exact public pins. No private reviewer inputs, parent memory, historical solutions or other workers read; no children.

- Attempt 1: 0.419s, exit 1 at served-asset verification. `/index.html` was a file-server redirect; fixture locator needed `/`. Wrong pin/model/token/path controls had already rejected without forwarding. The partial receipt, exact old scripts/policy, error and GUI exit remain. This was harness setup, not a product bug. The first diagnostic lacked the asset name; the repaired harness records every route/status/hash.
- Attempt 2: 5.704s, exit 0 after that locator repair.
- Before sealing, source inspection added socket abort before exceptional server teardown; this was preventive harness cleanup, not an observed production failure. The preceding passing revision remains retained.
- Final exact-pin rehearsal: **5.695s, exit 0**. Real frozen GUI completed its two-request local tool turn in 4.862s; first streamed response took **3.802s**, so a valid response beyond three seconds passed. All nine served assets matched. Two additional simultaneous forwards filled slots 3–4; the fifth was rejected. Wrong pin/model/token/path and a 302 redirect were rejected. Upstream body hashes matched unchanged bytes. No fabricated dispatch/completion/native evidence.
- Actual local branches used fixed `127.0.0.1` destinations only; key-read/external-request counters are zero. This plus source inspection is not a general OS network audit. The paid and affirmative serving/approval/browser/native paths remain **unexecuted**, hence the later gate is still required. Final sockets, threads and GUI exited; no job remains.

No production/source/binary changes, rebuilds, new source archive, full suites, grades or lifetime reruns. Prior source/evidence/artifacts remain unchanged. Final disk check was 4.7 GiB free; no cache deletion. This sealed UTF-8/non-NUL packet includes exact commands, hashes, partials and local-only raw-log indexes; no tar/binary/pycache/raw API/private material or credential. Do not append future runs to it. Root should inspect this packet and request the three explicit permissions above; preparation is not overall acceptance or an acceptance tag.
