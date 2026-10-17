package ensemble_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	ensemble "ensemble"
)

func offlineAgent(t *testing.T) ensemble.Agent {
	t.Helper()
	root := ensemble.New(log.New(io.Discard, "", 0))
	a, err := root.NewAgent(ensemble.Config{System: "stable instruction", MaxTokens: 256})
	if err != nil {
		t.Fatal(err)
	}
	if a.Log().Parent() != a || a.Engine().Parent() != a || a.Parent() != root {
		t.Fatal("broken owner chain")
	}
	return a
}
func appendFact(t *testing.T, a ensemble.Agent, event ensemble.Event) {
	t.Helper()
	if err := a.Log().Append(event); err != nil {
		t.Fatal(err)
	}
}
func issuer(v ensemble.Vendor) ensemble.Provenance {
	s := ensemble.Messages
	if v == ensemble.OpenAI {
		s = ensemble.ChatCompletions
	}
	if v == ensemble.Gemini {
		s = ensemble.GenerateContent
	}
	return ensemble.Provenance{Vendor: v, Model: "exact-model", Surface: s}
}
func user(text string) ensemble.Event {
	return ensemble.Event{Type: "message_received", Message: &ensemble.MessageEvent{Actor: "human", Parts: []ensemble.Part{{Type: "text", Text: text}}}}
}
func response(parts ...ensemble.Part) ensemble.Event {
	return ensemble.Event{Type: "response_ended", Response: &ensemble.ResponseEvent{From: issuer(ensemble.Anthropic), Parts: parts, Usage: ensemble.Usage{Input: 3, CacheWrite: 4, CacheRead: 5, Output: 6}}}
}

// TestChapter2HistoryLifecycle exercises the public owned reducer, not a copy of
// applyEvent in the test. Timestamps are deliberately backwards: seq is order.
func TestChapter2HistoryLifecycle(t *testing.T) {
	a := offlineAgent(t)
	if got := a.Log().Context().Turn; got != ensemble.Idle {
		t.Fatal(got)
	}
	ephemeral := ensemble.Event{Type: "message_received", Time: time.Unix(20, 0), Message: &ensemble.MessageEvent{Actor: "system", Parts: []ensemble.Part{{Type: "text", Text: "one shot"}}}}
	appendFact(t, a, ephemeral)
	appendFact(t, a, ensemble.Event{Type: "response_started", Response: &ensemble.ResponseEvent{From: issuer(ensemble.Anthropic)}})
	if c := a.Log().Context(); c.Turn != ensemble.InputPending || len(c.Dialogue) != 0 || len(c.Ephemera) != 1 {
		t.Fatalf("pending %+v", c)
	}
	appendFact(t, a, user("first"))
	appendFact(t, a, ensemble.Event{Type: "request_sent", Request: &ensemble.RequestEvent{To: issuer(ensemble.Anthropic)}})
	if c := a.Log().Context(); c.Turn != ensemble.InFlight || len(c.Ephemera) != 0 {
		t.Fatalf("inflight %+v", c)
	}
	appendFact(t, a, ensemble.Event{Type: "tool_called", Tool: &ensemble.ToolEvent{CallID: "a"}})
	if a.Log().Context().Turn != ensemble.InFlight {
		t.Fatal("known irrelevant event changed state")
	}
	appendFact(t, a, response(ensemble.Part{Type: "tool_call", CallID: "a", Name: "lookup", Args: json.RawMessage(`{}`), From: issuer(ensemble.Anthropic)}, ensemble.Part{Type: "tool_call", CallID: "b", Name: "lookup", Args: json.RawMessage(`{}`), From: issuer(ensemble.Anthropic)}))
	if a.Log().Context().Turn != ensemble.ToolsPending {
		t.Fatal("tool call state")
	}
	for _, id := range []string{"a", "b"} {
		appendFact(t, a, ensemble.Event{Type: "tool_returned", Tool: &ensemble.ToolEvent{CallID: id, Parts: []ensemble.Part{{Type: "text", Text: "done"}}}})
		want := ensemble.ToolsPending
		if id == "b" {
			want = ensemble.InputPending
		}
		if a.Log().Context().Turn != want {
			t.Fatal("outstanding tools state")
		}
	}
	appendFact(t, a, ensemble.Event{Type: "request_sent", Request: &ensemble.RequestEvent{To: issuer(ensemble.Anthropic)}})
	appendFact(t, a, ensemble.Event{Type: "error_occurred", Error: &ensemble.ErrorEvent{Message: "network failure"}})
	if c := a.Log().Context(); c.Turn != ensemble.Idle || c.Usage != (ensemble.Usage{Input: 3, CacheWrite: 4, CacheRead: 5, Output: 6}) {
		t.Fatalf("final %+v", c)
	}
	before := a.Log().Context()
	appendFact(t, a, ensemble.Event{Type: "error_occurred", Error: &ensemble.ErrorEvent{Message: "irrelevant"}})
	if !reflect.DeepEqual(before, a.Log().Context()) {
		t.Fatal("idle error was not identity")
	}
	var dump bytes.Buffer
	if err := a.Log().Dump(&dump); err != nil {
		t.Fatal(err)
	}
	b := offlineAgent(t)
	if err := b.Log().Load(bytes.NewReader(dump.Bytes())); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(a.Log().Context(), b.Log().Context()) {
		t.Fatal("reload changed context")
	}
	snapshot := a.Log().Events()
	snapshot[0].Message.Parts[0].Text = "mutated"
	c := a.Log().Context()
	c.Dialogue[0].Parts[0].Text = "mutated"
	if a.Log().Events()[0].Message.Parts[0].Text != "one shot" || a.Log().Context().Dialogue[0].Parts[0].Text != "first" {
		t.Fatal("snapshot aliases log")
	}
}

