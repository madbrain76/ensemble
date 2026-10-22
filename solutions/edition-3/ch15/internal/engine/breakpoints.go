package engine

import (
	"encoding/json"
	"ensemble/internal/common"
)

type aCacheControl struct {
	Type string `json:"type"`
}

func ephemeralCache() *aCacheControl { return &aCacheControl{Type: "ephemeral"} }

// One retained-history endpoint, captured before ephemera/coalescing. The other
// conversation slot belongs only to the surviving handoff, never an arbitrary
// older message. Thinking, signatures, media and unknown blocks are protected.
func conversationBreakpoints(cfg common.Config, messages []aMessage) {
	if ok, _ := common.CacheBreakpoints(cfg); !ok {
		return
	}
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].Role != "user" && messages[i].Role != "assistant" {
			continue
		}
		for j := len(messages[i].Content) - 1; j >= 0; j-- {
			var block map[string]json.RawMessage
			if json.Unmarshal(messages[i].Content[j], &block) != nil {
				continue
			}
			var kind string
			_ = json.Unmarshal(block["type"], &kind)
			switch kind {
			case "text":
				var text string
				_ = json.Unmarshal(block["text"], &text)
				if text == "" {
					continue
				}
			case "tool_use", "tool_result":
			default:
				continue
			}
			// Insert only metadata. Re-marshaling this map would reorder keys
			// and could disturb opaque tool data when the marker later moves.
			if _, exists := block["cache_control"]; !exists {
				original := messages[i].Content[j]
				start := wireSpace(original, 0)
				marked := append(json.RawMessage(nil), original[:start+1]...)
				marked = append(marked, []byte(`"cache_control":{"type":"ephemeral"},`)...)
				messages[i].Content[j] = append(marked, original[start+1:]...)
			}
			return
		}
	}
}

func cacheEnabled(cfg common.Config) bool { ok, _ := common.CacheBreakpoints(cfg); return ok }

// Only text message-content blocks are legal OpenAI marker targets here.
func markOpenAIMessage(m *oMessage) bool {
	switch content := m.Content.(type) {
	case *string:
		if content == nil || *content == "" {
			return false
		}
		m.Content = []any{map[string]any{"type": "text", "text": *content, "prompt_cache_breakpoint": map[string]string{"mode": "explicit"}}}
		return true
	case []any:
		for i := len(content) - 1; i >= 0; i-- {
			block, ok := content[i].(map[string]any)
			if !ok || block["type"] != "text" {
				continue
			}
			text, ok := block["text"].(string)
			if !ok || text == "" {
				continue
			}
			block["prompt_cache_breakpoint"] = map[string]string{"mode": "explicit"}
			return true
		}
	}
	return false
}

// projectContext puts the typed surviving handoff at the ephemera front. Split
// that exact boundary before flattening text; never recognize a lookalike string.
func projectedTail(c common.Context) (handoff, transient string, err error) {
	parts := c.Ephemera
	if c.Handoff != "" && len(parts) > 0 && parts[0].Type == "text" && parts[0].Text == c.Handoff {
		handoff = parts[0].Text
		parts = parts[1:]
	}
	transient, err = textParts(parts)
	return
}
