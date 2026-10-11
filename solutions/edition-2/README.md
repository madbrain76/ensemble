# Edition 2 solutions

Accepted cumulative solutions live directly in `ch01/`, `ch02/`, and so on.
Chapters **1–18 are accepted**. [`main/`](main/) is the authoritative working
tree. Chapter 19 (original 20, The Crossover) awaits Bill's model choice. Auto-Recall follows
Bill’s October 11 quota/corpus decisions, recorded in the chapter notes.
Do not put new solutions in `history` or create symlinks.

Original Chapter 5 was absorbed into Chapter 4. From that point, Edition 2
numbers are one lower than the original book's numbers. The remaining standalone
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

The prose/grader enhancement pass is complete through Edition 2 Chapter 8.
Later accepted code follows original exercises and independently recorded grader
exceptions; its presence does not claim that later manuscript enhancement is done.
Original Chapter 20 / Edition 2 Chapter 19 is prose-only. Before its crossover
demonstration, Bill must identify the intended model: requested GPT-5 Astra is
absent from the observed account catalog, while GPT-6 Astra is available.

The revised grader covers enhanced Chapters 1–8, for example:

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
