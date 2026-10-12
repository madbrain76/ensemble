# Edition 2 Chapter 22 author record

Author: `edition2_author`. Assigned October 12, 2026 after Chapter21 acceptance
`d62ab24c`, checkpoint `7e6dc5bc` / `edition-2-enhanced-ch21-r1`.
Status: coherent prose draft; exact examples, final revised grader and independent
combined review pending. No Chapter23 authorship is released by this record.
Mandatory technical skill: `.agents/skills/ensemble-coding/SKILL.md`.

## Reading and context

The author read both complete originals and all four complete original-numbered
student/comparative records before drafting: original24 1–300/301–565;
original25 1–300/301–582; ch24 notes315, review242; ch25 notes349, review243.
After compaction and before any Chapter22 writing, the author reloaded all four
required documents fully, without truncated output:

- Workflow2233: 1–350,351–700,701–1050,1051–1400,1401–1750,1751–2100,2101–2233.
- Carryover981: 1–340,341–680,681–981.
- Voice632: 1–340,341–632.
- Chapter procedure599: 1–320,321–599.

The full373-line coding skill was reread, alongside current AGENTS authority,
enhancement opening and the chapter map. Both original chapters were then reread
in full (24:1–290/291–565;25:1–300/301–582). All four full documents remained
active together during this initial draft. Any subsequent compaction requires
another full reload before writing or reconciliation.

Source bindings:

- Original24 revision `2ca63636c8ba592ea0ed190fb0a3c4731e2b0848`, SHA-256
  `c6ac58046ed8e5496c44931e20d126d36081325fc06bfc27f4276c73424976b1`.
- Original25 revision `009d74c70b108681b2107b4c3a08685508567cee`, SHA-256
  `eb09ec25fe8a681479310ffe22f4cf19d0b4eba34bf507dbbfd8181f7302616f`.
- Accepted combined source `e00d5bcaf1412c50787428ef1fa08bed59b99e11`, main tree
  `30fb05d154b9808eb88622aaff9c1424cd7e8af1`, direct `solutions/edition-2/ch22`.
- Intermediate original24 accepted source `0430df4340fffac5234a94aea843f9a7c9d02282`,
  tree `8555e3cb9ab5200f0a1c52691007a7fc4a1c6797` remains separate.

Read actual public Ensemble/fleet construction, internal fleet and shared/tool
interfaces, terminal construction/reply handling, actor completion correction,
GUI construction/observation/assets/relay and external reuse/browser fixture.
The existing root supplies construction and observations; the actual Agent owns
state. No callback-bag design or new runtime is proposed.

## Through-line, outline and preservation

Stake: Bill wants reusable views and more hands without losing the ability to
steer; the deployed application must connect the public boundaries its library
already exposes. Resolution: an actual outside GUI consumer and a terminal
fleet that can answer questions, preserve a running join and submit tested work,
with the observed failure to obey a delivered hint retained.

The merged outline follows reusable GUI → actual root/Agent construction →
embedded/components/relay proof → two kinds of fleet silence → spawning and
identity → asking versus submission → message wake → bounded correction →
historical product wiring failures → runnable combined exercise → actual runs.
Both assignments appear explicitly in TL;DR and exercise. The fourteen original
rubric rows retain GUI100 + fleet100, pending final predicate reconciliation.

Preserved original material:

- Homebrew-VTT's desired Artifact component and the inaccessible-public-boundary
  irony; the historical two paths to Gorilla and born-unused exports.
- Current-model getter motivation, binary-without-its-interface deployment,
  positive/negative import controls, outside event vocabulary and public relay.
- Build-tool wrapper holding inherited pipes, and affected harness paths after
  a package move. These are historical reports, not new observed failures.
- Bill's listening/steering motivation; delegating mistaken for completion;
  combined idle/wait hazard; full Agent children and filesystem ownership tree.
- User as an address rather than a managed child; why asking and submitting must
  differ; fan-in yielding a child question; background job lifetime; three wakes.
- All five historical executable wiring failures, plus the package-counter and
  misleading relative-path bisect stories, scoped to the reported tests.

There is no length target. Replaced material consists chiefly of the superseded
refactor instructions/AgentHooks design, unsupported code-size inventories,
universal success/security claims and duplicated original check tables. Useful
explanations remain; a new complete public consumer and terminal setup make both
halves reproducible. The original source-relative diff, not the word count,
governs the preservation review.

## Teaching corrections and their evidence

