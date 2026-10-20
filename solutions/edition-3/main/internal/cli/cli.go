// Package cli adapts process I/O to the shared core; it owns no conversation logic.
package cli

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"ensemble"
)

func Run(ctx context.Context, args []string, in io.Reader, out, diagnostics io.Writer) (result error) {
	chat := len(args) == 1 && args[0] == "chat"
	renderMode := len(args) == 2 && args[0] == "render"
	dumpMode := len(args) == 1 && args[0] == "dump"
	if len(args) != 0 && !chat && !renderMode && !dumpMode {
		return errors.New("usage: ch01 [chat | render LOG | dump]")
	}
	config, err := ensemble.ConfigFromEnv(os.Getenv)
	if err != nil {
		return err
	}
	if dumpMode && config.Model == "" {
		config.Model = "offline-dump"
	}
	path := os.Getenv("CH07_AGENT_LOG")
	if path == "" {
		path = os.Getenv("CH02_LOG")
	}
	if path == "" {
		path = "ch02.jsonl"
	}
	if !renderMode && !dumpMode {
		config.LogPath = path
	}
	config.DataDir = "cr"
	app := ensemble.New(log.New(diagnostics, "", log.LstdFlags))
	tracePath := os.Getenv("E3_TRACE")
	if tracePath == "" && !renderMode && !dumpMode {
		tracePath = "api.log"
	}
	if path := tracePath; path != "" {
		trace, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			return errors.New("opening trace failed")
		}
		defer trace.Close()
		app.SetTrace(trace)
	}
	defer func() {
		if err := app.Shutdown(); err != nil && result == nil {
			result = err
		}
	}()
	a, err := app.NewAgent(config)
	if err != nil {
		return err
	}
	if renderMode || dumpMode {
		if renderMode {
			path = args[1]
		}
		f, err := os.Open(path)
		if err != nil {
			return errors.New("opening history failed")
		}
		defer f.Close()
		if err := a.Log().Load(f); err != nil {
			return err
		}
		if dumpMode {
			return a.Log().Dump(out)
		}
		b, err := a.Render()
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(out, string(b))
		return err
	}
	if chat {
		return runChat(ctx, a, in, out, diagnostics)
	}
	// ReadString grows to fit a physical line instead of imposing Scanner's
	// default 64 KiB ceiling on otherwise valid questions.
	reader := bufio.NewReader(in)
	encoder := json.NewEncoder(out)
	progress := a.Observe()
	defer progress.Close()
	for {
		line, readErr := reader.ReadString('\n')
		if readErr != nil && readErr != io.EOF {
			return errors.New("reading input failed")
		}
		if readErr == io.EOF && len(line) == 0 {
			break
		}
		var input struct {
			User      string  `json:"user"`
			Kind      string  `json:"kind"`
			Text      string  `json:"text"`
			Ephemeral *string `json:"ephemeral"`
		}
		if err := json.Unmarshal([]byte(line), &input); err != nil {
			return errors.New("invalid input JSON")
		}
		if input.Kind == "prompt" {
			input.User = input.Text
		}
		if input.Kind == "hint" {
			if err := a.Post(ensemble.Hint{Text: input.Text}); err != nil {
				return err
			}
			if err := encoder.Encode(map[string]string{"ack": "hint"}); err != nil {
				return err
			}
			continue
		}
		if input.Ephemeral != nil {
			if err := a.Post(ensemble.Ephemeral{Text: *input.Ephemeral}); err != nil {
				return err
			}
			if err := encoder.Encode(struct {
				Ack string `json:"ack"`
			}{"ephemeral"}); err != nil {
				return err
			}
			continue
		}
		reply, err := protocolTurn(ctx, a, input.User, progress, encoder, input.Kind == "prompt" || a.Config().Streaming())
		if err != nil {
			return err
		}
		if err := encoder.Encode(map[string]string{"assistant": reply}); err != nil {
			return err
		}
	}
	return encoder.Encode(map[string]ensemble.Usage{"usage": a.Usage()})
}

// Deltas are optional progress; the answer still comes from Submit's private
// reply. Drain again after completion so scheduling cannot overtake queued data.
func protocolTurn(ctx context.Context, a ensemble.Agent, text string, q ensemble.ObservationQueue, encoder *json.Encoder, observations bool) (string, error) {
	drain := func() error {
		for _, o := range q.Drain() {
			if !observations {
				continue
			}
			var record any
			switch v := o.(type) {
			case ensemble.StateChanged:
				record = struct {
					Observation string `json:"observation"`
					ensemble.StateChanged
				}{"state_changed", v}
			case ensemble.PartDelta:
				record = struct {
					Observation string `json:"observation"`
					ensemble.PartDelta
				}{"part_delta", v}
			case ensemble.PartFinal:
				record = struct {
					Observation string `json:"observation"`
					ensemble.PartFinal
					Text string `json:"text,omitempty"`
					Tool string `json:"tool,omitempty"`
				}{"part_final", v, v.Part.Text, v.Part.Name}
			default:
				continue
			}
			if err := encoder.Encode(record); err != nil {
				return err
			}
		}
		return nil
	}
	reply := a.Submit(ctx, text)
	for {
		select {
		case <-q.Wake():
			if err := drain(); err != nil {
				return "", err
			}
		case result := <-reply:
			if err := drain(); err != nil {
				return "", err
			}
			return result.Text, result.Err
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
}
