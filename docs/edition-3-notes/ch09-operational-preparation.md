# Chapter 9 operational preparation — final v3

**Preparation complete; browser and paid exercises UNEXECUTED.** Source remains `251ce88fb1677897dbcc6b33661096c5cb39cd8e`, tree `a2add41c306d567c05d98ae760911bc13c18b75a`. All five frozen binaries, served assets and tracked source match their existing manifest. No source edit, rebuild, product test/grader, browser launch, key read or external/provider HTTP occurred. No owned jobs remain.

## Exact finite scope

- **`browser-local`: one muted owned Chrome/profile and frozen GUI**, loopback fake provider, normal viewport, ≤600 seconds and ≤8 requests. No restart, second page, settings/speech/MIC campaign. Show `base → code-tools → helper`, substituted helper marker, actual `write_file`/`read_file` of `note.txt = SKILLS_OK`, pending unload, unchanged system text and next declarations, then actual dynamic write of `after-unload.txt = SKILLS_OK`. Correlate typed call/result IDs with cards, verified answer, DOM and settled screenshots. Informational job notices are not errors. There is no invented answer sentinel.
- **`paid-cli`: one `claude-sonnet-4-6` frozen CLI session**, ≤8 total upstream reservations, ≤60 seconds setup, ≤180 active, ≤240 absolute. Write/read `note.txt = SKILLS_OK`, request pending unload and answer verified content. No paid post-unload repetition, paid GUI, other vendor, retry, discovery or fallback.

Modes are independently root-approved and one-shot; the controller never chains them. A failed case stops, preserving partial evidence. Source/manifest/assets/binary/script/policy pins precede artifact launch; Chrome's owned PID/profile/command and loopback CDP listener precede attach, and GUI ownership precedes client connection. Port 8088/shared profiles are forbidden. Read-aloud stays false; browser is muted. Chrome background networking or host writes are **not** claimed impossible. `--use-mock-keychain` is retained as an observed working configuration, not a root-cause theory.

Only the sealed proxy may read `~/.cr/settings.json` field `directClaudeAPIKey`, and only after paid authorization plus atomic request reservation. It enters process memory and required HTTPS header, never argv, child environment or a key file. Child HOME/environment are isolated and contain only a dummy key. Exact Anthropic POST `/v1/messages`, model and native token value are checked; accepted body bytes are forwarded unchanged without redirects. Credential echo is redacted and fails the run. Raw API traces remain local/hash-only.

## Measured local rehearsal and guards

| Final case | Exit | Seconds | Evidence |
|---|---:|---:|---|
| Frozen CLI via real loopback forwarder | 0 | 4.646 | `preparation/rehearsals/cli-v3/sanitized-effects.json` |
| Frozen streaming GUI via WebSocket, **no Chrome** | 0 | 4.705 | `preparation/rehearsals/gui-v3/sanitized-effects.json`, `typed-facts.json` |
| Compact real fail-closed guards | 0 | 4.66 | `preparation/guard-rehearsal/` |

CLI forwarded 5 requests, GUI 6; every body hash matches the fake upstream, and primary-system hashes stay identical across each session. First response headers were delayed 3.410s / 3.410s (>3 seconds), successfully forwarded. **Measured native CLI `max_tokens = 1024`**, unchanged; GUI uses 4096. Actual files and durable paired results passed. No GUI screenshot/visual credit is claimed yet.

Wrong pin and absent approval yielded nonzero entrypoint exits with audit evidence of **zero subprocess, connection or credential-open entry**. Wrong model/token returned 502 with zero upstream calls. Eight actual HTTP reservations forwarded; the ninth was rejected and upstream count stayed eight. A genuine two-command sequence stopped after the failed first command; the next-process witness was absent. Owned process exits/listener closure and exact argv/cwd/sanitized environment/duration/exit are retained. DOM JavaScript was syntax-checked only, not browser-executed.

## Pins and root execution

- Policy SHA256: `b8d8525398503564367003eaa3c34dc73dab68be971238c86377752e213a0e79`
- Preparation seal SHA256: `070e03c3703a629fe3b29cd1adadc9292563a6bcf07ce1839aaefe910fb43ab1`
- Executable harness: `/Users/bill/projects/ensemble-edition-3-student/.ch09-local/operational-preparation/controller.py`
- Authoritative code/policy/fixtures are the top-level files named in `preparation/seal.json`, **not superseded folders**.

Root supplies a JSON approval with these exact `policy_sha256` / `seal_sha256`, `modes` containing only the independently approved `browser-local` and/or `paid-cli`, `expires_unix` within one hour, and a fresh `run_directory` directly under `/Users/bill/projects/ensemble-edition-3-student/.ch09-local/operational-preparation/execution-runs`. No new Bill-permission request is needed. For each specifically approved exercise, run from `/Users/bill/projects/ensemble-edition-3-student`:

```sh
python3 -B /Users/bill/projects/ensemble-edition-3-student/.ch09-local/operational-preparation/controller.py execute --exercise browser-local --approval ROOT_APPROVAL_JSON --policy-sha256 b8d8525398503564367003eaa3c34dc73dab68be971238c86377752e213a0e79 --seal-sha256 070e03c3703a629fe3b29cd1adadc9292563a6bcf07ce1839aaefe910fb43ab1
python3 -B /Users/bill/projects/ensemble-edition-3-student/.ch09-local/operational-preparation/controller.py execute --exercise paid-cli --approval ROOT_APPROVAL_JSON --policy-sha256 b8d8525398503564367003eaa3c34dc73dab68be971238c86377752e213a0e79 --seal-sha256 070e03c3703a629fe3b29cd1adadc9292563a6bcf07ce1839aaefe910fb43ab1
```

These are separate commands, not instructions to proceed after a failure. Execution writes its own receipt and reserved exercise directory; a consumed exercise cannot silently restart.

## Failures, assistance and integrity

I initially drafted an unapproved provider matrix/paid GUI; root corrected that scope. Those branches **never ran**. The initial local CLI failed on a wrong `dependencies` fixture key before forwarding. An auxiliary source lookup and browser-edit anchor also failed. Later the guard driver expected 403 rather than the actual reject status 502; forwarding was still zero and the sequence stopped. Superseded scripts, scope-reduction diff, failed receipts and the exact corrected final run are preserved. The initial scratch directory was moved intact into the named directory; old argv paths were not rewritten. See `failure-and-scope-history.json` and `preparation-revisions.diff`.

Assistance was limited to the public assignment/guidance, own source/evidence, root's scope correction, and the three explicitly permitted nonsecret Chapter 8 scaffold files. `scaffold-reuse.json` records exact hashes and reused ownership/forwarding/CDP concepts; no prior scenario, permission or success credit was imported. No children/private reviewer/audit/grader input was used.

`index.json` is the exact self-excluded UTF-8/no-NUL portable inventory, including this external report. `local-artifact-index.json` references local artifacts and frozen binaries by hash only; no profiles were created. `preservation.json` binds current public-read hashes and prior seal references. Disk headroom remains above 1 GiB. No genuine blocker remains; only root inspection/authorization and the two bounded exercises are pending.
