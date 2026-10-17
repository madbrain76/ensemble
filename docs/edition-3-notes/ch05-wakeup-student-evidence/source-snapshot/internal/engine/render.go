package engine

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"ensemble/internal/common"
)

// Target preserves the routing name; exact replay identity is a separate setting.
func Target(config common.Config) common.Provenance {
	vendor := config.Vendor
	if vendor == 0 {
		vendor = common.Anthropic
	}
	surface := common.Messages
	switch vendor {
	case common.OpenAI:
		surface = common.ChatCompletions
	case common.Gemini:
		surface = common.GenerateContent
	}
	return common.Provenance{Vendor: vendor, Model: config.Model, Surface: surface}
}
func targetValid(to common.Provenance) error {
	expected := Target(common.Config{Vendor: to.Vendor, Model: to.Model})
	if to.Model == "" || to.Vendor < common.Anthropic || to.Vendor > common.Gemini || expected != to {
		return fmt.Errorf("unsupported target vendor/model/surface")
	}
	return nil
}

// Render preserves Context and History. It loads supported local images into a
// copied request projection only; URI fetching is deliberately unsupported.
func (e *Engine) Render(c common.Context, to common.Provenance) (json.RawMessage, error) {
	c.Dialogue = append([]common.Turn(nil), c.Dialogue...)
	for i := range c.Dialogue {
		c.Dialogue[i].Parts = cloneParts(c.Dialogue[i].Parts)
	}
	c.Ephemera = cloneParts(c.Ephemera)
	if err := targetValid(to); err != nil {
		return nil, err
	}
	replay := to
	if resolved := e.parent.Config().ResolvedModel; resolved != "" {
		replay.Model = resolved
	}
	// Copy before synthesizing correlation ids. They belong to this request, not
	// the facts. Seq prevents missing ids in distinct turns from colliding.
	data, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	for i := range c.Dialogue {
		if err := e.prepareMedia(c.Dialogue[i].Parts, to); err != nil {
			return nil, err
		}
	}
	if err := e.prepareMedia(c.Ephemera, to); err != nil {
		return nil, err
	}
	lastID := ""
	for i := range c.Dialogue {
		for j := range c.Dialogue[i].Parts {
			p := &c.Dialogue[i].Parts[j]
			if p.Type == "tool_call" {
				if p.CallID == "" && to.Vendor != common.Gemini {
					p.CallID = fmt.Sprintf("call-%d-%d", c.Dialogue[i].Seq, j)
				}
				lastID = p.CallID
			}
			if p.Type == "tool_result" && p.CallID == "" {
				p.CallID = lastID
			}
		}
	}
	switch to.Vendor {
	case common.Anthropic:
		return json.Marshal(e.renderAnthropic(c, to, replay))
	case common.OpenAI:
		return json.Marshal(e.renderOpenAI(c, to, replay))
	case common.Gemini:
		return json.Marshal(e.renderGemini(c, replay))
	}
	return nil, fmt.Errorf("unknown vendor")
}
func (e *Engine) prepareMedia(parts []common.Part, to common.Provenance) error {
	for i := range parts {
		p := &parts[i]
		for _, raw := range []json.RawMessage{p.Args, p.Opaque, p.Data} {
			if len(raw) > 0 && !json.Valid(raw) {
				return fmt.Errorf("part contains invalid JSON")
			}
		}
		if p.Type == "tool_call" && (p.Name == "" || !json.Valid(p.Args)) {
			return fmt.Errorf("tool call requires name and JSON arguments")
		}
		if p.Type == "opaque" && !json.Valid(p.Data) {
			return fmt.Errorf("opaque part requires a JSON value")
		}
		switch p.Type {
		case "blob":
			kind, flag := "unknown", common.Media(0)
			switch {
			case strings.HasPrefix(p.MIME, "image/"):
				kind, flag = "image", common.MediaImage
			case strings.HasPrefix(p.MIME, "audio/"):
				kind, flag = "audio", common.MediaAudio
			case strings.HasPrefix(p.MIME, "video/"):
				kind, flag = "video", common.MediaVideo
			case p.MIME == "application/pdf":
				kind, flag = "document", common.MediaDocument
			}
			if ModelCapabilities(to.Vendor, to.Model).Media&flag == 0 {
				return fmt.Errorf("model %q does not support %s media %q", to.Model, kind, p.MIME)
			}
			if p.Ref == nil || p.Ref.Kind != common.RefPath || p.Ref.Locator == "" {
				return fmt.Errorf("model %q media %q requires a local path reference", to.Model, p.MIME)
			}
			path := p.Ref.Locator
			if !filepath.IsAbs(path) {
				path = filepath.Join(e.parent.Config().Workspace, path)
			}
			info, err := os.Stat(path)
			if err != nil {
				return fmt.Errorf("model %q media %q: %w", to.Model, p.MIME, err)
			}
			if !info.Mode().IsRegular() || info.Size() > 20<<20 {
				return fmt.Errorf("media must be a regular file at most 20 MiB")
			}
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			detected := http.DetectContentType(data)
			if kind == "image" && detected != p.MIME {
				return fmt.Errorf("model %q media MIME %q does not match %q", to.Model, p.MIME, detected)
			}
			p.Data = encoded(base64.StdEncoding.EncodeToString(data))
		case "tool_result":
			if err := e.prepareMedia(p.Parts, to); err != nil {
				return err
			}
		case "text", "opaque", "tool_call", "redacted":
		default:
			return fmt.Errorf("unsupported part type %q", p.Type)
		}
	}
	return nil
}
func text(p common.Part) string {
	if p.Type == "redacted" {
		return p.Stub
	}
	return p.Text
}
func textParts(parts []common.Part) string {
	var b strings.Builder
	for _, p := range parts {
		if p.Type == "text" || p.Type == "redacted" {
			b.WriteString(text(p))
		}
	}
	return b.String()
}

