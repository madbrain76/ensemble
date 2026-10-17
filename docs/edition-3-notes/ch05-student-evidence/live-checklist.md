# Finite Chapter 5 live checklist — pending parent review; no calls made

Before paid execution: final source commit and complete source manifest/archive;
build CLI/workflow from that frozen tree; pin hashes in receipt launcher; prove
wrong-source and wrong-binary refusals and run launch/capture/verifier with local
HTTP fakes first. Keep raw safe stdout/stderr, submitted inputs, local side effects,
provider/model/config and measured usage. Never serialize credential values,
thinking or signed opaque content. Any later source change invalidates final-byte
claims and requires appropriate rerun. Prior failed fake grades remain retained.

| Feature | Exact action/provider | Required observable result |
|---|---|---|
| Human CLI not deaf during work | Anthropic `claude-sonnet-4-6`, actual CLI `chat`, one initial prompt explicitly requesting a shell job that prints START, sleeps 8 seconds, writes FINISHED to a file, then prints it. Once retained job output contains START, write `/hint use concise prose` while job still running. | Human stdout contains hint acknowledgement before real FINISHED; timestamps/input transcript and job files establish ordering. |
| Interrupt does not kill job, follow-up works | Same human session: after hint acknowledgement, write `/interrupt`. Wait for actual retained FINISHED side effect, then a second prompt asking to use wait_for_job/read_file on the retained job and report the observed FINISHED text. | Interruption reply/fact before FINISHED, no kill fact; same real job reaches done, bytes retained, second reliable response succeeds; shutdown and saved safe terminal history agree. |
| Real application-owned roles/handoffs | OpenAI `gpt-4.1-2025-04-14`, actual separate workflow executable; one 80-word prose topic prompt requiring author write_draft/word_count, editor read_draft/edit_draft, reviewer review accept/reject. | Real manuscript and review.json; role-tagged safe observations and author_done/editor_done/reviewer_done markers; per-role safe request/job/end summaries prove reliable completed-request handoffs. |
| Local image consumption | Gemini `gemini-3.8-flash`, actual CLI rich-parts input; create a small PNG with colored regions and ask for their arrangement without describing the contents or meaningful filename. | Actual file bytes consumed; provider identifies visual arrangement, measured usage and safe answer; persisted history contains reference not base64 bytes. |
| Unsupported media refusal before network | Same final CLI locally with Anthropic selected row and an unsupported video blob. Point only this refusal check at a local counting HTTP server. | Explicit model+video refusal; zero requests. This is fake/network-boundary evidence, not a paid provider refusal. |

Limits: one attempt for each provider scenario; no silent model/cap changes or
automatic retry. CLI/workflow `LLM_TOOL_BATCH_LIMIT=3`, MaxTokens=1024 per call:
human at most 2*4 calls, workflow at most 3*4 calls, image exactly 1 call => at most
21 model requests. Initial human fixture must finish within 90s, workflow within
120s, image within 60s; these harness watchdogs are evidence limits, not framework
request deadlines. Stop and ask for bounded rerun if required evidence is absent.
No provider readiness campaign. Prior selected supported models remain unchanged.
Keys read only inside launcher process from approved settings fields; only the
selected key enters the actual executable environment, never argv or child shell
environments. No credentials are printed/copied/committed.
