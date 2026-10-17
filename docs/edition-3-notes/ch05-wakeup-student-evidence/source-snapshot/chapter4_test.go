package ensemble_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"ensemble"
)

func jobsAgent(t *testing.T) (*ensemble.Ensemble, ensemble.Agent) {
	t.Helper()
	r := ensemble.New(log.New(io.Discard, "", 0))
	a, e := r.NewAgent(ensemble.Config{System: "test", MaxTokens: 100, EnableTools: true, Workspace: t.TempDir(), DataDir: t.TempDir()})
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() {
		if e := r.Shutdown(); e != nil {
			t.Error(e)
		}
	})
	return r, a
}
func dispatch(t *testing.T, a ensemble.Agent, name string, args any) string {
	t.Helper()
	b, e := json.Marshal(args)
	if e != nil {
		t.Fatal(e)
	}
	s, e := a.CallTool(context.Background(), name, b)
	if e != nil {
		t.Fatal(e)
	}
	return s
}
func jobAt(t *testing.T, r *ensemble.Ensemble, n int) ensemble.Job {
	t.Helper()
	j, e := r.Jobs().Get(n)
	if e != nil {
		t.Fatal(e)
	}
	return j
}
func register(t *testing.T, a ensemble.Agent, name string, h ensemble.ToolHandler) {
	t.Helper()
	if e := a.RegisterTool(ensemble.ToolDefinition{Name: name, Description: "test operation", Parameters: json.RawMessage(`{"type":"object","properties":{},"additionalProperties":false}`)}, h); e != nil {
		t.Fatal(e)
	}
}

