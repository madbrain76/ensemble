package agent

import "ensemble/internal/common"

// A retrieval worker owns no conversation. Its generation token prevents an
// interrupted judgment from attaching stale snippets to the next user's turn.
type recalled struct {
	turn *turn
	text string
	err  error
}

// Recaller returns this Agent's optional passive retrieval child.
func (a *agent) Recaller() common.Recaller { return a.recaller }

// AttachRecaller verifies both sides of the ownership link at composition time.
func (a *agent) AttachRecaller(recaller common.Recaller) {
	if recaller == nil || recaller.Agent() != a {
		panic("Recall must belong to this Agent")
	}
	a.recaller = recaller
}
func (r *runtime) recall() {
	if r.owner.recaller == nil {
		r.send()
		return
	}
	t := r.active
	snapshot := r.owner.history.Context()
	r.workers++
	go func() {
		text, err := r.owner.recaller.Recall(t.request.ctx, t.request.text, snapshot)
		r.owner.box.Post(recalled{t, text, err})
	}()
}
func (r *runtime) recalled(m recalled) {
	r.workers--
	if m.turn != r.active {
		return
	}
	if err := m.turn.request.ctx.Err(); err != nil {
		r.interrupt(err)
		return
	}
	if m.err != nil {
		r.owner.Ensemble().Logf("recall unavailable: %v", m.err)
	}
	if m.text != "" {
		if err := r.append(common.Event{Type: "recall_attached", Recall: &common.RecallData{Parts: []common.Part{{Type: "text", Text: m.text}}}}); err != nil {
			r.finish("", err)
			return
		}
	}
	r.send()
}