func TestChapter2LogValidation(t *testing.T) {
	a := offlineAgent(t)
	appendFact(t, a, user("kept"))
	var valid bytes.Buffer
	_ = a.Log().Dump(&valid)
	lines := strings.Split(strings.TrimSpace(valid.String()), "\n")
	event := lines[1]
	cases := []string{
		`{"log_version":2}` + "\n" + event,
		event + "\n" + event,
		strings.Replace(event, `"message_received"`, `"future_event"`, 1),
		strings.Replace(event, `"human"`, `"nobody"`, 1),
		strings.Replace(event, `"text"`, `"future_part"`, 1),
		`{"seq":1,"type":"message_received","message":{"actor":"human","parts":[{"type":"blob","path":"legacy"}]}}`,
		`{"seq":1,"type":"message_received","message":{"actor":"human","parts":[{"type":"blob","ref":{"kind":0,"locator":"bad"}}]}}`,
		`{"seq":1,"type":"request_sent","request":{"to":{"vendor":"other","surface":"messages","model":"x"}}}`,
		`{"seq":1,"type":"request_sent","request":{"to":{"vendor":"openai","surface":"messages","model":"x"}}}`,
		`{"seq":1,"type":"response_ended","response":{"from":{"vendor":"anthropic","surface":"messages","model":"x"},"usage":{"input":-1}}}`,
		`{"seq":1,"type":"message_received","message":{"actor":"human","parts":[]},"error":{"message":"two"}}`,
		`{"seq":1,"type":"redacted","redact":{"from":1,"to":1,"level":"redact_result"}}`,
		`{"seq":1,"type":"message_received","message":{"actor":"human","parts":[{"type":"tool_call","from":{"vendor":"anthropic","surface":"messages","model":"x"},"name":"x","args":"not-json-object-but-valid-string"}]}}` + "\n" + `{"seq":2,"type":"redacted","redact":{"from":1,"to":1,"level":"invented"}}`,
		`{`,
	}
	for i, input := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			if err := a.Log().Load(strings.NewReader(input)); err == nil {
				t.Fatal("invalid log accepted")
			}
			var after bytes.Buffer
			_ = a.Log().Dump(&after)
			if after.String() != valid.String() {
				t.Fatal("failed load changed history")
			}
		})
	}
	if err := a.Log().Load(strings.NewReader(event)); err != nil {
		t.Fatal("headerless current log", err)
	}
	var alias ensemble.Provenance
	if err := json.Unmarshal([]byte(`{"vendor":"gemini","model":"x","surface":"generatecontent"}`), &alias); err != nil {
		t.Fatal(err)
	}
	raw, _ := json.Marshal(alias)
	if !strings.Contains(string(raw), `"generate_content"`) {
		t.Fatal(string(raw))
	}
	for _, kind := range []ensemble.RefKind{ensemble.RefPath, ensemble.RefURI, ensemble.RefHandle} {
		fact := ensemble.Event{Type: "message_received", Message: &ensemble.MessageEvent{Actor: "human", Parts: []ensemble.Part{{Type: "blob", MIME: "text/plain", Ref: &ensemble.Ref{Kind: kind, Locator: "unchanged"}}}}}
		appendFact(t, a, fact)
	}
	var refs bytes.Buffer
	_ = a.Log().Dump(&refs)
	b := offlineAgent(t)
	if err := b.Log().Load(&refs); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(a.Log().Context(), b.Log().Context()) {
		t.Fatal("reference roundtrip")
	}
}

