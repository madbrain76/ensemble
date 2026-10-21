package ensemble

import "ensemble/internal/diagnostics"

// DiagnosticLog retains complete API or GUI trace records independently of
// recovery history. Positive limits keep newest records; zero is unlimited.
type DiagnosticLog = diagnostics.Log

// OpenDiagnosticLog wires an existing trace destination to its owner's writer.
// The callback is sampled by Retain at a safe startup or shutdown boundary.
func OpenDiagnosticLog(path string, limit func() int) (*DiagnosticLog, error) {
	return diagnostics.Open(path, limit)
}
