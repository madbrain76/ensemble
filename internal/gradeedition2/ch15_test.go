package gradeedition2

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/waywardgeek/ensemble/internal/grade"
)

// TestChapter15Memory exercises the actual accepted GUI/Agent, owned settings,
// memory files, history and process lifetime. Only provider responses are fake.
func TestChapter15Memory(t *testing.T) {
	bin, cleanup, err := grade.Build("../../agent/gui")
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	r := memoryChecks(bin)
	if r.Fatal != "" {
		t.Fatal(r.Fatal)
	}
	for id, e := range r.Errs {
		if e != "" {
			t.Errorf("%s: %s", id, e)
		}
	}
}

// TestMemoryBackground exercises the new scored scenario independently so a
// targeted runtime defect must fail that predicate rather than an older row.
func TestMemoryBackground(t *testing.T) {
	bin, cleanup, err := grade.Build("../../agent/gui")
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	if err = memoryBackground(bin, t.TempDir()); err != nil {
		t.Fatal(err)
	}
}

// TestMemoryRequestClassification keeps ordinary delegated submit traffic out of
// the memory route, including a task whose prose itself mentions memory.
func TestMemoryRequestClassification(t *testing.T) {
	for _, tc := range []struct {
		name, body string
		memory     bool
	}{
		{"memory", `{"tools":[{"name":"submit","input_schema":{"type":"object","properties":{"memory":{"type":"string"}}}}]}`, true},
		{"ordinary", `{"messages":[{"content":"remember memory"}],"tools":[{"name":"submit","input_schema":{"type":"object","properties":{"result":{"type":"string"}}}}]}`, false},
		{"wrong type", `{"tools":[{"name":"submit","input_schema":{"properties":{"memory":{"type":"number"}}}}]}`, false},
		{"extra tool", `{"tools":[{"name":"submit","input_schema":{"properties":{"memory":{"type":"string"}}}},{"name":"read_file"}]}`, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := memoryCompressor([]byte(tc.body)); got != tc.memory {
				t.Fatalf("memory classification = %v, want %v", got, tc.memory)
			}
		})
	}
}

// TestMemoryOutcomeBoundary checks actual snapshot/journal files: a truncated
// terminal line, main turn end or unrelated task cannot acknowledge a launch.
func TestMemoryOutcomeBoundary(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "save.json")
	if err := os.WriteFile(path, []byte(`{"log":[{"type":"compactor_launched","compactor":{"task_id":"selected"}}]}`), 0600); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name, tail string
		pending    bool
	}{
		{"main end", "{\"type\":\"turn_ended\"}\n", true},
		{"other task", "{\"type\":\"compaction_outcome\",\"compaction_outcome\":{\"task_id\":\"other\",\"state\":\"committed\"}}\n", true},
		{"partial tail", `{"type":"compaction_outcome","compaction_outcome":{"task_id":"selected","state":"committed"}}`, true},
		{"complete", "{\"type\":\"compaction_outcome\",\"compaction_outcome\":{\"task_id\":\"selected\",\"state\":\"committed\"}}\n", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := os.WriteFile(path+".journal", []byte(tc.tail), 0600); err != nil {
				t.Fatal(err)
			}
			pending, err := memoryPending(dir)
			if err != nil || (len(pending) == 1 && pending[0] == "selected") != tc.pending {
				t.Fatalf("pending = %v, error = %v", pending, err)
			}
		})
	}
}
