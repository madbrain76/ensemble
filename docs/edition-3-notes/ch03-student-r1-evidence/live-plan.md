# R1 pre-live plan — written before the single authorized paid session

Pinned commit: `28be8bffdbc1356eaaa50a97369e59d0e0556dc3`.
Pinned CLI: `/tmp/ensemble-ch03-r1/student-cli`, SHA256 `db9a31df81f88cace5edcdde6af50dfa35ffe25e96fbc7d7585e51dd6e264dde`.
Use actual `chat`, not JSON-lines protocol. Gemini `gemini-3.8-flash`, default GenerateContent endpoint, existing default max output tokens (1024), existing ten-round budget, existing `LLM_REQUIRE_TOOL_EXPLANATION=true`. No cap/model/options changes. Prior discovery and access were already recorded; this is the one parent-authorized R1 live session, not a new discovery campaign. No paid Anthropic/OpenAI reruns.

| Interactive feature | Concrete step | Independent observation |
|---|---|---|
| Human chat input/output | Feed two ordinary lines to `student-cli chat` in a fresh `/tmp` project; close stdin afterward | Both terminal normalized visible answers must be nonblank and present on actual stdout; no opaque output dump |
| Discover directory/files | First turn requests `list_directory`, `read_file`, and `search_files` | Named successful tool calls/results; no claims inferred merely from prose |
| Execute command and read failure | First turn requests `run_command` with `go test ./...` before repair | Independent initial test is exit1; model sees actual command feedback |
| Exact repair | Agent must `edit_file` only the subtraction bug in `add.go` | Real disk changes to addition; test and module hashes unchanged |
| Create/read persistent follow-up artifact | Agent must `write_file` an explanation marker; second turn reads it using tools | Marker exists on disk and is discovered in paired tool feedback |
| Verification from tool feedback | Agent runs `go test ./...` after repair | Independent post-run tests exit0; successful command result exists |
| Conversation preservation | First line contains a harmless recall token, explicitly not for project files; follow-up asks for it | Second visible final contains that token; token absent from project sources/marker; prior response precedes follow-up in durable History |
| Signed text fix | Model may emit signed ordinary text naturally | Record only whether terminal text is nonblank, signature attachment/provenance, and stdout agreement; never opaque content |
| Selected explanation gate | Existing CLI policy enabled, without model-name switch | Record refusal/retry and paired result counts honestly; deterministic local tests are the primary ordering/no-effect proof |
| Cancellation/finite work | External 240-second watchdog around this one complete two-turn interaction | Observation limit is not a provider/command work deadline. Timeout/error/incomplete answer is not a pass; no automatic retry |

Credential launcher loads only the authorized directGeminiAPIKey field for the API client environment in process memory; no key is put in arguments, prompts, files, or reports. Tools run only in the disposable project and inherit the production scrubbed environment. Benign-canary environment tests remain in the full suite. Scope instructions are not an OS sandbox.

Retain only sanitized metadata, hashes, tool names/IDs/error flags/part kinds/lengths, independent project test receipts, and project artifact identity in committed evidence. Raw provider blocks or assistant reasoning are not printed or copied to the report. The private local normalized History artifact is identified by path/hash for independent authorized review.
