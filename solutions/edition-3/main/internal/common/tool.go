package common

import (
	"context"
	"encoding/json"
)

// ToolDefinition is a wire-neutral declaration, not a handler or runtime state.
type ToolDefinition struct {
	Name        string
	Description string
	Parameters  json.RawMessage
}

// ToolRegistry owns one Agent's built-ins. Dispatch returns a result even on failure.
// Calls are synchronous and serialized by the Agent; there are no job handles.
type ToolRegistry interface {
	Parent() Agent
	Definitions() []ToolDefinition
	Dispatch(context.Context, string, json.RawMessage) (PartList, bool)
}
