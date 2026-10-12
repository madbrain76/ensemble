package gradeedition2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/waywardgeek/ensemble/internal/fakevendor"
	"github.com/waywardgeek/ensemble/internal/grade"
)

const managementTools = "spawn_sub_agent join_agents agents_status shutdown_agent send_message_to_parent respond_to_parent_message submit wait_for_agent_change send_message wait_for_agent ask_agent get_submission kill_agent get_sub_agent_status interrupt_agent run_command write_file read_file wait_for_job tool_limits"

// This launches the actual optional application. A free port is not readiness:
// check the launched child's exit channel before trusting the connected socket.
func managementStart(bin, work, base string) (*mcpCommand, *guiPeer, error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, nil, err
	}
	port := fmt.Sprint(l.Addr().(*net.TCPAddr).Port)
	l.Close()
	p, err := startMCPCommand(bin, work, []string{"--port", port}, []string{"LLM_MODEL=claude-sonnet-5", "LLM_BASE_URL=" + base, "EN_SKILLS_DIR=" + filepath.Join(work, "skills"), "EN_SKILL_DIR=" + filepath.Join(work, "skills"), "EN_PRIMARY_SKILL=base", "EN_RECALL_DISABLED=1"})
	if err != nil {
		return nil, nil, err
	}
	deadline := time.Now().Add(4 * time.Second)
	for time.Now().Before(deadline) {
		peer, e := openGUIPeer("ws://127.0.0.1:" + port + "/ws")
		select {
		case exit := <-p.done:
			p.done <- exit
			if peer != nil {
				peer.close()
			}
			p.close()
			return nil, nil, fmt.Errorf("GUI child exited before readiness: %v", exit)
		default:
		}
		if e == nil {
			return p, peer, nil
		}
		time.Sleep(20 * time.Millisecond)
	}
	p.close()
	return nil, nil, fmt.Errorf("GUI readiness deadline: %s", p.stderr.String())
}

func managementEOF(p *mcpCommand) error {
	p.stdin.Close()
	select {
	case err := <-p.done:
		p.done <- err
		if err != nil {
			return fmt.Errorf("GUI exit: %w; %s", err, p.stderr.String())
		}
		return nil
	case <-time.After(4 * time.Second):
		return fmt.Errorf("GUI did not finish normally after EOF")
	}
}

// Selection replay has an existing effective-settings trailer, not a new grader
// protocol. Retain all preceding historical records and the current state frame.
func managementSelect(peer *guiPeer, key string) ([]map[string]any, error) {
	if err := peer.conn.WriteJSON(map[string]string{"type": "select_agent", "key": key}); err != nil {
		return nil, err
	}
	return peer.until(2*time.Second, messageKind("current_settings"))
}
func managementNode(v any, id string) map[string]any {
	m, ok := v.(map[string]any)
	if !ok {
		return nil
	}
	if m["id"] == id {
		return m
	}
	for _, child := range list(m["children"]) {
		if got := managementNode(child, id); got != nil {
			return got
		}
	}
	return nil
}
func managementPhase(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprint(object(v)["phase"])
}
func managementState(frames []map[string]any) string {
	state := ""
	for _, m := range frames {
		if m["type"] == "state_changed" {
			state = fmt.Sprint(m["to"])
		}
	}
	return state
}
func managementContains(frames []map[string]any, text string) bool {
	// Only rendered result/text events qualify, never submitted call arguments.
	for _, m := range frames {
		content := ""
		switch m["type"] {
		case "part_final":
			part := object(m["part"])
			if part["type"] == "text" {
				content, _ = part["text"].(string)
			}
		case "tool_finished":
			content, _ = m["result"].(string)
		case "turn_ended":
			content, _ = m["text"].(string)
		}
		if strings.Contains(content, text) {
			return true
		}
	}
	return false
}

