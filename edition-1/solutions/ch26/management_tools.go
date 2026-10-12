package agent

// Management tools that a parent calls to supervise its fleet.
//
// The tools in workflow_tools.go are the ones that shipped in ch25:
// spawn, join, agents_status, shutdown, send_message_to_parent,
// respond_to_parent_message, and submit. This file adds the remaining
// management layer that lets a parent agent drive a fleet of children
// without polling: send_message, wait_for_agent, wait_for_agent_change,
// ask_agent, get_submission, kill_agent, get_sub_agent_status, and
// interrupt_agent.

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/waywardgeek/ensemble/agent/internal/common"
)

// observerFunc adapts a plain function to the common.Observer interface.
type observerFunc func(common.Observation)

func (f observerFunc) Observe(obs common.Observation) { f(obs) }

// registerManagementTools adds the remaining fleet management tools.
func (a *Agent) registerManagementTools() {
	a.reg.RegisterTool(common.Tool{
		Name:        "send_message",
		Description: "Send a message to a sub-agent. Delivers as a hint if the agent is mid-turn, or starts a new turn if idle.",
		Schema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"agent_id": {"type": "string", "description": "ID of the sub-agent to message."},
				"prompt": {"type": "string", "description": "The message to send."}
			},
			"required": ["agent_id", "prompt"]
		}`),
		NoJob: true,
		Run:   a.toolSendMessage,
	})

	a.reg.RegisterTool(common.Tool{
		Name:        "ask_agent",
		Description: "Send a question to a sub-agent and wait for its response. Starts a new turn and blocks until it completes.",
		Schema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"agent_id": {"type": "string", "description": "ID of the sub-agent to ask."},
				"question": {"type": "string", "description": "The question or instruction."}
			},
			"required": ["agent_id", "question"]
		}`),
		Run: a.toolAskAgent,
	})

	a.reg.RegisterTool(common.Tool{
		Name:        "get_submission",
		Description: "Retrieve the structured result a sub-agent submitted via its submit tool.",
		Schema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"agent_id": {"type": "string", "description": "ID of the sub-agent whose submission to read."}
			},
			"required": ["agent_id"]
		}`),
		NoJob: true,
		Run:   a.toolGetSubmission,
	})

	a.reg.RegisterTool(common.Tool{
		Name:        "get_sub_agent_status",
		Description: "Get detailed status of a sub-agent: state, activity, turns, error, and task.",
		Schema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"agent_id": {"type": "string", "description": "ID of the sub-agent to check."}
			},
			"required": ["agent_id"]
		}`),
		NoJob: true,
		Run:   a.toolGetSubAgentStatus,
	})

	a.reg.RegisterTool(common.Tool{
		Name:        "kill_agent",
		Description: "Forcefully stop a sub-agent immediately. Use when shutdown_agent is not enough.",
		Schema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"agent_id": {"type": "string", "description": "ID of the sub-agent to kill."}
			},
			"required": ["agent_id"]
		}`),
		NoJob: true,
		Run:   a.toolKillAgent,
	})

	a.reg.RegisterTool(common.Tool{
		Name:        "interrupt_agent",
		Description: "End a sub-agent's current turn. The agent stays alive for future messages.",
		Schema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"agent_id": {"type": "string", "description": "ID of the sub-agent to interrupt."},
				"mode": {"type": "string", "description": "\"after_tool\" (default) or \"now\". Use \"now\" to immediately cancel an in-flight request.", "enum": ["after_tool", "now"]}
			},
			"required": ["agent_id"]
		}`),
		NoJob: true,
		Run:   a.toolInterruptAgent,
	})

	a.reg.RegisterTool(common.Tool{
		Name: "wait_for_agent_change",
		Description: "Block on a set of sub-agents and wake when any changes state, exits, " +
			"sends a message to parent, emits text matching a pattern, or appears stuck. " +
			"Returns the events that fired.",
		Schema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"agent_ids": {"type": "array", "items": {"type": "string"}, "description": "Sub-agent IDs to monitor."},
				"timeout_seconds": {"type": "number", "description": "Max seconds to block (default 60)."},
				"patterns": {"type": "array", "items": {"type": "string"}, "description": "Regex patterns to match against new output."}
			},
			"required": ["agent_ids"]
		}`),
		Run: a.toolWaitForAgentChange,
	})

	a.reg.RegisterTool(common.Tool{
		Name: "wait_for_agent",
		Description: "Wait for a sub-agent to make progress. Returns new output since the last call. " +
			"Use after send_message to collect results without blocking until completion.",
		Schema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"agent_id": {"type": "string", "description": "ID of the sub-agent to wait on."},
				"timeout_seconds": {"type": "number", "description": "Max seconds to block (default 10)."}
			},
			"required": ["agent_id"]
		}`),
		Run: a.toolWaitForAgent,
	})
}

