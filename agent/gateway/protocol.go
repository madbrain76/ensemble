package gateway

import "encoding/json"

// ProtocolVersion is the prologue mixed into every Noise-KK handshake.
//
// Both sides must use the same prologue or the handshake fails, which makes
// this a cheap version gate: an old client and a new gateway cannot complete
// a handshake, let alone exchange a malformed frame.
const ProtocolVersion = "ensemble-gateway-v1"

// FrameType names the kind of message inside an encrypted frame.
type FrameType string

const (
	// --- agent to gateway ---

	// FrameHello identifies the agent and requests resume from a sequence.
	FrameHello FrameType = "hello"
	// FrameAck confirms delivery of a sequence number.
	FrameAck FrameType = "ack"
	// FrameResponse replies to a previously delivered prompt.
	FrameResponse FrameType = "response"
	// FramePing is a keepalive.
	FramePing FrameType = "ping"
	// FrameMCPRequest carries a gateway-administration tool call.
	FrameMCPRequest FrameType = "mcp_request"

	// --- gateway to agent ---

	// FrameWelcome confirms the connection and reports queue depth.
	FrameWelcome FrameType = "welcome"
	// FrameDeliver delivers a prompt envelope.
	FrameDeliver FrameType = "deliver"
	// FramePong answers a keepalive.
	FramePong FrameType = "pong"
	// FrameError reports a protocol or authorization failure.
	FrameError FrameType = "error"
	// FrameMCPResponse carries the result of a tool call.
	FrameMCPResponse FrameType = "mcp_response"
)

// Frame is the envelope for everything sent inside the encrypted stream.
//
// Payload stays as raw JSON so the frame can be read, routed on Type, and
// only then decoded into the specific payload the type implies. A reader
// that had to know every payload shape in advance could not tolerate a peer
// that learned a new one.
type Frame struct {
	Type    FrameType       `json:"type"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

// HelloPayload is the agent's opening frame.
type HelloPayload struct {
	AgentID string `json:"agent_id"`
	Version string `json:"version"`
	// LastAckedSeq is the highest sequence number this agent has durably
	// handled. The gateway replays from LastAckedSeq+1, which is what makes
	// delivery at-least-once across a reconnect.
	LastAckedSeq uint64 `json:"last_acked_seq"`
}

// WelcomePayload confirms a connection.
type WelcomePayload struct {
	AgentID       string `json:"agent_id"`
	QueuedCount   int    `json:"queued_count"`
	ServerVersion string `json:"server_version"`
}

// DeliverPayload carries one envelope to its recipient.
type DeliverPayload struct {
	Seq      uint64   `json:"seq"`
	Envelope Envelope `json:"envelope"`
}

// AckPayload confirms receipt through Seq.
type AckPayload struct {
	Seq uint64 `json:"seq"`
}

// ResponsePayload is an agent's reply to a delivered prompt.
type ResponsePayload struct {
	MessageID string            `json:"message_id"`
	Content   string            `json:"content"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// ErrorPayload reports a failure to the agent.
type ErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Error codes the gateway returns. These are stable strings because callers
// branch on them; an error that can only be matched by its prose is an error
// that cannot be handled.
const (
	// ErrCodeUnauthorized means the sender holds no group the recipient
	// accepts contact from.
	ErrCodeUnauthorized = "unauthorized"
	// ErrCodeUnknownAgent means the named destination is not configured.
	ErrCodeUnknownAgent = "unknown_agent"
	// ErrCodeForbidden means the caller lacks a required capability.
	ErrCodeForbidden = "forbidden"
	// ErrCodeBadRequest means the frame was malformed.
	ErrCodeBadRequest = "bad_request"
)

// MCPRequestPayload is a gateway-administration tool call.
//
// Administration is deliberately not a separate HTTP API or admin CLI. An
// agent already knows how to call tools, so the gateway exposes its own
// controls the same way, over the connection it already has.
type MCPRequestPayload struct {
	// ID correlates the response with this request.
	ID string `json:"id"`
	// Tool is the tool name, such as "gateway_register_child".
	Tool string `json:"tool"`
	// Arguments is the tool's JSON argument object.
	Arguments json.RawMessage `json:"arguments,omitempty"`
}

// MCPResponsePayload is the result of a tool call.
type MCPResponsePayload struct {
	ID      string          `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	IsError bool            `json:"is_error,omitempty"`
	Error   string          `json:"error,omitempty"`
}

// Gateway administration tool names.
const (
	ToolDeclareContactable = "gateway_declare_contactable"
	ToolRegisterChild      = "gateway_register_child"
	ToolGrantMembership    = "gateway_grant_membership"
	ToolRevokeMembership   = "gateway_revoke_membership"
	ToolSendMessage        = "gateway_send_message"
	ToolListAgents         = "gateway_list_agents"
	ToolCreateGroup        = "gateway_create_group"
	ToolEditGroup          = "gateway_edit_group"
	ToolDeleteGroup        = "gateway_delete_group"
)

// Capabilities a group may confer.
const (
	// CapGroupAdmin permits creating, editing and deleting groups. It is the
	// capability that protects every other capability, so it belongs to a
	// root group that is seeded in config and granted sparingly.
	CapGroupAdmin = "group-admin"
	// CapSendMessage permits sending messages to other agents.
	CapSendMessage = "send-message"
	// CapDeclareContactable permits declaring inbound contactability.
	CapDeclareContactable = "declare-contactable"
)
