# Edition 2 Chapter 13 grader: The Channel Nobody Tested

Coder checkpoint; independent acceptance is pending. Original Chapter14 and both
frozen implementations are unchanged. Read the full
`.agents/skills/ensemble-coding/SKILL.md`, carryover and workflow before resuming.

The separately versioned driver preserves the actual original95 speech points
plus5 corrected Chapter12 parity. It executes the existing public tts-harness
contract; it does not detect implementation names or introduce a persona protocol.
Assertions now require recognizable complete utterances, a fence announcement,
actual successful correlated file output before silence credit, failure-specific
speech, and a successful harness exit even when a partial log exists. Provider
nonstream mode is distinct from final-only renderer delivery. Typing pause/release
logs do not prove whitespace preservation or both-cause ordering. Those stronger
claims, five-tool Agent isolation, native speech and pending-sleep EOF cleanup
retain separately accepted source/browser evidence.

Focused receipts live in `solutions/edition-2/evidence/ch14/enhancement13/`.
The first accepted run could not import Playwright and scored0/95; installing the
locked GUI dependencies changed no tracked source. The next75/95 exposed a fixture
error: the harness launches its GUI from a scratch cwd, and read_file interprets
relative paths there, while WORKSPACE only supplies configuration for other tools.
The planted file now uses the existing absolute-path input; the actual correlated
result contains ZEBRAFISH, independently of the spoken PELICAN answer. Root and
reviewer independently confirmed this diagnosis. No runtime patch was made.

With the absolute path, accepted and legacy initially scored95/95. Root identified
that one isolated bullet missed the documented adjacent-line legacy defect. The
corrected pair gives accepted95/95 and legacy75/95: actual legacy speech retains
`Bullet words - Following bullet.` in both delivery modes. Initial measurements
remain in their original receipts. These are focused results, not full grades;
TestChapter13Probe exits successfully to collect failed checks, so its process
exit is not a score or pass assertion.

The real-command race test passed before the adjacent-bullet fixture addition;
the final pair was then exercised on both implementations. Vet, command tests
and formatting passed. Final full parity, targeted sensitivity, exact published
samples and source bindings are still pending at this stage. Historical original
95/100 remains its earlier-source result; no new paid model run is required.
