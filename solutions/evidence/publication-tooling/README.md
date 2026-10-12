# Layout migration: active tooling validation

These are new relocation checks, not rewritten historical chapter results.
The migration base is `564e316f3f8aeac2ec26fe8edfbd8cdf42aa350b`.
The mandatory `.agents/skills/ensemble-coding/SKILL.md` was read in full;
workflow/carryover remain active from the complete reload.

Current default Chapter 1 and Chapter 26 grades are 100/100; archived original
Chapter 1 also scores 100/100. Commands, exits and whole-command times are in
`commands.json`. The maintained tooling packages build and vet. Their test
packages compile, with the actual Chapter 1 header regression selected; this is
not a claim that every chapter suite ran again. All active literal snapshot paths
exist after relocation. No archived or frozen source was edited.

The fakevendor launcher now sets the child cwd to its selected Go module and
runs `go run .`; the actual Chapter 3 loop completes three local HTTP requests,
two tool results and normal EOF. No paid provider or personal credential is used.
The current exporter defaults to committed `agent/`, with explicit historical
`--source-prefix solutions/edition-2/main` support. Four tests exercise real Git
repositories, preserving committed bytes/modes/ignored receipts, refusing existing
outputs and checking historical-prefix provenance.

The root Makefile addresses maintained `cmd`/`internal` packages, avoiding the
mixed-package historical fixtures in `scripts/edition2` and the separate runtime
modules. Current grade targets use Edition 2 numbers. Historical mutation scripts
are no longer active make targets and remain preserved with their archived edition.
The scripts README labels retired acceptance tools honestly instead of making
their obsolete commands look current.

A tiny real publication built from an unrelated temporary directory exercises
Pandoc and Chrome: zero figure references plus the repository cover, both populated
EPUB navigation documents, chapter-level navigation, no nav in the spine and no
printed Contents. The link-only Lua filter sends local repository links to public
GitHub URLs, preserving external URLs and internal anchors; images remain bundled.
Print resolves inputs from its own script directory and maps its heading to a
stable page number. Tool errors and continued pagination changes now fail rather
than claiming a successful print build. The failed-convergence branch was inspected,
not simulated with mocks. This fixture is not the full publication build.

`source.json` binds this coherent tooling stage. Final full book artifact checks
wait for the coordinator's epilogue/manuscript freeze. No current application
runtime, user process, original result, or frozen snapshot is modified by this work.
