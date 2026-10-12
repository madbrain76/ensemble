# Edition 2 solutions

Accepted cumulative solutions live directly in `ch01/`, `ch02/`, and so on.
Chapters **1–24 are accepted** (24 is optional, ungraded and self-reviewed
under Bill's explicit instruction). [`main/`](main/) is the authoritative working
tree. Original24+25 are accepted together as Edition2 Chapter22, The Ensemble.
Original26 was complete at that review boundary and is now independently accepted
as direct Edition2 Chapter23 after one correction round. The 26 original chapters
therefore have 23 matching Edition2 solution directories. The new Chapter24
adds the temporary Codex endpoint and has no first-edition counterpart; Gateway
is still deferred.
Subsequent authorized coding uses GPT-6 Astra/high. Auto-Recall follows
Bill’s October 11 quota/corpus decisions, recorded in the chapter notes.
Do not put new solutions in `history` or create symlinks.

Use older supported Anthropic models such as Sonnet-5; 5.5 models are excluded
from current supported use under Bill's ephemeral-context restriction. Current
`main/` and revised examples follow this rule. Frozen snapshots and past results
retain their original model tables and identities. Chapter 5's revised media
example uses its existing Sonnet 4.6 row. See the
[supported-model correction](../../docs/edition-2-notes/supported-model-correction.md).

Original Chapter 5 was absorbed into Chapter 4. From that point, Edition 2
numbers are one lower through original Chapter 21; the remaining merges below
bring the final offset to three. The remaining standalone
`ch05/` is the actor chapter, not the deleted refactor chapter.

| Edition 2 chapter | Accepted solution | Original chapter | Original-numbered checkpoint |
| --- | --- | --- | --- |
| 1 | [ch01](ch01/) | 1 | `edition-2-ch01-r1` |
| 2 | [ch02](ch02/) | 2 | `edition-2-ch02-r3` |
| 3 | [ch03](ch03/) | 3 | `edition-2-ch03-r3` |
| 4: jobs and external consumer | [ch04](ch04/) | 5, including 4 | `edition-2-ch05-r2` |
| 5: actor and steering | [ch05](ch05/) | 6 | `edition-2-ch06-r3` |
| 6: streaming | [ch06](ch06/) | 7 | `edition-2-ch07-r4` |
| 7: GUI and replay | [ch07](ch07/) | 8 | `edition-2-ch08-r3` |
| 8: workbench and settings | [ch08](ch08/) | 9 | `edition-2-ch09-r2` |
| 9: Skills | [ch09](ch09/), [manifest](ch09-manifest.json) | 10 | `edition-2-ch10-r2` |
| 10: Persistence | [ch10](ch10/), [manifest](ch10-manifest.json) | 11 | `edition-2-ch11-r1` |
| 11: MCP | [ch11](ch11/), [manifest](ch11-manifest.json) | 12 | `edition-2-ch12-r1` |
| 12: The Agent Sees Itself | [ch12](ch12/), [manifest](ch12-manifest.json) | 13 | `edition-2-ch13-r1` |
| 13: Speech and blind persona | [ch13](ch13/), [manifest](ch13-manifest.json) | 14 | `edition-2-ch14-r1` |
| 14: Context management | [ch14](ch14/), [manifest](ch14-manifest.json) | 15 | `edition-2-ch15-r1` |
| 15: Memory bands | [ch15](ch15/), [manifest](ch15-manifest.json) | 16 | `edition-2-ch16-r1` |
| 16: Auto-Recall | [ch16](ch16/), [manifest](ch16-manifest.json) | 17 | `edition-2-ch17-r1` |
| 17: The Invisible Invoice | [ch17](ch17/), [manifest](ch17-manifest.json) | 18 | `edition-2-ch18-r1` |
| 18: Subscription authentication and reasoning | [ch18](ch18/), [manifest](ch18-manifest.json) | 19 | `edition-2-ch19-r1` |
| 19: The Crossover | [ch19](ch19/), [manifest](ch19-manifest.json) | 20 | `edition-2-ch20-r1` |
| 20: Web access, status and accounting | [ch20](ch20/), [manifest](ch20-manifest.json) | 21 + 22 | `edition-2-ch22-r1` |
| 21: The Sandbox | [ch21](ch21/), [manifest](ch21-manifest.json) | 23 | `edition-2-ch23-r1` |
| 22: The Ensemble, including reusable GUI | [ch22](ch22/), [manifest](ch22-manifest.json) | 24 + 25 | `edition-2-ch25-r1` |
| 23: Fleet management and GUI | [ch23](ch23/), [manifest](ch23-manifest.json) | 26 | `edition-2-ch26-r1` |
| 24: Optional Codex workaround | [ch24](ch24/), [manifest](ch24-manifest.json) | New in Edition 2 | `edition-2-optional-ch24-r1` |

The prose/grader enhancement pass is complete through Edition 2 Chapter 15.
Later accepted code follows original exercises and independently recorded grader
exceptions; its presence does not claim that later manuscript enhancement is done.
Original Chapter 20 / Edition 2 Chapter 19 is prose-only. Its bounded real coding
demonstration with GPT-6 Astra/high is accepted; no original20 grader is claimed.
Original21’s web-search exercise is independently accepted with two narrow
exceptions; its actual original grade remains 0/100 FAIL. Live search/fetch and
supplemental verification are recorded in the chapter review. The former original22/24 placement stops are lifted; follow the map below.

The revised grader covers enhanced Chapters 1–15, for example:

```sh
go run ./cmd/grade-edition2 -ch 4
go run ./cmd/grade-edition2 -ch 8
```

Later original graders retain original numbers and take the correctly mapped
submission directory. [Chapter notes](../../docs/edition-2-notes/) and
[`evidence/`](evidence/) also retain original numbering and actual results.
See the [enhancement record](../../docs/edition-2-enhancement.md) for revised
grader results. Do not treat accepted exceptions as passing original grades.

[`history/ch05/`](history/ch05/) preserves the retired standalone original
Chapter 5. History contains no other chapter directories. Earlier snapshot paths
and layouts remain available in their immutable Git tags; tags were not moved.
Old reports describe the original paths and numbering at the time of their runs.

Snapshot contents, nested example names and relative module replacements are
unchanged. This table supplies their current numbering. Keep accepted snapshots
immutable and extend `main/` for new work.

## Final chapter mapping

[The authoritative map](../../docs/edition-2-chapter-map.md) absorbs exactly
original Chapters 5, 22 and 24. Original 21 + 22 share `ch20`; original 23 uses
`ch21`; original 24 + 25 share `ch22`; original 26 uses `ch23`. Original26
was complete at original25 acceptance, and its implementation/review is complete.
The former original-21-only `ch20` remains
preserved in immutable `edition-2-ch21-r1`. The direct `ch20` now contains the
independently accepted combined source from original Chapters 21 and 22.
Future directories appear only when their combined exercises are accepted.
First-edition Chapters 1–26 therefore map to second-edition Chapters 1–23.
