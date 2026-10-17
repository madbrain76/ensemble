package engine

import (
	"bytes"
	"encoding/json"
	"fmt"

	"ensemble/internal/common"
)

// Parse translates a full provider response without storing or sending it.
// Returned model identifiers override routing aliases; otherwise caller provenance is retained.
func (e *Engine) Parse(data []byte, from common.Provenance) (common.ResponseEvent, error) {
	result := common.ResponseEvent{From: from, Parts: []common.Part{}}
	if err := targetValid(from); err != nil {
		return result, err
	}
	if !json.Valid(data) {
		return result, fmt.Errorf("decode %s response: invalid JSON", from.Vendor)
	}
	if err := measuredUsage(data, from.Vendor); err != nil {
		return result, err
	}
	switch from.Vendor {
	case common.Anthropic:
		return parseAnthropic(data, from)
	case common.OpenAI:
		return parseOpenAI(data, from)
	case common.Gemini:
		return parseGemini(data, from)
	}
	return result, fmt.Errorf("unsupported vendor %q", from.Vendor)
}

func parseAnthropic(data []byte, from common.Provenance) (common.ResponseEvent, error) {
	var wire struct {
		// Model is the provider routing or returned model identifier, as selected at this boundary.
		Model string `json:"model"`
		// Content carries ordered provider-native message content.
		Content []json.RawMessage `json:"content"`
		Usage   struct {
			// Input is measured input reported by this provider.
			Input int `json:"input_tokens"`
			// Output is measured output reported by this provider.
			Output int `json:"output_tokens"`
			// Write counts provider-reported newly cached input.
			Write int `json:"cache_creation_input_tokens"`
			// Read counts provider-reported reused cached input.
			Read int `json:"cache_read_input_tokens"`
		} `json:"usage"`
	}
	result := common.ResponseEvent{From: from, Parts: []common.Part{}}
	if err := json.Unmarshal(data, &wire); err != nil {
		return result, fmt.Errorf("decode anthropic response: %w", err)
	}
	if wire.Model != "" {
		from.Model = wire.Model
		result.From = from
	}
	for _, raw := range wire.Content {
		var block struct {
			// Type is the provider-native block or call discriminator.
			Type string `json:"type"`
			// Text preserves text, including a deliberately empty string.
			Text string `json:"text"`
			// ID preserves a provider call ID when present.
			ID string `json:"id"`
			// Name names the requested function without registering or executing it.
			Name string `json:"name"`
			// Input carries the JSON tool argument value.
			Input json.RawMessage `json:"input"`
		}
		if err := json.Unmarshal(raw, &block); err != nil {
			return result, err
		}
		switch block.Type {
		case "text":
			result.Parts = append(result.Parts, common.Part{Type: "text", Text: block.Text})
		case "tool_use":
			result.Parts = append(result.Parts, common.Part{Type: "tool_call", CallID: block.ID, Name: block.Name, Args: block.Input, From: from, Opaque: raw})
		default:
			result.Parts = append(result.Parts, common.Part{Type: "opaque", From: from, Data: raw})
		}
	}
	result.Usage = common.Usage{Input: wire.Usage.Input, Output: wire.Usage.Output, CacheWrite: wire.Usage.Write, CacheRead: wire.Usage.Read}
	if len(result.Parts) == 0 {
		return result, fmt.Errorf("anthropic response contained no content")
	}
	return result, nil
}

func parseOpenAI(data []byte, from common.Provenance) (common.ResponseEvent, error) {
	var wire struct {
		// Model is the provider routing or returned model identifier, as selected at this boundary.
		Model   string `json:"model"`
		Choices []struct {
			Message struct {
				// Content carries ordered provider-native message content.
				Content json.RawMessage `json:"content"`
				// ToolCalls contains assistant-authored function requests.
				ToolCalls []json.RawMessage `json:"tool_calls"`
				// Reasoning preserves provider reasoning bytes only for exact-provenance replay.
				Reasoning json.RawMessage `json:"reasoning_content"`
				// ReasoningAlt preserves the alternate provider reasoning field.
				ReasoningAlt json.RawMessage `json:"reasoning"`
			} `json:"message"`
		} `json:"choices"`
		Usage struct {
			// Input is total prompt input before disjoint cache subtraction.
			Input int `json:"prompt_tokens"`
			// Output already includes reasoning tokens on ChatCompletions.
			Output  int `json:"completion_tokens"`
			Details struct {
				// Cached is the reused portion of prompt input.
				Cached int `json:"cached_tokens"`
				// Written is the optionally reported newly cached portion of prompt input.
				Written int `json:"cache_write_tokens"`
			} `json:"prompt_tokens_details"`
		} `json:"usage"`
	}
	result := common.ResponseEvent{From: from, Parts: []common.Part{}}
	if err := json.Unmarshal(data, &wire); err != nil {
		return result, fmt.Errorf("decode openai response: %w", err)
	}
	if wire.Model != "" {
		from.Model = wire.Model
		result.From = from
	}
	// Attempt cost is measured even when no usable answer was returned.
	// Completion tokens already include reasoning tokens on this surface.
	result.Usage = common.Usage{Input: max(0, wire.Usage.Input-wire.Usage.Details.Cached-wire.Usage.Details.Written), CacheWrite: wire.Usage.Details.Written, CacheRead: wire.Usage.Details.Cached, Output: wire.Usage.Output}
	if len(wire.Choices) == 0 {
		return result, fmt.Errorf("openai response contained no choices")
	}
	message := wire.Choices[0].Message
	if len(message.Content) > 0 && !bytes.Equal(message.Content, []byte("null")) {
		var text string
		if err := json.Unmarshal(message.Content, &text); err != nil {
			return result, fmt.Errorf("openai content is not text")
		}
		result.Parts = append(result.Parts, common.Part{Type: "text", Text: text})
	}
	for _, raw := range message.ToolCalls {
		var call struct {
			// ID preserves a provider call ID when present.
			ID       string `json:"id"`
			Function struct {
				// Name names the requested function without registering or executing it.
				Name string `json:"name"`
				// Arguments is the ChatCompletions JSON-encoded argument string.
				Arguments string `json:"arguments"`
			} `json:"function"`
		}
		if err := json.Unmarshal(raw, &call); err != nil {
			return result, err
		}
		if !json.Valid([]byte(call.Function.Arguments)) {
			return result, fmt.Errorf("openai tool arguments are not JSON")
		}
		result.Parts = append(result.Parts, common.Part{Type: "tool_call", CallID: call.ID, Name: call.Function.Name, Args: json.RawMessage(call.Function.Arguments), From: from, Opaque: raw})
	}
	for _, item := range []struct {
		key string
		raw json.RawMessage
	}{{"reasoning_content", message.Reasoning}, {"reasoning", message.ReasoningAlt}} {
		key, raw := item.key, item.raw
		if len(raw) == 0 || bytes.Equal(raw, []byte("null")) {
			continue
		}
		opaque, _ := json.Marshal(map[string]json.RawMessage{key: raw})
		result.Parts = append(result.Parts, common.Part{Type: "opaque", From: from, Data: opaque})
	}
	if len(result.Parts) == 0 {
		return result, fmt.Errorf("openai response contained no content")
	}
	return result, nil
}

