// ch06 is application-owned collaboration, importing only the public facade.
package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"ensemble"
)

type manuscript struct {
	mu   sync.Mutex
	text string
	path string
}

func (m *manuscript) read() string { m.mu.Lock(); defer m.mu.Unlock(); return m.text }
func (m *manuscript) write(text string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if err := os.WriteFile(m.path, []byte(text), 0600); err != nil {
		return err
	}
	m.text = text
	return nil
}
func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
func run() (result error) {
	config, err := settings()
	if err != nil {
		return err
	}
	path := os.Getenv("CH06_LOG")
	if path == "" {
		path = "workflow.jsonl"
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	records := make(chan map[string]any, 256)
	var writer sync.WaitGroup
	writer.Add(1)
	var writeErr error
	go func() {
		defer writer.Done()
		out := json.NewEncoder(io.MultiWriter(os.Stdout, file))
		for record := range records {
			if err := out.Encode(record); err != nil && writeErr == nil {
				writeErr = err
			}
		}
	}()
	root := ensemble.New(log.New(os.Stderr, "", 0))
	defer func() {
		err := root.Shutdown()
		close(records)
		writer.Wait()
		if result == nil {
			result = err
		}
		if result == nil {
			result = writeErr
		}
	}()
	store := &manuscript{path: "manuscript.txt"}
	roles := []string{"author", "editor", "reviewer"}
	agents := map[string]ensemble.Agent{}
	for _, role := range roles {
		instructions := map[string]string{"author": "You are the author. Write the requested manuscript with write_draft, verify its word_count, and finish with a short summary.", "editor": "You are the editor. Read the stored draft with read_draft and improve it with edit_draft. Finish with a short summary.", "reviewer": "You are the reviewer. Read the supplied edited manuscript and call review with decision accept or reject and notes. Finish with a short verdict."}
		roleConfig := config
		roleConfig.System = instructions[role] + " Explain what you will do and why before each tool call."
		roleConfig.HistoryPath = path + "." + role + ".jsonl"
		roleConfig.EnableTools = false
		a, err := root.NewAgent(roleConfig)
		if err != nil {
			return err
		}
		agents[role] = a
		tag := role
		a.Observe(func(event ensemble.Observation) {
			record := observation(tag, event)
			if record != nil {
				select {
				case records <- record:
				default:
				}
			}
		})
		if err := register(a, role, store); err != nil {
			return err
		}
	}
	var pipeline sync.WaitGroup
	var pipelineErr error
	started := false
	defer pipeline.Wait()
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 4096), 4<<20)
	for scanner.Scan() {
		var command struct {
			Kind  string `json:"kind"`
			Text  string `json:"text"`
			Agent string `json:"agent"`
		}
		decoder := json.NewDecoder(strings.NewReader(scanner.Text()))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&command); err != nil {
			return err
		}
		var extra any
		if decoder.Decode(&extra) != io.EOF {
			return fmt.Errorf("trailing protocol data")
		}
		switch command.Kind {
		case "prompt":
			if started || strings.TrimSpace(command.Text) == "" {
				return fmt.Errorf("exactly one nonempty prompt is supported")
			}
			started = true
			pipeline.Add(1)
			go func(topic string) {
				defer pipeline.Done()
				for _, role := range roles {
					prompt := topic
					if role != "author" {
						prompt = "Stored manuscript for " + role + ":\n" + store.read()
					}
					if _, err := agents[role].Ask(context.Background(), prompt); err != nil {
						pipelineErr = fmt.Errorf("%s: %w", role, err)
						return
					}
					records <- map[string]any{"observation": "turn_ended", "agent": role}
					records <- map[string]any{"pipeline": role + "_done"}
				}
			}(command.Text)
		case "hint", "interrupt":
			a := agents[command.Agent]
			if a == nil {
				return fmt.Errorf("unknown agent %q", command.Agent)
			}
			var input ensemble.Input = ensemble.Interrupt{}
			if command.Kind == "hint" {
				if strings.TrimSpace(command.Text) == "" {
					return fmt.Errorf("hint requires text")
				}
				input = ensemble.Hint{Text: command.Text}
			}
			if err := a.Post(input); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown kind %q", command.Kind)
		}
	}
	pipeline.Wait()
	if err := scanner.Err(); err != nil {
		return err
	}
	return pipelineErr
}
func observation(role string, event ensemble.Observation) map[string]any {
	r := map[string]any{"agent": role}
	switch v := event.(type) {
	case ensemble.StateChanged:
		r["observation"] = "state_changed"
		r["from"] = v.From
		r["to"] = v.To
	case ensemble.PartDelta:
		r["observation"] = "part_delta"
		r["part_id"] = v.PartID
		r["chunk"] = v.Chunk
	case ensemble.PartFinal:
		if v.Part.Type == "thinking" || v.Part.Type == "opaque" {
			return nil
		}
		r["observation"] = "part_final"
		r["seq"] = v.Seq
		r["part_id"] = v.PartID
		r["part"] = v.Part
	default:
		return nil
	}
	return r
}
func register(a ensemble.Agent, role string, m *manuscript) error {
	type tool struct {
		name, description, schema string
		run                       ensemble.ToolHandler
	}
	tools := []tool{}
	textSchema := `{"type":"object","properties":{"text":{"type":"string"}},"required":["text"],"additionalProperties":false}`
	emptySchema := `{"type":"object","properties":{},"additionalProperties":false}`
	write := func(_ ensemble.ToolContext, raw json.RawMessage) (string, error) {
		var p struct {
			Text *string `json:"text"`
		}
		if err := decode(raw, &p); err != nil {
			return "", err
		}
		if p.Text == nil {
			return "", fmt.Errorf("text is required")
		}
		if err := m.write(*p.Text); err != nil {
			return "", err
		}
		return "stored manuscript", nil
	}
	switch role {
	case "author":
		tools = append(tools, tool{"write_draft", "Store the first draft text.", textSchema, write}, tool{"word_count", "Count whitespace-separated words in the stored draft.", emptySchema, func(_ ensemble.ToolContext, raw json.RawMessage) (string, error) {
			if err := decode(raw, &struct{}{}); err != nil {
				return "", err
			}
			return fmt.Sprint(len(strings.Fields(m.read()))), nil
		}}, tool{"think", "Wait seconds, then return thought.", `{"type":"object","properties":{"seconds":{"type":"number","minimum":0,"maximum":60},"thought":{"type":"string"}},"required":["seconds","thought"],"additionalProperties":false}`, func(_ ensemble.ToolContext, raw json.RawMessage) (string, error) {
			var p struct {
				Seconds *float64 `json:"seconds"`
				Thought *string  `json:"thought"`
			}
			if err := decode(raw, &p); err != nil {
				return "", err
			}
			if p.Seconds == nil || p.Thought == nil || math.IsNaN(*p.Seconds) || *p.Seconds < 0 || *p.Seconds > 60 {
				return "", fmt.Errorf("seconds 0..60 and thought required")
			}
			time.Sleep(time.Duration(*p.Seconds * float64(time.Second)))
			return *p.Thought, nil
		}})
	case "editor":
		tools = append(tools, tool{"read_draft", "Read the stored manuscript.", emptySchema, func(_ ensemble.ToolContext, raw json.RawMessage) (string, error) {
			if err := decode(raw, &struct{}{}); err != nil {
				return "", err
			}
			return m.read(), nil
		}}, tool{"edit_draft", "Replace the stored manuscript.", textSchema, write})
	case "reviewer":
		tools = append(tools, tool{"review", "Accept or reject the edited manuscript with notes.", `{"type":"object","properties":{"decision":{"type":"string","enum":["accept","reject"]},"notes":{"type":"string"}},"required":["decision","notes"],"additionalProperties":false}`, func(_ ensemble.ToolContext, raw json.RawMessage) (string, error) {
			var p struct {
				Decision string  `json:"decision"`
				Notes    *string `json:"notes"`
			}
			if err := decode(raw, &p); err != nil {
				return "", err
			}
			if (p.Decision != "accept" && p.Decision != "reject") || p.Notes == nil {
				return "", fmt.Errorf("accept/reject decision and notes required")
			}
			data, _ := json.Marshal(p)
			if err := os.WriteFile("review.json", data, 0600); err != nil {
				return "", err
			}
			return string(data), nil
		}})
	}
	for _, t := range tools {
		if err := a.RegisterTool(ensemble.ToolDefinition{Name: t.name, Description: t.description, Parameters: json.RawMessage(t.schema)}, t.run); err != nil {
			return err
		}
	}
	return nil
}
func decode(raw json.RawMessage, value any) error {
	d := json.NewDecoder(strings.NewReader(string(raw)))
	d.DisallowUnknownFields()
	if err := d.Decode(value); err != nil {
		return err
	}
	var extra any
	if d.Decode(&extra) != io.EOF {
		return fmt.Errorf("trailing arguments")
	}
	return nil
}

