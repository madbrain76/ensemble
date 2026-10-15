package gradeedition2

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/waywardgeek/ensemble/internal/fakevendor"
	"github.com/waywardgeek/ensemble/internal/grade"
)

// Chapter21 preserves the seven original fleet-management weights. Actual owner
// construction and the stronger interruption/slow-client lifetime obligations
// also require the separately recorded source review; a reply cannot prove them.
func Chapter21(bin, gui string) (grade.Report, error) {
	checks := grade.Ch26Checks(grade.Ch26Result{})
	for i := range checks {
		checks[i].Passed = true
		checks[i].Earned = checks[i].Points
		checks[i].Details = nil
	}
	checks[0].Title = "Correlated activity and validated structured submission state"
	checks[1].Title = "Named creation selects model/skill, restores conversation and deletes durably"
	checks[2].Title = "Owned typed retrieval and indistinguishable foreign/missing refusal"
	checks[3].Title = "Selected child progress wakes its change wait; command and submission complete"
	checks[4].Title = "Actual CLI recursively constructs and completes a three-generation fleet"
	checks[5].Title = "Live nested GUI tree, exact structured replay and truthful state"
	fleetMessaging(bin, &checks[4])
	fleetIdle(bin, &checks[6], &checks[6])
	guiBin, cleanup, err := grade.Build(gui)
	if err != nil {
		reuseFail(checks, err, 0, 1, 2, 3, 5, 6)
	} else {
		defer cleanup()
		fleetRootChat(guiBin, &checks[6])
		managementGUI(guiBin, checks)
	}
	return grade.NewTitledReport("Edition 2, Chapter 21 — Managing the Swarm", checks, ""), nil
}

// Tool results are matched to the actual call, including error status. Foreign
// and missing IDs differ by necessity; normalize only those requested names.
func managementResults(reqs []fakevendor.Recorded, checks []grade.Check) {
	status := fleetResult(&checks[0], reqs, "child-status", false)
	if !strings.Contains(strings.ToLower(status), "activity") || !strings.Contains(strings.ToLower(status), "idle") {
		fail(&checks[0], "completed child status lacks idle activity: %s", status)
	}
	for _, name := range []string{"get_sub_agent_status", "shutdown_agent", "get_submission"} {
		foreign := fleetResult(&checks[2], reqs, name+"-foreign", true)
		missing := fleetResult(&checks[2], reqs, name+"-missing", true)
		foreign = strings.ReplaceAll(foreign, "worker23", "TARGET")
		missing = strings.ReplaceAll(missing, "absent23", "TARGET")
		if foreign != missing {
			fail(&checks[2], "%s distinguishes foreign and absent targets: %q / %q", name, foreign, missing)
		}
	}
	fleetResult(&checks[2], reqs, "child-status", false)
	fleetResult(&checks[0], reqs, "child-invalid", true)
	var rejected, completed map[string]any
	if err := json.Unmarshal([]byte(fleetResult(&checks[0], reqs, "invalid-status", false)), &rejected); err != nil || rejected["agent_id"] != "worker23" || rejected["state"] != "running" || rejected["payload"] != nil {
		fail(&checks[0], "rejected submission completed the child: %v (%v)", rejected, err)
	}
	if err := json.Unmarshal([]byte(status), &completed); err != nil || completed["state"] != "submitted" || !reflect.DeepEqual(completed["payload"], fleetJSON(fleetReviewPayload)) {
		fail(&checks[0], "accepted status lost structured result: %s (%v)", status, err)
	}
	managementPayload(&checks[2], fleetResult(&checks[2], reqs, "child-retrieve", false), "worker23")
	changed := fleetResult(&checks[3], reqs, "wait-change", false)
	if !strings.Contains(changed, "worker23") || (!strings.Contains(changed, "state") && !strings.Contains(changed, "status")) {
		fail(&checks[3], "wait lost selected child transition: %s", changed)
	}
	output := fleetResult(&checks[3], reqs, "child-command", false)
	if !strings.Contains(output, "COMMAND_RESULT_23") {
		fail(&checks[3], "held child command did not complete: %s", output)
	}
	joined := fleetResult(&checks[3], reqs, "join-worker", false)
	var joinedRows []map[string]any
	if err := json.Unmarshal([]byte(joined), &joinedRows); err != nil || len(joinedRows) != 1 || !reflect.DeepEqual(joinedRows[0]["payload"], fleetJSON(fleetReviewPayload)) {
		fail(&checks[3], "selected child submission missing: %s", joined)
	}
}

