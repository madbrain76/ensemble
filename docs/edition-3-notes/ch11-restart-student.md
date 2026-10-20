# Chapter 11 restart — student compiling HOLD

Implemented in `solutions/edition-3/main`; no commit, installation, deployment, or running-behavior claim. STOP at initial compiling candidate for root plus retained-reviewer inspection. No children or behavioral tests were launched.

## Ownership and design

- Added one Agent-owned MCP connection manager and shared contracts/root wiring. JSON-RPC correlation, initialization, listing, calling, cancellation, errors/content and closure share the same peer over real stdio children, split streams, dedicated sockets, and the existing GUI WebSocket tunnel. Child stderr stays diagnostic; provider credentials are filtered using the same helper as local commands.
- Skills parse declared servers, discover/normalize schemas, enforce ordinary remote-tool allowlists/collisions, and reconcile connection lifetimes. Remote schemas replace reserved execution-control properties; dispatch strips those controls before forwarding. External tools stop on pending unload; retained instructions and compiled tools keep the inherited pending-unload policy.
- Historical `ephemeral: round|turn` tools are hidden from ordinary definitions/calls. Collection uses job admission and cancellable completion observation, then guarded actor delivery through SYSTEM history events before `RequestSent`. Automatic job facts are retained without adding tool dialogue. This is not advertised as a current MCP-standard extension. Reverse `tools/builtin/list` and `tools/builtin/call` use enabled local-tool admission; `$/cancelRequest` cancels observation, not the admitted job.
- Added Agent methods `RegisterExternalTool`, `AttachMCPConnection`, `ServeMCPConnection`; foreground `--mcp-pipe` and `--mcp-socket ADDRESS` branches in CLI and GUI. No second Agent executor or provider client.
- GUI adapters reuse its single socket writer. Supplied optional `skills/ensemble` and `skills/gui-debug` use existing skill configuration. `gui_snapshot` and `tts_queue` inspect current DOM/speech state only on request, with explicit bounds/identifiers and recent callback records. GUI discovery does not block its WebSocket reader.

## Actual compilation

All checks sequential, with `GOMAXPROCS=2`, `GOFLAGS='-mod=readonly -p=1'`, `GOTOOLCHAIN=local`, `GOPROXY=off`, `GOSUMDB=off`, `GOMODCACHE=/Users/bill/go/pkg/mod`; owned HOME/TMPDIR/GOCACHE, existing Go/Delve PATH, standard provider keys and generic base URL unset. Shared module cache untouched. Disk checked before compilation; final available 3,023,276 KiB, never below 1 GiB.

Let `B=/Users/bill/projects/cadence-memories/children/e3-restart-ch11-student/build`, `I=/Users/bill/projects/cadence-memories/children/e3-restart-ch11-student/io`.

- Main: `go build -o $B/out/ch01 ./cmd/ch01`; `go test -c -vet=off -o $B/out/main-tests/ ./...` — exit 0, final log `$I/3667`.
- GUI: `go build -o $B/out/gui ./cmd`; `go test -c -vet=off -o $B/out/gui-tests/ ./...` — exit 0, `$I/3668`.
- JS: `node --check gui/web/gui/app.js` — exit 0, `$I/3668`.
- Each of `ch06`, `ch07`, `examples/ch04-consumer`: `go build ./...`, then `go test -c -vet=off -o $B/out/<module>-tests/ ./...` (slash replaced with hyphen) — exit 0, `$I/3668`. These single-package builds emitted executables in their module directories; moved only those newly generated files into `$B/out` afterward (`$I/3670`). Nothing added to Git.
- Interim CLI/main-test/GUI/JS compilation also passed (`$I/3646`, `$I/3647`). Contract reread rejected an initial `context/get` design in favor of the chapter's hidden-tool semantics. The subsequent compile failed (exit 1, `$I/3662`) on a missing import, stale initialization/parser names and schema locals; corrected before the final successful checks. Early inspection/edit-anchor misses were setup mistakes, not runtime results.

Focused test sources cover peer correlation/disconnect, schema/content, cancellation, and a real stdio helper child. They compiled; **none ran**.

## Remaining / feedback

Root and retained reviewer inspect actual source before further gates. Runtime evidence is still needed for stdio/socket/foreground EOF and shutdown, reverse admission/cancellation, schema/control allowlisting, round/turn freshness and stale-turn interruption, load/unload/reconnect (especially browser tunnel reopening), DOM/TTS accuracy, and inherited Chapter 9/10 replay/backup/EOF/GUI behavior. No browser, speech, provider, race, vet, grader, full-suite, or live-server operations performed. No reviewer feedback received yet; root's post-review production-correction budget has not been used. No further validation will launch automatically.

## Production correction Round 1 of 10 — HOLD

Root/reviewer feedback addressed within the existing owners (no architecture restart):

- Public `MCPTransport` now has the promised `Send(json.RawMessage) error`, `Recv() (json.RawMessage, error)`, and `Close() error`. Added `Agent.MCP().Connect(ctx,name,transport)`, `NewMCPTransport(reader,writer,closer)`, and `MCPPipe()`; existing adapters share the same codec. Explicit application connections without skills register a discovered catalogue atomically through the registry; skills-enabled connections still require a declaration. Outgoing IDs are integers. The GUI uses the public `jsonrpc`/`payload` envelope through its existing serialized writer.
- Automatic and ordinary remote tools both require the owning server's declaring skill grant. Unrelated same-name grants do not widen that server. Admitted automatic tools remain absent from ordinary definitions, `$TOOLS`, and ordinary dispatch. Collision checks retain local capabilities and include automatic-mode drift.
- Browser peer closure now clears page readiness. GUI skill loading can establish/await a new tunnel off the socket reader; pending external skills offer Reload. A real `/ws` regression proves late connection and load → unload → reload on one socket, with fresh initialization/discovery, immediate external refusal while pending, and the compiled skill still retained.
- `gui_snapshot` returns Markdown capped at 4096 UTF-8 bytes, with stable queryable data-attribute selectors, actual control/visibility state and previews, plus explicit elided-visible-preview-character and omitted/hidden-artifact counts. The DOM fixture exercises actual production functions, Unicode accounting, anonymous control actions, secret-value withholding, and page tunnel reset. It is not browser/native-speech evidence; production speech callback ownership is unchanged.
- Incoming aliases now obey initialization. Reverse ordinary calls still use real Agent/job admission. Peer-owned cancellable write admission releases a canceled queued call without disrupting the live writer; whole-message serialization, active-write closure, child reaping, and observation/job-lifetime separation remain.