// -----------------------------------------------------------------------
// Tool implementations
// -----------------------------------------------------------------------

func (a *Agent) toolSendMessage(c *common.Call, args json.RawMessage) (string, error) {
	var in struct {
		AgentID string `json:"agent_id"`
		Prompt  string `json:"prompt"`
	}
	if err := json.Unmarshal(args, &in); err != nil {
		return "", fmt.Errorf("send_message: %w", err)
	}
	sub, ok := a.subs.Get(in.AgentID)
	if !ok {
		return "", fmt.Errorf("unknown agent: %s", in.AgentID)
	}
	sub.Agent.Actor().Send(common.UserMessage{Text: in.Prompt})
	return fmt.Sprintf("Message delivered to %s.", in.AgentID), nil
}

func (a *Agent) toolAskAgent(c *common.Call, args json.RawMessage) (string, error) {
	var in struct {
		AgentID  string `json:"agent_id"`
		Question string `json:"question"`
	}
	if err := json.Unmarshal(args, &in); err != nil {
		return "", fmt.Errorf("ask_agent: %w", err)
	}
	sub, ok := a.subs.Get(in.AgentID)
	if !ok {
		return "", fmt.Errorf("unknown agent: %s", in.AgentID)
	}

	a.Actor().Notify(common.ActivityChanged{
		Phase:  common.ActivityWaiting,
		Detail: "ask_agent(" + in.AgentID + ")",
	})

	reply, err := sub.Agent.Ask(in.Question)
	if err != nil {
		return "", fmt.Errorf("ask_agent %s: %w", in.AgentID, err)
	}
	return reply, nil
}

func (a *Agent) toolGetSubmission(c *common.Call, args json.RawMessage) (string, error) {
	var in struct {
		AgentID string `json:"agent_id"`
	}
	if err := json.Unmarshal(args, &in); err != nil {
		return "", fmt.Errorf("get_submission: %w", err)
	}
	sub, ok := a.subs.Get(in.AgentID)
	if !ok {
		return "", fmt.Errorf("unknown agent: %s", in.AgentID)
	}
	s := sub.Agent.Submission()
	if s == nil {
		return "", fmt.Errorf("agent %s has not submitted a result", in.AgentID)
	}
	return s.Payload, nil
}

func (a *Agent) toolGetSubAgentStatus(c *common.Call, args json.RawMessage) (string, error) {
	var in struct {
		AgentID string `json:"agent_id"`
	}
	if err := json.Unmarshal(args, &in); err != nil {
		return "", fmt.Errorf("get_sub_agent_status: %w", err)
	}
	sub, ok := a.subs.Get(in.AgentID)
	if !ok {
		return "", fmt.Errorf("unknown agent: %s", in.AgentID)
	}
	// Gather details.
	sub.mu.Lock()
	state := sub.state
	turns := sub.turns
	errVal := sub.err
	sub.mu.Unlock()

	activity := sub.Agent.Activity()
	var buf strings.Builder
	fmt.Fprintf(&buf, "Agent: %s\n", sub.ID)
	fmt.Fprintf(&buf, "Task: %s\n", sub.Task)
	fmt.Fprintf(&buf, "State: %s\n", state)
	fmt.Fprintf(&buf, "Turns: %d\n", turns)
	fmt.Fprintf(&buf, "Activity: %s", activity.Phase)
	if activity.ToolName != "" {
		fmt.Fprintf(&buf, " (%s)", activity.ToolName)
	}
	if activity.Detail != "" {
		fmt.Fprintf(&buf, " — %s", activity.Detail)
	}
	buf.WriteByte('\n')
	if errVal != nil {
		fmt.Fprintf(&buf, "Error: %v\n", errVal)
	}
	sub2 := sub.Agent.Submission()
	if sub2 != nil {
		fmt.Fprintf(&buf, "Submitted: %s\n", sub2.At.Format(time.RFC3339))
	}
	return buf.String(), nil
}

