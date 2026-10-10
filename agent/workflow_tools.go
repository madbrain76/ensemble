package agent

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/waywardgeek/ensemble/agent/internal/common"
)

// defaultJoinTimeout bounds a fan-in that nobody gave a deadline.
//
// A wait with no bound is how a fleet hangs forever: one child never
// finishes and the parent never looks at anything again. A default is not a
// policy about how long work should take, it is a guarantee that control
// comes back.
const defaultJoinTimeout = 10 * time.Minute

// defaultAskTimeout bounds a child waiting on its parent for an answer.
const defaultAskTimeout = 5 * time.Minute

// SubAgents returns this agent's fleet manager.
func (a *Agent) SubAgents() *SubAgentManager { return a.subs }

// SetAgentFactory supplies the function that builds this agent's children.
//
// It is set after construction rather than passed to NewAgent because the
// factory usually closes over the very agent being built, and a value cannot
// be passed to the call that produces it.
func (a *Agent) SetAgentFactory(f AgentFactory) {
	a.subs.mu.Lock()
	defer a.subs.mu.Unlock()
	a.subs.factory = f
}

// registerWorkflowTools adds the tools a model calls to run a fleet.
//
// They are registered here, in the root package, as closures over this
// agent. The alternative is for the tools package to know what an agent is,
// which would point a dependency edge backwards: the root composes the
// library, the library does not reach up into the root.
//
// Each agent registers its own, so a tool call always acts on the agent that
// made it. A shared registry with an agent argument would let one agent name
// another and operate on it.
func (a *Agent) registerWorkflowTools() {
	// Supervision tools that return immediately run inline with NoJob. The
	// two that BLOCK deliberately do not: a blocking tool on the inline
	// path would hold the actor loop, and the whole point of stage 3 is
	// that a message can interrupt exactly these two.
	a.reg.RegisterTool(common.Tool{
		Name: "spawn_sub_agent",
		Description: "Create a sub-agent to work on a task in the background. " +
			"The child is a full agent with its own files. Use join_agents to collect results.",
		Schema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"agent_id": {"type": "string", "description": "Unique name for the child, e.g. 'db-expert'."},
				"task": {"type": "string", "description": "What the child should do."}
			},
			"required": ["agent_id", "task"]
		}`),
		NoJob: true,
		Run:   a.toolSpawnSubAgent,
	})

	a.reg.RegisterTool(common.Tool{
		Name: "join_agents",
		Description: "Wait for sub-agents to finish and collect their results. " +
			"Returns early if a child asks a question, so you can answer it.",
		Schema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"agent_ids": {"type": "array", "items": {"type": "string"}, "description": "Children to wait for."},
				"timeout_seconds": {"type": "number", "description": "Maximum seconds to wait."}
			},
			"required": ["agent_ids"]
		}`),
		Run: a.toolJoinAgents,
	})

	a.reg.RegisterTool(common.Tool{
		Name:        "agents_status",
		Description: "List this agent's children, what each is doing, and any questions waiting for an answer.",
		Schema:      json.RawMessage(`{"type": "object", "properties": {}}`),
		NoJob:       true,
		Run:         a.toolAgentsStatus,
	})

	a.reg.RegisterTool(common.Tool{
		Name:        "shutdown_agent",
		Description: "Stop a sub-agent and remove it from the fleet.",
		Schema: json.RawMessage(`{
			"type": "object",
			"properties": {"agent_id": {"type": "string"}},
			"required": ["agent_id"]
		}`),
		NoJob: true,
		Run:   a.toolShutdownAgent,
	})

	a.reg.RegisterTool(common.Tool{
		Name: "send_message_to_parent",
		Description: "Ask whoever you report to a question and wait for the answer. " +
			"Use this when you need a decision you cannot infer from your task. " +
			"This does NOT mean you are finished; call submit for that.",
		Schema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"content": {"type": "string", "description": "Your message."},
				"timeout_seconds": {"type": "number", "description": "Maximum seconds to wait for a reply."}
			},
			"required": ["content"]
		}`),
		Run: a.toolSendMessageToParent,
	})

	a.reg.RegisterTool(common.Tool{
		Name:        "respond_to_parent_message",
		Description: "Answer a question a sub-agent asked you, which unblocks it.",
		Schema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"agent_id": {"type": "string", "description": "The child that asked."},
				"message_id": {"type": "string", "description": "Which message you are answering."},
				"response": {"type": "string"}
			},
			"required": ["agent_id", "message_id", "response"]
		}`),
		NoJob: true,
		Run:   a.toolRespondToParentMessage,
	})

	a.reg.RegisterTool(common.Tool{
		Name: "submit",
		Description: "Report your finished result. This is the only way to say you are done. " +
			"Stopping without calling this is treated as a mistake and you will be woken to continue.",
		Schema: json.RawMessage(`{
			"type": "object",
			"properties": {"result": {"type": "string", "description": "Your result."}},
			"required": ["result"]
		}`),
		NoJob: true,
		Run:   a.toolSubmit,
	})
}

