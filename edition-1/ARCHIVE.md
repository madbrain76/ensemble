# First edition archive

This directory preserves the first edition of *The Art of Building AI Coding
Agents*, frozen at Chapter 27, Gateway. The current second edition lives at the
repository's top level.

- [Read the first edition](book/the-art-of-building-ai-coding-agents.md).
- [Original repository README](README.md).
- [Original application](agent/) and [chapter solutions](solutions/).
- [Original graders](internal/grade/) and [grader command](cmd/grade/).

The original manuscript, application, solutions, grading code and publication
files were moved here without changing their tracked contents. Original build
and maintenance scripts are preserved as historical artifacts; some record the
former top-level layout or Bill's machine-specific paths. The current book's
build instructions are in [book/README.md](../book/README.md).

The archive has its original Go module and support packages. Run original
chapter grading from this directory so its relative solution defaults agree:

```sh
cd edition-1
go run ./cmd/grade -ch 1
```

For a particular source tree, supply the directory explicitly. Grader numbers
here are first-edition numbers. Historical failures remain recorded failures;
the second edition's revised graders and results are separate.

The parent repository retains shared grading helpers and fake providers because
its revised graders use them. That retained support does not change the original
copy here. Existing commits and tags were not rewritten or moved. See the
[relocation record](../docs/edition-2-publication.md) for source identity and checks.
