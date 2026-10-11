// Command interviewer is a standalone agent that interviews another agent
// through the gateway.
//
// It exists to prove the gateway is reusable infrastructure rather than a
// part of Ensemble. This binary shares no code with Ensemble: it imports the
// public client library, holds its own keypair, and reaches its subject only
// by routing. If this works, so does any other program anyone writes.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	gateway "github.com/waywardgeek/ensemble/agent/gateway"
)

// defaultQuestions are the self-improvement prompts the interviewer asks.
//
// They are fixed rather than model-generated so that an interview is
// reproducible and can be checked. Generating follow-up questions from the
// answers is the obvious extension, and it changes nothing about the
// gateway.
var defaultQuestions = []string{
	"What friction points slow you down most in your current workflow?",
	"Which tasks take the most turns or the most tokens?",
	"Where do you get stuck and need human intervention?",
	"What new tool or skill would most improve your work?",
}

func main() {
	var (
		gatewayURL = flag.String("gateway", "ws://127.0.0.1:8093/v1/ws", "gateway websocket URL")
		agentID    = flag.String("agent-id", "interviewer", "this agent's id")
		keyPath    = flag.String("key", "./data/interviewer.key", "path to this agent's private key")
		gatewayKey = flag.String("gateway-key", "", "the gateway's public key (x25519:base64)")
		target     = flag.String("target", "ensemble", "the agent to interview")
		outPath    = flag.String("out", "", "markdown file to write (default interviews/<target>-<date>.md)")
		timeout    = flag.Duration("timeout", 60*time.Second, "how long to wait for each answer")
		printKey   = flag.Bool("print-key", false, "print this agent's public key and exit")
	)
	flag.Parse()

	if err := run(runConfig{
		gatewayURL: *gatewayURL,
		agentID:    *agentID,
		keyPath:    *keyPath,
		gatewayKey: *gatewayKey,
		target:     *target,
		outPath:    *outPath,
		timeout:    *timeout,
		printKey:   *printKey,
	}); err != nil {
		fmt.Fprintf(os.Stderr, "interviewer: %v\n", err)
		os.Exit(1)
	}
}

type runConfig struct {
	gatewayURL string
	agentID    string
	keyPath    string
	gatewayKey string
	target     string
	outPath    string
	timeout    time.Duration
	printKey   bool
}

// answerHandler collects replies routed back from the subject.
//
// The handler does nothing but hand the envelope to the interview loop. Work
// in a Handler blocks the client's receive goroutine, so the rule is to move
// it somewhere else immediately.
type answerHandler struct {
	gateway.NopHandler
	answers chan gateway.Envelope
}

func (h *answerHandler) OnDeliver(env gateway.Envelope) { h.answers <- env }

func run(cfg runConfig) error {
	key, err := gateway.LoadOrCreateKeypair(cfg.keyPath)
	if err != nil {
		return err
	}
	if cfg.printKey {
		fmt.Println(gateway.EncodePublicKey(key.Public))
		return nil
	}
	if cfg.gatewayKey == "" {
		return fmt.Errorf("-gateway-key is required; get it from the gateway's startup log")
	}
	gwPub, err := gateway.ParsePublicKey(cfg.gatewayKey)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	handler := &answerHandler{answers: make(chan gateway.Envelope, 16)}
	client, err := gateway.Connect(ctx, gateway.ClientConfig{
		GatewayURL:    cfg.gatewayURL,
		AgentID:       cfg.agentID,
		StaticKey:     key,
		GatewayPubKey: gwPub,
	}, handler)
	if err != nil {
		return err
	}
	defer func() { _ = client.Close() }()

	// Accept replies from anyone the subject might answer as. Without this
	// the subject's answers would be refused: contact is opt-in in both
	// directions, including for the agent that started the conversation.
	if err := client.DeclareContactable(ctx, []string{"agents", "admin"}); err != nil {
		return fmt.Errorf("declare contactable: %w", err)
	}

	conversationID := "interview-" + time.Now().UTC().Format("20060102-150405")
	transcript := make([]exchange, 0, len(defaultQuestions))

	for _, question := range defaultQuestions {
		answer, err := ask(ctx, client, handler, cfg, conversationID, question)
		if err != nil {
			return err
		}
		transcript = append(transcript, exchange{Question: question, Answer: answer})
		fmt.Printf("interviewer: asked %q\n", question)
	}

	out := cfg.outPath
	if out == "" {
		out = filepath.Join("interviews",
			fmt.Sprintf("%s-%s.md", cfg.target, time.Now().UTC().Format("2006-01-02")))
	}
	if err := writeReport(out, cfg.target, transcript); err != nil {
		return err
	}
	fmt.Printf("interviewer: wrote %s\n", out)
	return nil
}

type exchange struct {
	Question string
	Answer   string
}

// ask sends one question and waits for the matching answer.
func ask(ctx context.Context, client *gateway.Client, handler *answerHandler,
	cfg runConfig, conversationID, question string) (string, error) {

	if err := client.Send(ctx, gateway.OutboundMessage{
		ToAgent:        cfg.target,
		Content:        question,
		ConversationID: conversationID,
		ResponseMode:   gateway.ResponseAsync,
	}); err != nil {
		return "", fmt.Errorf("send question: %w", err)
	}

	select {
	case env := <-handler.answers:
		return stripProvenance(env.Content), nil
	case <-time.After(cfg.timeout):
		return "", fmt.Errorf("timed out after %s waiting for an answer to %q", cfg.timeout, question)
	case <-client.Done():
		return "", fmt.Errorf("connection closed while waiting for an answer")
	}
}

// stripProvenance removes the gateway's header from a delivered message.
//
// The header is for the reader of a prompt, not for a transcript. Leaving it
// in would put routing metadata in the middle of the findings document.
func stripProvenance(content string) string {
	if !strings.HasPrefix(content, "[GATEWAY ") {
		return content
	}
	_, body, found := strings.Cut(content, "\n\n")
	if !found {
		return content
	}
	return body
}

// writeReport renders the interview as markdown.
func writeReport(path, subject string, transcript []exchange) error {
	if dir := filepath.Dir(path); dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}

	var b strings.Builder
	fmt.Fprintf(&b, "# Interview: %s - %s\n\n", subject, time.Now().UTC().Format("2006-01-02"))

	fmt.Fprintf(&b, "## Friction Points\n\n")
	for _, ex := range transcript[:min(2, len(transcript))] {
		fmt.Fprintf(&b, "- %s\n", firstLine(ex.Answer))
	}

	fmt.Fprintf(&b, "\n## Improvement Opportunities\n\n")
	if len(transcript) > 2 {
		for _, ex := range transcript[2:] {
			fmt.Fprintf(&b, "- %s\n", firstLine(ex.Answer))
		}
	}

	fmt.Fprintf(&b, "\n## Raw Transcript\n\n")
	for _, ex := range transcript {
		fmt.Fprintf(&b, "**Q:** %s\n\n", ex.Question)
		fmt.Fprintf(&b, "**A:** %s\n\n", ex.Answer)
	}

	return os.WriteFile(path, []byte(b.String()), 0o644)
}

func firstLine(s string) string {
	line, _, _ := strings.Cut(strings.TrimSpace(s), "\n")
	if line == "" {
		return "(no answer)"
	}
	return line
}
