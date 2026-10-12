package gradeedition2

import (
	"testing"

	"github.com/waywardgeek/ensemble/internal/grade"
)

// TestChapter16Recall exercises the accepted GUI/Agent with actual files and
// provider HTTP. The full command separately measures the real 15-second timeout;
// focused regressions reuse that evidence instead of paying the delay each time.
func TestChapter16Recall(t *testing.T) {
	bin, clean, err := grade.Build("../../solutions/edition-2/ch16/gui")
	if err != nil {
		t.Fatal(err)
	}
	defer clean()
	r := recallChecks(bin, "cr", false)
	if r.Fatal != "" {
		t.Fatal(r.Fatal)
	}
	for id, e := range r.Errs {
		if e != "" {
			t.Errorf("%s: %s", id, e)
		}
	}
	if err := recallTopology("../../solutions/edition-2/ch16", "../../solutions/edition-2/ch16/gui"); err != "" {
		t.Error(err)
	}
}
