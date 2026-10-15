package gradeedition2

import (
	"testing"

	"github.com/waywardgeek/ensemble/internal/grade"
)

// TestChapter18Combined drives both shipped clients against actual HTTP peers.
// It exercises current allowed models; historical two-priced-model evidence is
// retained separately rather than selecting the frozen unsupported 5.5 row.
func TestChapter18Combined(t *testing.T) {
	core := "../../solutions/ch18"
	bin, clean, err := grade.Build(core + "/gui")
	if err != nil {
		t.Fatal(err)
	}
	defer clean()
	report, err := Chapter18(bin, core, [4]float64{2, 2.5, .2, 10}, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range report.Checks {
		if !c.Passed {
			t.Errorf("%s: %v", c.ID, c.Details)
		}
	}
}