// Each wire shape has source-defined field order. RawMessage is reserved for
// opaque provider material and JSON values whose bytes must not be interpreted.
type anthropicBlock struct {
	// Type is the provider-native block or call discriminator.
	Type string `json:"type"`
	// Text preserves text, including a deliberately empty string.
	Text string `json:"text,omitempty"`
	// ID preserves a provider call ID when present.
	ID string `json:"id,omitempty"`
	// Name names the requested function without registering or executing it.
	Name string `json:"name,omitempty"`
	// Input carries the JSON tool argument value.
	Input json.RawMessage `json:"input,omitempty"`
	// ToolUseID pairs an Anthropic result with its preceding call.
	ToolUseID string `json:"tool_use_id,omitempty"`
	// Content carries ordered provider-native message content.
	Content []json.RawMessage `json:"content,omitempty"`
	// IsError preserves a failed tool return without executing or retrying it.
	IsError bool `json:"is_error,omitempty"`
}
type anthropicMessage struct {
	// Role encodes provider-specific authorship, not domain ownership.
	Role string `json:"role"`
	// Content carries ordered provider-native message content.
	Content []json.RawMessage `json:"content"`
}
type anthropicTool struct {
	// Name is the registered tool name.
	Name string `json:"name"`
	// Description explains the operation.
	Description string `json:"description"`
	// InputSchema is Anthropic's argument schema.
	InputSchema json.RawMessage `json:"input_schema"`
}
type openAITool struct {
	// Type is the ChatCompletions function discriminator.
	Type string `json:"type"`
	// Function carries the shared function definition.
	Function common.ToolDefinition `json:"function"`
}

