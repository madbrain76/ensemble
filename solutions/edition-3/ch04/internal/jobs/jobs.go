// Package jobs owns execution, output cursors, and supervised lifetimes.
package jobs

import (
	"bufio"
	"context"
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

type collection struct{ privateCollection }
type privateCollection struct {
	parent common.Agent
	mu     sync.Mutex
	closed bool
}

func New(parent common.Agent) common.Jobs {
	return &collection{privateCollection: privateCollection{parent: parent}}
}
func (c *collection) Parent() common.Agent { return c.parent }
func (c *collection) Get(handle int) (common.Job, error) {
	j, ok := c.parent.Parent().LookupJob(handle)
	if !ok || j.Parent() != c {
		return nil, fmt.Errorf("unknown job handle %d for this Agent", handle)
	}
	return j, nil
}
func (c *collection) Shutdown() error {
	c.mu.Lock()
	c.closed = true
	c.mu.Unlock()
	var errs []error
	for _, j := range c.parent.Parent().JobTable() {
		if j.Parent() == c {
			errs = append(errs, j.Kill("shutdown"))
		}
	}
	return errors.Join(errs...)
}

type job struct{ privateJob }
type privateJob struct {
	parent     *collection
	mu         sync.Mutex
	data       common.JobData
	file       *os.File
	cursor     int64
	changed    chan struct{}
	settled    chan struct{}
	cancel     context.CancelFunc
	input      io.Writer
	stop       func() error
	failed     bool
	setupError string
}

func (j *job) Parent() common.Jobs { return j.parent }
func (j *job) snapshot() common.JobData {
	d := j.data
	if d.ExitCode != nil {
		n := *d.ExitCode
		d.ExitCode = &n
	}
	return d
}
func (j *job) Data() common.JobData { j.mu.Lock(); defer j.mu.Unlock(); return j.snapshot() }
func (j *job) wake()                { close(j.changed); j.changed = make(chan struct{}) }

// Start registers before executing or returning. A killed job cannot later start work.
func (c *collection) Start(ctx context.Context, handler common.ToolHandler, args []byte, cwd string, called func(common.JobData) error) (common.Job, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return nil, errors.New("Agent jobs are shut down")
	}
	runCtx, cancel := context.WithCancel(ctx)
	j := &job{privateJob: privateJob{parent: c, changed: make(chan struct{}), settled: make(chan struct{}), cancel: cancel, data: common.JobData{Status: common.JobRunning}}}
	j.mu.Lock()
	handle, err := c.parent.Parent().RegisterJob(j)
	if err != nil {
		j.mu.Unlock()
		cancel()
		return nil, err
	}
	j.data.Handle = handle
	j.data.Cwd = cwd
	config := c.parent.Config()
	locator := filepath.Join(config.DataDir, "io", fmt.Sprint(handle))
	path := locator
	if !filepath.IsAbs(path) {
		path = filepath.Join(config.BaseDir, path)
	}
	j.data.Output = common.Ref{Kind: common.RefHandle, Locator: filepath.ToSlash(locator)}
	if err = os.MkdirAll(filepath.Dir(path), 0700); err == nil {
		j.file, err = os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0600)
	}
	if err != nil {
		j.setupError = "preparing job output: " + err.Error()
		j.data.Output = common.Ref{}
	}
	if err == nil && called != nil {
		err = called(j.snapshot())
	}
	if err != nil {
		j.failed = true
		j.data.Status = common.JobDone
		if j.file != nil {
			_ = j.file.Close()
			j.file = nil
		}
		j.wake()
		close(j.settled)
		j.mu.Unlock()
		cancel()
		return j, err
	}
	j.mu.Unlock()
	input := append([]byte(nil), args...)
	go func() {
		if err := runCtx.Err(); err != nil {
			j.finish(err)
			return
		}
		err := handler(runCtx, &call{parent: c.parent, job: j}, input)
		j.finish(err)
	}()
	return j, nil
}

// call carries its real Agent parent, not a copied sibling-service bag.
type call struct {
	parent common.Agent
	job    *job
}

func (c *call) Parent() common.Agent        { return c.parent }
func (c *call) Write(p []byte) (int, error) { return c.job.write(p) }
func (c *call) SetExitCode(code int) {
	c.job.mu.Lock()
	defer c.job.mu.Unlock()
	if c.job.data.Status == common.JobRunning {
		c.job.data.ExitCode = &code
	}
}
func (c *call) StartProcess(start func() (io.Writer, func() error, error)) error {
	j := c.job
	j.mu.Lock()
	defer j.mu.Unlock()
	if j.data.Status != common.JobRunning {
		return errors.New("job no longer running")
	}
	input, stop, err := start()
	if err != nil {
		return err
	}
	j.input, j.stop = input, stop
	return nil
}
func (j *job) write(p []byte) (int, error) {
	j.mu.Lock()
	defer j.mu.Unlock()
	if j.data.Status != common.JobRunning {
		return len(p), nil
	}
	n, err := j.file.Write(p)
	j.data.Bytes += int64(n)
	if err != nil {
		j.failed = true
	}
	j.wake()
	return n, err
}
func (j *job) finish(err error) {
	j.mu.Lock()
	defer j.mu.Unlock()
	if j.data.Status != common.JobRunning {
		return
	}
	if err != nil {
		j.failed = true
		message := err.Error()
		if j.data.Bytes > 0 {
			message = "\n" + message
		}
		n, _ := j.file.WriteString(message)
		j.data.Bytes += int64(n)
	}
	j.data.Status = common.JobDone
	j.cancel()
	// Observations open bounded read views; completion closes the writer.
	j.closeWriter()
	close(j.settled)
	j.wake()
}
func (j *job) closeWriter() {
	if j.file != nil {
		if err := j.file.Close(); err != nil {
			j.failed = true
		}
		j.file = nil
	}
}
func (j *job) Kill(reason string) error {
	if reason != "kill_job" && reason != "shutdown" {
		return errors.New("invalid kill reason")
	}
	j.mu.Lock()
	if j.data.Status != common.JobRunning {
		settled := j.settled
		j.mu.Unlock()
		<-settled
		return nil
	}
	j.data.Status = common.JobKilledStatus
	j.data.ExitCode = nil
	j.data.Reason = reason
	j.cancel()
	stop := j.stop
	j.closeWriter()
	j.wake()
	snapshot := j.snapshot()
	j.mu.Unlock()
	var stopErr error
	if stop != nil {
		stopErr = stop()
	}
	defer close(j.settled)
	snapshot = j.Data()
	logErr := j.parent.parent.Log().Append(common.Event{Type: common.JobKilled, Job: &snapshot})
	return errors.Join(stopErr, logErr)
}
func (j *job) outputPath() string {
	path := filepath.FromSlash(j.data.Output.Locator)
	if !filepath.IsAbs(path) {
		path = filepath.Join(j.parent.parent.Config().BaseDir, path)
	}
	return path
}

