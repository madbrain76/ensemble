# Original Chapter 16 / Edition 2 Chapter 15 — independent review

## October 11 preparation: original forcing and provider usage

Reviewer `enhancement_reviewer_resumed`. This is Bill-authorized, narrow
preimplementation inspection of the original answer, **not student review or
chapter acceptance**. Bill resolved the manuscript's cumulative-usage wording:
forcing uses total input of the prior request, including correct provider cache
accounting. The coordinator records the student-safe ruling and current official
provider sources in [ch16.md](ch16.md). Keep the original implementation details
below out of the student's handoff.

Read the entire mandatory `.agents/skills/ensemble-coding/SKILL.md`, carryover and
workflow. Inspected Anthropic and OpenAI first, then Gemini, as Bill requested.
No product, original source or grader changes; no tests, grader or paid calls.

**Bill's recollection is correct.** Original `internal/llm/forcing.go:85–100`
explicitly distinguishes the last request from running totals. It scans backward
to the latest `ResponseEnded` carrying a response and uses that response's
`Input + CacheRead + CacheWrite`, excluding output. This is provider-reported
input from the latest completed response, not a newly estimated context size.
The separate running total is accumulated by `internal/common/context.go:200–205`
and is not the forcing numerator. `forcing.go:38–51` divides the selected count
by the current model's context window; absent window or zero selected usage
returns without forcing. Missing usage is therefore not independently recognized
as unknown in this implementation.

The original test `forcing_test.go:101–117` records 3,932 input tokens and then
512, asserting that tools return after the latter. Their sum still exceeds the
4,096-token window, so the test's expectation supports the last-request reading.
This is source inspection of that test, not a newly executed result.

### Exact parser accounting

All paths ultimately emit canonical usage in `ResponseEnded`. The common usage
type documents disjoint input/cache-write/cache-read/output buckets
(`internal/common/event.go:315–332`).

| Provider | Original plain-response normalization and resulting forcing input | Streaming behavior |
| --- | --- | --- |
| Anthropic | `claude.go:375–380,448–454`: Input=`input_tokens`, CacheWrite=`cache_creation_input_tokens`, CacheRead=`cache_read_input_tokens`. Forcing adds those three disjoint fields. Neither output nor any nested cache breakdown is added. | `claude.go:587–595` takes input/cache counts from `message_start`; `640–644` changes only output from `message_delta`; `666–670` emits the final usage. Input/cache updates appearing later would not be incorporated by this parser. |
| OpenAI Chat Completions | `openai.go:281–288,344–356`: Input=`max(prompt_tokens−cached_tokens−cache_write_tokens,0)`; the latter two fields come from `prompt_tokens_details` and become the two cache buckets. With consistent subset counts, forcing's sum is exactly `prompt_tokens`, not prompt plus cached tokens. An absent optional cache-write field decodes to zero. | `openai.go:431–438` replaces local usage whenever a usage object arrives, including a no-choices usage chunk; `498–500` emits it. It does not add usage chunks together. |
| Gemini | `gemini.go:321–326,391–402`: Input=`max(promptTokenCount−cachedContentTokenCount,0)`, CacheRead=`cachedContentTokenCount`, CacheWrite=0. With consistent subset counts, forcing's sum is exactly `promptTokenCount`. Output is candidates plus thoughts, but is excluded from forcing. | `gemini.go:479–483` replaces local usage with each nonzero `usageMetadata`; `566–568` emits it. The comment's “cumulative” means within this response's stream, not lifetime accumulation across requests. |

The clamps in OpenAI/Gemini mean contradictory subset counts exceeding the
reported prompt count could make their canonical sum exceed that prompt count.
That is an algebraic limitation, not evidence that a live provider supplied such
data. Likewise, original comments about historical live cache fields are source
claims, not fresh verification of today's schema. The coordinator's official
documentation check must distinguish those claims from currently documented
fields, especially OpenAI `cache_write_tokens` and Anthropic stream updates.

The narrow contract conclusion is: use the previous completed request's total
input; add Anthropic's disjoint cache buckets, but do not add already-included
cached subsets again to OpenAI/Gemini prompt totals. A cached token still occupies
context. This confirms Bill's ruling without supplying an implementation recipe.

