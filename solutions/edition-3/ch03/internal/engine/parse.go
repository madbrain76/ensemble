package engine

import (
	"encoding/json"
	"ensemble/internal/common"
	"fmt"
	"strings"
)

func parse(status int, b []byte, c common.Config) ([]common.Event, error) {
	fail := func(message string) ([]common.Event, error) {
		return []common.Event{errorEvent(message)}, fmt.Errorf("%s", message)
	}
	if status != 200 {
		return fail(fmt.Sprintf("provider HTTP status %d", status))
	}
	// The route is not a producer identity. A missing returned model remains
	// unknown even if the caller supplied a resolved target for the next request.
	from := common.Provenance{Vendor: c.Vendor, Surface: c.Surface}
	var envelope map[string]json.RawMessage
	if json.Unmarshal(b, &envelope) != nil || envelope == nil {
		return fail("decoding response failed")
	}
	if c.Vendor == common.Gemini {
		_ = json.Unmarshal(envelope["modelVersion"], &from.Model)
	} else {
		_ = json.Unmarshal(envelope["model"], &from.Model)
	}
	usage, missing, usageErr := parseUsage(c.Vendor, envelope)
	var parts common.PartList
	var contentErr error
	switch c.Vendor {
	case common.Anthropic:
		parts, contentErr = parseAnthropic(b, from)
	case common.OpenAI:
		parts, contentErr = parseOpenAI(b, from)
	case common.Gemini:
		parts, contentErr = parseGemini(b, from)
	default:
		return fail("unsupported vendor")
	}
	// Even an unusable or partially parsed answer is an observed response. Its
	// measured usage survives; missing fields are explicitly not measured zero.
	events := []common.Event{{Type: common.ResponseStarted, Response: &common.ResponseData{From: from}}, {Type: common.ResponseEnded, Response: &common.ResponseData{Parts: parts, Usage: usage, UsageMissing: missing, From: from}}}
	usable := len(parts) > 0
	if c.Vendor == common.Anthropic {
		usable = false
		for _, p := range parts {
			if p.Type != "text" || strings.TrimSpace(p.Text) != "" {
				usable = true
			}
		}
	}
	err := usageErr
	if err == nil {
		err = contentErr
	}
	if err == nil && !usable {
		err = fmt.Errorf("response contained no usable text")
	}
	if err != nil {
		events = append(events, errorEvent(err.Error()))
		return events, err
	}
	return events, nil
}

