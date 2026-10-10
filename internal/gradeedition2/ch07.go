package gradeedition2

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/waywardgeek/ensemble/internal/grade"
)

// Chapter7 grades the optional GUI and the real core through their own entry
// points. Browser/speech quality and causal active-partial replay remain explicit
// source/integration obligations rather than claims made by a WebSocket score.
func Chapter7(bin, core, gui, guiModel, consumer, workflow, exercise, streamModel string, example ConsumerCase, media map[string]string) (grade.Report, error) {
	checks := guiChecks(core, gui, guiModel)
	previous, err := Chapter6(bin, core, consumer, workflow, exercise, streamModel, example, media)
	if err != nil {
		return grade.Report{}, err
	}
	parity := grade.Check{ID: "ch6-parity", Title: "All corrected Chapter 6 checks pass", Points: 20, Earned: 20, Passed: true}
	for _, c := range previous.Checks {
		if !c.Passed {
			fail(&parity, "%s: %s", c.ID, strings.Join(c.Details, "; "))
		}
	}
	return grade.NewTitledReport("Edition 2, Chapter 7 — Everything Is an Artifact", append(checks, parity), ""), nil
}

func guiChecks(core, dir, model string) []grade.Check {
	// The original launcher recursively sends this server into finite CLI tests.
	// Keep its four useful evaluators, but collect their evidence independently.
	r := &grade.Ch8Result{}
	bin, cleanup, err := grade.Build(dir)
	if err != nil {
		r.BuildErr = err.Error()
		return grade.Ch8Evaluate(r)[:4]
	}
	defer cleanup()
	r.BuildOK = true
	wire, replay, trace, err := driveGUI(bin, filepath.Join(dir, "web", "gui"), model)
	r.Messages = oldMessages(wire)
	r.ReplayMessages = oldMessages(replay)
	r.GuiLogContent = trace
	if err != nil {
		r.HelpersErr = err.Error()
		r.ReplayErr = err.Error()
	}
	for _, m := range replay {
		if m["type"] == "part_final" {
			for _, t := range replay {
				if t["type"] == "tool_finished" {
					r.ReplayOK = true
				}
			}
		}
	}
	pause, pauseErr := driveGUIPause(bin, filepath.Join(dir, "web", "gui"), model)
	r.PauseMessages = oldMessages(pause)
	r.PauseOK = pauseErr == nil
	if pauseErr != nil {
		r.PauseErr = pauseErr.Error()
	}
	checks := grade.Ch8Evaluate(r)[:4]
	if err == nil {
		checkGUIWire(&checks[0], wire)
		checkGUIReplay(&checks[1], wire, replay)
		checkGUITrace(&checks[2], trace)
	}
	// Preserve observed wire results even when a same-module legacy application
	// violates the binding optional-client boundary. Inspect the actual executable.
	target := "."
	if info, e := os.Stat(filepath.Join(dir, "cmd")); e == nil && info.IsDir() {
		target = "./cmd"
	}
	if e := publicDependencyTarget(core, dir, target); e != nil {
		wirePassed := checks[0].Passed
		fail(&checks[0], "GUI public dependency: %v", e)
		checks[0].Details = append(checks[0].Details, fmt.Sprintf("observed wire checks passed before dependency assessment: %t", wirePassed))
	}
	checks[1].Details = append(checks[1].Details, "completed-content/identity replay checked; active partials, handoff and fixed window retain required source/integration proof")
	return checks
}

func oldMessages(messages []map[string]any) []grade.Ch8WsMsg {
	var out []grade.Ch8WsMsg
	for _, m := range messages {
		data, _ := json.Marshal(m)
		var v grade.Ch8WsMsg
		json.Unmarshal(data, &v)
		out = append(out, v)
	}
	return out
}

func guiPart(m map[string]any) streamPart {
	// Both published envelopes are valid: neutral embedded parts and the older
	// flattened fields. Compare display content without imposing Go type names.
	p := object(m["part"])
	out := streamPart{kind: fmt.Sprint(m["kind"])}
	if m["opaque"] == true {
		out.kind = "opaque"
	}
	if p != nil {
		out.kind = fmt.Sprint(p["type"])
	}
	if text, ok := m["text"].(string); ok {
		out.kind = "text"
		out.text = text
	}
	if p["type"] == "text" {
		out.text, _ = p["text"].(string)
	}
	if p["type"] == "tool_call" {
		out.name, _ = p["name"].(string)
		out.args = p["args"]
	}
	if name, ok := m["tool"].(string); ok {
		out.kind = "tool_call"
		out.name = name
		out.args = m["args"]
		if s, ok := out.args.(string); ok {
			json.Unmarshal([]byte(s), &out.args)
		}
	}
	return out
}

// JSON identity can be a number or namespaced string. Missing IDs are failures,
// not the string "<nil>" accidentally becoming a valid shared artifact key.
func artifactID(m map[string]any) string {
	switch v := m["part_id"].(type) {
	case string:
		return v
	case float64:
		if v >= 0 && v == float64(uint64(v)) {
			return fmt.Sprint(v)
		}
	}
	return ""
}

