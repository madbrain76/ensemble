package cli

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"

	"ensemble"
)

type chatLine struct {
	text string
	err  error
}
type chatObserver struct{ events chan ensemble.Observation }

func (o chatObserver) Observe(v ensemble.Observation) {
	select {
	case o.events <- v:
	default:
	}
}
func chatInput(in io.Reader, lines chan<- chatLine, done <-chan struct{}) {
	defer close(lines)
	r := bufio.NewReader(in)
	for {
		line, err := r.ReadString('\n')
		if len(line) > 0 {
			line = strings.TrimSuffix(strings.TrimSuffix(line, "\n"), "\r")
			select {
			case lines <- chatLine{text: line}:
			case <-done:
				return
			}
		}
		if err != nil {
			if err != io.EOF {
				select {
				case lines <- chatLine{err: err}:
				case <-done:
				}
			}
			return
		}
	}
}
func runChat(ctx context.Context, a ensemble.Agent, in io.Reader, out, diagnostics io.Writer) error {
	lines := make(chan chatLine)
	done := make(chan struct{})
	defer close(done)
	go chatInput(in, lines, done)
	interrupts := make(chan os.Signal, 1)
	signal.Notify(interrupts, os.Interrupt)
	defer signal.Stop(interrupts)
	progress := make(chan ensemble.Observation, 64)
	a.Subscribe(chatObserver{progress})
	var pending []<-chan ensemble.Result
	var printedTurn uint64
	fmt.Fprintln(diagnostics, "Chat ready — /interrupt or Ctrl-C ends the turn; Ctrl-D quits")
	fmt.Fprint(diagnostics, "you> ")
	for lines != nil || len(pending) > 0 {
		var reply <-chan ensemble.Result
		if len(pending) > 0 {
			reply = pending[0]
		}
		select {
		case line, ok := <-lines:
			if !ok {
				lines = nil
				continue
			}
			if line.err != nil {
				return line.err
			}
			text := line.text
			if strings.TrimSpace(text) == "" {
				continue
			}
			if text == "/interrupt" {
				if err := a.Post(ensemble.Interrupt{}); err != nil {
					return err
				}
				fmt.Fprintln(diagnostics, "interrupt: turn ending; jobs remain available")
			} else {
				// Only the actor knows whether this is idle input or a live hint.
				// Keep every private completion so EOF also drains newly started turns.
				pending = append(pending, a.SubmitLive(ctx, text))
				fmt.Fprintln(diagnostics, "input queued: "+text)
			}
			fmt.Fprint(diagnostics, "you> ")
		case <-interrupts:
			if err := a.Post(ensemble.Interrupt{}); err != nil {
				return err
			}
			fmt.Fprintln(diagnostics, "interrupt: turn ending; jobs remain available")
		case result := <-reply:
			pending = pending[1:]
			if result.TurnID != 0 && result.TurnID == printedTurn {
				continue
			}
			if result.Err != nil {
				if !errors.Is(result.Err, ensemble.ErrInterrupted) {
					return result.Err
				}
				fmt.Fprintln(diagnostics, "turn interrupted")
			} else if _, err := fmt.Fprintf(out, "\nassistant> %s\n\n", result.Text); err != nil {
				return err
			}
			printedTurn = result.TurnID
			u := a.Usage()
			fmt.Fprintf(diagnostics, "[%d turns | %d input tokens | %d output tokens]\n", len(a.History()), u.Input, u.Output)
		case observation := <-progress:
			if final, ok := observation.(ensemble.PartFinal); ok && final.Part.Type == "tool_result" {
				for _, p := range final.Part.Parts {
					if p.Type == "text" {
						fmt.Fprintln(diagnostics, p.Text)
					}
				}
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}
