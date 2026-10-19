package ensemble

import (
	"fmt"
	"strconv"
	"strings"
)

// ConfigFromEnv is shared by live and offline adapters; the core owns no
// process environment or default log filename. Empty generic values fall back.
func ConfigFromEnv(get func(string) string) (Config, error) {
	var c Config
	prefix := "ANTHROPIC"
	switch strings.ToLower(get("LLM_VENDOR")) {
	case "", "anthropic":
		c.Vendor = Anthropic
	case "openai":
		c.Vendor = OpenAI
		prefix = "OPENAI"
	case "gemini":
		c.Vendor = Gemini
		prefix = "GEMINI"
	default:
		return c, fmt.Errorf("unknown LLM_VENDOR")
	}
	value := func(suffix string) string {
		if v := get("LLM_" + suffix); v != "" {
			return v
		}
		return get(prefix + "_" + suffix)
	}
	c.Model = value("MODEL")
	c.APIKey = value("API_KEY")
	c.BaseURL = value("BASE_URL")
	c.ResolvedModel = get("LLM_RESOLVED_MODEL")
	switch selected := strings.TrimSpace(get("LLM_TOOLS")); selected {
	case "", "default": // nil means the default six built-ins.
	case "none":
		c.Tools = []string{}
	default:
		c.Tools = strings.Split(selected, ",")
		for i := range c.Tools {
			c.Tools[i] = strings.TrimSpace(c.Tools[i])
		}
	}
	if policy := get("LLM_REQUIRE_TOOL_EXPLANATION"); policy != "" {
		enabled, err := strconv.ParseBool(policy)
		if err != nil {
			return Config{}, fmt.Errorf("LLM_REQUIRE_TOOL_EXPLANATION must be a boolean")
		}
		c.RequireToolExplanation = enabled
	}
	return c, nil
}
