package llm

import (
	"bytes"
	"encoding/json"
	"sync/atomic"
	"testing"

	"github.com/waywardgeek/ensemble/agent/internal/common"
	"github.com/waywardgeek/ensemble/agent/internal/jobs"
)

// recordingJobs is the real registry with a counter on Interrupt. Embedding
// the interface keeps the real behavior, so this records what happened rather
// than standing in for it.
type recordingJobs struct {
	common.JobManager
	interrupts atomic.Int64
}

func (r *recordingJobs) Interrupt() int {
	r.interrupts.Add(1)
	return r.JobManager.Interrupt()
}

// A message arriving while a tool is outstanding must end the wait, not just
// land in context. The hint alone leaves the agent parked until the tool
// returns, which is how a parent and the child it is waiting on end up
// waiting on each other.
func TestMessageDuringToolWaitInterruptsTheWait(t *testing.T) {
	e := lostCallEngine(t)
	rec := &recordingJobs{JobManager: jobs.NewJobs(e.parent)}
	e.Jobs = rec
	a := NewActor(e, e.parent)

	a.mb.Post(common.UserMessage{Text: "answer your child"})
	a.mb.Post(common.ToolCompleted{CallID: "current"})
	a.waitForTool("current")

	if got := rec.interrupts.Load(); got != 1 {
		t.Fatalf("Interrupt called %d times, want exactly 1 for one message", got)
	}
}

// hintReached reports whether the text made it into context. It checks the
// serialized form because a hint attached while calls are outstanding waits in
// Held before it lands in Dialogue, and the test cares that the text arrived,
// not which of the two it is resting in at the moment.
func hintReached(e *Engine, text string) bool {
	blob, err := json.Marshal(e.Ctx)
	if err != nil {
		return false
	}
	return bytes.Contains(blob, []byte(text))
}

// The hint path stays. The message has to reach context as well, or the agent
// wakes up with no idea what it was woken for.
func TestMessageDuringToolWaitStillReachesContext(t *testing.T) {
	e := lostCallEngine(t)
	e.Jobs = &recordingJobs{JobManager: jobs.NewJobs(e.parent)}
	a := NewActor(e, e.parent)

	a.mb.Post(common.UserMessage{Text: "steer the work"})
	a.mb.Post(common.ToolCompleted{CallID: "current"})
	a.waitForTool("current")

	if !hintReached(e, "steer the work") {
		t.Fatal("message did not reach context; interrupting must not replace the hint")
	}
}

// A tool completing on its own must not be reported as a message. Nothing
// arrived, and an agent told otherwise goes looking for a message that was
// never sent.
func TestToolWaitWithNoMessageDoesNotInterrupt(t *testing.T) {
	e := lostCallEngine(t)
	rec := &recordingJobs{JobManager: jobs.NewJobs(e.parent)}
	e.Jobs = rec
	a := NewActor(e, e.parent)

	a.mb.Post(common.ToolCompleted{CallID: "current"})
	a.waitForTool("current")

	if got := rec.interrupts.Load(); got != 0 {
		t.Fatalf("Interrupt called %d times with no message, want 0", got)
	}
}
