# Chapter 9 evidence

The student review is `docs/edition-2-notes/ch09.md` in the outer repository.
`original-grade.txt` is the sole running original grader's actual output;
`grader-source.json` records its input. No final score is claimed until it ends.
`review-source.json` describes the later review source (CSS correction, comments,
README and one additional public boundary test).

`browser-check.mjs` drives installed headed Chrome with real native speech and
WebSocket APIs, the built GUI application and a faithful external HTTP peer.
`live-check.mjs` discovers an available OpenAI model and drives one bounded live
turn, terminal hint and speaker. It reads only authorized credentials in memory;
do not print the private settings file. JSON files retain selected sanitized
observations, not raw private configuration or auth headers. Native events are
not a subjective listening result. Initial browser failures are retained with
explanations in the student review.

To reproduce, install `playwright@1.64.0` under `/tmp/ensemble-ch09-browser`, build
`solutions/edition-2/main/gui/cmd` into that directory as `ensemble-gui`, copy a
script there, and run it with Node. The scripts use installed Google Chrome,
ports 19391/19392 and disposable workspaces. `CH09_WEB` and `CH09_EVIDENCE`
redirect the deterministic browser check for an isolated static-source mutation.
The wrappers record speech API calls/events and invoke the actual native API.
They do not replace speech. No Node dependency belongs to the shipped GUI.

Mutation logs preserve each defect and its intended behavioral failure. Scratch
source/static copies are removed after testing. Unchanged Chapter 8 evidence
continues to cover earlier behavior; no repeated paid provider matrix was needed.
