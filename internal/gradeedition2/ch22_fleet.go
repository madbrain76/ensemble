package gradeedition2

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/waywardgeek/ensemble/internal/fakevendor"
	"github.com/waywardgeek/ensemble/internal/grade"
)

const fleetRoot = "FLEET_ROOT_22"
const fleetDone = "FLEET_COMPLETED_22"

// Route only external model replies. Every spawned child, tool job, question,
// saved conversation and idle correction belongs to the submitted executable.
func fleetRun(bin, work string, route func([]byte) *fakevendor.Reply, steer func(*mcpCommand) error) (reqs []fakevendor.Recorded, err error) {
	skills := filepath.Join(work, "skills")
	err = recallWrite(skills, map[string]string{"base/SKILL.md": "---\nname: base\ndescription: Terminal fleet exercise\ntype: primary\ntools: spawn_sub_agent join_agents agents_status shutdown_agent send_message_to_parent respond_to_parent_message submit run_command wait_for_job tool_limits write_file read_file\n---\nFinish work with submit. Questions are not submissions.\n"})
	if err != nil {
		return
	}
	provider := fakevendor.NewWithOptions(nil, fakevendor.Options{Route: func(body []byte) *fakevendor.Reply {
		var request map[string]any
		_ = json.Unmarshal(body, &request)
		if len(list(request["tools"])) == 0 {
			return &fakevendor.Reply{Text: "[]"}
		}
		return route(body)
	}})
	defer provider.Close()
	p, e := startMCPCommand(bin, work, []string{"chat"}, []string{"LLM_MODEL=claude-sonnet-5", "LLM_BASE_URL=" + provider.URL(), "EN_SKILLS_DIR=" + skills, "EN_SKILL_DIR=" + skills, "EN_PRIMARY_SKILL=base", "EN_RECALL_DISABLED=1"})
	if e != nil {
		return nil, e
	}
	defer p.close()
	_, err = fmt.Fprintln(p.stdin, fleetRoot)
	if err == nil && steer != nil {
		err = steer(p)
	}
	if err == nil {
		for {
			var line string
			line, err = p.read()
			if err != nil || strings.Contains(line, fleetDone) {
				break
			}
		}
	}
	p.stdin.Close()
	select {
	case exit := <-p.done:
		p.done <- exit
		if err == nil && exit != nil {
			err = fmt.Errorf("CLI exit: %w", exit)
		}
	case <-time.After(5 * time.Second):
		p.cancel()
		exit := <-p.done
		p.done <- exit
		if err == nil {
			err = fmt.Errorf("CLI did not finish after EOF")
		}
	}
	if err != nil {
		err = fmt.Errorf("%w; stderr: %s", err, p.stderr.String())
	}
	return provider.Requests(), err
}

func fleetCall(id, name string, args any) *fakevendor.Reply {
	b, _ := json.Marshal(args)
	r := skillCall(id, name, string(b))
	return &r
}
func fleetHas(body []byte, id string) bool {
	_, _, err := anthropicToolText(body, id)
	return err == nil
}
func fleetFinish() *fakevendor.Reply {
	return fleetCall("finish", "submit", map[string]any{"result": fleetDone})
}

// Result lookup is correlated to the scripted call, never a marker in its
// arguments. Repeated requests may replay the same immutable result.
func fleetResult(c *grade.Check, reqs []fakevendor.Recorded, id string, wantError bool) string {
	for _, r := range reqs {
		if fleetHas(r.Body, id) {
			return skillResult(c, r, id, wantError)
		}
	}
	fail(c, "no actual result for %s", id)
	return ""
}
func fleetTemp(checks []grade.Check, run func(string)) {
	work, err := os.MkdirTemp("", "fleet-grade-")
	if err != nil {
		for i := range checks {
			fail(&checks[i], "%v", err)
		}
		return
	}
	defer os.RemoveAll(work)
	run(work)
}
func fleetFailed(checks []grade.Check, err error) bool {
	if err == nil {
		return false
	}
	for i := range checks {
		fail(&checks[i], "fleet session: %v", err)
	}
	return true
}

