# Original Chapter 12 independent review

## Interim original-grader probe — October 10, 2026

Contributor: `enhancement_reviewer_resumed`. This is a requested interim probe,
not initial formal source review or chapter acceptance. The coder is still
validating and refining its implementation. No correction round has begun.
Read the full mandatory `.agents/skills/ensemble-coding/SKILL.md`, carryover,
workflow, original Chapter 12 and its coordinator launch erratum.

Original chapter: file revision `ec41c6e8cf0800e266f7a39ebadf255d94379082`,
SHA-256 `a754bd50ca7f3c7a2da868ff5acbdb98b4d5eddaa05349e7890c4a3dbf6cd4b9`.
The original full grader was stopped during the inherited GUI/headless parity
chain; the coder reports exit 143 and verified its owned processes were gone.
Its partial output remains **interrupted/unscored**, recorded in
`../../solutions/edition-2/evidence/ch12/grader-interrupted.json`.

After that harness release, I copied the original grader and interim student
tree into an isolated scratch directory. Repository HEAD was
`7636c43765b45641aaac43761eae34c229686ab6`; this is the base revision, not a
commit containing the uncommitted student implementation. All 89 selected source
and module-file hashes matched before copying, afterward, and in the copy.
The retained manifest identifies those actual bytes and the original grader
files. This result must not be attributed to later corrected source.

The scratch-only bridge calls unchanged `Ch12Run` and `Ch12Checks`. It explicitly
marks recursive Chapter 11 parity skipped. Command:
`go test ./internal/grade -run '^TestReviewerCh12Probe$' -count=1 -v -timeout 1m`,
with `REVIEW_SUBMISSION` pointing to the copied student and `REVIEW_REPORT` to
the scratch result. The collector exited 0 in 1.311 seconds; that means collection
succeeded, not that all grading checks passed.

| Current chapter check | Actual result |
| --- | --- |
| mcp-handshake | 15/15 |
| tool-discovery | 15/15 |
| mcp-tool-call | 15/15 |
| ephemeral-round | 20/20 |
| reverse-call | 15/15 |
| ws-tunnel | 0/10: could not find `gui/handler.go` or `internal/ws/handler.go` |

The measured current-chapter result is **80/90**. Recursive parity's remaining
10 points were **unrun**, not a measured failure. The raw check array retains
its explicit skipped marker; do not report this as a full 80/100 grade.

Independent inspection of the complete unchanged harness and check table finds
important limits:

- Handshake credit observes the initial request and ID; the harness treats the
  chapter's initialized notification as optional. Discovery credit observes a
  list request, not the resulting registry or model declarations.
- Tool-call and ephemeral credit is assigned after the protocol exchange,
  without an engine turn or actual tool invocation. The claimed separate static
  ephemeral check is absent here. These 35 points do not establish either
  behavior, delivery timing or exclusion from model-visible tools.
- Reverse credit observes a correlated response to an unknown tool; it does
  not establish successful enabled-tool execution or skill/trust enforcement.
- WebSocket grading inspects fixed file names and substrings, not a running
  browser tunnel. Its failure on the student's separate GUI layout is not proof
  of a broken tunnel; passing that source check would not prove one either.
  No dummy file, relocation or core GUI dependency should be added to appease it.

The final review must independently assess the real job/cancellation path,
round/turn ephemera, visibility and reverse dispatch, skill server lifecycle and
browser behavior through the submitted source and relevant actual evidence.
These are existing chapter obligations, not a new protocol or grader framework.
No exception or acceptance of the still-changing implementation is granted here.

Small retained artifacts are under `../../solutions/edition-2/evidence/ch12/`:
`reviewer-interim-source.json`, `reviewer-interim-bridge.go.txt`,
`reviewer-interim-result.json`, `reviewer-interim-stdout.txt` and
`reviewer-interim-exit.json`. Manifest SHA-256:
`526e412e397b7eb00e3e15a5e0624ee7717fcf7ea544b26923c11b405ae1c16c`.
The harness matches file revision `f9028379fadf8555232534d412bb2208e7c7682a`
and SHA-256 `a7c91fe21dc09412fd0bdb31427986f138ebfe45f8c3964697657aeaa6f3573d`;
the checks SHA-256 is
`9e3e74c172fb479376133a5a11b4e44cd3eb36e043502d74972fae95abeeb630`.
No original grader, chapter or student implementation was edited. No additional
paid call or broad test sweep was run. Scratch copies were removed after capture.
