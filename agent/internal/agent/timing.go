package agent

import (
	"sync"
	"time"

	"ensemble/internal/common"
)

// The turn owns this clock. HTTP and report workers only bracket their waits;
// sealing makes late completions harmless even after a newer turn has started.
type turnClock struct {
	mu                  sync.Mutex
	started, last       time.Time
	modelTime, toolTime time.Duration
	models, tools       int
	requests, calls     int
	closed              bool
}

func (clock *turnClock) advance(at time.Time) {
	elapsed := at.Sub(clock.last)
	if elapsed < 0 {
		return
	}
	// Attribute each instant once. A provider wait takes priority over any
	// overlapping foreground tool wait; background job lifetimes never enter here.
	if clock.models > 0 {
		clock.modelTime += elapsed
	} else if clock.tools > 0 {
		clock.toolTime += elapsed
	}
	clock.last = at
}

func (clock *turnClock) start(at time.Time, model bool) {
	// Settle the old state before changing its counters. Reversing this order
	// would charge queueing or a pause to the wait that only just started.
	clock.advance(at)
	if model {
		clock.models++
		clock.requests++
	} else {
		clock.tools++
		clock.calls++
	}
}

func (clock *turnClock) end(at time.Time, model bool) {
	// The completed wait owns the interval up to this boundary, not after it.
	clock.advance(at)
	if model {
		clock.models--
	} else {
		clock.tools--
	}
}

func (clock *turnClock) wait(model bool) func() {
	clock.mu.Lock()
	defer clock.mu.Unlock()
	if clock.closed {
		return func() {}
	}
	clock.start(time.Now(), model)
	ended := false
	return func() {
		clock.mu.Lock()
		defer clock.mu.Unlock()
		// Both repeated completion and completion after interruption are harmless.
		// The captured flag shares the clock lock rather than needing another lock.
		if !ended && !clock.closed {
			clock.end(time.Now(), model)
		}
		ended = true
	}
}

// ModelStarted observes the HTTP exchange, not rendering or credential renewal.
func (clock *turnClock) ModelStarted() func() { return clock.wait(true) }

// ToolStarted observes the report wait, not the continuing job's execution.
func (clock *turnClock) ToolStarted() func() { return clock.wait(false) }

func (clock *turnClock) finish(interrupted bool) *common.TurnTiming {
	clock.mu.Lock()
	defer clock.mu.Unlock()
	return clock.snapshot(time.Now(), interrupted)
}

func (clock *turnClock) snapshot(at time.Time, interrupted bool) *common.TurnTiming {
	// Open waits stop at logical turn completion, even if HTTP cleanup or a
	// detached tool continues. The stored result is final, not a live stopwatch.
	clock.advance(at)
	clock.closed = true
	timing := &common.TurnTiming{
		Started: clock.started.UTC(), TotalMS: at.Sub(clock.started).Milliseconds(),
		ModelMS: clock.modelTime.Milliseconds(), ToolMS: clock.toolTime.Milliseconds(),
		Calls: clock.calls, Requests: clock.requests, Interrupted: interrupted,
	}
	// Round the residual last so displayed millisecond categories always reconcile.
	timing.OtherMS = timing.TotalMS - timing.ModelMS - timing.ToolMS
	return timing
}
