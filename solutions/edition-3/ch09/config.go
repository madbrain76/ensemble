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
	case "", "default": // nil means the default ten built-ins.
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
	for _, option := range []struct {
		name string
		dst  *bool
	}{{"LLM_PAUSE_ON_START", &c.PauseOnStart}, {"LLM_SPEAKING", &c.Speaking}} {
		if value := get(option.name); value != "" {
			enabled, err := strconv.ParseBool(value)
			if err != nil {
				return Config{}, fmt.Errorf("%s must be a boolean", option.name)
			}
			*option.dst = enabled
		}
	}
	directory, primary := get("EN_SKILLS_DIR"), get("EN_PRIMARY_SKILL")
	if directory != "" || primary != "" {
		if directory == "" || primary == "" {
			return Config{}, fmt.Errorf("EN_SKILLS_DIR and EN_PRIMARY_SKILL must be supplied together")
		}
		custom := get("EN_CUSTOM_VAR")
		c.Skills = &SkillSpec{Directory: directory, Primary: primary, Variables: map[string]func() string{"CUSTOM_VAR": func() string { return custom }}}
	}
	c.LogPath = get("LLM_LOG_PATH")
	if rounds := get("LLM_MAX_TOOL_ROUNDS"); rounds != "" {
		n, err := strconv.Atoi(rounds)
		if err != nil || n < 0 {
			return Config{}, fmt.Errorf("LLM_MAX_TOOL_ROUNDS must be nonnegative")
		}
		c.MaxToolRounds = n
	}
	disabled := get("EN_DISABLE_STREAMING")
	if disabled == "" {
		disabled = get("LLM_DISABLE_STREAMING")
	}
	if disabled != "" {
		value, err := strconv.ParseBool(disabled)
		if err != nil {
			return Config{}, fmt.Errorf("EN_DISABLE_STREAMING must be a boolean")
		}
		c.DisableStreaming = value
	}
	c.StreamCapability = StreamCapability(get("LLM_STREAM_CAPABILITY"))
	if c.StreamCapability != StreamUnknown && c.StreamCapability != StreamPlain && c.StreamCapability != StreamSSE {
		return Config{}, fmt.Errorf("LLM_STREAM_CAPABILITY must be sse or plain")
	}
	return c, nil
}
