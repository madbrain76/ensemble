# Chapter 6 student — formal correction round 1

Status: R1 correction and required local gates completed; candidate for independent review, not acceptance. Original candidate and all prior
evidence remain immutable. No new paid calls or credential/settings reads are
authorized or performed. Concurrent main publication work is out of scope.

## Supplied finding and assistance

Parent supplied `/tmp/cadence-ch06-correction-r1-task.txt` in full after joining
independent review final712. I did not read reviewer code, reports, or private
evidence. Four supplied paired cases identify one Gemini boundary/provenance
defect: opaque-first rejection; signed ordinary text merging with unsigned text;
two signed texts merging; signed thought losing its signature to an opaque part.
The original numerical grades do not discharge these obligations. This report
does not claim an observed real-provider rejection.

## Ownership and test plan (recorded before source edits)

Keep the sole Parse seam, shared plain decoder, actor/history and replay renderer
unchanged. Repair only Gemini's existing response-local accumulator. Explicit
parts-array positions delimit parts inside a frame. Across frames, only the first
part may continue the previous compatible unsigned content; a signed prior part
is closed. A signature-only first part of a later frame may complete an unsigned
prior text/thought/tool part, but must not overwrite a signature, absorb an
explicit later array element, or disappear without a previous content part.
Other opaque material is carried, not interpreted. No new IDs/API/controllers or
feature rows are needed.

Add public-facade tests through actual local HTTP responses. Cover all four
supplied cases in single-frame and split-frame layouts, exact hand-written
canonical values/order/source/signature bindings, final/delta correlation with
zero valid, parser reuse, and legitimate continuation. Feed the resulting
canonical response through the existing public History/Render path to verify
same-source replay and cross-model/vendor/surface filtering while ordinary signed
text stays visible. Exercise the real Ask path for ordinary-answer privacy.
First compile and run these tests against unchanged 800da21 production; preserve
behavioral failures, then apply only the accumulator repair and require passing
restored/corrected runs. Freeze new source/binaries in the NEW correction temp
directory, run required gates sequentially to natural completion, and record
exact commands, exits and byte identities. No prior live evidence is automatically
carried forward to changed Gemini bytes.

## Implemented ownership rule and direct evidence

Only `internal/engine/stream.go`'s existing Gemini accumulator changed in
production. The other changed file is the new public test. Explicit array entries
get separate blocks; first-entry continuation is restricted to compatible
unsigned content, and signature presence closes a binding. A standalone opaque
part gets its own block/response-local ID; it emits no private-data delta. The
existing shared decoder supplies canonical source provenance. The plain decoder,
replay renderer, actor/history, CLI and public APIs did not change.

`original-boundary-corrected.stdout.txt` is the meaningful old-code control:
compiled public tests fail on all four supplied cases, in one-frame and
split-frame layouts. They also expose lost ordinary Ask output for opaque-first
streaming. The exact same test bytes pass after the repair. The test constructs
explicit canonical expected values rather than trusting plain-vs-stream equality
alone. Assertions cover part count/order/type, exact synthetic signature binding,
From, usage, final IDs including zero, text/thinking/tool delta correlation, two
plain/stream passes through one Engine, same-source replay, cross-model/vendor/
surface filtering, and ordinary Ask privacy. Extra bounded continuity controls
cover trailing text, thought and tool signatures; explicit adjacent unsigned
texts; and explicit text then standalone opaque material.

The supplied old-code failures are the directly demonstrated parent-identified
bug. No redundant defective source variant was needed. Legitimate trailing
thought/tool signatures already passed on original production and continue to
pass; merging text after its trailing signature was separately caught.

## Retained student mistakes

1. Initial test compile referenced nonexistent Context.Turns; corrected to the
   actual public Context.Dialogue. This was setup, not a behavioral control.
2. First compiled test run omitted required positive MaxTokens in its Agent
   helper. Added512; retained failure, not counted as a behavioral control.
3. An edit script used an incorrect Gemini source anchor (initially misdescribed as non-unique; inspection showed a missing anchor) and refused before
   writing. Its shell still ran old code, producing the misleadingly named
   `corrected-boundary.*` failure receipts. Those are explicitly old-code results.
   The scoped repair was applied only afterward; `corrected-boundary-applied.*`
   is the actual successful repaired run.

No R1 final gate was weakened or skipped; no R1 lifetime probe was shortened or
interrupted. No private reviewer source/report was read. The supplied task and
parent milestone messages are the assistance boundary.

## Frozen R1 identity

