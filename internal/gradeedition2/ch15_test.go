package gradeedition2

import (
	"testing"

	"github.com/waywardgeek/ensemble/internal/grade"
)

// TestChapter15Memory exercises the actual accepted GUI/Agent, owned settings,
// memory files, history and process lifetime. Only provider responses are fake.
func TestChapter15Memory(t *testing.T) {
	bin, cleanup, err := grade.Build("../../solutions/ch15/gui")
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
