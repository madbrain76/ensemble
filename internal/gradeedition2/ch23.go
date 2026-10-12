package gradeedition2

import (
	"fmt"
	"strings"

	"github.com/waywardgeek/ensemble/internal/fakevendor"
	"github.com/waywardgeek/ensemble/internal/grade"
)

// Chapter23 preserves the seven original fleet-management weights. Actual owner
// construction and the stronger interruption/slow-client lifetime obligations
// also require the separately recorded source review; a reply cannot prove them.
func Chapter23(bin, gui string) (grade.Report, error) {
	checks := grade.Ch26Checks(grade.Ch26Result{})
	for i := range checks {
		checks[i].Passed = true
		checks[i].Earned = checks[i].Points
		checks[i].Details = nil
	}
	checks[0].Title = "Correlated child status and ordered observed activity"
	checks[1].Title = "Named creation selects model/skill, restores conversation and deletes durably"
	checks[2].Title = "Owned success and indistinguishable foreign/missing management refusal"
	checks[3].Title = "Selected child progress wakes its change wait; command and submission complete"
	checks[4].Title = "Actual CLI recursively constructs and completes a three-generation fleet"
	checks[5].Title = "Live nested GUI tree, conversation selection and truthful replay state"
	fleetMessaging(bin, &checks[4])
	fleetIdle(bin, &checks[6], &checks[6])
	guiBin, cleanup, err := grade.Build(gui)
	if err != nil {
		reuseFail(checks, err, 0, 1, 2, 3, 5)
	} else {
		defer cleanup()
		managementGUI(guiBin, checks)
	}
	return grade.NewTitledReport("Edition 2, Chapter 23 — Managing the Swarm", checks, ""), nil
}

// Tool results are matched to the actual call, including error status. Foreign
// and missing IDs differ by necessity; normalize only those requested names.
func managementResults(reqs []fakevendor.Recorded, checks []grade.Check) {
	status := fleetResult(&checks[0], reqs, "child-status", false)
	if !strings.Contains(strings.ToLower(status), "activity") || !strings.Contains(strings.ToLower(status), "idle") {
		fail(&checks[0], "completed child status lacks idle activity: %s", status)
	}
	for _, name := range []string{"get_sub_agent_status", "shutdown_agent"} {
		foreign := fleetResult(&checks[2], reqs, name+"-foreign", true)
		missing := fleetResult(&checks[2], reqs, name+"-missing", true)
		foreign = strings.ReplaceAll(foreign, "worker23", "TARGET")
		missing = strings.ReplaceAll(missing, "absent23", "TARGET")
		if foreign != missing {
			fail(&checks[2], "%s distinguishes foreign and absent targets: %q / %q", name, foreign, missing)
		}
	}
	fleetResult(&checks[2], reqs, "child-status", false)
	changed := fleetResult(&checks[3], reqs, "wait-change", false)
	if !strings.Contains(changed, "worker23") || (!strings.Contains(changed, "state") && !strings.Contains(changed, "status")) {
		fail(&checks[3], "wait lost selected child transition: %s", changed)
	}
	output := fleetResult(&checks[3], reqs, "child-command", false)
	if !strings.Contains(output, "COMMAND_RESULT_23") {
		fail(&checks[3], "held child command did not complete: %s", output)
	}
	joined := fleetResult(&checks[3], reqs, "join-worker", false)
	if !strings.Contains(joined, "WORKER_RESULT_23") {
		fail(&checks[3], "selected child submission missing: %s", joined)
	}
}

func managementReply(body []byte, work string) *fakevendor.Reply {
	b := string(body)
	if strings.Contains(b, "PENDING_TASK_23") {
		if fleetHas(body, "pending-submit") {
			return &fakevendor.Reply{Text: "PENDING_DONE_23"}
		}
		return fleetCall("pending-submit", "submit", map[string]any{"result": "PENDING_DONE_23"})
	}
	// These are external scripted decisions, not replacement fleet machinery.
	for _, final := range [][2]string{{"peer-submit", "NAMED_CONVERSATION_23"}, {"manage-submit", "MANAGEMENT_DONE_23"}, {"child-submit", "WORKER_RESULT_23"}, {"intruder-submit", "INTRUDER_DONE_23"}} {
		if fleetHas(body, final[0]) {
			return &fakevendor.Reply{Text: final[1]}
		}
	}
	if !strings.Contains(b, "MANAGEMENT_ROOT_23") && strings.Contains(b, "INTRUDER_TASK_23") {
		for _, name := range []string{"get_sub_agent_status", "shutdown_agent"} {
			for _, pair := range [][2]string{{"foreign", "worker23"}, {"missing", "absent23"}} {
				id := name + "-" + pair[0]
				if !fleetHas(body, id) {
					return fleetCall(id, name, map[string]any{"agent_id": pair[1]})
				}
			}
		}
		return fleetCall("intruder-submit", "submit", map[string]any{"result": "INTRUDER_DONE_23"})
	}
	if !strings.Contains(b, "MANAGEMENT_ROOT_23") && strings.Contains(b, "WORKER_TASK_23") {
		if !fleetHas(body, "child-command") {
			return fleetCall("child-command", "run_command", map[string]any{"command": fmt.Sprintf("printf ready > %q; while [ ! -f %q ]; do sleep .02; done; printf COMMAND_RESULT_23", work+"/ready", work+"/release"), "ai_callback_delay": 10})
		}
		return fleetCall("child-submit", "submit", map[string]any{"result": "WORKER_RESULT_23"})
	}
	if strings.Contains(b, "PEER_TASK_23") {
		return fleetCall("peer-submit", "submit", map[string]any{"result": "NAMED_CONVERSATION_23"})
	}
	for _, p := range []sandboxProbe{
		{"spawn-worker", "spawn_sub_agent", map[string]any{"agent_id": "worker23", "task": "WORKER_TASK_23"}},
		{"wait-change", "wait_for_agent_change", map[string]any{"agent_ids": []string{"worker23"}, "timeout_seconds": 5, "ai_callback_delay": 10}},
		{"join-worker", "join_agents", map[string]any{"agent_ids": []string{"worker23"}, "ai_callback_delay": 5}},
		{"child-status", "get_sub_agent_status", map[string]any{"agent_id": "worker23"}},
		{"spawn-intruder", "spawn_sub_agent", map[string]any{"agent_id": "intruder23", "task": "INTRUDER_TASK_23"}},
		{"join-intruder", "join_agents", map[string]any{"agent_ids": []string{"intruder23"}, "ai_callback_delay": 5}},
	} {
		if !fleetHas(body, p.id) {
			r := p.reply()
			return &r
		}
	}
	return fleetCall("manage-submit", "submit", map[string]any{"result": "MANAGEMENT_DONE_23"})
}
