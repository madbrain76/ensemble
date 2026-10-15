package gradeedition2

import (
	"github.com/waywardgeek/ensemble/internal/grade"
	"testing"
)

// TestChapter15Invoice runs actual GUI clients, save/restart and provider parsers.
// The external fake supplies disjoint counts; no cache grant is inferred from it.
func TestChapter15Invoice(t *testing.T) {
	bin, clean, err := grade.Build("../../agent/gui")
	if err != nil {
		t.Fatal(err)
	}
	defer clean()
	r := invoiceChecks(bin, "cr/cache-current.txt", "cr/cache-prior.txt", [4]float64{2, 2.5, .2, 10})
	if r.Fatal != "" {
		t.Fatal(r.Fatal)
	}
	for id, e := range r.Errs {
		if e != "" {
			t.Errorf("%s: %s", id, e)
		}
	}
}