func managementGUI(bin string, checks []grade.Check) {
	fleetTemp([]*grade.Check{&checks[0], &checks[1], &checks[2], &checks[3], &checks[5]}, func(work string) {
		skill := "---\nname: %s\ndescription: Fleet management fixture\ntype: primary\ntools: " + managementTools + "\n---\n%s\n"
		if err := recallWrite(filepath.Join(work, "skills"), map[string]string{"base/SKILL.md": fmt.Sprintf(skill, "base", "Complete tasks using submit."), "peer/SKILL.md": fmt.Sprintf(skill, "peer", "PEER_SKILL_BODY_23. Complete tasks using submit.")}); err != nil {
			reuseFail(checks, err, 0, 1, 2, 3, 5)
			return
		}
		provider := fakevendor.NewWithOptions(nil, fakevendor.Options{Route: func(body []byte) *fakevendor.Reply {
			var r map[string]any
			_ = json.Unmarshal(body, &r)
			if len(list(r["tools"])) == 0 && !fleetHas(body, "peer-submit") && !fleetHas(body, "manage-submit") && !fleetHas(body, "child-submit") && !fleetHas(body, "intruder-submit") {
				return &fakevendor.Reply{Text: "[]"}
			}
			return managementReply(body, work)
		}})
		defer provider.Close()
		entered, release := make(chan struct{}), make(chan struct{})
		var enterOnce, releaseOnce sync.Once
		unpark := func() { releaseOnce.Do(func() { close(release) }) }
		upstream, _ := url.Parse(provider.URL())
		forward := httputil.NewSingleHostReverseProxy(upstream)
		// Hold only the external HTTP request, outside the fake vendor's capture
		// mutex. Actual replay stays observable while the model is genuinely busy.
		boundary := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			raw, _ := io.ReadAll(r.Body)
			r.Body.Close()
			r.Body = io.NopCloser(bytes.NewReader(raw))
			if strings.Contains(string(raw), "PENDING_TASK_23") {
				enterOnce.Do(func() { close(entered) })
				select {
				case <-release:
				case <-r.Context().Done():
					return
				}
			}
			// Child-start acknowledgment comes from the actual shell, not an
			// earlier activity event. Hold the parent's next response outside
			// the capture mutex so the child can receive its own tool call.
			if strings.Contains(string(raw), "MANAGEMENT_ROOT_23") && fleetHas(raw, "spawn-worker") && !fleetHas(raw, "wait-change") {
				deadline := time.NewTimer(2 * time.Second)
				defer deadline.Stop()
				tick := time.NewTicker(10 * time.Millisecond)
				defer tick.Stop()
				for {
					if _, err := os.Stat(filepath.Join(work, "ready")); err == nil {
						break
					}
					select {
					case <-tick.C:
					case <-r.Context().Done():
						return
					case <-deadline.C:
						http.Error(w, "child command never acknowledged startup", http.StatusGatewayTimeout)
						return
					}
				}
			}
			forward.ServeHTTP(w, r)
		}))
		defer boundary.Close()
		defer unpark()
		p, peer, err := managementStart(bin, work, boundary.URL)
		if err != nil {
			reuseFail(checks, err, 0, 1, 2, 3, 5)
			return
		}
		defer p.close()
		defer peer.close()
		peer.send("subscribe", "")
		initial, err := peer.until(2*time.Second, func(m map[string]any) bool { return m["type"] == "agent_tree" || m["type"] == "agents_tree" })
		if err != nil || len(initial) == 0 || len(object(initial[len(initial)-1]["root"])) == 0 {
			reuseFail(checks, fmt.Errorf("initial real tree: %v", err), 0, 1, 2, 3, 5)
			return
		}
		peer.send("prompt", "MANAGEMENT_ROOT_23")
		released := false
		phases := []string{}
		nested := false
		live, err := peer.until(9*time.Second, func(m map[string]any) bool {
			phase := ""
			if m["type"] == "agent_activity" && strings.HasSuffix(fmt.Sprint(m["key"]), "/worker23") {
				phase = managementPhase(m["activity"])
			}
			if m["type"] == "agent_tree" || m["type"] == "agents_tree" {
				if n := managementNode(m["root"], "worker23"); n != nil {
					nested = true
					phase = managementPhase(n["activity"])
				}
			}
			if phase != "" && (len(phases) == 0 || phases[len(phases)-1] != phase) {
				phases = append(phases, phase)
			}
			// Waiting activity precedes handler execution. It only starts this
			// bounded no-progress control; it is not an entry acknowledgment.
			// The child's shell is already held and has no output/state change.
			raw, _ := json.Marshal(m)
			if !released && strings.Contains(string(raw), "wait_for_agent_change") && strings.Contains(string(raw), "waiting") {
				quiet, e := peer.until(200*time.Millisecond, nil)
				if e != nil {
					fail(&checks[3], "held wait observation: %v", e)
				}
				for _, observed := range quiet {
					if observed["type"] == "tool_finished" && observed["call_id"] == "wait-change" {
						fail(&checks[3], "change wait completed while selected child had no progress")
					}
				}
				released = true
				if e = os.WriteFile(filepath.Join(work, "release"), nil, 0600); e != nil {
					fail(&checks[3], "release selected child: %v", e)
				}
			}
			return m["type"] == "turn_ended" && m["replay"] != true && strings.Contains(fmt.Sprint(m["text"]), "MANAGEMENT_DONE_23")
		})
		if err != nil {
			reuseFail(checks, fmt.Errorf("management turn: %v", err), 0, 2, 3)
			_ = os.WriteFile(filepath.Join(work, "release"), nil, 0600)
		}
		managementResults(provider.Requests(), checks)
		// Extra real transitions are allowed; require the promised phases in order,
		// not one implementation-specific number of observations.
		at := 0
		for _, phase := range phases {
			want := []string{"thinking", "tool_call", "idle"}
			if at < len(want) && strings.ReplaceAll(phase, "toolcall", "tool_call") == want[at] {
				at++
			}
		}
		if at != 3 {
			fail(&checks[0], "ordered child thinking/tool/idle phases absent: %v", phases)
		}
		if !released {
			fail(&checks[3], "no actual waiting observation before child release")
		}
		if !nested {
			fail(&checks[5], "no live nested child in GUI tree")
		}
		_ = live
		// Named peer creation is a real client command, not a planted registry file.
		peer.conn.WriteJSON(map[string]string{"type": "create_agent", "name": "named23", "model": "claude-opus-4-6", "skill": "peer"})
		_, err = peer.until(2*time.Second, func(m map[string]any) bool {
			return (m["type"] == "agent_tree" || m["type"] == "agents_tree") && managementNode(m["root"], "named23") != nil
		})
		if err == nil {
			_, err = managementSelect(peer, "named23")
		}
		if err != nil {
			reuseFail(checks, fmt.Errorf("named create/select: %v", err), 1, 5)
			return
		}
		peer.send("prompt", "PEER_TASK_23")
		_, err = peer.until(3*time.Second, func(m map[string]any) bool {
			return m["type"] == "turn_ended" && m["replay"] != true && strings.Contains(fmt.Sprint(m["text"]), "NAMED_CONVERSATION_23")
		})
		if err != nil {
			reuseFail(checks, err, 1, 5)
			return
		}
		replay, err := managementSelect(peer, "named23")
		if err != nil || managementState(replay) != "idle" || !managementContains(replay, "NAMED_CONVERSATION_23") {
			reuseFail(checks, fmt.Errorf("completed selected replay: state=%s history=%t err=%v", managementState(replay), managementContains(replay, "NAMED_CONVERSATION_23"), err), 1, 5)
		}
		// An acknowledgment alone cannot prove constructor inputs took effect.
		// Inspect the actual named Agent request, including its primary instruction.
		selected := false
		for _, r := range provider.Requests() {
			if strings.Contains(string(r.Body), "PEER_TASK_23") && r.JSON()["model"] == "claude-opus-4-6" && strings.Contains(string(r.Body), "PEER_SKILL_BODY_23") {
				selected = true
			}
		}
		if !selected {
			fail(&checks[1], "created model and primary-skill body did not reach actual provider request")
		}
		if err = managementEOF(p); err != nil {
			reuseFail(checks, err, 1, 5)
			return
		}
		// Keep the provider alive across restart: retained words must come from
		// the saved conversation, with no request to reconstruct them.
		count := len(provider.Requests())
		second, view, err := managementStart(bin, work, boundary.URL)
		if err != nil {
			reuseFail(checks, err, 1, 5)
			return
		}
		defer second.close()
		defer view.close()
		replay, err = managementSelect(view, "named23")
		if err != nil || managementState(replay) != "idle" || !managementContains(replay, "NAMED_CONVERSATION_23") || len(provider.Requests()) != count {
			reuseFail(checks, fmt.Errorf("restart did not restore idle named conversation without inference: state=%s err=%v", managementState(replay), err), 1, 5)
			return
		}
		view.send("prompt", "PENDING_TASK_23")
		select {
		case <-entered:
		case <-time.After(3 * time.Second):
			fail(&checks[5], "real pending request never entered provider boundary")
			return
		}
		pending, e := managementSelect(view, "named23")
		if e != nil || managementState(pending) != "in_flight" {
			fail(&checks[5], "genuinely pending selected state=%s err=%v", managementState(pending), e)
		}
		unpark()
		_, e = view.until(3*time.Second, func(m map[string]any) bool {
			return m["type"] == "turn_ended" && m["replay"] != true && strings.Contains(fmt.Sprint(m["text"]), "PENDING_DONE_23")
		})
		if e != nil {
			fail(&checks[5], "pending request did not complete: %v", e)
			return
		}
		completed, e := managementSelect(view, "named23")
		if e != nil || managementState(completed) != "idle" {
			fail(&checks[5], "completed pending request state=%s err=%v", managementState(completed), e)
		}
		view.conn.WriteJSON(map[string]string{"type": "delete_agent", "name": "named23"})
		_, err = view.until(2*time.Second, messageKind("agent_deleted"))
		if err != nil {
			fail(&checks[1], "named deletion: %v", err)
		}
		if err = managementEOF(second); err != nil {
			fail(&checks[1], "%v", err)
			return
		}
		third, last, err := managementStart(bin, work, boundary.URL)
		if err != nil {
			fail(&checks[1], "delete restart: %v", err)
			return
		}
		defer third.close()
		defer last.close()
		last.send("subscribe", "")
		frames, e := last.until(2*time.Second, func(m map[string]any) bool { return m["type"] == "agent_tree" || m["type"] == "agents_tree" })
		if e != nil || len(frames) == 0 || managementNode(frames[len(frames)-1]["root"], "named23") != nil {
			fail(&checks[1], "deleted peer returned after restart: %v", e)
		}
		if err = managementEOF(third); err != nil {
			fail(&checks[1], "%v", err)
		}
	})
}
