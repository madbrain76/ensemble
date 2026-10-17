# Chapter 5 entry audit evidence

See ../ch05-entry-audit.md. This is an entry audit at the accepted Chapter 4 predecessor, not Chapter 5 acceptance.

Every named control has a `.meta.json` (exact argv/cwd, measured seconds, exit code, outer watchdog status), `.stdout`, and `.stderr`. Empty streams are intentionally retained. A timeout/callback is not interpreted as a behavioral result. The first broad student suite timed out; its longer retry passed. Setup failures are retained separately from actual negative controls.

Paths use <AUDIT_SCRATCH>, <MAIN_REPO>, <STUDENT_REPO>, <USER_HOME>, and <SYSTEM_TMP>. Substitute a new scratch root when interpreting commands. Executables and private scratch drivers/probes are not distributed. The API adapter interface and semantics are documented in the report; it is not an installed repository command. Historical controls are private audit evidence, not student examples or repair fixtures.

provenance.json: exact original revisions, tag object and peeled commit, host/toolchain, required full reads and sanitized environment.
source-hashes.json: audited current documents, grader source and accepted student source/config bytes.
artifact-hashes.json: built executable/private-driver hashes and Go build information.
SHA256SUMS: SHA-256 of every exported evidence file except this manifest itself.

All provider traffic was to loopback fake endpoints. No inherited provider credentials were available. Local fixture credentials are inert test constants. No binaries, source fixtures, historical implementation, network payload fixture files, or paid/live-GUI evidence are included here.
