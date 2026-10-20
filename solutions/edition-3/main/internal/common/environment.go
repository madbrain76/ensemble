package common

import (
	"os"
	"strings"
)

// ProcessEnvironment is shared by local command and MCP child launches. Provider
// credentials belong to the Agent's HTTP owner, not arbitrary local processes.
func ProcessEnvironment() []string {
	env := []string{}
	for _, variable := range os.Environ() {
		name, _, _ := strings.Cut(variable, "=")
		switch name {
		case "LLM_API_KEY", "ANTHROPIC_API_KEY", "OPENAI_API_KEY", "GEMINI_API_KEY", "TERM":
			continue
		}
		env = append(env, variable)
	}
	return append(env, "TERM=dumb")
}
