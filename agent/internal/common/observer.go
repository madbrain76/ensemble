package common

import "encoding/json"

// The outbound seam: what leaves the agent.
//
// waywardgeest — the ghost in the machine watches, and what it sees leaves
// through here.
//
// Observe MUST NOT BLOCK. An observer that blocks parks the actor's loop and
// re-creates the exact deafness this chapter exists to remove — the GUI would
// become able to freeze the agent by being slow. A slow observer buffers, or
// drops, on its own time.

// AgentID names which agent an observation came from.
//
// Empty means "the one agent", which is the only case the basic engine builds.
// It is omitempty everywhere, so the single-agent log is byte-identical to one
// written before this field existed.
type AgentID string

// Observer watches one or more agents. It never calls back into the agent.
type Observer interface {
	Observe(Observation)
}

// Observation is what an observer receives. Sealed union, same discipline as
// Part: the interface method is unexported, so only types in this package can
// implement it.
type Observation interface{ isObservation() }

// PartDelta is content arriving incrementally.
//
// THE RULE THAT KEEPS THIS SEAM STABLE: streaming is not a mode. A
// non-streaming vendor emits exactly one delta and then a final. Adding real
// streaming later therefore adds NO new observation kinds — only a different
// chunk count.
type PartDelta struct {
	Agent  AgentID   `json:"agent,omitempty"`
	PartID uint64    `json:"part_id"`
	Kind   DeltaKind `json:"kind"`
	Chunk  string    `json:"chunk"`
}

// PartFinal is the authoritative, complete part.
type PartFinal struct {
	Agent  AgentID `json:"agent,omitempty"`
	Seq    Seq     `json:"seq"`
	PartID uint64  `json:"part_id"`
	Part   Part    `json:"part"`
}

// StateChanged is a transition of the state machine. Chapter 5 does not
// introduce agent state. It EXPOSES it.
type StateChanged struct {
	Agent AgentID   `json:"agent,omitempty"`
	From  TurnState `json:"from"`
	To    TurnState `json:"to"`
}

// TurnEnded signals that a turn completed. Observers use this to know when
// the agent stopped processing and what its final reply text was.
type TurnEnded struct {
	Agent AgentID `json:"agent,omitempty"`
	Text  string  `json:"text"`
	Err   string  `json:"error,omitempty"`
}

// ToolDispatched fires when a tool call begins execution.
//
// Named to avoid collision with ToolCalled (the event log entry) and
// ToolCompleted (the mailbox inbound). The observation is what leaves; the
// event is what stays; the inbound is what arrives.
type ToolDispatched struct {
	Agent  AgentID         `json:"agent,omitempty"`
	CallID string          `json:"call_id"`
	Name   string          `json:"name"`
	Input  json.RawMessage `json:"input"`
}

// ToolFinished fires when a tool call completes (success or error).
type ToolFinished struct {
	Agent   AgentID `json:"agent,omitempty"`
	CallID  string  `json:"call_id"`
	Result  string  `json:"result"`
	IsError bool    `json:"is_error"`
}

// ConversationCleared fires when the operator resets the conversation. The
// reset is recorded as an event, so a client that replays the log rebuilds the
// cleared screen on its own. A client that is already connected replays
// nothing, and so has to be told.
type ConversationCleared struct {
	Agent AgentID `json:"agent,omitempty"`
}

// ActivityPhase says what an agent is doing right now at a granularity the
// operator can act on. The four SubAgentState values govern lifecycle (spawn,
// correct, escalate, shutdown); this governs display. "Running" is a single
// bucket that hides everything interesting, and this replaces it with four
// that do not.
type ActivityPhase int

const (
	// ActivityIdle means the agent is between turns.
	ActivityIdle ActivityPhase = iota + 1

	// ActivityThinking means the LLM is streaming a response.
	ActivityThinking

	// ActivityToolCall means the agent is executing a tool.
	ActivityToolCall

	// ActivityWaiting means the agent is blocked on a fleet management
	// tool: join_agents, wait_for_agent_change, ask_agent, or
	// send_message_to_parent.
	ActivityWaiting
)

// String renders a phase for status output and JSON.
func (p ActivityPhase) String() string {
	switch p {
	case ActivityIdle:
		return "idle"
	case ActivityThinking:
		return "thinking"
	case ActivityToolCall:
		return "tool_call"
	case ActivityWaiting:
		return "waiting"
	default:
		return "unknown"
	}
}

// ActivityChanged fires when an agent's activity phase transitions. The
// combination of Phase + ToolName + Detail gives the GUI enough to render
// a one-line summary of what every agent in the tree is doing.
type ActivityChanged struct {
	Agent    AgentID       `json:"agent,omitempty"`
	Phase    ActivityPhase `json:"phase"`
	ToolName string        `json:"tool_name,omitempty"`
	Detail   string        `json:"detail,omitempty"`
}

func (ActivityChanged) isObservation()      {}
func (ConversationCleared) isObservation() {}
func (PartDelta) isObservation()           {}
func (PartFinal) isObservation()           {}
func (StateChanged) isObservation()        {}
func (TurnEnded) isObservation()           {}
func (ToolDispatched) isObservation()      {}
func (ToolFinished) isObservation()        {}
