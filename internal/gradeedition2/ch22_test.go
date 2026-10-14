package gradeedition2

import (
	"github.com/waywardgeek/ensemble/internal/grade"
	"os"
	"testing"
)

// TestFleetRootChat22 checks two completed real GUI turns with fleet enabled,
// so a root missing-submit loop cannot be concealed by early process EOF.
func TestFleetRootChat22(t *testing.T) {
	bin, cleanup, err := grade.Build("../../agent/gui")
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	c := grade.Ch25Checks(grade.Ch25Result{})[6]
	c.Passed = true
	c.Details = nil
	fleetRootChat(bin, &c)
	if !c.Passed {
		t.Fatal(c.Details)
	}
}

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

// The outside consumer and deployed GUI remain real programs. This entry point
// isolates their affected checks from the terminal fleet during sensitivity work.
func TestReuse22(t *testing.T) {
	core := os.Getenv("EDITION2_REUSE_CORE")
	if core == "" {
		t.Skip("set EDITION2_REUSE_CORE to a submitted core module")
	}
	checks := grade.Ch24Checks(grade.Ch24Result{})
	for i := range checks {
		checks[i].Passed = true
		checks[i].Earned = checks[i].Points
		checks[i].Details = nil
	}
	reuseChecks(core, core+"/gui", false, checks)
	for _, c := range checks {
		if !c.Passed {
			t.Errorf("%s: %v", c.ID, c.Details)
		}
	}
}
