package gradeedition2

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/waywardgeek/ensemble/internal/fakevendor"
	"github.com/waywardgeek/ensemble/internal/grade"
)

// Chapter12 grades startup GUI skills and their real MCP lifecycle on both clients.
// Native browser effects and speech/input provenance retain separate browser proof.
func Chapter12(bin, core, gui, guiModel, consumer, workflow, exercise, streamModel string, example ConsumerCase, media map[string]string) (grade.Report, error) {
	checks := guiDebugChecks(bin, gui, guiModel)
	previous, err := Chapter11(bin, core, gui, guiModel, consumer, workflow, exercise, streamModel, example, media)
	if err != nil {
		return grade.Report{}, err
	}
	parity := grade.Check{ID: "ch11-parity", Title: "All corrected Chapter 11 checks pass", Points: 10, Earned: 10, Passed: true}
	for _, c := range previous.Checks {
		if !c.Passed {
			fail(&parity, "%s: %s", c.ID, strings.Join(c.Details, "; "))
		}
	}
	return grade.NewTitledReport("Edition 2, Chapter 12 — The Agent Sees Itself", append(checks, parity), ""), nil
}

func guiDebugChecks(bin, gui, model string) []grade.Check {
	checks := grade.Ch13Checks(grade.Ch13Result{SkillLoadsOK: true, EphemeralInjectedOK: true, SnapshotUpdatesOK: true, TTSVisibilityOK: true, GUIInteractionOK: true, SkillUnloadOK: true})[:6]
	// The two commands must actually use their startup flags. A passing core
	// constructor does not prove the independently built GUI passes those inputs.
	guiDebugClient(bin, model, nil, "core", checks)
	guiBin, cleanup, err := grade.Build(gui)
	if err != nil {
		for i := range checks {
			fail(&checks[i], "GUI build: %v", err)
		}
		return checks
	}
	defer cleanup()
	static, err := filepath.Abs(filepath.Join(gui, "web", "gui"))
	if err != nil {
		for i := range checks {
			fail(&checks[i], "GUI assets: %v", err)
		}
		return checks
	}
	guiDebugClient(guiBin, model, []string{"--port", "0", "--gui-dir", static}, "GUI", checks)
	return checks
}

// This stateful external peer is the original exercise's simulated browser.
// Only the selected real method/arguments change its state. Unique output-only
// counters distinguish new observations from call arguments or retained history.
const guiDebugPeer = `import json,os,sys
open('browser.pid','w').write(str(os.getpid()))
counts={};button='enabled';value='<empty>'
for line in sys.stdin:
 m=json.loads(line);method=m.get('method')
 with open('browser-rpc.jsonl','a') as trace: trace.write(json.dumps(m)+'\n')
 if 'id' not in m: continue
 if method=='initialize': result={'protocolVersion':'2024-11-05','capabilities':{'tools':{}}}
 elif method=='tools/list':
  result={'tools':[{'name':name,'description':name,'inputSchema':{'type':'object','properties':{'selector':{'type':'string'},'text':{'type':'string'}}},'ephemeral':mode} for name,mode in [('gui_snapshot','round'),('tts_queue','round'),('gui_click',''),('gui_input','')]]}
 elif method=='tools/call':
  p=m['params'];name=p['name'];args=p.get('arguments',{});counts[name]=counts.get(name,0)+1;bad=False
  if name=='gui_snapshot': text='SNAPSHOT_ONLY_'+str(counts[name])+' button='+button+' input='+value
  elif name=='tts_queue': text=json.dumps([{'text':'TTS_ONLY_'+str(counts[name]),'status':'pending'}])
  elif name=='gui_click' and args.get('selector')=='#button-A': button='disabled';text='CLICK_RESULT_ONLY'
  elif name=='gui_input' and args.get('selector')=='#search-box': value=args.get('text','');text='INPUT_RESULT_ONLY:'+value
  else: text='BAD_BROWSER_OPERATION';bad=True
  result={'content':[{'type':'text','text':text}],'isError':bad}
 else: continue
 print(json.dumps({'jsonrpc':'2.0','id':m['id'],'result':result}),flush=True)
`