// Gemini parameters is its Schema subset, not arbitrary JSON Schema. In
// particular the live API rejects additionalProperties here. Local argument
// decoding remains strict regardless of vendor declaration representation.
type geminiSchema struct {
	// Type is the top-level object discriminator.
	Type string `json:"type"`
	// Properties retains the registry's supported per-argument schemas.
	Properties json.RawMessage `json:"properties"`
	// Required lists mandatory argument names.
	Required []string `json:"required,omitempty"`
}
type geminiDeclaration struct {
	// Name identifies the Agent-owned operation.
	Name string `json:"name"`
	// Description explains the operation to the model.
	Description string `json:"description"`
	// Parameters is the Gemini-supported schema representation.
	Parameters geminiSchema `json:"parameters"`
}
type geminiTools struct {
	// FunctionDeclarations lists Gemini functions.
	FunctionDeclarations []geminiDeclaration `json:"functionDeclarations"`
}

type anthropicRequest struct {
	// Tools advertises only this Agent's registered operations.
	Tools []anthropicTool `json:"tools,omitempty"`
	// Model is the provider routing or returned model identifier, as selected at this boundary.
	Model string `json:"model"`
	// MaxTokens bounds generated tokens rather than generation time.
	MaxTokens int `json:"max_tokens"`
	// System contains stable Config instructions plus any pending one-request ephemera.
	System string `json:"system"`
	// Messages contains the ordered provider request messages.
	Messages []anthropicMessage `json:"messages"`
}

func encoded(v any) json.RawMessage { data, _ := json.Marshal(v); return data }
func anthropicParts(parts []common.Part, to common.Provenance) []json.RawMessage {
	out := []json.RawMessage{}
	for _, p := range parts {
		switch p.Type {
		case "blob":
			var data string
			_ = json.Unmarshal(p.Data, &data)
			kind := "image"
			if p.MIME == "application/pdf" {
				kind = "document"
			}
			out = append(out, encoded(map[string]any{"type": kind, "source": map[string]any{"type": "base64", "media_type": p.MIME, "data": data}}))
		case "opaque":
			if p.From == to {
				out = append(out, p.Data)
			}
		case "tool_call":
			if p.From == to && len(p.Opaque) > 0 {
				var raw anthropicBlock
				_ = json.Unmarshal(p.Opaque, &raw)
				if raw.Type == "tool_use" && raw.ID != "" {
					out = append(out, p.Opaque)
					continue
				}
			}
			out = append(out, encoded(anthropicBlock{Type: "tool_use", ID: p.CallID, Name: p.Name, Input: p.Args}))
		case "tool_result":
			out = append(out, encoded(anthropicBlock{Type: "tool_result", ToolUseID: p.CallID, Content: anthropicParts(p.Parts, to), IsError: p.IsError}))
		default:
			out = append(out, encoded(struct {
				// Type is the provider-native block or call discriminator.
				Type string `json:"type"`
				// Text preserves text, including a deliberately empty string.
				Text string `json:"text"`
			}{"text", text(p)}))
		}
	}
	return out
}
func (e *Engine) renderAnthropic(c common.Context, to, replay common.Provenance) anthropicRequest {
	result := anthropicRequest{Model: to.Model, MaxTokens: e.parent.Config().MaxTokens, System: e.parent.Config().System, Messages: []anthropicMessage{}}
	for _, tool := range e.parent.Tools() {
		d := tool.Definition()
		result.Tools = append(result.Tools, anthropicTool{d.Name, d.Description, d.Parameters})
	}
	appendTurn := func(role string, parts []common.Part) {
		blocks := anthropicParts(parts, replay)
		if len(blocks) == 0 {
			return
		}
		n := len(result.Messages)
		if n > 0 && result.Messages[n-1].Role == role {
			result.Messages[n-1].Content = append(result.Messages[n-1].Content, blocks...)
		} else {
			result.Messages = append(result.Messages, anthropicMessage{Role: role, Content: blocks})
		}
	}
	for _, turn := range c.Dialogue {
		role := "user"
		if turn.Actor == "agent" {
			role = "assistant"
		}
		appendTurn(role, turn.Parts)
	}
	appendTurn("user", c.Ephemera)
	// Group only after the complete merge, retaining arrival order within both
	// groups (result one, human text, result two -> both results, then text).
	for i := range result.Messages {
		if result.Messages[i].Role != "user" {
			continue
		}
		results, rest := []json.RawMessage{}, []json.RawMessage{}
		for _, raw := range result.Messages[i].Content {
			var tag struct {
				// Type is the provider-native block or call discriminator.
				Type string `json:"type"`
			}
			_ = json.Unmarshal(raw, &tag)
			if tag.Type == "tool_result" {
				results = append(results, raw)
			} else {
				rest = append(rest, raw)
			}
		}
		result.Messages[i].Content = append(results, rest...)
	}
	return result
}

