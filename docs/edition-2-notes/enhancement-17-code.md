# Edition 2 Chapter 17 grader checkpoint

Coder: GPT-6 Astra/high. Original source: Chapter 18, *The Invisible Invoice*,
revision `ca9d807d9e344375037d0c75bb29185192920aab`, SHA-256
`2224c498d28cb1ba2aaec8b08eb3e21153f80f816842b4898929522e3f173d11`.
Mandatory `.agents/skills/ensemble-coding/SKILL.md`, full carryover/workflow and
current enhancement/map were reread after compaction. Student and comparative
reviews and the original grader were read completely. This is a tested stage,
not final independent acceptance; exact samples and sensitivity work remain.

The separately versioned grader preserves all seven original weights, totaling
100, without recursive parity. It uses the actual optional GUI, a typed primary
with five declared real tools, public completed-turn observations and fresh
subscription meter snapshots. Subscribers remain alive until normal EOF; the
process must save and exit successfully. The GUI's actual default save is
cwd/save.json, despite the accepted Agent's data directory being cr. No original,
runtime or frozen-source file is changed.

Run `go run ./cmd/grade-edition2 -ch 17 -json`. Grader-only artifact options
`-cache-current` and `-cache-prior` default to `cr/cache-current.txt` and
`cr/cache-prior.txt`. `-cache-prices 2,2.5,0.2,10` supplies the accepted October 11
Sonnet-5 rates per million fresh input/cache write/cache read/output tokens.
These are dated fixture rates, not a current pricing recommendation or values
read back from the implementation under test. Legacy uses:

```sh
go run ./cmd/grade-edition2 -ch 17 -json -gui solutions/ch18 \
  -cache-current request.json -cache-prior prior_request.json \
  -cache-prices 3,3.75,0.3,15 solutions/ch18
```

The independently agreed legacy tariff tests its declared historical arithmetic,
without treating different prices as a calculation defect. No provider call is
paid and no active fixture selects Anthropic 5.5.

The stage reports under `solutions/edition-2/evidence/ch18/enhancement17/` record
actual accepted **100/100**, exit 0, 0.756 seconds, and legacy **85/100**, exit 1,
0.625 seconds. Earlier actual full 100/85 measurements remain separately bound.
Legacy omits the tool-definition breakpoint: actual captured requests have zero
tool markers, one system marker and advancing history markers. Original Chapter
18's TL;DR and vendor-asymmetry section explicitly require that tool breakpoint;
this is missing valid coverage in the original grader, not a new feature.

The initial local test decoded `response_ended` as the payload key instead of
actual `response`; that setup failure is retained in `initial-fixture.json`.
A follow-up capture run crashed legacy with `send on closed channel` after a
short-lived meter subscriber disconnected. Its actual 0/100 panic report is
`legacy-wire-report.json`, separate from the successful 85/100 baseline. Keeping
observing subscribers open avoids making that unrelated known GUI lifetime
failure a prerequisite to cache grading; normal EOF is still checked. The first
line-based dump parser also rejected valid pretty JSON; delimiter-based receipt
parsing corrected only that observation. `legacy-markers.json` covers only the
one request captured before the panic, while `legacy-held-wire.txt` contains the
successful full run's two-turn marker evidence.

The new checks require actual rotating request content, nonempty stable tools
and primary instructions, direct markers on the final tool/system/stable-history
blocks, exact durable and live disjoint usage buckets, cumulative price and hit
rate arithmetic, unknown-price reporting, and retained dialogue alongside zero
initial restart usage. Gemini's separate thought/output counters are exercised.
Artifact JSON is decoded for content comparison, allowing report indentation and
object-key order; this is not byte-offset diagnostic proof. The stable wire
sections themselves are compared without reserializing their content.

The original **0/100 on 7436fcb** remains historical: omitted skill type blocked
startup and masked all seven checks. Twelve successful real browser turns and
their actual grants, dated model identities and earlier-source provenance remain
in ch18 notes/review. No new cache grant, native DOM observation or invoice is
claimed. Main-only lens routing, detailed byte diagnostics, surviving-memory and
ephemera boundaries, OpenAI caching variants, per-model/tier attribution and
interrupted-worker prediction ownership retain their accepted source, mutation
and live evidence separately; this small numeric fixture does not certify them.

Stage focused race and vet checks pass; raw receipts are `stage-race.txt` and
`stage-vet.txt`. Reviewer preflight strengthened legal marker placement and retained
original spend/charge monetary-key checks before this checkpoint. Final sensitivity,
exact examples, preservation identities and independent disposition are pending.
