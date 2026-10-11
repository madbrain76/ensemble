package jobs

import (
	"testing"
	"time"

	"github.com/waywardgeek/ensemble/agent/internal/common"
)

// longWait is far longer than any of these tests take. Using it means a test
// that sees WokeMessage cannot be seeing an expired timer, so the interrupt is
// the only thing that could have ended the wait.
const longWait = 30 * time.Second

func runningJob() *Job {
	return &Job{status: common.StatusRunning, changed: make(chan struct{})}
}

// interruptWhenWaiting retries until a waiter is actually parked. Interrupt
// reports false until then, so this also pins that a job with no waiter
// refuses to be interrupted.
func interruptWhenWaiting(t *testing.T, j *Job) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if j.Interrupt() {
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatal("no waiter ever registered, so the wait could not be interrupted")
}

func TestInterruptEndsTheWaitAndLeavesTheJobRunning(t *testing.T) {
	j := runningJob()
	got := make(chan common.WakeReason, 1)
	go func() { got <- j.Wait(common.Limits{Delay: longWait}) }()

	interruptWhenWaiting(t, j)

	select {
	case r := <-got:
		if r != common.WokeMessage {
			t.Fatalf("wake reason = %v, want WokeMessage", r)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("interrupt did not end the wait")
	}
	// The point of interrupting rather than killing: the work continues.
	if j.Status() != common.StatusRunning {
		t.Fatalf("status = %v, want still running after an interrupt", j.Status())
	}
}

// The backgrounded tool must still deliver. An agent that got control back
// early and then never received the result would have to run the work again.
func TestInterruptedJobStillDeliversItsResult(t *testing.T) {
	j := runningJob()
	go func() { j.Wait(common.Limits{Delay: longWait}) }()
	interruptWhenWaiting(t, j)

	j.Finish("the result", nil)

	if j.Status() != common.StatusDone {
		t.Fatalf("status = %v, want done", j.Status())
	}
	if out := j.Report(common.WokeDone, common.DefaultLimits()); out != "the result" {
		t.Fatalf("report = %q, want the result to survive the interrupt", out)
	}
}

// A job nobody is parked on must refuse, or the flag would sit on a
// background job until its next waiter consumed it as a message that never
// arrived.
func TestInterruptRefusesWhenNobodyIsWaiting(t *testing.T) {
	j := runningJob()
	if j.Interrupt() {
		t.Fatal("interrupted a job with no waiter")
	}

	// A later wait must not find a leftover flag. The job is finished first
	// so this cannot hang if the behavior regresses.
	j.Finish("done", nil)
	if r := j.Wait(common.Limits{Delay: longWait}); r != common.WokeDone {
		t.Fatalf("wake reason = %v, want WokeDone; a stale interrupt was consumed", r)
	}
}

// A tool that has already produced its result reports that result. Reporting
// an interrupt instead would throw away work that was done.
func TestFinishedJobReportsDoneEvenWhenInterrupted(t *testing.T) {
	j := runningJob()
	j.Finish("output", nil)
	j.interrupted = true

	if r := j.Wait(common.Limits{Delay: longWait}); r != common.WokeDone {
		t.Fatalf("wake reason = %v, want WokeDone to win over a pending interrupt", r)
	}
}

// One interrupt ends one wait. If the flag survived, the next wait would
// return at once on a message that was already handled.
func TestInterruptIsConsumedByASingleWait(t *testing.T) {
	j := runningJob()
	go func() { j.Wait(common.Limits{Delay: longWait}) }()
	interruptWhenWaiting(t, j)

	// Give the first waiter time to consume the flag and leave.
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		j.mu.Lock()
		gone := !j.interrupted && j.waiters == 0
		j.mu.Unlock()
		if gone {
			break
		}
		time.Sleep(time.Millisecond)
	}

	second := make(chan common.WakeReason, 1)
	go func() { second <- j.Wait(common.Limits{Delay: 100 * time.Millisecond}) }()
	select {
	case r := <-second:
		if r != common.WokeDelay {
			t.Fatalf("wake reason = %v, want WokeDelay; the interrupt was reused", r)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("second wait never returned")
	}
}

// The registry walks every running job, but only the ones with waiters end.
// Without that, interrupting would disturb work the agent deliberately left
// running in the background.
func TestRegistryInterruptSkipsJobsWithNoWaiter(t *testing.T) {
	waited := runningJob()
	background := runningJob()
	js := &Jobs{all: map[int]*Job{1: waited, 2: background}}

	go func() { waited.Wait(common.Limits{Delay: longWait}) }()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		waited.mu.Lock()
		parked := waited.waiters > 0
		waited.mu.Unlock()
		if parked {
			break
		}
		time.Sleep(time.Millisecond)
	}

	if n := js.Interrupt(); n != 1 {
		t.Fatalf("interrupted %d jobs, want exactly the one with a waiter", n)
	}
	background.mu.Lock()
	flagged := background.interrupted
	background.mu.Unlock()
	if flagged {
		t.Fatal("a background job with no waiter was left holding an interrupt")
	}
}