func fleetChecks(bin string, checks []grade.Check) {
	fleetTemp(checks[:4], func(work string) {
		reqs, err := fleetRun(bin, work, func(body []byte) *fakevendor.Reply {
			b := string(body)
			if !strings.Contains(b, fleetRoot) {
				payload := "CHILD_A_22"
				if strings.Contains(b, "TASK_B_22") {
					payload = "CHILD_B_22"
				}
				return fleetCall("child-submit", "submit", map[string]any{"result": payload})
			}
			for _, p := range []sandboxProbe{{"spawn-a", "spawn_sub_agent", map[string]any{"agent_id": "worker-a", "task": "TASK_A_22"}}, {"spawn-b", "spawn_sub_agent", map[string]any{"agent_id": "worker-b", "task": "TASK_B_22"}}, {"join-both", "join_agents", map[string]any{"agent_ids": []string{"worker-a", "worker-b"}, "ai_callback_delay": 5}}, {"status-before", "agents_status", map[string]any{}}, {"shutdown", "shutdown_agent", map[string]any{"agent_id": "worker-a"}}, {"status-after", "agents_status", map[string]any{}}} {
				if !fleetHas(body, p.id) {
					r := p.reply()
					return &r
				}
			}
			return fleetFinish()
		}, nil)
		if fleetFailed(checks[:4], err) {
			return
		}
		joined := fleetResult(&checks[0], reqs, "join-both", false)
		if !strings.Contains(joined, "CHILD_A_22") {
			fail(&checks[0], "join lost submitted child payload")
		}
		if !strings.Contains(joined, "CHILD_A_22") || !strings.Contains(joined, "CHILD_B_22") {
			fail(&checks[3], "one correlated join did not contain both submissions: %s", joined)
		}
		before := fleetResult(&checks[2], reqs, "status-before", false)
		after := fleetResult(&checks[2], reqs, "status-after", false)
		if !strings.Contains(before, "worker-a") || !strings.Contains(before, "worker-b") || strings.Contains(after, "worker-a") || !strings.Contains(after, "worker-b") {
			fail(&checks[2], "status did not remove exactly the shut-down child: before=%s after=%s", before, after)
		}
		fleetResult(&checks[2], reqs, "shutdown", false)
		// Validate separate real saves containing their own submitted facts. Merely
		// creating two empty files or observing two paths is not isolation evidence.
		owners := map[string]string{}
		_ = filepath.WalkDir(work, func(path string, d os.DirEntry, e error) error {
			if e != nil || d.IsDir() || d.Name() != "save.json" {
				return nil
			}
			raw, e := os.ReadFile(path)
			if e != nil {
				return nil
			}
			var save map[string]any
			if json.Unmarshal(raw, &save) != nil {
				return nil
			}
			for _, token := range []string{fleetRoot, "TASK_A_22", "TASK_B_22"} {
				if fleetExactText(save, token) {
					owners[token] = filepath.Dir(path)
				}
			}
			return nil
		})
		if owners[fleetRoot] == "" || owners["TASK_A_22"] == "" || owners["TASK_B_22"] == "" || owners[fleetRoot] == owners["TASK_A_22"] || owners[fleetRoot] == owners["TASK_B_22"] || owners["TASK_A_22"] == owners["TASK_B_22"] {
			fail(&checks[1], "distinct valid parent/child saved conversations absent: %v", owners)
		}
	})
	fleetMessaging(bin, &checks[4])
	fleetWake(bin, &checks[5])
	fleetIdle(bin, &checks[6])
}

