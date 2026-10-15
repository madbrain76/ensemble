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
	"time"

	"github.com/waywardgeek/ensemble/internal/fakevendor"
)

type invoiceOpts struct {
	dir, vendor, model string
	prompts            []string
	replies            []fakevendor.Reply
}
type invoiceOut struct {
	reqs   []fakevendor.Recorded
	events []map[string]json.RawMessage
	meters []map[string]any
	fatal  string
}

// Each snapshot comes from a fresh subscription after a completed turn. This
// observes the live meter without assuming its broadcast precedes turn_ended,
// or using the original arbitrary post-response sleep to guess at completion.
// Keep these subscribers alive through shutdown: mid-turn disconnect is a
// separately reviewed lifetime concern, not a prerequisite to meter arithmetic.
func invoiceSnapshot(url string) (map[string]any, *guiPeer, error) {
	p, err := openGUIPeer(url)
	if err != nil {
		return nil, nil, err
	}

	if err = p.send("subscribe", ""); err != nil {
		return nil, p, err
	}
	if _, err = readSettings(p, "current_settings"); err != nil {
		return nil, p, err
	}
	messages, err := p.until(5*time.Second, messageKind("usage"))
	if err != nil {
		return nil, p, err
	}
	return messages[len(messages)-1], p, nil
}

func invoiceLaunch(bin, skillsDir, guiDir string, o invoiceOpts) (out invoiceOut) {
	srv := fakevendor.New(o.replies)
	defer srv.Close()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		out.fatal = err.Error()
		return
	}
	port := fmt.Sprint(listener.Addr().(*net.TCPAddr).Port)
	listener.Close()
	url := "ws://localhost:" + port + "/ws"
	cmd := exec.Command(bin, "--port", port, "--gui-dir", guiDir)
	cmd.Dir = o.dir
	cmd.Env = append(os.Environ(), "LLM_VENDOR="+o.vendor, "LLM_MODEL="+o.model, "LLM_RESOLVED_MODEL=", "LLM_API_KEY=fake-key", "LLM_BASE_URL="+srv.URL(), "EN_DISABLE_STREAMING=1", "EN_SKILLS_DIR="+skillsDir, "EN_PRIMARY_SKILL=base")
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
		out.reqs = srv.Requests()
		out.events, err = recallReadEvents(o.dir)
		if err != nil && out.fatal == "" {
			out.fatal = err.Error()
		}
		if out.fatal != "" {
			out.fatal += ": " + stderr.String()
		}
	}()
	var peer *guiPeer
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		peer, err = openGUIPeer(url)
		if err == nil {
			break
		}
		select {
		case e := <-done:
			finished = true
			out.fatal = fmt.Sprintf("GUI exited before ready: %v", e)
			return
		default:
		}
		time.Sleep(20 * time.Millisecond)
	}
	if err != nil {
		out.fatal = "GUI socket did not start: " + err.Error()
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
	// The initial frame is measured even when the save contains earlier usage.
	initial, snapshot, err := invoiceSnapshot(url)
	if snapshot != nil {
		defer snapshot.close()
	}
	if err != nil {
		out.fatal = err.Error()
		return
	}
	out.meters = append(out.meters, initial)
	for i, prompt := range o.prompts {
		if err = json.NewEncoder(stdin).Encode(map[string]string{"kind": "prompt", "text": prompt}); err != nil {
			out.fatal = err.Error()
			return
		}
		messages, e := peer.until(10*time.Second, messageKind("turn_ended"))
		if e != nil {
			out.fatal = fmt.Sprintf("prompt %d completion: %v", i+1, e)
			return
		}
		if problem, ok := messages[len(messages)-1]["error"].(string); ok && problem != "" {
			out.fatal = "turn failed: " + problem
			return
		}
		meter, snapshot, e := invoiceSnapshot(url)
		if snapshot != nil {
			defer snapshot.close()
		}
		if e != nil {
			out.fatal = e.Error()
			return
		}
		out.meters = append(out.meters, meter)
		if n := len(srv.Requests()); n != i+1 {
			out.fatal = fmt.Sprintf("completed turn %d made %d requests", i+1, n)
			return
		}
	}
	stdin.Close()
	select {
	case err = <-done:
		finished = true
		if err != nil {
			out.fatal = "nonzero GUI exit after EOF: " + err.Error()
		}
	case <-time.After(5 * time.Second):
		out.fatal = "GUI did not save and exit after EOF"
	}
	return
}

// Artifact paths belong to the submitted client, not the core library. They
// are explicit grader fixture inputs and never become application flags.
func invoiceRead(dir, path string) ([]byte, error) {
	return os.ReadFile(filepath.Join(dir, filepath.FromSlash(path)))
}
