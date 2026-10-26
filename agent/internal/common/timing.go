package common

import "time"

// TurnTiming is a completed turn's wall-clock breakdown, not a sum of job runtimes.
type TurnTiming struct {
	// Started records server admission before mailbox or pending-turn queueing.
	Started time.Time `json:"started"`
	// TotalMS runs from admission to logical completion, before final persistence.
	TotalMS int64 `json:"total_ms"`
	// ModelMS covers provider waits, including auxiliary judgment and retries.
	ModelMS int64 `json:"model_ms"`
	// ToolMS covers blocking tool waits only while no model request is active.
	ToolMS int64 `json:"tool_ms"`
	// OtherMS includes queueing, pauses and local work outside the measured waits.
	OtherMS int64 `json:"other_ms"`
	// Calls counts dispatched invocations, including those returning background handles.
	Calls int `json:"calls"`
	// Requests counts provider HTTP attempts, including retries and recall judgment.
	Requests int `json:"requests"`
	// Interrupted distinguishes a bounded cancellation from normal completion.
	Interrupted bool `json:"interrupted"`
}

// WaitTimingKey identifies the turn-scoped observer carried by derived contexts.
type WaitTimingKey struct{}

// WaitObserver reports synchronous wait boundaries without owning turn state.
type WaitObserver interface {
	// ModelStarted counts an HTTP attempt and returns its idempotent completion.
	ModelStarted() func()
	// ToolStarted counts a dispatch and returns the end of its blocking wait.
	ToolStarted() func()
}