func fleetMessaging(bin string, c *grade.Check) {
	fleetTemp([]grade.Check{*c}, func(work string) {
		reqs, err := fleetRun(bin, work, func(body []byte) *fakevendor.Reply {
			b := string(body)
			if strings.Contains(b, fleetRoot) {
				if !fleetHas(body, "spawn-mid") {
					return fleetCall("spawn-mid", "spawn_sub_agent", map[string]any{"agent_id": "middle", "task": "MIDDLE_TASK_22"})
				}
				if !fleetHas(body, "join-mid") {
					return fleetCall("join-mid", "join_agents", map[string]any{"agent_ids": []string{"middle"}, "ai_callback_delay": 5})
				}
				return fleetFinish()
			}
			if strings.Contains(b, "MIDDLE_TASK_22") {
				if !fleetHas(body, "spawn-leaf") {
					return fleetCall("spawn-leaf", "spawn_sub_agent", map[string]any{"agent_id": "leaf", "task": "LEAF_TASK_22"})
				}
				if !fleetHas(body, "question-join") {
					return fleetCall("question-join", "join_agents", map[string]any{"agent_ids": []string{"leaf"}, "ai_callback_delay": 5})
				}
				if !fleetHas(body, "answer") {
					return fleetCall("answer", "respond_to_parent_message", map[string]any{"agent_id": "leaf", "response": "ANSWER_22"})
				}
				if !fleetHas(body, "leaf-join") {
					return fleetCall("leaf-join", "join_agents", map[string]any{"agent_ids": []string{"leaf"}, "ai_callback_delay": 5})
				}
				return fleetCall("mid-submit", "submit", map[string]any{"result": "MIDDLE_DONE_22"})
			}
			if !fleetHas(body, "ask") {
				return fleetCall("ask", "send_message_to_parent", map[string]any{"content": "QUESTION_22", "ai_callback_delay": 5})
			}
			return fleetCall("leaf-submit", "submit", map[string]any{"result": "LEAF_DONE_22"})
		}, nil)
		if err != nil {
			fail(c, "three-generation session: %v", err)
			return
		}
		for id, parts := range map[string][]string{"question-join": {"QUESTION_22", "leaf"}, "ask": {"ANSWER_22", "middle"}, "leaf-join": {"LEAF_DONE_22"}, "join-mid": {"MIDDLE_DONE_22"}} {
			text := fleetResult(c, reqs, id, false)
			for _, part := range parts {
				if !strings.Contains(text, part) {
					fail(c, "%s lost %s: %s", id, part, text)
				}
			}
		}
		// Actual human terminal delivery/reply and same-name question isolation retain
		// their accepted built-CLI/source tests; this numeric case observes agent hops.
	})
}

var fleetHandle = regexp.MustCompile(`(?i)\b(?:handle|job)[" :=]+([0-9]+)`)

