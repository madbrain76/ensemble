package gradeedition2

import (
	"github.com/waywardgeek/ensemble/internal/grade"
	"os"
	"testing"
)

// These opt-in checks launch the real submitted binary; only its remote provider
// is replaced. They also permit focused restored-defect checks without GUI work.
func TestFleet22(t *testing.T) {
	bin := os.Getenv("EDITION2_FLEET_BIN")
	if bin == "" {
		t.Skip("set EDITION2_FLEET_BIN to a built submitted CLI")
	}
	checks := grade.Ch25Checks(grade.Ch25Result{})
	for i := range checks {
		checks[i].Passed = true
		checks[i].Earned = checks[i].Points
		checks[i].Details = nil
	}
	fleetChecks(bin, checks)
	for _, c := range checks {
		if !c.Passed {
			t.Errorf("%s: %v", c.ID, c.Details)
		}
	}
}
