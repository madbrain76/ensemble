# Chapter 1 — student correction round 1

**Revised candidate, not acceptance:** `1aad3975842d502b9144b7b9bddfd30ca7309db7` in `/Users/bill/projects/ensemble-edition-3-student`.

Previous candidate: `d2baff64055ac33702a91c6d6ee5a4246e84c0c9`. Book/unchanged grader revision remains `7f3d069420112c2136bc95d1d551e12fad6964c1`; chapter SHA256 remains `ab016baaff3f6a9d83ecffcbfc3772cd692c9d33c06d9fa1c52ad9a8708d25df`.

## R1 correction and scope

The parent supplied one narrow public-contract finding from independent review: my successful response fixture incorrectly described an invented `thinking` block with an extra `text` member as valid for this text-only chapter. I removed that block. The comment now says: “Two text blocks must be concatenated exactly, without added separators.”

The **only committed diff** is two lines replaced in `ensemble_test.go`: that comment and the fixture response. The remaining blocks are exactly `{"type":"text","text":"hello "}` and `{"type":"text","text":"世界"}`; usage stays input=17/output=9. All assertions are untouched, including exact `hello 世界` concatenation, history replay, usage totals, wire/configuration checks, request count, owner links and snapshots. No unexpected-input case was added. No thinking support, production filtering change, new acceptance requirement or broad hardening was introduced.

## New validation

All commands ran from the student repository using Go 1.25.4 on darwin/arm64:

- `go test -count=1 -v -run '^TestConversation$' .`: passed before mutation and after restoration.
- Focused mutation: temporarily inserted `break` immediately after production `text.WriteString(block.Text)` and ran that same test command. It exited **1**, specifically at `all text blocks: "hello " <nil>`; this was a behavioral assertion failure, not a compile failure. Restored the Engine file byte-for-byte in a `finally` block and reran the passing test before committing.
- `go vet ./...`: exit **0**.
- `go test -race -count=1 -v ./...`: exit **0**, package elapsed **125.327s**. Both actual 62-second header/body waits passed, as did public API and built-CLI tests.
- `go build -o /tmp/e3-ch01-correction1-build.ZtEXNc/ch01 ./cmd/ch01`: exit **0**. `go version -m` confirmed revision `1aad3975842d502b9144b7b9bddfd30ca7309db7`, `vcs.modified=false`. The temporary build directory was removed.

## Production and evidence provenance

`provenance.json` records baseline/corrected SHA256 values for every tracked file and checks working-tree bytes against the committed bytes. **All five production Go files and `go.mod` are byte-identical to the previous candidate.** `README.md` and `cli_test.go` are also identical; only `ensemble_test.go` differs. The Git path set is unchanged, and `correction.patch` shows the complete two-line replacement.

Accordingly, the original student grader result (`passed=true`, 100/100, both lifetime gates passed, exit 0) and final default live PTY result remain evidence for the identical production source at `d2baff64055ac33702a91c6d6ee5a4246e84c0c9`. They are **not relabeled as new runs at this correction commit**. I did not rerun the grader, make provider calls or access credentials this round. New binaries may embed the new VCS revision; the proof here is source-byte identity, not a claim that VCS-stamped binary bytes match.

The original student review, its evidence, and `ch01-candidate.json` were hashed before correction and verified unchanged afterward. No original artifact was edited. New evidence is exclusively under `docs/edition-3-notes/ch01-student-correction-1-evidence/`:

- `correction.patch`, `provenance.json`, `original-artifact-hashes.json`: exact change and preservation/identity proof.
- `focused-before-mutation.txt`, `first-block-mutation.txt`, `focused-after-restoration.txt`: focused sensitivity and recovery.
- `vet.txt`/`.exit`, `tests-final.txt`/`.exit`, `build.txt`/`.exit`: final command outputs and status.
- `final-hygiene.txt`: final clean-tree and original-artifact recheck.

## Assistance and remaining disposition

Assistance was limited to the supplied R1 finding and instructions. I did not read the full private review, historical/reference evidence, grader implementation or future chapters, and did not delegate. No production defect was reported by the supplied review finding; I make no independent acceptance claim from that statement. No new failure occurred except the intentional mutation failure.

The revised source/tests are committed, and the student repository is clean. No production changes, binaries, credentials, push, tags, author edits or Chapter 2 work. Independent re-review and owner acceptance remain pending.
