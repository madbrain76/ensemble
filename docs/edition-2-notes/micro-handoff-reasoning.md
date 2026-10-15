# Micro-handoff reasoning cleanup

Bill reported on October 14 that `micro_handoff` must remove all encrypted
thinking from the context window. Root confirmed the defect on main
`389b83c4`: the completed-batch reducer removed calls/results and recalled
text, but retained opaque provider parts. The compaction branch had the same
defect. This is a context projection change; original events remain intact.

The correction removes opaque parts and clears private replay metadata from
retained parts when the existing handoff lands. Visible conversation, skills,
declaration updates and the handoff note survive. Reasoning generated after
the checkpoint remains eligible for normal replay until the next checkpoint.

The real HTTP Responses regression failed before the fix: both earlier and
same-batch encrypted blocks appeared in the next request. Independent reviewer
`compaction_runtime_audit` found a second representation: Gemini attaches
`thoughtSignature` to visible text. A Gemini wire regression reproduced that
remaining defect; clearing the copied part's `Opaque` field fixed it without
removing the visible text. The reviewer accepted the corrected source.

Checks on main: focused context, handoff, Responses and persistence tests pass
under `-race` (2.410s); affected-package `go vet` passes; current Chapter 14
grader against explicit `./agent` passes 100/100. Tests cover Anthropic signed
and redacted thinking, Responses encrypted reasoning, Gemini standalone and
text/call-bound signatures, unchanged original events, completed sibling tools,
and full-log/snapshot recovery. These use local provider fakes, not paid calls.
No live process was restarted. Previously saved projections require a new
handoff to clear retained reasoning; this change does not migrate old snapshots.

Carry this fix into the compaction branch and reconstructed Chapter 14 onward
before creating new exports. Preserve historical exports and test receipts.
