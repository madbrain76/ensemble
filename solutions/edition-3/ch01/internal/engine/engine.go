// Package engine owns raw HTTP exchange and measured token accounting.
package engine

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"ensemble/internal/common"
)

const systemPrompt = "You are a helpful assistant built from raw HTTP calls in Chapter 1 of The Art of Building AI Coding Agents. Answer briefly."

type engine struct {
	parent common.Agent
	client *http.Client
	usage  common.Usage
}

// New retains Agent as parent, not a copy of its configuration or logger.
func New(parent common.Agent) (common.Engine, error) {
	if parent == nil {
		return nil, errors.New("Engine requires an Agent parent")
	}
	// Default transport retains connection/TLS protections. No Client.Timeout:
	// model work may take over a minute before headers or while reading a body.
	// Refuse redirects rather than forwarding the API key to another endpoint.
	client := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	}}
	return &engine{parent: parent, client: client}, nil
}

// Parent exposes the Agent that owns this Engine.
func (e *engine) Parent() common.Agent { return e.parent }

// Usage returns a value copy of Engine's authoritative measured totals.
func (e *engine) Usage() common.Usage { return e.usage }

// failure follows real owners for diagnostics. Only locally chosen messages
// reach this helper: remote bodies, URLs and transport errors can contain secrets.
func (e *engine) failure(message string) error {
	e.parent.Parent().Logger().Print(message)
	return errors.New(message)
}

// Send sends the full pending conversation once. The caller's context controls
// both header and body waiting; there is no independent model-work deadline.
func (e *engine) Send(ctx context.Context, messages []common.Message) (string, error) {
	config := e.parent.Config()
	body, err := json.Marshal(struct {
		Model     string           `json:"model"`
		MaxTokens int              `json:"max_tokens"`
		System    string           `json:"system"`
		Messages  []common.Message `json:"messages"`
	}{config.Model, 1024, systemPrompt, messages})
	if err != nil {
		return "", e.failure("encoding request failed")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, config.BaseURL+"/v1/messages", bytes.NewReader(body))
	if err != nil {
		return "", e.failure("creating request failed")
	}
	req.Header.Set("x-api-key", config.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("content-type", "application/json")
	resp, err := e.client.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}
		return "", e.failure("HTTP request failed")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", e.failure(fmt.Sprintf("HTTP status %d", resp.StatusCode))
	}
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}
		return "", e.failure("reading response failed")
	}
	var parsed struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
		Usage *struct {
			Input  *int64 `json:"input_tokens"`
			Output *int64 `json:"output_tokens"`
		} `json:"usage"`
	}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return "", e.failure("decoding response failed")
	}
	// Pointers distinguish a measured zero from absent accounting.
	if parsed.Usage == nil || parsed.Usage.Input == nil || parsed.Usage.Output == nil {
		return "", e.failure("response is missing usage")
	}
	if *parsed.Usage.Input < 0 || *parsed.Usage.Output < 0 {
		return "", e.failure("response has negative usage")
	}
	// Even an unusable answer cost tokens. Count them before validating text.
	e.usage.Input += *parsed.Usage.Input
	e.usage.Output += *parsed.Usage.Output
	var text strings.Builder
	for _, block := range parsed.Content {
		if block.Type == "text" {
			text.WriteString(block.Text)
		}
	}
	if strings.TrimSpace(text.String()) == "" {
		return "", e.failure("response has no usable text")
	}
	return text.String(), nil
}
