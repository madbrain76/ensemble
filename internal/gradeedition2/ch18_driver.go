package gradeedition2

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/waywardgeek/ensemble/internal/grade"
)

type subscriptionOpts struct {
	prompts         []string
	oauth           grade.Edition2OAuthOptions
	responses       grade.Edition2ResponsesOptions
	skipLogin, gate bool
}
type subscriptionOut struct {
	oauth                          grade.Edition2OAuthObservations
	responses                      grade.Edition2ResponsesObservations
	messages                       []map[string]any
	errors                         []string
	loginOutput, output, artifacts string
	loginErr, exitErr, fatal       string
	gated                          bool
}

// The proxy forwards the original peer's real SSE frames. It withholds the
// remainder after the first summary delta until the client displays that delta.
// Mere summary-before-text ordering after a buffered response cannot open it.
type subscriptionGate struct {
	mu             sync.Mutex
	active, opened bool
	release        chan struct{}
	once           sync.Once
	delivered      string
}

func (g *subscriptionGate) observe(m map[string]any) {
	if m["type"] != "part_delta" || m["kind"] != "reasoning_summary" {
		return
	}
	g.mu.Lock()
	if g.active {
		chunk, _ := m["chunk"].(string)
		g.delivered += chunk
		// Delivery may split a provider delta across several observations.
		if g.delivered == "Reading " {
			g.opened = true
			g.once.Do(func() { close(g.release) })
		}
	}
	g.mu.Unlock()
}
func (g *subscriptionGate) proxy(target string, gated bool) *httptest.Server {
	var first sync.Once
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "request body", 400)
			return
		}
		var request map[string]any
		json.Unmarshal(body, &request)
		model, _ := request["model"].(string)
		req, err := http.NewRequestWithContext(r.Context(), r.Method, target+r.URL.RequestURI(), bytes.NewReader(body))
		if err != nil {
			http.Error(w, "request construction", 500)
			return
		}
		req.Header = r.Header.Clone()
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, "peer request failed", 502)
			return
		}
		defer resp.Body.Close()
		for k, v := range resp.Header {
			w.Header()[k] = v
		}
		// Retain the original absent-Content-Type case instead of letting net/http
		// invent a content type when the stream's first bytes are written.
		if resp.Header.Get("Content-Type") == "" {
			w.Header()["Content-Type"] = nil
		}
		w.WriteHeader(resp.StatusCode)
		scan := bufio.NewScanner(resp.Body)
		var frame strings.Builder
		for scan.Scan() {
			line := scan.Text()
			// Captured real terminal responses identify their model. The original
			// fake omitted it, which made valid provenance-aware replay impossible
			// for clients that trust response metadata. Fill only that absent fact.
			if strings.HasPrefix(line, "data: ") {
				var event map[string]any
				if json.Unmarshal([]byte(line[6:]), &event) == nil && event["type"] == "response.completed" {
					response, _ := event["response"].(map[string]any)
					if _, present := response["model"]; response != nil && !present {
						response["model"] = model
						data, _ := json.Marshal(event)
						line = "data: " + string(data)
					}
				}
			}
			frame.WriteString(line + "\n")
			if line != "" {
				continue
			}
			data := frame.String()
			frame.Reset()
			hold := false
			if gated && strings.Contains(data, "event: response.reasoning_summary_text.delta\n") {
				first.Do(func() { hold = true })
			}
			if hold {
				g.mu.Lock()
				g.active = true
				g.mu.Unlock()
			}
			io.WriteString(w, data)
			w.(http.Flusher).Flush()
			if hold {
				select {
				case <-g.release:
				case <-r.Context().Done():
				case <-time.After(3 * time.Second):
				}
				g.mu.Lock()
				g.active = false
				g.mu.Unlock()
			}
		}
		// Ordinary JSON error bodies have no blank SSE delimiter.
		io.WriteString(w, frame.String())
	}))
}

func subscriptionLogin(bin, dir string, env []string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin, "auth", "login", "--print-url")
	cmd.Dir = dir
	cmd.Env = env
	var output subscriptionBuffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Start(); err != nil {
		return "", err
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	visited := false
	for {
		select {
		case err := <-done:
			return output.String(), err
		case <-time.After(10 * time.Millisecond):
			if visited {
				continue
			}
			for _, field := range strings.Fields(output.String()) {
				raw := strings.Trim(field, "\"'<>()[],")
				if !strings.HasPrefix(raw, "http") || !strings.Contains(raw, "/authorize") {
					continue
				}
				visited = true
				client := http.Client{Timeout: 5 * time.Second}
				resp, err := client.Get(raw)
				if err == nil {
					io.Copy(io.Discard, resp.Body)
					resp.Body.Close()
				}
				// The callback may close its listener after receiving the code,
				// before the browser's response finishes. The bounded login process
				// and observed token exchange determine authentication success.
				break
			}
		}
	}
}