func TestChapter4LimitsPersistenceAndCursor(t *testing.T) {
	r, a := jobsAgent(t)
	big := "HEAD-MARK" + strings.Repeat("x", 20000) + "TAIL-MARK"
	if e := os.WriteFile(filepath.Join(a.Config().Workspace, "large"), []byte(big), 0600); e != nil {
		t.Fatal(e)
	}
	dispatch(t, a, "tool_limits", map[string]any{"max_output_bytes": 1024, "ai_callback_delay": 1})
	out := dispatch(t, a, "read_file", map[string]string{"path": "large"})
	if !strings.Contains(out, "limits consumed") || !strings.Contains(out, "HEAD-MARK") || !strings.Contains(out, "TAIL-MARK") || !strings.Contains(out, "bytes omitted") || len(out) > 1500 {
		t.Fatal(out)
	}
	j := jobAt(t, r, 1)
	rec := j.Snapshot()
	if rec.Status != "done" || rec.Output.Locator != "cr/io/1" {
		t.Fatal(rec)
	}
	full, e := os.ReadFile(filepath.Join(a.Config().DataDir, rec.Output.Locator))
	if e != nil || string(full) != big {
		t.Fatalf("full output %d %v", len(full), e)
	}
	out = dispatch(t, a, "wait_for_job", map[string]any{"handle": 1})
	if strings.Contains(out, "HEAD-MARK") || strings.Contains(out, "TAIL-MARK") {
		t.Fatal("cursor replay", out)
	}
	out = dispatch(t, a, "read_file", map[string]any{"path": "cr/io/1", "max_bytes": 9})
	if out != "HEAD-MARK\n[truncated at 9 bytes]" {
		t.Fatal(out)
	}
	dispatch(t, a, "tool_limits", map[string]any{"max_output_bytes": 2048})
	out = dispatch(t, a, "unknown", map[string]any{})
	if !strings.Contains(out, "limits consumed") || !strings.Contains(out, "unknown tool") {
		t.Fatal(out)
	}
	out = dispatch(t, a, "read_file", map[string]string{"path": "missing"})
	if strings.Contains(out, "limits consumed") {
		t.Fatal("failed call did not consume", out)
	}
	for _, name := range []string{"wait_for_job", "send_input", "kill_job", "tool_limits"} {
		dispatch(t, a, "tool_limits", map[string]any{"max_output_bytes": 1000})
		out = dispatch(t, a, name, map[string]any{"bogus": true})
		if !strings.Contains(out, "limits consumed") || !strings.Contains(out, "error:") {
			t.Fatal(name, out)
		}
	}
	// All invalid observer options fail before a process can create a marker.
	for _, opts := range []map[string]any{{"ai_callback_delay": -1}, {"max_output_bytes": 0}, {"ai_callback_pattern": "["}, {"ai_callback_delay": nil}} {
		opts["command"] = "touch must-not-exist"
		out = dispatch(t, a, "run_command", opts)
		if !strings.Contains(out, "error:") {
			t.Fatal(out)
		}
	}
	if _, e := os.Stat(filepath.Join(a.Config().Workspace, "must-not-exist")); !os.IsNotExist(e) {
		t.Fatal(e)
	}
}
func TestChapter4RegisteredKillLateResultAndOwner(t *testing.T) {
	r, a := jobsAgent(t)
	started := make(chan ensemble.Job, 1)
	release := make(chan struct{})
	ended := make(chan struct{})
	register(t, a, "blocked", func(c ensemble.ToolContext, b json.RawMessage) (string, error) {
		if c.Agent != a || c.Job.Parent() != a || c.Agent.Parent() != r || r.Jobs().Parent() != r {
			t.Error("owner chain")
		}
		started <- c.Job
		<-release
		defer close(ended)
		return "LATE-RESULT", nil
	})
	if len(a.Tools()) != 11 {
		t.Fatal("builtins replaced")
	}
	if e := a.RegisterTool(a.Tools()[0].Definition(), func(ensemble.ToolContext, json.RawMessage) (string, error) { return "", nil }); e == nil {
		t.Fatal("duplicate accepted")
	}
	dispatch(t, a, "tool_limits", map[string]any{"ai_callback_delay": 0.01, "max_output_bytes": 1024})
	out := dispatch(t, a, "blocked", map[string]any{})
	if !strings.Contains(out, "job 1 running") || !strings.Contains(out, "limits consumed") {
		t.Fatal(out)
	}
	j := <-started
	wait := make(chan string, 1)
	go func() {
		s, e := j.Wait(context.Background(), ensemble.WaitOptions{Delay: time.Minute, MaxBytes: 1024})
		if e != nil {
			wait <- e.Error()
		} else {
			wait <- s
		}
	}()
	out = dispatch(t, a, "kill_job", map[string]any{"handle": 1})
	if !strings.Contains(out, "cannot stop the goroutine") {
		t.Fatal(out)
	}
	select {
	case out := <-wait:
		if !strings.Contains(out, "killed") {
			t.Fatal(out)
		}
	case <-time.After(time.Second):
		t.Fatal("blocked waiter not notified")
	}
	close(release)
	<-ended
	time.Sleep(10 * time.Millisecond)
	if j.Snapshot().Status != "killed" {
		t.Fatal("late resurrection")
	}
	data, e := os.ReadFile(filepath.Join(a.Config().DataDir, "cr/io/1"))
	if e != nil || len(data) != 0 {
		t.Fatalf("late output %q %v", data, e)
	}
	kills := 0
	for _, ev := range a.Log().Events() {
		if ev.Type == "job_killed" {
			kills++
			if ev.Job.Reason != "kill_job" {
				t.Fatal(ev)
			}
		}
	}
	if kills != 1 {
		t.Fatal(kills)
	}
	var dump bytes.Buffer
	if e := a.Log().Dump(&dump); e != nil {
		t.Fatal(e)
	}
	_, loaded := jobsAgent(t)
	if e := loaded.Log().Load(&dump); e != nil {
		t.Fatal(e)
	}
	if _, e := loaded.Parent().Jobs().Get(1); e == nil {
		t.Fatal("load rehydrated live job")
	}
}
func TestChapter4PTYInputCWDPatternAndShutdown(t *testing.T) {
	r, a := jobsAgent(t)
	if e := os.Mkdir(filepath.Join(a.Config().Workspace, "sub"), 0700); e != nil {
		t.Fatal(e)
	}
	out := dispatch(t, a, "run_command", map[string]any{"command": "printf 'term=%s\\n' \"$TERM\"; stty size; test -t 0 && test -t 1 && test -t 2; printf 'ready>'; read line; printf 'answer=%s\\n' \"$line\"; exit 7", "cwd": "sub", "ai_callback_pattern": "ready>", "ai_callback_delay": 1})
	if !strings.Contains(out, "running") || !strings.Contains(out, "term=dumb") || !strings.Contains(out, "50 200") || strings.Contains(out, "\r\n") || !strings.Contains(out, filepath.Join(a.Config().Workspace, "sub")) {
		t.Fatal(out)
	}
	start := time.Now()
	out = dispatch(t, a, "wait_for_job", map[string]any{"handle": 1, "ai_callback_pattern": "ready>", "ai_callback_delay": 0.04})
	if time.Since(start) < 30*time.Millisecond || strings.Contains(out, "ready>") {
		t.Fatal("matched/replayed consumed bytes", out)
	}
	out = dispatch(t, a, "send_input", map[string]any{"handle": 1, "input": "hello", "ai_callback_delay": 1})
	if !strings.Contains(out, "answer=hello") || !strings.Contains(out, "exit_code: 7") || !strings.Contains(out, "done") {
		t.Fatal(out)
	}
	if code := jobAt(t, r, 1).Snapshot().ExitCode; code == nil || *code != 7 {
		t.Fatal(code)
	}
	out = dispatch(t, a, "run_command", map[string]any{"command": "pwd; cd /; export PRIVATE_SHELL=yes"})
	if !strings.Contains(out, a.Config().Workspace) {
		t.Fatal(out)
	}
	out = dispatch(t, a, "run_command", map[string]any{"command": "pwd; printf 'local=%s' \"$PRIVATE_SHELL\""})
	if !strings.Contains(out, "local=\n") || !strings.Contains(out, a.Config().Workspace) {
		t.Fatal(out)
	}
	out = dispatch(t, a, "run_command", map[string]any{"command": "sleep 60 & child=$!; printf 'child=%s\\n' \"$child\"; wait", "ai_callback_pattern": "child=[0-9]+", "ai_callback_delay": 1})
	startPID := strings.Index(out, "child=")
	if startPID < 0 {
		t.Fatal(out)
	}
	pidText := strings.Fields(out[startPID+6:])[0]
	pid, e := strconv.Atoi(pidText)
	if e != nil {
		t.Fatal(e)
	}
	j := jobAt(t, r, 4)
	if j.Snapshot().Status != "running" {
		t.Fatal("shutdown proof requires running job")
	}
	if e := r.Shutdown(); e != nil {
		t.Fatal(e)
	}
	if j.Snapshot().Status != "killed" {
		t.Fatal(j.Snapshot())
	}
	deadline := time.Now().Add(time.Second)
	for syscall.Kill(pid, 0) == nil && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	if e := syscall.Kill(pid, 0); e != syscall.ESRCH {
		t.Fatalf("descendant %d survives: %v", pid, e)
	}
	found := false
	for _, ev := range a.Log().Events() {
		if ev.Type == "job_killed" && ev.Job.Handle == 4 && ev.Job.Reason == "shutdown" {
			found = true
		}
	}
	if !found {
		t.Fatal("no shutdown event")
	}
	if _, e := r.Jobs().Create(a); e == nil {
		t.Fatal("job after shutdown")
	}
}
func TestChapter4RootIsolationAndMetadata(t *testing.T) {
	r, a := jobsAgent(t)
	other, b := jobsAgent(t)
	if r.Jobs() == other.Jobs() {
		t.Fatal("global registry")
	}
	dispatch(t, a, "write_file", map[string]string{"path": "one", "content": "content"})
	dispatch(t, b, "write_file", map[string]string{"path": "two", "content": "content"})
	if jobAt(t, r, 1).Parent() != a || jobAt(t, other, 1).Parent() != b {
		t.Fatal("wrong job parent")
	}
	for _, ev := range a.Log().Events() {
		if ev.Type == "tool_called" || ev.Type == "tool_returned" {
			if ev.Tool.Job == nil || ev.Tool.Job.Handle != 1 || ev.Tool.Job.Output.Locator != "cr/io/1" {
				t.Fatal(ev)
			}
			if ev.Type == "tool_returned" && ev.Tool.Job.Status != "done" {
				t.Fatal(ev.Tool.Job)
			}
			ev.Tool.Job.Handle = 99
		}
	}
	for _, ev := range a.Log().Events() {
		if ev.Tool != nil && ev.Tool.Job.Handle != 1 {
			t.Fatal("snapshot alias")
		}
	}
	register(t, a, "fail", func(ensemble.ToolContext, json.RawMessage) (string, error) {
		return "partial", fmt.Errorf("deliberate")
	})
	out := dispatch(t, a, "fail", map[string]any{})
	if !strings.Contains(out, "partial") || !strings.Contains(out, "error: deliberate") {
		t.Fatal(out)
	}
}

