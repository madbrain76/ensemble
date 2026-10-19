// Command ch01 defaults to JSON lines; human chat is explicitly opt-in.
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
	"strings"

	"ensemble"
)

func main() {
	// The process owns cancellation, not an arbitrary HTTP response timer.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	if err := run(ctx, os.Args[1:], os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string, in io.Reader, out, diagnostics io.Writer) error {
	chat := len(args) == 1 && args[0] == "chat"
	if len(args) != 0 && !chat {
		return errors.New("usage: ch01 [chat]")
	}
	app := ensemble.New(log.New(diagnostics, "", 0))
	a, err := app.NewAgent(ensemble.Config{
		BaseURL: os.Getenv("ANTHROPIC_BASE_URL"),
		APIKey:  os.Getenv("ANTHROPIC_API_KEY"),
		Model:   os.Getenv("ANTHROPIC_MODEL"),
	})
	if err != nil {
		return err
	}
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
				User string `json:"user"`
			}
			if err := json.Unmarshal([]byte(line), &input); err != nil {
				return errors.New("invalid input JSON")
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