func subscriptionLaunch(bin, skills, gui, dir string, o subscriptionOpts) (out subscriptionOut) {
	auth := grade.NewEdition2OAuthPeer(o.oauth)
	defer auth.Close()
	server := grade.NewEdition2ResponsesPeer(o.responses)
	defer server.Close()
	gate := &subscriptionGate{release: make(chan struct{})}
	proxy := gate.proxy(server.URL(), o.gate)
	defer proxy.Close()
	endpoint := proxy.URL
	defer func() {
		out.oauth = auth.Observations()
		out.responses = server.Observations()
		gate.mu.Lock()
		out.gated = gate.opened
		gate.mu.Unlock()
	}()
	home := filepath.Join(dir, "profile")
	os.MkdirAll(home, 0700)
	env := append(os.Environ(), "HOME="+home, "EN_HOME="+home, "OPENAI_OIDC_ISSUER="+auth.URL(), "LLM_BASE_URL="+endpoint, "LLM_VENDOR=openai", "LLM_MODEL=gpt-ch19-course", "LLM_RESOLVED_MODEL=", "LLM_API_KEY="+grade.Edition2APIKey, "EN_BROWSER_CMD=true", "EN_SKILLS_DIR="+skills, "EN_PRIMARY_SKILL=base", "EN_DISABLE_STREAMING=", "CH02_LOG="+filepath.Join(dir, "events.jsonl"))
	if !o.skipLogin {
		var err error
		out.loginOutput, err = subscriptionLogin(bin, dir, env)
		if err != nil {
			out.loginErr = err.Error()
			return
		}
	}
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		out.fatal = err.Error()
		return
	}
	port := fmt.Sprint(listener.Addr().(*net.TCPAddr).Port)
	listener.Close()
	cmd := exec.Command(bin, "--port", port, "--gui-dir", gui, "--save", filepath.Join(dir, "save.json"))
	cmd.Dir = dir
	cmd.Env = env
	var output subscriptionBuffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	stdin, err := cmd.StdinPipe()
	if err != nil {
		out.fatal = err.Error()
		return
	}
	if err = cmd.Start(); err != nil {
		stdin.Close()
		out.fatal = err.Error()
		return
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	finished := false
	defer func() {
		stdin.Close()
		if !finished {
			cmd.Process.Kill()
			<-done
		}
		out.output = output.String()
		// The private credential directory is deliberately excluded. A protected
		// token store is permitted; a save, event log or request capture is not.
		filepath.WalkDir(dir, func(path string, d os.DirEntry, e error) error {
			if e != nil {
				return nil
			}
			if path == home {
				return filepath.SkipDir
			}
			if d.IsDir() {
				return nil
			}
			data, e := os.ReadFile(path)
			if e == nil {
				out.artifacts += string(data)
			}
			return nil
		})
		if out.fatal != "" {
			out.fatal += ": " + out.output
		}
	}()
	var peer *guiPeer
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		select {
		case e := <-done:
			finished = true
			out.fatal = fmt.Sprintf("GUI exited before ready: %v", e)
			return
		default:
		}
		peer, err = openGUIPeer("ws://127.0.0.1:" + port + "/ws")
		if err == nil {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	if err != nil {
		out.fatal = "GUI did not start: " + err.Error()
		return
	}
	defer peer.close()
	if err = peer.send("subscribe", ""); err != nil {
		out.fatal = err.Error()
		return
	}
	if _, err = readSettings(peer, "current_settings"); err != nil {
		out.fatal = err.Error()
		return
	}
	for i, prompt := range o.prompts {
		if err = json.NewEncoder(stdin).Encode(map[string]string{"kind": "prompt", "text": prompt}); err != nil {
			out.fatal = err.Error()
			return
		}
		messages, e := peer.until(8*time.Second, func(m map[string]any) bool {
			if o.gate {
				gate.observe(m)
			}
			return m["type"] == "turn_ended"
		})
		out.messages = append(out.messages, messages...)
		if e != nil {
			out.fatal = fmt.Sprintf("turn %d: %v", i+1, e)
			return
		}
		problem, _ := messages[len(messages)-1]["error"].(string)
		out.errors = append(out.errors, problem)
	}
	// Keep the subscriber connected until the child exits: unrelated legacy
	// disconnect races must not determine these credential/transport observations.
	stdin.Close()
	select {
	case err = <-done:
		finished = true
		if err != nil {
			// A reported authentication/billing refusal may end a client with a
			// normal nonzero status. Positive sessions still require zero; signals
			// and unexplained startup/lifetime failures never establish refusal.
			if exit, ok := err.(*exec.ExitError); ok && exit.ExitCode() >= 0 {
				out.exitErr = err.Error()
			} else {
				out.fatal = "GUI exit after EOF: " + err.Error()
			}
		}
	case <-time.After(5 * time.Second):
		out.fatal = "GUI did not exit after EOF"
	}
	return
}

// Login prints while its child is still running; readers need a stable copy.
type subscriptionBuffer struct {
	mu   sync.Mutex
	text strings.Builder
}

func (b *subscriptionBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.text.Write(p)
}
func (b *subscriptionBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.text.String()
}
