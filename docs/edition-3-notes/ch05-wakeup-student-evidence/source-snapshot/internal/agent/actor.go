package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"ensemble/internal/common"
)

type request struct {
	ctx       context.Context
	parts     []common.Part
	message   *common.MessageEvent
	ephemera  []common.Part
	reply     chan common.Result
	direct    *common.Part
	toolReply chan common.Result
	barrier   bool
}
type turn struct {
	id             uint64
	req            *request
	ctx            context.Context
	cancel         context.CancelFunc
	calls          []common.Part
	index, batches int
	refusals       int
	job            common.Job
	dispatched     bool
	observation    uint64
	observeCancel  context.CancelCauseFunc
	wakeBatch      bool
	text           string
}
type modelReport struct {
	id       uint64
	response common.ResponseEvent
	err      error
}
type toolReport struct {
	id          uint64
	observation uint64
	result      common.ToolEvent
}
type jobComplete struct{ report common.JobRecord }
type cancelRequest struct {
	req *request
	err error
}
type stopInput struct{}
type finishInput struct{}
type observeInput struct{ observer common.Observer }

// Start idempotently starts the sole conversation owner.
func (a *Agent) Start() { a.once.Do(func() { go a.loop() }) }

// enqueue protects only admission and the mailbox, never network or tool work.
// The wake channel coalesces notifications; messages themselves are not dropped.
func (a *Agent) enqueue(m any, external bool) error {
	a.Start()
	a.mailMu.Lock()
	defer a.mailMu.Unlock()
	if external && !a.accepting {
		return common.ErrStopped
	}
	a.mail = append(a.mail, m)
	select {
	case a.wake <- struct{}{}:
	default:
	}
	return nil
}
func (a *Agent) internal(m any) { _ = a.enqueue(m, false) }

// Post copies input before admission; it does not promise completed work.
func (a *Agent) Post(m common.Input) error {
	switch v := m.(type) {
	case common.UserMessage:
		if len(v.Parts) == 0 && strings.TrimSpace(v.Text) == "" {
			return fmt.Errorf("user message is empty")
		}
		v.Parts = cloneParts(v.Parts)
		v.Ephemera = cloneParts(v.Ephemera)
		return a.enqueue(v, true)
	case common.Hint:
		if strings.TrimSpace(v.Text) == "" {
			return fmt.Errorf("hint is empty")
		}
		return a.enqueue(v, true)
	case common.Interrupt:
		return a.enqueue(v, true)
	default:
		return fmt.Errorf("unknown input %T", m)
	}
}

// Submit admits a separate request even while another turn is active.
func (a *Agent) Submit(ctx context.Context, text string) <-chan common.Result {
	return a.SubmitMessage(ctx, common.UserMessage{Text: text})
}

// SubmitMessage snapshots caller-owned parts and allocates the reliable reply.
func (a *Agent) SubmitMessage(ctx context.Context, m common.UserMessage) <-chan common.Result {
	if ctx == nil {
		ctx = context.Background()
	}
	parts := cloneParts(m.Parts)
	if len(parts) == 0 && strings.TrimSpace(m.Text) != "" {
		parts = []common.Part{{Type: "text", Text: m.Text}}
	}
	r := &request{ctx: ctx, parts: parts, message: addressed(m, "human", parts), ephemera: cloneParts(m.Ephemera), reply: make(chan common.Result, 1)}
	if len(parts) == 0 {
		r.reply <- common.Result{Err: fmt.Errorf("user message is empty")}
		close(r.reply)
	} else if err := a.enqueue(r, true); err != nil {
		r.reply <- common.Result{Err: err}
		close(r.reply)
	}
	return r.reply
}

// Ask is the blocking adapter over Submit, not a second turn loop.
func (a *Agent) Ask(ctx context.Context, text string) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	reply := a.Submit(ctx, text)
	select {
	case r := <-reply:
		return r.Text, r.Err
	default:
	}
	var err error
	select {
	case r := <-reply:
		return r.Text, r.Err
	case <-ctx.Done():
		err = ctx.Err()
	case <-a.done:
		err = common.ErrStopped
	}
	// Cancellation/shutdown does not discard an already-buffered durable result.
	// Otherwise the actor still owns completion; its private buffer never needs us.
	select {
	case r := <-reply:
		return r.Text, r.Err
	default:
		return "", err
	}
}