func redactionAgent(t *testing.T) ensemble.Agent {
	a := offlineAgent(t)
	appendFact(t, a, user("secret prose"))
	appendFact(t, a, response(ensemble.Part{Type: "text", Text: "assistant prose"}, ensemble.Part{Type: "tool_call", CallID: "call", Name: "lookup", Args: json.RawMessage(`{"q":"arg"}`), From: issuer(ensemble.Anthropic)}))
	parts := []ensemble.Part{{Type: "text", Text: "secret result"}}
	for _, kind := range []ensemble.RefKind{ensemble.RefPath, ensemble.RefURI, ensemble.RefHandle} {
		parts = append(parts, ensemble.Part{Type: "blob", MIME: "image/png", Ref: &ensemble.Ref{Kind: kind, Locator: "private-locator"}})
	}
	appendFact(t, a, ensemble.Event{Type: "tool_returned", Tool: &ensemble.ToolEvent{CallID: "call", Parts: parts}})
	return a
}
func TestChapter2RedactionLevels(t *testing.T) {
	for _, level := range []string{"redact_result", "redact_tool", "redact_dialogue", "redact_summary"} {
		t.Run(level, func(t *testing.T) {
			a := redactionAgent(t)
			before := a.Log().Events()
			r := &ensemble.Redaction{From: 1, To: 3, Level: level, Reason: "budget"}
			if level == "redact_summary" {
				r.Replacement = []ensemble.Part{{Type: "text", Text: "short summary"}}
			}
			appendFact(t, a, ensemble.Event{Type: "redacted", Redact: r})
			if !reflect.DeepEqual(before, a.Log().Events()[:3]) {
				t.Fatal("redaction edited facts")
			}
			c := a.Log().Context()
			switch level {
			case "redact_result":
				if c.Dialogue[0].Parts[0].Text != "secret prose" || c.Dialogue[1].Parts[1].Type != "tool_call" {
					t.Fatal("result filter affected other content")
				}
				stubs := c.Dialogue[2].Parts[0].Parts
				if stubs[0].Type != "redacted" || !strings.Contains(stubs[0].Stub, "13 bytes") {
					t.Fatalf("stub %+v", stubs[0])
				}
				for i, kind := range []ensemble.RefKind{ensemble.RefPath, ensemble.RefURI, ensemble.RefHandle} {
					if stubs[i+1].Ref == nil || stubs[i+1].Ref.Kind != kind || stubs[i+1].Ref.Locator != "private-locator" {
						t.Fatal("stub lost exact ref")
					}
				}
				for _, v := range []ensemble.Vendor{ensemble.Anthropic, ensemble.OpenAI, ensemble.Gemini} {
					wire, err := a.Engine().Render(c, issuer(v))
					if err != nil {
						t.Fatal(err)
					}
					if bytes.Contains(wire, []byte("secret result")) || bytes.Contains(wire, []byte("private-locator")) {
						t.Fatal("redacted bytes leaked")
					}
				}
			case "redact_tool":
				if len(c.Dialogue[1].Parts) != 1 || len(c.Dialogue[2].Parts) != 0 || c.Dialogue[0].Parts[0].Text != "secret prose" {
					t.Fatal("tool filter", c)
				}
			case "redact_dialogue":
				if c.Dialogue[0].Parts[0].Type != "redacted" || c.Dialogue[1].Parts[1].Type != "tool_call" || c.Dialogue[2].Parts[0].Parts[0].Text != "secret result" {
					t.Fatal("dialogue filter", c)
				}
			case "redact_summary":
				if len(c.Dialogue) != 1 || c.Dialogue[0].Actor != "system" || c.Dialogue[0].Parts[0].Text != "short summary" {
					t.Fatal("summary filter", c)
				}
			}
		})
	}
	a := redactionAgent(t)
	for _, v := range []ensemble.Vendor{ensemble.Anthropic, ensemble.OpenAI, ensemble.Gemini} {
		if _, err := a.Engine().Render(a.Log().Context(), issuer(v)); err == nil {
			t.Fatal("unsupported blob silently emitted")
		}
	}
}

