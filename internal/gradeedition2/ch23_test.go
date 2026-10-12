package gradeedition2

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/waywardgeek/ensemble/internal/grade"
)

// An opt-in real executable keeps mutation runs focused on changed management
// observations, without rerunning unchanged recursive CLI and idle-correction checks.
func TestManagement23(t *testing.T) {
	bin := os.Getenv("EDITION2_MANAGEMENT_BIN")
	if bin == "" {
		t.Skip("set EDITION2_MANAGEMENT_BIN to an actual optional GUI executable")
	}
	checks := grade.Ch26Checks(grade.Ch26Result{})
	for i := range checks {
		checks[i].Passed = true
		checks[i].Earned = checks[i].Points
		checks[i].Details = nil
	}
	managementGUI(bin, checks)
	data, _ := json.Marshal(checks)
	fmt.Println(string(data))
	for _, i := range []int{0, 1, 2, 3, 5} {
		if !checks[i].Passed {
			t.Errorf("%s: %v", checks[i].ID, checks[i].Details)
		}
	}
}
