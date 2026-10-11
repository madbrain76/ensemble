package gateway

import (
	"encoding/json"
	"fmt"

	pub "github.com/waywardgeek/ensemble/agent/gateway"
)

// Gateway administration is exposed as MCP tool calls tunneled over the
// same Noise-KK connection the agent already holds.
//
// There is no admin CLI and no second HTTP API. An agent already knows how
// to call tools, and a second control plane would be a second place to get
// authorization wrong. Every handler below starts from the caller's
// authenticated identity, which came from the handshake and cannot be
// supplied as an argument.

type declareContactableArgs struct {
	Groups []string `json:"groups"`
}

type registerChildArgs struct {
	AgentID   string   `json:"agent_id"`
	PublicKey string   `json:"public_key"`
	Groups    []string `json:"groups"`
}

type membershipArgs struct {
	AgentID string `json:"agent_id"`
	Group   string `json:"group"`
}

type sendMessageArgs struct {
	ToAgent        string `json:"to_agent"`
	Content        string `json:"content"`
	ConversationID string `json:"conversation_id,omitempty"`
	ResponseMode   string `json:"response_mode,omitempty"`
}

type groupArgs struct {
	Name         string    `json:"name"`
	Description  string    `json:"description,omitempty"`
	Capabilities []string  `json:"capabilities,omitempty"`
	Trust        pub.Trust `json:"trust,omitempty"`
}

// handleMCPRequest dispatches one administration tool call and replies.
func (s *Server) handleMCPRequest(c *conn, frame *pub.Frame) error {
	var req pub.MCPRequestPayload
	if err := json.Unmarshal(frame.Payload, &req); err != nil {
		return err
	}

	result, err := s.dispatchTool(c.agentID, req)
	if err != nil {
		return c.sendFrame(pub.FrameMCPResponse, pub.MCPResponsePayload{
			ID:      req.ID,
			IsError: true,
			Error:   err.Error(),
		})
	}
	encoded, err := json.Marshal(result)
	if err != nil {
		return err
	}
	return c.sendFrame(pub.FrameMCPResponse, pub.MCPResponsePayload{ID: req.ID, Result: encoded})
}

// dispatchTool runs one tool on behalf of callerID.
//
// callerID is the identity proved by the Noise handshake. No tool takes a
// caller argument, so no tool can be asked to act as somebody else.
func (s *Server) dispatchTool(callerID string, req pub.MCPRequestPayload) (any, error) {
	switch req.Tool {
	case pub.ToolDeclareContactable:
		var a declareContactableArgs
		if err := decodeArgs(req.Arguments, &a); err != nil {
			return nil, err
		}
		if !s.registry.HasCapability(callerID, pub.CapDeclareContactable) {
			return nil, fmt.Errorf("%s lacks the %s capability", callerID, pub.CapDeclareContactable)
		}
		if err := s.registry.DeclareContactable(callerID, a.Groups); err != nil {
			return nil, err
		}
		return map[string]any{"contactable_by": s.registry.ContactableBy(callerID)}, nil

	case pub.ToolRegisterChild:
		var a registerChildArgs
		if err := decodeArgs(req.Arguments, &a); err != nil {
			return nil, err
		}
		if err := s.registry.RegisterChild(callerID, a.AgentID, a.PublicKey, a.Groups); err != nil {
			return nil, err
		}
		return map[string]any{
			"agent_id": a.AgentID,
			"groups":   s.registry.GroupsOf(a.AgentID),
			"parent":   callerID,
		}, nil

	case pub.ToolGrantMembership:
		var a membershipArgs
		if err := decodeArgs(req.Arguments, &a); err != nil {
			return nil, err
		}
		if err := s.registry.GrantMembership(callerID, a.AgentID, a.Group); err != nil {
			return nil, err
		}
		return map[string]any{"agent_id": a.AgentID, "groups": s.registry.GroupsOf(a.AgentID)}, nil

	case pub.ToolRevokeMembership:
		var a membershipArgs
		if err := decodeArgs(req.Arguments, &a); err != nil {
			return nil, err
		}
		if err := s.registry.RevokeMembership(callerID, a.AgentID, a.Group); err != nil {
			return nil, err
		}
		return map[string]any{"agent_id": a.AgentID, "groups": s.registry.GroupsOf(a.AgentID)}, nil

	case pub.ToolSendMessage:
		var a sendMessageArgs
		if err := decodeArgs(req.Arguments, &a); err != nil {
			return nil, err
		}
		env, err := s.SendMessage(callerID, pub.OutboundMessage{
			ToAgent:        a.ToAgent,
			Content:        a.Content,
			ConversationID: a.ConversationID,
			ResponseMode:   pub.ResponseMode(a.ResponseMode),
		})
		if err != nil {
			return nil, err
		}
		return map[string]any{"message_id": env.MessageID, "seq": env.Seq}, nil

	case pub.ToolListAgents:
		return map[string]any{"agents": s.registry.ListAgents(), "online": s.conns.online()}, nil

	case pub.ToolCreateGroup:
		var a groupArgs
		if err := decodeArgs(req.Arguments, &a); err != nil {
			return nil, err
		}
		if err := s.registry.CreateGroup(callerID, a.Name, GroupConfig{
			Description:  a.Description,
			Capabilities: a.Capabilities,
			Trust:        a.Trust,
		}); err != nil {
			return nil, err
		}
		return map[string]any{"group": a.Name}, nil

	case pub.ToolEditGroup:
		var a groupArgs
		if err := decodeArgs(req.Arguments, &a); err != nil {
			return nil, err
		}
		if err := s.registry.EditGroup(callerID, a.Name, GroupConfig{
			Description:  a.Description,
			Capabilities: a.Capabilities,
			Trust:        a.Trust,
		}); err != nil {
			return nil, err
		}
		return map[string]any{"group": a.Name}, nil

	case pub.ToolDeleteGroup:
		var a groupArgs
		if err := decodeArgs(req.Arguments, &a); err != nil {
			return nil, err
		}
		if err := s.registry.DeleteGroup(callerID, a.Name); err != nil {
			return nil, err
		}
		return map[string]any{"deleted": a.Name}, nil

	default:
		return nil, fmt.Errorf("unknown gateway tool %q", req.Tool)
	}
}

func decodeArgs(raw json.RawMessage, dst any) error {
	if len(raw) == 0 {
		return nil
	}
	if err := json.Unmarshal(raw, dst); err != nil {
		return fmt.Errorf("invalid arguments: %w", err)
	}
	return nil
}
