# Ported feature: offer workspace work products for download

Branch `port/11-offer-download`, based on `fix/linux-path-case`. Coordinator, 2026-11-03. Ported from `go-checkpoint` (`71b9f7b4`).

`offer_download(path)` publishes a descriptor for one finished file — name, path, mime, size — and the GUI renders a single labelled link with the size. The browser never composes a path: the route re-opens the given path through the owning Agent sandbox, so workspace containment is the sandbox decision rather than a string check, and directories, FIFOs, empty files and unreadable names all fail the same way.

The extension whitelist is one function, `Tool.DownloadMime`, shared by the tool and the route so a type can never be offered by one and refused by the other. Executables and HTML are excluded on purpose: ensemble should not become a distribution point, and the browser is not made to sniff a file type. Origin, same-site and 64 MiB guards match the neighbouring job-output route.

Adding a built-in changed the pinned tool registry assertion in `agent/Tests/ToolTests.cs`; that list is the contract, so it was updated rather than the tool suppressed.
