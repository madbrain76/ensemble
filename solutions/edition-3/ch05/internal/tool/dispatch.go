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
	result := r.Prepare(ctx, common.Part{Name: name, Args: input}, "", 0).Await(ctx)
	return result.Parts, result.IsError
}

type operation struct {
	registry *registry
	call     common.Part
	args     arguments
	limits   common.ToolLimits
	custom   bool
	job      common.Job
	control  bool
	err      error
	ready    *common.ToolData
}

func (o *operation) Called() common.ToolData {
	d := common.ToolData{CallID: o.call.CallID, Name: o.call.Name, Args: o.call.Args}
	if o.job != nil {
		j := o.job.Data()
		d.Job = &j
	}
	return d
}
func (o *operation) Await(ctx context.Context) common.ToolData {
	if o.ready != nil {
		return *o.ready
	}
	d := o.Called()
	var text string
	err := o.err
	if err == nil && o.control {
		var report *common.JobReport
		text, report, err = o.registry.control(ctx, o.call.Name, o.args, o.limits)
		if report != nil {
			// Controls report an existing target; they never own a job.
			text = formatReport(*report)
			d.IsError = report.IsError
		}
	} else if err == nil {
		report, reportErr := o.job.InitialReport(ctx, o.limits)
		err = reportErr
		d.Job = &report.Job
		d.IsError = report.IsError
		text = report.Text
		if report.Job.Status != common.JobDone || report.Job.Bytes > int64(o.limits.MaxOutputBytes) || report.Job.Cwd != "" {
			text = formatReport(report)
		}
	}
	if err != nil {
		text = err.Error()
		d.IsError = true
	}
	if o.custom {
		pattern := "none"
		if o.limits.Pattern != "" {
			pattern = fmt.Sprintf("%q", o.limits.Pattern)
		}
		text = fmt.Sprintf("[tool_limits consumed by this %s call: ai_callback_delay %s, ai_callback_pattern %s, max_output_bytes %d]\n", o.call.Name, o.limits.Delay, pattern, o.limits.MaxOutputBytes) + text
	}
	d.Parts = common.PartList{common.Text(text)}
	return d
}
func (r *registry) Prepare(ctx context.Context, call common.Part, refusal string, turn uint64) common.ToolOperation {
	// Literal next means every position, including refusals and skipped setters.
	r.mu.Lock()
	o := &operation{registry: r, call: call, limits: common.DefaultToolLimits(), custom: r.next != nil}
	if o.custom {
		o.limits = *r.next
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
	o.args, o.err = decode(call.Args)
	if o.err == nil && (call.Name == "run_command" || call.Name == "wait_for_job" || call.Name == "send_input") {
		o.limits, o.err = parseLimits(o.args, o.limits)
	}
	if selected == nil {
		o.err = fmt.Errorf("unknown or disabled tool %q", call.Name)
	}
	if refusal != "" {
		o.err = errors.New(refusal)
	}
	// Rejection is a tool result, never execution or job allocation.
	if o.err != nil {
		return o
	}
	if selected.control {
		o.control = true
		if call.Name == "tool_limits" {
			d := o.Await(ctx)
			o.ready = &d
		}
		return o
	}
	handler := func(ctx context.Context, owner common.ToolCall, raw json.RawMessage) error {
		if selected.handler != nil {
			return selected.handler(ctx, owner, raw)
		}
		text, err := selected.execute(ctx, o.args)
		if err != nil {
			return err
		}
		_, err = io.WriteString(owner, text)
		return err
	}
	cwd := ""
	if call.Name == "run_command" {
		if override, err := o.args.text("cwd", false); err == nil && override != "" {
			cwd = r.path(override)
		}
	}
	origin := common.JobData{TurnID: turn, CallID: call.CallID, Name: call.Name, Args: append([]byte(nil), call.Args...)}
	o.job, o.err = r.parent.Jobs().Start(ctx, handler, call.Args, cwd, origin)
	return o
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
		definition("kill_job", "Kill and reap a managed process group. For Go-only work, mark killed and close output; Go cannot stop the handler. Repeated kills are safe.", `"handle":{"type":"integer","minimum":1}`, `"handle"`),
	}
}

// Invoke retains the direct Chapter 4 API without making the registry a History
// writer. The Agent admits, records, and joins this operation like other work.
func (r *registry) Invoke(ctx context.Context, call common.Part, refusal string) error {
	return r.parent.InvokeTool(ctx, call, refusal)
}