type openAIFunction struct {
	// Name names the requested function without registering or executing it.
	Name string `json:"name"`
	// Arguments is the ChatCompletions JSON-encoded argument string.
	Arguments string `json:"arguments"`
}
type openAICall struct {
	// ID preserves a provider call ID when present.
	ID string `json:"id"`
	// Type is the provider-native block or call discriminator.
	Type string `json:"type"`
	// Function holds ChatCompletions function-call details.
	Function openAIFunction `json:"function"`
}
type openAIMessage struct {
	// Role encodes provider-specific authorship, not domain ownership.
	Role string `json:"role"`
	// Content carries ordered provider-native message content.
	Content    any    `json:"content"`
	ToolCallID string `json:"tool_call_id,omitempty"`
	// ToolCalls contains assistant-authored function requests.
	ToolCalls []json.RawMessage `json:"tool_calls,omitempty"`
	// ReasoningContent preserves provider reasoning bytes only for exact-provenance replay.
	ReasoningContent json.RawMessage `json:"reasoning_content,omitempty"`
	// Reasoning preserves the alternate provider reasoning field.
	Reasoning json.RawMessage `json:"reasoning,omitempty"`
}
type openAIRequest struct {
	// Tools advertises only this Agent's registered operations.
	Tools []openAITool `json:"tools,omitempty"`
	// Model is the provider routing or returned model identifier, as selected at this boundary.
	Model string `json:"model"`
	// MaxTokens bounds generated tokens rather than generation time.
	MaxTokens int `json:"max_tokens"`
	// Messages contains the ordered provider request messages.
	Messages []openAIMessage `json:"messages"`
}

func (e *Engine) renderOpenAI(c common.Context, to, replay common.Provenance) openAIRequest {
	result := openAIRequest{Model: to.Model, MaxTokens: e.parent.Config().MaxTokens, Messages: []openAIMessage{{Role: "system", Content: e.parent.Config().System}}}
	for _, tool := range e.parent.Tools() {
		result.Tools = append(result.Tools, openAITool{Type: "function", Function: tool.Definition()})
	}
	for _, turn := range c.Dialogue {
		role := "user"
		if turn.Actor == "agent" {
			role = "assistant"
		}
		if turn.Actor == "system" {
			role = "system"
		}
		message := openAIMessage{Role: role, Content: ""}
		blocks := []any{}
		image := false
		results := []openAIMessage{}
		hasContent := false
		for _, p := range turn.Parts {
			switch p.Type {
			case "blob":
				var data string
				_ = json.Unmarshal(p.Data, &data)
				blocks = append(blocks, map[string]any{"type": "image_url", "image_url": map[string]string{"url": "data:" + p.MIME + ";base64," + data}})
				image = true
				hasContent = true
			case "opaque":
				if p.From == replay {
					var opaque struct {
						// ReasoningContent preserves provider reasoning bytes only for exact-provenance replay.
						ReasoningContent json.RawMessage `json:"reasoning_content"`
						// Reasoning preserves the alternate provider reasoning field.
						Reasoning json.RawMessage `json:"reasoning"`
					}
					_ = json.Unmarshal(p.Data, &opaque)
					if opaque.ReasoningContent != nil {
						message.ReasoningContent = opaque.ReasoningContent
					}
					if opaque.Reasoning != nil {
						message.Reasoning = opaque.Reasoning
					}
				}
			case "tool_call":
				raw := encoded(openAICall{ID: p.CallID, Type: "function", Function: openAIFunction{Name: p.Name, Arguments: string(p.Args)}})
				if p.From == replay && len(p.Opaque) > 0 {
					var call openAICall
					_ = json.Unmarshal(p.Opaque, &call)
					if call.ID != "" && call.Type == "function" {
						raw = p.Opaque
					}
				}
				message.ToolCalls = append(message.ToolCalls, raw)
			case "tool_result":
				results = append(results, openAIMessage{Role: "tool", ToolCallID: p.CallID, Content: textParts(p.Parts)})
			default:
				message.Content = message.Content.(string) + text(p)
				blocks = append(blocks, map[string]string{"type": "text", "text": text(p)})
				hasContent = true
			}
		}
		if image {
			message.Content = blocks
		}
		if hasContent || len(message.ToolCalls) > 0 || len(results) == 0 {
			result.Messages = append(result.Messages, message)
		}
		result.Messages = append(result.Messages, results...)
	}
	if len(c.Ephemera) > 0 {
		result.Messages = append(result.Messages, openAIMessage{Role: "system", Content: textParts(c.Ephemera)})
	}
	return result
}