func TestChapter4PersistentFilesAndConcurrentShutdown(t *testing.T) {
	workspace, data := t.TempDir(), t.TempDir()
	makeRoot := func() (*ensemble.Ensemble, ensemble.Agent) {
		r := ensemble.New(log.New(io.Discard, "", 0))
		a, e := r.NewAgent(ensemble.Config{System: "test", MaxTokens: 100, EnableTools: true, Workspace: workspace, DataDir: data})
		if e != nil {
			t.Fatal(e)
		}
		t.Cleanup(func() {
			if e := r.Shutdown(); e != nil {
				t.Error(e)
			}
		})
		return r, a
	}
	r, a := makeRoot()
	dispatch(t, a, "write_file", map[string]string{"path": "one", "content": "first"})
	before, e := os.ReadFile(filepath.Join(data, "cr/io/1"))
	if e != nil {
		t.Fatal(e)
	}
	if e := r.Shutdown(); e != nil {
		t.Fatal(e)
	}
	r2, b := makeRoot()
	dispatch(t, b, "read_file", map[string]string{"path": "one"})
	after, e := os.ReadFile(filepath.Join(data, "cr/io/1"))
	if e != nil || !bytes.Equal(before, after) {
		t.Fatal("output overwritten", e)
	}
	if jobAt(t, r2, 2).Snapshot().Output.Locator != "cr/io/2" {
		t.Fatal("reused output path")
	}
	// Race actual process launch against two root shutdown callers, not a fake process.
	for i := 0; i < 12; i++ {
		root, agent := makeRoot()
		job, e := root.Jobs().Create(agent)
		if e != nil {
			t.Fatal(e)
		}
		var shell ensemble.Tool
		for _, tool := range agent.Tools() {
			if tool.Definition().Name == "run_command" {
				shell = tool
			}
		}
		job.Start(context.Background(), shell, json.RawMessage(`{"command":"sleep 60"}`))
		done := make(chan error, 2)
		go func() { done <- root.Shutdown() }()
		go func() { done <- root.Shutdown() }()
		for n := 0; n < 2; n++ {
			select {
			case e := <-done:
				if e != nil {
					t.Fatal(e)
				}
			case <-time.After(3 * time.Second):
				t.Fatal("shutdown did not reap")
			}
		}
		if job.Snapshot().Status != "killed" {
			t.Fatal(job.Snapshot())
		}
	}
}
