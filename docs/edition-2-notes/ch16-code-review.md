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
