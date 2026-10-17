# Chapter 5 remaining-only live checklist — NOT YET AUTHORIZED

## Existing evidence and frozen identity

Production remains clean at `00489b05b65cfe5209f26cb655667df44b3847e2`. Manifest `/tmp/cadence-ch05-evidence-r1/frozen/identity.json`, SHA256 `950d5c3fb732ad49ff132bbb240e3e0405628c623313cc28da23baa458c6ac19`, fixes the source archive, source bytes and clean-VCS race-built CLI/workflow. Both final public grades are100/100 with exit0; local production gates/coverage/mutations are unchanged and not rerun.

The first authorized live command remains **failed** (`live-r1/summary.json`, launcher exit2). Its human CLI exited0 after four Anthropic HTTP200 requests. The failure was an evidence substring check matching the assistant's sentence “No KEY_ENV_LEAK was printed”; actual guarded process output had no sentinel, real job exit0, actual retained output inspected. New user turn completed18.958662seconds before the old job's actual terminal report. Parent accepted this diagnosis and authorized evidence reconciliation, not rewriting the failed run as PASS. Exact artifact hashes and timestamps: `live-r1-human-observed-facts.json`. **No human inference will be repeated.** Workflow and image paid requests so far: zero.

## Bounded evidence-only delta and validation

Current launcher SHA256 **`68066167187e3618beba4d9c8a49b59ee5f2ba105286dc3ee23235960c4cbfad`**. Preserved prior launcher and exact diff: `guard-route-r4/launch.py.before`, `launch.py.diff`.

Changes only: check the guard sentinel in actual `cr/io/1`, not model prose; add `remaining` to the existing scenario selector; skip human for that route; record `human_run:false`. Approved settings-field selection, all cap/forwarding code, model choices, and workflow/image/unsupported scenario bodies are byte-for-byte unchanged.

- `python3 /tmp/cadence-ch05-evidence-r1/guard-route-controls-r4.py` — exit0: benign assistant-prose positive, actual malicious-process-output negative, original human artifact/hash reconciliation with no inference, no-human route (human/key lookup forbidden during dispatch control), and cap1 producing one upstream request plus a localHTTP429 on attempt2.
- `python3 /tmp/cadence-ch05-evidence-r1/launch.py --manifest /tmp/cadence-ch05-evidence-r1/frozen/identity.json --scenario remaining --out /tmp/cadence-ch05-evidence-r1/local-remaining-r4` — exit0: real frozen executables against local fakes; only workflow/image/unsupported ran, all passed.
- One local argument-parser insertion failure was caught before those passes and is preserved as `guard-route-r4/first-control.stderr`; no real credential or provider was involved. No production change/build occurred.

## Exact proposed command (requires fresh parent approval)

```
python3 /tmp/cadence-ch05-evidence-r1/launch.py \
  --manifest /tmp/cadence-ch05-evidence-r1/frozen/identity.json \
  --scenario remaining --live \
  --out /tmp/cadence-ch05-evidence-r1/live-remaining-r1
```

One invocation, sequential, stop on first failure, no automatic retry, model substitution or extra call.

1. **OpenAI `gpt-4.1-2025-04-14` workflow: hard maximum12 upstream requests.** Require actual successful call-correlated author `write_draft`/`word_count`, editor `read_draft`/`edit_draft`, reviewer `review`; nonempty authored text; correct count/read; nonempty genuinely changed edit matching final manuscript; nonempty actual review notes and accept/reject result matching saved review; reliable completion markers in order. No done-only/empty-file acceptance.
2. **Gemini `gemini-3.8-flash` image: hard maximum1 upstream request.** Actual local PNG bytes through CLI media path, exact structured left:red/right:blue result, stable reference and no base64 body in saved history. Reversed orientation fails. Prompt does not supply the answer colors.
3. **Local-only unsupported-media step:** same frozen CLI, Anthropic selected row with video; explicit refusal and zero HTTP calls. Fake fixture key only, no paid call.

No Anthropic key selection or human process is reachable in the remaining route. Only the matching exact `directOpenAIAPIKey` and `directGeminiAPIKey` fields in approved `~/.cr/settings.json` are selected for paid scenarios, process memory/environment only. No other credential source or fallback. Model response limit1024tokens/request; tool batch limit3. Forwarder enforces caps independently, follows no redirect/retry, and records counts/status/lengths only. **Maximum additional requests13; including prior human4, cumulative maximum17.** Existing authorized ceiling21 is not a license for extra requests.

Same safe capture and finite watchdogs as the reviewed launcher. Preserve failures and partial stdout/stderr/events/exits/counts; strip thinking/opaque/signature data; no keys on argv or in artifacts. Actual human evidence stays attributed to the original run and its independent reconciliation; unchanged earlier controls stay attributed to their original hashes.
