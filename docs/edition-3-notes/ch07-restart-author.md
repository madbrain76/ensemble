# Chapter 7 restart — bounded author pass

## Changes and reasons

Kept the opening, voice, historical account, useful rationale and qualifications.
Six small insertions clarify demonstrated boundaries already within this chapter:

- Empty `part_delta` retains `"chunk":""`, including the key.
- Saving/publishing dispatch and releasing its handler share admission; cancellation
  cannot leave a published dispatch describing an unreleased reservation.
- Ping replies are payload-echoing WebSocket Pong control frames through the writer,
  rather than JSON messages.
- Same-card dispatch/job metadata updates are valid; identity reuse for another
  logical part is the collision.
- Count replay contents by renderable history entries, not envelopes or unique cards.
- Install real-speech listeners before playback, so a fast start is observable.

The replay window remains last 100 renderable events OR last hour, whichever reaches
farther back. Parts and live-owner contracts are unchanged. Existing distinctions
between callback evidence and human hearing, historical provider results and current
availability remain intact. The listener advice describes observation timing, not a
product failure. No first-person provider/hearing experience, independent visual
inspection, complete coverage or new capability is claimed. Bill's preface and all
other chapters, source, graders and guidance are untouched.

Read the requested current guidance, full chapter/voice/procedure and restart notes;
used their FINAL outcomes rather than superseded findings. Inspected only relevant
accepted admission, projection, replay, socket and speech source. This was not another
code review or application test campaign.

## Prose lint

Sequential baseline/final same-basename temporary copies; both commands exited 0:

```sh
env -i HOME="$HOME" PATH="$PATH" GOCACHE=/tmp/e3-ch05-build-cache \
  go run ./cmd/lintprose /tmp/ch07-restart-author/baseline/chapter-07.md
env -i HOME="$HOME" PATH="$PATH" GOCACHE=/tmp/e3-ch05-build-cache \
  go run ./cmd/lintprose /tmp/ch07-restart-author/final/chapter-07.md
```

Baseline output:

```text
/tmp/ch07-restart-author/baseline/chapter-07.md
  ok   confessions                       0   <=3
  ok   self-defense                      0   <=1
  WARN negation forms                   20   <=9 (1/500w)
  ok   inventory openers                 1   <=3
  ok   short closers                    35   <=35 of 105 paras
  ok   pointer closers                   1   <=2
  ok   scope superlatives                0   <=2
  ok   won't pretend                     0   <=1
  ok   em-dashes                         0   0
  ok   LLM crutches                      0   0
  WARN person gap (after L414)        2334   <=1200
  ok   heading shape run                 0   <3
```

Final output:

```text
/tmp/ch07-restart-author/final/chapter-07.md
  ok   confessions                       0   <=3
  ok   self-defense                      0   <=1
  WARN negation forms                   20   <=9 (1/500w)
  ok   inventory openers                 1   <=3
  ok   short closers                    35   <=35 of 107 paras
  ok   pointer closers                   1   <=2
  ok   scope superlatives                0   <=2
  ok   won't pretend                     0   <=1
  ok   em-dashes                         0   0
  ok   LLM crutches                      0   0
  WARN person gap (after L420)        2428   <=1200
  ok   heading shape run                 0   <3
```

Both retain soft warnings for 20 negation forms and the long technical person gap
(2,334 → 2,428 words). These do not justify padding the passage with invented personal
experience or rewriting its rationale. An initial exact-anchor edit refused before
writing; its following lint merely repeated the unchanged baseline. The final output
above is from the successful edited draft. No application builds/tests, browser or
provider runs, keys, children, commits or tags were used. Disk free space stayed above
8 GiB, safely above the 1 GiB floor; no unrelated cleanup occurred.

Read-only first draft for the single formal text review and possible one response:
`/tmp/ch07-restart-author/first-draft-chapter-07.md` (mode 0444). The baseline and lint
outputs remain in that same owned temporary directory. No response revision yet.
