package ensemble

import (
	"fmt"
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
	return c, nil
}
