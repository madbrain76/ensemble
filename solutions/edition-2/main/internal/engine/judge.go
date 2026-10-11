package engine

import (
	"context"
	"encoding/json"
	"strings"

	"ensemble/internal/common"
)

// Judge is one tool-free, history-free provider exchange. Recall supplies ALL
// task/format instructions; neither the primary skill nor compressor instructions
// belong here. The existing transport still owns credentials, logs and usage.
func (e *engine) Judge(ctx context.Context, prompt string) (string, error) {
	cfg := e.parent.Config()
	if cfg.Recall.JudgeModel != "" {
		cfg.Model = cfg.Recall.JudgeModel
	}
	cfg.DisableStreaming = true
	var request any
	switch cfg.Vendor {
	case common.Anthropic:
		request = map[string]any{"model": cfg.Model, "max_tokens": 512, "messages": []any{map[string]any{"role": "user", "content": prompt}}}
	case common.OpenAI:
		request = map[string]any{"model": cfg.Model, "messages": []any{map[string]any{"role": "user", "content": prompt}}}
	case common.Gemini:
		request = map[string]any{"contents": []any{map[string]any{"role": "user", "parts": []any{map[string]any{"text": prompt}}}}}
	}
	body, err := json.Marshal(request)
	if err != nil {
		return "", err
	}
	events, err := e.send(ctx, body, common.StreamCallbacks{}, cfg)
	if err != nil {
		return "", err
	}
	var out strings.Builder
	for _, event := range events {
		if event.Response != nil {
			out.WriteString(text(event.Response.Parts))
		}
	}
	return out.String(), nil
}
