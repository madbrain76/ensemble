package common

import (
	"os"
	"strings"
)

// ProcessEnvironment is shared by local command and MCP child launches. Provider
// credentials belong to the Agent's HTTP owner, not arbitrary local processes.
func ProcessEnvironment(agents ...Agent) []string {
	private := map[string]bool{}
	for _, a := range agents {
		servers := a.Config().MCPServers
		for _, skill := range a.Skills().ListSkills() {
			servers = append(servers, skill.MCPServers...)
		}
		for _, server := range servers {
			if server.AuthEnv != "" {
				private[server.AuthEnv] = true
			}
		}
	}
	env := []string{}
	for _, variable := range os.Environ() {
		name, _, _ := strings.Cut(variable, "=")
		if private[name] {
			continue
		}
		switch name {
		case "LLM_API_KEY", "ANTHROPIC_API_KEY", "OPENAI_API_KEY", "GEMINI_API_KEY", "FIRECRAWL_API_KEY", "TERM":
			continue
		}
		env = append(env, variable)
	}
	return append(env, "TERM=dumb")
}