func (a *Agent) toolSpawnSubAgent(c *common.Call, args json.RawMessage) (string, error) {
	var in struct {
		AgentID string `json:"agent_id"`
		Task    string `json:"task"`
	}
	if err := json.Unmarshal(args, &in); err != nil {
		return "", fmt.Errorf("spawn_sub_agent: %w", err)
	}

	sub, err := a.subs.Spawn(in.AgentID, in.Task, nil)
	if err != nil {
		return "", err
	}
	if err := sub.Start(in.Task); err != nil {
		return "", err
	}
	return fmt.Sprintf("Spawned %s. It is running; use join_agents to collect its result.", sub.ID), nil
}

func (a *Agent) toolJoinAgents(c *common.Call, args json.RawMessage) (string, error) {
	var in struct {
		AgentIDs []string `json:"agent_ids"`
		Timeout  float64  `json:"timeout_seconds"`
	}
	if err := json.Unmarshal(args, &in); err != nil {
		return "", fmt.Errorf("join_agents: %w", err)
	}
	if len(in.AgentIDs) == 0 {
		return "", fmt.Errorf("join_agents: agent_ids must not be empty")
	}

	timeout := defaultJoinTimeout
	if in.Timeout > 0 {
		timeout = time.Duration(in.Timeout * float64(time.Second))
	}
	return a.subs.Join(in.AgentIDs, timeout).Render(), nil
}

func (a *Agent) toolAgentsStatus(c *common.Call, args json.RawMessage) (string, error) {
	return a.subs.Status(), nil
}

func (a *Agent) toolShutdownAgent(c *common.Call, args json.RawMessage) (string, error) {
	var in struct {
		AgentID string `json:"agent_id"`
	}
	if err := json.Unmarshal(args, &in); err != nil {
		return "", fmt.Errorf("shutdown_agent: %w", err)
	}
	if err := a.subs.Shutdown(in.AgentID); err != nil {
		return "", err
	}
	return fmt.Sprintf("Shut down %s.", in.AgentID), nil
}

func (a *Agent) toolSendMessageToParent(c *common.Call, args json.RawMessage) (string, error) {
	var in struct {
		Content string  `json:"content"`
		Timeout float64 `json:"timeout_seconds"`
	}
	if err := json.Unmarshal(args, &in); err != nil {
		return "", fmt.Errorf("send_message_to_parent: %w", err)
	}

	timeout := defaultAskTimeout
	if in.Timeout > 0 {
		timeout = time.Duration(in.Timeout * float64(time.Second))
	}

	reply, err := a.MessageParent(in.Content, timeout)
	if err != nil {
		return "", err
	}
	// A top-level agent's message goes to the human, who does not answer.
	if reply == nil {
		return "Delivered to the user.", nil
	}
	// The answer is rendered WITH its sender, so the model can address a
	// follow-up rather than only knowing what it was told.
	return fmt.Sprintf("Reply from %s: %s", reply.From, reply.Content), nil
}

func (a *Agent) toolRespondToParentMessage(c *common.Call, args json.RawMessage) (string, error) {
	var in struct {
		AgentID   string `json:"agent_id"`
		MessageID string `json:"message_id"`
		Response  string `json:"response"`
	}
	if err := json.Unmarshal(args, &in); err != nil {
		return "", fmt.Errorf("respond_to_parent_message: %w", err)
	}
	if err := a.subs.Respond(in.AgentID, in.MessageID, in.Response); err != nil {
		return "", err
	}
	return fmt.Sprintf("Answered %s; it has resumed.", in.AgentID), nil
}

func (a *Agent) toolSubmit(c *common.Call, args json.RawMessage) (string, error) {
	var in struct {
		Result string `json:"result"`
	}
	if err := json.Unmarshal(args, &in); err != nil {
		return "", fmt.Errorf("submit: %w", err)
	}
	a.Submit(in.Result)
	return "Result submitted.", nil
}
