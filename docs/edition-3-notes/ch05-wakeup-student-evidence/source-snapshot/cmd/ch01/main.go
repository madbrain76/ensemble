// Command ch01 preserves the accepted executable path and adds Chapter 2 replay.
package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"

	ensemble "ensemble"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

// settings is the only CLI configuration path, shared by live and offline use.
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

func run(args []string) (result error) {
	if len(args) > 2 {
		return fmt.Errorf("usage: ch01 [chat | render LOG | dump]")
	}
	config, err := settings()
	if err != nil {
		return err
	}
	path := os.Getenv("CH02_LOG")
	if path == "" {
		path = "ch02.jsonl"
	}
	root := ensemble.New(log.New(os.Stderr, "", 0))
	config.HistoryPath = path
	agent, err := root.NewAgent(config)
	if err != nil {
		return err
	}
	if len(args) > 0 {
		switch args[0] {
		case "render":
			if len(args) != 2 {
				return fmt.Errorf("render requires a log path")
			}
			if err := load(agent, args[1]); err != nil {
				return err
			}
			body, err := agent.Engine().Render(agent.Log().Context(), agent.Engine().Target())
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(os.Stdout, string(body))
			return err
		case "dump":
			if len(args) != 1 {
				return fmt.Errorf("dump takes no arguments")
			}
			if err := load(agent, path); err != nil {
				return err
			}
			return agent.Log().Dump(os.Stdout)
		case "chat":
			if len(args) != 1 {
				return fmt.Errorf("chat takes no arguments")
			}
		default:
			return fmt.Errorf("unknown command %q", args[0])
		}
	}
	// Save failure facts as well as successful exchanges. No raw provider payloads,
	// requests, or credentials are stored, and sessions do not auto-resume.
	defer func() {
		if err := root.Shutdown(); result == nil {
			result = err
		}
		if err := save(agent, path); result == nil {
			result = err
		}
	}()
	if len(args) == 1 {
		return human(agent)
	}
	return protocol(agent)
}
func load(agent ensemble.Agent, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return agent.Log().Load(file)
}
func save(agent ensemble.Agent, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	err = agent.Log().Dump(file)
	closeErr := file.Close()
	if err != nil {
		return err
	}
	return closeErr
}

// session keeps input live while actor requests are outstanding.
func protocol(a ensemble.Agent) error { return session(a, false) }
func human(a ensemble.Agent) error {
	fmt.Fprintln(os.Stderr, "Conversation ready. Enter a message; EOF exits.")
	return session(a, true)
}

type inputLine struct {
	text string
	err  error
}