func TestChapter2ParseAndProvenance(t *testing.T) {
	fixtures := []struct {
		vendor ensemble.Vendor
		body   string
		usage  ensemble.Usage
	}{
		{ensemble.Anthropic, `{"model":"returned-model","content":[{"type":"text","text":"a"},{"type":"thinking","thinking":"secret thought","signature":"sig"},{"type":"tool_use","id":"c1","name":"lookup","input":{"q":"x"}}],"usage":{"input_tokens":7,"cache_creation_input_tokens":300,"cache_read_input_tokens":500,"output_tokens":11}}`, ensemble.Usage{Input: 7, CacheWrite: 300, CacheRead: 500, Output: 11}},
		{ensemble.OpenAI, `{"model":"returned-model","choices":[{"message":{"content":"a","reasoning_content":"secret thought","tool_calls":[{"id":"c1","type":"function","function":{"name":"lookup","arguments":"{\"q\":\"x\"}"}}]}}],"usage":{"prompt_tokens":807,"prompt_tokens_details":{"cached_tokens":500,"cache_write_tokens":300},"completion_tokens":11,"completion_tokens_details":{"reasoning_tokens":5}}}`, ensemble.Usage{Input: 7, CacheWrite: 300, CacheRead: 500, Output: 11}},
		{ensemble.Gemini, `{"modelVersion":"returned-model","candidates":[{"content":{"parts":[{"text":"a"},{"thought":true,"text":"secret thought"},{"functionCall":{"name":"lookup","args":{"q":"x"}},"thoughtSignature":"signature-exact"}]}}],"usageMetadata":{"promptTokenCount":507,"cachedContentTokenCount":500,"candidatesTokenCount":8,"thoughtsTokenCount":3}}`, ensemble.Usage{Input: 7, CacheRead: 500, Output: 11}},
	}
	for _, fixture := range fixtures {
		t.Run(fixture.vendor.String(), func(t *testing.T) {
			a := offlineAgent(t)
			from := issuer(fixture.vendor)
			parsed, err := a.Engine().Parse([]byte(fixture.body), from)
			if err != nil {
				t.Fatal(err)
			}
			if parsed.From.Model != "returned-model" || parsed.Usage != fixture.usage {
				t.Fatalf("parsed %+v", parsed)
			}
			if len(a.Log().Events()) != 0 || a.Engine().Usage() != (ensemble.Usage{}) {
				t.Fatal("Parse mutated owner")
			}
			var call ensemble.Part
			for _, p := range parsed.Parts {
				if p.Type == "tool_call" {
					call = p
				}
				if p.Type == "tool_call" || p.Type == "opaque" {
					if p.From != parsed.From {
						t.Fatal("part lost resolved provenance")
					}
				}
			}
			if string(call.Args) != `{"q":"x"}` {
				t.Fatalf("arguments %s", call.Args)
			}
			if fixture.vendor == ensemble.Gemini && (call.CallID != "" || string(call.Opaque) != `"signature-exact"`) {
				t.Fatalf("signature or absent id changed: %+v", call)
			}
			appendFact(t, a, ensemble.Event{Type: "response_ended", Response: &parsed})
			before := a.Log().Events()
			to := parsed.From
			first, err := a.Engine().Render(a.Log().Context(), to)
			if err != nil {
				t.Fatal(err)
			}
			second, err := a.Engine().Render(a.Log().Context(), to)
			if err != nil || !bytes.Equal(first, second) {
				t.Fatal("render is not deterministic")
			}
			if !bytes.Contains(first, []byte("secret thought")) {
				t.Fatalf("same issuer lost opaque %s", first)
			}
			if fixture.vendor == ensemble.Gemini && !bytes.Contains(first, []byte(`"thoughtSignature":"signature-exact"`)) {
				t.Fatalf("lost call signature %s", first)
			}
			to.Model = "other-model"
			foreign, err := a.Engine().Render(a.Log().Context(), to)
			if err != nil {
				t.Fatal(err)
			}
			if bytes.Contains(foreign, []byte("secret thought")) || bytes.Contains(foreign, []byte("signature-exact")) {
				t.Fatal("model-bound bytes leaked")
			}
			other := issuer(ensemble.OpenAI)
			if fixture.vendor == ensemble.OpenAI {
				other = issuer(ensemble.Anthropic)
			}
			foreign, err = a.Engine().Render(a.Log().Context(), other)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Contains(foreign, []byte("lookup")) || bytes.Contains(foreign, []byte("secret thought")) {
				t.Fatal("cross-vendor tool normalization")
			}
			if !reflect.DeepEqual(before, a.Log().Events()) {
				t.Fatal("renderer mutated facts")
			}
		})
	}
}