### Provenance

Inspected from repository HEAD `1878e64` with no modifications under original
`solutions/ch16`. Its latest directory revision is
`37697fe3f133cb3f85d86e7475375a654c23cb3d`. Individual last-changing revisions and
SHA-256 values below identify the exact files; paths are relative to that directory.

| File | Last-changing revision | SHA-256 |
| --- | --- | --- |
| `internal/llm/forcing.go` | `050f9c14d0c7328cff3b898a1d66c322edfc6589` | `fc72060d2bd4244264300f29d537f509cbb9ccaf67cdc3927d7ac3c32acc8e5d` |
| `internal/llm/forcing_test.go` | `050f9c14d0c7328cff3b898a1d66c322edfc6589` | `2a1fe886bb8ca70801e21134b092e32d6de33be9a6a418bc246c6507edf4910b` |
| `internal/llm/engine.go` | `37697fe3f133cb3f85d86e7475375a654c23cb3d` | `358c6e448f7832cd7798889cdedcc3e435916153d23c6bc44a6ac904ff5993d3` |
| `internal/common/context.go` | `dbe4a56a445d5724ff2fc9bd3281b57686fde5cd` | `2c2663c0b244cd603b83ec11351467841e4f92d7c52cf1917fbde1c2b49d9af0` |
| `internal/common/event.go` | `050f9c14d0c7328cff3b898a1d66c322edfc6589` | `38b24f31c39ce85f48441985c2638e848dadb8b1c0b106706c707cbf5202ad20` |
| `internal/llm/claude.go` | `050f9c14d0c7328cff3b898a1d66c322edfc6589` | `25da4f9940e9fb39e48a14ac0b8a5ff5c7fb2d73b4b250200b91f72f82a640af` |
| `internal/llm/openai.go` | `050f9c14d0c7328cff3b898a1d66c322edfc6589` | `92d68fed485f740dfa5335982bf839d2cd4825f0c116e00a8362c38da2e0900a` |
| `internal/llm/gemini.go` | `050f9c14d0c7328cff3b898a1d66c322edfc6589` | `532dca26aab706f7beef3f6abf270dfc2cd791b8c1b3df37a7153651739c0de9` |

A narrowly scoped later-file comparison found `forcing.go` byte-identical in
original Chapters 17, 18 and 22 (same SHA-256 above). The distinction was already
correct in Chapter 16; it was not introduced by those later snapshots. This is
not a review of their other implementation or provider behavior. No formal
student review round has begun.

## October 11 grader preflight — intermediate source, no acceptance

Reviewer reread the mandatory skill and the current student checkpoint in
[ch16.md](ch16.md). This assessment uses committed student source
`ccd63509c3b2e08ff5a231b0409423b236953022` (tree
`58dd0d51c9eb34e1096a1134fa2d045c97d07c97`), not the coder's changing working tree.
Final provider, mutation and live evidence remains pending. No formal review
round or chapter acceptance begins here; no grader/live call was repeated.

The original grader completed **20/100 FAIL**, grader exit1 / make exit2, on
that intermediate source. The committed [provenance](../../solutions/edition-2/evidence/ch16/original-grade-provenance.json)
and [output](../../solutions/edition-2/evidence/ch16/original-grade.txt) are from
`bd4213bf0d52adafc5d9bc6b5112416201da7aea`; output SHA-256 is
`6be427fd976941c65c215f09d8332ec0eba2fe83ef03a0a57132d422ce06ca5d`.

**The settings-path exception is independently justified, narrowly.** Original
`ch16_harness.go:50–72` writes `cwd/settings.json`, launches the GUI in that cwd
and supplies no Agent-directory override. The student's GUI command at
`gui/cmd/main.go:77` selects `DataDir="cr"`; its existing owner-based settings
path is `gui/settings.go:91–92`, and the loader applies memory settings to the
public Agent at lines85–88. The grader's small conversation/band budgets and
disabled flags therefore never reach the Agent. Its normal conversation budget
is 3,000 bytes (`internal/memory/memory.go:16–39`), rather than the fixture's 300.
The output consistently reports zero compressors/memories. This is the same
pre-existing ownership mismatch as original15, not grounds to relocate settings
to global cwd or waive compression behavior. The chapter expressly leaves disk
layout to the student; Bill's Agent ownership remains binding.

