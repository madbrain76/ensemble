# Chapter 1 student evidence

Submitted source: `d2baff64055ac33702a91c6d6ee5a4246e84c0c9` in `/Users/bill/projects/ensemble-edition-3-student`.
Book/grader: `7f3d069420112c2136bc95d1d551e12fad6964c1`.
See `../ch01-student-review.md` for observations, failures, assistance and limitations. These artifacts are evidence, not acceptance.

## Final candidate

- `revisions.txt`: baseline, submitted revision, chapter hash, Go version, final executable hash and clean status.
- `tests-final.txt`, `tests-final.exit`: `go vet ./... && go test -race -count=1 -v ./...`, student cwd. Exit 0; both real 62-second waits passed.
- `grader.json`, `grader.stderr`, `grader.exit`: main-repository cwd, `go run ./cmd/grade-edition3 -ch 1 -json /tmp/e3-ch01.VimxYH/ch01-final`. Exit 0, `passed: true`, 100/100 including both required lifetime gates. Empty stderr file is expected; the JSON contains the captured submission diagnostics.
- `live-default.txt`: actual final executable's four-turn PTY chat; model `claude-sonnet-5`; EOF exit 0. It was run by `python3 /tmp/e3-ch01.VimxYH/live-final.py`, which opened a PTY, launched `/tmp/e3-ch01.VimxYH/ch01-final chat`, sent each recorded line only after a prompt, then sent terminal EOF. Key resolution stayed in the local process and only the intended child environment; no shell credential export or argv credential.
- `models.txt`: successful live `GET /v1/models`, selected model present.
- `hygiene.txt`: final credential-presence and Git cleanliness checks.

## Explicitly intermediate or negative evidence

- `grader-initial.stderr`, `grader-initial.exit`: the required directory-target invocation failed before checks (`permission denied`, go-run exit 1 / grader status 2). The parent supplied the executable-target clarification; no grader changes were made.
- `grader-60bf58a.json`, `grader-60bf58a.exit`: first passing executable-target run, before the final invalid-root guard/comments.
- `tests-race.txt`: first full race suite, before the final guard test.
- `sensitivity.txt`: nine mutations at `60bf58a`, exact replacements/commands/assertion failures/exits and confirmation that originals were restored. They exercise the same production paths retained in the final candidate; the final guard has its own public regression test.
- `live-initial.txt`: original supplied credential path was absent; no API call occurred.
- `live-default-60bf58a.txt`: first successful default live chat.
- `live-terse.txt`: isolated copy of `60bf58a`, changing only the system string to the chapter's terse/annoyed/Waywardgeek example; two live turns, EOF exit 0. The variant was never committed.

Temporary binaries, copied source and one-off drivers were kept outside the student repository and removed after evidence collection. No actual credential value appears in these files or tracked student files. The student repository is clean and has no binaries.
