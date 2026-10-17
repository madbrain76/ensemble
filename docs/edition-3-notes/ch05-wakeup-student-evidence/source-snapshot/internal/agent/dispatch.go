package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"ensemble/internal/common"
)

func defaults() common.WaitOptions {
	return common.WaitOptions{Delay: 3 * time.Second, MaxBytes: 16384}
}
func limits(raw json.RawMessage, base common.WaitOptions) (common.WaitOptions, error) {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(raw, &fields); err != nil {
		return base, err
	}
	if fields == nil {
		return base, fmt.Errorf("arguments must be an object")
	}
	for _, key := range []string{"ai_callback_delay", "ai_callback_pattern", "max_output_bytes"} {
		if value, ok := fields[key]; ok && bytes.Equal(bytes.TrimSpace(value), []byte("null")) {
			return base, fmt.Errorf("%s must not be null", key)
		}
	}
	var args struct {
		Delay   *float64 `json:"ai_callback_delay"`
		Pattern *string  `json:"ai_callback_pattern"`
		Cap     *int     `json:"max_output_bytes"`
	}
	if err := json.Unmarshal(raw, &args); err != nil {
		return base, err
	}
	if args.Delay != nil {
		if *args.Delay < 0 || *args.Delay > float64(1<<63-1)/float64(time.Second) {
			return base, fmt.Errorf("ai_callback_delay must be nonnegative and finite")
		}
		base.Delay = time.Duration(*args.Delay * float64(time.Second))
	}
	if args.Pattern != nil {
		base.Pattern = *args.Pattern
	}
	if args.Cap != nil {
		base.MaxBytes = *args.Cap
	}
	if base.MaxBytes <= 0 {
		return base, fmt.Errorf("max_output_bytes must be positive")
	}
	if _, err := regexp.Compile(base.Pattern); err != nil {
		return base, err
	}
	return base, nil
}
func isControl(name string) bool {
	return name == "wait_for_job" || name == "send_input" || name == "kill_job" || name == "tool_limits"
}

// consumeLimits belongs to the actor. Even a call rejected before dispatch
// consumes the prior setting; an undispatched setter never installs a new one.
func (a *Agent) consumeLimits(call common.Part) (common.WaitOptions, string) {
	opts := defaults()
	consumed := a.pending != nil
	if consumed {
		opts = *a.pending
		a.pending = nil
	}
	annotation := ""
	if consumed {
		pattern := opts.Pattern
		if pattern == "" {
			pattern = "none"
		}
		annotation = fmt.Sprintf("[tool_limits consumed by this %s call: ai_callback_delay %s, ai_callback_pattern %s, max_output_bytes %d]\n", call.Name, opts.Delay, pattern, opts.MaxBytes)
	}
	return opts, annotation
}
func (a *Agent) notExecuted(call common.Part, reason string) common.ToolEvent {
	_, annotation := a.consumeLimits(call)
	return failed(call, annotation+reason)
}

// prepareCall owns all policy/one-shot mutations; work only observes an outcome.
func (a *Agent) prepareCall(call common.Part) (func(context.Context) common.ToolEvent, common.Job, error) {
	opts, annotation := a.consumeLimits(call)
	var callErr error
	if call.Name == "run_command" || isControl(call.Name) {
		opts, callErr = limits(call.Args, opts)
	}
	event := common.ToolEvent{CallID: call.CallID, Name: call.Name, Args: call.Args}
	var j common.Job
	if !isControl(call.Name) {
		var err error
		j, err = a.parent.Jobs().Create(a)
		if err != nil {
			return nil, nil, err
		}
		r := j.Snapshot()
		event.Job = &r
	}
	if err := a.append(common.Event{Type: "tool_called", Tool: &event}); err != nil {
		return nil, j, err
	}
	finish := func(out string, err error) common.ToolEvent {
		e := event
		e.Parts = []common.Part{{Type: "text", Text: annotation + out}}
		e.IsError = err != nil
		return e
	}
	if call.Name == "tool_limits" {
		out := ""
		if callErr == nil {
			out, callErr = a.control(context.Background(), call.Name, call.Args, opts)
		}
		if callErr != nil {
			out = "error: " + callErr.Error()
		}
		e := finish(out, callErr)
		return func(context.Context) common.ToolEvent { return e }, nil, nil
	}
	if isControl(call.Name) {
		return func(ctx context.Context) common.ToolEvent {
			out := ""
			err := callErr
			if err == nil {
				out, err = a.control(ctx, call.Name, call.Args, opts)
			}
			if err != nil && out == "" {
				out = "error: " + err.Error()
			}
			return finish(out, err)
		}, nil, nil
	}
	var tool common.Tool
	for _, candidate := range a.Tools() {
		if candidate.Definition().Name == call.Name {
			tool = candidate
			break
		}
	}
	if callErr != nil {
		j.Finish("", callErr, nil)
		opts = defaults()
	} else if tool == nil {
		j.Finish("", fmt.Errorf("unknown tool %q", call.Name), nil)
	} else {
		j.Start(context.Background(), tool, call.Args)
	}
	a.workers.Add(1)
	go func() { defer a.workers.Done(); <-j.Done(); a.internal(jobComplete{j.Snapshot()}) }()
	return func(ctx context.Context) common.ToolEvent {
		out, err := j.Wait(ctx, opts)
		if err != nil && out == "" {
			out = "error: " + err.Error()
		}
		e := finish(out, err)
		r := j.Snapshot()
		e.Job = &r
		return e
	}, j, nil
}

func (a *Agent) control(ctx context.Context, name string, raw json.RawMessage, opts common.WaitOptions) (string, error) {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(raw, &fields); err != nil {
		return "", err
	}
	if fields == nil {
		return "", fmt.Errorf("arguments must be an object")
	}
	for key, value := range fields {
		allowed := key == "handle" && name != "tool_limits" || key == "input" && name == "send_input" || name != "kill_job" && (key == "ai_callback_delay" || key == "ai_callback_pattern" || key == "max_output_bytes")
		if !allowed {
			return "", fmt.Errorf("unknown argument %s", key)
		}
		if bytes.Equal(bytes.TrimSpace(value), []byte("null")) {
			return "", fmt.Errorf("%s must not be null", key)
		}
	}
	if name != "tool_limits" {
		if _, ok := fields["handle"]; !ok {
			return "", fmt.Errorf("handle is required")
		}
	}
	if name == "tool_limits" {
		// The previous setting was consumed by this call. A new setting starts
		// from defaults, so omitted fields cannot inherit those spent limits.
		next, err := limits(raw, defaults())
		if err != nil {
			return "", err
		}
		a.pending = &next
		return "tool_limits set: the next call, whichever tool that is, will report the consumption.", nil
	}
	var args struct {
		Handle int     `json:"handle"`
		Input  *string `json:"input"`
	}
	if err := json.Unmarshal(raw, &args); err != nil {
		return "", err
	}
	j, err := a.parent.Jobs().Get(args.Handle)
	if err != nil {
		return "", err
	}
	if j.Parent() != a {
		return "", fmt.Errorf("job belongs to another agent")
	}
	if name == "kill_job" {
		if err := j.Kill("kill_job"); err != nil {
			return "", err
		}
		opts.Delay = 0
	}
	if name == "send_input" {
		if args.Input == nil {
			return "", fmt.Errorf("input is required")
		}
		if err := j.Input(*args.Input); err != nil {
			return "", err
		}
	}
	if name == "wait_for_job" {
		opts.ReturnAvailable = true
	}
	return j.Wait(ctx, opts)
}
