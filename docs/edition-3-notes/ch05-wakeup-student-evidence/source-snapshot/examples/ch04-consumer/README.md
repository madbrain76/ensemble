# Independent public-package consumer

This is a separate Go module with a local replacement for the submitted root
library. It imports only `ensemble`, retains the ten builtins, registers `shout`
through `Agent.RegisterTool`, and uses normal `Agent.Ask` dispatch. Its handler
uses `ToolContext.Agent.Parent().Logger()` (the actual owning root) for a
timestamped diagnostic. The CLI shuts down root-owned jobs before dumping its
neutral event log.

Build from this directory with `go build -o /absolute/path/consumer .`.
`consumer-case.json` documents the tool name, input, output, and diagnostic used
by the public Chapter 4 grader. Environment mapping is explicit in `main.go`:
`LLM_VENDOR`, `LLM_BASE_URL`, `LLM_MODEL`, `LLM_API_KEY`,
`LLM_REQUIRE_TOOL_EXPLANATION`, `ENSEMBLE_WORKSPACE`, `ENSEMBLE_DATA_DIR`, and
`CH02_LOG`. No credentials belong in this module or its case file.
