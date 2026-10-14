package gradeedition2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/waywardgeek/ensemble/internal/fakevendor"
	"github.com/waywardgeek/ensemble/internal/grade"
)

// The schema describes the whole direct argument object. The nullable field is
// deliberately present in the accepted result so nested JSON is compared as data.
const fleetReviewSchema = `{"type":"object","properties":{"summary":{"type":"string"},"findings":{"type":"array","items":{"type":"object","properties":{"path":{"type":"string"},"fixed":{"type":"boolean"}},"required":["path","fixed"],"additionalProperties":false}},"follow_up":{"type":["string","null"]}},"required":["summary","findings","follow_up"],"additionalProperties":false}`
const fleetReviewPayload = `{"summary":"STRUCTURED_RESULT_22","findings":[{"path":"main.go","fixed":false}],"follow_up":null}`
const fleetDefaultSchema = `{"type":"object","properties":{"result":{"type":"string"}},"required":["result"]}`

func fleetJSON(raw string) any {
	var value any
	if err := json.Unmarshal([]byte(raw), &value); err != nil {
		panic(err) // These are authored fixture constants, never model input.
	}
	return value
}

// fleetStructured adds wire and durable observations to the existing spawn,
// isolation and fan-in rows. Ordinary text siblings remain a positive control.
func fleetStructured(bin string, checks []grade.Check) {
	fleetTemp([]*grade.Check{&checks[0], &checks[1], &checks[3]}, func(work string) {
		const childTask = "STRUCTURED_TASK_22"
		const leafTask = "SCHEMA_GRANDCHILD_22"
		invalid := []sandboxProbe{
			{"missing-field", "submit", map[string]any{"summary": "missing"}},
			{"wrong-type", "submit", map[string]any{"summary": 7, "findings": []any{}, "follow_up": nil}},
			{"nested-field", "submit", map[string]any{"summary": "nested", "findings": []any{map[string]any{"path": "main.go"}}, "follow_up": nil}},
			{"extra-field", "submit", map[string]any{"summary": "extra", "findings": []any{}, "follow_up": nil, "unexpected": true}},
		}
		reqs, err := fleetRun(bin, work, func(body []byte) *fakevendor.Reply {
			// Inspect conversation text rather than quoted spawn arguments. Each real
			// child has only its own task; the parent also contains both call inputs.
			text := memoryConversation(body)
			if !strings.Contains(text, fleetRoot) {
				if strings.Contains(text, childTask) {
					if !fleetHas(body, "spawn-schema-leaf") {
						return fleetCall("spawn-schema-leaf", "spawn_sub_agent", map[string]any{"agent_id": "leaf", "task": leafTask})
					}
					if !fleetHas(body, "join-schema-leaf") {
						return fleetCall("join-schema-leaf", "join_agents", map[string]any{"agent_ids": []string{"leaf"}, "ai_callback_delay": 5})
					}
					for _, probe := range invalid {
						if !fleetHas(body, probe.id) {
							r := probe.reply()
							return &r
						}
					}
					return fleetCall("review-submit", "submit", fleetJSON(fleetReviewPayload))
				}
				if strings.Contains(text, leafTask) {
					return fleetCall("leaf-submit", "submit", map[string]any{"result": "DEFAULT_GRANDCHILD_22"})
				}
				return fleetCall("child-submit", "submit", map[string]any{"result": "DEFAULT_SIBLING_22"})
			}
			for _, probe := range []sandboxProbe{
				{"bad-required-schema", "spawn_sub_agent", map[string]any{"agent_id": "bad-required", "task": "MUST_NOT_LAUNCH", "submit_schema": map[string]any{"type": "object", "required": 7}}},
				{"bad-reference-schema", "spawn_sub_agent", map[string]any{"agent_id": "bad-reference", "task": "MUST_NOT_LAUNCH", "submit_schema": map[string]any{"type": "object", "$ref": "#/missing"}}},
				{"bad-keyword-schema", "spawn_sub_agent", map[string]any{"agent_id": "bad-keyword", "task": "MUST_NOT_LAUNCH", "submit_schema": map[string]any{"type": "object", "format": "email"}}},
				{"spawn-review", "spawn_sub_agent", map[string]any{"agent_id": "reviewer", "task": childTask, "submit_schema": fleetJSON(fleetReviewSchema)}},
				{"spawn-default", "spawn_sub_agent", map[string]any{"agent_id": "ordinary", "task": "DEFAULT_SIBLING_TASK_22"}},
				{"join-review", "join_agents", map[string]any{"agent_ids": []string{"reviewer", "ordinary"}, "ai_callback_delay": 5}},
				{"structured-status", "agents_status", map[string]any{}},
			} {
				if !fleetHas(body, probe.id) {
					r := probe.reply()
					return &r
				}
			}
			return fleetFinish()
		}, nil)
		if err != nil {
			for _, i := range []int{0, 1, 3} {
				fail(&checks[i], "structured fleet: %v", err)
			}
			return
		}
		for _, id := range []string{"bad-required-schema", "bad-reference-schema", "bad-keyword-schema"} {
			fleetResult(&checks[0], reqs, id, true)
		}
		for _, probe := range invalid {
			fleetResult(&checks[0], reqs, probe.id, true)
		}
		customSeen, defaultSeen := false, false
		// Declaration equality and runtime rejection are separate observations. A
		// correct advertised schema alone cannot establish that dispatch enforces it.
		for _, request := range reqs {
			text := memoryConversation(request.Body)
			if strings.Contains(text, fleetRoot) {
				continue
			}
			if strings.Contains(text, "MUST_NOT_LAUNCH") {
				fail(&checks[0], "invalid schema reached a child model request")
			}
			want := ""
			if strings.Contains(text, childTask) {
				want, customSeen = fleetReviewSchema, true
			} else if strings.Contains(text, leafTask) {
				want, defaultSeen = fleetDefaultSchema, true
			}
			if want == "" {
				continue
			}
			var actual any
			for _, value := range list(request.JSON()["tools"]) {
				tool := object(value)
				if tool["name"] == "submit" {
					actual = tool["input_schema"]
				}
			}
			if !reflect.DeepEqual(actual, fleetJSON(want)) {
				fail(&checks[0], "child/grandchild declaration did not preserve its own schema")
			}
		}
		if !customSeen || !defaultSeen {
			fail(&checks[1], "custom child or default grandchild never reached the provider")
		}
		leaf := fleetResult(&checks[1], reqs, "join-schema-leaf", false)
		if !strings.Contains(leaf, "DEFAULT_GRANDCHILD_22") {
			fail(&checks[1], "default grandchild result did not reach custom parent")
		}
		for _, id := range []string{"join-review", "structured-status"} {
			// Compare the nested value directly: finding its serialized marker inside
			// a quoted JSON string would falsely accept double encoding.
			text := fleetResult(&checks[3], reqs, id, false)
			var statuses []map[string]any
			if err := json.Unmarshal([]byte(text), &statuses); err != nil {
				fail(&checks[3], "%s result is not a JSON status array: %v", id, err)
				continue
			}
			found := map[string]bool{}
			for _, status := range statuses {
				switch status["agent_id"] {
				case "reviewer":
					found["reviewer"] = status["state"] == "submitted" && reflect.DeepEqual(status["payload"], fleetJSON(fleetReviewPayload))
				case "ordinary":
					found["ordinary"] = status["state"] == "submitted" && status["result"] == "DEFAULT_SIBLING_22"
				default:
					fail(&checks[0], "refused launch appeared in owned status: %v", status["agent_id"])
				}
			}
			if !found["reviewer"] || !found["ordinary"] {
				fail(&checks[3], "%s lost exact nested JSON or legacy text: %s", id, text)
			}
		}
		// Inspect actual durable rows after normal process drain. An echoed submit
		// argument or saved status label alone cannot establish durable typed data.
		stored := false
		err = filepath.WalkDir(work, func(path string, entry os.DirEntry, err error) error {
			if err != nil || entry.IsDir() || entry.Name() != "submit.md" {
				return err
			}
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			for _, line := range bytes.Split(bytes.TrimSpace(data), []byte{'\n'}) {
				var row map[string]any
				if err := json.Unmarshal(line, &row); err != nil {
					return fmt.Errorf("submission row: %w", err)
				}
				if reflect.DeepEqual(row["payload"], fleetJSON(fleetReviewPayload)) && reflect.DeepEqual(row["submit_schema"], fleetJSON(fleetReviewSchema)) && row["version"] == float64(1) {
					// Bind the row to its actual owner's saved conversation. Merely
					// finding these bytes anywhere under the root is not isolation.
					data, err := os.ReadFile(filepath.Join(filepath.Dir(path), "save.json"))
					if err != nil {
						return err
					}
					var save any
					if err := json.Unmarshal(data, &save); err != nil {
						return err
					}
					if !fleetExactText(save, childTask) {
						return fmt.Errorf("custom row stored outside its child's saved conversation")
					}
					stored = true
				}
			}
			return nil
		})
		if err != nil || !stored {
			fail(&checks[1], "durable custom schema/payload row absent: %v", err)
		}
	})
}
