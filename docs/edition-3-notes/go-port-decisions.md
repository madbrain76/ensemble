# Go port items deliberately not landed

Bill's list of nineteen features to port from `go-checkpoint` to the Edition 3
C# tree. Seventeen are on their own `port/NN-*` branches. These two are not, and
this records why, so the decision is reviewable rather than silent.

## 6 — pending unload retains capability until compaction: not ported

`go-checkpoint` `babd50d2` ("agent: keep a pending unload's capability until
compaction") split unload into `MarkUnload`, which stops a skill's MCP servers
and defers compiled-tool removal, and `PurgePending`, called only when compaction
commits. In Go a pending skill therefore still dispatches `write_file` until
compaction.

Edition 3 deliberately made the opposite choice, and its book is explicit. The
current contract is immediate revocation with retained instructions:

- `book/chapter-14.md:19` — "Immediate revocation, retained instructions.
  Pending skill unload removes authority, including deferred references and
  automatic snapshots. Its manual remains in SYSTEM. Retiring that manual is not
  this chapter's job."
- `book/chapter-11.md:139` — "retaining a manual does not retain executable
  authority"; clients, processes and remote Tools close immediately.
- `book/chapter-12.md:157` and `:441` — unloading closes the client and removes
  executable browser tools immediately; a model that claimed MCP tools "remained
  available until compaction" was treated as a defect, and the unload wording
  was clarified to separate immediate server removal from pending instruction
  removal.
- `book/chapter-25.md:123` — "its authority is revoked immediately, while its
  old manual can remain in canonical History pending later retirement."

`agent/README.md:709` states the same rule, and four sites pin it as a failing
test if capability survives: `agent/Tests/SkillTests.cs:80` ("Pending unload
retained capability"), `agent/Tests/SkillAuthorityTests.cs:6`, and the public
consumer check `agent/Examples/Skills/Program.cs:74` ("Pending unload did not
revoke the real writer"), which asserts the file is unchanged after unload.
Chapter 9's "tools remain active" wording is the earlier rule that Chapter 14
supersedes, so porting Go's behavior would regress a security property and turn
four pinned checks red. It stays in Go, where Bill owns that call, unless he
overrules the current book.

Separately observed, not changed: nothing in C# retires a pending manual during
a session. `agent/Agent/Compaction.cs` handles only the context byte budget, and
skill entries are replaced at `RuntimeConfigured` (restart) in
`agent/History/History.cs:281`. That matches the book deferring retirement, so
it is a future-chapter item rather than a porting gap.

## 7 — MCP incoming service role: already in upstream C#

`go-checkpoint` added `ServeIncoming` so a peer can drive an Agent as an MCP
server without both sides sending `initialize`. Edition 3 already has this
shape, so no branch exists.

- `agent/Mcp/Service.cs:7` `ServeAsync` adopts the receiving role and refuses to
  convert an established outgoing connection into a service ("An outgoing MCP
  connection cannot become a service"), which is the race the Go change fixed.
  The comment records the invariant: there is exactly one reader per mode and
  the service issues no outgoing RPC.
- `ServeCoreAsync` answers `initialize` for protocol `2024-11-05`, refuses a
  repeated or unsupported one, requires `notifications/initialized` before any
  request is served, and routes `tools/list` and `tools/call` through the same
  Agent admission, skill grants and reverse Job lane as ordinary dispatch.
- `agent/Cli/Program.cs:447` already exposes it as `--mcp-pipe` over
  `McpStreamTransport` on stdin/stdout, with no conversation or HTTP started.
- `agent/Tests/McpServiceTests.cs:7` drives a real socket pair and asserts
  pre-initialization refusal, ID zero correlation, unsupported protocol refusal,
  grant-filtered discovery, and that a denied tool has no effects. Both cases
  pass: "Chapter 11 incoming MCP service and real file dispatch" and "Chapter 11
  service malformed input and joined termination".
