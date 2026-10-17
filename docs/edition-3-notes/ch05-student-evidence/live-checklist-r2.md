# Chapter 5 finite live-evidence checklist — revised, NOT YET AUTHORIZED

Production is frozen at clean commit `00489b05b65cfe5209f26cb655667df44b3847e2`. No production edit or rebuild is proposed for this checklist. Parent approval is required before `--live`; no credentials have been read and no paid request has occurred.

## Identity and launcher

- Manifest: `/tmp/cadence-ch05-evidence-r1/frozen/identity.json`.
- CLI: `/tmp/cadence-ch05-evidence-r1/frozen/ensemble`, SHA256 `258233df606e71d50aa0103c5865fcd2e383efc471c3d3f8e71fe59f07609863`.
- Workflow: `/tmp/cadence-ch05-evidence-r1/frozen/workflow`, SHA256 `c5c2fee6aa8c463853c1c6210f7b2a699be8bf6806e550ba0b28d0aa44507114`.
- Launcher: `/tmp/cadence-ch05-evidence-r1/launch.py`, SHA256 `c582a4e322b9dd60b672a0133c673a54f09a59d80bb9947f2d6b396fc6613c82`.
- Launcher checks all source bytes, clean commit, archive and both binary hashes before any credential lookup or process launch and again before each scenario. Negative wrong-source/wrong-binary checks must pass on this launcher.
- Final local controls and complete fake rehearsal must pass before requesting approval. Earlier failed harness runs and their receipts remain intact.

## Credentials, privacy, and hard cost boundary

Read only `~/.cr/settings.json`, after approval and identity verification. Use only `directAnthropicApiKey`, `directOpenAIApiKey`, and `directGeminiApiKey` for their matching vendor. No environment/dotenv/alias fallback. Values stay in launcher memory and each selected executable's `LLM_API_KEY` environment; they never enter argv, printed commands, source, reports, manifests, or stored input. The shell-job environment is scrubbed by the executable and checked without printing a value. The selector is tested using an in-memory dummy settings fixture, never real credentials.

A loopback forwarder connects only to each fixed official provider origin. It records counts, status, durations and byte lengths, not auth headers, query strings, prompts or provider response bodies. No upstream retries or redirects are followed. **Hard maximum 21 upstream requests total:** human 8, workflow 12, image 1. The first over-budget attempt receives local HTTP429 without an upstream connection. A local negative control must establish two attempts but only one upstream request at cap1. Each application request has the existing 1024-output-token limit; tool batch limit is explicitly 3. The proxy ceiling is independent of the separate explanation-refusal guard.

Capture safe stdout/stderr, exact stdin with relative times, process exits, filtered event logs, actual role/job files, provider/model labels, request-count receipts and all artifact hashes. Omit thinking/opaque/signature payloads; do not copy raw sensitive provider history into evidence. Raw executable runtime files live in a launcher-owned temporary directory and are removed after safe capture. If a selected key appears in output, withhold that capture and fail. Do not silently change a model or retry a failed run.

## One finite run

Proposed command (no keys on command line):

```
python3 /tmp/cadence-ch05-evidence-r1/launch.py \
  --manifest /tmp/cadence-ch05-evidence-r1/frozen/identity.json \
  --live --out /tmp/cadence-ch05-evidence-r1/live-r1
```

Scenarios run sequentially. First failure stops the remainder and preserves safe partial receipts; any retry requires a new parent decision.

1. **Actual human CLI / Anthropic `claude-sonnet-4-6` (cap8, expected4 requests).** Launch the frozen CLI in `chat` mode. Ask for a real `run_command` job that prints START, sleeps20s, writes FINISHED to `retained.txt`, and prints it. After observing actual START in retained job output, send `/hint`, require visible acknowledgement while the job remains unfinished, send `/interrupt`, require reliable interruption acknowledgement, then send a genuinely new prompt. Require its reliable `NEW_TURN_OK` answer **before the old job's FINISHED side effect**. Only afterward wait for FINISHED, ask `wait_for_job` handle1, and require its actual successful tool result to contain FINISHED. Preserve ordering timestamps, saved events, job output, side effect and clean exit. No simulated/synthetic completion is accepted.
2. **Actual three-agent workflow / OpenAI `gpt-4.1-2025-04-14` (cap12, expected6 requests).** Run the frozen public-importing workflow with one roughly80-word rainy-library prompt. Require author/editor/reviewer reliable completion markers in order and successful call-correlated `write_draft`, `word_count`, `read_draft`, `edit_draft`, and `review` results under the correct roles. Require a nonempty authored text, correct word count, editor read of that authored text, a nonempty genuinely changed edited text, final file equal to that edit, and an accept/reject review with nonempty notes whose result and artifact match the actual tool arguments. Done markers, empty files, or fabricated uncorrelated observations cannot pass.
3. **Actual local-image path / Gemini `gemini-3.8-flash` (cap1).** Supply a generated PNG containing a red left half and blue right half using the real CLI rich-Part path reference. The prompt does not name the colors or their answer. Require structured JSON `{"left":"red","right":"blue"}` from the model, not merely the words red/blue/left/right. Require stable references and no base64 image body in persisted history. Preserve fixture bytes/hash, filtered history, request count and safe final answer. A reversed-orientation negative control must fail locally.
4. **Unsupported media / local-only boundary.** On the same frozen CLI, use Anthropic's selected row with a video Part and a local recording fake. Require explicit model/kind refusal and zero HTTP requests. This is fake/preflight evidence, not a live media-support claim, and uses no real key or paid request.

## Current readiness — ready for parent decision, NOT authorized

- Frozen source and both binary identities verified. Production remains unchanged.
- Root and all module vet/race/build gates passed; root coverage85.5%; six valid mutation detections. Earlier inconclusive runs preserved.
- Frozen public Chapter5 and Chapter1 grades each passed100/100 with exact process exit0. Commands and receipts: `frozen/grade-commands.txt`, `frozen/grade-exits.txt`, `frozen/grade-ch05-r2.json`, `frozen/grade-ch01-r2.json`.
- Final local controls: `launcher-controls-r2/results.json`, all passed/rejected as intended; `identity.json` records launcher/control/fixture hashes and zero real credential reads/paid calls.
- Complete revised launcher rehearsal: `local-launcher-r6/summary.json`, passed; process exit0. Its human ordering, nonempty/changed role artifacts and actual review, correct image orientation, and unsupported pre-HTTP refusal all passed.
- No build or test is being duplicated. No paid call has run. Await explicit parent authorization for exactly the finite command above.