func managementReply(body []byte, work string) *fakevendor.Reply {
	b := string(body)
	// A restored owner explicitly reconstructs its task child before using the
	// durable fallback. Registry restoration alone must not revive task children.
	if strings.Contains(memoryConversation(body), "RESTORE_ROOT_23") {
		for _, p := range []sandboxProbe{
			{"respawn-worker", "spawn_sub_agent", map[string]any{"agent_id": "worker23", "task": "REOPEN_TASK_23"}},
			{"durable-retrieve", "get_submission", map[string]any{"agent_id": "worker23"}},
			{"durable-status", "get_sub_agent_status", map[string]any{"agent_id": "worker23"}},
			{"stop-reopened", "kill_agent", map[string]any{"agent_id": "worker23"}},
		} {
			if !fleetHas(body, p.id) {
				r := p.reply()
				return &r
			}
		}
		return &fakevendor.Reply{Text: "RESTORED_RESULT_23"}
	}
	if strings.Contains(b, "PENDING_TASK_23") {
		return &fakevendor.Reply{Text: "PENDING_DONE_23"}
	}
	// These are external scripted decisions, not replacement fleet machinery.
	for _, final := range [][2]string{{"peer-submit", "NAMED_CONVERSATION_23"}, {"manage-submit", "MANAGEMENT_DONE_23"}, {"child-submit", "WORKER_RESULT_23"}, {"intruder-submit", "INTRUDER_DONE_23"}} {
		if fleetHas(body, final[0]) {
			return &fakevendor.Reply{Text: final[1]}
		}
	}
	if !strings.Contains(b, "MANAGEMENT_ROOT_23") && strings.Contains(b, "INTRUDER_TASK_23") {
		for _, name := range []string{"get_sub_agent_status", "shutdown_agent", "get_submission"} {
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
		if !fleetHas(body, "child-invalid") {
			return fleetCall("child-invalid", "submit", map[string]any{"summary": "missing required fields"})
		}
		return fleetCall("child-submit", "submit", fleetJSON(fleetReviewPayload))
	}
	if strings.Contains(b, "PEER_TASK_23") {
		return &fakevendor.Reply{Text: "NAMED_CONVERSATION_23"}
	}
	for _, p := range []sandboxProbe{
		{"spawn-worker", "spawn_sub_agent", map[string]any{"agent_id": "worker23", "task": "WORKER_TASK_23", "submit_schema": fleetJSON(fleetReviewSchema)}},
		{"wait-change", "wait_for_agent_change", map[string]any{"agent_ids": []string{"worker23"}, "timeout_seconds": 5, "ai_callback_delay": 10}},
		{"invalid-status", "get_sub_agent_status", map[string]any{"agent_id": "worker23"}},
		{"join-worker", "join_agents", map[string]any{"agent_ids": []string{"worker23"}, "ai_callback_delay": 5}},
		{"child-status", "get_sub_agent_status", map[string]any{"agent_id": "worker23"}},
		{"child-retrieve", "get_submission", map[string]any{"agent_id": "worker23"}},
		{"spawn-intruder", "spawn_sub_agent", map[string]any{"agent_id": "intruder23", "task": "INTRUDER_TASK_23"}},
		{"join-intruder", "join_agents", map[string]any{"agent_ids": []string{"intruder23"}, "ai_callback_delay": 5}},
	} {
		if !fleetHas(body, p.id) {
			r := p.reply()
			return &r
		}
	}
	return &fakevendor.Reply{Text: "MANAGEMENT_DONE_23"}
}
