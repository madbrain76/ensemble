from pathlib import Path
root=Path('/Users/bill/projects/ensemble-edition-3-student')
def edit(p,a,b):
 p=root/p;s=p.read_text();assert s.count(a)==1,(p,a[:80],s.count(a));p.write_text(s.replace(a,b))
edit('internal/common/jobs.go','"encoding/json"','"encoding/json"\n "errors"')
edit('internal/common/jobs.go','// JobRecord is', '// ErrObservationWoken releases a foreground observation without canceling work.\nvar ErrObservationWoken = errors.New("observation woken by input")\n\n// JobRecord is')
fields='''
 // Sender and Recipient preserve application-supplied addressing, not routing.
 Sender string `json:"sender,omitempty"`
 Recipient string `json:"recipient,omitempty"`
 // RequestID and TurnID preserve application-supplied correlation verbatim.
 RequestID string `json:"request_id,omitempty"`
 TurnID string `json:"turn_id,omitempty"`
'''
edit('internal/common/actor.go','type UserMessage struct {','type UserMessage struct {'+fields)
edit('internal/common/common.go','type MessageEvent struct {','type MessageEvent struct {'+fields)
edit('internal/history/history.go','func apply(c *common.Context, event common.Event) {','''// messageParts projects supplied addressing as explicit context without changing
// the sender's parts or confusing application names with provider roles.
func messageParts(m *common.MessageEvent) []common.Part {
 if m.Sender == "" && m.Recipient == "" && m.RequestID == "" && m.TurnID == "" { return m.Parts }
 data, _ := json.Marshal(struct {
  Sender string `json:"sender,omitempty"`
  Recipient string `json:"recipient,omitempty"`
  RequestID string `json:"request_id,omitempty"`
  TurnID string `json:"turn_id,omitempty"`
 }{m.Sender,m.Recipient,m.RequestID,m.TurnID})
 return append([]common.Part{{Type:"text",Text:"Message addressing: "+string(data)}},m.Parts...)
}
func apply(c *common.Context, event common.Event) {''')
p=root/'internal/history/history.go';s=p.read_text();s=s.replace('event.Message.Parts...', 'messageParts(event.Message)...').replace('Parts: event.Message.Parts}', 'Parts: messageParts(event.Message)}');p.write_text(s)
edit('internal/agent/actor.go','parts     []common.Part','parts     []common.Part\n message *common.MessageEvent')
edit('internal/agent/actor.go','dispatched     bool','dispatched     bool\n observation uint64\n observeCancel context.CancelCauseFunc\n wakeBatch bool')
edit('internal/agent/actor.go','type toolReport struct {\n\tid     uint64','type toolReport struct {\n id uint64\n observation uint64')
edit('internal/agent/actor.go','r := &request{ctx: ctx, parts: parts, ephemera:', 'r := &request{ctx: ctx, parts: parts, message: addressed(m,"human",parts), ephemera:')
edit('internal/agent/actor.go','func (a *Agent) next() {','''func addressed(m common.UserMessage, actor string, parts []common.Part) *common.MessageEvent {
 return &common.MessageEvent{Actor:actor,Parts:parts,Sender:m.Sender,Recipient:m.Recipient,RequestID:m.RequestID,TurnID:m.TurnID}
}

// wakeObservation runs only on the actor, after the actual input fact commits.
// The child context belongs to one report worker, never to the job or turn.
func (a *Agent) wakeObservation() {
 if t := a.active; t != nil && t.dispatched && t.observeCancel != nil {
  t.wakeBatch = true
  t.observeCancel(common.ErrObservationWoken)
 }
}

func (a *Agent) next() {''')
edit('internal/agent/actor.go','if err := a.append(common.Event{Type: "message_received", Message: &common.MessageEvent{Actor: "human", Parts: r.parts}}); err != nil {','''m := r.message
  if m == nil { m = &common.MessageEvent{Actor:"human",Parts:r.parts} }
  if err := a.append(common.Event{Type: "message_received", Message: m}); err != nil {''')
edit('internal/agent/actor.go','''a.workers.Add(1)
	go func(gen uint64, ctx context.Context) { defer a.workers.Done(); a.internal(toolReport{gen, work(ctx)}) }(t.id, t.ctx)''','''t.observation++
 ctx, cancel := context.WithCancelCause(t.ctx)
 t.observeCancel = cancel
 a.workers.Add(1)
 go func(gen, observation uint64) {
  defer a.workers.Done()
  defer cancel(nil)
  a.internal(toolReport{gen, observation, work(ctx)})
 }(t.id,t.observation)''')
