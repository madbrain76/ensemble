package main
import (
 "encoding/json"
 "os"
 _ "unsafe"
 "github.com/waywardgeek/ensemble/internal/grade"
 _ "github.com/waywardgeek/ensemble/internal/gradeedition2"
)
//go:linkname settingsChecks github.com/waywardgeek/ensemble/internal/gradeedition2.settingsChecks
func settingsChecks(string,string) []grade.Check
func main(){ checks:=settingsChecks(os.Args[1],"entry-audit-model"); json.NewEncoder(os.Stdout).Encode(checks); for _,c:=range checks { if !c.Passed {os.Exit(1)} } }