func settings() (ensemble.Config, error) {
	name := os.Getenv("LLM_VENDOR")
	if name == "" {
		name = "anthropic"
	}
	vendor := ensemble.Anthropic
	switch name {
	case "anthropic":
	case "openai":
		vendor = ensemble.OpenAI
	case "gemini":
		vendor = ensemble.Gemini
	default:
		return ensemble.Config{}, fmt.Errorf("unknown vendor %q", name)
	}
	value := func(generic, suffix string) string {
		if s := os.Getenv(generic); s != "" {
			return s
		}
		return os.Getenv(strings.ToUpper(name) + suffix)
	}
	base := value("LLM_BASE_URL", "_BASE_URL")
	if base == "" {
		base = map[string]string{"anthropic": "https://api.anthropic.com", "openai": "https://api.openai.com", "gemini": "https://generativelanguage.googleapis.com"}[name]
	}
	requireExplanation := false
	if value := os.Getenv("LLM_REQUIRE_TOOL_EXPLANATION"); value != "" {
		var err error
		requireExplanation, err = strconv.ParseBool(value)
		if err != nil {
			return ensemble.Config{}, fmt.Errorf("LLM_REQUIRE_TOOL_EXPLANATION must be a boolean")
		}
	}
	limit := 0
	if raw := os.Getenv("LLM_TOOL_BATCH_LIMIT"); raw != "" {
		var err error
		limit, err = strconv.Atoi(raw)
		if err != nil || limit < 0 {
			return ensemble.Config{}, fmt.Errorf("LLM_TOOL_BATCH_LIMIT must be nonnegative")
		}
	}
	return ensemble.Config{ToolBatchLimit: limit, Workspace: os.Getenv("ENSEMBLE_WORKSPACE"), DataDir: os.Getenv("ENSEMBLE_DATA_DIR"), EnableTools: true, RequireToolExplanation: requireExplanation, Vendor: vendor, BaseURL: base, APIKey: value("LLM_API_KEY", "_API_KEY"), Model: value("LLM_MODEL", "_MODEL"), ResolvedModel: os.Getenv("LLM_RESOLVED_MODEL"), System: "You are a helpful coding assistant. Use the available tools to inspect facts before claiming results. Before each tool batch, write visible text explaining what you will do and why. Tools run in the current working directory, which is not an OS sandbox.", MaxTokens: 1024}, nil
}
