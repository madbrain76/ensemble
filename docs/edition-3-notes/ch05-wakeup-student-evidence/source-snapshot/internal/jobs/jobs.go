// Package jobs owns invocation lifetimes and persistent consumable output.
package jobs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"ensemble/internal/common"
)

// Table belongs to a single Ensemble, not a global handle namespace.
type Table struct {
	parent common.Ensemble
	mu     sync.Mutex
	next   int
	jobs   map[int]*Job
	closed bool
}

func New(parent common.Ensemble) *Table  { return &Table{parent: parent, jobs: make(map[int]*Job)} }
func (t *Table) Parent() common.Ensemble { return t.parent }
func (t *Table) Create(a common.Agent) (common.Job, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.closed {
		return nil, fmt.Errorf("ensemble is shut down")
	}
	if a == nil || a.Parent() != t.parent {
		return nil, fmt.Errorf("job requires an agent owned by this root")
	}
	var locator, path string
	var f *os.File
	for {
		t.next++
		locator = fmt.Sprintf("cr/io/%d", t.next)
		path = filepath.Join(a.Config().DataDir, locator)
		if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			return nil, err
		}
		var err error
		f, err = os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0600)
		if errors.Is(err, os.ErrExist) {
			continue
		}
		if err != nil {
			return nil, err
		}
		break
	}
	j := &Job{parent: a, record: common.JobRecord{Handle: t.next, Status: "running", Output: common.Ref{Kind: common.RefHandle, Locator: locator}}, file: f, path: path, changed: make(chan struct{}), cleaned: make(chan struct{})}
	t.jobs[t.next] = j
	return j, nil
}
func (t *Table) Get(handle int) (common.Job, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	j, ok := t.jobs[handle]
	if !ok {
		return nil, fmt.Errorf("unknown job handle %d", handle)
	}
	return j, nil
}
func (t *Table) Shutdown() error {
	t.mu.Lock()
	t.closed = true
	jobs := make([]*Job, 0, len(t.jobs))
	for _, j := range t.jobs {
		jobs = append(jobs, j)
	}
	t.mu.Unlock()
	var first error
	for _, j := range jobs {
		j.mu.Lock()
		process := j.processTool || j.process
		started := j.started
		j.mu.Unlock()
		if !started {
			j.Finish("", fmt.Errorf("shutdown before start"), nil)
		}
		var err error
		if process {
			err = j.Kill("shutdown")
		}
		if err != nil && first == nil {
			first = err
		}
	}
	// Stop every managed process before waiting for unrelated ordinary handlers.
	// Their true returns still own completion and output; shutdown does not fake it.
	for _, j := range jobs {
		<-j.cleaned
		j.execution.Wait()
	}
	return first
}

// Job serializes all output, cursor and terminal transitions behind one mutex.
type Job struct {
	parent      common.Agent
	mu          sync.Mutex
	record      common.JobRecord
	file        *os.File
	path        string
	cursor      int64
	reports     int
	changed     chan struct{}
	cleaned     chan struct{}
	started     bool
	execution   sync.WaitGroup
	processTool bool
	cancel      context.CancelFunc
	input       io.WriteCloser
	stop        func() error
	process     bool
	failure     error
}

// Done closes after terminal process/output cleanup, not merely cancellation.
func (j *Job) Done() <-chan struct{} { return j.cleaned }
func (j *Job) Parent() common.Agent  { return j.parent }
func (j *Job) snapshot() common.JobRecord {
	r := j.record
	if r.ExitCode != nil {
		v := *r.ExitCode
		r.ExitCode = &v
	}
	return r
}
func (j *Job) Snapshot() common.JobRecord { j.mu.Lock(); defer j.mu.Unlock(); return j.snapshot() }
func (j *Job) notify()                    { close(j.changed); j.changed = make(chan struct{}) }
func (j *Job) Start(ctx context.Context, tool common.Tool, args json.RawMessage) {
	ctx, cancel := context.WithCancel(ctx)
	j.mu.Lock()
	if j.record.Status != "running" || j.started {
		j.mu.Unlock()
		cancel()
		return
	}
	j.started = true
	j.processTool = tool != nil && tool.Definition().Name == "run_command"
	j.execution.Add(1)
	j.cancel = cancel
	j.mu.Unlock()
	go func() { <-ctx.Done(); _ = j.Kill("cancellation") }()
	go func() {
		defer j.execution.Done()
		if tool == nil {
			j.Finish("", fmt.Errorf("unknown tool"), nil)
			return
		}
		if async, ok := tool.(common.JobTool); ok {
			out, owned, err := async.Start(ctx, j, args)
			if !owned {
				j.Finish(out, err, nil)
			}
			return
		}
		out, err := tool.Execute(ctx, args)
		j.Finish(out, err, nil)
	}()
}
func (j *Job) Write(p []byte) (int, error) {
	j.mu.Lock()
	defer j.mu.Unlock()
	if j.record.Status != "running" {
		return 0, fmt.Errorf("job is %s", j.record.Status)
	}
	n, err := j.file.Write(p)
	j.record.Bytes += int64(n)
	j.notify()
	return n, err
}
func (j *Job) SetCWD(cwd string) {
	j.mu.Lock()
	defer j.mu.Unlock()
	if j.record.Status == "running" {
		j.record.CWD = cwd
	}
}

// LaunchProcess makes process creation and attachment atomic against Kill.
func (j *Job) LaunchProcess(start func() (io.WriteCloser, func() error, error)) error {
	j.mu.Lock()
	defer j.mu.Unlock()
	if j.record.Status != "running" {
		return fmt.Errorf("job is %s", j.record.Status)
	}
	input, stop, err := start()
	if err != nil {
		return err
	}
	j.input, j.stop, j.process = input, stop, true
	return nil
}

