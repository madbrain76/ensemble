# Original Chapter18 / Edition2 Chapter17 — independent review

## October11 public-contract preparation — no implementation verdict

Reviewer `enhancement_reviewer_resumed` read all665 lines of original
`book/chapter-18.md`, revision `ca9d807d9e344375037d0c75bb29185192920aab`,
SHA-256 `2224c498d28cb1ba2aaec8b08eb3e21153f80f816842b4898929522e3f173d11`,
and the complete current [student-safe preparation](ch18.md). Read the mandatory
`.agents/skills/ensemble-coding/SKILL.md` and current governing documents; the
student extends accepted original17 tree
`83e2f4cc85eb91ef9f62a54f638baa4d538cacf8`. No original answer, grader or changing
implementation was inspected, and no test/provider call was made for this prep.

Scope is the actual-request cache lens, supported vendor breakpoints and visible
session usage/cost meter. Independent review will distinguish bytes from tokens,
stable payload from moving marker metadata, and fake grants from actual provider
cache reads. Main conversation requests belong in the lens; stateless judge and
compression calls must not create false conversation-prefix alarms. Genuine JSON
array growth and Gemini's `contents` name need the same correct classification as
`messages`. The byte-measurement warning prohibits silently substituting remarshal
sizes for the measured artifact; it does not prohibit parsing to classify shape.

Existing carryover resolves the source's obsolete memory-in-system assumption:
do not move mutable data into the prompt. Identify the actual surviving compaction
boundary in the accepted implementation. The stable end excludes ephemeral content
even when rendering merges it into a message with retained content. Marker limits
are budgets, not reasons to invent missing tool/system/handoff content. OpenAI's
explicit opt-in depends on a marker actually landing. These are chapter obligations,
not a prescribed implementation or additional cache protocol.

The meter must retain the distinction between session totals, latest response,
restored history and the effective model. Counts are durable facts, money is
derived, unknown pricing is visibly unknown. The coordinator's dated provider
references already cover model support, current field names, thresholds and pricing
tiers; this preparation does not repeat that research or independently claim live
grant. Historical numeric rates, warm-up patterns and precise provider lookback
counts are source observations, not guaranteed acceptance targets. No new product
ambiguity presently needs Bill's decision.

The public chapter supplies the conceptual lens methods and visible meter fields
but not a concrete external grading launch, artifact-path/schema or GUI message
envelope. If those facts are needed, the coordinator may recover only the missing
invocation/observation contract under the standing authorization and record exactly
what was supplied. It must not impose the original implementation's method names,
single-module layout or private source paths. A mismatch remains a review question,
not a reason to invent compatibility machinery before seeing an actual result.

Await frozen implementation, original grade and completed student evidence for
formal comparative review. Zero review rounds or acceptance claims begin here.
No future chapter or authorship is released; original22/24 stops remain. Eventual
accepted export is direct Edition2 `ch17`, owned by the coordinator.

## October11 original grader preflight — first green, not acceptance

Independently read complete unchanged `ch18_harness.go`, `ch18_run.go` and
`ch18_checks.go`, plus the relevant frozen launch/skill/lens/usage/save paths at
`7436fcb35ef54e2996114a2d3fe68fae4cfac82e`, main tree
`05e0d2e92a9f70f5d2433af76c1a3cbf5208b493`. The actual completed original result
is **0/100 FAIL**, grader exit1/make exit2, 13:34:32–13:34:54 UTC. See
[output](../../solutions/edition-2/evidence/ch18/original-grade.txt) and
[provenance](../../solutions/edition-2/evidence/ch18/original-grade-provenance.json),
output SHA-256 `82053233f0bfdee03295053fdbb7b400df3f1406bc056db4de4e26e1e3a3dbe1`.
No reviewer grader or behavioral probe ran. Final student evidence is pending.

