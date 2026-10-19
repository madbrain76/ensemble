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
	"strings"

	"ensemble"
)

func Run(ctx context.Context, args []string, in io.Reader, out, diagnostics io.Writer) (result error) {
	chat := len(args) == 1 && args[0] == "chat"
	renderMode := len(args) == 2 && args[0] == "render"
	dumpMode := len(args) == 1 && args[0] == "dump"
	if len(args) != 0 && !chat && !renderMode && !dumpMode {
		return errors.New("usage: ch02 [chat | render LOG | dump]")
	}
	config, err := ensemble.ConfigFromEnv(os.Getenv)
	if err != nil {
		return err
	}
	if dumpMode && config.Model == "" {
		config.Model = "offline-dump"
	}
	config.DataDir = "cr"
	app := ensemble.New(log.New(diagnostics, "", log.LstdFlags))
	defer func() {
		if err := app.Shutdown(); err != nil && result == nil {
			result = err
		}
	}()
	a, err := app.NewAgent(config)
	if err != nil {
		return err
	}
	path := os.Getenv("CH02_LOG")
	if path == "" {
		path = "ch02.jsonl"
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
	// Sessions start fresh. Persist even failed attempts, but never hide the
	// original session error behind a secondary disk error.
	defer func() {
		if err := app.Shutdown(); err != nil && result == nil {
			result = err
		}
		f, err := os.Create(path)
		if err == nil {
			err = a.Log().Dump(f)
			closeErr := f.Close()
			if err == nil {
				err = closeErr
			}
		}
		if err != nil && result == nil {
			result = errors.New("saving history failed")
		}
	}()
	if chat {
		fmt.Fprintln(diagnostics, "Chat ready — Ctrl-D to quit")
	}
	// ReadString grows to fit a physical line instead of imposing Scanner's
	// default 64 KiB ceiling on otherwise valid questions.
	reader := bufio.NewReader(in)
	encoder := json.NewEncoder(out)
	for {
		if chat {
			fmt.Fprint(diagnostics, "you> ")
		}
		line, readErr := reader.ReadString('\n')
		if readErr != nil && readErr != io.EOF {
			return errors.New("reading input failed")
		}
		if readErr == io.EOF && len(line) == 0 {
			break
		}
		// EOF may accompany a final line without a newline; process it once.
		question := strings.TrimSuffix(strings.TrimSuffix(line, "\n"), "\r")
		if chat && strings.TrimSpace(question) == "" {
			continue
		}
		if !chat {
			var input struct {
				User      string  `json:"user"`
				Ephemeral *string `json:"ephemeral"`
			}
			if err := json.Unmarshal([]byte(line), &input); err != nil {
				return errors.New("invalid input JSON")
			}
			if input.Ephemeral != nil {
				if err := a.Log().Append(ensemble.Event{Type: ensemble.MessageReceived, Message: &ensemble.MessageData{Actor: ensemble.System, Parts: ensemble.PartList{{Type: "text", Text: *input.Ephemeral}}}}); err != nil {
					return err
				}
				if err := encoder.Encode(struct {
					Ack string `json:"ack"`
				}{"ephemeral"}); err != nil {
					return err
				}
				continue
			}
			question = input.User
		}
		reply, err := a.Ask(ctx, question)
		if err != nil {
			return err
		}
		if chat {
			if _, err := fmt.Fprintf(out, "\nclaude> %s\n\n", reply); err != nil {
				return err
			}
			u := a.Usage()
			fmt.Fprintf(diagnostics, "[%d turns | %d input tokens | %d output tokens]\n", len(a.History()), u.Input, u.Output)
		} else if err := encoder.Encode(map[string]string{"assistant": reply}); err != nil {
			return err
		}
	}
	if !chat {
		return encoder.Encode(map[string]ensemble.Usage{"usage": a.Usage()})
	}
	return nil
}
