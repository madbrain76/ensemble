package gradeedition2

import (
	"testing"

	"github.com/waywardgeek/ensemble/internal/grade"
)

// TestChapter2AcceptedSolution exercises the public command path against the
// unchanged accepted solution.
// Targeted source mutations are recorded separately in the enhancement notes.
func TestChapter2AcceptedSolution(t *testing.T) {
	bin, cleanup, err := grade.Build("../../solutions/ch02")
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	report, err := Chapter2(bin)
	if err != nil {
		t.Fatal(err)
	}
	if !report.Passed || report.Score != 100 || report.MaxScore != 100 {
		t.Fatalf("corrected Chapter 2: %+v", report)
	}
}
