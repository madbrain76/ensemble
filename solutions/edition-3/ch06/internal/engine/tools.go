package engine

import (
	"encoding/json"
	"ensemble/internal/common"
)

// Only Engine knows the three provider declaration envelopes.
type aDefinition struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InputSchema json.RawMessage `json:"input_schema"`
}
type functionDefinition struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Parameters  json.RawMessage `json:"parameters"`
}
type oDefinition struct {
	Type     string             `json:"type"`
	Function functionDefinition `json:"function"`
}
type gDefinitions struct {
	Functions []functionDefinition `json:"functionDeclarations"`
}

func anthropicDefinitions(defs []common.ToolDefinition) []aDefinition {
	var out []aDefinition
	for _, d := range defs {
		out = append(out, aDefinition{d.Name, d.Description, d.Parameters})
	}
	return out
}
func openAIDefinitions(defs []common.ToolDefinition) []oDefinition {
	var out []oDefinition
	for _, d := range defs {
		out = append(out, oDefinition{"function", functionDefinition{d.Name, d.Description, d.Parameters}})
	}
	return out
}
func geminiDefinitions(defs []common.ToolDefinition) ([]gDefinitions, error) {
	if len(defs) == 0 {
		return nil, nil
	}
	var functions []functionDefinition
	for _, d := range defs {
		// Gemini's current parameters dialect rejects this top-level JSON Schema
		// keyword. Project only that known difference; do not mutate the common
		// declaration or imply support for arbitrary future schemas.
		var parameters map[string]json.RawMessage
		if err := json.Unmarshal(d.Parameters, &parameters); err != nil {
			return nil, err
		}
		delete(parameters, "additionalProperties")
		projected, err := json.Marshal(parameters)
		if err != nil {
			return nil, err
		}
		functions = append(functions, functionDefinition{d.Name, d.Description, projected})
	}
	return []gDefinitions{{Functions: functions}}, nil
}
