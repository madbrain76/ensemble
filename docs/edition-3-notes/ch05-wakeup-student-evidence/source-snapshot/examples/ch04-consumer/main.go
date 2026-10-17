// The consumer is an independent application using only Ensemble's public API.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"ensemble"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
func run() (result error) {
	vendor := ensemble.Anthropic
	switch os.Getenv("LLM_VENDOR") {
	case "", "anthropic":
	case "openai":
		vendor = ensemble.OpenAI
	case "gemini":
		vendor = ensemble.Gemini
	default:
		return fmt.Errorf("unsupported LLM_VENDOR")
	}
	base := os.Getenv("LLM_BASE_URL")
	if base == "" {
		switch vendor {
		case ensemble.Anthropic:
			base = "https://api.anthropic.com"
		case ensemble.OpenAI:
			base = "https://api.openai.com"
		case ensemble.Gemini:
			base = "https://generativelanguage.googleapis.com"
		}
	}
	root := ensemble.New(log.New(os.Stderr, "", log.LstdFlags))
	a, err := root.NewAgent(ensemble.Config{Vendor: vendor, BaseURL: base, Model: os.Getenv("LLM_MODEL"), APIKey: os.Getenv("LLM_API_KEY"), System: "Use shout to uppercase the requested text. Explain what you will do and why in visible prose before every tool call.", MaxTokens: 2048, EnableTools: true, RequireToolExplanation: os.Getenv("LLM_REQUIRE_TOOL_EXPLANATION") == "true", Workspace: os.Getenv("ENSEMBLE_WORKSPACE"), DataDir: os.Getenv("ENSEMBLE_DATA_DIR")})
	if err != nil {
		return err
	}
	defer func() {
		if err := root.Shutdown(); result == nil {
			result = err
		}
		path := os.Getenv("CH02_LOG")
		if path == "" {
			path = "consumer.jsonl"
		}
		file, err := os.Create(path)
		if err != nil {
			if result == nil {
				result = err
			}
			return
		}
		err = a.Log().Dump(file)
		closeErr := file.Close()
		if result == nil {
			result = err
		}
		if result == nil {
			result = closeErr
		}
	}()
	err = a.RegisterTool(ensemble.ToolDefinition{Name: "shout", Description: "Uppercase the supplied text.", Parameters: json.RawMessage(`{"type":"object","properties":{"text":{"type":"string"}},"required":["text"],"additionalProperties":false}`)}, func(call ensemble.ToolContext, raw json.RawMessage) (string, error) {
		var args struct {
			Text *string `json:"text"`
		}
		decoder := json.NewDecoder(strings.NewReader(string(raw)))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&args); err != nil {
			return "", err
		}
		if args.Text == nil {
			return "", fmt.Errorf("text is required")
		}
		if call.Job == nil || call.Job.Parent() != call.Agent {
			return "", fmt.Errorf("missing actual job owner")
		}
		call.Agent.Parent().Logger().Print("shout called")
		return strings.ToUpper(*args.Text), nil
	})
	if err != nil {
		return err
	}
	answer, err := a.Ask(context.Background(), "Call shout with text hello framework, then report its result.")
	if err != nil {
		return err
	}
	fmt.Println(answer)
	return nil
}