// observeLocked reads at most the payload cap; the cursor covers omitted bytes too.
func (j *job) observeLocked(l common.ToolLimits) (common.JobReport, error) {
	report := common.JobReport{Job: j.snapshot(), IsError: j.failed}
	if j.setupError != "" {
		report.Text = j.setupError
		return report, nil
	}
	count := j.data.Bytes - j.cursor
	if count == 0 {
		return report, nil
	}
	f, err := os.Open(j.outputPath())
	if err != nil {
		return report, err
	}
	defer f.Close()
	read := func(offset int64, n int) (string, error) {
		b := make([]byte, n)
		_, err := f.ReadAt(b, offset)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}
	if count <= int64(l.MaxOutputBytes) {
		report.Text, err = read(j.cursor, int(count))
	} else {
		head := l.MaxOutputBytes / 2
		tail := l.MaxOutputBytes - head
		var a, b string
		a, err = read(j.cursor, head)
		if err == nil {
			b, err = read(j.data.Bytes-int64(tail), tail)
		}
		report.Text = fmt.Sprintf("%s\n[output truncated; %d bytes omitted; %d bytes total; full output: %s]\n%s", a, count-int64(l.MaxOutputBytes), j.data.Bytes, j.data.Output.Locator, b)
	}
	if err == nil {
		j.cursor = j.data.Bytes
	}
	return report, err
}

// InitialReport preserves fast completed calls' original bytes. Without a
// pattern it waits for completion or the observation delay, not the first chunk.
func (j *job) InitialReport(ctx context.Context, l common.ToolLimits) (common.JobReport, error) {
	return j.wait(ctx, l, false)
}
func (j *job) Wait(ctx context.Context, l common.ToolLimits) (common.JobReport, error) {
	return j.wait(ctx, l, true)
}
func (j *job) wait(ctx context.Context, l common.ToolLimits, outputWake bool) (common.JobReport, error) {
	if l.Delay < 0 || l.MaxOutputBytes <= 0 {
		return common.JobReport{}, errors.New("invalid observation limits")
	}
	var pattern *regexp.Regexp
	var err error
	if l.Pattern != "" {
		pattern, err = regexp.Compile(l.Pattern)
		if err != nil {
			return common.JobReport{}, fmt.Errorf("invalid callback pattern: %w", err)
		}
	}
	timer := time.NewTimer(l.Delay)
	defer timer.Stop()
	expired := l.Delay == 0
	for {
		j.mu.Lock()
		ready := expired || j.data.Status != common.JobRunning
		if !ready && j.data.Bytes > j.cursor {
			if pattern == nil {
				ready = outputWake
			} else {
				f, e := os.Open(j.outputPath())
				if e != nil {
					j.mu.Unlock()
					return common.JobReport{}, e
				}
				ready = pattern.MatchReader(bufio.NewReader(io.NewSectionReader(f, j.cursor, j.data.Bytes-j.cursor)))
				_ = f.Close()
			}
		}
		if ready {
			result, e := j.observeLocked(l)
			j.mu.Unlock()
			return result, e
		}
		changed := j.changed
		j.mu.Unlock()
		select {
		case <-changed:
		case <-timer.C:
			expired = true
		case <-ctx.Done():
			return common.JobReport{}, ctx.Err()
		}
	}
}
func (j *job) SendInput(ctx context.Context, text string, l common.ToolLimits) (common.JobReport, error) {
	j.mu.Lock()
	input, status := j.input, j.data.Status
	j.mu.Unlock()
	if input == nil {
		return common.JobReport{}, errors.New("job does not support input")
	}
	if status != common.JobRunning {
		return common.JobReport{}, errors.New("job is not running")
	}
	// Submit one terminal line; a caller-supplied terminator is not doubled.
	if !strings.HasSuffix(text, "\n") && !strings.HasSuffix(text, "\r") {
		text += "\n"
	}
	if _, err := io.WriteString(input, text); err != nil {
		return common.JobReport{}, fmt.Errorf("sending input: %w", err)
	}
	// As with a new command, allow the response to arrive after terminal echo.
	// An explicit pattern still wakes promptly; subsequent Wait calls return
	// available unseen output immediately.
	return j.InitialReport(ctx, l)
}
