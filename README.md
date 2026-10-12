# The Art of Building AI Coding Agents

*Second edition. An agentic codebook: build the agent, test the teaching, improve the book.*

**[Read the book](book/the-art-of-building-ai-coding-agents.md)** ·
[Chapter guide](book/chapters.md) · [Solutions](solutions/README.md) ·
[First edition archive](edition-1/ARCHIVE.md)

Ensemble is the coding agent built by this book's exercises. It has tools and
supervised jobs, streaming CLI and browser interfaces, persistent conversations,
skills, memory, a fleet of agents and an authenticated gateway. The code is here,
along with the instructions, graders and reviews that shaped it.

The second edition was built with Codex under Bill Cox's direction, using a
student implementation followed by comparative review and revisions to both
teaching and graders. It catches up through the first edition's final exercise,
Gateway, now Chapter 26. Original Chapters 5, 22 and 24 were absorbed into related
chapters while preserving their capabilities.

**Crossover has begun.** Ensemble has performed real coding tasks with GPT-6
Astra at high reasoning. Making it the everyday tool for the next edition is
ongoing work. Chapter 25 is deliberately a placeholder for that experience;
this publication does not claim Ensemble has already produced Edition 3.
Chapter 24 is an optional, ungraded Codex-endpoint workaround that readers may
skip when their public OAuth route caches satisfactorily.

## Repository layout

| Path | Purpose |
| --- | --- |
| [`book/`](book/) | Current second-edition manuscript and publication files |
| [`agent/`](agent/) | Current application development tree, with its own Go module |
| [`agent/gui/`](agent/gui/) | Optional browser application, in a separate Go module |
| [`solutions/chNN/`](solutions/README.md) | Accepted cumulative chapter snapshots; use matching book numbers |
| [`solutions/evidence/`](solutions/evidence/) | Retained runs and checks, using original exercise numbers |
| [`cmd/grade-edition2/`](cmd/grade-edition2/) | Current chapter grader command |
| [`docs/edition-2-notes/`](docs/edition-2-notes/) | Student, author and independent review records |
| [`edition-1/`](edition-1/) | Historical first-edition book, application, solutions and original graders |

Develop in `agent/`; keep accepted snapshots intact. Chapter 25 has no accepted
solution directory yet. The retired standalone refactor snapshot is the only
chapter under `solutions/history/`. Historical manifests and receipts retain
the paths and source revisions at which they were measured; the
[relocation record](docs/edition-2-publication.md) explains the move.

## Build and run

Use the Go version declared in each module. Build from the module directory:

```sh
mkdir -p bin
(cd agent && go build -o ../bin/ensemble .)
(cd agent/gui && go build -o ../../bin/ensemble-gui ./cmd)
```

For API-key use, select `LLM_VENDOR`, `LLM_MODEL` and `LLM_API_KEY`, then run the
CLI from the workspace you want it to use:

```sh
/path/to/ensemble/bin/ensemble chat --skills-dir /path/to/ensemble/agent/skills
```

Choose a model supported by the implementation and available to your account.
The application [README](agent/README.md) documents configuration, authentication,
streaming and persistence; the [GUI README](agent/gui/README.md) covers the browser.

For the optional Chapter 24 route, an existing Codex sign-in supplies the
read-only credential file. This example keeps conversation state outside the
checkout while tools work in the repository:

```sh
ensemble_repo="$PWD"
mkdir -p "$HOME/ensemble-session"
(
  cd "$HOME/ensemble-session"
  unset LLM_BASE_URL OPENAI_BASE_URL LLM_RESOLVED_MODEL
  LLM_VENDOR=openai LLM_MODEL=gpt-6-astra \
  OPENAI_REASONING_EFFORT=high OPENAI_BILLING_MODE=codex \
    "$ensemble_repo/bin/ensemble-gui" \
    --workspace "$ensemble_repo" \
    --skills-dir "$ensemble_repo/agent/skills" --save save.json
)
```

Use the browser address printed by the application. The route reads Codex's
existing credentials without copying or renewing them. See
[Chapter 24](book/chapter-24.md) for the workaround's evidence and limitations.
A successful answer alone does not prove prompt-cache reuse.

## Grade an exercise

From the repository root:

```sh
make grade                              # current Chapter 1 snapshot
make grade-dir CH=26 DIR=./solutions/ch26
go run ./cmd/grade-edition2 -ch 8         # matching snapshot by default
go run ./cmd/grade-edition2 -ch 3 /path/to/your/submission
```

Revised graders cover Chapters 1–18, 20–23 and 26. Chapter 19 is explanatory,
Chapter 24 is optional and ungraded, and Chapter 25 is reserved. Checks run the
actual application against local fake services without paid inference. The
chapter states the public contract and any separate review or live-use work.

Original grader failures remain failures in the historical record, including
independently accepted exceptions. Revised scores are separate results. Read
[the workflow](docs/agentic-codebook-workflow.md) and
[enhancement records](docs/edition-2-enhancement.md) for the distinctions.
To run an original grader, enter the archive first:

```sh
(cd edition-1 && go run ./cmd/grade -ch 1)
```

## Build the book

```sh
make book
make ebook
python3 book/build-print.py
```

The assembled Markdown is committed so it can be read directly in the repository.
The [publication instructions](book/README.md) describe EPUB navigation, print
page verification and the cover-prefixed PDF. The first-edition publication
files remain in `edition-1/book/`.

The book and code are [Apache-2.0 licensed](LICENSE). Fork them, change what your
agent needs, and preserve the reasons for the changes. A fresh implementation
is a useful test of how much the instructions actually teach.
