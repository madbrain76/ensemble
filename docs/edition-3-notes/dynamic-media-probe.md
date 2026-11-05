# Ported feature: an unlisted model earns image capability from its endpoint

Branch `port/4-dynamic-media-probe`, based on `fix/linux-path-case`. Coordinator, 2026-11-03. Ported from `go-checkpoint` (`agent/internal/engine/media_probe.go`).

`Engine.LookupModel` is fail-closed: a model the curated table does not name gets `Media.None`, and `CheckMedia` refuses bytes for any model that is not `known`. That is correct as application knowledge and wrong as a claim about a live endpoint — Bill's local Qwen model receives images perfectly well and had them silently withheld. `ProbeImagesAsync` now closes that gap: a 1×1 solid-color PNG and a one-word question. A server that drops image parts has no way to name the color, so it fails instead of authorizing hallucinated answers.

The proof is two probes asking two different colors, issued in parallel so the first image send still waits for only one round trip. One probe is a one-in-three test of an endpoint that answers from a fixed word, and a positive is retained for the process, so the cheap certainty is worth the second small request. A reply must also not name the colors it was not asked about: a hedge over all three is not an observation. `Tests` pins both properties.

Three rules keep the probe from becoming a second authority:

- The table is never overridden in either direction. A listed row answers immediately and issues no request; only an unlisted route is asked. `Tests` pins both halves.
- Only the delivered answer counts. `TextPart` is the answer and `OpaquePart` is provider reasoning, so a color that appears only in `reasoning_content` refuses the probe. A thinking model can therefore not pass by guessing inside its chain of thought.
- The route key is vendor, surface, base URL and model. Proof over one dialect is
  not proof over another, because the image part and the response shape differ.
- Results are retained per route, failures included, so neither a text-only endpoint nor an unreachable server is charged again per attachment. A negative expires after ten minutes so one transient timeout never becomes a permanent refusal; a positive never needs refreshing. Teardown and the caller's own cancellation are never cached as a negative.

The probe runs only where waiting is already normal — the GUI image upload boundary and the async screenshot tool — and never inside `Render`, which stays pure. It reuses `SendCoreAsync` with `auxiliary: true` exactly as the recall judge does, so authentication, same-source plan retry, parsing and auxiliary accounting stay in one place; no parallel transport exists.

Defect found while landing this: the first draft fired a background probe on model selection and on peer construction. Both were wrong and were removed. They charged a request for capability nobody had asked for, and, more seriously, an in-flight probe was observable by the Agent's own invariants — `Chapter 18 actor-owned model selection`, `Chapter 21 ordinary peer chat across restart` and `memory pressure unknown` all failed while it held the Engine-send guard, and a probe must never be the reason a turn or a selection waits. On-demand probing from an async caller passes the full suite unchanged.

Verified in three layers. `agent/Tests` is 344 passed, 4 failed against an identical 341 passed, 4 failed baseline taken in a clean `fix/linux-path-case` worktree, so the only failures are the pre-existing macOS PTY ones and the three new cases pass; the image group was repeated five times because the color pair is random, and the fake provider decodes the request JSON rather than matching it as text, since `System.Text.Json` escapes `+` in base64 as `\u002B` and a raw substring assertion would have been a flake.

The changed GUI upload path was then driven in real headless Firefox against the actual `Examples/Gui` server with only the provider faked, attaching a real PNG through `#attach-images` and its file chooser — the composer ignores a file written straight to the hidden input, so the button is the only honest route. With a capable provider the browser showed two distinct color probes at the endpoint, "Images admitted", the attached image rendered in the transcript and the reply arrived. With a provider that answers "I cannot see an image" the same run showed two probes, "Image message refused; draft retained", no reply and no forwarded bytes. Three capable runs and one blind run, all stable.

Not verified: a real Qwen endpoint, which needs Bill's controller. The harness above is a scratch script rather than a landed `agent/GuiTests` case; it belongs with the Firefox tier on `port/21-firefox-harness` rather than in this branch.

Consequence worth Bill's attention: the *reasoning* half of the same curated table still refuses `--reasoning-effort` for unlisted models (`Tests/ResponsesTests.cs` pins that), and this port deliberately did not touch it. Capability-by-probe for reasoning would need a different proof than a color name.
