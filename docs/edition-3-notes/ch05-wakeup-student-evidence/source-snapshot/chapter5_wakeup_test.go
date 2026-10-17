package ensemble_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"ensemble"
)

func wakeReceive[T any](t *testing.T, ch <-chan T, label string) T {
	t.Helper()
	select {
	case v := <-ch:
		return v
	case <-time.After(2 * time.Second):
		t.Fatalf("WAKEUP %s did not arrive promptly", label)
		var zero T
		return zero
	}
}
func wakeCalled(t *testing.T, a ensemble.Agent, id string) ensemble.ToolEvent {
	t.Helper()
	until := time.Now().Add(2 * time.Second)
	for time.Now().Before(until) {
		for _, e := range a.Log().Events() {
			if e.Type == "tool_called" && e.Tool.CallID == id {
				return *e.Tool
			}
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatalf("WAKEUP call %s not admitted", id)
	return ensemble.ToolEvent{}
}
func wakeResults(t *testing.T, raw string, ids ...string) map[string]struct {
	Text  string
	Error bool
} {
	t.Helper()
	var body struct {
		Messages []struct {
			Content []struct {
				Type    string          `json:"type"`
				ID      string          `json:"tool_use_id"`
				Content json.RawMessage `json:"content"`
				Error   bool            `json:"is_error"`
			}
		}
	}
	if err := json.Unmarshal([]byte(raw), &body); err != nil {
		t.Fatal(err)
	}
	results := map[string]struct {
		Text  string
		Error bool
	}{}
	for _, m := range body.Messages {
		for _, p := range m.Content {
			if p.Type == "tool_result" {
				if _, ok := results[p.ID]; ok {
					t.Fatalf("duplicate result %s", p.ID)
				}
				results[p.ID] = struct {
					Text  string
					Error bool
				}{string(p.Content), p.Error}
			}
		}
	}
	for _, id := range ids {
		if _, ok := results[id]; !ok {
			t.Fatalf("WAKEUP missing paired result %s in %s", id, raw)
		}
	}
	return results
}
func wakeDefinition(name string) ensemble.ToolDefinition {
	return ensemble.ToolDefinition{Name: name, Description: name, Parameters: json.RawMessage(`{"type":"object"}`)}
}
func wakeAddressed(t *testing.T) ensemble.UserMessage {
	t.Helper()
	var m ensemble.UserMessage
	// JSON keeps this regression compilable on the previous public API: losing
	// the newly supplied fields must fail behaviorally, not at build time.
	if err := json.Unmarshal([]byte(`{"Text":"CHILD_QUESTION_CAN_I_WRITE","sender":"child:writer","recipient":"parent:editor","request_id":"request-17","turn_id":"turn-23"}`), &m); err != nil {
		t.Fatal(err)
	}
	return m
}

func TestChapter5WakeupInitialMixedBatch(t *testing.T) {
	held, release := releaseOnce()
	t.Cleanup(release)
	started := make(chan ensemble.ToolContext, 1)
	next := make(chan string, 1)
	var requests, sideEffects atomic.Int32
	a, root := actorFixture(t, func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		var parts []ensemble.Part
		switch requests.Add(1) {
		case 1:
			parts = []ensemble.Part{{Type: "opaque"}, prose("I will inspect then join."), callPart("done", "done", map[string]any{}), callPart("limit", "tool_limits", map[string]any{"ai_callback_delay": 3600, "max_output_bytes": 512}), callPart("join", "join_like", map[string]any{}), callPart("skip", "effect", map[string]any{}), callPart("skip-setter", "tool_limits", map[string]any{"ai_callback_delay": 3600, "max_output_bytes": 1})}
		case 2:
			next <- string(raw)
			parts = []ensemble.Part{prose("I can answer the child while the join continues."), callPart("answer", "effect", map[string]any{})}
		default:
			parts = []ensemble.Part{prose("child answered")}
		}
		json.NewEncoder(w).Encode(providerReply(ensemble.Anthropic, parts))
	}, func(c *ensemble.Config) { c.EnableTools = true; c.ResolvedModel = "test-model" })
	t.Cleanup(release)
	a.RegisterTool(wakeDefinition("done"), func(ensemble.ToolContext, json.RawMessage) (string, error) { return "COMPLETED_REAL", nil })
	a.RegisterTool(wakeDefinition("join_like"), func(c ensemble.ToolContext, _ json.RawMessage) (string, error) {
		started <- c
		<-held
		return "FINAL_JOIN:" + strings.Repeat("x", 2000), nil
	})
	a.RegisterTool(wakeDefinition("effect"), func(ensemble.ToolContext, json.RawMessage) (string, error) {
		sideEffects.Add(1)
		return "ANSWER_SENT", nil
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	reply := a.Submit(ctx, "join child")
	jobContext := wakeReceive(t, started, "ordinary handler start")
	id := jobContext.Job.Snapshot().Handle
	if err := a.Post(wakeAddressed(t)); err != nil {
		t.Fatal(err)
	}
	raw := wakeReceive(t, next, "same-turn initial observation continuation")
	for _, want := range []string{"CHILD_QUESTION_CAN_I_WRITE", "child:writer", "parent:editor", "request-17", "turn-23"} {
		if strings.Count(raw, want) != 1 {
			t.Fatalf("WAKEUP actual next HTTP lost/duplicated %s: %s", want, raw)
		}
	}
	results := wakeResults(t, raw, "done", "limit", "join", "skip", "skip-setter")
	if strings.Count(raw, "fixture-signature") != 1 {
		t.Fatal("original signed response block was rewritten or lost")
	}
	if !strings.Contains(results["done"].Text, "COMPLETED_REAL") || results["done"].Error {
		t.Fatal("completed real result lost")
	}
	if results["join"].Error || !strings.Contains(results["join"].Text, fmt.Sprintf("job %d running", id)) {
		t.Fatalf("WAKEUP untruthful active job result: %+v", results["join"])
	}
	for _, id := range []string{"skip", "skip-setter"} {
		if !results[id].Error || !strings.Contains(results[id].Text, "not executed") {
			t.Fatalf("unpaired skipped call %s: %+v", id, results[id])
		}
	}
	if r := wait5(t, reply); r.Err != nil {
		t.Fatal(r)
	}
	if sideEffects.Load() != 1 {
		t.Fatalf("WAKEUP undispatched effect ran or reissue lost: %d", sideEffects.Load())
	}
	if jobContext.Job.Snapshot().Status != "running" {
		t.Fatal("wake canceled underlying join")
	}
	select {
	case <-jobContext.Done():
		t.Fatal("wake canceled handler context")
	default:
	}
	for _, e := range a.Log().Events() {
		if e.Type == "tool_called" && (e.Tool.CallID == "skip" || e.Tool.CallID == "skip-setter") {
			t.Fatal("undispatched tail was executed")
		}
	}
	// The skipped setter cannot arm limits on the newly issued answer.
	for _, e := range a.Log().Events() {
		if e.Type == "tool_returned" && e.Tool.CallID == "answer" && strings.Contains(e.Tool.Parts[0].Text, "tool_limits consumed") {
			t.Fatal("undispatched setter leaked")
		}
	}
	release()
	<-jobContext.Job.Done()
	out, err := a.CallTool(context.Background(), "wait_for_job", json.RawMessage(fmt.Sprintf(`{"handle":%d,"ai_callback_delay":3600,"max_output_bytes":512}`, id)))
	if err != nil || !strings.Contains(out, fmt.Sprintf("job %d done", id)) || !strings.Contains(out, "omitted") {
		t.Fatalf("same retained job final report: %s %v", out, err)
	}
	full, err := os.ReadFile(filepath.Join(a.Config().DataDir, jobContext.Job.Snapshot().Output.Locator))
	if err != nil || !strings.Contains(string(full), "FINAL_JOIN:") || len(full) < 2000 {
		t.Fatalf("full retained output lost: %v", err)
	}
	repeat, _ := a.CallTool(context.Background(), "wait_for_job", json.RawMessage(fmt.Sprintf(`{"handle":%d,"ai_callback_delay":3600}`, id)))
	if strings.Contains(repeat, "FINAL_JOIN:") || !strings.Contains(repeat, fmt.Sprintf("job %d done", id)) {
		t.Fatal("terminal repeat/cursor contract", repeat)
	}
	if _, err := root.Jobs().Get(id); err != nil {
		t.Fatal(err)
	}
	if _, err := root.Jobs().Get(id + 2); err == nil {
		t.Fatal("control or undispatched call created a new job")
	}

	saved, err := os.ReadFile(a.Config().HistoryPath)
	if err != nil || !bytes.Contains(saved, []byte(`"sender":"child:writer"`)) || !bytes.Contains(saved, []byte("child answered")) {
		t.Fatalf("saved message/correlation/result lost: %v", err)
	}
}

func TestChapter5WakeupSetterDispatchBoundary(t *testing.T) {
	var requests, effects atomic.Int32
	next := make(chan string, 1)
	a, _ := actorFixture(t, func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		var p []ensemble.Part
		if requests.Add(1) == 1 {
			p = []ensemble.Part{prose("I will act."), callPart("before", "effect", map[string]any{}), callPart("armed", "tool_limits", map[string]any{"ai_callback_delay": 3600, "max_output_bytes": 73}), callPart("spent", "effect", map[string]any{}), callPart("unarmed", "tool_limits", map[string]any{"ai_callback_delay": 3600})}
		} else {
			next <- string(raw)
			p = []ensemble.Part{prose("woke")}
		}
		json.NewEncoder(w).Encode(providerReply(ensemble.Anthropic, p))
	}, nil)
	a.RegisterTool(wakeDefinition("effect"), func(ensemble.ToolContext, json.RawMessage) (string, error) { effects.Add(1); return "real effect", nil })
	var once sync.Once
	// Publishing the real preceding result happens on the owner, before it starts
	// the setter. The hint is queued before that setter's worker can report.
	a.Observe(func(e ensemble.Observation) {
		if p, ok := e.(ensemble.PartFinal); ok && p.Part.Type == "tool_result" && p.Part.CallID == "before" {
			once.Do(func() { _ = a.Post(ensemble.Hint{Text: "SETTER_WAKE"}) })
		}
	})
	reply := a.Submit(context.Background(), "begin")
	raw := wakeReceive(t, next, "setter dispatch boundary")
	results := wakeResults(t, raw, "before", "armed", "spent", "unarmed")
	if results["armed"].Error || !strings.Contains(results["armed"].Text, "tool_limits set") {
		t.Fatal("actually dispatched setter mislabeled unexecuted")
	}
	if !results["spent"].Error || !strings.Contains(results["spent"].Text, "max_output_bytes 73") {
		t.Fatal("WAKEUP literal next rejected call did not consume armed limits", results["spent"])
	}
	if effects.Load() != 1 {
		t.Fatalf("WAKEUP undispatched effect ran: %d", effects.Load())
	}
	if r := wait5(t, reply); r.Err != nil {
		t.Fatal(r)
	}
	out, err := a.CallTool(context.Background(), "effect", json.RawMessage(`{}`))
	if err != nil || strings.Contains(out, "tool_limits consumed") {
		t.Fatal("rejected/undispatched setter leaked", out, err)
	}
}

func TestChapter5WakeupLaterProcessWait(t *testing.T) {
	dir := t.TempDir()
	release := filepath.Join(dir, "release")
	var requests atomic.Int32
	next := make(chan string, 1)
	var a ensemble.Agent
	a, _ = actorFixture(t, func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		var p []ensemble.Part
		switch requests.Add(1) {
		case 1:
			p = []ensemble.Part{prose("I will launch."), callPart("launch", "run_command", map[string]any{"command": fmt.Sprintf("while [ ! -f '%s' ]; do sleep .02; done; printf 'PROCESS_FINAL\\n'", release), "ai_callback_delay": 0})}
		case 2:
			id := 0
			for _, e := range a.Log().Events() {
				if e.Type == "tool_returned" && e.Tool.CallID == "launch" {
					id = e.Tool.Job.Handle
				}
			}
			p = []ensemble.Part{prose("I will observe."), callPart("wait", "wait_for_job", map[string]any{"handle": id, "ai_callback_delay": 3600})}
		default:
			next <- string(raw)
			p = []ensemble.Part{prose("handled process hint")}
		}
		json.NewEncoder(w).Encode(providerReply(ensemble.Anthropic, p))
	}, func(c *ensemble.Config) { c.EnableTools = true })
	t.Cleanup(func() { _ = os.WriteFile(release, nil, 0600) })
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	reply := a.Submit(ctx, "start")
	launch := wakeCalled(t, a, "launch")
	wakeCalled(t, a, "wait")
	if err := a.Post(ensemble.Hint{Text: "PROCESS_WAIT_HINT"}); err != nil {
		t.Fatal(err)
	}
	raw := wakeReceive(t, next, "later 3600-second wait continuation")
	result := wakeResults(t, raw, "launch", "wait")["wait"]
	if result.Error || !strings.Contains(result.Text, fmt.Sprintf("job %d running", launch.Job.Handle)) || !strings.Contains(raw, "PROCESS_WAIT_HINT") {
		t.Fatal("WAKEUP later wait lost actual message/state", raw)
	}
	if r := wait5(t, reply); r.Err != nil {
		t.Fatal(r)
	}
	j, _ := a.Parent().Jobs().Get(launch.Job.Handle)
	if j.Snapshot().Status != "running" {
		t.Fatal("process canceled by hint")
	}
	_ = os.WriteFile(release, nil, 0600)
	<-j.Done()
	out, err := a.CallTool(context.Background(), "wait_for_job", json.RawMessage(fmt.Sprintf(`{"handle":%d,"ai_callback_delay":3600}`, launch.Job.Handle)))
	if err != nil || !strings.Contains(out, "PROCESS_FINAL") {
		t.Fatal("eventual real process result", out, err)
	}
}

func TestChapter5WakeupTimerAndAvailableOutput(t *testing.T) {
	a, root := actorFixture(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(providerReply(ensemble.Anthropic, []ensemble.Part{prose("done")}))
	}, nil)
	j, err := root.Jobs().Create(a)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { j.Finish("", nil, nil) })
	if _, err = j.Wait(context.Background(), ensemble.WaitOptions{Delay: time.Millisecond, MaxBytes: 512}); err != nil {
		t.Fatal(err)
	}
	j.Write([]byte("AVAILABLE_NOW"))
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	opts := ensemble.WaitOptions{Delay: time.Hour, MaxBytes: 512}
	// JSON keeps the amended option test compilable against the predecessor.
	json.Unmarshal([]byte(`{"ReturnAvailable":true}`), &opts)
	out, err := j.Wait(ctx, opts)
	if err != nil || !strings.Contains(out, "AVAILABLE_NOW") {
		t.Fatalf("WAKEUP later available bytes waited for callback: %s %v", out, err)
	}
	// Empty initial snapshot: later output must notify/recheck without returning
	// unconditionally before a pattern, completion, wake or silence boundary.
	pattern := ensemble.WaitOptions{Delay: time.Second, MaxBytes: 512, Pattern: "MATCH_LATER"}
	json.Unmarshal([]byte(`{"ReturnAvailable":true}`), &pattern)
	matched := make(chan string, 1)
	go func() {
		s, e := j.Wait(context.Background(), pattern)
		if e != nil {
			s = e.Error()
		}
		matched <- s
	}()
	time.Sleep(20 * time.Millisecond)
	j.Write([]byte("UNMATCHED\n"))
	select {
	case out := <-matched:
		t.Fatalf("later nonmatching output ended observation: %s", out)
	case <-time.After(40 * time.Millisecond):
	}
	j.Write([]byte("MATCH_LATER\n"))
	if out := wakeReceive(t, matched, "later output pattern notification"); !strings.Contains(out, "MATCH_LATER") || strings.Contains(out, "AVAILABLE_NOW") {
		t.Fatal(out)
	}
	j.Finish("TERMINAL", nil, nil)
	out, err = j.Wait(context.Background(), ensemble.WaitOptions{Delay: time.Hour, MaxBytes: 512})
	if err != nil || strings.Contains(out, "MATCH_LATER") || !strings.Contains(out, "TERMINAL") {
		t.Fatal("cursor/terminal", out, err)
	}
}