func TestChapter2AnthropicOrderingAndIDs(t *testing.T) {
	a := offlineAgent(t)
	from := issuer(ensemble.Anthropic)
	appendFact(t, a, response(ensemble.Part{Type: "tool_call", CallID: "one", Name: "lookup", Args: json.RawMessage(`{}`), From: from}, ensemble.Part{Type: "tool_call", CallID: "two", Name: "lookup", Args: json.RawMessage(`{}`), From: from}))
	appendFact(t, a, ensemble.Event{Type: "tool_returned", Tool: &ensemble.ToolEvent{CallID: "one", Parts: []ensemble.Part{{Type: "text", Text: "first"}}}})
	appendFact(t, a, user("human between results"))
	appendFact(t, a, ensemble.Event{Type: "tool_returned", Tool: &ensemble.ToolEvent{CallID: "two", Parts: []ensemble.Part{{Type: "text", Text: "second"}}}})
	appendFact(t, a, ensemble.Event{Type: "message_received", Message: &ensemble.MessageEvent{Actor: "system", Parts: []ensemble.Part{{Type: "text", Text: "ephemeral"}}}})
	wire, err := a.Engine().Render(a.Log().Context(), from)
	if err != nil {
		t.Fatal(err)
	}
	var request struct {
		System   string `json:"system"`
		Messages []struct {
			Role    string `json:"role"`
			Content []struct {
				Type string `json:"type"`
				ID   string `json:"tool_use_id"`
				Text string `json:"text"`
			} `json:"content"`
		} `json:"messages"`
	}
	if err := json.Unmarshal(wire, &request); err != nil {
		t.Fatal(err)
	}
	if request.System != "stable instruction" || len(request.Messages) != 2 {
		t.Fatalf("grouping %s", wire)
	}
	blocks := request.Messages[1].Content
	if len(blocks) != 4 || blocks[0].ID != "one" || blocks[1].ID != "two" || blocks[2].Text != "human between results" || blocks[3].Text != "ephemeral" {
		t.Fatalf("result ordering %s", wire)
	}
	b := offlineAgent(t)
	for i := 0; i < 2; i++ {
		appendFact(t, b, response(ensemble.Part{Type: "tool_call", Name: "no-id", Args: json.RawMessage(`{}`), From: issuer(ensemble.Gemini)}))
	}
	out, err := b.Engine().Render(b.Log().Context(), issuer(ensemble.OpenAI))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(out, []byte("call-1-0")) || !bytes.Contains(out, []byte("call-2-0")) {
		t.Fatalf("ids not seq-derived %s", out)
	}
	if b.Log().Events()[0].Response.Parts[0].CallID != "" {
		t.Fatal("synthetic id stored as fact")
	}
}

func TestChapter2OfflineCredentialsAndErrors(t *testing.T) {
	a := offlineAgent(t)
	if _, err := a.Ask(context.Background(), "hello"); err == nil {
		t.Fatal("missing model accepted for send")
	}
	for _, body := range []string{`{`, `{}`, `{"content":[],"usage":{"input_tokens":null,"output_tokens":1}}`, `{"choices":[],"usage":{"prompt_tokens":1,"completion_tokens":1}}`, `{"candidates":[],"usageMetadata":{"promptTokenCount":1,"candidatesTokenCount":1}}`} {
		from := issuer(ensemble.Anthropic)
		if strings.Contains(body, "choices") {
			from = issuer(ensemble.OpenAI)
		}
		if strings.Contains(body, "candidates") {
			from = issuer(ensemble.Gemini)
		}
		if _, err := a.Engine().Parse([]byte(body), from); err == nil {
			t.Fatal("invalid response accepted", body)
		}
	}
	for _, target := range []ensemble.Provenance{{}, {Vendor: ensemble.OpenAI, Model: "x", Surface: ensemble.Messages}} {
		if _, err := a.Engine().Render(a.Log().Context(), target); err == nil {
			t.Fatal("invalid target accepted")
		}
	}
}

