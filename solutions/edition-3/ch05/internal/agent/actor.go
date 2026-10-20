package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"ensemble/internal/common"
)

// The slice is the lossless FIFO; wake is only a coalescing signal.
type mailbox struct {
	mu     sync.Mutex
	queue  []any
	wake   chan struct{}
	closed bool
}

func (b *mailbox) post(v any, public bool) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if public && b.closed {
		return common.ErrStopped
	}
	b.queue = append(b.queue, v)
	select {
	case b.wake <- struct{}{}:
	default:
	}
	return nil
}
func (b *mailbox) drain() []any {
	b.mu.Lock()
	defer b.mu.Unlock()
	q := b.queue
	b.queue = nil
	return q
}

type request struct {
	live    bool
	turnID  uint64
	message *common.MessageData
	ctx     context.Context
	text    string
	reply   chan common.Result
	done    chan struct{}
}
type cancelRequest struct {
	request *request
	err     error
}
type response struct {
	turn   uint64
	events []common.Event
	err    error
}
type report struct {
	turn uint64
	call common.ToolData
}
type stop struct{}
type invoke struct {
	ctx     context.Context
	call    common.Part
	refusal string
	reply   chan error
}
type invoked struct {
	data  common.ToolData
	reply chan error
}

func (a *agent) InvokeTool(ctx context.Context, call common.Part, refusal string) error {
	reply := make(chan error, 1)
	if err := a.box.post(invoke{ctx, call, refusal, reply}, true); err != nil {
		return err
	}
	select {
	case err := <-reply:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

type stopped struct{ err error }
type turn struct {
	id            uint64
	request       *request
	followers     []*request
	ctx           context.Context
	cancel        context.CancelFunc
	httpCancel    context.CancelFunc
	observeCancel context.CancelFunc
	calls         common.PartList
	index         int
	operation     common.ToolOperation
	refusal       string
	skip          bool
	hinted        bool
	rounds, limit int
}

// Start is called once by the composition root, after all child links exist.
func (a *agent) Start() { go a.run() }
func (a *agent) Post(input common.Inbound) error {
	if input == nil {
		return errors.New("nil input")
	}
	return a.box.post(input, true)
}
func (a *agent) Submit(ctx context.Context, text string) <-chan common.Result {
	return a.submit(ctx, text, false)
}
func (a *agent) SubmitLive(ctx context.Context, text string) <-chan common.Result {
	return a.submit(ctx, text, true)
}
func (a *agent) submit(ctx context.Context, text string, live bool) <-chan common.Result {
	r := &request{live: live, ctx: ctx, text: text, reply: make(chan common.Result, 1), done: make(chan struct{})}
	if err := a.box.post(r, true); err != nil {
		r.reply <- common.Result{Err: err}
		close(r.reply)
		close(r.done)
		return r.reply
	}
	go func() {
		select {
		case <-ctx.Done():
			_ = a.box.post(cancelRequest{r, ctx.Err()}, false)
		case <-r.done:
		case <-a.done:
		}
	}()
	return r.reply
}
func (a *agent) Ask(ctx context.Context, text string) (string, error) {
	result := <-a.Submit(ctx, text)
	return result.Text, result.Err
}

type killedReport struct {
	report   common.JobReport
	recorded chan struct{}
}

func (a *agent) ReportJob(r common.JobReport) {
	if r.Job.Status != common.JobKilledStatus {
		_ = a.box.post(r, false)
		return
	}
	recorded := make(chan struct{})
	_ = a.box.post(killedReport{r, recorded}, false)
	// A direct Kill must not return before its terminal fact is recorded.
	// This acknowledges the actor, never the uncooperative handler.
	select {
	case <-recorded:
	case <-a.done:
	}
}
func (a *agent) Shutdown() error {
	a.box.mu.Lock()
	if !a.box.closed {
		a.box.closed = true
		a.box.queue = append(a.box.queue, stop{})
		select {
		case a.box.wake <- struct{}{}:
		default:
		}
	}
	a.box.mu.Unlock()
	<-a.done
	return a.shutdownErr
}

func (a *agent) recordJob(r common.JobReport) error {
	typ := common.JobReported
	if r.Job.Status == common.JobKilledStatus {
		typ = common.JobKilled
	}
	_, err := a.record(common.Event{Type: typ, TurnID: r.Job.TurnID, Job: &r.Job})
	return err
}
func (a *agent) run() {
	for {
		<-a.box.wake
		for _, message := range a.box.drain() {
			switch m := message.(type) {
			case *request:
				if a.stopping {
					a.reply(m, "", common.ErrStopped)
				} else if err := m.ctx.Err(); err != nil {
					a.reply(m, "", err)
				} else if m.live && a.active != nil {
					m.turnID = a.active.id
					a.active.followers = append(a.active.followers, m)
					a.hint(common.MessageData{Actor: common.Human, Parts: common.PartList{common.Text(m.text)}})
				} else {
					a.pending = append(a.pending, m)
				}
			case cancelRequest:
				if a.active != nil && a.active.request == m.request {
					a.interrupt(m.err)
				} else {
					if a.active != nil {
						for i, r := range a.active.followers {
							if r == m.request {
								a.active.followers = append(a.active.followers[:i], a.active.followers[i+1:]...)
								a.reply(r, "", m.err)
								break
							}
						}
					}
					for i, r := range a.pending {
						if r == m.request {
							a.pending = append(a.pending[:i], a.pending[i+1:]...)
							a.reply(r, "", m.err)
							break
						}
					}
				}
			case common.UserMessage:
				if !a.stopping {
					if a.active == nil {
						a.pending = append(a.pending, &request{ctx: context.Background(), text: m.Text})
					} else {
						a.hint(common.MessageData{Actor: common.Human, Parts: common.PartList{common.Text(m.Text)}})
					}
				}
			case common.Ephemeral:
				_, err := a.record(common.Event{Type: common.MessageReceived, Message: &common.MessageData{Actor: common.System, Parts: common.PartList{common.Text(m.Text)}}})
				if err != nil {
					a.finish("", err)
				}
			case common.Hint:
				if !a.stopping {
					a.hint(common.MessageData{Actor: common.Human, Parts: common.PartList{common.Text(m.Text)}})
				}
			case common.Question:
				if !a.stopping {
					a.hint(common.MessageData{Actor: common.Human, Parts: common.PartList{common.Text(m.Text)}, Sender: m.Sender, Recipient: m.Recipient, RequestID: m.RequestID, TurnID: m.TurnID})
				}
			case common.Interrupt:
				a.interrupt(common.ErrInterrupted)
			case invoke:
				if a.stopping {
					m.reply <- common.ErrStopped
					break
				}
				if err := m.ctx.Err(); err != nil {
					m.reply <- err
					break
				}
				if a.active != nil {
					m.reply <- errors.New("direct Invoke requires an idle conversation")
					break
				}
				op := a.tools.Prepare(m.ctx, m.call, m.refusal, 0)
				d := op.Called()
				if _, err := a.record(common.Event{Type: common.ToolCalled, Tool: &d}); err != nil {
					m.reply <- err
					break
				}
				a.workers.Add(1)
				go func() { defer a.workers.Done(); _ = a.box.post(invoked{op.Await(m.ctx), m.reply}, false) }()
			case invoked:
				m.reply <- a.returned(0, m.data)
			case response:
				a.receive(m)
			case report:
				a.report(m)
			case common.JobReport:
				if err := a.recordJob(m); err != nil {
					a.parent.Logger().Print(err)
				}
			case killedReport:
				if err := a.recordJob(m.report); err != nil {
					a.parent.Logger().Print(err)
				}
				close(m.recorded)
			case stop:
				a.stopping = true
				a.interrupt(common.ErrStopped)
				for _, r := range a.pending {
					a.reply(r, "", common.ErrStopped)
				}
				a.pending = nil
				go func() { err := a.jobs.Shutdown(); a.workers.Wait(); _ = a.box.post(stopped{err}, false) }()
			case stopped:
				a.shutdownErr = errors.Join(m.err, a.save())
				close(a.done)
				return
			}
			a.startNext()
		}
	}
}
func (a *agent) startNext() {
	for a.active == nil && !a.stopping && len(a.pending) > 0 {
		r := a.pending[0]
		a.pending = a.pending[1:]
		if err := r.ctx.Err(); err != nil {
			a.reply(r, "", err)
			continue
		}
		if strings.TrimSpace(r.text) == "" {
			a.reply(r, "", errors.New("user message is empty"))
			continue
		}
		a.nextTurn++
		ctx, cancel := context.WithCancel(r.ctx)
		a.active = &turn{id: a.nextTurn, request: r, ctx: ctx, cancel: cancel, limit: a.config.MaxToolRounds}
		r.turnID = a.active.id
		event := common.Event{Type: common.MessageReceived, TurnID: a.active.id, Message: &common.MessageData{Actor: common.Human, Parts: common.PartList{common.Text(r.text)}}}
		if r.message != nil {
			event.Type = common.HintReceived
			event.Message = r.message
		}
		if _, err := a.record(event); err != nil {
			a.finish("", err)
			continue
		}
		a.transition(common.InputPending)
		a.send()
	}
}
func (a *agent) hint(m common.MessageData) {
	if m.Sender != "" || m.Recipient != "" || m.RequestID != "" || m.TurnID != "" {
		metadata, _ := json.Marshal(map[string]string{"sender": m.Sender, "recipient": m.Recipient, "request_id": m.RequestID, "turn_id": m.TurnID})
		m.Parts = append(common.PartList{common.Text("Addressed question " + string(metadata) + "\n")}, m.Parts...)
	}
	if a.active == nil {
		var text strings.Builder
		for _, p := range m.Parts {
			if p.Type == "text" {
				text.WriteString(p.Text)
			}
		}
		a.pending = append(a.pending, &request{ctx: context.Background(), text: text.String(), message: &m})
		return
	}
	id := a.active.id
	if _, err := a.record(common.Event{Type: common.HintReceived, TurnID: id, Message: &m}); err != nil {
		if a.active != nil {
			a.finish("", err)
		}
		return
	}
	if a.active != nil {
		a.active.hinted = true
		if a.active.observeCancel != nil {
			a.active.skip = true
			a.active.observeCancel()
		}
	}
}
func (a *agent) send() {
	t := a.active
	body, err := a.engine.Render(a.history.Context())
	if err != nil {
		a.finish("", err)
		return
	}
	if _, err = a.record(common.Event{Type: common.RequestSent, TurnID: t.id, Request: &common.RequestData{To: a.config.Target()}}); err != nil {
		a.finish("", err)
		return
	}
	t.hinted = false
	a.transition(common.InFlight)
	ctx, cancel := context.WithCancel(t.ctx)
	t.httpCancel = cancel
	a.workers.Add(1)
	go func() {
		defer a.workers.Done()
		defer cancel()
		events, err := a.engine.ExchangeRequest(ctx, body)
		_ = a.box.post(response{t.id, events, err}, false)
	}()
}
func (a *agent) receive(r response) {
	stale := a.active == nil || a.active.id != r.turn
	if !stale {
		a.active.httpCancel = nil
	}
	idErr := a.identifyCalls(r.events)
	var parts common.PartList
	for _, e := range r.events {
		e.TurnID = r.turn
		e.Stale = stale
		seq, err := a.record(e)
		if err != nil {
			if !stale {
				a.finish("", err)
			}
			return
		}
		if e.Type == common.ResponseEnded && !stale {
			parts = e.Response.Parts
			for _, p := range e.Response.Parts {
				a.final(seq, p)
			}
		}
	}
	if stale {
		return
	}
	if r.err != nil {
		a.finish("", r.err)
		return
	}
	if idErr != nil {
		a.finish("", idErr)
		return
	}
	var text strings.Builder
	explained := false
	t := a.active
	t.calls = nil
	t.index = 0
	t.refusal = ""
	t.skip = false
	for _, p := range parts {
		if p.Type == "text" {
			text.WriteString(p.Text)
			if len(t.calls) == 0 && strings.TrimSpace(p.Text) != "" {
				explained = true
			}
		}
		if p.Type == "tool_call" {
			t.calls = append(t.calls, p)
		}
	}
	if len(t.calls) == 0 {
		if t.hinted {
			a.send()
		} else {
			a.finish(text.String(), nil)
		}
		return
	}
	a.transition(common.ToolsPending)
	if a.config.RequireToolExplanation && !explained {
		t.refusal = "No tools ran: provide nonblank ordinary assistant text explaining what will run and why before the first tool call, then retry."
	}
	// Input admitted during HTTP makes this entire batch undispatched.
	// Pair it without spending an allowance for a dispatched batch.
	if t.hinted {
		a.skipTail("interrupted, not executed: live input arrived before dispatch")
		a.send()
		return
	}
	if t.rounds >= t.limit {
		a.skipTail("not executed: tool round limit reached")
		a.finish("", common.ErrToolRoundLimit)
		return
	}
	t.rounds++
	a.dispatch()
}
func (a *agent) dispatch() {
	t := a.active
	if t.index == len(t.calls) {
		a.transition(common.InputPending)
		a.send()
		return
	}
	call := t.calls[t.index]
	op := a.tools.Prepare(t.ctx, call, t.refusal, t.id)
	d := op.Called()
	if _, err := a.record(common.Event{Type: common.ToolCalled, TurnID: t.id, Tool: &d}); err != nil {
		a.finish("", err)
		return
	}
	t.operation = op
	ctx, cancel := context.WithCancel(context.Background())
	t.observeCancel = cancel
	a.workers.Add(1)
	go func() {
		defer a.workers.Done()
		defer cancel()
		d := op.Await(ctx)
		_ = a.box.post(report{t.id, d}, false)
	}()
}
func (a *agent) returned(turn uint64, d common.ToolData) error {
	seq, err := a.record(common.Event{Type: common.ToolReturned, TurnID: turn, Tool: &d})
	if err == nil {
		a.final(seq, common.Part{Type: "tool_result", CallID: d.CallID, Name: d.Name, Parts: d.Parts, IsError: d.IsError})
	}
	return err
}
func (a *agent) report(r report) {
	t := a.active
	if t == nil || t.id != r.turn || t.operation == nil || t.calls[t.index].CallID != r.call.CallID {
		_, err := a.record(common.Event{Type: common.JobReported, TurnID: r.turn, Tool: &r.call})
		if err != nil {
			a.parent.Logger().Print(err)
		}
		return
	}
	t.observeCancel = nil
	t.operation = nil
	if err := a.returned(t.id, r.call); err != nil {
		a.finish("", err)
		return
	}
	t.index++
	if t.skip {
		a.skipTail("interrupted, not executed: live input interrupted this pending batch")
	}
	a.dispatch()
}
func (a *agent) skipTail(reason string) {
	t := a.active
	for t.index < len(t.calls) {
		op := a.tools.Prepare(t.ctx, t.calls[t.index], reason, t.id)
		d := op.Called()
		if _, err := a.record(common.Event{Type: common.ToolCalled, TurnID: t.id, Tool: &d}); err != nil {
			a.parent.Logger().Print(err)
		}
		// A refused operation has no worker, effects, timer, or job.
		if err := a.returned(t.id, op.Await(context.Background())); err != nil {
			a.parent.Logger().Print(err)
		}
		t.index++
	}
}
func (a *agent) interrupt(err error) {
	t := a.active
	if t == nil {
		return
	}
	if t.httpCancel != nil {
		t.httpCancel()
	}
	if t.observeCancel != nil {
		t.observeCancel()
	}
	if t.operation != nil {
		d := t.operation.Called()
		d.IsError = true
		d.Parts = common.PartList{common.Text("observation interrupted; dispatched work may still be running")}
		if recordErr := a.returned(t.id, d); recordErr != nil {
			err = errors.Join(err, recordErr)
		}
		t.operation = nil
		t.index++
	}
	a.skipTail("interrupted, not executed: turn ended before dispatch")
	if _, recordErr := a.record(common.Event{Type: common.TurnInterrupted, TurnID: t.id, Error: &common.ErrorData{Message: err.Error()}}); recordErr != nil {
		err = errors.Join(err, recordErr)
	}
	a.transition(common.Interrupted)
	a.finish("", err)
}
func (a *agent) finish(text string, err error) {
	t := a.active
	if t == nil {
		return
	}
	t.cancel()
	if err != nil {
		_, recordErr := a.record(common.Event{Type: common.ErrorOccurred, TurnID: t.id, Error: &common.ErrorData{Message: err.Error()}})
		err = errors.Join(err, recordErr)
	}
	_, recordErr := a.record(common.Event{Type: common.TurnEnded, TurnID: t.id, Message: &common.MessageData{Actor: common.Assistant, Parts: common.PartList{}}})
	err = errors.Join(err, recordErr)
	a.active = nil
	a.transition(common.Idle)
	a.reply(t.request, text, err)
	for _, r := range t.followers {
		a.reply(r, text, err)
	}
}
func (a *agent) reply(r *request, text string, err error) {
	if r.reply != nil {
		r.reply <- common.Result{TurnID: r.turnID, Text: text, Err: err}
		close(r.reply)
		close(r.done)
	}
}
func (a *agent) record(e common.Event) (common.Seq, error) {
	if err := a.history.Append(e); err != nil {
		return 0, err
	}
	events := a.history.Events()
	seq := events[len(events)-1].Seq
	return seq, a.save()
}
func (a *agent) save() error {
	path := a.config.LogPath
	if path == "" {
		return nil
	}
	f, err := os.CreateTemp(filepath.Dir(path), ".ensemble-log-*")
	if err != nil {
		return fmt.Errorf("saving log: %w", err)
	}
	name := f.Name()
	defer os.Remove(name)
	err = a.history.Dump(f)
	err = errors.Join(err, f.Close())
	if err == nil {
		err = os.Rename(name, path)
	}
	return err
}
func (a *agent) transition(to common.TurnState) {
	if a.state == to {
		return
	}
	from := a.state
	a.state = to
	a.publish(common.StateChanged{From: string(from), To: string(to)})
}
func (a *agent) final(seq common.Seq, p common.Part) {
	a.partID++
	if p.Type == "text" {
		a.publish(common.PartDelta{PartID: a.partID, Chunk: p.Text})
	}
	a.publish(common.PartFinal{Seq: uint64(seq), PartID: a.partID, Part: p})
}

// Each observer gets a bounded queue. Neither a slow callback nor a slow Wait
// predicate can block the actor; neither queue is a completion channel.
func (a *agent) subscribe() (uint64, chan common.Observation) {
	a.observersMu.Lock()
	defer a.observersMu.Unlock()
	a.nextObserver++
	id := a.nextObserver
	ch := make(chan common.Observation, 64)
	a.observers[id] = ch
	return id, ch
}
func (a *agent) unsubscribe(id uint64) {
	a.observersMu.Lock()
	delete(a.observers, id)
	a.observersMu.Unlock()
}
func (a *agent) publish(o common.Observation) {
	a.parent.Publish(a, o)
	a.observersMu.Lock()
	defer a.observersMu.Unlock()
	for _, ch := range a.observers {
		select {
		case ch <- o:
		default:
		}
	}
}
func (a *agent) Subscribe(o common.Observer) {
	if o == nil {
		return
	}
	id, ch := a.subscribe()
	go func() {
		defer a.unsubscribe(id)
		for {
			select {
			case v := <-ch:
				o.Observe(v)
			case <-a.done:
				return
			}
		}
	}()
}
func (a *agent) Wait(ctx context.Context, predicate func(common.Observation) bool) (common.Observation, error) {
	id, ch := a.subscribe()
	defer a.unsubscribe(id)
	for {
		select {
		case o := <-ch:
			if predicate == nil || predicate(o) {
				return o, nil
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-a.done:
			return nil, common.ErrStopped
		}
	}
}