type geminiCall struct {
	// Name names the requested function without registering or executing it.
	Name string `json:"name"`
	// Args carries Gemini function arguments as a JSON value.
	Args json.RawMessage `json:"args"`
	// ID preserves a provider call ID when present.
	ID string `json:"id,omitempty"`
}
type geminiResult struct {
	// ID pairs a result with the exact native or captured call.
	ID string `json:"id,omitempty"`
	// Name names the requested function without registering or executing it.
	Name     string `json:"name"`
	Response struct {
		// Content carries ordered provider-native message content.
		Content string `json:"content"`
		// IsError preserves a failed tool return without executing or retrying it.
		IsError bool `json:"is_error,omitempty"`
	} `json:"response"`
}
type geminiPart struct {
	// Text preserves text, including a deliberately empty string.
	Text *string `json:"text,omitempty"`
	// FunctionCall carries Gemini model-authored call details.
	FunctionCall *geminiCall `json:"functionCall,omitempty"`
	// FunctionResponse carries Gemini user-authored result details.
	FunctionResponse *geminiResult `json:"functionResponse,omitempty"`
	// ThoughtSignature is opaque, exact-model replay material attached to its original part.
	ThoughtSignature json.RawMessage `json:"thoughtSignature,omitempty"`
}
type geminiContent struct {
	// Role encodes provider-specific authorship, not domain ownership.
	Role string `json:"role,omitempty"`
	// Parts preserves provider part ordering.
	Parts []json.RawMessage `json:"parts"`
}
type geminiRequest struct {
	// Tools advertises only this Agent's registered operations.
	Tools []geminiTools `json:"tools,omitempty"`
	// SystemInstruction keeps Gemini system instruction outside conversational contents.
	SystemInstruction geminiContent `json:"systemInstruction"`
	// Contents contains Gemini conversational contents.
	Contents         []geminiContent `json:"contents"`
	GenerationConfig struct {
		// MaxOutputTokens bounds generated tokens rather than elapsed time.
		MaxOutputTokens int `json:"maxOutputTokens"`
	} `json:"generationConfig"`
}

