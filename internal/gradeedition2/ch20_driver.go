package gradeedition2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/waywardgeek/ensemble/internal/fakevendor"
)

func webSkills(work, endpoint string) (string, error) {
	dir := filepath.Join(work, "skills")
	err := recallWrite(dir, map[string]string{
		"base/SKILL.md":       "---\nname: base\ndescription: Web transport fixture\ntype: primary\ntools: agent_status wait_for_job load_skill unload_skill\nloadable-skills: web-search\n---\nUse the granted web tools when asked.\n",
		"web-search/SKILL.md": "---\nname: web-search\ndescription: Course HTTP search and fetch\ntype: loadable\ntools: firecrawl_search firecrawl_scrape\nmcp_servers:\n  - name: web\n    transport: url\n    url: " + endpoint + "\n    auth-env: EDITION2_WEB_TOKEN\n---\nSearch for a page and fetch it as evidence.\n",
	})
	return dir, err
}

// The final-only, uniquely named answer is an actual terminal observation. It
// precedes EOF; closing stdin immediately after sending the prompt would instead
// test interruption. Exit status is checked even when partial artifacts exist.
func webLaunch(bin, work, skills, token string, replies []fakevendor.Reply, observe func()) ([]fakevendor.Recorded, string, error) {
	provider := fakevendor.NewWithOptions(replies, fakevendor.Options{Route: func([]byte) *fakevendor.Reply {
		if observe != nil {
			observe()
		}
		return nil
	}})
	defer provider.Close()
	p, err := startMCPCommand(bin, work, []string{"chat"}, []string{"LLM_MODEL=claude-sonnet-5", "LLM_BASE_URL=" + provider.URL(), "EN_SKILLS_DIR=" + skills, "EN_PRIMARY_SKILL=base", "EDITION2_WEB_TOKEN=" + token})
	if err != nil {
		return nil, "", err
	}
	defer p.close()
	if _, err = fmt.Fprintln(p.stdin, "Exercise the web transport."); err == nil {
		for {
			var line string
			line, err = p.read()
			if err != nil || strings.Contains(line, "WEB_EXERCISE_COMPLETE_891") {
				break
			}
		}
	}
	p.stdin.Close()
	select {
	case exit := <-p.done:
		p.done <- exit // shared cleanup still owns its one join
		if exit != nil && err == nil {
			err = fmt.Errorf("CLI exit: %w", exit)
		}
	case <-time.After(5 * time.Second):
		if err == nil {
			err = fmt.Errorf("CLI did not exit after EOF")
		}
		p.cancel()
		exit := <-p.done
		p.done <- exit
	}
	return provider.Requests(), p.stderr.String(), err
}

type statusRun struct {
	reqs          []fakevendor.Recorded
	before, after []map[string]any
	fatal         string
}

var statusModels = []string{"claude-sonnet-5", "claude-opus-4-6", "claude-sonnet-5"}
var statusUsage = []fakevendor.Canonical{
	{Input: 10000, Output: 1000, CacheRead: 2000, CacheWrite: 3000},
	{Input: 2000, Output: 200, CacheRead: 3000, CacheWrite: 1000},
	{Input: 1000, Output: 100, CacheRead: 500, CacheWrite: 250},
}

// Exercise the shipped primary, not a fabricated status-only skill that could
// conceal a dead shipped feature. The settings acknowledgment orders selection;
// the next captured provider model independently verifies its effect.
func statusLaunch(bin, core string) (out statusRun) {
	var replies []fakevendor.Reply
	for i, usage := range statusUsage {
		replies = append(replies, fakevendor.Reply{ToolID: fmt.Sprintf("status-%d", i), ToolName: "agent_status", ToolArgs: `{}`, Usage: usage}, fakevendor.Reply{Text: fmt.Sprintf("STATUS_TURN_%d_COMPLETE", i)})
	}
	provider := fakevendor.New(replies)
	defer provider.Close()
	work, err := os.MkdirTemp("", "status-grade-")
	if err != nil {
		out.fatal = err.Error()
		return
	}
	defer os.RemoveAll(work)
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		out.fatal = err.Error()
		return
	}
	port := fmt.Sprint(listener.Addr().(*net.TCPAddr).Port)
	listener.Close()
	assets := filepath.Join(work, "web")
	os.MkdirAll(assets, 0755)
	skills, _ := filepath.Abs(filepath.Join(core, "skills"))
	cmd := exec.Command(bin, "--port", port, "--gui-dir", assets)
	cmd.Dir = work
	cmd.Env = append(os.Environ(), "LLM_VENDOR=anthropic", "LLM_MODEL="+statusModels[0], "LLM_RESOLVED_MODEL=", "LLM_API_KEY=fake-key", "LLM_BASE_URL="+provider.URL(), "EN_DISABLE_STREAMING=1", "EN_SKILLS_DIR="+skills, "EN_PRIMARY_SKILL=ensemble")
	var stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = io.Discard, &stderr
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
		out.reqs = provider.Requests()
		if out.fatal != "" {
			out.fatal += ": " + stderr.String()
		}
	}()
	url := "ws://127.0.0.1:" + port + "/ws"
	var peer *guiPeer
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		peer, err = openGUIPeer(url)
		select {
		case e := <-done:
			finished = true
			if peer != nil {
				peer.close()
			}
			out.fatal = fmt.Sprintf("GUI exited before ready: %v", e)
			return
		default:
		}
		if err == nil {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	if err != nil {
		out.fatal = "GUI readiness: " + err.Error()
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
	// Keep observation connections alive until the child exits; subscriber
	// disconnect races have separate historical proof and are not meter fixtures.
	for i, model := range statusModels {
		if i > 0 {
			err = peer.conn.WriteJSON(map[string]any{"type": "update_settings", "settings": map[string]string{"model": model}})
			var settings map[string]any
			if err == nil {
				settings, err = readSettings(peer, "settings_changed")
			}
			if err != nil || settings["model"] != model {
				out.fatal = fmt.Sprintf("model selection %s: %v %v", model, settings, err)
				return
			}
		}
		before, p, e := invoiceSnapshot(url)
		if p != nil {
			defer p.close()
		}
		if e != nil {
			out.fatal = e.Error()
			return
		}
		out.before = append(out.before, before)
		if err = json.NewEncoder(stdin).Encode(map[string]string{"kind": "prompt", "text": fmt.Sprintf("Report status turn %d.", i)}); err != nil {
			out.fatal = err.Error()
			return
		}
		messages, e := peer.until(10*time.Second, messageKind("turn_ended"))
		if e != nil {
			out.fatal = e.Error()
			return
		}
		if problem, _ := messages[len(messages)-1]["error"].(string); problem != "" {
			out.fatal = problem
			return
		}
		if len(provider.Requests()) != 2*(i+1) {
			out.fatal = "status turn did not make exactly a tool request and continuation"
			return
		}
		after, p, e := invoiceSnapshot(url)
		if p != nil {
			defer p.close()
		}
		if e != nil {
			out.fatal = e.Error()
			return
		}
		out.after = append(out.after, after)
	}
	stdin.Close()
	select {
	case err = <-done:
		finished = true
		if err != nil {
			out.fatal = "GUI EOF exit: " + err.Error()
		}
	case <-time.After(5 * time.Second):
		out.fatal = "GUI did not exit after EOF"
	}
	return
}
