package common

import (
	"context"
	"encoding/json"
	"io"
)

// ToolDefinition is a wire-neutral declaration, not a handler or runtime state.
type ToolDefinition struct {
	Name        string
	Description string
	Parameters  json.RawMessage
}

// ToolRegistry owns one Agent's built-ins. Dispatch returns a result even on failure.
// Every ordinary dispatch uses the Agent-owned Jobs collection.
type ToolRegistry interface {
	Parent() Agent
	Definitions() []ToolDefinition
	Dispatch(context.Context, string, json.RawMessage) (PartList, bool)
	Invoke(context.Context, Part, string) error
	Register(ToolDefinition, ToolHandler) error
}

// ToolCall retains the actual Agent and streams output into its ordinary Job.
// A handler must cooperate with cancellation; late writes after kill are ignored.
type ToolCall interface {
	io.Writer
	Parent() Agent
	StartProcess(func() (io.Writer, func() error, error)) error
	SetExitCode(int)
}

type ToolHandler func(context.Context, ToolCall, json.RawMessage) error
