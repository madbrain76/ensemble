# Chapter 10 operational harness v2 — partial handoff

**NOT overall preparation or operational success.** Local CLI/GUI/guard checks completed. The single owned-browser attempt failed during setup and is consumed. **Paid Claude flow: NOTEXECUTED.** No browser/provider retry, credential read, source edit, rebuild, product suite, grade, commit, or tag follows this failure. Root directed this partial seal; broader acceptance and any recovery remain root-owned.

## Authority and immutable inputs

The original static-only task remains exactly 7335 bytes / SHA256 `6fd180b5d69f5427e048c8960b5f42a2be9c4a4c481b2cd67872dafadf080d37`. It was not rewritten or backdated. A **later prospective root message** (identified by root as send result job4087), following Bill’s authorization of necessary operational access, lifted the execution prohibition/approval gate and authorized local rehearsals, one bounded muted GUI/browser flow, then a conditional Claude flow. Root also permitted materializing finite, matching approvals. Subsequent KISS/continuation steering did not expand scope. Parent confirmation `main-34` in job4102 establishes this chronology; no private execution log or authorization/protection manifest was read. Exact delivery wall time of the original authorizing message was not captured; approval issue/expiry timestamps are retained. See `authority-transition.json` and `approval-pins.json`.

Frozen student head `729bc919e734405307b4ac0060e65c40e33652e5`, tree `025ef3744b42c8f9daec34e76661edc40583ca6a` remain unchanged and tracked-clean. CLI SHA256 `a18f824da32a225bd6960ebf5d6128420da2a065d640ea38269a31fa2a13cc33`; GUI SHA256 `fb506081f98c09d2e6da628edbb6cb922f067f9db24983ae856d2d1222dcbf53`. Chrome executable `/Applications/Google Chrome.app/Contents/MacOS/Google Chrome` was pinned and its owned launch identity checked. `final-preservation.json` verifies **147 distinct named source/artifact/tool files and all254 prior immutable entries**. Old unaccepted bundles, scope reconciliation, blocked seals, reports, failures and real TAR remain unchanged. No formal student production-code correction round occurred.

## What was built and actually run

The current helper is `/Users/bill/projects/ensemble-edition-3-student/.ch10-local/operational-preparation-v2/launch.py`; explicit modes are `cli-rehearsal`, `gui-rehearsal`, `guards`, `browser`, and `paid-cli`. No mode implicitly launches another campaign. The exact command interface is:

```
<pinned-python> launch.py --mode <exact-mode> --run-root <fresh-mode-root> \
  --approval <finite-root-approved-record> --expect-seal <exact-sha256>
```

Policy/seal/run-root/expiry and source/helper/executable pins precede entries. The shared driver launches the actual frozen client twice with the same workspace/config/save override. Its fake Anthropic server requests the real `read_file` tool; it does not synthesize a tool result. Assertions cover actual returned fixture content, typed call/result linkage, the four-field archive, retained history prefix, larger `as_of`, startup request counts, and a new nonempty post-restart reply. The raw-GUI adapter uses one connection per launch and checks old IDs versus new final IDs.

The one-provider proxy validates **Claude `claude-sonnet-4-6`, `POST /v1/messages`, and native max_tokens 1024 (CLI) or4096 (GUI)**, rather than rewriting payloads. A shared atomic four-request allowance covers both launches, counting attempts before upstream/key access. Bodies and legitimate tool responses are forwarded unchanged. The paid path is conditional on a root-accepted successful browser receipt plus cleanup, uses the narrow authorized key field, supplies only a dummy value to the application, permits only the normal Anthropic origin, and has no short upstream response-header deadline. That paid path **has not been executed or validated against a real provider**.

### Existing receipts

| Attempt | Start UTC | Wrapper duration | Exit | Outcome |
|---|---|---:|---:|---|
| CLI setup v1 | 2026-10-19T04:39:26.975036+00:00 | 0.249s | 1 | FAIL: early-cli-exit |
| CLI restart pair v2 | 2026-10-19T04:40:38.661816+00:00 | 5.025s | 0 | PASS: checks completed |
| Raw-GUI restart pair v2 | 2026-10-19T04:42:06.820177+00:00 | 0.553s | 0 | PASS: checks completed |
| Compact guards v3 | 2026-10-19T04:43:40.121670+00:00 | 0.737s | 0 | PASS: checks completed |
| Owned browser A v3 | 2026-10-19T04:45:17.513570+00:00 | 67.328s | 1 | FAIL: setup-deadline |