The primary fixture at run.go:80–97 omits `type`; inherited parse.go:96–99
requires primary/loadable/dependency, as original10 documents. `root.NewAgent`
fails before GUI startup (gui/cmd/main.go:102–104). The prefix launch assigns
that failure to six checks, then the vendor scenario fails on Anthropic and
returns before OpenAI/Gemini. There were two attempted launches, not seven
independent failing implementations. There is no recursive parity in Chapter18.
The missing-type fixture exception is independently justified; all seven valid
behavioral obligations remain unmeasured by this result.

Three unreached assumptions need separate treatment:

- The collector only reads cwd `request.json`/`prior_request.json`
  (harness.go, ch18ReadLensFiles), whereas the real Engine-owned lens writes
  `cache-current.txt`/`cache-prior.txt` below its Agent DataDir
  (cachelens/lens.go:139–148). GUI selects DataDir `cr`. The chapter explicitly
  permits artifact framing choices and requires inspectable consecutive content,
  not these filenames; ownership is binding. This is a latent naming/path mismatch.
- Session check run.go:357 requires cwd `save.json`. The Agent's inherited save
  path correctly uses DataDir/save.json, so a real saved `cr/save.json` would be
  missed. Correct session reset must still be proven against actual restored
  history; a path mismatch does not establish the meter's behavior.
- The cost check requires `claude-sonnet-5` to report priced/positive cost.
  Frozen student pricing currently covers the three independently verified live
  rows and displays other prices as unknown. This is not an observed failure or
  automatically accepted exception: final review must assess the student's scope
  explanation and actual pricing evidence. Inventing a price to satisfy a fake is
  inappropriate; a truthful unknown display alone is not proof of correct pricing.

The fixture's undefined `think` grant is **not** another observed startup failure.
No Chapter18 fake calls it, and declarations iterate implemented tools and their
visibility, rather than demanding that every granted name have an implementation
(agent.go:74–81). Do not invent a think tool or imply its missing contract caused
the seven failures. Valid nonempty tool-prefix checks can use real existing tools.

| Original check | Still-required independent evidence |
| --- | --- |
| cache-lens-exists,15 | Actual main-request rotation, readable sections and truthful byte measurement; marker stripping, append versus edit, and auxiliary-call exclusion. Mere file presence does not prove these. |
| prefix-is-stable,20 | Independently compared rendered tool/system bytes across real turns and correct change detection. Original predicate skips missing sections; do not infer coverage from an empty prefix. |
| request-has-breakpoints,15 | Actual supported vendor markers, advancing stable end, surviving compaction boundary, tool/system placement, budget and ephemera exclusion; OpenAI opt-in only when a marker lands, implicit routes preserved. Original check tests only Anthropic system presence/history advancement. |
| cost-is-computed-not-stored,15 | Durable integer counts without money; GUI derives correct nonzero price, distinguishes unknown, uses effective model and appropriate request tier. Original greater-than-zero check alone cannot prove arithmetic. |
| usage-is-session-scoped,10 | Real restart retains history but clears process spend/totals; latest-response values and session tooltip/rate remain distinct. Original check primarily compares latest input and does not establish cost/rate reset. |
| cache-read-rate,10 | Exact normalized provider reads, correct denominator and visible reporting, plus bounded live cache grant. Fake cold/warm values cannot establish provider eligibility or successful caching. |
| all-vendors-report-usage,15 | Actual parser paths preserve disjoint input/write/read/output counts, including streaming usage; absent provider write fields are not invented. Only the first failed Anthropic launch occurred in this original run. |

Harness/run hashes match the coordinator's public erratum. Checks revision
`9bd6b14a033db347212bc7f1aae8a837816bcced`, SHA-256
`00ffe9ddfb7745c71325efd0230d14ebb9232b49fb80a3a678a682bc689b9bb7`.
The original stripper itself recognizes only a few exact Anthropic spellings;
its presence is not broad artifact-format coverage. Preserve valid checks and
source expectations for later grader revision, without changing this grader now.
Await frozen final student review, actual live/browser/mutation evidence and
comparative source assessment. Zero formal correction rounds; no acceptance yet.
Mandatory handoff: `.agents/skills/ensemble-coding/SKILL.md`.
