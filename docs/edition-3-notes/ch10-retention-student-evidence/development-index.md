# Development chronology and limitations of receipts

All raw development outputs are retained without replacing failed attempts.
These statuses are transcribed from supervised tool completion plus the raw outputs;
older wrapper outputs do not themselves contain explicit exit footers.

| Receipt | Actual result |
|---|---|
| development-core | exit 1: undefined `s` in loader; fixed to `save`; ownership location was initially misdescribed in chat, then corrected |
| development-core-repair | exit 0: existing Chapter 10 focused tests |
| development-delta | exit 1: new test helper returned value instead of pointer |
| development-delta-repair | exit 1: missing fixture provenance Surface |
| development-delta-fixtures | exit 1: tool-call fixture provenance absent and wrong CLI save path |
| development-delta-core | exit 0: repaired focused core/CLI tests |
| development-delta-gui | exit 0 BUT INVALID GUI CREDIT: wrapper silently cd'd into core; only core ran |
| development-delta-gui-correct-module | exit 1: new test used nonexistent GUI fixture helper names |
| development-delta-gui-fixtures | exit 0: repaired focused GUI/settings tests in correct module |
| development-core-complete | exit 1: new proof test used nonexistent MessageEvent.Text and Engine.Project API |
| development-core-complete-repair | exit 0: repaired full core run, inadvertently INCLUDING TestRequestLifetime (two 62-second waits), 132.717s core package |
| development-gui-complete | exit 0: full GUI run |
| negative-retention | exit 1: compile-valid owned-copy mutation; intended count assertion failed, 251 instead of 200 |
| restored-retention | exit 0: exact restored copy passes the same focused test |

The accidental lifetime selection is a finite-scope deviation, not a fresh live-provider,
consumer-binary or same-binary acceptance claim. It is explicitly excluded from final
root race validation. No paid/browser/legacy grader campaigns were launched.

A new CLI test initially wrote its synthetic settings into cmd/ch01/cr; it was moved
into the local evidence area and the test changed to use t.Chdir(t.TempDir()). The
preserved harmless fixture is published as development-misplaced-cli-settings.txt,
not as settings.json. The final tracked source is clean.

Non-test failed probes: rg is unavailable (subsequent searches used grep); guessed
paths internal/history.go and gui/http.go did not exist. No sibling was inspected.
Several large combined prerequisite/source reads were clipped and recovered with
bounded sections/separate reads. Initial source replacements missed exact anchors;
source inspection caught the missing encoder/startup changes before the repaired run.

Final receipts preserve start UTC, exact argv text, selected module PWD, isolated env,
Go version, free bytes, and outputs. The intended footer insertion did not match the
wrapper's printf timestamp line; no EXIT trap was present. Parent-observed completion
and supervised job exit statuses were appended with an explicit correction, not
invented exact end times. Root race then vet job 154 and GUI race then vet job 162
both exited 0. No test was rerun solely to improve receipt cosmetics.
