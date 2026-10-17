# Chapter 6 artifact portability

The committed report and safe text receipts preserve commands, source/artifact identities, failures and qualified results. Pinned executables, source archives, private runtime material and Python bytecode remain local-only; a fresh clone does not acquire those artifacts or reproduce their executions merely by reading their hashes. Materialize matching authorized artifacts before replaying local-only checks. Never copy private provider traces or credentials into Git.

The student inventories explicitly classify portable and local-only entries. Historical reviewer checksum files are preserved byte-for-byte, even where they also pin a local-only file. Thus a full original on-disk checksum check is distinct from verification of the portable clone subset. In particular, the original code-review checksum index pins the bytecode file listed below; it is intentionally excluded from the commit, not silently removed from the historical index.

## Local-only files physically under the main evidence directories

These files are preserved locally and explicitly excluded from staging:

| Path | Bytes | SHA256 |
|---|---:|---|
| `docs/edition-3-notes/ch06-student-evidence/candidate-source.tar` | 665600 | `6e678c85217ffa88cf32133971057606c60209c25da6b9b96c89b7e012cb9203` |
| `docs/edition-3-notes/ch06-student-evidence/candidate1/source.tar` | 665600 | `518c09084aaba668cf7de3d29dd55b1fd976f6b1060e67a47d5498880dd0da7a` |
| `docs/edition-3-notes/ch06-student-evidence/__pycache__/live_harness.cpython-313.pyc` | 17476 | `e0aa664ea1a1ed9e9da55afaf7cb76a2f195813c18c5293a0905198a8ab7d0c4` |
| `docs/edition-3-notes/ch06-student-evidence/__pycache__/shape_fixture.cpython-313.pyc` | 6404 | `cebe4309bcfc44aec73baf01609e032ff4619642e1d03ef8462f5fcfb565247d` |
| `docs/edition-3-notes/ch06-code-review-evidence/__pycache__/cli-gated-control.cpython-313.pyc` | 12318 | `100debc04d52b6cff774eccecce5b2ab2eec66af383a83155998ab6d2e5fec83` |

All other local-only entries named in the student inventories retain their recorded external paths. Private recovery copies are not repository evidence or imported computation.