- Commit `f595c787675e2a6343053ec42ea6c7e1e02a930e`
- Tree `6fb36168020a22b8993412812987338aa4b6e010`
- Source manifest `ch06-student-r1-source-manifest.json` (immutable freeze pins)
- candidate-source.tar: `0d6f30006d21ca1a2d128ede8f98c533b905a15740d6a0626f7061a19ad01cd8` — 686080 bytes, **local-only**
- core: `06c26681699a1cc67bfb1eb28280bd4c323e22923395e7714aefeb47edd5a0d4` — 9264162 bytes, **local-only**
- consumer: `16a155f246a703ad63bb510edeac6e4821ad9aa2f55bb6bd063d2e143b8dadc7` — 9194114 bytes, **local-only**
- workflow: `55df1c6a84413b104343955fd881ecbdad1c0435e42baac5406640b9ad888602` — 9229730 bytes, **local-only**
- shape: `d52f16f4174e72559b0bce9cd72d74878e38d5003a1f44ec2166e7c2b0df8451` — 9193970 bytes, **local-only**

All archives and executables are local-only regardless of whether their bytes
can be decoded as text. No original source manifest, review, archive, executable,
paid receipt or artifact index is overwritten. New temp artifacts are confined
to `/tmp/e3-ch06-correction-r1/`; new evidence uses only the R1 output paths.

## Live status

No R1 provider call or credential/settings read is authorized or performed. The
previous original Anthropic marker failure and the separately successful
remaining-only live run remain historical evidence for800da21, not fresh proof on
changed Gemini bytes. R1's real executables exercise local HTTP fixtures only.
No actual provider rejection is claimed from the supplied structural defect.

## Final-byte validation outcomes

Every command below completed on the clean frozen R1 source. Exact cwd, argv,
environment, timeout, elapsed duration, exits and output hashes are recorded in
`ch06-student-r1-evidence/final-command-receipts.json`. The original compiled
behavioral negatives and test setup/edit failures remain separate, unmodified
receipts. No gate was skipped or interrupted.

| Gate | Exit | Elapsed seconds |
|---|---:|---:|
| `root-test` | 0 | 128.681 |
| `root-race` | 0 | 140.744 |
| `root-vet` | 0 | 0.334 |
| `consumer-test` | 0 | 0.126 |
| `consumer-vet` | 0 | 0.122 |
| `workflow-test` | 0 | 0.662 |
| `workflow-vet` | 0 | 0.128 |
| `shape-test` | 0 | 0.125 |
| `shape-vet` | 0 | 0.124 |
| `focused-final` | 0 | 0.597 |
| `shape-fixture` | 0 | 1.153 |
| `grade-ch06` | 0 | 13.066 |
| `grade-ch01` | 0 | 124.187 |

Chapter6: `passed:true`, 100/100; Chapter1:
`passed:true`, 100/100. Both used the same frozen core,
with every flag before the executable. `paired-grade-command-receipt.json` carries
the full actual argv, exact model-map path/bytes/hash, source identity, executable
pin, stdout/stderr hashes, exits and passed booleans. The public wrapper's emitted
title still says Edition2; the invoked route is the required `grade-edition3`.
This grading result does not override the direct boundary/provenance obligations
or independent review. Both genuine62-second Chapter1 phases completed naturally.

The retained actual CLI gate is in the final focused tests; it builds its child
from the frozen clean source. The independent Gemini shape fixture executes the
exact frozen shape artifact. Streaming: first delta1.248ms, request157.890ms,
three finalized parts, two visible text deltas, answer `shape one two`. Plain:
streaming:false, first delta3.087ms, request4.809ms, the same answer and three
parts. Both modes retain known-zero usage. These timings are local HTTP fixture
observations, not provider latency or paid-provider evidence.

Only `internal/engine/stream.go` and `chapter6_gemini_boundary_test.go` differ from
800da21. The shared plain decoder and renderer hashes equal their original pins.
All474 entries in the original evidence inventory were re-hashed successfully,
including every original paid receipt and retained candidate1 artifact. The
original report/index/source manifest were not rewritten; their exact hashes
are in `original-boundary-preservation.json`. Concurrent unrelated publication
edits were not staged, reset, cleaned, or modified.

## Submission boundary and qualifications

R1 source/artifact pins remain exact and the worktree is clean. No tests, builds,
fixtures, graders, model requests or owned approval jobs are pending. Zero new
paid requests and zero credential/settings reads were performed. Historical live
results are explicitly NOT rerun on this changed candidate. No provider rejection
is asserted; the repaired bug is demonstrated by local canonical, provenance,
replay and visibility boundaries. No acceptance tag was created.

`ch06-student-r1-evidence-manifest.json` indexes the new safe evidence and explicit
new temp artifacts. Every archive, frozen executable and raw local runtime file
is LOCAL-ONLY regardless of UTF-8 decodability. It records path, bytes, SHA256 and
portability for every artifact, excludes only its own recursive self-hash, and
does not merge or overwrite the original inventory. The final report/index/source
manifest hashes accompany the student submission. Context is retained for the
parent's independent re-review and any bounded later correction finding.