// TestChapter2CLIReplay runs the actual executable and verifies offline commands
// with a dead endpoint after capture. No production constructor is replaced.
func TestChapter2CLIReplay(t *testing.T) {
	binary := buildCLI(t)
	dir := t.TempDir()
	path := filepath.Join(dir, "session.jsonl")
	root, a, _ := newAgent(t, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"model":"resolved","content":[{"type":"text","text":"answer"}],"usage":{"input_tokens":9,"output_tokens":4}}`)
	})
	_ = root
	env := append(cliEnv(a.Config().BaseURL), "CH02_LOG="+path)
	cmd := exec.Command(binary)
	cmd.Env = env
	cmd.Stdin = strings.NewReader("{\"ephemeral\":\"hint\"}\n{\"user\":\"question\"}\n")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("live fake %v: %s", err, &stderr)
	}
	if !strings.HasPrefix(stdout.String(), "{\"ack\":\"ephemeral\"}\n{\"assistant\":\"answer\"}") || !strings.Contains(stdout.String(), `"cache_read":0`) {
		t.Fatalf("protocol %s", &stdout)
	}
	saved, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(saved, []byte("fixture-key")) || bytes.Contains(saved, []byte("x-api-key")) {
		t.Fatal("credential captured")
	}
	if err := a.Log().Load(bytes.NewReader(saved)); err != nil {
		t.Fatal(err)
	}
	if len(a.Log().Events()) != 5 || a.Log().Events()[4].Type != "turn_ended" || len(a.Log().Context().Ephemera) != 0 {
		t.Fatal("capture facts", string(saved))
	}
	env = append(cliEnv("http://127.0.0.1:1"), "CH02_LOG="+path)
	var rendered []byte
	for i := 0; i < 2; i++ {
		cmd = exec.Command(binary, "render", path)
		cmd.Env = env
		out, err := cmd.Output()
		if err != nil {
			t.Fatal("offline render", err)
		}
		if i == 0 {
			rendered = out
		} else if !bytes.Equal(rendered, out) {
			t.Fatal("fresh-process replay differs")
		}
	}
	cmd = exec.Command(binary, "dump")
	cmd.Env = env
	dump, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(saved, dump) {
		t.Fatal("dump rewrote original facts")
	}
	cmd = exec.Command(binary, "--protocol")
	cmd.Env = env
	if err := cmd.Run(); err == nil {
		t.Fatal("invented flag accepted")
	}
}

// buildCLI builds the only shipped application constructor path.
func buildCLI(t *testing.T) string {
	t.Helper()
	binary := filepath.Join(t.TempDir(), "student-cli")
	if out, err := exec.Command("go", "build", "-o", binary, "./cmd/ch01").CombinedOutput(); err != nil {
		t.Fatalf("build: %v %s", err, out)
	}
	return binary
}
func cliEnv(base string) []string {
	env := []string{}
	for _, entry := range os.Environ() {
		if !strings.HasPrefix(entry, "LLM_") && !strings.HasPrefix(entry, "ANTHROPIC_") && !strings.HasPrefix(entry, "CH02_LOG=") {
			env = append(env, entry)
		}
	}
	return append(env, "ANTHROPIC_BASE_URL="+base, "ANTHROPIC_API_KEY=fixture-key", "ANTHROPIC_MODEL=test-model")
}

func TestChapter2RenderRefusesMalformedParts(t *testing.T) {
	a := offlineAgent(t)
	for _, p := range []ensemble.Part{{Type: "future"}, {Type: "tool_call", Name: "tool", Args: json.RawMessage(`{`)}, {Type: "tool_call", Args: json.RawMessage(`{}`)}, {Type: "opaque", From: ensemble.Provenance{Vendor: ensemble.Anthropic, Model: "model-a", Surface: ensemble.Messages}, Data: json.RawMessage(`{`)}, {Type: "opaque", From: ensemble.Provenance{Vendor: ensemble.Anthropic, Model: "model-a", Surface: ensemble.Messages}}} {
		if _, err := a.Engine().Render(ensemble.Context{Dialogue: []ensemble.Turn{{Actor: "agent", Parts: []ensemble.Part{p}}}}, ensemble.Provenance{Vendor: ensemble.Anthropic, Model: "model-a", Surface: ensemble.Messages}); err == nil {
			t.Fatalf("accepted %+v", p)
		}
	}
	for _, raw := range []string{`"unknown"`, `null`, `5`} {
		var v ensemble.Vendor
		if err := json.Unmarshal([]byte(raw), &v); err == nil {
			t.Fatalf("vendor %s", raw)
		}
		var surface ensemble.Surface
		if err := json.Unmarshal([]byte(raw), &surface); err == nil {
			t.Fatalf("surface %s", raw)
		}
	}
	if _, err := json.Marshal(ensemble.Vendor(0)); err == nil {
		t.Fatal("zero vendor")
	}
	if _, err := json.Marshal(ensemble.Surface(0)); err == nil {
		t.Fatal("zero surface")
	}
}
