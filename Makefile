.PHONY: test vet fmt lint-prose grade grade-dir fake fake-serve book ebook

# The root module contains grading/tooling; runtime modules remain separate.
# Historical first-edition graders and mutation tools live under edition-1/.
test: vet
	go test ./cmd/... ./internal/...

vet:
	go vet ./cmd/... ./internal/...

fmt:
	gofmt -l -w cmd internal

lint-prose:
	go run ./cmd/lintprose $(LINTFLAGS) book/preface.md book/chapter-[0-9][0-9].md

# Edition 2 numbering. Chapters 19/24 are ungraded; 25 is a placeholder.
CH ?= 1
GRADE_CHAPTERS := 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 20 21 22 23 26
.PHONY: $(addprefix grade,$(GRADE_CHAPTERS))
grade:
	go run ./cmd/grade-edition2 -ch $(CH) $(GRADEFLAGS)

$(addprefix grade,$(GRADE_CHAPTERS)): grade%:
	go run ./cmd/grade-edition2 -ch $* $(GRADEFLAGS)

# Example: make grade-dir CH=26 DIR=./agent
grade-dir:
	@test -n "$(DIR)" || (echo 'Set DIR to the submission directory' >&2; exit 2)
	go run ./cmd/grade-edition2 -ch $(CH) $(GRADEFLAGS) $(DIR)

VENDOR ?= anthropic
FAKE_CH ?= 3
fake:
	go run ./cmd/fakevendor -ch $(FAKE_CH) -vendor $(VENDOR) chat

fake-serve:
	go run ./cmd/fakevendor -vendor $(VENDOR)

# The assembled manuscript is a publication deliverable, committed deliberately.
book:
	scripts/assemble-book.sh > book/the-art-of-building-ai-coding-agents.md
	@wc -w < book/the-art-of-building-ai-coding-agents.md | xargs -I{} echo "book/the-art-of-building-ai-coding-agents.md: {} words"

ebook: book
	scripts/build-ebook.sh