func TestChapter5WakeupShortTimerCannotReleaseLaterWait(t *testing.T) {
	held, release := releaseOnce()
	started := make(chan ensemble.ToolContext, 1)
	var count atomic.Int32
	httpBodies := make(chan string, 4)
	var a ensemble.Agent
	a, _ = actorFixture(t, func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		httpBodies <- string(raw)
		var p []ensemble.Part
		switch count.Add(1) {
		case 1:
			p = []ensemble.Part{prose("I will join."), callPart("short-limit", "tool_limits", map[string]any{"ai_callback_delay": 0.12}), callPart("short", "join_like", map[string]any{})}
		case 2:
			id := 0
			for _, e := range a.Log().Events() {
				if e.Type == "tool_returned" && e.Tool.CallID == "short" {
					id = e.Tool.Job.Handle
				}
			}
			p = []ensemble.Part{prose("I will wait again."), callPart("later", "wait_for_job", map[string]any{"handle": id, "ai_callback_delay": 3600})}
		default:
			p = []ensemble.Part{prose("both messages handled")}
		}
		json.NewEncoder(w).Encode(providerReply(ensemble.Anthropic, p))
	}, nil)
	t.Cleanup(release)
	a.RegisterTool(wakeDefinition("join_like"), func(c ensemble.ToolContext, _ json.RawMessage) (string, error) {
		started <- c
		<-held
		return "ACTUAL_LATE", nil
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	reply := a.Submit(ctx, "begin")
	j := wakeReceive(t, started, "short timer handler").Job
	<-httpBodies
	_ = a.Post(ensemble.Hint{Text: "FIRST_SHORT_WAKE"})
	raw := wakeReceive(t, httpBodies, "short observation wake")
	if !strings.Contains(raw, "FIRST_SHORT_WAKE") {
		t.Fatal("missing first message")
	}
	wakeResults(t, raw, "short-limit", "short")
	wakeCalled(t, a, "later")
	select {
	case b := <-httpBodies:
		t.Fatalf("WAKEUP old timer released later observation: %s", b)
	case <-time.After(250 * time.Millisecond):
	}
	_ = a.Post(wakeAddressed(t))
	raw = wakeReceive(t, httpBodies, "addressed question wakes later wait")
	for _, want := range []string{"CHILD_QUESTION_CAN_I_WRITE", "child:writer", "parent:editor", "request-17", "turn-23"} {
		if !strings.Contains(raw, want) {
			t.Fatalf("actual later HTTP lost %s", want)
		}
	}
	if strings.Contains(raw, "FIRST_SHORT_WAKE") {
		t.Fatal("old message replayed")
	}
	rs := wakeResults(t, raw, "short", "later")
	if rs["later"].Error || !strings.Contains(rs["later"].Text, fmt.Sprintf("job %d running", j.Snapshot().Handle)) {
		t.Fatal("later wait not truthful", rs)
	}
	if r := wait5(t, reply); r.Err != nil {
		t.Fatal(r)
	}
	for _, id := range []string{"short", "later"} {
		n := 0
		for _, e := range a.Log().Events() {
			if e.Type == "tool_returned" && e.Tool.CallID == id {
				n++
			}
		}
		if n != 1 {
			t.Fatalf("WAKEUP duplicate report for %s: %d", id, n)
		}
	}
	release()
	<-j.Done()
}

func TestChapter5WakeupCompletionRaces(t *testing.T) {
	for _, order := range []string{"message-first", "completion-first", "concurrent"} {
		t.Run(order, func(t *testing.T) {
			held, release := releaseOnce()
			gate, open := releaseOnce()
			started := make(chan ensemble.ToolContext, 1)
			var count atomic.Int32
			payloads := make(chan string, 8)
			a, _ := actorFixture(t, func(w http.ResponseWriter, r *http.Request) {
				raw, _ := io.ReadAll(r.Body)
				payloads <- string(raw)
				p := []ensemble.Part{prose("finished")}
				if count.Add(1) == 1 {
					p = []ensemble.Part{prose("I will join."), callPart("limits", "tool_limits", map[string]any{"ai_callback_delay": 3600}), callPart("race", "join_like", map[string]any{})}
				} else {
					<-gate
				}
				json.NewEncoder(w).Encode(providerReply(ensemble.Anthropic, p))
			}, nil)
			t.Cleanup(func() { release(); open() })
			a.RegisterTool(wakeDefinition("join_like"), func(c ensemble.ToolContext, _ json.RawMessage) (string, error) {
				started <- c
				<-held
				return "RACE_REAL_FINAL", nil
			})
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			reply := a.Submit(ctx, "begin")
			j := wakeReceive(t, started, "race handler").Job
			post := func() { _ = a.Post(ensemble.Hint{Text: "RACE_MESSAGE"}) }
			switch order {
			case "message-first":
				post()
				wakeReceive(t, payloads, "first request")
				raw := wakeReceive(t, payloads, "wake before completion")
				if !strings.Contains(raw, "RACE_MESSAGE") {
					t.Fatal("message lost")
				}
				release()
			case "completion-first":
				release()
				<-j.Done()
				post()
			default:
				done := make(chan struct{})
				go func() { post(); close(done) }()
				release()
				<-done
			}
			ackctx, ackcancel := context.WithTimeout(context.Background(), time.Second)
			defer ackcancel()
			if _, err := a.Wait(ackctx, func(e ensemble.Observation) bool {
				p, ok := e.(ensemble.PartDelta)
				return ok && p.Chunk == "hint: RACE_MESSAGE"
			}); err != nil {
				t.Fatal(err)
			}
			open()
			if r := wait5(t, reply); r.Err != nil {
				t.Fatal(r)
			}
			n := 0
			for _, e := range a.Log().Events() {
				if e.Type == "tool_returned" && e.Tool.CallID == "race" {
					n++
					out := e.Tool.Parts[0].Text
					if e.Tool.IsError || (!strings.Contains(out, "RACE_REAL_FINAL") && !strings.Contains(out, fmt.Sprintf("job %d running", j.Snapshot().Handle))) {
						t.Fatal("untruthful completion race", out)
					}
				}
			}
			if n != 1 {
				t.Fatalf("WAKEUP race returned %d results", n)
			}
			if len(a.Log().Context().Ephemera) != 0 {
				t.Fatal("race stranded message")
			}
			if order != "message-first" {
				found := 0
				for len(payloads) > 0 {
					if strings.Contains(<-payloads, "RACE_MESSAGE") {
						found++
					}
				}
				if found != 1 {
					t.Fatalf("WAKEUP race message consumed %d times", found)
				}
			}
		})
	}
}

func TestChapter5WakeupFailedRetainedReport(t *testing.T) {
	held, release := releaseOnce()
	started := make(chan ensemble.ToolContext, 1)
	next := make(chan string, 1)
	var count atomic.Int32
	var a ensemble.Agent
	a, _ = actorFixture(t, func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		var p []ensemble.Part
		switch count.Add(1) {
		case 1:
			p = []ensemble.Part{prose("I will start the join."), callPart("limit", "tool_limits", map[string]any{"ai_callback_delay": 0}), callPart("start", "join_like", map[string]any{})}
		case 2:
			id := wakeCalled(t, a, "start").Job.Handle
			p = []ensemble.Part{prose("I will wait."), callPart("failed-wait", "wait_for_job", map[string]any{"handle": id, "ai_callback_delay": 3600})}
		default:
			next <- string(raw)
			p = []ensemble.Part{prose("handled failure")}
		}
		json.NewEncoder(w).Encode(providerReply(ensemble.Anthropic, p))
	}, nil)
	t.Cleanup(release)
	a.RegisterTool(wakeDefinition("join_like"), func(c ensemble.ToolContext, _ json.RawMessage) (string, error) {
		started <- c
		<-held
		return "REAL_FAILURE_OUTPUT", fmt.Errorf("ordinary failure")
	})
	reply := a.Submit(context.Background(), "begin")
	j := wakeReceive(t, started, "failure handler").Job
	wakeCalled(t, a, "failed-wait")
	// Completion wins this observation. Its useful report must not be replaced
	// by a bare error before the next actual model request.
	release()
	<-j.Done()
	raw := wakeReceive(t, next, "real failed wait outcome")
	r := wakeResults(t, raw, "failed-wait")["failed-wait"]
	if !r.Error || !strings.Contains(r.Text, "REAL_FAILURE_OUTPUT") || !strings.Contains(r.Text, fmt.Sprintf("job %d done", j.Snapshot().Handle)) {
		t.Fatalf("WAKEUP failed observation discarded retained report: %+v", r)
	}
	if r := wait5(t, reply); r.Err != nil {
		t.Fatal(r)
	}
}