// CallTool uses the same owner and retained-job path as model tool calls.
func (a *Agent) CallTool(ctx context.Context, name string, args json.RawMessage) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	r := &request{ctx: ctx, direct: &common.Part{Name: name, Args: append(json.RawMessage(nil), args...)}, toolReply: make(chan common.Result, 1)}
	if err := a.enqueue(r, true); err != nil {
		return "", err
	}
	result := <-r.toolReply
	return result.Text, result.Err
}

// Observe installs a callback on the owner loop; the callback must not block.
func (a *Agent) Observe(o common.Observer) {
	if o != nil {
		_ = a.enqueue(observeInput{o}, true)
	}
}

// Wait consumes lossy progress, optionally filtering without growing a call stack.
func (a *Agent) Wait(ctx context.Context, pred ...func(common.HubEvent) bool) (common.HubEvent, error) {
	a.Start()
	for {
		select {
		case e, ok := <-a.observations:
			if !ok {
				return nil, common.ErrStopped
			}
			if len(pred) > 0 && pred[0] != nil && !pred[0](e) {
				continue
			}
			return e, nil
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

// State reads only the synchronized reduced history snapshot.
func (a *Agent) State() common.TurnState {
	a.stateMu.RLock()
	defer a.stateMu.RUnlock()
	return a.state
}
func cloneHub(e common.HubEvent) common.HubEvent {
	switch v := e.(type) {
	case common.PartFinal:
		v.Part = cloneParts([]common.Part{v.Part})[0]
		return v
	case common.PartDelta:
		v.Part = cloneParts([]common.Part{v.Part})[0]
		return v
	}
	return e
}

// emit publishes copied payloads after reduction. A full observation channel
// drops progress; it never delays reliable replies or changes the event log.
func (a *Agent) emit(e common.HubEvent) {
	select {
	case a.observations <- cloneHub(e):
	default:
	}
	for _, o := range a.observers {
		o(cloneHub(e))
	}
}
func (a *Agent) changed() {
	s := a.history.Context().Turn
	a.stateMu.Lock()
	old := a.state
	a.state = s
	a.stateMu.Unlock()
	if old != s {
		a.emit(common.StateChanged{State: s, From: old, To: s})
	}
}
func (a *Agent) append(e common.Event) error {
	err := a.history.Append(e)
	if err == nil {
		a.changed()
	}
	return err
}

// StopAdmission establishes the boundary before root-owned jobs are joined.
func (a *Agent) StopAdmission() {
	a.Start()
	a.mailMu.Lock()
	if a.accepting {
		a.accepting = false
		a.mail = append(a.mail, stopInput{})
		select {
		case a.wake <- struct{}{}:
		default:
		}
	}
	a.mailMu.Unlock()
	<-a.stopAck
}

// Shutdown waits for workers, commits their reports, then persists the final log.
func (a *Agent) Shutdown() error {
	a.StopAdmission()
	// Root calls StopAdmission on every Agent before the shared job ledger is closed.
	_ = a.parent.Jobs().Shutdown()
	a.finishOnce.Do(func() { a.workers.Wait(); a.internal(finishInput{}) })
	<-a.done
	return a.shutdownErr
}
func (a *Agent) save() error {
	path := a.config.HistoryPath
	if path == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	f, err := os.CreateTemp(filepath.Dir(path), ".history-*")
	if err != nil {
		return err
	}
	name := f.Name()
	defer os.Remove(name)
	err = a.history.Dump(f)
	closeErr := f.Close()
	if err != nil {
		return err
	}
	if closeErr != nil {
		return closeErr
	}
	return os.Rename(name, path)
}
func (a *Agent) complete(r *request, text string, err error) {
	if r.toolReply != nil {
		r.toolReply <- common.Result{Text: text, Err: err}
		close(r.toolReply)
	}
	if r.reply != nil {
		r.reply <- common.Result{Text: text, Err: err}
		close(r.reply)
	}
}

// end saves the reduced terminal fact before replying. A failed save is a request
// failure, not a successful completion followed by an advisory warning.
func (a *Agent) end(err error) {
	t := a.active
	if t == nil {
		return
	}
	t.cancel()
	outcome := "completed"
	if err != nil {
		outcome = err.Error()
	}
	_ = a.append(common.Event{Type: "turn_ended", Error: &common.ErrorEvent{Message: outcome}})
	if saveErr := a.save(); saveErr != nil {
		err = fmt.Errorf("save history: %w", saveErr)
	}
	a.complete(t.req, t.text, err)
	a.active = nil
}

// abort pairs outstanding calls with synthetic outcomes but never kills their
// retained jobs. Generation checks keep late workers from resuming this turn.
func (a *Agent) abort(err error) {
	t := a.active
	if t == nil {
		return
	}
	if t.req.direct != nil && t.index == 0 {
		r := failed(*t.req.direct, err.Error())
		if t.job != nil {
			snap := t.job.Snapshot()
			r.Job = &snap
		}
		a.recordResult(r)
		t.text = toolText(r)
		t.index++
	}
	for t.req.direct == nil && t.index < len(t.calls) {
		c := t.calls[t.index]
		msg := "not executed: " + err.Error()
		if t.dispatched {
			msg = "observation detached: " + err.Error() + "; retained job may continue"
		}
		r := failed(c, msg)
		if !t.dispatched {
			r = a.notExecuted(c, msg)
		}
		if t.job != nil {
			snap := t.job.Snapshot()
			r.Job = &snap
		}
		a.recordResult(r)
		t.index++
		t.dispatched = false
		t.job = nil
	}
	a.end(err)
}
func addressed(m common.UserMessage, actor string, parts []common.Part) *common.MessageEvent {
	return &common.MessageEvent{Actor: actor, Parts: parts, Sender: m.Sender, Recipient: m.Recipient, RequestID: m.RequestID, TurnID: m.TurnID}
}

// wakeObservation runs only on the actor, after the actual input fact commits.
// The child context belongs to one report worker, never to the job or turn.
func (a *Agent) wakeObservation() {
	if t := a.active; t != nil && t.dispatched && t.observeCancel != nil {
		t.wakeBatch = true
		t.observeCancel(common.ErrObservationWoken)
	}
}

func (a *Agent) next() {
	if a.active != nil || a.stopping {
		return
	}
	for len(a.queue) > 0 {
		r := a.queue[0]
		a.queue = a.queue[1:]
		if err := r.ctx.Err(); err != nil {

			a.complete(r, "", err)
			continue
		}
		if r.barrier {
			a.complete(r, "", nil)
			continue
		}
		a.generation++
		ctx, cancel := context.WithCancel(r.ctx)
		a.active = &turn{id: a.generation, req: r, ctx: ctx, cancel: cancel}
		a.workers.Add(1)
		go func() { defer a.workers.Done(); <-ctx.Done(); a.internal(cancelRequest{r, ctx.Err()}) }()
		if r.direct != nil {
			r.direct.CallID = fmt.Sprintf("local-%d", a.generation)
			a.active.calls = []common.Part{*r.direct}
			a.dispatchNext()
			return
		}
		if len(r.ephemera) > 0 {
			_ = a.append(common.Event{Type: "hint_received", Message: &common.MessageEvent{Actor: "system", Parts: r.ephemera}})
		}
		m := r.message
		if m == nil {
			m = &common.MessageEvent{Actor: "human", Parts: r.parts}
		}
		if err := a.append(common.Event{Type: "message_received", Message: m}); err != nil {
			a.end(err)
			continue
		}
		a.sendModel()
		return
	}
}

// sendModel snapshots context before workers can observe it. Rendering (including
// local media reads), HTTP and parsing stay outside the owner loop.
func (a *Agent) sendModel() {
	t := a.active
	c := a.history.Context()
	_ = a.append(common.Event{Type: "request_sent", Request: &common.RequestEvent{To: a.engine.Target()}})
	a.workers.Add(1)
	go func(id uint64, ctx context.Context) {
		defer a.workers.Done()
		response, err := a.engine.Send(ctx, c)
		a.internal(modelReport{id, response, err})
	}(t.id, t.ctx)
}
func (a *Agent) recordResult(r common.ToolEvent) {
	_ = a.append(common.Event{Type: "tool_returned", Tool: &r})
	a.emit(common.PartFinal{Part: common.Part{Type: "tool_result", CallID: r.CallID, Name: r.Name, Parts: r.Parts, IsError: r.IsError}, Seq: uint64(len(a.history.Events())), PartID: 1})
}
func (a *Agent) dispatchNext() {
	t := a.active
	if t.index == len(t.calls) {
		a.sendModel()
		return
	}
	c := t.calls[t.index]
	work, j, err := a.prepareCall(c)
	if err != nil {
		a.abort(err)
		return
	}
	t.dispatched = true
	t.job = j
	t.observation++
	ctx, cancel := context.WithCancelCause(t.ctx)
	t.observeCancel = cancel
	a.workers.Add(1)
	go func(gen, observation uint64) {
		defer a.workers.Done()
		defer cancel(nil)
		a.internal(toolReport{gen, observation, work(ctx)})
	}(t.id, t.observation)
}
func (a *Agent) onModel(m modelReport) {
	t := a.active
	if t == nil || t.id != m.id {
		return
	}
	if m.err != nil {
		_ = a.append(common.Event{Type: "error_occurred", Error: &common.ErrorEvent{Message: m.err.Error()}})
		a.end(m.err)
		return
	}
	for i := range m.response.Parts {
		p := &m.response.Parts[i]
		if p.Type == "tool_call" && p.CallID == "" {
			p.CallID = fmt.Sprintf("call-%d-%d-%d", t.id, t.batches, i)
		}
	}
	_ = a.append(common.Event{Type: "response_ended", Response: &m.response})
	visible, explained := false, true
	var text strings.Builder
	t.calls = nil
	t.wakeBatch = false
	t.index = 0
	t.job = nil
	t.dispatched = false
	for i, p := range m.response.Parts {
		if p.Type == "text" {
			a.emit(common.PartDelta{Part: p, PartID: i + 1, Chunk: p.Text})
		}
		a.emit(common.PartFinal{Part: p, Seq: uint64(len(a.history.Events())), PartID: i + 1})
		if p.Type == "text" && p.Type != "redacted" {
			text.WriteString(p.Text)
			if strings.TrimSpace(p.Text) != "" {
				visible = true
			}
		}
		if p.Type == "tool_call" && p.Type != "redacted" {
			if len(t.calls) == 0 {
				explained = visible
			}
			t.calls = append(t.calls, p)
		}
	}
	t.text = text.String()
	if len(t.calls) == 0 {
		if strings.TrimSpace(t.text) == "" {
			a.end(fmt.Errorf("%s response has no text parts", a.config.Vendor))
			return
		}
		// An admitted hint cannot change the request that produced this response.
		// Consume it in this turn's next immutable snapshot before claiming success.
		if len(a.history.Context().Ephemera) > 0 {
			a.sendModel()
			return
		}
		a.end(nil)
		return
	}
	limit := a.config.ToolBatchLimit
	if limit == 0 {
		limit = common.DefaultToolBatchLimit
	}
	if t.batches >= limit {
		a.abort(common.ErrToolBatchLimit)
		return
	}
	if a.config.RequireToolExplanation && !explained {
		t.refusals++
		if t.refusals > common.DefaultToolBatchLimit {
			a.abort(common.ErrToolBatchLimit)
			return
		}
		for _, c := range t.calls {
			e := common.ToolEvent{CallID: c.CallID, Name: c.Name, Args: c.Args}
			_ = a.append(common.Event{Type: "tool_called", Tool: &e})
			a.recordResult(a.notExecuted(c, "tool batch refused: provide visible explanation text BEFORE the first tool call, then retry the calls; no tool in this batch ran"))
		}
		t.calls = nil
		a.sendModel()
		return
	}
	t.batches++
	a.dispatchNext()
}
func (a *Agent) loop() {
	defer close(a.done)
	defer close(a.observations)
	for {
		<-a.wake
		for {
			a.mailMu.Lock()
			if len(a.mail) == 0 {
				a.mailMu.Unlock()
				break
			}
			m := a.mail[0]
			a.mail[0] = nil
			a.mail = a.mail[1:]
			a.mailMu.Unlock()
			switch v := m.(type) {
			case *request:
				if a.stopping {

					a.complete(v, "", common.ErrStopped)
				} else {
					a.queue = append(a.queue, v)
				}
			case common.UserMessage:
				parts := v.Parts
				if len(parts) == 0 {
					parts = []common.Part{{Type: "text", Text: v.Text}}
				}
				if a.active == nil {
					a.queue = append(a.queue, &request{ctx: context.Background(), parts: parts, message: addressed(v, "human", parts)})
				} else {
					_ = a.append(common.Event{Type: "hint_received", Message: addressed(v, "system", parts)})
					a.wakeObservation()
					for _, p := range parts {
						if p.Type == "text" {
							a.emit(common.PartDelta{Chunk: "hint: " + p.Text, Part: common.Part{Type: "text", Text: "hint: " + p.Text}})
						}
					}
				}
			case common.Hint:
				if a.active == nil {
					a.queue = append(a.queue, &request{ctx: context.Background(), parts: []common.Part{{Type: "text", Text: v.Text}}})
				} else {
					_ = a.append(common.Event{Type: "hint_received", Message: &common.MessageEvent{Actor: "system", Parts: []common.Part{{Type: "text", Text: v.Text}}}})
					a.wakeObservation()
					a.emit(common.PartDelta{Chunk: "hint: " + v.Text, Part: common.Part{Type: "text", Text: "hint: " + v.Text}})
				}
			case common.Interrupt:
				a.abort(common.ErrInterrupted)
			case observeInput:
				a.observers = append(a.observers, v.observer)
			case cancelRequest:
				if a.active != nil && a.active.req == v.req {
					a.abort(v.err)
				}
			case modelReport:
				a.onModel(v)
			case toolReport:
				if a.active != nil && a.active.id == v.id && a.active.dispatched && a.active.observation == v.observation {
					t := a.active
					t.observeCancel(nil)
					t.observeCancel = nil
					a.recordResult(v.result)
					t.index++
					t.dispatched = false
					t.job = nil
					if t.req.direct != nil {
						t.text = toolText(v.result)
						a.end(nil)
					} else {
						if t.wakeBatch {
							// Calls remain in the original assistant fact. Only their exact paired
							// errors are synthesized; no tool lookup, job creation or replay occurs.
							for t.index < len(t.calls) {
								c := t.calls[t.index]
								a.recordResult(a.notExecuted(c, "interrupted by incoming message; not executed"))
								t.index++
							}
							t.wakeBatch = false
						}
						a.dispatchNext()
					}
				}
			case jobKilled:
				r := v.record
				_ = a.append(common.Event{Type: "job_killed", Job: &r})
			case jobComplete:
				r := v.report
				_ = a.append(common.Event{Type: "job_report", Job: &r})
			case stopInput:
				a.stopping = true
				a.abort(common.ErrStopped)
				for _, r := range a.queue {

					a.complete(r, "", common.ErrStopped)
				}
				a.queue = nil
				close(a.stopAck)
			case finishInput:
				a.shutdownErr = a.save()
				return
			}
			a.next()
		}
	}
}

func failed(c common.Part, message string) common.ToolEvent {
	return common.ToolEvent{CallID: c.CallID, Name: c.Name, Args: c.Args, IsError: true, Parts: []common.Part{{Type: "text", Text: message}}}
}
func toolText(e common.ToolEvent) string {
	var b strings.Builder
	for _, p := range e.Parts {
		b.WriteString(p.Text)
	}
	return b.String()
}

// RecordJob routes an explicit job-kill fact back to the history owner.
func (a *Agent) RecordJob(r common.JobRecord) { a.internal(jobKilled{r}) }

type jobKilled struct{ record common.JobRecord }

// Drain reliably waits for work admitted before the barrier, without a new prompt.
func (a *Agent) Drain(ctx context.Context) error {
	r := &request{ctx: ctx, barrier: true, reply: make(chan common.Result, 1)}
	if err := a.enqueue(r, true); err != nil {
		return err
	}
	select {
	case result := <-r.reply:
		return result.Err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// cloneParts detaches borrowed data before this owner mutates or publishes it.
func cloneParts(parts []common.Part) []common.Part {
	if parts == nil {
		return nil
	}
	out := make([]common.Part, len(parts))
	for i, p := range parts {
		p.Args = append([]byte(nil), p.Args...)
		p.Data = append([]byte(nil), p.Data...)
		p.Opaque = append([]byte(nil), p.Opaque...)
		p.Parts = cloneParts(p.Parts)
		if p.Ref != nil {
			r := *p.Ref
			p.Ref = &r
		}
		out[i] = p
	}
	return out
}
