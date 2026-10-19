# Chapter 10 retention amendment — incomplete independent delta review

## Verdict

**INCOMPLETE — not approved; no student-product defect established.** The reviewer hit a setup failure before launching any behavioral control. The governing finite task requires stopping and sealing on setup failure. No correction round is charged to the student (zero formal rounds used).

Frozen candidate: `b03c28fdb1586f16dfb53849f8c763576ebb39ae`  
Tree: `b2f86ba9fd48ff9a923a8a924a9612549aa762c0`  
Accepted base: `729bc919e734405307b4ac0060e65c40e33652e5`

## What happened

The reviewer verified governing/input/protection pins, read the task and delta, and extracted the candidate with `git archive` into `/tmp/ch10-retention-independent-review-20261019/source`. During setup, a Python script incorrectly assumed that the private input manifest had a `candidate_tree` key. It raised `KeyError` before writing the intended copy manifest or test plan. This is a **reviewer harness error**, not a failed assertion against the implementation.

A subsequent custody-only script also incorrectly assumed a `tree` key and failed. Both raw failure outputs are preserved. Inspection then confirmed that the manifest pins the commit and per-file hashes; the tree pin lives in the task. Final custody verification uses the actual fields and task tree pin. This verification is not a fixture repair or a behavioral rerun.

Earlier read-only inspection mistakes (two nonexistent assignment filenames and mistyped revision strings), clipped/recovered reads, and missing `--no-pager` on some Git inspection calls are disclosed in the evidence ledger. No failing output was replaced by a claimed passing test. The required exact test command/environment/timing records do not exist because no tests launched; this report does not claim otherwise.

## Scope and unrun obligations

The exact 18-file delta was inspected, with persistence/retention/settings/lifecycle ownership paths and fixture declarations. Inspection did not mature into a complete independent contract verdict. In particular, **all new behavioral evidence remains unrun**:

- Focused root retention controls under `-race`.
- Focused GUI retention/settings controls under `-race`.
- Any small independent valid-input probe.
- A green baseline followed by one meaningful compile-valid negative.
- Exact mutation restoration and the focused restored positive.

No Go test, vet, benchmark, original Chapter 10 lifetime/grade/whole suite, browser, live provider, external HTTP, or child agent ran. No source mutation or test injection occurred, so there is no negative/restored result to claim. Preliminary commentary about a different selection was not a compliant sealed run plan and was never executed.

## Carried evidence and boundaries

The 34 portable student-evidence members and their input-pinned inventory/source manifest/parent verification were hash/mode checked. Their root race (excluding `TestRequestLifetime`), root vet, GUI race/vet, and student diagnostic-bypass negative/restored receipts remain **student evidence**, not independently rerun acceptance. Old Chapter 10 numerical/live acceptance is historical and not new-binary proof.

The student report's earlier unintended lifetime execution, incorrectly directed GUI invocation, missing wrapper completion footer, and setup/compile mistakes remain disclosed in their unchanged evidence. They were not used as new product findings or excuses to rerun suites.

Final custody checks establish all 142 candidate members and all 142 owned-copy bytes/modes match the pins; candidate HEAD/tree/parent and clean status match the task. Pinned governing/current/protected inputs and the student seal remain unchanged. The shared main worktree had unrelated pre-existing changes; no blanket clean-main claim is made. No frozen source, guide, chapter, or paused Chapter 11 content was edited. No commit/tag/acceptance marker was created.

## Seal and handoff

Evidence: `ch10-retention-code-review-evidence/inventory.json` (self-excluded, report included), `preservation.json`, `read-ledger.json`, `execution-ledger.json`, and the two raw failure outputs. The small sealing script is included for custody reproducibility; it is not a test harness. UTF-8/no-NUL and closed-world membership were checked when sealing.

No owned jobs remain active. No cleanup was performed; the exact owned source and unused private HOME/TMP directories remain for custody. A parent-authorized follow-up review is needed to complete the unrun obligations; this sealed incomplete result must not be treated as approval or a request for student production changes. Root retains acceptance, prose review, correction-round ownership, and later Chapter 11 integration.
