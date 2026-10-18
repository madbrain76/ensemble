# Narrow reconciliation excerpts

These are quotations from existing files, not new execution output. Complete originals remain in place. File hashes and line counts are recorded in pins.json.

## Existing sensitivity.stdout: complete output (7 lines)

```text
=== RUN   TestSkillsThroughLoop
    skills_test.go:133: tool not enabled: write_file
    skills_test.go:159: dynamic provenance: {Source: EventSeq:0}, expected 7
--- FAIL: TestSkillsThroughLoop (0.01s)
FAIL
FAIL	ensemble/ensemble	0.013s
FAIL
```

## Current solutions/ch09/ensemble/skills_test.go

Lines 48-51:

```go
func toolNames(a ensemble.Agent) []string {
	var names []string
	for _, declaration := range a.Declarations() {
		names = append(names, declaration.Name)
```

Lines 86-88 (write precedes unload):

```go
			call("write", "write_file", map[string]any{"path": path, "content": "enabled"})
			call("next", "load_skill", map[string]any{"name": "next"})
			call("unload", "unload_skill", map[string]any{"name": "code"})
```

Lines 131-134:

```go
	for _, name := range []string{"write_file", "list_directory", "search_files"} {
		if !slices.Contains(toolNames(a), name) {
			t.Errorf("tool not enabled: %s", name)
		}
```

Lines 158-160:

```go
	if meta, ok := a.Skills().ToolMeta("write_file"); !ok || meta.Source != ensemble.SourceDynamic || meta.EventSeq != seq || seq == 0 {
		t.Fatalf("dynamic provenance: %+v, expected %d", meta, seq)
	}
```

## Current pinned book/edition-3/chapter-09.md

Lines 437-442:

```text
The focused source test instead checks that the dynamically added `write_file`
remains declared after pending unload. Declaration and dispatch share the same
admission predicate, and restoring early removal makes that check fail. The
revised grader also exercises a dynamic write after marking unload. Match the
tool in the assertion to the capability whose lifetime is in question; another
successful paid run of the same read would prove no more.
```

Lines 452-454 (actual EOF):

```text
The agent can now acquire a manual and its tools without rewriting its
constitution. The next chapter gives the conversation a life beyond the
process that created it.
```

## Current pinned docs/edition-2-carryover.md

Lines 1109-1110 (actual EOF):

```text
grounded in actual outcomes; preserve the first-edition epilogue. A successful
regeneration is an experiment to demonstrate, not a result this draft assumes.
```

No range request beyond EOF is treated as additional content.