### Actual focused checks and compilation

Same sequential offline environment and owned HOME/TMPDIR/GOCACHE as above; shared module cache untouched. No full suites, race, vet, grader, browser, native speech, external provider, children, deployment, commits, or snapshots. Fake peer/local HTTP fixtures only; the GUI case asserts zero provider requests.

- **Before FAIL:** `go test -vet=off ./internal/mcp -run '^(TestIncomingAliasesRequireHandshake|TestQueuedWriterCancellation|TestOutgoingIntegerIDs)$' -count=1` — exit 1, all three expected behavioral failures (`$I/3762`). Public compatibility `go test -c -vet=off -o <owned-output> .` — exit 1 on missing/wrong promised signatures (`$I/3767`), a distinct compile failure.
- **After PASS:** focused protocol run adding `TestCorrelationAndDisconnect`, `TestRealStdioChild`, and `TestCancelledCallDoesNotWaitForProvider` — exit 0 (`$I/3804`, before that batch's later test-source failure). Root tests `TestMCPPublicCompatibility`, `TestMCPExplicitOwningGrants`, `TestMCPIncomingAliasesUseJobs`, and `TestMCPApplicationConnectAndCollision`, followed by GUI `TestMCPBrowserReloadSameSocket` — all exit 0 (`$I/3807`). Commands used `go test -vet=off <package> -run '^(<listed names>)$' -count=1`.
- `node --test gui/mcp_dom_test.cjs` — 3 passed, exit 0 (`$I/3795`, before that batch's GUI fixture compile failure). `node --check gui/web/gui/app.js` — exit 0 (`$I/3809`). Interim interface/fixture compile mistakes and a wrong-working-directory edit were corrected, not counted as behavioral regressions (examples `$I/3795`, `$I/3798`, `$I/3804`, `$I/3806`). No unrelated product failure encountered.
- Final main `go build -o $B/round1/ch01 ./cmd/ch01`, GUI `go build -o $B/round1/gui-command ./cmd`, and each of ch06/ch07/examples/ch04-consumer `go build -o $B/round1/<module>-bin/ ./...`; every module then `go test -c -vet=off -o $B/round1/<module>-tests/ ./...` — all exit 0 (`$I/3809`). Artifacts stayed outside the repository. Disk checked before every build/compile; final available 2,649,132 KiB, above the 1 GiB floor.

**HOLD for root and the same retained reviewer.** Scoped evidence above does not replace full validation of automatic round/turn freshness, interruption and shutdown interleavings, foreground EOF behavior, inherited Chapter 9/10 recovery/replay/GUI behavior, or real-browser/native-voice observation. No further validation will launch automatically. Round 1 of the root-owned ten-round correction budget is submitted.

## Released full validation — STOP at Chapter 11 failure

Source/fixtures remained frozen after root + same-reviewer acceptance of R1 **for validation only**. Used `env -i` with the approved offline Go variables, existing Go/Delve PATH, owned HOME/TMPDIR/GOCACHE, and no inherited provider keys or generic base URL. Ran gates sequentially, checking the 1 GiB disk floor; no cache cleanup, fixes, or operational work.

- Main `go test -race ./...` — exit 0 (`$I/3839`); main `go vet ./...` — exit 0 (`$I/3840`).
- GUI `go test -race ./...`, then `go vet ./...` — both exit 0 (`$I/3840`).
- Separate ch06, ch07, and examples/ch04-consumer: each `go test ./...`, then `go vet ./...` — all exit 0 (`$I/3841`). Ch07 and consumer reported no test files: their test commands establish compilation only.
- Built once: from main, `go build -o /tmp/e3-ch11-final/frozen-cli ./cmd/ch01` — exit 0. SHA-256 (`shasum -a 256`): **`7e949c75dd8e46d60d6072502c8a63c3338ec06304c945f309e3bb1f34df58ac`** (`$I/3843`). No rebuild afterward.
- Public map shape confirmed from `docs/edition-3-student.md` and Chapter 5, without inspecting grader code. Input `$B/media-models-validation.json`: `{"anthropic":"claude-sonnet-4-6","openai":"gpt-5.4","gemini":"gemini-3.8-flash"}`.
- From book repository root, ran exactly: `go run ./cmd/grade-edition3 -ch 11 -json -core /Users/bill/projects/ensemble/solutions/edition-3/main -consumer /Users/bill/projects/ensemble/solutions/edition-3/main/examples/ch04-consumer -workflow /Users/bill/projects/ensemble/solutions/edition-3/main/ch06 -stream-exercise /Users/bill/projects/ensemble/solutions/edition-3/main/ch07 -stream-model claude-sonnet-4-6 -gui /Users/bill/projects/ensemble/solutions/edition-3/main/gui -gui-model claude-sonnet-4-6 -media-models /Users/bill/projects/cadence-memories/children/e3-restart-ch11-student/build/media-models-validation.json /tmp/e3-ch11-final/frozen-cli`.

**Chapter 11 composite FAILED:** `passed:false`, score 25/100, process exit 1; full output retained at `$I/3845`. Reported failures: protocol-pipe handshake command-output deadline; discovered description/schema changed (five details); failed automatic output reached requests 1–5 and incorrect once-per-turn context for requests 2–4; reverse exercise prevented by the handshake deadline; GUI HTTP startup failed with `MCP browser: MCP transport unavailable: websocket`. MCP tool call and Chapter 10 parity checks reported pass. These are report facts, not a source diagnosis or a numerical acceptance claim.

Stopped immediately after this failed gate. **Chapter 1 was not invoked.** No automatic fixes, fixture changes, later gates, or operations proposal. Final disk available 2,352,996 KiB, above the floor. Frozen CLI remains available with the hash above. Root direction requested on the retained failure output before any further work.

## Remaining grading after root fixture correction — FINAL/HOLD

Root authorized one fresh Chapter 11 run after its reported public-contract fixture corrections and retained-reviewer acceptance. No production Round 2, student source/fixture edits, CLI rebuild, local-suite repeats, private grader inspection, or operational work. Prior full-suite passes remain valid; the original `$I/3845` result remains **25/100, passed:false, exit 1**, not overwritten or reclassified as a pass.

- Verified existing `/tmp/e3-ch11-final/frozen-cli` SHA-256 **`7e949c75dd8e46d60d6072502c8a63c3338ec06304c945f309e3bb1f34df58ac`** before grading. Reused the same `$B/media-models-validation.json` input.
- From book repository root, repeated the exact Chapter 11 command recorded above: same public adapter, five separate module paths, `claude-sonnet-4-6` stream/GUI models, existing media map, and frozen executable. **100/100, passed:true, process exit 0**; all seven checks passed. Fresh complete output/exit/hash receipt: `$I/3966`.
- Only after that pass, ran the **first** Chapter 1 grade from book root: `go run ./cmd/grade-edition3 -ch 1 -json /tmp/e3-ch11-final/frozen-cli`. **100/100, passed:true, process exit 0** (`$I/3967`). Report explicitly confirms complete answer and usage after **1m2s response-header delay** and **1m2s response-body delay**. Observed command elapsed **124 seconds**, with no execution deadline or cancellation. Final CLI hash is identical.
- Both commands used the same clean `env -i`, approved offline toolchain/cache variables, existing Go/Delve PATH, and owned HOME/TMPDIR/GOCACHE. No inherited keys or generic base URL. Sequential disk checks stayed above 1 GiB; final available **2,192,400 KiB**. New outputs preserve, rather than replace, the original failed grade.

### Operations proposal only — not authorized or executed

1. After explicit root release, use root-designated isolated foreground instances and an approved connected browser; do not restart or modify the live shared GUI. Observe a real stdio MCP child and the existing-GUI tunnel through initialization, declared-skill load, call, unload, same-page reload, and close. Limit reverse local-tool effects to an owned temporary file.
2. On demand in that browser, inspect `gui_snapshot`, use a listed selector for a benign reversible action/input, then inspect again to verify actual DOM state and identity. Separately, with explicit native-speech permission, inspect `tts_queue` against real queue/start/end or cancellation callbacks. Do not substitute fixture results for native observation.
3. Keep any paid-provider exercise, credential access, process restart, or publication outside this proposal's execution authority; root must separately specify and release it if required. Record actual observations and limitations concisely in this note, then return to review rather than treating grading as chapter acceptance.

**FINAL/HOLD:** grading is complete on the one frozen executable. No operations scripts, rehearsals, browsers, native speech, provider/key discovery, children, commits, tags, or snapshots launched. Await root and the retained reviewer before any operational release.

## Authorized operations pass — STOP at launcher setup failure

Prepared one owned controller in `/tmp/e3-ch11-operations.rnujz1` (macOS resolves it to `/private/tmp/...`): `control.py`, `browser.mjs`, `peer.py`. Read only the three explicitly permitted old scaffold files, not adjacent captures/cases/secrets; used current source/public Chapter 11 for the intended workflow. Static Python parsing and `node --check` passed. No production/source/fixture edit, binary rebuild, or additional grade/suite.

Pinned identities at launch (`$I/3993`):
- Source digest `c1f70e0172e9b77e71c92ff52625da3096ed4b38f266b13d0245d70a90f0f5a7` (controller's sorted relative-path/NUL/content/NUL digest over current Go/module/JS/HTML/CSS/CJS source plus the two current skill files).
- Frozen CLI SHA-256 unchanged: `7e949c75dd8e46d60d6072502c8a63c3338ec06304c945f309e3bb1f34df58ac`.
- Reused GUI from held R1 final build `$I/3809`, build metadata `ensemble/gui/cmd`, predecessor revision with modified working tree: SHA-256 `ebe33db890a0ba628d89745ef32837d2a66e83cdb18fb8be761156a0cb87a60d`. No GUI rebuild.

Actual invocation: clean approved `env -i ... /opt/homebrew/opt/python@3.13/libexec/bin/python3 /tmp/e3-ch11-operations.rnujz1/control.py`, cwd that owned directory — **exit 1** (`$I/3993`). Its first GUI child (PID/PGID 81128) was launched with `gui-command -listen 127.0.0.1:61944 -settings <owned>/local/settings.json -save <owned>/local/save.json` and **exited 2** immediately. Full stderr: `<owned>/local/gui.stderr`: **`flag provided but not defined: -listen`**. The actual help lists `-port int` and `-gui-dir string` (default `web/gui`). This was an operations-launcher assumption error despite the preparatory source reads, **not evidence of a GUI/MCP product defect**. No flag correction or automatic retry followed.

Controller immediately stopped and closed its owned fake-provider listener; ports 61943/61944/61945 were observed closed. GUI process was reaped. Chrome was **not launched or attached**; only an empty owned profile directory had been prepared. No shared browser/8088, focus, volume, MIC, native speech, live discovery, or credential access. Actual counts: **0 fake-provider requests; 0 native utterances; 0 models GET; 0 paid POST slots/completions; 0 key reads**. No local browser/MCP/native assertions were reached, and the live gate stayed closed. Prepared local settings/skills and raw startup/controller logs remain only in the owned temporary directory. Free disk at stop: **2,180,767,744 bytes**, above the floor.

Partial ordinary receipts: `<owned>/operations.jsonl`, `<owned>/local/gui.stdout`, `<owned>/local/gui.stderr`, plus `$I/3993`; planned browser/live receipts do not exist because those stages did not run. **FINAL/HOLD pending root direction** on this setup failure. Prior suite/grading passes remain unchanged; operational acceptance is not claimed.

## Controller-only repair v2 — STATIC HOLD, not production Round 2

Per root's bounded authorization, copied **only** `control.py`, `browser.mjs`, and unchanged `peer.py` into `/tmp/e3-ch11-operations-v2.qjUn0M`. No raw captures/settings/profile were copied. Original `/tmp/e3-ch11-operations.rnujz1` scripts were hash-checked unchanged; its failed zero-request attempt and artifacts remain intact.

V2 changes are limited to the requested launch/lifetime/settings/pause/TLS/cleanup repairs:
- GUI now uses `-port <int> -gui-dir /Users/bill/projects/ensemble/solutions/edition-3/main/gui/web/gui`, with `stdin=PIPE` retained throughout the exercise. On successful browser completion it closes stdin, waits for exit 0, and reads the resulting save; failure cleanup remains owned-process-only.
- Settings are the direct field object. Chrome adds the approved `--use-mock-keychain`, with no new variants/cases.
- Fake request 2 issues `gui_click` first, `gui_input` last. The driver observes the actual Paused display, then explicitly Continues after the first native cancellation; it does not wait for a later dispatch behind the typing veto. Native/other budgets remain unchanged.
- TLS explicitly uses `/etc/ssl/cert.pem`, SHA-256 `9dae8d76e55cb08991f2b672d58999ea15560d910759c16b544f843bdffbb994`, checked before controller networking/key access and again at pins/upstream connections. Hostname verification and `CERT_REQUIRED` remain mandatory; no fallback.
- Only the new controller's owned profile is eligible for removal, after its Chrome has exited and all owned local listeners are closed. The original failed-attempt profile is untouched. Local failure is published before cleanup so blocked fixture work can stop.

Static receipts: `$I/4008` versioned exact-anchor edits/original hash checks — exit 0. `$I/4009`: `diff -u` original/v2 controller and driver — exit 1 each (expected differences); `cmp` peer — exit 0; clean-environment Python `compile`/`ast.parse` (no controller import/execution) — exit 0; `node --check /tmp/e3-ch11-operations-v2.qjUn0M/browser.mjs` — exit 0. Current source verifies the literal Paused/Continue controls. Source and both binary digests remain the held values recorded above.

The actual controller interpreter resolves to `/opt/homebrew/opt/python@3.13/bin/python3.13`. Static SSL statistics in the clean environment: default context **191 CAs**, explicitly pinned system bundle **128 CAs**; explicit context `check_hostname=True`, `verify_mode=CERT_REQUIRED`. This differs from root's default-root observation for `/opt/homebrew/bin/python3`; v2 does not depend on either interpreter's default roots. Certificate checks/statistics made no HTTP connection or key read.

V2 SHA-256:
- `control.py`: `24ae6ad403be903f6aa07c11ad1a0eac73636d25610742f8b8b25a878db9f1c9`
- `browser.mjs`: `af86aaf115c4f6d4b90f2da04cacf712a0b4c4c58315b5095140dd50129b34c9`
- `peer.py`: `76138faae13b3188c7f4d075238952f7f90c9ff98fecd77cc11295c8bee53bf0`

**FINAL/HOLD for root inspection. V2 has NOT been executed.** No browser/GUI/native launch, networking/GET/POST, key read, new operation/rehearsal, production edit, rebuild, grade, suite, or cleanup occurred in this repair turn. Eventual budgets (fake ≤6; native ≤2/90s; GET ≤1; paid POST ≤4; key reads ≤1) remain unused. No operational pass is claimed.

## Exact v2 execution — STOP at native second-START observation timeout

Root released the exact v2 script hashes recorded above once. Verified them and the 1 GiB floor, then ran clean approved `env -i ... /opt/homebrew/opt/python@3.13/libexec/bin/python3 /tmp/e3-ch11-operations-v2.qjUn0M/control.py`, cwd that directory — **exit 1**, full output `$I/4016`. No script/production edit, rebuild, suite, grade, retry, or additional diagnostic execution. Source/CLI/GUI launch pins matched the held digests above. Preserve the TLS qualification: actual Python 3.13 has 191 default roots; root's earlier zero-root observation was another interpreter, **not this failure's cause**. V2 used the explicitly pinned 128-root bundle, hostname verification and `CERT_REQUIRED`; no live connection was reached.

Actual local ownership/effects:
- GUI PID/PGID 81236, loopback 61963, correct `-port`/absolute `-gui-dir`, owned settings/save/cwd and held stdin PIPE. Command/listener verified before attach.
- Headed, muted Chrome PID/PGID 81240, owned private profile and CDP 61964, approved mock-keychain flag. PID/command/profile/listener and Chrome/CDP identity verified before the Node driver (PID 81259) attached. Page receipt: client `2`, session `1792528711193999000`, agent `agent-1`. No shared 8088/browser, MIC, OS focus or volume manipulation. Browser background writes remain possible; no host-wide-zero-effects or hearing/locality claim.
- **4 fake-provider requests**, within the six-request cap. Model-issued MCP calls reached the real Agent/jobs/GUI WebSocket/DOM: requested `gui_snapshot` twice, `tts_queue` twice, `gui_click` and `gui_input` once each. Skill was loaded after page connection; no tool receipt existed before explicit requests. Typing pause was observed and explicitly Continued.
- Actual snapshot receipts passed independent byte/count/selector checks: **3947 UTF-8 bytes**, elided 0, omitted 27/58 (27 hidden); then **4059 bytes**, elided 4093 visible-preview characters, omitted 67/98 (27 hidden). All listed selectors were queryable; the second snapshot showed the model-edited draft with the stable textarea selector. The intended final snapshot after the normal long assistant text was **not reached**.
- The first short native phrase produced real START/speaking and fresh cancellation observations. Requested `tts_queue` recorded production `queued`/`speaking`/`cancelled` events and empty current/pending/native state afterward. Two UI Speak-action attempts were made. The driver then stopped at **`deadline: native second START`** (30-second observation bound); no second START/END or terminal native-state dump was saved, so an exact second-utterance/END outcome is **not claimed**. No callbacks were faked/substituted. This remains an unclassified operations observation failure, not an established product or OS defect.

Stop/cleanup and unreached gates:
- Browser driver exit **1**; owned Chrome and GUI exited **0** during failure cleanup. This is not proof of the successful-path controlled-EOF/save gate, which was not reached.
- Owned fake-provider 61962, GUI 61963, CDP 61964 all observed closed. Only v2's owned profile was removed after Chrome exit and closed-listener checks; the original failed-attempt artifacts/profile were preserved. Total local attempt about **35.55 seconds**, below ten minutes; free disk at stop **2,100,527,104 bytes**.
- Same-page unload/reload, final snapshot/queue assertions, and LOCAL_PASS were not reached. **Live remained closed: 0 models GET, 0 paid POST reservations/completions, 0 key reads.** No paid retry/fallback or live CLI launch.

Partial ordinary receipts remain in `/tmp/e3-ch11-operations-v2.qjUn0M`: `operations.jsonl`, `local/browser.jsonl` (precise completed checks and STOP), `local/browser.stdout`, empty `local/browser.stderr`, `local/mcp-results.jsonl`, `local/api/{1..4}.*`, and owned GUI/Chrome logs. Only saved receipts/current driver's catch block were read after stop; no further executable diagnosis was run. **FINAL/HOLD for root and retained reviewer. No operational acceptance claimed.**

## Remaining-only continuation — LOCAL PASS; live pre-forward guard STOP

Root/reviewer classification carried forward: v2 had only the first delivered Speak/request/claim. Its second coordinate-click return did **not** prove delivery; no second native-production defect or rAF/OS cause was established. The valid first RPC/DOM observations and real START/fresh cancellation are prior evidence, not the new browser's history.

Authorized remaining-only scripts: `/tmp/e3-ch11-operations-remaining.sp7uhA/{control.py,browser.mjs,peer.py}`. No production/binary changes or rebuilds/gates. Compared against v2; copied only scripts and the explicitly authorized save. Save SHA `20833e0824b83c32dd3cc762c799950f14d85f3f83b9b25e6205aab20c9c5621` contained the phrase and 1628-codepoint long assistant answer at sequence 38; copied byte-equal, not reconstructed. Current GUI successfully validated that copy with `backup=false`, then rendered both required parts.

Preparation receipts: `$I/4061` adaptation; `$I/4062` full version diffs (exit 1 expected), Python AST/compile and JS syntax (exit 0), original v2/source-binary-input pins; `$I/4065` local IDs numbered cumulatively 5/6 to avoid reuse in restored history, final diff and syntax (exit 0 overall). Live proxy/TLS/EOF suffix and peer were byte-compared **unchanged** from v2. Final script SHA-256:
- controller `88e0844687850ffe0717769f99c7239f1a10deb1789209710462cf68a7a67f6a`
- browser `dd82f6832b16abc251bc06c4b1ad2971ce30eeda79fd38828bc785ab9e1a04dd`
- peer `76138faae13b3188c7f4d075238952f7f90c9ff98fecd77cc11295c8bee53bf0`

Executed **once**: clean approved `env -i ... /opt/homebrew/opt/python@3.13/libexec/bin/python3 /tmp/e3-ch11-operations-remaining.sp7uhA/control.py`, cwd same — **exit 1**, full output `$I/4067`. Held source `c1f70e…f5a7`, CLI `7e949c…58ac`, GUI `ebe33d…a60d` matched (full hashes above). No automatic retry or executable diagnosis after stop.

Completed remaining local observations (~4.82 seconds):
- Owned GUI PID 81457/port 62086, Chrome 81461/CDP 62087, Node 81479; command/profile/listener/CDP identities verified before attach. Headed, muted, private profile, mock-keychain; no shared browser/8088, MIC, OS focus-volume action, voice discovery, hearing/locality claim, or host-wide-zero-effects claim.
- New page client `2`, session `1792529414990940000`, agent `agent-1`: declared skill loaded, unloaded with external peer closure/instructions retained, then reloaded on the **same page/socket**. No tool ran before explicit request.
- Exactly **one** new short phrase via the actual existing Speak button's programmatic `.click()` with CDP user gesture—not a physical mouse claim. Saved outgoing Speak, speech_requested, claim_speech, and owner-matched speech_claimed witnesses. Passive observations recorded `submitted`, actual **START**, actual **END**, with existing production callbacks still functions. Completion/cleared queue and pause observed **3130 ms from delivery**; no error/cancel or second native call in this session. Evidence is saved even on failure by this driver; no callbacks/results substituted.
- **Two new fake-provider requests, six cumulative**. Real Agent/jobs/reloaded GUIWS requested one `gui_snapshot` and one `tts_queue`, then normal final answer. Snapshot exactly **4096 UTF-8 bytes**, independently checked **8345 elided visible-preview characters**, **95/122 omitted artifacts (26 hidden)**, 27 unique queryable listed selectors. It used the saved long DOM, not newly manufactured content. Requested queue reported this session's actual queued/speaking/chunk-ended/done callbacks and no current/pending/delayed/request/native work; Agent pause returned to idle. Prior cancellation was not asserted as current history.
- GUI stdin then received controlled EOF: exit **0**, saved sequence **54**. Node and Chrome exit **0**. Owned fake/GUI/CDP listeners 62085/62086/62087 closed; only this continuation's owned profile removed afterward. LOCAL_PASS gated live work only after cleanup.

Live attempt then stopped **before any paid reservation/forwarding**:
- One authorized memory-only key read; one fresh verified-TLS models GET, HTTP **200**, confirmed `claude-sonnet-4-6`. Explicit pinned `/etc/ssl/cert.pem`/128 roots, hostname verification and `CERT_REQUIRED` unchanged. Actual Python 3.13's 191 default roots remain irrelevant to that explicit choice; root's prior zero-root observation was another interpreter, not a failure cause.
- Frozen CLI PID 81492 launched once with owned cwd/save, dummy/minimal environment. Real stdio peer receipts show initialize/initialized/list, `environment_snapshot` **once** (`Snapshot number 1`), `session_brief` **once**. No `read_marker` call occurred.
- STOP reason: **`automatic tool menu exposure`**. The unchanged proxy predicate requires `read_marker` present **and** automatic names absent; the label alone does not identify which clause failed. Guard refused the first CLI request before its numbered request receipt or upstream slot; **0 paid POST reservations, 0 completed paid POSTs**. CLI received proxy HTTP 502 and exited **1**. Marker read/paired result/final answer/live EOF→usage/save proof remain **uncompleted**, not accepted. No production/OS/TLS diagnosis asserted from the guard label.
- Owned live proxy listener 62132 closed and CLI reaped. Entire continuation ~6.22 seconds; final free disk **2,030,628,864 bytes**. **Cumulative fake 6; models GET 1; key reads 1; paid POST 0.** No key retry/read, fallback, redirect, body rewrite, new native action, or further operation followed.

Ordinary receipts: continuation `operations.jsonl`; `local/browser.jsonl`, `native-observed.jsonl`, `delivered-state.json`, `native-completed-state.json`, `final-state.json`, `gui-wire.jsonl`, `mcp-results.jsonl`, `api/{1,2}.*`, GUI EOF save/logs; `live/cli.stderr` (502), discovery receipt and `live/work/peer.jsonl`. Raw artifacts stay local/excluded Git. Only saved receipts/current guard/peer source were read after stop. **FINAL/HOLD for root + SAME reviewer: local missing gates complete; live exercise failed at pre-forward guard and remains incomplete.**

## Production correction Round 2/10 — automatic-tool projection; compiling HOLD

Root + SAME reviewer corrected/withdrew the initial flags/exact-name-guard diagnosis: the executed peer's ephemeral fields and the existing live guard were already correct. The supported Anthropic projection was the defect. No drift/unauthorized-edit claim, operation-script repair, or paid retry was made.

Production change is confined to `main/internal/engine/engine.go`, `Engine.Render`: build the supported-model catalogue and enabling-call reference candidates from ordinary (`Ephemeral == ""`) tools only. The filter precedes **both** projections. Full `Snapshot.Tools`, admission/grants, automatic gathering, dormant ordinary catalogue entries, genuine enabling-event checks, fallback rendering, and canonical history are unchanged. No new owner/framework/provider layer.

Focused regression `main/ch11_projection_test.go` uses real Agents, loopback MCP startup discovery, real automatic jobs, and local Anthropic HTTP. Supported 4.6 and conservative-fallback 4 subcases each perform two actual requests, including a genuine model-issued/paired `load_skill` job. It checks initial round/turn context, disabled automatic entries, dynamic automatic activation, hidden declarations/references, unchanged ordinary dormant/schema/defer behavior, actual ordinary enabling reference/provenance, canonical handler status/pairs, missing-pair fallback, and nonmutating replay. Automatic calls remain round=2, turn=1, dynamic round=1, newly enabled application round=1 per subcase. The existing Go test peer responder in `ch11_round1_test.go` was factored for reuse over the loopback startup transport; operation peer/controller scripts were not changed.

Actual commands/results (clean `env -i`, approved offline Go/cache/owned HOME/TMP, `GOMAXPROCS=2`, readonly `-p=1`, no inherited keys/base URL):
- **Before production edit**, `go test -vet=off -run '^TestChapter11AutomaticAnthropicProjection$' -v .` — **exit 1**, `$I/4144`. Newly compiled regression on the held production code showed initial, dormant and dynamic automatic declarations plus dynamic automatic `tool_reference` disclosure on supported 4.6. Fallback subcase passed. This is a behavioral before-failure, not a compile/setup failure.
- After the single projection fix, `go test -vet=off -run '^(TestChapter11AutomaticAnthropicProjection|TestChapter09EnablingEventAndRendering|TestChapter09Round1DirectLoadExposure|TestMCPExplicitOwningGrants)$' -v .` — **exit 0**. Then `go test -vet=off -run '^(TestSkillToolReferencesAreContentNotInstructions|TestChapter09Round2OrdinaryResults)$' -v ./internal/engine` — **exit 0**. Both in `$I/4149`; no full suite/race/vet.
- Sequential builds/test-source compilation — **all exit 0**, `$I/4150`: main `go build -o …/main-cli ./cmd/ch01`; GUI `go build -o …/gui-command ./cmd`; ch06/ch07/consumer `go build -o …/<module>-build/ ./...`; each module `go test -vet=off -c -o …/<module>-tests/ ./...`. ch07 and consumer report no test files (compilation only). All outputs are scratch artifacts under **`/tmp/e3-ch11-r2-compile.2XSZ2W`**, not frozen-pin replacements. Disk checked throughout; ended **1,651,316 KiB available**, above 1 GiB; no cleanup.

Setup-only missteps, separate from the regression: early inspection guessed absent `anthropic*.go`, `ch09*`/`context.go`, and a main-relative book path (e.g. `$I/4132`, `$I/4134`); actual owners were then read. One exact-edit anchor was refused with no file change; the actual small `anthropicDisclosure` block was reread before the successful edit. No compile failure or unrelated behavioral failure occurred.

Old frozen hashes rechecked in `$I/4150`: CLI **`7e949c75dd8e46d60d6072502c8a63c3338ec06304c945f309e3bb1f34df58ac`**; GUI **`ebe33db890a0ba628d89745ef32837d2a66e83cdb18fb8be761156a0cb87a60d`**. They and all operations receipts remain preserved; they are **not** binaries of this revised source. No new aligned frozen CLI, grades, browser/native action, GET, upstream POST, key read, script rewrite, children, commit/tag, or manuscript/guidance change. Completed local RPC/native evidence carries under root's qualification; a general hearing report does not identify a particular voice/run. Live usage remains **POST 0 / GET 1 / key read 1**, with no retry authority.

**HOLD** after focused fix + five-module compilation. Root + SAME reviewer to inspect this delta before any aligned full validation/new CLI/live-only release. Acceptance remains pending.

## R2 aligned validation PASS — live-only controller PREPARED, NOT EXECUTED

Following root + SAME reviewer's R2 acceptance for validation, source/tests stayed held. Sequential clean `env -i` with approved PATH/owned HOME/TMP/GOCACHE, `GOMAXPROCS=2`, `GOFLAGS='-mod=readonly -p=1'`, `GOTOOLCHAIN=local`, `GOPROXY=off`, `GOSUMDB=off`, existing readonly module cache; no inherited real keys/generic base URL:
- cwd `/Users/bill/projects/ensemble/solutions/edition-3/main`: `go test -race ./...` — **exit 0**, `$I/4178`; then `go vet ./...` — **exit 0**, `$I/4179`.
- **One new frozen build only:** `go build -o /tmp/e3-ch11-r2-final/frozen-cli ./cmd/ch01` — **exit 0**, `$I/4179`. SHA-256 **`a8af82a8cdc1be9e92b335777d640713441062e4a69fa219385fc6f306bbb11c`**. Old CLI `7e949c…58ac` and GUI `ebe33d…a60d` rechecked unchanged. No GUI/module gates or local browser/native exercise repeated; those qualified results carry as authorized.
- From book repository root, exactly:
  `go run ./cmd/grade-edition3 -ch 11 -json -core /Users/bill/projects/ensemble/solutions/edition-3/main -consumer /Users/bill/projects/ensemble/solutions/edition-3/main/examples/ch04-consumer -workflow /Users/bill/projects/ensemble/solutions/edition-3/main/ch06 -stream-exercise /Users/bill/projects/ensemble/solutions/edition-3/main/ch07 -stream-model claude-sonnet-4-6 -gui /Users/bill/projects/ensemble/solutions/edition-3/main/gui -gui-model claude-sonnet-4-6 -media-models /Users/bill/projects/cadence-memories/children/e3-restart-ch11-student/build/media-models-validation.json /tmp/e3-ch11-r2-final/frozen-cli`
  — **process exit 0 AND JSON-parsed `passed: true`**, 100/100, `$I/4180`.
- Only after that pass: `go run ./cmd/grade-edition3 -ch 1 -json /tmp/e3-ch11-r2-final/frozen-cli` — **process exit 0 AND JSON-parsed `passed: true`**, 100/100, **124 seconds elapsed**, `$I/4182`. Genuine 62-second headers and 62-second body checks each reported complete answer/usage; ordinary cumulative usage matched 511 input/65 output. No execution deadline/cancel or rebuild between grades; same new CLI hash rechecked afterward.
- Versioned reports: `/tmp/e3-ch11-r2-final/{grade-ch11.json,grade-ch11.stderr,grade-ch01.json,grade-ch01.stderr}`. All older reports/failures retained. No grader/private source inspection or edits. No unexpected gate failure. Disk remained above 1 GiB; final static-preparation observation **1,568,124 KiB available**. No cleanup.

### Prepared only for future explicit live-only release

Directory `/tmp/e3-ch11-r2-live-only.7ui5ezji/` contains only prepared `control.py` and unchanged `peer.py` (no browser driver). Based on `remaining.sp7uhA`: all local/fake/GUI/Chrome/native execution and all models GET code removed; root-carried local results are not repeated. Updated hard source pin **`12bd5c3c6c388c525c72e3464f9bce6f21e404241e4efe3e0a0009241c6e1487`** and new CLI pin above. Reuses only the existing successful discovery `/tmp/e3-ch11-operations-remaining.sp7uhA/live/models.json`, SHA **`02a9188989a823990a206aff15d9764786d298d249cd7b4f355ab834361fcd9e`**, with the requested model present; this was a local receipt read, **not a GET**.

Controller retains the correct POST guard, atomic pre-upstream four-total-POST cap, fixed native 1024/model, untouched-body/SSE forwarding, explicit pinned CA/hostname/CERT_REQUIRED, 60s setup/120s active/180s absolute bounds, one memory-only key read **only after future execution release**, actual peer/count/result checks, complete answer→EOF→exit/usage/save ordering, owned cleanup and stop/no-retry behavior. Guard/upstream/process/record functions were text-compared unchanged; peer byte-equal. No new harness case.

Preparation/syntax only `$I/4186` **exit 0**: AST/compile without importing/executing either file; existing discovery/model and source digest read. Exact version diff `$I/4187`: `diff -u /tmp/e3-ch11-operations-remaining.sp7uhA/control.py /tmp/e3-ch11-r2-live-only.7ui5ezji/control.py` — **exit 1 expected differences**; peer `cmp` **exit 0**, overall command **exit 0**.
- Controller SHA **`cb48d48bb098bfbd1ab7249f229b330d0361824c12141665ede1938449504576`**.
- Peer SHA **`76138faae13b3188c7f4d075238952f7f90c9ff98fecd77cc11295c8bee53bf0`**.

**Proposed exact command—printed for review, NOT RUN; requires a future explicit root release** (cwd `/tmp/e3-ch11-r2-live-only.7ui5ezji`):
```sh
env -i PATH=/Users/bill/go/bin:/usr/local/go/bin:/opt/homebrew/bin:/usr/bin:/bin:/usr/sbin:/sbin HOME=/Users/bill/projects/cadence-memories/children/e3-restart-ch11-student/build/home TMPDIR=/Users/bill/projects/cadence-memories/children/e3-restart-ch11-student/build/tmp GOCACHE=/Users/bill/projects/cadence-memories/children/e3-restart-ch11-student/build/cache GOMAXPROCS=2 GOFLAGS='-mod=readonly -p=1' GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off GOMODCACHE=/Users/bill/go/pkg/mod /opt/homebrew/opt/python@3.13/libexec/bin/python3 /tmp/e3-ch11-r2-live-only.7ui5ezji/control.py
```

**HOLD.** No prepared controller execution/import, key read, upstream GET/POST, browser/native action, source/guidance/manuscript edit, children, commit/tag, or operational retry occurred in this turn. Operational usage remains **POST 0 / GET 1 / key read 1**. Aligned validation is complete; live semantic proof and chapter acceptance remain pending root + SAME reviewer disposition.

## Authorized R2 live-only execution PASS — final HOLD, acceptance still pending

Executed **once** the exact reviewed `control.py` (`cb48d48b…04576`) and unchanged peer (`76138faa…3bf0`) using the full clean-env Python 3.13 command/cwd printed in the preceding preparation section. Pre-execution SHA checks matched controller, peer and new frozen CLI; no file repair/edit, rebuild, grade, local fake, GUI/browser/native operation or agent delegation. Actual command/full controller output: **`$I/4202`, exit 0**.

Pinned CLI **`/tmp/e3-ch11-r2-final/frozen-cli`**, SHA **`a8af82a8cdc1be9e92b335777d640713441062e4a69fa219385fc6f306bbb11c`**; source **`12bd5c3c6c388c525c72e3464f9bce6f21e404241e4efe3e0a0009241c6e1487`**. Successful previous models discovery reused by its pinned hash; **zero new GETs**. Pins/reservation readiness preceded the **one new authorized memory-only key read**. Real key stayed out of argv/files/logs/child environments; CLI used dummy credentials to the owned proxy. Explicit pinned `/etc/ssl/cert.pem`, hostname verification/CERT_REQUIRED and reviewed no-redirect/no-retry/unchanged-body guard retained.

Actual live results:
- **2 POSTs reserved and completed**, both HTTP 200, `claude-sonnet-4-6`, native `max_tokens: 1024`, absent `stream` treated as plain. Four-slot cap was not exhausted: **2 unused slots remain unused**. No fallback/body rewrite/retry. Both requests passed the unchanged pre-upstream declaration guard (ordinary `read_marker` visible; `environment_snapshot` and `session_brief` absent from model declarations). First saved request also directly inspected: ordinary dormant catalogue retained, no automatic declarations/references.
- Real stdio peer `initialize` → `notifications/initialized` → `tools/list`; schema declarations retained `environment_snapshot: round`, `session_brief: turn`. Peer receipt `live/work/peer.jsonl`: **environment_snapshot=2**, texts “Snapshot number 1” then “Snapshot number 2”; **session_brief=1**; **read_marker=1**, returning actual owned `marker.txt` text **`five cobalt keys\n`** (17 bytes).
- Saved canonical dialogue pairs call/result **`toolu_016yu7UpZX8mgYAusZt9kSN4`**, name `read_marker`, args `{}`. Real job handle **3**, `tool_called` **seq 14** → `tool_returned` **seq 17**, same call ID and exact marker bytes. `live/save.json` lines 118–147 and 484–590 retain that ordinary pair; automatic job events are separately marked `automatic:true`, with snapshot text recorded in the arrival/history path. No forged model call/pair or canonical rewrite.
- Completed answer observed **before EOF**: `The exact marker is **"five cobalt keys"** and the current automatic snapshot number is **2**.` Then controlled stdin EOF, **CLI PID 82664 exit 0**, complete final usage and save. Usage **1798 input / 62 output / 0 cache-write / 0 cache-read**; archive **as_of 26**. This was one CLI conversation, no restart exercise.
- Controller reported `LIVE_PASS` at **2.792s active / 2.814s absolute**, within the 60s setup/120s active/180s absolute bounds. Owned proxy listener **62681 closed**; controller `OPERATIONS_PASS`, `finished.error=null`, **exit 0**. Final disk **1,594,937,344 bytes free**, above 1 GiB. Secret reference cleared during controller cleanup. No broad/shared cleanup or changes to old binaries/receipts.

Receipts stay local under **`/tmp/e3-ch11-r2-live-only.7ui5ezji/`**: `operations.jsonl`, `live/1.request.json`, `live/2.request.json`, retained response/API files, `live/cli.stdout`, `live/cli.stderr`, `live/cli-events.jsonl`, `live/save.json`, and `live/work/peer.jsonl`/`marker.txt`. No raw wire/settings/profile content added to Git.

Post-pass receipt-inspection distinction: one offline summary command **`$I/4203`, exit 1**, guessed nonexistent `live/1.response.bin` (and initially looked for nested peer names rather than the actual top-level fields). It stopped before its archive summary; this was **not** a controller/provider/product failure, request or rerun. It was not repaired/re-executed. Direct reads of the known peer receipt and targeted canonical save entries supplied the facts above; source/controller/peer stayed unchanged.

**Cumulative operational usage: POST 2 / GET 1 / key reads 2.** No further capacity used. Qualified prior local RPC/native results carry without repetition or fresh hearing/locality claims. **FINAL/HOLD** for root + SAME reviewer: all released work has stopped; this execution grant and its pass are not chapter acceptance.