func (j *Job) Finish(out string, err error, code *int) {
	j.mu.Lock()
	defer j.mu.Unlock()
	if j.record.Status != "running" {
		return
	}
	if err != nil {
		if out != "" {
			out += "\n"
		}
		out += "error: " + err.Error()
	}
	if out != "" {
		n, writeErr := io.WriteString(j.file, out)
		j.record.Bytes += int64(n)
		if writeErr != nil {
			err = writeErr
		}
	}
	j.failure = err
	j.record.Status = "done"
	close(j.cleaned)
	j.record.ExitCode = code
	if closeErr := j.file.Close(); err == nil && closeErr != nil {
		j.failure = closeErr
	}
	if j.cancel != nil {
		j.cancel()
	}
	j.notify()
}
func (j *Job) Input(input string) error {
	j.mu.Lock()
	if j.record.Status != "running" || j.input == nil {
		j.mu.Unlock()
		return fmt.Errorf("job is not an interactive running process")
	}
	terminal := j.input
	j.mu.Unlock()
	if !strings.HasSuffix(input, "\n") {
		input += "\n"
	}
	_, err := io.WriteString(terminal, input)
	return err
}

func (j *Job) Kill(reason string) error {
	j.mu.Lock()
	if j.record.Status != "running" {
		j.mu.Unlock()
		<-j.cleaned
		return nil
	}
	j.record.Status = "killed"
	j.record.Reason = reason
	if j.cancel != nil {
		j.cancel()
	}
	stop := j.stop
	closeErr := j.file.Close()
	record := j.snapshot()
	j.notify()
	j.mu.Unlock()
	var stopErr error
	if stop != nil {
		stopErr = stop()
	}
	defer close(j.cleaned)
	var logErr error
	if owner, ok := j.parent.(interface{ RecordJob(common.JobRecord) }); ok {
		owner.RecordJob(record)
	} else {
		logErr = j.parent.Log().Append(common.Event{Type: "job_killed", Job: &record})
	}
	if stopErr != nil {
		return stopErr
	}
	if closeErr != nil {
		return closeErr
	}
	return logErr
}
func (j *Job) Wait(ctx context.Context, opts common.WaitOptions) (string, error) {
	if opts.Delay < 0 || opts.MaxBytes <= 0 {
		return "", fmt.Errorf("invalid observation limits")
	}
	var pattern *regexp.Regexp
	if opts.Pattern != "" {
		var err error
		pattern, err = regexp.Compile(opts.Pattern)
		if err != nil {
			return "", err
		}
	}
	timer := time.NewTimer(opts.Delay)
	defer timer.Stop()
	firstSnapshot := true
	for {
		j.mu.Lock()
		data, err := os.ReadFile(j.path)
		if err != nil {
			j.mu.Unlock()
			return "", err
		}
		unseen := data[j.cursor:]
		woken := context.Cause(ctx) == common.ErrObservationWoken
		// A resumed wait may return bytes already present at admission. New output
		// still notifies/rechecks, but does not bypass its pattern/silence limit.
		// Wake and completion use this same locked cursor/report path.
		ready := woken || j.record.Status != "running" || (opts.ReturnAvailable && firstSnapshot && len(unseen) > 0) || (pattern != nil && pattern.Match(unseen))
		firstSnapshot = false
		changed := j.changed
		if ready {
			out, err := j.report(unseen, opts, woken)
			j.mu.Unlock()
			return out, err
		}
		j.mu.Unlock()
		select {
		case <-changed:
			continue
		case <-timer.C:
			j.mu.Lock()
			data, err = os.ReadFile(j.path)
			if err != nil {
				j.mu.Unlock()
				return "", err
			}
			out, err := j.report(data[j.cursor:], opts, context.Cause(ctx) == common.ErrObservationWoken)
			j.mu.Unlock()
			return out, err
		case <-ctx.Done():
			if context.Cause(ctx) == common.ErrObservationWoken {
				continue
			}
			return "", ctx.Err()
		}
	}
}
func (j *Job) report(unseen []byte, opts common.WaitOptions, identified bool) (string, error) {
	first := j.reports == 0
	j.reports++
	j.cursor = j.record.Bytes
	if first && !identified && j.record.Status == "done" && len(unseen) <= opts.MaxBytes && j.record.CWD == "" {
		return string(unseen), j.failure
	}
	header := fmt.Sprintf("job %d %s; %d bytes; output %s", j.record.Handle, j.record.Status, j.record.Bytes, j.record.Output.Locator)
	if j.record.CWD != "" {
		header += "; cwd " + j.record.CWD
	}
	if j.record.Status == "killed" && !j.process {
		header += "; Go cannot stop the goroutine; late completion is discarded"
	}
	payload := string(unseen)
	if len(header)+1+len(unseen) > opts.MaxBytes {
		budget := opts.MaxBytes - len(header) - 128
		if budget < 0 {
			budget = 0
		}
		if budget > len(unseen) {
			budget = len(unseen)
		}
		head, tail := budget/2, budget-budget/2
		payload = string(unseen[:head]) + fmt.Sprintf("\n[%d bytes omitted; total %d bytes; full output %s]\n", len(unseen)-budget, j.record.Bytes, j.record.Output.Locator) + string(unseen[len(unseen)-tail:])
	}
	return header + "\n" + payload, j.failure
}