The proposed exception covers the **fixture location**, and hence the following
eight failed check outcomes (77 points), not their valid behavioral obligations:

| Check (actual earned/possible) | Still requires independent supplemental proof on final source |
| --- | --- |
| `micro-handoff-compresses` 0/15 | A real checkpoint crossing an effective budget compresses oldest complete segments and removes replaced raw work. |
| `planted-fact-survives` 0/10 | The actual compressor receives the selected work; its returned memory reaches later provider context. |
| `graduation-fires-oldest-first` 0/12 | Crossing effective watermarks folds oldest files, publishes their replacement, removes those originals from active context, and leaves newer work. |
| `disable-enable-idempotent` 0/12 | Disabled bands leave context and stop receiving compression; re-enable rescans edited disk files, preserves canonical order, and ordinary restart respects disk edits. |
| `disabled-neighbor-refused` 0/8 | An actually full source refuses a disabled destination without losing old content, with the required warning. |
| `replay-needs-no-llm` 0/10 | Actual recorded compactions rebuild memory without another compressor call. |
| `memory-is-data` 0/5 | First-person compressor instruction and actual conversation memory, separate from the frozen system prefix. |
| `fresh-start-populates` 0/5 | Real disk memories populate before the first request with no prior save/log needed. |

Each fails before its later meaningful assertions. In particular, “no memory”
does not establish a replay, ordering or data-placement defect; neither does a
fixture exception establish those behaviors passed. Final review must assess
the student's real-runtime tests, captured effects, targeted mutations and live
quality, without treating names of passing tests as proof.

**The two awarded checks need different interpretations.** `forced-handoff`
12/12 actually presents high reported input and observes only `micro_handoff`
on the later request (`ch16_run.go:740–803`). Its fixture does not establish
Bill's latest-versus-lifetime distinction, provider cache arithmetic, 90% warning
or eventual restoration; those remain separate obligations. `abandon-on-restart`
8/8 (`863–917`) does not require a compressor to run or return its scripted500.
With no compression, old work remains and zero copies satisfy its at-most-one
memory test. It also does not launch a SIGKILL case. Thus the actual awarded
eight points are retained, but this run proves neither failed/abandoned-fold
recovery nor crash safety. Supplemental proof needs an actually launched fold,
its failure/interruption, retained source information and nonduplicating recovery.

**Inherited parity remains 0/3**, with exactly the four reported original15
failures: `skill-survives-the-ladder`, `ladder-is-recorded`,
`replay-equals-snapshot`, and `frozen-prefix`. Reinspection confirms the first
three retain the cwd-settings mismatch (`ch15_harness.go:366–369`) and stop on
missing cuts/uncut content before their later comparisons. The fourth passes
the preceding system/tools prefix comparisons, then stops at undefined `think`
before the MCP declaration predicate (`ch15_run.go:500–549`). The independently
accepted [original15 exceptions](ch15-code-review.md) continue to explain those
specific failures. They do not waive memory's interactions with recorded cuts,
skills, declarations, journal recovery or request replay; reuse unchanged valid
evidence and check affected behavior at the final gate.

Source provenance: original `ch16_harness.go` and `ch16_run.go` still match the
revisions/hashes recorded in the student-safe launch erratum. `ch16_checks.go`
is revision `174eb9778243914db251aa99f4ef9d61889ce5cd`, SHA-256
`4e99b4427dc115ced53c92747e73b3299de5f1995122a51fb226949ec9497521`.
Original15 harness/run are respectively `107c59d45ee2977d82910f1c1e734eabb73d538a`
and `374916a10263c4b32cc3dc22999043f07202f95c` (unchanged original files).
No source correction is requested by this preflight. Final comparative review,
complete student response and authentic provider evidence are still required.
