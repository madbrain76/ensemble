package gradeedition2

import (
	"github.com/waywardgeek/ensemble/internal/grade"
	"testing"
)

// TestChapter18Subscription exercises the actual GUI, login command and private
// profile with only the external issuer/provider replaced by faithful peers.
func TestChapter18Subscription(t *testing.T) {
	bin, clean, err := grade.Build("../../solutions/edition-2/ch18/gui")
	if err != nil {
		t.Fatal(err)
	}
	defer clean()
	r := subscriptionChecks(bin)
	if r.Fatal != "" {
		t.Fatal(r.Fatal)
	}
	for id, e := range r.Errs {
		if e != "" {
			t.Errorf("%s: %s", id, e)
		}
	}
}
