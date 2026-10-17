// Package engine implements full-response provider transport and translation.
package engine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"ensemble/internal/common"
)

// Engine owns the HTTP connection policy and live usage. It never edits History.
type Engine struct {
	parent common.Agent
	client *http.Client
	mu     sync.Mutex
	usage  common.Usage
}

// New requires an owning Agent; network setup is bounded, model work is not.
func New(parent common.Agent) (*Engine, error) {
	if parent == nil {
		return nil, fmt.Errorf("engine requires an agent parent")
	}
	transport := &http.Transport{Proxy: http.ProxyFromEnvironment, DialContext: (&net.Dialer{Timeout: 30 * time.Second, KeepAlive: 30 * time.Second}).DialContext, TLSHandshakeTimeout: 10 * time.Second, ExpectContinueTimeout: time.Second, IdleConnTimeout: 90 * time.Second, MaxIdleConns: 100}
	client := &http.Client{Transport: transport, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	return &Engine{parent: parent, client: client}, nil
}

// Parent exposes its concrete owner's interface, not a configuration callback bag.
func (e *Engine) Parent() common.Agent { return e.parent }

// Target returns the configured exact vendor surface.
func (e *Engine) Target() common.Provenance { return Target(e.parent.Config()) }

// Usage reports measured live replies, including a measured but unusable answer.
func (e *Engine) Usage() common.Usage { e.mu.Lock(); defer e.mu.Unlock(); return e.usage }

// Send performs exactly one full-response request; only caller context can impose
// a generation deadline. Provider bodies and credential-bearing errors are not logged.
func (e *Engine) Send(ctx context.Context, c common.Context) (common.ResponseEvent, error) {
	payload, err := e.Render(c, e.Target())
	if err != nil {
		return common.ResponseEvent{}, err
	}
	return e.SendPrepared(ctx, payload)
}

func (e *Engine) SendPrepared(ctx context.Context, payload json.RawMessage) (common.ResponseEvent, error) {
	config := e.parent.Config()
	target := Target(config)
	if config.APIKey == "" {
		return common.ResponseEvent{}, fmt.Errorf("API key is required")
	}
	if err := targetValid(target); err != nil {
		return common.ResponseEvent{}, err
	}
	base, err := url.Parse(config.BaseURL)
	if err != nil || base == nil || base.Host == "" || (base.Scheme != "http" && base.Scheme != "https") || base.User != nil || base.RawQuery != "" || base.Fragment != "" {
		return common.ResponseEvent{}, fmt.Errorf("invalid provider base URL")
	}
	path := "/v1/messages"
	switch target.Vendor {
	case common.OpenAI:
		path = "/v1/chat/completions"
	case common.Gemini:
		path = "/v1beta/models/" + url.PathEscape(target.Model) + ":generateContent"
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(config.BaseURL, "/")+path, bytes.NewReader(payload))
	if err != nil {
		return common.ResponseEvent{}, fmt.Errorf("create provider request")
	}
	request.Header.Set("Content-Type", "application/json")
	switch target.Vendor {
	case common.Anthropic:
		request.Header.Set("x-api-key", config.APIKey)
		request.Header.Set("anthropic-version", "2023-06-01")
	case common.OpenAI:
		request.Header.Set("Authorization", "Bearer "+config.APIKey)
	case common.Gemini:
		request.Header.Set("x-goog-api-key", config.APIKey)
	}
	response, err := e.client.Do(request)
	if err != nil {
		if ctx.Err() != nil {
			return common.ResponseEvent{}, ctx.Err()
		}
		return common.ResponseEvent{}, fmt.Errorf("%s request failed", target.Vendor)
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return common.ResponseEvent{}, fmt.Errorf("%s HTTP %d", target.Vendor, response.StatusCode)
	}
	data, err := io.ReadAll(response.Body)
	if err != nil {
		if ctx.Err() != nil {
			return common.ResponseEvent{}, ctx.Err()
		}
		return common.ResponseEvent{}, fmt.Errorf("read %s response body", target.Vendor)
	}
	source := target
	if config.ResolvedModel != "" {
		source.Model = config.ResolvedModel
	}
	parsed, err := e.Parse(data, source)
	e.mu.Lock()
	e.usage.Input += parsed.Usage.Input
	e.usage.CacheWrite += parsed.Usage.CacheWrite
	e.usage.CacheRead += parsed.Usage.CacheRead
	e.usage.Output += parsed.Usage.Output
	e.mu.Unlock()
	e.parent.Parent().Logger().Printf("input=%d output=%d cache_write=%d cache_read=%d", parsed.Usage.Input, parsed.Usage.Output, parsed.Usage.CacheWrite, parsed.Usage.CacheRead)
	return parsed, err
}
