package agent

import (
	"fmt"

	"ensemble/internal/common"
)

// Settings acknowledgments follow their recorded content changes, even when idle.
type memorySettings struct {
	value common.MemorySettings
	reply chan error
}

// A compression worker owns only its HTTP call. The actor owns commits and
// drops a late result after interruption, leaving the old source information live.
type compressed struct {
	turn *turn
	task *common.Compression
	text string
	err  error
}

// Running startup and re-enable reread disk so human corrections replace logged
// bytes with the same identity. Offline history replay itself never scans files.
func (r *runtime) syncMemory() error {
	for b := common.BandSoul; b <= common.BandSession; b++ {
		disabled := r.owner.memory.Settings(b).Disabled
		wasDisabled := r.memoryDisabled[b]
		if !r.memoryReady {
			wasDisabled = r.owner.history.Context().DisabledBands[b]
		}
		if !r.memoryReady || wasDisabled != disabled {
			if disabled {
				if err := r.append(common.Event{Type: "band_depopulated", Depopulated: &common.BandDepopulatedData{Band: b, Reason: "disabled"}}); err != nil {
					return err
				}
			} else {
				files, err := r.owner.memory.Scan(b)
				if err != nil {
					return err
				}
				for i := range files {
					if err := r.append(common.Event{Type: "band_populated", Populated: &files[i]}); err != nil {
						return err
					}
				}
			}
		}
		r.memoryDisabled[b] = disabled
	}
	if !r.memoryReady && r.owner.history.Context().Compression != nil {
		r.owner.Ensemble().Logf("memory: abandoned compressor; source retained until a later checkpoint")
	}
	r.memoryReady = true
	return nil
}

// Only successful micro_handoff batches reach this hook. An under-budget
// conversation can still have a band waiting for its next checkpoint to fold.
func (r *runtime) checkpoint(conversation bool) {
	var task *common.Compression
	var err error
	if conversation {
		task, err = r.owner.memory.Conversation()
	}
	if task == nil && err == nil {
		task, err = r.owner.memory.Graduation()
	}
	if err != nil {
		r.compressionWarning(err)
		r.send()
		return
	}
	if task == nil {
		r.send()
		return
	}
	launch := &common.CompactorLaunchedData{Band: task.Destination, Thru: task.Thru}
	if err := r.append(common.Event{Type: "compactor_launched", Compactor: launch}); err != nil {
		r.finish("", err)
		return
	}
	t := r.active
	r.workers++
	go func() {
		text, err := r.owner.engine.Compress(t.request.ctx, task.Text, task.Ratio)
		r.owner.box.Post(compressed{t, task, text, err})
	}()
}
func (r *runtime) compressionWarning(err error) {
	r.owner.Ensemble().Logf("memory: compression abandoned: %v", err)
}
func (r *runtime) compressed(m compressed) {
	r.workers--
	if m.turn != r.active {
		return
	}
	t := m.task
	// A settings change during the HTTP call must not smuggle data into a band
	// that the user just disabled. No cut has taken effect at this point.
	if m.err == nil && (r.owner.memory.Settings(t.Destination).Disabled || r.owner.memory.Settings(t.Source).Disabled) {
		m.err = fmt.Errorf("source or destination disabled during compression")
	}
	if m.err == nil {
		m.err = r.owner.memory.Write(t, m.text)
	}
	if m.err != nil {
		r.compressionWarning(m.err)
		r.send()
		return
	}
	var event common.Event
	if t.Source == 0 {
		event = common.Event{Type: "redacted", Redact: &common.RedactData{From: t.From, To: t.To, Level: "redact_summary", Reason: "memory compression"}}
	} else {
		event = common.Event{Type: "band_depopulated", Depopulated: &common.BandDepopulatedData{Band: t.Source, Thru: &t.Thru, Reason: "graduation"}}
	}
	if err := r.append(event); err != nil {
		r.finish("", err)
		return
	}
	populated := &common.BandPopulatedData{Band: t.Destination, File: t.File, Thru: t.Thru, Text: m.text, Source: "compression"}
	if err := r.append(common.Event{Type: "band_populated", Populated: populated}); err != nil {
		r.finish("", err)
		return
	}
	if err := r.owner.memory.Archive(t); err != nil {
		r.compressionWarning(err)
	}
	r.checkpoint(false)
}

// Occupancy is the prior conversation request, never lifetime usage or the
// small compressor response. All three normalized input categories are disjoint.
func (r *runtime) forceCheckpoint() error {
	if r.owner.spec.Compaction {
		return nil
	}
	c := r.owner.history.Context()
	u := c.LastUsage
	window := r.owner.engine.Features().ContextWindow
	input := u.Input + u.CacheRead + u.CacheWrite
	if window == 0 || r.forcingChecked == c.SentThrough {
		return nil
	}
	r.forcingChecked = c.SentThrough
	if input*100 < window*90 {
		r.warned = false
		return nil
	}
	if !r.warned {
		r.warned = true
		if err := r.append(common.Event{Type: "message_received", Message: &common.MessageData{Actor: common.System, Parts: []common.Part{{Type: "text", Text: "Context is at least 90% full. Finish this micro-goal and call micro_handoff now to preserve what you learned."}}}}); err != nil {
			return err
		}
	}
	if input*100 >= window*95 {
		r.owner.forceHandoff.Store(true)
	}
	return nil
}