1. `go list -deps` measures import reachability, not every linked machine-code
   byte. Internal-package restrictions really apply to an outside consumer;
   importing internal/llm is no escape hatch. Keyed struct literals can omit
   fields. Official Go pages checked October12:
   https://pkg.go.dev/cmd/go#hdr-Internal_Directories and
   https://go.dev/ref/spec#Composite_literals. Existing carryover records all
   three corrections. Actual public root/Agent construction replaces AgentHooks.
2. Embedded contents are read-only, but an exported variable binding is not
   immutable. Accepted gui/assets.go uses a private anchor and Assets() fs.FS;
   https://pkg.go.dev/embed checked October12. Explicit relative overrides still
   depend on cwd; the chapter recommends an absolute override for relocation.
3. An asset HTTP200 can be a directory listing. The original24 student sensitivity
   correction required actual shell identity; this positive-result lesson stays
   in prose. Component-only serving is distinct from actual Chrome rendering.
4. Public relay uses real listeners, original caller IDs and actual browser
   forwarding. No-browser error, connected-silent wait and client/browser/server
   closure have different conditions. Accepted gui/relay.go and its outside
   reuse tests supply the mechanism, without imposing a new timeout.
5. Routing IDs do not replace owner interfaces. Actual factory receives parent
   Agent and spec, descendants reach the shared root, and state directories are
   separate while task Workspace remains shared. Chapter21 narrowing remains.
6. Question IDs belong to the recipient pending map. Same-name parent/child
   overwrote a human question before the accepted correction; omission of a reply
   ID is allowed only when selection is unique. Accepted ch25 review/probes bind it.
7. Message wake is independent of callback delay. tool_limits remains the prior
   one-shot next-call mechanism, not the original25's unsupported per-name map.
   Preserve running jobs and later collection under the same handle.
8. Factory installation enables library correction for every fleet root/child.
   Earlier single-Agent embedders without a factory keep normal response completion.
   Exactly three framework wakes precede escalation; successful submit resets.
   Exhausted child state must be failed, with no invented result. Failed turns
   do not enter the correction loop. Actual actor finish and fleet Idle bind this.
9. Remove universals that every chapter/mutant passed, every failure mode is
   prevented, or a broad six-row failure proves no narrower mutation exists.
   Preserve the original real wiring lessons without adding their old repair
   sequence to the new assignment.

## Historical results and live provenance

Original24 grade remains0/100 on its recorded launch/import assumptions;
original25 grade remains0/100 before inference with no supplied model. Accepted
exceptions did not turn those into passes. Neither score describes the revised
combined grader. The original24 source review accepted three narrow fixture
exceptions; original25 accepted missing-model and factory-enabled compatibility
exceptions, both at initial review, zero formal rounds.

Directly read coordinator live receipts under evidence/ch24 and ch25. Original24
live source `35df549210130f824fdc2a4eb632ea878b2c529b`, tree
`5b36f28cbc8bac320febe1a9afa810efd27452e3`:9.649seconds, exit0,
three Astra/high requests/completions, seven assets, actual42-byte file and
matching browser DOM. A large combined live-results read truncated; only its
fully read coordinator receipt supplies the published aggregate claims. Component,
relay and switching proof retain their separate actual test/live sources.

Fleet live source `891d513bf142013d6888760cfae4fcaab1edd680`, tree
`6cbf18933f82ee450d91750442161b176e469f8e`:48.198seconds, exit0,
20Astra/high requests/completions, children add/multiply/subtract with saved
conversations and successful independent combined assertions. Join job12 starts
18:30:55.751093Z and returns running18:30:55.887642Z after steering; the same job's
completed result18:31:10.340664Z contains all three submissions. Request15 includes
the steering text; final submit omits STEER_COBALT. Delivery is demonstrated,
compliance fails. The malformed initial child shell command and corrected success
are retained; the concurrent wire trace is not claimed to pair requests by ordinal.
Final accepted production matches that live implementation; later README/test
changes do not relabel the paid run as a final-source rerun.

## Initial validation and pending gate

Direct hard prose lint passed on the coherent draft (6527 prose words);
soft person-gap warning is considered without inventing a story. The chapter
will be reread and checked after final sample/rubric reconciliation. No runtime,
frozen export, first-edition or concurrent main/crossover file has been edited.
Exact new public consumer, embedded/disk/component/relay guidance and typed fleet
launch are being handed to the grader coder for technical validation against
external peers. No new paid run or autonomous-model quality result is claimed.
Independent whole-chapter review and the final combined grader freeze remain.