- **CLI-r1:** harness used an invented `protocol` positional subcommand. Frozen CLI correctly rejected it before any request. Its sole child exited1; cleanup completed. Bundle v2 removes only that argument and uses the default protocol entry.
- **CLI-r2:** two actual CLI launches, three fake provider requests. Real tool/result content and restored continuation passed. The first successful fake response withheld headers for3.5s; the measured observation is retained in the completion/proxy records, not inferred from configuration alone.
- **GUI-r1:** two actual GUI launches and one raw-WebSocket connection per launch, three fake provider requests, actual tool card/result and saved-history identity checks, and a post-restart final passed. This is **transport/application evidence, not browser rendering evidence**.
- **Guards-r1:** actual entry-point subprocesses exited65 for wrong seal and missing/mismatched authority before application entry. Actual proxy requests demonstrated wrong model/token denial before upstream, four forwards with the fifth denied, and no later phase after a first failure. Seven broker POSTs / four fake upstream calls; no credential opens.

There were **7 top-level application children** (3 CLI,3 GUI,1 Chrome), plus2 deliberate preflight Python probes. Across all modes:13 broker POSTs and10 upstream connect calls, all to owned loopback fixtures. Those counts do not include GUI asset/CDP HTTP traffic and are not host-wide network accounting. All five Python credential-open audits reported0 attempts/0 authorized opens. This is not an OS-wide file-open audit of inherited applications or Chrome.

## Consumed browser attempt A — failed setup

`browser-r1` launched the owned GUI PID98223 and fresh-profile Chrome PID98226. The saved `chrome-owned-identity` record contains the observed executable/profile/command, debug listener ownership and helper ancestry. Dynamic owned ports were used, not live8088. Read-aloud was disabled in the newly owned configuration and Chrome was launched muted; no theme/settings/TTS test was conducted.

The attempt timed out at its60s **setup** bound. There is no `browser-page` completion milestone, phase DOM JSON, phase screenshot, user/provider request, or GUI restart for this attempt. **No visual inspection or successful UI observation is claimed.** A persistent CDP command/frame trace was not retained; the exact missing acknowledgement/navigation/replay step and underlying cause are not inferred. Root can inspect the preserved scripts, event timestamps, Chrome/GUI stderr and profile references without replay. This handoff does not attribute the failure to production code.

A CDP close request was recorded, but graceful Chrome exit was not confirmed. Cleanup then used TERM followed by KILL on the owned Chrome group; Chrome exited-9. GUI exited0. The harness recorded owned groups absent, owned HTTP resources closed, cleanup complete, and exit1 after67.272s (wrapper67.328s). Those are bounded recorded cleanup observations, **not** a new account-wide absence proof. No further diagnosis or runtime probe was performed after root’s sealing instruction.

**B remains NOTEXECUTED:** no paid approval was materialized, no key was obtained, and no real-provider request or usage result exists. A consumed browser attempt will not be retried or silently reclassified under the standing authority; root will decide any narrowly justified recovery.

## Versions, assistance, preservation and evidence limits

Complete scripts/policy/schema/fixture/manifests are copied byte-for-byte in `bundle-v1/`, `bundle-v2/`, and `bundle-v3/`. Their seal SHA256 values are:

- v1 `62642cad94c25bdc4c4f6468b6cb2683ae3d2830c3444e0689667f0b9b19b4e7`
- v2 `87e7bdac62f77b7971b51a908ad1d477f6d7b0d19bb7b660c342eea6281c2115`
- v3/current `7bb5f4c00943675c022dd6cd15cf3ae76219e9ead4cc71fc3140234d1f4a30f6`

The `NOTREHEARSED_NOTEXECUTED` field in each historical seal describes creation time, not current outcomes. Bundle v3 changes only the browser observer from v2: await fonts/two animation frames before reading DOM and recheck errors after capture. The passed CLI/raw-GUI paths are byte-identical; no completed local pair was repeated for that change. Bundle v3 guards and failed browser run have their own exact pins and approvals.

Own prior unaccepted drafts supplied framing/ownership ideas, but their generic three-vendor proposal and cancellation/sentinel campaign were not adopted as authority or successful evidence. v2 replaces those with the bounded two-flow driver and real guarded entry points. Full assistance and non-operational tool/read/sealing mistakes are recorded in `assistance-and-failures.json` and `sealing-metadata-failure.json`, alongside the retained CLI and browser failures. The only functional harness correction was versioned local setup work, not a production-code correction.

Each `*.receipt.json` contains actual argv, cwd, clean/offline environment, timestamps, duration, exit, process group, approval pin, helper pins before/after and stdout/stderr hashes. `attempt-records/` copies existing completions/events and generated approvals with provenance; `attempt-summary.json` is a concise roll-up. **Raw traces, save files, profile data, captures and binaries remain local-only/excluded from Git**, referenced by `local-artifact-index.json`; no binary placeholder or source archive copy was created. Private/forbidden filenames are not exported. Local artifact hashes were measured at sealing and are distinguished from contemporaneous receipt pins. The final exact UTF-8/no-NUL inventory includes this report and excludes itself.

Root’s final instruction was to seal this partial state, preserve all failures and old seals, and stop. Existing jobs have finished. **Preparation/operations are not accepted or complete.**