func TestChapter5WakeupAddressedProcessInitialAndLater(t *testing.T) {
	dir := t.TempDir()
	ready, release := filepath.Join(dir, "ready"), filepath.Join(dir, "release")
	t.Cleanup(func() { os.WriteFile(release, []byte("release"), 0600) })
	next := make(chan string, 2)
	var n atomic.Int32
	var a ensemble.Agent
	var id int
	a, _ = actorFixture(t, func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		var p []ensemble.Part
		switch n.Add(1) {
		case 1:
			p = []ensemble.Part{prose("I will run a held process."), callPart("process", "run_command", map[string]any{"command": fmt.Sprintf("printf ready > '%s'; while [ ! -f '%s' ]; do sleep .02; done; printf ACTUAL_TERMINAL", ready, release), "ai_callback_delay": 3600})}
		case 2:
			next <- string(raw)
			id = wakeCalled(t, a, "process").Job.Handle
			p = []ensemble.Part{prose("I will wait again."), callPart("later", "wait_for_job", map[string]any{"handle": id, "ai_callback_delay": 3600})}
		default:
			next <- string(raw)
			p = []ensemble.Part{prose("answered both process questions")}
		}
		json.NewEncoder(w).Encode(providerReply(ensemble.Anthropic, p))
	}, func(c *ensemble.Config) { c.EnableTools = true })
	reply := a.Submit(context.Background(), "start")
	deadline := time.Now().Add(2 * time.Second)
	for {
		if _, e := os.Stat(ready); e == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("process readiness")
		}
		time.Sleep(time.Millisecond)
	}
	first := wakeAddressed(t)
	first.Text = "INITIAL_PROCESS_QUESTION"
	if e := a.Post(first); e != nil {
		t.Fatal(e)
	}
	raw := wakeReceive(t, next, "addressed initial process next HTTP")
	for _, token := range []string{first.Text, "child:writer", "parent:editor", "request-17", "turn-23"} {
		if !strings.Contains(raw, token) {
			t.Fatal("WAKEUP missing initial process question", token)
		}
	}
	result := wakeResults(t, raw, "process")["process"]
	if result.Error || !strings.Contains(result.Text, "running") {
		t.Fatal(result)
	}
	wakeCalled(t, a, "later")
	second := wakeAddressed(t)
	second.Text = "LATER_PROCESS_QUESTION"
	if e := a.Post(second); e != nil {
		t.Fatal(e)
	}
	raw = wakeReceive(t, next, "addressed later process next HTTP")
	for _, token := range []string{second.Text, "child:writer", "parent:editor", "request-17", "turn-23"} {
		if !strings.Contains(raw, token) {
			t.Fatal("WAKEUP missing later process question", token)
		}
	}
	if strings.Contains(raw, first.Text) {
		t.Fatal("consumed question replayed")
	}
	result = wakeResults(t, raw, "later")["later"]
	if result.Error || !strings.Contains(result.Text, fmt.Sprintf("job %d running", id)) {
		t.Fatal(result)
	}
	if r := wait5(t, reply); r.Err != nil {
		t.Fatal(r)
	}
	j, e := a.Parent().Jobs().Get(id)
	if e != nil {
		t.Fatal(e)
	}
	if j.Snapshot().Status != "running" {
		t.Fatal("question killed process")
	}
	os.WriteFile(release, []byte("release"), 0600)
	<-j.Done()
	out, e := a.CallTool(context.Background(), "wait_for_job", json.RawMessage(fmt.Sprintf(`{"handle":%d,"ai_callback_delay":3600}`, id)))
	if e != nil || !strings.Contains(out, "ACTUAL_TERMINAL") {
		t.Fatal(out, e)
	}
}
