# Edition 3 C# application; Go is retained for grading and publishing tools.
.DEFAULT_GOAL := all
DOTNET ?= dotnet
DOTNET_PATH := $(shell command -v $(DOTNET) 2>/dev/null)
GO ?= go
NODE ?= node
CH ?= 1

.PHONY: all build ensemble cli clean test test-csharp test-go test-gui test-js test-legacy-go vet fmt lint-prose grade grade-dir fake fake-serve book ebook print common-diagram
all: build
build: ensemble

ensemble:
	$(DOTNET) build agent/Examples/Gui/Gui.csproj
	@printf '%s\n' '#!/bin/sh' 'set -eu' 'export DOTNET_HOST_PATH="$${DOTNET_HOST_PATH:-$(DOTNET_PATH)}"' 'here=$$(CDPATH= cd -- "$$(dirname -- "$$0")" && pwd)' 'exec "$$here/agent/start" "$$@"' > ensemble
	@chmod +x ensemble

cli:
	$(DOTNET) build agent/Cli/Ch04.csproj

clean:
	$(DOTNET) clean agent/Examples/Gui/Gui.csproj
	$(DOTNET) clean agent/Cli/Ch04.csproj
	rm -f ensemble

test: test-csharp test-go
test-csharp:
	$(DOTNET) run --project agent/Tests/Tests.csproj --no-launch-profile

# The GUI's own tiers, kept out of `test` on purpose: the browser tier needs an
# installed Firefox, and both are far slower than the unit suites. They were
# previously reachable only by remembering the invocation, which is how a stale
# assertion sat failing unnoticed inside the external consumer suite.
# `make test-js SUITES=image-picker` runs one browser suite.
# Twelve test classes were reachable only when named: the consumer program's own
# no-argument flow covers 11 classes and leaves these out, which is how a failing
# assertion sat inside AssetCacheTests with nothing to invoke it.
#
# The harness's naming convention is the registry rather than a list kept here:
# every test selector branch ends in `-tests`, and the branches that do not are
# long-running fixture servers for the browser tier, or the default flow itself.
# Enumerating them instead would have meant remembering to edit this file, which is
# the same failure that hid the twelve. A new class therefore needs no Makefile
# change, and if the convention ever breaks the target says so rather than passing.
# Every bare selector branch the harness offers. Two guards read this list: one
# insists some test selectors exist, the other insists none of them escaped the
# convention. A single future branch named `--foo-checks` would otherwise be
# skipped in silence, which is the residue the empty-list guard alone cannot catch.
GUI_BRANCHES := $(shell grep -oE 'args is \["--[a-z0-9-]+"\]' agent/GuiTests/Program.cs | sed -E 's/.*"(--[a-z0-9-]+)".*/\1/' | sort -u)
GUI_SELECTORS := $(filter %-tests,$(GUI_BRANCHES))
GUI_UNCLASSIFIED := $(filter-out %-tests %-host --chat-reasoning-projection,$(GUI_BRANCHES))
# Each selector is its own process, so no shared static state can order the results.
#
# The built dll is spawned rather than `dotnet run`: dotnet run pins the working
# directory to the project directory and this SDK offers no --no-cwd-change, so
# GoalControlsTests, which resolves agent/skills/goals/SKILL.md against the caller's
# directory, fails under dotnet run while passing wherever the dll is spawned - as
# the tier runner already spawns it.
GUI_TESTS_DLL := agent/GuiTests/bin/Debug/net10.0/GuiTests.dll

test-gui: build
	@test -z '$(GUI_UNCLASSIFIED)' || { echo 'test-gui: selector branches outside the -tests/-host convention: $(GUI_UNCLASSIFIED)' >&2; exit 1; }
	@test -n '$(GUI_SELECTORS)' || { echo 'test-gui: discovered no --*-tests selectors; the harness naming convention has changed' >&2; exit 1; }
	$(DOTNET) build agent/GuiTests/GuiTests.csproj -v q --nologo
	$(DOTNET) $(GUI_TESTS_DLL)
	@for selector in $(GUI_SELECTORS); do \
	  printf 'GuiTests %s\n' "$$selector"; \
	  $(DOTNET) $(GUI_TESTS_DLL) $$selector || exit 1; \
	done

test-js: build
	$(NODE) agent/GuiTests/run-browser-tests.cjs $(SUITES)

test-go: vet
	$(GO) test ./cmd/... ./internal/gradeedition3 ./internal/fakevendor ./internal/fakeanthropic

# Historical source-layout/pricing expectations are not V3 acceptance.
test-legacy-go:
	$(GO) test ./cmd/... ./internal/...
vet:
	$(GO) vet ./cmd/... ./internal/...
fmt:
	$(GO) fmt ./cmd/... ./internal/...
lint-prose:
	$(GO) run ./cmd/lintprose $(LINTFLAGS) book/preface.md book/chapter-[0-9][0-9].md

# Always supply an explicit C# executable/wrapper and chapter-specific flags.
grade grade-dir:
	@test -n "$(DIR)" || (echo 'Set DIR to the explicit submission or CLI wrapper; see the chapter grading scope' >&2; exit 2)
	$(GO) run ./cmd/grade-edition3 -ch $(CH) $(GRADEFLAGS) $(DIR)

VENDOR ?= anthropic
FAKE_CH ?= 3
fake:
	$(GO) run ./cmd/fakevendor -ch $(FAKE_CH) -vendor $(VENDOR) chat
fake-serve:
	$(GO) run ./cmd/fakevendor -vendor $(VENDOR)

book:
	scripts/assemble-book.sh > book/the-art-of-building-ai-coding-agents.md
ebook: book
	scripts/build-ebook.sh
print: book
	python3 book/build-print.py

# Source-only architecture diagrams; requires existing SDK, Python and Graphviz.
common-diagram:
	scripts/common-diagram/regenerate.sh "$(DOTNET)"
