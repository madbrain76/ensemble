package gradeedition2

import (
	"github.com/waywardgeek/ensemble/internal/grade"
	"path/filepath"
	"testing"
)

// TestChapter14Context drives the real optional GUI and Agent through settings,
// completed turns, crashes and reloads; only provider/MCP peers are external fakes.
func TestChapter14Context(t *testing.T) {
	gui, err := filepath.Abs("../../solutions/edition-2/ch14/gui")
	if err != nil {
		t.Fatal(err)
	}
	bin, cleanup, err := grade.Build(gui)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	report, err := Chapter14(bin)
	if err != nil {
		t.Fatal(err)
	}
	for _, check := range report.Checks {
		if !check.Passed {
			t.Errorf("%s: %v", check.ID, check.Details)
		}
	}
}