// Usage is decoded before content and field-by-field: one missing or malformed
// count cannot discard another measured charge. Optional cache/thought fields
// have documented zero defaults; the vendors' base input/output counts do not.
func parseUsage(v common.Vendor, envelope map[string]json.RawMessage) (common.Usage, []string, error) {
	key := "usage"
	if v == common.Gemini {
		key = "usageMetadata"
	}
	var fields map[string]json.RawMessage
	_ = json.Unmarshal(envelope[key], &fields)
	var firstErr error
	missing := []string{}
	mark := func(bucket, message string) {
		found := false
		for _, s := range missing {
			if s == bucket {
				found = true
			}
		}
		if !found {
			missing = append(missing, bucket)
		}
		if firstErr == nil {
			firstErr = fmt.Errorf("%s", message)
		}
	}
	read := func(m map[string]json.RawMessage, key, bucket string, required bool) (int64, bool) {
		b, exists := m[key]
		if !exists {
			if required {
				mark(bucket, "response missing usage")
			}
			return 0, !required
		}
		var n *int64
		if json.Unmarshal(b, &n) != nil || n == nil {
			mark(bucket, "response contained invalid usage")
			return 0, false
		}
		if *n < 0 {
			mark(bucket, "response contained negative usage")
			return 0, false
		}
		return *n, true
	}
	nested := func(key string) (map[string]json.RawMessage, bool) {
		var m map[string]json.RawMessage
		if b, ok := fields[key]; ok {
			if json.Unmarshal(b, &m) != nil || m == nil {
				return nil, false
			}
		}
		return m, true // Omitted optional details retain their zero defaults.
	}
	var u common.Usage
	switch v {
	case common.Anthropic:
		u.Input, _ = read(fields, "input_tokens", "input", true)
		u.Output, _ = read(fields, "output_tokens", "output", true)
		u.CacheRead, _ = read(fields, "cache_read_input_tokens", "cache_read", false)
		u.CacheWrite, _ = read(fields, "cache_creation_input_tokens", "cache_write", false)
	case common.OpenAI:
		total, ok := read(fields, "prompt_tokens", "input", true)
		u.Output, _ = read(fields, "completion_tokens", "output", true)
		details, detailsOK := nested("prompt_tokens_details")
		var readOK, writeOK bool
		if detailsOK {
			u.CacheRead, readOK = read(details, "cached_tokens", "cache_read", false)
			u.CacheWrite, writeOK = read(details, "cache_write_tokens", "cache_write", false)
		} else {
			mark("cache_read", "response contained invalid usage")
			mark("cache_write", "response contained invalid usage")
		}
		// Fresh input is derived, not an independent measured count. Unknown
		// cache operands must not turn their numeric placeholders into charges.
		if !readOK || !writeOK {
			mark("input", "response contained incomplete input usage")
		} else if ok {
			if total < u.CacheRead+u.CacheWrite {
				mark("input", "response contained inconsistent usage")
			} else {
				u.Input = total - u.CacheRead - u.CacheWrite
			}
		}
	case common.Gemini:
		total, ok := read(fields, "promptTokenCount", "input", true)
		var readOK bool
		u.CacheRead, readOK = read(fields, "cachedContentTokenCount", "cache_read", false)
		visible, _ := read(fields, "candidatesTokenCount", "output", true)
		thoughts, _ := read(fields, "thoughtsTokenCount", "output", false)
		u.Output = visible + thoughts
		if !readOK {
			mark("input", "response contained incomplete input usage")
		} else if ok {
			if total < u.CacheRead {
				mark("input", "response contained inconsistent usage")
			} else {
				u.Input = total - u.CacheRead
			}
		}
	}
	return u, missing, firstErr
}
func parseAnthropic(b []byte, from common.Provenance) (common.PartList, error) {
	var response struct {
		Content []json.RawMessage `json:"content"`
	}
	if json.Unmarshal(b, &response) != nil {
		return nil, fmt.Errorf("decoding response failed")
	}
	parts := common.PartList{}
	for _, block := range response.Content {
		var p struct {
			Type  string          `json:"type"`
			Text  *string         `json:"text"`
			ID    string          `json:"id"`
			Name  string          `json:"name"`
			Input json.RawMessage `json:"input"`
		}
		if json.Unmarshal(block, &p) != nil {
			return parts, fmt.Errorf("invalid provider content")
		}
		switch p.Type {
		case "text":
			if p.Text == nil {
				return parts, fmt.Errorf("missing text")
			}
			parts = append(parts, common.Text(*p.Text))
		case "tool_use":
			args, err := canonical(p.Input)
			if err != nil {
				return parts, fmt.Errorf("invalid tool arguments")
			}
			parts = append(parts, common.Part{Type: "tool_call", CallID: p.ID, Name: p.Name, Args: args, From: &from})
		case "thinking", "redacted_thinking":
			parts = append(parts, common.Part{Type: "opaque", From: &from, Data: block})
		}
	}
	return parts, nil
}
func parseOpenAI(b []byte, from common.Provenance) (common.PartList, error) {
	var response struct {
		Choices []struct {
			Message struct {
				Content   *string `json:"content"`
				ToolCalls []oCall `json:"tool_calls"`
				Refusal   *string `json:"refusal"`
			} `json:"message"`
		} `json:"choices"`
	}
	if json.Unmarshal(b, &response) != nil {
		return nil, fmt.Errorf("decoding response failed")
	}
	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("missing choice")
	}
	message := response.Choices[0].Message
	parts := common.PartList{}
	if message.Content != nil {
		parts = append(parts, common.Text(*message.Content))
	}
	if message.Refusal != nil && *message.Refusal != "" {
		parts = append(parts, common.Text(*message.Refusal))
	}
	for _, call := range message.ToolCalls {
		args, err := canonical(json.RawMessage(call.Function.Arguments))
		if err != nil {
			return parts, fmt.Errorf("invalid tool arguments")
		}
		parts = append(parts, common.Part{Type: "tool_call", CallID: call.ID, Name: call.Function.Name, Args: args, From: &from})
	}
	return parts, nil
}
func parseGemini(b []byte, from common.Provenance) (common.PartList, error) {
	var response struct {
		Candidates []struct {
			Content struct {
				Parts []json.RawMessage `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}
	if json.Unmarshal(b, &response) != nil {
		return nil, fmt.Errorf("decoding response failed")
	}
	if len(response.Candidates) == 0 {
		return nil, fmt.Errorf("missing candidate")
	}
	parts := common.PartList{}
	for _, block := range response.Candidates[0].Content.Parts {
		var p struct {
			Text         *string `json:"text"`
			Thought      bool    `json:"thought"`
			FunctionCall *struct {
				ID   string          `json:"id"`
				Name string          `json:"name"`
				Args json.RawMessage `json:"args"`
			} `json:"functionCall"`
			Signature json.RawMessage `json:"thoughtSignature"`
		}
		if json.Unmarshal(block, &p) != nil {
			return parts, fmt.Errorf("invalid provider content")
		}
		if p.FunctionCall != nil {
			args, err := canonical(p.FunctionCall.Args)
			if err != nil {
				return parts, fmt.Errorf("invalid tool arguments")
			}
			parts = append(parts, common.Part{Type: "tool_call", CallID: p.FunctionCall.ID, Name: p.FunctionCall.Name, Args: args, From: &from, Opaque: p.Signature})
		} else if p.Thought {
			parts = append(parts, common.Part{Type: "opaque", From: &from, Data: block})
		} else if p.Text != nil {
			part := common.Text(*p.Text)
			if len(p.Signature) > 0 {
				part.From = &from
				part.Opaque = p.Signature
			}
			parts = append(parts, part)
		} else {
			return parts, fmt.Errorf("unknown Gemini part")
		}
	}
	return parts, nil
}