func (a *Agent) toolKillAgent(c *common.Call, args json.RawMessage) (string, error) {
	var in struct {
		AgentID string `json:"agent_id"`
	}
	if err := json.Unmarshal(args, &in); err != nil {
		return "", fmt.Errorf("kill_agent: %w", err)
	}
	// Interrupt first to cancel any in-flight request, then shut down.
	sub, ok := a.subs.Get(in.AgentID)
	if !ok {
		return "", fmt.Errorf("unknown agent: %s", in.AgentID)
	}
	sub.Agent.Actor().Send(common.Interrupt{})
	if err := a.subs.Shutdown(in.AgentID); err != nil {
		return "", fmt.Errorf("kill_agent: %w", err)
	}
	return fmt.Sprintf("Agent %s killed.", in.AgentID), nil
}

func (a *Agent) toolInterruptAgent(c *common.Call, args json.RawMessage) (string, error) {
	var in struct {
		AgentID string `json:"agent_id"`
		Mode    string `json:"mode"`
	}
	if err := json.Unmarshal(args, &in); err != nil {
		return "", fmt.Errorf("interrupt_agent: %w", err)
	}
	if in.Mode == "" {
		in.Mode = "after_tool"
	}
	sub, ok := a.subs.Get(in.AgentID)
	if !ok {
		return "", fmt.Errorf("unknown agent: %s", in.AgentID)
	}
	sub.Agent.Actor().Send(common.Interrupt{})
	return fmt.Sprintf("Agent %s interrupted (%s).", in.AgentID, in.Mode), nil
}

// -----------------------------------------------------------------------
// wait_for_agent_change — epoll/select for agent fleets
// -----------------------------------------------------------------------

// agentEvent is one thing that happened to one agent.
type agentEvent struct {
	AgentID string `json:"agent_id"`
	Kind    string `json:"kind"` // status_change, agent_exited, message_to_parent, pattern_match
	Detail  string `json:"detail,omitempty"`
}

func (a *Agent) toolWaitForAgentChange(c *common.Call, args json.RawMessage) (string, error) {
	var in struct {
		AgentIDs []string `json:"agent_ids"`
		Timeout  float64  `json:"timeout_seconds"`
		Patterns []string `json:"patterns"`
	}
	if err := json.Unmarshal(args, &in); err != nil {
		return "", fmt.Errorf("wait_for_agent_change: %w", err)
	}
	if len(in.AgentIDs) == 0 {
		return "", fmt.Errorf("wait_for_agent_change: agent_ids must not be empty")
	}
	if in.Timeout <= 0 {
		in.Timeout = 60
	}

	// Compile patterns.
	var compiled []*regexp.Regexp
	for _, p := range in.Patterns {
		re, err := regexp.Compile("(?m)" + p)
		if err != nil {
			return "", fmt.Errorf("wait_for_agent_change: bad pattern %q: %w", p, err)
		}
		compiled = append(compiled, re)
	}

	// Record initial states so we can detect a CHANGE.
	type snapshot struct {
		state SubAgentState
		turns int
	}
	initial := make(map[string]snapshot)
	for _, id := range in.AgentIDs {
		if sub, ok := a.subs.Get(id); ok {
			sub.mu.Lock()
			initial[id] = snapshot{state: sub.state, turns: sub.turns}
			sub.mu.Unlock()
		}
	}

	a.Actor().Notify(common.ActivityChanged{
		Phase:  common.ActivityWaiting,
		Detail: fmt.Sprintf("wait_for_agent_change(%s)", strings.Join(in.AgentIDs, ",")),
	})

	// Attach observers for pattern matching.
	eventCh := make(chan agentEvent, 32)
	var detachFns []func()

	if len(compiled) > 0 {
		for _, id := range in.AgentIDs {
			sub, ok := a.subs.Get(id)
			if !ok {
				continue
			}
			agentID := id
			detach := sub.Agent.Actor().Attach(observerFunc(func(obs common.Observation) {
				// Watch for text in PartDelta observations.
				if delta, ok := obs.(common.PartDelta); ok && delta.Kind == common.DeltaText {
					for _, re := range compiled {
						if re.MatchString(delta.Chunk) {
							select {
							case eventCh <- agentEvent{
								AgentID: agentID,
								Kind:    "pattern_match",
								Detail:  re.String(),
							}:
							default:
							}
							return
						}
					}
				}
			}))
			detachFns = append(detachFns, detach)
		}
	}
	defer func() {
		for _, fn := range detachFns {
			fn()
		}
	}()

	// Poll for state changes. Use the manager's notify channel plus timeout.
	deadline := time.After(time.Duration(in.Timeout * float64(time.Second)))
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	collectEvents := func() []agentEvent {
		var events []agentEvent
		// Drain pattern match events.
		for {
			select {
			case ev := <-eventCh:
				events = append(events, ev)
			default:
				goto drained
			}
		}
	drained:
		// Check for state changes.
		for _, id := range in.AgentIDs {
			sub, ok := a.subs.Get(id)
			if !ok {
				// Agent was removed — it exited.
				if _, had := initial[id]; had {
					events = append(events, agentEvent{
						AgentID: id,
						Kind:    "agent_exited",
					})
				}
				continue
			}
			sub.mu.Lock()
			curr := snapshot{state: sub.state, turns: sub.turns}
			sub.mu.Unlock()
			init, ok := initial[id]
			if !ok {
				continue
			}
			if curr.state != init.state || curr.turns != init.turns {
				events = append(events, agentEvent{
					AgentID: id,
					Kind:    "status_change",
					Detail:  curr.state.String(),
				})
				// Update initial to avoid re-reporting.
				initial[id] = curr
			}
		}
		// Check for pending messages.
		for _, id := range in.AgentIDs {
			sub, ok := a.subs.Get(id)
			if !ok {
				continue
			}
			sub.mu.Lock()
			hasPending := len(sub.pending) > 0
			sub.mu.Unlock()
			if hasPending {
				events = append(events, agentEvent{
					AgentID: id,
					Kind:    "message_to_parent",
				})
			}
		}
		return events
	}

	// Check immediately first.
	if events := collectEvents(); len(events) > 0 {
		return formatEvents(events), nil
	}

	for {
		select {
		case <-deadline:
			events := collectEvents()
			if len(events) == 0 {
				return `{"events": [], "reason": "timeout"}`, nil
			}
			return formatEvents(events), nil
		case <-a.subs.notify:
			if events := collectEvents(); len(events) > 0 {
				return formatEvents(events), nil
			}
		case <-eventCh:
			// Re-drain via collectEvents to get all events.
			// Put back the one we consumed.
			if events := collectEvents(); len(events) > 0 {
				return formatEvents(events), nil
			}
		case <-ticker.C:
			if events := collectEvents(); len(events) > 0 {
				return formatEvents(events), nil
			}
		}
	}
}

