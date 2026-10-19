# Chapter 2 restart: bounded TLDR clarification

The root's verified first-candidate finding was implementation noncompliance:
the candidate substituted `message_added`/`ephemeral_added` and string IDs for
the public frozen contract, omitted required Turn/context Usage and the four
usage buckets, and lacked file-based `render LOG`, `dump`/`CH02_LOG`, and distinct
replay-target configuration. It compiled but was held before tests. The chapter
already explicitly teaches these behaviors; neither body ambiguity nor defective
examples explain this failure. Architecture and owner packages were substantially
sound.

The root corrected its initial assumption that a TLDR already existed. This edit
adds `## TL;DR` after the historical opener, before §2.1, with one bullet: implement
the frozen event/content/JSONL vocabulary and exercise CLI as interoperability
contracts, including causal turn state, four-bucket cumulative usage, and the
ephemeral lifecycle, rather than substituting a smaller message-only protocol.
APIs marked as pseudocode remain flexible. All existing prose and architecture
remain unchanged; no schema, feature, or requirement was added.

Prose lint (`go run ./cmd/lintprose book/edition-3/chapter-02.md`) exited 0 both
before and after the edit. Soft warnings remain: negation forms increased from
22 to 23 for the explicit contract contrast; the unchanged 5,540-word person gap
moved from line 625 to 634. `git diff --check` passed for the chapter. This is a
pre-reexecution clarification, not the formal text review; that single review
remains after code experience. No candidate code, reviewer report, private probe,
or reference implementation was used as input.
