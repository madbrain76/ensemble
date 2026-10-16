# Edition 3 student inputs

Read this with the shared [execution procedure](edition-3-execution.md), the complete mandatory `.agents/skills/ensemble-coding/SKILL.md`, and the assigned chapter before starting and after compaction. These current paths/phases replace historical edition instructions. Architecture, KISS, testing, documentation and review obligations remain.

## Permitted inputs

- The assigned `book/edition-3/chapter-NN.md`, Chapter 00 for orientation, and earlier chapters as needed.
- Your own working tree and accepted predecessor history in `/Users/bill/projects/ensemble-edition-3-student`.
- The current role, shared procedure, this public brief, and narrow public-contract clarifications explicitly supplied by the orchestrator. Mandatory governing instructions may contain incidental historical path/revision references; reading those instructions is allowed, but following those pointers into old solutions or repair evidence is not.
- Public dependency/provider documentation needed for the assigned exercise.

Do not read `agent/`, `solutions/`, `edition-1/`, future chapters, grader implementation, historical answer/repair notes, or their Git blobs to solve the exercise. Do not read the coordinator's readiness/carryover/evidence records unless a specific student-safe excerpt is supplied. The historical skill's direction to read Edition 2 carryover is replaced for the fresh Edition 3 student by this brief; authors/reviewers still read their complete governing documents. Do not follow answer-revealing links from permitted documents. These are input-separation rules, not a claim of OS sandboxing.

## Chapter boundary

Chapter 1 starts from an empty source tree. Later chapters extend your exact accepted predecessor; do not copy the full reference. The orchestrator supplies the assigned chapter revision, accepted predecessor and grader revision. No later capabilities early. Implement required behavior with the simplest adequate design, with useful explanations of non-obvious choices.

Run the assigned chapter grader with an explicit positional target. Chapter 1 accepts a supported solution directory or a built CLI executable. If your public root is a library with a nested CLI, follow the chapter's “Grade yourself” build instructions rather than grading the library directory. For the executable, from `/Users/bill/projects/ensemble` (the repository-root Go module, not `agent/`):

```sh
go run ./cmd/grade-edition3 -ch 1 -json /absolute/path/to/your/student-cli
```

Substitute the executable you actually built; its package layout is your implementation choice.

For Chapter 2, follow its revised exercise and run both existing graders against
that same built executable from the book repository root:

```sh
go run ./cmd/grade-edition2 -ch 2 -json /absolute/path/to/your/student-cli
go run ./cmd/grade-edition3 -ch 1 -json /absolute/path/to/your/student-cli
```

Both must pass and exit successfully. The second command preserves inherited
Chapter 1 obligations; it does not replace the Chapter 2 pass.

The Edition 3 command requires an explicit target and currently supports Chapter 1 only; it preserves the inherited conversation checks and adds delayed-header/body lifetime gates. Check the report's `passed` field and exit status, not score alone: required zero-point checks can fail a 100/100 score. Those real delayed probes take about two minutes, so a 90-second suite timeout is invalid. Omitting the path from the older Edition 2 command silently targets an old solution and is not a valid student result. A grader result does not replace required live exercises, test-sensitivity evidence or independent acceptance. Credentials remain in the authorized provider loader/test process, never prompts, command arguments, evidence or commits.

Submit one concise chapter review identifying the chapter/code/grader revisions, what you built, actual tests and live observations, failures, ambiguity, assistance and useful teaching clarifications. You may explain a suspected grader defect; you may not silently redefine success or edit the grader. The orchestrator supplies any missing public contract without old answer code. Follow the shared bounded review loop.