func fleetWake(bin string, c *grade.Check) {
	fleetTemp([]grade.Check{*c}, func(work string) {
		start, release := filepath.Join(work, "started"), filepath.Join(work, "release")
		command := fmt.Sprintf("printf started > %q; while [ ! -f %q ]; do sleep 0.02; done; printf 'WAKE_RESULT_%%s' $((700+22))", start, release)
		var mu sync.Mutex
		woke := false
		handle := ""
		var defect string
		reqs, err := fleetRun(bin, work, func(body []byte) *fakevendor.Reply {
			if !fleetHas(body, "limits") {
				return fleetCall("limits", "tool_limits", map[string]any{"ai_callback_delay": 60})
			}
			if !fleetHas(body, "slow") {
				return fleetCall("slow", "run_command", map[string]any{"command": command})
			}
			text, _, _ := anthropicToolText(body, "slow")
			mu.Lock()
			defer mu.Unlock()
			if handle == "" {
				m := fleetHandle.FindStringSubmatch(text)
				if len(m) == 2 {
					handle = m[1]
				}
			}
			if strings.Contains(string(body), "STEER_WAKE_22") && !strings.Contains(text, "WAKE_RESULT_722") {
				woke = true
				_ = os.WriteFile(release, []byte("release"), 0600)
			}
			if fleetHas(body, "wait-same") {
				return fleetFinish()
			}
			if handle == "" {
				defect = "no running job handle"
				return fleetFinish()
			}
			var h int
			fmt.Sscan(handle, &h)
			return fleetCall("wait-same", "wait_for_job", map[string]any{"handle": h, "ai_callback_delay": 5})
		}, func(p *mcpCommand) error {
			deadline := time.Now().Add(3 * time.Second)
			for time.Now().Before(deadline) {
				if _, e := os.Stat(start); e == nil {
					_, e = fmt.Fprintln(p.stdin, "STEER_WAKE_22")
					return e
				}
				time.Sleep(10 * time.Millisecond)
			}
			return fmt.Errorf("real command never reached its pending barrier")
		})
		if err != nil {
			fail(c, "message-wake session: %v", err)
			return
		}
		mu.Lock()
		defer mu.Unlock()
		if !woke || defect != "" {
			fail(c, "no causal message wake before real command release: woke=%v %s", woke, defect)
		}
		text := fleetResult(c, reqs, "wait-same", false)
		if !strings.Contains(text, "WAKE_RESULT_722") {
			fail(c, "same handle %s lost completed output: %s", handle, text)
		}
	})
}

func fleetIdle(bin string, c *grade.Check) {
	fleetTemp([]grade.Check{*c}, func(work string) {
		release := filepath.Join(work, "idle-release")
		count := 0
		reqs, err := fleetRun(bin, work, func(body []byte) *fakevendor.Reply {
			if !strings.Contains(string(body), fleetRoot) {
				if !fleetHas(body, "child-hold") {
					return fleetCall("child-hold", "run_command", map[string]any{"command": fmt.Sprintf("while [ ! -f %q ]; do sleep 0.02; done", release), "ai_callback_delay": 60})
				}
				count++
				return &fakevendor.Reply{Text: "I stopped without submitting."}
			}
			for _, p := range []sandboxProbe{{"spawn-idle", "spawn_sub_agent", map[string]any{"agent_id": "drifter", "task": "IDLE_TASK_22"}}, {"premature", "submit", map[string]any{"result": "PREMATURE_22"}}, {"release-idle", "write_file", map[string]any{"path": release, "content": "release"}}, {"idle-join", "join_agents", map[string]any{"agent_ids": []string{"drifter"}, "ai_callback_delay": 5}}, {"stop-idle", "shutdown_agent", map[string]any{"agent_id": "drifter"}}} {
				if !fleetHas(body, p.id) {
					r := p.reply()
					return &r
				}
			}
			return fleetFinish()
		}, nil)
		if err != nil {
			fail(c, "idle session: %v", err)
			return
		}
		fleetResult(c, reqs, "premature", true)
		if count != 4 {
			fail(c, "unsubmitted child received %d completed prose responses, want its initial response plus exactly three wakes", count)
		}
		joined := ""
		for _, r := range reqs {
			if text, _, e := anthropicToolText(r.Body, "idle-join"); e == nil {
				joined = text
				break
			}
		}
		if !strings.Contains(joined, "drifter") || (!strings.Contains(strings.ToLower(joined), "fail") && !strings.Contains(strings.ToLower(joined), "without submit")) {
			fail(c, "join lost actual exhausted-child escalation: %s", joined)
		}
	})
}

// Match stored text values, not a parent's echoed spawn task argument. The save
// remains public JSON; this does not impose a runtime struct or directory name.
func fleetExactText(v any, want string) bool {
	switch x := v.(type) {
	case map[string]any:
		if x["text"] == want {
			return true
		}
		for _, e := range x {
			if fleetExactText(e, want) {
				return true
			}
		}
	case []any:
		for _, e := range x {
			if fleetExactText(e, want) {
				return true
			}
		}
	}
	return false
}