edit('internal/agent/actor.go','t.calls = nil\n\tt.index = 0','t.calls = nil\n t.wakeBatch = false\n\tt.index = 0')
edit('internal/agent/actor.go','a.queue = append(a.queue, &request{ctx: context.Background(), parts: parts})','a.queue = append(a.queue, &request{ctx: context.Background(), parts: parts, message: addressed(v,"human",parts)})')
edit('internal/agent/actor.go','_ = a.append(common.Event{Type: "hint_received", Message: &common.MessageEvent{Actor: "system", Parts: parts}})','_ = a.append(common.Event{Type: "hint_received", Message: addressed(v,"system",parts)})\n a.wakeObservation()')
edit('internal/agent/actor.go','''_ = a.append(common.Event{Type: "hint_received", Message: &common.MessageEvent{Actor: "system", Parts: []common.Part{{Type: "text", Text: v.Text}}}})''','''_ = a.append(common.Event{Type: "hint_received", Message: &common.MessageEvent{Actor: "system", Parts: []common.Part{{Type: "text", Text: v.Text}}}})
 a.wakeObservation()''')
edit('internal/agent/actor.go','if a.active != nil && a.active.id == v.id {','if a.active != nil && a.active.id == v.id && a.active.dispatched && a.active.observation == v.observation {')
edit('internal/agent/actor.go','''t := a.active
					a.recordResult(v.result)''','''t := a.active
     t.observeCancel(nil)
     t.observeCancel = nil
					a.recordResult(v.result)''')
edit('internal/agent/actor.go','''} else {
						a.dispatchNext()
					}''','''} else {
      if t.wakeBatch {
       // Calls remain in the original assistant fact. Only their exact paired
       // errors are synthesized; no tool lookup, job creation or replay occurs.
       for t.index < len(t.calls) {
        c := t.calls[t.index]
        a.recordResult(a.notExecuted(c,"interrupted by incoming message; not executed"))
        t.index++
       }
       t.wakeBatch = false
      }
						a.dispatchNext()
					}''')
# share literal next-call consumption with synthetic errors, including policy/abort.
edit('internal/agent/dispatch.go','''// prepareCall owns all policy/one-shot mutations; work only observes an outcome.
func (a *Agent) prepareCall(call common.Part) (func(context.Context) common.ToolEvent, common.Job, error) {''','''// consumeLimits belongs to the actor. Even a call rejected before dispatch
// consumes the prior setting; an undispatched setter never installs a new one.
func (a *Agent) consumeLimits(call common.Part) (common.WaitOptions, string) {''')
edit('internal/agent/dispatch.go','\tvar callErr error',''' return opts, annotation
}
func (a *Agent) notExecuted(call common.Part, reason string) common.ToolEvent {
 _, annotation := a.consumeLimits(call)
 return failed(call,annotation+reason)
}
// prepareCall owns all policy/one-shot mutations; work only observes an outcome.
func (a *Agent) prepareCall(call common.Part) (func(context.Context) common.ToolEvent, common.Job, error) {
 opts, annotation := a.consumeLimits(call)
 var callErr error''')
edit('internal/agent/actor.go','r := failed(c, msg)','r := failed(c, msg)\n  if !t.dispatched { r = a.notExecuted(c,msg) }')
edit('internal/agent/actor.go','a.recordResult(failed(c, "tool batch refused:', 'a.recordResult(a.notExecuted(c, "tool batch refused:')
# same locked report path, force stable identity for a wake even if terminal wins.
edit('internal/jobs/jobs.go','ready := j.record.Status != "running" || (pattern != nil && pattern.Match(unseen))','''woken := context.Cause(ctx) == common.ErrObservationWoken
  ready := woken || j.record.Status != "running" || (j.reports > 0 && len(unseen) > 0) || (pattern != nil && pattern.Match(unseen))''')
edit('internal/jobs/jobs.go','j.report(unseen, opts)','j.report(unseen, opts, woken)')
edit('internal/jobs/jobs.go','j.report(data[j.cursor:], opts)','j.report(data[j.cursor:], opts, context.Cause(ctx) == common.ErrObservationWoken)')
edit('internal/jobs/jobs.go','case <-ctx.Done():\n\t\t\treturn "", ctx.Err()','''case <-ctx.Done():
   if context.Cause(ctx) == common.ErrObservationWoken { continue }
   return "", ctx.Err()''')
edit('internal/jobs/jobs.go','func (j *Job) report(unseen []byte, opts common.WaitOptions)', 'func (j *Job) report(unseen []byte, opts common.WaitOptions, identified bool)')
edit('internal/jobs/jobs.go','if first && j.record.Status == "done"', 'if first && !identified && j.record.Status == "done"')
