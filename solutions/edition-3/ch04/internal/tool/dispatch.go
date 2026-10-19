package tool

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"regexp"
	"strings"
	"time"

	"ensemble/internal/common"
)

// Register extends only this Agent's registry, using the same ordinary job path.
func (r *registry) Register(d common.ToolDefinition, handler common.ToolHandler) error {
	if strings.TrimSpace(d.Name) == "" || strings.TrimSpace(d.Description) == "" || handler == nil {
		return errors.New("tool requires name, description, and handler")
	}
	schema, err := decode(d.Parameters)
	if err != nil {
		return fmt.Errorf("tool schema: %w", err)
	}
	kind, err := schema.text("type", true)
	if err != nil || kind != "object" {
		return errors.New("tool schema must have type object")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, e := range r.builtins() {
		if e.definition.Name == d.Name {
			return fmt.Errorf("reserved built-in name %q", d.Name)
		}
	}
	for _, e := range r.entries {
		if e.definition.Name == d.Name {
			return fmt.Errorf("duplicate tool %q", d.Name)
		}
	}
	d.Parameters = append(json.RawMessage(nil), d.Parameters...)
	r.entries = append(r.entries, entry{definition: d, handler: handler})
	return nil
}
func (r *registry) Dispatch(ctx context.Context, name string, input json.RawMessage) (common.PartList, bool) {
	parts, failed, err := r.dispatch(ctx, common.Part{Name: name, Args: input}, "", false)
	if err != nil {
		return common.PartList{common.Text(err.Error())}, true
	}
	return parts, failed
}
func (r *registry) Invoke(ctx context.Context, call common.Part, refusal string) error {
	_, _, err := r.dispatch(ctx, call, refusal, true)
	return err
}
func (r *registry) dispatch(ctx context.Context, call common.Part, refusal string, record bool) (common.PartList, bool, error) {
	// Consume before lookup, validation, refusal, setters, and supervision alike.
	r.mu.Lock()
	limits := common.DefaultToolLimits()
	custom := r.next != nil
	if custom {
		limits = *r.next
	}
	r.next = nil
	var selected *entry
	for _, e := range r.entries {
		if e.definition.Name == call.Name {
			copy := e
			selected = &copy
			break
		}
	}
	r.mu.Unlock()
	args, argErr := decode(call.Args)
	if argErr == nil && (call.Name == "run_command" || call.Name == "wait_for_job" || call.Name == "send_input") {
		limits, argErr = parseLimits(args, limits)
	}
	if selected == nil {
		argErr = fmt.Errorf("unknown or disabled tool %q", call.Name)
	}
	if refusal != "" {
		argErr = errors.New(refusal)
	}
	calledAlready := false
	var eventErr error
	called := func(data *common.JobData) error {
		if record {
			eventErr = r.parent.Log().Append(common.Event{Type: common.ToolCalled, Tool: &common.ToolData{CallID: call.CallID, Name: call.Name, Args: call.Args, Job: data}})
		}
		calledAlready = eventErr == nil
		return eventErr
	}
	result := func(text string, failed bool, data *common.JobData) (common.PartList, bool, error) {
		if custom {
			pattern := "none"
			if limits.Pattern != "" {
				pattern = fmt.Sprintf("%q", limits.Pattern)
			}
			text = fmt.Sprintf("[tool_limits consumed by this %s call: ai_callback_delay %s, ai_callback_pattern %s, max_output_bytes %d]\n", call.Name, limits.Delay, pattern, limits.MaxOutputBytes) + text
		}
		parts := common.PartList{common.Text(text)}
		if record {
			if err := r.parent.Log().Append(common.Event{Type: common.ToolReturned, Tool: &common.ToolData{CallID: call.CallID, Name: call.Name, Parts: parts, IsError: failed, Job: data}}); err != nil {
				return nil, true, err
			}
		}
		return parts, failed, nil
	}
	if selected != nil && selected.control {
		if err := called(nil); err != nil {
			return nil, true, err
		}
		var text string
		var report *common.JobReport
		if argErr == nil {
			text, report, argErr = r.control(ctx, call.Name, args, limits)
		}
		failed := argErr != nil
		if report != nil {
			text = formatReport(*report)
			failed = failed || report.IsError
		}
		if argErr != nil {
			text = argErr.Error()
		}
		return result(text, failed, nil) // Supervision observes another job, never owns one.
	}
	handler := func(ctx context.Context, owner common.ToolCall, raw json.RawMessage) error {
		if argErr != nil {
			return argErr
		}
		if selected.handler != nil {
			return selected.handler(ctx, owner, raw)
		}
		text, err := selected.execute(ctx, args)
		if err != nil {
			return err
		}
		_, err = io.WriteString(owner, text)
		return err
	}
	cwd := ""
	if call.Name == "run_command" && args != nil {
		if override, err := args.text("cwd", false); err == nil && override != "" {
			cwd = r.path(override)
		}
	}
	job, err := r.parent.Jobs().Start(ctx, handler, call.Args, cwd, func(data common.JobData) error { return called(&data) })
	if err != nil {
		// Failed admission/output setup must still pair with this exact call. No
		// handler ran without its output owner; the remaining batch can continue.
		if eventErr != nil {
			return nil, true, eventErr
		}
		var data *common.JobData
		if job != nil {
			snapshot := job.Data()
			data = &snapshot
		}
		if !calledAlready {
			if eventErr = called(data); eventErr != nil {
				return nil, true, eventErr
			}
		}
		return result("preparing job: "+err.Error(), true, data)
	}
	observation := limits
	if argErr != nil {
		observation = common.DefaultToolLimits()
	}
	report, err := job.InitialReport(ctx, observation)
	if err != nil {
		data := job.Data()
		return result(err.Error(), true, &data)
	}
	text := report.Text
	if report.Job.Status != common.JobDone || report.Job.Bytes > int64(observation.MaxOutputBytes) || report.Job.Cwd != "" {
		text = formatReport(report)
	}
	return result(text, report.IsError, &report.Job)
}
func formatReport(r common.JobReport) string {
	d := r.Job
	output := d.Output.Locator
	if output == "" {
		output = "unavailable"
	}
	header := fmt.Sprintf("job %d %s; %d bytes total; output: %s", d.Handle, d.Status, d.Bytes, output)
	if d.Cwd != "" {
		header += "; cwd: " + d.Cwd
	}
	if d.ExitCode != nil && d.Status != common.JobKilledStatus {
		header += fmt.Sprintf("; exit_code: %d", *d.ExitCode)
	}
	if d.Reason != "" {
		header += "; reason: " + d.Reason
	}
	if r.Text != "" {
		header += "\n" + r.Text
	}
	return header
}
func parseLimits(a arguments, l common.ToolLimits) (common.ToolLimits, error) {
	if raw, ok := a["ai_callback_delay"]; ok {
		var seconds float64
		if string(raw) == "null" || json.Unmarshal(raw, &seconds) != nil || math.IsNaN(seconds) || math.IsInf(seconds, 0) || seconds < 0 || seconds >= float64(math.MaxInt64)/float64(time.Second) {
			return l, errors.New("ai_callback_delay must be finite nonnegative seconds within duration range")
		}
		l.Delay = time.Duration(seconds * float64(time.Second))
	}
	if _, ok := a["ai_callback_pattern"]; ok {
		pattern, err := a.text("ai_callback_pattern", false)
		if err != nil {
			return l, err
		}
		if _, err = regexp.Compile(pattern); err != nil {
			return l, fmt.Errorf("ai_callback_pattern: %w", err)
		}
		l.Pattern = pattern
	}
	if _, ok := a["max_output_bytes"]; ok {
		n, err := a.integer("max_output_bytes", 0)
		if err != nil {
			return l, err
		}
		if n <= 0 {
			return l, errors.New("max_output_bytes must be a positive integer")
		}
		l.MaxOutputBytes = n
	}
	return l, nil
}
func (r *registry) control(ctx context.Context, name string, a arguments, limits common.ToolLimits) (string, *common.JobReport, error) {
	if name == "tool_limits" {
		if err := a.only("ai_callback_delay", "ai_callback_pattern", "max_output_bytes"); err != nil {
			return "", nil, err
		}
		next, err := parseLimits(a, common.DefaultToolLimits())
		if err != nil {
			return "", nil, err
		}
		r.mu.Lock()
		r.next = &next
		r.mu.Unlock()
		return fmt.Sprintf("Limits for the next call, whichever tool that is: delay=%gs pattern=%q max_output_bytes=%d. That call will report the consumption; these settings are then spent.", next.Delay.Seconds(), next.Pattern, next.MaxOutputBytes), nil, nil
	}
	allowed := []string{"handle"}
	if name != "kill_job" {
		allowed = append(allowed, "ai_callback_delay", "ai_callback_pattern", "max_output_bytes")
	}
	if name == "send_input" {
		allowed = append(allowed, "input")
	}
	if err := a.only(allowed...); err != nil {
		return "", nil, err
	}
	handle, err := a.integer("handle", 0)
	if err != nil || handle <= 0 {
		return "", nil, errors.New("handle must be a positive integer")
	}
	// Validate everything before any input or kill side effect.
	limits, err = parseLimits(a, limits)
	if err != nil {
		return "", nil, err
	}
	input := ""
	if name == "send_input" {
		input, err = a.text("input", true)
		if err != nil {
			return "", nil, err
		}
	}
	job, err := r.parent.Jobs().Get(handle)
	if err != nil {
		return "", nil, err
	}
	var report common.JobReport
	switch name {
	case "wait_for_job":
		report, err = job.Wait(ctx, limits)
	case "send_input":
		report, err = job.SendInput(ctx, input, limits)
	case "kill_job":
		err = job.Kill("kill_job")
		if err == nil {
			limits.Delay = 0
			limits.Pattern = ""
			report, err = job.Wait(ctx, limits)
		}
	default:
		err = fmt.Errorf("unknown control %q", name)
	}
	if err != nil {
		return "", nil, err
	}
	// A successful kill operation is not itself a failed tool call.
	if name == "kill_job" {
		report.IsError = false
	}
	return "", &report, nil
}
func controlDefinitions() []entry {
	limits := `"ai_callback_delay":{"type":"number","minimum":0},"ai_callback_pattern":{"type":"string"},"max_output_bytes":{"type":"integer","minimum":1}`
	definition := func(name, description, properties, required string) entry {
		return entry{control: true, definition: common.ToolDefinition{Name: name, Description: description, Parameters: json.RawMessage(`{"type":"object","properties":{` + properties + `},"required":[` + required + `],"additionalProperties":false}`)}}
	}
	return []entry{
		definition("tool_limits", "Set observation limits for the literal next tool call only, including errors, setters, and supervision. Omitted values are fresh defaults: 3 seconds, no pattern, 16384 bytes. Limits never kill jobs.", limits, ""),
		definition("wait_for_job", "Observe new output from this Agent's job, including a finished job. Delay or a pattern on unseen output ends observation, never execution. Omitted delay is 3 seconds; cap is 16384 bytes.", `"handle":{"type":"integer","minimum":1},`+limits, `"handle"`),
		definition("send_input", "Submit a line to a running PTY job, appending a terminal newline if absent and preserving an existing terminator. Then observe with non-destructive callback limits; no raw-input mode.", `"handle":{"type":"integer","minimum":1},"input":{"type":"string"},`+limits, `"handle","input"`),
		definition("kill_job", "Explicitly stop a job and its process group, close output, reap the process, and record kill_job. Repeated kills are safe.", `"handle":{"type":"integer","minimum":1}`, `"handle"`),
	}
}