func session(a ensemble.Agent, human bool) error {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	defer signal.Stop(signals)
	observations := make(chan ensemble.Observation, 128)
	a.Observe(func(o ensemble.Observation) {
		select {
		case observations <- o:
		default:
		}
	})
	liveMode := false
	lines := make(chan inputLine)
	stop := make(chan struct{})
	var reader sync.WaitGroup
	reader.Add(1)
	go func() {
		defer reader.Done()
		defer close(lines)
		s := bufio.NewScanner(os.Stdin)
		s.Buffer(make([]byte, 4096), 4*1024*1024)
		for s.Scan() {
			select {
			case lines <- inputLine{text: s.Text()}:
			case <-stop:
				return
			}
		}
		if err := s.Err(); err != nil {
			select {
			case lines <- inputLine{err: err}:
			case <-stop:
			}
		}
	}()
	defer func() { close(stop); _ = os.Stdin.Close(); reader.Wait() }()
	results := make(chan ensemble.Result)
	var requests sync.WaitGroup
	pending := 0
	var first error
	var ephemera []ensemble.Part
	var backlog []ensemble.UserMessage
	defer requests.Wait()
	submit := func(m ensemble.UserMessage) { m.Ephemera = ephemera; ephemera = nil; backlog = append(backlog, m) }
	for lines != nil || pending > 0 || len(backlog) > 0 {
		if first != nil {
			backlog = nil
			lines = nil
			if pending == 0 {
				break
			}
		}
		if pending == 0 && len(backlog) > 0 {
			m := backlog[0]
			backlog = backlog[1:]
			reply := a.SubmitMessage(context.Background(), m)
			pending++
			requests.Add(1)
			go func() { defer requests.Done(); results <- <-reply }()
		}
		select {
		case <-signals:
			if err := a.Post(ensemble.Interrupt{}); err != nil {
				first = err
			}
		case observation := <-observations:
			if delta, ok := observation.(ensemble.PartDelta); ok && strings.HasPrefix(delta.Chunk, "hint:") {
				if human {
					fmt.Fprintln(os.Stderr, "Hint:", strings.TrimPrefix(delta.Chunk, "hint:"))
				}
			}
			if liveMode && !human {
				if record := observationRecord(observation); record != nil {
					if err := json.NewEncoder(os.Stdout).Encode(record); err != nil {
						first = err
					}
				}
			}
		case line, ok := <-lines:
			if !ok {
				if liveMode {
					pending++
					requests.Add(1)
					go func() {
						defer requests.Done()
						err := a.Drain(context.Background())
						results <- ensemble.Result{Err: err}
					}()
				}
				lines = nil
				continue
			}
			if line.err != nil {
				first = line.err
				lines = nil
				continue
			}
			if human {
				if strings.TrimSpace(line.text) == "" {
					continue
				}
				if line.text == "/interrupt" {
					if err := a.Post(ensemble.Interrupt{}); err != nil {
						first = err
					}
					continue
				}
				if strings.HasPrefix(line.text, "/hint ") {
					if err := a.Post(ensemble.Hint{Text: strings.TrimPrefix(line.text, "/hint ")}); err != nil {
						first = err
					}
					continue
				}
				submit(ensemble.UserMessage{Text: line.text})
				continue
			}
			var record map[string]json.RawMessage
			d := json.NewDecoder(strings.NewReader(line.text))
			err := d.Decode(&record)
			var trailing any
			if err == nil && record["kind"] != nil {
				liveMode = true
				var command struct {
					Kind string `json:"kind"`
					Text string `json:"text"`
				}
				decoder := json.NewDecoder(strings.NewReader(line.text))
				decoder.DisallowUnknownFields()
				if err := decoder.Decode(&command); err != nil {
					first = err
					lines = nil
					continue
				}
				var extra any
				if decoder.Decode(&extra) != io.EOF {
					first = fmt.Errorf("trailing protocol data")
					lines = nil
					continue
				}
				switch command.Kind {
				case "prompt", "hint":
					if strings.TrimSpace(command.Text) == "" {
						first = fmt.Errorf("%s requires text", command.Kind)
					} else {
						first = a.Post(ensemble.UserMessage{Text: command.Text})
					}
				case "interrupt":
					if command.Text != "" {
						first = fmt.Errorf("interrupt does not accept text")
					} else {
						first = a.Post(ensemble.Interrupt{})
					}
				default:
					first = fmt.Errorf("unknown command kind %q", command.Kind)
				}
				continue
			}
			if err != nil || d.Decode(&trailing) != io.EOF || len(record) != 1 {
				first = fmt.Errorf("protocol requires one JSON object containing user, parts, ephemeral, hint, or interrupt")
				lines = nil
				continue
			}
			for kind, raw := range record {
				switch kind {
				case "interrupt":
					var value bool
					if err := json.Unmarshal(raw, &value); err != nil || !value {
						first = fmt.Errorf("interrupt must be true")
					} else if err := a.Post(ensemble.Interrupt{}); err != nil {
						first = err
					}
				case "parts":
					var parts []ensemble.Part
					if err := json.Unmarshal(raw, &parts); err != nil || len(parts) == 0 {
						first = fmt.Errorf("parts must be a nonempty Part array")
					} else {
						submit(ensemble.UserMessage{Parts: parts})
					}
				case "user", "hint", "ephemeral":
					var text string
					if err := json.Unmarshal(raw, &text); err != nil || strings.TrimSpace(text) == "" {
						first = fmt.Errorf("%s requires nonempty text", kind)
						break
					}
					switch kind {
					case "user":
						submit(ensemble.UserMessage{Text: text})
					case "hint":
						if err := a.Post(ensemble.Hint{Text: text}); err != nil {
							first = err
						}
					case "ephemeral":
						ephemera = append(ephemera, ensemble.Part{Type: "text", Text: text})
						first = json.NewEncoder(os.Stdout).Encode(map[string]string{"ack": "ephemeral"})
					}
				default:
					first = fmt.Errorf("unknown protocol record %q", kind)
				}
			}
			if first != nil {
				lines = nil
			}
		case result := <-results:
			pending--
			if result.Err != nil {
				if errors.Is(result.Err, ensemble.ErrInterrupted) {
					if human {
						fmt.Fprintln(os.Stderr, "Interrupted; ready for another prompt.")
					}
					continue
				}
				if first == nil {
					first = result.Err
				}
				continue
			}
			if result.Text == "" {
				continue
			}
			if human {
				_, err := fmt.Fprintln(os.Stdout, "> "+result.Text)
				if first == nil {
					first = err
				}
			} else {
				err := json.NewEncoder(os.Stdout).Encode(map[string]any{"assistant": result.Text})
				if first == nil {
					first = err
				}
			}
		}
	}
	if first == nil {
		if human {
			u := a.Engine().Usage()
			fmt.Fprintf(os.Stderr, "cumulative usage: input=%d output=%d cache_write=%d cache_read=%d\n", u.Input, u.Output, u.CacheWrite, u.CacheRead)
		} else {
			first = json.NewEncoder(os.Stdout).Encode(map[string]any{"usage": a.Engine().Usage()})
		}
	}
	return first
}

func observationRecord(o ensemble.Observation) map[string]any {
	switch v := o.(type) {
	case ensemble.StateChanged:
		return map[string]any{"observation": "state_changed", "from": v.From, "to": v.To}
	case ensemble.PartDelta:
		return map[string]any{"observation": "part_delta", "part_id": v.PartID, "chunk": v.Chunk}
	case ensemble.PartFinal:
		if v.Part.Type == "thinking" || v.Part.Type == "opaque" {
			return nil
		}
		return map[string]any{"observation": "part_final", "seq": v.Seq, "part_id": v.PartID, "part": v.Part}
	}
	return nil
}
