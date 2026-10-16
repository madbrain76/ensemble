// Package gradeedition3 adds Edition 3 checks without changing historical graders.
package gradeedition3

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"sync/atomic"
	"time"

	"github.com/waywardgeek/ensemble/internal/fakeanthropic"
	"github.com/waywardgeek/ensemble/internal/grade"
	"github.com/waywardgeek/ensemble/internal/gradeedition2"
)

// This exceeds the observed historical 60-second client timeout. It is not a
// suggested application deadline, nor proof against every possible deadline.
const generationDelay = 62 * time.Second

// Chapter1 preserves all Edition 2 points and adds two required, zero-point
// lifetime checks. Score 100 alone is insufficient: callers must check Passed.
// The two extra real-process probes take approximately two minutes in total.
func Chapter1(bin string) (grade.Report, error) {
	res, err := grade.Run(bin)
	if err != nil {
		return grade.Report{}, err
	}
	checks := gradeedition2.Chapter1(res).Checks
	for _, body := range []bool{false, true} {
		check, err := checkGenerationLifetime(bin, body, generationDelay)
		if err != nil {
			return grade.Report{}, err
		}
		checks = append(checks, check)
	}
	return grade.NewTitledReport("Chapter 1 — Edition 3 conversation and request lifetime", checks, res.Stderr), nil
}

// Each phase uses a new process and provider. Keeping body delay separate from
// header delay catches clients that return from Do correctly but stop reading.
func checkGenerationLifetime(bin string, body bool, delay time.Duration) (grade.Check, error) {
	phase := "headers"
	if body {
		phase = "body"
	}
	check := grade.Check{
		ID:     "generation_" + phase,
		Title:  "Completes generation with delayed response " + phase,
		Passed: true,
	}
	const answer = "delayed generation completed"
	const response = `{"id":"msg_delayed","type":"message","role":"assistant","content":[{"type":"text","text":"` + answer + `"}],"stop_reason":"end_turn","usage":{"input_tokens":7,"output_tokens":3}}`
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		// Drain the request before flushing: net/http must not confuse an unread
		// request body with the response-body lifetime being exercised here.
		_, _ = io.Copy(io.Discard, r.Body)
		_ = r.Body.Close()
		w.Header().Set("Content-Type", "application/json")
		if body {
			w.WriteHeader(http.StatusOK)
			w.(http.Flusher).Flush()
		}
		timer := time.NewTimer(delay)
		defer timer.Stop()
		select {
		case <-timer.C:
			_, _ = io.WriteString(w, response)
		case <-r.Context().Done():
			// A timed-out or killed client must not strand the fake's handler.
		}
	}))
	defer server.Close()

	dir, err := os.MkdirTemp("", "edition3-ch01-lifetime-*")
	if err != nil {
		return grade.Check{}, err
	}
	defer os.RemoveAll(dir)
	// This is only the grader's process watchdog. The student's request is
	// allowed to wait through the entire delayed response, then exit normally.
	ctx, cancel := context.WithTimeout(context.Background(), delay+grade.RoundTimeout+grade.ExitTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		"ANTHROPIC_BASE_URL="+server.URL,
		"ANTHROPIC_API_URL="+server.URL,
		"ANTHROPIC_API_KEY="+fakeanthropic.ExpectedAPIKey,
		"ANTHROPIC_MODEL="+fakeanthropic.ExpectedModel,
		"LLM_MODEL="+fakeanthropic.ExpectedModel,
	)
	cmd.Stdin = strings.NewReader("{\"user\":\"Please complete this delayed generation.\"}\n")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	// Do not echo submission diagnostics: they may contain credentials even
	// though this probe supplies only the fake provider's test key.
	cmd.Stderr = io.Discard
	cmd.WaitDelay = grade.ExitTimeout
	if err := cmd.Start(); err != nil {
		return grade.Check{}, fmt.Errorf("starting %s lifetime probe: %w", phase, err)
	}
	runErr := cmd.Wait()
	fail := func(detail string) {
		check.Passed = false
		check.Details = append(check.Details, detail)
	}
	if ctx.Err() != nil {
		fail("grader watchdog expired after allowing the full provider delay and exit grace")
	} else if runErr != nil {
		fail("submission did not exit successfully after the delayed response")
	}
	if got := requests.Load(); got != 1 {
		fail(fmt.Sprintf("delayed provider received %d requests, want exactly one", got))
	}

	// The ordinary run already grades the complete protocol. This probe adds
	// only the observable lifetime consequence: the eventual reply and usage
	// must survive the delay, rather than being replaced by a timeout error.
	var answers, usages int
	scanner := bufio.NewScanner(&stdout)
	for scanner.Scan() {
		var line struct {
			Assistant *string `json:"assistant"`
			Usage     *struct {
				Input  int `json:"input"`
				Output int `json:"output"`
			} `json:"usage"`
		}
		if err := json.Unmarshal(scanner.Bytes(), &line); err != nil {
			fail("delayed run emitted a non-JSON protocol line")
			continue
		}
		if line.Assistant != nil && *line.Assistant == answer {
			answers++
		} else if line.Usage != nil && line.Usage.Input == 7 && line.Usage.Output == 3 {
			usages++
		} else {
			fail("delayed run emitted an unexpected reply or usage line")
		}
	}
	if scanner.Err() != nil {
		fail("could not read delayed run protocol output")
	}
	if answers != 1 || usages != 1 {
		fail(fmt.Sprintf("after %s delay of %s: got %d complete answers and %d matching usage lines, want one each; inspect whole-request and response-read deadlines", phase, delay, answers, usages))
	}
	if check.Passed {
		check.Details = []string{fmt.Sprintf("received the complete answer and usage after %s %s delay; longer arbitrary deadlines still require source review", delay, phase)}
	}
	return check, nil
}