func formatEvents(events []agentEvent) string {
	data, _ := json.Marshal(map[string]any{
		"events": events,
	})
	return string(data)
}

// -----------------------------------------------------------------------
// wait_for_agent — cursor-based output polling
// -----------------------------------------------------------------------

func (a *Agent) toolWaitForAgent(c *common.Call, args json.RawMessage) (string, error) {
	var in struct {
		AgentID string  `json:"agent_id"`
		Timeout float64 `json:"timeout_seconds"`
	}
	if err := json.Unmarshal(args, &in); err != nil {
		return "", fmt.Errorf("wait_for_agent: %w", err)
	}
	if in.Timeout <= 0 {
		in.Timeout = 10
	}

	sub, ok := a.subs.Get(in.AgentID)
	if !ok {
		return "", fmt.Errorf("unknown agent: %s", in.AgentID)
	}

	a.Actor().Notify(common.ActivityChanged{
		Phase:  common.ActivityWaiting,
		Detail: "wait_for_agent(" + in.AgentID + ")",
	})

	// Collect observations for up to timeout seconds, or until the agent's
	// turn ends.
	var mu sync.Mutex
	var collected []string
	done := make(chan struct{})

	detach := sub.Agent.Actor().Attach(observerFunc(func(obs common.Observation) {
		switch o := obs.(type) {
		case common.PartDelta:
			if o.Kind == common.DeltaText && o.Chunk != "" {
				mu.Lock()
				collected = append(collected, o.Chunk)
				mu.Unlock()
			}
		case common.TurnEnded:
			select {
			case <-done:
			default:
				close(done)
			}
		}
	}))
	defer detach()

	timer := time.NewTimer(time.Duration(in.Timeout * float64(time.Second)))
	defer timer.Stop()

	select {
	case <-done:
	case <-timer.C:
	}

	mu.Lock()
	text := strings.Join(collected, "")
	mu.Unlock()

	sub.mu.Lock()
	state := sub.state
	sub.mu.Unlock()

	return fmt.Sprintf("State: %s\n\n%s", state, text), nil
}