func geminiParts(parts []common.Part, to common.Provenance, names map[string]string) []json.RawMessage {
	out := []json.RawMessage{}
	for _, p := range parts {
		switch p.Type {
		case "blob":
			var data string
			_ = json.Unmarshal(p.Data, &data)
			out = append(out, encoded(map[string]any{"inlineData": map[string]string{"mimeType": p.MIME, "data": data}}))
		case "text":
			part := geminiPart{Text: &p.Text}
			// Keep the signature attached to its single original visible part,
			// and only on the exact issuing vendor, surface, and model.
			if p.From == to {
				part.ThoughtSignature = p.Opaque
			}
			out = append(out, encoded(part))
		case "opaque":
			if p.From == to {
				out = append(out, p.Data)
			}
		case "tool_call":
			// Resolve only calls already visited, including sequential absent IDs.
			names[p.CallID] = p.Name
			if p.From == to && len(p.Opaque) > 0 {
				var raw geminiPart
				_ = json.Unmarshal(p.Opaque, &raw)
				if raw.FunctionCall != nil {
					out = append(out, p.Opaque)
					continue
				}
			}
			part := geminiPart{FunctionCall: &geminiCall{Name: p.Name, Args: p.Args}}
			// A synthesized correlation id is needed by other wires, not by Gemini
			// models that omitted one. Native ids survive unchanged.
			part.FunctionCall.ID = p.CallID
			if p.From == to && len(p.Opaque) > 0 {
				var opaque struct {
					// Signature is opaque, exact-model replay material attached to its original part.
					Signature json.RawMessage `json:"thoughtSignature"`
				}
				_ = json.Unmarshal(p.Opaque, &opaque)
				part.ThoughtSignature = opaque.Signature
				if len(part.ThoughtSignature) == 0 {
					var signature string
					if json.Unmarshal(p.Opaque, &signature) == nil {
						part.ThoughtSignature = p.Opaque
					}
				}
			}
			out = append(out, encoded(part))
		case "tool_result":
			name := names[p.CallID]
			if name == "" {
				name = p.CallID
			}
			result := &geminiResult{Name: name, ID: p.CallID}
			result.Response.Content = textParts(p.Parts)
			result.Response.IsError = p.IsError
			out = append(out, encoded(geminiPart{FunctionResponse: result}))
		default:
			value := text(p)
			out = append(out, encoded(geminiPart{Text: &value}))
		}
	}
	return out
}
func (e *Engine) renderGemini(c common.Context, replay common.Provenance) geminiRequest {
	system := e.parent.Config().System
	result := geminiRequest{SystemInstruction: geminiContent{Parts: []json.RawMessage{encoded(geminiPart{Text: &system})}}, Contents: []geminiContent{}}
	definitions := []geminiDeclaration{}
	for _, tool := range e.parent.Tools() {
		definition := tool.Definition()
		var parameters geminiSchema
		// Agent-owned built-ins supply valid schemas; retain the supported fields.
		_ = json.Unmarshal(definition.Parameters, &parameters)
		definitions = append(definitions, geminiDeclaration{Name: definition.Name, Description: definition.Description, Parameters: parameters})
	}
	if len(definitions) > 0 {
		result.Tools = []geminiTools{{FunctionDeclarations: definitions}}
	}
	result.GenerationConfig.MaxOutputTokens = e.parent.Config().MaxTokens
	names := map[string]string{}
	for _, turn := range c.Dialogue {
		role := "user"
		if turn.Actor == "agent" {
			role = "model"
		}
		parts := geminiParts(turn.Parts, replay, names)
		if len(parts) > 0 {
			result.Contents = append(result.Contents, geminiContent{Role: role, Parts: parts})
		}
	}
	if len(c.Ephemera) > 0 {
		result.Contents = append(result.Contents, geminiContent{Role: "user", Parts: geminiParts(c.Ephemera, replay, names)})
	}
	return result
}

// cloneParts detaches borrowed data before this owner mutates or publishes it.
func cloneParts(parts []common.Part) []common.Part {
	if parts == nil {
		return nil
	}
	out := make([]common.Part, len(parts))
	for i, p := range parts {
		p.Args = append([]byte(nil), p.Args...)
		p.Data = append([]byte(nil), p.Data...)
		p.Opaque = append([]byte(nil), p.Opaque...)
		p.Parts = cloneParts(p.Parts)
		if p.Ref != nil {
			r := *p.Ref
			p.Ref = &r
		}
		out[i] = p
	}
	return out
}