func guiDebugClient(bin, model string, clientArgs []string, label string, checks []grade.Check) {
	badSetup := func(err error) {
		for i := range checks {
			fail(&checks[i], "%s session: %v", label, err)
		}
	}
	work, err := os.MkdirTemp("", "edition2-gui-debug-")
	if err != nil {
		badSetup(err)
		return
	}
	defer os.RemoveAll(work)
	python, err := exec.LookPath("python3")
	if err != nil {
		badSetup(err)
		return
	}
	script := filepath.Join(work, "browser.py")
	if err = os.WriteFile(script, []byte(guiDebugPeer), 0600); err != nil {
		badSetup(err)
		return
	}
	// Valid public metadata corrects the old missing-kind fixture. The primary
	// admits the loadable skill, which explicitly grants all four discovered tools.
	args, _ := json.Marshal([]string{"-u", script})
	skills := filepath.Join(work, "skills")
	for name, body := range map[string]string{
		"ensemble":  "---\nname: ensemble\ndescription: GUI exercise primary\ntype: primary\nloadable-skills: gui-debug\n---\nPRIMARY_MANUAL_ONLY\n",
		"gui-debug": "---\nname: gui-debug\ndescription: GUI observation exercise\ntype: loadable\ndepends: ensemble\ntools: gui_snapshot tts_queue gui_click gui_input\nmcp_servers:\n  - name: browser-debug\n    transport: stdio\n    command: " + python + "\n    args: " + string(args) + "\n---\nSTARTUP_MANUAL_ONLY\n",
	} {
		dir := filepath.Join(skills, name)
		if err = os.MkdirAll(dir, 0700); err == nil {
			err = os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(body), 0600)
		}
		if err != nil {
			badSetup(err)
			return
		}
	}
	provider := fakevendor.New([]fakevendor.Reply{
		skillCall("debug-click", "gui_click", `{"selector":"#button-A"}`),
		skillCall("debug-input", "gui_input", `{"selector":"#search-box","text":"INPUT_ARGUMENT_ONLY"}`),
		skillCall("debug-unload", "unload_skill", `{"name":"gui-debug"}`),
		skillCall("debug-guessed", "gui_click", `{"selector":"#button-A"}`),
		{Text: "GUI_DEBUG_DONE"},
	})
	defer provider.Close()
	launch := append(append([]string{}, clientArgs...), "--gui-debug", "--skills-dir", skills)
	p, err := startMCPCommand(bin, work, launch, []string{"LLM_MODEL=" + model, "LLM_BASE_URL=" + provider.URL(), "EN_SKILLS_DIR=" + filepath.Join(work, "missing-env-skills"), "EN_PRIMARY_SKILL=ensemble"})
	if err != nil {
		badSetup(err)
		return
	}
	defer p.close()
	// Cleanup must also remove a deliberately surviving mutant's external peer.
	// Read its self-reported PID, never a name-based process search or another job.
	defer func() {
		if pid := guiDebugPID(work); pid > 0 {
			process, _ := os.FindProcess(pid)
			_ = process.Kill()
		}
	}()
	if err = p.send(map[string]string{"kind": "prompt", "text": "Exercise the browser tools."}); err != nil {
		badSetup(err)
		return
	}
	for {
		line, e := p.read()
		if e != nil {
			badSetup(e)
			return
		}
		if strings.Contains(line, "GUI_DEBUG_DONE") {
			break
		}
	}
	requests := provider.Requests()
	if len(requests) != 5 {
		badSetup(fmt.Errorf("want five actual requests, got %d", len(requests)))
		return
	}
	// Keep the parent command waiting on stdin while assessing unload. Otherwise
	// ordinary process shutdown could hide a missing skill-specific disconnect.
	if err = p.process.Signal(syscall.Signal(0)); err != nil {
		fail(&checks[5], "%s Agent command ended before unload observation: %v", label, err)
	}
	pid := guiDebugPID(work)
	if pid <= 0 {
		fail(&checks[0], "%s no actual peer PID", label)
		fail(&checks[5], "%s cannot observe peer lifetime", label)
	} else if err = syscall.Kill(pid, 0); err != syscall.ESRCH {
		fail(&checks[5], "%s peer still present after unload: %v", label, err)
	}
	system := requests[0].JSON()["system"]
	if !strings.Contains(fmt.Sprint(system), "PRIMARY_MANUAL_ONLY") || !strings.Contains(fmt.Sprint(system), "STARTUP_MANUAL_ONLY") {
		fail(&checks[0], "%s startup instructions missing", label)
	}
	for i, r := range requests {
		if r.JSON()["system"] != system {
			fail(&checks[0], "%s request%d changed frozen startup instructions", label, i+1)
		}
		declared := []string{"load_skill", "unload_skill"}
		target := &checks[5]
		if i < 3 {
			declared = append(declared, "gui_click", "gui_input")
			target = &checks[0]
		}
		skillDeclarations(target, r, declared)
		content := guiDebugText(r)
		// Check each observation once, then its absence on every later request.
		// The event log may retain it; retained dialogue and current context may not.
		for n := 1; n <= 3; n++ {
			want := 0
			if i < 3 && n == i+1 {
				want = 1
			}
			if strings.Count(content, fmt.Sprintf("SNAPSHOT_ONLY_%d", n)) != want {
				target := &checks[1]
				if i >= 3 {
					target = &checks[5]
				}
				fail(target, "%s request%d stale/missing snapshot%d", label, i+1, n)
			}
			if strings.Count(content, fmt.Sprintf("TTS_ONLY_%d", n)) != want {
				target := &checks[3]
				if i >= 3 {
					target = &checks[5]
				}
				fail(target, "%s request%d stale/missing TTS%d", label, i+1, n)
			}
		}
		if i < 3 && !strings.Contains(content, `"status": "pending"`) {
			fail(&checks[3], "%s request%d lost utterance status", label, i+1)
		}
		if i < 3 {
			button, value := "disabled", "<empty>"
			if i == 0 {
				button = "enabled"
			}
			if i == 2 {
				value = "INPUT_ARGUMENT_ONLY"
			}
			want := fmt.Sprintf("SNAPSHOT_ONLY_%d button=%s input=%s", i+1, button, value)
			if !strings.Contains(content, want) {
				fail(&checks[2], "%s request%d lacks changed peer state %q", label, i+1, want)
			}
		}
	}
	// Return-only sentinels cannot match the earlier model's selector or text.
	if !strings.Contains(skillResult(&checks[4], requests[1], "debug-click", false), "CLICK_RESULT_ONLY") {
		fail(&checks[4], "%s actual click result missing", label)
	}
	if !strings.Contains(skillResult(&checks[4], requests[2], "debug-input", false), "INPUT_RESULT_ONLY:INPUT_ARGUMENT_ONLY") {
		fail(&checks[4], "%s actual input result/arguments missing", label)
	}
	skillResult(&checks[5], requests[3], "debug-unload", false)
	skillResult(&checks[5], requests[4], "debug-guessed", true)
	// The peer log independently confirms startup and that no post-unload guessed
	// call ran. The mere presence/absence of a model declaration is insufficient.
	data, err := os.ReadFile(filepath.Join(work, "browser-rpc.jsonl"))
	if err != nil {
		badSetup(err)
		return
	}
	methods, calls := map[string]int{}, map[string]int{}
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		var m map[string]any
		if json.Unmarshal([]byte(line), &m) != nil {
			fail(&checks[0], "%s malformed peer trace", label)
			continue
		}
		method, _ := m["method"].(string)
		methods[method]++
		if method == "tools/call" {
			name, _ := object(m["params"])["name"].(string)
			calls[name]++
		}
	}
	for _, method := range []string{"initialize", "notifications/initialized", "tools/list"} {
		if methods[method] != 1 {
			fail(&checks[0], "%s %s count=%d", label, method, methods[method])
		}
	}
	for name, want := range map[string]int{"gui_snapshot": 3, "tts_queue": 3, "gui_click": 1, "gui_input": 1} {
		if calls[name] != want {
			fail(&checks[5], "%s external %s calls=%d want%d", label, name, calls[name], want)
		}
	}
}

func guiDebugPID(work string) int {
	data, err := os.ReadFile(filepath.Join(work, "browser.pid"))
	if err != nil {
		return 0
	}
	pid, _ := strconv.Atoi(strings.TrimSpace(string(data)))
	return pid
}

// Only text blocks are automatic context in this fixture. Tool result text is
// checked separately by call identity, preventing earlier arguments from passing.
func guiDebugText(r fakevendor.Recorded) string {
	var text []string
	for _, message := range list(r.JSON()["messages"]) {
		for _, part := range list(object(message)["content"]) {
			p := object(part)
			if p["type"] == "text" {
				text = append(text, fmt.Sprint(p["text"]))
			}
		}
	}
	return strings.Join(text, "\n")
}
