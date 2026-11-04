# Fix: a rejected launch flag names itself

Branch `fix/launch-flag-error`, based on `fix/linux-path-case`. Coordinator, 2026-11-04.

`Examples/Gui` rejected an unsupported flag with `ArgumentException("Unknown or
incomplete GUI flag")`, and `ReportFailure` deliberately prints only the exception
type plus fixed advice for that class of failure. The result was a launch that
died with "Invalid launch configuration. Check supported flags and value
formats." and no indication of which of eight flags was wrong. Bill hit exactly
this running `run.sh`, whose spellings came from the Go launcher (`-model`,
`--api-key`, `--wol-mac`), and could not tell a typo from an unavailable feature.

`LaunchFlagException` now marks the two throws that reject the operator's own
command line, and `ReportFailure` prints that one type's message, which carries
the flag name and never the value that followed it — `--api-key some-secret`
reports `Unknown or incomplete GUI flag: --api-key`. Everything else keeps the
existing privacy-preserving behaviour, because arbitrary `ArgumentException` text
from the library can still echo paths and provider detail.

The CLI parses its own flags and already prints `ArgumentException.Message`
verbatim, so the three vague "Missing ... option value" strings there now name
the flag directly and needed no new type.

Verified against the built binary: `--wol-mac`, `-model`, `--api-key` and a
value-less `--wake-on-lan` each report their own name. No test or document
depended on the old exact strings.

## Follow-up: an aggregate identifies nothing

The same handler then hid a second class of failure. Composition joins the serving
and fleet-restore tasks, so a launch that died from unreadable persisted state
reported only `GUI failed: AggregateException. Check private diagnostics...` —
correctly private, but useless. `ReportFailure` now walks the whole cause chain
and prints the type names it finds (`AggregateException -> InvalidDataException`),
which carries no private text, and takes a rejected flag's
message from any depth rather than only the outermost exception. Diagnosing that particular
launch had been possible only by temporarily printing `error.ToString()`.
