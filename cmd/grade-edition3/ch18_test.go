package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestChapter18GuardsAndRoute(t *testing.T) {
	dir := t.TempDir()
	core := filepath.Join(dir, "core")
	if err := os.Mkdir(core, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(core, "go.mod"), []byte("module course.invalid/status\n\ngo 1.25\n"), 0600); err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(dir, "launched")
	target := filepath.Join(dir, "target")
	if err := os.WriteFile(target, []byte("#!/bin/sh\nprintf launched >> "+strconv.Quote(marker)+"\nexit 71\n"), 0700); err != nil {
		t.Fatal(err)
	}
	base := []string{"-ch", "18", "-core", core, "-gui", target, "-status-prices", "3,3.75,.3,15"}
	bad := [][]string{{"-ch", "18", target}, {"-ch", "17", target}, {"-ch", "19", target}}
	for _, pair := range [][2]string{{"core", ""}, {"gui", ""}, {"gui", filepath.Join(dir, "absent")}, {"status-prices", ""}, {"status-prices", "NaN,1,1,1"}, {"status-prices", "1,-1,1,1"}, {"status-second-prices", ""}, {"status-second-prices", "1,2,3"}} {
		args := append([]string{}, base...)
		args = append(args, "-"+pair[0], pair[1], target)
		bad = append(bad, args)
	}
	for _, name := range []string{"consumer", "consumer-case", "workflow", "media-models", "stream-exercise", "stream-model", "gui-model", "cache-prices", "cache-current", "cache-prior"} {
		args := append([]string{}, base...)
		bad = append(bad, append(args, "-"+name, "", target))
	}
	for _, ch := range []string{"1", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16"} {
		for _, name := range []string{"status-prices", "status-second-prices"} {
			bad = append(bad, []string{"-ch", ch, "-" + name, "", target})
		}
	}
	for _, args := range bad {
		if code := run(args); code != 2 {
			t.Fatalf("%v exit=%d", args, code)
		}
		if _, err := os.Stat(marker); !os.IsNotExist(err) {
			t.Fatalf("invalid options launched target: %v", args)
		}
	}
	out, err := os.CreateTemp(dir, "report-")
	if err != nil {
		t.Fatal(err)
	}
	old := os.Stdout
	os.Stdout = out
	code := run(append(base, "-json", target))
	os.Stdout = old
	if err = out.Close(); err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(out.Name())
	if err != nil {
		t.Fatal(err)
	}
	var report struct {
		Title    string
		MaxScore int `json:"max_score"`
		Checks   []struct{ ID string }
	}
	if err = json.Unmarshal(body, &report); err != nil {
		t.Fatalf("%v %s", err, body)
	}
	if code != 1 || !strings.Contains(report.Title, "Chapter 18") || report.MaxScore != 155 || len(report.Checks) != 11 {
		t.Fatalf("wrong route: exit%d %s", code, body)
	}
	if _, err = os.Stat(marker); err != nil {
		t.Fatal(err)
	}
}
