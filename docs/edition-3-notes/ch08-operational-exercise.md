# Chapter 8 — bounded operational attempt (sealed partial evidence)

## Disposition

**Stopped at the first unresolved owned-browser setup failure. Not an application defect finding, operational pass, or overall acceptance.** R1 code remains independently PASS/root code-review accepted as reported by root. Candidate `b8f22ba5b3c9d1fc52d1a114810b90b77928684d` / tree `fa896caa404b87b08de61eb38ce079eb38138fbe` remains clean. All 126 tracked members and nine retained artifacts are byte-identical. Existing reports, source manifests, and sealed plan/evidence were verified unchanged. No production/test/source edits, rebuilds, grades/lifetime reruns, commits/tags or child agents occurred.

## Authority, scope and actual chronology

Read the complete 7696-byte execution task (SHA-256 `0709ca122d7dc35a5ac4ad6817288d3c0759ace23d25a08bc968cd36f81eca77`), loaded/read GUI-debug guidance, and verified both scoped approvals against the immutable launcher/policy before launch. Earlier full prerequisites were retained only after exact public-file pin verification. No private reviewer/coordinator material was used. Root's execution authority was followed without another routine permission request.

Created a separately pinned, owned CDP controller; the sealed launcher and policy were **not modified**. Before execution, source inspection corrected the new controller's subscriber ordering and added a single default-voice qualification read to its then-unexecuted native branch. The pre-protocol version is retained and was not exercised. Executed controller SHA-256: `255d0df1645986480539fcfeae14db1420066b798e248683fba589b130ccbb3c`. Full script, command/receipt, version and pin records are in this package; this run does not validate the controller's unexecuted UI/native/paid branches.

The only execution command was the retained receipt's exact `python3 -B .../exercise.py --phase local-native --script-sha256 ... --run-id local-01`, through the existing clean-environment receipt wrapper. It invoked the exact sealed local serving command and supplied approval. Frozen GUI hash: `6b513f3b10c5906bf66c6acf89e645903c95b0a397b711d93ae067b5b35ec0ae`. The launcher verified nine served asset hashes before browser attachment. GUI used its owned directory, owned HOME, dummy Anthropic key, and read-aloud off.

One new **headed, muted** Chrome/profile was started, with loopback dynamic CDP. Version `Google Chrome 155.0.8059.80`; executable SHA-256 `5f423d624250e67f2ae453400d485e7547a533a2f872232c54bb551a79bf552d`. Browser PID62717, profile under the owned `local-01/browser-profile-visual`, CDP127.0.0.1:52107. Exact argv/profile/executable/`ps` and PID-specific `lsof` ownership were durably recorded **before** WebSocket attachment. No8088, shared/personal browser attachment, focus automation, OS-volume/voice/MIC changes, or native callback replacement was performed.

CDP successfully enumerated the owned target, attached a flattened session, and acknowledged Page/Runtime/Network enablement. `Page.navigate` to the owned GUI127.0.0.1:52097 emitted navigation-start/loading events but did not return within 30 seconds. There was **no completed application navigation or application DOM/geometry/control assertion**. No screenshot was captured or visually inspected. The run stopped, issued Browser.close and shut down the owned serving process. Total measured command time: **41.056 seconds**, exit1. No second browser style/profile or paid retry followed.

## Failure, host-side observations and limits

The actual failure is `RuntimeError: CDP timeout Page.navigate`. Chrome later reported a teardown-watchdog expiry and exited **2**. Its retained log also contains an allocator multiple-load warning. These are observed symptoms; the root cause is **unproven**. Nothing here establishes that the frozen application caused the failure.

**Unexpected browser-internal side effect disclosed:** despite the controller's background/component-update disabling flags, Chrome automatically spawned its bundled `GoogleUpdater --install` child (PID62764, parent62717), which reported return3; a later wake-all reported return0. The controller did not invoke an installer/updater or change updater settings. This is not proof that no host-side writes/installation/network activity occurred. No host installed-state audit or undo was attempted, and no shared updater/cache/profile was cleaned. The main Chrome executable's hash remained unchanged. Raw diagnostics remain local and hash-indexed; only this safe observation is exported.

A packaging-only failure attempted to hash a nonexistent guessed GUI-debug filesystem path. The skill was actually loaded via `load_skill`; no filesystem skill hash is claimed. The unsealed attempt and original packager were preserved; a versioned packager corrected provenance without any operational rerun. An earlier pre-browser filesystem lookup of that path also failed.

A subsequent diagnostic guessed the wrong completion filename, then located the actual `completion.json`; that lookup failure is retained as a qualification, not a new exercise attempt. The sealed launcher's `browser_launched=false` field refers only to that serving process: the separate owned controller really did launch Chrome. Browser-specific ownership/attachment/exit receipts govern that fact.

## Evidence gates and counters

| Gate | Actual result |
| --- | --- |
| Frozen pins, scoped approvals, owned serving setup, nine asset hashes | Verified |
| Owned Chrome PID/profile/listener before CDP attachment | Verified |
| Initial application navigation | Failed to complete; unresolved setup/Chrome failure |
| Rendered panes, controls, themes, font/contrast, resizing/narrow layout, two-page synchronization, streaming/tool/pause/scroll/reconnect | **Not exercised**; no DOM/screenshot pass claimed |
| One allowed GUI restart | **0 used** |
| Local fake upstream requests / completed GUI model turns | **0 / 0** |
| Native rate/start/end/cancellation and audible window | **Not attempted**; 0 nonempty utterances requested |
| Human listening | **No report**; no hearing/silence claim |
| Paid model turn / external model requests / real proxy key reads | **Not attempted / 0 / 0** |

Native and paid were correctly gated off by the unresolved earlier failure. The proxy's zero external/model count is not a claim about every possible Chrome background network operation. No provider discovery, fallback, standalone key reading, provider settings export, or credential handling occurred. The owned GUI's dummy credential is not a real provider credential. No paid request budget was consumed.

## Cleanup and preservation

Actual exits: GUI PID62716 **0**, serving launcher PID62709 **0**, Chrome PID62717 **2**, outer controller command **1**. All recorded PIDs are absent, all eight Chrome-log helper PIDs are absent, no process matches the owned profile, and an additional executable/start-window check found no surviving Chrome/updater process from this launch window. Proxy52096, GUI52097 and CDP52107 all reject connection with code61. No broad process kill was used. Profiles, raw traces, logs and crash material remain on disk; no shared cache/profile deletion occurred. The task has **zero surviving owned jobs/processes/listeners**.

## Seal and next boundary

`ch08-operational-exercise-evidence/MANIFEST.json` specifies the exact safe-text portable set, notes-directory-relative base, every included file's size/SHA-256, and explicit self-hash exclusion. Raw CDP/browser/API/history/debug logs, profiles, binary/cache/crash material and any images are local-only and separately hash-indexed; symlinks are not dereferenced. No TAR/PNG/binary/raw trace enters the portable set. The report and package are sealed: no later append/reseal or silent reinterpretation of the failed receipt.

Root can inspect this partial evidence and reconcile the owned Chrome startup/navigation failure before deciding any bounded continuation. Actual GUI, native and model evidence remains open. This submission does not reopen the closed code round, claim operational success, accept/tag the chapter, or start author work.