func parseGemini(data []byte, from common.Provenance) (common.ResponseEvent, error) {
	var wire struct {
		// Model is Gemini's returned model identity when supplied.
		Model      string `json:"modelVersion"`
		Candidates []struct {
			Content struct {
				// Parts preserves provider part ordering.
				Parts []json.RawMessage `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
		Usage struct {
			// Input is prompt input before cached-content subtraction.
			Input int `json:"promptTokenCount"`
			// Cached is reused cached prompt input.
			Cached int `json:"cachedContentTokenCount"`
			// Output is visible generated output.
			Output int `json:"candidatesTokenCount"`
			// Thoughts is separately reported generated thinking.
			Thoughts int `json:"thoughtsTokenCount"`
		} `json:"usageMetadata"`
	}
	result := common.ResponseEvent{From: from, Parts: []common.Part{}}
	if err := json.Unmarshal(data, &wire); err != nil {
		return result, fmt.Errorf("decode gemini response: %w", err)
	}
	if wire.Model != "" {
		from.Model = wire.Model
		result.From = from
	}
	// Attempt cost is measured even when no usable answer was returned.
	result.Usage = common.Usage{Input: max(0, wire.Usage.Input-wire.Usage.Cached), CacheRead: wire.Usage.Cached, Output: wire.Usage.Output + wire.Usage.Thoughts}
	if len(wire.Candidates) == 0 {
		return result, fmt.Errorf("gemini response contained no candidates")
	}
	for _, raw := range wire.Candidates[0].Content.Parts {
		var block struct {
			// Text preserves text, including a deliberately empty string.
			Text *string `json:"text"`
			// Thought marks a model-bound reasoning part.
			Thought bool `json:"thought"`
			// Signature is opaque, exact-model replay material attached to its original part.
			Signature json.RawMessage `json:"thoughtSignature"`
			Function  *struct {
				// Name names the requested function without registering or executing it.
				Name string `json:"name"`
				// Args carries Gemini function arguments as a JSON value.
				Args json.RawMessage `json:"args"`
				// ID preserves a provider call ID when present.
				ID string `json:"id"`
			} `json:"functionCall"`
		}
		if err := json.Unmarshal(raw, &block); err != nil {
			return result, err
		}
		if block.Function != nil {
			id := block.Function.ID
			args := block.Function.Args
			if len(args) == 0 {
				args = json.RawMessage(`{}`)
			}
			result.Parts = append(result.Parts, common.Part{Type: "tool_call", CallID: id, Name: block.Function.Name, Args: args, From: from, Opaque: block.Signature})
		} else if block.Thought || block.Text == nil {
			result.Parts = append(result.Parts, common.Part{Type: "opaque", From: from, Data: raw})
		} else {
			// A signature is replay metadata, not evidence that text is hidden.
			part := common.Part{Type: "text", Text: *block.Text}
			if len(block.Signature) > 0 {
				part.From, part.Opaque = from, block.Signature
			}
			result.Parts = append(result.Parts, part)
		}
	}
	if len(result.Parts) == 0 {
		return result, fmt.Errorf("gemini response contained no content")
	}
	return result, nil
}

// measuredUsage distinguishes measured zero from absent accounting. Provider
// totals are accepted as facts only when the required counters are present.
func measuredUsage(data []byte, vendor common.Vendor) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	key := "usage"
	required := []string{"input_tokens", "output_tokens"}
	if vendor == common.OpenAI {
		required = []string{"prompt_tokens", "completion_tokens"}
	}
	if vendor == common.Gemini {
		key = "usageMetadata"
		required = []string{"promptTokenCount", "candidatesTokenCount"}
	}
	var usage map[string]json.RawMessage
	if err := json.Unmarshal(fields[key], &usage); err != nil {
		return fmt.Errorf("%s response missing measured usage", vendor)
	}
	for _, name := range required {
		var count *int
		if err := json.Unmarshal(usage[name], &count); err != nil || count == nil || *count < 0 {
			return fmt.Errorf("%s response missing valid %s", vendor, name)
		}
	}
	return nil
}