func checkGUIWire(c *grade.Check, messages []map[string]any) {
	// Type labels alone can hide empty tool errors or unrelated final content.
	// Correlate the actual call and each text card against the provider fixture.
	texts := map[string]string{}
	finals := map[string]streamPart{}
	dispatched, finished, ended, toolFinal := false, false, false, false
	for _, m := range messages {
		switch m["type"] {
		case "part_delta":
			id := artifactID(m)
			chunk, ok := m["chunk"].(string)
			if id == "" || !ok || (m["kind"] != "text" && m["kind"] != "thinking" && m["kind"] != "tool_call") {
				fail(c, "delta lacks identity/kind/chunk")
				continue
			}
			if m["kind"] == "text" {
				texts[id] += chunk
			}
		case "part_final":
			id := artifactID(m)
			if id == "" {
				fail(c, "final lacks artifact identity")
				continue
			}
			if _, ok := finals[id]; ok {
				fail(c, "final artifact identity reused: %s", id)
			}
			finals[id] = guiPart(m)
			p := finals[id]
			if p.kind != "text" && p.kind != "opaque" && p.kind != "thinking" && p.kind != "tool_call" {
				fail(c, "final lacks recognizable content kind")
			}
			if p.kind == "tool_call" && p.name == "run_command" && reflect.DeepEqual(p.args, map[string]any{"command": "printf GUI_TOOL_RESULT"}) {
				toolFinal = true
			}
		case "tool_dispatched":
			if m["call_id"] == "gui_call" {
				dispatched = m["name"] == "run_command" && reflect.DeepEqual(m["input"], map[string]any{"command": "printf GUI_TOOL_RESULT"})
			}
		case "tool_finished":
			if m["call_id"] == "gui_call" {
				bad, ok := m["is_error"].(bool)
				result, text := m["result"].(string)
				finished = ok && !bad && text && strings.Contains(result, "GUI_TOOL_RESULT")
			}
		case "turn_ended":
			ended = m["text"] == guiAnswer && (m["error"] == nil || m["error"] == "")
		}
	}
	for id, text := range texts {
		if f, ok := finals[id]; !ok || f.kind != "text" || f.text != text {
			fail(c, "text delta group %s differs from its final", id)
		}
	}
	for _, want := range []string{guiIntro, guiAnswer} {
		found := false
		for id, p := range finals {
			if p.kind == "text" && p.text == want && texts[id] == want {
				found = true
			}
		}
		if !found {
			fail(c, "missing actual text content %q", want)
		}
	}
	if !toolFinal {
		fail(c, "finalized tool name/arguments missing or changed")
	}
	if !dispatched || !finished || !ended {
		fail(c, "actual call input/result or final reply missing: dispatch=%t result=%t ended=%t", dispatched, finished, ended)
	}
}

func checkGUIReplay(c *grade.Check, live, replay []map[string]any) {
	// Compare actual completed facts, not just type labels. Replay may add metadata
	// but it must address the same cards and retain their observable content.
	for _, m := range live {
		kind := m["type"]
		if kind != "part_final" && kind != "tool_dispatched" && kind != "tool_finished" {
			continue
		}
		found := false
		for _, n := range replay {
			if n["type"] != kind {
				continue
			}
			if kind == "part_final" {
				if artifactID(m) != "" && artifactID(n) == artifactID(m) && reflect.DeepEqual(guiPart(m), guiPart(n)) {
					found = true
				}
			} else if n["call_id"] == m["call_id"] {
				keys := []string{"name", "input"}
				if kind == "tool_finished" {
					keys = []string{"result", "is_error"}
				}
				same := true
				for _, key := range keys {
					same = same && reflect.DeepEqual(m[key], n[key])
				}
				found = found || same
			}
		}
		if !found {
			fail(c, "replay lost %s content/identity (part=%v call=%v)", kind, m["part_id"], m["call_id"])
		}
	}
	for _, m := range replay {
		if m["type"] == "part_delta" || m["type"] == "part_partial" {
			fail(c, "completed turn replay retains a stale partial")
		}
		if m["type"] == "state_changed" && (m["from"] == nil || m["to"] == nil || m["from"] == "" || m["to"] == "") {
			fail(c, "replay state lacks from/to")
		}
	}
}

func checkGUITrace(c *grade.Check, trace string) {
	inbound, outbound := false, false
	for i, line := range strings.Split(strings.TrimSpace(trace), "\n") {
		// WriteJSON may retain a trailing newline inside the trace. Empty
		// separators carry no record; each actual message is still parsed fully.
		if strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.SplitN(line, " ", 3)
		if len(fields) != 3 {
			fail(c, "gui.log line %d lacks timestamp/direction/JSON", i+1)
			continue
		}
		_, err := time.Parse(time.RFC3339Nano, fields[0])
		var message map[string]any
		if err != nil || json.Unmarshal([]byte(fields[2]), &message) != nil || message["type"] == nil {
			fail(c, "gui.log line %d has invalid timestamp/message JSON", i+1)
		}
		switch fields[1] {
		case ">":
			inbound = true
		case "<":
			outbound = true
		default:
			fail(c, "gui.log line %d has invalid direction", i+1)
		}
	}
	if !inbound || !outbound {
		fail(c, "gui.log needs both actual traffic directions")
	}
}
